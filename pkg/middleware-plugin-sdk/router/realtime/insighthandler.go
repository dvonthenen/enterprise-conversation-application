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

func NewInsightHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = InsightHandler{
		manager:  options.Manager,
		callback: options.Callback,
	}
	return &handler
}

func (ih InsightHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("InsightHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var ir shared.InsightResponse
	err = json.Unmarshal(byData, &ir)
	if err != nil {
		klog.V(1).Infof("[InsightHandler] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	// invoke callback
	err = (*ih.callback).InsightResponseMessage(&ir)
	if err == nil {
		klog.V(5).Infof("[InsightHandler] Callback succeeded\n")
	} else {
		klog.V(1).Infof("[InsightHandler] Callback failed. Err: %v\n", err)
		return err
	}

	return nil
}
