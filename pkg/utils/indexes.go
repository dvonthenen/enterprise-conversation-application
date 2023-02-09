// Copyright 2023 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package utils

import (
	"strings"

	shared "github.com/dvonthenen/enterprise-reference-implementation/pkg/shared"
)

var (
	indexReplace = map[string]string{
		"#conversation_index#": shared.DatabaseIndexConversation,
		"#message_index#":      shared.DatabaseIndexMessage,
		"#user_index#":         shared.DatabaseIndexUser,
		"#topic_index#":        shared.DatabaseIndexTopic,
		"#tracker_index#":      shared.DatabaseIndexTracker,
		"#insight_index#":      shared.DatabaseIndexInsight,
		"#entity_index#":       shared.DatabaseIndexEntity,
		"#match_index#":        shared.DatabaseIndexEntityMatch,
	}
)

func ReplaceIndexes(str string) string {
	for key, value := range indexReplace {
		str = strings.ReplaceAll(str, key, value)
	}
	return str
}
