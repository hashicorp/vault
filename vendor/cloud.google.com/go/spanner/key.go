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
	"fmt"
	"math/big"
	"time"

	"cloud.google.com/go/civil"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/reflect/protoreflect"
	proto3 "google.golang.org/protobuf/types/known/structpb"
)

// A Key can be either a Cloud Spanner row's primary key or a secondary index
// key. It is essentially an interface{} array, which represents a set of Cloud
// Spanner columns. A Key can be used as:
//
//   - A primary key which uniquely identifies a Cloud Spanner row.
//   - A secondary index key which maps to a set of Cloud Spanner rows indexed under it.
//   - An endpoint of primary key/secondary index ranges; see the KeyRange type.
//
// Rows that are identified by the Key type are outputs of read operation or
// targets of delete operation in a mutation. Note that for
// Insert/Update/InsertOrUpdate/Update mutation types, although they don't
// require a primary key explicitly, the column list provided must contain
// enough columns that can comprise a primary key.
//
// Keys are easy to construct.  For example, suppose you have a table with a
// primary key of username and product ID.  To make a key for this table:
//
//	key := spanner.Key{"john", 16}
//
// See the description of Row and Mutation types for how Go types are mapped to
// Cloud Spanner types. For convenience, Key type supports a wide range of Go
// types:
//   - int, int8, int16, int32, int64, and NullInt64 are mapped to Cloud Spanner's INT64 type.
//   - uint8, uint16 and uint32 are also mapped to Cloud Spanner's INT64 type.
//   - float32, float64, NullFloat64 are mapped to Cloud Spanner's FLOAT64 type.
//   - bool and NullBool are mapped to Cloud Spanner's BOOL type.
//   - []byte is mapped to Cloud Spanner's BYTES type.
//   - string and NullString are mapped to Cloud Spanner's STRING type.
//   - time.Time and NullTime are mapped to Cloud Spanner's TIMESTAMP type.
//   - civil.Date and NullDate are mapped to Cloud Spanner's DATE type.
//   - protoreflect.Enum and NullProtoEnum are mapped to Cloud Spanner's ENUM type.
type Key []interface{}

// errInvdKeyPartType returns error for unsupported key part type.
func errInvdKeyPartType(part interface{}) error {
	return spannerErrorf(codes.InvalidArgument, "key part has unsupported type %T", part)
}

// keyPartValue converts a part of the Key (which is a valid Cloud Spanner type)
// into a proto3.Value. Used for encoding Key type into protobuf.
func keyPartValue(part interface{}) (pb *proto3.Value, err error) {
	switch v := part.(type) {
	case int:
		pb, _, err = encodeValue(int64(v))
	case int8:
		pb, _, err = encodeValue(int64(v))
	case int16:
		pb, _, err = encodeValue(int64(v))
	case int32:
		pb, _, err = encodeValue(int64(v))
	case uint8:
		pb, _, err = encodeValue(int64(v))
	case uint16:
		pb, _, err = encodeValue(int64(v))
	case uint32:
		pb, _, err = encodeValue(int64(v))
	case int64, float64, float32, NullInt64, NullFloat64, NullFloat32, bool, NullBool, []byte, string, NullString, time.Time, civil.Date, NullTime, NullDate, big.Rat, NullNumeric, protoreflect.Enum, NullProtoEnum:
		pb, _, err = encodeValue(v)
	case Encoder:
		part, err = v.EncodeSpanner()
		if err != nil {
			return nil, err
		}
		pb, err = keyPartValue(part)
	default:
		return nil, errInvdKeyPartType(v)
	}
	return pb, err
}

// proto converts a spanner.Key into a proto3.ListValue.
func (key Key) proto() (*proto3.ListValue, error) {
	lv := &proto3.ListValue{}
	lv.Values = make([]*proto3.Value, 0, len(key))
	for _, part := range key {
		v, err := keyPartValue(part)
		if err != nil {
			return nil, err
		}
		lv.Values = append(lv.Values, v)
	}
	return lv, nil
}

// keySetProto lets a single Key act as a KeySet.
func (key Key) keySetProto() (*sppb.KeySet, error) {
	kp, err := key.proto()
	if err != nil {
		return nil, err
	}
	return &sppb.KeySet{Keys: []*proto3.ListValue{kp}}, nil
}

// String implements fmt.Stringer for Key. For string, []byte and NullString, it
// prints the uninterpreted bytes of their contents, leaving caller with the
// opportunity to escape the output.
func (key Key) String() string {
	b := &bytes.Buffer{}
	fmt.Fprint(b, "(")
	for i, part := range []interface{}(key) {
		if i != 0 {
			fmt.Fprint(b, ",")
		}
		key.elemString(b, part)
	}
	fmt.Fprint(b, ")")
	return b.String()
}

