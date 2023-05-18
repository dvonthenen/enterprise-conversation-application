// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package handlers

import (
	"encoding/json"
	"fmt"

	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	klog "k8s.io/klog/v2"

	pluginsdkmsg "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
	interfacessdk "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/interfaces"
	shared "github.com/dvonthenen/enterprise-reference-implementation/pkg/shared"
	utils "github.com/dvonthenen/enterprise-reference-implementation/pkg/utils"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/cmd/example-realtime-plugin/interfaces"
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

func (h *Handler) InitializedConversation(im *shared.InitializationResponse) error {
	conversationId := im.InitializationMessage.Message.Data.ConversationID
	klog.V(2).Infof("InitializedConversation - conversationID: %s\n", conversationId)
	h.cache[conversationId] = utils.NewMessageCache()
	return nil
}

func (h *Handler) RecognitionResultMessage(rr *shared.RecognitionResponse) error {
	// This is usually only needed for very very special cases for manipulating the transcription
	// Otherwise, just use the native mechanism for transcription related stuff

	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) MessageResponseMessage(mr *shared.MessageResponse) error {
	// this is needed to keep a message cache which might be valuable for message look up later on
	cache := h.cache[mr.ConversationID]
	if cache != nil {
		for _, msg := range mr.MessageResponse.Messages {
			cache.Push(msg.ID, msg.Payload.Content, msg.From.ID, msg.From.Name, msg.From.UserID)
		}
	} else {
		klog.V(1).Infof("MessageCache for ConversationID(%s) not found.", mr.ConversationID)
	}

	// This is usually only needed for special cases where you are triggering off specific keywords
	// Otherwise, just use the native mechanism for chat message related stuff

	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) InsightResponseMessage(ir *shared.InsightResponse) error {
	// TODO: if your plugin cares about insights, implement the work below
	// TODO: Otherwise, just return nil

	for _, curInsight := range ir.InsightResponse.Insights {
		// TODO: do some work

		// TODO: build your application specific message
		// TODO: this is just an example. Define your own message in ../interfaces/types.go
		msg := &interfaces.AppSpecificScaffold{
			AppSpecificType: &pluginsdkmsg.AppSpecificType{
				Type: shared.MessageTypeUserDefined,
				Metadata: pluginsdkmsg.Metadata{
					Type: interfaces.AppSpecificMessageTypeScaffold,
				},
			},
			Data: interfaces.Data{
				Type: interfaces.UserScaffoldTypeInsight,
				Data: fmt.Sprintf("Hey! Someone mentioned a \"%s\" with this content: \"%s\"", curInsight.Type, curInsight.Payload.Content),
			},
		}

		// convert to JSON
		data, err := json.Marshal(*msg)
		if err != nil {
			klog.V(1).Infof("[Insight] json.Marshal failed. Err: %v\n", err)
		}

		// send the message to the client
		err = (*h.msgPublisher).PublishMessage(ir.ConversationID, data)
		if err != nil {
			klog.V(1).Infof("[Insight] PublishMessage failed. Err: %v\n", err)
		}
	}

	return nil
}

func (h *Handler) TopicResponseMessage(tr *shared.TopicResponse) error {
	// TODO: if your plugin cares about topics, implement the work below
	// TODO: Otherwise, just return nil

	for _, curTopic := range tr.TopicResponse.Topics {
		// TODO: do some work

		// TODO: build your application specific message
		// TODO: this is just an example. Define your own message in ../interfaces/types.go
		msg := &interfaces.AppSpecificScaffold{
			AppSpecificType: &pluginsdkmsg.AppSpecificType{
				Type: shared.MessageTypeUserDefined,
				Metadata: pluginsdkmsg.Metadata{
					Type: interfaces.AppSpecificMessageTypeScaffold,
				},
			},
			Data: interfaces.Data{
				Type: interfaces.UserScaffoldTypeTopic,
				Data: fmt.Sprintf("Hey! Someone mentioned this Topic \"%s\" with this context: %s", curTopic.Phrases, h.convertMessageReferenceToSlice(tr.ConversationID, curTopic.MessageReferences)),
			},
		}

		// convert to JSON
		data, err := json.Marshal(*msg)
		if err != nil {
			klog.V(1).Infof("[Topics] json.Marshal failed. Err: %v\n", err)
		}

		// send the message to the client
		err = (*h.msgPublisher).PublishMessage(tr.ConversationID, data)
		if err != nil {
			klog.V(1).Infof("[Topics] PublishMessage failed. Err: %v\n", err)
		}
	}

	return nil
}

func (h *Handler) TrackerResponseMessage(tr *shared.TrackerResponse) error {
	// TODO: if your plugin cares about trackers, implement the work below
	// TODO: Otherwise, just return nil

	for _, curTracker := range tr.TrackerResponse.Trackers {
		for _, match := range curTracker.Matches {
			// TODO: do some work

			// TODO: build your application specific message
			// TODO: this is just an example. Define your own message in ../interfaces/types.go
			msg := &interfaces.AppSpecificScaffold{
				AppSpecificType: &pluginsdkmsg.AppSpecificType{
					Type: shared.MessageTypeUserDefined,
					Metadata: pluginsdkmsg.Metadata{
						Type: interfaces.AppSpecificMessageTypeScaffold,
					},
				},
				Data: interfaces.Data{
					Type: interfaces.UserScaffoldTypeTracker,
					Data: fmt.Sprintf("Hey! Someone mentioned this Tracker \"%s\" with this context: %s", curTracker.Name, h.convertMessageRefsToSlice(tr.ConversationID, match.MessageRefs)),
				},
			}

			// convert to JSON
			data, err := json.Marshal(*msg)
			if err != nil {
				klog.V(1).Infof("[Tracker] json.Marshal failed. Err: %v\n", err)
			}

			// send the message to the client
			err = (*h.msgPublisher).PublishMessage(tr.ConversationID, data)
			if err != nil {
				klog.V(1).Infof("[Tracker] PublishMessage failed. Err: %v\n", err)
			}
		}
	}

	return nil
}

func (h *Handler) EntityResponseMessage(er *shared.EntityResponse) error {
	// TODO: if your plugin cares about entities, implement the work below
	// TODO: Otherwise, just return nil

	for _, entity := range er.EntityResponse.Entities {
		for _, match := range entity.Matches {
			// TODO: do some work

			// TODO: build your application specific message
			// TODO: this is just an example. Define your own message in ../interfaces/types.go
			msg := &interfaces.AppSpecificScaffold{
				AppSpecificType: &pluginsdkmsg.AppSpecificType{
					Type: shared.MessageTypeUserDefined,
					Metadata: pluginsdkmsg.Metadata{
						Type: interfaces.AppSpecificMessageTypeScaffold,
					},
				},
				Data: interfaces.Data{
					Type: interfaces.UserScaffoldTypeEntity,
					Data: fmt.Sprintf("Hey! Someone mentioned this Entity \"%s/%s/%s/%s\" with this context: %s", entity.Category, entity.Type, entity.SubType, match.DetectedValue, h.convertMessageRefsToSlice(er.ConversationID, match.MessageRefs)),
				},
			}

			// convert to JSON
			data, err := json.Marshal(*msg)
			if err != nil {
				klog.V(1).Infof("[Entity] json.Marshal failed. Err: %v\n", err)
			}

			// send the message to the client
			err = (*h.msgPublisher).PublishMessage(er.ConversationID, data)
			if err != nil {
				klog.V(1).Infof("[Entity] PublishMessage failed. Err: %v\n", err)
			}
		}
	}

	return nil
}

func (h *Handler) TeardownConversation(tm *shared.TeardownResponse) error {
	conversationId := tm.TeardownMessage.Message.Data.ConversationID
	klog.V(2).Infof("TeardownConversation - conversationID: %s\n", conversationId)
	delete(h.cache, conversationId)
	return nil
}

func (h *Handler) UserDefinedMessage(data []byte) error {
	// This is only needed on the client side and not on the plugin side.
	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) UnhandledMessage(byMsg []byte) error {
	klog.Errorf("\n\n-------------------------------\n")
	klog.Errorf("UnhandledMessage:\n%v\n", string(byMsg))
	klog.Errorf("-------------------------------\n\n")
	return ErrUnhandledMessage
}

func (h *Handler) convertMessageReferenceToSlice(conversationId string, msgRefs []sdkinterfaces.MessageReference) []string {
	tmp := make([]string, 0)

	cache := h.cache[conversationId]
	if cache == nil {
		tmp = append(tmp, interfaces.MessageNotFound)
		return tmp
	}

	for _, msgRef := range msgRefs {
		cacheMessage, err := cache.Find(msgRef.ID)
		if err != nil {
			klog.V(4).Infof("Msg ID not found: %s\n", msgRef.ID)
			tmp = append(tmp, interfaces.MessageNotFound)
			continue
		}

		tmp = append(tmp, cacheMessage.Text)
	}

	return tmp
}

func (h *Handler) convertMessageRefsToSlice(conversationId string, msgRefs []sdkinterfaces.MessageRef) []string {
	tmp := make([]string, 0)

	cache := h.cache[conversationId]
	if cache == nil {
		tmp = append(tmp, interfaces.MessageNotFound)
		return tmp
	}

	for _, msgRef := range msgRefs {
		cacheMessage, err := cache.Find(msgRef.ID)
		if err != nil {
			klog.V(4).Infof("Msg ID not found: %s\n", msgRef.ID)
			tmp = append(tmp, interfaces.MessageNotFound)
			continue
		}

		tmp = append(tmp, cacheMessage.Text)
	}

	return tmp
}
