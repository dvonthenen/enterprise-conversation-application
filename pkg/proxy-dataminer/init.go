// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package dataminer

import (
	"flag"
	"strconv"

	klog "k8s.io/klog/v2"
)

type LogLevel int64

const (
	LogLevelDefault   LogLevel = iota
	LogLevelErrorOnly          = 1
	LogLevelStandard           = 2
	LogLevelElevated           = 3
	LogLevelFull               = 4
	LogLevelDebug              = 5
	LogLevelTrace              = 6
	LogLevelVerbose            = 7
)

type EnterpriseInit struct {
	LogLevel      LogLevel
	DebugFilePath string
}

func Init(init EnterpriseInit) {
	if init.LogLevel == LogLevelDefault {
		init.LogLevel = LogLevelStandard
	}

	klog.InitFlags(nil)
	flag.Set("v", strconv.FormatInt(int64(init.LogLevel), 10))
	if init.DebugFilePath != "" {
		flag.Set("logtostderr", "false")
		flag.Set("log_file", init.DebugFilePath)
	}
	flag.Parse()
}
