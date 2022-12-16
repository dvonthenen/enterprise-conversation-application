// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	callback "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
)

func NewEntityHandler(options HandlerOptions) *callback.RabbitMessageHandler {
	var handler callback.RabbitMessageHandler
	handler = EntityHandler{
		session:      options.Session,
		symblClient:  options.SymblClient,
		pushCallback: options.PushCallback,
	}
	return &handler
}

func (eh EntityHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("[EntityHandler] prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("EntityHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// TODO: template for add your businesss logic

	return nil
}
