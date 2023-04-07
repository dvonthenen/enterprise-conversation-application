// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

const (
	// rabbit message names/exchanges
	RabbitExchangeConversationInit     string = "conversation-created"
	RabbitExchangeMessage              string = "message-created"
	RabbitExchangeTopic                string = "topic-created"
	RabbitExchangeTracker              string = "tracker-created"
	RabbitExchangeEntity               string = "entity-created"
	RabbitExchangeInsight              string = "insight-created"
	RabbitExchangeConversationTeardown string = "conversation-teardown"
	RabbitClientNotifications          string = "client-notification"

	// neo4j node ID names/index/uniqueIds
	DatabaseIndexConversation string = "conversationId"
	DatabaseIndexMessage      string = "messageId"
	DatabaseIndexUser         string = "userId"
	DatabaseIndexTopic        string = "topicId"
	DatabaseIndexTracker      string = "trackerId"
	DatabaseIndexInsight      string = "insightId"
	DatabaseIndexEntity       string = "entityId" // = entity.Type + "_" + entity.SubType + "_" + entity.Category
	DatabaseIndexEntityMatch  string = "matchId"  // = conversationId + "_" + entityId
)
