// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"fmt"
	"reflect"
	"strconv"

	ParticleType "github.com/aerospike/aerospike-client-go/v5/internal/particle_type"
	"github.com/aerospike/aerospike-client-go/v5/types"

	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

// this function will be set in value_slow file if included
var newValueReflect func(interface{}) Value

// MapPair is used when the client returns sorted maps from the server
// Since the default map in Go is a hash map, we will use a slice
// to return the results in server order
type MapPair struct{ Key, Value interface{} }

// Value interface is used to efficiently serialize objects into the wire protocol.
type Value interface {

	// Calculate number of vl.bytes necessary to serialize the value in the wire protocol.
	EstimateSize() (int, Error)

	// Serialize the value in the wire protocol.
	write(cmd BufferEx) (int, Error)

	// Serialize the value using MessagePack.
	pack(cmd BufferEx) (int, Error)

	// GetType returns wire protocol value type.
	GetType() int

	// GetObject returns original value as an interface{}.
	GetObject() interface{}

	// String implements Stringer interface.
	String() string
}

//revive:disable

// AerospikeBlob interface allows the user to write a conversion function from their value to []bytes.
type AerospikeBlob interface {
	// EncodeBlob returns a byte slice representing the encoding of the
	// receiver for transmission to a Decoder, usually of the same
	// concrete type.
	EncodeBlob() ([]byte, error)
}

//revive:enable

