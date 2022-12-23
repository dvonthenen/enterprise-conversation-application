// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfaces

import (
	interfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

/*
	Higher Level Application Message
*/
type Author struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type Message struct {
	Correlation     string `json:"correlation,omitempty"`
	CurrentContent  string `json:"currentContent,omitempty"`
	CurrentMatch    string `json:"currentMatch,omitempty"`
	PreviousContent string `json:"previousContent,omitempty"`
	PreviousMatch   string `json:"previousMatch,omitempty"`
}

type Data struct {
	Type    string  `json:"type,omitempty"`
	Author  Author  `json:"author,omitempty"`
	Message Message `json:"message,omitempty"`
}

type ClientTrackerMessage struct {
	Type string `json:"type,omitempty"`
	Data Data   `json:"data,omitempty"`
}

type ClientMessageType struct {
	Type string `json:"type,omitempty"`
}

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
	ConversationID string              `json:"conversationId,omitempty"`
	Insight        *interfaces.Insight `json:"insight,omitempty"`
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
