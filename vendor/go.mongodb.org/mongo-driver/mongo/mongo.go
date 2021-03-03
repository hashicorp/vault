// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo // import "go.mongodb.org/mongo-driver/mongo"

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Dialer is used to make network connections.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// BSONAppender is an interface implemented by types that can marshal a
// provided type into BSON bytes and append those bytes to the provided []byte.
// The AppendBSON can return a non-nil error and non-nil []byte. The AppendBSON
// method may also write incomplete BSON to the []byte.
type BSONAppender interface {
	AppendBSON([]byte, interface{}) ([]byte, error)
}

// BSONAppenderFunc is an adapter function that allows any function that
// satisfies the AppendBSON method signature to be used where a BSONAppender is
// used.
type BSONAppenderFunc func([]byte, interface{}) ([]byte, error)

// AppendBSON implements the BSONAppender interface
func (baf BSONAppenderFunc) AppendBSON(dst []byte, val interface{}) ([]byte, error) {
	return baf(dst, val)
}

// MarshalError is returned when attempting to transform a value into a document
// results in an error.
type MarshalError struct {
	Value interface{}
	Err   error
}

// Error implements the error interface.
func (me MarshalError) Error() string {
	return fmt.Sprintf("cannot transform type %s to a BSON Document: %v", reflect.TypeOf(me.Value), me.Err)
}

// Pipeline is a type that makes creating aggregation pipelines easier. It is a
// helper and is intended for serializing to BSON.
//
// Example usage:
//
//		mongo.Pipeline{
//			{{"$group", bson.D{{"_id", "$state"}, {"totalPop", bson.D{{"$sum", "$pop"}}}}}},
//			{{"$match", bson.D{{"totalPop", bson.D{{"$gte", 10*1000*1000}}}}}},
//		}
//
type Pipeline []bson.D

// transformAndEnsureID is a hack that makes it easy to get a RawValue as the _id value. This will
// be removed when we switch from using bsonx to bsoncore for the driver package.
func transformAndEnsureID(registry *bsoncodec.Registry, val interface{}) (bsonx.Doc, interface{}, error) {
	// TODO: performance is going to be pretty bad for bsonx.Doc here since we turn it into a []byte
	// only to turn it back into a bsonx.Doc. We can fix this post beta1 when we refactor the driver
	// package to use bsoncore.Document instead of bsonx.Doc.
	if registry == nil {
		registry = bson.NewRegistryBuilder().Build()
	}
	switch tt := val.(type) {
	case nil:
		return nil, nil, ErrNilDocument
	case bsonx.Doc:
		val = tt.Copy()
	case []byte:
		// Slight optimization so we'll just use MarshalBSON and not go through the codec machinery.
		val = bson.Raw(tt)
	}

	// TODO(skriptble): Use a pool of these instead.
	buf := make([]byte, 0, 256)
	b, err := bson.MarshalAppendWithRegistry(registry, buf, val)
	if err != nil {
		return nil, nil, MarshalError{Value: val, Err: err}
	}

	d, err := bsonx.ReadDoc(b)
	if err != nil {
		return nil, nil, err
	}

	var id interface{}

	idx := d.IndexOf("_id")
	var idElem bsonx.Elem
	switch idx {
	case -1:
		idElem = bsonx.Elem{"_id", bsonx.ObjectID(primitive.NewObjectID())}
		d = append(d, bsonx.Elem{})
		copy(d[1:], d)
		d[0] = idElem
	default:
		idElem = d[idx]
		copy(d[1:idx+1], d[0:idx])
		d[0] = idElem
	}

	idBuf := make([]byte, 0, 256)
	t, data, err := idElem.Value.MarshalAppendBSONValue(idBuf[:0])
	if err != nil {
		return nil, nil, err
	}

	err = bson.RawValue{Type: t, Value: data}.UnmarshalWithRegistry(registry, &id)
	if err != nil {
		return nil, nil, err
	}

	return d, id, nil
}

