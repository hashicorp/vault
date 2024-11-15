// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package writeconcern defines write concerns for MongoDB operations.
//
// For more information about MongoDB write concerns, see
// https://www.mongodb.com/docs/manual/reference/write-concern/
package writeconcern // import "go.mongodb.org/mongo-driver/mongo/writeconcern"

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const majority = "majority"

// ErrInconsistent indicates that an inconsistent write concern was specified.
//
// Deprecated: ErrInconsistent will be removed in Go Driver 2.0.
var ErrInconsistent = errors.New("a write concern cannot have both w=0 and j=true")

// ErrEmptyWriteConcern indicates that a write concern has no fields set.
//
// Deprecated: ErrEmptyWriteConcern will be removed in Go Driver 2.0.
var ErrEmptyWriteConcern = errors.New("a write concern must have at least one field set")

// ErrNegativeW indicates that a negative integer `w` field was specified.
//
// Deprecated: ErrNegativeW will be removed in Go Driver 2.0.
var ErrNegativeW = errors.New("write concern `w` field cannot be a negative number")

// ErrNegativeWTimeout indicates that a negative WTimeout was specified.
//
// Deprecated: ErrNegativeWTimeout will be removed in Go Driver 2.0.
var ErrNegativeWTimeout = errors.New("write concern `wtimeout` field cannot be negative")

// A WriteConcern defines a MongoDB write concern, which describes the level of acknowledgment
// requested from MongoDB for write operations to a standalone mongod, to replica sets, or to
// sharded clusters.
//
// For more information about MongoDB write concerns, see
// https://www.mongodb.com/docs/manual/reference/write-concern/
type WriteConcern struct {
	// W requests acknowledgment that the write operation has propagated to a
	// specified number of mongod instances or to mongod instances with
	// specified tags. It sets the "w" option in a MongoDB write concern.
	//
	// W values must be a string or an int.
	//
	// Common values are:
	//   - "majority": requests acknowledgment that write operations have been
	//     durably committed to the calculated majority of the data-bearing
	//     voting members.
	//   - 1: requests acknowledgment that write operations have been written
	//     to 1 node.
	//   - 0: requests no acknowledgment of write operations
	//
	// For more information about the "w" option, see
	// https://www.mongodb.com/docs/manual/reference/write-concern/#w-option
	W interface{}

	// Journal requests acknowledgment from MongoDB that the write operation has
	// been written to the on-disk journal. It sets the "j" option in a MongoDB
	// write concern.
	//
	// For more information about the "j" option, see
	// https://www.mongodb.com/docs/manual/reference/write-concern/#j-option
	Journal *bool

	// WTimeout specifies a time limit for the write concern. It sets the
	// "wtimeout" option in a MongoDB write concern.
	//
	// It is only applicable for "w" values greater than 1. Using a WTimeout and
	// setting Timeout on the Client at the same time will result in undefined
	// behavior.
	//
	// For more information about the "wtimeout" option, see
	// https://www.mongodb.com/docs/manual/reference/write-concern/#wtimeout
	WTimeout time.Duration
}

// Unacknowledged returns a WriteConcern that requests no acknowledgment of
// write operations.
//
// For more information about write concern "w: 0", see
// https://www.mongodb.com/docs/manual/reference/write-concern/#mongodb-writeconcern-writeconcern.-number-
func Unacknowledged() *WriteConcern {
	return &WriteConcern{W: 0}
}

// W1 returns a WriteConcern that requests acknowledgment that write operations
// have been written to memory on one node (e.g. the standalone mongod or the
// primary in a replica set).
//
// For more information about write concern "w: 1", see
// https://www.mongodb.com/docs/manual/reference/write-concern/#mongodb-writeconcern-writeconcern.-number-
func W1() *WriteConcern {
	return &WriteConcern{W: 1}
}

// Journaled returns a WriteConcern that requests acknowledgment that write
// operations have been written to the on-disk journal on MongoDB.
//
// The database's default value for "w" determines how many nodes must write to
// their on-disk journal before the write operation is acknowledged.
//
// For more information about write concern "j: true", see
// https://www.mongodb.com/docs/manual/reference/write-concern/#mongodb-writeconcern-ournal
func Journaled() *WriteConcern {
	journal := true
	return &WriteConcern{Journal: &journal}
}

