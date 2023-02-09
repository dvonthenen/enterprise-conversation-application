// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	prettyjson "github.com/hokaccha/go-prettyjson"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
	shared "github.com/dvonthenen/enterprise-reference-implementation/pkg/shared"
	utils "github.com/dvonthenen/enterprise-reference-implementation/pkg/utils"
)

func NewHandler(options MessageHandlerOptions) (*MessageHandler, error) {
	if len(options.ConversationId) == 0 {
		klog.V(1).Infof("conversationId is empty\n")
		return nil, ErrInvalidInput
	}

	mh := &MessageHandler{
		conversationId: options.ConversationId,
		callback:       options.Callback,
		options:        options,
		neo4jMgr:       options.Neo4jMgr,
		rabbitMgr:      options.RabbitMgr,
	}
	return mh, nil
}

func (mh *MessageHandler) Init() error {
	klog.V(6).Infof("MessageHandler.Init ENTER\n")

	// init all rabbit channels
	err := mh.setupRabbitChannels()
	if err != nil {
		klog.V(1).Infof("setupRabbitChannels failed. Err: %v\n", err)
		klog.V(6).Infof("MessageHandler.Init LEAVE\n")
		return err
	}

	// context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// neo4j create conversation object
	_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			createConversationQuery := utils.ReplaceIndexes(`
				MERGE (c:Conversation { #conversation_index#: $conversation_id })
					ON CREATE SET
						c.created = datetime(),
						c.lastAccessed = datetime()
					ON MATCH SET
						c.lastAccessed = datetime()
				SET c = { #conversation_index#: $conversation_id }
				`)
			result, err := tx.Run(ctx, createConversationQuery, map[string]any{
				"conversation_id": mh.conversationId,
			})
			if err != nil {
				klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
				return nil, err
			}
			return result.Collect(ctx)
		})
	if err != nil {
		klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
		klog.V(6).Infof("MessageHandler.Init LEAVE\n")
		return err
	}

	klog.V(4).Infof("MessageHandler.Init Succeeded\n")
	klog.V(6).Infof("MessageHandler.Init LEAVE\n")

	return nil
}

func (mh *MessageHandler) setupRabbitChannels() error {
	/*
		Setup Publishers...
	*/
	_, err := (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitExchangeConversationInit,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitExchangeConversationInit, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitExchangeMessage,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitExchangeMessage, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitExchangeTopic,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitExchangeTopic, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitExchangeTracker,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitExchangeTracker, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitExchangeEntity,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitExchangeEntity, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitExchangeInsight,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitExchangeInsight, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitExchangeConversationTeardown,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitExchangeConversationTeardown, err)
		return err
	}

	return nil
}

func (mh *MessageHandler) Teardown() error {
	klog.V(6).Infof("MessageHandler.Teardown ENTER\n")
	klog.V(4).Infof("MessageHandler.Teardown Succeeded\n") // This Teardown() currently is a NOOP
	klog.V(6).Infof("MessageHandler.Teardown LEAVE\n")

	return nil
}

