// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import (
	shared "github.com/dvonthenen/enterprise-conversation-application/pkg/shared"
)

/*
	Common
*/
/*
	This provides an interface for the implementing struct to sent messages to the client.
*/
type MessagePublisher interface {
	PublishMessage(name string, data []byte) error
}

/*
	Real-time Streaming
*/
/*
	This implementer of this interface handles much of the routing to your middleware-analyzer component
	As Symbl insights are received from the platform, these callbacks will be invoked
*/
type InsightCallback interface {
	/*
		This ensures there is a 1-to-1 mapping between insights the Symbl Platform provides and
		which events are possible to be notified to
	*/
	// streaminginterfaces.InsightCallback

	InitializedConversation(im *shared.InitializationResponse) error
	RecognitionResultMessage(rr *shared.RecognitionResponse) error
	MessageResponseMessage(mr *shared.MessageResponse) error
	InsightResponseMessage(ir *shared.InsightResponse) error
	TopicResponseMessage(tr *shared.TopicResponse) error
	TrackerResponseMessage(tr *shared.TrackerResponse) error
	EntityResponseMessage(er *shared.EntityResponse) error
	TeardownConversation(tm *shared.TeardownResponse) error
	UserDefinedMessage(data []byte) error
	UnhandledMessage(byMsg []byte) error

	/*
		The Client Publisher interface will be set before messages trigger functions in the
		streaminginterfaces.InsightCallback
	*/
	SetClientPublisher(mp *MessagePublisher)
}

/*
	Interface to the RealtimeManager which receives Rabbit messages and then calls the
	appropriate callback function in the InsightCallback interface above
*/
type RealtimeManager interface {
	Init() error
	Teardown() error
}

/*
	Asynchronous
*/
/*
	This implementer of this interface handles much of the routing to your middleware-analyzer component
	As Symbl insights are received from the platform, these callbacks will be invoked when using Async APIs
*/
type AsynchronousCallback interface {
	/*
		This ensures there is a 1-to-1 mapping between insights the Symbl Platform provides and
		which events are possible to be notified to
	*/
	// asyncinterfaces.InsightCallback

	InitializedConversation(ci *shared.InitializationResult) error
	MessageResult(mr *shared.MessageResult) error
	QuestionResult(qr *shared.QuestionResult) error
	FollowUpResult(fr *shared.FollowUpResult) error
	ActionItemResult(air *shared.ActionItemResult) error
	TopicResult(tr *shared.TopicResult) error
	TrackerResult(tr *shared.TrackerResult) error
	EntityResult(er *shared.EntityResult) error
	TeardownConversation(ct *shared.TeardownResult) error

	/*
		The Client Publisher interface will be set before messages trigger functions in the
		sdkinterfaces.InsightCallback
	*/
	SetClientPublisher(mp *MessagePublisher)
}

/*
	Interface to the AsynchronousManager which receives Rabbit messages and then calls the
	appropriate callback function in the AsynchronousCallback interface above
*/
type AsynchronousManager interface {
	Init() error
	Teardown() error
}
