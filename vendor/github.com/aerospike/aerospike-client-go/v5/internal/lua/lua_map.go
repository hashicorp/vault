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

	lua "github.com/yuin/gopher-lua"
)

// Map is used internally for the Lua instance
type Map struct {
	m map[interface{}]interface{}
}

const luaLuaMapTypeName = "LuaMap"

func registerLuaMapType(L *lua.LState) {
	// Map package
	mt := L.NewTypeMetatable(luaLuaMapTypeName)

	L.SetGlobal("Map", mt)

	// static attributes
	L.SetField(mt, "__call", L.NewFunction(newLuaMap))

	L.SetField(mt, "create", L.NewFunction(luaMapCreate))

	// methods
	L.SetMetatable(mt, mt)

	// map package
	mt = L.NewTypeMetatable(luaLuaMapTypeName)

	L.SetGlobal("map", mt)

	// static attributes
	L.SetField(mt, "__call", L.NewFunction(newLuaMap))

	L.SetField(mt, "create", L.NewFunction(luaMapCreate))

	L.SetField(mt, "pairs", L.NewFunction(luaMapPairs))
	L.SetField(mt, "size", L.NewFunction(luaMapSize))
	L.SetField(mt, "keys", L.NewFunction(luaMapKeys))
	L.SetField(mt, "values", L.NewFunction(luaMapValues))
	L.SetField(mt, "remove", L.NewFunction(luaMapRemove))
	L.SetField(mt, "clone", L.NewFunction(luaMapClone))
	L.SetField(mt, "merge", L.NewFunction(luaMapMerge))
	L.SetField(mt, "diff", L.NewFunction(luaMapDiff))

	// methods
	L.SetFuncs(mt, map[string]lua.LGFunction{
		"__index":    luaMapIndex,
		"__newindex": luaMapNewIndex,
		"__len":      luaMapSize,
		"__tostring": luaMapToString,
	})

	L.SetMetatable(mt, mt)
}

// Constructor
func luaMapCreate(L *lua.LState) int {
	if L.GetTop() == 1 {
		luaMap := &Map{m: map[interface{}]interface{}{}}
		ud := L.NewUserData()
		ud.Value = luaMap
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaMapTypeName))
		L.Push(ud)
		return 1
	} else if L.GetTop() == 2 {
		L.CheckTable(1)
		sz := L.CheckInt(2)
		luaMap := &Map{m: make(map[interface{}]interface{}, sz)}
		ud := L.NewUserData()
		ud.Value = luaMap
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaMapTypeName))
		L.Push(ud)
		return 1
	}
	L.ArgError(1, "Only one argument expected for map create method")
	return 0
}

func newLuaMap(L *lua.LState) int {
	if L.GetTop() == 1 {
		luaMap := &Map{m: make(map[interface{}]interface{}, 4)}
		ud := L.NewUserData()
		ud.Value = luaMap
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaMapTypeName))
		L.Push(ud)
		return 1
	} else if L.GetTop() == 2 {
		L.CheckTable(1)
		t := L.CheckTable(2)
		m := make(map[interface{}]interface{}, t.Len())
		t.ForEach(func(k, v lua.LValue) { m[LValueToInterface(k)] = LValueToInterface(v) })

		luaMap := &Map{m: m}
		ud := L.NewUserData()
		ud.Value = luaMap
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaMapTypeName))
		L.Push(ud)
		return 1
	}
	L.ArgError(1, "Only one argument expected for map create method")
	return 0
}

// Checks whether the first lua argument is a *LUserData with *LuaMap and returns this *LuaMap.
func checkLuaMap(L *lua.LState, arg int) *Map {
	ud := L.CheckUserData(arg)
	if v, ok := ud.Value.(*Map); ok {
		return v
	}
	L.ArgError(1, "luaMap expected")
	return nil
}

func luaMapRemove(L *lua.LState) int {
	p := checkLuaMap(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for remove method")
		return 0
	}
	key := L.CheckAny(2)

	delete(p.m, LValueToInterface(key))
	return 0
}

func luaMapClone(L *lua.LState) int {
	p := checkLuaMap(L, 1)
	if L.GetTop() != 1 {
		L.ArgError(1, "No arguments expected for clone method")
		return 0
	}

	newMap := &Map{m: make(map[interface{}]interface{}, len(p.m))}
	for k, v := range p.m {
		newMap.m[k] = v
	}

	ud := L.NewUserData()
	ud.Value = newMap
	L.SetMetatable(ud, L.GetTypeMetatable(luaLuaMapTypeName))
	L.Push(ud)
	return 1
}

