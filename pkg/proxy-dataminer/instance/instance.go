// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"fmt"
	"net/http"
	"net/url"

	rabbit "github.com/dvonthenen/rabbitmq-manager/pkg"
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	symblinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	common "github.com/koding/websocketproxy/pkg/common"
	halfproxy "github.com/koding/websocketproxy/pkg/half-duplex"
	sse "github.com/r3labs/sse/v2"
	klog "k8s.io/klog/v2"

	routing "github.com/dvonthenen/enterprise-reference-implementation/pkg/proxy-dataminer/routing"
)

func New(options ProxyOptions) *Proxy {
	server := &Proxy{
		options:  options,
		neo4jMgr: options.Neo4jMgr,
		proxyMgr: options.ProxyMgr,
	}
	return server
}

func (p *Proxy) GetRedirectAddress() string {
	return p.options.RedirectAddress
}

func (p *Proxy) GetProxyPort() int {
	return p.options.ProxyPort
}

func (p *Proxy) GetNotifyPort() int {
	return p.options.NotifyPort
}

func (p *Proxy) Init() error {
	klog.V(6).Infof("Proxy.Init ENTER\n")

	// init rabbit manager
	rabbitMgr, err := rabbit.New(rabbitinterfaces.ManagerOptions{
		RabbitURI: p.options.RabbitURI,
	})
	if err != nil {
		klog.V(1).Infof("rabbit.New failed. Err: %v\n", err)
		klog.V(6).Infof("Proxy.Init LEAVE\n")
		return err
	}

	/*
		Create Application Channel Subscriber

		This implements the rabbit channel for receiving your High-level Application messages
		sent by the Analyzer component
	*/
	var rabbitHandler rabbitinterfaces.RabbitMessageHandler
	rabbitHandler = p

	_, err = (*rabbitMgr).CreateSubscriber(rabbitinterfaces.SubscriberOptions{
		Name:        p.options.ConversationId,
		AutoDeleted: true,
		IfUnused:    true,
		Handler:     &rabbitHandler,
	})
	if err != nil {
		klog.V(1).Infof("CreateSubscriber %s failed. Err: %v\n", p.options.ConversationId, err)
		klog.V(6).Infof("Proxy.Init LEAVE\n")
		return err
	}

	// create message router
	options := routing.MessageHandlerOptions{
		ConversationId: p.options.ConversationId,
		Neo4jMgr:       p.neo4jMgr,
		RabbitMgr:      rabbitMgr,
	}
	messageMgr, err := routing.NewHandler(options)
	if err != nil {
		klog.V(1).Infof("routing.NewHandler failed. Err: %v\n", err)
		klog.V(6).Infof("Proxy.Init LEAVE\n")
		return err
	}

	// init callback
	err = messageMgr.Init()
	if err != nil {
		klog.V(1).Infof("messageMgr.Init failed. Err: %v\n", err)
		klog.V(6).Infof("Proxy.Init LEAVE\n")
		return err
	}

	// init rabbit
	err = (*rabbitMgr).Init()
	if err != nil {
		klog.V(1).Infof("rabbitMgr.Init failed. Err: %v\n", err)
		klog.V(6).Infof("Proxy.Init LEAVE\n")
		return err
	}

	// housekeeping
	p.rabbitMgr = rabbitMgr
	p.messageMgr = messageMgr

	klog.V(4).Infof("Proxy.Init Succeeded\n")
	klog.V(6).Infof("Proxy.Init LEAVE\n")

	return nil
}

func (p *Proxy) Start() error {
	klog.V(6).Infof("Proxy.Start ENTER\n")

	u, err := url.Parse(DefaultSymblWebSocket)
	if err != nil {
		klog.V(1).Infof("New failed. Err: %v\n", err)
		klog.V(6).Infof("Proxy.Start LEAVE\n")
		return err
	}

	/*
		This is where the hooks are placed for this Proxy service. This implements
		a listener that sits in between the Client Application (ie CPaaS) and the Symbl Platform
	*/
	p.symblChan = make(chan struct{})
	symblServerFunc := func(stopChan chan struct{}) {
		select {
		default:
			var callback symblinterfaces.InsightCallback
			callback = p.messageMgr

			proxy := halfproxy.NewProxy(common.ProxyOptions{
				UniqueID:      p.options.ConversationId,
				Url:           u,
				NaturalTunnel: true,
				Viewer: routing.NewRouter(routing.MessageRouterOptions{
					Callback: &callback,
				}),
				Manager: *p.proxyMgr,
			})
			p.proxy = proxy

			p.serverSymbl = &http.Server{
				Addr:    p.options.ProxyBindAddress,
				Handler: proxy,
			}

			err = p.serverSymbl.ListenAndServeTLS(p.options.CrtFile, p.options.KeyFile)
			if err != nil {
				klog.V(1).Infof("ListenAndServeTLS server stopped. Err: %v\n", err)
			}
		case <-stopChan:
			klog.V(6).Infof("Exiting symblServer Loop\n")
			return
		}
	}
	go symblServerFunc(p.symblChan)

	/*
		This is the listener for the Higher-level Application-specific Messages destined for the
		Client Application
	*/
	if p.options.NotifyType == ClientNotifyTypeServerSendEvent {
		p.notifyChan = make(chan struct{})
		notifyServerFunc := func(stopChan chan struct{}) {
			select {
			default:
				p.notifySse = sse.New()
				p.notifySse.CreateStream("messages")

				notifyPath := fmt.Sprintf("/%s/notifications", p.options.ConversationId)
				klog.V(3).Infof("notifyPath: %s\n", notifyPath)

				mux := http.NewServeMux()
				mux.HandleFunc(notifyPath, func(w http.ResponseWriter, r *http.Request) {
					go func() {
						// Received Browser Disconnection
						<-r.Context().Done()
						klog.V(3).Infof("Received client disconnect notice")
						return
					}()

					p.notifySse.ServeHTTP(w, r)
				})

				p.notifyServer = &http.Server{
					Addr:    p.options.NotifyBindAddress,
					Handler: mux,
				}

				err = p.notifyServer.ListenAndServeTLS(p.options.CrtFile, p.options.KeyFile)
				if err != nil {
					klog.V(1).Infof("ListenAndServeTLS server stopped. Err: %v\n", err)
				}
			case <-stopChan:
				klog.V(6).Infof("Exiting notifyServer Loop\n")
				return
			}
		}
		go notifyServerFunc(p.notifyChan)
	}

	klog.V(4).Infof("Proxy.Start Succeeded\n")
	klog.V(6).Infof("Proxy.Start LEAVE\n")

	return nil
}

