// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driverutil

// Operation Names should be sourced from the command reference documentation:
// https://www.mongodb.com/docs/manual/reference/command/
const (
	AbortTransactionOp  = "abortTransaction"  // AbortTransactionOp is the name for aborting a transaction
	AggregateOp         = "aggregate"         // AggregateOp is the name for aggregating
	CommitTransactionOp = "commitTransaction" // CommitTransactionOp is the name for committing a transaction
	CountOp             = "count"             // CountOp is the name for counting
	CreateOp            = "create"            // CreateOp is the name for creating
	CreateIndexesOp     = "createIndexes"     // CreateIndexesOp is the name for creating indexes
	DeleteOp            = "delete"            // DeleteOp is the name for deleting
	DistinctOp          = "distinct"          // DistinctOp is the name for distinct
	DropOp              = "drop"              // DropOp is the name for dropping
	DropDatabaseOp      = "dropDatabase"      // DropDatabaseOp is the name for dropping a database
	DropIndexesOp       = "dropIndexes"       // DropIndexesOp is the name for dropping indexes
	EndSessionsOp       = "endSessions"       // EndSessionsOp is the name for ending sessions
	FindAndModifyOp     = "findAndModify"     // FindAndModifyOp is the name for finding and modifying
	FindOp              = "find"              // FindOp is the name for finding
	InsertOp            = "insert"            // InsertOp is the name for inserting
	ListCollectionsOp   = "listCollections"   // ListCollectionsOp is the name for listing collections
	ListIndexesOp       = "listIndexes"       // ListIndexesOp is the name for listing indexes
	ListDatabasesOp     = "listDatabases"     // ListDatabasesOp is the name for listing databases
	UpdateOp            = "update"            // UpdateOp is the name for updating
)
