// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Val represents a BSON value.
type Val struct {
	// NOTE: The bootstrap is a small amount of space that'll be on the stack. At 15 bytes this
	// doesn't make this type any larger, since there are 7 bytes of padding and we want an int64 to
	// store small values (e.g. boolean, double, int64, etc...). The primitive property is where all
	// of the larger values go. They will use either Go primitives or the primitive.* types.
	t         bsontype.Type
	bootstrap [15]byte
	primitive interface{}
}

func (v Val) reset() Val {
	v.primitive = nil // clear out any pointers so we don't accidentally stop them from being garbage collected.
	v.t = bsontype.Type(0)
	v.bootstrap[0] = 0x00
	v.bootstrap[1] = 0x00
	v.bootstrap[2] = 0x00
	v.bootstrap[3] = 0x00
	v.bootstrap[4] = 0x00
	v.bootstrap[5] = 0x00
	v.bootstrap[6] = 0x00
	v.bootstrap[7] = 0x00
	v.bootstrap[8] = 0x00
	v.bootstrap[9] = 0x00
	v.bootstrap[10] = 0x00
	v.bootstrap[11] = 0x00
	v.bootstrap[12] = 0x00
	v.bootstrap[13] = 0x00
	v.bootstrap[14] = 0x00
	return v
}

func (v Val) string() string {
	if v.primitive != nil {
		return v.primitive.(string)
	}
	// The string will either end with a null byte or it fills the entire bootstrap space.
	length := uint8(v.bootstrap[0])
	return string(v.bootstrap[1 : length+1])
}

func (v Val) writestring(str string) Val {
	switch {
	case len(str) < 15:
		v.bootstrap[0] = uint8(len(str))
		copy(v.bootstrap[1:], str)
	default:
		v.primitive = str
	}
	return v
}

func (v Val) i64() int64 {
	return int64(v.bootstrap[0]) | int64(v.bootstrap[1])<<8 | int64(v.bootstrap[2])<<16 |
		int64(v.bootstrap[3])<<24 | int64(v.bootstrap[4])<<32 | int64(v.bootstrap[5])<<40 |
		int64(v.bootstrap[6])<<48 | int64(v.bootstrap[7])<<56
}

func (v Val) writei64(i64 int64) Val {
	v.bootstrap[0] = byte(i64)
	v.bootstrap[1] = byte(i64 >> 8)
	v.bootstrap[2] = byte(i64 >> 16)
	v.bootstrap[3] = byte(i64 >> 24)
	v.bootstrap[4] = byte(i64 >> 32)
	v.bootstrap[5] = byte(i64 >> 40)
	v.bootstrap[6] = byte(i64 >> 48)
	v.bootstrap[7] = byte(i64 >> 56)
	return v
}

// IsZero returns true if this value is zero or a BSON null.
func (v Val) IsZero() bool { return v.t == bsontype.Type(0) || v.t == bsontype.Null }

func (v Val) String() string {
	// TODO(GODRIVER-612): When bsoncore has appenders for extended JSON use that here.
	return fmt.Sprintf("%v", v.Interface())
}

// Interface returns the Go value of this Value as an empty interface.
//
// This method will return nil if it is empty, otherwise it will return a Go primitive or a
// primitive.* instance.
func (v Val) Interface() interface{} {
	switch v.Type() {
	case bsontype.Double:
		return v.Double()
	case bsontype.String:
		return v.StringValue()
	case bsontype.EmbeddedDocument:
		switch v.primitive.(type) {
		case Doc:
			return v.primitive.(Doc)
		case MDoc:
			return v.primitive.(MDoc)
		default:
			return primitive.Null{}
		}
	case bsontype.Array:
		return v.Array()
	case bsontype.Binary:
		return v.primitive.(primitive.Binary)
	case bsontype.Undefined:
		return primitive.Undefined{}
	case bsontype.ObjectID:
		return v.ObjectID()
	case bsontype.Boolean:
		return v.Boolean()
	case bsontype.DateTime:
		return v.DateTime()
	case bsontype.Null:
		return primitive.Null{}
	case bsontype.Regex:
		return v.primitive.(primitive.Regex)
	case bsontype.DBPointer:
		return v.primitive.(primitive.DBPointer)
	case bsontype.JavaScript:
		return v.JavaScript()
	case bsontype.Symbol:
		return v.Symbol()
	case bsontype.CodeWithScope:
		return v.primitive.(primitive.CodeWithScope)
	case bsontype.Int32:
		return v.Int32()
	case bsontype.Timestamp:
		t, i := v.Timestamp()
		return primitive.Timestamp{T: t, I: i}
	case bsontype.Int64:
		return v.Int64()
	case bsontype.Decimal128:
		return v.Decimal128()
	case bsontype.MinKey:
		return primitive.MinKey{}
	case bsontype.MaxKey:
		return primitive.MaxKey{}
	default:
		return primitive.Null{}
	}
}

