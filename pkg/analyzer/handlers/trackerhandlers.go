// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package router

import (
	"context"
	"encoding/json"
	"fmt"

	prettyjson "github.com/hokaccha/go-prettyjson"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"

	callback "github.com/dvonthenen/enterprise-reference-implementation/pkg/analyzer/rabbit/interfaces"
	interfaces "github.com/dvonthenen/enterprise-reference-implementation/pkg/interfaces"
)

func NewTrackerHandler(options HandlerOptions) *callback.RabbitMessageHandler {
	var handler callback.RabbitMessageHandler
	handler = TrackerHandler{
		session:      options.Session,
		symblClient:  options.SymblClient,
		pushCallback: options.PushCallback,
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
				RETURN m, u`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": tr.ConversationID,
				"tracker_name":    tracker.Name,
			})
			if err != nil {
				return nil, err
			}

			// message
			for result.Next(ctx) {
				message := result.Record().Values[0].(neo4j.Node)
				user := result.Record().Values[1].(neo4j.Node)

				messageTemplate := `
					-------------------------------<br>
					User: %s<br>
					Message: %s<br>
					-------------------------------<br><br>
				`
				msg := fmt.Sprintf(messageTemplate, user.Props["name"].(string), message.Props["content"].(string))

				err := (*ch.pushCallback).PushNotification(msg)
				if err != nil {
					klog.V(1).Infof("PushNotification failed. Err: %v\n", err)
				}
			}

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
				RETURN i, u`)
			result, err := tx.Run(ctx, myQuery, map[string]any{
				"conversation_id": tr.ConversationID,
				"tracker_name":    tracker.Name,
			})
			if err != nil {
				return nil, err
			}

			// insight
			for result.Next(ctx) {
				insight := result.Record().Values[0].(neo4j.Node)
				user := result.Record().Values[1].(neo4j.Node)

				messageTemplate := `
				-------------------------------<br>
				User: %s<br>
				Message: %s<br>
				-------------------------------<br><br>
				`
				msg := fmt.Sprintf(messageTemplate, user.Props["name"].(string), insight.Props["content"].(string))

				err := (*ch.pushCallback).PushNotification(msg)
				if err != nil {
					klog.V(1).Infof("PushNotification failed. Err: %v\n", err)
				}
			}

			return nil, result.Err()
		})
		if err != nil {
			klog.V(1).Infof("[EntityHandler] ExecuteRead failed. Err: %v\n", err)
			return err
		}
	}

	return nil
}
