// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

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

type Metadata struct {
	Type string `json:"type"`
}

type AppSpecificHistorical struct {
	Type       string   `json:"type"`
	Metadata   Metadata `json:"metadata"`
	Historical []Data   `json:"historical,omitempty"`
}