// MarshalBSONValue implements the bsoncodec.ValueMarshaler interface.
func (v Val) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return v.MarshalAppendBSONValue(nil)
}

// MarshalAppendBSONValue is similar to MarshalBSONValue, but allows the caller to specify a slice
// to add the bytes to.
func (v Val) MarshalAppendBSONValue(dst []byte) (bsontype.Type, []byte, error) {
	t := v.Type()
	switch v.Type() {
	case bsontype.Double:
		dst = bsoncore.AppendDouble(dst, v.Double())
	case bsontype.String:
		dst = bsoncore.AppendString(dst, v.String())
	case bsontype.EmbeddedDocument:
		switch v.primitive.(type) {
		case Doc:
			t, dst, _ = v.primitive.(Doc).MarshalBSONValue() // Doc.MarshalBSONValue never returns an error.
		case MDoc:
			t, dst, _ = v.primitive.(MDoc).MarshalBSONValue() // MDoc.MarshalBSONValue never returns an error.
		}
	case bsontype.Array:
		t, dst, _ = v.Array().MarshalBSONValue() // Arr.MarshalBSON never returns an error.
	case bsontype.Binary:
		subtype, bindata := v.Binary()
		dst = bsoncore.AppendBinary(dst, subtype, bindata)
	case bsontype.Undefined:
	case bsontype.ObjectID:
		dst = bsoncore.AppendObjectID(dst, v.ObjectID())
	case bsontype.Boolean:
		dst = bsoncore.AppendBoolean(dst, v.Boolean())
	case bsontype.DateTime:
		dst = bsoncore.AppendDateTime(dst, int64(v.DateTime()))
	case bsontype.Null:
	case bsontype.Regex:
		pattern, options := v.Regex()
		dst = bsoncore.AppendRegex(dst, pattern, options)
	case bsontype.DBPointer:
		ns, ptr := v.DBPointer()
		dst = bsoncore.AppendDBPointer(dst, ns, ptr)
	case bsontype.JavaScript:
		dst = bsoncore.AppendJavaScript(dst, string(v.JavaScript()))
	case bsontype.Symbol:
		dst = bsoncore.AppendSymbol(dst, string(v.Symbol()))
	case bsontype.CodeWithScope:
		code, doc := v.CodeWithScope()
		var scope []byte
		scope, _ = doc.MarshalBSON() // Doc.MarshalBSON never returns an error.
		dst = bsoncore.AppendCodeWithScope(dst, code, scope)
	case bsontype.Int32:
		dst = bsoncore.AppendInt32(dst, v.Int32())
	case bsontype.Timestamp:
		t, i := v.Timestamp()
		dst = bsoncore.AppendTimestamp(dst, t, i)
	case bsontype.Int64:
		dst = bsoncore.AppendInt64(dst, v.Int64())
	case bsontype.Decimal128:
		dst = bsoncore.AppendDecimal128(dst, v.Decimal128())
	case bsontype.MinKey:
	case bsontype.MaxKey:
	default:
		panic(fmt.Errorf("invalid BSON type %v", t))
	}

	return t, dst, nil
}

