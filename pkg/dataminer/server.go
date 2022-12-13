// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package dataminer

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"

	symblinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	wsinterfaces "github.com/koding/websocketproxy/pkg/interfaces"
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

	// server
	server := &Server{
		options:        options,
		creds:          creds,
		instanceById:   make(map[string]*instance.ServerInstance),
		instanceByPort: make(map[int]*instance.ServerInstance),
		driver:         nil,
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
		session := (*se.driver).NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

		// create server
		options := routing.MessageHandlerOptions{
			ConversationId:   conversationId,
			Session:          &session,
			RabbitConnection: se.rabbitConn,
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
	klog.V(2).Infof("Bind Address: %s\n", newServer)

	// redirect address
	redirect := r.URL.Host
	if len(redirect) == 0 {
		redirect = "127.0.0.1"
	}
	newRedirect := fmt.Sprintf("https://%s:%d", redirect, random)
	klog.V(2).Infof("New Proxy Server: %s\n", newRedirect)

	var callback symblinterfaces.InsightCallback
	callback = <-chanCallback

	var manager wsinterfaces.ManageCallback
	manager = se

	server := instance.New(instance.InstanceOptions{
		Port:            random,
		BindAddress:     newServer,
		RedirectAddress: newRedirect,
		ConversationId:  conversationId,
		CrtFile:         se.options.CrtFile,
		KeyFile:         se.options.KeyFile,
		Callback:        &callback,
		Manager:         &manager,
	})

	err := server.Start()
	if err != nil {
		klog.V(1).Infof("server.Start failed. Err: %v\n", err)
		http.Error(w, "Failed to start server instance", http.StatusBadRequest)
		return
	}

	se.mu.Lock()
	se.instanceById[conversationId] = server
	se.instanceByPort[random] = server
	se.mu.Unlock()

	// wait for everyone to finish
	wg.Wait()

	// redirect
	http.Redirect(w, r, newRedirect, http.StatusSeeOther)
}

func (se *Server) Start() error {
	// neo4j
	err := se.RebuildDatabase()
	if err != nil {
		klog.V(6).Infof("RebuildDatabase failed. Err: %v\n", err)
		return err
	}

	// rabbitmq
	err = se.RebuildMessageBus()
	if err != nil {
		klog.V(6).Infof("RebuildDatabase failed. Err: %v\n", err)
		return err
	}

	// redirect
	mux := http.NewServeMux()
	mux.HandleFunc("/", se.redirectToInstance)

	se.server = &http.Server{
		Addr:    ":443",
		Handler: mux,
	}

	go func() {
		// this is a blocking call
		klog.V(2).Infof("Starting server...\n")
		err = se.server.ListenAndServeTLS(se.options.CrtFile, se.options.KeyFile)
		if err != nil {
			klog.V(6).Infof("ListenAndServeTLS failed. Err: %v\n", err)
		}
	}()

	return nil
}

func (se *Server) RebuildDatabase() error {
	//teardown
	if se.driver != nil {
		ctx := context.Background()
		(*se.driver).Close(ctx)
		se.driver = nil
	}

	// init neo4j
	// auth
	auth := neo4j.BasicAuth(se.creds.Username, se.creds.Password, "")

	// You typically have one driver instance for the entire application. The
	// driver maintains a pool of database connections to be used by the sessions.
	// The driver is thread safe.
	driver, err := neo4j.NewDriverWithContext(se.creds.ConnectionStr, auth)
	if err != nil {
		klog.V(1).Infof("NewDriverWithContext failed. Err: %v\n", err)
		return err
	}

	se.driver = &driver

	return err
}

func (se *Server) RebuildMessageBus() error {
	// teardown
	if se.rabbitConn != nil {
		se.rabbitConn.Close()
		se.rabbitConn = nil
	}

	// init rabbitmq
	conn, err := amqp.Dial(se.options.RabbitMQURI)
	if err != nil {
		klog.V(1).Infof("amqp.Dial failed. Err: %v\n", err)
		return err
	}

	// save to pass onto instances
	se.rabbitConn = conn

	return nil
}

func (se *Server) RemoveConnection(uniqueId string) {
	se.mu.Lock()
	instance := se.instanceById[uniqueId]
	if instance == nil {
		klog.V(1).Infof("RemoveConnection(%s) instance not found\n", uniqueId)
		se.mu.Unlock()
		return
	}

	port := instance.Options.Port

	delete(se.instanceById, uniqueId)
	delete(se.instanceByPort, port)

	klog.V(2).Infof("RemoveConnection(%s) Successful\n", uniqueId)
	se.mu.Unlock()
}

func (se *Server) CheckForDeadInstances() {
	for _, instance := range se.instanceById {
		if !instance.IsConnected() {
			se.mu.Lock()
			err := instance.Stop()
			if err != nil {
				klog.V(1).Infof("instance.Stop() failed. Err: %v\n", err)
			}

			delete(se.instanceById, instance.Options.ConversationId)
			delete(se.instanceByPort, instance.Options.Port)
			se.mu.Unlock()
		}
	}
}

func (se *Server) Stop() error {
	// stop all instances
	for _, instance := range se.instanceById {
		err := instance.Stop()
		if err != nil {
			klog.V(1).Infof("instance.Stop() failed. Err: %v\n", err)
		}
	}
	se.instanceById = make(map[string]*instance.ServerInstance)
	se.instanceByPort = make(map[int]*instance.ServerInstance)

	// clean up neo4j driver
	ctx := context.Background()
	(*se.driver).Close(ctx)

	// clean up rabbitmq
	se.rabbitConn.Close()

	// stop this endpoint
	err := se.server.Close()
	if err != nil {
		klog.V(1).Infof("server.Close() failed. Err: %v\n", err)
	}

	return nil
}