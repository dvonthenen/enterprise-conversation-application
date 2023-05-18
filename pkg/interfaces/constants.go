// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

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
