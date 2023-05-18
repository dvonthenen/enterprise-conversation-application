// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package handlers

import (
	"fmt"

	klog "k8s.io/klog/v2"

	interfacessdk "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/interfaces"
	shared "github.com/dvonthenen/enterprise-reference-implementation/pkg/shared"
	utils "github.com/dvonthenen/enterprise-reference-implementation/pkg/utils"
)

func NewHandler(options HandlerOptions) *Handler {
	handler := Handler{
		session:     options.Session,
		symblClient: options.SymblClient,
		cache:       make(map[string]*utils.MessageCache),
	}
	return &handler
}

func (h *Handler) SetClientPublisher(mp *interfacessdk.MessagePublisher) {
	klog.V(4).Infof("SetClientPublisher called...\n")
	h.msgPublisher = mp
}

func (h *Handler) InitializedConversation(im *shared.InitializationResult) error {
	conversationId := im.InitializationMessage.ConversationID
	klog.V(2).Infof("InitializedConversation - conversationID: %s\n", conversationId)
	h.cache[conversationId] = utils.NewMessageCache()
	return nil
}

func (h *Handler) MessageResult(mr *shared.MessageResult) error {
	// this is needed to keep a message cache which might be valuable for message look up later on

	cache := h.cache[mr.ConversationID]
	if cache != nil {
		for _, msg := range mr.MessageResult.Messages {
			cache.Push(msg.ID, msg.Text, msg.From.ID, msg.From.Name, "")
		}
	} else {
		klog.V(1).Infof("MessageCache for ConversationID(%s) not found.", mr.ConversationID)
	}

	// This is usually only needed for special cases where you are triggering off specific keywords
	// Otherwise, just use the native mechanism for chat message related stuff

	// No implementation required. Return Succeess!

	return nil
}

func (h *Handler) QuestionResult(qr *shared.QuestionResult) error {
	// TODO: if your plugin cares about insights, implement the work below
	// TODO: Otherwise, just return nil

	for _, question := range qr.QuestionResult.Questions {
		// TODO: do some work
		fmt.Printf("QUESTION:\n%v\n", question)
	}

	return nil
}

func (h *Handler) FollowUpResult(fur *shared.FollowUpResult) error {
	// TODO: if your plugin cares about insights, implement the work below
	// TODO: Otherwise, just return nil

	for _, followUp := range fur.FollowUpResult.FollowUps {
		// TODO: do some work
		fmt.Printf("FOLLOWUP:\n%v\n", followUp)
	}

	return nil
}

func (h *Handler) ActionItemResult(air *shared.ActionItemResult) error {
	// TODO: if your plugin cares about insights, implement the work below
	// TODO: Otherwise, just return nil

	for _, actionItem := range air.ActionItemResult.ActionItems {
		// TODO: do some work
		fmt.Printf("ACTIONITEM:\n%v\n", actionItem)
	}

	return nil
}

func (h *Handler) TopicResult(tr *shared.TopicResult) error {
	// TODO: if your plugin cares about topics, implement the work below
	// TODO: Otherwise, just return nil

	for _, topic := range tr.TopicResult.Topics {
		// TODO: do some work
		fmt.Printf("TOPIC:\n%v\n", topic)
	}

	return nil
}

func (h *Handler) TrackerResult(tr *shared.TrackerResult) error {
	// TODO: if your plugin cares about trackers, implement the work below
	// TODO: Otherwise, just return nil

	for _, trackerMatch := range tr.TrackerResult.Matches {
		// TODO: do some work
		fmt.Printf("TRACKER:\n%v\n", trackerMatch)
	}

	return nil
}

func (h *Handler) EntityResult(er *shared.EntityResult) error {
	// TODO: if your plugin cares about entities, implement the work below
	// TODO: Otherwise, just return nil

	for _, entity := range er.EntityResult.Entities {
		// for _, match := range entity.Matches {
		// 	// TODO: do some work
		// }

		// TODO: do some work
		fmt.Printf("ENTITY:\n%s\n", entity)
	}

	return nil
}

func (h *Handler) TeardownConversation(tm *shared.TeardownResult) error {
	conversationId := tm.TeardownMessage.ConversationID
	klog.V(2).Infof("TeardownConversation - conversationID: %s\n", conversationId)
	delete(h.cache, conversationId)
	return nil
}
