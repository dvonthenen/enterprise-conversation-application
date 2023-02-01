// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"

	interfaces "github.com/dvonthenen/enterprise-reference-implementation/cmd/example-middleware-analyzer/interfaces"
	commoninterfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
	middlewareinterfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/middleware-analyzer/interfaces"
)

func NewHandler(options HandlerOptions) *Handler {
	handler := Handler{
		session:     options.Session,
		symblClient: options.SymblClient,
		cache:       NewMessageCache(),
	}
	return &handler
}

func (h *Handler) SetClientPublisher(mp *middlewareinterfaces.MessagePublisher) {
	klog.V(4).Infof("SetClientPublisher called...\n")
	h.msgPublisher = mp
}

func (h *Handler) InitializedConversation(im *sdkinterfaces.InitializationMessage) error {
	h.conversationID = im.Message.Data.ConversationID
	klog.V(2).Infof("conversationID: %s\n", h.conversationID)
	return nil
}

func (h *Handler) RecognitionResultMessage(rr *sdkinterfaces.RecognitionResult) error {
	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) MessageResponseMessage(mr *sdkinterfaces.MessageResponse) error {
	for _, msg := range mr.Messages {
		h.cache.Push(msg.ID, msg.Payload.Content)
	}
	return nil
}

func (h *Handler) InsightResponseMessage(ir *sdkinterfaces.InsightResponse) error {
	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) TopicResponseMessage(tr *sdkinterfaces.TopicResponse) error {
	ctx := context.Background()

	for _, curTopic := range tr.Topics {
		// housekeeping
		atLeastOnce := false
		var msg *interfaces.AppSpecificHistorical

		// get past instances
		_, err := (*h.session).ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			myQuery := commoninterfaces.ReplaceIndexes(`
				MATCH (t:Topic)-[x:TOPIC_MESSAGE_REF]-(m:Message)-[y:SPOKE]-(u:User)
				WHERE x.#conversation_index# <> $conversation_id AND y.#conversation_index# <> $conversation_id AND t.value = $topic_phrases
				RETURN t, x, m, y, u ORDER BY x.created DESC LIMIT 5`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": h.conversationID,
				"topic_phrases":   curTopic.Phrases,
			})
			if err != nil {
				return nil, err
			}

			for result.Next(ctx) {
				// only init once!
				if !atLeastOnce {
					klog.V(2).Infof("Check Previous Topics\n")
					klog.V(2).Infof("----------------------------------------\n")

					msg = &interfaces.AppSpecificHistorical{
						Type: sdkinterfaces.MessageTypeUserDefined,
					}
					atLeastOnce = true
				}

				prevTopic := result.Record().Values[0].(neo4j.Node)
				message := result.Record().Values[2].(neo4j.Node)
				relationship := result.Record().Values[1].(neo4j.Relationship)
				user := result.Record().Values[4].(neo4j.Node)

				klog.V(2).Infof("Previous Topic\n")
				klog.V(2).Infof("Author: %s / %s\n", user.Props["name"].(string), user.Props["email"].(string))
				klog.V(2).Infof("Topic Match: %s\n", relationship.Props["value"].(string))
				klog.V(2).Infof("Corresponding sentence: %s\n", message.Props["content"].(string))

				content := ""
				for _, msgRef := range curTopic.MessageReferences {
					contentTmp, err := h.cache.Find(msgRef.ID)
					if err != nil {
						klog.V(4).Infof("Msg ID not found: %s\n", msgRef.ID)
						continue
					}

					if len(content) > 0 {
						content += "/"
					}
					content += contentTmp
				}

				msg.Data = append(msg.Data, &interfaces.Data{
					Type: interfaces.UserMessageTypeEntityAssociation,
					Author: interfaces.Author{
						Name:  user.Props["name"].(string),
						Email: user.Props["email"].(string),
					},
					Message: interfaces.Message{
						Correlation:     curTopic.Phrases,
						CurrentContent:  content,
						CurrentMatch:    convertRootWordToString(curTopic.RootWords),
						PreviousContent: message.Props["content"].(string),
						PreviousMatch:   prevTopic.Props["rootWords"].(string),
					},
				})
			}

			if atLeastOnce {
				klog.V(2).Infof("----------------------------------------\n")
			}

			return nil, result.Err()
		})
		if err != nil {
			klog.V(1).Infof("[Entities] ExecuteRead failed. Err: %v\n", err)
			return err
		}

		/*
			If there is at least one Message or Insight that has triggered this Tracker, then
			send your High-level Application message back to the Dataminer component
		*/
		if atLeastOnce {
			data, err := json.Marshal(*msg)
			if err != nil {
				klog.V(1).Infof("[Entities] json.Marshal failed. Err: %v\n", err)
			}

			err = (*h.msgPublisher).PublishMessage(h.conversationID, data)
			if err != nil {
				klog.V(1).Infof("[Entities] PublishMessage failed. Err: %v\n", err)
			}
		}
	}

	return nil
}

