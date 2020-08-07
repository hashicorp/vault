// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// UpdateOptions represents all possible options to the UpdateOne() and UpdateMany() functions.
type UpdateOptions struct {
	ArrayFilters             *ArrayFilters // A set of filters specifying to which array elements an update should apply
	BypassDocumentValidation *bool         // If true, allows the write to opt-out of document level validation
	Collation                *Collation    // Specifies a collation
	Upsert                   *bool         // When true, creates a new document if no document matches the query
}

// Update returns a pointer to a new UpdateOptions
func Update() *UpdateOptions {
	return &UpdateOptions{}
}

// SetArrayFilters specifies a set of filters specifying to which array elements an update should apply
// Valid for server versions >= 3.6.
func (uo *UpdateOptions) SetArrayFilters(af ArrayFilters) *UpdateOptions {
	uo.ArrayFilters = &af
	return uo
}

// SetBypassDocumentValidation allows the write to opt-out of document level validation.
// Valid for server versions >= 3.2. For servers < 3.2, this option is ignored.
func (uo *UpdateOptions) SetBypassDocumentValidation(b bool) *UpdateOptions {
	uo.BypassDocumentValidation = &b
	return uo
}

// SetCollation specifies a collation.
// Valid for server versions >= 3.4.
func (uo *UpdateOptions) SetCollation(c *Collation) *UpdateOptions {
	uo.Collation = c
	return uo
}

// SetUpsert allows the creation of a new document if not document matches the query
func (uo *UpdateOptions) SetUpsert(b bool) *UpdateOptions {
	uo.Upsert = &b
	return uo
}

// MergeUpdateOptions combines the argued UpdateOptions into a single UpdateOptions in a last-one-wins fashion
func MergeUpdateOptions(opts ...*UpdateOptions) *UpdateOptions {
	uOpts := Update()
	for _, uo := range opts {
		if uo == nil {
			continue
		}
		if uo.ArrayFilters != nil {
			uOpts.ArrayFilters = uo.ArrayFilters
		}
		if uo.BypassDocumentValidation != nil {
			uOpts.BypassDocumentValidation = uo.BypassDocumentValidation
		}
		if uo.Collation != nil {
			uOpts.Collation = uo.Collation
		}
		if uo.Upsert != nil {
			uOpts.Upsert = uo.Upsert
		}
	}

	return uOpts
}
