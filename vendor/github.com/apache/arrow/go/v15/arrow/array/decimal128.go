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
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/bitutil"
	"github.com/apache/arrow/go/v15/arrow/decimal128"
	"github.com/apache/arrow/go/v15/arrow/internal/debug"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/apache/arrow/go/v15/internal/json"
)

// A type which represents an immutable sequence of 128-bit decimal values.
type Decimal128 struct {
	array

	values []decimal128.Num
}

func NewDecimal128Data(data arrow.ArrayData) *Decimal128 {
	a := &Decimal128{}
	a.refCount = 1
	a.setData(data.(*Data))
	return a
}

func (a *Decimal128) Value(i int) decimal128.Num { return a.values[i] }

func (a *Decimal128) ValueStr(i int) string {
	if a.IsNull(i) {
		return NullValueStr
	}
	return a.GetOneForMarshal(i).(string)
}

func (a *Decimal128) Values() []decimal128.Num { return a.values }

func (a *Decimal128) String() string {
	o := new(strings.Builder)
	o.WriteString("[")
	for i := 0; i < a.Len(); i++ {
		if i > 0 {
			fmt.Fprintf(o, " ")
		}
		switch {
		case a.IsNull(i):
			o.WriteString(NullValueStr)
		default:
			fmt.Fprintf(o, "%v", a.Value(i))
		}
	}
	o.WriteString("]")
	return o.String()
}

func (a *Decimal128) setData(data *Data) {
	a.array.setData(data)
	vals := data.buffers[1]
	if vals != nil {
		a.values = arrow.Decimal128Traits.CastFromBytes(vals.Bytes())
		beg := a.array.data.offset
		end := beg + a.array.data.length
		a.values = a.values[beg:end]
	}
}

func (a *Decimal128) GetOneForMarshal(i int) interface{} {
	if a.IsNull(i) {
		return nil
	}

	typ := a.DataType().(*arrow.Decimal128Type)
	f := (&big.Float{}).SetInt(a.Value(i).BigInt())
	f.Quo(f, big.NewFloat(math.Pow10(int(typ.Scale))))
	return f.Text('g', int(typ.Precision))
}

// ["1.23", ]
func (a *Decimal128) MarshalJSON() ([]byte, error) {
	vals := make([]interface{}, a.Len())
	for i := 0; i < a.Len(); i++ {
		vals[i] = a.GetOneForMarshal(i)
	}
	return json.Marshal(vals)
}

func arrayEqualDecimal128(left, right *Decimal128) bool {
	for i := 0; i < left.Len(); i++ {
		if left.IsNull(i) {
			continue
		}
		if left.Value(i) != right.Value(i) {
			return false
		}
	}
	return true
}

type Decimal128Builder struct {
	builder

	dtype   *arrow.Decimal128Type
	data    *memory.Buffer
	rawData []decimal128.Num
}

func NewDecimal128Builder(mem memory.Allocator, dtype *arrow.Decimal128Type) *Decimal128Builder {
	return &Decimal128Builder{
		builder: builder{refCount: 1, mem: mem},
		dtype:   dtype,
	}
}

func (b *Decimal128Builder) Type() arrow.DataType { return b.dtype }

// Release decreases the reference count by 1.
// When the reference count goes to zero, the memory is freed.
func (b *Decimal128Builder) Release() {
	debug.Assert(atomic.LoadInt64(&b.refCount) > 0, "too many releases")

	if atomic.AddInt64(&b.refCount, -1) == 0 {
		if b.nullBitmap != nil {
			b.nullBitmap.Release()
			b.nullBitmap = nil
		}
		if b.data != nil {
			b.data.Release()
			b.data = nil
			b.rawData = nil
		}
	}
}

func (b *Decimal128Builder) Append(v decimal128.Num) {
	b.Reserve(1)
	b.UnsafeAppend(v)
}

func (b *Decimal128Builder) UnsafeAppend(v decimal128.Num) {
	bitutil.SetBit(b.nullBitmap.Bytes(), b.length)
	b.rawData[b.length] = v
	b.length++
}

func (b *Decimal128Builder) AppendNull() {
	b.Reserve(1)
	b.UnsafeAppendBoolToBitmap(false)
}

func (b *Decimal128Builder) AppendNulls(n int) {
	for i := 0; i < n; i++ {
		b.AppendNull()
	}
}

func (b *Decimal128Builder) AppendEmptyValue() {
	b.Append(decimal128.Num{})
}

func (b *Decimal128Builder) AppendEmptyValues(n int) {
	for i := 0; i < n; i++ {
		b.AppendEmptyValue()
	}
}

func (b *Decimal128Builder) UnsafeAppendBoolToBitmap(isValid bool) {
	if isValid {
		bitutil.SetBit(b.nullBitmap.Bytes(), b.length)
	} else {
		b.nulls++
	}
	b.length++
}

