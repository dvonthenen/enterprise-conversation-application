// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package dataminer

import (
	"net/http"
	"sync"
	"time"

	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	instance "github.com/dvonthenen/enterprise-reference-implementation/pkg/proxy-dataminer/instance"
)

// Credentials is the input needed to login to neo4j
type Credentials struct {
	ConnectionStr string
	Username      string
	Password      string
}

// ServerOptions for the main HTTP endpoint
type ServerOptions struct {
	CrtFile   string
	KeyFile   string
	StartPort int
	EndPort   int
	RabbitURI string
}

type Server struct {
	// server versions
	options ServerOptions
	creds   Credentials

	// bookkeeping
	instanceById   map[string]*instance.Proxy
	instanceByPort map[int]*instance.Proxy
	server         *http.Server
	mu             sync.Mutex
	ticker         *time.Ticker
	stopPoll       chan struct{}

	// neo4j
	driver *neo4j.DriverWithContext
}
