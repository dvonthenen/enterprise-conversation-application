// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
)

type MyMessageRouter struct {
	TranscriptionDemo    bool
	TranscriptionDisable bool

	ChatmessageDemo    bool
	ChatmessageDisable bool
}

func NewMyMessageRouter() *MyMessageRouter {
	var transcriptionDemoStr string
	if v := os.Getenv("SYMBL_TRANSCRIPTION_DEMO"); v != "" {
		klog.V(4).Info("SYMBL_TRANSCRIPTION_DEMO found")
		transcriptionDemoStr = v
	}
	var transcriptionDisableStr string
	if v := os.Getenv("SYMBL_TRANSCRIPTION_DISABLE"); v != "" {
		klog.V(4).Info("SYMBL_TRANSCRIPTION_DISABLE found")
		transcriptionDisableStr = v
	}
	var chatmessageDemoStr string
	if v := os.Getenv("SYMBL_CHAT_MESSAGE_DEMO"); v != "" {
		klog.V(4).Info("SYMBL_CHAT_MESSAGE_DEMO found")
		chatmessageDemoStr = v
	}
	var chatmessageDisableStr string
	if v := os.Getenv("SYMBL_CHAT_MESSAGE_DISABLE"); v != "" {
		klog.V(4).Info("SYMBL_CHAT_MESSAGE_DISABLE found")
		chatmessageDisableStr = v
	}

	transcriptionDemo := strings.EqualFold(strings.ToLower(transcriptionDemoStr), "true")
	transcriptionDisable := strings.EqualFold(strings.ToLower(transcriptionDisableStr), "true")
	chatmessageDemo := strings.EqualFold(strings.ToLower(chatmessageDemoStr), "true")
	chatmessageDisable := strings.EqualFold(strings.ToLower(chatmessageDisableStr), "true")

	return &MyMessageRouter{
		TranscriptionDemo:    transcriptionDemo,
		TranscriptionDisable: transcriptionDisable,
		ChatmessageDemo:      chatmessageDemo,
		ChatmessageDisable:   chatmessageDisable,
	}
}

