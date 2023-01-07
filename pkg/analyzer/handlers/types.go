// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
)

/*
	Subscriber handlers
*/
type HandlerOptions struct {
	Session     *neo4j.SessionWithContext // retrieve insights
	SymblClient *symbl.RestClient
	Manager     *rabbitinterfaces.Manager
}

type ConversationInitHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *rabbitinterfaces.Manager
}

type EntityHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *rabbitinterfaces.Manager
}

type InsightHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *rabbitinterfaces.Manager
}

type MessageHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *rabbitinterfaces.Manager
}

type TopicHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *rabbitinterfaces.Manager
}

type TrackerHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *rabbitinterfaces.Manager
}

type ConversationTeardownHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *rabbitinterfaces.Manager
}

/*
	Notification mmanager
*/
type NotificationManagerOption struct {
	Driver        *neo4j.DriverWithContext
	RabbitManager *rabbitinterfaces.Manager
	SymblClient   *symbl.RestClient
}

type NotificationManager struct {
	// neo4j
	driver *neo4j.DriverWithContext

	// symbl
	symblClient *symbl.RestClient

	// rabbit
	rabbitManager *rabbitinterfaces.Manager
}
