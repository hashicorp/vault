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

// DatabaseOptions represent all possible options to configure a Database.
type DatabaseOptions struct {
	ReadConcern    *readconcern.ReadConcern   // The read concern for operations in the database.
	WriteConcern   *writeconcern.WriteConcern // The write concern for operations in the database.
	ReadPreference *readpref.ReadPref         // The read preference for operations in the database.
	Registry       *bsoncodec.Registry        // The registry to be used to construct BSON encoders and decoders for the database.
}

// Database creates a new DatabaseOptions instance
func Database() *DatabaseOptions {
	return &DatabaseOptions{}
}

// SetReadConcern sets the read concern for the database.
func (d *DatabaseOptions) SetReadConcern(rc *readconcern.ReadConcern) *DatabaseOptions {
	d.ReadConcern = rc
	return d
}

// SetWriteConcern sets the write concern for the database.
func (d *DatabaseOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *DatabaseOptions {
	d.WriteConcern = wc
	return d
}

// SetReadPreference sets the read preference for the database.
func (d *DatabaseOptions) SetReadPreference(rp *readpref.ReadPref) *DatabaseOptions {
	d.ReadPreference = rp
	return d
}

// SetRegistry sets the bsoncodec Registry for the database.
func (d *DatabaseOptions) SetRegistry(r *bsoncodec.Registry) *DatabaseOptions {
	d.Registry = r
	return d
}

// MergeDatabaseOptions combines the *DatabaseOptions arguments into a single *DatabaseOptions in a last one wins
// fashion.
func MergeDatabaseOptions(opts ...*DatabaseOptions) *DatabaseOptions {
	d := Database()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ReadConcern != nil {
			d.ReadConcern = opt.ReadConcern
		}
		if opt.WriteConcern != nil {
			d.WriteConcern = opt.WriteConcern
		}
		if opt.ReadPreference != nil {
			d.ReadPreference = opt.ReadPreference
		}
		if opt.Registry != nil {
			d.Registry = opt.Registry
		}
	}

	return d
}
