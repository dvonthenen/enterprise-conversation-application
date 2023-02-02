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

func NewEntityHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = EntityHandler{
		manager:  options.Manager,
		callback: options.Callback,
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

	// reform struct
	var er interfaces.EntityResponse
	err = json.Unmarshal(byData, &er)
	if err != nil {
		klog.V(1).Infof("[EntityHandler] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	// invoke callback
	err = (*eh.callback).EntityResponseMessage(er.EntityResponse)
	if err == nil {
		klog.V(5).Infof("[EntityHandler] Callback succeeded\n")
	} else {
		klog.V(1).Infof("[EntityHandler] Callback failed. Err: %v\n", err)
		return err
	}

	return nil
}