func (h *Handler) TrackerResponseMessage(tr *sdkinterfaces.TrackerResponse) error {
	ctx := context.Background()

	for _, curTracker := range tr.Trackers {
		// housekeeping
		atLeastOnceMessage := false
		atLeastOnceInsight := false
		var msg *interfaces.AppSpecificHistorical

		// get past messages
		_, err := (*h.session).ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			myQuery := commoninterfaces.ReplaceIndexes(`
				MATCH (t:Tracker)-[x:TRACKER_MESSAGE_REF]-(m:Message)-[y:SPOKE]-(u:User)
				WHERE x.#conversation_index# <> $conversation_id AND y.#conversation_index# <> $conversation_id AND t.name = $tracker_name
				RETURN t, x, m, y, u ORDER BY x.created DESC LIMIT 5`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": h.conversationID,
				"tracker_name":    curTracker.Name,
			})
			if err != nil {
				return nil, err
			}

			for result.Next(ctx) {
				// only init once!
				if !atLeastOnceMessage {
					klog.V(2).Infof("Check Previous Tracker [Message]\n")
					klog.V(2).Infof("----------------------------------------\n")

					msg = &interfaces.AppSpecificHistorical{
						Type: sdkinterfaces.MessageTypeUserDefined,
					}
					atLeastOnceMessage = true
				}

				message := result.Record().Values[2].(neo4j.Node)
				relationship := result.Record().Values[1].(neo4j.Relationship)
				user := result.Record().Values[4].(neo4j.Node)

				klog.V(2).Infof("Previous Tracker [Message]\n")
				klog.V(2).Infof("Author: %s / %s\n", user.Props["name"].(string), user.Props["email"].(string))
				klog.V(2).Infof("Tracker Match: %s\n", relationship.Props["value"].(string))
				klog.V(2).Infof("Corresponding sentence: %s\n", message.Props["content"].(string))

				for _, match := range curTracker.Matches {
					for _, refs := range match.MessageRefs {
						msg.Data = append(msg.Data, &interfaces.Data{
							Type: interfaces.UserMessageTypeTrackerAssociation,
							Author: interfaces.Author{
								Name:  user.Props["name"].(string),
								Email: user.Props["email"].(string),
							},
							Message: interfaces.Message{
								Correlation:     curTracker.Name,
								CurrentContent:  refs.Text,
								CurrentMatch:    match.Value,
								PreviousContent: message.Props["content"].(string),
								PreviousMatch:   relationship.Props["value"].(string),
							},
						})
					}
				}
			}

			if atLeastOnceMessage {
				klog.V(2).Infof("----------------------------------------\n")
			}

			return nil, result.Err()
		})
		if err != nil {
			klog.V(1).Infof("[Tracker] ExecuteRead failed. Err: %v\n", err)
			return err
		}

		// get past insights
		_, err = (*h.session).ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			myQuery := commoninterfaces.ReplaceIndexes(`
				MATCH (t:Tracker)-[x:TRACKER_INSIGHT_REF]-(i:Insight)-[y:SPOKE]-(u:User)
				WHERE x.#conversation_index# <> $conversation_id AND y.#conversation_index# <> $conversation_id AND t.name = $tracker_name
				RETURN t, x, i, y, u ORDER BY x.created DESC LIMIT 5`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": h.conversationID,
				"tracker_name":    curTracker.Name,
			})
			if err != nil {
				return nil, err
			}

			// insight
			for result.Next(ctx) {
				// only init once!
				if !atLeastOnceInsight {
					klog.V(2).Infof("Check Previous Tracker [Insight]\n")
					klog.V(2).Infof("----------------------------------------\n")

					msg = &interfaces.AppSpecificHistorical{
						Type: sdkinterfaces.MessageTypeUserDefined,
					}
					atLeastOnceInsight = true
				}

				insight := result.Record().Values[2].(neo4j.Node)
				relationship := result.Record().Values[1].(neo4j.Relationship)
				user := result.Record().Values[4].(neo4j.Node)

				klog.V(2).Infof("Previous Tracker [Insight]\n")
				klog.V(2).Infof("Author: %s / %s\n", user.Props["name"].(string), user.Props["email"].(string))
				klog.V(2).Infof("Tracker Match: %s\n", relationship.Props["value"].(string))
				klog.V(2).Infof("Corresponding sentence: %s\n", insight.Props["content"].(string))

				for _, match := range curTracker.Matches {
					for _, refs := range match.InsightRefs {
						// send a message per match

						msg.Data = append(msg.Data, &interfaces.Data{
							Type: interfaces.UserMessageTypeTrackerAssociation,
							Author: interfaces.Author{
								Name:  user.Props["name"].(string),
								Email: user.Props["email"].(string),
							},
							Message: interfaces.Message{
								Correlation:     curTracker.Name,
								CurrentContent:  refs.Text,
								CurrentMatch:    match.Value,
								PreviousContent: insight.Props["content"].(string),
								PreviousMatch:   relationship.Props["value"].(string),
							},
						})
					}
				}
			}

			if atLeastOnceInsight {
				klog.V(2).Infof("----------------------------------------\n")
			}

			return nil, result.Err()
		})
		if err != nil {
			klog.V(1).Infof("[Tracker] ExecuteRead failed. Err: %v\n", err)
			return err
		}

		/*
			If there is at least one Message or Insight that has triggered this Tracker, then
			send your High-level Application message back to the Dataminer component
		*/
		if atLeastOnceMessage || atLeastOnceInsight {
			data, err := json.Marshal(*msg)
			if err != nil {
				klog.V(1).Infof("[Tracker] json.Marshal failed. Err: %v\n", err)
			}

			err = (*h.msgPublisher).PublishMessage(h.conversationID, data)
			if err != nil {
				klog.V(1).Infof("[Tracker] PublishMessage failed. Err: %v\n", err)
			}
		}
	}

	return nil
}

