// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package router

import (
	"encoding/json"

	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	shared "github.com/dvonthenen/enterprise-conversation-application/pkg/shared"
)

func NewActionItemHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = ActionItemHandler{
		manager:  options.Manager,
		callback: options.Callback,
	}
	return &handler
}

func (aih ActionItemHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("ActionItemHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var air shared.ActionItemResult
	err = json.Unmarshal(byData, &air)
	if err != nil {
		klog.V(1).Infof("[ActionItemHandler] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	// invoke callback
	err = (*aih.callback).ActionItemResult(&air)
	if err == nil {
		klog.V(5).Infof("[ActionItemHandler] Callback succeeded\n")
	} else {
		klog.V(1).Infof("[ActionItemHandler] Callback failed. Err: %v\n", err)
		return err
	}

	return nil
}