// tryConcreteValue will return an aerospike value.
// If the encoder does not exist, it will not try to use reflection.
func tryConcreteValue(v interface{}) Value {
	switch val := v.(type) {
	case Value:
		return val
	case int:
		return IntegerValue(val)
	case int64:
		return LongValue(val)
	case string:
		return StringValue(val)
	case []interface{}:
		return ListValue(val)
	case map[string]interface{}:
		return JsonValue(val)
	case map[interface{}]interface{}:
		return NewMapValue(val)
	case nil:
		return nullValue
	case []Value:
		return NewValueArray(val)
	case []byte:
		return BytesValue(val)
	case int8:
		return IntegerValue(int(val))
	case int16:
		return IntegerValue(int(val))
	case int32:
		return IntegerValue(int(val))
	case uint8: // byte supported here
		return IntegerValue(int(val))
	case uint16:
		return IntegerValue(int(val))
	case uint32:
		return IntegerValue(int(val))
	case float32:
		return FloatValue(float64(val))
	case float64:
		return FloatValue(val)
	case uint:
		// if it doesn't overflow int64, it is OK
		if int64(val) >= 0 {
			return LongValue(int64(val))
		}
	case bool:
		return BoolValue(val)
	case MapIter:
		return NewMapperValue(val)
	case ListIter:
		return NewListerValue(val)
	case AerospikeBlob:
		return NewBlobValue(val)

	/*
		The following cases will try to avoid using reflection by matching against the
		internal generic types.
		If you have custom type aliases in your code, you can use the same aerospike types to cast your type into,
		to avoid hitting the reflection.
	*/
	case []string:
		return NewListerValue(stringSlice(val))
	case []int:
		return NewListerValue(intSlice(val))
	case []int8:
		return NewListerValue(int8Slice(val))
	case []int16:
		return NewListerValue(int16Slice(val))
	case []int32:
		return NewListerValue(int32Slice(val))
	case []int64:
		return NewListerValue(int64Slice(val))
	case []uint16:
		return NewListerValue(uint16Slice(val))
	case []uint32:
		return NewListerValue(uint32Slice(val))
	case []uint64:
		return NewListerValue(uint64Slice(val))
	case []float32:
		return NewListerValue(float32Slice(val))
	case []float64:
		return NewListerValue(float64Slice(val))
	case map[string]string:
		return NewMapperValue(stringStringMap(val))
	case map[string]int:
		return NewMapperValue(stringIntMap(val))
	case map[string]int8:
		return NewMapperValue(stringInt8Map(val))
	case map[string]int16:
		return NewMapperValue(stringInt16Map(val))
	case map[string]int32:
		return NewMapperValue(stringInt32Map(val))
	case map[string]int64:
		return NewMapperValue(stringInt64Map(val))
	case map[string]uint16:
		return NewMapperValue(stringUint16Map(val))
	case map[string]uint32:
		return NewMapperValue(stringUint32Map(val))
	case map[string]float32:
		return NewMapperValue(stringFloat32Map(val))
	case map[string]float64:
		return NewMapperValue(stringFloat64Map(val))
	case map[int]string:
		return NewMapperValue(intStringMap(val))
	case map[int]int:
		return NewMapperValue(intIntMap(val))
	case map[int]int8:
		return NewMapperValue(intInt8Map(val))
	case map[int]int16:
		return NewMapperValue(intInt16Map(val))
	case map[int]int32:
		return NewMapperValue(intInt32Map(val))
	case map[int]int64:
		return NewMapperValue(intInt64Map(val))
	case map[int]uint16:
		return NewMapperValue(intUint16Map(val))
	case map[int]uint32:
		return NewMapperValue(intUint32Map(val))
	case map[int]float32:
		return NewMapperValue(intFloat32Map(val))
	case map[int]float64:
		return NewMapperValue(intFloat64Map(val))
	case map[int]interface{}:
		return NewMapperValue(intInterfaceMap(val))
	case map[int8]string:
		return NewMapperValue(int8StringMap(val))
	case map[int8]int:
		return NewMapperValue(int8IntMap(val))
	case map[int8]int8:
		return NewMapperValue(int8Int8Map(val))
	case map[int8]int16:
		return NewMapperValue(int8Int16Map(val))
	case map[int8]int32:
		return NewMapperValue(int8Int32Map(val))
	case map[int8]int64:
		return NewMapperValue(int8Int64Map(val))
	case map[int8]uint16:
		return NewMapperValue(int8Uint16Map(val))
	case map[int8]uint32:
		return NewMapperValue(int8Uint32Map(val))
	case map[int8]float32:
		return NewMapperValue(int8Float32Map(val))
	case map[int8]float64:
		return NewMapperValue(int8Float64Map(val))
	case map[int8]interface{}:
		return NewMapperValue(int8InterfaceMap(val))
	case map[int16]string:
		return NewMapperValue(int16StringMap(val))
	case map[int16]int:
		return NewMapperValue(int16IntMap(val))
	case map[int16]int8:
		return NewMapperValue(int16Int8Map(val))
	case map[int16]int16:
		return NewMapperValue(int16Int16Map(val))
	case map[int16]int32:
		return NewMapperValue(int16Int32Map(val))
	case map[int16]int64:
		return NewMapperValue(int16Int64Map(val))
	case map[int16]uint16:
		return NewMapperValue(int16Uint16Map(val))
	case map[int16]uint32:
		return NewMapperValue(int16Uint32Map(val))
	case map[int16]float32:
		return NewMapperValue(int16Float32Map(val))
	case map[int16]float64:
		return NewMapperValue(int16Float64Map(val))
	case map[int16]interface{}:
		return NewMapperValue(int16InterfaceMap(val))
	case map[int32]string:
		return NewMapperValue(int32StringMap(val))
	case map[int32]int:
		return NewMapperValue(int32IntMap(val))
	case map[int32]int8:
		return NewMapperValue(int32Int8Map(val))
	case map[int32]int16:
		return NewMapperValue(int32Int16Map(val))
	case map[int32]int32:
		return NewMapperValue(int32Int32Map(val))
	case map[int32]int64:
		return NewMapperValue(int32Int64Map(val))
	case map[int32]uint16:
		return NewMapperValue(int32Uint16Map(val))
	case map[int32]uint32:
		return NewMapperValue(int32Uint32Map(val))
	case map[int32]float32:
		return NewMapperValue(int32Float32Map(val))
	case map[int32]float64:
		return NewMapperValue(int32Float64Map(val))
	case map[int32]interface{}:
		return NewMapperValue(int32InterfaceMap(val))
	case map[int64]string:
		return NewMapperValue(int64StringMap(val))
	case map[int64]int:
		return NewMapperValue(int64IntMap(val))
	case map[int64]int8:
		return NewMapperValue(int64Int8Map(val))
	case map[int64]int16:
		return NewMapperValue(int64Int16Map(val))
	case map[int64]int32:
		return NewMapperValue(int64Int32Map(val))
	case map[int64]int64:
		return NewMapperValue(int64Int64Map(val))
	case map[int64]uint16:
		return NewMapperValue(int64Uint16Map(val))
	case map[int64]uint32:
		return NewMapperValue(int64Uint32Map(val))
	case map[int64]float32:
		return NewMapperValue(int64Float32Map(val))
	case map[int64]float64:
		return NewMapperValue(int64Float64Map(val))
	case map[int64]interface{}:
		return NewMapperValue(int64InterfaceMap(val))
	case map[uint16]string:
		return NewMapperValue(uint16StringMap(val))
	case map[uint16]int:
		return NewMapperValue(uint16IntMap(val))
	case map[uint16]int8:
		return NewMapperValue(uint16Int8Map(val))
	case map[uint16]int16:
		return NewMapperValue(uint16Int16Map(val))
	case map[uint16]int32:
		return NewMapperValue(uint16Int32Map(val))
	case map[uint16]int64:
		return NewMapperValue(uint16Int64Map(val))
	case map[uint16]uint16:
		return NewMapperValue(uint16Uint16Map(val))
	case map[uint16]uint32:
		return NewMapperValue(uint16Uint32Map(val))
	case map[uint16]float32:
		return NewMapperValue(uint16Float32Map(val))
	case map[uint16]float64:
		return NewMapperValue(uint16Float64Map(val))
	case map[uint16]interface{}:
		return NewMapperValue(uint16InterfaceMap(val))
	case map[uint32]string:
		return NewMapperValue(uint32StringMap(val))
	case map[uint32]int:
		return NewMapperValue(uint32IntMap(val))
	case map[uint32]int8:
		return NewMapperValue(uint32Int8Map(val))
	case map[uint32]int16:
		return NewMapperValue(uint32Int16Map(val))
	case map[uint32]int32:
		return NewMapperValue(uint32Int32Map(val))
	case map[uint32]int64:
		return NewMapperValue(uint32Int64Map(val))
	case map[uint32]uint16:
		return NewMapperValue(uint32Uint16Map(val))
	case map[uint32]uint32:
		return NewMapperValue(uint32Uint32Map(val))
	case map[uint32]float32:
		return NewMapperValue(uint32Float32Map(val))
	case map[uint32]float64:
		return NewMapperValue(uint32Float64Map(val))
	case map[uint32]interface{}:
		return NewMapperValue(uint32InterfaceMap(val))
	case map[float32]string:
		return NewMapperValue(float32StringMap(val))
	case map[float32]int:
		return NewMapperValue(float32IntMap(val))
	case map[float32]int8:
		return NewMapperValue(float32Int8Map(val))
	case map[float32]int16:
		return NewMapperValue(float32Int16Map(val))
	case map[float32]int32:
		return NewMapperValue(float32Int32Map(val))
	case map[float32]int64:
		return NewMapperValue(float32Int64Map(val))
	case map[float32]uint16:
		return NewMapperValue(float32Uint16Map(val))
	case map[float32]uint32:
		return NewMapperValue(float32Uint32Map(val))
	case map[float32]float32:
		return NewMapperValue(float32Float32Map(val))
	case map[float32]float64:
		return NewMapperValue(float32Float64Map(val))
	case map[float32]interface{}:
		return NewMapperValue(float32InterfaceMap(val))
	case map[float64]string:
		return NewMapperValue(float64StringMap(val))
	case map[float64]int:
		return NewMapperValue(float64IntMap(val))
	case map[float64]int8:
		return NewMapperValue(float64Int8Map(val))
	case map[float64]int16:
		return NewMapperValue(float64Int16Map(val))
	case map[float64]int32:
		return NewMapperValue(float64Int32Map(val))
	case map[float64]int64:
		return NewMapperValue(float64Int64Map(val))
	case map[float64]uint16:
		return NewMapperValue(float64Uint16Map(val))
	case map[float64]uint32:
		return NewMapperValue(float64Uint32Map(val))
	case map[float64]float32:
		return NewMapperValue(float64Float32Map(val))
	case map[float64]float64:
		return NewMapperValue(float64Float64Map(val))
	case map[float64]interface{}:
		return NewMapperValue(float64InterfaceMap(val))
	case map[string]uint64:
		return NewMapperValue(stringUint64Map(val))
	case map[int]uint64:
		return NewMapperValue(intUint64Map(val))
	case map[int8]uint64:
		return NewMapperValue(int8Uint64Map(val))
	case map[int16]uint64:
		return NewMapperValue(int16Uint64Map(val))
	case map[int32]uint64:
		return NewMapperValue(int32Uint64Map(val))
	case map[int64]uint64:
		return NewMapperValue(int64Uint64Map(val))
	case map[uint16]uint64:
		return NewMapperValue(uint16Uint64Map(val))
	case map[uint32]uint64:
		return NewMapperValue(uint32Uint64Map(val))
	case map[float32]uint64:
		return NewMapperValue(float32Uint64Map(val))
	case map[float64]uint64:
		return NewMapperValue(float64Uint64Map(val))
	case map[uint64]string:
		return NewMapperValue(uint64StringMap(val))
	case map[uint64]int:
		return NewMapperValue(uint64IntMap(val))
	case map[uint64]int8:
		return NewMapperValue(uint64Int8Map(val))
	case map[uint64]int16:
		return NewMapperValue(uint64Int16Map(val))
	case map[uint64]int32:
		return NewMapperValue(uint64Int32Map(val))
	case map[uint64]int64:
		return NewMapperValue(uint64Int64Map(val))
	case map[uint64]uint16:
		return NewMapperValue(uint64Uint16Map(val))
	case map[uint64]uint32:
		return NewMapperValue(uint64Uint32Map(val))
	case map[uint64]uint64:
		return NewMapperValue(uint64Uint64Map(val))
	case map[uint64]float32:
		return NewMapperValue(uint64Float32Map(val))
	case map[uint64]float64:
		return NewMapperValue(uint64Float64Map(val))
	case map[uint64]interface{}:
		return NewMapperValue(uint64InterfaceMap(val))
	}

	return nil
}

