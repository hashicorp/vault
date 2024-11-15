// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

// LegacyOperationKind indicates if an operation is a legacy find, getMore, or killCursors. This is used
// in Operation.Execute, which will create legacy OP_QUERY, OP_GET_MORE, or OP_KILL_CURSORS instead
// of sending them as a command.
type LegacyOperationKind uint

// These constants represent the three different kinds of legacy operations.
const (
	LegacyNone LegacyOperationKind = iota
	LegacyFind
	LegacyGetMore
	LegacyKillCursors
	LegacyListCollections
	LegacyListIndexes
	LegacyHandshake
)
