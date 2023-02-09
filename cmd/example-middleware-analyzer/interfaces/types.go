// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

/*
	Higher Level Application Message
*/
type Author struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type Message struct {
	ID     string `json:"id,omitempty"`
	Text   string `json:"text,omitempty"`
	Author Author `json:"author,omitempty"`
}

type Insight struct {
	Correlation     string `json:"correlation,omitempty"`
	Messages    []Message `json:"message,omitempty"`
}

type Data struct {
	Type    string  `json:"type,omitempty"`
	Correlation string    `json:"correlation,omitempty"`
	Current     []Insight `json:"current,omitempty"`
	Previous    []Insight `json:"previous,omitempty"`
}

type Metadata struct {
	Type string `json:"type"`
}

type AppSpecificHistorical struct {
	Type       string   `json:"type"`
	Metadata   Metadata `json:"metadata"`
	Historical Data     `json:"historical,omitempty"`
}