// OpResults encapsulates the results of batch read operations
type OpResults []interface{}

// NewValue generates a new Value object based on the type.
// If the type is not supported, NewValue will panic.
// This method is a convenience method, and should not be used
// when absolute performance is required unless for the reason mentioned below.
//
// If you have custom maps or slices like:
//     type MyMap map[primitive1]primitive2, eg: map[int]string
// or
//     type MySlice []primitive, eg: []float64
// cast them to their primitive type when passing them to this method:
//     v := NewValue(map[int]string(myVar))
//     v := NewValue([]float64(myVar))
// This way you will avoid hitting reflection.
// To completely avoid reflection in the library,
// use the build tag: as_performance while building your program.
func NewValue(v interface{}) Value {
	if value := tryConcreteValue(v); value != nil {
		return value
	}

	if newValueReflect != nil {
		if res := newValueReflect(v); res != nil {
			return res
		}
	}

	// panic for anything that is not supported.
	panic(newError(types.TYPE_NOT_SUPPORTED, fmt.Sprintf("Value type '%v' (%s) not supported (if you are compiling via 'as_performance' tag, use cast either to primitives, or use ListIter or MapIter interfaces.)", v, reflect.TypeOf(v).String())))
}

// NullValue is an empty value.
type NullValue struct{}

