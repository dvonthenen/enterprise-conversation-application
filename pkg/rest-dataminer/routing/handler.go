// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	rabbitinterfaces "github.com/dvonthenen/rabbitmq-manager/pkg/interfaces"
	asyncinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/async/v1/interfaces"
	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/async/v1/interfaces"
	prettyjson "github.com/hokaccha/go-prettyjson"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"

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

	klog.V(4).Infof("MessageHandler.Init Succeeded\n")
	klog.V(6).Infof("MessageHandler.Init LEAVE\n")

	return nil
}

func (mh *MessageHandler) setupRabbitChannels() error {
	/*
		Setup Publishers...
	*/
	_, err := (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitAsyncConversationInit,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitAsyncConversationInit, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitAsyncMessage,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitAsyncMessage, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitAsyncQuestion,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitAsyncQuestion, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitAsyncFollowUp,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitAsyncFollowUp, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitAsyncActionItem,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitAsyncActionItem, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitAsyncTopic,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitAsyncTopic, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitAsyncTracker,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitAsyncTracker, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitAsyncEntity,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitAsyncEntity, err)
		return err
	}
	_, err = (*mh.rabbitMgr).CreatePublisher(rabbitinterfaces.PublisherOptions{
		Name:        shared.RabbitAsyncConversationTeardown,
		Type:        rabbitinterfaces.ExchangeTypeFanout,
		AutoDeleted: true,
		IfUnused:    true,
	})
	if err != nil {
		klog.V(1).Infof("CreatePublisher %s failed. Err: %v\n", shared.RabbitAsyncConversationTeardown, err)
		return err
	}

	return nil
}