// Majority returns a WriteConcern that requests acknowledgment that write
// operations have been durably committed to the calculated majority of the
// data-bearing voting members.
//
// Write concern "w: majority" typically requires write operations to be written
// to the on-disk journal before they are acknowledged, unless journaling is
// disabled on MongoDB or the "writeConcernMajorityJournalDefault" replica set
// configuration is set to false.
//
// For more information about write concern "w: majority", see
// https://www.mongodb.com/docs/manual/reference/write-concern/#mongodb-writeconcern-writeconcern.-majority-
func Majority() *WriteConcern {
	return &WriteConcern{W: majority}
}

// Custom returns a WriteConcern that requests acknowledgment that write
// operations have propagated to tagged members that satisfy the custom write
// concern defined in "settings.getLastErrorModes".
//
// For more information about custom write concern names, see
// https://www.mongodb.com/docs/manual/reference/write-concern/#mongodb-writeconcern-writeconcern.-custom-write-concern-name-
func Custom(tag string) *WriteConcern {
	return &WriteConcern{W: tag}
}

// Option is an option to provide when creating a WriteConcern.
//
// Deprecated: Use the WriteConcern convenience functions or define a struct literal instead.
// For example:
//
//	writeconcern.Majority()
//
// or
//
//	journal := true
//	&writeconcern.WriteConcern{
//		W:       2,
//		Journal: &journal,
//	}
type Option func(concern *WriteConcern)

// New constructs a new WriteConcern.
//
// Deprecated: Use the WriteConcern convenience functions or define a struct literal instead.
// For example:
//
//	writeconcern.Majority()
//
// or
//
//	journal := true
//	&writeconcern.WriteConcern{
//		W:       2,
//		Journal: &journal,
//	}
func New(options ...Option) *WriteConcern {
	concern := &WriteConcern{}

	for _, option := range options {
		option(concern)
	}

	return concern
}

// W requests acknowledgement that write operations propagate to the specified number of mongod
// instances.
//
// Deprecated: Use the Unacknowledged or W1 functions or define a struct literal instead.
// For example:
//
//	writeconcern.Unacknowledged()
//
// or
//
//	journal := true
//	&writeconcern.WriteConcern{
//		W:       2,
//		Journal: &journal,
//	}
func W(w int) Option {
	return func(concern *WriteConcern) {
		concern.W = w
	}
}

// WMajority requests acknowledgement that write operations propagate to the majority of mongod
// instances.
//
// Deprecated: Use [Majority] instead.
func WMajority() Option {
	return func(concern *WriteConcern) {
		concern.W = majority
	}
}

// WTagSet requests acknowledgement that write operations propagate to the specified mongod
// instance.
//
// Deprecated: Use [Custom] instead.
func WTagSet(tag string) Option {
	return func(concern *WriteConcern) {
		concern.W = tag
	}
}

// J requests acknowledgement from MongoDB that write operations are written to
// the journal.
//
// Deprecated: Use the Journaled function or define a struct literal instead.
// For example:
//
//	writeconcern.Journaled()
//
// or
//
//	journal := true
//	&writeconcern.WriteConcern{
//		W:       2,
//		Journal: &journal,
//	}
func J(j bool) Option {
	return func(concern *WriteConcern) {
		// To maintain backward compatible behavior (now that the J field is a
		// bool pointer), only set a value for J if the input is true. If the
		// input is false, do not set a value, which omits "j" from the
		// marshaled write concern.
		if j {
			concern.Journal = &j
		}
	}
}

