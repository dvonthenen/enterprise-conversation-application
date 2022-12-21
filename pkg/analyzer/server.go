// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package analyzer

import (
	"context"
	"os"

	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	amqp "github.com/rabbitmq/amqp091-go"
	klog "k8s.io/klog/v2"

	handlers "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/handlers"
	rabbit "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit"
	common "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
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

func (s *Server) Init() error {
	klog.V(6).Infof("Server.Init ENTER\n")

	// symbl
	err := s.RebuildSymblClient()
	if err != nil {
		klog.V(1).Infof("RebuildSymblClient failed. Err: %v\n", err)
		klog.V(6).Infof("Server.Init LEAVE\n")
		return err
	}

	// neo4j
	err = s.RebuildDatabase()
	if err != nil {
		klog.V(1).Infof("RebuildDatabase failed. Err: %v\n", err)
		klog.V(6).Infof("Server.Init LEAVE\n")
		return err
	}

	// rabbitmq
	err = s.RebuildMessageBus()
	if err != nil {
		klog.V(1).Infof("RebuildDatabase failed. Err: %v\n", err)
		klog.V(6).Infof("Server.Init LEAVE\n")
		return err
	}

	klog.V(4).Infof("Server.Init Succeeded\n")
	klog.V(6).Infof("Server.Init LEAVE\n")

	return nil
}

func (s *Server) Start() error {
	klog.V(6).Infof("Server.Start ENTER\n")

	// start receiving notifications
	klog.V(2).Infof("Starting server...\n")
	err := s.notificationMgr.Start()
	if err != nil {
		klog.V(1).Infof("ListenAndServeTLS failed. Err: %v\n", err)
		klog.V(6).Infof("Server.Start LEAVE\n")
		return err
	}

	klog.V(4).Infof("Server.Start Succeeded\n")
	klog.V(6).Infof("Server.Start LEAVE\n")

	return nil
}

func (s *Server) RebuildSymblClient() error {
	klog.V(6).Infof("Server.RebuildSymblClient ENTER\n")

	ctx := context.Background()

	symblClient, err := symbl.NewRestClient(ctx)
	if err != nil {
		klog.V(1).Infof("RebuildSymblClient failed. Err: %v\n", err)
		klog.V(6).Infof("Server.RebuildSymblClient LEAVE\n")
		return err
	}
	s.symblClient = symblClient

	klog.V(4).Infof("Server.RebuildSymblClient Succeded\n")
	klog.V(6).Infof("Server.RebuildSymblClient LEAVE\n")

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
	s.driver = &driver

	klog.V(4).Infof("Server.RebuildDatabase Succeeded\n")
	klog.V(6).Infof("Server.RebuildDatabase LEAVE\n")

	return err
}

func (s *Server) RebuildMessageBus() error {
	klog.V(6).Infof("Server.RebuildMessageBus ENTER\n")

	// teardown
	if s.notificationMgr != nil {
		err := s.notificationMgr.Teardown()
		if err != nil {
			klog.V(1).Infof("notificationMgr.Teardown failed. Err: %v\n", err)
		}
		s.rabbitMgr = nil
	}
	if s.rabbitMgr != nil {
		err := s.rabbitMgr.Teardown()
		if err != nil {
			klog.V(1).Infof("rabbitMgr.DeleteAll failed. Err: %v\n", err)
		}
		s.rabbitMgr = nil
	}
	if s.rabbitConn != nil {
		s.rabbitConn.Close()
		s.rabbitConn = nil
	}

	// init
	conn, err := amqp.Dial(s.options.RabbitMQURI)
	if err != nil {
		klog.V(1).Infof("amqp.Dial failed. Err: %v\n", err)
		klog.V(6).Infof("Server.RebuildMessageBus LEAVE\n")
		return err
	}

	// notification client messages
	ch, err := conn.Channel()
	if err != nil {
		klog.V(1).Infof("conn.Channel failed. Err: %v\n", err)
		klog.V(6).Infof("Server.RebuildMessageBus LEAVE\n")
		return err
	}
	err = ch.ExchangeDeclare(
		common.RabbitClientNotifications, // name
		"fanout",                         // type
		true,                             // durable
		true,                             // auto-deleted
		false,                            // internal
		false,                            // no-wait
		nil,                              // arguments
	)
	if err != nil {
		klog.V(1).Infof("ch.ExchangeDeclare failed. Err: %v\n", err)
		klog.V(6).Infof("Server.RebuildMessageBus LEAVE\n")
		return err
	}

	// rabbitmgr
	if s.symblClient == nil {
		err := s.RebuildSymblClient()
		if err != nil {
			klog.V(1).Infof("RebuildSymblClient failed. Err: %v\n", err)
			klog.V(6).Infof("Server.RebuildMessageBus LEAVE\n")
			return err
		}
	}
	rabbitMgr := rabbit.New(rabbit.RabbitManagerOptions{
		Connection: conn,
	})

	notificationManager := handlers.NewNotificationManager(handlers.NotificationManagerOption{
		Driver:        s.driver,
		RabbitManager: rabbitMgr,
		SymblClient:   s.symblClient,
	})

	err = notificationManager.Init()
	if err != nil {
		klog.V(1).Infof("notificationManager.Init failed. Err: %v\n", err)
		klog.V(6).Infof("Server.RebuildMessageBus LEAVE\n")
		return err
	}

	// housekeeping
	s.rabbitConn = conn
	s.rabbitMgr = rabbitMgr
	s.notificationMgr = notificationManager

	klog.V(4).Infof("Server.RebuildMessageBus Succeeded\n")
	klog.V(6).Infof("Server.RebuildMessageBus LEAVE\n")

	return nil
}

func (s *Server) Stop() error {
	klog.V(6).Infof("Server.Stop ENTER\n")

	// clean up neo4j driver
	ctx := context.Background()
	if s.driver != nil {
		(*s.driver).Close(ctx)
	}
	s.driver = nil

	// clean up rabbitmq
	if s.rabbitConn != nil {
		s.rabbitConn.Close()
	}
	s.rabbitConn = nil

	klog.V(4).Infof("Server.Stop Succeeded\n")
	klog.V(6).Infof("Server.Stop LEAVE\n")

	return nil
}
