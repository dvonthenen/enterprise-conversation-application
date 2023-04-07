// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package server

import (
	"context"
	"os"

	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"

	middlewaresdk "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk"
	interfacessdk "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/interfaces"

	handlers "github.com/dvonthenen/enterprise-reference-implementation/cmd/example-middleware-plugin/handlers"
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

	// middleware analyzer
	err = s.RebuildMiddlewareAnalyzer()
	if err != nil {
		klog.V(1).Infof("RebuildMiddlewareAnalyzer failed. Err: %v\n", err)
		klog.V(6).Infof("Server.Init LEAVE\n")
		return err
	}

	klog.V(4).Infof("Server.Init Succeeded\n")
	klog.V(6).Infof("Server.Init LEAVE\n")

	return nil
}

func (s *Server) Start() error {
	klog.V(6).Infof("Server.Start ENTER\n")

	// rebuild neo4j driver if needed
	if s.driver == nil {
		klog.V(4).Infof("Calling RebuildDatabase...\n")
		err := s.RebuildDatabase()
		if err != nil {
			klog.V(1).Infof("RebuildDatabase failed. Err: %v\n", err)
			klog.V(6).Infof("Server.Start LEAVE\n")
			return err
		}
	}

	// rebuild symbl client if needed
	if s.symblClient == nil {
		klog.V(4).Infof("Calling RebuildSymblClient...\n")
		err := s.RebuildSymblClient()
		if err != nil {
			klog.V(1).Infof("RebuildSymblClient failed. Err: %v\n", err)
			klog.V(6).Infof("Server.Start LEAVE\n")
			return err
		}
	}

	// rebuild middleware if needed
	if s.middlewareAnalyzer == nil {
		klog.V(4).Infof("Calling RebuildMiddlewareAnalyzer...\n")
		err := s.RebuildMiddlewareAnalyzer()
		if err != nil {
			klog.V(1).Infof("RebuildMiddlewareAnalyzer failed. Err: %v\n", err)
			klog.V(6).Infof("Server.Start LEAVE\n")
			return err
		}
	}

	// start middleware
	err := s.middlewareAnalyzer.Init()
	if err != nil {
		klog.V(1).Infof("middlewareAnalyzer.Init() failed. Err: %v\n", err)
		klog.V(6).Infof("Server.Start LEAVE\n")
		return err
	}

	// TODO: start metrics and tracing

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

	// housekeeping
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

	// housekeeping
	s.driver = &driver

	klog.V(4).Infof("Server.RebuildDatabase Succeeded\n")
	klog.V(6).Infof("Server.RebuildDatabase LEAVE\n")

	return err
}

func (s *Server) RebuildMiddlewareAnalyzer() error {
	klog.V(6).Infof("Server.RebuildMiddlewareAnalyzer ENTER\n")

	// teardown
	if s.middlewareAnalyzer != nil {
		err := (*s.middlewareAnalyzer).Teardown()
		if err != nil {
			klog.V(1).Infof("middlewareAnalyzer.Teardown failed. Err: %v\n", err)
		}
		s.middlewareAnalyzer = nil
	}

	// create database session
	ctx := context.Background()
	session := (*s.driver).NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

	// create handler
	messageHandler := handlers.NewHandler(handlers.HandlerOptions{
		Session:     &session,
		SymblClient: s.symblClient,
	})

	// create middleware
	var callback interfacessdk.InsightCallback
	callback = messageHandler

	middlewareAnalyzer, err := middlewaresdk.NewMiddlewareAnalyzer(middlewaresdk.MiddlewareAnalyzerOption{
		RabbitURI: s.options.RabbitURI,
		Callback:  &callback,
	})
	if err != nil {
		klog.V(1).Infof("NewMiddlewareAnalyzer failed. Err: %v\n", err)
		klog.V(6).Infof("Server.RebuildMiddlewareAnalyzer LEAVE\n")
		return err
	}

	// housekeeping
	s.middlewareAnalyzer = middlewareAnalyzer

	klog.V(4).Infof("Server.RebuildMiddlewareAnalyzer Succeeded\n")
	klog.V(6).Infof("Server.RebuildMiddlewareAnalyzer LEAVE\n")

	return nil
}

func (s *Server) Stop() error {
	klog.V(6).Infof("Server.Stop ENTER\n")

	// TODO: stop metrics and tracing

	// clean up middleware
	if s.middlewareAnalyzer != nil {
		err := s.middlewareAnalyzer.Teardown()
		if err != nil {
			klog.V(1).Infof("middlewareAnalyzer.Teardown failed. Err: %v\n", err)
		}
	}
	s.middlewareAnalyzer = nil

	// clean up symbl client
	s.symblClient = nil

	// clean up neo4j driver
	if s.driver != nil {
		ctx := context.Background()
		(*s.driver).Close(ctx)
	}
	s.driver = nil

	klog.V(4).Infof("Server.Stop Succeeded\n")
	klog.V(6).Infof("Server.Stop LEAVE\n")

	return nil
}