// AppendValues will append the values in the v slice. The valid slice determines which values
// in v are valid (not null). The valid slice must either be empty or be equal in length to v. If empty,
// all values in v are appended and considered valid.
func (b *Decimal128Builder) AppendValues(v []decimal128.Num, valid []bool) {
	if len(v) != len(valid) && len(valid) != 0 {
		panic("len(v) != len(valid) && len(valid) != 0")
	}

	if len(v) == 0 {
		return
	}

	b.Reserve(len(v))
	if len(v) > 0 {
		arrow.Decimal128Traits.Copy(b.rawData[b.length:], v)
	}
	b.builder.unsafeAppendBoolsToBitmap(valid, len(v))
}

func (b *Decimal128Builder) init(capacity int) {
	b.builder.init(capacity)

	b.data = memory.NewResizableBuffer(b.mem)
	bytesN := arrow.Decimal128Traits.BytesRequired(capacity)
	b.data.Resize(bytesN)
	b.rawData = arrow.Decimal128Traits.CastFromBytes(b.data.Bytes())
}

// Reserve ensures there is enough space for appending n elements
// by checking the capacity and calling Resize if necessary.
func (b *Decimal128Builder) Reserve(n int) {
	b.builder.reserve(n, b.Resize)
}

// Resize adjusts the space allocated by b to n elements. If n is greater than b.Cap(),
// additional memory will be allocated. If n is smaller, the allocated memory may reduced.
func (b *Decimal128Builder) Resize(n int) {
	nBuilder := n
	if n < minBuilderCapacity {
		n = minBuilderCapacity
	}

	if b.capacity == 0 {
		b.init(n)
	} else {
		b.builder.resize(nBuilder, b.init)
		b.data.Resize(arrow.Decimal128Traits.BytesRequired(n))
		b.rawData = arrow.Decimal128Traits.CastFromBytes(b.data.Bytes())
	}
}

// NewArray creates a Decimal128 array from the memory buffers used by the builder and resets the Decimal128Builder
// so it can be used to build a new array.
func (b *Decimal128Builder) NewArray() arrow.Array {
	return b.NewDecimal128Array()
}

// NewDecimal128Array creates a Decimal128 array from the memory buffers used by the builder and resets the Decimal128Builder
// so it can be used to build a new array.
func (b *Decimal128Builder) NewDecimal128Array() (a *Decimal128) {
	data := b.newData()
	a = NewDecimal128Data(data)
	data.Release()
	return
}

func (b *Decimal128Builder) newData() (data *Data) {
	bytesRequired := arrow.Decimal128Traits.BytesRequired(b.length)
	if bytesRequired > 0 && bytesRequired < b.data.Len() {
		// trim buffers
		b.data.Resize(bytesRequired)
	}
	data = NewData(b.dtype, b.length, []*memory.Buffer{b.nullBitmap, b.data}, nil, b.nulls, 0)
	b.reset()

	if b.data != nil {
		b.data.Release()
		b.data = nil
		b.rawData = nil
	}

	return
}

func (b *Decimal128Builder) AppendValueFromString(s string) error {
	if s == NullValueStr {
		b.AppendNull()
		return nil
	}
	val, err := decimal128.FromString(s, b.dtype.Precision, b.dtype.Scale)
	if err != nil {
		b.AppendNull()
		return err
	}
	b.Append(val)
	return nil
}

func (b *Decimal128Builder) UnmarshalOne(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}

	switch v := t.(type) {
	case float64:
		val, err := decimal128.FromFloat64(v, b.dtype.Precision, b.dtype.Scale)
		if err != nil {
			return err
		}
		b.Append(val)
	case string:
		val, err := decimal128.FromString(v, b.dtype.Precision, b.dtype.Scale)
		if err != nil {
			return err
		}
		b.Append(val)
	case json.Number:
		val, err := decimal128.FromString(v.String(), b.dtype.Precision, b.dtype.Scale)
		if err != nil {
			return err
		}
		b.Append(val)
	case nil:
		b.AppendNull()
		return nil
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprint(t),
			Type:   reflect.TypeOf(decimal128.Num{}),
			Offset: dec.InputOffset(),
		}
	}

	return nil
}

func (b *Decimal128Builder) Unmarshal(dec *json.Decoder) error {
	for dec.More() {
		if err := b.UnmarshalOne(dec); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalJSON will add the unmarshalled values to this builder.
//
// If the values are strings, they will get parsed with big.ParseFloat using
// a rounding mode of big.ToNearestAway currently.
func (b *Decimal128Builder) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("decimal128 builder must unpack from json array, found %s", delim)
	}

	return b.Unmarshal(dec)
}

var (
	_ arrow.Array = (*Decimal128)(nil)
	_ Builder     = (*Decimal128Builder)(nil)
)
