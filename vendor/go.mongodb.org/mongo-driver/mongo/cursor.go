// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"io"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Cursor is used to iterate over a stream of documents. Each document can be decoded into a Go type via the Decode
// method or accessed as raw BSON via the Current field.
type Cursor struct {
	// Current contains the BSON bytes of the current change document. This property is only valid until the next call
	// to Next or TryNext. If continued access is required, a copy must be made.
	Current bson.Raw

	bc            batchCursor
	batch         *bsoncore.DocumentSequence
	registry      *bsoncodec.Registry
	clientSession *session.Client

	err error
}

func newCursor(bc batchCursor, registry *bsoncodec.Registry) (*Cursor, error) {
	return newCursorWithSession(bc, registry, nil)
}

func newCursorWithSession(bc batchCursor, registry *bsoncodec.Registry, clientSession *session.Client) (*Cursor, error) {
	if registry == nil {
		registry = bson.DefaultRegistry
	}
	if bc == nil {
		return nil, errors.New("batch cursor must not be nil")
	}
	c := &Cursor{
		bc:            bc,
		registry:      registry,
		clientSession: clientSession,
	}
	if bc.ID() == 0 {
		c.closeImplicitSession()
	}
	return c, nil
}

func newEmptyCursor() *Cursor {
	return &Cursor{bc: driver.NewEmptyBatchCursor()}
}

// ID returns the ID of this cursor, or 0 if the cursor has been closed or exhausted.
func (c *Cursor) ID() int64 { return c.bc.ID() }

// Next gets the next document for this cursor. It returns true if there were no errors and the cursor has not been
// exhausted.
//
// Next blocks until a document is available, an error occurs, or ctx expires. If ctx expires, the
// error will be set to ctx.Err(). In an error case, Next will return false.
//
// If Next returns false, subsequent calls will also return false.
func (c *Cursor) Next(ctx context.Context) bool {
	return c.next(ctx, false)
}

// TryNext attempts to get the next document for this cursor. It returns true if there were no errors and the next
// document is available. This is only recommended for use with tailable cursors as a non-blocking alternative to
// Next. See https://docs.mongodb.com/manual/core/tailable-cursors/ for more information about tailable cursors.
//
// TryNext returns false if the cursor is exhausted, an error occurs when getting results from the server, the next
// document is not yet available, or ctx expires. If ctx expires, the error will be set to ctx.Err().
//
// If TryNext returns false and an error occurred or the cursor has been exhausted (i.e. c.Err() != nil || c.ID() == 0),
// subsequent attempts will also return false. Otherwise, it is safe to call TryNext again until a document is
// available.
//
// This method requires driver version >= 1.2.0.
func (c *Cursor) TryNext(ctx context.Context) bool {
	return c.next(ctx, true)
}

func (c *Cursor) next(ctx context.Context, nonBlocking bool) bool {
	// return false right away if the cursor has already errored.
	if c.err != nil {
		return false
	}

	if ctx == nil {
		ctx = context.Background()
	}
	doc, err := c.batch.Next()
	switch err {
	case nil:
		c.Current = bson.Raw(doc)
		return true
	case io.EOF: // Need to do a getMore
	default:
		c.err = err
		return false
	}

	// call the Next method in a loop until at least one document is returned in the next batch or
	// the context times out.
	for {
		// If we don't have a next batch
		if !c.bc.Next(ctx) {
			// Do we have an error? If so we return false.
			c.err = c.bc.Err()
			if c.err != nil {
				return false
			}
			// Is the cursor ID zero?
			if c.bc.ID() == 0 {
				c.closeImplicitSession()
				return false
			}
			// empty batch, but cursor is still valid.
			// use nonBlocking to determine if we should continue or return control to the caller.
			if nonBlocking {
				return false
			}
			continue
		}

		// close the implicit session if this was the last getMore
		if c.bc.ID() == 0 {
			c.closeImplicitSession()
		}

		c.batch = c.bc.Batch()
		doc, err = c.batch.Next()
		switch err {
		case nil:
			c.Current = bson.Raw(doc)
			return true
		case io.EOF: // Empty batch so we continue
		default:
			c.err = err
			return false
		}
	}
}

// Decode will unmarshal the current document into val and return any errors from the unmarshalling process without any
// modification. If val is nil or is a typed nil, an error will be returned.
func (c *Cursor) Decode(val interface{}) error {
	return bson.UnmarshalWithRegistry(c.registry, c.Current, val)
}

// Err returns the last error seen by the Cursor, or nil if no error has occurred.
func (c *Cursor) Err() error { return c.err }

// Close closes this cursor. Next and TryNext must not be called after Close has been called. Close is idempotent. After
// the first call, any subsequent calls will not change the state.
func (c *Cursor) Close(ctx context.Context) error {
	defer c.closeImplicitSession()
	return c.bc.Close(ctx)
}

// All iterates the cursor and decodes each document into results. The results parameter must be a pointer to a slice.
// The slice pointed to by results will be completely overwritten. This method will close the cursor after retrieving
// all documents. If the cursor has been iterated, any previously iterated documents will not be included in results.
//
// This method requires driver version >= 1.1.0.
func (c *Cursor) All(ctx context.Context, results interface{}) error {
	resultsVal := reflect.ValueOf(results)
	if resultsVal.Kind() != reflect.Ptr {
		return errors.New("results argument must be a pointer to a slice")
	}

	sliceVal := resultsVal.Elem()
	elementType := sliceVal.Type().Elem()
	var index int
	var err error

	defer c.Close(ctx)

	batch := c.batch // exhaust the current batch before iterating the batch cursor
	for {
		sliceVal, index, err = c.addFromBatch(sliceVal, elementType, batch, index)
		if err != nil {
			return err
		}

		if !c.bc.Next(ctx) {
			break
		}

		batch = c.bc.Batch()
	}

	if err = c.bc.Err(); err != nil {
		return err
	}

	resultsVal.Elem().Set(sliceVal.Slice(0, index))
	return nil
}

// addFromBatch adds all documents from batch to sliceVal starting at the given index. It returns the new slice value,
// the next empty index in the slice, and an error if one occurs.
func (c *Cursor) addFromBatch(sliceVal reflect.Value, elemType reflect.Type, batch *bsoncore.DocumentSequence,
	index int) (reflect.Value, int, error) {

	docs, err := batch.Documents()
	if err != nil {
		return sliceVal, index, err
	}

	for _, doc := range docs {
		if sliceVal.Len() == index {
			// slice is full
			newElem := reflect.New(elemType)
			sliceVal = reflect.Append(sliceVal, newElem.Elem())
			sliceVal = sliceVal.Slice(0, sliceVal.Cap())
		}

		currElem := sliceVal.Index(index).Addr().Interface()
		if err = bson.UnmarshalWithRegistry(c.registry, doc, currElem); err != nil {
			return sliceVal, index, err
		}

		index++
	}

	return sliceVal, index, nil
}

func (c *Cursor) closeImplicitSession() {
	if c.clientSession != nil && c.clientSession.SessionType == session.Implicit {
		c.clientSession.EndSession()
	}
}

// BatchCursorFromCursor returns a driver.BatchCursor for the given Cursor. If there is no underlying
// driver.BatchCursor, nil is returned. This method is deprecated and does not have any stability guarantees. It may be
// removed in the future.
func BatchCursorFromCursor(c *Cursor) *driver.BatchCursor {
	bc, _ := c.bc.(*driver.BatchCursor)
	return bc
}
