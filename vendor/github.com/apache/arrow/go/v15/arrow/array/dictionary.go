// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package array

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/bits"
	"sync/atomic"
	"unsafe"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/bitutil"
	"github.com/apache/arrow/go/v15/arrow/decimal128"
	"github.com/apache/arrow/go/v15/arrow/decimal256"
	"github.com/apache/arrow/go/v15/arrow/float16"
	"github.com/apache/arrow/go/v15/arrow/internal/debug"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/apache/arrow/go/v15/internal/hashing"
	"github.com/apache/arrow/go/v15/internal/json"
	"github.com/apache/arrow/go/v15/internal/utils"
)

// Dictionary represents the type for dictionary-encoded data with a data
// dependent dictionary.
//
// A dictionary array contains an array of non-negative integers (the "dictionary"
// indices") along with a data type containing a "dictionary" corresponding to
// the distinct values represented in the data.
//
// For example, the array:
//
//	["foo", "bar", "foo", "bar", "foo", "bar"]
//
// with dictionary ["bar", "foo"], would have the representation of:
//
//	indices: [1, 0, 1, 0, 1, 0]
//	dictionary: ["bar", "foo"]
//
// The indices in principle may be any integer type.
type Dictionary struct {
	array

	indices arrow.Array
	dict    arrow.Array
}

// NewDictionaryArray constructs a dictionary array with the provided indices
// and dictionary using the given type.
func NewDictionaryArray(typ arrow.DataType, indices, dict arrow.Array) *Dictionary {
	a := &Dictionary{}
	a.array.refCount = 1
	dictdata := NewData(typ, indices.Len(), indices.Data().Buffers(), indices.Data().Children(), indices.NullN(), indices.Data().Offset())
	dictdata.dictionary = dict.Data().(*Data)
	dict.Data().Retain()

	defer dictdata.Release()
	a.setData(dictdata)
	return a
}

// checkIndexBounds returns an error if any value in the provided integer
// arraydata is >= the passed upperlimit or < 0. otherwise nil
func checkIndexBounds(indices *Data, upperlimit uint64) error {
	if indices.length == 0 {
		return nil
	}

	var maxval uint64
	switch indices.dtype.ID() {
	case arrow.UINT8:
		maxval = math.MaxUint8
	case arrow.UINT16:
		maxval = math.MaxUint16
	case arrow.UINT32:
		maxval = math.MaxUint32
	case arrow.UINT64:
		maxval = math.MaxUint64
	}
	// for unsigned integers, if the values array is larger than the maximum
	// index value (especially for UINT8/UINT16), then there's no need to
	// boundscheck. for signed integers we still need to bounds check
	// because a value could be < 0.
	isSigned := maxval == 0
	if !isSigned && upperlimit > maxval {
		return nil
	}

	start := indices.offset
	end := indices.offset + indices.length

	// TODO(ARROW-15950): lift BitSetRunReader from parquet to utils
	// and use it here for performance improvement.

	switch indices.dtype.ID() {
	case arrow.INT8:
		data := arrow.Int8Traits.CastFromBytes(indices.buffers[1].Bytes())
		min, max := utils.GetMinMaxInt8(data[start:end])
		if min < 0 || max >= int8(upperlimit) {
			return fmt.Errorf("contains out of bounds index: min: %d, max: %d", min, max)
		}
	case arrow.UINT8:
		data := arrow.Uint8Traits.CastFromBytes(indices.buffers[1].Bytes())
		_, max := utils.GetMinMaxUint8(data[start:end])
		if max >= uint8(upperlimit) {
			return fmt.Errorf("contains out of bounds index: max: %d", max)
		}
	case arrow.INT16:
		data := arrow.Int16Traits.CastFromBytes(indices.buffers[1].Bytes())
		min, max := utils.GetMinMaxInt16(data[start:end])
		if min < 0 || max >= int16(upperlimit) {
			return fmt.Errorf("contains out of bounds index: min: %d, max: %d", min, max)
		}
	case arrow.UINT16:
		data := arrow.Uint16Traits.CastFromBytes(indices.buffers[1].Bytes())
		_, max := utils.GetMinMaxUint16(data[start:end])
		if max >= uint16(upperlimit) {
			return fmt.Errorf("contains out of bounds index: max: %d", max)
		}
	case arrow.INT32:
		data := arrow.Int32Traits.CastFromBytes(indices.buffers[1].Bytes())
		min, max := utils.GetMinMaxInt32(data[start:end])
		if min < 0 || max >= int32(upperlimit) {
			return fmt.Errorf("contains out of bounds index: min: %d, max: %d", min, max)
		}
	case arrow.UINT32:
		data := arrow.Uint32Traits.CastFromBytes(indices.buffers[1].Bytes())
		_, max := utils.GetMinMaxUint32(data[start:end])
		if max >= uint32(upperlimit) {
			return fmt.Errorf("contains out of bounds index: max: %d", max)
		}
	case arrow.INT64:
		data := arrow.Int64Traits.CastFromBytes(indices.buffers[1].Bytes())
		min, max := utils.GetMinMaxInt64(data[start:end])
		if min < 0 || max >= int64(upperlimit) {
			return fmt.Errorf("contains out of bounds index: min: %d, max: %d", min, max)
		}
	case arrow.UINT64:
		data := arrow.Uint64Traits.CastFromBytes(indices.buffers[1].Bytes())
		_, max := utils.GetMinMaxUint64(data[indices.offset : indices.offset+indices.length])
		if max >= upperlimit {
			return fmt.Errorf("contains out of bounds value: max: %d", max)
		}
	default:
		return fmt.Errorf("invalid type for bounds checking: %T", indices.dtype)
	}

	return nil
}

// NewValidatedDictionaryArray constructs a dictionary array from the provided indices
// and dictionary arrays, while also performing validation checks to ensure correctness
// such as bounds checking at are usually skipped for performance.
func NewValidatedDictionaryArray(typ *arrow.DictionaryType, indices, dict arrow.Array) (*Dictionary, error) {
	if indices.DataType().ID() != typ.IndexType.ID() {
		return nil, fmt.Errorf("dictionary type index (%T) does not match indices array type (%T)", typ.IndexType, indices.DataType())
	}

	if !arrow.TypeEqual(typ.ValueType, dict.DataType()) {
		return nil, fmt.Errorf("dictionary value type (%T) does not match dict array type (%T)", typ.ValueType, dict.DataType())
	}

	if err := checkIndexBounds(indices.Data().(*Data), uint64(dict.Len())); err != nil {
		return nil, err
	}

	return NewDictionaryArray(typ, indices, dict), nil
}

// NewDictionaryData creates a strongly typed Dictionary array from
// an ArrayData object with a datatype of arrow.Dictionary and a dictionary
func NewDictionaryData(data arrow.ArrayData) *Dictionary {
	a := &Dictionary{}
	a.refCount = 1
	a.setData(data.(*Data))
	return a
}

func (d *Dictionary) Retain() {
	atomic.AddInt64(&d.refCount, 1)
}

func (d *Dictionary) Release() {
	debug.Assert(atomic.LoadInt64(&d.refCount) > 0, "too many releases")

	if atomic.AddInt64(&d.refCount, -1) == 0 {
		d.data.Release()
		d.data, d.nullBitmapBytes = nil, nil
		d.indices.Release()
		d.indices = nil
		if d.dict != nil {
			d.dict.Release()
			d.dict = nil
		}
	}
}

