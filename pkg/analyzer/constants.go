// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package analyzer

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
