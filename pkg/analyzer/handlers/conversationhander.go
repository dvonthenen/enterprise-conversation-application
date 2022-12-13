// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

func NewConversationHandler(options HandlerOptions) *interfaces.RabbitMessageHandler {
	var handler interfaces.RabbitMessageHandler
	handler = ConversationHandler{
		session: options.Session,
	}
	return &handler
}

func (ch ConversationHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("ConversationHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// TODO

	return nil
}
