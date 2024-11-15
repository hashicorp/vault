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

type stringSlice []string

// PackList packs StringSlice as msgpack.
func (ts stringSlice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackString(buf, elem)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of StringSlice
func (ts stringSlice) Len() int {
	return len(ts)
}

type intSlice []int

// PackList packs IntSlice as msgpack.
func (ts intSlice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackInt64(buf, int64(elem))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of IntSlice
func (ts intSlice) Len() int {
	return len(ts)
}

type int8Slice []int8

// PackList packs Int8Slice as msgpack.
func (ts int8Slice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackInt64(buf, int64(elem))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of Int8Slice
func (ts int8Slice) Len() int {
	return len(ts)
}

type int16Slice []int16

// PackList packs Int16Slice as msgpack.
func (ts int16Slice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackInt64(buf, int64(elem))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of Int16Slice
func (ts int16Slice) Len() int {
	return len(ts)
}

type int32Slice []int32

// PackList packs Int32Slice as msgpack.
func (ts int32Slice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackInt64(buf, int64(elem))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of Int32Slice
func (ts int32Slice) Len() int {
	return len(ts)
}

type int64Slice []int64

// PackList packs Int64Slice as msgpack.
func (ts int64Slice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackInt64(buf, elem)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of Int64Slice
func (ts int64Slice) Len() int {
	return len(ts)
}

type uint16Slice []uint16

// PackList packs Uint16Slice as msgpack.
func (ts uint16Slice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackInt64(buf, int64(elem))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of Uint16Slice
func (ts uint16Slice) Len() int {
	return len(ts)
}

type uint32Slice []uint32

// PackList packs Uint32Slice as msgpack.
func (ts uint32Slice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackInt64(buf, int64(elem))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of Uint32Slice
func (ts uint32Slice) Len() int {
	return len(ts)
}

type uint64Slice []uint64

// PackList packs Uint64Slice as msgpack.
func (ts uint64Slice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackUInt64(buf, elem)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of Uint64Slice
func (ts uint64Slice) Len() int {
	return len(ts)
}

type float32Slice []float32

// PackList packs Float32Slice as msgpack.
func (ts float32Slice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackFloat32(buf, elem)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of Float32Slice
func (ts float32Slice) Len() int {
	return len(ts)
}

type float64Slice []float64

// PackList packs Float64Slice as msgpack.
func (ts float64Slice) PackList(buf BufferEx) (int, error) {
	size := 0
	for _, elem := range ts {
		n, err := PackFloat64(buf, elem)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of Float64Slice
func (ts float64Slice) Len() int {
	return len(ts)
}

///////////////////////////////////////////////////////////////////////////////////////////

type stringStringMap map[string]string

//PackMap packs TypeMap as msgpack.
func (tm stringStringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringStringMap) Len() int {
	return len(tm)
}

type stringIntMap map[string]int

//PackMap packs TypeMap as msgpack.
func (tm stringIntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringIntMap) Len() int {
	return len(tm)
}

type stringInt8Map map[string]int8

//PackMap packs TypeMap as msgpack.
func (tm stringInt8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringInt8Map) Len() int {
	return len(tm)
}

type stringInt16Map map[string]int16

//PackMap packs TypeMap as msgpack.
func (tm stringInt16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringInt16Map) Len() int {
	return len(tm)
}

type stringInt32Map map[string]int32

//PackMap packs TypeMap as msgpack.
func (tm stringInt32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringInt32Map) Len() int {
	return len(tm)
}

type stringInt64Map map[string]int64

//PackMap packs TypeMap as msgpack.
func (tm stringInt64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringInt64Map) Len() int {
	return len(tm)
}

type stringUint16Map map[string]uint16

//PackMap packs TypeMap as msgpack.
func (tm stringUint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringUint16Map) Len() int {
	return len(tm)
}

type stringUint32Map map[string]uint32

//PackMap packs TypeMap as msgpack.
func (tm stringUint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringUint32Map) Len() int {
	return len(tm)
}

type stringFloat32Map map[string]float32

//PackMap packs TypeMap as msgpack.
func (tm stringFloat32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringFloat32Map) Len() int {
	return len(tm)
}

type stringFloat64Map map[string]float64

//PackMap packs TypeMap as msgpack.
func (tm stringFloat64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringFloat64Map) Len() int {
	return len(tm)
}

type intStringMap map[int]string

//PackMap packs TypeMap as msgpack.
func (tm intStringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intStringMap) Len() int {
	return len(tm)
}

type intIntMap map[int]int

//PackMap packs TypeMap as msgpack.
func (tm intIntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intIntMap) Len() int {
	return len(tm)
}

type intInt8Map map[int]int8

//PackMap packs TypeMap as msgpack.
func (tm intInt8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intInt8Map) Len() int {
	return len(tm)
}

type intInt16Map map[int]int16

//PackMap packs TypeMap as msgpack.
func (tm intInt16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intInt16Map) Len() int {
	return len(tm)
}

type intInt32Map map[int]int32

//PackMap packs TypeMap as msgpack.
func (tm intInt32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intInt32Map) Len() int {
	return len(tm)
}

type intInt64Map map[int]int64

//PackMap packs TypeMap as msgpack.
func (tm intInt64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intInt64Map) Len() int {
	return len(tm)
}

type intUint16Map map[int]uint16

//PackMap packs TypeMap as msgpack.
func (tm intUint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intUint16Map) Len() int {
	return len(tm)
}

type intUint32Map map[int]uint32

//PackMap packs TypeMap as msgpack.
func (tm intUint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intUint32Map) Len() int {
	return len(tm)
}

type intFloat32Map map[int]float32

//PackMap packs TypeMap as msgpack.
func (tm intFloat32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intFloat32Map) Len() int {
	return len(tm)
}

type intFloat64Map map[int]float64

//PackMap packs TypeMap as msgpack.
func (tm intFloat64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intFloat64Map) Len() int {
	return len(tm)
}

type intInterfaceMap map[int]interface{}

//PackMap packs TypeMap as msgpack.
func (tm intInterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intInterfaceMap) Len() int {
	return len(tm)
}

type int8StringMap map[int8]string

//PackMap packs TypeMap as msgpack.
func (tm int8StringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8StringMap) Len() int {
	return len(tm)
}

type int8IntMap map[int8]int

//PackMap packs TypeMap as msgpack.
func (tm int8IntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8IntMap) Len() int {
	return len(tm)
}

type int8Int8Map map[int8]int8

//PackMap packs TypeMap as msgpack.
func (tm int8Int8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8Int8Map) Len() int {
	return len(tm)
}

type int8Int16Map map[int8]int16

//PackMap packs TypeMap as msgpack.
func (tm int8Int16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8Int16Map) Len() int {
	return len(tm)
}

type int8Int32Map map[int8]int32

//PackMap packs TypeMap as msgpack.
func (tm int8Int32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8Int32Map) Len() int {
	return len(tm)
}

type int8Int64Map map[int8]int64

//PackMap packs TypeMap as msgpack.
func (tm int8Int64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8Int64Map) Len() int {
	return len(tm)
}

type int8Uint16Map map[int8]uint16

//PackMap packs TypeMap as msgpack.
func (tm int8Uint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8Uint16Map) Len() int {
	return len(tm)
}

type int8Uint32Map map[int8]uint32

//PackMap packs TypeMap as msgpack.
func (tm int8Uint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8Uint32Map) Len() int {
	return len(tm)
}

type int8Float32Map map[int8]float32

//PackMap packs TypeMap as msgpack.
func (tm int8Float32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8Float32Map) Len() int {
	return len(tm)
}

type int8Float64Map map[int8]float64

//PackMap packs TypeMap as msgpack.
func (tm int8Float64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8Float64Map) Len() int {
	return len(tm)
}

type int8InterfaceMap map[int8]interface{}

//PackMap packs TypeMap as msgpack.
func (tm int8InterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8InterfaceMap) Len() int {
	return len(tm)
}

type int16StringMap map[int16]string

//PackMap packs TypeMap as msgpack.
func (tm int16StringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16StringMap) Len() int {
	return len(tm)
}

type int16IntMap map[int16]int

//PackMap packs TypeMap as msgpack.
func (tm int16IntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16IntMap) Len() int {
	return len(tm)
}

type int16Int8Map map[int16]int8

//PackMap packs TypeMap as msgpack.
func (tm int16Int8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16Int8Map) Len() int {
	return len(tm)
}

type int16Int16Map map[int16]int16

//PackMap packs TypeMap as msgpack.
func (tm int16Int16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16Int16Map) Len() int {
	return len(tm)
}

type int16Int32Map map[int16]int32

//PackMap packs TypeMap as msgpack.
func (tm int16Int32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16Int32Map) Len() int {
	return len(tm)
}

type int16Int64Map map[int16]int64

//PackMap packs TypeMap as msgpack.
func (tm int16Int64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16Int64Map) Len() int {
	return len(tm)
}

type int16Uint16Map map[int16]uint16

//PackMap packs TypeMap as msgpack.
func (tm int16Uint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16Uint16Map) Len() int {
	return len(tm)
}

type int16Uint32Map map[int16]uint32

//PackMap packs TypeMap as msgpack.
func (tm int16Uint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16Uint32Map) Len() int {
	return len(tm)
}

type int16Float32Map map[int16]float32

//PackMap packs TypeMap as msgpack.
func (tm int16Float32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16Float32Map) Len() int {
	return len(tm)
}

type int16Float64Map map[int16]float64

//PackMap packs TypeMap as msgpack.
func (tm int16Float64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16Float64Map) Len() int {
	return len(tm)
}

type int16InterfaceMap map[int16]interface{}

//PackMap packs TypeMap as msgpack.
func (tm int16InterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16InterfaceMap) Len() int {
	return len(tm)
}

type int32StringMap map[int32]string

//PackMap packs TypeMap as msgpack.
func (tm int32StringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32StringMap) Len() int {
	return len(tm)
}

type int32IntMap map[int32]int

//PackMap packs TypeMap as msgpack.
func (tm int32IntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32IntMap) Len() int {
	return len(tm)
}

type int32Int8Map map[int32]int8

//PackMap packs TypeMap as msgpack.
func (tm int32Int8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32Int8Map) Len() int {
	return len(tm)
}

type int32Int16Map map[int32]int16

//PackMap packs TypeMap as msgpack.
func (tm int32Int16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32Int16Map) Len() int {
	return len(tm)
}

type int32Int32Map map[int32]int32

//PackMap packs TypeMap as msgpack.
func (tm int32Int32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32Int32Map) Len() int {
	return len(tm)
}

type int32Int64Map map[int32]int64

//PackMap packs TypeMap as msgpack.
func (tm int32Int64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32Int64Map) Len() int {
	return len(tm)
}

type int32Uint16Map map[int32]uint16

//PackMap packs TypeMap as msgpack.
func (tm int32Uint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32Uint16Map) Len() int {
	return len(tm)
}

type int32Uint32Map map[int32]uint32

//PackMap packs TypeMap as msgpack.
func (tm int32Uint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32Uint32Map) Len() int {
	return len(tm)
}

type int32Float32Map map[int32]float32

//PackMap packs TypeMap as msgpack.
func (tm int32Float32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32Float32Map) Len() int {
	return len(tm)
}

type int32Float64Map map[int32]float64

//PackMap packs TypeMap as msgpack.
func (tm int32Float64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32Float64Map) Len() int {
	return len(tm)
}

type int32InterfaceMap map[int32]interface{}

//PackMap packs TypeMap as msgpack.
func (tm int32InterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32InterfaceMap) Len() int {
	return len(tm)
}

type int64StringMap map[int64]string

//PackMap packs TypeMap as msgpack.
func (tm int64StringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64StringMap) Len() int {
	return len(tm)
}

type int64IntMap map[int64]int

//PackMap packs TypeMap as msgpack.
func (tm int64IntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64IntMap) Len() int {
	return len(tm)
}

type int64Int8Map map[int64]int8

//PackMap packs TypeMap as msgpack.
func (tm int64Int8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64Int8Map) Len() int {
	return len(tm)
}

type int64Int16Map map[int64]int16

//PackMap packs TypeMap as msgpack.
func (tm int64Int16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64Int16Map) Len() int {
	return len(tm)
}

type int64Int32Map map[int64]int32

//PackMap packs TypeMap as msgpack.
func (tm int64Int32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64Int32Map) Len() int {
	return len(tm)
}

type int64Int64Map map[int64]int64

//PackMap packs TypeMap as msgpack.
func (tm int64Int64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64Int64Map) Len() int {
	return len(tm)
}

type int64Uint16Map map[int64]uint16

//PackMap packs TypeMap as msgpack.
func (tm int64Uint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64Uint16Map) Len() int {
	return len(tm)
}

type int64Uint32Map map[int64]uint32

//PackMap packs TypeMap as msgpack.
func (tm int64Uint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64Uint32Map) Len() int {
	return len(tm)
}

type int64Float32Map map[int64]float32

//PackMap packs TypeMap as msgpack.
func (tm int64Float32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64Float32Map) Len() int {
	return len(tm)
}

type int64Float64Map map[int64]float64

//PackMap packs TypeMap as msgpack.
func (tm int64Float64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64Float64Map) Len() int {
	return len(tm)
}

type int64InterfaceMap map[int64]interface{}

//PackMap packs TypeMap as msgpack.
func (tm int64InterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64InterfaceMap) Len() int {
	return len(tm)
}

type uint16StringMap map[uint16]string

//PackMap packs TypeMap as msgpack.
func (tm uint16StringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16StringMap) Len() int {
	return len(tm)
}

type uint16IntMap map[uint16]int

//PackMap packs TypeMap as msgpack.
func (tm uint16IntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16IntMap) Len() int {
	return len(tm)
}

type uint16Int8Map map[uint16]int8

//PackMap packs TypeMap as msgpack.
func (tm uint16Int8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16Int8Map) Len() int {
	return len(tm)
}

type uint16Int16Map map[uint16]int16

//PackMap packs TypeMap as msgpack.
func (tm uint16Int16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16Int16Map) Len() int {
	return len(tm)
}

type uint16Int32Map map[uint16]int32

//PackMap packs TypeMap as msgpack.
func (tm uint16Int32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16Int32Map) Len() int {
	return len(tm)
}

type uint16Int64Map map[uint16]int64

//PackMap packs TypeMap as msgpack.
func (tm uint16Int64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16Int64Map) Len() int {
	return len(tm)
}

type uint16Uint16Map map[uint16]uint16

//PackMap packs TypeMap as msgpack.
func (tm uint16Uint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16Uint16Map) Len() int {
	return len(tm)
}

type uint16Uint32Map map[uint16]uint32

//PackMap packs TypeMap as msgpack.
func (tm uint16Uint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16Uint32Map) Len() int {
	return len(tm)
}

type uint16Float32Map map[uint16]float32

//PackMap packs TypeMap as msgpack.
func (tm uint16Float32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16Float32Map) Len() int {
	return len(tm)
}

type uint16Float64Map map[uint16]float64

//PackMap packs TypeMap as msgpack.
func (tm uint16Float64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16Float64Map) Len() int {
	return len(tm)
}

type uint16InterfaceMap map[uint16]interface{}

//PackMap packs TypeMap as msgpack.
func (tm uint16InterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16InterfaceMap) Len() int {
	return len(tm)
}

type uint32StringMap map[uint32]string

//PackMap packs TypeMap as msgpack.
func (tm uint32StringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32StringMap) Len() int {
	return len(tm)
}

type uint32IntMap map[uint32]int

//PackMap packs TypeMap as msgpack.
func (tm uint32IntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32IntMap) Len() int {
	return len(tm)
}

type uint32Int8Map map[uint32]int8

//PackMap packs TypeMap as msgpack.
func (tm uint32Int8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32Int8Map) Len() int {
	return len(tm)
}

type uint32Int16Map map[uint32]int16

//PackMap packs TypeMap as msgpack.
func (tm uint32Int16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32Int16Map) Len() int {
	return len(tm)
}

type uint32Int32Map map[uint32]int32

//PackMap packs TypeMap as msgpack.
func (tm uint32Int32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32Int32Map) Len() int {
	return len(tm)
}

type uint32Int64Map map[uint32]int64

//PackMap packs TypeMap as msgpack.
func (tm uint32Int64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32Int64Map) Len() int {
	return len(tm)
}

type uint32Uint16Map map[uint32]uint16

//PackMap packs TypeMap as msgpack.
func (tm uint32Uint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32Uint16Map) Len() int {
	return len(tm)
}

type uint32Uint32Map map[uint32]uint32

//PackMap packs TypeMap as msgpack.
func (tm uint32Uint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32Uint32Map) Len() int {
	return len(tm)
}

type uint32Float32Map map[uint32]float32

//PackMap packs TypeMap as msgpack.
func (tm uint32Float32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32Float32Map) Len() int {
	return len(tm)
}

type uint32Float64Map map[uint32]float64

//PackMap packs TypeMap as msgpack.
func (tm uint32Float64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32Float64Map) Len() int {
	return len(tm)
}

type uint32InterfaceMap map[uint32]interface{}

//PackMap packs TypeMap as msgpack.
func (tm uint32InterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32InterfaceMap) Len() int {
	return len(tm)
}

type float32StringMap map[float32]string

//PackMap packs TypeMap as msgpack.
func (tm float32StringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32StringMap) Len() int {
	return len(tm)
}

type float32IntMap map[float32]int

//PackMap packs TypeMap as msgpack.
func (tm float32IntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32IntMap) Len() int {
	return len(tm)
}

type float32Int8Map map[float32]int8

//PackMap packs TypeMap as msgpack.
func (tm float32Int8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32Int8Map) Len() int {
	return len(tm)
}

type float32Int16Map map[float32]int16

//PackMap packs TypeMap as msgpack.
func (tm float32Int16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32Int16Map) Len() int {
	return len(tm)
}

type float32Int32Map map[float32]int32

//PackMap packs TypeMap as msgpack.
func (tm float32Int32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32Int32Map) Len() int {
	return len(tm)
}

type float32Int64Map map[float32]int64

//PackMap packs TypeMap as msgpack.
func (tm float32Int64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32Int64Map) Len() int {
	return len(tm)
}

type float32Uint16Map map[float32]uint16

//PackMap packs TypeMap as msgpack.
func (tm float32Uint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32Uint16Map) Len() int {
	return len(tm)
}

type float32Uint32Map map[float32]uint32

//PackMap packs TypeMap as msgpack.
func (tm float32Uint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32Uint32Map) Len() int {
	return len(tm)
}

type float32Float32Map map[float32]float32

//PackMap packs TypeMap as msgpack.
func (tm float32Float32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32Float32Map) Len() int {
	return len(tm)
}

type float32Float64Map map[float32]float64

//PackMap packs TypeMap as msgpack.
func (tm float32Float64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32Float64Map) Len() int {
	return len(tm)
}

type float32InterfaceMap map[float32]interface{}

//PackMap packs TypeMap as msgpack.
func (tm float32InterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32InterfaceMap) Len() int {
	return len(tm)
}

type float64StringMap map[float64]string

//PackMap packs TypeMap as msgpack.
func (tm float64StringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64StringMap) Len() int {
	return len(tm)
}

type float64IntMap map[float64]int

//PackMap packs TypeMap as msgpack.
func (tm float64IntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64IntMap) Len() int {
	return len(tm)
}

type float64Int8Map map[float64]int8

//PackMap packs TypeMap as msgpack.
func (tm float64Int8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64Int8Map) Len() int {
	return len(tm)
}

type float64Int16Map map[float64]int16

//PackMap packs TypeMap as msgpack.
func (tm float64Int16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64Int16Map) Len() int {
	return len(tm)
}

type float64Int32Map map[float64]int32

//PackMap packs TypeMap as msgpack.
func (tm float64Int32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64Int32Map) Len() int {
	return len(tm)
}

type float64Int64Map map[float64]int64

//PackMap packs TypeMap as msgpack.
func (tm float64Int64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64Int64Map) Len() int {
	return len(tm)
}

type float64Uint16Map map[float64]uint16

//PackMap packs TypeMap as msgpack.
func (tm float64Uint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64Uint16Map) Len() int {
	return len(tm)
}

type float64Uint32Map map[float64]uint32

//PackMap packs TypeMap as msgpack.
func (tm float64Uint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64Uint32Map) Len() int {
	return len(tm)
}

type float64Float32Map map[float64]float32

//PackMap packs TypeMap as msgpack.
func (tm float64Float32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64Float32Map) Len() int {
	return len(tm)
}

type float64Float64Map map[float64]float64

//PackMap packs TypeMap as msgpack.
func (tm float64Float64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64Float64Map) Len() int {
	return len(tm)
}

type float64InterfaceMap map[float64]interface{}

//PackMap packs TypeMap as msgpack.
func (tm float64InterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64InterfaceMap) Len() int {
	return len(tm)
}

type stringUint64Map map[string]uint64

//PackMap packs TypeMap as msgpack.
func (tm stringUint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackString(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm stringUint64Map) Len() int {
	return len(tm)
}

type intUint64Map map[int]uint64

//PackMap packs TypeMap as msgpack.
func (tm intUint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm intUint64Map) Len() int {
	return len(tm)
}

type int8Uint64Map map[int8]uint64

//PackMap packs TypeMap as msgpack.
func (tm int8Uint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int8Uint64Map) Len() int {
	return len(tm)
}

type int16Uint64Map map[int16]uint64

//PackMap packs TypeMap as msgpack.
func (tm int16Uint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int16Uint64Map) Len() int {
	return len(tm)
}

type int32Uint64Map map[int32]uint64

//PackMap packs TypeMap as msgpack.
func (tm int32Uint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int32Uint64Map) Len() int {
	return len(tm)
}

type int64Uint64Map map[int64]uint64

//PackMap packs TypeMap as msgpack.
func (tm int64Uint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm int64Uint64Map) Len() int {
	return len(tm)
}

type uint16Uint64Map map[uint16]uint64

//PackMap packs TypeMap as msgpack.
func (tm uint16Uint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint16Uint64Map) Len() int {
	return len(tm)
}

type uint32Uint64Map map[uint32]uint64

//PackMap packs TypeMap as msgpack.
func (tm uint32Uint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackInt64(buf, int64(k))
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint32Uint64Map) Len() int {
	return len(tm)
}

type float32Uint64Map map[float32]uint64

//PackMap packs TypeMap as msgpack.
func (tm float32Uint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat32(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float32Uint64Map) Len() int {
	return len(tm)
}

type float64Uint64Map map[float64]uint64

//PackMap packs TypeMap as msgpack.
func (tm float64Uint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackFloat64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm float64Uint64Map) Len() int {
	return len(tm)
}

type uint64StringMap map[uint64]string

//PackMap packs TypeMap as msgpack.
func (tm uint64StringMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackString(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64StringMap) Len() int {
	return len(tm)
}

type uint64IntMap map[uint64]int

//PackMap packs TypeMap as msgpack.
func (tm uint64IntMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64IntMap) Len() int {
	return len(tm)
}

type uint64Int8Map map[uint64]int8

//PackMap packs TypeMap as msgpack.
func (tm uint64Int8Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64Int8Map) Len() int {
	return len(tm)
}

type uint64Int16Map map[uint64]int16

//PackMap packs TypeMap as msgpack.
func (tm uint64Int16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64Int16Map) Len() int {
	return len(tm)
}

type uint64Int32Map map[uint64]int32

//PackMap packs TypeMap as msgpack.
func (tm uint64Int32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64Int32Map) Len() int {
	return len(tm)
}

type uint64Int64Map map[uint64]int64

//PackMap packs TypeMap as msgpack.
func (tm uint64Int64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64Int64Map) Len() int {
	return len(tm)
}

type uint64Uint16Map map[uint64]uint16

//PackMap packs TypeMap as msgpack.
func (tm uint64Uint16Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64Uint16Map) Len() int {
	return len(tm)
}

type uint64Uint32Map map[uint64]uint32

//PackMap packs TypeMap as msgpack.
func (tm uint64Uint32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackInt64(buf, int64(v))
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64Uint32Map) Len() int {
	return len(tm)
}

type uint64Uint64Map map[uint64]uint64

//PackMap packs TypeMap as msgpack.
func (tm uint64Uint64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackUInt64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64Uint64Map) Len() int {
	return len(tm)
}

type uint64Float32Map map[uint64]float32

//PackMap packs TypeMap as msgpack.
func (tm uint64Float32Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat32(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64Float32Map) Len() int {
	return len(tm)
}

type uint64Float64Map map[uint64]float64

//PackMap packs TypeMap as msgpack.
func (tm uint64Float64Map) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = PackFloat64(buf, v)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64Float64Map) Len() int {
	return len(tm)
}

type uint64InterfaceMap map[uint64]interface{}

//PackMap packs TypeMap as msgpack.
func (tm uint64InterfaceMap) PackMap(buf BufferEx) (int, error) {
	size := 0
	for k, v := range tm {
		n, err := PackUInt64(buf, k)
		size += n
		if err != nil {
			return size, err
		}

		n, err = packObject(buf, v, false)
		size += n
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

// Len return the length of TypeSlice
func (tm uint64InterfaceMap) Len() int {
	return len(tm)
}
