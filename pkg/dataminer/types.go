// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package dataminer

import (
	"net/http"
	"sync"
	"time"

	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	amqp "github.com/rabbitmq/amqp091-go"

	instance "github.com/dvonthenen/enterprise-reference-implementation/pkg/dataminer/instance"
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
	StartPort   int
	EndPort     int
	RabbitMQURI string
}

type Server struct {
	// server versions
	options ServerOptions
	creds   Credentials

	// bookkeeping
	instanceById   map[string]*instance.ServerInstance
	instanceByPort map[int]*instance.ServerInstance
	server         *http.Server
	mu             sync.Mutex
	ticker         *time.Ticker
	stopPoll       chan struct{}

	// neo4j
	driver *neo4j.DriverWithContext

	// rabbitmq
	rabbitConn *amqp.Connection
}
