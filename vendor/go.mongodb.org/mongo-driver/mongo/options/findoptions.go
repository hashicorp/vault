// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"
)

// FindOptions represent all possible options to the Find() function.
type FindOptions struct {
	AllowPartialResults *bool          // If true, allows partial results to be returned if some shards are down.
	BatchSize           *int32         // Specifies the number of documents to return in every batch.
	Collation           *Collation     // Specifies a collation to be used
	Comment             *string        // Specifies a string to help trace the operation through the database.
	CursorType          *CursorType    // Specifies the type of cursor to use
	Hint                interface{}    // Specifies the index to use.
	Limit               *int64         // Sets a limit on the number of results to return.
	Max                 interface{}    // Sets an exclusive upper bound for a specific index
	MaxAwaitTime        *time.Duration // Specifies the maximum amount of time for the server to wait on new documents.
	MaxTime             *time.Duration // Specifies the maximum amount of time to allow the query to run.
	Min                 interface{}    // Specifies the inclusive lower bound for a specific index.
	NoCursorTimeout     *bool          // If true, prevents cursors from timing out after an inactivity period.
	OplogReplay         *bool          // Adds an option for internal use only and should not be set.
	Projection          interface{}    // Limits the fields returned for all documents.
	ReturnKey           *bool          // If true, only returns index keys for all result documents.
	ShowRecordID        *bool          // If true, a $recordId field with the record identifier will be added to the returned documents.
	Skip                *int64         // Specifies the number of documents to skip before returning
	Snapshot            *bool          // If true, prevents the cursor from returning a document more than once because of an intervening write operation.
	Sort                interface{}    // Specifies the order in which to return results.
}

// Find creates a new FindOptions instance.
func Find() *FindOptions {
	return &FindOptions{}
}

// SetAllowPartialResults sets whether partial results can be returned if some shards are down.
// For server versions < 3.2, this defaults to false.
func (f *FindOptions) SetAllowPartialResults(b bool) *FindOptions {
	f.AllowPartialResults = &b
	return f
}

// SetBatchSize sets the number of documents to return in each batch.
func (f *FindOptions) SetBatchSize(i int32) *FindOptions {
	f.BatchSize = &i
	return f
}

// SetCollation specifies a Collation to use for the Find operation.
// Valid for server versions >= 3.4
func (f *FindOptions) SetCollation(collation *Collation) *FindOptions {
	f.Collation = collation
	return f
}

// SetComment specifies a string to help trace the operation through the database.
func (f *FindOptions) SetComment(comment string) *FindOptions {
	f.Comment = &comment
	return f
}

// SetCursorType specifes the type of cursor to use.
func (f *FindOptions) SetCursorType(ct CursorType) *FindOptions {
	f.CursorType = &ct
	return f
}

// SetHint specifies the index to use.
func (f *FindOptions) SetHint(hint interface{}) *FindOptions {
	f.Hint = hint
	return f
}

// SetLimit specifies a limit on the number of results.
// A negative limit implies that only 1 batch should be returned.
func (f *FindOptions) SetLimit(i int64) *FindOptions {
	f.Limit = &i
	return f
}

// SetMax specifies an exclusive upper bound for a specific index.
func (f *FindOptions) SetMax(max interface{}) *FindOptions {
	f.Max = max
	return f
}

// SetMaxAwaitTime specifies the max amount of time for the server to wait on new documents.
// If the cursor type is not TailableAwait, this option is ignored.
// For server versions < 3.2, this option is ignored.
func (f *FindOptions) SetMaxAwaitTime(d time.Duration) *FindOptions {
	f.MaxAwaitTime = &d
	return f
}

// SetMaxTime specifies the max time to allow the query to run.
func (f *FindOptions) SetMaxTime(d time.Duration) *FindOptions {
	f.MaxTime = &d
	return f
}

// SetMin specifies the inclusive lower bound for a specific index.
func (f *FindOptions) SetMin(min interface{}) *FindOptions {
	f.Min = min
	return f
}

// SetNoCursorTimeout specifies whether or not cursors should time out after a period of inactivity.
// For server versions < 3.2, this defaults to false.
func (f *FindOptions) SetNoCursorTimeout(b bool) *FindOptions {
	f.NoCursorTimeout = &b
	return f
}

// SetOplogReplay adds an option for internal use only and should not be set.
// For server versions < 3.2, this defaults to false.
func (f *FindOptions) SetOplogReplay(b bool) *FindOptions {
	f.OplogReplay = &b
	return f
}

