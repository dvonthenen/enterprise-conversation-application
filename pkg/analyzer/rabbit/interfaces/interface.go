// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfaces

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SubscriberOptions struct {
	Name    string
	Channel *amqp.Channel
	Queue   *amqp.Queue
	Handler *RabbitMessageHandler
}

type PublisherOptions struct {
	Name    string
	Channel *amqp.Channel
}

type RabbitMessageHandler interface {
	ProcessMessage(byData []byte) error
}

type RabbitManagerHandler interface {
	CreatePublisher(options PublisherOptions) error
	PublishMessageByChannelName(name string, data []byte) error
	DeletePublisher(name string) error
}
