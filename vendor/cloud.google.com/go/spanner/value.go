/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/internal/fields"
	"github.com/golang/protobuf/proto"
	proto3 "github.com/golang/protobuf/ptypes/struct"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
)

const (
	// nullString is returned by the String methods of NullableValues when the
	// underlying database value is null.
	nullString                       = "<null>"
	commitTimestampPlaceholderString = "spanner.commit_timestamp()"

	// NumericPrecisionDigits is the maximum number of digits in a NUMERIC
	// value.
	NumericPrecisionDigits = 38

	// NumericScaleDigits is the maximum number of digits after the decimal
	// point in a NUMERIC value.
	NumericScaleDigits = 9
)

// NumericString returns a string representing a *big.Rat in a format compatible
// with Spanner SQL. It returns a floating-point literal with 9 digits after the
// decimal point.
func NumericString(r *big.Rat) string {
	return r.FloatString(NumericScaleDigits)
}

var (
	// CommitTimestamp is a special value used to tell Cloud Spanner to insert
	// the commit timestamp of the transaction into a column. It can be used in
	// a Mutation, or directly used in InsertStruct or InsertMap. See
	// ExampleCommitTimestamp. This is just a placeholder and the actual value
	// stored in this variable has no meaning.
	CommitTimestamp = commitTimestamp
	commitTimestamp = time.Unix(0, 0).In(time.FixedZone("CommitTimestamp placeholder", 0xDB))

	jsonNullBytes = []byte("null")
)

// Encoder is the interface implemented by a custom type that can be encoded to
// a supported type by Spanner. A code example:
//
// type customField struct {
//     Prefix string
//     Suffix string
// }
//
// // Convert a customField value to a string
// func (cf customField) EncodeSpanner() (interface{}, error) {
//     var b bytes.Buffer
//     b.WriteString(cf.Prefix)
//     b.WriteString("-")
//     b.WriteString(cf.Suffix)
//     return b.String(), nil
// }
type Encoder interface {
	EncodeSpanner() (interface{}, error)
}

// Decoder is the interface implemented by a custom type that can be decoded
// from a supported type by Spanner. A code example:
//
// type customField struct {
//     Prefix string
//     Suffix string
// }
//
// // Convert a string to a customField value
// func (cf *customField) DecodeSpanner(val interface{}) (err error) {
//     strVal, ok := val.(string)
//     if !ok {
//         return fmt.Errorf("failed to decode customField: %v", val)
//     }
//     s := strings.Split(strVal, "-")
//     if len(s) > 1 {
//         cf.Prefix = s[0]
//         cf.Suffix = s[1]
//     }
//     return nil
// }
type Decoder interface {
	DecodeSpanner(input interface{}) error
}

// NullableValue is the interface implemented by all null value wrapper types.
type NullableValue interface {
	// IsNull returns true if the underlying database value is null.
	IsNull() bool
}

// NullInt64 represents a Cloud Spanner INT64 that may be NULL.
type NullInt64 struct {
	Int64 int64 // Int64 contains the value when it is non-NULL, and zero when NULL.
	Valid bool  // Valid is true if Int64 is not NULL.
}

// IsNull implements NullableValue.IsNull for NullInt64.
func (n NullInt64) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullInt64
func (n NullInt64) String() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", n.Int64)
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullInt64.
func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%v", n.Int64)), nil
	}
	return jsonNullBytes, nil
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullInt64.
func (n *NullInt64) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Int64 = int64(0)
		n.Valid = false
		return nil
	}
	num, err := strconv.ParseInt(string(payload), 10, 64)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to int64: got %v", string(payload))
	}
	n.Int64 = num
	n.Valid = true
	return nil
}

// NullString represents a Cloud Spanner STRING that may be NULL.
type NullString struct {
	StringVal string // StringVal contains the value when it is non-NULL, and an empty string when NULL.
	Valid     bool   // Valid is true if StringVal is not NULL.
}

// IsNull implements NullableValue.IsNull for NullString.
func (n NullString) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullString
func (n NullString) String() string {
	if !n.Valid {
		return nullString
	}
	return n.StringVal
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullString.
func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%q", n.StringVal)), nil
	}
	return jsonNullBytes, nil
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullString.
func (n *NullString) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.StringVal = ""
		n.Valid = false
		return nil
	}
	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	n.StringVal = string(payload)
	n.Valid = true
	return nil
}

// NullFloat64 represents a Cloud Spanner FLOAT64 that may be NULL.
type NullFloat64 struct {
	Float64 float64 // Float64 contains the value when it is non-NULL, and zero when NULL.
	Valid   bool    // Valid is true if Float64 is not NULL.
}

// IsNull implements NullableValue.IsNull for NullFloat64.
func (n NullFloat64) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullFloat64
func (n NullFloat64) String() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", n.Float64)
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullFloat64.
func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%v", n.Float64)), nil
	}
	return jsonNullBytes, nil
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullFloat64.
func (n *NullFloat64) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Float64 = float64(0)
		n.Valid = false
		return nil
	}
	num, err := strconv.ParseFloat(string(payload), 64)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to float64: got %v", string(payload))
	}
	n.Float64 = num
	n.Valid = true
	return nil
}

// NullBool represents a Cloud Spanner BOOL that may be NULL.
type NullBool struct {
	Bool  bool // Bool contains the value when it is non-NULL, and false when NULL.
	Valid bool // Valid is true if Bool is not NULL.
}

// IsNull implements NullableValue.IsNull for NullBool.
func (n NullBool) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullBool
func (n NullBool) String() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", n.Bool)
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullBool.
func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%v", n.Bool)), nil
	}
	return jsonNullBytes, nil
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullBool.
func (n *NullBool) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Bool = false
		n.Valid = false
		return nil
	}
	b, err := strconv.ParseBool(string(payload))
	if err != nil {
		return fmt.Errorf("payload cannot be converted to bool: got %v", string(payload))
	}
	n.Bool = b
	n.Valid = true
	return nil
}

// NullTime represents a Cloud Spanner TIMESTAMP that may be null.
type NullTime struct {
	Time  time.Time // Time contains the value when it is non-NULL, and a zero time.Time when NULL.
	Valid bool      // Valid is true if Time is not NULL.
}

// IsNull implements NullableValue.IsNull for NullTime.
func (n NullTime) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullTime
func (n NullTime) String() string {
	if !n.Valid {
		return nullString
	}
	return n.Time.Format(time.RFC3339Nano)
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullTime.
func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%q", n.String())), nil
	}
	return jsonNullBytes, nil
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullTime.
func (n *NullTime) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Time = time.Time{}
		n.Valid = false
		return nil
	}
	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	s := string(payload)
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to time.Time: got %v", string(payload))
	}
	n.Time = t
	n.Valid = true
	return nil
}

// NullDate represents a Cloud Spanner DATE that may be null.
type NullDate struct {
	Date  civil.Date // Date contains the value when it is non-NULL, and a zero civil.Date when NULL.
	Valid bool       // Valid is true if Date is not NULL.
}

// IsNull implements NullableValue.IsNull for NullDate.
func (n NullDate) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullDate
func (n NullDate) String() string {
	if !n.Valid {
		return nullString
	}
	return n.Date.String()
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullDate.
func (n NullDate) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%q", n.String())), nil
	}
	return jsonNullBytes, nil
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullDate.
func (n *NullDate) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Date = civil.Date{}
		n.Valid = false
		return nil
	}
	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	s := string(payload)
	t, err := civil.ParseDate(s)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to civil.Date: got %v", string(payload))
	}
	n.Date = t
	n.Valid = true
	return nil
}

// NullNumeric represents a Cloud Spanner Numeric that may be NULL.
type NullNumeric struct {
	Numeric big.Rat // Numeric contains the value when it is non-NULL, and a zero big.Rat when NULL.
	Valid   bool    // Valid is true if Numeric is not NULL.
}

