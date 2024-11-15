// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrNoDocuments is returned by SingleResult methods when the operation that created the SingleResult did not return
// any documents.
var ErrNoDocuments = errors.New("mongo: no documents in result")

// SingleResult represents a single document returned from an operation. If the operation resulted in an error, all
// SingleResult methods will return that error. If the operation did not return any documents, all SingleResult methods
// will return ErrNoDocuments.
type SingleResult struct {
	ctx      context.Context
	err      error
	cur      *Cursor
	rdr      bson.Raw
	bsonOpts *options.BSONOptions
	reg      *bsoncodec.Registry
}

// NewSingleResultFromDocument creates a SingleResult with the provided error, registry, and an underlying Cursor pre-loaded with
// the provided document, error and registry. If no registry is provided, bson.DefaultRegistry will be used. If an error distinct
// from the one provided occurs during creation of the SingleResult, that error will be stored on the returned SingleResult.
//
// The document parameter must be a non-nil document.
func NewSingleResultFromDocument(document interface{}, err error, registry *bsoncodec.Registry) *SingleResult {
	if document == nil {
		return &SingleResult{err: ErrNilDocument}
	}
	if registry == nil {
		registry = bson.DefaultRegistry
	}

	cur, createErr := NewCursorFromDocuments([]interface{}{document}, err, registry)
	if createErr != nil {
		return &SingleResult{err: createErr}
	}

	return &SingleResult{
		cur: cur,
		err: err,
		reg: registry,
	}
}

// Decode will unmarshal the document represented by this SingleResult into v. If there was an error from the operation
// that created this SingleResult, that error will be returned. If the operation returned no documents, Decode will
// return ErrNoDocuments.
//
// If the operation was successful and returned a document, Decode will return any errors from the unmarshalling process
// without any modification. If v is nil or is a typed nil, an error will be returned.
func (sr *SingleResult) Decode(v interface{}) error {
	if sr.err != nil {
		return sr.err
	}
	if sr.reg == nil {
		return bson.ErrNilRegistry
	}

	if sr.err = sr.setRdrContents(); sr.err != nil {
		return sr.err
	}

	dec, err := getDecoder(sr.rdr, sr.bsonOpts, sr.reg)
	if err != nil {
		return fmt.Errorf("error configuring BSON decoder: %w", err)
	}

	return dec.Decode(v)
}

// Raw returns the document represented by this SingleResult as a bson.Raw. If
// there was an error from the operation that created this SingleResult, both
// the result and that error will be returned. If the operation returned no
// documents, this will return (nil, ErrNoDocuments).
func (sr *SingleResult) Raw() (bson.Raw, error) {
	if sr.err != nil {
		return sr.rdr, sr.err
	}

	if sr.err = sr.setRdrContents(); sr.err != nil {
		return nil, sr.err
	}
	return sr.rdr, nil
}

// DecodeBytes will return the document represented by this SingleResult as a bson.Raw. If there was an error from the
// operation that created this SingleResult, both the result and that error will be returned. If the operation returned
// no documents, this will return (nil, ErrNoDocuments).
//
// Deprecated: Use [SingleResult.Raw] instead.
func (sr *SingleResult) DecodeBytes() (bson.Raw, error) {
	return sr.Raw()
}

// setRdrContents will set the contents of rdr by iterating the underlying cursor if necessary.
func (sr *SingleResult) setRdrContents() error {
	switch {
	case sr.err != nil:
		return sr.err
	case sr.rdr != nil:
		return nil
	case sr.cur != nil:
		defer sr.cur.Close(sr.ctx)

		if !sr.cur.Next(sr.ctx) {
			if err := sr.cur.Err(); err != nil {
				return err
			}

			return ErrNoDocuments
		}
		sr.rdr = sr.cur.Current
		return nil
	}

	return ErrNoDocuments
}

// Err provides a way to check for query errors without calling Decode. Err returns the error, if
// any, that was encountered while running the operation. If the operation was successful but did
// not return any documents, Err returns ErrNoDocuments. If this error is not nil, this error will
// also be returned from Decode.
func (sr *SingleResult) Err() error {
	sr.err = sr.setRdrContents()

	return sr.err
}
