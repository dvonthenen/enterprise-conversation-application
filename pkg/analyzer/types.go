// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package analyzer

import (
	"net/http"

	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	amqp "github.com/rabbitmq/amqp091-go"

	handlers "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/handlers"
	rabbit "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit"
)

// Credentials is the input needed to login to neo4j
type Credentials struct {
	ConnectionStr string
	Username      string
	Password      string
}

// ServerOptions for the main HTTP endpoint
type ServerOptions struct {
	CrtFile     string
	KeyFile     string
	BindAddress string
	BindPort    int
	RabbitMQURI string
}

type Server struct {
	// server versions
	options ServerOptions
	creds   Credentials

	// bookkeeping
	server          *http.Server
	rabbitMgr       *rabbit.RabbitManager
	notificationMgr *handlers.NotificationManager

	// neo4j
	driver *neo4j.DriverWithContext

	// rabbitmq
	rabbitConn *amqp.Connection
}
