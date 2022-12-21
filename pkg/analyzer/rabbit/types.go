// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package rabbit

import (
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
)

/*
	Subscriber
*/
type Subscriber struct {
	options  interfaces.SubscriberOptions
	channel  *amqp.Channel
	queue    *amqp.Queue
	stopChan chan struct{}
	handler  *interfaces.RabbitMessageHandler
	running  bool
}

/*
	Publisher
*/
type Publisher struct {
	options interfaces.PublisherOptions
	name    string
	channel *amqp.Channel
}

/*
	The one that manages everything
*/
type RabbitManagerOptions struct {
	Connection *amqp.Connection
}

type RabbitManager struct {
	// housekeeping
	publishers  map[string]*Publisher
	subscribers map[string]*Subscriber
	mu          sync.Mutex

	// rabbitmq
	rabbitConn *amqp.Connection
}
