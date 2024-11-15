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
	"github.com/apache/arrow/go/v15/arrow/decimal256"
	"github.com/apache/arrow/go/v15/arrow/internal/debug"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/apache/arrow/go/v15/internal/json"
)

// Decimal256 is a type that represents an immutable sequence of 256-bit decimal values.
type Decimal256 struct {
	array

	values []decimal256.Num
}

func NewDecimal256Data(data arrow.ArrayData) *Decimal256 {
	a := &Decimal256{}
	a.refCount = 1
	a.setData(data.(*Data))
	return a
}

func (a *Decimal256) Value(i int) decimal256.Num { return a.values[i] }

func (a *Decimal256) ValueStr(i int) string {
	if a.IsNull(i) {
		return NullValueStr
	}
	return a.GetOneForMarshal(i).(string)
}

func (a *Decimal256) Values() []decimal256.Num { return a.values }

func (a *Decimal256) String() string {
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

func (a *Decimal256) setData(data *Data) {
	a.array.setData(data)
	vals := data.buffers[1]
	if vals != nil {
		a.values = arrow.Decimal256Traits.CastFromBytes(vals.Bytes())
		beg := a.array.data.offset
		end := beg + a.array.data.length
		a.values = a.values[beg:end]
	}
}

func (a *Decimal256) GetOneForMarshal(i int) interface{} {
	if a.IsNull(i) {
		return nil
	}

	typ := a.DataType().(*arrow.Decimal256Type)
	f := (&big.Float{}).SetInt(a.Value(i).BigInt())
	f.Quo(f, big.NewFloat(math.Pow10(int(typ.Scale))))
	return f.Text('g', int(typ.Precision))
}

func (a *Decimal256) MarshalJSON() ([]byte, error) {
	vals := make([]interface{}, a.Len())
	for i := 0; i < a.Len(); i++ {
		vals[i] = a.GetOneForMarshal(i)
	}
	return json.Marshal(vals)
}

func arrayEqualDecimal256(left, right *Decimal256) bool {
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

type Decimal256Builder struct {
	builder

	dtype   *arrow.Decimal256Type
	data    *memory.Buffer
	rawData []decimal256.Num
}

func NewDecimal256Builder(mem memory.Allocator, dtype *arrow.Decimal256Type) *Decimal256Builder {
	return &Decimal256Builder{
		builder: builder{refCount: 1, mem: mem},
		dtype:   dtype,
	}
}

// Release decreases the reference count by 1.
// When the reference count goes to zero, the memory is freed.
func (b *Decimal256Builder) Release() {
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

func (b *Decimal256Builder) Append(v decimal256.Num) {
	b.Reserve(1)
	b.UnsafeAppend(v)
}

func (b *Decimal256Builder) UnsafeAppend(v decimal256.Num) {
	bitutil.SetBit(b.nullBitmap.Bytes(), b.length)
	b.rawData[b.length] = v
	b.length++
}

func (b *Decimal256Builder) AppendNull() {
	b.Reserve(1)
	b.UnsafeAppendBoolToBitmap(false)
}

func (b *Decimal256Builder) AppendNulls(n int) {
	for i := 0; i < n; i++ {
		b.AppendNull()
	}
}

func (b *Decimal256Builder) AppendEmptyValue() {
	b.Append(decimal256.Num{})
}

func (b *Decimal256Builder) AppendEmptyValues(n int) {
	for i := 0; i < n; i++ {
		b.AppendEmptyValue()
	}
}

func (b *Decimal256Builder) Type() arrow.DataType { return b.dtype }

func (b *Decimal256Builder) UnsafeAppendBoolToBitmap(isValid bool) {
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
func (b *Decimal256Builder) AppendValues(v []decimal256.Num, valid []bool) {
	if len(v) != len(valid) && len(valid) != 0 {
		panic("arrow/array: len(v) != len(valid) && len(valid) != 0")
	}

	if len(v) == 0 {
		return
	}

	b.Reserve(len(v))
	if len(v) > 0 {
		arrow.Decimal256Traits.Copy(b.rawData[b.length:], v)
	}
	b.builder.unsafeAppendBoolsToBitmap(valid, len(v))
}

func (b *Decimal256Builder) init(capacity int) {
	b.builder.init(capacity)

	b.data = memory.NewResizableBuffer(b.mem)
	bytesN := arrow.Decimal256Traits.BytesRequired(capacity)
	b.data.Resize(bytesN)
	b.rawData = arrow.Decimal256Traits.CastFromBytes(b.data.Bytes())
}

// Reserve ensures there is enough space for appending n elements
// by checking the capacity and calling Resize if necessary.
func (b *Decimal256Builder) Reserve(n int) {
	b.builder.reserve(n, b.Resize)
}

// Resize adjusts the space allocated by b to n elements. If n is greater than b.Cap(),
// additional memory will be allocated. If n is smaller, the allocated memory may reduced.
func (b *Decimal256Builder) Resize(n int) {
	nBuilder := n
	if n < minBuilderCapacity {
		n = minBuilderCapacity
	}

	if b.capacity == 0 {
		b.init(n)
	} else {
		b.builder.resize(nBuilder, b.init)
		b.data.Resize(arrow.Decimal256Traits.BytesRequired(n))
		b.rawData = arrow.Decimal256Traits.CastFromBytes(b.data.Bytes())
	}
}

// NewArray creates a Decimal256 array from the memory buffers used by the builder and resets the Decimal256Builder
// so it can be used to build a new array.
func (b *Decimal256Builder) NewArray() arrow.Array {
	return b.NewDecimal256Array()
}

// NewDecimal256Array creates a Decimal256 array from the memory buffers used by the builder and resets the Decimal256Builder
// so it can be used to build a new array.
func (b *Decimal256Builder) NewDecimal256Array() (a *Decimal256) {
	data := b.newData()
	a = NewDecimal256Data(data)
	data.Release()
	return
}

func (b *Decimal256Builder) newData() (data *Data) {
	bytesRequired := arrow.Decimal256Traits.BytesRequired(b.length)
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

func (b *Decimal256Builder) AppendValueFromString(s string) error {
	if s == NullValueStr {
		b.AppendNull()
		return nil
	}
	val, err := decimal256.FromString(s, b.dtype.Precision, b.dtype.Scale)
	if err != nil {
		b.AppendNull()
		return err
	}
	b.Append(val)
	return nil
}

func (b *Decimal256Builder) UnmarshalOne(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}

	switch v := t.(type) {
	case float64:
		val, err := decimal256.FromFloat64(v, b.dtype.Precision, b.dtype.Scale)
		if err != nil {
			return err
		}
		b.Append(val)
	case string:
		out, err := decimal256.FromString(v, b.dtype.Precision, b.dtype.Scale)
		if err != nil {
			return err
		}
		b.Append(out)
	case json.Number:
		out, err := decimal256.FromString(v.String(), b.dtype.Precision, b.dtype.Scale)
		if err != nil {
			return err
		}
		b.Append(out)
	case nil:
		b.AppendNull()
		return nil
	default:
		return &json.UnmarshalTypeError{
			Value:  fmt.Sprint(t),
			Type:   reflect.TypeOf(decimal256.Num{}),
			Offset: dec.InputOffset(),
		}
	}

	return nil
}

func (b *Decimal256Builder) Unmarshal(dec *json.Decoder) error {
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
func (b *Decimal256Builder) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("arrow/array: decimal256 builder must unpack from json array, found %s", delim)
	}

	return b.Unmarshal(dec)
}

var (
	_ arrow.Array = (*Decimal256)(nil)
	_ Builder     = (*Decimal256Builder)(nil)
)