// IsNull implements NullableValue.IsNull for NullNumeric.
func (n NullNumeric) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullNumeric
func (n NullNumeric) String() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", NumericString(&n.Numeric))
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullNumeric.
func (n NullNumeric) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%q", NumericString(&n.Numeric))), nil
	}
	return jsonNullBytes, nil
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullNumeric.
func (n *NullNumeric) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Numeric = big.Rat{}
		n.Valid = false
		return nil
	}
	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	s := string(payload)
	val, ok := (&big.Rat{}).SetString(s)
	if !ok {
		return fmt.Errorf("payload cannot be converted to big.Rat: got %v", string(payload))
	}
	n.Numeric = *val
	n.Valid = true
	return nil
}

// NullRow represents a Cloud Spanner STRUCT that may be NULL.
// See also the document for Row.
// Note that NullRow is not a valid Cloud Spanner column Type.
type NullRow struct {
	Row   Row  // Row contains the value when it is non-NULL, and a zero Row when NULL.
	Valid bool // Valid is true if Row is not NULL.
}

// GenericColumnValue represents the generic encoded value and type of the
// column.  See google.spanner.v1.ResultSet proto for details.  This can be
// useful for proxying query results when the result types are not known in
// advance.
//
// If you populate a GenericColumnValue from a row using Row.Column or related
// methods, do not modify the contents of Type and Value.
type GenericColumnValue struct {
	Type  *sppb.Type
	Value *proto3.Value
}

// Decode decodes a GenericColumnValue. The ptr argument should be a pointer
// to a Go value that can accept v.
func (v GenericColumnValue) Decode(ptr interface{}) error {
	return decodeValue(v.Value, v.Type, ptr)
}

// NewGenericColumnValue creates a GenericColumnValue from Go value that is
// valid for Cloud Spanner.
func newGenericColumnValue(v interface{}) (*GenericColumnValue, error) {
	value, typ, err := encodeValue(v)
	if err != nil {
		return nil, err
	}
	return &GenericColumnValue{Value: value, Type: typ}, nil
}

// errTypeMismatch returns error for destination not having a compatible type
// with source Cloud Spanner type.
func errTypeMismatch(srcCode, elCode sppb.TypeCode, dst interface{}) error {
	s := srcCode.String()
	if srcCode == sppb.TypeCode_ARRAY {
		s = fmt.Sprintf("%v[%v]", srcCode, elCode)
	}
	return spannerErrorf(codes.InvalidArgument, "type %T cannot be used for decoding %s", dst, s)
}

// errNilSpannerType returns error for nil Cloud Spanner type in decoding.
func errNilSpannerType() error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected nil Cloud Spanner data type in decoding")
}

// errNilSrc returns error for decoding from nil proto value.
func errNilSrc() error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected nil Cloud Spanner value in decoding")
}

// errNilDst returns error for decoding into nil interface{}.
func errNilDst(dst interface{}) error {
	return spannerErrorf(codes.InvalidArgument, "cannot decode into nil type %T", dst)
}

// errNilArrElemType returns error for input Cloud Spanner data type being a array but without a
// non-nil array element type.
func errNilArrElemType(t *sppb.Type) error {
	return spannerErrorf(codes.FailedPrecondition, "array type %v is with nil array element type", t)
}

func errUnsupportedEmbeddedStructFields(fname string) error {
	return spannerErrorf(codes.InvalidArgument, "Embedded field: %s. Embedded and anonymous fields are not allowed "+
		"when converting Go structs to Cloud Spanner STRUCT values. To create a STRUCT value with an "+
		"unnamed field, use a `spanner:\"\"` field tag.", fname)
}

// errDstNotForNull returns error for decoding a SQL NULL value into a destination which doesn't
// support NULL values.
func errDstNotForNull(dst interface{}) error {
	return spannerErrorf(codes.InvalidArgument, "destination %T cannot support NULL SQL values", dst)
}

// errBadEncoding returns error for decoding wrongly encoded types.
func errBadEncoding(v *proto3.Value, err error) error {
	return spannerErrorf(codes.FailedPrecondition, "%v wasn't correctly encoded: <%v>", v, err)
}

func parseNullTime(v *proto3.Value, p *NullTime, code sppb.TypeCode, isNull bool) error {
	if p == nil {
		return errNilDst(p)
	}
	if code != sppb.TypeCode_TIMESTAMP {
		return errTypeMismatch(code, sppb.TypeCode_TYPE_CODE_UNSPECIFIED, p)
	}
	if isNull {
		*p = NullTime{}
		return nil
	}
	x, err := getStringValue(v)
	if err != nil {
		return err
	}
	y, err := time.Parse(time.RFC3339Nano, x)
	if err != nil {
		return errBadEncoding(v, err)
	}
	p.Valid = true
	p.Time = y
	return nil
}