func (key Key) elemString(b *bytes.Buffer, part interface{}) {
	switch v := part.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64, bool, protoreflect.Enum:
		// Use %v to print numeric types and bool.
		fmt.Fprintf(b, "%v", v)
	case string:
		fmt.Fprintf(b, "%q", v)
	case []byte:
		if v != nil {
			fmt.Fprintf(b, "%q", v)
		} else {
			fmt.Fprint(b, nullString)
		}
	case NullInt64, NullFloat64, NullBool, NullNumeric, NullProtoEnum:
		// The above types implement fmt.Stringer.
		fmt.Fprintf(b, "%s", v)
	case NullString, NullDate, NullTime:
		// Quote the returned string if it is not null.
		if v.(NullableValue).IsNull() {
			fmt.Fprintf(b, "%s", nullString)
		} else {
			fmt.Fprintf(b, "%q", v)
		}
	case civil.Date:
		fmt.Fprintf(b, "%q", v)
	case time.Time:
		fmt.Fprintf(b, "%q", v.Format(time.RFC3339Nano))
	case big.Rat:
		fmt.Fprintf(b, "%v", NumericString(&v))
	case Encoder:
		var err error
		part, err = v.EncodeSpanner()
		if err != nil {
			fmt.Fprintf(b, "error")
		} else {
			key.elemString(b, part)
		}
	default:
		fmt.Fprintf(b, "%v", v)
	}
}

// AsPrefix returns a KeyRange for all keys where k is the prefix.
func (key Key) AsPrefix() KeyRange {
	return KeyRange{
		Start: key,
		End:   key,
		Kind:  ClosedClosed,
	}
}

// KeyRangeKind describes the kind of interval represented by a KeyRange:
// whether it is open or closed on the left and right.
type KeyRangeKind int

const (
	// ClosedOpen is closed on the left and open on the right: the Start
	// key is included, the End key is excluded.
	ClosedOpen KeyRangeKind = iota

	// ClosedClosed is closed on the left and the right: both keys are included.
	ClosedClosed

	// OpenClosed is open on the left and closed on the right: the Start
	// key is excluded, the End key is included.
	OpenClosed

	// OpenOpen is open on the left and the right: neither key is included.
	OpenOpen
)

// A KeyRange represents a range of rows in a table or index.
//
// A range has a Start key and an End key.  IncludeStart and IncludeEnd
// indicate whether the Start and End keys are included in the range.
//
// For example, consider the following table definition:
//
//	CREATE TABLE UserEvents (
//	  UserName STRING(MAX),
//	  EventDate STRING(10),
//	) PRIMARY KEY(UserName, EventDate);
//
// The following keys name rows in this table:
//
//	spanner.Key{"Bob", "2014-09-23"}
//	spanner.Key{"Alfred", "2015-06-12"}
//
// Since the UserEvents table's PRIMARY KEY clause names two columns, each
// UserEvents key has two elements; the first is the UserName, and the second
// is the EventDate.
//
// Key ranges with multiple components are interpreted lexicographically by
// component using the table or index key's declared sort order. For example,
// the following range returns all events for user "Bob" that occurred in the
// year 2015:
//
//	spanner.KeyRange{
//		Start: spanner.Key{"Bob", "2015-01-01"},
//		End:   spanner.Key{"Bob", "2015-12-31"},
//		Kind:  spanner.ClosedClosed,
//	}
//
// Start and end keys can omit trailing key components. This affects the
// inclusion and exclusion of rows that exactly match the provided key
// components: if IncludeStart is true, then rows that exactly match the
// provided components of the Start key are included; if IncludeStart is false
// then rows that exactly match are not included.  IncludeEnd and End key
// behave in the same fashion.
//
// For example, the following range includes all events for "Bob" that occurred
// during and after the year 2000:
//
//	spanner.KeyRange{
//		Start: spanner.Key{"Bob", "2000-01-01"},
//		End:   spanner.Key{"Bob"},
//		Kind:  spanner.ClosedClosed,
//	}
//
// The next example retrieves all events for "Bob":
//
//	spanner.Key{"Bob"}.AsPrefix()
//
// To retrieve events before the year 2000:
//
//	spanner.KeyRange{
//		Start: spanner.Key{"Bob"},
//		End:   spanner.Key{"Bob", "2000-01-01"},
//		Kind:  spanner.ClosedOpen,
//	}
//
// Although we specified a Kind for this KeyRange, we didn't need to, because
// the default is ClosedOpen. In later examples we'll omit Kind if it is
// ClosedOpen.
//
// The following range includes all rows in a table or under a
// index:
//
//	spanner.AllKeys()
//
// This range returns all users whose UserName begins with any
// character from A to C:
//
//	spanner.KeyRange{
//		Start: spanner.Key{"A"},
//		End:   spanner.Key{"D"},
//	}
//
// This range returns all users whose UserName begins with B:
//
//	spanner.KeyRange{
//		Start: spanner.Key{"B"},
//		End:   spanner.Key{"C"},
//	}
//
// Key ranges honor column sort order. For example, suppose a table is defined
// as follows:
//
//	CREATE TABLE DescendingSortedTable {
//	  Key INT64,
//	  ...
//	) PRIMARY KEY(Key DESC);
//
// The following range retrieves all rows with key values between 1 and 100
// inclusive:
//
//	spanner.KeyRange{
//		Start: spanner.Key{100},
//		End:   spanner.Key{1},
//		Kind:  spanner.ClosedClosed,
//	}
//
// Note that 100 is passed as the start, and 1 is passed as the end, because
// Key is a descending column in the schema.
type KeyRange struct {
	// Start specifies the left boundary of the key range; End specifies
	// the right boundary of the key range.
	Start, End Key

	// Kind describes whether the boundaries of the key range include
	// their keys.
	Kind KeyRangeKind
}