// WTimeout specifies a time limit for the write concern.
//
// It is only applicable for "w" values greater than 1. Using a WTimeout and setting Timeout on the
// Client at the same time will result in undefined behavior.
//
// Deprecated: Use the WriteConcern convenience functions or define a struct literal instead.
// For example:
//
//	wc := writeconcern.W1()
//	wc.WTimeout = 30 * time.Second
//
// or
//
//	journal := true
//	&writeconcern.WriteConcern{
//		W:        "majority",
//		WTimeout: 30 * time.Second,
//	}
func WTimeout(d time.Duration) Option {
	return func(concern *WriteConcern) {
		concern.WTimeout = d
	}
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
//
// Deprecated: Marshaling a WriteConcern to BSON will not be supported in Go
// Driver 2.0.
func (wc *WriteConcern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if wc == nil {
		return 0, nil, ErrEmptyWriteConcern
	}

	var elems []byte
	if wc.W != nil {
		// Only support string or int values for W. That aligns with the
		// documentation and the behavior of other functions, like Acknowledged.
		switch w := wc.W.(type) {
		case int:
			if w < 0 {
				return 0, nil, ErrNegativeW
			}

			// If Journal=true and W=0, return an error because that write
			// concern is ambiguous.
			if wc.Journal != nil && *wc.Journal && w == 0 {
				return 0, nil, ErrInconsistent
			}

			elems = bsoncore.AppendInt32Element(elems, "w", int32(w))
		case string:
			elems = bsoncore.AppendStringElement(elems, "w", w)
		default:
			return 0,
				nil,
				fmt.Errorf("WriteConcern.W must be a string or int, but is a %T", wc.W)
		}
	}

	if wc.Journal != nil {
		elems = bsoncore.AppendBooleanElement(elems, "j", *wc.Journal)
	}

	if wc.WTimeout < 0 {
		return 0, nil, ErrNegativeWTimeout
	}

	if wc.WTimeout != 0 {
		elems = bsoncore.AppendInt64Element(elems, "wtimeout", int64(wc.WTimeout/time.Millisecond))
	}

	if len(elems) == 0 {
		return 0, nil, ErrEmptyWriteConcern
	}
	return bson.TypeEmbeddedDocument, bsoncore.BuildDocument(nil, elems), nil
}

// AcknowledgedValue returns true if a BSON RawValue for a write concern represents an acknowledged write concern.
// The element's value must be a document representing a write concern.
//
// Deprecated: AcknowledgedValue will not be supported in Go Driver 2.0.
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
	// Only {w: 0} or {w: 0, j: false} are an unacknowledged write concerns. All other values are
	// acknowledged.
	return wc == nil || wc.W != 0 || (wc.Journal != nil && *wc.Journal)
}

// IsValid returns true if the WriteConcern is valid.
func (wc *WriteConcern) IsValid() bool {
	if wc == nil {
		return true
	}

	switch w := wc.W.(type) {
	case int:
		// A write concern with {w: int} must have a non-negative value and
		// cannot have the combination {w: 0, j: true}.
		return w >= 0 && (w > 0 || wc.Journal == nil || !*wc.Journal)
	case string, nil:
		// A write concern with {w: string} or no w specified is always valid.
		return true
	default:
		// A write concern with an unsupported w type is not valid.
		return false
	}
}

// GetW returns the write concern w level.
//
// Deprecated: Use the WriteConcern.W field instead.
func (wc *WriteConcern) GetW() interface{} {
	return wc.W
}

// GetJ returns the write concern journaling level.
//
// Deprecated: Use the WriteConcern.Journal field instead.
func (wc *WriteConcern) GetJ() bool {
	// Treat a nil Journal as false. That maintains backward compatibility with the existing
	// behavior of GetJ where unset is false. If users want the real value of Journal, they can
	// access the Journal field.
	return wc.Journal != nil && *wc.Journal
}

// GetWTimeout returns the write concern timeout.
//
// Deprecated: Use the WriteConcern.WTimeout field instead.
func (wc *WriteConcern) GetWTimeout() time.Duration {
	return wc.WTimeout
}

// WithOptions returns a copy of this WriteConcern with the options set.
//
// Deprecated: Use the WriteConcern convenience functions or define a struct literal instead.
// For example:
//
//	writeconcern.Majority()
//
// or
//
//	journal := true
//	&writeconcern.WriteConcern{
//		W:       2,
//		Journal: &journal,
//	}
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
//
// Deprecated: Use [WriteConcern.Acknowledged] instead.
func AckWrite(wc *WriteConcern) bool {
	return wc == nil || wc.Acknowledged()
}
