// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	"encoding/json"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	callback "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
)

func NewEntityHandler(options HandlerOptions) *callback.RabbitMessageHandler {
	var handler callback.RabbitMessageHandler
	handler = EntityHandler{
		session: options.Session,
	}
	return &handler
}

func (eh EntityHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("EntityHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var er interfaces.EntityResponse
	err = json.Unmarshal(byData, &er)
	if err != nil {
		klog.V(1).Infof("EntityResponseMessage json.Unmarshal failed. Err: %v\n", err)
		klog.V(6).Infof("EntityResponseMessage LEAVE\n")
		return err
	}

	// query for any other references to this entity
	// TODO

	return nil
}
