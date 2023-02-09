// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package utils

import "errors"

const (
	// default number of message to cache at any given time
	DefaultNumOfMsgToCache int = 50
)

var (
	// ErrItemNotFound item not found
	ErrItemNotFound = errors.New("item not found")
)
