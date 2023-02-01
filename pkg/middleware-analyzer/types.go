// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package middleware

import (
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-analyzer/interfaces"
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