// transformAndEnsureIDv2 is a hack that makes it easy to get a RawValue as the _id value. This will
// be removed when we switch from using bsonx to bsoncore for the driver package.
func transformAndEnsureIDv2(registry *bsoncodec.Registry, val interface{}) (bsoncore.Document, interface{}, error) {
	if registry == nil {
		registry = bson.NewRegistryBuilder().Build()
	}
	switch tt := val.(type) {
	case nil:
		return nil, nil, ErrNilDocument
	case bsonx.Doc:
		val = tt.Copy()
	case []byte:
		// Slight optimization so we'll just use MarshalBSON and not go through the codec machinery.
		val = bson.Raw(tt)
	}

	// TODO(skriptble): Use a pool of these instead.
	doc := make(bsoncore.Document, 0, 256)
	doc, err := bson.MarshalAppendWithRegistry(registry, doc, val)
	if err != nil {
		return nil, nil, MarshalError{Value: val, Err: err}
	}

	var id interface{}

	value := doc.Lookup("_id")
	switch value.Type {
	case bsontype.Type(0):
		value = bsoncore.Value{Type: bsontype.ObjectID, Data: bsoncore.AppendObjectID(nil, primitive.NewObjectID())}
		olddoc := doc
		doc = make(bsoncore.Document, 0, len(olddoc)+17) // type byte + _id + null byte + object ID
		_, doc = bsoncore.ReserveLength(doc)
		doc = bsoncore.AppendValueElement(doc, "_id", value)
		doc = append(doc, olddoc[4:]...) // remove the length
		doc = bsoncore.UpdateLength(doc, 0, int32(len(doc)))
	default:
		// We copy the bytes here to ensure that any bytes returned to the user aren't modified
		// later.
		buf := make([]byte, len(value.Data))
		copy(buf, value.Data)
		value.Data = buf
	}

	err = bson.RawValue{Type: value.Type, Value: value.Data}.UnmarshalWithRegistry(registry, &id)
	if err != nil {
		return nil, nil, err
	}

	return doc, id, nil
}

func transformDocument(registry *bsoncodec.Registry, val interface{}) (bsonx.Doc, error) {
	if doc, ok := val.(bsonx.Doc); ok {
		return doc.Copy(), nil
	}
	b, err := transformBsoncoreDocument(registry, val)
	if err != nil {
		return nil, err
	}
	return bsonx.ReadDoc(b)
}

func transformBsoncoreDocument(registry *bsoncodec.Registry, val interface{}) (bsoncore.Document, error) {
	if registry == nil {
		registry = bson.DefaultRegistry
	}
	if val == nil {
		return nil, ErrNilDocument
	}
	if bs, ok := val.([]byte); ok {
		// Slight optimization so we'll just use MarshalBSON and not go through the codec machinery.
		val = bson.Raw(bs)
	}

	// TODO(skriptble): Use a pool of these instead.
	buf := make([]byte, 0, 256)
	b, err := bson.MarshalAppendWithRegistry(registry, buf[:0], val)
	if err != nil {
		return nil, MarshalError{Value: val, Err: err}
	}
	return b, nil
}

func ensureID(d bsonx.Doc) (bsonx.Doc, interface{}) {
	var id interface{}

	elem, err := d.LookupElementErr("_id")
	switch err.(type) {
	case nil:
		id = elem
	default:
		oid := primitive.NewObjectID()
		d = append(d, bsonx.Elem{"_id", bsonx.ObjectID(oid)})
		id = oid
	}
	return d, id
}

func ensureDollarKey(doc bsonx.Doc) error {
	if len(doc) == 0 {
		return errors.New("update document must have at least one element")
	}
	if !strings.HasPrefix(doc[0].Key, "$") {
		return errors.New("update document must contain key beginning with '$'")
	}
	return nil
}

func ensureDollarKeyv2(doc bsoncore.Document) error {
	firstElem, err := doc.IndexErr(0)
	if err != nil {
		return errors.New("update document must have at least one element")
	}

	if !strings.HasPrefix(firstElem.Key(), "$") {
		return errors.New("update document must contain key beginning with '$'")
	}
	return nil
}

func ensureNoDollarKey(doc bsoncore.Document) error {
	if elem, err := doc.IndexErr(0); err == nil && strings.HasPrefix(elem.Key(), "$") {
		return errors.New("replacement document cannot contains keys beginning with '$")
	}

	return nil
}

