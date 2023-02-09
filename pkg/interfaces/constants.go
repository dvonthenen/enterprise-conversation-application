// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfaces

import (
	interfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

const (
	// these are convenience pass-through consts
	MessageTypeRecognitionResult string = interfaces.MessageTypeRecognitionResult
	MessageTypeMessageResponse   string = interfaces.MessageTypeMessageResponse

	// user-defined messages
	MessageTypeUserDefined string = interfaces.MessageTypeUserDefined
)
