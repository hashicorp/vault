// +build !app_engine

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

package lua

import (
	"fmt"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

// NewValue creates a value from interface{} in the interpreter
func NewValue(L *lua.LState, value interface{}) lua.LValue {
	// Nils should return immediately
	if value == nil {
		return lua.LNil
	}

	// if it is a LValue already, return it without delay
	if lval, ok := value.(lua.LValue); ok {
		return lval
	}

	switch v := value.(type) {
	case string:
		return lua.LString(v)
	case int:
		return lua.LNumber(float64(v))
	case uint:
		return lua.LNumber(float64(v))
	case int8:
		return lua.LNumber(float64(v))
	case uint8:
		return lua.LNumber(float64(v))
	case int16:
		return lua.LNumber(float64(v))
	case uint16:
		return lua.LNumber(float64(v))
	case int32:
		return lua.LNumber(float64(v))
	case uint32:
		return lua.LNumber(float64(v))
	case int64:
		return lua.LNumber(float64(v))
	case uint64:
		return lua.LNumber(float64(v))
	case float32:
		return lua.LNumber(float64(v))
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case map[interface{}]interface{}:
		luaMap := &Map{m: v}
		ud := L.NewUserData()
		ud.Value = luaMap
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaMapTypeName))
		return ud

	case []interface{}:
		luaList := &List{l: v}
		ud := L.NewUserData()
		ud.Value = luaList
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaListTypeName))
		return ud
	}

	// check for array and map
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		l := rv.Len()
		arr := make([]interface{}, l)
		for i := 0; i < l; i++ {
			arr[i] = rv.Index(i).Interface()
		}

		return NewValue(L, arr)
	case reflect.Map:
		l := rv.Len()
		amap := make(map[interface{}]interface{}, l)
		for _, i := range rv.MapKeys() {
			amap[i.Interface()] = rv.MapIndex(i).Interface()
		}

		return NewValue(L, amap)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewValue(L, reflect.ValueOf(value).Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return NewValue(L, int64(reflect.ValueOf(value).Uint()))
	case reflect.String:
		return NewValue(L, rv.String())
	case reflect.Float32, reflect.Float64:
		return NewValue(L, rv.Float())
	case reflect.Bool:
		return NewValue(L, rv.Bool())
	}

	panic(fmt.Sprintf("unrecognized data type for lua: %#v\n", value))
}

// LValueToInterface converts a generic LValue to a native type
func LValueToInterface(val lua.LValue) interface{} {
	switch val.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return lua.LVAsBool(val)
	case lua.LTNumber:
		return float64(lua.LVAsNumber(val))
	case lua.LTString:
		return lua.LVAsString(val)
	case lua.LTUserData:
		ud := val.(*lua.LUserData).Value
		switch v := ud.(type) {
		case *Map:
			return v.m
		case *List:
			return v.l
		default:
			return v
		}

	case lua.LTTable:
		t := val.(*lua.LTable)
		m := make(map[interface{}]interface{}, t.Len())
		t.ForEach(func(k, v lua.LValue) { m[k] = v })
		return m
	default:
		panic(fmt.Sprintf("unrecognized data type %#v", val))
	}
}

func allToString(L *lua.LState) int {
	ud := L.CheckUserData(1)
	value := ud.Value
	if stringer, ok := value.(fmt.Stringer); ok {
		L.Push(lua.LString(stringer.String()))
	} else {
		L.Push(lua.LString(fmt.Sprintf("%v", value)))
	}
	return 1
}
