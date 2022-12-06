// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"net/http"

	interfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

type InstanceOptions struct {
	ConversationId  string
	CrtFile         string
	KeyFile         string
	BindAddress     string
	RedirectAddress string
	Port            int
	Callback        interfaces.InsightCallback
}

type ServerInstance struct {
	Options InstanceOptions

	server *http.Server
}
