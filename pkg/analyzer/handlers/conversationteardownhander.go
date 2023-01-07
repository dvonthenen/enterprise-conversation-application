// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	"encoding/json"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
)

func NewConversationTeardownHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = ConversationTeardownHandler{
		session:     options.Session,
		symblClient: options.SymblClient,
		manager:     options.Manager,
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
	klog.V(2).Infof("ConversationTeardownHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var tm interfaces.TeardownMessage
	err = json.Unmarshal(byData, &tm)
	if err != nil {
		klog.V(1).Infof("[ConversationTeardown] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	/*
		Delete Application Channel Publisher

		This cleans up the rabbit channel for sending your High-level Application messages
		sent by this component based on the conversationId
	*/
	err = (*ch.manager).DeletePublisher(tm.Message.Data.ConversationID)
	if err == nil {
		klog.V(4).Infof("DeletePublisher succeeded\n")
	} else {
		klog.V(1).Infof("[ConversationTeardown] DeletePublisher failed. Err: %v\n", err)
		return err
	}

	// TODO: template for add your businesss logic

	return nil
}
