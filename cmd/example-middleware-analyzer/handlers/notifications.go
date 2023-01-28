// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package handlers

import (
	"context"
	"encoding/json"

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
	return nil
}

func (h *Handler) MessageResponseMessage(mr *sdkinterfaces.MessageResponse) error {
	return nil
}

func (h *Handler) InsightResponseMessage(ir *sdkinterfaces.InsightResponse) error {
	return nil
}

func (h *Handler) TopicResponseMessage(tr *sdkinterfaces.TopicResponse) error {
	return nil
}

func (h *Handler) TrackerResponseMessage(tr *sdkinterfaces.TrackerResponse) error {
	/*
		Example: See what other conversations mentioned certain trackers
	*/
	ctx := context.Background()

	for _, tracker := range tr.Trackers {
		// messages
		_, err := (*h.session).ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			myQuery := commoninterfaces.ReplaceIndexes(`
				MATCH (t:Tracker)-[x:TRACKER_MESSAGE_REF]-(m:Message)-[y:SPOKE]-(u:User)
				WHERE x.#conversation_index# <> $conversation_id AND y.#conversation_index# <> $conversation_id AND t.name = $tracker_name
				RETURN t, x, m, y, u ORDER BY x.created DESC LIMIT 5`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": h.conversationID,
				"tracker_name":    tracker.Name,
			})
			if err != nil {
				return nil, err
			}

			// message
			klog.V(2).Infof("Check Previous Tracker [Message]\n")
			klog.V(2).Infof("----------------------------------------\n")
			for result.Next(ctx) {
				message := result.Record().Values[2].(neo4j.Node)
				relationship := result.Record().Values[1].(neo4j.Relationship)
				user := result.Record().Values[4].(neo4j.Node)

				klog.V(2).Infof("Previous Tracker [Message]\n")
				klog.V(2).Infof("Author: %s / %s\n", user.Props["name"].(string), user.Props["email"].(string))
				klog.V(2).Infof("Tracker Match: %s\n", relationship.Props["value"].(string))
				klog.V(2).Infof("Corresponding sentence: %s\n", message.Props["content"].(string))

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

						/*
							This sends your High-level Application message back to the Dataminer component
						*/
						err = (*h.msgPublisher).PublishMessage(h.conversationID, data)
						if err != nil {
							klog.V(1).Infof("PushNotification failed. Err: %v\n", err)
						}
					}
				}
			}
			klog.V(2).Infof("----------------------------------------\n")

			return nil, result.Err()
		})
		if err != nil {
			klog.V(1).Infof("[EntityHandler] ExecuteRead failed. Err: %v\n", err)
			return err
		}

		// insights
		_, err = (*h.session).ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			myQuery := commoninterfaces.ReplaceIndexes(`
				MATCH (t:Tracker)-[x:TRACKER_INSIGHT_REF]-(i:Insight)-[y:SPOKE]-(u:User)
				WHERE x.#conversation_index# <> $conversation_id AND y.#conversation_index# <> $conversation_id AND t.name = $tracker_name
				RETURN t, x, i, y, u ORDER BY x.created DESC LIMIT 5`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": h.conversationID,
				"tracker_name":    tracker.Name,
			})
			if err != nil {
				return nil, err
			}

			// insight
			klog.V(2).Infof("Check Previous Tracker [Insight]\n")
			klog.V(2).Infof("----------------------------------------\n")
			for result.Next(ctx) {
				insight := result.Record().Values[2].(neo4j.Node)
				relationship := result.Record().Values[1].(neo4j.Relationship)
				user := result.Record().Values[4].(neo4j.Node)

				klog.V(2).Infof("Previous Tracker [Insight]\n")
				klog.V(2).Infof("Author: %s / %s\n", user.Props["name"].(string), user.Props["email"].(string))
				klog.V(2).Infof("Tracker Match: %s\n", relationship.Props["value"].(string))
				klog.V(2).Infof("Corresponding sentence: %s\n", insight.Props["content"].(string))

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

						/*
							This sends your High-level Application message back to the Dataminer component
						*/
						err = (*h.msgPublisher).PublishMessage(h.conversationID, data)
						if err != nil {
							klog.V(1).Infof("PushNotification failed. Err: %v\n", err)
						}
					}
				}
			}
			klog.V(2).Infof("----------------------------------------\n")

			return nil, result.Err()
		})
		if err != nil {
			klog.V(1).Infof("[EntityHandler] ExecuteRead failed. Err: %v\n", err)
			return err
		}
	}

	return nil
}

func (h *Handler) EntityResponseMessage(tr *sdkinterfaces.EntityResponse) error {
	return nil
}

func (h *Handler) TeardownConversation(tm *sdkinterfaces.TeardownMessage) error {
	return nil
}

func (h *Handler) UserDefinedMessage(data []byte) error {
	return nil
}

func (h *Handler) UnhandledMessage(byMsg []byte) error {
	// TODO
	return nil
}