// decodeValue decodes a protobuf Value into a pointer to a Go value, as
// specified by sppb.Type.
func decodeValue(v *proto3.Value, t *sppb.Type, ptr interface{}) error {
	if v == nil {
		return errNilSrc()
	}
	if t == nil {
		return errNilSpannerType()
	}
	code := t.Code
	acode := sppb.TypeCode_TYPE_CODE_UNSPECIFIED
	if code == sppb.TypeCode_ARRAY {
		if t.ArrayElementType == nil {
			return errNilArrElemType(t)
		}
		acode = t.ArrayElementType.Code
	}
	_, isNull := v.Kind.(*proto3.Value_NullValue)

	// Do the decoding based on the type of ptr.
	switch p := ptr.(type) {
	case nil:
		return errNilDst(nil)
	case *string:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		*p = x
	case *NullString, **string:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullString:
				*sp = NullString{}
			case **string:
				*sp = nil
			}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *NullString:
			sp.Valid = true
			sp.StringVal = x
		case **string:
			*sp = &x
		}
	case *[]NullString, *[]*string:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullString:
				*sp = nil
			case *[]*string:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullString:
			y, err := decodeNullStringArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*string:
			y, err := decodeStringPointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]string:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeStringArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *[]byte:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := base64.StdEncoding.DecodeString(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		*p = y
	case *[][]byte:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeByteArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *int64:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return errBadEncoding(v, err)
		}
		*p = y
	case *NullInt64, **int64:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullInt64:
				*sp = NullInt64{}
			case **int64:
				*sp = nil
			}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return errBadEncoding(v, err)
		}
		switch sp := ptr.(type) {
		case *NullInt64:
			sp.Valid = true
			sp.Int64 = y
		case **int64:
			*sp = &y
		}
	case *[]NullInt64, *[]*int64:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullInt64:
				*sp = nil
			case *[]*int64:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullInt64:
			y, err := decodeNullInt64Array(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*int64:
			y, err := decodeInt64PointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]int64:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeInt64Array(x)
		if err != nil {
			return err
		}
		*p = y
	case *bool:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getBoolValue(v)
		if err != nil {
			return err
		}
		*p = x
	case *NullBool, **bool:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullBool:
				*sp = NullBool{}
			case **bool:
				*sp = nil
			}
			break
		}
		x, err := getBoolValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *NullBool:
			sp.Valid = true
			sp.Bool = x
		case **bool:
			*sp = &x
		}
	case *[]NullBool, *[]*bool:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullBool:
				*sp = nil
			case *[]*bool:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullBool:
			y, err := decodeNullBoolArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*bool:
			y, err := decodeBoolPointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]bool:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeBoolArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *float64:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getFloat64Value(v)
		if err != nil {
			return err
		}
		*p = x
	case *NullFloat64, **float64:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullFloat64:
				*sp = NullFloat64{}
			case **float64:
				*sp = nil
			}
			break
		}
		x, err := getFloat64Value(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *NullFloat64:
			sp.Valid = true
			sp.Float64 = x
		case **float64:
			*sp = &x
		}
	case *[]NullFloat64, *[]*float64:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullFloat64:
				*sp = nil
			case *[]*float64:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullFloat64:
			y, err := decodeNullFloat64Array(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*float64:
			y, err := decodeFloat64PointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]float64:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeFloat64Array(x)
		if err != nil {
			return err
		}
		*p = y
	case *big.Rat:
		if code != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x := v.GetStringValue()
		y, ok := (&big.Rat{}).SetString(x)
		if !ok {
			return errUnexpectedNumericStr(x)
		}
		*p = *y
	case *NullNumeric:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = NullNumeric{}
			break
		}
		x := v.GetStringValue()
		y, ok := (&big.Rat{}).SetString(x)
		if !ok {
			return errUnexpectedNumericStr(x)
		}
		*p = NullNumeric{*y, true}
	case **big.Rat:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x := v.GetStringValue()
		y, ok := (&big.Rat{}).SetString(x)
		if !ok {
			return errUnexpectedNumericStr(x)
		}
		*p = y
	case *[]NullNumeric, *[]*big.Rat:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullNumeric:
				*sp = nil
			case *[]*big.Rat:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullNumeric:
			y, err := decodeNullNumericArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*big.Rat:
			y, err := decodeNumericPointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]big.Rat:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeNumericArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *time.Time:
		var nt NullTime
		if isNull {
			return errDstNotForNull(ptr)
		}
		err := parseNullTime(v, &nt, code, isNull)
		if err != nil {
			return err
		}
		*p = nt.Time
	case *NullTime:
		err := parseNullTime(v, p, code, isNull)
		if err != nil {
			return err
		}
	case **time.Time:
		var nt NullTime
		if isNull {
			*p = nil
			break
		}
		err := parseNullTime(v, &nt, code, isNull)
		if err != nil {
			return err
		}
		*p = &nt.Time
	case *[]NullTime, *[]*time.Time:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_TIMESTAMP {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullTime:
				*sp = nil
			case *[]*time.Time:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullTime:
			y, err := decodeNullTimeArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*time.Time:
			y, err := decodeTimePointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]time.Time:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_TIMESTAMP {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeTimeArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *civil.Date:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := civil.ParseDate(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		*p = y
	case *NullDate, **civil.Date:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullDate:
				*sp = NullDate{}
			case **civil.Date:
				*sp = nil
			}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := civil.ParseDate(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		switch sp := ptr.(type) {
		case *NullDate:
			sp.Valid = true
			sp.Date = y
		case **civil.Date:
			*sp = &y
		}
	case *[]NullDate, *[]*civil.Date:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullDate:
				*sp = nil
			case *[]*civil.Date:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullDate:
			y, err := decodeNullDateArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*civil.Date:
			y, err := decodeDatePointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]civil.Date:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeDateArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *[]NullRow:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_STRUCT {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeRowArray(t.ArrayElementType.StructType, x)
		if err != nil {
			return err
		}
		*p = y
	case *GenericColumnValue:
		*p = GenericColumnValue{Type: t, Value: v}
	default:
		// Check if the pointer is a custom type that implements spanner.Decoder
		// interface.
		if decodedVal, ok := ptr.(Decoder); ok {
			x, err := getGenericValue(v)
			if err != nil {
				return err
			}
			return decodedVal.DecodeSpanner(x)
		}

		// Check if the pointer is a variant of a base type.
		decodableType := getDecodableSpannerType(ptr, true)
		if decodableType != spannerTypeUnknown {
			if isNull && !decodableType.supportsNull() {
				return errDstNotForNull(ptr)
			}
			return decodableType.decodeValueToCustomType(v, t, acode, ptr)
		}

		// Check if the proto encoding is for an array of structs.
		if !(code == sppb.TypeCode_ARRAY && acode == sppb.TypeCode_STRUCT) {
			return errTypeMismatch(code, acode, ptr)
		}
		vp := reflect.ValueOf(p)
		if !vp.IsValid() {
			return errNilDst(p)
		}
		if !isPtrStructPtrSlice(vp.Type()) {
			// The container is not a pointer to a struct pointer slice.
			return errTypeMismatch(code, acode, ptr)
		}
		// Only use reflection for nil detection on slow path.
		// Also, IsNil panics on many types, so check it after the type check.
		if vp.IsNil() {
			return errNilDst(p)
		}
		if isNull {
			// The proto Value is encoding NULL, set the pointer to struct
			// slice to nil as well.
			vp.Elem().Set(reflect.Zero(vp.Elem().Type()))
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		if err = decodeStructArray(t.ArrayElementType.StructType, x, p); err != nil {
			return err
		}
	}
	return nil
}

// decodableSpannerType represents the Go types that a value from a Spanner
// database can be converted to.
type decodableSpannerType uint

const (
	spannerTypeUnknown decodableSpannerType = iota
	spannerTypeInvalid
	spannerTypeNonNullString
	spannerTypeByteArray
	spannerTypeNonNullInt64
	spannerTypeNonNullBool
	spannerTypeNonNullFloat64
	spannerTypeNonNullNumeric
	spannerTypeNonNullTime
	spannerTypeNonNullDate
	spannerTypeNullString
	spannerTypeNullInt64
	spannerTypeNullBool
	spannerTypeNullFloat64
	spannerTypeNullTime
	spannerTypeNullDate
	spannerTypeNullNumeric
	spannerTypeArrayOfNonNullString
	spannerTypeArrayOfByteArray
	spannerTypeArrayOfNonNullInt64
	spannerTypeArrayOfNonNullBool
	spannerTypeArrayOfNonNullFloat64
	spannerTypeArrayOfNonNullNumeric
	spannerTypeArrayOfNonNullTime
	spannerTypeArrayOfNonNullDate
	spannerTypeArrayOfNullString
	spannerTypeArrayOfNullInt64
	spannerTypeArrayOfNullBool
	spannerTypeArrayOfNullFloat64
	spannerTypeArrayOfNullNumeric
	spannerTypeArrayOfNullTime
	spannerTypeArrayOfNullDate
)

// supportsNull returns true for the Go types that can hold a null value from
// Spanner.
func (d decodableSpannerType) supportsNull() bool {
	switch d {
	case spannerTypeNonNullString, spannerTypeNonNullInt64, spannerTypeNonNullBool, spannerTypeNonNullFloat64, spannerTypeNonNullTime, spannerTypeNonNullDate, spannerTypeNonNullNumeric:
		return false
	default:
		return true
	}
}

// The following list of types represent the struct types that represent a
// specific Spanner data type in Go. If a pointer to one of these types is
// passed to decodeValue, the client library will decode one column value into
// the struct. For pointers to all other struct types, the client library will
// treat it as a generic struct that should contain a field for each column in
// the result set that is being decoded.

var typeOfNonNullTime = reflect.TypeOf(time.Time{})
var typeOfNonNullDate = reflect.TypeOf(civil.Date{})
var typeOfNonNullNumeric = reflect.TypeOf(big.Rat{})
var typeOfNullString = reflect.TypeOf(NullString{})
var typeOfNullInt64 = reflect.TypeOf(NullInt64{})
var typeOfNullBool = reflect.TypeOf(NullBool{})
var typeOfNullFloat64 = reflect.TypeOf(NullFloat64{})
var typeOfNullTime = reflect.TypeOf(NullTime{})
var typeOfNullDate = reflect.TypeOf(NullDate{})
var typeOfNullNumeric = reflect.TypeOf(NullNumeric{})

// getDecodableSpannerType returns the corresponding decodableSpannerType of
// the given pointer.
func getDecodableSpannerType(ptr interface{}, isPtr bool) decodableSpannerType {
	var val reflect.Value
	var kind reflect.Kind
	if isPtr {
		val = reflect.Indirect(reflect.ValueOf(ptr))
	} else {
		val = reflect.ValueOf(ptr)
	}
	kind = val.Kind()
	if kind == reflect.Invalid {
		return spannerTypeInvalid
	}
	switch kind {
	case reflect.Invalid:
		return spannerTypeInvalid
	case reflect.String:
		return spannerTypeNonNullString
	case reflect.Int64:
		return spannerTypeNonNullInt64
	case reflect.Bool:
		return spannerTypeNonNullBool
	case reflect.Float64:
		return spannerTypeNonNullFloat64
	case reflect.Ptr:
		t := val.Type()
		if t.ConvertibleTo(typeOfNullNumeric) {
			return spannerTypeNullNumeric
		}
	case reflect.Struct:
		t := val.Type()
		if t.ConvertibleTo(typeOfNonNullNumeric) {
			return spannerTypeNonNullNumeric
		}
		if t.ConvertibleTo(typeOfNonNullTime) {
			return spannerTypeNonNullTime
		}
		if t.ConvertibleTo(typeOfNonNullDate) {
			return spannerTypeNonNullDate
		}
		if t.ConvertibleTo(typeOfNullString) {
			return spannerTypeNullString
		}
		if t.ConvertibleTo(typeOfNullInt64) {
			return spannerTypeNullInt64
		}
		if t.ConvertibleTo(typeOfNullBool) {
			return spannerTypeNullBool
		}
		if t.ConvertibleTo(typeOfNullFloat64) {
			return spannerTypeNullFloat64
		}
		if t.ConvertibleTo(typeOfNullTime) {
			return spannerTypeNullTime
		}
		if t.ConvertibleTo(typeOfNullDate) {
			return spannerTypeNullDate
		}
		if t.ConvertibleTo(typeOfNullNumeric) {
			return spannerTypeNullNumeric
		}
	case reflect.Slice:
		kind := val.Type().Elem().Kind()
		switch kind {
		case reflect.Invalid:
			return spannerTypeUnknown
		case reflect.String:
			return spannerTypeArrayOfNonNullString
		case reflect.Uint8:
			return spannerTypeByteArray
		case reflect.Int64:
			return spannerTypeArrayOfNonNullInt64
		case reflect.Bool:
			return spannerTypeArrayOfNonNullBool
		case reflect.Float64:
			return spannerTypeArrayOfNonNullFloat64
		case reflect.Ptr:
			t := val.Type().Elem()
			if t.ConvertibleTo(typeOfNullNumeric) {
				return spannerTypeArrayOfNullNumeric
			}
		case reflect.Struct:
			t := val.Type().Elem()
			if t.ConvertibleTo(typeOfNonNullNumeric) {
				return spannerTypeArrayOfNonNullNumeric
			}
			if t.ConvertibleTo(typeOfNonNullTime) {
				return spannerTypeArrayOfNonNullTime
			}
			if t.ConvertibleTo(typeOfNonNullDate) {
				return spannerTypeArrayOfNonNullDate
			}
			if t.ConvertibleTo(typeOfNullString) {
				return spannerTypeArrayOfNullString
			}
			if t.ConvertibleTo(typeOfNullInt64) {
				return spannerTypeArrayOfNullInt64
			}
			if t.ConvertibleTo(typeOfNullBool) {
				return spannerTypeArrayOfNullBool
			}
			if t.ConvertibleTo(typeOfNullFloat64) {
				return spannerTypeArrayOfNullFloat64
			}
			if t.ConvertibleTo(typeOfNullTime) {
				return spannerTypeArrayOfNullTime
			}
			if t.ConvertibleTo(typeOfNullDate) {
				return spannerTypeArrayOfNullDate
			}
			if t.ConvertibleTo(typeOfNullNumeric) {
				return spannerTypeArrayOfNullNumeric
			}
		case reflect.Slice:
			// The only array-of-array type that is supported is [][]byte.
			kind := val.Type().Elem().Elem().Kind()
			switch kind {
			case reflect.Uint8:
				return spannerTypeArrayOfByteArray
			}
		}
	}
	// Not convertible to a known base type.
	return spannerTypeUnknown
}

// decodeValueToCustomType decodes a protobuf Value into a pointer to a Go
// value. It must be possible to convert the value to the type pointed to by
// the pointer.
func (dsc decodableSpannerType) decodeValueToCustomType(v *proto3.Value, t *sppb.Type, acode sppb.TypeCode, ptr interface{}) error {
	code := t.Code
	_, isNull := v.Kind.(*proto3.Value_NullValue)
	if dsc == spannerTypeInvalid {
		return errNilDst(ptr)
	}
	if isNull && !dsc.supportsNull() {
		return errDstNotForNull(ptr)
	}

	var result interface{}
	switch dsc {
	case spannerTypeNonNullString, spannerTypeNullString:
		if code != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullString{}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		if dsc == spannerTypeNonNullString {
			result = &x
		} else {
			result = &NullString{x, !isNull}
		}
	case spannerTypeByteArray:
		if code != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = []byte(nil)
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := base64.StdEncoding.DecodeString(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		result = y
	case spannerTypeNonNullInt64, spannerTypeNullInt64:
		if code != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullInt64{}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return errBadEncoding(v, err)
		}
		if dsc == spannerTypeNonNullInt64 {
			result = &y
		} else {
			result = &NullInt64{y, !isNull}
		}
	case spannerTypeNonNullBool, spannerTypeNullBool:
		if code != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullBool{}
			break
		}
		x, err := getBoolValue(v)
		if err != nil {
			return err
		}
		if dsc == spannerTypeNonNullBool {
			result = &x
		} else {
			result = &NullBool{x, !isNull}
		}
	case spannerTypeNonNullFloat64, spannerTypeNullFloat64:
		if code != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullFloat64{}
			break
		}
		x, err := getFloat64Value(v)
		if err != nil {
			return err
		}
		if dsc == spannerTypeNonNullFloat64 {
			result = &x
		} else {
			result = &NullFloat64{x, !isNull}
		}
	case spannerTypeNonNullNumeric, spannerTypeNullNumeric:
		if code != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullNumeric{}
			break
		}
		x := v.GetStringValue()
		y, ok := (&big.Rat{}).SetString(x)
		if !ok {
			return errUnexpectedNumericStr(x)
		}
		if dsc == spannerTypeNonNullNumeric {
			result = y
		} else {
			result = &NullNumeric{*y, true}
		}
	case spannerTypeNonNullTime, spannerTypeNullTime:
		var nt NullTime
		err := parseNullTime(v, &nt, code, isNull)
		if err != nil {
			return err
		}
		if dsc == spannerTypeNonNullTime {
			result = &nt.Time
		} else {
			result = &nt
		}
	case spannerTypeNonNullDate, spannerTypeNullDate:
		if code != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullDate{}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := civil.ParseDate(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		if dsc == spannerTypeNonNullDate {
			result = &y
		} else {
			result = &NullDate{y, !isNull}
		}
	case spannerTypeArrayOfNonNullString, spannerTypeArrayOfNullString:
		if acode != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, stringType(), "STRING")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfByteArray:
		if acode != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, bytesType(), "BYTES")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullInt64, spannerTypeArrayOfNullInt64:
		if acode != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, intType(), "INT64")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullBool, spannerTypeArrayOfNullBool:
		if acode != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, boolType(), "BOOL")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullFloat64, spannerTypeArrayOfNullFloat64:
		if acode != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, floatType(), "FLOAT64")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullNumeric, spannerTypeArrayOfNullNumeric:
		if acode != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, numericType(), "NUMERIC")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullTime, spannerTypeArrayOfNullTime:
		if acode != sppb.TypeCode_TIMESTAMP {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, timeType(), "TIMESTAMP")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullDate, spannerTypeArrayOfNullDate:
		if acode != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, dateType(), "DATE")
		if err != nil {
			return err
		}
		result = y
	default:
		// This should not be possible.
		return fmt.Errorf("unknown decodable type found: %v", dsc)
	}
	source := reflect.Indirect(reflect.ValueOf(result))
	destination := reflect.Indirect(reflect.ValueOf(ptr))
	destination.Set(source.Convert(destination.Type()))
	return nil
}

