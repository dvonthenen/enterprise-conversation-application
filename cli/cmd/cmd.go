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
		LogLevel: server.LogLevelTrace,
	})

	server := server.New(server.ServerOptions{
		CrtFile: "localhost.crt",
		KeyFile: "localhost.key",
	})

	// start... blocking call
	go func() {
		server.Start()
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	fmt.Printf("Succeeded!\n\n")
}