// UnmarshalBSONValue implements the bsoncodec.ValueUnmarshaler interface.
func (v *Val) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if v == nil {
		return errors.New("cannot unmarshal into nil Value")
	}
	var err error
	var ok = true
	var rem []byte
	switch t {
	case bsontype.Double:
		var f64 float64
		f64, rem, ok = bsoncore.ReadDouble(data)
		*v = Double(f64)
	case bsontype.String:
		var str string
		str, rem, ok = bsoncore.ReadString(data)
		*v = String(str)
	case bsontype.EmbeddedDocument:
		var raw []byte
		var doc Doc
		raw, rem, ok = bsoncore.ReadDocument(data)
		doc, err = ReadDoc(raw)
		*v = Document(doc)
	case bsontype.Array:
		var raw []byte
		arr := make(Arr, 0)
		raw, rem, ok = bsoncore.ReadArray(data)
		err = arr.UnmarshalBSONValue(t, raw)
		*v = Array(arr)
	case bsontype.Binary:
		var subtype byte
		var bindata []byte
		subtype, bindata, rem, ok = bsoncore.ReadBinary(data)
		*v = Binary(subtype, bindata)
	case bsontype.Undefined:
		*v = Undefined()
	case bsontype.ObjectID:
		var oid primitive.ObjectID
		oid, rem, ok = bsoncore.ReadObjectID(data)
		*v = ObjectID(oid)
	case bsontype.Boolean:
		var b bool
		b, rem, ok = bsoncore.ReadBoolean(data)
		*v = Boolean(b)
	case bsontype.DateTime:
		var dt int64
		dt, rem, ok = bsoncore.ReadDateTime(data)
		*v = DateTime(dt)
	case bsontype.Null:
		*v = Null()
	case bsontype.Regex:
		var pattern, options string
		pattern, options, rem, ok = bsoncore.ReadRegex(data)
		*v = Regex(pattern, options)
	case bsontype.DBPointer:
		var ns string
		var ptr primitive.ObjectID
		ns, ptr, rem, ok = bsoncore.ReadDBPointer(data)
		*v = DBPointer(ns, ptr)
	case bsontype.JavaScript:
		var js string
		js, rem, ok = bsoncore.ReadJavaScript(data)
		*v = JavaScript(js)
	case bsontype.Symbol:
		var symbol string
		symbol, rem, ok = bsoncore.ReadSymbol(data)
		*v = Symbol(symbol)
	case bsontype.CodeWithScope:
		var raw []byte
		var code string
		var scope Doc
		code, raw, rem, ok = bsoncore.ReadCodeWithScope(data)
		scope, err = ReadDoc(raw)
		*v = CodeWithScope(code, scope)
	case bsontype.Int32:
		var i32 int32
		i32, rem, ok = bsoncore.ReadInt32(data)
		*v = Int32(i32)
	case bsontype.Timestamp:
		var i, t uint32
		t, i, rem, ok = bsoncore.ReadTimestamp(data)
		*v = Timestamp(t, i)
	case bsontype.Int64:
		var i64 int64
		i64, rem, ok = bsoncore.ReadInt64(data)
		*v = Int64(i64)
	case bsontype.Decimal128:
		var d128 primitive.Decimal128
		d128, rem, ok = bsoncore.ReadDecimal128(data)
		*v = Decimal128(d128)
	case bsontype.MinKey:
		*v = MinKey()
	case bsontype.MaxKey:
		*v = MaxKey()
	default:
		err = fmt.Errorf("invalid BSON type %v", t)
	}

	if !ok && err == nil {
		err = bsoncore.NewInsufficientBytesError(data, rem)
	}

	return err
}

// Type returns the BSON type of this value.
func (v Val) Type() bsontype.Type {
	if v.t == bsontype.Type(0) {
		return bsontype.Null
	}
	return v.t
}

// IsNumber returns true if the type of v is a numberic BSON type.
func (v Val) IsNumber() bool {
	switch v.Type() {
	case bsontype.Double, bsontype.Int32, bsontype.Int64, bsontype.Decimal128:
		return true
	default:
		return false
	}
}

// Double returns the BSON double value the Value represents. It panics if the value is a BSON type
// other than double.
func (v Val) Double() float64 {
	if v.t != bsontype.Double {
		panic(ElementTypeError{"bson.Value.Double", v.t})
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(v.bootstrap[0:8]))
}

// DoubleOK is the same as Double, but returns a boolean instead of panicking.
func (v Val) DoubleOK() (float64, bool) {
	if v.t != bsontype.Double {
		return 0, false
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(v.bootstrap[0:8])), true
}

// StringValue returns the BSON string the Value represents. It panics if the value is a BSON type
// other than string.
//
// NOTE: This method is called StringValue to avoid it implementing the
// fmt.Stringer interface.
func (v Val) StringValue() string {
	if v.t != bsontype.String {
		panic(ElementTypeError{"bson.Value.StringValue", v.t})
	}
	return v.string()
}

