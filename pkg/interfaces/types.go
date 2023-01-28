// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfaces

import (
	interfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

/*
	Conversation Insight with Metadata
*/
type InitializationMessage struct {
	interfaces.InitializationMessage
}

type RecognitionResult struct {
	ConversationID    string                        `json:"conversationId,omitempty"`
	RecognitionResult *interfaces.RecognitionResult `json:"recognitionResult,omitempty"`
}

type MessageResponse struct {
	ConversationID  string                      `json:"conversationId,omitempty"`
	MessageResponse *interfaces.MessageResponse `json:"messageResponse,omitempty"`
}

type InsightResponse struct {
	ConversationID  string                      `json:"conversationId,omitempty"`
	InsightResponse *interfaces.InsightResponse `json:"insightResponse,omitempty"`
}

type TopicResponse struct {
	ConversationID string                    `json:"conversationId,omitempty"`
	TopicResponse  *interfaces.TopicResponse `json:"topicResponse,omitempty"`
}

type TrackerResponse struct {
	ConversationID  string                      `json:"conversationId,omitempty"`
	TrackerResponse *interfaces.TrackerResponse `json:"trackerResponse,omitempty"`
}

type EntityResponse struct {
	ConversationID string                     `json:"conversationId,omitempty"`
	EntityResponse *interfaces.EntityResponse `json:"entityResponse,omitempty"`
}

type TeardownMessage struct {
	interfaces.TeardownMessage
}
