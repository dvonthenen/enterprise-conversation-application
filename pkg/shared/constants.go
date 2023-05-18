// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

const (
	// rabbit message names/exchanges
	RabbitRealTimeConversationInit     string = "realtime-conversation-created"
	RabbitRealTimeMessage              string = "realtime-message-created"
	RabbitRealTimeTopic                string = "realtime-topic-created"
	RabbitRealTimeTracker              string = "realtime-tracker-created"
	RabbitRealTimeEntity               string = "realtime-entity-created"
	RabbitRealTimeInsight              string = "realtime-insight-created"
	RabbitRealTimeConversationTeardown string = "realtime-conversation-teardown"
	RabbitRealTimeClientNotifications  string = "realtime-client-notification"

	RabbitAsyncConversationInit     string = "async-conversation-created"
	RabbitAsyncMessage              string = "async-message-created"
	RabbitAsyncQuestion             string = "async-question-created"
	RabbitAsyncFollowUp             string = "async-followup-created"
	RabbitAsyncActionItem           string = "async-actionitem-created"
	RabbitAsyncTopic                string = "async-topic-created"
	RabbitAsyncTracker              string = "async-tracker-created"
	RabbitAsyncEntity               string = "async-entity-created"
	RabbitAsyncConversationTeardown string = "async-conversation-teardown"

	// user-defined messages
	MessageTypeUserDefined string = "user_defined"

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
