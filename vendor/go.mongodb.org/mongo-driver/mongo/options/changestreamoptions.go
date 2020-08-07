// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// ChangeStreamOptions represents all possible options to a change stream
type ChangeStreamOptions struct {
	BatchSize            *int32               // The number of documents to return per batch
	Collation            *Collation           // Specifies a collation
	FullDocument         *FullDocument        // When set to ‘updateLookup’, the change notification for partial updates will include both a delta describing the changes to the document, as well as a copy of the entire document that was changed from some time after the change occurred.
	MaxAwaitTime         *time.Duration       // The maximum amount of time for the server to wait on new documents to satisfy a change stream query
	ResumeAfter          interface{}          // Specifies the logical starting point for the new change stream
	StartAtOperationTime *primitive.Timestamp // Ensures that a change stream will only provide changes that occurred after a timestamp.
	StartAfter           interface{}          // Specifies a resume token. The started change stream will return the first notification after the token.
}

// ChangeStream returns a pointer to a new ChangeStreamOptions
func ChangeStream() *ChangeStreamOptions {
	cso := &ChangeStreamOptions{}
	cso.SetFullDocument(Default)
	return cso
}

// SetBatchSize specifies the number of documents to return per batch
func (cso *ChangeStreamOptions) SetBatchSize(i int32) *ChangeStreamOptions {
	cso.BatchSize = &i
	return cso
}

// SetCollation specifies a collation
func (cso *ChangeStreamOptions) SetCollation(c Collation) *ChangeStreamOptions {
	cso.Collation = &c
	return cso
}

// SetFullDocument specifies the fullDocument option.
// When set to ‘updateLookup’, the change notification for partial updates will
// include both a delta describing the changes to the document, as well as a
// copy of the entire document that was changed from some time after the change
// occurred.
func (cso *ChangeStreamOptions) SetFullDocument(fd FullDocument) *ChangeStreamOptions {
	cso.FullDocument = &fd
	return cso
}

// SetMaxAwaitTime specifies the maximum amount of time for the server to wait on new documents to satisfy a change stream query
func (cso *ChangeStreamOptions) SetMaxAwaitTime(d time.Duration) *ChangeStreamOptions {
	cso.MaxAwaitTime = &d
	return cso
}

// SetResumeAfter specifies the logical starting point for the new change stream
func (cso *ChangeStreamOptions) SetResumeAfter(rt interface{}) *ChangeStreamOptions {
	cso.ResumeAfter = rt
	return cso
}

// SetStartAtOperationTime ensures that a change stream will only provide changes that occurred after a specified timestamp.
func (cso *ChangeStreamOptions) SetStartAtOperationTime(t *primitive.Timestamp) *ChangeStreamOptions {
	cso.StartAtOperationTime = t
	return cso
}

// SetStartAfter specifies a resume token. The resulting change stream will return the first notification after the token.
// Cannot be used in conjunction with ResumeAfter.
func (cso *ChangeStreamOptions) SetStartAfter(sa interface{}) *ChangeStreamOptions {
	cso.StartAfter = sa
	return cso
}

// MergeChangeStreamOptions combines the argued ChangeStreamOptions into a single ChangeStreamOptions in a last-one-wins fashion
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
