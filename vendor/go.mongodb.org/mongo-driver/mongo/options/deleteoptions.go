// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// DeleteOptions represents all possible options to the DeleteOne() and DeleteMany() functions.
type DeleteOptions struct {
	Collation *Collation // Specifies a collation
}

// Delete returns a pointer to a new DeleteOptions
func Delete() *DeleteOptions {
	return &DeleteOptions{}
}

// SetCollation specifies a collation
// Valid for servers >= 3.4.
func (do *DeleteOptions) SetCollation(c *Collation) *DeleteOptions {
	do.Collation = c
	return do
}

// MergeDeleteOptions combines the argued DeleteOptions into a single DeleteOptions in a last-one-wins fashion
func MergeDeleteOptions(opts ...*DeleteOptions) *DeleteOptions {
	dOpts := Delete()
	for _, do := range opts {
		if do == nil {
			continue
		}
		if do.Collation != nil {
			dOpts.Collation = do.Collation
		}
	}

	return dOpts
}