func (d *Dictionary) setData(data *Data) {
	d.array.setData(data)

	dictType := data.dtype.(*arrow.DictionaryType)
	if data.dictionary == nil {
		if data.length > 0 {
			panic("arrow/array: no dictionary set in Data for Dictionary array")
		}
	} else {
		debug.Assert(arrow.TypeEqual(dictType.ValueType, data.dictionary.DataType()), "mismatched dictionary value types")
	}

	indexData := NewData(dictType.IndexType, data.length, data.buffers, data.childData, data.nulls, data.offset)
	defer indexData.Release()
	d.indices = MakeFromData(indexData)
}

// Dictionary returns the values array that makes up the dictionary for this
// array.
func (d *Dictionary) Dictionary() arrow.Array {
	if d.dict == nil {
		d.dict = MakeFromData(d.data.dictionary)
	}
	return d.dict
}

// Indices returns the underlying array of indices as it's own array
func (d *Dictionary) Indices() arrow.Array {
	return d.indices
}

// CanCompareIndices returns true if the dictionary arrays can be compared
// without having to unify the dictionaries themselves first.
// This means that the index types are equal too.
func (d *Dictionary) CanCompareIndices(other *Dictionary) bool {
	if !arrow.TypeEqual(d.indices.DataType(), other.indices.DataType()) {
		return false
	}

	minlen := int64(min(d.data.dictionary.length, other.data.dictionary.length))
	return SliceEqual(d.Dictionary(), 0, minlen, other.Dictionary(), 0, minlen)
}

func (d *Dictionary) ValueStr(i int) string {
	if d.IsNull(i) {
		return NullValueStr
	}
	return d.Dictionary().ValueStr(d.GetValueIndex(i))
}

func (d *Dictionary) String() string {
	return fmt.Sprintf("{ dictionary: %v\n  indices: %v }", d.Dictionary(), d.Indices())
}

// GetValueIndex returns the dictionary index for the value at index i of the array.
// The actual value can be retrieved by using d.Dictionary().(valuetype).Value(d.GetValueIndex(i))
func (d *Dictionary) GetValueIndex(i int) int {
	indiceData := d.data.buffers[1].Bytes()
	// we know the value is non-negative per the spec, so
	// we can use the unsigned value regardless.
	switch d.indices.DataType().ID() {
	case arrow.UINT8, arrow.INT8:
		return int(uint8(indiceData[d.data.offset+i]))
	case arrow.UINT16, arrow.INT16:
		return int(arrow.Uint16Traits.CastFromBytes(indiceData)[d.data.offset+i])
	case arrow.UINT32, arrow.INT32:
		idx := arrow.Uint32Traits.CastFromBytes(indiceData)[d.data.offset+i]
		debug.Assert(bits.UintSize == 64 || idx <= math.MaxInt32, "arrow/dictionary: truncation of index value")
		return int(idx)
	case arrow.UINT64, arrow.INT64:
		idx := arrow.Uint64Traits.CastFromBytes(indiceData)[d.data.offset+i]
		debug.Assert((bits.UintSize == 32 && idx <= math.MaxInt32) || (bits.UintSize == 64 && idx <= math.MaxInt64), "arrow/dictionary: truncation of index value")
		return int(idx)
	}
	debug.Assert(false, "unreachable dictionary index")
	return -1
}

func (d *Dictionary) GetOneForMarshal(i int) interface{} {
	if d.IsNull(i) {
		return nil
	}
	vidx := d.GetValueIndex(i)
	return d.Dictionary().GetOneForMarshal(vidx)
}

func (d *Dictionary) MarshalJSON() ([]byte, error) {
	vals := make([]interface{}, d.Len())
	for i := 0; i < d.Len(); i++ {
		vals[i] = d.GetOneForMarshal(i)
	}
	return json.Marshal(vals)
}

func arrayEqualDict(l, r *Dictionary) bool {
	return Equal(l.Dictionary(), r.Dictionary()) && Equal(l.indices, r.indices)
}

func arrayApproxEqualDict(l, r *Dictionary, opt equalOption) bool {
	return arrayApproxEqual(l.Dictionary(), r.Dictionary(), opt) && arrayApproxEqual(l.indices, r.indices, opt)
}

// helper for building the properly typed indices of the dictionary builder
type IndexBuilder struct {
	Builder
	Append func(int)
}

func createIndexBuilder(mem memory.Allocator, dt arrow.FixedWidthDataType) (ret IndexBuilder, err error) {
	ret = IndexBuilder{Builder: NewBuilder(mem, dt)}
	switch dt.ID() {
	case arrow.INT8:
		ret.Append = func(idx int) {
			ret.Builder.(*Int8Builder).Append(int8(idx))
		}
	case arrow.UINT8:
		ret.Append = func(idx int) {
			ret.Builder.(*Uint8Builder).Append(uint8(idx))
		}
	case arrow.INT16:
		ret.Append = func(idx int) {
			ret.Builder.(*Int16Builder).Append(int16(idx))
		}
	case arrow.UINT16:
		ret.Append = func(idx int) {
			ret.Builder.(*Uint16Builder).Append(uint16(idx))
		}
	case arrow.INT32:
		ret.Append = func(idx int) {
			ret.Builder.(*Int32Builder).Append(int32(idx))
		}
	case arrow.UINT32:
		ret.Append = func(idx int) {
			ret.Builder.(*Uint32Builder).Append(uint32(idx))
		}
	case arrow.INT64:
		ret.Append = func(idx int) {
			ret.Builder.(*Int64Builder).Append(int64(idx))
		}
	case arrow.UINT64:
		ret.Append = func(idx int) {
			ret.Builder.(*Uint64Builder).Append(uint64(idx))
		}
	default:
		debug.Assert(false, "dictionary index type must be integral")
		err = fmt.Errorf("dictionary index type must be integral, not %s", dt)
	}

	return
}

// helper function to construct an appropriately typed memo table based on
// the value type for the dictionary
func createMemoTable(mem memory.Allocator, dt arrow.DataType) (ret hashing.MemoTable, err error) {
	switch dt.ID() {
	case arrow.INT8:
		ret = hashing.NewInt8MemoTable(0)
	case arrow.UINT8:
		ret = hashing.NewUint8MemoTable(0)
	case arrow.INT16:
		ret = hashing.NewInt16MemoTable(0)
	case arrow.UINT16:
		ret = hashing.NewUint16MemoTable(0)
	case arrow.INT32:
		ret = hashing.NewInt32MemoTable(0)
	case arrow.UINT32:
		ret = hashing.NewUint32MemoTable(0)
	case arrow.INT64:
		ret = hashing.NewInt64MemoTable(0)
	case arrow.UINT64:
		ret = hashing.NewUint64MemoTable(0)
	case arrow.DURATION, arrow.TIMESTAMP, arrow.DATE64, arrow.TIME64:
		ret = hashing.NewInt64MemoTable(0)
	case arrow.TIME32, arrow.DATE32, arrow.INTERVAL_MONTHS:
		ret = hashing.NewInt32MemoTable(0)
	case arrow.FLOAT16:
		ret = hashing.NewUint16MemoTable(0)
	case arrow.FLOAT32:
		ret = hashing.NewFloat32MemoTable(0)
	case arrow.FLOAT64:
		ret = hashing.NewFloat64MemoTable(0)
	case arrow.BINARY, arrow.FIXED_SIZE_BINARY, arrow.DECIMAL128, arrow.DECIMAL256, arrow.INTERVAL_DAY_TIME, arrow.INTERVAL_MONTH_DAY_NANO:
		ret = hashing.NewBinaryMemoTable(0, 0, NewBinaryBuilder(mem, arrow.BinaryTypes.Binary))
	case arrow.STRING:
		ret = hashing.NewBinaryMemoTable(0, 0, NewBinaryBuilder(mem, arrow.BinaryTypes.String))
	case arrow.NULL:
	default:
		err = fmt.Errorf("unimplemented dictionary value type, %s", dt)
	}

	return
}

