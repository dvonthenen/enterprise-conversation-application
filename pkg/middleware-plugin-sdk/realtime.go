// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package middleware

import (
	rabbit "github.com/dvonthenen/rabbitmq-manager/pkg"
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	klog "k8s.io/klog/v2"

	middlewareinterfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/interfaces"
	router "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/router/realtime"
	shared "github.com/dvonthenen/enterprise-reference-implementation/pkg/shared"
)

func NewRealtimeAnalyzer(options RealtimeAnalyzerOption) (*RealtimeAnalyzer, error) {
	// setup rabbit manager
	rabbitMgr, err := rabbit.New(rabbitinterfaces.ManagerOptions{
		RabbitURI: options.RabbitURI,
	})
	if err != nil {
		klog.V(1).Infof("rabbit.New failed. Err: %v\n", err)
		return nil, err
	}

	// create middleware
	mgr := &RealtimeAnalyzer{
		rabbitManager: rabbitMgr,
		callback:      options.Callback,
	}

	// set publisher
	var messagePublisher middlewareinterfaces.MessagePublisher
	messagePublisher = mgr
	(*mgr.callback).SetClientPublisher(&messagePublisher)

	return mgr, nil
}

/*
	This initializes all of the subscribers to the Symbl Proxy/Dataminer component

	Each rabbit subscriber listens for a specific Symbl derived/discovered conversation insight and
	is then notified through a callback handler with the original Symbl RealTime API message struct
*/
func (ma *RealtimeAnalyzer) Init() error {
	klog.V(6).Infof("NotificationManager.Init ENTER\n")

	type InitFunc func(router.HandlerOptions) *rabbitinterfaces.RabbitMessageHandler
	type MyHandler struct {
		Name string
		Func InitFunc
	}

	// init rabbit clients
	myHandlers := make([]*MyHandler, 0)
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitRealTimeConversationInit,
		Func: router.NewConversationInitHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitRealTimeEntity,
		Func: router.NewEntityHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitRealTimeInsight,
		Func: router.NewInsightHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitRealTimeMessage,
		Func: router.NewMessageHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitRealTimeTopic,
		Func: router.NewTopicHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitRealTimeTracker,
		Func: router.NewTrackerHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitRealTimeConversationTeardown,
		Func: router.NewConversationTeardownHandler,
	})

	for _, myHandler := range myHandlers {
		// create subscriber
		handler := myHandler.Func(router.HandlerOptions{
			Manager:  ma.rabbitManager,
			Callback: ma.callback,
		})

		_, err := (*ma.rabbitManager).CreateSubscriber(rabbitinterfaces.SubscriberOptions{
			Name:        myHandler.Name,
			Type:        rabbitinterfaces.ExchangeTypeFanout,
			AutoDeleted: true,
			IfUnused:    true,
			Handler:     handler,
		})
		if err != nil {
			klog.V(1).Infof("CreateSubscription failed. Err: %v\n", err)
		}
	}

	// init the system
	err := (*ma.rabbitManager).Init()
	if err != nil {
		klog.V(1).Infof("rabbitManager.Init failed. Err: %v\n", err)
		klog.V(6).Infof("NotificationManager.Init LEAVE\n")
		return err
	}

	klog.V(4).Infof("Init Succeeded\n")
	klog.V(6).Infof("NotificationManager.Init LEAVE\n")

	return nil
}

func (ma *RealtimeAnalyzer) PublishMessage(name string, data []byte) error {
	klog.V(6).Infof("NotificationManager.PublishMessage ENTER\n")

	publisher, err := (*ma.rabbitManager).GetPublisherByName(name)
	if err != nil {
		klog.V(1).Infof("GetPublisherByName failed. Err: %v\n", err)
		klog.V(6).Infof("NotificationManager.PublishMessage LEAVE\n")
		return err
	}

	err = (*publisher).SendMessage(data)
	if err != nil {
		klog.V(1).Infof("SendMessage failed. Err: %v\n", err)
		klog.V(6).Infof("NotificationManager.PublishMessage LEAVE\n")
		return err
	}

	klog.V(6).Infof("NotificationManager.PublishMessage LEAVE\n")
	return nil
}

func (ma *RealtimeAnalyzer) Teardown() error {
	klog.V(6).Infof("NotificationManager.Teardown ENTER\n")

	err := (*ma.rabbitManager).Teardown()
	if err != nil {
		klog.V(1).Infof("rabbitManager.Teardown failed. Err: %v\n", err)
		klog.V(6).Infof("NotificationManager.Stop LEAVE\n")
		return err
	}

	klog.V(4).Infof("Teardown Succeeded\n")
	klog.V(6).Infof("NotificationManager.Teardown LEAVE\n")

	return nil
}
