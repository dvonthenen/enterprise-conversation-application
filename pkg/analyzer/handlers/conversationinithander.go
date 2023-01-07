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

func NewConversationInitHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = ConversationInitHandler{
		session:     options.Session,
		symblClient: options.SymblClient,
		manager:     options.Manager,
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
	klog.V(2).Infof("ConversationInitHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var im interfaces.InitializationMessage
	err = json.Unmarshal(byData, &im)
	if err != nil {
		klog.V(1).Infof("[ConversationInit] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	/*
		Create Application Channel Publisher

		This implements the rabbit channel for sending your High-level Application messages
		sent by this component based on the conversationId
	*/
	_, err = (*ch.manager).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        im.Message.Data.ConversationID,
		AutoDeleted: true,
		IfUnused:    true,
		// NoWait:      true,
	})
	if err == nil {
		klog.V(4).Infof("CreatePublisher succeeded\n")
	} else {
		klog.V(1).Infof("[ConversationInit] CreatePublisher failed. Err: %v\n", err)
		return err
	}

	// TODO: template for add your businesss logic

	return nil
}
