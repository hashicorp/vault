// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// CountOptions represents options that can be used to configure a CountDocuments operation.
type CountOptions struct {
	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// TODO(GODRIVER-2386): CountOptions executor uses aggregation under the hood, which means this type has to be
	// TODO a string for now.  This can be replaced with `Comment interface{}` once 2386 is implemented.

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default is nil, which means that no comment will be included in the logs.
	Comment *string

	// The index to use for the aggregation. This should either be the index name as a string or the index specification
	// as a document. The driver will return an error if the hint parameter is a multi-key map. The default value is nil,
	// which means that no hint will be sent.
	Hint interface{}

	// The maximum number of documents to count. The default value is 0, which means that there is no limit and all
	// documents matching the filter will be counted.
	Limit *int64

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there is
	// no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used in
	// its place to control the amount of time that a single operation can run before returning an error. MaxTime is
	// ignored if Timeout is set on the client.
	MaxTime *time.Duration

	// The number of documents to skip before counting. The default value is 0.
	Skip *int64
}

// Count creates a new CountOptions instance.
func Count() *CountOptions {
	return &CountOptions{}
}

// SetCollation sets the value for the Collation field.
func (co *CountOptions) SetCollation(c *Collation) *CountOptions {
	co.Collation = c
	return co
}

// SetComment sets the value for the Comment field.
func (co *CountOptions) SetComment(c string) *CountOptions {
	co.Comment = &c
	return co
}

// SetHint sets the value for the Hint field.
func (co *CountOptions) SetHint(h interface{}) *CountOptions {
	co.Hint = h
	return co
}

// SetLimit sets the value for the Limit field.
func (co *CountOptions) SetLimit(i int64) *CountOptions {
	co.Limit = &i
	return co
}

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (co *CountOptions) SetMaxTime(d time.Duration) *CountOptions {
	co.MaxTime = &d
	return co
}

// SetSkip sets the value for the Skip field.
func (co *CountOptions) SetSkip(i int64) *CountOptions {
	co.Skip = &i
	return co
}

// MergeCountOptions combines the given CountOptions instances into a single CountOptions in a last-one-wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeCountOptions(opts ...*CountOptions) *CountOptions {
	countOpts := Count()
	for _, co := range opts {
		if co == nil {
			continue
		}
		if co.Collation != nil {
			countOpts.Collation = co.Collation
		}
		if co.Comment != nil {
			countOpts.Comment = co.Comment
		}
		if co.Hint != nil {
			countOpts.Hint = co.Hint
		}
		if co.Limit != nil {
			countOpts.Limit = co.Limit
		}
		if co.MaxTime != nil {
			countOpts.MaxTime = co.MaxTime
		}
		if co.Skip != nil {
			countOpts.Skip = co.Skip
		}
	}

	return countOpts
}