var nullValue NullValue

// NewNullValue generates a NullValue instance.
func NewNullValue() NullValue {
	return nullValue
}

// EstimateSize returns the size of the NullValue in wire protocol.
func (vl NullValue) EstimateSize() (int, Error) {
	return 0, nil
}

func (vl NullValue) write(cmd BufferEx) (int, Error) {
	return 0, nil
}

func (vl NullValue) pack(cmd BufferEx) (int, Error) {
	return packNil(cmd)
}

// GetType returns wire protocol value type.
func (vl NullValue) GetType() int {
	return ParticleType.NULL
}

// GetObject returns original value as an interface{}.
func (vl NullValue) GetObject() interface{} {
	return nil
}

func (vl NullValue) String() string {
	return ""
}

///////////////////////////////////////////////////////////////////////////////

// InfinityValue is an empty value.
type InfinityValue struct{}

var infinityValue InfinityValue

// NewInfinityValue generates a InfinityValue instance.
func NewInfinityValue() InfinityValue {
	return infinityValue
}

// EstimateSize returns the size of the InfinityValue in wire protocol.
func (vl InfinityValue) EstimateSize() (int, Error) {
	return 0, nil
}

func (vl InfinityValue) write(cmd BufferEx) (int, Error) {
	return 0, nil
}

func (vl InfinityValue) pack(cmd BufferEx) (int, Error) {
	return packInfinity(cmd)
}

// GetType returns wire protocol value type.
func (vl InfinityValue) GetType() int {
	panic("Invalid particle type: INF")
}

// GetObject returns original value as an interface{}.
func (vl InfinityValue) GetObject() interface{} {
	return nil
}

func (vl InfinityValue) String() string {
	return "INF"
}

///////////////////////////////////////////////////////////////////////////////

// WildCardValue is an empty value.
type WildCardValue struct{}

var wildCardValue WildCardValue

// NewWildCardValue generates a WildCardValue instance.
func NewWildCardValue() WildCardValue {
	return wildCardValue
}

// EstimateSize returns the size of the WildCardValue in wire protocol.
func (vl WildCardValue) EstimateSize() (int, Error) {
	return 0, nil
}

func (vl WildCardValue) write(cmd BufferEx) (int, Error) {
	return 0, nil
}

func (vl WildCardValue) pack(cmd BufferEx) (int, Error) {
	return packWildCard(cmd)
}

// GetType returns wire protocol value type.
func (vl WildCardValue) GetType() int {
	panic("Invalid particle type: WildCard")
}

// GetObject returns original value as an interface{}.
func (vl WildCardValue) GetObject() interface{} {
	return nil
}

func (vl WildCardValue) String() string {
	return "*"
}

///////////////////////////////////////////////////////////////////////////////

