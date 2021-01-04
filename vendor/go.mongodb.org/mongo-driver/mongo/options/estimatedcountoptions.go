// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// EstimatedDocumentCountOptions represents options that can be used to configure an EstimatedDocumentCount operation.
type EstimatedDocumentCountOptions struct {
	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	MaxTime *time.Duration
}

// EstimatedDocumentCount creates a new EstimatedDocumentCountOptions instance.
func EstimatedDocumentCount() *EstimatedDocumentCountOptions {
	return &EstimatedDocumentCountOptions{}
}

// SetMaxTime sets the value for the MaxTime field.
func (eco *EstimatedDocumentCountOptions) SetMaxTime(d time.Duration) *EstimatedDocumentCountOptions {
	eco.MaxTime = &d
	return eco
}

// MergeEstimatedDocumentCountOptions combines the given EstimatedDocumentCountOptions instances into a single
// EstimatedDocumentCountOptions in a last-one-wins fashion.
func MergeEstimatedDocumentCountOptions(opts ...*EstimatedDocumentCountOptions) *EstimatedDocumentCountOptions {
	e := EstimatedDocumentCount()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.MaxTime != nil {
			e.MaxTime = opt.MaxTime
		}
	}

	return e
}