type DictionaryBuilder interface {
	Builder

	NewDictionaryArray() *Dictionary
	NewDelta() (indices, delta arrow.Array, err error)
	AppendArray(arrow.Array) error
	AppendIndices([]int, []bool)
	ResetFull()
}

type dictionaryBuilder struct {
	builder

	dt          *arrow.DictionaryType
	deltaOffset int
	memoTable   hashing.MemoTable
	idxBuilder  IndexBuilder
}

// NewDictionaryBuilderWithDict initializes a dictionary builder and inserts the values from `init` as the first
// values in the dictionary, but does not insert them as values into the array.
func NewDictionaryBuilderWithDict(mem memory.Allocator, dt *arrow.DictionaryType, init arrow.Array) DictionaryBuilder {
	if init != nil && !arrow.TypeEqual(dt.ValueType, init.DataType()) {
		panic(fmt.Errorf("arrow/array: cannot initialize dictionary type %T with array of type %T", dt.ValueType, init.DataType()))
	}

	idxbldr, err := createIndexBuilder(mem, dt.IndexType.(arrow.FixedWidthDataType))
	if err != nil {
		panic(fmt.Errorf("arrow/array: unsupported builder for index type of %T", dt))
	}

	memo, err := createMemoTable(mem, dt.ValueType)
	if err != nil {
		panic(fmt.Errorf("arrow/array: unsupported builder for value type of %T", dt))
	}

	bldr := dictionaryBuilder{
		builder:    builder{refCount: 1, mem: mem},
		idxBuilder: idxbldr,
		memoTable:  memo,
		dt:         dt,
	}

	switch dt.ValueType.ID() {
	case arrow.NULL:
		ret := &NullDictionaryBuilder{bldr}
		debug.Assert(init == nil, "arrow/array: doesn't make sense to init a null dictionary")
		return ret
	case arrow.UINT8:
		ret := &Uint8DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Uint8)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.INT8:
		ret := &Int8DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Int8)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.UINT16:
		ret := &Uint16DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Uint16)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.INT16:
		ret := &Int16DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Int16)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.UINT32:
		ret := &Uint32DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Uint32)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.INT32:
		ret := &Int32DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Int32)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.UINT64:
		ret := &Uint64DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Uint64)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.INT64:
		ret := &Int64DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Int64)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.FLOAT16:
		ret := &Float16DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Float16)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.FLOAT32:
		ret := &Float32DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Float32)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.FLOAT64:
		ret := &Float64DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Float64)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.STRING:
		ret := &BinaryDictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertStringDictValues(init.(*String)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.BINARY:
		ret := &BinaryDictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Binary)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.FIXED_SIZE_BINARY:
		ret := &FixedSizeBinaryDictionaryBuilder{
			bldr, dt.ValueType.(*arrow.FixedSizeBinaryType).ByteWidth,
		}
		if init != nil {
			if err = ret.InsertDictValues(init.(*FixedSizeBinary)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.DATE32:
		ret := &Date32DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Date32)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.DATE64:
		ret := &Date64DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Date64)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.TIMESTAMP:
		ret := &TimestampDictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Timestamp)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.TIME32:
		ret := &Time32DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Time32)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.TIME64:
		ret := &Time64DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Time64)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.INTERVAL_MONTHS:
		ret := &MonthIntervalDictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*MonthInterval)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.INTERVAL_DAY_TIME:
		ret := &DayTimeDictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*DayTimeInterval)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.DECIMAL128:
		ret := &Decimal128DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Decimal128)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.DECIMAL256:
		ret := &Decimal256DictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Decimal256)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.LIST:
	case arrow.STRUCT:
	case arrow.SPARSE_UNION:
	case arrow.DENSE_UNION:
	case arrow.DICTIONARY:
	case arrow.MAP:
	case arrow.EXTENSION:
	case arrow.FIXED_SIZE_LIST:
	case arrow.DURATION:
		ret := &DurationDictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*Duration)); err != nil {
				panic(err)
			}
		}
		return ret
	case arrow.LARGE_STRING:
	case arrow.LARGE_BINARY:
	case arrow.LARGE_LIST:
	case arrow.INTERVAL_MONTH_DAY_NANO:
		ret := &MonthDayNanoDictionaryBuilder{bldr}
		if init != nil {
			if err = ret.InsertDictValues(init.(*MonthDayNanoInterval)); err != nil {
				panic(err)
			}
		}
		return ret
	}

	panic("arrow/array: unimplemented dictionary key type")
}

func NewDictionaryBuilder(mem memory.Allocator, dt *arrow.DictionaryType) DictionaryBuilder {
	return NewDictionaryBuilderWithDict(mem, dt, nil)
}

func (b *dictionaryBuilder) Type() arrow.DataType { return b.dt }

func (b *dictionaryBuilder) Release() {
	debug.Assert(atomic.LoadInt64(&b.refCount) > 0, "too many releases")

	if atomic.AddInt64(&b.refCount, -1) == 0 {
		b.idxBuilder.Release()
		b.idxBuilder.Builder = nil
		if binmemo, ok := b.memoTable.(*hashing.BinaryMemoTable); ok {
			binmemo.Release()
		}
		b.memoTable = nil
	}
}

func (b *dictionaryBuilder) AppendNull() {
	b.length += 1
	b.nulls += 1
	b.idxBuilder.AppendNull()
}

func (b *dictionaryBuilder) AppendNulls(n int) {
	for i := 0; i < n; i++ {
		b.AppendNull()
	}
}

func (b *dictionaryBuilder) AppendEmptyValue() {
	b.length += 1
	b.idxBuilder.AppendEmptyValue()
}

func (b *dictionaryBuilder) AppendEmptyValues(n int) {
	for i := 0; i < n; i++ {
		b.AppendEmptyValue()
	}
}

func (b *dictionaryBuilder) Reserve(n int) {
	b.idxBuilder.Reserve(n)
}

func (b *dictionaryBuilder) Resize(n int) {
	b.idxBuilder.Resize(n)
	b.length = b.idxBuilder.Len()
}

func (b *dictionaryBuilder) ResetFull() {
	b.builder.reset()
	b.idxBuilder.NewArray().Release()
	b.memoTable.Reset()
}

func (b *dictionaryBuilder) Cap() int { return b.idxBuilder.Cap() }

func (b *dictionaryBuilder) IsNull(i int) bool { return b.idxBuilder.IsNull(i) }

func (b *dictionaryBuilder) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("dictionary builder must unpack from json array, found %s", delim)
	}

	return b.Unmarshal(dec)
}

