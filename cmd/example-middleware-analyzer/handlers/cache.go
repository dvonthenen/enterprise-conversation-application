// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package handlers

import (
	"container/list"
)

func NewMessageCache() *MessageCache {
	cache := MessageCache{
		rotatingWindowOfMsg: list.New(),
		mapIdToMsg:          make(map[string]string),
	}
	return &cache
}

func (mc *MessageCache) Push(ID, message string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if len(mc.mapIdToMsg) >= DefaultNumOfMsgToCache {
		e := mc.rotatingWindowOfMsg.Front()
		if e != nil {
			itemMessage := Message(e.Value.(Message))
			delete(mc.mapIdToMsg, itemMessage.ID)
		}
		mc.rotatingWindowOfMsg.Remove(e)
	}

	mc.mapIdToMsg[ID] = message
	mc.rotatingWindowOfMsg.PushBack(Message{
		ID:  ID,
		Msg: message,
	})

	return nil
}

func (mc *MessageCache) Find(ID string) (string, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	msg := mc.mapIdToMsg[ID]

	if msg == "" {
		return "", ErrItemNotFound
	}

	return msg, nil
}
