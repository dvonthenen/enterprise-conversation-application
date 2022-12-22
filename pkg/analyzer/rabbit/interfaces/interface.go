// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfaces

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

/*
	Configuration Options for Rabbit Subscriber/Publisher
*/
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

/*
	Each message handler (Topic, Entity, etc) implements this callback which is called when the
	Rabbit Subscriber receives a conversation insight from the Symbl dataminer component
*/
type RabbitMessageHandler interface {
	ProcessMessage(byData []byte) error
}

/*
	Interface to the Rabbit Manager which allows the handlers to create and delete
	rabbit queues for each Real-Time Streaming connection
*/
type RabbitManagerHandler interface {
	CreatePublisher(options PublisherOptions) error
	PublishMessageByChannelName(name string, data []byte) error
	DeletePublisher(name string) error
}
