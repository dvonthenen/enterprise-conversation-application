// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package analyzer

import (
	"context"
	"fmt"
	"net/http"
	"os"

	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	amqp "github.com/rabbitmq/amqp091-go"
	klog "k8s.io/klog/v2"

	handlers "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/handlers"
	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/interfaces"
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

// TODO: Example... probably should do something better with this
func (se *Server) displayStaticPage(w http.ResponseWriter, r *http.Request) {
	webpageTemplate := `
		<html>
		<head>
		<title>Example</title>
		<script>
		<!--
		function timedRefresh(timeoutPeriod) {
			setTimeout("location.reload(true);",timeoutPeriod);
		}

		window.onload = timedRefresh(5000);
		//   -->
		</script>
		</head>
		<body>
		<p>Data below:</p>
		<p>%s</p>
		</body>
		</html>
	`
	webpageHtml := fmt.Sprintf(webpageTemplate, se.pushData)

	// send page to client
	fmt.Fprintf(w, webpageHtml)
}

func (se *Server) PushNotification(msg string) error {
	se.pushData += msg
	return nil
}

func (se *Server) Init() error {
	// symbl
	err := se.RebuildSymblClient()
	if err != nil {
		klog.V(1).Infof("RebuildSymblClient failed. Err: %v\n", err)
		return err
	}

	// neo4j
	err = se.RebuildDatabase()
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
	mux.HandleFunc("/", se.displayStaticPage)

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

func (se *Server) RebuildSymblClient() error {
	ctx := context.Background()

	symblClient, err := symbl.NewRestClient(ctx)
	if err != nil {
		klog.V(1).Infof("RebuildSymblClient failed. Err: %v\n", err)
		return err
	}
	se.symblClient = symblClient

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
	if se.symblClient == nil {
		err := se.RebuildSymblClient()
		if err != nil {
			klog.V(1).Infof("RebuildSymblClient failed. Err: %v\n", err)
			return err
		}
	}
	rabbitMgr := rabbit.New(rabbit.RabbitManagerOptions{
		Connection: conn,
	})

	var callback interfaces.PushNotificationCallback
	callback = se

	notificationManager := handlers.NewNotificationManager(handlers.NotificationManagerOption{
		Driver:        se.driver,
		RabbitManager: rabbitMgr,
		SymblClient:   se.symblClient,
		PushCallback:  &callback,
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
