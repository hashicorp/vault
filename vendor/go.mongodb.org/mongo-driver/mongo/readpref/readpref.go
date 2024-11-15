// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package readpref defines read preferences for MongoDB queries.
package readpref // import "go.mongodb.org/mongo-driver/mongo/readpref"

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/tag"
)

var (
	errInvalidReadPreference = errors.New("can not specify tags, max staleness, or hedge with mode primary")
)

// Primary constructs a read preference with a PrimaryMode.
func Primary() *ReadPref {
	return &ReadPref{mode: PrimaryMode}
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
		if opt == nil {
			continue
		}
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
	hedgeEnabled    *bool
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

// HedgeEnabled returns whether or not hedged reads are enabled for this read preference. If this option was not
// specified during read preference construction, nil is returned.
func (r *ReadPref) HedgeEnabled() *bool {
	return r.hedgeEnabled
}

// String returns a human-readable description of the read preference.
func (r *ReadPref) String() string {
	var b bytes.Buffer
	b.WriteString(r.mode.String())
	delim := "("
	if r.maxStalenessSet {
		fmt.Fprintf(&b, "%smaxStaleness=%v", delim, r.maxStaleness)
		delim = " "
	}
	for _, tagSet := range r.tagSets {
		fmt.Fprintf(&b, "%stagSet=%s", delim, tagSet.String())
		delim = " "
	}
	if r.hedgeEnabled != nil {
		fmt.Fprintf(&b, "%shedgeEnabled=%v", delim, *r.hedgeEnabled)
		delim = " "
	}
	if delim != "(" {
		b.WriteString(")")
	}
	return b.String()
}
