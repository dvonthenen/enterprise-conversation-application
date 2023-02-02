// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import (
	"strings"
)

var (
	indexReplace = map[string]string{
		"#conversation_index#": DatabaseIndexConversation,
		"#message_index#":      DatabaseIndexMessage,
		"#user_index#":         DatabaseIndexUser,
		"#topic_index#":        DatabaseIndexTopic,
		"#tracker_index#":      DatabaseIndexTracker,
		"#insight_index#":      DatabaseIndexInsight,
		"#entity_index#":       DatabaseIndexEntity,
		"#match_index#":        DatabaseIndexEntityMatch,
	}
)

func ReplaceIndexes(str string) string {
	for key, value := range indexReplace {
		str = strings.ReplaceAll(str, key, value)
	}
	return str
}