// errSrvVal returns an error for getting a wrong source protobuf value in decoding.
func errSrcVal(v *proto3.Value, want string) error {
	return spannerErrorf(codes.FailedPrecondition, "cannot use %v(Kind: %T) as %s Value",
		v, v.GetKind(), want)
}

// getStringValue returns the string value encoded in proto3.Value v whose
// kind is proto3.Value_StringValue.
func getStringValue(v *proto3.Value) (string, error) {
	if x, ok := v.GetKind().(*proto3.Value_StringValue); ok && x != nil {
		return x.StringValue, nil
	}
	return "", errSrcVal(v, "String")
}

// getBoolValue returns the bool value encoded in proto3.Value v whose
// kind is proto3.Value_BoolValue.
func getBoolValue(v *proto3.Value) (bool, error) {
	if x, ok := v.GetKind().(*proto3.Value_BoolValue); ok && x != nil {
		return x.BoolValue, nil
	}
	return false, errSrcVal(v, "Bool")
}

// getListValue returns the proto3.ListValue contained in proto3.Value v whose
// kind is proto3.Value_ListValue.
func getListValue(v *proto3.Value) (*proto3.ListValue, error) {
	if x, ok := v.GetKind().(*proto3.Value_ListValue); ok && x != nil {
		return x.ListValue, nil
	}
	return nil, errSrcVal(v, "List")
}

