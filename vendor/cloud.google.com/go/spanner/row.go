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
	"errors"
	"fmt"
	"reflect"
	"strings"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/grpc/codes"
	proto3 "google.golang.org/protobuf/types/known/structpb"
)

// A Row is a view of a row of data returned by a Cloud Spanner read.
// It consists of a number of columns; the number depends on the columns
// used to construct the read.
//
// The column values can be accessed by index. For instance, if the read specified
// []string{"photo_id", "caption"}, then each row will contain two
// columns: "photo_id" with index 0, and "caption" with index 1.
//
// Column values are decoded by using one of the Column, ColumnByName, or
// Columns methods. The valid values passed to these methods depend on the
// column type. For example:
//
//	var photoID int64
//	err := row.Column(0, &photoID) // Decode column 0 as an integer.
//
//	var caption string
//	err := row.Column(1, &caption) // Decode column 1 as a string.
//
//	// Decode all the columns.
//	err := row.Columns(&photoID, &caption)
//
// Supported types and their corresponding Cloud Spanner column type(s) are:
//
//	*string(not NULL), *NullString - STRING
//	*[]string, *[]NullString - STRING ARRAY
//	*[]byte - BYTES
//	*[][]byte - BYTES ARRAY
//	*int64(not NULL), *NullInt64 - INT64
//	*[]int64, *[]NullInt64 - INT64 ARRAY
//	*bool(not NULL), *NullBool - BOOL
//	*[]bool, *[]NullBool - BOOL ARRAY
//	*float32(not NULL), *NullFloat32 - FLOAT32
//	*[]float32, *[]NullFloat32 - FLOAT32 ARRAY
//	*float64(not NULL), *NullFloat64 - FLOAT64
//	*[]float64, *[]NullFloat64 - FLOAT64 ARRAY
//	*big.Rat(not NULL), *NullNumeric - NUMERIC
//	*[]big.Rat, *[]NullNumeric - NUMERIC ARRAY
//	*time.Time(not NULL), *NullTime - TIMESTAMP
//	*[]time.Time, *[]NullTime - TIMESTAMP ARRAY
//	*Date(not NULL), *NullDate - DATE
//	*[]civil.Date, *[]NullDate - DATE ARRAY
//	*[]*some_go_struct, *[]NullRow - STRUCT ARRAY
//	*NullJSON - JSON
//	*[]NullJSON - JSON ARRAY
//	*GenericColumnValue - any Cloud Spanner type
//
// For TIMESTAMP columns, the returned time.Time object will be in UTC.
//
// To fetch an array of BYTES, pass a *[][]byte. To fetch an array of (sub)rows, pass
// a *[]spanner.NullRow or a *[]*some_go_struct where some_go_struct holds all
// information of the subrow, see spanner.Row.ToStruct for the mapping between a
// Cloud Spanner row and a Go struct. To fetch an array of other types, pass a
// *[]spanner.NullXXX type of the appropriate type. Use GenericColumnValue when you
// don't know in advance what column type to expect.
//
// Row decodes the row contents lazily; as a result, each call to a getter has
// a chance of returning an error.
//
// A column value may be NULL if the corresponding value is not present in
// Cloud Spanner. The spanner.NullXXX types (spanner.NullInt64 et al.) allow fetching
// values that may be null. A NULL BYTES can be fetched into a *[]byte as nil.
// It is an error to fetch a NULL value into any other type.
type Row struct {
	fields []*sppb.StructType_Field
	vals   []*proto3.Value // keep decoded for now
}

// String implements fmt.stringer.
func (r *Row) String() string {
	return fmt.Sprintf("{fields: %s, values: %s}", r.fields, r.vals)
}

// errNamesValuesMismatch returns error for when columnNames count is not equal
// to columnValues count.
func errNamesValuesMismatch(columnNames []string, columnValues []interface{}) error {
	return spannerErrorf(codes.FailedPrecondition,
		"different number of names(%v) and values(%v)", len(columnNames), len(columnValues))
}

