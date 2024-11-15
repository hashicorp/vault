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
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/internal/fields"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
	"google.golang.org/protobuf/reflect/protoreflect"
	proto3 "google.golang.org/protobuf/types/known/structpb"
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

// LossOfPrecisionHandlingOption describes the option to deal with loss of
// precision on numeric values.
type LossOfPrecisionHandlingOption int

const (
	// NumericRound automatically rounds a numeric value that has a higher
	// precision than what is supported by Spanner, e.g., 0.1234567895 rounds
	// to 0.123456790.
	NumericRound LossOfPrecisionHandlingOption = iota
	// NumericError returns an error for numeric values that have a higher
	// precision than what is supported by Spanner. E.g. the client returns an
	// error if the application tries to insert the value 0.1234567895.
	NumericError
)

// LossOfPrecisionHandling configures how to deal with loss of precision on
// numeric values. The value of this configuration is global and will be used
// for all Spanner clients.
var LossOfPrecisionHandling LossOfPrecisionHandlingOption

// NumericString returns a string representing a *big.Rat in a format compatible
// with Spanner SQL. It returns a floating-point literal with 9 digits after the
// decimal point.
func NumericString(r *big.Rat) string {
	return r.FloatString(NumericScaleDigits)
}

// validateNumeric returns nil if there are no errors. It will return an error
// when the numeric number is not valid.
func validateNumeric(r *big.Rat) error {
	if r == nil {
		return nil
	}
	// Add one more digit to the scale component to find out if there are more
	// digits than required.
	strRep := r.FloatString(NumericScaleDigits + 1)
	strRep = strings.TrimRight(strRep, "0")
	strRep = strings.TrimLeft(strRep, "-")
	s := strings.Split(strRep, ".")
	whole := s[0]
	scale := s[1]
	if len(scale) > NumericScaleDigits {
		return fmt.Errorf("max scale for a numeric is %d. The requested numeric has more", NumericScaleDigits)
	}
	if len(whole) > NumericPrecisionDigits-NumericScaleDigits {
		return fmt.Errorf("max precision for the whole component of a numeric is %d. The requested numeric has a whole component with precision %d", NumericPrecisionDigits-NumericScaleDigits, len(whole))
	}
	return nil
}

var (
	// CommitTimestamp is a special value used to tell Cloud Spanner to insert
	// the commit timestamp of the transaction into a column. It can be used in
	// a Mutation, or directly used in InsertStruct or InsertMap. See
	// ExampleCommitTimestamp. This is just a placeholder and the actual value
	// stored in this variable has no meaning.
	CommitTimestamp = commitTimestamp
	commitTimestamp = time.Unix(0, 0).In(time.FixedZone("CommitTimestamp placeholder", 0xDB))

	jsonUseNumber bool

	protoMsgReflectType  = reflect.TypeOf((*proto.Message)(nil)).Elem()
	protoEnumReflectType = reflect.TypeOf((*protoreflect.Enum)(nil)).Elem()

	errPayloadNil = errors.New("payload should not be nil")
)

// UseNumberWithJSONDecoderEncoder specifies whether Cloud Spanner JSON numbers are decoded
// as Number (preserving precision) or float64 (risking loss).
// Defaults to the same behavior as the standard Go library, which means decoding to float64.
// Call this method to enable lossless precision.
// NOTE 1: Calling this method affects the behavior of all clients created by this library, both existing and future instances.
// NOTE 2: This method sets a global variable that is used by the client to encode/decode JSON numbers. Access to the global variable is not synchronized. You should only call this method when there are no goroutines encoding/decoding Cloud Spanner JSON values. It is recommended to only call this method during the initialization of your application, and preferably before you create any Cloud Spanner clients, and/or in tests when there are no queries being executed.
func UseNumberWithJSONDecoderEncoder(useNumber bool) {
	jsonUseNumber = useNumber
}

func jsonUnmarshal(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	if jsonUseNumber {
		dec.UseNumber()
	}
	return dec.Decode(v)
}

// jsonIsNull returns whether v matches JSON null literal
func jsonIsNull(v []byte) bool {
	return string(v) == "null"
}

// Encoder is the interface implemented by a custom type that can be encoded to
// a supported type by Spanner. A code example:
//
//	type customField struct {
//	    Prefix string
//	    Suffix string
//	}
//
//	// Convert a customField value to a string
//	func (cf customField) EncodeSpanner() (interface{}, error) {
//	    var b bytes.Buffer
//	    b.WriteString(cf.Prefix)
//	    b.WriteString("-")
//	    b.WriteString(cf.Suffix)
//	    return b.String(), nil
//	}
type Encoder interface {
	EncodeSpanner() (interface{}, error)
}

// Decoder is the interface implemented by a custom type that can be decoded
// from a supported type by Spanner. A code example:
//
//	type customField struct {
//	    Prefix string
//	    Suffix string
//	}
//
//	// Convert a string to a customField value
//	func (cf *customField) DecodeSpanner(val interface{}) (err error) {
//	    strVal, ok := val.(string)
//	    if !ok {
//	        return fmt.Errorf("failed to decode customField: %v", val)
//	    }
//	    s := strings.Split(strVal, "-")
//	    if len(s) > 1 {
//	        cf.Prefix = s[0]
//	        cf.Suffix = s[1]
//	    }
//	    return nil
//	}
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
	return nulljson(n.Valid, n.Int64)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullInt64.
func (n *NullInt64) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
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

// Value implements the driver.Valuer interface.
func (n NullInt64) Value() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Int64, nil
}

// Scan implements the sql.Scanner interface.
func (n *NullInt64) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullInt64: %v", p)
	case *int64:
		n.Int64 = *p
	case int64:
		n.Int64 = p
	case *NullInt64:
		n.Int64 = p.Int64
		n.Valid = p.Valid
	case NullInt64:
		n.Int64 = p.Int64
		n.Valid = p.Valid
	}
	return nil
}

// GormDataType is used by gorm to determine the default data type for fields with this type.
func (n NullInt64) GormDataType() string {
	return "INT64"
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
	return nulljson(n.Valid, n.StringVal)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullString.
func (n *NullString) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
		n.StringVal = ""
		n.Valid = false
		return nil
	}
	var s *string
	if err := jsonUnmarshal(payload, &s); err != nil {
		return err
	}
	if s != nil {
		n.StringVal = *s
		n.Valid = true
	} else {
		n.StringVal = ""
		n.Valid = false
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (n NullString) Value() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.StringVal, nil
}

// Scan implements the sql.Scanner interface.
func (n *NullString) Scan(value interface{}) error {
	if value == nil {
		n.StringVal, n.Valid = "", false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullString: %v", p)
	case *string:
		n.StringVal = *p
	case string:
		n.StringVal = p
	case *NullString:
		n.StringVal = p.StringVal
		n.Valid = p.Valid
	case NullString:
		n.StringVal = p.StringVal
		n.Valid = p.Valid
	}
	return nil
}

// GormDataType is used by gorm to determine the default data type for fields with this type.
func (n NullString) GormDataType() string {
	return "STRING(MAX)"
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
	return nulljson(n.Valid, n.Float64)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullFloat64.
func (n *NullFloat64) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
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

// Value implements the driver.Valuer interface.
func (n NullFloat64) Value() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Float64, nil
}

