// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package router

import (
	"encoding/json"

	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
)

func NewConversationTeardownHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = ConversationTeardownHandler{
		manager:  options.Manager,
		callback: options.Callback,
	}
	return &handler
}

func (ch ConversationTeardownHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("TeardownHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var tm interfaces.TeardownMessage
	err = json.Unmarshal(byData, &tm)
	if err != nil {
		klog.V(1).Infof("[TeardownHandler] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	/*
		Delete Application Channel Publisher

		This cleans up the rabbit channel for sending your High-level Application messages
		sent by this component based on the conversationId
	*/
	err = (*ch.manager).DeletePublisher(tm.Message.Data.ConversationID)
	if err == nil {
		klog.V(4).Infof("[TeardownHandler] DeletePublisher succeeded\n")
	} else {
		klog.V(1).Infof("[TeardownHandler] DeletePublisher failed. Err: %v\n", err)
		return err
	}

	// invoke callback
	err = (*ch.callback).TeardownConversation(&tm.TeardownMessage)
	if err == nil {
		klog.V(5).Infof("[TeardownHandler] Callback succeeded\n")
	} else {
		klog.V(1).Infof("[TeardownHandler] Callback failed. Err: %v\n", err)
		return err
	}

	return nil
}
