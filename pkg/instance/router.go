// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package instance

import (
	"fmt"
	// "github.com/davecgh/go-spew/spew"
)

type MessageRouter struct{}

func NewRouter() *MessageRouter {
	mr := &MessageRouter{}
	return mr
}

func (mr *MessageRouter) HandleMessage(byMsg []byte) error {
	fmt.Printf("---------------------------\n\n")
	// spew.Dump(topicsResult)
	// json.MarshalIndent(byMsg, "", "    ")
	fmt.Printf("%s\n", string(byMsg))
	fmt.Printf("---------------------------\n\n")

	return nil
}
