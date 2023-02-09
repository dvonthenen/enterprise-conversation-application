// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package utils

import (
	"container/list"
	"sync"
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
