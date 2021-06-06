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

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
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

// SetMaxTime sets the value for the MaxTime field.
func (do *DistinctOptions) SetMaxTime(d time.Duration) *DistinctOptions {
	do.MaxTime = &d
	return do
}

// MergeDistinctOptions combines the given DistinctOptions instances into a single DistinctOptions in a last-one-wins
// fashion.
func MergeDistinctOptions(opts ...*DistinctOptions) *DistinctOptions {
	distinctOpts := Distinct()
	for _, do := range opts {
		if do == nil {
			continue
		}
		if do.Collation != nil {
			distinctOpts.Collation = do.Collation
		}
		if do.MaxTime != nil {
			distinctOpts.MaxTime = do.MaxTime
		}
	}

	return distinctOpts
}
