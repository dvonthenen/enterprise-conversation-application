// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"net/http"

	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	halfproxy "github.com/koding/websocketproxy/pkg/half-duplex"
	wsinterfaces "github.com/koding/websocketproxy/pkg/interfaces"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	sse "github.com/r3labs/sse/v2"

	routing "github.com/dvonthenen/enterprise-reference-implementation/pkg/proxy-dataminer/routing"
)

type ProxyOptions struct {
	// housekeeping
	ConversationId    string
	RabbitURI         string
	ProxyBindAddress  string
	NotifyBindAddress string
	RedirectAddress   string
	ProxyPort         int
	NotifyPort        int
	NotifyType        ClientNotifyType

	// SSL Serve
	CrtFile string
	KeyFile string

	// objects
	Neo4jMgr *neo4j.SessionWithContext
	ProxyMgr *wsinterfaces.ManageCallback
}

type Proxy struct {
	options ProxyOptions

	// housekeeping
	// callback *symblinterfaces.InsightCallback
	neo4jMgr   *neo4j.SessionWithContext
	messageMgr *routing.MessageHandler
	proxyMgr   *wsinterfaces.ManageCallback

	// symbl proxy housekeeping
	proxy       *halfproxy.HalfDuplexWebsocketProxy
	serverSymbl *http.Server
	symblChan   chan struct{}

	// rabbit notifications
	rabbitMgr *rabbitinterfaces.Manager

	// server send events
	notifySse    *sse.Server
	notifyServer *http.Server
	notifyChan   chan struct{}
}