// getGenericValue returns the interface{} value encoded in proto3.Value.
func getGenericValue(v *proto3.Value) (interface{}, error) {
	switch x := v.GetKind().(type) {
	case *proto3.Value_NumberValue:
		return x.NumberValue, nil
	case *proto3.Value_BoolValue:
		return x.BoolValue, nil
	case *proto3.Value_StringValue:
		return x.StringValue, nil
	default:
		return 0, errSrcVal(v, "Number, Bool, String")
	}
}

// errUnexpectedNumericStr returns error for decoder getting an unexpected
// string for representing special numeric values.
func errUnexpectedNumericStr(s string) error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected string value %q for numeric number", s)
}

// errUnexpectedFloat64Str returns error for decoder getting an unexpected
// string for representing special float values.
func errUnexpectedFloat64Str(s string) error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected string value %q for float64 number", s)
}

// getFloat64Value returns the float64 value encoded in proto3.Value v whose
// kind is proto3.Value_NumberValue / proto3.Value_StringValue.
// Cloud Spanner uses string to encode NaN, Infinity and -Infinity.
func getFloat64Value(v *proto3.Value) (float64, error) {
	switch x := v.GetKind().(type) {
	case *proto3.Value_NumberValue:
		if x == nil {
			break
		}
		return x.NumberValue, nil
	case *proto3.Value_StringValue:
		if x == nil {
			break
		}
		switch x.StringValue {
		case "NaN":
			return math.NaN(), nil
		case "Infinity":
			return math.Inf(1), nil
		case "-Infinity":
			return math.Inf(-1), nil
		default:
			return 0, errUnexpectedFloat64Str(x.StringValue)
		}
	}
	return 0, errSrcVal(v, "Number")
}

// errNilListValue returns error for unexpected nil ListValue in decoding Cloud Spanner ARRAYs.
func errNilListValue(sqlType string) error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected nil ListValue in decoding %v array", sqlType)
}

// errDecodeArrayElement returns error for failure in decoding single array element.
func errDecodeArrayElement(i int, v proto.Message, sqlType string, err error) error {
	var se *Error
	if !errorAs(err, &se) {
		return spannerErrorf(codes.Unknown,
			"cannot decode %v(array element %v) as %v, error = <%v>", v, i, sqlType, err)
	}
	se.decorate(fmt.Sprintf("cannot decode %v(array element %v) as %v", v, i, sqlType))
	return se
}

// decodeGenericArray decodes proto3.ListValue pb into a slice which type is
// determined through reflection.
func decodeGenericArray(tp reflect.Type, pb *proto3.ListValue, t *sppb.Type, sqlType string) (interface{}, error) {
	if pb == nil {
		return nil, errNilListValue(sqlType)
	}
	a := reflect.MakeSlice(tp, len(pb.Values), len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, t, a.Index(i).Addr().Interface()); err != nil {
			return nil, errDecodeArrayElement(i, v, "STRING", err)
		}
	}
	return a.Interface(), nil
}

// decodeNullStringArray decodes proto3.ListValue pb into a NullString slice.
func decodeNullStringArray(pb *proto3.ListValue) ([]NullString, error) {
	if pb == nil {
		return nil, errNilListValue("STRING")
	}
	a := make([]NullString, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, stringType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "STRING", err)
		}
	}
	return a, nil
}

// decodeStringPointerArray decodes proto3.ListValue pb into a *string slice.
func decodeStringPointerArray(pb *proto3.ListValue) ([]*string, error) {
	if pb == nil {
		return nil, errNilListValue("STRING")
	}
	a := make([]*string, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, stringType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "STRING", err)
		}
	}
	return a, nil
}

// decodeStringArray decodes proto3.ListValue pb into a string slice.
func decodeStringArray(pb *proto3.ListValue) ([]string, error) {
	if pb == nil {
		return nil, errNilListValue("STRING")
	}
	a := make([]string, len(pb.Values))
	st := stringType()
	for i, v := range pb.Values {
		if err := decodeValue(v, st, &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "STRING", err)
		}
	}
	return a, nil
}

// decodeNullInt64Array decodes proto3.ListValue pb into a NullInt64 slice.
func decodeNullInt64Array(pb *proto3.ListValue) ([]NullInt64, error) {
	if pb == nil {
		return nil, errNilListValue("INT64")
	}
	a := make([]NullInt64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, intType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "INT64", err)
		}
	}
	return a, nil
}

// decodeInt64PointerArray decodes proto3.ListValue pb into a *int64 slice.
func decodeInt64PointerArray(pb *proto3.ListValue) ([]*int64, error) {
	if pb == nil {
		return nil, errNilListValue("INT64")
	}
	a := make([]*int64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, intType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "INT64", err)
		}
	}
	return a, nil
}

// decodeInt64Array decodes proto3.ListValue pb into a int64 slice.
func decodeInt64Array(pb *proto3.ListValue) ([]int64, error) {
	if pb == nil {
		return nil, errNilListValue("INT64")
	}
	a := make([]int64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, intType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "INT64", err)
		}
	}
	return a, nil
}

// decodeNullBoolArray decodes proto3.ListValue pb into a NullBool slice.
func decodeNullBoolArray(pb *proto3.ListValue) ([]NullBool, error) {
	if pb == nil {
		return nil, errNilListValue("BOOL")
	}
	a := make([]NullBool, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, boolType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "BOOL", err)
		}
	}
	return a, nil
}

// decodeBoolPointerArray decodes proto3.ListValue pb into a *bool slice.
func decodeBoolPointerArray(pb *proto3.ListValue) ([]*bool, error) {
	if pb == nil {
		return nil, errNilListValue("BOOL")
	}
	a := make([]*bool, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, boolType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "BOOL", err)
		}
	}
	return a, nil
}

// decodeBoolArray decodes proto3.ListValue pb into a bool slice.
func decodeBoolArray(pb *proto3.ListValue) ([]bool, error) {
	if pb == nil {
		return nil, errNilListValue("BOOL")
	}
	a := make([]bool, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, boolType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "BOOL", err)
		}
	}
	return a, nil
}

// decodeNullFloat64Array decodes proto3.ListValue pb into a NullFloat64 slice.
func decodeNullFloat64Array(pb *proto3.ListValue) ([]NullFloat64, error) {
	if pb == nil {
		return nil, errNilListValue("FLOAT64")
	}
	a := make([]NullFloat64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, floatType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "FLOAT64", err)
		}
	}
	return a, nil
}

// decodeFloat64PointerArray decodes proto3.ListValue pb into a *float slice.
func decodeFloat64PointerArray(pb *proto3.ListValue) ([]*float64, error) {
	if pb == nil {
		return nil, errNilListValue("FLOAT64")
	}
	a := make([]*float64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, floatType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "FLOAT64", err)
		}
	}
	return a, nil
}

// decodeFloat64Array decodes proto3.ListValue pb into a float64 slice.
func decodeFloat64Array(pb *proto3.ListValue) ([]float64, error) {
	if pb == nil {
		return nil, errNilListValue("FLOAT64")
	}
	a := make([]float64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, floatType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "FLOAT64", err)
		}
	}
	return a, nil
}

// decodeNullNumericArray decodes proto3.ListValue pb into a NullNumeric slice.
func decodeNullNumericArray(pb *proto3.ListValue) ([]NullNumeric, error) {
	if pb == nil {
		return nil, errNilListValue("NUMERIC")
	}
	a := make([]NullNumeric, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, numericType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "NUMERIC", err)
		}
	}
	return a, nil
}

