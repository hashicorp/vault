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

// DefaultCausalConsistency is the default value for the CausalConsistency option.
var DefaultCausalConsistency = true

// SessionOptions represents options that can be used to configure a Session.
type SessionOptions struct {
	// If true, causal consistency will be enabled for the session. This option cannot be set to true if Snapshot is
	// set to true. The default value is true unless Snapshot is set to true. See
	// https://www.mongodb.com/docs/manual/core/read-isolation-consistency-recency/#sessions for more information.
	CausalConsistency *bool

	// The default read concern for transactions started in the session. The default value is nil, which means that
	// the read concern of the client used to start the session will be used.
	DefaultReadConcern *readconcern.ReadConcern

	// The default read preference for transactions started in the session. The default value is nil, which means that
	// the read preference of the client used to start the session will be used.
	DefaultReadPreference *readpref.ReadPref

	// The default write concern for transactions started in the session. The default value is nil, which means that
	// the write concern of the client used to start the session will be used.
	DefaultWriteConcern *writeconcern.WriteConcern

	// The default maximum amount of time that a CommitTransaction operation executed in the session can run on the
	// server. The default value is nil, which means that that there is no time limit for execution.
	//
	// NOTE(benjirewis): DefaultMaxCommitTime will be deprecated in a future release. The more general Timeout option
	// may be used in its place to control the amount of time that a single operation can run before returning an
	// error. DefaultMaxCommitTime is ignored if Timeout is set on the client.
	DefaultMaxCommitTime *time.Duration

	// If true, all read operations performed with this session will be read from the same snapshot. This option cannot
	// be set to true if CausalConsistency is set to true. Transactions and write operations are not allowed on
	// snapshot sessions and will error. The default value is false.
	Snapshot *bool
}

// Session creates a new SessionOptions instance.
func Session() *SessionOptions {
	return &SessionOptions{}
}

// SetCausalConsistency sets the value for the CausalConsistency field.
func (s *SessionOptions) SetCausalConsistency(b bool) *SessionOptions {
	s.CausalConsistency = &b
	return s
}

// SetDefaultReadConcern sets the value for the DefaultReadConcern field.
func (s *SessionOptions) SetDefaultReadConcern(rc *readconcern.ReadConcern) *SessionOptions {
	s.DefaultReadConcern = rc
	return s
}

// SetDefaultReadPreference sets the value for the DefaultReadPreference field.
func (s *SessionOptions) SetDefaultReadPreference(rp *readpref.ReadPref) *SessionOptions {
	s.DefaultReadPreference = rp
	return s
}

// SetDefaultWriteConcern sets the value for the DefaultWriteConcern field.
func (s *SessionOptions) SetDefaultWriteConcern(wc *writeconcern.WriteConcern) *SessionOptions {
	s.DefaultWriteConcern = wc
	return s
}

// SetDefaultMaxCommitTime sets the value for the DefaultMaxCommitTime field.
//
// NOTE(benjirewis): DefaultMaxCommitTime will be deprecated in a future release. The more
// general Timeout option may be used in its place to control the amount of time that a
// single operation can run before returning an error. DefaultMaxCommitTime is ignored if
// Timeout is set on the client.
func (s *SessionOptions) SetDefaultMaxCommitTime(mct *time.Duration) *SessionOptions {
	s.DefaultMaxCommitTime = mct
	return s
}

// SetSnapshot sets the value for the Snapshot field.
func (s *SessionOptions) SetSnapshot(b bool) *SessionOptions {
	s.Snapshot = &b
	return s
}

// MergeSessionOptions combines the given SessionOptions instances into a single SessionOptions in a last-one-wins
// fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeSessionOptions(opts ...*SessionOptions) *SessionOptions {
	s := Session()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.CausalConsistency != nil {
			s.CausalConsistency = opt.CausalConsistency
		}
		if opt.DefaultReadConcern != nil {
			s.DefaultReadConcern = opt.DefaultReadConcern
		}
		if opt.DefaultReadPreference != nil {
			s.DefaultReadPreference = opt.DefaultReadPreference
		}
		if opt.DefaultWriteConcern != nil {
			s.DefaultWriteConcern = opt.DefaultWriteConcern
		}
		if opt.DefaultMaxCommitTime != nil {
			s.DefaultMaxCommitTime = opt.DefaultMaxCommitTime
		}
		if opt.Snapshot != nil {
			s.Snapshot = opt.Snapshot
		}
	}
	if s.CausalConsistency == nil && (s.Snapshot == nil || !*s.Snapshot) {
		s.CausalConsistency = &DefaultCausalConsistency
	}

	return s
}
