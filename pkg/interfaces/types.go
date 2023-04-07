// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

/*
	AppSpecificType is common metadata relating to all application-level messages with
	the "Data" property and aspects removed from it. This is just to pluck out the
	"Type" property to inspect which type of application message this is.

	This struct is needed to:
	- route and differentiate application-level messages in the UserDefinedMessage callback
	- using the Server Sent Events (SSE) mode for application-level messages
*/
type Metadata struct {
	Type string `json:"type"`
}

type AppSpecificType struct {
	Type     string   `json:"type"`
	Metadata Metadata `json:"metadata"`
}

/*
	Shared structs relating to Recognition and Message insights
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

/*
	Structs relating to Recognition and Message insights
*/
type UserDefinedRecognition struct {
	Type        string      `json:"type"`
	Metadata    Metadata    `json:"metadata"`
	Recognition Recognition `json:"recognition,omitempty"`
}

type UserDefinedMessages struct {
	Type     string    `json:"type"`
	Metadata Metadata  `json:"metadata"`
	Messages []Message `json:"message,omitempty"`
}
