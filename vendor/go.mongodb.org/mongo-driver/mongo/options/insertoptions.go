// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// InsertOneOptions represents all possible options to the InsertOne() function.
type InsertOneOptions struct {
	BypassDocumentValidation *bool // If true, allows the write to opt-out of document level validation
}

// InsertOne returns a pointer to a new InsertOneOptions
func InsertOne() *InsertOneOptions {
	return &InsertOneOptions{}
}

// SetBypassDocumentValidation allows the write to opt-out of document level validation.
// Valid for server versions >= 3.2. For servers < 3.2, this option is ignored.
func (ioo *InsertOneOptions) SetBypassDocumentValidation(b bool) *InsertOneOptions {
	ioo.BypassDocumentValidation = &b
	return ioo
}

// MergeInsertOneOptions combines the argued InsertOneOptions into a single InsertOneOptions in a last-one-wins fashion
func MergeInsertOneOptions(opts ...*InsertOneOptions) *InsertOneOptions {
	ioOpts := InsertOne()
	for _, ioo := range opts {
		if ioo == nil {
			continue
		}
		if ioo.BypassDocumentValidation != nil {
			ioOpts.BypassDocumentValidation = ioo.BypassDocumentValidation
		}
	}

	return ioOpts
}

// InsertManyOptions represents all possible options to the InsertMany() function.
type InsertManyOptions struct {
	BypassDocumentValidation *bool // If true, allows the write to opt-out of document level validation
	Ordered                  *bool // If true, when an insert fails, return without performing the remaining inserts. Defaults to true.
}

// InsertMany returns a pointer to a new InsertManyOptions
func InsertMany() *InsertManyOptions {
	return &InsertManyOptions{
		Ordered: &DefaultOrdered,
	}
}

// SetBypassDocumentValidation allows the write to opt-out of document level validation.
// Valid for server versions >= 3.2. For servers < 3.2, this option is ignored.
func (imo *InsertManyOptions) SetBypassDocumentValidation(b bool) *InsertManyOptions {
	imo.BypassDocumentValidation = &b
	return imo
}

// SetOrdered configures the ordered option. If true, when a write fails, the function will return without attempting
// remaining writes. Defaults to true.
func (imo *InsertManyOptions) SetOrdered(b bool) *InsertManyOptions {
	imo.Ordered = &b
	return imo
}

// MergeInsertManyOptions combines the argued InsertManyOptions into a single InsertManyOptions in a last-one-wins fashion
func MergeInsertManyOptions(opts ...*InsertManyOptions) *InsertManyOptions {
	imOpts := InsertMany()
	for _, imo := range opts {
		if imo == nil {
			continue
		}
		if imo.BypassDocumentValidation != nil {
			imOpts.BypassDocumentValidation = imo.BypassDocumentValidation
		}
		if imo.Ordered != nil {
			imOpts.Ordered = imo.Ordered
		}
	}

	return imOpts
}