func (p *Proxy) ProcessMessage(byData []byte) error {
	/*
		This forwards the Higher-level Application Messages from the Analyzer component to the
		Client Application (Web UI, CPaaS, etc). These are NOT Symbl conversation messages, but
		rather your custom application messages.
	*/
	klog.V(5).Infof(" [x] %s\n", string(byData))

	switch p.options.NotifyType {
	case ClientNotifyTypeWebSocket:
		if p.options.NotifyType != ClientNotifyTypeWebSocket || p.proxy == nil {
			return ErrInvalidNotifyConfig
		}
		err := p.proxy.SendMessage(byData)
		if err != nil {
			klog.V(1).Infof("SendMessage failed. Err: %v\n", err)
		}
		return err
	case ClientNotifyTypeServerSendEvent:
		if p.options.NotifyType != ClientNotifyTypeServerSendEvent || p.notifySse == nil {
			return ErrInvalidNotifyConfig
		}
		p.notifySse.Publish("messages", &sse.Event{
			Data: []byte(byData),
		})
		return nil
	default:
		klog.V(1).Infof("Unknown Message Type: %d\n", p.options.NotifyType)
		return ErrUnknownNotifyType
	}

	return nil
}

func (p *Proxy) IsConnected() bool {
	if p.serverSymbl == nil || p.serverSymbl.Handler == nil {
		return false
	}
	isConnected := p.proxy.IsConnected()
	klog.V(3).Infof("Proxy.IsConnected: %t\n", isConnected)
	return isConnected
}

func (p *Proxy) Stop() error {
	klog.V(6).Infof("Proxy.Stop ENTER\n")

	/*
		fire off a teardown message explicitly in case connection terminated abnormally
		there is logic in TeardownConversation() to only fire once
	*/
	if p.messageMgr != nil {
		teardownMsg := sdkinterfaces.TeardownMessage{}
		teardownMsg.Type = MessageTypeMessage
		teardownMsg.Message.Type = MessageTypeTeardownConversation
		teardownMsg.Message.Data.ConversationID = p.options.ConversationId

		err := p.messageMgr.TeardownConversation(&teardownMsg)
		if err != nil {
			klog.V(1).Infof("TeardownConversation failed. Err: %v\n", err)
		}
	}

	// wait for threads
	if p.symblChan != nil {
		close(p.symblChan)
		<-p.symblChan
	}
	if p.notifyChan != nil {
		close(p.notifyChan)
		<-p.notifyChan
	}

	// delete the handler
	if p.messageMgr != nil {
		err := p.messageMgr.Teardown()
		if err != nil {
			klog.V(1).Infof("messageMgr.Teardown() failed. Err: %v\n", err)
		}
	}

	/*
		Delete Application Channel Subscriber

		This delete the rabbit channel for receiving your High-level Application messages
		sent by the Analyzer component
	*/
	if p.rabbitMgr != nil {
		err := (*p.rabbitMgr).Teardown()
		if err != nil {
			klog.V(1).Infof("rabbitMgr.Teardown() failed. Err: %v\n", err)
		}
	}
	p.rabbitMgr = nil

	// close HTTP server
	if p.notifyServer != nil {
		err := p.notifyServer.Close()
		if err != nil {
			klog.V(1).Infof("notifyServer.Close() failed. Err: %v\n", err)
		}
	}
	p.notifyServer = nil

	if p.serverSymbl != nil {
		err := p.serverSymbl.Close()
		if err != nil {
			klog.V(1).Infof("serverSymbl.Close() failed. Err: %v\n", err)
		}
	}
	p.serverSymbl = nil

	klog.V(4).Infof("Proxy.Stop Succeeded\n")
	klog.V(6).Infof("Proxy.Stop LEAVE\n")

	return nil
}
