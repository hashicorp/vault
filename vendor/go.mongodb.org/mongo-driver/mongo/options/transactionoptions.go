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

// TransactionOptions represents all possible options for starting a transaction.
type TransactionOptions struct {
	ReadConcern    *readconcern.ReadConcern   // The read concern for the transaction. Defaults to the session's read concern.
	ReadPreference *readpref.ReadPref         // The read preference for the transaction. Defaults to the session's read preference.
	WriteConcern   *writeconcern.WriteConcern // The write concern for the transaction. Defaults to the session's write concern.
	MaxCommitTime  *time.Duration             // The maximum amount of time to allow a single commitTransaction command to run.
}

// Transaction creates a new *TransactionOptions
func Transaction() *TransactionOptions {
	return &TransactionOptions{}
}

// SetReadConcern sets the read concern for the transaction.
func (t *TransactionOptions) SetReadConcern(rc *readconcern.ReadConcern) *TransactionOptions {
	t.ReadConcern = rc
	return t
}

// SetReadPreference sets the read preference for the transaction.
func (t *TransactionOptions) SetReadPreference(rp *readpref.ReadPref) *TransactionOptions {
	t.ReadPreference = rp
	return t
}

// SetWriteConcern sets the write concern for the transaction.
func (t *TransactionOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *TransactionOptions {
	t.WriteConcern = wc
	return t
}

// SetMaxCommitTime sets the max commit time for the transaction.
func (t *TransactionOptions) SetMaxCommitTime(mct *time.Duration) *TransactionOptions {
	t.MaxCommitTime = mct
	return t
}

// MergeTransactionOptions combines the given *TransactionOptions into a single *TransactionOptions in a last one wins
// fashion.
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
