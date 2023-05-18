// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package middleware

import (
	rabbit "github.com/dvonthenen/rabbitmq-manager/pkg"
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	klog "k8s.io/klog/v2"

	middlewareinterfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/interfaces"
	router "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk/router/asynchronous"
	shared "github.com/dvonthenen/enterprise-reference-implementation/pkg/shared"
)

func NewAsynchronousAnalyzer(options AsynchronousAnalyzerOption) (*AsynchronousAnalyzer, error) {
	// setup rabbit manager
	rabbitMgr, err := rabbit.New(rabbitinterfaces.ManagerOptions{
		RabbitURI: options.RabbitURI,
	})
	if err != nil {
		klog.V(1).Infof("rabbit.New failed. Err: %v\n", err)
		return nil, err
	}

	// create middleware
	mgr := &AsynchronousAnalyzer{
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
func (aa *AsynchronousAnalyzer) Init() error {
	klog.V(6).Infof("AsynchronousAnalyzer.Init ENTER\n")

	type InitFunc func(router.HandlerOptions) *rabbitinterfaces.RabbitMessageHandler
	type MyHandler struct {
		Name string
		Func InitFunc
	}

	// init rabbit clients
	myHandlers := make([]*MyHandler, 0)
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitAsyncConversationInit,
		Func: router.NewConversationInitHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitAsyncMessage,
		Func: router.NewMessageHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitAsyncQuestion,
		Func: router.NewQuestionHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitAsyncFollowUp,
		Func: router.NewFollowUpHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitAsyncActionItem,
		Func: router.NewActionItemHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitAsyncTopic,
		Func: router.NewTopicHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitAsyncTracker,
		Func: router.NewTrackerHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitAsyncEntity,
		Func: router.NewEntityHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: shared.RabbitAsyncConversationTeardown,
		Func: router.NewConversationTeardownHandler,
	})

	for _, myHandler := range myHandlers {
		// create subscriber
		handler := myHandler.Func(router.HandlerOptions{
			Manager:  aa.rabbitManager,
			Callback: aa.callback,
		})

		_, err := (*aa.rabbitManager).CreateSubscriber(rabbitinterfaces.SubscriberOptions{
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
	err := (*aa.rabbitManager).Init()
	if err != nil {
		klog.V(1).Infof("rabbitManager.Init failed. Err: %v\n", err)
		klog.V(6).Infof("AsynchronousAnalyzer.Init LEAVE\n")
		return err
	}

	klog.V(4).Infof("Init Succeeded\n")
	klog.V(6).Infof("AsynchronousAnalyzer.Init LEAVE\n")

	return nil
}

func (aa *AsynchronousAnalyzer) PublishMessage(name string, data []byte) error {
	klog.V(6).Infof("AsynchronousAnalyzer.PublishMessage ENTER\n")

	publisher, err := (*aa.rabbitManager).GetPublisherByName(name)
	if err != nil {
		klog.V(1).Infof("GetPublisherByName failed. Err: %v\n", err)
		klog.V(6).Infof("AsynchronousAnalyzer.PublishMessage LEAVE\n")
		return err
	}

	err = (*publisher).SendMessage(data)
	if err != nil {
		klog.V(1).Infof("SendMessage failed. Err: %v\n", err)
		klog.V(6).Infof("AsynchronousAnalyzer.PublishMessage LEAVE\n")
		return err
	}

	klog.V(6).Infof("AsynchronousAnalyzer.PublishMessage LEAVE\n")
	return nil
}

func (aa *AsynchronousAnalyzer) Teardown() error {
	klog.V(6).Infof("AsynchronousAnalyzer.Teardown ENTER\n")

	err := (*aa.rabbitManager).Teardown()
	if err != nil {
		klog.V(1).Infof("rabbitManager.Teardown failed. Err: %v\n", err)
		klog.V(6).Infof("AsynchronousAnalyzer.Stop LEAVE\n")
		return err
	}

	klog.V(4).Infof("Teardown Succeeded\n")
	klog.V(6).Infof("AsynchronousAnalyzer.Teardown LEAVE\n")

	return nil
}