// StringValueOK is the same as StringValue, but returns a boolean instead of
// panicking.
func (v Val) StringValueOK() (string, bool) {
	if v.t != bsontype.String {
		return "", false
	}
	return v.string(), true
}

func (v Val) asDoc() Doc {
	doc, ok := v.primitive.(Doc)
	if ok {
		return doc
	}
	mdoc := v.primitive.(MDoc)
	for k, v := range mdoc {
		doc = append(doc, Elem{k, v})
	}
	return doc
}

func (v Val) asMDoc() MDoc {
	mdoc, ok := v.primitive.(MDoc)
	if ok {
		return mdoc
	}
	mdoc = make(MDoc)
	doc := v.primitive.(Doc)
	for _, elem := range doc {
		mdoc[elem.Key] = elem.Value
	}
	return mdoc
}

// Document returns the BSON embedded document value the Value represents. It panics if the value
// is a BSON type other than embedded document.
func (v Val) Document() Doc {
	if v.t != bsontype.EmbeddedDocument {
		panic(ElementTypeError{"bson.Value.Document", v.t})
	}
	return v.asDoc()
}

// DocumentOK is the same as Document, except it returns a boolean
// instead of panicking.
func (v Val) DocumentOK() (Doc, bool) {
	if v.t != bsontype.EmbeddedDocument {
		return nil, false
	}
	return v.asDoc(), true
}

// MDocument returns the BSON embedded document value the Value represents. It panics if the value
// is a BSON type other than embedded document.
func (v Val) MDocument() MDoc {
	if v.t != bsontype.EmbeddedDocument {
		panic(ElementTypeError{"bson.Value.MDocument", v.t})
	}
	return v.asMDoc()
}

// MDocumentOK is the same as Document, except it returns a boolean
// instead of panicking.
func (v Val) MDocumentOK() (MDoc, bool) {
	if v.t != bsontype.EmbeddedDocument {
		return nil, false
	}
	return v.asMDoc(), true
}

// Array returns the BSON array value the Value represents. It panics if the value is a BSON type
// other than array.
func (v Val) Array() Arr {
	if v.t != bsontype.Array {
		panic(ElementTypeError{"bson.Value.Array", v.t})
	}
	return v.primitive.(Arr)
}

// ArrayOK is the same as Array, except it returns a boolean
// instead of panicking.
func (v Val) ArrayOK() (Arr, bool) {
	if v.t != bsontype.Array {
		return nil, false
	}
	return v.primitive.(Arr), true
}

// Binary returns the BSON binary value the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) Binary() (byte, []byte) {
	if v.t != bsontype.Binary {
		panic(ElementTypeError{"bson.Value.Binary", v.t})
	}
	bin := v.primitive.(primitive.Binary)
	return bin.Subtype, bin.Data
}

// BinaryOK is the same as Binary, except it returns a boolean instead of
// panicking.
func (v Val) BinaryOK() (byte, []byte, bool) {
	if v.t != bsontype.Binary {
		return 0x00, nil, false
	}
	bin := v.primitive.(primitive.Binary)
	return bin.Subtype, bin.Data, true
}

// Undefined returns the BSON undefined the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) Undefined() {
	if v.t != bsontype.Undefined {
		panic(ElementTypeError{"bson.Value.Undefined", v.t})
	}
	return
}

// UndefinedOK is the same as Undefined, except it returns a boolean instead of
// panicking.
func (v Val) UndefinedOK() bool {
	if v.t != bsontype.Undefined {
		return false
	}
	return true
}

// ObjectID returns the BSON ObjectID the Value represents. It panics if the value is a BSON type
// other than ObjectID.
func (v Val) ObjectID() primitive.ObjectID {
	if v.t != bsontype.ObjectID {
		panic(ElementTypeError{"bson.Value.ObjectID", v.t})
	}
	var oid primitive.ObjectID
	copy(oid[:], v.bootstrap[:12])
	return oid
}

// ObjectIDOK is the same as ObjectID, except it returns a boolean instead of
// panicking.
func (v Val) ObjectIDOK() (primitive.ObjectID, bool) {
	if v.t != bsontype.ObjectID {
		return primitive.ObjectID{}, false
	}
	var oid primitive.ObjectID
	copy(oid[:], v.bootstrap[:12])
	return oid, true
}

