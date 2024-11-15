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
	"reflect"
)

func init() {
	newValueReflect = concreteNewValueReflect
}

// if the returned value is nil, the caller will panic
func concreteNewValueReflect(v interface{}) Value {
	// check for array and map
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		l := rv.Len()
		arr := make([]interface{}, l)
		for i := 0; i < l; i++ {
			arr[i] = rv.Index(i).Interface()
		}

		return NewListValue(arr)
	case reflect.Map:
		l := rv.Len()
		amap := make(map[interface{}]interface{}, l)
		for _, i := range rv.MapKeys() {
			amap[i.Interface()] = rv.MapIndex(i).Interface()
		}

		return NewMapValue(amap)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewLongValue(reflect.ValueOf(v).Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return NewLongValue(int64(reflect.ValueOf(v).Uint()))
	case reflect.String:
		return NewStringValue(rv.String())
	}

	return nil
}
