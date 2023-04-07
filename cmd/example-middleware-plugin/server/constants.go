// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package server

import (
	"errors"
)

const (
	DefaultPort int = 40000
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")

	// ErrRabbitMgrNil the rabbit manager has not been initialized
	ErrRabbitMgrNil = errors.New("the rabbit manager has not been initialized")
)
