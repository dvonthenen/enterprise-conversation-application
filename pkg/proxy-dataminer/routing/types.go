// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package routing

import (
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
)

/*
	MessagePassthrough used mainly for chat and closed captioning
*/
type MessagePassthrough interface {
	SendRecognition(r *interfaces.UserDefinedRecognition) error
	SendMessages(m *interfaces.UserDefinedMessages) error
}

/*
	MessageRouter objects...
*/
// MessageRouterOptions to init the router
type MessageRouterOptions struct {
	Callback *sdkinterfaces.InsightCallback
}

// MessageRouter converts messages to Symbl objects
type MessageRouter struct {
	callback *sdkinterfaces.InsightCallback
}

/*
	MessageHandler objects...
*/

// MessageHandlerOptions to init the handler
type MessageHandlerOptions struct {
	//housekeeping
	ConversationId string

	// features
	TranscriptionEnabled bool
	MessagingEnabled     bool

	// callback
	Callback *MessagePassthrough

	// neo4j
	Neo4jMgr *neo4j.SessionWithContext

	// neo4j session
	RabbitMgr *rabbitinterfaces.Manager
}

// MessageHandler takes the Symbl objects and performs an action with them
type MessageHandler struct {
	// general
	conversationId  string
	terminationSent bool

	// features
	options MessageHandlerOptions

	// callback
	callback *MessagePassthrough

	// neo4j
	neo4jMgr *neo4j.SessionWithContext

	// rabbitmq
	rabbitMgr *rabbitinterfaces.Manager
}
