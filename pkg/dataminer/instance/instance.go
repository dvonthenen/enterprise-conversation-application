// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"net/http"
	"net/url"

	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	common "github.com/koding/websocketproxy/pkg/common"
	halfproxy "github.com/koding/websocketproxy/pkg/half-duplex"
	klog "k8s.io/klog/v2"

	routing "github.com/dvonthenen/enterprise-reference-implementation/pkg/dataminer/routing"
)

func New(options InstanceOptions) *ServerInstance {
	server := &ServerInstance{
		Options: options,
	}
	return server
}

func (si *ServerInstance) Start() error {
	u, err := url.Parse(DefaultSymblWebSocket)
	if err != nil {
		klog.V(1).Infof("New failed. Err: %v\n", err)
		return err
	}

	go func() {
		proxy := halfproxy.NewProxy(common.ProxyOptions{
			UniqueID:      si.Options.ConversationId,
			Url:           u,
			NaturalTunnel: true,
			Viewer:        routing.NewRouter(si.Options.Callback),
			Manager:       *si.Options.Manager,
		})
		si.proxy = proxy

		si.server = &http.Server{
			Addr:    si.Options.BindAddress,
			Handler: proxy,
		}

		err = si.server.ListenAndServeTLS(si.Options.CrtFile, si.Options.KeyFile)
		if err != nil {
			klog.V(1).Infof("ListenAndServeTLS server stopped. Err: %v\n", err)
		}
	}()

	return err
}

func (si *ServerInstance) IsConnected() bool {
	if si.server == nil || si.server.Handler == nil {
		return false
	}
	return si.proxy.IsConnected()
}

func (si *ServerInstance) Stop() error {
	// fire teardown message explicitly in case connection terminated abnormally
	// there is logic in TeardownConversation() to only fire once
	callback := si.Options.Callback
	if callback != nil {
		teardownMsg := sdkinterfaces.TeardownMessage{}
		teardownMsg.Type = MessageTypeMessage
		teardownMsg.Message.Type = MessageTypeTeardownConversation
		teardownMsg.Message.Data.ConversationID = si.Options.ConversationId

		err := (*callback).TeardownConversation(&teardownMsg)
		if err != nil {
			klog.V(1).Infof("TeardownConversation failed. Err: %v\n", err)
		}
	}

	// close HTTP server
	err := si.server.Close()
	if err != nil {
		klog.V(1).Infof("server.Close() failed. Err: %v\n", err)
	}

	return nil
}
