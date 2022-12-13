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
type SubscribeOptions struct {
	Name    string
	Channel *amqp.Channel
	Queue   *amqp.Queue
	Handler *interfaces.RabbitMessageHandler
}

type Subscriber struct {
	options  SubscribeOptions
	channel  *amqp.Channel
	queue    *amqp.Queue
	stopChan chan struct{}
	handler  *interfaces.RabbitMessageHandler
	running  bool
}

/*
	The one that manages everything
*/
type RabbitManagerOptions struct {
	Connection *amqp.Connection
}

type RabbitManager struct {
	// housekeeping
	subscribers map[string]*Subscriber
	mu          sync.Mutex

	// rabbitmq
	rabbitConn *amqp.Connection
}
