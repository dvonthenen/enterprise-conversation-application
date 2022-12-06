// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	"github.com/dvonthenen/enterprise-reference-implementation/pkg/server"
)

func main() {
	// os hooks
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// init
	server.Init(server.EnterpriseInit{
		LogLevel: server.LogLevelFull,
	})

	server, err := server.New(server.ServerOptions{
		CrtFile: "localhost.crt",
		KeyFile: "localhost.key",
	})
	if err != nil {
		fmt.Printf("server.New failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start... blocking call
	go func() {
		err := server.Start()
		if err != nil {
			fmt.Printf("server.Start() failed. Err: %v\n", err)
		}
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	err = server.Stop()
	if err != nil {
		fmt.Printf("server.Stop() failed. Err: %v\n", err)
	}

	fmt.Printf("Succeeded!\n\n")
}