// NewRow returns a Row containing the supplied data.  This can be useful for
// mocking Cloud Spanner Read and Query responses for unit testing.
func NewRow(columnNames []string, columnValues []interface{}) (*Row, error) {
	if len(columnValues) != len(columnNames) {
		return nil, errNamesValuesMismatch(columnNames, columnValues)
	}
	r := Row{
		fields: make([]*sppb.StructType_Field, len(columnValues)),
		vals:   make([]*proto3.Value, len(columnValues)),
	}
	for i := range columnValues {
		val, typ, err := encodeValue(columnValues[i])
		if err != nil {
			return nil, err
		}
		r.fields[i] = &sppb.StructType_Field{
			Name: columnNames[i],
			Type: typ,
		}
		r.vals[i] = val
	}
	return &r, nil
}

// Size is the number of columns in the row.
func (r *Row) Size() int {
	return len(r.fields)
}

// ColumnName returns the name of column i, or empty string for invalid column.
func (r *Row) ColumnName(i int) string {
	if i < 0 || i >= len(r.fields) {
		return ""
	}
	return r.fields[i].Name
}

// ColumnIndex returns the index of the column with the given name. The
// comparison is case-sensitive.
func (r *Row) ColumnIndex(name string) (int, error) {
	found := false
	var index int
	if len(r.vals) != len(r.fields) {
		return 0, errFieldsMismatchVals(r)
	}
	for i, f := range r.fields {
		if f == nil {
			return 0, errNilColType(i)
		}
		if name == f.Name {
			if found {
				return 0, errDupColName(name)
			}
			found = true
			index = i
		}
	}
	if !found {
		return 0, errColNotFound(name)
	}
	return index, nil
}

// ColumnNames returns all column names of the row.
func (r *Row) ColumnNames() []string {
	var n []string
	for _, c := range r.fields {
		n = append(n, c.Name)
	}
	return n
}

// ColumnType returns the Cloud Spanner Type of column i, or nil for invalid column.
func (r *Row) ColumnType(i int) *sppb.Type {
	if i < 0 || i >= len(r.fields) {
		return nil
	}
	return r.fields[i].Type
}

// ColumnValue returns the Cloud Spanner Value of column i, or nil for invalid column.
func (r *Row) ColumnValue(i int) *proto3.Value {
	if i < 0 || i >= len(r.vals) {
		return nil
	}
	return r.vals[i]
}

// errColIdxOutOfRange returns error for requested column index is out of the
// range of the target Row's columns.
func errColIdxOutOfRange(i int, r *Row) error {
	return spannerErrorf(codes.OutOfRange, "column index %d out of range [0,%d)", i, len(r.vals))
}

// errDecodeColumn returns error for not being able to decode a indexed column.
func errDecodeColumn(i int, err error) error {
	if err == nil {
		return nil
	}
	var se *Error
	if !errors.As(err, &se) {
		return spannerErrorf(codes.InvalidArgument, "failed to decode column %v, error = <%v>", i, err)
	}
	se.decorate(fmt.Sprintf("failed to decode column %v", i))
	return se
}

// errFieldsMismatchVals returns error for field count isn't equal to value count in a Row.
func errFieldsMismatchVals(r *Row) error {
	return spannerErrorf(codes.FailedPrecondition, "row has different number of fields(%v) and values(%v)",
		len(r.fields), len(r.vals))
}

// errNilColType returns error for column type for column i being nil in the row.
func errNilColType(i int) error {
	return spannerErrorf(codes.FailedPrecondition, "column(%v)'s type is nil", i)
}

