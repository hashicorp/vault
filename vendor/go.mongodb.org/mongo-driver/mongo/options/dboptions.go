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

// DatabaseOptions represents options that can be used to configure a Database.
type DatabaseOptions struct {
	// The read concern to use for operations executed on the Database. The default value is nil, which means that
	// the read concern of the client used to configure the Database will be used.
	ReadConcern *readconcern.ReadConcern

	// The write concern to use for operations executed on the Database. The default value is nil, which means that the
	// write concern of the client used to configure the Database will be used.
	WriteConcern *writeconcern.WriteConcern

	// The read preference to use for operations executed on the Database. The default value is nil, which means that
	// the read preference of the client used to configure the Database will be used.
	ReadPreference *readpref.ReadPref

	// The BSON registry to marshal and unmarshal documents for operations executed on the Database. The default value
	// is nil, which means that the registry of the client used to configure the Database will be used.
	Registry *bsoncodec.Registry
}

// Database creates a new DatabaseOptions instance.
func Database() *DatabaseOptions {
	return &DatabaseOptions{}
}

// SetReadConcern sets the value for the ReadConcern field.
func (d *DatabaseOptions) SetReadConcern(rc *readconcern.ReadConcern) *DatabaseOptions {
	d.ReadConcern = rc
	return d
}

// SetWriteConcern sets the value for the WriteConcern field.
func (d *DatabaseOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *DatabaseOptions {
	d.WriteConcern = wc
	return d
}

// SetReadPreference sets the value for the ReadPreference field.
func (d *DatabaseOptions) SetReadPreference(rp *readpref.ReadPref) *DatabaseOptions {
	d.ReadPreference = rp
	return d
}

// SetRegistry sets the value for the Registry field.
func (d *DatabaseOptions) SetRegistry(r *bsoncodec.Registry) *DatabaseOptions {
	d.Registry = r
	return d
}

// MergeDatabaseOptions combines the given DatabaseOptions instances into a single DatabaseOptions in a last-one-wins
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