func luaMapMerge(L *lua.LState) int {
	p := checkLuaMap(L, 1)
	if L.GetTop() < 2 || L.GetTop() > 3 {
		L.ArgError(1, "Only 2 or 3 argument expected for merge method")
		return 0
	}

	if L.GetTop() == 2 {
		sp := checkLuaMap(L, 2)

		newMap := &Map{m: make(map[interface{}]interface{}, len(p.m)+len(sp.m))}
		for k, v := range p.m {
			newMap.m[k] = v
		}

		for k, v := range sp.m {
			newMap.m[k] = v
		}

		ud := L.NewUserData()
		ud.Value = newMap
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaMapTypeName))
		L.Push(ud)
	} else {
		sp := checkLuaMap(L, 2)
		fn := L.CheckFunction(3)

		newMap := &Map{m: make(map[interface{}]interface{}, len(p.m)+len(sp.m))}
		for k, v := range p.m {
			if v2, exists := sp.m[k]; exists {
				L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true, Handler: nil}, NewValue(L, v), NewValue(L, v2))
				ret := L.CheckAny(-1)
				L.Pop(1) // remove received value
				newMap.m[k] = LValueToInterface(ret)
			} else {
				newMap.m[k] = v
			}
		}

		for k, v := range sp.m {
			// only add keys that haven't been processed already
			if _, exists := newMap.m[k]; !exists {
				newMap.m[k] = v
			}
		}

		ud := L.NewUserData()
		ud.Value = newMap
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaMapTypeName))
		L.Push(ud)
	}

	return 1
}

func luaMapDiff(L *lua.LState) int {
	p := checkLuaMap(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for diff method")
		return 0
	}

	sp := checkLuaMap(L, 2)

	newMap := &Map{m: make(map[interface{}]interface{}, len(p.m)+len(sp.m))}

	for k, v := range p.m {
		if _, exists := sp.m[k]; !exists {
			newMap.m[k] = v
		}
	}

	for k, v := range sp.m {
		if _, exists := p.m[k]; !exists {
			newMap.m[k] = v
		}
	}

	ud := L.NewUserData()
	ud.Value = newMap
	L.SetMetatable(ud, L.GetTypeMetatable(luaLuaMapTypeName))
	L.Push(ud)

	return 1
}

func luaMapToString(L *lua.LState) int {
	p := checkLuaMap(L, 1)
	if L.GetTop() != 1 {
		L.ArgError(1, "No arguments expected for tostring method")
		return 0
	}
	L.Push(lua.LString(fmt.Sprintf("%v", p.m)))
	return 1
}

func luaMapSize(L *lua.LState) int {
	p := checkLuaMap(L, 1)
	if L.GetTop() != 1 {
		L.ArgError(1, "No arguments expected for __size method")
		return 0
	}
	L.Push(lua.LNumber(len(p.m)))
	return 1
}

func luaMapIndex(L *lua.LState) int {
	ref := checkMap(L)
	key := LValueToInterface(L.CheckAny(2))

	v := ref.m[key]
	if v == nil {
		v = lua.LNil
	}

	L.Push(NewValue(L, v))
	return 1
}

func luaMapNewIndex(L *lua.LState) int {
	ref := checkMap(L)
	key := LValueToInterface(L.CheckAny(2))
	value := LValueToInterface(L.CheckAny(3))

	ref.m[key] = value
	return 0
}

func luaMapLen(L *lua.LState) int {
	ref := checkMap(L)
	L.Push(lua.LNumber(len(ref.m)))
	return 1
}

func luaMapPairs(L *lua.LState) int {
	ref := checkMap(L)

	// make an iterator
	iter := make(chan *struct{ k, v interface{} })

	go func() {
		for k, v := range ref.m {
			iter <- &struct{ k, v interface{} }{k, v}
		}
		close(iter)
	}()

	fn := func(L *lua.LState) int {
		tuple := <-iter
		if tuple == nil {
			return 0
		}

		L.Push(NewValue(L, tuple.k))
		L.Push(NewValue(L, tuple.v))
		return 2
	}
	L.Push(L.NewFunction(fn))
	return 1
}

func luaMapKeys(L *lua.LState) int {
	ref := checkMap(L)

	// make an iterator
	iter := make(chan interface{})

	go func() {
		for k := range ref.m {
			iter <- k
		}
		close(iter)
	}()

	fn := func(L *lua.LState) int {
		tuple := <-iter
		if tuple == nil {
			return 0
		}

		L.Push(NewValue(L, tuple))
		return 1
	}
	L.Push(L.NewFunction(fn))
	return 1
}

func luaMapValues(L *lua.LState) int {
	ref := checkMap(L)

	// make an iterator
	iter := make(chan interface{})

	go func() {
		for _, v := range ref.m {
			iter <- v
		}
		close(iter)
	}()

	fn := func(L *lua.LState) int {
		tuple := <-iter
		if tuple == nil {
			return 0
		}

		L.Push(NewValue(L, tuple))
		return 1
	}
	L.Push(L.NewFunction(fn))
	return 1
}

func luaMapEq(L *lua.LState) int {
	map1 := checkMap(L)
	map2 := checkMap(L)
	L.Push(lua.LBool(map1 == map2))
	return 1
}

func checkMap(L *lua.LState) *Map {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Map); ok {
		return v
	}
	L.ArgError(1, "luaMap expected")
	return nil
}
