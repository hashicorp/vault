// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// CollectionOptions represent all possible options to configure a Collection.
type CollectionOptions struct {
	ReadConcern    *readconcern.ReadConcern   // The read concern for operations in the collection.
	WriteConcern   *writeconcern.WriteConcern // The write concern for operations in the collection.
	ReadPreference *readpref.ReadPref         // The read preference for operations in the collection.
	Registry       *bsoncodec.Registry        // The registry to be used to construct BSON encoders and decoders for the collection.
}

// Collection creates a new CollectionOptions instance
func Collection() *CollectionOptions {
	return &CollectionOptions{}
}

// SetReadConcern sets the read concern for the collection.
func (c *CollectionOptions) SetReadConcern(rc *readconcern.ReadConcern) *CollectionOptions {
	c.ReadConcern = rc
	return c
}

// SetWriteConcern sets the write concern for the collection.
func (c *CollectionOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *CollectionOptions {
	c.WriteConcern = wc
	return c
}

// SetReadPreference sets the read preference for the collection.
func (c *CollectionOptions) SetReadPreference(rp *readpref.ReadPref) *CollectionOptions {
	c.ReadPreference = rp
	return c
}

// SetRegistry sets the bsoncodec Registry for the collection.
func (c *CollectionOptions) SetRegistry(r *bsoncodec.Registry) *CollectionOptions {
	c.Registry = r
	return c
}

// MergeCollectionOptions combines the *CollectionOptions arguments into a single *CollectionOptions in a last one wins
// fashion.
func MergeCollectionOptions(opts ...*CollectionOptions) *CollectionOptions {
	c := Collection()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ReadConcern != nil {
			c.ReadConcern = opt.ReadConcern
		}
		if opt.WriteConcern != nil {
			c.WriteConcern = opt.WriteConcern
		}
		if opt.ReadPreference != nil {
			c.ReadPreference = opt.ReadPreference
		}
		if opt.Registry != nil {
			c.Registry = opt.Registry
		}
	}

	return c
}