func (b *dictionaryBuilder) Unmarshal(dec *json.Decoder) error {
	bldr := NewBuilder(b.mem, b.dt.ValueType)
	defer bldr.Release()

	if err := bldr.Unmarshal(dec); err != nil {
		return err
	}

	arr := bldr.NewArray()
	defer arr.Release()
	return b.AppendArray(arr)
}

func (b *dictionaryBuilder) AppendValueFromString(s string) error {
	bldr := NewBuilder(b.mem, b.dt.ValueType)
	defer bldr.Release()

	if err := bldr.AppendValueFromString(s); err != nil {
		return err
	}

	arr := bldr.NewArray()
	defer arr.Release()
	return b.AppendArray(arr)
}

func (b *dictionaryBuilder) UnmarshalOne(dec *json.Decoder) error {
	bldr := NewBuilder(b.mem, b.dt.ValueType)
	defer bldr.Release()

	if err := bldr.UnmarshalOne(dec); err != nil {
		return err
	}

	arr := bldr.NewArray()
	defer arr.Release()
	return b.AppendArray(arr)
}

func (b *dictionaryBuilder) NewArray() arrow.Array {
	return b.NewDictionaryArray()
}

func (b *dictionaryBuilder) newData() *Data {
	indices, dict, err := b.newWithDictOffset(0)
	if err != nil {
		panic(err)
	}

	indices.dtype = b.dt
	indices.dictionary = dict
	return indices
}

func (b *dictionaryBuilder) NewDictionaryArray() *Dictionary {
	a := &Dictionary{}
	a.refCount = 1

	indices := b.newData()
	a.setData(indices)
	indices.Release()
	return a
}

func (b *dictionaryBuilder) newWithDictOffset(offset int) (indices, dict *Data, err error) {
	idxarr := b.idxBuilder.NewArray()
	defer idxarr.Release()

	indices = idxarr.Data().(*Data)

	b.deltaOffset = b.memoTable.Size()
	dict, err = GetDictArrayData(b.mem, b.dt.ValueType, b.memoTable, offset)
	b.reset()
	indices.Retain()
	return
}

// NewDelta returns the dictionary indices and a delta dictionary since the
// last time NewArray or NewDictionaryArray were called, and resets the state
// of the builder (except for the dictionary / memotable)
func (b *dictionaryBuilder) NewDelta() (indices, delta arrow.Array, err error) {
	indicesData, deltaData, err := b.newWithDictOffset(b.deltaOffset)
	if err != nil {
		return nil, nil, err
	}

	defer indicesData.Release()
	defer deltaData.Release()
	indices, delta = MakeFromData(indicesData), MakeFromData(deltaData)
	return
}

func (b *dictionaryBuilder) insertDictValue(val interface{}) error {
	_, _, err := b.memoTable.GetOrInsert(val)
	return err
}

func (b *dictionaryBuilder) insertDictBytes(val []byte) error {
	_, _, err := b.memoTable.GetOrInsertBytes(val)
	return err
}

func (b *dictionaryBuilder) appendValue(val interface{}) error {
	idx, _, err := b.memoTable.GetOrInsert(val)
	b.idxBuilder.Append(idx)
	b.length += 1
	return err
}

func (b *dictionaryBuilder) appendBytes(val []byte) error {
	idx, _, err := b.memoTable.GetOrInsertBytes(val)
	b.idxBuilder.Append(idx)
	b.length += 1
	return err
}

func getvalFn(arr arrow.Array) func(i int) interface{} {
	switch typedarr := arr.(type) {
	case *Int8:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Uint8:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Int16:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Uint16:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Int32:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Uint32:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Int64:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Uint64:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Float16:
		return func(i int) interface{} { return typedarr.Value(i).Uint16() }
	case *Float32:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Float64:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Duration:
		return func(i int) interface{} { return int64(typedarr.Value(i)) }
	case *Timestamp:
		return func(i int) interface{} { return int64(typedarr.Value(i)) }
	case *Date64:
		return func(i int) interface{} { return int64(typedarr.Value(i)) }
	case *Time64:
		return func(i int) interface{} { return int64(typedarr.Value(i)) }
	case *Time32:
		return func(i int) interface{} { return int32(typedarr.Value(i)) }
	case *Date32:
		return func(i int) interface{} { return int32(typedarr.Value(i)) }
	case *MonthInterval:
		return func(i int) interface{} { return int32(typedarr.Value(i)) }
	case *Binary:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *FixedSizeBinary:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *String:
		return func(i int) interface{} { return typedarr.Value(i) }
	case *Decimal128:
		return func(i int) interface{} {
			val := typedarr.Value(i)
			return (*(*[arrow.Decimal128SizeBytes]byte)(unsafe.Pointer(&val)))[:]
		}
	case *Decimal256:
		return func(i int) interface{} {
			val := typedarr.Value(i)
			return (*(*[arrow.Decimal256SizeBytes]byte)(unsafe.Pointer(&val)))[:]
		}
	case *DayTimeInterval:
		return func(i int) interface{} {
			val := typedarr.Value(i)
			return (*(*[arrow.DayTimeIntervalSizeBytes]byte)(unsafe.Pointer(&val)))[:]
		}
	case *MonthDayNanoInterval:
		return func(i int) interface{} {
			val := typedarr.Value(i)
			return (*(*[arrow.MonthDayNanoIntervalSizeBytes]byte)(unsafe.Pointer(&val)))[:]
		}
	}

	panic("arrow/array: invalid dictionary value type")
}

