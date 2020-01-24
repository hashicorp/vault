// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package writeconcern // import "go.mongodb.org/mongo-driver/mongo/writeconcern"

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ErrInconsistent indicates that an inconsistent write concern was specified.
var ErrInconsistent = errors.New("a write concern cannot have both w=0 and j=true")

// ErrEmptyWriteConcern indicates that a write concern has no fields set.
var ErrEmptyWriteConcern = errors.New("a write concern must have at least one field set")

// ErrNegativeW indicates that a negative integer `w` field was specified.
var ErrNegativeW = errors.New("write concern `w` field cannot be a negative number")

// ErrNegativeWTimeout indicates that a negative WTimeout was specified.
var ErrNegativeWTimeout = errors.New("write concern `wtimeout` field cannot be negative")

// WriteConcern describes the level of acknowledgement requested from MongoDB for write operations
// to a standalone mongod or to replica sets or to sharded clusters.
type WriteConcern struct {
	w        interface{}
	j        bool
	wTimeout time.Duration
}

// Option is an option to provide when creating a WriteConcern.
type Option func(concern *WriteConcern)

// New constructs a new WriteConcern.
func New(options ...Option) *WriteConcern {
	concern := &WriteConcern{}

	for _, option := range options {
		option(concern)
	}

	return concern
}

// W requests acknowledgement that write operations propagate to the specified number of mongod
// instances.
func W(w int) Option {
	return func(concern *WriteConcern) {
		concern.w = w
	}
}

// WMajority requests acknowledgement that write operations propagate to the majority of mongod
// instances.
func WMajority() Option {
	return func(concern *WriteConcern) {
		concern.w = "majority"
	}
}

// WTagSet requests acknowledgement that write operations propagate to the specified mongod
// instance.
func WTagSet(tag string) Option {
	return func(concern *WriteConcern) {
		concern.w = tag
	}
}

// J requests acknowledgement from MongoDB that write operations are written to
// the journal.
func J(j bool) Option {
	return func(concern *WriteConcern) {
		concern.j = j
	}
}

// WTimeout specifies specifies a time limit for the write concern.
func WTimeout(d time.Duration) Option {
	return func(concern *WriteConcern) {
		concern.wTimeout = d
	}
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (wc *WriteConcern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !wc.IsValid() {
		return bsontype.Type(0), nil, ErrInconsistent
	}

	var elems []byte

	if wc.w != nil {
		switch t := wc.w.(type) {
		case int:
			if t < 0 {
				return bsontype.Type(0), nil, ErrNegativeW
			}

			elems = bsoncore.AppendInt32Element(elems, "w", int32(t))
		case string:
			elems = bsoncore.AppendStringElement(elems, "w", string(t))
		}
	}

	if wc.j {
		elems = bsoncore.AppendBooleanElement(elems, "j", wc.j)
	}

	if wc.wTimeout < 0 {
		return bsontype.Type(0), nil, ErrNegativeWTimeout
	}

	if wc.wTimeout != 0 {
		elems = bsoncore.AppendInt64Element(elems, "wtimeout", int64(wc.wTimeout/time.Millisecond))
	}

	if len(elems) == 0 {
		return bsontype.Type(0), nil, ErrEmptyWriteConcern
	}
	return bsontype.EmbeddedDocument, bsoncore.BuildDocument(nil, elems), nil
}

// AcknowledgedValue returns true if a BSON RawValue for a write concern represents an acknowledged write concern.
// The element's value must be a document representing a write concern.
func AcknowledgedValue(rawv bson.RawValue) bool {
	doc, ok := bsoncore.Value{Type: rawv.Type, Data: rawv.Value}.DocumentOK()
	if !ok {
		return false
	}

	val, err := doc.LookupErr("w")
	if err != nil {
		// key w not found --> acknowledged
		return true
	}

	i32, ok := val.Int32OK()
	if !ok {
		return false
	}
	return i32 != 0
}

// Acknowledged indicates whether or not a write with the given write concern will be acknowledged.
func (wc *WriteConcern) Acknowledged() bool {
	if wc == nil || wc.j {
		return true
	}

	switch v := wc.w.(type) {
	case int:
		if v == 0 {
			return false
		}
	}

	return true
}

// IsValid checks whether the write concern is invalid.
func (wc *WriteConcern) IsValid() bool {
	if !wc.j {
		return true
	}

	switch v := wc.w.(type) {
	case int:
		if v == 0 {
			return false
		}
	}

	return true
}

// GetW returns the write concern w level.
func (wc *WriteConcern) GetW() interface{} {
	return wc.w
}

// GetJ returns the write concern journaling level.
func (wc *WriteConcern) GetJ() bool {
	return wc.j
}

// GetWTimeout returns the write concern timeout.
func (wc *WriteConcern) GetWTimeout() time.Duration {
	return wc.wTimeout
}

// WithOptions returns a copy of this WriteConcern with the options set.
func (wc *WriteConcern) WithOptions(options ...Option) *WriteConcern {
	if wc == nil {
		return New(options...)
	}
	newWC := &WriteConcern{}
	*newWC = *wc

	for _, option := range options {
		option(newWC)
	}

	return newWC
}

// AckWrite returns true if a write concern represents an acknowledged write
func AckWrite(wc *WriteConcern) bool {
	return wc == nil || wc.Acknowledged()
}
