// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package handlers

import "errors"

const (
	// default number of message to cache at any given time
	DefaultNumOfMsgToCache int = 50
)

var (
	// ErrUnhandledMessage runhandled message from symbl-proxy-dataminer
	ErrUnhandledMessage = errors.New("unhandled message from symbl-proxy-dataminer")

	// ErrItemNotFound item not found
	ErrItemNotFound = errors.New("item not found")
)
