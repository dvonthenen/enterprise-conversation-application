// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	klog "k8s.io/klog/v2"

	callback "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

func NewMessageHandler(options HandlerOptions) *callback.RabbitMessageHandler {
	var handler callback.RabbitMessageHandler
	handler = MessageHandler{
		session: options.Session,
	}
	return &handler
}

func (ch MessageHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("MessageHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// TODO

	return nil
}
