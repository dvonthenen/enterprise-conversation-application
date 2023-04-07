// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

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

	// ErrInvalidNotifyConfig invalid notify configuration
	ErrInvalidNotifyConfig = errors.New("invalid notify configuration")

	// ErrUnknownNotifyType unknown notify message type
	ErrUnknownNotifyType = errors.New("unknown notify message type")
)
