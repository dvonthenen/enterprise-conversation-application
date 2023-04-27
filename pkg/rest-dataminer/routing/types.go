// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package routing

import (
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/*
	MessageHandler objects...
*/
// MessageHandlerOptions to init the handler
type MessageHandlerOptions struct {
	//housekeeping
	ConversationId string

	// neo4j
	Neo4jMgr *neo4j.SessionWithContext

	// neo4j session
	RabbitMgr *rabbitinterfaces.Manager
}

// MessageHandler takes the Symbl objects and performs an action with them
type MessageHandler struct {
	// general
	conversationId     string
	terminationSent    bool
	conversationExists bool
	existsSet          bool

	// features
	options MessageHandlerOptions

	// neo4j
	neo4jMgr *neo4j.SessionWithContext

	// rabbitmq
	rabbitMgr *rabbitinterfaces.Manager
}
