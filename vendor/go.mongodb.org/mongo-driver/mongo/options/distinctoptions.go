// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// DistinctOptions represents options that can be used to configure a Distinct operation.
type DistinctOptions struct {
	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation. The default value is nil, which means that no comment will be included in the logs.
	Comment interface{}

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be
	// used in its place to control the amount of time that a single operation can run before returning an error.
	// MaxTime is ignored if Timeout is set on the client.
	MaxTime *time.Duration
}

// Distinct creates a new DistinctOptions instance.
func Distinct() *DistinctOptions {
	return &DistinctOptions{}
}

// SetCollation sets the value for the Collation field.
func (do *DistinctOptions) SetCollation(c *Collation) *DistinctOptions {
	do.Collation = c
	return do
}

// SetComment sets the value for the Comment field.
func (do *DistinctOptions) SetComment(comment interface{}) *DistinctOptions {
	do.Comment = comment
	return do
}

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (do *DistinctOptions) SetMaxTime(d time.Duration) *DistinctOptions {
	do.MaxTime = &d
	return do
}

// MergeDistinctOptions combines the given DistinctOptions instances into a single DistinctOptions in a last-one-wins
// fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeDistinctOptions(opts ...*DistinctOptions) *DistinctOptions {
	distinctOpts := Distinct()
	for _, do := range opts {
		if do == nil {
			continue
		}
		if do.Collation != nil {
			distinctOpts.Collation = do.Collation
		}
		if do.Comment != nil {
			distinctOpts.Comment = do.Comment
		}
		if do.MaxTime != nil {
			distinctOpts.MaxTime = do.MaxTime
		}
	}

	return distinctOpts
}
