// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"fmt"
	"net/http"
	"net/url"

	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	common "github.com/koding/websocketproxy/pkg/common"
	halfproxy "github.com/koding/websocketproxy/pkg/half-duplex"
	sse "github.com/r3labs/sse/v2"
	klog "k8s.io/klog/v2"

	routing "github.com/dvonthenen/enterprise-reference-implementation/pkg/dataminer/routing"
)

func New(options InstanceOptions) *ServerInstance {
	server := &ServerInstance{
		Options:     options,
		rabbitConn:  options.RabbitConn,
		callback:    options.Callback,
		manager:     options.Manager,
		messageChan: make(chan struct{}),
		symblChan:   make(chan struct{}),
		notifyChan:  make(chan struct{}),
	}
	return server
}

func (si *ServerInstance) Init() error {
	klog.V(6).Infof("ServerInstance.Init ENTER\n")

	ch, err := si.rabbitConn.Channel()
	if err != nil {
		klog.V(1).Infof("Channel failed. Err: %v\n", err)
		klog.V(6).Infof("ServerInstance.Init LEAVE\n")
		return err
	}

	klog.V(3).Infof("ExchangeDeclare: %v\n", si.Options.ConversationId)
	err = ch.ExchangeDeclare(
		si.Options.ConversationId, // name
		"fanout",                  // type
		true,                      // durable
		true,                      // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		klog.V(1).Infof("ExchangeDeclare(%s) failed. Err: %v\n", si.Options.ConversationId, err)
		klog.V(6).Infof("ServerInstance.Init LEAVE\n")
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		klog.V(1).Infof("QueueDeclare() failed. Err: %v\n", err)
		klog.V(6).Infof("ServerInstance.Init LEAVE\n")
		return err
	}

	klog.V(3).Infof("QueueBind: %v\n", si.Options.ConversationId)
	err = ch.QueueBind(
		q.Name,                    // queue name
		"",                        // routing key
		si.Options.ConversationId, // exchange
		false,
		nil)
	if err != nil {
		klog.V(1).Infof("QueueBind() failed. Err: %v\n", err)
		klog.V(6).Infof("ServerInstance.Init LEAVE\n")
		return err
	}

	// housekeeping
	si.channel = ch
	si.queue = &q

	klog.V(4).Infof("ServerInstance.Init Succeeded\n")
	klog.V(6).Infof("ServerInstance.Init LEAVE\n")

	return nil
}

func (si *ServerInstance) Start() error {
	klog.V(6).Infof("ServerInstance.Start ENTER\n")

	u, err := url.Parse(DefaultSymblWebSocket)
	if err != nil {
		klog.V(1).Infof("New failed. Err: %v\n", err)
		klog.V(6).Infof("ServerInstance.Start LEAVE\n")
		return err
	}

	// symbl proxy
	symblServerFunc := func(stopChan chan struct{}) {
		select {
		default:
			proxy := halfproxy.NewProxy(common.ProxyOptions{
				UniqueID:      si.Options.ConversationId,
				Url:           u,
				NaturalTunnel: true,
				Viewer:        routing.NewRouter(si.callback),
				Manager:       *si.manager,
			})
			si.proxy = proxy

			si.serverSymbl = &http.Server{
				Addr:    si.Options.ProxyBindAddress,
				Handler: proxy,
			}

			err = si.serverSymbl.ListenAndServeTLS(si.Options.CrtFile, si.Options.KeyFile)
			if err != nil {
				klog.V(1).Infof("ListenAndServeTLS server stopped. Err: %v\n", err)
			}
		case <-stopChan:
			klog.V(6).Infof("Exiting symblServer Loop\n")
			return
		}
	}
	go symblServerFunc(si.symblChan)

	// async notifications
	notifyServerFunc := func(stopChan chan struct{}) {
		select {
		default:
			si.notification = sse.New()
			si.notification.CreateStream("messages")

			notifyPath := fmt.Sprintf("/%s/notifications", si.Options.ConversationId)
			klog.V(3).Infof("notifyPath: %s\n", notifyPath)

			mux := http.NewServeMux()
			mux.HandleFunc(notifyPath, func(w http.ResponseWriter, r *http.Request) {
				go func() {
					// Received Browser Disconnection
					<-r.Context().Done()
					println("The client is disconnected here")
					return
				}()

				si.notification.ServeHTTP(w, r)
			})

			si.notifyServer = &http.Server{
				Addr:    si.Options.NotifyBindAddress,
				Handler: mux,
			}

			err = si.notifyServer.ListenAndServeTLS(si.Options.CrtFile, si.Options.KeyFile)
			if err != nil {
				klog.V(1).Infof("ListenAndServeTLS server stopped. Err: %v\n", err)
			}
		case <-stopChan:
			klog.V(6).Infof("Exiting notifyServer Loop\n")
			return
		}
	}
	go notifyServerFunc(si.notifyChan)

	// async notifications
	msgs, err := si.channel.Consume(
		si.queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		klog.V(1).Infof("Consume failed. Err: %v\n", err)
		return err
	}

	klog.V(5).Infof("Server.Start Running message loop...\n")
	processNotifyMsg := func(stopChan chan struct{}) {
		for {
			select {
			default:
				for d := range msgs {
					klog.V(5).Infof(" [x] %s\n", d.Body)

					si.notification.Publish("messages", &sse.Event{
						Data: []byte(d.Body),
					})
				}
			case <-stopChan:
				klog.V(6).Infof("Exiting Messaging Loop\n")
				return
			}
		}
	}
	go processNotifyMsg(si.messageChan)

	klog.V(4).Infof("ServerInstance.Start Succeeded\n")
	klog.V(6).Infof("ServerInstance.Start LEAVE\n")

	return nil
}

func (si *ServerInstance) IsConnected() bool {
	if si.serverSymbl == nil || si.serverSymbl.Handler == nil {
		return false
	}
	isConnected := si.proxy.IsConnected()
	klog.V(3).Infof("ServerInstance.IsConnected: %t\n", isConnected)
	return isConnected
}

func (si *ServerInstance) Stop() error {
	klog.V(6).Infof("ServerInstance.Stop ENTER\n")

	// fire teardown message explicitly in case connection terminated abnormally
	// there is logic in TeardownConversation() to only fire once
	callback := si.callback
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

	// wait for threads
	close(si.messageChan)
	<-si.messageChan
	close(si.symblChan)
	<-si.symblChan
	close(si.notifyChan)
	<-si.notifyChan

	// rabbit
	if si.channel != nil {
		si.channel.Close()
		si.channel = nil
	}

	// close HTTP server
	err := si.notifyServer.Close()
	if err != nil {
		klog.V(1).Infof("notifyServer.Close() failed. Err: %v\n", err)
	}
	err = si.serverSymbl.Close()
	if err != nil {
		klog.V(1).Infof("serverSymbl.Close() failed. Err: %v\n", err)
	}

	klog.V(4).Infof("ServerInstance.Stop Succeeded\n")
	klog.V(6).Infof("ServerInstance.Stop LEAVE\n")

	return nil
}
