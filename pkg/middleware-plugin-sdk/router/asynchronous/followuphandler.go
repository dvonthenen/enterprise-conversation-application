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

func NewFollowUpHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = FollowUpHandler{
		manager:  options.Manager,
		callback: options.Callback,
	}
	return &handler
}

func (fuh FollowUpHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("FollowUpHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var fur shared.FollowUpResult
	err = json.Unmarshal(byData, &fur)
	if err != nil {
		klog.V(1).Infof("[FollowUpHandler] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	// invoke callback
	err = (*fuh.callback).FollowUpResult(&fur)
	if err == nil {
		klog.V(5).Infof("[FollowUpHandler] Callback succeeded\n")
	} else {
		klog.V(1).Infof("[FollowUpHandler] Callback failed. Err: %v\n", err)
		return err
	}

	return nil
}
