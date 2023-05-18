// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package middleware

import (
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"

	interfaces "github.com/dvonthenen/enterprise-conversation-application/pkg/middleware-plugin-sdk/interfaces"
)

/*
	RealtimeAnalyzer struct
*/
type RealtimeAnalyzerOption struct {
	RabbitURI string
	Callback  *interfaces.InsightCallback
}

type RealtimeAnalyzer struct {
	// rabbit
	rabbitManager *rabbitinterfaces.Manager

	// callback
	callback *interfaces.InsightCallback
}

/*
	AsynchronousAnalyzer struct
*/
type AsynchronousAnalyzerOption struct {
	RabbitURI string
	Callback  *interfaces.AsynchronousCallback
}

type AsynchronousAnalyzer struct {
	// rabbit
	rabbitManager *rabbitinterfaces.Manager

	// callback
	callback *interfaces.AsynchronousCallback
}