func (mh *MessageHandler) DoesConversationExist() (bool, error) {
	if mh.existsSet {
		return mh.conversationExists, nil
	}

	ctx := context.Background()

	var retValue int64

	_, err := (*mh.neo4jMgr).ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		myQuery := utils.ReplaceIndexes(`
			MATCH (c:Conversation)
			WHERE c.#conversation_index# = $conversation_id
			RETURN count(c)`)
		result, err := tx.Run(ctx, myQuery, map[string]any{
			"conversation_id": mh.conversationId,
		})
		if err != nil {
			return false, err
		}

		for result.Next(ctx) {
			retValue = result.Record().Values[0].(int64)
		}

		return nil, result.Err()
	})
	if err != nil {
		klog.V(1).Infof("ExecuteRead failed. Err: %v\n", err)
		return false, err
	}

	mh.conversationExists = (retValue > 0)
	mh.existsSet = true

	return mh.conversationExists, nil
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
	im.ConversationID = mh.conversationId

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

	// insights already created
	exists, err := mh.DoesConversationExist()
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}

	// only add to DB if new
	if !exists {
		// context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// neo4j create conversation object
		_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				createConversationQuery := utils.ReplaceIndexes(`
					MERGE (c:Conversation { #conversation_index#: $conversation_id })
						ON CREATE SET
							c.createdAt = datetime(),
							c.lastAccessed = datetime()
						ON MATCH SET
							c.lastAccessed = datetime()
					SET c = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime() }
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
	}

	// rabbitmq
	wrapperStruct := shared.InitializationResult{
		Duplicate: exists,
		InitializationMessage: &asyncinterfaces.InitializationMessage{
			ConversationID: mh.conversationId,
		},
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("InitializedConversation json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("InitializedConversation LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitAsyncConversationInit, data)
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

func (mh *MessageHandler) MessageResult(mr *sdkinterfaces.MessageResult) error {
	klog.V(6).Infof("MessageResult ENTER\n")

	data, err := json.Marshal(mr)
	if err != nil {
		klog.V(1).Infof("MessageResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("MessageResult LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("MessageResult LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("MessageResult:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// any insights?
	if len(mr.Messages) == 0 {
		klog.V(4).Infof("mr.Messages is empty. Nothing to process.\n")
		klog.V(6).Infof("MessageResult LEAVE\n")
		return nil
	}

	// insights already created
	exists, err := mh.DoesConversationExist()
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("MessageResult LEAVE\n")
		return err
	}

	// only add to DB if new
	if !exists {
		// write the object to the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// process messages
		for _, message := range mr.Messages {
			_, err := (*mh.neo4jMgr).ExecuteWrite(ctx,
				func(tx neo4j.ManagedTransaction) (any, error) {
					createMessageToPeopleQuery := utils.ReplaceIndexes(`
						MATCH (c:Conversation { #conversation_index#: $conversation_id })
						MERGE (m:Message { #message_index#: $message_id })
							ON CREATE SET
								m.createdAt = datetime(),
								m.lastAccessed = datetime()
							ON MATCH SET
								m.lastAccessed = datetime()
						SET m = { #message_index#: $message_id, content: $content, startTime: $start_time, endTime: $end_time, timeOffset: $time_offset, duration: $duration, sequenceNumber: $sequence_number, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						MERGE (u:User { #user_index#: $user_id })
							ON CREATE SET
								u.createdAt = datetime(),
								u.lastAccessed = datetime()
							ON MATCH SET
								u.lastAccessed = datetime()
						SET u = { realId: $user_real_id, #user_index#: $user_id, name: $user_name, email: $user_id, createdAt: datetime(), lastAccessed: datetime() }
						MERGE (c)-[x:MESSAGES { #conversation_index#: $conversation_id }]-(m)
							ON CREATE SET
								x.createdAt = datetime(),
								x.lastAccessed = datetime()
							ON MATCH SET
								x.lastAccessed = datetime()
						SET x = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						MERGE (m)-[y:SPOKE { #conversation_index#: $conversation_id }]-(u)
							ON CREATE SET
								y.createdAt = datetime(),
								y.lastAccessed = datetime()
							ON MATCH SET
								y.lastAccessed = datetime()
						SET y = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						`)
					result, err := tx.Run(ctx, createMessageToPeopleQuery, map[string]any{
						"conversation_id": mh.conversationId,
						"message_id":      message.ID,
						"content":         message.Text,
						"start_time":      message.StartTime,
						"end_time":        message.EndTime,
						"time_offset":     message.TimeOffset,
						"duration":        message.Duration,
						"sequence_number": "TODO", // TODO
						"user_real_id":    message.From.ID,
						"user_name":       message.From.Name,
						"user_id":         message.From.ID, // TODO: Look into it
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
				klog.V(6).Infof("MessageResult LEAVE\n")
				return err
			}
		}
	}

	// rabbitmq
	wrapperStruct := shared.MessageResult{
		Duplicate:      exists,
		ConversationID: mh.conversationId,
		MessageResult:  mr,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("MessageResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("MessageResult LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitAsyncMessage, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("MessageResult LEAVE\n")
		return err
	}
	klog.V(3).Infof("MessageResult.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("MessageResult Succeeded\n")
	klog.V(6).Infof("MessageResult LEAVE\n")

	return nil
}

func (mh *MessageHandler) QuestionResult(qr *sdkinterfaces.QuestionResult) error {
	klog.V(6).Infof("QuestionResult ENTER\n")

	data, err := json.Marshal(qr)
	if err != nil {
		klog.V(1).Infof("QuestionResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("QuestionResult LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("QuestionResult LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("QuestionResult:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// any insights?
	if len(qr.Questions) == 0 {
		klog.V(4).Infof("qr.Questions is empty. Nothing to process.\n")
		klog.V(6).Infof("QuestionResult LEAVE\n")
		return nil
	}

	// insights already created
	exists, err := mh.DoesConversationExist()
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("QuestionResult LEAVE\n")
		return err
	}

	// only add to DB if new
	if !exists {
		// write the object to the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// if we need to do something with them
		cnt := 0

		for _, question := range qr.Questions {
			_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
				func(tx neo4j.ManagedTransaction) (any, error) {
					createInsightQuery := utils.ReplaceIndexes(`
						MATCH (c:Conversation { #conversation_index#: $conversation_id })
						MERGE (i:Insight { #insight_index#: $insight_id })
							ON CREATE SET
								i.createdAt = datetime(),
								i.lastAccessed = datetime()
							ON MATCH SET
								i.lastAccessed = datetime()
						SET i = { #insight_index#: $insight_id, type: $type, content: $content, sequenceNumber: $sequence_number, assigneeId: $assignee_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						MERGE (u:User { #user_index#: $user_id })
							ON CREATE SET
								u.createdAt = datetime(),
								u.lastAccessed = datetime()
							ON MATCH SET
								u.lastAccessed = datetime()
						SET u = { realId: $user_real_id, #user_index#: $user_id, name: $user_name, email: $user_id, createdAt: datetime(), lastAccessed: datetime() }
						MERGE (c)-[x:INSIGHT { #conversation_index#: $conversation_id }]-(i)
							ON CREATE SET
								x.createdAt = datetime(),
								x.lastAccessed = datetime()
							ON MATCH SET
								x.lastAccessed = datetime()
						SET x = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						MERGE (i)-[y:SPOKE { #conversation_index#: $conversation_id }]-(u)
							ON CREATE SET
								y.createdAt = datetime(),
								y.lastAccessed = datetime()
							ON MATCH SET
								y.lastAccessed = datetime()
						SET y = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime() }
						`)
					result, err := tx.Run(ctx, createInsightQuery, map[string]any{
						"conversation_id": mh.conversationId,
						"insight_id":      question.ID,
						"type":            strings.ToLower(question.Type),
						"content":         question.Text,
						"sequence_number": cnt,
						"assignee_id":     "TODO", // TODO
						"user_real_id":    question.From.ID,
						"user_id":         question.From.ID, // TODO: Look into it
						"user_name":       question.From.Name,
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
				klog.V(6).Infof("QuestionResult LEAVE\n")
				return err
			}

			cnt++
		}
	}

	// rabbitmq
	wrapperStruct := shared.QuestionResult{
		Duplicate:      exists,
		ConversationID: mh.conversationId,
		QuestionResult: qr,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("QuestionResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("QuestionResult LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitAsyncQuestion, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("QuestionResult LEAVE\n")
		return err
	}
	klog.V(3).Infof("QuestionResult.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("QuestionResult Succeeded\n")
	klog.V(6).Infof("QuestionResult LEAVE\n")

	return nil
}

func (mh *MessageHandler) FollowUpResult(fur *sdkinterfaces.FollowUpResult) error {
	klog.V(6).Infof("FollowUpResult ENTER\n")

	data, err := json.Marshal(fur)
	if err != nil {
		klog.V(1).Infof("FollowUpResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("FollowUpResult LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("FollowUpResult LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("FollowUpResult:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// any insights?
	if len(fur.FollowUps) == 0 {
		klog.V(4).Infof("fur.FollowUps is empty. Nothing to process.\n")
		klog.V(6).Infof("FollowUpResult LEAVE\n")
		return nil
	}

	// insights already created
	exists, err := mh.DoesConversationExist()
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("FollowUpResult LEAVE\n")
		return err
	}

	// only add to DB if new
	if !exists {
		// write the object to the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// if we need to do something with them
		cnt := 0

		for _, followUps := range fur.FollowUps {
			_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
				func(tx neo4j.ManagedTransaction) (any, error) {
					createInsightQuery := utils.ReplaceIndexes(`
						MATCH (c:Conversation { #conversation_index#: $conversation_id })
						MERGE (i:Insight { #insight_index#: $insight_id })
							ON CREATE SET
								i.createdAt = datetime(),
								i.lastAccessed = datetime()
							ON MATCH SET
								i.lastAccessed = datetime()
						SET i = { #insight_index#: $insight_id, type: $type, content: $content, sequenceNumber: $sequence_number, assigneeId: $assignee_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						MERGE (u:User { #user_index#: $user_id })
							ON CREATE SET
								u.createdAt = datetime(),
								u.lastAccessed = datetime()
							ON MATCH SET
								u.lastAccessed = datetime()
						SET u = { realId: $user_real_id, #user_index#: $user_id, name: $user_name, email: $user_id, createdAt: datetime(), lastAccessed: datetime() }
						MERGE (c)-[x:INSIGHT { #conversation_index#: $conversation_id }]-(i)
							ON CREATE SET
								x.createdAt = datetime(),
								x.lastAccessed = datetime()
							ON MATCH SET
								x.lastAccessed = datetime()
						SET x = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						MERGE (i)-[y:SPOKE { #conversation_index#: $conversation_id }]-(u)
							ON CREATE SET
								y.createdAt = datetime(),
								y.lastAccessed = datetime()
							ON MATCH SET
								y.lastAccessed = datetime()
						SET y = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime() }
						`)
					result, err := tx.Run(ctx, createInsightQuery, map[string]any{
						"conversation_id": mh.conversationId,
						"insight_id":      followUps.ID,
						"type":            strings.ToLower(followUps.Type),
						"content":         followUps.Text,
						"sequence_number": cnt,
						"assignee_id":     followUps.Assignee.ID, // TODO: Look into it ID or Name
						"user_real_id":    followUps.From.ID,
						"user_id":         followUps.From.ID, // TODO: Look into it
						"user_name":       followUps.From.Name,
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
				klog.V(6).Infof("FollowUpResult LEAVE\n")
				return err
			}

			cnt++
		}
	}

	// rabbitmq
	wrapperStruct := shared.FollowUpResult{
		Duplicate:      exists,
		ConversationID: mh.conversationId,
		FollowUpResult: fur,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("FollowUpResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("FollowUpResult LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitAsyncFollowUp, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("FollowUpResult LEAVE\n")
		return err
	}
	klog.V(3).Infof("FollowUpResult.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("FollowUpResult Succeeded\n")
	klog.V(6).Infof("FollowUpResult LEAVE\n")

	return nil
}

func (mh *MessageHandler) ActionItemResult(air *sdkinterfaces.ActionItemResult) error {
	klog.V(6).Infof("ActionItemResult ENTER\n")

	data, err := json.Marshal(air)
	if err != nil {
		klog.V(1).Infof("ActionItemResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("ActionItemResult LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("ActionItemResult LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("ActionItemResult:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// any insights?
	if len(air.ActionItems) == 0 {
		klog.V(4).Infof("air.ActionItems is empty. Nothing to process.\n")
		klog.V(6).Infof("ActionItemResult LEAVE\n")
		return nil
	}

	// insights already created
	exists, err := mh.DoesConversationExist()
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("ActionItemResult LEAVE\n")
		return err
	}

	// only add to DB if new
	if !exists {
		// write the object to the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// if we need to do something with them
		cnt := 0
		for _, actionItem := range air.ActionItems {
			_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
				func(tx neo4j.ManagedTransaction) (any, error) {
					createInsightQuery := utils.ReplaceIndexes(`
						MATCH (c:Conversation { #conversation_index#: $conversation_id })
						MERGE (i:Insight { #insight_index#: $insight_id })
							ON CREATE SET
								i.createdAt = datetime(),
								i.lastAccessed = datetime()
							ON MATCH SET
								i.lastAccessed = datetime()
						SET i = { #insight_index#: $insight_id, type: $type, content: $content, sequenceNumber: $sequence_number, assigneeId: $assignee_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						MERGE (u:User { #user_index#: $user_id })
							ON CREATE SET
								u.createdAt = datetime(),
								u.lastAccessed = datetime()
							ON MATCH SET
								u.lastAccessed = datetime()
						SET u = { realId: $user_real_id, #user_index#: $user_id, name: $user_name, email: $user_id, createdAt: datetime(), lastAccessed: datetime() }
						MERGE (c)-[x:INSIGHT { #conversation_index#: $conversation_id }]-(i)
							ON CREATE SET
								x.createdAt = datetime(),
								x.lastAccessed = datetime()
							ON MATCH SET
								x.lastAccessed = datetime()
						SET x = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						MERGE (i)-[y:SPOKE { #conversation_index#: $conversation_id }]-(u)
							ON CREATE SET
								y.createdAt = datetime(),
								y.lastAccessed = datetime()
							ON MATCH SET
								y.lastAccessed = datetime()
						SET y = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime() }
						`)
					result, err := tx.Run(ctx, createInsightQuery, map[string]any{
						"conversation_id": mh.conversationId,
						"insight_id":      actionItem.ID,
						"type":            strings.ToLower(actionItem.Type),
						"content":         actionItem.Text,
						"sequence_number": cnt,
						"assignee_id":     actionItem.Assignee.ID, // TODO: Look into it ID or Name
						"user_real_id":    actionItem.From.ID,
						"user_id":         actionItem.From.ID, // TODO: Look into it
						"user_name":       actionItem.From.Name,
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
				klog.V(6).Infof("ActionItemResult LEAVE\n")
				return err
			}

			cnt++
		}
	}

	// rabbitmq
	wrapperStruct := shared.ActionItemResult{
		Duplicate:        exists,
		ConversationID:   mh.conversationId,
		ActionItemResult: air,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("ActionItemResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("ActionItemResult LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitAsyncActionItem, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("ActionItemResult LEAVE\n")
		return err
	}
	klog.V(3).Infof("ActionItemResult.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("ActionItemResult Succeeded\n")
	klog.V(6).Infof("ActionItemResult LEAVE\n")

	return nil
}

func (mh *MessageHandler) TopicResult(tr *sdkinterfaces.TopicResult) error {
	klog.V(6).Infof("TopicResult ENTER\n")

	data, err := json.Marshal(tr)
	if err != nil {
		klog.V(1).Infof("TopicResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TopicResult LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TopicResult LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("TopicResult:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// any insights?
	if len(tr.Topics) == 0 {
		klog.V(4).Infof("tr.Topics is empty. Nothing to process.\n")
		klog.V(6).Infof("TopicResult LEAVE\n")
		return nil
	}

	// insights already created
	exists, err := mh.DoesConversationExist()
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TopicResult LEAVE\n")
		return err
	}

	// only add to DB if new
	if !exists {
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
								t.createdAt = datetime(),
								t.lastAccessed = datetime()
							ON MATCH SET
								t.lastAccessed = datetime()
						SET t = { #topic_index#: $topic_id, phrases: $phrases, score: $score, type: $type, messageIndex: $symbl_message_index, rootWords: $root_words, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						MERGE (c)-[x:TOPICS { #conversation_index#: $conversation_id }]-(t)
							ON CREATE SET
								x.createdAt = datetime(),
								x.lastAccessed = datetime()
							ON MATCH SET
								x.lastAccessed = datetime()
						SET x = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
						`)
					result, err := tx.Run(ctx, createTopicsQuery, map[string]any{
						"conversation_id":     mh.conversationId,
						"topic_id":            "TODO", // TODO
						"phrases":             strings.ToLower(topic.Text),
						"score":               topic.Score,
						"type":                topic.Type,
						"symbl_message_index": "TODO", // TODO: topic.MessageIndex,
						"root_words":          "TODO", // TODO: convertRootWordToString(topic.RootWords),
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
				klog.V(6).Infof("TopicResult LEAVE\n")
				return err
			}

			// associate topic to message
			for _, msgId := range topic.MessageIds {
				_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						createTopicsQuery := utils.ReplaceIndexes(`
							MATCH (t:Topic { topicId: $topic_id })
							MATCH (m:Message { #message_index#: $message_id })
							MERGE (t)-[x:TOPIC_MESSAGE_REF { #conversation_index#: $conversation_id }]-(m)
								ON CREATE SET
									x.createdAt = datetime(),
									x.lastAccessed = datetime()
								ON MATCH SET
									x.lastAccessed = datetime()
							SET x = { #conversation_index#: $conversation_id, value: $value, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
							`)
						result, err := tx.Run(ctx, createTopicsQuery, map[string]any{
							"conversation_id": mh.conversationId,
							"topic_id":        "TODO", // TODO
							"message_id":      msgId,
							"value":           strings.ToLower(topic.Text),
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
					klog.V(6).Infof("TopicResult LEAVE\n")
					return err
				}
			}
		}
	}

	// rabbitmq
	wrapperStruct := shared.TopicResult{
		Duplicate:      exists,
		ConversationID: mh.conversationId,
		TopicResult:    tr,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("TopicResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TopicResult LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitAsyncTopic, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("TopicResult LEAVE\n")
		return err
	}
	klog.V(3).Infof("TopicResult.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("TopicResult Succeeded\n")
	klog.V(6).Infof("TopicResult LEAVE\n")

	return nil
}
func (mh *MessageHandler) TrackerResult(tr *sdkinterfaces.TrackerResult) error {
	klog.V(6).Infof("TrackerResult ENTER\n")

	data, err := json.Marshal(tr)
	if err != nil {
		klog.V(1).Infof("TrackerResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TrackerResult LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TrackerResult LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("TrackerResult:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// any insights?
	if len(tr.Matches) == 0 {
		klog.V(4).Infof("tr.Matches is empty. Nothing to process.\n")
		klog.V(6).Infof("TrackerResult LEAVE\n")
		return nil
	}

	// insights already created
	exists, err := mh.DoesConversationExist()
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TrackerResult LEAVE\n")
		return err
	}

	// only add to DB if new
	if !exists {
		// write the object to the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				createTrackersQuery := utils.ReplaceIndexes(`
					MATCH (c:Conversation { #conversation_index#: $conversation_id })
					MERGE (t:Tracker { #tracker_index#: $tracker_id })
						ON CREATE SET
							t.createdAt = datetime(),
							t.lastAccessed = datetime()
						ON MATCH SET
							t.lastAccessed = datetime()
					SET t = { #tracker_index#: $tracker_id, name: $tracker_name, createdAt: datetime(), lastAccessed: datetime() }
					MERGE (c)-[x:TRACKER { #conversation_index#: $conversation_id }]-(t)
						ON CREATE SET
							x.createdAt = datetime(),
							x.lastAccessed = datetime()
						ON MATCH SET
							x.lastAccessed = datetime()
					SET x = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
					`)
				result, err := tx.Run(ctx, createTrackersQuery, map[string]any{
					"conversation_id": mh.conversationId,
					"tracker_id":      tr.ID,
					"tracker_name":    strings.ToLower(tr.Name),
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
			klog.V(6).Infof("TrackerResult LEAVE\n")
			return err
		}

		// associate tracker to messages and insights
		for _, match := range tr.Matches {

			// messages
			for _, msgRef := range match.MessageRefs {
				_, err = (*mh.neo4jMgr).ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						createTopicsQuery := utils.ReplaceIndexes(`
							MATCH (t:Tracker { #tracker_index#: $tracker_id })
							MATCH (m:Message { #message_index#: $message_id })
							MERGE (t)-[x:TRACKER_MESSAGE_REF { #conversation_index#: $conversation_id }]-(m)
								ON CREATE SET
									x.createdAt = datetime(),
									x.lastAccessed = datetime()
								ON MATCH SET
									x.lastAccessed = datetime()
							SET x = { #conversation_index#: $conversation_id, name: $tracker_name, value: $value, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
							`)
						result, err := tx.Run(ctx, createTopicsQuery, map[string]any{
							"conversation_id": mh.conversationId,
							"tracker_id":      tr.ID,
							"message_id":      msgRef.ID,
							"tracker_name":    strings.ToLower(tr.Name),
							"value":           strings.ToLower(match.Value),
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
					klog.V(6).Infof("TrackerResult LEAVE\n")
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
									x.createdAt = datetime(),
									x.lastAccessed = datetime()
								ON MATCH SET
									x.lastAccessed = datetime()
							SET x = { #conversation_index#: $conversation_id, name: $tracker_name, value: $value, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
							`)
						result, err := tx.Run(ctx, createTrackerMatchQuery, map[string]any{
							"conversation_id": mh.conversationId,
							"tracker_id":      tr.ID,
							"insight_id":      inRef.ID,
							"tracker_name":    strings.ToLower(tr.Name),
							"value":           strings.ToLower(match.Value),
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
					klog.V(6).Infof("TrackerResult LEAVE\n")
					return err
				}
			}
		}
	}

	// rabbitmq
	wrapperStruct := shared.TrackerResult{
		Duplicate:      exists,
		ConversationID: mh.conversationId,
		TrackerResult:  tr,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("TrackerResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TrackerResult LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitAsyncTracker, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("TrackerResult LEAVE\n")
		return err
	}
	klog.V(3).Infof("TrackerResult.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("TrackerResult Succeeded\n")
	klog.V(6).Infof("TrackerResult LEAVE\n")

	return nil
}

func (mh *MessageHandler) EntityResult(er *sdkinterfaces.EntityResult) error {
	klog.V(6).Infof("EntityResult ENTER\n")

	data, err := json.Marshal(er)
	if err != nil {
		klog.V(1).Infof("EntityResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("EntityResult LEAVE\n")
		return err
	}

	// pretty print
	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("EntityResult LEAVE\n")
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("EntityResult:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// any insights?
	if len(er.Entities) == 0 {
		klog.V(4).Infof("er.Entities is empty. Nothing to process.\n")
		klog.V(6).Infof("EntityResult LEAVE\n")
		return nil
	}

	// insights already created
	exists, err := mh.DoesConversationExist()
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("EntityResult LEAVE\n")
		return err
	}

	// only add to DB if new
	if !exists {
		// write the object to the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for _, entity := range er.Entities {

			// associate tracker to messages and insights
			for _, match := range entity.Matches {

				// entity id
				entityCategory := strings.ToLower(strings.ReplaceAll(entity.Category, " ", "_"))
				entityType := strings.ToLower(strings.ReplaceAll(entity.Type, " ", "_"))
				entitySubType := strings.ToLower(strings.ReplaceAll(entity.SubType, " ", "_"))
				entityValue := strings.ToLower(strings.ReplaceAll(match.DetectedValue, " ", "_"))
				entityId := fmt.Sprintf("%s/%s/%s/%s", entityCategory, entityType, entitySubType, entityValue)

				// entity
				_, err := (*mh.neo4jMgr).ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						createEntitiesQuery := utils.ReplaceIndexes(`
							MATCH (c:Conversation { #conversation_index#: $conversation_id })
							MERGE (e:Entity { #entity_index#: $entity_id })
								ON CREATE SET
									e.createdAt = datetime(),
									e.lastAccessed = datetime()
								ON MATCH SET
									e.lastAccessed = datetime()
							SET e = { #entity_index#: $entity_id, type: $type, subType: $sub_type, category: $category, value: $value, createdAt: datetime(), lastAccessed: datetime() }
							MERGE (c)-[x:ENTITY { #conversation_index#: $conversation_id }]-(e)
								ON CREATE SET
									x.createdAt = datetime(),
									x.lastAccessed = datetime()
								ON MATCH SET
									x.lastAccessed = datetime()
							SET x = { #conversation_index#: $conversation_id, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
							`)
						result, err := tx.Run(ctx, createEntitiesQuery, map[string]any{
							"conversation_id": mh.conversationId,
							"entity_id":       entityId,
							"type":            strings.ToLower(entity.Type),
							"sub_type":        strings.ToLower(entity.SubType),
							"category":        strings.ToLower(entity.Category),
							"value":           strings.ToLower(match.DetectedValue),
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
					klog.V(6).Infof("EntityResult LEAVE\n")
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
										x.createdAt = datetime(),
										x.lastAccessed = datetime()
									ON MATCH SET
										x.lastAccessed = datetime()
								SET x = { #conversation_index#: $conversation_id, value: $value, createdAt: datetime(), lastAccessed: datetime(), raw: $raw }
								`)
							result, err := tx.Run(ctx, createEntitiesQuery, map[string]any{
								"conversation_id": mh.conversationId,
								"entity_id":       entityId,
								"message_id":      msgRef.ID,
								"value":           strings.ToLower(match.DetectedValue),
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
						klog.V(6).Infof("EntityResult LEAVE\n")
						return err
					}
				}
			}
		}
	}

	// rabbitmq
	wrapperStruct := shared.EntityResult{
		Duplicate:      exists,
		ConversationID: mh.conversationId,
		EntityResult:   er,
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("EntityResult json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("EntityResult LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitAsyncEntity, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("EntityResult LEAVE\n")
		return err
	}
	klog.V(3).Infof("EntityResult.PublishWithContext:\n%s\n", string(data))

	klog.V(4).Infof("EntityResult Succeeded\n")
	klog.V(6).Infof("EntityResult LEAVE\n")

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

	// insights already created
	exists, err := mh.DoesConversationExist()
	if err != nil {
		klog.V(1).Infof("json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TeardownConversation LEAVE\n")
		return err
	}

	// rabbitmq
	wrapperStruct := shared.TeardownResult{
		Duplicate: exists,
		TeardownMessage: &asyncinterfaces.TeardownMessage{
			ConversationID: mh.conversationId,
		},
	}

	data, err = json.Marshal(wrapperStruct)
	if err != nil {
		klog.V(1).Infof("TeardownConversation json.Marshal failed. Err: %v\n", err)
		klog.V(6).Infof("TeardownConversation LEAVE\n")
		return err
	}

	err = (*mh.rabbitMgr).PublishMessageByName(shared.RabbitAsyncConversationTeardown, data)
	if err != nil {
		klog.V(1).Infof("PublishMessageByName failed. Err: %v\n", err)
		klog.V(6).Infof("TeardownConversation LEAVE\n")
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