func (mmr *MyMessageRouter) InitializedConversation(im *sdkinterfaces.InitializationMessage) error {
	data, err := json.Marshal(im)
	if err != nil {
		klog.V(1).Infof("InitializationMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nInitializationMessage Object DUMP:\n%s\n\n", prettyJson)
	return nil
}

func (mmr *MyMessageRouter) RecognitionResultMessage(rr *sdkinterfaces.RecognitionResult) error {
	if mmr.TranscriptionDisable {
		return nil // disable all output
	}

	if mmr.TranscriptionDemo {
		klog.Infof("TRANSCRIPTION: %s\n", rr.Message.Punctuated.Transcript)
		return nil
	}

	data, err := json.Marshal(rr)
	if err != nil {
		klog.V(1).Infof("RecognitionResult json.Marshal failed. Err: %v\n", err)
		return err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nRecognitionResult Object DUMP:\n%s\n\n", prettyJson)

	return nil
}

func (mmr *MyMessageRouter) MessageResponseMessage(mr *sdkinterfaces.MessageResponse) error {
	if mmr.ChatmessageDisable {
		return nil // disable chat output
	}

	if mmr.ChatmessageDemo {
		for _, msg := range mr.Messages {
			klog.Infof("\n\nChat Message [%s]: %s\n\n", msg.From.Name, msg.Payload.Content)
		}
		return nil
	}

	data, err := json.Marshal(mr)
	if err != nil {
		klog.V(1).Infof("MessageResponse json.Marshal failed. Err: %v\n", err)
		return err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nMessageResponse Object DUMP:\n%s\n\n", prettyJson)

	return nil
}

func (mmr *MyMessageRouter) InsightResponseMessage(ir *sdkinterfaces.InsightResponse) error {
	data, err := json.Marshal(ir)
	if err != nil {
		klog.V(1).Infof("InsightResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nInsightResponseMessage Object DUMP:\n%s\n\n", prettyJson)
	return nil
}

func (mmr *MyMessageRouter) TopicResponseMessage(tr *sdkinterfaces.TopicResponse) error {
	data, err := json.Marshal(tr)
	if err != nil {
		klog.V(1).Infof("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nTopicResponseMessage Object DUMP:\n%s\n\n", prettyJson)
	return nil
}
func (mmr *MyMessageRouter) TrackerResponseMessage(tr *sdkinterfaces.TrackerResponse) error {
	data, err := json.Marshal(tr)
	if err != nil {
		klog.V(1).Infof("TrackerResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nTrackerResponseMessage Object DUMP:\n%s\n\n", prettyJson)
	return nil
}

func (mmr *MyMessageRouter) EntityResponseMessage(tr *sdkinterfaces.EntityResponse) error {
	data, err := json.Marshal(tr)
	if err != nil {
		klog.V(1).Infof("EntityResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nEntityResponseMessage Object DUMP:\n%s\n\n", prettyJson)
	return nil
}

func (mmr *MyMessageRouter) TeardownConversation(tm *sdkinterfaces.TeardownMessage) error {
	data, err := json.Marshal(tm)
	if err != nil {
		klog.V(1).Infof("TeardownConversation json.Marshal failed. Err: %v\n", err)
		return err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nTeardownConversation Object DUMP:\n%s\n\n", prettyJson)
	return nil
}

func (mmr *MyMessageRouter) UserDefinedMessage(byMsg []byte) error {
	// reform struct
	var ast interfaces.AppSpecificType
	err := json.Unmarshal(byMsg, &ast)
	if err != nil {
		klog.V(1).Infof("[UserDefinedMessage] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	switch ast.Metadata.Type {
	case interfaces.MessageTypeRecognitionResult:
		return mmr.processUserDefinedRecognition(byMsg)
	case interfaces.MessageTypeMessageResponse:
		return mmr.processUserDefinedMessage(byMsg)
	}

	prettyJson, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nUserDefinedMessage Object DUMP:\n%s\n\n", prettyJson)
	return nil
}

func (mmr *MyMessageRouter) UnhandledMessage(byMsg []byte) error {
	prettyJson, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nUnhandledMessage Object DUMP:\n%s\n\n", prettyJson)
	return nil
}

func (mmr *MyMessageRouter) processUserDefinedRecognition(byMsg []byte) error {
	if mmr.TranscriptionDisable {
		return nil // disable all output
	}

	// reform struct
	var udr interfaces.UserDefinedRecognition
	err := json.Unmarshal(byMsg, &udr)
	if err != nil {
		klog.V(1).Infof("[processUserDefinedRecognition] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	if mmr.TranscriptionDemo {
		fmt.Printf("TRANSCRIPTION [%s]: %s\n", udr.Recognition.From.Name, udr.Recognition.Content)
		return nil
	}

	prettyJson, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nUserDefinedMessage Object DUMP:\n%s\n\n", prettyJson)

	return nil
}

func (mmr *MyMessageRouter) processUserDefinedMessage(byMsg []byte) error {
	if mmr.ChatmessageDisable {
		return nil // disable chat output
	}

	// reform struct
	var udm interfaces.UserDefinedMessages
	err := json.Unmarshal(byMsg, &udm)
	if err != nil {
		klog.V(1).Infof("[processUserDefinedMessage] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	if mmr.ChatmessageDemo {
		for _, msg := range udm.Messages {
			fmt.Printf("\nChat Message [%s]: %s\n\n", msg.From.Name, msg.Content)
		}
		return nil
	}

	prettyJson, err := prettyjson.Format(byMsg)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}

	klog.Infof("\n\nUserDefinedMessage Object DUMP:\n%s\n\n", prettyJson)

	return nil
}
