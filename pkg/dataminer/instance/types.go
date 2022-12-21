// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"net/http"

	symblinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	halfproxy "github.com/koding/websocketproxy/pkg/half-duplex"
	wsinterfaces "github.com/koding/websocketproxy/pkg/interfaces"
	sse "github.com/r3labs/sse/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

type InstanceOptions struct {
	ConversationId    string
	CrtFile           string
	KeyFile           string
	ProxyBindAddress  string
	NotifyBindAddress string
	RedirectAddress   string
	ProxyPort         int
	NotifyPort        int

	RabbitConn *amqp.Connection
	Callback   *symblinterfaces.InsightCallback
	Manager    *wsinterfaces.ManageCallback
}

type ServerInstance struct {
	Options InstanceOptions

	// housekeeping
	callback    *symblinterfaces.InsightCallback
	manager     *wsinterfaces.ManageCallback
	channel     *amqp.Channel
	queue       *amqp.Queue
	messageChan chan struct{}

	// symbl proxy housekeeping
	proxy       *halfproxy.HalfDuplexWebsocketProxy
	serverSymbl *http.Server
	symblChan   chan struct{}

	// rabbit notifications
	rabbitConn *amqp.Connection

	// server send events
	notification *sse.Server
	notifyServer *http.Server
	notifyChan   chan struct{}
}
