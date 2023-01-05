// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package routing

import (
	streaming "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1"
	klog "k8s.io/klog/v2"
)

func NewRouter(options MessageRouterOptions) *MessageRouter {
	mr := &MessageRouter{
		callback: options.Callback,
	}
	return mr
}

func (mr *MessageRouter) HandleMessage(byMsg []byte) error {
	klog.V(6).Infof("MessageRouter.HandleMessage ENTER\n")

	router := streaming.New(*mr.callback)
	err := router.Message(byMsg)
	if err != nil {
		klog.V(1).Infof("HandleMessage Failed. Err: %v\n", err)
		klog.V(6).Infof("MessageRouter.HandleMessage LEAVE\n")

		return err
	}

	klog.V(6).Infof("HandleMessage Succeeded\n")
	klog.V(6).Infof("MessageRouter.HandleMessage LEAVE\n")

	return nil
}
