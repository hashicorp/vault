// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"errors"
	"io"
	"strings"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ListCollectionsBatchCursor is a special batch cursor returned from ListCollections that properly
// handles current and legacy ListCollections operations.
type ListCollectionsBatchCursor struct {
	legacy       bool // server version < 3.0
	bc           *BatchCursor
	currentBatch *bsoncore.DocumentSequence
	err          error
}

// NewListCollectionsBatchCursor creates a new non-legacy ListCollectionsCursor.
func NewListCollectionsBatchCursor(bc *BatchCursor) (*ListCollectionsBatchCursor, error) {
	if bc == nil {
		return nil, errors.New("batch cursor must not be nil")
	}
	return &ListCollectionsBatchCursor{bc: bc, currentBatch: new(bsoncore.DocumentSequence)}, nil
}

// NewLegacyListCollectionsBatchCursor creates a new legacy ListCollectionsCursor.
func NewLegacyListCollectionsBatchCursor(bc *BatchCursor) (*ListCollectionsBatchCursor, error) {
	if bc == nil {
		return nil, errors.New("batch cursor must not be nil")
	}
	return &ListCollectionsBatchCursor{legacy: true, bc: bc, currentBatch: new(bsoncore.DocumentSequence)}, nil
}

// ID returns the cursor ID for this batch cursor.
func (lcbc *ListCollectionsBatchCursor) ID() int64 {
	return lcbc.bc.ID()
}

// Next indicates if there is another batch available. Returning false does not necessarily indicate
// that the cursor is closed. This method will return false when an empty batch is returned.
//
// If Next returns true, there is a valid batch of documents available. If Next returns false, there
// is not a valid batch of documents available.
func (lcbc *ListCollectionsBatchCursor) Next(ctx context.Context) bool {
	if !lcbc.bc.Next(ctx) {
		return false
	}

	if !lcbc.legacy {
		lcbc.currentBatch.Style = lcbc.bc.currentBatch.Style
		lcbc.currentBatch.Data = lcbc.bc.currentBatch.Data
		lcbc.currentBatch.ResetIterator()
		return true
	}

	lcbc.currentBatch.Style = bsoncore.SequenceStyle
	lcbc.currentBatch.Data = lcbc.currentBatch.Data[:0]

	var doc bsoncore.Document
	for {
		doc, lcbc.err = lcbc.bc.currentBatch.Next()
		if lcbc.err != nil {
			if lcbc.err == io.EOF {
				lcbc.err = nil
				break
			}
			return false
		}
		doc, lcbc.err = lcbc.projectNameElement(doc)
		if lcbc.err != nil {
			return false
		}
		lcbc.currentBatch.Data = append(lcbc.currentBatch.Data, doc...)
	}

	return true
}

// Batch will return a DocumentSequence for the current batch of documents. The returned
// DocumentSequence is only valid until the next call to Next or Close.
func (lcbc *ListCollectionsBatchCursor) Batch() *bsoncore.DocumentSequence { return lcbc.currentBatch }

// Server returns a pointer to the cursor's server.
func (lcbc *ListCollectionsBatchCursor) Server() Server { return lcbc.bc.server }

// Err returns the latest error encountered.
func (lcbc *ListCollectionsBatchCursor) Err() error {
	if lcbc.err != nil {
		return lcbc.err
	}
	return lcbc.bc.Err()
}

// Close closes this batch cursor.
func (lcbc *ListCollectionsBatchCursor) Close(ctx context.Context) error { return lcbc.bc.Close(ctx) }

// project out the database name for a legacy server
func (*ListCollectionsBatchCursor) projectNameElement(rawDoc bsoncore.Document) (bsoncore.Document, error) {
	elems, err := rawDoc.Elements()
	if err != nil {
		return nil, err
	}

	var filteredElems []byte
	for _, elem := range elems {
		key := elem.Key()
		if key != "name" {
			filteredElems = append(filteredElems, elem...)
			continue
		}

		name := elem.Value().StringValue()
		collName := name[strings.Index(name, ".")+1:]
		filteredElems = bsoncore.AppendStringElement(filteredElems, "name", collName)
	}

	var filteredDoc []byte
	filteredDoc = bsoncore.BuildDocument(filteredDoc, filteredElems)
	return filteredDoc, nil
}
