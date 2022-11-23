// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package server

import (
	"errors"
)

const (
	DefaultStartPort int = 50000
	DefaultEndPort   int = 60000
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")
)
