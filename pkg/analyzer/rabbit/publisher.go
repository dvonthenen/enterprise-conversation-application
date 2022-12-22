// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package rabbit

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
)

func NewPublisher(options interfaces.PublisherOptions) *Publisher {
	rabbit := &Publisher{
		options: options,
		name:    options.Name,
		channel: options.Channel,
	}
	return rabbit
}

func (p *Publisher) Init() error {
	klog.V(6).Infof("Publisher.Init ENTER\n")

	err := p.channel.ExchangeDeclare(
		p.name,   // name
		"fanout", // type
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		klog.V(1).Infof("ExchangeDeclare failed. Err: %v\n", err)
		klog.V(6).Infof("Publisher.Init LEAVE\n")
		return err
	}

	klog.V(4).Infof("Publisher.Init Succeeded\n")
	klog.V(6).Infof("Publisher.Init LEAVE\n")

	return nil
}

func (p *Publisher) SendMessage(data []byte) error {
	klog.V(6).Infof("Publisher.SendMessage ENTER\n")

	ctx := context.Background()

	klog.V(3).Infof("Publishing to: %s\n", p.options.Name)
	klog.V(3).Infof("Data: %s\n", string(data))
	err := p.channel.PublishWithContext(ctx,
		p.options.Name, // exchange
		"",             // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	if err != nil {
		klog.V(1).Infof("PublishWithContext failed. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("Publisher.SendMessage succeeded\n%s\n", string(data))
	klog.V(6).Infof("Publisher.SendMessage LEAVE\n")

	return nil
}

func (p *Publisher) Teardown() error {
	klog.V(6).Infof("Publisher.Teardown ENTER\n")

	if p.channel != nil {
		p.channel.Close()
		p.channel = nil
	}

	klog.V(4).Infof("Publisher.Teardown Succeeded\n")
	klog.V(6).Infof("Publisher.Teardown LEAVE\n")

	return nil
}