// decodeNumericPointerArray decodes proto3.ListValue pb into a *big.Rat slice.
func decodeNumericPointerArray(pb *proto3.ListValue) ([]*big.Rat, error) {
	if pb == nil {
		return nil, errNilListValue("NUMERIC")
	}
	a := make([]*big.Rat, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, numericType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "NUMERIC", err)
		}
	}
	return a, nil
}

// decodeNumericArray decodes proto3.ListValue pb into a big.Rat slice.
func decodeNumericArray(pb *proto3.ListValue) ([]big.Rat, error) {
	if pb == nil {
		return nil, errNilListValue("NUMERIC")
	}
	a := make([]big.Rat, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, numericType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "NUMERIC", err)
		}
	}
	return a, nil
}

// decodeByteArray decodes proto3.ListValue pb into a slice of byte slice.
func decodeByteArray(pb *proto3.ListValue) ([][]byte, error) {
	if pb == nil {
		return nil, errNilListValue("BYTES")
	}
	a := make([][]byte, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, bytesType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "BYTES", err)
		}
	}
	return a, nil
}

// decodeNullTimeArray decodes proto3.ListValue pb into a NullTime slice.
func decodeNullTimeArray(pb *proto3.ListValue) ([]NullTime, error) {
	if pb == nil {
		return nil, errNilListValue("TIMESTAMP")
	}
	a := make([]NullTime, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, timeType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "TIMESTAMP", err)
		}
	}
	return a, nil
}

// decodeTimePointerArray decodes proto3.ListValue pb into a NullTime slice.
func decodeTimePointerArray(pb *proto3.ListValue) ([]*time.Time, error) {
	if pb == nil {
		return nil, errNilListValue("TIMESTAMP")
	}
	a := make([]*time.Time, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, timeType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "TIMESTAMP", err)
		}
	}
	return a, nil
}

// decodeTimeArray decodes proto3.ListValue pb into a time.Time slice.
func decodeTimeArray(pb *proto3.ListValue) ([]time.Time, error) {
	if pb == nil {
		return nil, errNilListValue("TIMESTAMP")
	}
	a := make([]time.Time, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, timeType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "TIMESTAMP", err)
		}
	}
	return a, nil
}

// decodeNullDateArray decodes proto3.ListValue pb into a NullDate slice.
func decodeNullDateArray(pb *proto3.ListValue) ([]NullDate, error) {
	if pb == nil {
		return nil, errNilListValue("DATE")
	}
	a := make([]NullDate, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, dateType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "DATE", err)
		}
	}
	return a, nil
}

// decodeDatePointerArray decodes proto3.ListValue pb into a *civil.Date slice.
func decodeDatePointerArray(pb *proto3.ListValue) ([]*civil.Date, error) {
	if pb == nil {
		return nil, errNilListValue("DATE")
	}
	a := make([]*civil.Date, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, dateType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "DATE", err)
		}
	}
	return a, nil
}

// decodeDateArray decodes proto3.ListValue pb into a civil.Date slice.
func decodeDateArray(pb *proto3.ListValue) ([]civil.Date, error) {
	if pb == nil {
		return nil, errNilListValue("DATE")
	}
	a := make([]civil.Date, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, dateType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "DATE", err)
		}
	}
	return a, nil
}

func errNotStructElement(i int, v *proto3.Value) error {
	return errDecodeArrayElement(i, v, "STRUCT",
		spannerErrorf(codes.FailedPrecondition, "%v(type: %T) doesn't encode Cloud Spanner STRUCT", v, v))
}

// decodeRowArray decodes proto3.ListValue pb into a NullRow slice according to
// the structural information given in sppb.StructType ty.
func decodeRowArray(ty *sppb.StructType, pb *proto3.ListValue) ([]NullRow, error) {
	if pb == nil {
		return nil, errNilListValue("STRUCT")
	}
	a := make([]NullRow, len(pb.Values))
	for i := range pb.Values {
		switch v := pb.Values[i].GetKind().(type) {
		case *proto3.Value_ListValue:
			a[i] = NullRow{
				Row: Row{
					fields: ty.Fields,
					vals:   v.ListValue.Values,
				},
				Valid: true,
			}
		// Null elements not currently supported by the server, see
		// https://cloud.google.com/spanner/docs/query-syntax#using-structs-with-select
		case *proto3.Value_NullValue:
			// no-op, a[i] is NullRow{} already
		default:
			return nil, errNotStructElement(i, pb.Values[i])
		}
	}
	return a, nil
}

// errNilSpannerStructType returns error for unexpected nil Cloud Spanner STRUCT
// schema type in decoding.
func errNilSpannerStructType() error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected nil StructType in decoding Cloud Spanner STRUCT")
}

// errUnnamedField returns error for decoding a Cloud Spanner STRUCT with
// unnamed field into a Go struct.
func errUnnamedField(ty *sppb.StructType, i int) error {
	return spannerErrorf(codes.InvalidArgument, "unnamed field %v in Cloud Spanner STRUCT %+v", i, ty)
}

// errNoOrDupGoField returns error for decoding a Cloud Spanner
// STRUCT into a Go struct which is either missing a field, or has duplicate
// fields.
func errNoOrDupGoField(s interface{}, f string) error {
	return spannerErrorf(codes.InvalidArgument, "Go struct %+v(type %T) has no or duplicate fields for Cloud Spanner STRUCT field %v", s, s, f)
}

// errDupColNames returns error for duplicated Cloud Spanner STRUCT field names
// found in decoding a Cloud Spanner STRUCT into a Go struct.
func errDupSpannerField(f string, ty *sppb.StructType) error {
	return spannerErrorf(codes.InvalidArgument, "duplicated field name %q in Cloud Spanner STRUCT %+v", f, ty)
}

// errDecodeStructField returns error for failure in decoding a single field of
// a Cloud Spanner STRUCT.
func errDecodeStructField(ty *sppb.StructType, f string, err error) error {
	var se *Error
	if !errorAs(err, &se) {
		return spannerErrorf(codes.Unknown,
			"cannot decode field %v of Cloud Spanner STRUCT %+v, error = <%v>", f, ty, err)
	}
	se.decorate(fmt.Sprintf("cannot decode field %v of Cloud Spanner STRUCT %+v", f, ty))
	return se
}

// decodeStruct decodes proto3.ListValue pb into struct referenced by pointer
// ptr, according to
// the structural information given in sppb.StructType ty.
func decodeStruct(ty *sppb.StructType, pb *proto3.ListValue, ptr interface{}) error {
	if reflect.ValueOf(ptr).IsNil() {
		return errNilDst(ptr)
	}
	if ty == nil {
		return errNilSpannerStructType()
	}
	// t holds the structural information of ptr.
	t := reflect.TypeOf(ptr).Elem()
	// v is the actual value that ptr points to.
	v := reflect.ValueOf(ptr).Elem()

	fields, err := fieldCache.Fields(t)
	if err != nil {
		return ToSpannerError(err)
	}
	seen := map[string]bool{}
	for i, f := range ty.Fields {
		if f.Name == "" {
			return errUnnamedField(ty, i)
		}
		sf := fields.Match(f.Name)
		if sf == nil {
			return errNoOrDupGoField(ptr, f.Name)
		}
		if seen[f.Name] {
			// We don't allow duplicated field name.
			return errDupSpannerField(f.Name, ty)
		}
		// Try to decode a single field.
		if err := decodeValue(pb.Values[i], f.Type, v.FieldByIndex(sf.Index).Addr().Interface()); err != nil {
			return errDecodeStructField(ty, f.Name, err)
		}
		// Mark field f.Name as processed.
		seen[f.Name] = true
	}
	return nil
}

// isPtrStructPtrSlice returns true if ptr is a pointer to a slice of struct pointers.
func isPtrStructPtrSlice(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice {
		// t is not a pointer to a slice.
		return false
	}
	if t = t.Elem(); t.Elem().Kind() != reflect.Ptr || t.Elem().Elem().Kind() != reflect.Struct {
		// the slice that t points to is not a slice of struct pointers.
		return false
	}
	return true
}

