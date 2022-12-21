// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	rabbit "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit"
	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
)

/*
	Subscriber handlers
*/
type HandlerOptions struct {
	Session     *neo4j.SessionWithContext // retrieve insights
	SymblClient *symbl.RestClient
	Manager     *interfaces.RabbitManagerHandler
}

type ConversationInitHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *interfaces.RabbitManagerHandler
}

type EntityHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *interfaces.RabbitManagerHandler
}

type InsightHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *interfaces.RabbitManagerHandler
}

type MessageHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *interfaces.RabbitManagerHandler
}

type TopicHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *interfaces.RabbitManagerHandler
}

type TrackerHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *interfaces.RabbitManagerHandler
}

type ConversationTeardownHandler struct {
	session     *neo4j.SessionWithContext
	symblClient *symbl.RestClient
	manager     *interfaces.RabbitManagerHandler
}

/*
	Notification mmanager
*/
type NotificationManagerOption struct {
	Driver        *neo4j.DriverWithContext
	RabbitManager *rabbit.RabbitManager
	SymblClient   *symbl.RestClient
}

type NotificationManager struct {
	// neo4j
	driver *neo4j.DriverWithContext

	// symbl
	symblClient *symbl.RestClient

	// rabbit
	rabbitManager *rabbit.RabbitManager
}
