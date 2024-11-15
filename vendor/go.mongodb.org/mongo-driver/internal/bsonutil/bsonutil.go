// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonutil

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// StringSliceFromRawValue decodes the provided BSON value into a []string. This function returns an error if the value
// is not an array or any of the elements in the array are not strings. The name parameter is used to add context to
// error messages.
func StringSliceFromRawValue(name string, val bson.RawValue) ([]string, error) {
	arr, ok := val.ArrayOK()
	if !ok {
		return nil, fmt.Errorf("expected '%s' to be an array but it's a BSON %s", name, val.Type)
	}

	arrayValues, err := arr.Values()
	if err != nil {
		return nil, err
	}

	strs := make([]string, 0, len(arrayValues))
	for _, arrayVal := range arrayValues {
		str, ok := arrayVal.StringValueOK()
		if !ok {
			return nil, fmt.Errorf("expected '%s' to be an array of strings, but found a BSON %s", name, arrayVal.Type)
		}
		strs = append(strs, str)
	}
	return strs, nil
}

// RawToDocuments converts a bson.Raw that is internally an array of documents to []bson.Raw.
func RawToDocuments(doc bson.Raw) []bson.Raw {
	values, err := doc.Values()
	if err != nil {
		panic(fmt.Sprintf("error converting BSON document to values: %v", err))
	}

	out := make([]bson.Raw, len(values))
	for i := range values {
		out[i] = values[i].Document()
	}

	return out
}

// RawToInterfaces takes one or many bson.Raw documents and returns them as a []interface{}.
func RawToInterfaces(docs ...bson.Raw) []interface{} {
	out := make([]interface{}, len(docs))
	for i := range docs {
		out[i] = docs[i]
	}
	return out
}
