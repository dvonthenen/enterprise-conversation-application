// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"errors"
)

const (
	DefaultSymblWebSocket string = "wss://api.symbl.ai"

	MessageTypeMessage              string = "message"
	MessageTypeInitConversation     string = "conversation_created"
	MessageTypeTeardownConversation string = "conversation_completed"
)

type ClientNotifyType int

const (
	ClientNotifyTypeWebSocket ClientNotifyType = iota
	ClientNotifyTypeServerSendEvent
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")
)
