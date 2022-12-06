// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package server

// streaming
import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"

	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	amqp "github.com/rabbitmq/amqp091-go"
	klog "k8s.io/klog/v2"

	instance "github.com/dvonthenen/enterprise-reference-implementation/pkg/dataminer/instance"
	routing "github.com/dvonthenen/enterprise-reference-implementation/pkg/dataminer/routing"
)

func New(options ServerOptions) (*Server, error) {
	if options.StartPort == 0 {
		options.StartPort = DefaultStartPort
	}
	if options.EndPort == 0 {
		options.EndPort = DefaultEndPort
	}

	var connectionStr string
	if v := os.Getenv("NEO4J_CONNECTION"); v != "" {
		klog.V(4).Info("NEO4J_CONNECTION found")
		connectionStr = v
	} else {
		klog.Errorf("NEO4J_CONNECTION not found\n")
		return nil, ErrInvalidInput
	}
	var username string
	if v := os.Getenv("NEO4J_USERNAME"); v != "" {
		klog.V(4).Info("NEO4J_USERNAME found")
		username = v
	} else {
		klog.Errorf("NEO4J_USERNAME not found\n")
		return nil, ErrInvalidInput
	}
	var password string
	if v := os.Getenv("NEO4J_PASSWORD"); v != "" {
		klog.V(4).Info("NEO4J_PASSWORD found")
		password = v
	} else {
		klog.Errorf("NEO4J_PASSWORD not found\n")
		return nil, ErrInvalidInput
	}

	creds := Credentials{
		ConnectionStr: connectionStr,
		Username:      username,
		Password:      password,
	}

	// auth
	auth := neo4j.BasicAuth(creds.Username, creds.Password, "")

	// You typically have one driver instance for the entire application. The
	// driver maintains a pool of database connections to be used by the sessions.
	// The driver is thread safe.
	driver, err := neo4j.NewDriverWithContext(creds.ConnectionStr, auth)
	if err != nil {
		klog.V(1).Infof("NewDriverWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("MessageHandler.Init ENTER\n")
		return nil, err
	}

	// server
	server := &Server{
		options:        options,
		creds:          creds,
		instanceById:   make(map[string]*instance.ServerInstance),
		instanceByPort: make(map[int]*instance.ServerInstance),
		Driver:         driver,
	}
	return server, nil
}

func (se *Server) redirectToInstance(w http.ResponseWriter, r *http.Request) {
	// conversationId
	conversationId := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	klog.V(2).Infof("conversationId: %s\n", conversationId)

	// does the server already exist, return the serverInstance
	serverInstance := se.instanceById[conversationId]
	if serverInstance != nil {
		klog.V(2).Infof("Server for conversationId (%s) already exists\n", serverInstance.Options.ConversationId)
		http.Redirect(w, r, serverInstance.Options.RedirectAddress, http.StatusSeeOther)
	}

	// need to do this concurrently because the session create takes time
	chanCallback := make(chan *routing.MessageHandler)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// signal done
		defer wg.Done()

		// Create a neo4j session to run transactions in. Sessions are lightweight to
		// create and use. Sessions are NOT thread safe.
		ctx := context.Background()
		session := se.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

		// create server
		options := routing.MessageHandlerOptions{
			ConversationId: conversationId,
			Session:        session,
			RabbitChan:     se.rabbitChan,
		}
		callback, err := routing.NewHandler(options)
		if err != nil {
			klog.V(1).Infof("server.Start failed. Err: %v\n", err)
			http.Error(w, "Failed to create message handler", http.StatusBadRequest)
			return
		}

		// init
		err = callback.Init()
		if err != nil {
			klog.V(1).Infof("callback.Init failed. Err: %v\n", err)
			http.Error(w, "Failed to init message handler", http.StatusBadRequest)
			return
		}

		chanCallback <- callback
	}()

	// get random port
	diff := se.options.EndPort - se.options.StartPort
	var random int
	for {
		random = se.options.StartPort + rand.Intn(diff)
		if se.instanceByPort[random] == nil {
			break // found an unused port
		}
	}

	// bind address
	newServer := fmt.Sprintf("%s:%d", r.URL.Host, random)
	klog.V(2).Infof("newServer: %s\n", newServer)

	// redirect address
	redirect := r.URL.Host
	if len(redirect) == 0 {
		redirect = "127.0.0.1"
	}
	newRedirect := fmt.Sprintf("https://%s:%d", redirect, random)
	klog.V(2).Infof("newRedirect: %s\n", newRedirect)

	server := instance.New(instance.InstanceOptions{
		Port:            random,
		BindAddress:     newServer,
		RedirectAddress: newRedirect,
		ConversationId:  conversationId,
		CrtFile:         se.options.CrtFile,
		KeyFile:         se.options.KeyFile,
		Callback:        <-chanCallback,
	})

	err := server.Start()
	if err != nil {
		klog.V(1).Infof("server.Start failed. Err: %v\n", err)
		http.Error(w, "Failed to start server instance", http.StatusBadRequest)
		return
	}
	se.instanceById[conversationId] = server
	se.instanceByPort[random] = server

	// wait for everyone to finish
	wg.Wait()

	// redirect
	http.Redirect(w, r, newRedirect, http.StatusSeeOther)
}

func (se *Server) Start() error {
	// redirect
	mux := http.NewServeMux()
	mux.HandleFunc("/", se.redirectToInstance)

	se.server = &http.Server{
		Addr:    ":443",
		Handler: mux,
	}

	//

	// this is a blocking call
	klog.V(2).Infof("Starting server...\n")
	err := se.server.ListenAndServeTLS(se.options.CrtFile, se.options.KeyFile)
	if err != nil {
		klog.V(6).Infof("ListenAndServeTLS failed. Err: %v\n", err)
	}

	return nil
}

func (se *Server) RebuildMessageBus() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		klog.V(1).Infof("amqp.Dial failed. Err: %v\n", err)
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		klog.V(1).Infof("conn.Channel failed. Err: %v\n", err)
		return err
	}

	err = ch.ExchangeDeclare(
		"create-conversation", // name
		"fanout",              // type
		true,                  // durable
		true,                  // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		klog.V(1).Infof("ch.ExchangeDeclare failed. Err: %v\n", err)
		return err
	}

	// save to pass onto instances
	se.rabbitChan = ch
	se.rabbitConn = conn

	return nil
}

// TODO check for dead instances

func (se *Server) Stop() error {
	// stop all instances
	for _, instance := range se.instanceById {
		err := instance.Stop()
		if err != nil {
			fmt.Printf("instance.Stop() failed. Err: %v\n", err)
		}
	}
	se.instanceById = make(map[string]*instance.ServerInstance)
	se.instanceByPort = make(map[int]*instance.ServerInstance)

	// clean up neo4j driver
	ctx := context.Background()
	se.Driver.Close(ctx)

	// clean up rabbitmq
	se.rabbitConn.Close()
	se.rabbitChan.Close()

	// stop this endpoint
	err := se.server.Close()
	if err != nil {
		fmt.Printf("server.Close() failed. Err: %v\n", err)
	}

	return nil
}
