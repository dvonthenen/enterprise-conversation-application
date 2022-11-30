// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"net/http"
	"net/url"

	common "github.com/koding/websocketproxy/pkg/common"
	halfproxy "github.com/koding/websocketproxy/pkg/half-duplex"
	klog "k8s.io/klog/v2"

	routing "github.com/dvonthenen/enterprise-reference-implementation/pkg/routing"
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
			Url:           u,
			NaturalTunnel: true,
			Viewer:        routing.NewRouter(&si.Options.Callback),
		})

		si.server = &http.Server{
			Addr:    si.Options.BindAddress,
			Handler: proxy,
		}

		err = si.server.ListenAndServeTLS(si.Options.CrtFile, si.Options.KeyFile)
		if err != nil {
			klog.V(6).Infof("ListenAndServeTLS failed. Err: %v\n", err)
		}
	}()

	return err
}

func (si *ServerInstance) Stop() error {
	err := si.server.Close()
	if err != nil {
		klog.V(1).Infof("New failed. Err: %v\n", err)
		return err
	}

	return nil
}
