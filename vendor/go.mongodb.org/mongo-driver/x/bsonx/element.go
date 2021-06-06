// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const validateMaxDepthDefault = 2048

// ElementTypeError specifies that a method to obtain a BSON value an incorrect type was called on a bson.Value.
//
// TODO: rename this ValueTypeError.
type ElementTypeError struct {
	Method string
	Type   bsontype.Type
}

// Error implements the error interface.
func (ete ElementTypeError) Error() string {
	return "Call of " + ete.Method + " on " + ete.Type.String() + " type"
}

// Elem represents a BSON element.
//
// NOTE: Element cannot be the value of a map nor a property of a struct without special handling.
// The default encoders and decoders will not process Element correctly. To do so would require
// information loss since an Element contains a key, but the keys used when encoding a struct are
// the struct field names. Instead of using an Element, use a Value as a value in a map or a
// property of a struct.
type Elem struct {
	Key   string
	Value Val
}

// Equal compares e and e2 and returns true if they are equal.
func (e Elem) Equal(e2 Elem) bool {
	if e.Key != e2.Key {
		return false
	}
	return e.Value.Equal(e2.Value)
}

func (e Elem) String() string {
	// TODO(GODRIVER-612): When bsoncore has appenders for extended JSON use that here.
	return fmt.Sprintf(`bson.Element{"%s": %v}`, e.Key, e.Value)
}