// Scan implements the sql.Scanner interface.
func (n *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		n.Float64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullFloat64: %v", p)
	case *float64:
		n.Float64 = *p
	case float64:
		n.Float64 = p
	case *NullFloat64:
		n.Float64 = p.Float64
		n.Valid = p.Valid
	case NullFloat64:
		n.Float64 = p.Float64
		n.Valid = p.Valid
	}
	return nil
}

// GormDataType is used by gorm to determine the default data type for fields with this type.
func (n NullFloat64) GormDataType() string {
	return "FLOAT64"
}

// NullFloat32 represents a Cloud Spanner FLOAT32 that may be NULL.
type NullFloat32 struct {
	Float32 float32 // Float32 contains the value when it is non-NULL, and zero when NULL.
	Valid   bool    // Valid is true if FLOAT32 is not NULL.
}

// IsNull implements NullableValue.IsNull for NullFloat32.
func (n NullFloat32) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullFloat32
func (n NullFloat32) String() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", n.Float32)
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullFloat32.
func (n NullFloat32) MarshalJSON() ([]byte, error) {
	return nulljson(n.Valid, n.Float32)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullFloat32.
func (n *NullFloat32) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
		n.Float32 = float32(0)
		n.Valid = false
		return nil
	}
	num, err := strconv.ParseFloat(string(payload), 32)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to float32: got %v", string(payload))
	}
	n.Float32 = float32(num)
	n.Valid = true
	return nil
}

// Value implements the driver.Valuer interface.
func (n NullFloat32) Value() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Float32, nil
}

// Scan implements the sql.Scanner interface.
func (n *NullFloat32) Scan(value interface{}) error {
	if value == nil {
		n.Float32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullFloat32: %v", p)
	case *float32:
		n.Float32 = *p
	case float32:
		n.Float32 = p
	case *NullFloat32:
		n.Float32 = p.Float32
		n.Valid = p.Valid
	case NullFloat32:
		n.Float32 = p.Float32
		n.Valid = p.Valid
	}
	return nil
}

// GormDataType is used by gorm to determine the default data type for fields with this type.
func (n NullFloat32) GormDataType() string {
	return "FLOAT32"
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
	return nulljson(n.Valid, n.Bool)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullBool.
func (n *NullBool) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
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

// Value implements the driver.Valuer interface.
func (n NullBool) Value() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Bool, nil
}

// Scan implements the sql.Scanner interface.
func (n *NullBool) Scan(value interface{}) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullBool: %v", p)
	case *bool:
		n.Bool = *p
	case bool:
		n.Bool = p
	case *NullBool:
		n.Bool = p.Bool
		n.Valid = p.Valid
	case NullBool:
		n.Bool = p.Bool
		n.Valid = p.Valid
	}
	return nil
}

// GormDataType is used by gorm to determine the default data type for fields with this type.
func (n NullBool) GormDataType() string {
	return "BOOL"
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
	return nulljson(n.Valid, n.Time)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullTime.
func (n *NullTime) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
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

// Value implements the driver.Valuer interface.
func (n NullTime) Value() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Time, nil
}

// Scan implements the sql.Scanner interface.
func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		n.Time, n.Valid = time.Time{}, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullTime: %v", p)
	case *time.Time:
		n.Time = *p
	case time.Time:
		n.Time = p
	case *NullTime:
		n.Time = p.Time
		n.Valid = p.Valid
	case NullTime:
		n.Time = p.Time
		n.Valid = p.Valid
	}
	return nil
}

// GormDataType is used by gorm to determine the default data type for fields with this type.
func (n NullTime) GormDataType() string {
	return "TIMESTAMP"
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
	return nulljson(n.Valid, n.Date)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullDate.
func (n *NullDate) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
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

// Value implements the driver.Valuer interface.
func (n NullDate) Value() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Date, nil
}

// Scan implements the sql.Scanner interface.
func (n *NullDate) Scan(value interface{}) error {
	if value == nil {
		n.Date, n.Valid = civil.Date{}, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullDate: %v", p)
	case *civil.Date:
		n.Date = *p
	case civil.Date:
		n.Date = p
	case *NullDate:
		n.Date = p.Date
		n.Valid = p.Valid
	case NullDate:
		n.Date = p.Date
		n.Valid = p.Valid
	}
	return nil
}

// GormDataType is used by gorm to determine the default data type for fields with this type.
func (n NullDate) GormDataType() string {
	return "DATE"
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
	return nulljson(n.Valid, NumericString(&n.Numeric))
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullNumeric.
func (n *NullNumeric) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
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

// Value implements the driver.Valuer interface.
func (n NullNumeric) Value() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Numeric, nil
}

// Scan implements the sql.Scanner interface.
func (n *NullNumeric) Scan(value interface{}) error {
	if value == nil {
		n.Numeric, n.Valid = big.Rat{}, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullNumeric: %v", p)
	case *big.Rat:
		n.Numeric = *p
	case big.Rat:
		n.Numeric = p
	case *NullNumeric:
		n.Numeric = p.Numeric
		n.Valid = p.Valid
	case NullNumeric:
		n.Numeric = p.Numeric
		n.Valid = p.Valid
	}
	return nil
}

// GormDataType is used by gorm to determine the default data type for fields with this type.
func (n NullNumeric) GormDataType() string {
	return "NUMERIC"
}

// NullJSON represents a Cloud Spanner JSON that may be NULL.
//
// This type must always be used when encoding values to a JSON column in Cloud
// Spanner.
//
// NullJSON does not implement the driver.Valuer and sql.Scanner interfaces, as
// the underlying value can be anything. This means that the type NullJSON must
// also be used when calling sql.Row#Scan(dest ...interface{}) for a JSON
// column.
type NullJSON struct {
	Value interface{} // Val contains the value when it is non-NULL, and nil when NULL.
	Valid bool        // Valid is true if Json is not NULL.
}

// IsNull implements NullableValue.IsNull for NullJSON.
func (n NullJSON) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullJSON.
func (n NullJSON) String() string {
	if !n.Valid {
		return nullString
	}
	b, err := json.Marshal(n.Value)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%v", string(b))
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullJSON.
func (n NullJSON) MarshalJSON() ([]byte, error) {
	return nulljson(n.Valid, n.Value)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullJSON.
func (n *NullJSON) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
		n.Valid = false
		return nil
	}
	var v interface{}
	err := jsonUnmarshal(payload, &v)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to a struct: got %v, err: %w", string(payload), err)
	}
	n.Value = v
	n.Valid = true
	return nil
}

// GormDataType is used by gorm to determine the default data type for fields with this type.
func (n NullJSON) GormDataType() string {
	return "JSON"
}

// PGNumeric represents a Cloud Spanner PG Numeric that may be NULL.
type PGNumeric struct {
	Numeric string // Numeric contains the value when it is non-NULL, and an empty string when NULL.
	Valid   bool   // Valid is true if Numeric is not NULL.
}

// IsNull implements NullableValue.IsNull for PGNumeric.
func (n PGNumeric) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for PGNumeric
func (n PGNumeric) String() string {
	if !n.Valid {
		return nullString
	}
	return n.Numeric
}