// Column fetches the value from the ith column, decoding it into ptr.
// See the Row documentation for the list of acceptable argument types.
// see Client.ReadWriteTransaction for an example.
func (r *Row) Column(i int, ptr interface{}) error {
	if len(r.vals) != len(r.fields) {
		return errFieldsMismatchVals(r)
	}
	if i < 0 || i >= len(r.fields) {
		return errColIdxOutOfRange(i, r)
	}
	if r.fields[i] == nil {
		return errNilColType(i)
	}
	if err := decodeValue(r.vals[i], r.fields[i].Type, ptr); err != nil {
		return errDecodeColumn(i, err)
	}
	return nil
}

// errDupColName returns error for duplicated column name in the same row.
func errDupColName(n string) error {
	return spannerErrorf(codes.FailedPrecondition, "ambiguous column name %q", n)
}

// errColNotFound returns error for not being able to find a named column.
func errColNotFound(n string) error {
	return spannerErrorf(codes.NotFound, "column %q not found", n)
}

func errNotASlicePointer() error {
	return spannerErrorf(codes.InvalidArgument, "destination must be a pointer to a slice")
}

func errNilSlicePointer() error {
	return spannerErrorf(codes.InvalidArgument, "destination must be a non nil pointer")
}

func errTooManyColumns() error {
	return spannerErrorf(codes.InvalidArgument, "too many columns returned for primitive slice")
}

// ColumnByName fetches the value from the named column, decoding it into ptr.
// See the Row documentation for the list of acceptable argument types.
func (r *Row) ColumnByName(name string, ptr interface{}) error {
	index, err := r.ColumnIndex(name)
	if err != nil {
		return err
	}
	return r.Column(index, ptr)
}

// errNumOfColValue returns error for providing wrong number of values to Columns.
func errNumOfColValue(n int, r *Row) error {
	return spannerErrorf(codes.InvalidArgument,
		"Columns(): number of arguments (%d) does not match row size (%d)", n, len(r.vals))
}

// Columns fetches all the columns in the row at once.
//
// The value of the kth column will be decoded into the kth argument to Columns. See
// Row for the list of acceptable argument types. The number of arguments must be
// equal to the number of columns. Pass nil to specify that a column should be
// ignored.
func (r *Row) Columns(ptrs ...interface{}) error {
	if len(ptrs) != len(r.vals) {
		return errNumOfColValue(len(ptrs), r)
	}
	if len(r.vals) != len(r.fields) {
		return errFieldsMismatchVals(r)
	}
	for i, p := range ptrs {
		if p == nil {
			continue
		}
		if err := r.Column(i, p); err != nil {
			return err
		}
	}
	return nil
}

// errToStructArgType returns error for p not having the correct data type(pointer to Go struct) to
// be the argument of Row.ToStruct.
func errToStructArgType(p interface{}) error {
	return spannerErrorf(codes.InvalidArgument, "ToStruct(): type %T is not a valid pointer to Go struct", p)
}

// ToStruct fetches the columns in a row into the fields of a struct.
// The rules for mapping a row's columns into a struct's exported fields
// are:
//
//  1. If a field has a `spanner: "column_name"` tag, then decode column
//     'column_name' into the field. A special case is the `spanner: "-"`
//     tag, which instructs ToStruct to ignore the field during decoding.
//
//  2. Otherwise, if the name of a field matches the name of a column (ignoring case),
//     decode the column into the field.
//
//  3. The number of columns in the row must match the number of exported fields in the struct.
//     There must be exactly one match for each column in the row. The method will return an error
//     if a column in the row cannot be assigned to a field in the struct.
//
// The fields of the destination struct can be of any type that is acceptable
// to spanner.Row.Column.
//
// Slice and pointer fields will be set to nil if the source column is NULL, and a
// non-nil value if the column is not NULL. To decode NULL values of other types, use
// one of the spanner.NullXXX types as the type of the destination field.
//
// If ToStruct returns an error, the contents of p are undefined. Some fields may
// have been successfully populated, while others were not; you should not use any of
// the fields.
func (r *Row) ToStruct(p interface{}) error {
	// Check if p is a pointer to a struct
	if t := reflect.TypeOf(p); t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errToStructArgType(p)
	}
	if len(r.vals) != len(r.fields) {
		return errFieldsMismatchVals(r)
	}
	// Call decodeStruct directly to decode the row as a typed proto.ListValue.
	return decodeStruct(
		&sppb.StructType{Fields: r.fields},
		&proto3.ListValue{Values: r.vals},
		p,
		false,
	)
}

