// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package routing

import (
	symblinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	amqp "github.com/rabbitmq/amqp091-go"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/dataminer/interfaces"
)

// MessageHandlerOptions to init the handler
type MessageHandlerOptions struct {
	//housekeeping
	ConversationId string
	NotifyClient   *interfaces.PushNotificationCallback

	// neo4j
	Session *neo4j.SessionWithContext

	// neo4j session
	RabbitConnection *amqp.Connection
}

// MessageRouter converts messages to Symbl objects
type MessageRouter struct {
	callback *symblinterfaces.InsightCallback
}

// MessageHandler takes the Symbl objects and performs an action with them
type MessageHandler struct {
	// general
	ConversationId  string
	notifyClient    *interfaces.PushNotificationCallback
	terminationSent bool

	// neo4j
	session *neo4j.SessionWithContext

	// rabbitmq
	rabbitConnection *amqp.Connection
	rabbitPublish    map[string]*amqp.Channel
}