func (b *dictionaryBuilder) AppendArray(arr arrow.Array) error {
	debug.Assert(arrow.TypeEqual(b.dt.ValueType, arr.DataType()), "wrong value type of array to append to dict")

	valfn := getvalFn(arr)
	for i := 0; i < arr.Len(); i++ {
		if arr.IsNull(i) {
			b.AppendNull()
		} else {
			if err := b.appendValue(valfn(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *dictionaryBuilder) IndexBuilder() IndexBuilder {
	return b.idxBuilder
}

func (b *dictionaryBuilder) AppendIndices(indices []int, valid []bool) {
	b.length += len(indices)
	switch idxbldr := b.idxBuilder.Builder.(type) {
	case *Int8Builder:
		vals := make([]int8, len(indices))
		for i, v := range indices {
			vals[i] = int8(v)
		}
		idxbldr.AppendValues(vals, valid)
	case *Int16Builder:
		vals := make([]int16, len(indices))
		for i, v := range indices {
			vals[i] = int16(v)
		}
		idxbldr.AppendValues(vals, valid)
	case *Int32Builder:
		vals := make([]int32, len(indices))
		for i, v := range indices {
			vals[i] = int32(v)
		}
		idxbldr.AppendValues(vals, valid)
	case *Int64Builder:
		vals := make([]int64, len(indices))
		for i, v := range indices {
			vals[i] = int64(v)
		}
		idxbldr.AppendValues(vals, valid)
	case *Uint8Builder:
		vals := make([]uint8, len(indices))
		for i, v := range indices {
			vals[i] = uint8(v)
		}
		idxbldr.AppendValues(vals, valid)
	case *Uint16Builder:
		vals := make([]uint16, len(indices))
		for i, v := range indices {
			vals[i] = uint16(v)
		}
		idxbldr.AppendValues(vals, valid)
	case *Uint32Builder:
		vals := make([]uint32, len(indices))
		for i, v := range indices {
			vals[i] = uint32(v)
		}
		idxbldr.AppendValues(vals, valid)
	case *Uint64Builder:
		vals := make([]uint64, len(indices))
		for i, v := range indices {
			vals[i] = uint64(v)
		}
		idxbldr.AppendValues(vals, valid)
	}
}

type NullDictionaryBuilder struct {
	dictionaryBuilder
}

func (b *NullDictionaryBuilder) NewArray() arrow.Array {
	return b.NewDictionaryArray()
}

func (b *NullDictionaryBuilder) NewDictionaryArray() *Dictionary {
	idxarr := b.idxBuilder.NewArray()
	defer idxarr.Release()

	out := idxarr.Data().(*Data)
	dictarr := NewNull(0)
	defer dictarr.Release()

	dictarr.data.Retain()
	out.dtype = b.dt
	out.dictionary = dictarr.data

	return NewDictionaryData(out)
}

func (b *NullDictionaryBuilder) AppendArray(arr arrow.Array) error {
	if arr.DataType().ID() != arrow.NULL {
		return fmt.Errorf("cannot append non-null array to null dictionary")
	}

	for i := 0; i < arr.(*Null).Len(); i++ {
		b.AppendNull()
	}
	return nil
}

type Int8DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Int8DictionaryBuilder) Append(v int8) error { return b.appendValue(v) }
func (b *Int8DictionaryBuilder) InsertDictValues(arr *Int8) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type Uint8DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Uint8DictionaryBuilder) Append(v uint8) error { return b.appendValue(v) }
func (b *Uint8DictionaryBuilder) InsertDictValues(arr *Uint8) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type Int16DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Int16DictionaryBuilder) Append(v int16) error { return b.appendValue(v) }
func (b *Int16DictionaryBuilder) InsertDictValues(arr *Int16) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type Uint16DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Uint16DictionaryBuilder) Append(v uint16) error { return b.appendValue(v) }
func (b *Uint16DictionaryBuilder) InsertDictValues(arr *Uint16) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type Int32DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Int32DictionaryBuilder) Append(v int32) error { return b.appendValue(v) }
func (b *Int32DictionaryBuilder) InsertDictValues(arr *Int32) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type Uint32DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Uint32DictionaryBuilder) Append(v uint32) error { return b.appendValue(v) }
func (b *Uint32DictionaryBuilder) InsertDictValues(arr *Uint32) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type Int64DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Int64DictionaryBuilder) Append(v int64) error { return b.appendValue(v) }
func (b *Int64DictionaryBuilder) InsertDictValues(arr *Int64) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type Uint64DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Uint64DictionaryBuilder) Append(v uint64) error { return b.appendValue(v) }
func (b *Uint64DictionaryBuilder) InsertDictValues(arr *Uint64) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type DurationDictionaryBuilder struct {
	dictionaryBuilder
}

func (b *DurationDictionaryBuilder) Append(v arrow.Duration) error { return b.appendValue(int64(v)) }
func (b *DurationDictionaryBuilder) InsertDictValues(arr *Duration) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(int64(v)); err != nil {
			break
		}
	}
	return
}

type TimestampDictionaryBuilder struct {
	dictionaryBuilder
}

func (b *TimestampDictionaryBuilder) Append(v arrow.Timestamp) error { return b.appendValue(int64(v)) }
func (b *TimestampDictionaryBuilder) InsertDictValues(arr *Timestamp) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(int64(v)); err != nil {
			break
		}
	}
	return
}

type Time32DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Time32DictionaryBuilder) Append(v arrow.Time32) error { return b.appendValue(int32(v)) }
func (b *Time32DictionaryBuilder) InsertDictValues(arr *Time32) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(int32(v)); err != nil {
			break
		}
	}
	return
}

type Time64DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Time64DictionaryBuilder) Append(v arrow.Time64) error { return b.appendValue(int64(v)) }
func (b *Time64DictionaryBuilder) InsertDictValues(arr *Time64) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(int64(v)); err != nil {
			break
		}
	}
	return
}

type Date32DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Date32DictionaryBuilder) Append(v arrow.Date32) error { return b.appendValue(int32(v)) }
func (b *Date32DictionaryBuilder) InsertDictValues(arr *Date32) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(int32(v)); err != nil {
			break
		}
	}
	return
}

type Date64DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Date64DictionaryBuilder) Append(v arrow.Date64) error { return b.appendValue(int64(v)) }
func (b *Date64DictionaryBuilder) InsertDictValues(arr *Date64) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(int64(v)); err != nil {
			break
		}
	}
	return
}

type MonthIntervalDictionaryBuilder struct {
	dictionaryBuilder
}

func (b *MonthIntervalDictionaryBuilder) Append(v arrow.MonthInterval) error {
	return b.appendValue(int32(v))
}
func (b *MonthIntervalDictionaryBuilder) InsertDictValues(arr *MonthInterval) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(int32(v)); err != nil {
			break
		}
	}
	return
}

type Float16DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Float16DictionaryBuilder) Append(v float16.Num) error { return b.appendValue(v.Uint16()) }
func (b *Float16DictionaryBuilder) InsertDictValues(arr *Float16) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v.Uint16()); err != nil {
			break
		}
	}
	return
}

type Float32DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Float32DictionaryBuilder) Append(v float32) error { return b.appendValue(v) }
func (b *Float32DictionaryBuilder) InsertDictValues(arr *Float32) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type Float64DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Float64DictionaryBuilder) Append(v float64) error { return b.appendValue(v) }
func (b *Float64DictionaryBuilder) InsertDictValues(arr *Float64) (err error) {
	for _, v := range arr.values {
		if err = b.insertDictValue(v); err != nil {
			break
		}
	}
	return
}

type BinaryDictionaryBuilder struct {
	dictionaryBuilder
}

func (b *BinaryDictionaryBuilder) Append(v []byte) error {
	if v == nil {
		b.AppendNull()
		return nil
	}

	return b.appendBytes(v)
}

func (b *BinaryDictionaryBuilder) AppendString(v string) error { return b.appendBytes([]byte(v)) }
func (b *BinaryDictionaryBuilder) InsertDictValues(arr *Binary) (err error) {
	if !arrow.TypeEqual(arr.DataType(), b.dt.ValueType) {
		return fmt.Errorf("dictionary insert type mismatch: cannot insert values of type %T to dictionary type %T", arr.DataType(), b.dt.ValueType)
	}

	for i := 0; i < arr.Len(); i++ {
		if err = b.insertDictBytes(arr.Value(i)); err != nil {
			break
		}
	}
	return
}
func (b *BinaryDictionaryBuilder) InsertStringDictValues(arr *String) (err error) {
	if !arrow.TypeEqual(arr.DataType(), b.dt.ValueType) {
		return fmt.Errorf("dictionary insert type mismatch: cannot insert values of type %T to dictionary type %T", arr.DataType(), b.dt.ValueType)
	}

	for i := 0; i < arr.Len(); i++ {
		if err = b.insertDictValue(arr.Value(i)); err != nil {
			break
		}
	}
	return
}

func (b *BinaryDictionaryBuilder) GetValueIndex(i int) int {
	switch b := b.idxBuilder.Builder.(type) {
	case *Uint8Builder:
		return int(b.Value(i))
	case *Int8Builder:
		return int(b.Value(i))
	case *Uint16Builder:
		return int(b.Value(i))
	case *Int16Builder:
		return int(b.Value(i))
	case *Uint32Builder:
		return int(b.Value(i))
	case *Int32Builder:
		return int(b.Value(i))
	case *Uint64Builder:
		return int(b.Value(i))
	case *Int64Builder:
		return int(b.Value(i))
	default:
		return -1
	}
}