// ToStructLenient fetches the columns in a row into the fields of a struct.
// The rules for mapping a row's columns into a struct's exported fields
// are:
//
//  1. If a field has a `spanner: "column_name"` tag, then decode column
//     'column_name' into the field. A special case is the `spanner: "-"`
//     tag, which instructs ToStruct to ignore the field during decoding.
//
//  2. Otherwise, if the name of a field matches the name of a column (ignoring case),
//     decode the column into the field.
//
//  3. The number of columns in the row and exported fields in the struct do not need to match.
//     Any field in the struct that cannot not be assigned a value from the row is assigned its default value.
//     Any column in the row that does not have a corresponding field in the struct is ignored.
//
// The fields of the destination struct can be of any type that is acceptable
// to spanner.Row.Column.
//
// Slice and pointer fields will be set to nil if the source column is NULL, and a
// non-nil value if the column is not NULL. To decode NULL values of other types, use
// one of the spanner.NullXXX types as the type of the destination field.
//
// If ToStructLenient returns an error, the contents of p are undefined. Some fields may
// have been successfully populated, while others were not; you should not use any of
// the fields.
func (r *Row) ToStructLenient(p interface{}) error {
	// Check if p is a pointer to a struct
	if t := reflect.TypeOf(p); t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errToStructArgType(p)
	}
	if len(r.vals) != len(r.fields) {
		return errFieldsMismatchVals(r)
	}
	// Call decodeStruct directly to decode the row as a typed proto.ListValue.
	return decodeStruct(
		&sppb.StructType{Fields: r.fields},
		&proto3.ListValue{Values: r.vals},
		p,
		true,
	)
}

// SelectAll iterates all rows to the end. After iterating it closes the rows
// and propagates any errors that could pop up with destination slice partially filled.
// It expects that destination should be a slice. For each row, it scans data and appends it to the destination slice.
// SelectAll supports both types of slices: slice of pointers and slice of structs or primitives by value,
// for example:
//
//	type Singer struct {
//	    ID    string
//	    Name  string
//	}
//
//	var singersByPtr []*Singer
//	var singersByValue []Singer
//
// Both singersByPtr and singersByValue are valid destinations for SelectAll function.
//
// Add the option `spanner.WithLenient()` to instruct SelectAll to ignore additional columns in the rows that are not present in the destination struct.
// example:
//
//	var singersByPtr []*Singer
//	err := spanner.SelectAll(row, &singersByPtr, spanner.WithLenient())
func SelectAll(rows rowIterator, destination interface{}, options ...DecodeOptions) error {
	if rows == nil {
		return errors.New("rows is nil")
	}
	if destination == nil {
		return errors.New("destination is nil")
	}
	dstVal := reflect.ValueOf(destination)
	if !dstVal.IsValid() || (dstVal.Kind() == reflect.Ptr && dstVal.IsNil()) {
		return errNilSlicePointer()
	}
	if dstVal.Kind() != reflect.Ptr {
		return errNotASlicePointer()
	}
	dstVal = dstVal.Elem()
	dstType := dstVal.Type()
	if k := dstType.Kind(); k != reflect.Slice {
		return errNotASlicePointer()
	}

	itemType := dstType.Elem()
	var itemByPtr bool
	// If it's a slice of pointers to structs,
	// we handle it the same way as it would be slice of struct by value
	// and dereference pointers to values,
	// because eventually we work with fields.
	// But if it's a slice of primitive type e.g. or []string or []*string,
	// we must leave and pass elements as is.
	if itemType.Kind() == reflect.Ptr {
		elementBaseTypeElem := itemType.Elem()
		if elementBaseTypeElem.Kind() == reflect.Struct {
			itemType = elementBaseTypeElem
			itemByPtr = true
		}
	}
	s := &decodeSetting{}
	for _, opt := range options {
		opt.Apply(s)
	}

	isPrimitive := itemType.Kind() != reflect.Struct
	var pointers []interface{}
	isFirstRow := true
	var err error
	return rows.Do(func(row *Row) error {
		sliceItem := reflect.New(itemType)
		if isFirstRow && !isPrimitive {
			defer func() {
				isFirstRow = false
			}()
			if pointers, err = structPointers(sliceItem.Elem(), row.fields, s.Lenient); err != nil {
				return err
			}
		} else if isPrimitive {
			if len(row.fields) > 1 && !s.Lenient {
				return errTooManyColumns()
			}
			pointers = []interface{}{sliceItem.Interface()}
		}
		if len(pointers) == 0 {
			return nil
		}
		err = row.Columns(pointers...)
		if err != nil {
			return err
		}
		if !isPrimitive {
			e := sliceItem.Elem()
			idx := 0
			for _, p := range pointers {
				if p == nil {
					continue
				}
				e.Field(idx).Set(reflect.ValueOf(p).Elem())
				idx++
			}
		}
		var elemVal reflect.Value
		if itemByPtr {
			if isFirstRow {
				// create a new pointer to the struct with all the values copied from sliceItem
				// because same underlying pointers array will be used for next rows
				elemVal = reflect.New(itemType)
				elemVal.Elem().Set(sliceItem.Elem())
			} else {
				elemVal = sliceItem
			}
		} else {
			elemVal = sliceItem.Elem()
		}
		dstVal.Set(reflect.Append(dstVal, elemVal))
		return nil
	})
}

