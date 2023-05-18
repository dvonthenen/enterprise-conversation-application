// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package utils

import (
	"container/list"
	"sync"
)

/*
	MessageCache
*/
type Author struct {
	ID    string
	Name  string
	Email string
}
type Message struct {
	ID     string
	Text   string
	Author Author
}

type MessageCache struct {
	rotatingWindowOfMsg *list.List
	mapIdToMsg          map[string]*Message
	mu                  sync.Mutex
}
