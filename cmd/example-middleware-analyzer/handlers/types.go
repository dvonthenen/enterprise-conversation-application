// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package handlers

import (
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-analyzer/interfaces"
	utils "github.com/dvonthenen/enterprise-reference-implementation/pkg/utils"
)

/*
	Handler for messages
*/
type HandlerOptions struct {
	Session     *neo4j.SessionWithContext // retrieve insights
	SymblClient *symbl.RestClient
}

type Handler struct {
	// properties
	conversationID string
	cache          *utils.MessageCache

	// housekeeping
	session      *neo4j.SessionWithContext
	symblClient  *symbl.RestClient
	msgPublisher *interfaces.MessagePublisher
}