// Boolean returns the BSON boolean the Value represents. It panics if the value is a BSON type
// other than boolean.
func (v Val) Boolean() bool {
	if v.t != bsontype.Boolean {
		panic(ElementTypeError{"bson.Value.Boolean", v.t})
	}
	return v.bootstrap[0] == 0x01
}

// BooleanOK is the same as Boolean, except it returns a boolean instead of
// panicking.
func (v Val) BooleanOK() (bool, bool) {
	if v.t != bsontype.Boolean {
		return false, false
	}
	return v.bootstrap[0] == 0x01, true
}

// DateTime returns the BSON datetime the Value represents. It panics if the value is a BSON type
// other than datetime.
func (v Val) DateTime() int64 {
	if v.t != bsontype.DateTime {
		panic(ElementTypeError{"bson.Value.DateTime", v.t})
	}
	return v.i64()
}

// DateTimeOK is the same as DateTime, except it returns a boolean instead of
// panicking.
func (v Val) DateTimeOK() (int64, bool) {
	if v.t != bsontype.DateTime {
		return 0, false
	}
	return v.i64(), true
}

// Time returns the BSON datetime the Value represents as time.Time. It panics if the value is a BSON
// type other than datetime.
func (v Val) Time() time.Time {
	if v.t != bsontype.DateTime {
		panic(ElementTypeError{"bson.Value.Time", v.t})
	}
	i := v.i64()
	return time.Unix(int64(i)/1000, int64(i)%1000*1000000)
}

// TimeOK is the same as Time, except it returns a boolean instead of
// panicking.
func (v Val) TimeOK() (time.Time, bool) {
	if v.t != bsontype.DateTime {
		return time.Time{}, false
	}
	i := v.i64()
	return time.Unix(int64(i)/1000, int64(i)%1000*1000000), true
}

// Null returns the BSON undefined the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) Null() {
	if v.t != bsontype.Null && v.t != bsontype.Type(0) {
		panic(ElementTypeError{"bson.Value.Null", v.t})
	}
	return
}

// NullOK is the same as Null, except it returns a boolean instead of
// panicking.
func (v Val) NullOK() bool {
	if v.t != bsontype.Null && v.t != bsontype.Type(0) {
		return false
	}
	return true
}

// Regex returns the BSON regex the Value represents. It panics if the value is a BSON type
// other than regex.
func (v Val) Regex() (pattern, options string) {
	if v.t != bsontype.Regex {
		panic(ElementTypeError{"bson.Value.Regex", v.t})
	}
	regex := v.primitive.(primitive.Regex)
	return regex.Pattern, regex.Options
}

// RegexOK is the same as Regex, except that it returns a boolean
// instead of panicking.
func (v Val) RegexOK() (pattern, options string, ok bool) {
	if v.t != bsontype.Regex {
		return "", "", false
	}
	regex := v.primitive.(primitive.Regex)
	return regex.Pattern, regex.Options, true
}

// DBPointer returns the BSON dbpointer the Value represents. It panics if the value is a BSON type
// other than dbpointer.
func (v Val) DBPointer() (string, primitive.ObjectID) {
	if v.t != bsontype.DBPointer {
		panic(ElementTypeError{"bson.Value.DBPointer", v.t})
	}
	dbptr := v.primitive.(primitive.DBPointer)
	return dbptr.DB, dbptr.Pointer
}

// DBPointerOK is the same as DBPoitner, except that it returns a boolean
// instead of panicking.
func (v Val) DBPointerOK() (string, primitive.ObjectID, bool) {
	if v.t != bsontype.DBPointer {
		return "", primitive.ObjectID{}, false
	}
	dbptr := v.primitive.(primitive.DBPointer)
	return dbptr.DB, dbptr.Pointer, true
}

// JavaScript returns the BSON JavaScript the Value represents. It panics if the value is a BSON type
// other than JavaScript.
func (v Val) JavaScript() string {
	if v.t != bsontype.JavaScript {
		panic(ElementTypeError{"bson.Value.JavaScript", v.t})
	}
	return v.string()
}

// JavaScriptOK is the same as Javascript, except that it returns a boolean
// instead of panicking.
func (v Val) JavaScriptOK() (string, bool) {
	if v.t != bsontype.JavaScript {
		return "", false
	}
	return v.string(), true
}

