// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package routing

import (
	"context"
	"encoding/json"
	"fmt"

	interfaces "github.com/dvonthenen/symbl-go-sdk/pkg/api/streaming/v1/interfaces"
	neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	klog "k8s.io/klog/v2"
)

func NewHandler(options MessageHandlerOptions) (*MessageHandler, error) {
	if len(options.ConversationId) == 0 {
		klog.Errorf("conversationId is empty\n")
		return nil, ErrInvalidInput
	}

	mh := &MessageHandler{
		ConversationId: options.ConversationId,
		Session:        options.Session,
	}
	return mh, nil
}

func (mh *MessageHandler) Init() error {
	klog.V(6).Infof("MessageHandler.Init ENTER\n")

	// context
	ctx := context.Background()

	// create conversation object
	_, err := mh.Session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			createConversationQuery := `
				MERGE (c:Conversation { conversationId: $conversation_id })
				SET c = { conversationId: $conversation_id }
				`
			result, err := tx.Run(ctx, createConversationQuery, map[string]any{
				"conversation_id": mh.ConversationId,
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

	klog.V(4).Infof("Init Succeeded\n")
	klog.V(6).Infof("MessageHandler.Init LEAVE\n")

	return nil
}

func (mh *MessageHandler) Teardown() error {
	klog.V(6).Infof("MessageHandler.Teardown ENTER\n")

	// close the session
	ctx := context.Background()
	mh.Session.Close(ctx)

	klog.V(4).Infof("Teardown Succeeded\n")
	klog.V(6).Infof("MessageHandler.Teardown LEAVE\n")

	return nil
}

func (mh *MessageHandler) RecognitionResultMessage(rr *interfaces.RecognitionResult) error {
	data, err := json.Marshal(rr)
	if err != nil {
		klog.V(1).Infof("RecognitionResult json.Marshal failed. Err: %v\n", err)
		return err
	}

	// We probably don't actually need this. Will just leave the debug statements here
	// for future use
	// TODO: fix level
	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("RecognitionResultMessage Object DUMP:\n%v\n\n", string(data))
	klog.V(3).Infof("\nMessage:\n%v\n\n", rr.Message.Punctuated.Transcript)
	klog.V(3).Infof("-------------------------------\n\n")

	return nil
}

func (mh *MessageHandler) MessageResponseMessage(mr *interfaces.MessageResponse) error {
	data, err := json.Marshal(mr)
	if err != nil {
		klog.V(1).Infof("MessageResponse json.Marshal failed. Err: %v\n", err)
		return err
	}

	// TODO: fix level
	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("MessageResponseMessage Object DUMP:\n%v\n", string(data))
	klog.V(3).Infof("-------------------------------\n\n")

	// write the object to the database
	ctx := context.Background()

	// if we need to do something with them
	// for records, message := range mr.Messages {
	for _, message := range mr.Messages {
		_, err := mh.Session.ExecuteWrite(ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				createMessageToPeopleQuery := `
					MATCH (c:Conversation { conversationId: $conversation_id })
					MERGE (m:Message { messageId: $message_id })
						ON CREATE SET
							m.lastAccessed = timestamp()
						ON MATCH SET
							m.lastAccessed = timestamp()
					SET m = { messageId: $message_id, content: $content, startTime: $start_time, endTime: $end_time, timeOffset: $time_offset, duration: $duration, raw: $raw }
					MERGE (u:User { userId: $user_id })
						ON CREATE SET
							u.lastAccessed = timestamp()
						ON MATCH SET
							u.lastAccessed = timestamp()
					SET u = { realId: $user_real_id, userId: $user_id, name: $user_name, email: $user_id }
					MERGE (c)-[:MESSAGES]-(m)
					MERGE (m)-[:SPOKE]-(u)
					`
				result, err := tx.Run(ctx, createMessageToPeopleQuery, map[string]any{
					"conversation_id": mh.ConversationId,
					"message_id":      message.ID,
					"content":         message.Payload.Content,
					"start_time":      message.Duration.StartTime,
					"end_time":        message.Duration.EndTime,
					"time_offset":     message.Duration.TimeOffset,
					"duration":        message.Duration.Duration,
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
			return err
		}
	}

	// TODO: pub/sub

	return nil
}

func (mh *MessageHandler) InsightResponseMessage(ir *interfaces.InsightResponse) error {
	for _, insight := range ir.Insights {
		switch insight.Type {
		case interfaces.InsightTypeQuestion:
			return mh.HandleQuestion(&insight)
		case interfaces.InsightTypeFollowUp:
			return mh.HandleFollowUp(&insight)
		case interfaces.InsightTypeActionItem:
			return mh.HandleActionItem(&insight)
		default:
			data, err := json.Marshal(ir)
			if err != nil {
				klog.V(1).Infof("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
				return err
			}

			// TODO: fix level
			klog.V(3).Infof("\n\n-------------------------------\n")
			klog.V(3).Infof("TopicResponseMessage Object DUMP:\n%v\n", string(data))
			klog.V(3).Infof("-------------------------------\n\n")
			return nil
		}
	}

	return nil
}

func (mh *MessageHandler) TopicResponseMessage(tr *interfaces.TopicResponse) error {
	data, err := json.Marshal(tr)
	if err != nil {
		klog.V(1).Infof("TopicResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	// TODO: fix level
	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("TopicResponseMessage Object DUMP:\n%v\n", string(data))
	klog.V(3).Infof("-------------------------------\n\n")

	// write the object to the database
	ctx := context.Background()

	for _, topic := range tr.Topics {
		_, err := mh.Session.ExecuteWrite(ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				createTopicsQuery := `
					MATCH (c:Conversation { conversationId: $conversation_id })
					MERGE (t:Topic { topicId: $topic_id })
						ON CREATE SET
							t.lastAccessed = timestamp()
						ON MATCH SET
							t.lastAccessed = timestamp()
					SET t = { topicId: $topic_id, phrases: $phrases, score: $score, type: $type, messageIndex: $message_index, rootWords: $root_words, raw: $raw }
					MERGE (c)-[:TOPICS]-(t)
					`
				result, err := tx.Run(ctx, createTopicsQuery, map[string]any{
					"conversation_id": mh.ConversationId,
					"topic_id":        topic.ID,
					"phrases":         topic.Phrases,
					"score":           topic.Score,
					"type":            topic.Type,
					"message_index":   topic.MessageIndex,
					"root_words":      ConvertRootWordToSlice(topic.RootWords),
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
			return err
		}

		// associate topic to message
		for _, ref := range topic.MessageReferences {
			_, err = mh.Session.ExecuteWrite(ctx,
				func(tx neo4j.ManagedTransaction) (any, error) {
					createTopicsQuery := `
						MATCH (t:Topic { topicId: $topic_id })
						MATCH (m:Message { messageId: $message_id })
						MERGE (t)-[:TOPIC_MESSAGE_REF]-(m)
						`
					result, err := tx.Run(ctx, createTopicsQuery, map[string]any{
						"topic_id":   topic.ID,
						"message_id": ref.ID,
					})
					if err != nil {
						klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
						return nil, err
					}
					return result.Collect(ctx)
				})
			if err != nil {
				klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
				return err
			}
		}
	}

	// TODO: pub/sub

	return nil
}
func (mh *MessageHandler) TrackerResponseMessage(tr *interfaces.TrackerResponse) error {
	data, err := json.Marshal(tr)
	if err != nil {
		klog.V(1).Infof("TrackerResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	// TODO: fix level
	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("TrackerResponseMessage Object DUMP:\n%v\n", string(data))
	klog.V(3).Infof("-------------------------------\n\n")

	// write the object to the database
	ctx := context.Background()

	for _, tracker := range tr.Trackers {
		_, err := mh.Session.ExecuteWrite(ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				createTrackersQuery := `
					MATCH (c:Conversation { conversationId: $conversation_id })
					MERGE (t:Tracker { trackerId: $tracker_id })
						ON CREATE SET
							t.lastAccessed = timestamp()
						ON MATCH SET
							t.lastAccessed = timestamp()
					SET t = { trackerId: $tracker_id, name: $tracker_name, raw: $raw }
					MERGE (c)-[:TRACKER]-(t)
					`
				result, err := tx.Run(ctx, createTrackersQuery, map[string]any{
					"conversation_id": mh.ConversationId,
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
			return err
		}

		// associate tracker to messages and insights
		for _, match := range tracker.Matches {

			// messages
			for _, msgRef := range match.MessageRefs {
				_, err = mh.Session.ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						createTopicsQuery := `
							MATCH (t:Tracker { trackerId: $tracker_id })
							MATCH (m:Message { messageId: $message_id })
							MERGE (t)-[:TRACKER_MESSAGE_REF]-(m)
							`
						result, err := tx.Run(ctx, createTopicsQuery, map[string]any{
							"tracker_id": tracker.ID,
							"message_id": msgRef.ID,
						})
						if err != nil {
							klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
							return nil, err
						}
						return result.Collect(ctx)
					})
				if err != nil {
					klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
					return err
				}
			}

			// insights
			for _, inRef := range match.InsightRefs {
				_, err = mh.Session.ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						createTrackerMatchQuery := `
							MATCH (t:TrackerMatch { trackerId: $tracker_id })
							MATCH (i:Insight { insightId: $insight_id })
							MERGE (t)-[:TRACKER_INSIGHT_REF]-(i)
							`
						result, err := tx.Run(ctx, createTrackerMatchQuery, map[string]any{
							"tracker_id": tracker.ID,
							"insight_id": inRef.ID,
						})
						if err != nil {
							klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
							return nil, err
						}
						return result.Collect(ctx)
					})
				if err != nil {
					klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
					return err
				}
			}
		}
	}

	// TODO: pub/sub

	return nil
}

func (mh *MessageHandler) EntityResponseMessage(er *interfaces.EntityResponse) error {
	data, err := json.Marshal(er)
	if err != nil {
		klog.V(1).Infof("EntityResponseMessage json.Marshal failed. Err: %v\n", err)
		return err
	}

	// TODO: fix level
	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("EntityResponseMessage Object DUMP:\n%v\n", string(data))
	klog.V(3).Infof("-------------------------------\n\n")

	// write the object to the database
	ctx := context.Background()

	for _, entity := range er.Entities {

		// match id
		entityId := fmt.Sprintf("%s_%s_%s", entity.Type, entity.SubType, entity.Category)

		// entity
		_, err := mh.Session.ExecuteWrite(ctx,
			func(tx neo4j.ManagedTransaction) (any, error) {
				createEntitiesQuery := `
					MATCH (c:Conversation { conversationId: $conversation_id })
					MERGE (e:Entity { entityId: $entityId })
						ON CREATE SET
							e.lastAccessed = timestamp()
						ON MATCH SET
							e.lastAccessed = timestamp()
					SET e = { entityId: $entityId, type: $type, subType: $sub_type, category: $category, raw: $raw }
					MERGE (c)-[:ENTITY]-(e)
					`
				result, err := tx.Run(ctx, createEntitiesQuery, map[string]any{
					"conversation_id": mh.ConversationId,
					"entityId":        entityId,
					"type":            entity.Type,
					"sub_type":        entity.SubType,
					"category":        entity.Category,
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
			return err
		}

		// associate tracker to messages and insights
		for _, match := range entity.Matches {

			// generate a unique ID
			matchId := fmt.Sprintf("%s_%s", mh.ConversationId, entityId)

			// match
			_, err = mh.Session.ExecuteWrite(ctx,
				func(tx neo4j.ManagedTransaction) (any, error) {
					createEntitiesQuery := `
						MATCH (e:Entity { entityId: $entityId })
						MERGE (m:EntityMatch { matchId: $matchId })
							ON CREATE SET
								m.lastAccessed = timestamp()
							ON MATCH SET
								m.lastAccessed = timestamp()
						SET m = { matchId: $matchId, value: $value }
						MERGE (e)-[:ENTITY_MATCH_REF]-(m)
						`
					result, err := tx.Run(ctx, createEntitiesQuery, map[string]any{
						"entityId": entityId,
						"matchId":  matchId,
						"value":    match.DetectedValue,
					})
					if err != nil {
						klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
						return nil, err
					}
					return result.Collect(ctx)
				})
			if err != nil {
				klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
				return err
			}

			// message
			for _, msgRef := range match.MessageRefs {
				_, err = mh.Session.ExecuteWrite(ctx,
					func(tx neo4j.ManagedTransaction) (any, error) {
						createEntitiesQuery := `
							MATCH (e:EntityMatch { matchId: $matchId })
							MATCH (m:Message { messageId: $message_id })
							MERGE (e)-[:ENTITY_MESSAGE_REF]-(m)
							`
						result, err := tx.Run(ctx, createEntitiesQuery, map[string]any{
							"matchId":    matchId,
							"value":      match.DetectedValue,
							"message_id": msgRef.ID,
						})
						if err != nil {
							klog.V(1).Infof("neo4j.Run failed create conversation object. Err: %v\n", err)
							return nil, err
						}
						return result.Collect(ctx)
					})
				if err != nil {
					klog.V(1).Infof("neo4j.ExecuteWrite failed. Err: %v\n", err)
					return err
				}
			}
		}
	}

	// TODO: pub/sub

	return nil
}

func (mh *MessageHandler) UnhandledMessage(byMsg []byte) error {
	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("UnhandledMessage Object DUMP:\n%v\n", string(byMsg))
	klog.V(3).Infof("-------------------------------\n\n")
	return nil
}

func (mh *MessageHandler) HandleQuestion(insight *interfaces.Insight) error {
	return mh.handleInsight(insight)
}

func (mh *MessageHandler) HandleActionItem(insight *interfaces.Insight) error {
	return mh.handleInsight(insight)
}

func (mh *MessageHandler) HandleFollowUp(insight *interfaces.Insight) error {
	return mh.handleInsight(insight)
}

func (mh *MessageHandler) handleInsight(insight *interfaces.Insight) error {
	data, err := json.Marshal(insight)
	if err != nil {
		klog.V(1).Infof("handleInsight json.Marshal failed. Err: %v\n", err)
		return err
	}

	// TODO: fix level
	klog.V(3).Infof("\n\n-------------------------------\n")
	klog.V(3).Infof("handleInsight Object DUMP:\n%v\n", string(data))
	klog.V(3).Infof("-------------------------------\n\n")

	// write the object to the database
	ctx := context.Background()

	// if we need to do something with them
	// for records, message := range mr.Messages {
	_, err = mh.Session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			createInsightQuery := `
				MATCH (c:Conversation { conversationId: $conversation_id })
				MERGE (i:Insight { insightId: $insight_id })
					ON CREATE SET
						i.lastAccessed = timestamp()
					ON MATCH SET
						i.lastAccessed = timestamp()
				SET i = { insightId: $insight_id, type: $type, content: $content, assigneeId: $assignee_id, userId: $user_id, raw: $raw }
				MERGE (u:User { userId: $user_id })
					ON CREATE SET
						u.lastAccessed = timestamp()
					ON MATCH SET
						u.lastAccessed = timestamp()
				SET u = { realId: $user_real_id, userId: $user_id, name: $user_name, email: $user_id }
				MERGE (c)-[:INSIGHT]-(i)
				MERGE (i)-[:SPOKE]-(u)
				`
			result, err := tx.Run(ctx, createInsightQuery, map[string]any{
				"conversation_id": mh.ConversationId,
				"insight_id":      insight.ID,
				"type":            insight.Type,
				"content":         insight.Payload.Content,
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
		return err
	}

	// TODO: pub/sub

	return nil
}

func ConvertRootWordToSlice(words []interfaces.RootWord) []string {
	var arr []string
	for _, word := range words {
		arr = append(arr, word.Text)
	}
	return arr
}
