// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package rabbit

import (
	"strings"

	klog "k8s.io/klog/v2"
)

func New(options RabbitManagerOptions) *RabbitManager {
	rabbit := &RabbitManager{
		rabbitConn:  options.Connection,
		subscribers: make(map[string]*Subscriber),
	}
	return rabbit
}

func (rm *RabbitManager) CreateSubscription(options CreateOptions) (*Subscriber, error) {
	klog.V(6).Infof("rabbit.CreateSubscription ENTER\n")

	ch, err := rm.rabbitConn.Channel()
	if err != nil {
		klog.V(1).Infof("Channel() failed. Err: %v\n", err)
		klog.V(6).Infof("rabbit.CreateSubscription LEAVE\n")
		return nil, err
	}

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
		klog.V(6).Infof("rabbit.CreateSubscription LEAVE\n")
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
		klog.V(6).Infof("rabbit.CreateSubscription LEAVE\n")
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		options.Name, // exchange
		false,
		nil)
	if err != nil {
		klog.V(1).Infof("QueueBind() failed. Err: %v\n", err)
		klog.V(6).Infof("rabbit.CreateSubscription LEAVE\n")
		return nil, err
	}

	subscribe := NewSubscribe(SubscribeOptions{
		Name:    options.Name,
		Channel: ch,
		Queue:   &q,
		Handler: options.Handler,
	})

	rm.mu.Lock()
	rm.subscribers[options.Name] = subscribe
	rm.mu.Unlock()

	klog.V(3).Infof("CreateSubscription(%s) Succeeded\n", options.Name)
	klog.V(6).Infof("rabbit.CreateSubscription LEAVE\n")

	return subscribe, nil
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

	klog.V(3).Infof("RabbitManager.Start Succeeded\n")
	klog.V(6).Infof("RabbitManager.Start LEAVE\n")

	return nil
}

func (rm *RabbitManager) StartByName(name string) error {
	klog.V(6).Infof("RabbitManager.StartByName(%s) ENTER\n", name)

	for msgType, subscriber := range rm.subscribers {
		if !strings.EqualFold(msgType, name) {
			continue
		}

		err := subscriber.Start()
		if err == nil {
			klog.V(3).Infof("subscriber.StartByName(%s) Succeeded\n", msgType)
		} else {
			klog.V(1).Infof("subscriber.StartByName() failed. Err: %v\n", err)
			klog.V(6).Infof("RabbitManager.StartByName LEAVE\n")
			return err
		}
	}

	klog.V(3).Infof("RabbitManager.StartByName Succeeded\n")
	klog.V(6).Infof("RabbitManager.StartByName LEAVE\n")

	return nil
}

func (rm *RabbitManager) StopByName(name string) error {
	klog.V(6).Infof("RabbitManager.StopByName ENTER\n")

	for msgType, subscriber := range rm.subscribers {
		if !strings.EqualFold(msgType, name) {
			continue
		}

		err := subscriber.Stop()
		if err == nil {
			klog.V(3).Infof("subscriber.StopByName(%s) Succeeded\n", msgType)
		} else {
			klog.V(1).Infof("subscriber.StopByName() failed. Err: %v\n", err)
			klog.V(6).Infof("RabbitManager.StopByName ENTER\n")
			return err
		}
	}

	klog.V(3).Infof("RabbitManager.StopByName Succeeded\n")
	klog.V(6).Infof("RabbitManager.StopByName ENTER\n")

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

	klog.V(3).Infof("RabbitManager.Stop Succeeded\n")
	klog.V(6).Infof("RabbitManager.Stop LEAVE\n")

	return nil
}

func (rm *RabbitManager) DeleteByName(name string) error {
	klog.V(6).Infof("RabbitManager.DeleteByName(%s) ENTER\n", name)

	for key, subscriber := range rm.subscribers {
		if !strings.EqualFold(key, name) {
			continue
		}
		err := subscriber.Stop()
		if err != nil {
			klog.V(1).Infof("subscriber.Stop() failed. Err: %v\n", err)
			klog.V(6).Infof("RabbitManager.DeleteByName LEAVE\n")
			return err
		}

		rm.mu.Lock()
		delete(rm.subscribers, key)
		rm.mu.Unlock()
	}

	klog.V(3).Infof("RabbitManager.DeleteByName Succeeded\n")
	klog.V(6).Infof("RabbitManager.DeleteByName LEAVE\n")

	return nil
}

func (rm *RabbitManager) Delete() error {
	klog.V(6).Infof("RabbitManager.Delete ENTER\n")

	for _, subscriber := range rm.subscribers {
		err := subscriber.Stop()
		if err != nil {
			klog.V(1).Infof("subscriber.Stop() failed. Err: %v\n", err)
		}
	}

	rm.mu.Lock()
	rm.subscribers = make(map[string]*Subscriber)
	rm.mu.Unlock()

	klog.V(3).Infof("RabbitManager.Delete Succeeded\n")
	klog.V(6).Infof("RabbitManager.Delete LEAVE\n")

	return nil
}
