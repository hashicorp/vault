// +build !as_performance

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

	"github.com/aerospike/aerospike-client-go/v5/types"
)

func init() {
	packObjectReflect = concretePackObjectReflect
}

func concretePackObjectReflect(cmd BufferEx, obj interface{}, mapKey bool) (int, Error) {
	// check for array and map
	rv := reflect.ValueOf(obj)
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Array, reflect.Slice:
		if mapKey && reflect.TypeOf(obj).Kind() == reflect.Slice {
			return 0, newError(types.SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", obj))
		}
		// pack bounded array of bytes differently
		if reflect.TypeOf(obj).Kind() == reflect.Array && reflect.TypeOf(obj).Elem().Kind() == reflect.Uint8 {
			l := rv.Len()
			arr := make([]byte, l)
			for i := 0; i < l; i++ {
				arr[i] = rv.Index(i).Interface().(uint8)
			}
			return packBytes(cmd, arr)
		}

		l := rv.Len()
		arr := make([]interface{}, l)
		for i := 0; i < l; i++ {
			arr[i] = rv.Index(i).Interface()
		}
		return packIfcList(cmd, arr)
	case reflect.Map:
		if mapKey {
			return 0, newError(types.SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", obj))
		}
		l := rv.Len()
		amap := make(map[interface{}]interface{}, l)
		for _, i := range rv.MapKeys() {
			amap[i.Interface()] = rv.MapIndex(i).Interface()
		}
		return packIfcMap(cmd, amap)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return packObject(cmd, rv.Int(), false)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return packObject(cmd, rv.Uint(), false)
	case reflect.Bool:
		return packObject(cmd, rv.Bool(), false)
	case reflect.String:
		return packObject(cmd, rv.String(), false)
	case reflect.Float32, reflect.Float64:
		return packObject(cmd, rv.Float(), false)
	}

	return 0, newError(types.SERIALIZE_ERROR, fmt.Sprintf("Type `%#v` not supported to pack.", obj))
}
