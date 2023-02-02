// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

/*
Package provides reference implementation for Go to handle streaming
audio generically from CPaaS platforms.

GitHub repo: https://github.com/dvonthenen/enterprise-reference-implementation
*/
package sdk

import (
	_ "github.com/dvonthenen/enterprise-reference-implementation/cmd/example-your-middleware"
	_ "github.com/dvonthenen/enterprise-reference-implementation/cmd/symbl-proxy-dataminer"
	_ "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-analyzer"
	_ "github.com/dvonthenen/enterprise-reference-implementation/pkg/proxy-dataminer"
)
