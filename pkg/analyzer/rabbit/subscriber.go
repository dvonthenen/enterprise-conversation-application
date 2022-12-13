// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package rabbit

import (
	klog "k8s.io/klog/v2"
)

func NewSubscribe(options SubscribeOptions) *Subscriber {
	rabbit := &Subscriber{
		options: options,
		channel: options.Channel,
		queue:   options.Queue,
		handler: options.Handler,
		running: false,
	}
	return rabbit
}

func (rs *Subscriber) Start() error {
	klog.V(6).Infof("subscribe.Start ENTER\n")

	if rs.running {
		klog.V(1).Infof("Subscribe already running\n")
		klog.V(6).Infof("subscribe.Start LEAVE\n")
		return nil
	}

	msgs, err := rs.channel.Consume(
		rs.queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		klog.V(1).Infof("QueueBind() failed. Err: %v\n", err)
		klog.V(6).Infof("subscribe.Start LEAVE\n")
		return err
	}

	klog.V(5).Infof("subscribe.Start Running message loop...\n")
	rs.running = true
	go func() {
		for {
			select {
			default:
				for d := range msgs {
					klog.V(5).Infof(" [x] %s\n", d.Body)

					err := (*rs.handler).ProcessMessage(d.Body)
					if err != nil {
						klog.V(1).Infof("ProcessMessage() failed. Err: %v\n", err)
					}
				}
			case <-rs.stopChan:
				klog.V(6).Infof("Exiting Subscription Loop\n")
				return
			}
		}
	}()

	return nil
}

func (rs *Subscriber) Stop() error {
	rs.running = false
	close(rs.stopChan)
	return nil
}
