// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfaces

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

type AppSpecificHistorical struct {
	Type string  `json:"type,omitempty"`
	Data []*Data `json:"data,omitempty"`
}

/*
	AppSpecificHistoricalType is just AppSpecificHistorical with the "Data" property removed from it.
	This is simply just to pluck out the "Type" property in order to inspect the value.

	This is only needed when using the Server Sent Events (SSE) mode for Application-level messages
*/
type AppSpecificHistoricalType struct {
	Type string `json:"type,omitempty"`
}
