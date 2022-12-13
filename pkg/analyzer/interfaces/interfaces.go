// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package interfaces

type PushNotificationCallback interface {
	PushNotification(msg string) error
}
