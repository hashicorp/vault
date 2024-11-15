// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// ReplaceOptions represents options that can be used to configure a ReplaceOne operation.
type ReplaceOptions struct {
	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://www.mongodb.com/docs/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool

	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default value is nil, which means that no comment will be included in the logs.
	Comment interface{}

	// The index to use for the operation. This should either be the index name as a string or the index specification
	// as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will return an error
	// if this option is specified. For server versions < 3.4, the driver will return a client-side error if this option
	// is specified. The driver will return an error if this option is specified during an unacknowledged write
	// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil,
	// which means that no hint will be sent.
	Hint interface{}

	// If true, a new document will be inserted if the filter does not match any documents in the collection. The
	// default value is false.
	Upsert *bool

	// Specifies parameters for the aggregate expression. This option is only valid for MongoDB versions >= 5.0. Older
	// servers will report an error for using this option. This must be a document mapping parameter names to values.
	// Values must be constant or closed expressions that do not reference document fields. Parameters can then be
	// accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface{}
}

// Replace creates a new ReplaceOptions instance.
func Replace() *ReplaceOptions {
	return &ReplaceOptions{}
}

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (ro *ReplaceOptions) SetBypassDocumentValidation(b bool) *ReplaceOptions {
	ro.BypassDocumentValidation = &b
	return ro
}

// SetCollation sets the value for the Collation field.
func (ro *ReplaceOptions) SetCollation(c *Collation) *ReplaceOptions {
	ro.Collation = c
	return ro
}

// SetComment sets the value for the Comment field.
func (ro *ReplaceOptions) SetComment(comment interface{}) *ReplaceOptions {
	ro.Comment = comment
	return ro
}

// SetHint sets the value for the Hint field.
func (ro *ReplaceOptions) SetHint(h interface{}) *ReplaceOptions {
	ro.Hint = h
	return ro
}

// SetUpsert sets the value for the Upsert field.
func (ro *ReplaceOptions) SetUpsert(b bool) *ReplaceOptions {
	ro.Upsert = &b
	return ro
}

// SetLet sets the value for the Let field.
func (ro *ReplaceOptions) SetLet(l interface{}) *ReplaceOptions {
	ro.Let = l
	return ro
}

// MergeReplaceOptions combines the given ReplaceOptions instances into a single ReplaceOptions in a last-one-wins
// fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
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
		if ro.Comment != nil {
			rOpts.Comment = ro.Comment
		}
		if ro.Hint != nil {
			rOpts.Hint = ro.Hint
		}
		if ro.Upsert != nil {
			rOpts.Upsert = ro.Upsert
		}
		if ro.Let != nil {
			rOpts.Let = ro.Let
		}
	}

	return rOpts
}
