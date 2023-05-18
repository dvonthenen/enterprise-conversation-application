// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package dataminer

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	rabbit "github.com/dvonthenen/rabbitmq-manager/pkg"
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	async "github.com/dvonthenen/symbl-go-sdk/pkg/api/async/v1"
	interfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/async/v1/interfaces"
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"

	routing "github.com/dvonthenen/enterprise-reference-implementation/pkg/rest-dataminer/routing"
)

func New(options ServerOptions) (*Server, error) {
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

	// DB Creds
	creds := Credentials{
		ConnectionStr: connectionStr,
		Username:      username,
		Password:      password,
	}

	if options.AuthMethod == AuthTypeDefault {
		options.AuthMethod = AuthTypeReuseToken
	}

	// server
	server := &Server{
		options: options,
		creds:   creds,
	}
	return server, nil
}

func (s *Server) processConversation(w http.ResponseWriter, r *http.Request) {
	conversationId := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	klog.V(3).Infof("URL: %s\n", r.URL.String())
	klog.V(3).Infof("conversationId: %s\n", conversationId)

	// init rabbit manager
	rabbitMgr, err := rabbit.New(rabbitinterfaces.ManagerOptions{
		RabbitURI: s.options.RabbitURI,
	})
	if err != nil {
		str := fmt.Sprintf("rabbit.New failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// init the db
	ctx := context.Background()
	session := (*s.driver).NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

	// init message handler
	handler, err := routing.NewHandler(routing.MessageHandlerOptions{
		ConversationId: conversationId,
		Neo4jMgr:       &session,
		RabbitMgr:      rabbitMgr,
	})
	if err != nil {
		str := fmt.Sprintf("NewHandler failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	exists, err := handler.DoesConversationExist()
	if err != nil {
		str := fmt.Sprintf("handler.Init failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}
	if s.options.DisableDuplicate && exists {
		klog.V(4).Infof("Conversation %s already processed. Exiting.\n", conversationId)
		w.WriteHeader(http.StatusOK)
		return
	}

	err = handler.Init()
	if err != nil {
		str := fmt.Sprintf("handler.Init failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// auth client
	var restClient *symbl.RestClient
	if s.options.AuthMethod == AuthTypeReuseToken {
		accessToken := r.Header.Get("X-API-KEY")
		// klog.V(7).Infof("X-API-KEY: %s\n", accessToken)

		restClient, err = symbl.NewRestClientWithToken(ctx, accessToken)
		if err != nil {
			str := fmt.Sprintf("NewRestClient failed. Err: %v\n", err)
			klog.V(1).Infof(str)
			http.Error(w, str, http.StatusBadRequest)
			return
		}
	} else {
		restClient, err = symbl.NewRestClient(ctx)
		if err != nil {
			str := fmt.Sprintf("NewRestClient failed. Err: %v\n", err)
			klog.V(1).Infof(str)
			http.Error(w, str, http.StatusBadRequest)
			return
		}
	}

	// create client
	asyncClient := async.New(restClient)

	// init
	err = handler.InitializedConversation(&interfaces.InitializationMessage{
		ConversationID: conversationId,
	})
	if err != nil {
		str := fmt.Sprintf("InitializedConversation failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// topics
	topicResult, err := asyncClient.GetTopics(ctx, conversationId)
	if err != nil {
		str := fmt.Sprintf("GetTopics failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	err = handler.TopicResult(topicResult)
	if err != nil {
		str := fmt.Sprintf("TopicResult failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// questions
	questionResult, err := asyncClient.GetQuestions(ctx, conversationId)
	if err != nil {
		str := fmt.Sprintf("GetQuestions failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	err = handler.QuestionResult(questionResult)
	if err != nil {
		str := fmt.Sprintf("QuestionResult failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// follow ups
	followupResult, err := asyncClient.GetFollowUps(ctx, conversationId)
	if err != nil {
		str := fmt.Sprintf("GetFollowUps failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	err = handler.FollowUpResult(followupResult)
	if err != nil {
		str := fmt.Sprintf("FollowUpResult failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// entity
	entityResult, err := asyncClient.GetEntities(ctx, conversationId)
	if err != nil {
		str := fmt.Sprintf("GetEntities failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	err = handler.EntityResult(entityResult)
	if err != nil {
		str := fmt.Sprintf("EntityResult failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// action items
	actionitemResult, err := asyncClient.GetActionItems(ctx, conversationId)
	if err != nil {
		str := fmt.Sprintf("GetActionItems failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	err = handler.ActionItemResult(actionitemResult)
	if err != nil {
		str := fmt.Sprintf("ActionItemResult failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// messages
	messageResult, err := asyncClient.GetMessages(ctx, conversationId)
	if err != nil {
		str := fmt.Sprintf("GetMessages failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	err = handler.MessageResult(messageResult)
	if err != nil {
		str := fmt.Sprintf("MessageResult failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// tracker
	trackerResult, err := asyncClient.GetTracker(ctx, conversationId)
	if err != nil {
		str := fmt.Sprintf("GetTracker failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	err = handler.TrackerResult(trackerResult)
	if err != nil {
		str := fmt.Sprintf("TrackerResult failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	// teardown
	err = handler.TeardownConversation(&interfaces.TeardownMessage{
		ConversationID: conversationId,
	})
	if err != nil {
		str := fmt.Sprintf("TeardownConversation failed. Err: %v\n", err)
		klog.V(1).Infof(str)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) Start() error {
	klog.V(6).Infof("Server.Start ENTER\n")

	// neo4j
	if s.driver == nil {
		klog.V(4).Infof("Calling RebuildDatabase...\n")
		err := s.RebuildDatabase()
		if err != nil {
			klog.V(1).Infof("RebuildDatabase failed. Err: %v\n", err)
			klog.V(6).Infof("Server.Start LEAVE\n")
			return err
		}
	}

	// redirect
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.processConversation)

	s.server = &http.Server{
		Addr:    ":443",
		Handler: mux,
	}

	// start the main entry endpoint to direct traffic
	go func() {
		// this is a blocking call
		err := s.server.ListenAndServeTLS(s.options.CrtFile, s.options.KeyFile)
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

func (s *Server) Stop() error {
	klog.V(6).Infof("Server.Stop ENTER\n")

	// clean up neo4j driver
	if s.driver != nil {
		ctx := context.Background()
		(*s.driver).Close(ctx)
	}
	s.driver = nil

	// stop this endpoint
	if s.server != nil {
		err := s.server.Close()
		if err != nil {
			klog.V(1).Infof("server.Close() failed. Err: %v\n", err)
		}
	}
	s.server = nil

	klog.V(4).Infof("Server.Stop Succeeded\n")
	klog.V(6).Infof("Server.Stop LEAVE\n")

	return nil
}
