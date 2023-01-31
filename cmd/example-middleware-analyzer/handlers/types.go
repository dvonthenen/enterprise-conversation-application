// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package handlers

import (
	"container/list"
	"sync"

	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-analyzer/interfaces"
)

/*
	MessageCache
*/
type Message struct {
	ID  string
	Msg string
}

type MessageCache struct {
	rotatingWindowOfMsg *list.List
	mapIdToMsg          map[string]string
	mu                  sync.Mutex
}

/*
	Handler for messages
*/
type HandlerOptions struct {
	Session     *neo4j.SessionWithContext // retrieve insights
	SymblClient *symbl.RestClient
}

type Handler struct {
	// properties
	conversationID string
	cache          *MessageCache

	// housekeeping
	session      *neo4j.SessionWithContext
	symblClient  *symbl.RestClient
	msgPublisher *interfaces.MessagePublisher
}
