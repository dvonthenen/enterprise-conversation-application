// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package routing

import (
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	symblinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/*
	Options to init objects...
*/
// MessageRouterOptions to init the router
type MessageRouterOptions struct {
	Callback *symblinterfaces.InsightCallback
}

// MessageHandlerOptions to init the handler
type MessageHandlerOptions struct {
	//housekeeping
	ConversationId string

	// neo4j
	Neo4jMgr *neo4j.SessionWithContext

	// neo4j session
	RabbitMgr *rabbitinterfaces.Manager
}

/*
	The objects...
*/
// MessageRouter converts messages to Symbl objects
type MessageRouter struct {
	callback *symblinterfaces.InsightCallback
}

// MessageHandler takes the Symbl objects and performs an action with them
type MessageHandler struct {
	// general
	conversationId  string
	terminationSent bool

	// neo4j
	neo4jMgr *neo4j.SessionWithContext

	// rabbitmq
	rabbitMgr *rabbitinterfaces.Manager
}
