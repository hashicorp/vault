// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package readpref

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/tag"
)

// ErrInvalidTagSet indicates that an invalid set of tags was specified.
var ErrInvalidTagSet = errors.New("an even number of tags must be specified")

// Option configures a read preference
type Option func(*ReadPref) error

// WithMaxStaleness sets the maximum staleness a
// server is allowed.
func WithMaxStaleness(ms time.Duration) Option {
	return func(rp *ReadPref) error {
		rp.maxStaleness = ms
		rp.maxStalenessSet = true
		return nil
	}
}

// WithTags specifies a single tag set used to match replica set members. If no members match the
// tag set, read operations will return an error. To avoid errors if no members match the tag set, use
// [WithTagSets] and include an empty tag set as the last tag set in the list.
//
// The last call to [WithTags] or [WithTagSets] overrides all previous calls to either method.
//
// For more information about read preference tags, see
// https://www.mongodb.com/docs/manual/core/read-preference-tags/
func WithTags(tags ...string) Option {
	return func(rp *ReadPref) error {
		length := len(tags)
		if length < 2 || length%2 != 0 {
			return ErrInvalidTagSet
		}

		tagset := make(tag.Set, 0, length/2)

		for i := 1; i < length; i += 2 {
			tagset = append(tagset, tag.Tag{Name: tags[i-1], Value: tags[i]})
		}

		return WithTagSets(tagset)(rp)
	}
}

// WithTagSets specifies a list of tag sets used to match replica set members. If the list contains
// multiple tag sets, members are matched against each tag set in succession until a match is found.
// Once a match is found, the remaining tag sets are ignored. If no members match any of the tag
// sets, the read operation returns with an error. To avoid an error if no members match any of the
// tag sets, include an empty tag set as the last tag set in the list.
//
// The last call to [WithTags] or [WithTagSets] overrides all previous calls to either method.
//
// For more information about read preference tags, see
// https://www.mongodb.com/docs/manual/core/read-preference-tags/
func WithTagSets(tagSets ...tag.Set) Option {
	return func(rp *ReadPref) error {
		rp.tagSets = tagSets
		return nil
	}
}

// WithHedgeEnabled specifies whether or not hedged reads should be enabled in the server. This feature requires MongoDB
// server version 4.4 or higher. For more information about hedged reads, see
// https://www.mongodb.com/docs/manual/core/sharded-cluster-query-router/#mongos-hedged-reads. If not specified, the default
// is to not send a value to the server, which will result in the server defaults being used.
func WithHedgeEnabled(hedgeEnabled bool) Option {
	return func(rp *ReadPref) error {
		rp.hedgeEnabled = &hedgeEnabled
		return nil
	}
}