func structPointers(sliceItem reflect.Value, cols []*sppb.StructType_Field, lenient bool) ([]interface{}, error) {
	pointers := make([]interface{}, 0, len(cols))
	fieldTag := make(map[string]reflect.Value, len(cols))
	initFieldTag(sliceItem, &fieldTag)

	for i, colName := range cols {
		if colName.Name == "" {
			return nil, errColNotFound(fmt.Sprintf("column %d", i))
		}

		var fieldVal reflect.Value
		if v, ok := fieldTag[strings.ToLower(colName.GetName())]; ok {
			fieldVal = v
		} else {
			if !lenient {
				return nil, errNoOrDupGoField(sliceItem, colName.GetName())
			}
			fieldVal = sliceItem.FieldByName(colName.GetName())
		}
		if !fieldVal.IsValid() || !fieldVal.CanSet() {
			// have to add if we found a column because Columns() requires
			// len(cols) arguments or it will error. This way we can scan to
			// a useless pointer
			pointers = append(pointers, nil)
			continue
		}

		pointers = append(pointers, fieldVal.Addr().Interface())
	}
	return pointers, nil
}

// Initialization the tags from struct.
func initFieldTag(sliceItem reflect.Value, fieldTagMap *map[string]reflect.Value) {
	typ := sliceItem.Type()

	for i := 0; i < sliceItem.NumField(); i++ {
		fieldType := typ.Field(i)
		exported := (fieldType.PkgPath == "")
		// If a named field is unexported, ignore it. An anonymous
		// unexported field is processed, because it may contain
		// exported fields, which are visible.
		if !exported && !fieldType.Anonymous {
			continue
		}
		if fieldType.Type.Kind() == reflect.Struct {
			// found an embedded struct
			if fieldType.Anonymous {
				sliceItemOfAnonymous := sliceItem.Field(i)
				initFieldTag(sliceItemOfAnonymous, fieldTagMap)
				continue
			}
		}
		name, keep, _, _ := spannerTagParser(fieldType.Tag)
		if !keep {
			continue
		}
		if name == "" {
			name = fieldType.Name
		}
		(*fieldTagMap)[strings.ToLower(name)] = sliceItem.Field(i)
	}
}