// BytesValue encapsulates an array of bytes.
type BytesValue []byte

// NewBytesValue generates a ByteValue instance.
func NewBytesValue(bytes []byte) BytesValue {
	return BytesValue(bytes)
}

// NewBlobValue accepts an AerospikeBlob interface, and automatically
// converts it to a BytesValue.
// If Encode returns an err, it will panic.
func NewBlobValue(object AerospikeBlob) BytesValue {
	buf, err := object.EncodeBlob()
	if err != nil {
		panic(err)
	}

	return NewBytesValue(buf)
}

// EstimateSize returns the size of the BytesValue in wire protocol.
func (vl BytesValue) EstimateSize() (int, Error) {
	return len(vl), nil
}

func (vl BytesValue) write(cmd BufferEx) (int, Error) {
	return cmd.Write(vl)
}

func (vl BytesValue) pack(cmd BufferEx) (int, Error) {
	return packBytes(cmd, vl)
}

// GetType returns wire protocol value type.
func (vl BytesValue) GetType() int {
	return ParticleType.BLOB
}

// GetObject returns original value as an interface{}.
func (vl BytesValue) GetObject() interface{} {
	return []byte(vl)
}

// String implements Stringer interface.
func (vl BytesValue) String() string {
	return Buffer.BytesToHexString(vl)
}

///////////////////////////////////////////////////////////////////////////////

// StringValue encapsulates a string value.
type StringValue string

// NewStringValue generates a StringValue instance.
func NewStringValue(value string) StringValue {
	return StringValue(value)
}

// EstimateSize returns the size of the StringValue in wire protocol.
func (vl StringValue) EstimateSize() (int, Error) {
	return len(vl), nil
}

func (vl StringValue) write(cmd BufferEx) (int, Error) {
	return cmd.WriteString(string(vl))
}

func (vl StringValue) pack(cmd BufferEx) (int, Error) {
	return packString(cmd, string(vl))
}

// GetType returns wire protocol value type.
func (vl StringValue) GetType() int {
	return ParticleType.STRING
}

// GetObject returns original value as an interface{}.
func (vl StringValue) GetObject() interface{} {
	return string(vl)
}

// String implements Stringer interface.
func (vl StringValue) String() string {
	return string(vl)
}

///////////////////////////////////////////////////////////////////////////////

// IntegerValue encapsulates an integer value.
type IntegerValue int

// NewIntegerValue generates an IntegerValue instance.
func NewIntegerValue(value int) IntegerValue {
	return IntegerValue(value)
}

// EstimateSize returns the size of the IntegerValue in wire protocol.
func (vl IntegerValue) EstimateSize() (int, Error) {
	return 8, nil
}

func (vl IntegerValue) write(cmd BufferEx) (int, Error) {
	n := cmd.WriteInt64(int64(vl))
	return n, nil
}

func (vl IntegerValue) pack(cmd BufferEx) (int, Error) {
	return packAInt64(cmd, int64(vl))
}

// GetType returns wire protocol value type.
func (vl IntegerValue) GetType() int {
	return ParticleType.INTEGER
}

// GetObject returns original value as an interface{}.
func (vl IntegerValue) GetObject() interface{} {
	return int(vl)
}

// String implements Stringer interface.
func (vl IntegerValue) String() string {
	return strconv.Itoa(int(vl))
}

///////////////////////////////////////////////////////////////////////////////

// LongValue encapsulates an int64 value.
type LongValue int64

// NewLongValue generates a LongValue instance.
func NewLongValue(value int64) LongValue {
	return LongValue(value)
}

// EstimateSize returns the size of the LongValue in wire protocol.
func (vl LongValue) EstimateSize() (int, Error) {
	return 8, nil
}

func (vl LongValue) write(cmd BufferEx) (int, Error) {
	n := cmd.WriteInt64(int64(vl))
	return n, nil
}

func (vl LongValue) pack(cmd BufferEx) (int, Error) {
	return packAInt64(cmd, int64(vl))
}

// GetType returns wire protocol value type.
func (vl LongValue) GetType() int {
	return ParticleType.INTEGER
}

// GetObject returns original value as an interface{}.
func (vl LongValue) GetObject() interface{} {
	return int64(vl)
}

// String implements Stringer interface.
func (vl LongValue) String() string {
	return strconv.Itoa(int(vl))
}

///////////////////////////////////////////////////////////////////////////////

// FloatValue encapsulates an float64 value.
type FloatValue float64

// NewFloatValue generates a FloatValue instance.
func NewFloatValue(value float64) FloatValue {
	return FloatValue(value)
}

