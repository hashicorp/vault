// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// AggregateOptions represents options that can be used to configure an Aggregate operation.
type AggregateOptions struct {
	// If true, the operation can write to temporary files in the _tmp subdirectory of the database directory path on
	// the server. The default value is false.
	AllowDiskUse *bool

	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://docs.mongodb.com/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool

	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	MaxTime *time.Duration

	// The maximum amount of time that the server should wait for new documents to satisfy a tailable cursor query.
	// This option is only valid for MongoDB versions >= 3.2 and is ignored for previous server versions.
	MaxAwaitTime *time.Duration

	// A string that will be included in server logs, profiling logs, and currentOp queries to help trace the operation.
	// The default is the empty string, which means that no comment will be included in the logs.
	Comment *string

	// The index to use for the aggregation. This should either be the index name as a string or the index specification
	// as a document. The hint does not apply to $lookup and $graphLookup aggregation stages. The default value is nil,
	// which means that no hint will be sent.
	Hint interface{}
}

// Aggregate creates a new AggregateOptions instance.
func Aggregate() *AggregateOptions {
	return &AggregateOptions{}
}

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (ao *AggregateOptions) SetAllowDiskUse(b bool) *AggregateOptions {
	ao.AllowDiskUse = &b
	return ao
}

// SetBatchSize sets the value for the BatchSize field.
func (ao *AggregateOptions) SetBatchSize(i int32) *AggregateOptions {
	ao.BatchSize = &i
	return ao
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (ao *AggregateOptions) SetBypassDocumentValidation(b bool) *AggregateOptions {
	ao.BypassDocumentValidation = &b
	return ao
}

// SetCollation sets the value for the Collation field.
func (ao *AggregateOptions) SetCollation(c *Collation) *AggregateOptions {
	ao.Collation = c
	return ao
}

// SetMaxTime sets the value for the MaxTime field.
func (ao *AggregateOptions) SetMaxTime(d time.Duration) *AggregateOptions {
	ao.MaxTime = &d
	return ao
}

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (ao *AggregateOptions) SetMaxAwaitTime(d time.Duration) *AggregateOptions {
	ao.MaxAwaitTime = &d
	return ao
}

// SetComment sets the value for the Comment field.
func (ao *AggregateOptions) SetComment(s string) *AggregateOptions {
	ao.Comment = &s
	return ao
}

// SetHint sets the value for the Hint field.
func (ao *AggregateOptions) SetHint(h interface{}) *AggregateOptions {
	ao.Hint = h
	return ao
}

// MergeAggregateOptions combines the given AggregateOptions instances into a single AggregateOptions in a last-one-wins
// fashion.
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
