// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	dataminer "github.com/dvonthenen/enterprise-reference-implementation/pkg/dataminer"
)

func main() {
	// os hooks
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// init
	dataminer.Init(dataminer.EnterpriseInit{
		LogLevel: dataminer.LogLevelStandard,
	})

	dataminer, err := dataminer.New(dataminer.ServerOptions{
		CrtFile:     "localhost.crt",
		KeyFile:     "localhost.key",
		RabbitMQURI: "amqp://guest:guest@localhost:5672",
	})
	if err != nil {
		fmt.Printf("dataminer.New failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start
	err = dataminer.Start()
	if err != nil {
		fmt.Printf("dataminer.Start() failed. Err: %v\n", err)
	}

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	err = dataminer.Stop()
	if err != nil {
		fmt.Printf("dataminer.Stop() failed. Err: %v\n", err)
	}

	fmt.Printf("Succeeded!\n\n")
}