// EstimateSize returns the size of the FloatValue in wire protocol.
func (vl FloatValue) EstimateSize() (int, Error) {
	return 8, nil
}

func (vl FloatValue) write(cmd BufferEx) (int, Error) {
	n := cmd.WriteFloat64(float64(vl))
	return n, nil
}

func (vl FloatValue) pack(cmd BufferEx) (int, Error) {
	return packFloat64(cmd, float64(vl))
}

// GetType returns wire protocol value type.
func (vl FloatValue) GetType() int {
	return ParticleType.FLOAT
}

// GetObject returns original value as an interface{}.
func (vl FloatValue) GetObject() interface{} {
	return float64(vl)
}

// String implements Stringer interface.
func (vl FloatValue) String() string {
	return (fmt.Sprintf("%f", vl))
}

///////////////////////////////////////////////////////////////////////////////

// BoolValue encapsulates a boolean value.
// Supported by Aerospike server v5.6+ only.
type BoolValue bool

// EstimateSize returns the size of the BoolValue in wire protocol.
func (vb BoolValue) EstimateSize() (int, Error) {
	return PackBool(nil, bool(vb))
}

func (vb BoolValue) write(cmd BufferEx) (int, Error) {
	n := cmd.WriteBool(bool(vb))
	return n, nil
}

func (vb BoolValue) pack(cmd BufferEx) (int, Error) {
	return PackBool(cmd, bool(vb))
}

// GetType returns wire protocol value type.
func (vb BoolValue) GetType() int {
	return ParticleType.BOOL
}

// GetObject returns original value as an interface{}.
func (vb BoolValue) GetObject() interface{} {
	return bool(vb)
}

// String implements Stringer interface.
func (vb BoolValue) String() string {
	return (fmt.Sprintf("%v", bool(vb)))
}

///////////////////////////////////////////////////////////////////////////////

// ValueArray encapsulates an array of Value.
// Supported by Aerospike 3+ servers only.
type ValueArray []Value

// NewValueArray generates a ValueArray instance.
func NewValueArray(array []Value) *ValueArray {
	// return &ValueArray{*NewListerValue(valueList(array))}
	res := ValueArray(array)
	return &res
}

// EstimateSize returns the size of the ValueArray in wire protocol.
func (va ValueArray) EstimateSize() (int, Error) {
	return packValueArray(nil, va)
}

func (va ValueArray) write(cmd BufferEx) (int, Error) {
	return packValueArray(cmd, va)
}

func (va ValueArray) pack(cmd BufferEx) (int, Error) {
	return packValueArray(cmd, []Value(va))
}

// GetType returns wire protocol value type.
func (va ValueArray) GetType() int {
	return ParticleType.LIST
}

// GetObject returns original value as an interface{}.
func (va ValueArray) GetObject() interface{} {
	return []Value(va)
}

// String implements Stringer interface.
func (va ValueArray) String() string {
	return fmt.Sprintf("%v", []Value(va))
}

///////////////////////////////////////////////////////////////////////////////

// ListValue encapsulates any arbitrary array.
// Supported by Aerospike 3+ servers only.
type ListValue []interface{}

// NewListValue generates a ListValue instance.
func NewListValue(list []interface{}) ListValue {
	return ListValue(list)
}

// EstimateSize returns the size of the ListValue in wire protocol.
func (vl ListValue) EstimateSize() (int, Error) {
	return packIfcList(nil, vl)
}

func (vl ListValue) write(cmd BufferEx) (int, Error) {
	return packIfcList(cmd, vl)
}

func (vl ListValue) pack(cmd BufferEx) (int, Error) {
	return packIfcList(cmd, []interface{}(vl))
}

// GetType returns wire protocol value type.
func (vl ListValue) GetType() int {
	return ParticleType.LIST
}

// GetObject returns original value as an interface{}.
func (vl ListValue) GetObject() interface{} {
	return []interface{}(vl)
}

// String implements Stringer interface.
func (vl ListValue) String() string {
	return fmt.Sprintf("%v", []interface{}(vl))
}

///////////////////////////////////////////////////////////////////////////////

// ListerValue encapsulates any arbitrary array.
// Supported by Aerospike 3+ servers only.
type ListerValue struct {
	list ListIter
}

// NewListerValue generates a NewListerValue instance.
func NewListerValue(list ListIter) *ListerValue {
	res := &ListerValue{
		list: list,
	}

	return res
}

// EstimateSize returns the size of the ListerValue in wire protocol.
func (vl *ListerValue) EstimateSize() (int, Error) {
	return packList(nil, vl.list)
}

func (vl *ListerValue) write(cmd BufferEx) (int, Error) {
	return packList(cmd, vl.list)
}