func transformAggregatePipeline(registry *bsoncodec.Registry, pipeline interface{}) (bsonx.Arr, error) {
	pipelineArr := bsonx.Arr{}
	switch t := pipeline.(type) {
	case bsoncodec.ValueMarshaler:
		btype, val, err := t.MarshalBSONValue()
		if err != nil {
			return nil, err
		}
		if btype != bsontype.Array {
			return nil, fmt.Errorf("ValueMarshaler returned a %v, but was expecting %v", btype, bsontype.Array)
		}
		err = pipelineArr.UnmarshalBSONValue(btype, val)
		if err != nil {
			return nil, err
		}
	default:
		val := reflect.ValueOf(t)
		if !val.IsValid() || (val.Kind() != reflect.Slice && val.Kind() != reflect.Array) {
			return nil, fmt.Errorf("can only transform slices and arrays into aggregation pipelines, but got %v", val.Kind())
		}
		for idx := 0; idx < val.Len(); idx++ {
			elem, err := transformDocument(registry, val.Index(idx).Interface())
			if err != nil {
				return nil, err
			}
			pipelineArr = append(pipelineArr, bsonx.Document(elem))
		}
	}

	return pipelineArr, nil
}

func transformAggregatePipelinev2(registry *bsoncodec.Registry, pipeline interface{}) (bsoncore.Document, bool, error) {
	switch t := pipeline.(type) {
	case bsoncodec.ValueMarshaler:
		btype, val, err := t.MarshalBSONValue()
		if err != nil {
			return nil, false, err
		}
		if btype != bsontype.Array {
			return nil, false, fmt.Errorf("ValueMarshaler returned a %v, but was expecting %v", btype, bsontype.Array)
		}

		var hasOutputStage bool
		pipelineDoc := bsoncore.Document(val)
		values, _ := pipelineDoc.Values()
		if pipelineLen := len(values); pipelineLen > 0 {
			if finalDoc, ok := values[pipelineLen-1].DocumentOK(); ok {
				if elem, err := finalDoc.IndexErr(0); err == nil && (elem.Key() == "$out" || elem.Key() == "$merge") {
					hasOutputStage = true
				}
			}
		}

		return pipelineDoc, hasOutputStage, nil
	default:
		val := reflect.ValueOf(t)
		if !val.IsValid() || (val.Kind() != reflect.Slice && val.Kind() != reflect.Array) {
			return nil, false, fmt.Errorf("can only transform slices and arrays into aggregation pipelines, but got %v", val.Kind())
		}

		aidx, arr := bsoncore.AppendArrayStart(nil)
		var hasOutputStage bool
		valLen := val.Len()
		for idx := 0; idx < valLen; idx++ {
			doc, err := transformBsoncoreDocument(registry, val.Index(idx).Interface())
			if err != nil {
				return nil, false, err
			}

			if idx == valLen-1 {
				if elem, err := doc.IndexErr(0); err == nil && (elem.Key() == "$out" || elem.Key() == "$merge") {
					hasOutputStage = true
				}
			}
			arr = bsoncore.AppendDocumentElement(arr, strconv.Itoa(idx), doc)
		}
		arr, _ = bsoncore.AppendArrayEnd(arr, aidx)
		return arr, hasOutputStage, nil
	}
}

