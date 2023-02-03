// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

/*
	Shared structs
*/
type From struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	UserID string `json:"userId,omitempty"`
}

type Duration struct {
	StartTime  string  `json:"startTime,omitempty"`
	EndTime    string  `json:"endTime,omitempty"`
	TimeOffset float64 `json:"timeOffset,omitempty"`
	Duration   float64 `json:"duration,omitempty"`
}

type Recognition struct {
	Type    string `json:"type,omitempty"`
	IsFinal bool   `json:"isFinal,omitempty"`
	From    From   `json:"user,omitempty"`
	Content string `json:"transcript,omitempty"`
}

type Message struct {
	From     From     `json:"from,omitempty"`
	Content  string   `json:"content,omitempty"`
	ID       string   `json:"id,omitempty"`
	Duration Duration `json:"duration,omitempty"`
}

type Fragment struct {
	Type     string    `json:"type"`
	Messages []Message `json:"message,omitempty"`
}

/*
	Conversation Insight with Metadata
*/
type UserDefinedRecognition struct {
	Type        string      `json:"type,omitempty"`
	Recognition Recognition `json:"recognition,omitempty"`
}

type UserDefinedMessages struct {
	Type     string   `json:"type"`
	Fragment Fragment `json:"fragment,omitempty"`
}
