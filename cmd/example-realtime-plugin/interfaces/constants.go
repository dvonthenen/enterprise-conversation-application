// Copyright 2023 Enterprise Reference Implementation contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

const (
	// app specific message type
	AppSpecificMessageTypeScaffold string = "realtime_scaffold"

	// user/app level scaffold type
	UserScaffoldTypeInsight string = "realtime_scaffold_insight"
	UserScaffoldTypeTopic   string = "realtime_scaffold_topic"
	UserScaffoldTypeTracker string = "realtime_scaffold_tracker"
	UserScaffoldTypeEntity  string = "realtime_scaffold_entity"

	// app specific message type
	MessageNotFound string = "**MESSAGE NOT FOUND**"
)