func (vl *ListerValue) pack(cmd BufferEx) (int, Error) {
	return packList(cmd, vl.list)
}

// GetType returns wire protocol value type.
func (vl *ListerValue) GetType() int {
	return ParticleType.LIST
}

// GetObject returns original value as an interface{}.
func (vl *ListerValue) GetObject() interface{} {
	return vl.list
}

// String implements Stringer interface.
func (vl *ListerValue) String() string {
	return fmt.Sprintf("%v", vl.list)
}

///////////////////////////////////////////////////////////////////////////////

// MapValue encapsulates an arbitrary map.
// Supported by Aerospike 3+ servers only.
type MapValue map[interface{}]interface{}

// NewMapValue generates a MapValue instance.
func NewMapValue(vmap map[interface{}]interface{}) MapValue {
	return MapValue(vmap)
}

// EstimateSize returns the size of the MapValue in wire protocol.
func (vl MapValue) EstimateSize() (int, Error) {
	return packIfcMap(nil, vl)
}

func (vl MapValue) write(cmd BufferEx) (int, Error) {
	return packIfcMap(cmd, vl)
}

func (vl MapValue) pack(cmd BufferEx) (int, Error) {
	return packIfcMap(cmd, vl)
}

// GetType returns wire protocol value type.
func (vl MapValue) GetType() int {
	return ParticleType.MAP
}

// GetObject returns original value as an interface{}.
func (vl MapValue) GetObject() interface{} {
	return map[interface{}]interface{}(vl)
}

func (vl MapValue) String() string {
	return fmt.Sprintf("%v", map[interface{}]interface{}(vl))
}

///////////////////////////////////////////////////////////////////////////////

// JsonValue encapsulates a Json map.
// Supported by Aerospike 3+ servers only.
type JsonValue map[string]interface{}

// NewJsonValue generates a JsonValue instance.
func NewJsonValue(vmap map[string]interface{}) JsonValue {
	return JsonValue(vmap)
}

// EstimateSize returns the size of the JsonValue in wire protocol.
func (vl JsonValue) EstimateSize() (int, Error) {
	return packJsonMap(nil, vl)
}

func (vl JsonValue) write(cmd BufferEx) (int, Error) {
	return packJsonMap(cmd, vl)
}

func (vl JsonValue) pack(cmd BufferEx) (int, Error) {
	return packJsonMap(cmd, vl)
}

// GetType returns wire protocol value type.
func (vl JsonValue) GetType() int {
	return ParticleType.MAP
}

// GetObject returns original value as an interface{}.
func (vl JsonValue) GetObject() interface{} {
	return map[string]interface{}(vl)
}

func (vl JsonValue) String() string {
	return fmt.Sprintf("%v", map[string]interface{}(vl))
}

///////////////////////////////////////////////////////////////////////////////

// MapperValue encapsulates an arbitrary map which implements a MapIter interface.
// Supported by Aerospike 3+ servers only.
type MapperValue struct {
	vmap MapIter
}

// NewMapperValue generates a MapperValue instance.
func NewMapperValue(vmap MapIter) *MapperValue {
	res := &MapperValue{
		vmap: vmap,
	}

	return res
}

// EstimateSize returns the size of the MapperValue in wire protocol.
func (vl *MapperValue) EstimateSize() (int, Error) {
	return packMap(nil, vl.vmap)
}

func (vl *MapperValue) write(cmd BufferEx) (int, Error) {
	return packMap(cmd, vl.vmap)
}

func (vl *MapperValue) pack(cmd BufferEx) (int, Error) {
	return packMap(cmd, vl.vmap)
}

// GetType returns wire protocol value type.
func (vl *MapperValue) GetType() int {
	return ParticleType.MAP
}

// GetObject returns original value as an interface{}.
func (vl *MapperValue) GetObject() interface{} {
	return vl.vmap
}

func (vl *MapperValue) String() string {
	return fmt.Sprintf("%v", vl.vmap)
}

///////////////////////////////////////////////////////////////////////////////

// GeoJSONValue encapsulates a 2D Geo point.
// Supported by Aerospike 3.6.1 servers and later only.
type GeoJSONValue string

// NewGeoJSONValue generates a GeoJSONValue instance.
func NewGeoJSONValue(value string) GeoJSONValue {
	res := GeoJSONValue(value)
	return res
}

// EstimateSize returns the size of the GeoJSONValue in wire protocol.
func (vl GeoJSONValue) EstimateSize() (int, Error) {
	// flags + ncells + jsonstr
	return 1 + 2 + len(string(vl)), nil
}