func (b *BinaryDictionaryBuilder) Value(i int) []byte {
	switch mt := b.memoTable.(type) {
	case *hashing.BinaryMemoTable:
		return mt.Value(i)
	}
	return nil
}

func (b *BinaryDictionaryBuilder) ValueStr(i int) string {
	return string(b.Value(i))
}

type FixedSizeBinaryDictionaryBuilder struct {
	dictionaryBuilder
	byteWidth int
}

func (b *FixedSizeBinaryDictionaryBuilder) Append(v []byte) error {
	return b.appendValue(v[:b.byteWidth])
}
func (b *FixedSizeBinaryDictionaryBuilder) InsertDictValues(arr *FixedSizeBinary) (err error) {
	var (
		beg = arr.array.data.offset * b.byteWidth
		end = (arr.array.data.offset + arr.data.length) * b.byteWidth
	)
	data := arr.valueBytes[beg:end]
	for len(data) > 0 {
		if err = b.insertDictValue(data[:b.byteWidth]); err != nil {
			break
		}
		data = data[b.byteWidth:]
	}
	return
}

type Decimal128DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Decimal128DictionaryBuilder) Append(v decimal128.Num) error {
	return b.appendValue((*(*[arrow.Decimal128SizeBytes]byte)(unsafe.Pointer(&v)))[:])
}
func (b *Decimal128DictionaryBuilder) InsertDictValues(arr *Decimal128) (err error) {
	data := arrow.Decimal128Traits.CastToBytes(arr.values)
	for len(data) > 0 {
		if err = b.insertDictValue(data[:arrow.Decimal128SizeBytes]); err != nil {
			break
		}
		data = data[arrow.Decimal128SizeBytes:]
	}
	return
}

type Decimal256DictionaryBuilder struct {
	dictionaryBuilder
}

func (b *Decimal256DictionaryBuilder) Append(v decimal256.Num) error {
	return b.appendValue((*(*[arrow.Decimal256SizeBytes]byte)(unsafe.Pointer(&v)))[:])
}
func (b *Decimal256DictionaryBuilder) InsertDictValues(arr *Decimal256) (err error) {
	data := arrow.Decimal256Traits.CastToBytes(arr.values)
	for len(data) > 0 {
		if err = b.insertDictValue(data[:arrow.Decimal256SizeBytes]); err != nil {
			break
		}
		data = data[arrow.Decimal256SizeBytes:]
	}
	return
}

type MonthDayNanoDictionaryBuilder struct {
	dictionaryBuilder
}

func (b *MonthDayNanoDictionaryBuilder) Append(v arrow.MonthDayNanoInterval) error {
	return b.appendValue((*(*[arrow.MonthDayNanoIntervalSizeBytes]byte)(unsafe.Pointer(&v)))[:])
}
func (b *MonthDayNanoDictionaryBuilder) InsertDictValues(arr *MonthDayNanoInterval) (err error) {
	data := arrow.MonthDayNanoIntervalTraits.CastToBytes(arr.values)
	for len(data) > 0 {
		if err = b.insertDictValue(data[:arrow.MonthDayNanoIntervalSizeBytes]); err != nil {
			break
		}
		data = data[arrow.MonthDayNanoIntervalSizeBytes:]
	}
	return
}

type DayTimeDictionaryBuilder struct {
	dictionaryBuilder
}

func (b *DayTimeDictionaryBuilder) Append(v arrow.DayTimeInterval) error {
	return b.appendValue((*(*[arrow.DayTimeIntervalSizeBytes]byte)(unsafe.Pointer(&v)))[:])
}
func (b *DayTimeDictionaryBuilder) InsertDictValues(arr *DayTimeInterval) (err error) {
	data := arrow.DayTimeIntervalTraits.CastToBytes(arr.values)
	for len(data) > 0 {
		if err = b.insertDictValue(data[:arrow.DayTimeIntervalSizeBytes]); err != nil {
			break
		}
		data = data[arrow.DayTimeIntervalSizeBytes:]
	}
	return
}

func IsTrivialTransposition(transposeMap []int32) bool {
	for i, t := range transposeMap {
		if t != int32(i) {
			return false
		}
	}
	return true
}

func TransposeDictIndices(mem memory.Allocator, data arrow.ArrayData, inType, outType arrow.DataType, dict arrow.ArrayData, transposeMap []int32) (arrow.ArrayData, error) {
	// inType may be different from data->dtype if data is ExtensionType
	if inType.ID() != arrow.DICTIONARY || outType.ID() != arrow.DICTIONARY {
		return nil, errors.New("arrow/array: expected dictionary type")
	}

	var (
		inDictType   = inType.(*arrow.DictionaryType)
		outDictType  = outType.(*arrow.DictionaryType)
		inIndexType  = inDictType.IndexType
		outIndexType = outDictType.IndexType.(arrow.FixedWidthDataType)
	)

	if inIndexType.ID() == outIndexType.ID() && IsTrivialTransposition(transposeMap) {
		// index type and values will be identical, we can reuse the existing buffers
		return NewDataWithDictionary(outType, data.Len(), []*memory.Buffer{data.Buffers()[0], data.Buffers()[1]},
			data.NullN(), data.Offset(), dict.(*Data)), nil
	}

	// default path: compute the transposed indices as a new buffer
	outBuf := memory.NewResizableBuffer(mem)
	outBuf.Resize(data.Len() * int(bitutil.BytesForBits(int64(outIndexType.BitWidth()))))
	defer outBuf.Release()

	// shift null buffer if original offset is non-zero
	var nullBitmap *memory.Buffer
	if data.Offset() != 0 && data.NullN() != 0 {
		nullBitmap = memory.NewResizableBuffer(mem)
		nullBitmap.Resize(int(bitutil.BytesForBits(int64(data.Len()))))
		bitutil.CopyBitmap(data.Buffers()[0].Bytes(), data.Offset(), data.Len(), nullBitmap.Bytes(), 0)
		defer nullBitmap.Release()
	} else {
		nullBitmap = data.Buffers()[0]
	}

	outData := NewDataWithDictionary(outType, data.Len(),
		[]*memory.Buffer{nullBitmap, outBuf}, data.NullN(), 0, dict.(*Data))
	err := utils.TransposeIntsBuffers(inIndexType, outIndexType,
		data.Buffers()[1].Bytes(), outBuf.Bytes(), data.Offset(), outData.offset, data.Len(), transposeMap)
	return outData, err
}

// DictionaryUnifier defines the interface used for unifying, and optionally producing
// transposition maps for, multiple dictionary arrays incrementally.
type DictionaryUnifier interface {
	// Unify adds the provided array of dictionary values to be unified.
	Unify(arrow.Array) error
	// UnifyAndTranspose adds the provided array of dictionary values,
	// just like Unify but returns an allocated buffer containing a mapping
	// to transpose dictionary indices.
	UnifyAndTranspose(dict arrow.Array) (transposed *memory.Buffer, err error)
	// GetResult returns the dictionary type (choosing the smallest index type
	// that can represent all the values) and the new unified dictionary.
	//
	// Calling GetResult clears the existing dictionary from the unifier so it
	// can be reused by calling Unify/UnifyAndTranspose again with new arrays.
	GetResult() (outType arrow.DataType, outDict arrow.Array, err error)
	// GetResultWithIndexType is like GetResult, but allows specifying the type
	// of the dictionary indexes rather than letting the unifier pick. If the
	// passed in index type isn't large enough to represent all of the dictionary
	// values, an error will be returned instead. The new unified dictionary
	// is returned.
	GetResultWithIndexType(indexType arrow.DataType) (arrow.Array, error)
	// Release should be called to clean up any allocated scratch memo-table used
	// for building the unified dictionary.
	Release()
}

