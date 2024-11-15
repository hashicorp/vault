// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	// A string that will be included in server logs, profiling logs, and currentOp queries to help trace the operation.
	// The default is nil, which means that no comment will be included in the logs.
	Comment *string

	// Specifies how the updated document should be returned in change notifications for update operations. The default
	// is options.Default, which means that only partial update deltas will be included in the change notification.
	FullDocument *FullDocument

	// Specifies how the pre-update document should be returned in change notifications for update operations. The default
	// is options.Off, which means that the pre-update document will not be included in the change notification.
	FullDocumentBeforeChange *FullDocument

	// The maximum amount of time that the server should wait for new documents to satisfy a tailable cursor query.
	MaxAwaitTime *time.Duration

	// A document specifying the logical starting point for the change stream. Only changes corresponding to an oplog
	// entry immediately after the resume token will be returned. If this is specified, StartAtOperationTime and
	// StartAfter must not be set.
	ResumeAfter interface{}

	// ShowExpandedEvents specifies whether the server will return an expanded list of change stream events. Additional
	// events include: createIndexes, dropIndexes, modify, create, shardCollection, reshardCollection and
	// refineCollectionShardKey. This option is only valid for MongoDB versions >= 6.0.
	ShowExpandedEvents *bool

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

	// Custom options to be added to the initial aggregate for the change stream. Key-value pairs of the BSON map should
	// correlate with desired option names and values. Values must be Marshalable. Custom options may conflict with
	// non-custom options, and custom options bypass client-side validation. Prefer using non-custom options where possible.
	Custom bson.M

	// Custom options to be added to the $changeStream stage in the initial aggregate. Key-value pairs of the BSON map should
	// correlate with desired option names and values. Values must be Marshalable. Custom pipeline options bypass client-side
	// validation. Prefer using non-custom options where possible.
	CustomPipeline bson.M
}

// ChangeStream creates a new ChangeStreamOptions instance.
func ChangeStream() *ChangeStreamOptions {
	cso := &ChangeStreamOptions{}
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

// SetComment sets the value for the Comment field.
func (cso *ChangeStreamOptions) SetComment(comment string) *ChangeStreamOptions {
	cso.Comment = &comment
	return cso
}

// SetFullDocument sets the value for the FullDocument field.
func (cso *ChangeStreamOptions) SetFullDocument(fd FullDocument) *ChangeStreamOptions {
	cso.FullDocument = &fd
	return cso
}

// SetFullDocumentBeforeChange sets the value for the FullDocumentBeforeChange field.
func (cso *ChangeStreamOptions) SetFullDocumentBeforeChange(fdbc FullDocument) *ChangeStreamOptions {
	cso.FullDocumentBeforeChange = &fdbc
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

// SetShowExpandedEvents sets the value for the ShowExpandedEvents field.
func (cso *ChangeStreamOptions) SetShowExpandedEvents(see bool) *ChangeStreamOptions {
	cso.ShowExpandedEvents = &see
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

// SetCustom sets the value for the Custom field. Key-value pairs of the BSON map should correlate
// with desired option names and values. Values must be Marshalable. Custom options may conflict
// with non-custom options, and custom options bypass client-side validation. Prefer using non-custom
// options where possible.
func (cso *ChangeStreamOptions) SetCustom(c bson.M) *ChangeStreamOptions {
	cso.Custom = c
	return cso
}

// SetCustomPipeline sets the value for the CustomPipeline field. Key-value pairs of the BSON map
// should correlate with desired option names and values. Values must be Marshalable. Custom pipeline
// options bypass client-side validation. Prefer using non-custom options where possible.
func (cso *ChangeStreamOptions) SetCustomPipeline(cp bson.M) *ChangeStreamOptions {
	cso.CustomPipeline = cp
	return cso
}

// MergeChangeStreamOptions combines the given ChangeStreamOptions instances into a single ChangeStreamOptions in a
// last-one-wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
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
		if cso.Comment != nil {
			csOpts.Comment = cso.Comment
		}
		if cso.FullDocument != nil {
			csOpts.FullDocument = cso.FullDocument
		}
		if cso.FullDocumentBeforeChange != nil {
			csOpts.FullDocumentBeforeChange = cso.FullDocumentBeforeChange
		}
		if cso.MaxAwaitTime != nil {
			csOpts.MaxAwaitTime = cso.MaxAwaitTime
		}
		if cso.ResumeAfter != nil {
			csOpts.ResumeAfter = cso.ResumeAfter
		}
		if cso.ShowExpandedEvents != nil {
			csOpts.ShowExpandedEvents = cso.ShowExpandedEvents
		}
		if cso.StartAtOperationTime != nil {
			csOpts.StartAtOperationTime = cso.StartAtOperationTime
		}
		if cso.StartAfter != nil {
			csOpts.StartAfter = cso.StartAfter
		}
		if cso.Custom != nil {
			csOpts.Custom = cso.Custom
		}
		if cso.CustomPipeline != nil {
			csOpts.CustomPipeline = cso.CustomPipeline
		}
	}

	return csOpts
}
