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
	"time"

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
		ticker:         time.NewTicker(time.Minute),
		stopPoll:       make(chan struct{}),
	}
	return server, nil
}

func (s *Server) redirectProxy(w http.ResponseWriter, r *http.Request) {
	// conversationId
	conversationId := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	klog.V(3).Infof("[redirectProxy] conversationId: %s\n", conversationId)

	// does the server already exist, return the serverInstance
	serverInstance := s.instanceById[conversationId]
	if serverInstance != nil {
		klog.V(3).Infof("Server for conversationId (%s) already exists\n", serverInstance.Options.ConversationId)
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
		// create and us. Sessions are NOT thread safe.
		ctx := context.Background()
		session := (*s.driver).NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

		// create server
		options := routing.MessageHandlerOptions{
			ConversationId:   conversationId,
			Session:          &session,
			RabbitConnection: s.rabbitConn,
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
	diff := s.options.EndPort - s.options.StartPort
	var random int
	for {
		random = s.options.StartPort + rand.Intn(diff)
		if s.instanceByPort[random] == nil {
			break // found an unused port
		}
	}

	// bind address
	newProxyServer := fmt.Sprintf("%s:%d", r.URL.Host, random)
	klog.V(3).Infof("Proxy Bind Address: %s\n", newProxyServer)
	newNotifyServer := fmt.Sprintf("%s:%d", r.URL.Host, (random + DefaultNotificationPortOffset))
	klog.V(3).Infof("Notify Bind Address: %s\n", newNotifyServer)

	// redirect address
	redirect := r.URL.Host
	if len(redirect) == 0 {
		redirect = "127.0.0.1"
	}
	newRedirect := fmt.Sprintf("https://%s:%d", redirect, random)
	klog.V(2).Infof("Proxy Redirect: %s\n", newRedirect)

	var callback symblinterfaces.InsightCallback
	callback = <-chanCallback

	var manager wsinterfaces.ManageCallback
	manager = s

	server := instance.New(instance.InstanceOptions{
		ProxyPort:         random,
		NotifyPort:        (random + DefaultNotificationPortOffset),
		ProxyBindAddress:  newProxyServer,
		NotifyBindAddress: newNotifyServer,
		RedirectAddress:   newRedirect,
		ConversationId:    conversationId,
		CrtFile:           s.options.CrtFile,
		KeyFile:           s.options.KeyFile,
		RabbitConn:        s.rabbitConn,
		Callback:          &callback,
		Manager:           &manager,
	})

	err := server.Init()
	if err != nil {
		klog.V(1).Infof("server.Init failed. Err: %v\n", err)
		http.Error(w, "Failed to init server instance", http.StatusBadRequest)
		return
	}

	err = server.Start()
	if err != nil {
		klog.V(1).Infof("server.Start failed. Err: %v\n", err)
		http.Error(w, "Failed to start server instance", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.instanceById[conversationId] = server
	s.instanceByPort[random] = server
	s.mu.Unlock()

	// wait for everyone to finish
	wg.Wait()

	// redirect
	http.Redirect(w, r, newRedirect, http.StatusSeeOther)
}

func (s *Server) redirectNotification(w http.ResponseWriter, r *http.Request) {
	// conversationId
	klog.V(3).Infof("[redirectNotification] URL Path: %s\n", r.URL.Path)
	split := strings.Split(r.URL.Path, "/")
	tokenCnt := len(split)
	conversationId := split[tokenCnt-2]
	klog.V(3).Infof("[redirectNotification] conversationId: %s\n", conversationId)

	// does the server already exist, return the serverInstance
	serverInstance := s.instanceById[conversationId]
	if serverInstance == nil {
		klog.V(2).Infof("Server for conversationId (%s) doesn't exists\n", conversationId)
		http.Error(w, "Failed to find conversationId instance", http.StatusNotFound)
		return
	}

	// redirect address
	redirect := r.URL.Host
	if len(redirect) == 0 {
		redirect = "127.0.0.1"
	}
	newRedirect := fmt.Sprintf("https://%s:%d%s", redirect, serverInstance.Options.NotifyPort, r.URL.Path)
	klog.V(3).Infof("Notify Redirect: %s\n", newRedirect)

	http.Redirect(w, r, newRedirect, http.StatusSeeOther)
}

func (s *Server) redirectToInstance(w http.ResponseWriter, r *http.Request) {
	lastToken := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	klog.V(3).Infof("URL: %s\n", r.URL.String())
	klog.V(3).Infof("Last Token: %s\n", lastToken)

	if lastToken == DefaultNotificationPath {
		s.redirectNotification(w, r)
	} else {
		s.redirectProxy(w, r)
	}
}

func (s *Server) Start() error {
	klog.V(6).Infof("Server.Start ENTER\n")

	// neo4j
	err := s.RebuildDatabase()
	if err != nil {
		klog.V(6).Infof("RebuildDatabase failed. Err: %v\n", err)
		klog.V(6).Infof("Server.Start LEAVE\n")
		return err
	}

	// rabbitmq
	err = s.RebuildMessageBus()
	if err != nil {
		klog.V(6).Infof("RebuildDatabase failed. Err: %v\n", err)
		klog.V(6).Infof("Server.Start LEAVE\n")
		return err
	}

	// redirect
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.redirectToInstance)

	s.server = &http.Server{
		Addr:    ":443",
		Handler: mux,
	}

	// poll for dead instances
	checkForDeadSessions := func(stopChan chan struct{}) {
		for {
			select {
			case <-s.ticker.C:
				s.CheckForDeadInstances()
			case <-stopChan:
				return
			}
		}
	}
	go checkForDeadSessions(s.stopPoll)

	// start the main entry endpoint to direct traffic
	go func() {
		// this is a blocking call
		klog.V(2).Infof("Starting server...\n")
		err = s.server.ListenAndServeTLS(s.options.CrtFile, s.options.KeyFile)
		if err != nil {
			klog.V(6).Infof("ListenAndServeTLS server stopped. Err: %v\n", err)
		}
	}()

	klog.V(4).Infof("Server.Start Succeeded\n")
	klog.V(6).Infof("Server.Start LEAVE\n")

	return nil
}

func (s *Server) RebuildDatabase() error {
	klog.V(6).Infof("Server.RebuildDatabase ENTER\n")

	//teardown
	if s.driver != nil {
		ctx := context.Background()
		(*s.driver).Close(ctx)
		s.driver = nil
	}

	// init neo4j
	// auth
	auth := neo4j.BasicAuth(s.creds.Username, s.creds.Password, "")

	// You typically have one driver instance for the entire application. The
	// driver maintains a pool of database connections to be used by the sessions.
	// The driver is thread safe.
	driver, err := neo4j.NewDriverWithContext(s.creds.ConnectionStr, auth)
	if err != nil {
		klog.V(1).Infof("NewDriverWithContext failed. Err: %v\n", err)
		klog.V(6).Infof("Server.RebuildDatabase LEAVE\n")
		return err
	}

	// save to pass onto instances
	s.driver = &driver

	klog.V(4).Infof("Server.RebuildDatabase Succeeded\n")
	klog.V(6).Infof("Server.RebuildDatabase LEAVE\n")

	return err
}

func (s *Server) RebuildMessageBus() error {
	klog.V(6).Infof("Server.RebuildMessageBus LEAVE\n")

	// teardown
	if s.rabbitConn != nil {
		s.rabbitConn.Close()
		s.rabbitConn = nil
	}

	// init rabbitmq
	conn, err := amqp.Dial(s.options.RabbitMQURI)
	if err != nil {
		klog.V(1).Infof("amqp.Dial failed. Err: %v\n", err)
		klog.V(6).Infof("Server.RebuildMessageBus LEAVE\n")
		return err
	}

	// save to pass onto instances
	s.rabbitConn = conn

	klog.V(4).Infof("Server.RebuildMessageBus Succeeded\n")
	klog.V(6).Infof("Server.RebuildMessageBus LEAVE\n")

	return nil
}

func (s *Server) RemoveConnection(uniqueId string) {
	klog.V(6).Infof("Server.RemoveConnection ENTER\n")

	s.mu.Lock()
	instance := s.instanceById[uniqueId]
	if instance == nil {
		klog.V(3).Infof("RemoveConnection(%s) instance not found\n", uniqueId)
		klog.V(6).Infof("Server.RemoveConnection LEAVE\n")
		s.mu.Unlock()
		return
	}

	// stop instance cleanly
	err := instance.Stop()
	if err != nil {
		klog.V(1).Infof("instance.Stop() failed. Err: %v\n", err)
	}

	// remove from record keeping
	port := instance.Options.ProxyPort

	delete(s.instanceById, uniqueId)
	delete(s.instanceByPort, port)

	klog.V(3).Infof("RemoveConnection(%s) Successful\n", uniqueId)
	klog.V(6).Infof("Server.RemoveConnection LEAVE\n")

	s.mu.Unlock()
}

func (s *Server) CheckForDeadInstances() {
	s.mu.Lock()
	for _, instance := range s.instanceById {
		if !instance.IsConnected() {
			err := instance.Stop()
			if err != nil {
				klog.V(1).Infof("instance.Stop() failed. Err: %v\n", err)
			}

			delete(s.instanceById, instance.Options.ConversationId)
			delete(s.instanceByPort, instance.Options.ProxyPort)
		}
	}
	s.mu.Unlock()
}

func (s *Server) Stop() error {
	klog.V(6).Infof("Server.Stop ENTER\n")

	// stop thread
	close(s.stopPoll)
	<-s.stopPoll

	// stop all instances
	s.mu.Lock()
	for _, instance := range s.instanceById {
		err := instance.Stop()
		if err != nil {
			klog.V(1).Infof("instance.Stop() failed. Err: %v\n", err)
		}
	}
	s.instanceById = make(map[string]*instance.ServerInstance)
	s.instanceByPort = make(map[int]*instance.ServerInstance)
	s.mu.Unlock()

	// clean up neo4j driver
	ctx := context.Background()
	(*s.driver).Close(ctx)

	// clean up rabbitmq
	s.rabbitConn.Close()

	// stop this endpoint
	err := s.server.Close()
	if err != nil {
		klog.V(1).Infof("server.Close() failed. Err: %v\n", err)
	}

	klog.V(4).Infof("Server.Stop Succeeded\n")
	klog.V(6).Infof("Server.Stop LEAVE\n")

	return nil
}