type unifier struct {
	mem       memory.Allocator
	valueType arrow.DataType
	memoTable hashing.MemoTable
}

// NewDictionaryUnifier constructs and returns a new dictionary unifier for dictionaries
// of valueType, using the provided allocator for allocating the unified dictionary
// and the memotable used for building it.
//
// This will only work for non-nested types currently. a nested valueType or dictionary type
// will result in an error.
func NewDictionaryUnifier(alloc memory.Allocator, valueType arrow.DataType) (DictionaryUnifier, error) {
	memoTable, err := createMemoTable(alloc, valueType)
	if err != nil {
		return nil, err
	}
	return &unifier{
		mem:       alloc,
		valueType: valueType,
		memoTable: memoTable,
	}, nil
}

func (u *unifier) Release() {
	if bin, ok := u.memoTable.(*hashing.BinaryMemoTable); ok {
		bin.Release()
	}
}

func (u *unifier) Unify(dict arrow.Array) (err error) {
	if !arrow.TypeEqual(u.valueType, dict.DataType()) {
		return fmt.Errorf("dictionary type different from unifier: %s, expected: %s", dict.DataType(), u.valueType)
	}

	valFn := getvalFn(dict)
	for i := 0; i < dict.Len(); i++ {
		if dict.IsNull(i) {
			u.memoTable.GetOrInsertNull()
			continue
		}

		if _, _, err = u.memoTable.GetOrInsert(valFn(i)); err != nil {
			return err
		}
	}
	return
}

func (u *unifier) UnifyAndTranspose(dict arrow.Array) (transposed *memory.Buffer, err error) {
	if !arrow.TypeEqual(u.valueType, dict.DataType()) {
		return nil, fmt.Errorf("dictionary type different from unifier: %s, expected: %s", dict.DataType(), u.valueType)
	}

	transposed = memory.NewResizableBuffer(u.mem)
	transposed.Resize(arrow.Int32Traits.BytesRequired(dict.Len()))

	newIdxes := arrow.Int32Traits.CastFromBytes(transposed.Bytes())
	valFn := getvalFn(dict)
	for i := 0; i < dict.Len(); i++ {
		if dict.IsNull(i) {
			idx, _ := u.memoTable.GetOrInsertNull()
			newIdxes[i] = int32(idx)
			continue
		}

		idx, _, err := u.memoTable.GetOrInsert(valFn(i))
		if err != nil {
			transposed.Release()
			return nil, err
		}
		newIdxes[i] = int32(idx)
	}
	return
}

func (u *unifier) GetResult() (outType arrow.DataType, outDict arrow.Array, err error) {
	dictLen := u.memoTable.Size()
	var indexType arrow.DataType
	switch {
	case dictLen <= math.MaxInt8:
		indexType = arrow.PrimitiveTypes.Int8
	case dictLen <= math.MaxInt16:
		indexType = arrow.PrimitiveTypes.Int16
	case dictLen <= math.MaxInt32:
		indexType = arrow.PrimitiveTypes.Int32
	default:
		indexType = arrow.PrimitiveTypes.Int64
	}
	outType = &arrow.DictionaryType{IndexType: indexType, ValueType: u.valueType}

	dictData, err := GetDictArrayData(u.mem, u.valueType, u.memoTable, 0)
	if err != nil {
		return nil, nil, err
	}

	u.memoTable.Reset()

	defer dictData.Release()
	outDict = MakeFromData(dictData)
	return
}

func (u *unifier) GetResultWithIndexType(indexType arrow.DataType) (arrow.Array, error) {
	dictLen := u.memoTable.Size()
	var toobig bool
	switch indexType.ID() {
	case arrow.UINT8:
		toobig = dictLen > math.MaxUint8
	case arrow.INT8:
		toobig = dictLen > math.MaxInt8
	case arrow.UINT16:
		toobig = dictLen > math.MaxUint16
	case arrow.INT16:
		toobig = dictLen > math.MaxInt16
	case arrow.UINT32:
		toobig = uint(dictLen) > math.MaxUint32
	case arrow.INT32:
		toobig = dictLen > math.MaxInt32
	case arrow.UINT64:
		toobig = uint64(dictLen) > uint64(math.MaxUint64)
	case arrow.INT64:
	default:
		return nil, fmt.Errorf("arrow/array: invalid dictionary index type: %s, must be integral", indexType)
	}
	if toobig {
		return nil, errors.New("arrow/array: cannot combine dictionaries. unified dictionary requires a larger index type")
	}

	dictData, err := GetDictArrayData(u.mem, u.valueType, u.memoTable, 0)
	if err != nil {
		return nil, err
	}

	u.memoTable.Reset()

	defer dictData.Release()
	return MakeFromData(dictData), nil
}

type binaryUnifier struct {
	mem       memory.Allocator
	memoTable *hashing.BinaryMemoTable
}

// NewBinaryDictionaryUnifier constructs and returns a new dictionary unifier for dictionaries
// of binary values, using the provided allocator for allocating the unified dictionary
// and the memotable used for building it.
func NewBinaryDictionaryUnifier(alloc memory.Allocator) DictionaryUnifier {
	return &binaryUnifier{
		mem:       alloc,
		memoTable: hashing.NewBinaryMemoTable(0, 0, NewBinaryBuilder(alloc, arrow.BinaryTypes.Binary)),
	}
}

func (u *binaryUnifier) Release() {
	u.memoTable.Release()
}

func (u *binaryUnifier) Unify(dict arrow.Array) (err error) {
	if !arrow.TypeEqual(arrow.BinaryTypes.Binary, dict.DataType()) {
		return fmt.Errorf("dictionary type different from unifier: %s, expected: %s", dict.DataType(), arrow.BinaryTypes.Binary)
	}

	typedDict := dict.(*Binary)
	for i := 0; i < dict.Len(); i++ {
		if dict.IsNull(i) {
			u.memoTable.GetOrInsertNull()
			continue
		}

		if _, _, err = u.memoTable.GetOrInsertBytes(typedDict.Value(i)); err != nil {
			return err
		}
	}
	return
}

func (u *binaryUnifier) UnifyAndTranspose(dict arrow.Array) (transposed *memory.Buffer, err error) {
	if !arrow.TypeEqual(arrow.BinaryTypes.Binary, dict.DataType()) {
		return nil, fmt.Errorf("dictionary type different from unifier: %s, expected: %s", dict.DataType(), arrow.BinaryTypes.Binary)
	}

	transposed = memory.NewResizableBuffer(u.mem)
	transposed.Resize(arrow.Int32Traits.BytesRequired(dict.Len()))

	newIdxes := arrow.Int32Traits.CastFromBytes(transposed.Bytes())
	typedDict := dict.(*Binary)
	for i := 0; i < dict.Len(); i++ {
		if dict.IsNull(i) {
			idx, _ := u.memoTable.GetOrInsertNull()
			newIdxes[i] = int32(idx)
			continue
		}

		idx, _, err := u.memoTable.GetOrInsertBytes(typedDict.Value(i))
		if err != nil {
			transposed.Release()
			return nil, err
		}
		newIdxes[i] = int32(idx)
	}
	return
}

