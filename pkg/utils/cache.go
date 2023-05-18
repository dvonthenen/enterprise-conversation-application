// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package utils

import (
	"container/list"
)

func NewMessageCache() *MessageCache {
	cache := MessageCache{
		rotatingWindowOfMsg: list.New(),
		mapIdToMsg:          make(map[string]*Message),
	}
	return &cache
}

func (mc *MessageCache) Push(msgId, text, userId, name, email string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if len(mc.mapIdToMsg) >= DefaultNumOfMsgToCache {
		e := mc.rotatingWindowOfMsg.Front()
		if e != nil {
			itemMessage := Message(e.Value.(Message))
			delete(mc.mapIdToMsg, itemMessage.ID)
			mc.rotatingWindowOfMsg.Remove(e)
		}
	}

	message := Message{
		ID:   msgId,
		Text: text,
		Author: Author{
			ID:    userId,
			Name:  name,
			Email: email,
		},
	}
	mc.mapIdToMsg[msgId] = &message
	mc.rotatingWindowOfMsg.PushBack(message)

	return nil
}

func (mc *MessageCache) Find(msgId string) (*Message, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	message := mc.mapIdToMsg[msgId]

	if message == nil {
		return nil, ErrItemNotFound
	}

	return message, nil
}