// String implements fmt.Stringer for KeyRange type.
func (r KeyRange) String() string {
	var left, right string
	switch r.Kind {
	case ClosedClosed:
		left, right = "[", "]"
	case ClosedOpen:
		left, right = "[", ")"
	case OpenClosed:
		left, right = "(", "]"
	case OpenOpen:
		left, right = "(", ")"
	default:
		left, right = "?", "?"
	}
	return fmt.Sprintf("%s%s,%s%s", left, r.Start, r.End, right)
}

// proto converts KeyRange into sppb.KeyRange.
func (r KeyRange) proto() (*sppb.KeyRange, error) {
	var err error
	var start, end *proto3.ListValue
	pb := &sppb.KeyRange{}
	if start, err = r.Start.proto(); err != nil {
		return nil, err
	}
	if end, err = r.End.proto(); err != nil {
		return nil, err
	}
	if r.Kind == ClosedClosed || r.Kind == ClosedOpen {
		pb.StartKeyType = &sppb.KeyRange_StartClosed{StartClosed: start}
	} else {
		pb.StartKeyType = &sppb.KeyRange_StartOpen{StartOpen: start}
	}
	if r.Kind == ClosedClosed || r.Kind == OpenClosed {
		pb.EndKeyType = &sppb.KeyRange_EndClosed{EndClosed: end}
	} else {
		pb.EndKeyType = &sppb.KeyRange_EndOpen{EndOpen: end}
	}
	return pb, nil
}

// keySetProto lets a KeyRange act as a KeySet.
func (r KeyRange) keySetProto() (*sppb.KeySet, error) {
	rp, err := r.proto()
	if err != nil {
		return nil, err
	}
	return &sppb.KeySet{Ranges: []*sppb.KeyRange{rp}}, nil
}

// A KeySet defines a collection of Cloud Spanner keys and/or key ranges. All
// the keys are expected to be in the same table or index. The keys need not be
// sorted in any particular way.
//
// An individual Key can act as a KeySet, as can a KeyRange. Use the KeySets
// function to create a KeySet consisting of multiple Keys and KeyRanges. To
// obtain an empty KeySet, call KeySets with no arguments.
//
// If the same key is specified multiple times in the set (for example if two
// ranges, two keys, or a key and a range overlap), the Cloud Spanner backend
// behaves as if the key were only specified once.
type KeySet interface {
	keySetProto() (*sppb.KeySet, error)
}

// AllKeys returns a KeySet that represents all Keys of a table or a index.
func AllKeys() KeySet {
	return all{}
}

type all struct{}

func (all) keySetProto() (*sppb.KeySet, error) {
	return &sppb.KeySet{All: true}, nil
}

// KeySets returns the union of the KeySets. If any of the KeySets is AllKeys,
// then the resulting KeySet will be equivalent to AllKeys.
func KeySets(keySets ...KeySet) KeySet {
	u := make(union, len(keySets))
	copy(u, keySets)
	return u
}

// KeySetFromKeys returns a KeySet containing the given slice of keys.
func KeySetFromKeys(keys ...Key) KeySet {
	u := make(union, len(keys))
	for i, k := range keys {
		u[i] = k
	}
	return u
}

type union []KeySet

func (u union) keySetProto() (*sppb.KeySet, error) {
	upb := &sppb.KeySet{}
	for _, ks := range u {
		pb, err := ks.keySetProto()
		if err != nil {
			return nil, err
		}
		if pb.All {
			return pb, nil
		}
		upb.Keys = append(upb.Keys, pb.Keys...)
		upb.Ranges = append(upb.Ranges, pb.Ranges...)
	}
	return upb, nil
}
