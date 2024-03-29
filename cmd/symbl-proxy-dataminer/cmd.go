// Copyright 2023 Enterprise Conversation Application contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	dataminer "github.com/dvonthenen/enterprise-conversation-application/pkg/proxy-dataminer"
)

func main() {
	// os hooks
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// init
	dataminer.Init(dataminer.EnterpriseInit{
		LogLevel: dataminer.LogLevelStandard, // LogLevelStandard / LogLevelFull / LogLevelTrace / LogLevelVerbose
	})

	dataminer, err := dataminer.New(dataminer.ServerOptions{
		CrtFile:   "localhost.crt",
		KeyFile:   "localhost.key",
		RabbitURI: "amqp://guest:guest@localhost:5672",
	})
	if err != nil {
		fmt.Printf("dataminer.New failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start
	fmt.Printf("Starting server...\n")
	err = dataminer.Start()
	if err != nil {
		fmt.Printf("dataminer.Start() failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	err = dataminer.Stop()
	if err != nil {
		fmt.Printf("dataminer.Stop() failed. Err: %v\n", err)
	}

	fmt.Printf("Server stopped...\n\n")
}