func (vl GeoJSONValue) write(cmd BufferEx) (int, Error) {
	cmd.WriteByte(0) // flags
	cmd.WriteByte(0) // flags
	cmd.WriteByte(0) // flags

	return cmd.WriteString(string(vl))
}

func (vl GeoJSONValue) pack(cmd BufferEx) (int, Error) {
	return packGeoJson(cmd, string(vl))
}

// GetType returns wire protocol value type.
func (vl GeoJSONValue) GetType() int {
	return ParticleType.GEOJSON
}

// GetObject returns original value as an interface{}.
func (vl GeoJSONValue) GetObject() interface{} {
	return string(vl)
}

// String implements Stringer interface.
func (vl GeoJSONValue) String() string {
	return string(vl)
}

///////////////////////////////////////////////////////////////////////////////

// HLLValue encapsulates a HyperLogLog value.
type HLLValue []byte

// NewHLLValue generates a ByteValue instance.
func NewHLLValue(bytes []byte) HLLValue {
	return HLLValue(bytes)
}

// EstimateSize returns the size of the HLLValue in wire protocol.
func (vl HLLValue) EstimateSize() (int, Error) {
	return len(vl), nil
}

func (vl HLLValue) write(cmd BufferEx) (int, Error) {
	return cmd.Write(vl)
}

func (vl HLLValue) pack(cmd BufferEx) (int, Error) {
	return packBytes(cmd, vl)
}

// GetType returns wire protocol value type.
func (vl HLLValue) GetType() int {
	return ParticleType.HLL
}

// GetObject returns original value as an interface{}.
func (vl HLLValue) GetObject() interface{} {
	return []byte(vl)
}

// String implements Stringer interface.
func (vl HLLValue) String() string {
	return Buffer.BytesToHexString([]byte(vl))
}

//////////////////////////////////////////////////////////////////////////////

func bytesToParticle(ptype int, buf []byte, offset int, length int) (interface{}, Error) {

	switch ptype {
	case ParticleType.INTEGER:
		// return `int` for 64bit platforms for compatibility reasons
		if Buffer.Arch64Bits {
			return int(Buffer.VarBytesToInt64(buf, offset, length)), nil
		}
		return Buffer.VarBytesToInt64(buf, offset, length), nil

	case ParticleType.STRING:
		return string(buf[offset : offset+length]), nil

	case ParticleType.FLOAT:
		return Buffer.BytesToFloat64(buf, offset), nil

	case ParticleType.BOOL:
		return Buffer.BytesToBool(buf, offset, length), nil

	case ParticleType.MAP:
		return newUnpacker(buf, offset, length).UnpackMap()

	case ParticleType.LIST:
		return newUnpacker(buf, offset, length).UnpackList()

	case ParticleType.GEOJSON:
		ncells := int(Buffer.BytesToInt16(buf, offset+1))
		headerSize := 1 + 2 + (ncells * 8)
		return string(buf[offset+headerSize : offset+length]), nil

	case ParticleType.HLL:
		newObj := make([]byte, length)
		copy(newObj, buf[offset:offset+length])
		return newObj, nil

	case ParticleType.BLOB:
		newObj := make([]byte, length)
		copy(newObj, buf[offset:offset+length])
		return newObj, nil

	case ParticleType.LDT:
		return newUnpacker(buf, offset, length).unpackObjects()

	}
	return nil, nil
}

func bytesToKeyValue(pType int, buf []byte, offset int, length int) (Value, Error) {

	switch pType {
	case ParticleType.STRING:
		return NewStringValue(string(buf[offset : offset+length])), nil

	case ParticleType.INTEGER:
		return NewLongValue(Buffer.VarBytesToInt64(buf, offset, length)), nil

	case ParticleType.FLOAT:
		return NewFloatValue(Buffer.BytesToFloat64(buf, offset)), nil

	case ParticleType.BLOB:
		bytes := make([]byte, length)
		copy(bytes, buf[offset:offset+length])
		return NewBytesValue(bytes), nil

	case ParticleType.LIST:
		v, err := newUnpacker(buf, offset, length).UnpackList()
		if err != nil {
			return nil, err
		}
		return ListValue(v), nil

	default:
		return nil, newError(types.PARSE_ERROR, fmt.Sprintf("ParticleType %d not recognized. Please file a github issue.", pType))
	}
}

func unwrapValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	if uv, ok := v.(Value); ok {
		return unwrapValue(uv.GetObject())
	} else if uv, ok := v.([]Value); ok {
		a := make([]interface{}, len(uv))
		for i := range uv {
			a[i] = unwrapValue(uv[i].GetObject())
		}
		return a
	}

	return v
}