// MarshalJSON implements json.Marshaler.MarshalJSON for PGNumeric.
func (n PGNumeric) MarshalJSON() ([]byte, error) {
	return nulljson(n.Valid, n.Numeric)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for PGNumeric.
func (n *PGNumeric) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
		n.Numeric = ""
		n.Valid = false
		return nil
	}
	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	n.Numeric = string(payload)
	n.Valid = true
	return nil
}

// NullProtoMessage represents a Cloud Spanner PROTO that may be NULL.
// To write a NULL value using NullProtoMessage set ProtoMessageVal to typed nil and set Valid to true.
type NullProtoMessage struct {
	ProtoMessageVal proto.Message // ProtoMessageVal contains the value when Valid is true, and nil when NULL.
	Valid           bool          // Valid is true if ProtoMessageVal is not NULL.
}

// IsNull implements NullableValue.IsNull for NullProtoMessage.
func (n NullProtoMessage) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullProtoMessage.
func (n NullProtoMessage) String() string {
	if !n.Valid {
		return nullString
	}
	return protoadapt.MessageV1Of(n.ProtoMessageVal).String()
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullProtoMessage.
func (n NullProtoMessage) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.ProtoMessageVal)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullProtoMessage.
func (n *NullProtoMessage) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
		n.ProtoMessageVal = nil
		n.Valid = false
		return nil
	}
	err := jsonUnmarshal(payload, n.ProtoMessageVal)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to a proto message: err: %s", err)
	}
	n.Valid = true
	return nil
}

// NullProtoEnum represents a Cloud Spanner ENUM that may be NULL.
// To write a NULL value using NullProtoEnum set ProtoEnumVal to typed nil and set Valid to true.
type NullProtoEnum struct {
	ProtoEnumVal protoreflect.Enum // ProtoEnumVal contains the value when Valid is true, and nil when NULL.
	Valid        bool              // Valid is true if ProtoEnumVal is not NULL.
}

// IsNull implements NullableValue.IsNull for NullProtoEnum.
func (n NullProtoEnum) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for NullProtoEnum.
func (n NullProtoEnum) String() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", n.ProtoEnumVal)
}

// MarshalJSON implements json.Marshaler.MarshalJSON for NullProtoEnum.
func (n NullProtoEnum) MarshalJSON() ([]byte, error) {
	if n.Valid && n.ProtoEnumVal != nil {
		return []byte(fmt.Sprintf("%v", n.ProtoEnumVal.Number())), nil
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for NullProtoEnum.
func (n *NullProtoEnum) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
		n.ProtoEnumVal = nil
		n.Valid = false
		return nil
	}
	if reflect.ValueOf(n.ProtoEnumVal).Kind() != reflect.Ptr {
		return errNotAPointerField(n, n.ProtoEnumVal)
	}
	num, err := strconv.ParseInt(string(payload), 10, 64)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to Enum: got %v", string(payload))
	}
	reflect.ValueOf(n.ProtoEnumVal).Elem().SetInt(num)
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

// PGJsonB represents a Cloud Spanner PGJsonB that may be NULL.
type PGJsonB struct {
	Value interface{} // Val contains the value when it is non-NULL, and nil when NULL.
	Valid bool        // Valid is true if PGJsonB is not NULL.
	// This is here to support customer wrappers around PGJsonB type, this will help during getDecodableSpannerType
	// to differentiate between PGJsonB and NullJSON types.
	_ bool
}

// IsNull implements NullableValue.IsNull for PGJsonB.
func (n PGJsonB) IsNull() bool {
	return !n.Valid
}

// String implements Stringer.String for PGJsonB.
func (n PGJsonB) String() string {
	if !n.Valid {
		return nullString
	}
	b, err := json.Marshal(n.Value)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%v", string(b))
}

// MarshalJSON implements json.Marshaler.MarshalJSON for PGJsonB.
func (n PGJsonB) MarshalJSON() ([]byte, error) {
	return nulljson(n.Valid, n.Value)
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON for PGJsonB.
func (n *PGJsonB) UnmarshalJSON(payload []byte) error {
	if payload == nil {
		return errPayloadNil
	}
	if jsonIsNull(payload) {
		n.Valid = false
		return nil
	}
	var v interface{}
	err := jsonUnmarshal(payload, &v)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to a struct: got %v, err: %w", string(payload), err)
	}
	n.Value = v
	n.Valid = true
	return nil
}

