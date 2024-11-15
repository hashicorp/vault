// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// EstimatedDocumentCountOptions represents options that can be used to configure an EstimatedDocumentCount operation.
type EstimatedDocumentCountOptions struct {
	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default is nil, which means that no comment will be included in the logs.
	Comment interface{}

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration
}

// EstimatedDocumentCount creates a new EstimatedDocumentCountOptions instance.
func EstimatedDocumentCount() *EstimatedDocumentCountOptions {
	return &EstimatedDocumentCountOptions{}
}

// SetComment sets the value for the Comment field.
func (eco *EstimatedDocumentCountOptions) SetComment(comment interface{}) *EstimatedDocumentCountOptions {
	eco.Comment = comment
	return eco
}

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option
// may be used in its place to control the amount of time that a single operation can run before
// returning an error. MaxTime is ignored if Timeout is set on the client.
func (eco *EstimatedDocumentCountOptions) SetMaxTime(d time.Duration) *EstimatedDocumentCountOptions {
	eco.MaxTime = &d
	return eco
}

// MergeEstimatedDocumentCountOptions combines the given EstimatedDocumentCountOptions instances into a single
// EstimatedDocumentCountOptions in a last-one-wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeEstimatedDocumentCountOptions(opts ...*EstimatedDocumentCountOptions) *EstimatedDocumentCountOptions {
	e := EstimatedDocumentCount()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Comment != nil {
			e.Comment = opt.Comment
		}
		if opt.MaxTime != nil {
			e.MaxTime = opt.MaxTime
		}
	}

	return e
}