// decodeStructArray decodes proto3.ListValue pb into struct slice referenced by
// pointer ptr, according to the
// structural information given in a sppb.StructType.
func decodeStructArray(ty *sppb.StructType, pb *proto3.ListValue, ptr interface{}) error {
	if pb == nil {
		return errNilListValue("STRUCT")
	}
	// Type of the struct pointers stored in the slice that ptr points to.
	ts := reflect.TypeOf(ptr).Elem().Elem()
	// The slice that ptr points to, might be nil at this point.
	v := reflect.ValueOf(ptr).Elem()
	// Allocate empty slice.
	v.Set(reflect.MakeSlice(v.Type(), 0, len(pb.Values)))
	// Decode every struct in pb.Values.
	for i, pv := range pb.Values {
		// Check if pv is a NULL value.
		if _, isNull := pv.Kind.(*proto3.Value_NullValue); isNull {
			// Append a nil pointer to the slice.
			v.Set(reflect.Append(v, reflect.New(ts).Elem()))
			continue
		}
		// Allocate empty struct.
		s := reflect.New(ts.Elem())
		// Get proto3.ListValue l from proto3.Value pv.
		l, err := getListValue(pv)
		if err != nil {
			return errDecodeArrayElement(i, pv, "STRUCT", err)
		}
		// Decode proto3.ListValue l into struct referenced by s.Interface().
		if err = decodeStruct(ty, l, s.Interface()); err != nil {
			return errDecodeArrayElement(i, pv, "STRUCT", err)
		}
		// Append the decoded struct back into the slice.
		v.Set(reflect.Append(v, s))
	}
	return nil
}

// errEncoderUnsupportedType returns error for not being able to encode a value
// of certain type.
func errEncoderUnsupportedType(v interface{}) error {
	return spannerErrorf(codes.InvalidArgument, "client doesn't support type %T", v)
}

// encodeValue encodes a Go native type into a proto3.Value.
func encodeValue(v interface{}) (*proto3.Value, *sppb.Type, error) {
	pb := &proto3.Value{
		Kind: &proto3.Value_NullValue{NullValue: proto3.NullValue_NULL_VALUE},
	}
	var pt *sppb.Type
	var err error
	switch v := v.(type) {
	case nil:
	case string:
		pb.Kind = stringKind(v)
		pt = stringType()
	case NullString:
		if v.Valid {
			return encodeValue(v.StringVal)
		}
		pt = stringType()
	case []string:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(stringType())
	case []NullString:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(stringType())
	case *string:
		if v != nil {
			return encodeValue(*v)
		}
		pt = stringType()
	case []*string:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(stringType())
	case []byte:
		if v != nil {
			pb.Kind = stringKind(base64.StdEncoding.EncodeToString(v))
		}
		pt = bytesType()
	case [][]byte:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(bytesType())
	case int:
		pb.Kind = stringKind(strconv.FormatInt(int64(v), 10))
		pt = intType()
	case []int:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(intType())
	case int64:
		pb.Kind = stringKind(strconv.FormatInt(v, 10))
		pt = intType()
	case []int64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(intType())
	case NullInt64:
		if v.Valid {
			return encodeValue(v.Int64)
		}
		pt = intType()
	case []NullInt64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(intType())
	case *int64:
		if v != nil {
			return encodeValue(*v)
		}
		pt = intType()
	case []*int64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(intType())
	case bool:
		pb.Kind = &proto3.Value_BoolValue{BoolValue: v}
		pt = boolType()
	case []bool:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(boolType())
	case NullBool:
		if v.Valid {
			return encodeValue(v.Bool)
		}
		pt = boolType()
	case []NullBool:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(boolType())
	case *bool:
		if v != nil {
			return encodeValue(*v)
		}
		pt = boolType()
	case []*bool:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(boolType())
	case float64:
		pb.Kind = &proto3.Value_NumberValue{NumberValue: v}
		pt = floatType()
	case []float64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(floatType())
	case NullFloat64:
		if v.Valid {
			return encodeValue(v.Float64)
		}
		pt = floatType()
	case []NullFloat64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(floatType())
	case *float64:
		if v != nil {
			return encodeValue(*v)
		}
		pt = floatType()
	case []*float64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(floatType())
	case big.Rat:
		pb.Kind = stringKind(NumericString(&v))
		pt = numericType()
	case []big.Rat:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(numericType())
	case NullNumeric:
		if v.Valid {
			return encodeValue(v.Numeric)
		}
		pt = numericType()
	case []NullNumeric:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(numericType())
	case *big.Rat:
		if v != nil {
			pb.Kind = stringKind(NumericString(v))
		}
		pt = numericType()
	case []*big.Rat:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(numericType())
	case time.Time:
		if v == commitTimestamp {
			pb.Kind = stringKind(commitTimestampPlaceholderString)
		} else {
			pb.Kind = stringKind(v.UTC().Format(time.RFC3339Nano))
		}
		pt = timeType()
	case []time.Time:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(timeType())
	case NullTime:
		if v.Valid {
			return encodeValue(v.Time)
		}
		pt = timeType()
	case []NullTime:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(timeType())
	case *time.Time:
		if v != nil {
			return encodeValue(*v)
		}
		pt = timeType()
	case []*time.Time:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(timeType())
	case civil.Date:
		pb.Kind = stringKind(v.String())
		pt = dateType()
	case []civil.Date:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(dateType())
	case NullDate:
		if v.Valid {
			return encodeValue(v.Date)
		}
		pt = dateType()
	case []NullDate:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(dateType())
	case *civil.Date:
		if v != nil {
			return encodeValue(*v)
		}
		pt = dateType()
	case []*civil.Date:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(dateType())
	case GenericColumnValue:
		// Deep clone to ensure subsequent changes to v before
		// transmission don't affect our encoded value.
		pb = proto.Clone(v.Value).(*proto3.Value)
		pt = proto.Clone(v.Type).(*sppb.Type)
	case []GenericColumnValue:
		return nil, nil, errEncoderUnsupportedType(v)
	default:
		// Check if the value is a custom type that implements spanner.Encoder
		// interface.
		if encodedVal, ok := v.(Encoder); ok {
			nv, err := encodedVal.EncodeSpanner()
			if err != nil {
				return nil, nil, err
			}
			return encodeValue(nv)
		}

		// Check if the value is a variant of a base type.
		decodableType := getDecodableSpannerType(v, false)
		if decodableType != spannerTypeUnknown && decodableType != spannerTypeInvalid {
			converted, err := convertCustomTypeValue(decodableType, v)
			if err != nil {
				return nil, nil, err
			}
			return encodeValue(converted)
		}

		if !isStructOrArrayOfStructValue(v) {
			return nil, nil, errEncoderUnsupportedType(v)
		}
		typ := reflect.TypeOf(v)

		// Value is a Go struct value/ptr.
		if (typ.Kind() == reflect.Struct) ||
			(typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct) {
			return encodeStruct(v)
		}

		// Value is a slice of Go struct values/ptrs.
		if typ.Kind() == reflect.Slice {
			return encodeStructArray(v)
		}
	}
	return pb, pt, nil
}

