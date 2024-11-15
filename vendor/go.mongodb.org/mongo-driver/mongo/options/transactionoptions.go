// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// TransactionOptions represents options that can be used to configure a transaction.
type TransactionOptions struct {
	// The read concern for operations in the transaction. The default value is nil, which means that the default
	// read concern of the session used to start the transaction will be used.
	ReadConcern *readconcern.ReadConcern

	// The read preference for operations in the transaction. The default value is nil, which means that the default
	// read preference of the session used to start the transaction will be used.
	ReadPreference *readpref.ReadPref

	// The write concern for operations in the transaction. The default value is nil, which means that the default
	// write concern of the session used to start the transaction will be used.
	WriteConcern *writeconcern.WriteConcern

	// The default maximum amount of time that a CommitTransaction operation executed in the session can run on the
	// server. The default value is nil, meaning that there is no time limit for execution.

	// The maximum amount of time that a CommitTransaction operation can executed in the transaction can run on the
	// server. The default value is nil, which means that the default maximum commit time of the session used to
	// start the transaction will be used.
	//
	// NOTE(benjirewis): MaxCommitTime will be deprecated in a future release. The more general Timeout option may
	// be used in its place to control the amount of time that a single operation can run before returning an error.
	// MaxCommitTime is ignored if Timeout is set on the client.
	MaxCommitTime *time.Duration
}

// Transaction creates a new TransactionOptions instance.
func Transaction() *TransactionOptions {
	return &TransactionOptions{}
}

// SetReadConcern sets the value for the ReadConcern field.
func (t *TransactionOptions) SetReadConcern(rc *readconcern.ReadConcern) *TransactionOptions {
	t.ReadConcern = rc
	return t
}

// SetReadPreference sets the value for the ReadPreference field.
func (t *TransactionOptions) SetReadPreference(rp *readpref.ReadPref) *TransactionOptions {
	t.ReadPreference = rp
	return t
}

// SetWriteConcern sets the value for the WriteConcern field.
func (t *TransactionOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *TransactionOptions {
	t.WriteConcern = wc
	return t
}

// SetMaxCommitTime sets the value for the MaxCommitTime field.
//
// NOTE(benjirewis): MaxCommitTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can run before
// returning an error. MaxCommitTime is ignored if Timeout is set on the client.
func (t *TransactionOptions) SetMaxCommitTime(mct *time.Duration) *TransactionOptions {
	t.MaxCommitTime = mct
	return t
}

// MergeTransactionOptions combines the given TransactionOptions instances into a single TransactionOptions in a
// last-one-wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeTransactionOptions(opts ...*TransactionOptions) *TransactionOptions {
	t := Transaction()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ReadConcern != nil {
			t.ReadConcern = opt.ReadConcern
		}
		if opt.ReadPreference != nil {
			t.ReadPreference = opt.ReadPreference
		}
		if opt.WriteConcern != nil {
			t.WriteConcern = opt.WriteConcern
		}
		if opt.MaxCommitTime != nil {
			t.MaxCommitTime = opt.MaxCommitTime
		}
	}

	return t
}
