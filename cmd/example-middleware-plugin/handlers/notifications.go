// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package handlers

import (
	"encoding/json"
	"fmt"

	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	klog "k8s.io/klog/v2"

	pluginsdkmsg "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
	interfacessdk "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/interfaces"
	utils "github.com/dvonthenen/enterprise-reference-implementation/pkg/utils"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/cmd/example-middleware-plugin/interfaces"
)

func NewHandler(options HandlerOptions) *Handler {
	handler := Handler{
		session:     options.Session,
		symblClient: options.SymblClient,
		cache:       utils.NewMessageCache(),
	}
	return &handler
}

func (h *Handler) SetClientPublisher(mp *interfacessdk.MessagePublisher) {
	klog.V(4).Infof("SetClientPublisher called...\n")
	h.msgPublisher = mp
}

func (h *Handler) InitializedConversation(im *sdkinterfaces.InitializationMessage) error {
	h.conversationID = im.Message.Data.ConversationID
	klog.V(2).Infof("conversationID: %s\n", h.conversationID)
	return nil
}

func (h *Handler) RecognitionResultMessage(rr *sdkinterfaces.RecognitionResult) error {
	// This is usually only needed for very very special cases for manipulating the transcription
	// Otherwise, just use the native mechanism for transcription related stuff

	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) MessageResponseMessage(mr *sdkinterfaces.MessageResponse) error {
	// this is needed to keep a message cache which might be valuable for message look up later on
	for _, msg := range mr.Messages {
		h.cache.Push(msg.ID, msg.Payload.Content, msg.From.ID, msg.From.Name, msg.From.UserID)
	}

	// This is usually only needed for special cases where you are triggering off specific keywords
	// Otherwise, just use the native mechanism for chat message related stuff

	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) InsightResponseMessage(ir *sdkinterfaces.InsightResponse) error {
	// TODO: if your plugin cares about insights, implement the work below
	// TODO: Otherwise, just return nil

	for _, curInsight := range ir.Insights {
		// TODO: do some work

		// TODO: build your application specific message
		// TODO: this is just an example. Define your own message in ../interfaces/types.go
		msg := &interfaces.AppSpecificScaffold{
			AppSpecificType: &pluginsdkmsg.AppSpecificType{
				Type: sdkinterfaces.MessageTypeUserDefined,
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
		err = (*h.msgPublisher).PublishMessage(h.conversationID, data)
		if err != nil {
			klog.V(1).Infof("[Insight] PublishMessage failed. Err: %v\n", err)
		}
	}

	return nil
}

func (h *Handler) TopicResponseMessage(tr *sdkinterfaces.TopicResponse) error {
	// TODO: if your plugin cares about topics, implement the work below
	// TODO: Otherwise, just return nil

	for _, curTopic := range tr.Topics {
		// TODO: do some work

		// TODO: build your application specific message
		// TODO: this is just an example. Define your own message in ../interfaces/types.go
		msg := &interfaces.AppSpecificScaffold{
			AppSpecificType: &pluginsdkmsg.AppSpecificType{
				Type: sdkinterfaces.MessageTypeUserDefined,
				Metadata: pluginsdkmsg.Metadata{
					Type: interfaces.AppSpecificMessageTypeScaffold,
				},
			},
			Data: interfaces.Data{
				Type: interfaces.UserScaffoldTypeTopic,
				Data: fmt.Sprintf("Hey! Someone mentioned this Topic \"%s\" with this context: %s", curTopic.Phrases, h.convertMessageReferenceToSlice(curTopic.MessageReferences)),
			},
		}

		// convert to JSON
		data, err := json.Marshal(*msg)
		if err != nil {
			klog.V(1).Infof("[Topics] json.Marshal failed. Err: %v\n", err)
		}

		// send the message to the client
		err = (*h.msgPublisher).PublishMessage(h.conversationID, data)
		if err != nil {
			klog.V(1).Infof("[Topics] PublishMessage failed. Err: %v\n", err)
		}
	}

	return nil
}

func (h *Handler) TrackerResponseMessage(tr *sdkinterfaces.TrackerResponse) error {
	// TODO: if your plugin cares about trackers, implement the work below
	// TODO: Otherwise, just return nil

	for _, curTracker := range tr.Trackers {
		for _, match := range curTracker.Matches {
			// TODO: do some work

			// TODO: build your application specific message
			// TODO: this is just an example. Define your own message in ../interfaces/types.go
			msg := &interfaces.AppSpecificScaffold{
				AppSpecificType: &pluginsdkmsg.AppSpecificType{
					Type: sdkinterfaces.MessageTypeUserDefined,
					Metadata: pluginsdkmsg.Metadata{
						Type: interfaces.AppSpecificMessageTypeScaffold,
					},
				},
				Data: interfaces.Data{
					Type: interfaces.UserScaffoldTypeTracker,
					Data: fmt.Sprintf("Hey! Someone mentioned this Tracker \"%s\" with this context: %s", curTracker.Name, h.convertMessageRefsToSlice(match.MessageRefs)),
				},
			}

			// convert to JSON
			data, err := json.Marshal(*msg)
			if err != nil {
				klog.V(1).Infof("[Tracker] json.Marshal failed. Err: %v\n", err)
			}

			// send the message to the client
			err = (*h.msgPublisher).PublishMessage(h.conversationID, data)
			if err != nil {
				klog.V(1).Infof("[Tracker] PublishMessage failed. Err: %v\n", err)
			}
		}
	}

	return nil
}

func (h *Handler) EntityResponseMessage(er *sdkinterfaces.EntityResponse) error {
	// TODO: if your plugin cares about entities, implement the work below
	// TODO: Otherwise, just return nil

	for _, entity := range er.Entities {
		for _, match := range entity.Matches {
			// TODO: do some work

			// TODO: build your application specific message
			// TODO: this is just an example. Define your own message in ../interfaces/types.go
			msg := &interfaces.AppSpecificScaffold{
				AppSpecificType: &pluginsdkmsg.AppSpecificType{
					Type: sdkinterfaces.MessageTypeUserDefined,
					Metadata: pluginsdkmsg.Metadata{
						Type: interfaces.AppSpecificMessageTypeScaffold,
					},
				},
				Data: interfaces.Data{
					Type: interfaces.UserScaffoldTypeEntity,
					Data: fmt.Sprintf("Hey! Someone mentioned this Entity \"%s/%s/%s/%s\" with this context: %s", entity.Category, entity.Type, entity.SubType, match.DetectedValue, h.convertMessageRefsToSlice(match.MessageRefs)),
				},
			}

			// convert to JSON
			data, err := json.Marshal(*msg)
			if err != nil {
				klog.V(1).Infof("[Entity] json.Marshal failed. Err: %v\n", err)
			}

			// send the message to the client
			err = (*h.msgPublisher).PublishMessage(h.conversationID, data)
			if err != nil {
				klog.V(1).Infof("[Entity] PublishMessage failed. Err: %v\n", err)
			}
		}
	}

	return nil
}

func (h *Handler) TeardownConversation(tm *sdkinterfaces.TeardownMessage) error {
	// No implementation required. Return Succeess!
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

func (h *Handler) convertMessageReferenceToSlice(msgRefs []sdkinterfaces.MessageReference) []string {
	tmp := make([]string, 0)

	for _, msgRef := range msgRefs {
		cacheMessage, err := h.cache.Find(msgRef.ID)
		if err != nil {
			klog.V(4).Infof("Msg ID not found: %s\n", msgRef.ID)
			tmp = append(tmp, interfaces.MessageNotFound)
			continue
		}

		tmp = append(tmp, cacheMessage.Text)
	}

	return tmp
}

func (h *Handler) convertMessageRefsToSlice(msgRefs []sdkinterfaces.MessageRef) []string {
	tmp := make([]string, 0)

	for _, msgRef := range msgRefs {
		cacheMessage, err := h.cache.Find(msgRef.ID)
		if err != nil {
			klog.V(4).Infof("Msg ID not found: %s\n", msgRef.ID)
			tmp = append(tmp, interfaces.MessageNotFound)
			continue
		}

		tmp = append(tmp, cacheMessage.Text)
	}

	return tmp
}