// Symbol returns the BSON symbol the Value represents. It panics if the value is a BSON type
// other than symbol.
func (v Val) Symbol() string {
	if v.t != bsontype.Symbol {
		panic(ElementTypeError{"bson.Value.Symbol", v.t})
	}
	return v.string()
}

// SymbolOK is the same as Javascript, except that it returns a boolean
// instead of panicking.
func (v Val) SymbolOK() (string, bool) {
	if v.t != bsontype.Symbol {
		return "", false
	}
	return v.string(), true
}

// CodeWithScope returns the BSON code with scope value the Value represents. It panics if the
// value is a BSON type other than code with scope.
func (v Val) CodeWithScope() (string, Doc) {
	if v.t != bsontype.CodeWithScope {
		panic(ElementTypeError{"bson.Value.CodeWithScope", v.t})
	}
	cws := v.primitive.(primitive.CodeWithScope)
	return string(cws.Code), cws.Scope.(Doc)
}

// CodeWithScopeOK is the same as JavascriptWithScope,
// except that it returns a boolean instead of panicking.
func (v Val) CodeWithScopeOK() (string, Doc, bool) {
	if v.t != bsontype.CodeWithScope {
		return "", nil, false
	}
	cws := v.primitive.(primitive.CodeWithScope)
	return string(cws.Code), cws.Scope.(Doc), true
}

// Int32 returns the BSON int32 the Value represents. It panics if the value is a BSON type
// other than int32.
func (v Val) Int32() int32 {
	if v.t != bsontype.Int32 {
		panic(ElementTypeError{"bson.Value.Int32", v.t})
	}
	return int32(v.bootstrap[0]) | int32(v.bootstrap[1])<<8 |
		int32(v.bootstrap[2])<<16 | int32(v.bootstrap[3])<<24
}

// Int32OK is the same as Int32, except that it returns a boolean instead of
// panicking.
func (v Val) Int32OK() (int32, bool) {
	if v.t != bsontype.Int32 {
		return 0, false
	}
	return int32(v.bootstrap[0]) | int32(v.bootstrap[1])<<8 |
			int32(v.bootstrap[2])<<16 | int32(v.bootstrap[3])<<24,
		true
}

// Timestamp returns the BSON timestamp the Value represents. It panics if the value is a
// BSON type other than timestamp.
func (v Val) Timestamp() (t, i uint32) {
	if v.t != bsontype.Timestamp {
		panic(ElementTypeError{"bson.Value.Timestamp", v.t})
	}
	return uint32(v.bootstrap[4]) | uint32(v.bootstrap[5])<<8 |
			uint32(v.bootstrap[6])<<16 | uint32(v.bootstrap[7])<<24,
		uint32(v.bootstrap[0]) | uint32(v.bootstrap[1])<<8 |
			uint32(v.bootstrap[2])<<16 | uint32(v.bootstrap[3])<<24
}

// TimestampOK is the same as Timestamp, except that it returns a boolean
// instead of panicking.
func (v Val) TimestampOK() (t uint32, i uint32, ok bool) {
	if v.t != bsontype.Timestamp {
		return 0, 0, false
	}
	return uint32(v.bootstrap[4]) | uint32(v.bootstrap[5])<<8 |
			uint32(v.bootstrap[6])<<16 | uint32(v.bootstrap[7])<<24,
		uint32(v.bootstrap[0]) | uint32(v.bootstrap[1])<<8 |
			uint32(v.bootstrap[2])<<16 | uint32(v.bootstrap[3])<<24,
		true
}

// Int64 returns the BSON int64 the Value represents. It panics if the value is a BSON type
// other than int64.
func (v Val) Int64() int64 {
	if v.t != bsontype.Int64 {
		panic(ElementTypeError{"bson.Value.Int64", v.t})
	}
	return v.i64()
}

// Int64OK is the same as Int64, except that it returns a boolean instead of
// panicking.
func (v Val) Int64OK() (int64, bool) {
	if v.t != bsontype.Int64 {
		return 0, false
	}
	return v.i64(), true
}

// Decimal128 returns the BSON decimal128 value the Value represents. It panics if the value is a
// BSON type other than decimal128.
func (v Val) Decimal128() primitive.Decimal128 {
	if v.t != bsontype.Decimal128 {
		panic(ElementTypeError{"bson.Value.Decimal128", v.t})
	}
	return v.primitive.(primitive.Decimal128)
}

