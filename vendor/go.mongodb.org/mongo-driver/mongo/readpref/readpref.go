// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package readpref // import "go.mongodb.org/mongo-driver/mongo/readpref"

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/tag"
)

var (
	errInvalidReadPreference = errors.New("can not specify tags or max staleness on primary")
)

var primary = ReadPref{mode: PrimaryMode}

// Primary constructs a read preference with a PrimaryMode.
func Primary() *ReadPref {
	return &primary
}

// PrimaryPreferred constructs a read preference with a PrimaryPreferredMode.
func PrimaryPreferred(opts ...Option) *ReadPref {
	// New only returns an error with a mode of Primary
	rp, _ := New(PrimaryPreferredMode, opts...)
	return rp
}

// SecondaryPreferred constructs a read preference with a SecondaryPreferredMode.
func SecondaryPreferred(opts ...Option) *ReadPref {
	// New only returns an error with a mode of Primary
	rp, _ := New(SecondaryPreferredMode, opts...)
	return rp
}

// Secondary constructs a read preference with a SecondaryMode.
func Secondary(opts ...Option) *ReadPref {
	// New only returns an error with a mode of Primary
	rp, _ := New(SecondaryMode, opts...)
	return rp
}

// Nearest constructs a read preference with a NearestMode.
func Nearest(opts ...Option) *ReadPref {
	// New only returns an error with a mode of Primary
	rp, _ := New(NearestMode, opts...)
	return rp
}

// New creates a new ReadPref.
func New(mode Mode, opts ...Option) (*ReadPref, error) {
	rp := &ReadPref{
		mode: mode,
	}

	if mode == PrimaryMode && len(opts) != 0 {
		return nil, errInvalidReadPreference
	}

	for _, opt := range opts {
		err := opt(rp)
		if err != nil {
			return nil, err
		}
	}

	return rp, nil
}

// ReadPref determines which servers are considered suitable for read operations.
type ReadPref struct {
	maxStaleness    time.Duration
	maxStalenessSet bool
	mode            Mode
	tagSets         []tag.Set
}

// MaxStaleness is the maximum amount of time to allow
// a server to be considered eligible for selection. The
// second return value indicates if this value has been set.
func (r *ReadPref) MaxStaleness() (time.Duration, bool) {
	return r.maxStaleness, r.maxStalenessSet
}

// Mode indicates the mode of the read preference.
func (r *ReadPref) Mode() Mode {
	return r.mode
}

// TagSets are multiple tag sets indicating
// which servers should be considered.
func (r *ReadPref) TagSets() []tag.Set {
	return r.tagSets
}
