// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import (
	asyncinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/async/v1/interfaces"
	streaminginterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

/*
	Conversation Insight Responses (RealTime) with Metadata
*/
type InitializationResponse struct {
	InitializationMessage *streaminginterfaces.InitializationMessage `json:"initializationMessage,omitempty"`
}

type RecognitionResponse struct {
	ConversationID    string                                 `json:"conversationId,omitempty"`
	RecognitionResult *streaminginterfaces.RecognitionResult `json:"recognitionResult,omitempty"`
}

type MessageResponse struct {
	ConversationID  string                               `json:"conversationId,omitempty"`
	MessageResponse *streaminginterfaces.MessageResponse `json:"messageResponse,omitempty"`
}

type InsightResponse struct {
	ConversationID  string                               `json:"conversationId,omitempty"`
	InsightResponse *streaminginterfaces.InsightResponse `json:"insightResponse,omitempty"`
}

type TopicResponse struct {
	ConversationID string                             `json:"conversationId,omitempty"`
	TopicResponse  *streaminginterfaces.TopicResponse `json:"topicResponse,omitempty"`
}

type TrackerResponse struct {
	ConversationID  string                               `json:"conversationId,omitempty"`
	TrackerResponse *streaminginterfaces.TrackerResponse `json:"trackerResponse,omitempty"`
}

type EntityResponse struct {
	ConversationID string                              `json:"conversationId,omitempty"`
	EntityResponse *streaminginterfaces.EntityResponse `json:"entityResponse,omitempty"`
}

type TeardownResponse struct {
	TeardownMessage *streaminginterfaces.TeardownMessage `json:"teardownMessage,omitempty"`
}

/*
	Conversation Asynchronous Results (Asynchronous) with Metadata
*/
type InitializationResult struct {
	Duplicate             bool                                   `json:"duplicate,omitempty"`
	InitializationMessage *asyncinterfaces.InitializationMessage `json:"initializationMessage,omitempty"`
}

type MessageResult struct {
	Duplicate      bool                           `json:"duplicate,omitempty"`
	ConversationID string                         `json:"conversationId,omitempty"`
	MessageResult  *asyncinterfaces.MessageResult `json:"messageResult,omitempty"`
}

type QuestionResult struct {
	Duplicate      bool                            `json:"duplicate,omitempty"`
	ConversationID string                          `json:"conversationId,omitempty"`
	QuestionResult *asyncinterfaces.QuestionResult `json:"questionResult,omitempty"`
}

type FollowUpResult struct {
	Duplicate      bool                            `json:"duplicate,omitempty"`
	ConversationID string                          `json:"conversationId,omitempty"`
	FollowUpResult *asyncinterfaces.FollowUpResult `json:"followUpResult,omitempty"`
}

type ActionItemResult struct {
	Duplicate        bool                              `json:"duplicate,omitempty"`
	ConversationID   string                            `json:"conversationId,omitempty"`
	ActionItemResult *asyncinterfaces.ActionItemResult `json:"acitonItemResult,omitempty"`
}

type TopicResult struct {
	Duplicate      bool                         `json:"duplicate,omitempty"`
	ConversationID string                       `json:"conversationId,omitempty"`
	TopicResult    *asyncinterfaces.TopicResult `json:"topicResult,omitempty"`
}

type TrackerResult struct {
	Duplicate      bool                           `json:"duplicate,omitempty"`
	ConversationID string                         `json:"conversationId,omitempty"`
	TrackerResult  *asyncinterfaces.TrackerResult `json:"trackerResult,omitempty"`
}

type EntityResult struct {
	Duplicate      bool                          `json:"duplicate,omitempty"`
	ConversationID string                        `json:"conversationId,omitempty"`
	EntityResult   *asyncinterfaces.EntityResult `json:"entityResult,omitempty"`
}

type TeardownResult struct {
	Duplicate       bool                             `json:"duplicate,omitempty"`
	TeardownMessage *asyncinterfaces.TeardownMessage `json:"teardownResult,omitempty"`
}
