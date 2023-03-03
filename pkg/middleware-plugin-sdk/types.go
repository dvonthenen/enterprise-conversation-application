// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package middleware

import (
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/interfaces"
)

/*
	MiddlewareAnalyzer struct
*/
type MiddlewareAnalyzerOption struct {
	RabbitURI string
	Callback  *interfaces.InsightCallback
}

type MiddlewareAnalyzer struct {
	// rabbit
	rabbitManager *rabbitinterfaces.Manager

	// callback
	callback *interfaces.InsightCallback
}
