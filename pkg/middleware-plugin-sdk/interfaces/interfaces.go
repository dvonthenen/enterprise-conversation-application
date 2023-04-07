// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import (
	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

/*
	This implementer of this interface handles much of the routing to your middleware-analyzer component
	As Symbl insights are received from the platform, these callbacks will be invoked
*/
type InsightCallback interface {
	/*
		This ensures there is a 1-to-1 mapping between insights the Symbl Platform provides and
		which events are possible to be notified to
	*/
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
	appropriate callback function in the InsightCallback interface above
*/
type InsightManager interface {
	Init() error
	Teardown() error
}
