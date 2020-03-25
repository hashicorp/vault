// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// DistinctOptions represents all possible options to the Distinct() function.
type DistinctOptions struct {
	Collation *Collation     // Specifies a collation
	MaxTime   *time.Duration // The maximum amount of time to allow the operation to run
}

// Distinct returns a pointer to a new DistinctOptions
func Distinct() *DistinctOptions {
	return &DistinctOptions{}
}

// SetCollation specifies a collation
// Valid for server versions >= 3.4
func (do *DistinctOptions) SetCollation(c *Collation) *DistinctOptions {
	do.Collation = c
	return do
}

// SetMaxTime specifies the maximum amount of time to allow the operation to run
func (do *DistinctOptions) SetMaxTime(d time.Duration) *DistinctOptions {
	do.MaxTime = &d
	return do
}

// MergeDistinctOptions combines the argued DistinctOptions into a single DistinctOptions in a last-one-wins fashion
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
