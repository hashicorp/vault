// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package readconcern // import "go.mongodb.org/mongo-driver/mongo/readconcern"

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ReadConcern for replica sets and replica set shards determines which data to return from a query.
type ReadConcern struct {
	level string
}

// Option is an option to provide when creating a ReadConcern.
type Option func(concern *ReadConcern)

// Level creates an option that sets the level of a ReadConcern.
func Level(level string) Option {
	return func(concern *ReadConcern) {
		concern.level = level
	}
}

// Local specifies that the query should return the instance’s most recent data.
func Local() *ReadConcern {
	return New(Level("local"))
}

// Majority specifies that the query should return the instance’s most recent data acknowledged as
// having been written to a majority of members in the replica set.
func Majority() *ReadConcern {
	return New(Level("majority"))
}

// Linearizable specifies that the query should return data that reflects all successful writes
// issued with a write concern of "majority" and acknowledged prior to the start of the read operation.
func Linearizable() *ReadConcern {
	return New(Level("linearizable"))
}

// Available specifies that the query should return data from the instance with no guarantee
// that the data has been written to a majority of the replica set members (i.e. may be rolled back).
func Available() *ReadConcern {
	return New(Level("available"))
}

// Snapshot is only available for operations within multi-document transactions.
func Snapshot() *ReadConcern {
	return New(Level("snapshot"))
}

// New constructs a new read concern from the given string.
func New(options ...Option) *ReadConcern {
	concern := &ReadConcern{}

	for _, option := range options {
		option(concern)
	}

	return concern
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (rc *ReadConcern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	var elems []byte

	if len(rc.level) > 0 {
		elems = bsoncore.AppendStringElement(elems, "level", rc.level)
	}

	return bsontype.EmbeddedDocument, bsoncore.BuildDocument(nil, elems), nil
}

// GetLevel returns the read concern level.
func (rc *ReadConcern) GetLevel() string {
	return rc.level
}
