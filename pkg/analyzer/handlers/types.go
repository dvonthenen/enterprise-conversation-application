// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	rabbit "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit"
)

/*
	Subscriber handlers
*/
type HandlerOptions struct {
	Session *neo4j.SessionWithContext
}

type ConversationHandler struct {
	session *neo4j.SessionWithContext
}

type EntityHandler struct {
	session *neo4j.SessionWithContext
}

type InsightHandler struct {
	session *neo4j.SessionWithContext
}

type MessageHandler struct {
	session *neo4j.SessionWithContext
}

type TopicHandler struct {
	session *neo4j.SessionWithContext
}

type TrackerHandler struct {
	session *neo4j.SessionWithContext
}

/*
	Notification mmanager
*/
type NotificationManagerOption struct {
	Driver        *neo4j.DriverWithContext
	RabbitManager *rabbit.RabbitManager
}

type NotificationManager struct {
	// neo4j
	driver *neo4j.DriverWithContext

	// rabbit
	rabbitManager *rabbit.RabbitManager
}
