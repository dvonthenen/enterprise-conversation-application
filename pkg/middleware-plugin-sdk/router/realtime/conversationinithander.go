// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package router

import (
	"encoding/json"

	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	shared "github.com/dvonthenen/enterprise-reference-implementation/pkg/shared"
)

func NewConversationInitHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = ConversationInitHandler{
		manager:  options.Manager,
		callback: options.Callback,
	}
	return &handler
}

func (ch ConversationInitHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("InitializationHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var ir shared.InitializationResponse
	err = json.Unmarshal(byData, &ir)
	if err != nil {
		klog.V(1).Infof("[InitializationHandler] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	/*
		Create Application Channel Publisher

		This implements the rabbit channel for sending your High-level Application messages
		sent by this component based on the conversationId
	*/
	_, err = (*ch.manager).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        ir.InitializationMessage.Message.Data.ConversationID,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err == nil {
		klog.V(4).Infof("[InitializationHandler] CreatePublisher succeeded\n")
	} else {
		klog.V(1).Infof("[InitializationHandler] CreatePublisher failed. Err: %v\n", err)
		return err
	}

	// invoke callback
	err = (*ch.callback).InitializedConversation(&ir)
	if err == nil {
		klog.V(5).Infof("[InitializationHandler] Callback succeeded\n")
	} else {
		klog.V(1).Infof("[InitializationHandler] Callback failed. Err: %v\n", err)
		return err
	}

	return nil
}
