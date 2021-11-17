// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChangeStreamOptions represents options that can be used to configure a Watch operation.
type ChangeStreamOptions struct {
	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// Specifies whether the updated document should be returned in change notifications for update operations along
	// with the deltas describing the changes made to the document. The default is options.Default, which means that
	// the updated document will not be included in the change notification.
	FullDocument *FullDocument

	// The maximum amount of time that the server should wait for new documents to satisfy a tailable cursor query.
	MaxAwaitTime *time.Duration

	// A document specifying the logical starting point for the change stream. Only changes corresponding to an oplog
	// entry immediately after the resume token will be returned. If this is specified, StartAtOperationTime and
	// StartAfter must not be set.
	ResumeAfter interface{}

	// If specified, the change stream will only return changes that occurred at or after the given timestamp. This
	// option is only valid for MongoDB versions >= 4.0. If this is specified, ResumeAfter and StartAfter must not be
	// set.
	StartAtOperationTime *primitive.Timestamp

	// A document specifying the logical starting point for the change stream. This is similar to the ResumeAfter
	// option, but allows a resume token from an "invalidate" notification to be used. This allows a change stream on a
	// collection to be resumed after the collection has been dropped and recreated or renamed. Only changes
	// corresponding to an oplog entry immediately after the specified token will be returned. If this is specified,
	// ResumeAfter and StartAtOperationTime must not be set. This option is only valid for MongoDB versions >= 4.1.1.
	StartAfter interface{}
}

// ChangeStream creates a new ChangeStreamOptions instance.
func ChangeStream() *ChangeStreamOptions {
	cso := &ChangeStreamOptions{}
	cso.SetFullDocument(Default)
	return cso
}

// SetBatchSize sets the value for the BatchSize field.
func (cso *ChangeStreamOptions) SetBatchSize(i int32) *ChangeStreamOptions {
	cso.BatchSize = &i
	return cso
}

// SetCollation sets the value for the Collation field.
func (cso *ChangeStreamOptions) SetCollation(c Collation) *ChangeStreamOptions {
	cso.Collation = &c
	return cso
}

// SetFullDocument sets the value for the FullDocument field.
func (cso *ChangeStreamOptions) SetFullDocument(fd FullDocument) *ChangeStreamOptions {
	cso.FullDocument = &fd
	return cso
}

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (cso *ChangeStreamOptions) SetMaxAwaitTime(d time.Duration) *ChangeStreamOptions {
	cso.MaxAwaitTime = &d
	return cso
}

// SetResumeAfter sets the value for the ResumeAfter field.
func (cso *ChangeStreamOptions) SetResumeAfter(rt interface{}) *ChangeStreamOptions {
	cso.ResumeAfter = rt
	return cso
}

// SetStartAtOperationTime sets the value for the StartAtOperationTime field.
func (cso *ChangeStreamOptions) SetStartAtOperationTime(t *primitive.Timestamp) *ChangeStreamOptions {
	cso.StartAtOperationTime = t
	return cso
}

// SetStartAfter sets the value for the StartAfter field.
func (cso *ChangeStreamOptions) SetStartAfter(sa interface{}) *ChangeStreamOptions {
	cso.StartAfter = sa
	return cso
}

// MergeChangeStreamOptions combines the given ChangeStreamOptions instances into a single ChangeStreamOptions in a
// last-one-wins fashion.
func MergeChangeStreamOptions(opts ...*ChangeStreamOptions) *ChangeStreamOptions {
	csOpts := ChangeStream()
	for _, cso := range opts {
		if cso == nil {
			continue
		}
		if cso.BatchSize != nil {
			csOpts.BatchSize = cso.BatchSize
		}
		if cso.Collation != nil {
			csOpts.Collation = cso.Collation
		}
		if cso.FullDocument != nil {
			csOpts.FullDocument = cso.FullDocument
		}
		if cso.MaxAwaitTime != nil {
			csOpts.MaxAwaitTime = cso.MaxAwaitTime
		}
		if cso.ResumeAfter != nil {
			csOpts.ResumeAfter = cso.ResumeAfter
		}
		if cso.StartAtOperationTime != nil {
			csOpts.StartAtOperationTime = cso.StartAtOperationTime
		}
		if cso.StartAfter != nil {
			csOpts.StartAfter = cso.StartAfter
		}
	}

	return csOpts
}
