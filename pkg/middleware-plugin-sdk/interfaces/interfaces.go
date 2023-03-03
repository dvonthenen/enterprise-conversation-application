// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import (
	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

/*
	This implementer of this interface handle much of the routing to your middleware-analyzer component
	As Symbl insights are receive from the platform, these callbacks will be invoked
*/
type InsightCallback interface {
	sdkinterfaces.InsightCallback

	/*
		The Client Publisher interface will be set before messages trigger functions in the
		sdkinterfaces.InsightCallback
	*/
	SetClientPublisher(mp *MessagePublisher)
}

/*
	This provides an interface for the implementing struct to sent messages to the client.
*/
type MessagePublisher interface {
	PublishMessage(name string, data []byte) error
}

/*
	Interface to the InsightManager which receives Rabbit messages and then calls the
	appropriate callback function
*/
type InsightManager interface {
	Init() error
	Teardown() error
}
