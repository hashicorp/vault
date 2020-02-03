// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"bytes"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// MDoc is an unordered, type safe, concise BSON document representation. This type should not be
// used if you require ordering of values or duplicate keys.
type MDoc map[string]Val

// ReadMDoc will create a Doc using the provided slice of bytes. If the
// slice of bytes is not a valid BSON document, this method will return an error.
func ReadMDoc(b []byte) (MDoc, error) {
	doc := make(MDoc, 0)
	err := doc.UnmarshalBSON(b)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Copy makes a shallow copy of this document.
func (d MDoc) Copy() MDoc {
	d2 := make(MDoc, len(d))
	for k, v := range d {
		d2[k] = v
	}
	return d2
}

// Lookup searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Value if they key does not exist. To know if they key actually
// exists, use LookupErr.
func (d MDoc) Lookup(key ...string) Val {
	val, _ := d.LookupErr(key...)
	return val
}

// LookupErr searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
func (d MDoc) LookupErr(key ...string) (Val, error) {
	elem, err := d.LookupElementErr(key...)
	return elem.Value, err
}

// LookupElement searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Element if they key does not exist. To know if they key actually
// exists, use LookupElementErr.
func (d MDoc) LookupElement(key ...string) Elem {
	elem, _ := d.LookupElementErr(key...)
	return elem
}

// LookupElementErr searches the document and potentially subdocuments for the
// provided key. Each key provided to this method represents a layer of depth.
func (d MDoc) LookupElementErr(key ...string) (Elem, error) {
	// KeyNotFound operates by being created where the error happens and then the depth is
	// incremented by 1 as each function unwinds. Whenever this function returns, it also assigns
	// the Key slice to the key slice it has. This ensures that the proper depth is identified and
	// the proper keys.
	if len(key) == 0 {
		return Elem{}, KeyNotFound{Key: key}
	}

	var elem Elem
	var err error
	val, ok := d[key[0]]
	if !ok {
		return Elem{}, KeyNotFound{Key: key}
	}

	if len(key) == 1 {
		return Elem{Key: key[0], Value: val}, nil
	}

	switch val.Type() {
	case bsontype.EmbeddedDocument:
		switch tt := val.primitive.(type) {
		case Doc:
			elem, err = tt.LookupElementErr(key[1:]...)
		case MDoc:
			elem, err = tt.LookupElementErr(key[1:]...)
		}
	default:
		return Elem{}, KeyNotFound{Type: val.Type()}
	}
	switch tt := err.(type) {
	case KeyNotFound:
		tt.Depth++
		tt.Key = key
		return Elem{}, tt
	case nil:
		return elem, nil
	default:
		return Elem{}, err // We can't actually hit this.
	}
}

// MarshalBSONValue implements the bsoncodec.ValueMarshaler interface.
//
// This method will never return an error.
func (d MDoc) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if d == nil {
		// TODO: Should we do this?
		return bsontype.Null, nil, nil
	}
	data, _ := d.MarshalBSON()
	return bsontype.EmbeddedDocument, data, nil
}

// MarshalBSON implements the Marshaler interface.
//
// This method will never return an error.
func (d MDoc) MarshalBSON() ([]byte, error) { return d.AppendMarshalBSON(nil) }

// AppendMarshalBSON marshals Doc to BSON bytes, appending to dst.
//
// This method will never return an error.
func (d MDoc) AppendMarshalBSON(dst []byte) ([]byte, error) {
	idx, dst := bsoncore.ReserveLength(dst)
	for k, v := range d {
		t, data, _ := v.MarshalBSONValue() // Value.MarshalBSONValue never returns an error.
		dst = append(dst, byte(t))
		dst = append(dst, k...)
		dst = append(dst, 0x00)
		dst = append(dst, data...)
	}
	dst = append(dst, 0x00)
	dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	return dst, nil
}

// UnmarshalBSON implements the Unmarshaler interface.
func (d *MDoc) UnmarshalBSON(b []byte) error {
	if d == nil {
		return ErrNilDocument
	}

	if err := bsoncore.Document(b).Validate(); err != nil {
		return err
	}

	elems, err := bsoncore.Document(b).Elements()
	if err != nil {
		return err
	}
	var val Val
	for _, elem := range elems {
		rawv := elem.Value()
		err = val.UnmarshalBSONValue(rawv.Type, rawv.Data)
		if err != nil {
			return err
		}
		(*d)[elem.Key()] = val
	}
	return nil
}

// Equal compares this document to another, returning true if they are equal.
func (d MDoc) Equal(id IDoc) bool {
	switch tt := id.(type) {
	case MDoc:
		d2 := tt
		if len(d) != len(d2) {
			return false
		}
		for key, value := range d {
			value2, ok := d2[key]
			if !ok {
				return false
			}
			if !value.Equal(value2) {
				return false
			}
		}
	case Doc:
		unique := make(map[string]struct{}, 0)
		for _, elem := range tt {
			unique[elem.Key] = struct{}{}
			val, ok := d[elem.Key]
			if !ok {
				return false
			}
			if !val.Equal(elem.Value) {
				return false
			}
		}
		if len(unique) != len(d) {
			return false
		}
	case nil:
		return d == nil
	default:
		return false
	}

	return true
}

// String implements the fmt.Stringer interface.
func (d MDoc) String() string {
	var buf bytes.Buffer
	buf.Write([]byte("bson.Document{"))
	first := true
	for key, value := range d {
		if !first {
			buf.Write([]byte(", "))
		}
		fmt.Fprintf(&buf, "%v", Elem{Key: key, Value: value})
		first = false
	}
	buf.WriteByte('}')

	return buf.String()
}

func (MDoc) idoc() {}
