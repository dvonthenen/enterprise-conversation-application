// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package instance

import (
	"net/http"

	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	halfproxy "github.com/dvonthenen/websocketproxy/pkg/half-duplex"
	wsinterfaces "github.com/dvonthenen/websocketproxy/pkg/interfaces"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	sse "github.com/r3labs/sse/v2"

	routing "github.com/dvonthenen/enterprise-conversation-application/pkg/proxy-dataminer/routing"
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

	// functionality
	TranscriptionEnabled bool
	MessagingEnabled     bool

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
