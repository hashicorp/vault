// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// DefaultOrdered is the default order for a BulkWriteOptions struct created from BulkWrite.
var DefaultOrdered = true

// BulkWriteOptions represent all possible options for a bulkWrite operation.
type BulkWriteOptions struct {
	BypassDocumentValidation *bool // If true, allows the write to opt out of document-level validation.
	Ordered                  *bool // If true, when a write fails, return without performing remaining writes. Defaults to true.
}

// BulkWrite creates a new *BulkWriteOptions
func BulkWrite() *BulkWriteOptions {
	return &BulkWriteOptions{
		Ordered: &DefaultOrdered,
	}
}

// SetOrdered configures the ordered option. If true, when a write fails, the function will return without attempting
// remaining writes. Defaults to true.
func (b *BulkWriteOptions) SetOrdered(ordered bool) *BulkWriteOptions {
	b.Ordered = &ordered
	return b
}

// SetBypassDocumentValidation specifies if the write should opt out of document-level validation.
// Valid for server versions >= 3.2. For servers < 3.2, this option is ignored.
func (b *BulkWriteOptions) SetBypassDocumentValidation(bypass bool) *BulkWriteOptions {
	b.BypassDocumentValidation = &bypass
	return b
}

// MergeBulkWriteOptions combines the given *BulkWriteOptions into a single *BulkWriteOptions in a last one wins fashion.
func MergeBulkWriteOptions(opts ...*BulkWriteOptions) *BulkWriteOptions {
	b := BulkWrite()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Ordered != nil {
			b.Ordered = opt.Ordered
		}
		if opt.BypassDocumentValidation != nil {
			b.BypassDocumentValidation = opt.BypassDocumentValidation
		}
	}

	return b
}
