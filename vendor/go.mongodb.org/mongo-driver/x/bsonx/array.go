// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx // import "go.mongodb.org/mongo-driver/x/bsonx"

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ErrNilArray indicates that an operation was attempted on a nil *Array.
var ErrNilArray = errors.New("array is nil")

// Arr represents an array in BSON.
type Arr []Val

// String implements the fmt.Stringer interface.
func (a Arr) String() string {
	var buf bytes.Buffer
	buf.Write([]byte("bson.Array["))
	for idx, val := range a {
		if idx > 0 {
			buf.Write([]byte(", "))
		}
		fmt.Fprintf(&buf, "%s", val)
	}
	buf.WriteByte(']')

	return buf.String()
}

// MarshalBSONValue implements the bsoncodec.ValueMarshaler interface.
func (a Arr) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if a == nil {
		// TODO: Should we do this?
		return bsontype.Null, nil, nil
	}

	idx, dst := bsoncore.ReserveLength(nil)
	for idx, value := range a {
		t, data, _ := value.MarshalBSONValue() // marshalBSONValue never returns an error.
		dst = append(dst, byte(t))
		dst = append(dst, strconv.Itoa(idx)...)
		dst = append(dst, 0x00)
		dst = append(dst, data...)
	}
	dst = append(dst, 0x00)
	dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	return bsontype.Array, dst, nil
}

// UnmarshalBSONValue implements the bsoncodec.ValueUnmarshaler interface.
func (a *Arr) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if a == nil {
		return ErrNilArray
	}
	*a = (*a)[:0]

	elements, err := bsoncore.Document(data).Elements()
	if err != nil {
		return err
	}

	for _, elem := range elements {
		var val Val
		rawval := elem.Value()
		err = val.UnmarshalBSONValue(rawval.Type, rawval.Data)
		if err != nil {
			return err
		}
		*a = append(*a, val)
	}
	return nil
}

// Equal compares this document to another, returning true if they are equal.
func (a Arr) Equal(a2 Arr) bool {
	if len(a) != len(a2) {
		return false
	}
	for idx := range a {
		if !a[idx].Equal(a2[idx]) {
			return false
		}
	}
	return true
}

func (Arr) idoc() {}
