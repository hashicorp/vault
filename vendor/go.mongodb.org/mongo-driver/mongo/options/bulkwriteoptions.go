// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// DefaultOrdered is the default value for the Ordered option in BulkWriteOptions.
var DefaultOrdered = true

// BulkWriteOptions represents options that can be used to configure a BulkWrite operation.
type BulkWriteOptions struct {
	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://www.mongodb.com/docs/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default value is nil, which means that no comment will be included in the logs.
	Comment interface{}

	// If true, no writes will be executed after one fails. The default value is true.
	Ordered *bool

	// Specifies parameters for all update and delete commands in the BulkWrite. This option is only valid for MongoDB
	// versions >= 5.0. Older servers will report an error for using this option. This must be a document mapping
	// parameter names to values. Values must be constant or closed expressions that do not reference document fields.
	// Parameters can then be accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface{}
}

// BulkWrite creates a new *BulkWriteOptions instance.
func BulkWrite() *BulkWriteOptions {
	return &BulkWriteOptions{
		Ordered: &DefaultOrdered,
	}
}

// SetComment sets the value for the Comment field.
func (b *BulkWriteOptions) SetComment(comment interface{}) *BulkWriteOptions {
	b.Comment = comment
	return b
}

// SetOrdered sets the value for the Ordered field.
func (b *BulkWriteOptions) SetOrdered(ordered bool) *BulkWriteOptions {
	b.Ordered = &ordered
	return b
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (b *BulkWriteOptions) SetBypassDocumentValidation(bypass bool) *BulkWriteOptions {
	b.BypassDocumentValidation = &bypass
	return b
}

// SetLet sets the value for the Let field. Let specifies parameters for all update and delete commands in the BulkWrite.
// This option is only valid for MongoDB versions >= 5.0. Older servers will report an error for using this option.
// This must be a document mapping parameter names to values. Values must be constant or closed expressions that do not
// reference document fields. Parameters can then be accessed as variables in an aggregate expression context (e.g. "$$var").
func (b *BulkWriteOptions) SetLet(let interface{}) *BulkWriteOptions {
	b.Let = &let
	return b
}

// MergeBulkWriteOptions combines the given BulkWriteOptions instances into a single BulkWriteOptions in a last-one-wins
// fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeBulkWriteOptions(opts ...*BulkWriteOptions) *BulkWriteOptions {
	b := BulkWrite()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Comment != nil {
			b.Comment = opt.Comment
		}
		if opt.Ordered != nil {
			b.Ordered = opt.Ordered
		}
		if opt.BypassDocumentValidation != nil {
			b.BypassDocumentValidation = opt.BypassDocumentValidation
		}
		if opt.Let != nil {
			b.Let = opt.Let
		}
	}

	return b
}
