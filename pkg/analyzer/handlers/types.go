// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/interfaces"
	rabbit "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit"
)

/*
	Subscriber handlers
*/
type HandlerOptions struct {
	Session      *neo4j.SessionWithContext // retrieve insights
	SymblClient  *symbl.RestClient
	PushCallback *interfaces.PushNotificationCallback
}

type ConversationInitHandler struct {
	session      *neo4j.SessionWithContext
	symblClient  *symbl.RestClient
	pushCallback *interfaces.PushNotificationCallback
}

type EntityHandler struct {
	session      *neo4j.SessionWithContext
	symblClient  *symbl.RestClient
	pushCallback *interfaces.PushNotificationCallback
}

type InsightHandler struct {
	session      *neo4j.SessionWithContext
	symblClient  *symbl.RestClient
	pushCallback *interfaces.PushNotificationCallback
}

type MessageHandler struct {
	session      *neo4j.SessionWithContext
	symblClient  *symbl.RestClient
	pushCallback *interfaces.PushNotificationCallback
}

type TopicHandler struct {
	session      *neo4j.SessionWithContext
	symblClient  *symbl.RestClient
	pushCallback *interfaces.PushNotificationCallback
}

type TrackerHandler struct {
	session      *neo4j.SessionWithContext
	symblClient  *symbl.RestClient
	pushCallback *interfaces.PushNotificationCallback
}

type ConversationTeardownHandler struct {
	session      *neo4j.SessionWithContext
	symblClient  *symbl.RestClient
	pushCallback *interfaces.PushNotificationCallback
}

/*
	Notification mmanager
*/
type NotificationManagerOption struct {
	Driver        *neo4j.DriverWithContext
	RabbitManager *rabbit.RabbitManager
	SymblClient   *symbl.RestClient
	PushCallback  *interfaces.PushNotificationCallback
}

type NotificationManager struct {
	// neo4j
	driver *neo4j.DriverWithContext

	// symbl
	symblClient *symbl.RestClient

	// rabbit
	rabbitManager *rabbit.RabbitManager

	// TODO: Example... probably should do something better with this
	pushCallback *interfaces.PushNotificationCallback
}
