// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package router

import (
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/interfaces"
)

/*
	Subscriber handlers
*/
type HandlerOptions struct {
	Manager  *rabbitinterfaces.Manager
	Callback *interfaces.AsynchronousCallback
}

type ConversationInitHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback

	appMessage *rabbitinterfaces.Publisher
}

type MessageHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}

type QuestionHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}

type FollowUpHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}

type ActionItemHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}

type TopicHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}

type SummaryHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}

type TrackerHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}

type EntityHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}

type SummaryUiHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}

type ConversationTeardownHandler struct {
	manager  *rabbitinterfaces.Manager
	callback *interfaces.AsynchronousCallback
}
