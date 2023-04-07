// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	middlewaresdk "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-plugin-sdk"

	server "github.com/dvonthenen/enterprise-reference-implementation/cmd/example-middleware-plugin/server"
)

func main() {
	// os hooks
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// init
	middlewaresdk.Init(middlewaresdk.EnterpriseInit{
		LogLevel: middlewaresdk.LogLevelStandard, // LogLevelStandard / LogLevelFull / LogLevelTrace / LogLevelVerbose
	})

	middlewareServer, err := server.New(server.ServerOptions{
		CrtFile:   "localhost.crt",
		KeyFile:   "localhost.key",
		RabbitURI: "amqp://guest:guest@localhost:5672",
	})
	if err != nil {
		fmt.Printf("server.New failed. Err: %v\n", err)
		os.Exit(1)
	}

	// init
	err = middlewareServer.Init()
	if err != nil {
		fmt.Printf("middlewareServer.Init() failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start
	fmt.Printf("Starting server...\n")
	err = middlewareServer.Start()
	if err != nil {
		fmt.Printf("middlewareServer.Start() failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	err = middlewareServer.Stop()
	if err != nil {
		fmt.Printf("middlewareServer.Stop() failed. Err: %v\n", err)
	}

	fmt.Printf("Server stopped...\n\n")
}