// SetProjection adds an option to limit the fields returned for all documents.
func (f *FindOptions) SetProjection(projection interface{}) *FindOptions {
	f.Projection = projection
	return f
}

// SetReturnKey adds an option to only return index keys for all result documents.
func (f *FindOptions) SetReturnKey(b bool) *FindOptions {
	f.ReturnKey = &b
	return f
}

// SetShowRecordID adds an option to determine whether to return the record identifier for each document.
// If true, a $recordId field will be added to each returned document.
func (f *FindOptions) SetShowRecordID(b bool) *FindOptions {
	f.ShowRecordID = &b
	return f
}

// SetSkip specifies the number of documents to skip before returning.
// For server versions < 3.2, this defaults to 0.
func (f *FindOptions) SetSkip(i int64) *FindOptions {
	f.Skip = &i
	return f
}

// SetSnapshot prevents the cursor from returning a document more than once because of an intervening write operation.
func (f *FindOptions) SetSnapshot(b bool) *FindOptions {
	f.Snapshot = &b
	return f
}

// SetSort specifies the order in which to return documents.
func (f *FindOptions) SetSort(sort interface{}) *FindOptions {
	f.Sort = sort
	return f
}

// MergeFindOptions combines the argued FindOptions into a single FindOptions in a last-one-wins fashion
func MergeFindOptions(opts ...*FindOptions) *FindOptions {
	fo := Find()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AllowPartialResults != nil {
			fo.AllowPartialResults = opt.AllowPartialResults
		}
		if opt.BatchSize != nil {
			fo.BatchSize = opt.BatchSize
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.Comment != nil {
			fo.Comment = opt.Comment
		}
		if opt.CursorType != nil {
			fo.CursorType = opt.CursorType
		}
		if opt.Hint != nil {
			fo.Hint = opt.Hint
		}
		if opt.Limit != nil {
			fo.Limit = opt.Limit
		}
		if opt.Max != nil {
			fo.Max = opt.Max
		}
		if opt.MaxAwaitTime != nil {
			fo.MaxAwaitTime = opt.MaxAwaitTime
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Min != nil {
			fo.Min = opt.Min
		}
		if opt.NoCursorTimeout != nil {
			fo.NoCursorTimeout = opt.NoCursorTimeout
		}
		if opt.OplogReplay != nil {
			fo.OplogReplay = opt.OplogReplay
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnKey != nil {
			fo.ReturnKey = opt.ReturnKey
		}
		if opt.ShowRecordID != nil {
			fo.ShowRecordID = opt.ShowRecordID
		}
		if opt.Skip != nil {
			fo.Skip = opt.Skip
		}
		if opt.Snapshot != nil {
			fo.Snapshot = opt.Snapshot
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
	}

	return fo
}

// FindOneOptions represent all possible options to the FindOne() function.
type FindOneOptions struct {
	AllowPartialResults *bool          // If true, allows partial results to be returned if some shards are down.
	BatchSize           *int32         // Specifies the number of documents to return in every batch.
	Collation           *Collation     // Specifies a collation to be used
	Comment             *string        // Specifies a string to help trace the operation through the database.
	CursorType          *CursorType    // Specifies the type of cursor to use
	Hint                interface{}    // Specifies the index to use.
	Max                 interface{}    // Sets an exclusive upper bound for a specific index
	MaxAwaitTime        *time.Duration // Specifies the maximum amount of time for the server to wait on new documents.
	MaxTime             *time.Duration // Specifies the maximum amount of time to allow the query to run.
	Min                 interface{}    // Specifies the inclusive lower bound for a specific index.
	NoCursorTimeout     *bool          // If true, prevents cursors from timing out after an inactivity period.
	OplogReplay         *bool          // Adds an option for internal use only and should not be set.
	Projection          interface{}    // Limits the fields returned for all documents.
	ReturnKey           *bool          // If true, only returns index keys for all result documents.
	ShowRecordID        *bool          // If true, a $recordId field with the record identifier will be added to the returned documents.
	Skip                *int64         // Specifies the number of documents to skip before returning
	Snapshot            *bool          // If true, prevents the cursor from returning a document more than once because of an intervening write operation.
	Sort                interface{}    // Specifies the order in which to return results.
}

// FindOne creates a new FindOneOptions instance.
func FindOne() *FindOneOptions {
	return &FindOneOptions{}
}

// SetAllowPartialResults sets whether partial results can be returned if some shards are down.
func (f *FindOneOptions) SetAllowPartialResults(b bool) *FindOneOptions {
	f.AllowPartialResults = &b
	return f
}

// SetBatchSize sets the number of documents to return in each batch.
func (f *FindOneOptions) SetBatchSize(i int32) *FindOneOptions {
	f.BatchSize = &i
	return f
}

// SetCollation specifies a Collation to use for the Find operation.
func (f *FindOneOptions) SetCollation(collation *Collation) *FindOneOptions {
	f.Collation = collation
	return f
}

// SetComment specifies a string to help trace the operation through the database.
func (f *FindOneOptions) SetComment(comment string) *FindOneOptions {
	f.Comment = &comment
	return f
}

// SetCursorType specifes the type of cursor to use.
func (f *FindOneOptions) SetCursorType(ct CursorType) *FindOneOptions {
	f.CursorType = &ct
	return f
}

// SetHint specifies the index to use.
func (f *FindOneOptions) SetHint(hint interface{}) *FindOneOptions {
	f.Hint = hint
	return f
}

// SetMax specifies an exclusive upper bound for a specific index.
func (f *FindOneOptions) SetMax(max interface{}) *FindOneOptions {
	f.Max = max
	return f
}

// SetMaxAwaitTime specifies the max amount of time for the server to wait on new documents.
// For server versions < 3.2, this option is ignored.
func (f *FindOneOptions) SetMaxAwaitTime(d time.Duration) *FindOneOptions {
	f.MaxAwaitTime = &d
	return f
}

// SetMaxTime specifies the max time to allow the query to run.
func (f *FindOneOptions) SetMaxTime(d time.Duration) *FindOneOptions {
	f.MaxTime = &d
	return f
}

// SetMin specifies the inclusive lower bound for a specific index.
func (f *FindOneOptions) SetMin(min interface{}) *FindOneOptions {
	f.Min = min
	return f
}

// SetNoCursorTimeout specifies whether or not cursors should time out after a period of inactivity.
func (f *FindOneOptions) SetNoCursorTimeout(b bool) *FindOneOptions {
	f.NoCursorTimeout = &b
	return f
}

// SetOplogReplay adds an option for internal use only and should not be set.
func (f *FindOneOptions) SetOplogReplay(b bool) *FindOneOptions {
	f.OplogReplay = &b
	return f
}

// SetProjection adds an option to limit the fields returned for all documents.
func (f *FindOneOptions) SetProjection(projection interface{}) *FindOneOptions {
	f.Projection = projection
	return f
}

// SetReturnKey adds an option to only return index keys for all result documents.
func (f *FindOneOptions) SetReturnKey(b bool) *FindOneOptions {
	f.ReturnKey = &b
	return f
}

// SetShowRecordID adds an option to determine whether to return the record identifier for each document.
// If true, a $recordId field will be added to each returned document.
func (f *FindOneOptions) SetShowRecordID(b bool) *FindOneOptions {
	f.ShowRecordID = &b
	return f
}

// SetSkip specifies the number of documents to skip before returning.
func (f *FindOneOptions) SetSkip(i int64) *FindOneOptions {
	f.Skip = &i
	return f
}

// SetSnapshot prevents the cursor from returning a document more than once because of an intervening write operation.
func (f *FindOneOptions) SetSnapshot(b bool) *FindOneOptions {
	f.Snapshot = &b
	return f
}

// SetSort specifies the order in which to return documents.
func (f *FindOneOptions) SetSort(sort interface{}) *FindOneOptions {
	f.Sort = sort
	return f
}

// MergeFindOneOptions combines the argued FindOneOptions into a single FindOneOptions in a last-one-wins fashion
func MergeFindOneOptions(opts ...*FindOneOptions) *FindOneOptions {
	fo := FindOne()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AllowPartialResults != nil {
			fo.AllowPartialResults = opt.AllowPartialResults
		}
		if opt.BatchSize != nil {
			fo.BatchSize = opt.BatchSize
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.Comment != nil {
			fo.Comment = opt.Comment
		}
		if opt.CursorType != nil {
			fo.CursorType = opt.CursorType
		}
		if opt.Hint != nil {
			fo.Hint = opt.Hint
		}
		if opt.Max != nil {
			fo.Max = opt.Max
		}
		if opt.MaxAwaitTime != nil {
			fo.MaxAwaitTime = opt.MaxAwaitTime
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Min != nil {
			fo.Min = opt.Min
		}
		if opt.NoCursorTimeout != nil {
			fo.NoCursorTimeout = opt.NoCursorTimeout
		}
		if opt.OplogReplay != nil {
			fo.OplogReplay = opt.OplogReplay
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnKey != nil {
			fo.ReturnKey = opt.ReturnKey
		}
		if opt.ShowRecordID != nil {
			fo.ShowRecordID = opt.ShowRecordID
		}
		if opt.Skip != nil {
			fo.Skip = opt.Skip
		}
		if opt.Snapshot != nil {
			fo.Snapshot = opt.Snapshot
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
	}

	return fo
}

// FindOneAndReplaceOptions represent all possible options to the FindOneAndReplace() function.
type FindOneAndReplaceOptions struct {
	BypassDocumentValidation *bool           // If true, allows the write to opt out of document-level validation.
	Collation                *Collation      // Specifies a collation to be used
	MaxTime                  *time.Duration  // Specifies the maximum amount of time to allow the query to run.
	Projection               interface{}     // Limits the fields returned for all documents.
	ReturnDocument           *ReturnDocument // Specifies whether the original or updated document should be returned.
	Sort                     interface{}     // Specifies the order in which to return results.
	Upsert                   *bool           // If true, creates a a new document if no document matches the query.
}

// FindOneAndReplace creates a new FindOneAndReplaceOptions instance.
func FindOneAndReplace() *FindOneAndReplaceOptions {
	return &FindOneAndReplaceOptions{}
}

// SetBypassDocumentValidation specifies whether or not the write should opt out of document-level validation.
// Valid for server versions >= 3.2. For servers < 3.2, this option is ignored.
func (f *FindOneAndReplaceOptions) SetBypassDocumentValidation(b bool) *FindOneAndReplaceOptions {
	f.BypassDocumentValidation = &b
	return f
}

// SetCollation specifies a Collation to use for the Find operation.
func (f *FindOneAndReplaceOptions) SetCollation(collation *Collation) *FindOneAndReplaceOptions {
	f.Collation = collation
	return f
}

// SetMaxTime specifies the max time to allow the query to run.
func (f *FindOneAndReplaceOptions) SetMaxTime(d time.Duration) *FindOneAndReplaceOptions {
	f.MaxTime = &d
	return f
}

// SetProjection adds an option to limit the fields returned for all documents.
func (f *FindOneAndReplaceOptions) SetProjection(projection interface{}) *FindOneAndReplaceOptions {
	f.Projection = projection
	return f
}

// SetReturnDocument specifies whether the original or updated document should be returned.
// If set to Before, the original document will be returned. If set to After, the updated document
// will be returned.
func (f *FindOneAndReplaceOptions) SetReturnDocument(rd ReturnDocument) *FindOneAndReplaceOptions {
	f.ReturnDocument = &rd
	return f
}

// SetSort specifies the order in which to return documents.
func (f *FindOneAndReplaceOptions) SetSort(sort interface{}) *FindOneAndReplaceOptions {
	f.Sort = sort
	return f
}

// SetUpsert specifies if a new document should be created if no document matches the query.
func (f *FindOneAndReplaceOptions) SetUpsert(b bool) *FindOneAndReplaceOptions {
	f.Upsert = &b
	return f
}

// MergeFindOneAndReplaceOptions combines the argued FindOneAndReplaceOptions into a single FindOneAndReplaceOptions in a last-one-wins fashion
func MergeFindOneAndReplaceOptions(opts ...*FindOneAndReplaceOptions) *FindOneAndReplaceOptions {
	fo := FindOneAndReplace()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.BypassDocumentValidation != nil {
			fo.BypassDocumentValidation = opt.BypassDocumentValidation
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnDocument != nil {
			fo.ReturnDocument = opt.ReturnDocument
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
		if opt.Upsert != nil {
			fo.Upsert = opt.Upsert
		}
	}

	return fo
}

// FindOneAndUpdateOptions represent all possible options to the FindOneAndUpdate() function.
type FindOneAndUpdateOptions struct {
	ArrayFilters             *ArrayFilters   // A set of filters specifying to which array elements an update should apply.
	BypassDocumentValidation *bool           // If true, allows the write to opt out of document-level validation.
	Collation                *Collation      // Specifies a collation to be used
	MaxTime                  *time.Duration  // Specifies the maximum amount of time to allow the query to run.
	Projection               interface{}     // Limits the fields returned for all documents.
	ReturnDocument           *ReturnDocument // Specifies whether the original or updated document should be returned.
	Sort                     interface{}     // Specifies the order in which to return results.
	Upsert                   *bool           // If true, creates a a new document if no document matches the query.
}

// FindOneAndUpdate creates a new FindOneAndUpdateOptions instance.
func FindOneAndUpdate() *FindOneAndUpdateOptions {
	return &FindOneAndUpdateOptions{}
}

// SetBypassDocumentValidation sets filters that specify to which array elements an update should apply.
func (f *FindOneAndUpdateOptions) SetBypassDocumentValidation(b bool) *FindOneAndUpdateOptions {
	f.BypassDocumentValidation = &b
	return f
}

// SetArrayFilters specifies a set of filters, which
func (f *FindOneAndUpdateOptions) SetArrayFilters(filters ArrayFilters) *FindOneAndUpdateOptions {
	f.ArrayFilters = &filters
	return f
}

// SetCollation specifies a Collation to use for the Find operation.
func (f *FindOneAndUpdateOptions) SetCollation(collation *Collation) *FindOneAndUpdateOptions {
	f.Collation = collation
	return f
}

// SetMaxTime specifies the max time to allow the query to run.
func (f *FindOneAndUpdateOptions) SetMaxTime(d time.Duration) *FindOneAndUpdateOptions {
	f.MaxTime = &d
	return f
}

// SetProjection adds an option to limit the fields returned for all documents.
func (f *FindOneAndUpdateOptions) SetProjection(projection interface{}) *FindOneAndUpdateOptions {
	f.Projection = projection
	return f
}

// SetReturnDocument specifies whether the original or updated document should be returned.
// If set to Before, the original document will be returned. If set to After, the updated document
// will be returned.
func (f *FindOneAndUpdateOptions) SetReturnDocument(rd ReturnDocument) *FindOneAndUpdateOptions {
	f.ReturnDocument = &rd
	return f
}

// SetSort specifies the order in which to return documents.
func (f *FindOneAndUpdateOptions) SetSort(sort interface{}) *FindOneAndUpdateOptions {
	f.Sort = sort
	return f
}

// SetUpsert specifies if a new document should be created if no document matches the query.
func (f *FindOneAndUpdateOptions) SetUpsert(b bool) *FindOneAndUpdateOptions {
	f.Upsert = &b
	return f
}

// MergeFindOneAndUpdateOptions combines the argued FindOneAndUpdateOptions into a single FindOneAndUpdateOptions in a last-one-wins fashion
func MergeFindOneAndUpdateOptions(opts ...*FindOneAndUpdateOptions) *FindOneAndUpdateOptions {
	fo := FindOneAndUpdate()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ArrayFilters != nil {
			fo.ArrayFilters = opt.ArrayFilters
		}
		if opt.BypassDocumentValidation != nil {
			fo.BypassDocumentValidation = opt.BypassDocumentValidation
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnDocument != nil {
			fo.ReturnDocument = opt.ReturnDocument
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
		if opt.Upsert != nil {
			fo.Upsert = opt.Upsert
		}
	}

	return fo
}

// FindOneAndDeleteOptions represent all possible options to the FindOneAndDelete() function.
type FindOneAndDeleteOptions struct {
	Collation  *Collation     // Specifies a collation to be used
	MaxTime    *time.Duration // Specifies the maximum amount of time to allow the query to run.
	Projection interface{}    // Limits the fields returned for all documents.
	Sort       interface{}    // Specifies the order in which to return results.
}

// FindOneAndDelete creates a new FindOneAndDeleteOptions instance.
func FindOneAndDelete() *FindOneAndDeleteOptions {
	return &FindOneAndDeleteOptions{}
}

// SetCollation specifies a Collation to use for the Find operation.
// Valid for server versions >= 3.4
func (f *FindOneAndDeleteOptions) SetCollation(collation *Collation) *FindOneAndDeleteOptions {
	f.Collation = collation
	return f
}

// SetMaxTime specifies the max time to allow the query to run.
func (f *FindOneAndDeleteOptions) SetMaxTime(d time.Duration) *FindOneAndDeleteOptions {
	f.MaxTime = &d
	return f
}

// SetProjection adds an option to limit the fields returned for all documents.
func (f *FindOneAndDeleteOptions) SetProjection(projection interface{}) *FindOneAndDeleteOptions {
	f.Projection = projection
	return f
}

// SetSort specifies the order in which to return documents.
func (f *FindOneAndDeleteOptions) SetSort(sort interface{}) *FindOneAndDeleteOptions {
	f.Sort = sort
	return f
}

// MergeFindOneAndDeleteOptions combines the argued FindOneAndDeleteOptions into a single FindOneAndDeleteOptions in a last-one-wins fashion
func MergeFindOneAndDeleteOptions(opts ...*FindOneAndDeleteOptions) *FindOneAndDeleteOptions {
	fo := FindOneAndDelete()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
	}

	return fo
}
