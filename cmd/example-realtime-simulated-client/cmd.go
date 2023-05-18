// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package main

// streaming
import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	// sse "github.com/r3labs/sse/v2"
	// streaming "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1"
	microphone "github.com/dvonthenen/symbl-go-sdk/pkg/audio/microphone"
	symbl "github.com/dvonthenen/symbl-go-sdk/pkg/client"
	interfaces "github.com/dvonthenen/symbl-go-sdk/pkg/client/interfaces"

	handlers "github.com/dvonthenen/enterprise-reference-implementation/cmd/example-realtime-simulated-client/handlers"
)

type HeadersContext struct{}

func main() {
	// init
	symbl.Init(symbl.SybmlInit{
		LogLevel: symbl.LogLevelStandard, // LogLevelStandard / LogLevelFull / LogLevelTrace / LogLevelVerbose
	})

	ctx := context.Background()

	// custom headers to enable options
	headers := http.Header{}
	headers.Add("X-ERI-TRANSCRIPTION", "true")
	headers.Add("X-ERI-MESSAGING", "true")
	ctx = interfaces.WithCustomHeaders(ctx, headers)

	// create a new websocket/streaming client
	cfg := symbl.GetDefaultConfig()
	cfg.Speaker.Name = "Jane Doe"
	cfg.Speaker.UserID = "jane.doe@mymail.com"

	/*
		Using WebSocket based Application-Level messages is done here...

		if you wanted to something more meaningful with the output besides just print to the console,
		you would create your own implementation of a message router instead of using the default one
	*/
	// router := streaming.NewDefaultMessageRouter()
	router := handlers.NewMyMessageRouter()

	options := symbl.StreamingOptions{
		SymblConfig:     cfg,
		SymblEndpoint:   "127.0.0.1",
		RedirectService: true,
		SkipServerAuth:  true,
		Callback:        router,
	}

	client, err := symbl.NewStreamClient(ctx, options)
	if err == nil {
		fmt.Println("Login Succeeded!")
	} else {
		fmt.Printf("New failed. Err: %v\n", err)
		os.Exit(1)
	}
	conversationId := client.GetConversationId()
	fmt.Printf("conversationId: %s\n", conversationId)

	err = client.Start()
	if err == nil {
		fmt.Printf("Streaming Session Started!\n")
	} else {
		fmt.Printf("client.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	/*
		Using Server Sent Events (SSE) option

		This is an example of would be how we would consume SSE events...
	*/
	// notificationURI := fmt.Sprintf("https://127.0.0.1/%s/notifications", conversationId)
	// fmt.Printf("notificationURI: %s\n", notificationURI)
	// notifications := sse.NewClient(notificationURI, func(c *sse.Client) {
	// 	/* #nosec G402 */
	// 	c.Connection.Transport = &http.Transport{
	// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// 	}
	// 	c.RedirectService = true
	// 	c.SkipServerAuth = true
	// })

	// // listen for events
	// stopPoll := make(chan struct{})
	// eventsChan := make(chan *sse.Event)
	// go func() {
	// 	notifications.SubscribeChan("messages", eventsChan)
	// }()

	// waitForMessages := func(eventsChan chan *sse.Event, stopChan chan struct{}) {
	// 	for {
	// 		select {
	// 		default:
	// 			fmt.Printf("Waiting for event...\n")
	// 			event := <-eventsChan

	// 			// you typically would unmarshal for ClientMessageType, then use MessageType to determine
	// 			// which struct to use to rebuild the object
	// 			// in this case, we are only using ClientTrackerMessage since that is the only app-level
	// 			// message that is implemented
	// 			var cn interfaces.ClientTrackerMessage
	// 			err = json.Unmarshal(event.Data, &cn)
	// 			if err != nil {
	// 				fmt.Printf("json.Unmarshal failed. Err: %v\n", err)
	// 				continue
	// 			}

	// 			prettyJson, err := prettyjson.Format(event.Data)
	// 			if err != nil {
	// 				fmt.Printf("prettyjson.Marshal failed. Err: %v\n", err)
	// 				continue
	// 			}

	// 			fmt.Printf("\n\n")
	// 			fmt.Printf("Application Event:\n")
	// 			fmt.Printf("%s\n", prettyJson)
	// 			fmt.Printf("\n\n")
	// 			fmt.Printf("Human Readable:\n")
	// 			fmt.Printf("This tracker (%s) was previous mentioned, here are the details:\n", cn.Data.Message.Correlation)
	// 			fmt.Printf("What you said:\n%s\n", cn.Data.Message.CurrentContent)
	// 			fmt.Printf("Previously mentioned by: %s / %s \n", cn.Data.Author.Name, cn.Data.Author.Email)
	// 			fmt.Printf("What they said:\n%s\n", cn.Data.Message.PreviousContent)
	// 			fmt.Printf("Commonality: %s == %s\n", cn.Data.Message.CurrentMatch, cn.Data.Message.PreviousMatch)
	// 			fmt.Printf("\n\n")

	// 		case <-stopChan:
	// 			return
	// 		}
	// 	}
	// }

	// fmt.Print("Listening for events\n")
	// go waitForMessages(eventsChan, stopPoll)

	// delay...
	time.Sleep(time.Second)

	// mic stuf
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// init microphone library
	microphone.Initialize()

	mic, err := microphone.New(microphone.AudioConfig{
		InputChannels: 1,
		SamplingRate:  16000,
	})
	if err != nil {
		fmt.Printf("Initialize failed. Err: %v\n", err)
		os.Exit(1)
	}

	// start the mic
	err = mic.Start()
	if err != nil {
		fmt.Printf("mic.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	go func() {
		// this is a blocking call
		mic.Stream(client)
	}()

	fmt.Print("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// close stream
	err = mic.Stop()
	if err != nil {
		fmt.Printf("mic.Stop failed. Err: %v\n", err)
		os.Exit(1)
	}
	microphone.Teardown()

	// close client
	client.Stop()

	fmt.Printf("Succeeded!\n\n")
}
