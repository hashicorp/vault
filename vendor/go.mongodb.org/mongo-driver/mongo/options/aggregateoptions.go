// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// AggregateOptions represents all possible options to the Aggregate() function.
type AggregateOptions struct {
	AllowDiskUse             *bool          // Enables writing to temporary files. When set to true, aggregation stages can write data to the _tmp subdirectory in the dbPath directory
	BatchSize                *int32         // The number of documents to return per batch
	BypassDocumentValidation *bool          // If true, allows the write to opt-out of document level validation. This only applies when the $out stage is specified
	Collation                *Collation     // Specifies a collation
	MaxTime                  *time.Duration // The maximum amount of time to allow the query to run
	MaxAwaitTime             *time.Duration // The maximum amount of time for the server to wait on new documents to satisfy a tailable cursor query
	Comment                  *string        // Enables users to specify an arbitrary string to help trace the operation through the database profiler, currentOp and logs.
	Hint                     interface{}    // The index to use for the aggregation. The hint does not apply to $lookup and $graphLookup stages
}

// Aggregate returns a pointer to a new AggregateOptions
func Aggregate() *AggregateOptions {
	return &AggregateOptions{}
}

// SetAllowDiskUse enables writing to temporary files. When set to true,
// aggregation stages can write data to the _tmp subdirectory in the
// dbPath directory
func (ao *AggregateOptions) SetAllowDiskUse(b bool) *AggregateOptions {
	ao.AllowDiskUse = &b
	return ao
}

// SetBatchSize specifies the number of documents to return per batch
func (ao *AggregateOptions) SetBatchSize(i int32) *AggregateOptions {
	ao.BatchSize = &i
	return ao
}

// SetBypassDocumentValidation allows the write to opt-out of document level
// validation. This only applies when the $out stage is specified
// Valid for server versions >= 3.2. For servers < 3.2, this option is ignored.
func (ao *AggregateOptions) SetBypassDocumentValidation(b bool) *AggregateOptions {
	ao.BypassDocumentValidation = &b
	return ao
}

// SetCollation specifies a collation.
// Valid for server versions >= 3.4
func (ao *AggregateOptions) SetCollation(c *Collation) *AggregateOptions {
	ao.Collation = c
	return ao
}

// SetMaxTime specifies the maximum amount of time to allow the query to run
func (ao *AggregateOptions) SetMaxTime(d time.Duration) *AggregateOptions {
	ao.MaxTime = &d
	return ao
}

// SetMaxAwaitTime specifies the maximum amount of time for the server to
// wait on new documents to satisfy a tailable cursor query
// For servers < 3.2, this option is ignored
func (ao *AggregateOptions) SetMaxAwaitTime(d time.Duration) *AggregateOptions {
	ao.MaxAwaitTime = &d
	return ao
}

// SetComment enables users to specify an arbitrary string to help trace the
// operation through the database profiler, currentOp and logs.
func (ao *AggregateOptions) SetComment(s string) *AggregateOptions {
	ao.Comment = &s
	return ao
}

// SetHint specifies the index to use for the aggregation. The hint does not
// apply to $lookup and $graphLookup stages
func (ao *AggregateOptions) SetHint(h interface{}) *AggregateOptions {
	ao.Hint = h
	return ao
}

// MergeAggregateOptions combines the argued AggregateOptions into a single AggregateOptions in a last-one-wins fashion
func MergeAggregateOptions(opts ...*AggregateOptions) *AggregateOptions {
	aggOpts := Aggregate()
	for _, ao := range opts {
		if ao == nil {
			continue
		}
		if ao.AllowDiskUse != nil {
			aggOpts.AllowDiskUse = ao.AllowDiskUse
		}
		if ao.BatchSize != nil {
			aggOpts.BatchSize = ao.BatchSize
		}
		if ao.BypassDocumentValidation != nil {
			aggOpts.BypassDocumentValidation = ao.BypassDocumentValidation
		}
		if ao.Collation != nil {
			aggOpts.Collation = ao.Collation
		}
		if ao.MaxTime != nil {
			aggOpts.MaxTime = ao.MaxTime
		}
		if ao.MaxAwaitTime != nil {
			aggOpts.MaxAwaitTime = ao.MaxAwaitTime
		}
		if ao.Comment != nil {
			aggOpts.Comment = ao.Comment
		}
		if ao.Hint != nil {
			aggOpts.Hint = ao.Hint
		}
	}

	return aggOpts
}
