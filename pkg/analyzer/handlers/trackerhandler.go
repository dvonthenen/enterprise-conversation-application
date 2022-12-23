// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	"context"
	"encoding/json"

	prettyjson "github.com/hokaccha/go-prettyjson"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"

	rabbitinterfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
	sdkinterfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
)

func NewTrackerHandler(options HandlerOptions) *rabbitinterfaces.RabbitMessageHandler {
	var handler rabbitinterfaces.RabbitMessageHandler
	handler = TrackerHandler{
		session:     options.Session,
		symblClient: options.SymblClient,
		manager:     options.Manager,
	}
	return &handler
}

func (ch TrackerHandler) ProcessMessage(byData []byte) error {
	// pretty print
	prettyJson, err := prettyjson.Format(byData)
	if err != nil {
		klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
		return err
	}
	klog.V(6).Infof("\n\n-------------------------------\n")
	klog.V(2).Infof("TrackerHandler:\n%v\n", string(prettyJson))
	klog.V(6).Infof("-------------------------------\n\n")

	// reform struct
	var tr interfaces.TrackerResponse
	err = json.Unmarshal(byData, &tr)
	if err != nil {
		klog.V(1).Infof("[TrackerHandler] json.Unmarshal failed. Err: %v\n", err)
		return err
	}

	// TODO: template for add your businesss logic

	/*
		Example: See what other conversations mentioned certain trackers
	*/
	ctx := context.Background()

	for _, tracker := range tr.TrackerResponse.Trackers {
		// messages
		_, err := (*ch.session).ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			myQuery := interfaces.ReplaceIndexes(`
				MATCH (t:Tracker)-[x:TRACKER_MESSAGE_REF]-(m:Message)-[y:SPOKE]-(u:User)
				WHERE x.#conversation_index# <> $conversation_id AND y.#conversation_index# <> $conversation_id AND t.name = $tracker_name
				RETURN t, x, m, y, u`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": tr.ConversationID,
				"tracker_name":    tracker.Name,
			})
			if err != nil {
				return nil, err
			}

			// message
			klog.V(3).Infof("Check Previous Tracker [Message]\n")
			klog.V(3).Infof("----------------------------------------\n")
			for result.Next(ctx) {
				message := result.Record().Values[2].(neo4j.Node)
				relationship := result.Record().Values[1].(neo4j.Relationship)
				user := result.Record().Values[4].(neo4j.Node)

				klog.V(3).Infof("Previous Tracker [Message]\n")
				klog.V(3).Infof("Author: %s / %s\n", user.Props["name"].(string), user.Props["email"].(string))
				klog.V(3).Infof("Tracker Match: %s\n", relationship.Props["value"].(string))
				klog.V(3).Infof("Corresponding sentence: %s\n", message.Props["content"].(string))

				for _, match := range tracker.Matches {
					for _, refs := range match.MessageRefs {
						// send a message per match
						msg := interfaces.ClientTrackerMessage{
							Type: sdkinterfaces.MessageTypeUserDefined,
							Data: interfaces.Data{
								Type: interfaces.UserMessageTypeAssociation,
								Author: interfaces.Author{
									Name:  user.Props["name"].(string),
									Email: user.Props["email"].(string),
								},
								Message: interfaces.Message{
									Correlation:     tracker.Name,
									CurrentContent:  refs.Text,
									CurrentMatch:    match.Value,
									PreviousContent: message.Props["content"].(string),
									PreviousMatch:   relationship.Props["value"].(string),
								},
							},
						}

						data, err := json.Marshal(msg)
						if err != nil {
							klog.V(1).Infof("MessageResponse json.Marshal failed. Err: %v\n", err)
						}

						err = (*ch.manager).PublishMessageByChannelName(tr.ConversationID, data)
						if err != nil {
							klog.V(1).Infof("PushNotification failed. Err: %v\n", err)
						}
					}
				}
			}
			klog.V(3).Infof("----------------------------------------\n")

			return nil, result.Err()
		})
		if err != nil {
			klog.V(1).Infof("[EntityHandler] ExecuteRead failed. Err: %v\n", err)
			return err
		}

		// insights
		_, err = (*ch.session).ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			myQuery := interfaces.ReplaceIndexes(`
				MATCH (t:Tracker)-[x:TRACKER_INSIGHT_REF]-(i:Insight)-[y:SPOKE]-(u:User)
				WHERE x.#conversation_index# <> $conversation_id AND y.#conversation_index# <> $conversation_id AND t.name = $tracker_name
				RETURN t, x, i, y, u`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": tr.ConversationID,
				"tracker_name":    tracker.Name,
			})
			if err != nil {
				return nil, err
			}

			// insight
			klog.V(3).Infof("Check Previous Tracker [Insight]\n")
			klog.V(3).Infof("----------------------------------------\n")
			for result.Next(ctx) {
				insight := result.Record().Values[2].(neo4j.Node)
				relationship := result.Record().Values[1].(neo4j.Relationship)
				user := result.Record().Values[4].(neo4j.Node)

				klog.V(3).Infof("Previous Tracker [Insight]\n")
				klog.V(3).Infof("Author: %s / %s\n", user.Props["name"].(string), user.Props["email"].(string))
				klog.V(3).Infof("Tracker Match: %s\n", relationship.Props["value"].(string))
				klog.V(3).Infof("Corresponding sentence: %s\n", insight.Props["content"].(string))

				for _, match := range tracker.Matches {
					for _, refs := range match.InsightRefs {
						// send a message per match

						msg := interfaces.ClientTrackerMessage{
							Type: sdkinterfaces.MessageTypeUserDefined,
							Data: interfaces.Data{
								Type: interfaces.UserMessageTypeAssociation,
								Author: interfaces.Author{
									Name:  user.Props["name"].(string),
									Email: user.Props["email"].(string),
								},
								Message: interfaces.Message{
									Correlation:     tracker.Name,
									CurrentContent:  refs.Text,
									CurrentMatch:    match.Value,
									PreviousContent: insight.Props["content"].(string),
									PreviousMatch:   relationship.Props["value"].(string),
								},
							},
						}

						data, err := json.Marshal(msg)
						if err != nil {
							klog.V(1).Infof("MessageResponse json.Marshal failed. Err: %v\n", err)
						}

						err = (*ch.manager).PublishMessageByChannelName(tr.ConversationID, data)
						if err != nil {
							klog.V(1).Infof("PushNotification failed. Err: %v\n", err)
						}
					}
				}
			}
			klog.V(3).Infof("----------------------------------------\n")

			return nil, result.Err()
		})
		if err != nil {
			klog.V(1).Infof("[EntityHandler] ExecuteRead failed. Err: %v\n", err)
			return err
		}
	}

	return nil
}