func convertCustomTypeValue(sourceType decodableSpannerType, v interface{}) (interface{}, error) {
	// destination will be initialized to a base type. The input value will be
	// converted to this type and copied to destination.
	var destination reflect.Value
	switch sourceType {
	case spannerTypeInvalid:
		return nil, fmt.Errorf("cannot encode a value to type spannerTypeInvalid")
	case spannerTypeNonNullString:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf("")))
	case spannerTypeNullString:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullString{})))
	case spannerTypeByteArray:
		// Return a nil array directly if the input value is nil instead of
		// creating an empty slice and returning that.
		if reflect.ValueOf(v).IsNil() {
			return []byte(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]byte{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeNonNullInt64:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(int64(0))))
	case spannerTypeNullInt64:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullInt64{})))
	case spannerTypeNonNullBool:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(false)))
	case spannerTypeNullBool:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullBool{})))
	case spannerTypeNonNullFloat64:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(float64(0.0))))
	case spannerTypeNullFloat64:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullFloat64{})))
	case spannerTypeNonNullTime:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(time.Time{})))
	case spannerTypeNullTime:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullTime{})))
	case spannerTypeNonNullDate:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(civil.Date{})))
	case spannerTypeNullDate:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullDate{})))
	case spannerTypeNonNullNumeric:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(big.Rat{})))
	case spannerTypeNullNumeric:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullNumeric{})))
	case spannerTypeArrayOfNonNullString:
		if reflect.ValueOf(v).IsNil() {
			return []string(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]string{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullString:
		if reflect.ValueOf(v).IsNil() {
			return []NullString(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullString{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfByteArray:
		if reflect.ValueOf(v).IsNil() {
			return [][]byte(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([][]byte{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullInt64:
		if reflect.ValueOf(v).IsNil() {
			return []int64(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]int64{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullInt64:
		if reflect.ValueOf(v).IsNil() {
			return []NullInt64(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullInt64{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullBool:
		if reflect.ValueOf(v).IsNil() {
			return []bool(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]bool{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullBool:
		if reflect.ValueOf(v).IsNil() {
			return []NullBool(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullBool{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullFloat64:
		if reflect.ValueOf(v).IsNil() {
			return []float64(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]float64{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullFloat64:
		if reflect.ValueOf(v).IsNil() {
			return []NullFloat64(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullFloat64{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullTime:
		if reflect.ValueOf(v).IsNil() {
			return []time.Time(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]time.Time{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullTime:
		if reflect.ValueOf(v).IsNil() {
			return []NullTime(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullTime{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullDate:
		if reflect.ValueOf(v).IsNil() {
			return []civil.Date(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]civil.Date{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullDate:
		if reflect.ValueOf(v).IsNil() {
			return []NullDate(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullDate{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullNumeric:
		if reflect.ValueOf(v).IsNil() {
			return []big.Rat(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]big.Rat{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullNumeric:
		if reflect.ValueOf(v).IsNil() {
			return []NullNumeric(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullNumeric{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	default:
		// This should not be possible.
		return nil, fmt.Errorf("unknown decodable type found: %v", sourceType)
	}
	// destination has been initialized. Convert and copy the input value to
	// destination. That must be done per element if the input type is a slice
	// or an array.
	if destination.Kind() == reflect.Slice || destination.Kind() == reflect.Array {
		sourceSlice := reflect.ValueOf(v)
		for i := 0; i < destination.Len(); i++ {
			source := sourceSlice.Index(i)
			destination.Index(i).Set(source.Convert(destination.Type().Elem()))
		}
	} else {
		source := reflect.ValueOf(v)
		destination.Set(source.Convert(destination.Type()))
	}
	// Return the converted value.
	return destination.Interface(), nil
}

// Encodes a Go struct value/ptr in v to the spanner Value and Type protos. v
// itself must be non-nil.
func encodeStruct(v interface{}) (*proto3.Value, *sppb.Type, error) {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	// Pointer to struct.
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		typ = typ.Elem()
		if val.IsNil() {
			// nil pointer to struct, representing a NULL STRUCT value. Use a
			// dummy value to get the type.
			_, st, err := encodeStruct(reflect.Zero(typ).Interface())
			if err != nil {
				return nil, nil, err
			}
			return nullProto(), st, nil
		}
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, nil, errEncoderUnsupportedType(v)
	}

	stf := make([]*sppb.StructType_Field, 0, typ.NumField())
	stv := make([]*proto3.Value, 0, typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		// If the field has a 'spanner' tag, use the value of that tag as the field name.
		// This is used to build STRUCT types with unnamed/duplicate fields.
		sf := typ.Field(i)
		fval := val.Field(i)

		// Embedded fields are not allowed.
		if sf.Anonymous {
			return nil, nil, errUnsupportedEmbeddedStructFields(sf.Name)
		}

		// Unexported fields are ignored.
		if !fval.CanInterface() {
			continue
		}

		fname, ok := sf.Tag.Lookup("spanner")
		if !ok {
			fname = sf.Name
		}

		eval, etype, err := encodeValue(fval.Interface())
		if err != nil {
			return nil, nil, err
		}

		stf = append(stf, mkField(fname, etype))
		stv = append(stv, eval)
	}

	return listProto(stv...), structType(stf...), nil
}

// Encodes a slice of Go struct values/ptrs in v to the spanner Value and Type
// protos. v itself must be non-nil.
func encodeStructArray(v interface{}) (*proto3.Value, *sppb.Type, error) {
	etyp := reflect.TypeOf(v).Elem()
	sliceval := reflect.ValueOf(v)

	// Slice of pointers to structs.
	if etyp.Kind() == reflect.Ptr {
		etyp = etyp.Elem()
	}

	// Use a dummy struct value to get the element type.
	_, elemTyp, err := encodeStruct(reflect.Zero(etyp).Interface())
	if err != nil {
		return nil, nil, err
	}

	// nil slice represents a NULL array-of-struct.
	if sliceval.IsNil() {
		return nullProto(), listType(elemTyp), nil
	}

	values := make([]*proto3.Value, 0, sliceval.Len())

	for i := 0; i < sliceval.Len(); i++ {
		ev, _, err := encodeStruct(sliceval.Index(i).Interface())
		if err != nil {
			return nil, nil, err
		}
		values = append(values, ev)
	}
	return listProto(values...), listType(elemTyp), nil
}

func isStructOrArrayOfStructValue(v interface{}) bool {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.Kind() == reflect.Struct
}

func isSupportedMutationType(v interface{}) bool {
	switch v.(type) {
	case nil, string, *string, NullString, []string, []*string, []NullString,
		[]byte, [][]byte,
		int, []int, int64, *int64, []int64, []*int64, NullInt64, []NullInt64,
		bool, *bool, []bool, []*bool, NullBool, []NullBool,
		float64, *float64, []float64, []*float64, NullFloat64, []NullFloat64,
		time.Time, *time.Time, []time.Time, []*time.Time, NullTime, []NullTime,
		civil.Date, *civil.Date, []civil.Date, []*civil.Date, NullDate, []NullDate,
		big.Rat, *big.Rat, []big.Rat, []*big.Rat, NullNumeric, []NullNumeric,
		GenericColumnValue:
		return true
	default:
		// Check if the custom type implements spanner.Encoder interface.
		if _, ok := v.(Encoder); ok {
			return true
		}

		decodableType := getDecodableSpannerType(v, false)
		return decodableType != spannerTypeUnknown && decodableType != spannerTypeInvalid
	}
}

// encodeValueArray encodes a Value array into a proto3.ListValue.
func encodeValueArray(vs []interface{}) (*proto3.ListValue, error) {
	lv := &proto3.ListValue{}
	lv.Values = make([]*proto3.Value, 0, len(vs))
	for _, v := range vs {
		if !isSupportedMutationType(v) {
			return nil, errEncoderUnsupportedType(v)
		}
		pb, _, err := encodeValue(v)
		if err != nil {
			return nil, err
		}
		lv.Values = append(lv.Values, pb)
	}
	return lv, nil
}

// encodeArray assumes that all values of the array element type encode without
// error.
func encodeArray(len int, at func(int) interface{}) (*proto3.Value, error) {
	vs := make([]*proto3.Value, len)
	var err error
	for i := 0; i < len; i++ {
		vs[i], _, err = encodeValue(at(i))
		if err != nil {
			return nil, err
		}
	}
	return listProto(vs...), nil
}

func spannerTagParser(t reflect.StructTag) (name string, keep bool, other interface{}, err error) {
	if s := t.Get("spanner"); s != "" {
		if s == "-" {
			return "", false, nil, nil
		}
		return s, true, nil, nil
	}
	return "", true, nil, nil
}

var fieldCache = fields.NewCache(spannerTagParser, nil, nil)

func trimDoubleQuotes(payload []byte) ([]byte, error) {
	if len(payload) <= 1 || payload[0] != '"' || payload[len(payload)-1] != '"' {
		return nil, fmt.Errorf("payload is too short or not wrapped with double quotes: got %q", string(payload))
	}
	// Remove the double quotes at the beginning and the end.
	return payload[1 : len(payload)-1], nil
}
