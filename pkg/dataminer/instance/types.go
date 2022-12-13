// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"net/http"

	symblinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	halfproxy "github.com/koding/websocketproxy/pkg/half-duplex"
	wsinterfaces "github.com/koding/websocketproxy/pkg/interfaces"
)

type InstanceOptions struct {
	ConversationId  string
	CrtFile         string
	KeyFile         string
	BindAddress     string
	RedirectAddress string
	Port            int
	Callback        *symblinterfaces.InsightCallback
	Manager         *wsinterfaces.ManageCallback
}

type ServerInstance struct {
	Options InstanceOptions

	proxy  *halfproxy.HalfDuplexWebsocketProxy
	server *http.Server
}
