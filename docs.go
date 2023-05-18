// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

/*
Package provides an off the shell implementation for Go to handle streaming and async
conversations generically from a variety of platforms, including CPaaS, CRM, Salesforce, etc

GitHub repo: https://github.com/dvonthenen/enterprise-conversation-application
*/
package sdk

import (
	_ "github.com/dvonthenen/enterprise-conversation-application/cmd/example-asynchronous-plugin"
	_ "github.com/dvonthenen/enterprise-conversation-application/cmd/example-realtime-plugin"
	_ "github.com/dvonthenen/enterprise-conversation-application/cmd/symbl-proxy-dataminer"
	_ "github.com/dvonthenen/enterprise-conversation-application/cmd/symbl-rest-dataminer"
	_ "github.com/dvonthenen/enterprise-conversation-application/pkg/middleware-plugin-sdk"
	_ "github.com/dvonthenen/enterprise-conversation-application/pkg/proxy-dataminer"
	_ "github.com/dvonthenen/enterprise-conversation-application/pkg/rest-dataminer"
)
