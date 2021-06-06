// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"encoding/binary"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IDoc is the interface implemented by Doc and MDoc. It allows either of these types to be provided
// to the Document function to create a Value.
type IDoc interface {
	idoc()
}

// Double constructs a BSON double Value.
func Double(f64 float64) Val {
	v := Val{t: bsontype.Double}
	binary.LittleEndian.PutUint64(v.bootstrap[0:8], math.Float64bits(f64))
	return v
}

// String constructs a BSON string Value.
func String(str string) Val { return Val{t: bsontype.String}.writestring(str) }

// Document constructs a Value from the given IDoc. If nil is provided, a BSON Null value will be
// returned.
func Document(doc IDoc) Val {
	var v Val
	switch tt := doc.(type) {
	case Doc:
		if tt == nil {
			v.t = bsontype.Null
			break
		}
		v.t = bsontype.EmbeddedDocument
		v.primitive = tt
	case MDoc:
		if tt == nil {
			v.t = bsontype.Null
			break
		}
		v.t = bsontype.EmbeddedDocument
		v.primitive = tt
	default:
		v.t = bsontype.Null
	}
	return v
}

// Array constructs a Value from arr. If arr is nil, a BSON Null value is returned.
func Array(arr Arr) Val {
	if arr == nil {
		return Val{t: bsontype.Null}
	}
	return Val{t: bsontype.Array, primitive: arr}
}

// Binary constructs a BSON binary Value.
func Binary(subtype byte, data []byte) Val {
	return Val{t: bsontype.Binary, primitive: primitive.Binary{Subtype: subtype, Data: data}}
}

// Undefined constructs a BSON binary Value.
func Undefined() Val { return Val{t: bsontype.Undefined} }

// ObjectID constructs a BSON objectid Value.
func ObjectID(oid primitive.ObjectID) Val {
	v := Val{t: bsontype.ObjectID}
	copy(v.bootstrap[0:12], oid[:])
	return v
}

// Boolean constructs a BSON boolean Value.
func Boolean(b bool) Val {
	v := Val{t: bsontype.Boolean}
	if b {
		v.bootstrap[0] = 0x01
	}
	return v
}

// DateTime constructs a BSON datetime Value.
func DateTime(dt int64) Val { return Val{t: bsontype.DateTime}.writei64(dt) }

// Time constructs a BSON datetime Value.
func Time(t time.Time) Val {
	return Val{t: bsontype.DateTime}.writei64(t.Unix()*1e3 + int64(t.Nanosecond()/1e6))
}

// Null constructs a BSON binary Value.
func Null() Val { return Val{t: bsontype.Null} }

// Regex constructs a BSON regex Value.
func Regex(pattern, options string) Val {
	regex := primitive.Regex{Pattern: pattern, Options: options}
	return Val{t: bsontype.Regex, primitive: regex}
}

// DBPointer constructs a BSON dbpointer Value.
func DBPointer(ns string, ptr primitive.ObjectID) Val {
	dbptr := primitive.DBPointer{DB: ns, Pointer: ptr}
	return Val{t: bsontype.DBPointer, primitive: dbptr}
}

// JavaScript constructs a BSON javascript Value.
func JavaScript(js string) Val {
	return Val{t: bsontype.JavaScript}.writestring(js)
}

// Symbol constructs a BSON symbol Value.
func Symbol(symbol string) Val {
	return Val{t: bsontype.Symbol}.writestring(symbol)
}

// CodeWithScope constructs a BSON code with scope Value.
func CodeWithScope(code string, scope IDoc) Val {
	cws := primitive.CodeWithScope{Code: primitive.JavaScript(code), Scope: scope}
	return Val{t: bsontype.CodeWithScope, primitive: cws}
}

// Int32 constructs a BSON int32 Value.
func Int32(i32 int32) Val {
	v := Val{t: bsontype.Int32}
	v.bootstrap[0] = byte(i32)
	v.bootstrap[1] = byte(i32 >> 8)
	v.bootstrap[2] = byte(i32 >> 16)
	v.bootstrap[3] = byte(i32 >> 24)
	return v
}

// Timestamp constructs a BSON timestamp Value.
func Timestamp(t, i uint32) Val {
	v := Val{t: bsontype.Timestamp}
	v.bootstrap[0] = byte(i)
	v.bootstrap[1] = byte(i >> 8)
	v.bootstrap[2] = byte(i >> 16)
	v.bootstrap[3] = byte(i >> 24)
	v.bootstrap[4] = byte(t)
	v.bootstrap[5] = byte(t >> 8)
	v.bootstrap[6] = byte(t >> 16)
	v.bootstrap[7] = byte(t >> 24)
	return v
}

// Int64 constructs a BSON int64 Value.
func Int64(i64 int64) Val { return Val{t: bsontype.Int64}.writei64(i64) }

// Decimal128 constructs a BSON decimal128 Value.
func Decimal128(d128 primitive.Decimal128) Val {
	return Val{t: bsontype.Decimal128, primitive: d128}
}

// MinKey constructs a BSON minkey Value.
func MinKey() Val { return Val{t: bsontype.MinKey} }

// MaxKey constructs a BSON maxkey Value.
func MaxKey() Val { return Val{t: bsontype.MaxKey} }
