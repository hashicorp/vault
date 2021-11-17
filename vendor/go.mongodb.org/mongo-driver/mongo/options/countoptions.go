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

	// The index to use for the aggregation. This should either be the index name as a string or the index specification
	// as a document. The driver will return an error if the hint parameter is a multi-key map. The default value is nil,
	// which means that no hint will be sent.
	Hint interface{}

	// The maximum number of documents to count. The default value is 0, which means that there is no limit and all
	// documents matching the filter will be counted.
	Limit *int64

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there is
	// no time limit for query execution.
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
func MergeCountOptions(opts ...*CountOptions) *CountOptions {
	countOpts := Count()
	for _, co := range opts {
		if co == nil {
			continue
		}
		if co.Collation != nil {
			countOpts.Collation = co.Collation
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
