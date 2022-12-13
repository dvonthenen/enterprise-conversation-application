// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

/*
Package provides reference implementation for Go to handle streaming
audio generically from CPaaS platforms.

GitHub repo: https://github.com/dvonthenen/enterprise-reference-implementation
*/
package sdk

import (
	_ "github.com/dvonthenen/enterprise-reference-implementation/cli/cmd/example-your-middleware"
	_ "github.com/dvonthenen/enterprise-reference-implementation/cli/cmd/symbl-dataminer"
	_ "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer"
	_ "github.com/dvonthenen/enterprise-reference-implementation/pkg/dataminer"
)
