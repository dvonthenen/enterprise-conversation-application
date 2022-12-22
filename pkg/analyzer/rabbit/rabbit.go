// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package rabbit

import (
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
)

func New(options RabbitManagerOptions) *RabbitManager {
	rabbit := &RabbitManager{
		rabbitConn:  options.Connection,
		subscribers: make(map[string]*Subscriber),
		publishers:  make(map[string]*Publisher),
	}
	return rabbit
}

func (rm *RabbitManager) CreateSubscriber(options interfaces.SubscriberOptions) (*Subscriber, error) {
	klog.V(6).Infof("RabbitManager.CreateSubscriber ENTER\n")

	ch, err := rm.rabbitConn.Channel()
	if err != nil {
		klog.V(1).Infof("Channel() failed. Err: %v\n", err)
		klog.V(6).Infof("RabbitManager.CreateSubscriber LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("ExchangeDeclare: %n\n", options.Name)
	err = ch.ExchangeDeclare(
		options.Name, // name
		"fanout",     // type
		true,         // durable
		true,         // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		klog.V(1).Infof("ExchangeDeclare(%s) failed. Err: %v\n", options.Name, err)
		klog.V(6).Infof("RabbitManager.CreateSubscriber LEAVE\n")
		return nil, err
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
		klog.V(6).Infof("RabbitManager.CreateSubscriber LEAVE\n")
		return nil, err
	}

	klog.V(3).Infof("QueueBind: %n\n", options.Name)
	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		options.Name, // exchange
		false,
		nil)
	if err != nil {
		klog.V(1).Infof("QueueBind() failed. Err: %v\n", err)
		klog.V(6).Infof("RabbitManager.CreateSubscriber LEAVE\n")
		return nil, err
	}

	// setup subscriber
	options.Channel = ch
	options.Queue = &q
	subscriber := NewSubscriber(options)

	rm.mu.Lock()
	rm.subscribers[options.Name] = subscriber
	rm.mu.Unlock()

	klog.V(4).Infof("RabbitManager.CreateSubscriber(%s) Succeeded\n", options.Name)
	klog.V(6).Infof("RabbitManager.CreateSubscriber LEAVE\n")

	return subscriber, nil
}

func (rm *RabbitManager) CreatePublisher(options interfaces.PublisherOptions) error {
	klog.V(6).Infof("RabbitManager.CreatePublisher ENTER\n")

	ch, err := rm.rabbitConn.Channel()
	if err != nil {
		klog.V(1).Infof("Channel() failed. Err: %v\n", err)
		klog.V(6).Infof("RabbitManager.CreatePublisher LEAVE\n")
		return err
	}

	// setup publisher
	options.Channel = ch
	publisher := NewPublisher(options)

	err = publisher.Init()
	if err != nil {
		klog.V(1).Infof("Init() failed. Err: %v\n", err)
		klog.V(6).Infof("RabbitManager.CreatePublisher LEAVE\n")
		return err
	}

	rm.mu.Lock()
	rm.publishers[options.Name] = publisher
	rm.mu.Unlock()

	klog.V(4).Infof("RabbitManager.CreatePublisher(%s) Succeeded\n", options.Name)
	klog.V(6).Infof("RabbitManager.CreatePublisher LEAVE\n")

	return nil
}

func (rm *RabbitManager) Start() error {
	klog.V(6).Infof("RabbitManager.Start ENTER\n")

	for msgType, subscriber := range rm.subscribers {
		err := subscriber.Start()
		if err == nil {
			klog.V(3).Infof("subscriber.Start(%s) Succeeded\n", msgType)
		} else {
			klog.V(1).Infof("subscriber.Start() failed. Err: %v\n", err)
			klog.V(6).Infof("RabbitManager.Start LEAVE\n")
		}
	}

	klog.V(4).Infof("RabbitManager.Start Succeeded\n")
	klog.V(6).Infof("RabbitManager.Start LEAVE\n")

	return nil
}

func (rm *RabbitManager) Stop() error {
	klog.V(6).Infof("RabbitManager.Stop ENTER\n")

	for _, subscriber := range rm.subscribers {
		err := subscriber.Stop()
		if err != nil {
			klog.V(1).Infof("subscriber.Stop() failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("RabbitManager.Stop Succeeded\n")
	klog.V(6).Infof("RabbitManager.Stop LEAVE\n")

	return nil
}

func (rm *RabbitManager) PublishMessageByChannelName(name string, data []byte) error {
	klog.V(6).Infof("RabbitManager.PublishMessageByChannelName ENTER\n")

	klog.V(3).Infof("Publishing to: %s\n", name)
	klog.V(3).Infof("Data: %s\n", string(data))
	publisher := rm.publishers[name]
	if publisher == nil {
		klog.V(1).Infof("Publisher(%s) not found\n", name)
		klog.V(6).Infof("RabbitManager.PublishMessageByChannelName LEAVE\n")
		return ErrPublisherNotFound
	}

	err := publisher.SendMessage(data)
	if err != nil {
		klog.V(1).Infof("SendMessage() failed. Err: %v\n", err)
		klog.V(6).Infof("RabbitManager.PublishMessageByChannelName LEAVE\n")
		return err
	}

	klog.V(4).Infof("RabbitManager.PublishMessageByChannelName Succeeded\n")
	klog.V(6).Infof("RabbitManager.PublishMessageByChannelName LEAVE\n")

	return nil
}

func (rm *RabbitManager) DeletePublisher(name string) error {
	klog.V(6).Infof("RabbitManager.DeletePublisher ENTER\n")

	// find
	klog.V(3).Infof("Deleting Publisher: %s\n", name)
	publisher := rm.publishers[name]
	if publisher == nil {
		klog.V(1).Infof("Publisher(%s) not found\n", name)
		klog.V(6).Infof("RabbitManager.DeletePublisher LEAVE\n")
		return ErrPublisherNotFound
	}

	// clean up
	err := publisher.Teardown()
	if err != nil {
		klog.V(1).Infof("Teardown failed. Err: %v\n", err)
	}

	rm.mu.Lock()
	delete(rm.publishers, name)
	rm.mu.Unlock()

	klog.V(4).Infof("RabbitManager.DeletePublisher Succeeded\n")
	klog.V(6).Infof("RabbitManager.DeletePublisher LEAVE\n")

	return nil
}

func (rm *RabbitManager) Teardown() error {
	klog.V(6).Infof("RabbitManager.Teardown ENTER\n")

	for _, subscriber := range rm.subscribers {
		err := subscriber.Stop()
		if err != nil {
			klog.V(1).Infof("subscriber.Stop() failed. Err: %v\n", err)
		}
	}
	for _, publisher := range rm.publishers {
		err := publisher.Teardown()
		if err != nil {
			klog.V(1).Infof("subscriber.Stop() failed. Err: %v\n", err)
		}
	}

	rm.mu.Lock()
	rm.subscribers = make(map[string]*Subscriber)
	rm.publishers = make(map[string]*Publisher)
	rm.mu.Unlock()

	klog.V(4).Infof("RabbitManager.Teardown Succeeded\n")
	klog.V(6).Infof("RabbitManager.Teardown LEAVE\n")

	return nil
}
