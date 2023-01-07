// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	"context"

	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
)

func NewNotificationManager(options NotificationManagerOption) *NotificationManager {
	mgr := &NotificationManager{
		driver:        options.Driver,
		rabbitManager: options.RabbitManager,
		symblClient:   options.SymblClient,
	}
	return mgr
}

/*
	This initializes all of the subscribers to the Symbl Dataminer component

	Each rabbit subscriber listens for a specific Symbl derived/discovered conversation insight and
	is then notified through a callback handler with the original Symbl RealTime API message struct
*/
func (nm *NotificationManager) Init() error {
	klog.V(6).Infof("NotificationManager.Init ENTER\n")

	type InitFunc func(HandlerOptions) *rabbitinterfaces.RabbitMessageHandler
	type MyHandler struct {
		Name string
		Func InitFunc
	}

	// init rabbit clients
	myHandlers := make([]*MyHandler, 0)
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeConversationInit,
		Func: NewConversationInitHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeEntity,
		Func: NewEntityHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeInsight,
		Func: NewInsightHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeMessage,
		Func: NewMessageHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeTopic,
		Func: NewTopicHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeTracker,
		Func: NewTrackerHandler,
	})
	myHandlers = append(myHandlers, &MyHandler{
		Name: interfaces.RabbitExchangeConversationTeardown,
		Func: NewConversationTeardownHandler,
	})

	for _, myHandler := range myHandlers {
		// create session
		ctx := context.Background()
		session := (*nm.driver).NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

		// signal
		handler := myHandler.Func(HandlerOptions{
			Session:     &session,
			SymblClient: nm.symblClient,
			Manager:     nm.rabbitManager,
		})

		_, err := (*nm.rabbitManager).CreateSubscriber(rabbitinterfaces.SubscriberOptions{
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

	err := (*nm.rabbitManager).Init()
	if err != nil {
		klog.V(1).Infof("rabbitManager.Init failed. Err: %v\n", err)
		klog.V(6).Infof("NotificationManager.Init LEAVE\n")
		return err
	}

	klog.V(4).Infof("Init Succeeded\n")
	klog.V(6).Infof("NotificationManager.Init LEAVE\n")

	return nil
}

func (nm *NotificationManager) Teardown() error {
	klog.V(6).Infof("NotificationManager.Teardown ENTER\n")

	err := (*nm.rabbitManager).Teardown()
	if err != nil {
		klog.V(1).Infof("rabbitManager.Teardown failed. Err: %v\n", err)
		klog.V(6).Infof("NotificationManager.Stop LEAVE\n")
		return err
	}

	klog.V(4).Infof("Teardown Succeeded\n")
	klog.V(6).Infof("NotificationManager.Teardown LEAVE\n")

	return nil
}