func (h *Handler) EntityResponseMessage(er *sdkinterfaces.EntityResponse) error {
	ctx := context.Background()

	for _, entity := range er.Entities {
		// housekeeping
		atLeastOnce := false
		var msg *interfaces.AppSpecificHistorical

		// get past instances
		_, err := (*h.session).ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			myQuery := commoninterfaces.ReplaceIndexes(`
				MATCH (e:Entity)-[x:ENTITY_MESSAGE_REF]-(m:Message)-[y:SPOKE]-(u:User)
				WHERE x.#conversation_index# <> $conversation_id AND y.#conversation_index# <> $conversation_id AND e.type = $entity_type AND e.subType = $entity_subtype
				RETURN e, x, m, y, u ORDER BY x.created DESC LIMIT 5`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": h.conversationID,
				"entity_type":     entity.Type,
				"entity_subtype":  entity.SubType,
			})
			if err != nil {
				return nil, err
			}

			for result.Next(ctx) {
				// only init once!
				if !atLeastOnce {
					klog.V(2).Infof("Check Previous Entities\n")
					klog.V(2).Infof("----------------------------------------\n")

					msg = &interfaces.AppSpecificHistorical{
						Type: sdkinterfaces.MessageTypeUserDefined,
					}
					atLeastOnce = true
				}

				message := result.Record().Values[2].(neo4j.Node)
				relationship := result.Record().Values[1].(neo4j.Relationship)
				user := result.Record().Values[4].(neo4j.Node)

				klog.V(2).Infof("Previous Entity\n")
				klog.V(2).Infof("Author: %s / %s\n", user.Props["name"].(string), user.Props["email"].(string))
				klog.V(2).Infof("Entity Match: %s\n", relationship.Props["value"].(string))
				klog.V(2).Infof("Corresponding sentence: %s\n", message.Props["content"].(string))

				for _, match := range entity.Matches {
					for _, refs := range match.MessageRefs {
						msg.Data = append(msg.Data, &interfaces.Data{
							Type: interfaces.UserMessageTypeEntityAssociation,
							Author: interfaces.Author{
								Name:  user.Props["name"].(string),
								Email: user.Props["email"].(string),
							},
							Message: interfaces.Message{
								Correlation:     fmt.Sprintf("%s/%s", entity.Type, entity.SubType),
								CurrentContent:  refs.Text,
								CurrentMatch:    match.DetectedValue,
								PreviousContent: message.Props["content"].(string),
								PreviousMatch:   relationship.Props["value"].(string),
							},
						})
					}
				}
			}

			if atLeastOnce {
				klog.V(2).Infof("----------------------------------------\n")
			}

			return nil, result.Err()
		})
		if err != nil {
			klog.V(1).Infof("[Entities] ExecuteRead failed. Err: %v\n", err)
			return err
		}

		/*
			If there is at least one Message or Insight that has triggered this Tracker, then
			send your High-level Application message back to the Dataminer component
		*/
		if atLeastOnce {
			data, err := json.Marshal(*msg)
			if err != nil {
				klog.V(1).Infof("[Entities] json.Marshal failed. Err: %v\n", err)
			}

			err = (*h.msgPublisher).PublishMessage(h.conversationID, data)
			if err != nil {
				klog.V(1).Infof("[Entities] PublishMessage failed. Err: %v\n", err)
			}
		}
	}

	return nil
}

func (h *Handler) TeardownConversation(tm *sdkinterfaces.TeardownMessage) error {
	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) UserDefinedMessage(data []byte) error {
	// No implementation required. Return Succeess!
	return nil
}

func (h *Handler) UnhandledMessage(byMsg []byte) error {
	klog.V(1).Infof("\n\n-------------------------------\n")
	klog.V(1).Infof("UnhandledMessage:\n%v\n", string(byMsg))
	klog.V(1).Infof("-------------------------------\n\n")
	return ErrUnhandledMessage
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