// Decimal128OK is the same as Decimal128, except that it returns a boolean
// instead of panicking.
func (v Val) Decimal128OK() (primitive.Decimal128, bool) {
	if v.t != bsontype.Decimal128 {
		return primitive.Decimal128{}, false
	}
	return v.primitive.(primitive.Decimal128), true
}

// MinKey returns the BSON minkey the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) MinKey() {
	if v.t != bsontype.MinKey {
		panic(ElementTypeError{"bson.Value.MinKey", v.t})
	}
	return
}

// MinKeyOK is the same as MinKey, except it returns a boolean instead of
// panicking.
func (v Val) MinKeyOK() bool {
	if v.t != bsontype.MinKey {
		return false
	}
	return true
}

// MaxKey returns the BSON maxkey the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) MaxKey() {
	if v.t != bsontype.MaxKey {
		panic(ElementTypeError{"bson.Value.MaxKey", v.t})
	}
	return
}

// MaxKeyOK is the same as MaxKey, except it returns a boolean instead of
// panicking.
func (v Val) MaxKeyOK() bool {
	if v.t != bsontype.MaxKey {
		return false
	}
	return true
}

// Equal compares v to v2 and returns true if they are equal. Unknown BSON types are
// never equal. Two empty values are equal.
func (v Val) Equal(v2 Val) bool {
	if v.Type() != v2.Type() {
		return false
	}
	if v.IsZero() && v2.IsZero() {
		return true
	}

	switch v.Type() {
	case bsontype.Double, bsontype.DateTime, bsontype.Timestamp, bsontype.Int64:
		return bytes.Equal(v.bootstrap[0:8], v2.bootstrap[0:8])
	case bsontype.String:
		return v.string() == v2.string()
	case bsontype.EmbeddedDocument:
		return v.equalDocs(v2)
	case bsontype.Array:
		return v.Array().Equal(v2.Array())
	case bsontype.Binary:
		return v.primitive.(primitive.Binary).Equal(v2.primitive.(primitive.Binary))
	case bsontype.Undefined:
		return true
	case bsontype.ObjectID:
		return bytes.Equal(v.bootstrap[0:12], v2.bootstrap[0:12])
	case bsontype.Boolean:
		return v.bootstrap[0] == v2.bootstrap[0]
	case bsontype.Null:
		return true
	case bsontype.Regex:
		return v.primitive.(primitive.Regex).Equal(v2.primitive.(primitive.Regex))
	case bsontype.DBPointer:
		return v.primitive.(primitive.DBPointer).Equal(v2.primitive.(primitive.DBPointer))
	case bsontype.JavaScript:
		return v.JavaScript() == v2.JavaScript()
	case bsontype.Symbol:
		return v.Symbol() == v2.Symbol()
	case bsontype.CodeWithScope:
		code1, scope1 := v.primitive.(primitive.CodeWithScope).Code, v.primitive.(primitive.CodeWithScope).Scope
		code2, scope2 := v2.primitive.(primitive.CodeWithScope).Code, v2.primitive.(primitive.CodeWithScope).Scope
		return code1 == code2 && v.equalInterfaceDocs(scope1, scope2)
	case bsontype.Int32:
		return v.Int32() == v2.Int32()
	case bsontype.Decimal128:
		h, l := v.Decimal128().GetBytes()
		h2, l2 := v2.Decimal128().GetBytes()
		return h == h2 && l == l2
	case bsontype.MinKey:
		return true
	case bsontype.MaxKey:
		return true
	default:
		return false
	}
}

func (v Val) equalDocs(v2 Val) bool {
	_, ok1 := v.primitive.(MDoc)
	_, ok2 := v2.primitive.(MDoc)
	if ok1 || ok2 {
		return v.asMDoc().Equal(v2.asMDoc())
	}
	return v.asDoc().Equal(v2.asDoc())
}

func (Val) equalInterfaceDocs(i, i2 interface{}) bool {
	switch d := i.(type) {
	case MDoc:
		d2, ok := i2.(IDoc)
		if !ok {
			return false
		}
		return d.Equal(d2)
	case Doc:
		d2, ok := i2.(IDoc)
		if !ok {
			return false
		}
		return d.Equal(d2)
	case nil:
		return i2 == nil
	default:
		return false
	}
}