func (u *binaryUnifier) GetResult() (outType arrow.DataType, outDict arrow.Array, err error) {
	dictLen := u.memoTable.Size()
	var indexType arrow.DataType
	switch {
	case dictLen <= math.MaxInt8:
		indexType = arrow.PrimitiveTypes.Int8
	case dictLen <= math.MaxInt16:
		indexType = arrow.PrimitiveTypes.Int16
	case dictLen <= math.MaxInt32:
		indexType = arrow.PrimitiveTypes.Int32
	default:
		indexType = arrow.PrimitiveTypes.Int64
	}
	outType = &arrow.DictionaryType{IndexType: indexType, ValueType: arrow.BinaryTypes.Binary}

	dictData, err := GetDictArrayData(u.mem, arrow.BinaryTypes.Binary, u.memoTable, 0)
	if err != nil {
		return nil, nil, err
	}

	u.memoTable.Reset()

	defer dictData.Release()
	outDict = MakeFromData(dictData)
	return
}

func (u *binaryUnifier) GetResultWithIndexType(indexType arrow.DataType) (arrow.Array, error) {
	dictLen := u.memoTable.Size()
	var toobig bool
	switch indexType.ID() {
	case arrow.UINT8:
		toobig = dictLen > math.MaxUint8
	case arrow.INT8:
		toobig = dictLen > math.MaxInt8
	case arrow.UINT16:
		toobig = dictLen > math.MaxUint16
	case arrow.INT16:
		toobig = dictLen > math.MaxInt16
	case arrow.UINT32:
		toobig = uint(dictLen) > math.MaxUint32
	case arrow.INT32:
		toobig = dictLen > math.MaxInt32
	case arrow.UINT64:
		toobig = uint64(dictLen) > uint64(math.MaxUint64)
	case arrow.INT64:
	default:
		return nil, fmt.Errorf("arrow/array: invalid dictionary index type: %s, must be integral", indexType)
	}
	if toobig {
		return nil, errors.New("arrow/array: cannot combine dictionaries. unified dictionary requires a larger index type")
	}

	dictData, err := GetDictArrayData(u.mem, arrow.BinaryTypes.Binary, u.memoTable, 0)
	if err != nil {
		return nil, err
	}

	u.memoTable.Reset()

	defer dictData.Release()
	return MakeFromData(dictData), nil
}

func unifyRecursive(mem memory.Allocator, typ arrow.DataType, chunks []*Data) (changed bool, err error) {
	debug.Assert(len(chunks) != 0, "must provide non-zero length chunk slice")
	var extType arrow.DataType

	if typ.ID() == arrow.EXTENSION {
		extType = typ
		typ = typ.(arrow.ExtensionType).StorageType()
	}

	if nestedTyp, ok := typ.(arrow.NestedType); ok {
		children := make([]*Data, len(chunks))
		for i, f := range nestedTyp.Fields() {
			for j, c := range chunks {
				children[j] = c.childData[i].(*Data)
			}

			childChanged, err := unifyRecursive(mem, f.Type, children)
			if err != nil {
				return false, err
			}
			if childChanged {
				// only when unification actually occurs
				for j := range chunks {
					chunks[j].childData[i] = children[j]
				}
				changed = true
			}
		}
	}

	if typ.ID() == arrow.DICTIONARY {
		dictType := typ.(*arrow.DictionaryType)
		var (
			uni     DictionaryUnifier
			newDict arrow.Array
		)
		// unify any nested dictionaries first, but the unifier doesn't support
		// nested dictionaries yet so this would fail.
		uni, err = NewDictionaryUnifier(mem, dictType.ValueType)
		if err != nil {
			return changed, err
		}
		defer uni.Release()
		transposeMaps := make([]*memory.Buffer, len(chunks))
		for i, c := range chunks {
			debug.Assert(c.dictionary != nil, "missing dictionary data for dictionary array")
			arr := MakeFromData(c.dictionary)
			defer arr.Release()
			if transposeMaps[i], err = uni.UnifyAndTranspose(arr); err != nil {
				return
			}
			defer transposeMaps[i].Release()
		}

		if newDict, err = uni.GetResultWithIndexType(dictType.IndexType); err != nil {
			return
		}
		defer newDict.Release()

		for j := range chunks {
			chnk, err := TransposeDictIndices(mem, chunks[j], typ, typ, newDict.Data(), arrow.Int32Traits.CastFromBytes(transposeMaps[j].Bytes()))
			if err != nil {
				return changed, err
			}
			chunks[j].Release()
			chunks[j] = chnk.(*Data)
			if extType != nil {
				chunks[j].dtype = extType
			}
		}
		changed = true
	}

	return
}

// UnifyChunkedDicts takes a chunked array of dictionary type and will unify
// the dictionary across all of the chunks with the returned chunked array
// having all chunks share the same dictionary.
//
// The return from this *must* have Release called on it unless an error is returned
// in which case the *arrow.Chunked will be nil.
//
// If there is 1 or fewer chunks, then nothing is modified and this function will just
// call Retain on the passed in Chunked array (so Release can safely be called on it).
// The same is true if the type of the array is not a dictionary or if no changes are
// needed for all of the chunks to be using the same dictionary.
func UnifyChunkedDicts(alloc memory.Allocator, chnkd *arrow.Chunked) (*arrow.Chunked, error) {
	if len(chnkd.Chunks()) <= 1 {
		chnkd.Retain()
		return chnkd, nil
	}

	chunksData := make([]*Data, len(chnkd.Chunks()))
	for i, c := range chnkd.Chunks() {
		c.Data().Retain()
		chunksData[i] = c.Data().(*Data)
	}
	changed, err := unifyRecursive(alloc, chnkd.DataType(), chunksData)
	if err != nil || !changed {
		for _, c := range chunksData {
			c.Release()
		}
		if err == nil {
			chnkd.Retain()
		} else {
			chnkd = nil
		}
		return chnkd, err
	}

	chunks := make([]arrow.Array, len(chunksData))
	for i, c := range chunksData {
		chunks[i] = MakeFromData(c)
		defer chunks[i].Release()
		c.Release()
	}

	return arrow.NewChunked(chnkd.DataType(), chunks), nil
}

// UnifyTableDicts performs UnifyChunkedDicts on each column of the table so that
// any dictionary column will have the dictionaries of its chunks unified.
//
// The returned Table should always be Release'd unless a non-nil error was returned,
// in which case the table returned will be nil.
func UnifyTableDicts(alloc memory.Allocator, table arrow.Table) (arrow.Table, error) {
	cols := make([]arrow.Column, table.NumCols())
	for i := 0; i < int(table.NumCols()); i++ {
		chnkd, err := UnifyChunkedDicts(alloc, table.Column(i).Data())
		if err != nil {
			return nil, err
		}
		defer chnkd.Release()
		cols[i] = *arrow.NewColumn(table.Schema().Field(i), chnkd)
		defer cols[i].Release()
	}
	return NewTable(table.Schema(), cols, table.NumRows()), nil
}

var (
	_ arrow.Array = (*Dictionary)(nil)
	_ Builder     = (*dictionaryBuilder)(nil)
)
