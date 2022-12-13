// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package analyzer

import (
	"context"
	"fmt"
	"net/http"
	"os"

	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	amqp "github.com/rabbitmq/amqp091-go"
	klog "k8s.io/klog/v2"

	handlers "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/handlers"
	rabbit "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit"
)

func New(options ServerOptions) (*Server, error) {
	if options.BindPort == 0 {
		options.BindPort = DefaultPort
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
		options: options,
		creds:   creds,
	}
	return server, nil
}

func (se *Server) doSomething(w http.ResponseWriter, r *http.Request) {
	// http.Redirect(w, r, newRedirect, http.StatusSeeOther)
}

func (se *Server) Init() error {
	// neo4j
	err := se.RebuildDatabase()
	if err != nil {
		klog.V(1).Infof("RebuildDatabase failed. Err: %v\n", err)
		return err
	}

	// rabbitmq
	err = se.RebuildMessageBus()
	if err != nil {
		klog.V(1).Infof("RebuildDatabase failed. Err: %v\n", err)
		return err
	}

	return nil
}

func (se *Server) Start() error {
	// start receiving notifications
	err := se.notificationMgr.Start()
	if err != nil {
		klog.V(1).Infof("ListenAndServeTLS failed. Err: %v\n", err)
		return err
	}

	// redirect
	newServer := fmt.Sprintf("%s:%d", se.options.BindAddress, se.options.BindPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", se.doSomething)

	se.server = &http.Server{
		Addr:    newServer,
		Handler: mux,
	}

	go func() {
		// this is a blocking call
		klog.V(2).Infof("Starting server...\n")
		err := se.server.ListenAndServeTLS(se.options.CrtFile, se.options.KeyFile)
		if err != nil {
			klog.V(1).Infof("ListenAndServeTLS failed. Err: %v\n", err)
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
	if se.notificationMgr != nil {
		err := se.notificationMgr.Teardown()
		if err != nil {
			klog.V(1).Infof("notificationMgr.Teardown failed. Err: %v\n", err)
		}
		se.rabbitMgr = nil
	}
	if se.rabbitMgr != nil {
		err := se.rabbitMgr.Delete()
		if err != nil {
			klog.V(1).Infof("rabbitMgr.DeleteAll failed. Err: %v\n", err)
		}
		se.rabbitMgr = nil
	}
	if se.rabbitConn != nil {
		se.rabbitConn.Close()
		se.rabbitConn = nil
	}

	// init
	conn, err := amqp.Dial(se.options.RabbitMQURI)
	if err != nil {
		klog.V(1).Infof("amqp.Dial failed. Err: %v\n", err)
		return err
	}

	// rabbitmgr
	rabbitMgr := rabbit.New(rabbit.RabbitManagerOptions{
		Connection: conn,
	})

	notificationManager := handlers.NewNotificationManager(handlers.NotificationManagerOption{
		Driver:        se.driver,
		RabbitManager: rabbitMgr,
	})

	err = notificationManager.Init()
	if err != nil {
		klog.V(1).Infof("notificationManager.Init failed. Err: %v\n", err)
		return err
	}

	// housekeeping
	se.rabbitConn = conn
	se.rabbitMgr = rabbitMgr
	se.notificationMgr = notificationManager

	return nil
}

func (se *Server) Stop() error {
	// clean up neo4j driver
	ctx := context.Background()
	if se.driver != nil {
		(*se.driver).Close(ctx)
	}
	se.driver = nil

	// clean up rabbitmq
	if se.rabbitConn != nil {
		se.rabbitConn.Close()
	}
	se.rabbitConn = nil

	// stop this endpoint
	err := se.server.Close()
	if err != nil {
		klog.V(1).Infof("server.Close() failed. Err: %v\n", err)
	}

	return nil
}