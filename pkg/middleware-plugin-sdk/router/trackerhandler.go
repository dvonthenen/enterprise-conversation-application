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

func NewTrackerHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = TrackerHandler{
		manager:  options.Manager,
		callback: options.Callback,
	}
	return &handler
}

func (th TrackerHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("TrackerHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var tr shared.TrackerResponse
	err = json.Unmarshal(byData, &tr)
	if err != nil {
		klog.V(1).Infof("[TrackerHandler] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	// invoke callback
	err = (*th.callback).TrackerResponseMessage(tr.TrackerResponse)
	if err == nil {
		klog.V(5).Infof("[TrackerHandler] Callback succeeded\n")
	} else {
		klog.V(1).Infof("[TrackerHandler] Callback failed. Err: %v\n", err)
		return err
	}

	return nil
}
