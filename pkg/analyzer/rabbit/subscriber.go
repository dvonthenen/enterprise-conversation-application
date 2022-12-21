// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package rabbit

import (
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
)

func NewSubscriber(options interfaces.SubscriberOptions) *Subscriber {
	rabbit := &Subscriber{
		options: options,
		channel: options.Channel,
		queue:   options.Queue,
		handler: options.Handler,
		running: false,
	}
	return rabbit
}

func (s *Subscriber) Start() error {
	klog.V(6).Infof("Subscriber.Start ENTER\n")

	if s.running {
		klog.V(1).Infof("Subscribe already running\n")
		klog.V(6).Infof("Subscriber.Start LEAVE\n")
		return nil
	}

	msgs, err := s.channel.Consume(
		s.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		klog.V(1).Infof("Consume failed. Err: %v\n", err)
		klog.V(6).Infof("Subscriber.Start LEAVE\n")
		return err
	}

	klog.V(3).Infof("Subscriber.Start Running message loop...\n")
	s.running = true
	go func() {
		for {
			select {
			default:
				for d := range msgs {
					klog.V(5).Infof(" [x] %s\n", d.Body)

					err := (*s.handler).ProcessMessage(d.Body)
					if err != nil {
						klog.V(1).Infof("ProcessMessage() failed. Err: %v\n", err)
					}
				}
			case <-s.stopChan:
				klog.V(5).Infof("Exiting Subscriber Loop\n")
				return
			}
		}
	}()

	klog.V(4).Infof("Subscriber.Start Succeeded\n")
	klog.V(6).Infof("Subscriber.Start LEAVE\n")

	return nil
}

func (s *Subscriber) Stop() error {
	klog.V(6).Infof("Subscriber.Stop ENTER\n")

	s.running = false

	close(s.stopChan)
	<-s.stopChan

	if s.channel != nil {
		s.channel.Close()
		s.channel = nil
	}

	klog.V(4).Infof("Subscriber.Stop Succeeded\n")
	klog.V(6).Infof("Subscriber.Stop LEAVE\n")

	return nil
}