func nulljson(valid bool, v interface{}) ([]byte, error) {
	if !valid {
		return []byte("null"), nil
	}
	return json.Marshal(v)
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

// errNilDstField returns error for decoding into nil interface{} of Value field in NullProtoMessage or NullProtoEnum.
func errNilDstField(dst interface{}, field string) error {
	return spannerErrorf(codes.InvalidArgument, "field %s in %T cannot be nil", field, dst)
}

// errNilArrElemType returns error for input Cloud Spanner data type being a array but without a
// non-nil array element type.
func errNilArrElemType(t *sppb.Type) error {
	return spannerErrorf(codes.FailedPrecondition, "array type %v is with nil array element type", t)
}

// errNotValidSrc returns error if Valid field is false for NullProtoMessage and NullProtoEnum
func errNotValidSrc(dst interface{}) error {
	return spannerErrorf(codes.InvalidArgument, "field \"Valid\" of %T cannot be set to false when writing data to Cloud Spanner. Use typed nil in %T to write null values to Cloud Spanner", dst, dst)
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

// errNotAPointer returns error for decoding a non pointer type.
func errNotAPointer(dst interface{}) error {
	return spannerErrorf(codes.InvalidArgument, "destination %T must be a pointer", dst)
}

// errNotAPointerField returns error for decoding a non pointer type.
func errNotAPointerField(dst interface{}, dstField interface{}) error {
	return spannerErrorf(codes.InvalidArgument, "destination %T in %T must be a pointer", dstField, dst)
}

func errNilNotAllowed(dst interface{}, name string) error {
	return spannerErrorf(codes.InvalidArgument, "destination %T does not support Null values. Use %s, an array with pointer type elements to read Null values", dst, name)
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
func decodeValue(v *proto3.Value, t *sppb.Type, ptr interface{}, opts ...DecodeOptions) error {
	if v == nil {
		return errNilSrc()
	}
	if t == nil {
		return errNilSpannerType()
	}
	code := t.Code
	typeAnnotation := t.TypeAnnotation
	acode := sppb.TypeCode_TYPE_CODE_UNSPECIFIED
	atypeAnnotation := sppb.TypeAnnotationCode_TYPE_ANNOTATION_CODE_UNSPECIFIED
	if code == sppb.TypeCode_ARRAY {
		if t.ArrayElementType == nil {
			return errNilArrElemType(t)
		}
		acode = t.ArrayElementType.Code
		atypeAnnotation = t.ArrayElementType.TypeAnnotation
	}

	if code == sppb.TypeCode_PROTO && reflect.TypeOf(ptr).Elem().Kind() == reflect.Ptr {
		pve := reflect.ValueOf(ptr).Elem()
		if pve.IsNil() {
			pve.Set(reflect.New(pve.Type().Elem()))
		}
		ptr = pve.Interface()
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
	case *NullString, **string, *sql.NullString:
		// Most Null* types are automatically supported for both spanner.Null* and sql.Null* types, except for
		// NullString, and we need to add explicit support for it here. The reason that the other types are
		// automatically supported is that they use the same field names (e.g. spanner.NullBool and sql.NullBool both
		// contain the fields Valid and Bool). spanner.NullString has a field StringVal, sql.NullString has a field
		// String.
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
			case *sql.NullString:
				*sp = sql.NullString{}
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
		case *sql.NullString:
			sp.Valid = true
			sp.String = x
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
		if code != sppb.TypeCode_BYTES && code != sppb.TypeCode_PROTO {
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
		if acode != sppb.TypeCode_BYTES && acode != sppb.TypeCode_PROTO {
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
		if code != sppb.TypeCode_INT64 && code != sppb.TypeCode_ENUM {
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
		if code != sppb.TypeCode_INT64 && code != sppb.TypeCode_ENUM {
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
		if acode != sppb.TypeCode_INT64 && acode != sppb.TypeCode_ENUM {
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
		if acode != sppb.TypeCode_INT64 && acode != sppb.TypeCode_ENUM {
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
	case *float32:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_FLOAT32 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getFloat32Value(v)
		if err != nil {
			return err
		}
		*p = x
	case *NullFloat32, **float32:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_FLOAT32 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullFloat32:
				*sp = NullFloat32{}
			case **float32:
				*sp = nil
			}
			break
		}
		x, err := getFloat32Value(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *NullFloat32:
			sp.Valid = true
			sp.Float32 = x
		case **float32:
			*sp = &x
		}
	case *[]NullFloat32, *[]*float32:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_FLOAT32 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullFloat32:
				*sp = nil
			case *[]*float32:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullFloat32:
			y, err := decodeNullFloat32Array(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*float32:
			y, err := decodeFloat32PointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]float32:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_FLOAT32 {
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
		y, err := decodeFloat32Array(x)
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
	case *NullJSON:
		if p == nil {
			return errNilDst(p)
		}
		if code == sppb.TypeCode_ARRAY {
			if acode != sppb.TypeCode_JSON {
				return errTypeMismatch(code, acode, ptr)
			}
			x, err := getListValue(v)
			if err != nil {
				return err
			}
			y, err := decodeNullJSONArrayToNullJSON(x)
			if err != nil {
				return err
			}
			*p = *y
		} else {
			if code != sppb.TypeCode_JSON {
				return errTypeMismatch(code, acode, ptr)
			}
			if isNull {
				*p = NullJSON{}
				break
			}
			x := v.GetStringValue()
			var y interface{}
			err := jsonUnmarshal([]byte(x), &y)
			if err != nil {
				return err
			}
			*p = NullJSON{y, true}
		}
	case *[]NullJSON:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_JSON {
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
		y, err := decodeNullJSONArray(x)
		if err != nil {
			return err
		}
		*p = y
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
	case *PGNumeric:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_NUMERIC || typeAnnotation != sppb.TypeAnnotationCode_PG_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = PGNumeric{}
			break
		}
		*p = PGNumeric{v.GetStringValue(), true}
	case *[]PGNumeric:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_NUMERIC || atypeAnnotation != sppb.TypeAnnotationCode_PG_NUMERIC {
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
		y, err := decodePGNumericArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *PGJsonB:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_JSON || typeAnnotation != sppb.TypeAnnotationCode_PG_JSONB {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = PGJsonB{}
			break
		}
		x := v.GetStringValue()
		var y interface{}
		err := jsonUnmarshal([]byte(x), &y)
		if err != nil {
			return err
		}
		*p = PGJsonB{Value: y, Valid: true}
	case *[]PGJsonB:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_JSON || typeAnnotation != sppb.TypeAnnotationCode_PG_JSONB {
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
		y, err := decodePGJsonBArray(x)
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
	case protoreflect.Enum:
		if p == nil {
			return errNilDst(p)
		}
		if reflect.ValueOf(p).Kind() != reflect.Ptr {
			return errNotAPointer(p)
		}
		if code != sppb.TypeCode_ENUM && code != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		y, err := getIntegerFromStringValue(v)
		if err != nil {
			return err
		}
		reflect.ValueOf(p).Elem().SetInt(y)
	case *NullProtoEnum:
		if p == nil {
			return errNilDst(p)
		}
		if p.ProtoEnumVal == nil {
			return errNilDstField(p, "ProtoEnumVal")
		}
		if reflect.ValueOf(p.ProtoEnumVal).Kind() != reflect.Ptr {
			return errNotAPointer(p)
		}
		if code != sppb.TypeCode_ENUM && code != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = NullProtoEnum{}
			break
		}
		y, err := getIntegerFromStringValue(v)
		if err != nil {
			return err
		}
		reflect.ValueOf(p.ProtoEnumVal).Elem().SetInt(y)
		p.Valid = true
	case proto.Message:
		// Check if the pointer is a custom type that implements spanner.Decoder
		// interface.
		if decodedVal, ok := ptr.(Decoder); ok {
			x, err := getGenericValue(t, v)
			if err != nil {
				return err
			}
			return decodedVal.DecodeSpanner(x)
		}
		if p == nil {
			return errNilDst(p)
		}
		if reflect.ValueOf(p).Kind() != reflect.Ptr {
			return errNotAPointer(p)
		}
		if code != sppb.TypeCode_PROTO && code != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		y, err := getBytesFromStringValue(v)
		if err != nil {
			return err
		}
		err = proto.Unmarshal(y, p)
		if err != nil {
			return err
		}
	case *NullProtoMessage:
		if p == nil {
			return errNilDst(p)
		}
		if p.ProtoMessageVal == nil {
			return errNilDstField(p, "ProtoMessageVal")
		}
		if reflect.ValueOf(p.ProtoMessageVal).Kind() != reflect.Ptr {
			return errNotAPointer(p.ProtoMessageVal)
		}
		if code != sppb.TypeCode_PROTO && code != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = NullProtoMessage{}
			break
		}
		y, err := getBytesFromStringValue(v)
		if err != nil {
			return err
		}
		err = proto.Unmarshal(y, p.ProtoMessageVal)
		if err != nil {
			return err
		}
		p.Valid = true
	default:
		// Check if the pointer is a custom type that implements spanner.Decoder
		// interface.
		if decodedVal, ok := ptr.(Decoder); ok {
			x, err := getGenericValue(t, v)
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
			return decodableType.decodeValueToCustomType(v, t, acode, atypeAnnotation, ptr)
		}

		rv := reflect.ValueOf(ptr)
		typ := rv.Type()
		// Check if the interface{} is a pointer and is of type array of proto columns
		if typ.Kind() == reflect.Ptr && isAnArrayOfProtoColumn(ptr) && code == sppb.TypeCode_ARRAY {
			if isNull {
				rv.Elem().Set(reflect.Zero(rv.Elem().Type()))
				break
			}
			// Get the user-defined type of the proto array
			etyp := typ.Elem().Elem()
			switch acode {
			case sppb.TypeCode_PROTO, sppb.TypeCode_BYTES:
				if etyp.Implements(protoMsgReflectType) {
					if etyp.Kind() == reflect.Ptr {
						x, err := getListValue(v)
						if err != nil {
							return err
						}
						return decodeProtoMessagePtrArray(x, t.ArrayElementType, rv)
					}
					return errTypeMismatch(code, acode, ptr)
				}
			case sppb.TypeCode_ENUM, sppb.TypeCode_INT64:
				if etyp.Implements(protoEnumReflectType) {
					x, err := getListValue(v)
					if err != nil {
						return err
					}
					if etyp.Kind() == reflect.Ptr {
						return decodeProtoEnumPtrArray(x, t.ArrayElementType, rv)
					}
					return decodeProtoEnumArray(x, t.ArrayElementType, rv, ptr)
				}
			}
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
			// The container is not a slice of struct pointers.
			return fmt.Errorf("the container is not a slice of struct pointers: %v", errTypeMismatch(code, acode, ptr))
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
		s := decodeSetting{
			Lenient: false,
		}
		for _, opt := range opts {
			opt.Apply(&s)
		}
		if err = decodeStructArray(t.ArrayElementType.StructType, x, p, s.Lenient); err != nil {
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
	spannerTypeNonNullFloat32
	spannerTypeNonNullNumeric
	spannerTypeNonNullTime
	spannerTypeNonNullDate
	spannerTypeNullString
	spannerTypeNullInt64
	spannerTypeNullBool
	spannerTypeNullFloat64
	spannerTypeNullFloat32
	spannerTypeNullTime
	spannerTypeNullDate
	spannerTypeNullNumeric
	spannerTypeNullJSON
	spannerTypePGNumeric
	spannerTypePGJsonB
	spannerTypeArrayOfNonNullString
	spannerTypeArrayOfByteArray
	spannerTypeArrayOfNonNullInt64
	spannerTypeArrayOfNonNullBool
	spannerTypeArrayOfNonNullFloat64
	spannerTypeArrayOfNonNullFloat32
	spannerTypeArrayOfNonNullNumeric
	spannerTypeArrayOfNonNullTime
	spannerTypeArrayOfNonNullDate
	spannerTypeArrayOfNullString
	spannerTypeArrayOfNullInt64
	spannerTypeArrayOfNullBool
	spannerTypeArrayOfNullFloat64
	spannerTypeArrayOfNullFloat32
	spannerTypeArrayOfNullNumeric
	spannerTypeArrayOfNullJSON
	spannerTypeArrayOfNullTime
	spannerTypeArrayOfNullDate
	spannerTypeArrayOfPGNumeric
	spannerTypeArrayOfPGJsonB
)

// supportsNull returns true for the Go types that can hold a null value from
// Spanner.
func (d decodableSpannerType) supportsNull() bool {
	switch d {
	case spannerTypeNonNullString, spannerTypeNonNullInt64, spannerTypeNonNullBool, spannerTypeNonNullFloat64, spannerTypeNonNullFloat32, spannerTypeNonNullTime, spannerTypeNonNullDate, spannerTypeNonNullNumeric:
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
var typeOfNullFloat32 = reflect.TypeOf(NullFloat32{})
var typeOfNullTime = reflect.TypeOf(NullTime{})
var typeOfNullDate = reflect.TypeOf(NullDate{})
var typeOfNullNumeric = reflect.TypeOf(NullNumeric{})
var typeOfNullJSON = reflect.TypeOf(NullJSON{})
var typeOfPGNumeric = reflect.TypeOf(PGNumeric{})
var typeOfPGJsonB = reflect.TypeOf(PGJsonB{})

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
	case reflect.Float32:
		return spannerTypeNonNullFloat32
	case reflect.Float64:
		return spannerTypeNonNullFloat64
	case reflect.Ptr:
		t := val.Type()
		if t.ConvertibleTo(typeOfNullNumeric) {
			return spannerTypeNullNumeric
		}
		if t.ConvertibleTo(typeOfNullJSON) {
			return spannerTypeNullJSON
		}
		if t.ConvertibleTo(typeOfPGJsonB) {
			return spannerTypePGJsonB
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
		if t.ConvertibleTo(typeOfNullFloat32) {
			return spannerTypeNullFloat32
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
		if t.ConvertibleTo(typeOfNullJSON) {
			return spannerTypeNullJSON
		}
		if t.ConvertibleTo(typeOfPGNumeric) {
			return spannerTypePGNumeric
		}
		if t.ConvertibleTo(typeOfPGJsonB) {
			return spannerTypePGJsonB
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
		case reflect.Float32:
			return spannerTypeArrayOfNonNullFloat32
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
			if t.ConvertibleTo(typeOfNullFloat32) {
				return spannerTypeArrayOfNullFloat32
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
			if t.ConvertibleTo(typeOfNullJSON) {
				return spannerTypeArrayOfNullJSON
			}
			if t.ConvertibleTo(typeOfPGNumeric) {
				return spannerTypeArrayOfPGNumeric
			}
			if t.ConvertibleTo(typeOfPGJsonB) {
				return spannerTypeArrayOfPGJsonB
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
func (dsc decodableSpannerType) decodeValueToCustomType(v *proto3.Value, t *sppb.Type, acode sppb.TypeCode, atypeAnnotation sppb.TypeAnnotationCode, ptr interface{}) error {
	code := t.Code
	typeAnnotation := t.TypeAnnotation
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
		if code != sppb.TypeCode_BYTES && code != sppb.TypeCode_PROTO {
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
		if code != sppb.TypeCode_INT64 && code != sppb.TypeCode_ENUM {
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
	case spannerTypeNonNullFloat32, spannerTypeNullFloat32:
		if code != sppb.TypeCode_FLOAT32 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullFloat32{}
			break
		}
		x, err := getFloat32Value(v)
		if err != nil {
			return err
		}
		if dsc == spannerTypeNonNullFloat32 {
			result = &x
		} else {
			result = &NullFloat32{x, !isNull}
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
	case spannerTypePGNumeric:
		if code != sppb.TypeCode_NUMERIC || typeAnnotation != sppb.TypeAnnotationCode_PG_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &PGNumeric{}
			break
		}
		result = &PGNumeric{v.GetStringValue(), true}
	case spannerTypeNullJSON:
		if code != sppb.TypeCode_JSON {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullJSON{}
			break
		}
		x := v.GetStringValue()
		var y interface{}
		err := jsonUnmarshal([]byte(x), &y)
		if err != nil {
			return err
		}
		result = &NullJSON{y, true}
	case spannerTypePGJsonB:
		if code != sppb.TypeCode_JSON || typeAnnotation != sppb.TypeAnnotationCode_PG_JSONB {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &PGJsonB{}
			break
		}
		x := v.GetStringValue()
		var y interface{}
		err := jsonUnmarshal([]byte(x), &y)
		if err != nil {
			return err
		}
		result = &PGJsonB{Value: y, Valid: true}
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
		if acode != sppb.TypeCode_BYTES && acode != sppb.TypeCode_PROTO {
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
		if acode != sppb.TypeCode_INT64 && acode != sppb.TypeCode_ENUM {
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
	case spannerTypeArrayOfNonNullFloat32, spannerTypeArrayOfNullFloat32:
		if acode != sppb.TypeCode_FLOAT32 {
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
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, float32Type(), "FLOAT32")
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
	case spannerTypeArrayOfPGNumeric:
		if acode != sppb.TypeCode_NUMERIC || atypeAnnotation != sppb.TypeAnnotationCode_PG_NUMERIC {
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
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, pgNumericType(), "PGNUMERIC")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNullJSON:
		if acode != sppb.TypeCode_JSON {
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
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, jsonType(), "JSON")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfPGJsonB:
		if acode != sppb.TypeCode_JSON || atypeAnnotation != sppb.TypeAnnotationCode_PG_JSONB {
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
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, pgJsonbType(), "PGJSONB")
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

// getIntegerFromStringValue returns the integer value of the string value encoded in proto3.Value v
func getIntegerFromStringValue(v *proto3.Value) (int64, error) {
	x, err := getStringValue(v)
	if err != nil {
		return 0, err
	}
	y, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return 0, errBadEncoding(v, err)
	}
	return y, nil
}

// getBytesFromStringValue returns the bytes value of the string value encoded in proto3.Value v
func getBytesFromStringValue(v *proto3.Value) ([]byte, error) {
	x, err := getStringValue(v)
	if err != nil {
		return nil, err
	}
	y, err := base64.StdEncoding.DecodeString(x)
	if err != nil {
		return nil, errBadEncoding(v, err)
	}
	return y, nil
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
func getGenericValue(t *sppb.Type, v *proto3.Value) (interface{}, error) {
	switch x := v.GetKind().(type) {
	case *proto3.Value_NumberValue:
		return x.NumberValue, nil
	case *proto3.Value_BoolValue:
		return x.BoolValue, nil
	case *proto3.Value_StringValue:
		return x.StringValue, nil
	case *proto3.Value_ListValue:
		return x.ListValue, nil
	case *proto3.Value_NullValue:
		return getTypedNil(t)
	default:
		return 0, errSrcVal(v, "Number, Bool, String, List")
	}
}

func getTypedNil(t *sppb.Type) (interface{}, error) {
	switch t.Code {
	case sppb.TypeCode_FLOAT64:
		var f *float64
		return f, nil
	case sppb.TypeCode_BOOL:
		var b *bool
		return b, nil
	default:
		// The encoding for most types is string, except for the ones listed
		// above.
		var s *string
		return s, nil
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

// errUnexpectedFloat32Str returns error for decoder getting an unexpected
// string for representing special float values.
func errUnexpectedFloat32Str(s string) error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected string value %q for float32 number", s)
}

// getFloat32Value returns the float32 value encoded in proto3.Value v whose
// kind is proto3.Value_NumberValue / proto3.Value_StringValue.
// Cloud Spanner uses string to encode NaN, Infinity and -Infinity.
func getFloat32Value(v *proto3.Value) (float32, error) {
	switch x := v.GetKind().(type) {
	case *proto3.Value_NumberValue:
		if x == nil {
			break
		}
		return float32(x.NumberValue), nil
	case *proto3.Value_StringValue:
		if x == nil {
			break
		}
		switch x.StringValue {
		case "NaN":
			return float32(math.NaN()), nil
		case "Infinity":
			return float32(math.Inf(1)), nil
		case "-Infinity":
			return float32(math.Inf(-1)), nil
		default:
			return 0, errUnexpectedFloat32Str(x.StringValue)
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
	if !errors.As(err, &se) {
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

// decodeNullFloat32Array decodes proto3.ListValue pb into a NullFloat32 slice.
func decodeNullFloat32Array(pb *proto3.ListValue) ([]NullFloat32, error) {
	if pb == nil {
		return nil, errNilListValue("FLOAT32")
	}
	a := make([]NullFloat32, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, float32Type(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "FLOAT32", err)
		}
	}
	return a, nil
}

// decodeFloat32PointerArray decodes proto3.ListValue pb into a *float32 slice.
func decodeFloat32PointerArray(pb *proto3.ListValue) ([]*float32, error) {
	if pb == nil {
		return nil, errNilListValue("FLOAT32")
	}
	a := make([]*float32, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, float32Type(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "FLOAT32", err)
		}
	}
	return a, nil
}

// decodeFloat32Array decodes proto3.ListValue pb into a float32 slice.
func decodeFloat32Array(pb *proto3.ListValue) ([]float32, error) {
	if pb == nil {
		return nil, errNilListValue("FLOAT32")
	}
	a := make([]float32, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, float32Type(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "FLOAT32", err)
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

// decodeNullJSONArray decodes proto3.ListValue pb into a NullJSON slice.
func decodeNullJSONArray(pb *proto3.ListValue) ([]NullJSON, error) {
	if pb == nil {
		return nil, errNilListValue("JSON")
	}
	a := make([]NullJSON, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, jsonType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "JSON", err)
		}
	}
	return a, nil
}

// decodeJsonBArray decodes proto3.ListValue pb into a JsonB slice.
func decodePGJsonBArray(pb *proto3.ListValue) ([]PGJsonB, error) {
	if pb == nil {
		return nil, errNilListValue("PGJSONB")
	}
	a := make([]PGJsonB, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, pgJsonbType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "PGJSONB", err)
		}
	}
	return a, nil
}

// decodeNullJSONArray decodes proto3.ListValue pb into a NullJSON pointer.
func decodeNullJSONArrayToNullJSON(pb *proto3.ListValue) (*NullJSON, error) {
	if pb == nil {
		return nil, errNilListValue("JSON")
	}
	strs := []string{}
	for _, v := range pb.Values {
		if _, ok := v.Kind.(*proto3.Value_NullValue); ok {
			strs = append(strs, "null")
		} else {
			strs = append(strs, v.GetStringValue())
		}
	}
	s := fmt.Sprintf("[%s]", strings.Join(strs, ","))
	var y interface{}
	err := jsonUnmarshal([]byte(s), &y)
	if err != nil {
		return nil, err
	}
	return &NullJSON{y, true}, nil
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

// decodePGNumericArray decodes proto3.ListValue pb into a PGNumeric slice.
func decodePGNumericArray(pb *proto3.ListValue) ([]PGNumeric, error) {
	if pb == nil {
		return nil, errNilListValue("PGNUMERIC")
	}
	a := make([]PGNumeric, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, pgNumericType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "PGNUMERIC", err)
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

// decodeProtoMessagePtrArray decodes proto3.ListValue pb into a *proto.Message slice.
// The elements in the array implements proto.Message interface only if the element is a pointer (e.g. *ProtoMessage).
// However, if the element is a value (e.g. ProtoMessage), then it does not implement proto.Message.
// Therefore, decodeProtoMessagePtrArray allows decoding of proto message array if the array element is a pointer only.
func decodeProtoMessagePtrArray(pb *proto3.ListValue, t *sppb.Type, rv reflect.Value) error {
	if pb == nil {
		return errNilListValue("PROTO")
	}
	etyp := rv.Type().Elem().Elem().Elem()
	a := reflect.MakeSlice(rv.Type().Elem(), len(pb.Values), len(pb.Values))
	for i, v := range pb.Values {
		_, isNull := v.Kind.(*proto3.Value_NullValue)
		if isNull {
			continue
		}
		msg := reflect.New(etyp).Interface().(proto.Message)
		if err := decodeValue(v, t, msg); err != nil {
			return errDecodeArrayElement(i, v, "PROTO", err)
		}
		a.Index(i).Set(reflect.ValueOf(msg))
	}
	rv.Elem().Set(a)
	return nil
}

// decodeProtoEnumPtrArray decodes proto3.ListValue pb into a *protoreflect.Enum slice.
func decodeProtoEnumPtrArray(pb *proto3.ListValue, t *sppb.Type, rv reflect.Value) error {
	if pb == nil {
		return errNilListValue("ENUM")
	}
	etyp := rv.Type().Elem().Elem().Elem()
	a := reflect.MakeSlice(rv.Type().Elem(), len(pb.Values), len(pb.Values))
	for i, v := range pb.Values {
		_, isNull := v.Kind.(*proto3.Value_NullValue)
		if isNull {
			continue
		}
		enum := reflect.New(etyp).Interface().(protoreflect.Enum)
		if err := decodeValue(v, t, enum); err != nil {
			return errDecodeArrayElement(i, v, "ENUM", err)
		}
		a.Index(i).Set(reflect.ValueOf(enum))
	}
	rv.Elem().Set(a)
	return nil
}

// decodeProtoEnumArray decodes proto3.ListValue pb into a protoreflect.Enum slice.
func decodeProtoEnumArray(pb *proto3.ListValue, t *sppb.Type, rv reflect.Value, ptr interface{}) error {
	if pb == nil {
		return errNilListValue("ENUM")
	}
	a := reflect.MakeSlice(rv.Type().Elem(), len(pb.Values), len(pb.Values))
	// decodeValue method can decode only if ENUM is a pointer type.
	// As the ENUM element in the Array is not a pointer type we cannot use decodeValue method
	// and hence handle it separately.
	for i, v := range pb.Values {
		_, isNull := v.Kind.(*proto3.Value_NullValue)
		// As the ENUM elements in the array are value type and not pointer type,
		// we cannot support NULL values in the array
		if isNull {
			return errNilNotAllowed(ptr, "*[]*protoreflect.Enum")
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return errBadEncoding(v, err)
		}
		a.Index(i).SetInt(y)
	}
	rv.Elem().Set(a)
	return nil
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

// errDupGoField returns error for duplicated Go STRUCT field names
func errDupGoField(s interface{}, name string) error {
	return spannerErrorf(codes.InvalidArgument, "Go struct %+v(type %T) has duplicate fields for GO STRUCT field %s", s, s, name)
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
	if !errors.As(err, &se) {
		return spannerErrorf(codes.Unknown,
			"cannot decode field %v of Cloud Spanner STRUCT %+v, error = <%v>", f, ty, err)
	}
	se.decorate(fmt.Sprintf("cannot decode field %v of Cloud Spanner STRUCT %+v", f, ty))
	return se
}

// decodeSetting contains all the settings for decoding from spanner struct
type decodeSetting struct {
	Lenient bool
}

// DecodeOptions is the interface to change decode struct settings
type DecodeOptions interface {
	Apply(s *decodeSetting)
}

type withLenient struct{ lenient bool }

func (w withLenient) Apply(s *decodeSetting) {
	s.Lenient = w.lenient
}

// WithLenient returns a DecodeOptions that allows decoding into a struct with missing fields in database.
func WithLenient() DecodeOptions {
	return withLenient{lenient: true}
}

// decodeStruct decodes proto3.ListValue pb into struct referenced by pointer
// ptr, according to
// the structural information given in sppb.StructType ty.
func decodeStruct(ty *sppb.StructType, pb *proto3.ListValue, ptr interface{}, lenient bool) error {
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
	// return error if lenient is true and destination has duplicate exported columns
	if lenient {
		fieldNames := getAllFieldNames(v)
		for _, f := range fieldNames {
			if fields.Match(f) == nil {
				return errDupGoField(ptr, f)
			}
		}
	}
	seen := map[string]bool{}
	for i, f := range ty.Fields {
		if f.Name == "" {
			return errUnnamedField(ty, i)
		}
		sf := fields.Match(f.Name)
		if sf == nil {
			if lenient {
				continue
			}
			return errNoOrDupGoField(ptr, f.Name)
		}
		if seen[f.Name] {
			// We don't allow duplicated field name.
			return errDupSpannerField(f.Name, ty)
		}
		opts := []DecodeOptions{withLenient{lenient: lenient}}
		// Try to decode a single field.
		if err := decodeValue(pb.Values[i], f.Type, v.FieldByIndex(sf.Index).Addr().Interface(), opts...); err != nil {
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
func decodeStructArray(ty *sppb.StructType, pb *proto3.ListValue, ptr interface{}, lenient bool) error {
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
		if err = decodeStruct(ty, l, s.Interface(), lenient); err != nil {
			return errDecodeArrayElement(i, pv, "STRUCT", err)
		}
		// Append the decoded struct back into the slice.
		v.Set(reflect.Append(v, s))
	}
	return nil
}

func getAllFieldNames(v reflect.Value) []string {
	var names []string
	typeOfT := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fieldType := typeOfT.Field(i)
		exported := (fieldType.PkgPath == "")
		// If a named field is unexported, ignore it. An anonymous
		// unexported field is processed, because it may contain
		// exported fields, which are visible.
		if !exported && !fieldType.Anonymous {
			continue
		}
		if f.Kind() == reflect.Struct {
			if fieldType.Anonymous {
				names = append(names, getAllFieldNames(reflect.ValueOf(f.Interface()))...)
			}
			continue
		}
		name, keep, _, _ := spannerTagParser(fieldType.Tag)
		if !keep {
			continue
		}
		if name == "" {
			name = fieldType.Name
		}
		names = append(names, name)
	}
	return names
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
	case sql.NullString:
		if v.Valid {
			return encodeValue(v.String)
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
	case float32:
		pb.Kind = &proto3.Value_NumberValue{NumberValue: float64(v)}
		pt = float32Type()
	case []float32:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(float32Type())
	case NullFloat32:
		if v.Valid {
			return encodeValue(v.Float32)
		}
		pt = float32Type()
	case []NullFloat32:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(float32Type())
	case *float32:
		if v != nil {
			return encodeValue(*v)
		}
		pt = float32Type()
	case []*float32:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(float32Type())
	case big.Rat:
		switch LossOfPrecisionHandling {
		case NumericError:
			err = validateNumeric(&v)
			if err != nil {
				return nil, nil, err
			}
		case NumericRound:
			// pass
		}
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
	case PGNumeric:
		if v.Valid {
			pb.Kind = stringKind(v.Numeric)
		}
		return pb, pgNumericType(), nil
	case []PGNumeric:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(pgNumericType())
	case NullJSON:
		if v.Valid {
			b, err := json.Marshal(v.Value)
			if err != nil {
				return nil, nil, err
			}
			pb.Kind = stringKind(string(b))
		}
		return pb, jsonType(), nil
	case []NullJSON:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(jsonType())
	case PGJsonB:
		if v.Valid {
			b, err := json.Marshal(v.Value)
			if err != nil {
				return nil, nil, err
			}
			pb.Kind = stringKind(string(b))
		}
		return pb, pgJsonbType(), nil
	case []PGJsonB:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface{} { return v[i] })
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(pgJsonbType())
	case *big.Rat:
		switch LossOfPrecisionHandling {
		case NumericError:
			err = validateNumeric(v)
			if err != nil {
				return nil, nil, err
			}
		case NumericRound:
			// pass
		}
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
	case protoreflect.Enum:
		// Check if the value is of protoreflect.Enum type that implements spanner.Encoder
		// interface.
		if encodedVal, ok := v.(Encoder); ok {
			nv, err := encodedVal.EncodeSpanner()
			if err != nil {
				return nil, nil, err
			}
			return encodeValue(nv)
		}

		if v != nil {
			var protoEnumfqn string
			rv := reflect.ValueOf(v)
			if rv.Kind() != reflect.Ptr || !rv.IsNil() {
				pb.Kind = stringKind(strconv.FormatInt(int64(v.Number()), 10))
				protoEnumfqn = string(v.Descriptor().FullName())
			} else {
				defaultType := reflect.Zero(rv.Type().Elem()).Interface().(protoreflect.Enum)
				protoEnumfqn = string(defaultType.Descriptor().FullName())
			}
			pt = protoEnumType(protoEnumfqn)
		}
	case NullProtoEnum:
		if v.Valid {
			return encodeValue(v.ProtoEnumVal)
		}
		return nil, nil, errNotValidSrc(v)
	case proto.Message:
		// Check if the value is of proto.Message type that implements spanner.Encoder
		// interface.
		if encodedVal, ok := v.(Encoder); ok {
			nv, err := encodedVal.EncodeSpanner()
			if err != nil {
				return nil, nil, err
			}
			return encodeValue(nv)
		}

		if v != nil {
			if v.ProtoReflect().IsValid() {
				bytes, err := proto.Marshal(v)
				if err != nil {
					return nil, nil, err
				}
				pb.Kind = stringKind(base64.StdEncoding.EncodeToString(bytes))
			}
			protoMessagefqn := string(v.ProtoReflect().Descriptor().FullName())
			pt = protoMessageType(protoMessagefqn)
		}
	case NullProtoMessage:
		if v.Valid {
			return encodeValue(v.ProtoMessageVal)
		}
		return nil, nil, errNotValidSrc(v)
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

		if !isStructOrArrayOfStructValue(v) && !isAnArrayOfProtoColumn(v) {
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
			if isAnArrayOfProtoColumn(v) {
				return encodeProtoArray(v)
			}
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
		return nil, errors.New("cannot encode a value to type spannerTypeInvalid")
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
	case spannerTypeNonNullFloat32:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(float32(0.0))))
	case spannerTypeNullFloat32:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullFloat32{})))
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
	case spannerTypeNullJSON:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullJSON{})))
	case spannerTypePGJsonB:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(PGJsonB{})))
	case spannerTypePGNumeric:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(PGNumeric{})))
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
	case spannerTypeArrayOfNonNullFloat32:
		if reflect.ValueOf(v).IsNil() {
			return []float32(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]float32{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullFloat32:
		if reflect.ValueOf(v).IsNil() {
			return []NullFloat32(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullFloat32{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
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
	case spannerTypeArrayOfNullJSON:
		if reflect.ValueOf(v).IsNil() {
			return []NullJSON(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullJSON{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfPGJsonB:
		if reflect.ValueOf(v).IsNil() {
			return []PGJsonB(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]PGJsonB{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfPGNumeric:
		if reflect.ValueOf(v).IsNil() {
			return []PGNumeric(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]PGNumeric{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
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

// Encodes a slice of proto messages or enum in v to the spanner Value and Type
// protos.
func encodeProtoArray(v interface{}) (*proto3.Value, *sppb.Type, error) {
	pb := nullProto()
	var pt *sppb.Type
	var err error
	sliceval := reflect.ValueOf(v)
	etyp := reflect.TypeOf(v).Elem()

	if etyp.Implements(protoMsgReflectType) {
		if !sliceval.IsNil() {
			pb, err = encodeProtoMessageArray(sliceval.Len(), func(i int) reflect.Value { return sliceval.Index(i) })
			if err != nil {
				return nil, nil, err
			}
		}
		defaultInstance := reflect.Zero(etyp).Interface().(proto.Message)
		protoMessagefqn := string(defaultInstance.ProtoReflect().Descriptor().FullName())
		pt = listType(protoMessageType(protoMessagefqn))
	} else if etyp.Implements(protoEnumReflectType) {
		if !sliceval.IsNil() {
			pb, err = encodeProtoEnumArray(sliceval.Len(), func(i int) reflect.Value { return sliceval.Index(i) })
			if err != nil {
				return nil, nil, err
			}
		}
		if etyp.Kind() == reflect.Ptr {
			etyp = etyp.Elem()
		}
		defaultInstance := reflect.Zero(etyp).Interface().(protoreflect.Enum)
		protoEnumfqn := string(defaultInstance.Descriptor().FullName())
		pt = listType(protoEnumType(protoEnumfqn))
	}
	return pb, pt, nil
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

func isAnArrayOfProtoColumn(v interface{}) bool {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	return typ.Implements(protoMsgReflectType) || typ.Implements(protoEnumReflectType)
}

func isSupportedMutationType(v interface{}) bool {
	switch v.(type) {
	case nil, string, *string, NullString, []string, []*string, []NullString,
		[]byte, [][]byte,
		int, []int, int64, *int64, []int64, []*int64, NullInt64, []NullInt64,
		bool, *bool, []bool, []*bool, NullBool, []NullBool,
		float64, *float64, []float64, []*float64, NullFloat64, []NullFloat64,
		float32, *float32, []float32, []*float32, NullFloat32, []NullFloat32,
		time.Time, *time.Time, []time.Time, []*time.Time, NullTime, []NullTime,
		civil.Date, *civil.Date, []civil.Date, []*civil.Date, NullDate, []NullDate,
		big.Rat, *big.Rat, []big.Rat, []*big.Rat, NullNumeric, []NullNumeric,
		GenericColumnValue, proto.Message, protoreflect.Enum, NullProtoMessage, NullProtoEnum:
		return true
	default:
		// Check if the custom type implements spanner.Encoder interface.
		if _, ok := v.(Encoder); ok {
			return true
		}

		if isAnArrayOfProtoColumn(v) {
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

func encodeProtoMessageArray(len int, at func(int) reflect.Value) (*proto3.Value, error) {
	vs := make([]*proto3.Value, len)
	var err error
	for i := 0; i < len; i++ {
		v := at(i).Interface().(proto.Message)
		vs[i], _, err = encodeValue(v)
		if err != nil {
			return nil, err
		}
	}
	return listProto(vs...), nil
}

func encodeProtoEnumArray(len int, at func(int) reflect.Value) (*proto3.Value, error) {
	vs := make([]*proto3.Value, len)
	var err error
	for i := 0; i < len; i++ {
		v := at(i).Interface().(protoreflect.Enum)
		vs[i], _, err = encodeValue(v)
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
