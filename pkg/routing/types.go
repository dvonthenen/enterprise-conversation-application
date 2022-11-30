// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package routing

import (
	interfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// MessageHandlerOptions to init the handler
type MessageHandlerOptions struct {
	ConversationId string
	Session        neo4j.SessionWithContext
}

// MessageRouter converts messages to Symbl objects
type MessageRouter struct {
	callback *interfaces.InsightCallback
}

// MessageHandler takes the Symbl objects and performs an action with them
type MessageHandler struct {
	// general
	ConversationId string

	// neo4j
	Session neo4j.SessionWithContext
}
