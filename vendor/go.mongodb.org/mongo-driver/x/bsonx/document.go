// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"bytes"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ErrNilDocument indicates that an operation was attempted on a nil *bson.Document.
var ErrNilDocument = errors.New("document is nil")

// KeyNotFound is an error type returned from the Lookup methods on Document. This type contains
// information about which key was not found and if it was actually not found or if a component of
// the key except the last was not a document nor array.
type KeyNotFound struct {
	Key   []string      // The keys that were searched for.
	Depth uint          // Which key either was not found or was an incorrect type.
	Type  bsontype.Type // The type of the key that was found but was an incorrect type.
}

func (knf KeyNotFound) Error() string {
	depth := knf.Depth
	if depth >= uint(len(knf.Key)) {
		depth = uint(len(knf.Key)) - 1
	}

	if len(knf.Key) == 0 {
		return "no keys were provided for lookup"
	}

	if knf.Type != bsontype.Type(0) {
		return fmt.Sprintf(`key "%s" was found but was not valid to traverse BSON type %s`, knf.Key[depth], knf.Type)
	}

	return fmt.Sprintf(`key "%s" was not found`, knf.Key[depth])
}

// Doc is a type safe, concise BSON document representation.
type Doc []Elem

// ReadDoc will create a Document using the provided slice of bytes. If the
// slice of bytes is not a valid BSON document, this method will return an error.
func ReadDoc(b []byte) (Doc, error) {
	doc := make(Doc, 0)
	err := doc.UnmarshalBSON(b)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Copy makes a shallow copy of this document.
func (d Doc) Copy() Doc {
	d2 := make(Doc, len(d))
	copy(d2, d)
	return d2
}

// Append adds an element to the end of the document, creating it from the key and value provided.
func (d Doc) Append(key string, val Val) Doc {
	return append(d, Elem{Key: key, Value: val})
}

// Prepend adds an element to the beginning of the document, creating it from the key and value provided.
func (d Doc) Prepend(key string, val Val) Doc {
	// TODO: should we just modify d itself instead of doing an alloc here?
	return append(Doc{{Key: key, Value: val}}, d...)
}

// Set replaces an element of a document. If an element with a matching key is
// found, the element will be replaced with the one provided. If the document
// does not have an element with that key, the element is appended to the
// document instead.
func (d Doc) Set(key string, val Val) Doc {
	idx := d.IndexOf(key)
	if idx == -1 {
		return append(d, Elem{Key: key, Value: val})
	}
	d[idx] = Elem{Key: key, Value: val}
	return d
}

// IndexOf returns the index of the first element with a key of key, or -1 if no element with a key
// was found.
func (d Doc) IndexOf(key string) int {
	for i, e := range d {
		if e.Key == key {
			return i
		}
	}
	return -1
}

// Delete removes the element with key if it exists and returns the updated Doc.
func (d Doc) Delete(key string) Doc {
	idx := d.IndexOf(key)
	if idx == -1 {
		return d
	}
	return append(d[:idx], d[idx+1:]...)
}

// Lookup searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Value if they key does not exist. To know if they key actually
// exists, use LookupErr.
func (d Doc) Lookup(key ...string) Val {
	val, _ := d.LookupErr(key...)
	return val
}

// LookupErr searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
func (d Doc) LookupErr(key ...string) (Val, error) {
	elem, err := d.LookupElementErr(key...)
	return elem.Value, err
}

// LookupElement searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Element if they key does not exist. To know if they key actually
// exists, use LookupElementErr.
func (d Doc) LookupElement(key ...string) Elem {
	elem, _ := d.LookupElementErr(key...)
	return elem
}

// LookupElementErr searches the document and potentially subdocuments for the
// provided key. Each key provided to this method represents a layer of depth.
func (d Doc) LookupElementErr(key ...string) (Elem, error) {
	// KeyNotFound operates by being created where the error happens and then the depth is
	// incremented by 1 as each function unwinds. Whenever this function returns, it also assigns
	// the Key slice to the key slice it has. This ensures that the proper depth is identified and
	// the proper keys.
	if len(key) == 0 {
		return Elem{}, KeyNotFound{Key: key}
	}

	var elem Elem
	var err error
	idx := d.IndexOf(key[0])
	if idx == -1 {
		return Elem{}, KeyNotFound{Key: key}
	}

	elem = d[idx]
	if len(key) == 1 {
		return elem, nil
	}

	switch elem.Value.Type() {
	case bsontype.EmbeddedDocument:
		switch tt := elem.Value.primitive.(type) {
		case Doc:
			elem, err = tt.LookupElementErr(key[1:]...)
		case MDoc:
			elem, err = tt.LookupElementErr(key[1:]...)
		}
	default:
		return Elem{}, KeyNotFound{Type: elem.Value.Type()}
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
func (d Doc) MarshalBSONValue() (bsontype.Type, []byte, error) {
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
func (d Doc) MarshalBSON() ([]byte, error) { return d.AppendMarshalBSON(nil) }

// AppendMarshalBSON marshals Doc to BSON bytes, appending to dst.
//
// This method will never return an error.
func (d Doc) AppendMarshalBSON(dst []byte) ([]byte, error) {
	idx, dst := bsoncore.ReserveLength(dst)
	for _, elem := range d {
		t, data, _ := elem.Value.MarshalBSONValue() // Value.MarshalBSONValue never returns an error.
		dst = append(dst, byte(t))
		dst = append(dst, elem.Key...)
		dst = append(dst, 0x00)
		dst = append(dst, data...)
	}
	dst = append(dst, 0x00)
	dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	return dst, nil
}

// UnmarshalBSON implements the Unmarshaler interface.
func (d *Doc) UnmarshalBSON(b []byte) error {
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
		*d = d.Append(elem.Key(), val)
	}
	return nil
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (d *Doc) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.EmbeddedDocument {
		return fmt.Errorf("cannot unmarshal %s into a bsonx.Doc", t)
	}
	return d.UnmarshalBSON(data)
}

// Equal compares this document to another, returning true if they are equal.
func (d Doc) Equal(id IDoc) bool {
	switch tt := id.(type) {
	case Doc:
		d2 := tt
		if len(d) != len(d2) {
			return false
		}
		for idx := range d {
			if !d[idx].Equal(d2[idx]) {
				return false
			}
		}
	case MDoc:
		unique := make(map[string]struct{}, 0)
		for _, elem := range d {
			unique[elem.Key] = struct{}{}
			val, ok := tt[elem.Key]
			if !ok {
				return false
			}
			if !val.Equal(elem.Value) {
				return false
			}
		}
		if len(unique) != len(tt) {
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
func (d Doc) String() string {
	var buf bytes.Buffer
	buf.Write([]byte("bson.Document{"))
	for idx, elem := range d {
		if idx > 0 {
			buf.Write([]byte(", "))
		}
		fmt.Fprintf(&buf, "%v", elem)
	}
	buf.WriteByte('}')

	return buf.String()
}

func (Doc) idoc() {}
