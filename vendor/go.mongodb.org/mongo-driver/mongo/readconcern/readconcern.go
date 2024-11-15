// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package readconcern defines read concerns for MongoDB operations.
//
// For more information about MongoDB read concerns, see
// https://www.mongodb.com/docs/manual/reference/read-concern/
package readconcern // import "go.mongodb.org/mongo-driver/mongo/readconcern"

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// A ReadConcern defines a MongoDB read concern, which allows you to control the consistency and
// isolation properties of the data read from replica sets and replica set shards.
//
// For more information about MongoDB read concerns, see
// https://www.mongodb.com/docs/manual/reference/read-concern/
type ReadConcern struct {
	Level string
}

// Option is an option to provide when creating a ReadConcern.
//
// Deprecated: Use the ReadConcern literal declaration instead. For example:
//
//	&readconcern.ReadConcern{Level: "local"}
type Option func(concern *ReadConcern)

// Level creates an option that sets the level of a ReadConcern.
//
// Deprecated: Use the ReadConcern literal declaration instead. For example:
//
//	&readconcern.ReadConcern{Level: "local"}
func Level(level string) Option {
	return func(concern *ReadConcern) {
		concern.Level = level
	}
}

// Local returns a ReadConcern that requests data from the instance with no guarantee that the data
// has been written to a majority of the replica set members (i.e. may be rolled back).
//
// For more information about read concern "local", see
// https://www.mongodb.com/docs/manual/reference/read-concern-local/
func Local() *ReadConcern {
	return New(Level("local"))
}

// Majority returns a ReadConcern that requests data that has been acknowledged by a majority of the
// replica set members (i.e. the documents read are durable and guaranteed not to roll back).
//
// For more information about read concern "majority", see
// https://www.mongodb.com/docs/manual/reference/read-concern-majority/
func Majority() *ReadConcern {
	return New(Level("majority"))
}

// Linearizable returns a ReadConcern that requests data that reflects all successful
// majority-acknowledged writes that completed prior to the start of the read operation.
//
// For more information about read concern "linearizable", see
// https://www.mongodb.com/docs/manual/reference/read-concern-linearizable/
func Linearizable() *ReadConcern {
	return New(Level("linearizable"))
}

// Available returns a ReadConcern that requests data from an instance with no guarantee that the
// data has been written to a majority of the replica set members (i.e. may be rolled back).
//
// For more information about read concern "available", see
// https://www.mongodb.com/docs/manual/reference/read-concern-available/
func Available() *ReadConcern {
	return New(Level("available"))
}

// Snapshot returns a ReadConcern that requests majority-committed data as it appears across shards
// from a specific single point in time in the recent past.
//
// For more information about read concern "snapshot", see
// https://www.mongodb.com/docs/manual/reference/read-concern-snapshot/
func Snapshot() *ReadConcern {
	return New(Level("snapshot"))
}

// New constructs a new read concern from the given string.
//
// Deprecated: Use the ReadConcern literal declaration instead. For example:
//
//	&readconcern.ReadConcern{Level: "local"}
func New(options ...Option) *ReadConcern {
	concern := &ReadConcern{}

	for _, option := range options {
		option(concern)
	}

	return concern
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
//
// Deprecated: Marshaling a ReadConcern to BSON will not be supported in Go Driver 2.0.
func (rc *ReadConcern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if rc == nil {
		return 0, nil, errors.New("cannot marshal nil ReadConcern")
	}

	var elems []byte

	if len(rc.Level) > 0 {
		elems = bsoncore.AppendStringElement(elems, "level", rc.Level)
	}

	return bsontype.EmbeddedDocument, bsoncore.BuildDocument(nil, elems), nil
}

// GetLevel returns the read concern level.
//
// Deprecated: Use the ReadConcern.Level field instead.
func (rc *ReadConcern) GetLevel() string {
	return rc.Level
}
