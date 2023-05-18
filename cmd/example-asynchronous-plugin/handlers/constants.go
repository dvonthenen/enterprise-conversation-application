// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package handlers

import "errors"

var (
	// ErrUnhandledMessage runhandled message from example-asynchronous-plugin
	ErrUnhandledMessage = errors.New("unhandled message from example-asynchronous-plugin")
)
