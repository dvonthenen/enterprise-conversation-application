// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package server

// streaming
import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	klog "k8s.io/klog/v2"

	instance "github.com/dvonthenen/enterprise-reference-implementation/pkg/instance"
)

type ServerOptions struct {
	CrtFile   string
	KeyFile   string
	StartPort int
	EndPort   int
}

type ServerEntry struct {
	options ServerOptions

	instance map[string]*instance.ServerInstance
}

func New(options ServerOptions) *ServerEntry {
	if options.StartPort == 0 {
		options.StartPort = DefaultStartPort
	}
	if options.EndPort == 0 {
		options.EndPort = DefaultEndPort
	}
	server := &ServerEntry{
		options:  options,
		instance: make(map[string]*instance.ServerInstance),
	}
	return server
}

func (se *ServerEntry) redirectToInstance(w http.ResponseWriter, r *http.Request) {
	// conversationId
	conversationId := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	klog.V(2).Infof("conversationId: %s\n", conversationId)

	// does the server already exist, return the serverInstance
	serverInstance := se.instance[conversationId]
	if serverInstance != nil {
		klog.V(2).Infof("Server for conversationId (%s) already exists\n", serverInstance.Options.ConversationId)

		http.Redirect(w, r, serverInstance.Options.RedirectAddress, http.StatusSeeOther)
	}

	// get random port
	diff := se.options.EndPort - se.options.StartPort
	random := se.options.StartPort + rand.Intn(diff)

	// bind address
	newServer := fmt.Sprintf("%s:%d", r.URL.Host, random)
	klog.V(2).Infof("newServer: %s\n", newServer)

	// redirect address
	redirect := r.URL.Host
	if len(redirect) == 0 {
		redirect = "127.0.0.1"
	}
	newRedirect := fmt.Sprintf("https://%s:%d", redirect, random)
	klog.V(2).Infof("newRedirect: %s\n", newRedirect)

	// create server
	server := instance.New(instance.InstanceOptions{
		Port:            random,
		BindAddress:     newServer,
		RedirectAddress: newRedirect,
		ConversationId:  conversationId,
		CrtFile:         se.options.CrtFile,
		KeyFile:         se.options.KeyFile,
	})

	err := server.Start()
	if err != nil {
		klog.V(1).Infof("server.Start failed. Err: %v\n", err)
	}
	se.instance[conversationId] = server

	// small delay
	time.Sleep(time.Millisecond * 250)

	// redirect
	http.Redirect(w, r, newRedirect, http.StatusSeeOther)
}

func (se *ServerEntry) Start() {
	// redirect
	mux := http.NewServeMux()
	mux.HandleFunc("/", se.redirectToInstance)

	// this is a blocking cal
	klog.V(2).Infof("Starting server...\n")
	err := http.ListenAndServeTLS(":443", se.options.CrtFile, se.options.KeyFile, mux)
	if err != nil {
		fmt.Printf("New failed. Err: %v\n", err)
	}
}