func transformUpdateValue(registry *bsoncodec.Registry, update interface{}, dollarKeysAllowed bool) (bsoncore.Value, error) {
	documentCheckerFunc := ensureDollarKeyv2
	if !dollarKeysAllowed {
		documentCheckerFunc = ensureNoDollarKey
	}

	var u bsoncore.Value
	var err error
	switch t := update.(type) {
	case nil:
		return u, ErrNilDocument
	case primitive.D, bsonx.Doc:
		u.Type = bsontype.EmbeddedDocument
		u.Data, err = transformBsoncoreDocument(registry, update)
		if err != nil {
			return u, err
		}

		return u, documentCheckerFunc(u.Data)
	case bson.Raw:
		u.Type = bsontype.EmbeddedDocument
		u.Data = t
		return u, documentCheckerFunc(u.Data)
	case bsoncore.Document:
		u.Type = bsontype.EmbeddedDocument
		u.Data = t
		return u, documentCheckerFunc(u.Data)
	case []byte:
		u.Type = bsontype.EmbeddedDocument
		u.Data = t
		return u, documentCheckerFunc(u.Data)
	case bsoncodec.Marshaler:
		u.Type = bsontype.EmbeddedDocument
		u.Data, err = t.MarshalBSON()
		if err != nil {
			return u, err
		}

		return u, documentCheckerFunc(u.Data)
	case bsoncodec.ValueMarshaler:
		u.Type, u.Data, err = t.MarshalBSONValue()
		if err != nil {
			return u, err
		}
		if u.Type != bsontype.Array && u.Type != bsontype.EmbeddedDocument {
			return u, fmt.Errorf("ValueMarshaler returned a %v, but was expecting %v or %v", u.Type, bsontype.Array, bsontype.EmbeddedDocument)
		}
		return u, err
	default:
		val := reflect.ValueOf(t)
		if !val.IsValid() {
			return u, fmt.Errorf("can only transform slices and arrays into update pipelines, but got %v", val.Kind())
		}
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			u.Type = bsontype.EmbeddedDocument
			u.Data, err = transformBsoncoreDocument(registry, update)
			if err != nil {
				return u, err
			}

			return u, documentCheckerFunc(u.Data)
		}

		u.Type = bsontype.Array
		aidx, arr := bsoncore.AppendArrayStart(nil)
		valLen := val.Len()
		for idx := 0; idx < valLen; idx++ {
			doc, err := transformBsoncoreDocument(registry, val.Index(idx).Interface())
			if err != nil {
				return u, err
			}

			if err := documentCheckerFunc(doc); err != nil {
				return u, err
			}

			arr = bsoncore.AppendDocumentElement(arr, strconv.Itoa(idx), doc)
		}
		u.Data, _ = bsoncore.AppendArrayEnd(arr, aidx)
		return u, err
	}
}

func transformValue(registry *bsoncodec.Registry, val interface{}) (bsoncore.Value, error) {
	if registry == nil {
		registry = bson.DefaultRegistry
	}
	if val == nil {
		return bsoncore.Value{}, ErrNilValue
	}

	buf := make([]byte, 0, 256)
	bsonType, bsonValue, err := bson.MarshalValueAppendWithRegistry(registry, buf[:0], val)
	if err != nil {
		return bsoncore.Value{}, MarshalError{Value: val, Err: err}
	}

	return bsoncore.Value{Type: bsonType, Data: bsonValue}, nil
}

// Build the aggregation pipeline for the CountDocument command.
func countDocumentsAggregatePipeline(registry *bsoncodec.Registry, filter interface{}, opts *options.CountOptions) (bsoncore.Document, error) {
	filterDoc, err := transformBsoncoreDocument(registry, filter)
	if err != nil {
		return nil, err
	}

	aidx, arr := bsoncore.AppendArrayStart(nil)
	didx, arr := bsoncore.AppendDocumentElementStart(arr, strconv.Itoa(0))
	arr = bsoncore.AppendDocumentElement(arr, "$match", filterDoc)
	arr, _ = bsoncore.AppendDocumentEnd(arr, didx)

	index := 1
	if opts != nil {
		if opts.Skip != nil {
			didx, arr = bsoncore.AppendDocumentElementStart(arr, strconv.Itoa(index))
			arr = bsoncore.AppendInt64Element(arr, "$skip", *opts.Skip)
			arr, _ = bsoncore.AppendDocumentEnd(arr, didx)
			index++
		}
		if opts.Limit != nil {
			didx, arr = bsoncore.AppendDocumentElementStart(arr, strconv.Itoa(index))
			arr = bsoncore.AppendInt64Element(arr, "$limit", *opts.Limit)
			arr, _ = bsoncore.AppendDocumentEnd(arr, didx)
			index++
		}
	}

	didx, arr = bsoncore.AppendDocumentElementStart(arr, strconv.Itoa(index))
	iidx, arr := bsoncore.AppendDocumentElementStart(arr, "$group")
	arr = bsoncore.AppendInt32Element(arr, "_id", 1)
	iiidx, arr := bsoncore.AppendDocumentElementStart(arr, "n")
	arr = bsoncore.AppendInt32Element(arr, "$sum", 1)
	arr, _ = bsoncore.AppendDocumentEnd(arr, iiidx)
	arr, _ = bsoncore.AppendDocumentEnd(arr, iidx)
	arr, _ = bsoncore.AppendDocumentEnd(arr, didx)

	return bsoncore.AppendArrayEnd(arr, aidx)
}
