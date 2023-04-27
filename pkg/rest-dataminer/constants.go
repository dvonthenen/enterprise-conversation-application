// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package dataminer

import (
	"errors"
)

type AuthType int64

const (
	AuthTypeDefault    AuthType = iota
	AuthTypeReuseToken          = 1
	AuthTypeEnvVars             = 2
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")
)
