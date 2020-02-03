// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// EstimatedDocumentCountOptions represents all possible options to the EstimatedDocumentCount() function.
type EstimatedDocumentCountOptions struct {
	MaxTime *time.Duration // The maximum amount of time to allow the operation to run
}

// EstimatedDocumentCount returns a pointer to a new EstimatedDocumentCountOptions
func EstimatedDocumentCount() *EstimatedDocumentCountOptions {
	return &EstimatedDocumentCountOptions{}
}

// SetMaxTime specifies the maximum amount of time to allow the operation to run
func (eco *EstimatedDocumentCountOptions) SetMaxTime(d time.Duration) *EstimatedDocumentCountOptions {
	eco.MaxTime = &d
	return eco
}

// MergeEstimatedDocumentCountOptions combines the given *EstimatedDocumentCountOptions into a single
// *EstimatedDocumentCountOptions in a last one wins fashion.
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
