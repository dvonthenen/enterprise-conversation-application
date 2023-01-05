// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package analyzer

import (
	rabbit "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	handlers "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/handlers"
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

	// bookkeeping
	rabbitMgr       *rabbit.Manager
	notificationMgr *handlers.NotificationManager

	// neo4j
	driver *neo4j.DriverWithContext

	// symbl client
	symblClient *symbl.RestClient
}
