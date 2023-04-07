// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import (
	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

/*
	Conversation Insight with Metadata
*/
type InitializationMessage struct {
	sdkinterfaces.InitializationMessage
}

type RecognitionResult struct {
	ConversationID    string                           `json:"conversationId,omitempty"`
	RecognitionResult *sdkinterfaces.RecognitionResult `json:"recognitionResult,omitempty"`
}

type MessageResponse struct {
	ConversationID  string                         `json:"conversationId,omitempty"`
	MessageResponse *sdkinterfaces.MessageResponse `json:"messageResponse,omitempty"`
}

type InsightResponse struct {
	ConversationID  string                         `json:"conversationId,omitempty"`
	InsightResponse *sdkinterfaces.InsightResponse `json:"insightResponse,omitempty"`
}

type TopicResponse struct {
	ConversationID string                       `json:"conversationId,omitempty"`
	TopicResponse  *sdkinterfaces.TopicResponse `json:"topicResponse,omitempty"`
}

type TrackerResponse struct {
	ConversationID  string                         `json:"conversationId,omitempty"`
	TrackerResponse *sdkinterfaces.TrackerResponse `json:"trackerResponse,omitempty"`
}

type EntityResponse struct {
	ConversationID string                        `json:"conversationId,omitempty"`
	EntityResponse *sdkinterfaces.EntityResponse `json:"entityResponse,omitempty"`
}

type TeardownMessage struct {
	sdkinterfaces.TeardownMessage
}
