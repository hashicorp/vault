// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// ReplaceOptions represents all possible options to the ReplaceOne() function.
type ReplaceOptions struct {
	BypassDocumentValidation *bool      // If true, allows the write to opt-out of document level validation
	Collation                *Collation // Specifies a collation
	Upsert                   *bool      // When true, creates a new document if no document matches the query
}

// Replace returns a pointer to a new ReplaceOptions
func Replace() *ReplaceOptions {
	return &ReplaceOptions{}
}

// SetBypassDocumentValidation allows the write to opt-out of document level validation.
// Valid for server versions >= 3.2. For servers < 3.2, this option is ignored.
func (ro *ReplaceOptions) SetBypassDocumentValidation(b bool) *ReplaceOptions {
	ro.BypassDocumentValidation = &b
	return ro
}

// SetCollation specifies a collation.
// Valid for servers >= 3.4
func (ro *ReplaceOptions) SetCollation(c *Collation) *ReplaceOptions {
	ro.Collation = c
	return ro
}

// SetUpsert allows the creation of a new document if not document matches the query
func (ro *ReplaceOptions) SetUpsert(b bool) *ReplaceOptions {
	ro.Upsert = &b
	return ro
}

// MergeReplaceOptions combines the argued ReplaceOptions into a single ReplaceOptions in a last-one-wins fashion
func MergeReplaceOptions(opts ...*ReplaceOptions) *ReplaceOptions {
	rOpts := Replace()
	for _, ro := range opts {
		if ro == nil {
			continue
		}
		if ro.BypassDocumentValidation != nil {
			rOpts.BypassDocumentValidation = ro.BypassDocumentValidation
		}
		if ro.Collation != nil {
			rOpts.Collation = ro.Collation
		}
		if ro.Upsert != nil {
			rOpts.Upsert = ro.Upsert
		}
	}

	return rOpts
}
