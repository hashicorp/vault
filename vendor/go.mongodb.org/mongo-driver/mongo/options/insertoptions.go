// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// InsertOneOptions represents options that can be used to configure an InsertOne operation.
type InsertOneOptions struct {
	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://docs.mongodb.com/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool
}

// InsertOne creates a new InsertOneOptions instance.
func InsertOne() *InsertOneOptions {
	return &InsertOneOptions{}
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (ioo *InsertOneOptions) SetBypassDocumentValidation(b bool) *InsertOneOptions {
	ioo.BypassDocumentValidation = &b
	return ioo
}

// MergeInsertOneOptions combines the given InsertOneOptions instances into a single InsertOneOptions in a last-one-wins
// fashion.
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

// InsertManyOptions represents options that can be used to configure an InsertMany operation.
type InsertManyOptions struct {
	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://docs.mongodb.com/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool

	// If true, no writes will be executed after one fails. The default value is true.
	Ordered *bool
}

// InsertMany creates a new InsertManyOptions instance.
func InsertMany() *InsertManyOptions {
	return &InsertManyOptions{
		Ordered: &DefaultOrdered,
	}
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (imo *InsertManyOptions) SetBypassDocumentValidation(b bool) *InsertManyOptions {
	imo.BypassDocumentValidation = &b
	return imo
}

// SetOrdered sets the value for the Ordered field.
func (imo *InsertManyOptions) SetOrdered(b bool) *InsertManyOptions {
	imo.Ordered = &b
	return imo
}

// MergeInsertManyOptions combines the given InsertManyOptions instances into a single InsertManyOptions in a last one
// wins fashion.
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
