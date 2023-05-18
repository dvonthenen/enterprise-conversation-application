// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package dataminer

import (
	"net/http"
	"sync"

	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Credentials is the input needed to login to neo4j
type Credentials struct {
	ConnectionStr string
	Username      string
	Password      string
}

// ServerOptions for the main HTTP endpoint
type ServerOptions struct {
	AuthMethod       AuthType
	CrtFile          string
	KeyFile          string
	RabbitURI        string
	DisableDuplicate bool
}

type Server struct {
	// server versions
	options ServerOptions
	creds   Credentials

	// bookkeeping
	server *http.Server
	mu     sync.Mutex

	// neo4j
	driver *neo4j.DriverWithContext
}
