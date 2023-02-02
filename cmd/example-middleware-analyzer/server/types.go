// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package server

import (
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	middleware "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-analyzer"
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
	RabbitURI   string
}

type Server struct {
	// server versions
	options ServerOptions
	creds   Credentials

	// middleware
	middlewareAnalyzer *middleware.MiddlewareAnalyzer

	// neo4j
	driver *neo4j.DriverWithContext

	// symbl client
	symblClient *symbl.RestClient
}