func (mh *MessageHandler) InitializedConversation(im *sdkinterfaces.InitializationMessage) error {
	klog.V(6).Infof("InitializedConversation ENTER\n")

	// set conversation id
	im.Message.Data.ConversationID = mh.conversationId

	data, err := json.Marshal(im)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("InitializationMessage:\n%v\n\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// rabbitmq
	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitExchangeConversationInit, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}
	klog.V(3).Infof("InitializedConversation.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("InitializedConversation Succeeded\n")
	klog.V(6).Infof("InitializedConversation LEAVE\n")

	return nil
}

func (mh *MessageHandler) RecognitionResultMessage(rr *sdkinterfaces.RecognitionResult) error {
	klog.V(6).Infof("RecognitionResultMessage ENTER\n")

	prettyJson, err := prettyjson.Marshal(rr)
	if err != nil {
		klog.V(1).Infof("RecognitionResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("RecognitionResultMessage LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(5).Infof("RecognitionResultMessage:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// pass-through?
	if mh.options.TranscriptionEnabled {
		recognitionMsg := &interfaces.UserDefinedRecognition{
			Type: interfaces.MessageTypeUserDefined,
			Metadata: interfaces.Metadata{
				Type: rr.Message.Type,
			},
			Recognition: interfaces.Recognition{
				IsFinal: rr.Message.IsFinal,
				From: interfaces.From{
					ID:     rr.Message.User.ID,
					Name:   rr.Message.User.Name,
					UserID: rr.Message.User.UserID,
				},
				Content: rr.Message.Punctuated.Transcript,
			},
		}

		err := (*mh.callback).SendRecognition(recognitionMsg)
		if err != nil {
			klog.V(1).Infof("SendRecognition failed. Err: %v\n", err)
		}
	}

	klog.V(4).Infof("RecognitionResultMessage Succeeded\n")
	klog.V(6).Infof("RecognitionResultMessage LEAVE\n")

	return nil
}

func (mh *MessageHandler) MessageResponseMessage(mr *sdkinterfaces.MessageResponse) error {
	klog.V(6).Infof("MessageResponseMessage ENTER\n")

	data, err := json.Marshal(mr)
	if err != nil {
		klog.V(1).Infof("MessageResponse json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("MessageResponseMessage LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("MessageResponseMessage LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("MessageResponseMessage:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// pass-through?
	if mh.options.TranscriptionEnabled {
		recognitionMsg := &interfaces.UserDefinedMessages{
			Type: interfaces.MessageTypeUserDefined,
			Metadata: interfaces.Metadata{
				Type: mr.Type,
			},
			Messages: make([]interfaces.Message, 0),
		}

		for _, message := range mr.Messages {
			recognitionMsg.Messages = append(recognitionMsg.Messages, interfaces.Message{
				ID: message.ID,
				From: interfaces.From{
					ID:     message.From.ID,
					Name:   message.From.Name,
					UserID: message.From.UserID,
				},
				Content: message.Payload.Content,
				Duration: interfaces.Duration{
					StartTime:  message.Duration.StartTime,
					EndTime:    message.Duration.EndTime,
					TimeOffset: message.Duration.TimeOffset,
					Duration:   message.Duration.Duration,
				},
			})
		}

		err := (*mh.callback).SendMessages(recognitionMsg)
		if err != nil {
			klog.V(1).Infof("SendMessages failed. Err: %v\n", err)
		}
	}

	// write the object to the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// if we need to do something with them
	// for records, message := range mr.Messages {
	for _, message := range mr.Messages {
		_, err := (*mh.neo4jMgr).ExecuteWrite(ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				createMessageToPeopleQuery := utils.ReplaceIndexes(`
					MATCH (c:Conversation { #conversation_index#: $conversation_id })
					MERGE (m:Message { #message_index#: $message_id })
						ON CREATE SET
							m.created = datetime(),
							m.lastAccessed = datetime()
						ON MATCH SET
							m.lastAccessed = datetime()
					SET m = { #message_index#: $message_id, content: $content, startTime: $start_time, endTime: $end_time, timeOffset: $time_offset, duration: $duration, sequenceNumber: $sequence_number, raw: $raw }
					MERGE (u:User { #user_index#: $user_id })
						ON CREATE SET
							u.created = datetime(),
							u.lastAccessed = datetime()
						ON MATCH SET
							u.lastAccessed = datetime()
					SET u = { realId: $user_real_id, #user_index#: $user_id, name: $user_name, email: $user_id }
					MERGE (c)-[x:MESSAGES { #conversation_index#: $conversation_id }]-(m)
						ON CREATE SET
							x.created = datetime(),
							x.lastAccessed = datetime()
						ON MATCH SET
							x.lastAccessed = datetime()
					SET x = { #conversation_index#: $conversation_id }
					MERGE (m)-[y:SPOKE { #conversation_index#: $conversation_id }]-(u)
						ON CREATE SET
							y.created = datetime(),
							y.lastAccessed = datetime()
						ON MATCH SET
							y.lastAccessed = datetime()
					SET y = { #conversation_index#: $conversation_id }
					`)
				result, err := tx.Run(ctx, createMessageToPeopleQuery, map[string]any{
					"conversation_id": mh.conversationId,
					"message_id":      message.ID,
					"content":         message.Payload.Content,
					"start_time":      message.Duration.StartTime,
					"end_time":        message.Duration.EndTime,
					"time_offset":     message.Duration.TimeOffset,
					"duration":        message.Duration.Duration,
					"sequence_number": mr.SequenceNumber,
					"user_real_id":    message.From.ID,
					"user_name":       message.From.Name,
					"user_id":         message.From.UserID,
					"raw":             string(data),
				})
				if err != nil {
					klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
					return nil, err
				}
				return result.Collect(ctx)
			})
		if err != nil {
			klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
			klog.V(6).Infof("MessageResponseMessage LEAVE\n")
			return err
		}
	}

	// rabbitmq
	wrapperStruct := shared.MessageResponse{
		ConversationID:  mh.conversationId,
		MessageResponse: mr,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("MessageResponse json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("MessageResponseMessage LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitExchangeMessage, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}
	klog.V(3).Infof("MessageResponseMessage.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("MessageResponseMessage Succeeded\n")
	klog.V(6).Infof("MessageResponseMessage LEAVE\n")

	return nil
}

func (mh *MessageHandler) InsightResponseMessage(ir *sdkinterfaces.InsightResponse) error {
	for _, insight := range ir.Insights {
		switch insight.Type {
		case sdkinterfaces.InsightTypeQuestion:
			err := mh.HandleQuestion(&insight, ir.SequenceNumber)
			if err != nil {
				klog.V(1).Infof("HandleQuestion failed. Err: %v\n", err)
				return err
			}
		case sdkinterfaces.InsightTypeFollowUp:
			err := mh.HandleFollowUp(&insight, ir.SequenceNumber)
			if err != nil {
				klog.V(1).Infof("HandleFollowUp failed. Err: %v\n", err)
				return err
			}
		case sdkinterfaces.InsightTypeActionItem:
			err := mh.HandleActionItem(&insight, ir.SequenceNumber)
			if err != nil {
				klog.V(1).Infof("HandleActionItem failed. Err: %v\n", err)
				return err
			}
		default:
			data, err := json.Marshal(ir)
			if err != nil {
				klog.V(1).Infof("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
				return err
			}

			klog.V(1).Infof("\n\n-------------------------------\n")
			klog.V(1).Infof("Unknown InsightResponseMessage:\n\n")
			klog.V(1).Infof("Object DUMP:\n%v\n\n", string(data))
			klog.V(1).Infof("-------------------------------\n\n")
			return nil
		}
	}

	// rabbitmq
	wrapperStruct := shared.InsightResponse{
		ConversationID:  mh.conversationId,
		InsightResponse: ir,
	}

	data, err := json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("InsightResponse json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("handleInsight LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitExchangeInsight, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}
	klog.V(3).Infof("handleInsight.PublishWithContext:\n%s\n", string(data))

	return nil
}

func (mh *MessageHandler) TopicResponseMessage(tr *sdkinterfaces.TopicResponse) error {
	klog.V(6).Infof("TopicResponseMessage ENTER\n")

	data, err := json.Marshal(tr)
	if err != nil {
		klog.V(1).Infof("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TopicResponseMessage LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TopicResponseMessage LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("TopicResponseMessage:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// write the object to the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, topic := range tr.Topics {
		_, err := (*mh.neo4jMgr).ExecuteWrite(ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				createTopicsQuery := utils.ReplaceIndexes(`
					MATCH (c:Conversation { #conversation_index#: $conversation_id })
					MERGE (t:Topic { #topic_index#: $topic_id })
						ON CREATE SET
							t.created = datetime(),
							t.lastAccessed = datetime()
						ON MATCH SET
							t.lastAccessed = datetime()
					SET t = { #topic_index#: $topic_id, phrases: $phrases, score: $score, type: $type, messageIndex: $symbl_message_index, rootWords: $root_words, raw: $raw }
					MERGE (c)-[x:TOPICS { #conversation_index#: $conversation_id }]-(t)
						ON CREATE SET
							x.created = datetime(),
							x.lastAccessed = datetime()
						ON MATCH SET
							x.lastAccessed = datetime()
					SET x = { #conversation_index#: $conversation_id }
					`)
				result, err := tx.Run(ctx, createTopicsQuery, map[string]any{
					"conversation_id":     mh.conversationId,
					"topic_id":            topic.ID,
					"phrases":             topic.Phrases,
					"score":               topic.Score,
					"type":                topic.Type,
					"symbl_message_index": topic.MessageIndex,
					"root_words":          convertRootWordToString(topic.RootWords),
					"raw":                 string(data),
				})
				if err != nil {
					klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
					return nil, err
				}
				return result.Collect(ctx)
			})
		if err != nil {
			klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
			klog.V(6).Infof("TopicResponseMessage LEAVE\n")
			return err
		}

		// associate topic to message
		for _, ref := range topic.MessageReferences {
			_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
				func(tx neo4j.ManagedTransaction) (any, error) {
					createTopicsQuery := utils.ReplaceIndexes(`
						MATCH (t:Topic { topicId: $topic_id })
						MATCH (m:Message { #message_index#: $message_id })
						MERGE (t)-[x:TOPIC_MESSAGE_REF { #conversation_index#: $conversation_id }]-(m)
							ON CREATE SET
								x.created = datetime(),
								x.lastAccessed = datetime()
							ON MATCH SET
								x.lastAccessed = datetime()
						SET x = { #conversation_index#: $conversation_id, value: $value }
						`)
					result, err := tx.Run(ctx, createTopicsQuery, map[string]any{
						"conversation_id": mh.conversationId,
						"topic_id":        topic.ID,
						"message_id":      ref.ID,
						"value":           topic.Phrases,
					})
					if err != nil {
						klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
						return nil, err
					}
					return result.Collect(ctx)
				})
			if err != nil {
				klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
				klog.V(6).Infof("TopicResponseMessage LEAVE\n")
				return err
			}
		}
	}

	// rabbitmq
	wrapperStruct := shared.TopicResponse{
		ConversationID: mh.conversationId,
		TopicResponse:  tr,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("TopicResponse json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TopicResponseMessage LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitExchangeTopic, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}
	klog.V(3).Infof("TopicResponseMessage.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("TopicResponseMessage Succeeded\n")
	klog.V(6).Infof("TopicResponseMessage LEAVE\n")

	return nil
}
func (mh *MessageHandler) TrackerResponseMessage(tr *sdkinterfaces.TrackerResponse) error {
	klog.V(6).Infof("TrackerResponseMessage ENTER\n")

	data, err := json.Marshal(tr)
	if err != nil {
		klog.V(1).Infof("TrackerResponseMessage json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TrackerResponseMessage LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TrackerResponseMessage LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("TrackerResponseMessage:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// write the object to the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, tracker := range tr.Trackers {
		_, err := (*mh.neo4jMgr).ExecuteWrite(ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				createTrackersQuery := utils.ReplaceIndexes(`
					MATCH (c:Conversation { #conversation_index#: $conversation_id })
					MERGE (t:Tracker { #tracker_index#: $tracker_id })
						ON CREATE SET
							t.created = datetime(),
							t.lastAccessed = datetime()
						ON MATCH SET
							t.lastAccessed = datetime()
					SET t = { #tracker_index#: $tracker_id, name: $tracker_name, raw: $raw }
					MERGE (c)-[x:TRACKER { #conversation_index#: $conversation_id }]-(t)
						ON CREATE SET
							x.created = datetime(),
							x.lastAccessed = datetime()
						ON MATCH SET
							x.lastAccessed = datetime()
					SET x = { #conversation_index#: $conversation_id }
					`)
				result, err := tx.Run(ctx, createTrackersQuery, map[string]any{
					"conversation_id": mh.conversationId,
					"tracker_id":      tracker.ID,
					"tracker_name":    tracker.Name,
					"raw":             string(data),
				})
				if err != nil {
					klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
					return nil, err
				}
				return result.Collect(ctx)
			})
		if err != nil {
			klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
			klog.V(6).Infof("TrackerResponseMessage LEAVE\n")
			return err
		}

		// associate tracker to messages and insights
		for _, match := range tracker.Matches {

			// messages
			for _, msgRef := range match.MessageRefs {
				_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						createTopicsQuery := utils.ReplaceIndexes(`
							MATCH (t:Tracker { #tracker_index#: $tracker_id })
							MATCH (m:Message { #message_index#: $message_id })
							MERGE (t)-[x:TRACKER_MESSAGE_REF { #conversation_index#: $conversation_id }]-(m)
								ON CREATE SET
									x.created = datetime(),
									x.lastAccessed = datetime()
								ON MATCH SET
									x.lastAccessed = datetime()
							SET x = { #conversation_index#: $conversation_id, value: $value }
							`)
						result, err := tx.Run(ctx, createTopicsQuery, map[string]any{
							"conversation_id": mh.conversationId,
							"tracker_id":      tracker.ID,
							"message_id":      msgRef.ID,
							"value":           match.Value,
						})
						if err != nil {
							klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
							return nil, err
						}
						return result.Collect(ctx)
					})
				if err != nil {
					klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
					klog.V(6).Infof("TrackerResponseMessage LEAVE\n")
					return err
				}
			}

			// insights
			for _, inRef := range match.InsightRefs {
				_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						createTrackerMatchQuery := utils.ReplaceIndexes(`
							MATCH (t:Tracker { #tracker_index#: $tracker_id })
							MATCH (i:Insight { #insight_index#: $insight_id })
							MERGE (t)-[x:TRACKER_INSIGHT_REF { #conversation_index#: $conversation_id }]-(i)
								ON CREATE SET
									x.created = datetime(),
									x.lastAccessed = datetime()
								ON MATCH SET
									x.lastAccessed = datetime()
							SET x = { #conversation_index#: $conversation_id, value: $value }
							`)
						result, err := tx.Run(ctx, createTrackerMatchQuery, map[string]any{
							"conversation_id": mh.conversationId,
							"tracker_id":      tracker.ID,
							"insight_id":      inRef.ID,
							"value":           match.Value,
						})
						if err != nil {
							klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
							return nil, err
						}
						return result.Collect(ctx)
					})
				if err != nil {
					klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
					klog.V(6).Infof("TrackerResponseMessage LEAVE\n")
					return err
				}
			}
		}
	}

	// rabbitmq
	wrapperStruct := shared.TrackerResponse{
		ConversationID:  mh.conversationId,
		TrackerResponse: tr,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("TrackerResponse json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TrackerResponseMessage LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitExchangeTracker, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}
	klog.V(3).Infof("TrackerResponseMessage.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("TrackerResponseMessage Succeeded\n")
	klog.V(6).Infof("TrackerResponseMessage LEAVE\n")

	return nil
}

func (mh *MessageHandler) EntityResponseMessage(er *sdkinterfaces.EntityResponse) error {
	klog.V(6).Infof("EntityResponseMessage ENTER\n")

	data, err := json.Marshal(er)
	if err != nil {
		klog.V(1).Infof("EntityResponseMessage json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("EntityResponseMessage LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("EntityResponseMessage LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("EntityResponseMessage:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// write the object to the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, entity := range er.Entities {

		// associate tracker to messages and insights
		for _, match := range entity.Matches {

			// entity id
			entityCategory := strings.ReplaceAll(entity.Category, " ", "_")
			entityType := strings.ReplaceAll(entity.Type, " ", "_")
			entitySubType := strings.ReplaceAll(entity.SubType, " ", "_")
			entityValue := strings.ReplaceAll(match.DetectedValue, " ", "_")
			entityId := fmt.Sprintf("%s/%s/%s/%s", entityCategory, entityType, entitySubType, entityValue)

			// entity
			_, err := (*mh.neo4jMgr).ExecuteWrite(ctx,
				func(tx neo4j.ManagedTransaction) (any, error) {
					createEntitiesQuery := utils.ReplaceIndexes(`
						MATCH (c:Conversation { #conversation_index#: $conversation_id })
						MERGE (e:Entity { #entity_index#: $entity_id })
							ON CREATE SET
								e.created = datetime(),
								e.lastAccessed = datetime()
							ON MATCH SET
								e.lastAccessed = datetime()
						SET e = { #entity_index#: $entity_id, type: $type, subType: $sub_type, category: $category, value: $value, raw: $raw }
						MERGE (c)-[x:ENTITY { #conversation_index#: $conversation_id }]-(e)
							ON CREATE SET
								x.created = datetime(),
								x.lastAccessed = datetime()
							ON MATCH SET
								x.lastAccessed = datetime()
						SET x = { #conversation_index#: $conversation_id }
						`)
					result, err := tx.Run(ctx, createEntitiesQuery, map[string]any{
						"conversation_id": mh.conversationId,
						"entity_id":       entityId,
						"type":            entity.Type,
						"sub_type":        entity.SubType,
						"category":        entity.Category,
						"value":           match.DetectedValue,
						"raw":             string(data),
					})
					if err != nil {
						klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
						return nil, err
					}
					return result.Collect(ctx)
				})
			if err != nil {
				klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
				klog.V(6).Infof("EntityResponseMessage LEAVE\n")
				return err
			}

			// message
			for _, msgRef := range match.MessageRefs {
				_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						createEntitiesQuery := utils.ReplaceIndexes(`
							MATCH (e:Entity { #entity_index#: $entity_id })
							MATCH (m:Message { #message_index#: $message_id })
							MERGE (e)-[x:ENTITY_MESSAGE_REF { #conversation_index#: $conversation_id }]-(m)
								ON CREATE SET
									x.created = datetime(),
									x.lastAccessed = datetime()
								ON MATCH SET
									x.lastAccessed = datetime()
							SET x = { #conversation_index#: $conversation_id, value: $value }
							`)
						result, err := tx.Run(ctx, createEntitiesQuery, map[string]any{
							"conversation_id": mh.conversationId,
							"entity_id":       entityId,
							"message_id":      msgRef.ID,
							"value":           match.DetectedValue,
						})
						if err != nil {
							klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
							return nil, err
						}
						return result.Collect(ctx)
					})
				if err != nil {
					klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
					klog.V(6).Infof("EntityResponseMessage LEAVE\n")
					return err
				}
			}
		}
	}

	// rabbitmq
	wrapperStruct := shared.EntityResponse{
		ConversationID: mh.conversationId,
		EntityResponse: er,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("EntityResponse json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("EntityResponseMessage LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitExchangeEntity, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}
	klog.V(3).Infof("EntityResponseMessage.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("EntityResponseMessage Succeeded\n")
	klog.V(6).Infof("EntityResponseMessage LEAVE\n")

	return nil
}

func (mh *MessageHandler) TeardownConversation(tm *sdkinterfaces.TeardownMessage) error {
	klog.V(6).Infof("TeardownConversation ENTER\n")

	if mh.terminationSent {
		klog.V(1).Infof("TeardownConversation already handled\n")
		klog.V(6).Infof("TeardownConversation LEAVE\n")
		return nil
	}

	data, err := json.Marshal(tm)
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TeardownConversation LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TeardownConversation LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("TeardownConversation:\n%v\n\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// rabbitmq
	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitExchangeConversationTeardown, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}
	klog.V(3).Infof("TeardownConversation.PublishWithContext:\n%s\n", string(data))

	// mark as teardown message sent
	mh.terminationSent = true

	klog.V(4).Infof("TeardownConversation Succeeded\n")
	klog.V(6).Infof("TeardownConversation LEAVE\n")

	return nil
}

// Not used
func (mh *MessageHandler) UnhandledMessage(byMsg []byte) error {
	klog.V(1).Infof("\n\n-------------------------------\n")
	klog.V(1).Infof("UnhandledMessage:\n%v\n", string(byMsg))
	klog.V(1).Infof("-------------------------------\n\n")
	return nil
}

// Not used
func (mh *MessageHandler) UserDefinedMessage(byMsg []byte) error {
	klog.V(1).Infof("\n\n-------------------------------\n")
	klog.V(1).Infof("UserDefinedMessage:\n%v\n", string(byMsg))
	klog.V(1).Infof("-------------------------------\n\n")
	return nil
}

func (mh *MessageHandler) HandleQuestion(insight *sdkinterfaces.Insight, number int) error {
	return mh.handleInsight(insight, number)
}

func (mh *MessageHandler) HandleActionItem(insight *sdkinterfaces.Insight, number int) error {
	return mh.handleInsight(insight, number)
}

func (mh *MessageHandler) HandleFollowUp(insight *sdkinterfaces.Insight, number int) error {
	return mh.handleInsight(insight, number)
}

func (mh *MessageHandler) handleInsight(insight *sdkinterfaces.Insight, squenceNumber int) error {
	klog.V(6).Infof("handleInsight ENTER\n")

	data, err := json.Marshal(insight)
	if err != nil {
		klog.V(1).Infof("handleInsight json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("handleInsight LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("handleInsight LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("handleInsight:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// write the object to the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// if we need to do something with them
	// for records, message := range mr.Messages {
	_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			createInsightQuery := utils.ReplaceIndexes(`
				MATCH (c:Conversation { #conversation_index#: $conversation_id })
				MERGE (i:Insight { #insight_index#: $insight_id })
					ON CREATE SET
						i.created = datetime(),
						i.lastAccessed = datetime()
					ON MATCH SET
						i.lastAccessed = datetime()
				SET i = { #insight_index#: $insight_id, type: $type, content: $content, sequenceNumber: $sequence_number, assigneeId: $assignee_id, raw: $raw }
				MERGE (u:User { #user_index#: $user_id })
					ON CREATE SET
						u.created = datetime(),
						u.lastAccessed = datetime()
					ON MATCH SET
						u.lastAccessed = datetime()
				SET u = { realId: $user_real_id, #user_index#: $user_id, name: $user_name, email: $user_id }
				MERGE (c)-[x:INSIGHT { #conversation_index#: $conversation_id }]-(i)
					ON CREATE SET
						x.created = datetime(),
						x.lastAccessed = datetime()
					ON MATCH SET
						x.lastAccessed = datetime()
				SET x = { #conversation_index#: $conversation_id }
				MERGE (i)-[y:SPOKE { #conversation_index#: $conversation_id }]-(u)
					ON CREATE SET
						y.created = datetime(),
						y.lastAccessed = datetime()
					ON MATCH SET
						y.lastAccessed = datetime()
				SET y = { #conversation_index#: $conversation_id }
				`)
			result, err := tx.Run(ctx, createInsightQuery, map[string]any{
				"conversation_id": mh.conversationId,
				"insight_id":      insight.ID,
				"type":            insight.Type,
				"content":         insight.Payload.Content,
				"sequence_number": squenceNumber,
				"assignee_id":     insight.Assignee.UserID,
				"user_real_id":    insight.From.ID,
				"user_id":         insight.From.UserID,
				"user_name":       insight.From.Name,
				"raw":             string(data),
			})
			if err != nil {
				klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
				return nil, err
			}
			return result.Collect(ctx)
		})
	if err != nil {
		klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
		klog.V(6).Infof("handleInsight LEAVE\n")
		return err
	}

	klog.V(4).Infof("handleInsight Succeeded\n")
	klog.V(6).Infof("handleInsight LEAVE\n")

	return nil
}

func convertRootWordToString(words []sdkinterfaces.RootWord) string {
	tmp := ""
	for _, word := range words {
		if len(tmp) > 0 {
			tmp = "," + tmp
		}
		tmp = tmp + word.Text
	}
	return tmp
}
