// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	analyzer "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer"
)

func main() {
	// os hooks
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// init
	analyzer.Init(analyzer.EnterpriseInit{
		LogLevel: analyzer.LogLevelStandard, // LogLevelStandard / LogLevelFull
	})

	analyzer, err := analyzer.New(analyzer.ServerOptions{
		CrtFile:     "localhost.crt",
		KeyFile:     "localhost.key",
		RabbitMQURI: "amqp://guest:guest@localhost:5672",
	})
	if err != nil {
		fmt.Printf("analyzer.New failed. Err: %v\n", err)
		os.Exit(1)
	}

	// init
	err = analyzer.Init()
	if err != nil {
		fmt.Printf("Starting analyzer. Err: %v\n", err)
		os.Exit(1)
	}

	// start
	err = analyzer.Start()
	if err != nil {
		fmt.Printf("analyzer.Start() failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	err = analyzer.Stop()
	if err != nil {
		fmt.Printf("analyzer.Stop() failed. Err: %v\n", err)
	}

	fmt.Printf("Succeeded!\n\n")
}
