// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"net/http"
	"net/url"

	"github.com/koding/websocketproxy"
	klog "k8s.io/klog/v2"
)

type InstanceOptions struct {
	ConversationId  string
	CrtFile         string
	KeyFile         string
	BindAddress     string
	RedirectAddress string
	Port            int
}

type ServerInstance struct {
	Options InstanceOptions
}

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
		proxy := websocketproxy.NewProxy(websocketproxy.ProxyOptions{
			Url:           u,
			NaturalTunnel: true,
			Viewer:        NewRouter(),
		})

		err = http.ListenAndServeTLS(si.Options.BindAddress, si.Options.CrtFile, si.Options.KeyFile, proxy)
		if err != nil {
			klog.V(1).Infof("New failed. Err: %v\n", err)
		}
	}()

	return err
}
