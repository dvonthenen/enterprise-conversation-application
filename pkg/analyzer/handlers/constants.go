// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	"errors"
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")

	// ErrUserCallbackNotDefined user callback object not defined
	ErrUserCallbackNotDefined = errors.New("user callback object not defined")
)