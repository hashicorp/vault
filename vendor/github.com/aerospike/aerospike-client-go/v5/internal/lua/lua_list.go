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

// List is the list data type to be used with a Lua instance
type List struct {
	l []interface{}
}

const luaLuaListTypeName = "LuaList"

// Registers my luaList type to given L.
func registerLuaListType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaLuaListTypeName)

	// List package
	L.SetGlobal("List", mt)

	// static attributes

	L.SetMetatable(mt, mt)

	// list package
	mt = L.NewTypeMetatable(luaLuaListTypeName)
	L.SetGlobal("list", mt)

	// static attributes
	L.SetField(mt, "__call", L.NewFunction(newLuaList))
	L.SetField(mt, "create", L.NewFunction(createLuaList))

	L.SetField(mt, "size", L.NewFunction(luaListSize))
	L.SetField(mt, "insert", L.NewFunction(luaListInsert))
	L.SetField(mt, "append", L.NewFunction(luaListAppend))
	L.SetField(mt, "prepend", L.NewFunction(luaListPrepend))
	L.SetField(mt, "take", L.NewFunction(luaListTake))
	L.SetField(mt, "remove", L.NewFunction(luaListRemove))
	L.SetField(mt, "drop", L.NewFunction(luaListDrop))
	L.SetField(mt, "trim", L.NewFunction(luaListTrim))
	L.SetField(mt, "clone", L.NewFunction(luaListClone))
	L.SetField(mt, "concat", L.NewFunction(luaListConcat))
	L.SetField(mt, "merge", L.NewFunction(luaListMerge))
	L.SetField(mt, "iterator", L.NewFunction(luaListIterator))

	// methods
	L.SetFuncs(mt, map[string]lua.LGFunction{
		"__index":    luaListIndex,
		"__newindex": luaListNewIndex,
		"__len":      luaListLen,
		"__tostring": luaListToString,
	})

	L.SetMetatable(mt, mt)
}

// Constructor
func createLuaList(L *lua.LState) int {
	if L.GetTop() == 0 {
		luaList := &List{l: []interface{}{}}
		ud := L.NewUserData()
		ud.Value = luaList
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaListTypeName))
		L.Push(ud)
		return 1
	} else if L.GetTop() == 1 || L.GetTop() == 2 {
		cp := L.CheckInt(1)
		l := make([]interface{}, 0, cp)

		luaList := &List{l: l}
		ud := L.NewUserData()
		ud.Value = luaList
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaListTypeName))
		L.Push(ud)
		return 1
	}
	L.ArgError(1, "Only one argument expected for list#create method")
	return 0
}

// Constructor
func newLuaList(L *lua.LState) int {
	if L.GetTop() == 1 {
		luaList := &List{l: []interface{}{}}
		ud := L.NewUserData()
		ud.Value = luaList
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaListTypeName))
		L.Push(ud)
		return 1
	} else if L.GetTop() == 2 {
		t := L.CheckTable(2)
		l := make([]interface{}, t.Len())
		for i := 1; i <= t.Len(); i++ {
			l[i-1] = LValueToInterface(t.RawGetInt(i))
		}

		luaList := &List{l: l}
		ud := L.NewUserData()
		ud.Value = luaList
		L.SetMetatable(ud, L.GetTypeMetatable(luaLuaListTypeName))
		L.Push(ud)
		return 1
	}
	L.ArgError(1, "Only one argument expected for list#create method")
	return 0
}

// Checks whether the first lua argument is a *LUserData with *LuaList and returns this *LuaList.
func checkLuaList(L *lua.LState, arg int) *List {
	ud := L.CheckUserData(arg)
	if v, ok := ud.Value.(*List); ok {
		return v
	}
	L.ArgError(1, "luaList expected")
	return nil
}

func luaListRemove(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for remove method")
		return 0
	}
	index := L.CheckInt(2) - 1

	if index < 0 || index >= len(p.l) {
		L.ArgError(1, "index out of range for list#remove")
		return 0
	}

	for i := index; i < len(p.l)-1; i++ {
		p.l[i] = p.l[i+1]
	}
	p.l = p.l[:len(p.l)-1]

	return 0
}

func luaListInsert(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 3 {
		L.ArgError(1, "Only two arguments expected for insert method")
		return 0
	}
	index := L.CheckInt(2)
	value := LValueToInterface(L.CheckAny(3))

	if cap(p.l) > len(p.l) {
		for i := len(p.l); i >= index; i-- {
			p.l[i] = p.l[i-1]
		}
		p.l[index-1] = value
	} else {
		ln := len(p.l) * 2
		if ln > 256 {
			ln = 256
		}
		newList := make([]interface{}, len(p.l)+1, ln)

		copy(newList, p.l[:index-1])
		newList[index-1] = value
		copy(newList[index:], p.l[index-1:len(p.l)])
		p.l = newList
	}

	return 0
}

func luaListAppend(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for append method")
		return 0
	}
	value := LValueToInterface(L.CheckAny(2))
	p.l = append(p.l, value)

	return 0
}

func luaListPrepend(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for append method")
		return 0
	}
	value := LValueToInterface(L.CheckAny(2))

	if cap(p.l) > len(p.l) {
		p.l = append(p.l, nil)
		for i := len(p.l) - 1; i > 0; i-- {
			p.l[i] = p.l[i-1]
		}
		p.l[0] = value
	} else {
		ln := len(p.l) * 2
		if ln > 256 {
			ln = 256
		}
		newList := make([]interface{}, len(p.l)+1, ln)

		copy(newList[1:], p.l)
		newList[0] = value
		p.l = newList
	}

	return 0
}

func luaListTake(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for take method")
		return 0
	}

	count := L.CheckInt(2)
	items := p.l
	if count <= len(p.l) {
		items = p.l[:count]
	}

	luaList := &List{l: items}
	ud := L.NewUserData()
	ud.Value = luaList
	L.SetMetatable(ud, L.GetTypeMetatable(luaLuaListTypeName))
	L.Push(ud)
	return 1
}

func luaListDrop(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for take method")
		return 0
	}

	count := L.CheckInt(2)
	var items []interface{}
	if count < len(p.l) {
		items = p.l[count:]
	} else {
		items = []interface{}{}
	}

	luaList := &List{l: items}
	ud := L.NewUserData()
	ud.Value = luaList
	L.SetMetatable(ud, L.GetTypeMetatable(luaLuaListTypeName))
	L.Push(ud)
	return 1
}

func luaListTrim(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for list#trim method")
		return 0
	}

	count := L.CheckInt(2)
	p.l = p.l[:count-1]

	return 0
}

func luaListConcat(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for list#concat method")
		return 0
	}

	sp := checkLuaList(L, 2)
	p.l = append(p.l, sp.l...)
	return 0
}

// LuaList#clone()
func luaListClone(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 1 {
		L.ArgError(1, "only one argument expected for list#clone method")
		return 0
	}

	newList := &List{l: make([]interface{}, len(p.l))}
	copy(newList.l, p.l)

	ud := L.NewUserData()
	ud.Value = newList
	L.SetMetatable(ud, L.GetTypeMetatable(luaLuaListTypeName))
	L.Push(ud)
	return 1
}

// LuaList#merge()
func luaListMerge(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for merge method")
		return 0
	}

	sp := checkLuaList(L, 2)

	newList := &List{l: make([]interface{}, 0, len(p.l)+len(sp.l))}
	newList.l = append(newList.l, p.l...)
	newList.l = append(newList.l, sp.l...)

	ud := L.NewUserData()
	ud.Value = newList
	L.SetMetatable(ud, L.GetTypeMetatable(luaLuaListTypeName))
	L.Push(ud)

	return 1
}

func luaListToString(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 1 {
		L.ArgError(1, "No arguments expected for tostring method")
		return 0
	}

	L.Push(lua.LString(fmt.Sprintf("%v", p.l)))
	return 1
}

func luaListSize(L *lua.LState) int {
	p := checkLuaList(L, 1)
	if L.GetTop() != 1 {
		L.ArgError(1, "No arguments expected for __size method")
		return 0
	}
	L.Push(lua.LNumber(len(p.l)))
	return 1
}

func luaListIndex(L *lua.LState) int {
	ref := checkLuaList(L, 1)
	index := L.CheckInt(2)

	if index <= 0 || index > len(ref.l) {
		L.Push(lua.LNil)
		return 1
	}

	item := ref.l[index-1]
	L.Push(NewValue(L, item))
	return 1
}

func luaListNewIndex(L *lua.LState) int {
	ref := checkLuaList(L, 1)
	index := L.CheckInt(2)
	value := L.CheckAny(3)

	ref.l[index-1] = LValueToInterface(value)
	return 0
}

func luaListLen(L *lua.LState) int {
	ref := checkLuaList(L, 1)
	L.Push(lua.LNumber(len(ref.l)))
	return 1
}

func luaListIterator(L *lua.LState) int {
	ref := checkLuaList(L, 1)

	// make an iterator
	idx := 0
	llen := len(ref.l)
	fn := func(L *lua.LState) int {
		if idx < llen {
			L.Push(NewValue(L, ref.l[idx]))
			idx++
			return 1
		}
		return 0
	}
	L.Push(L.NewFunction(fn))
	return 1
}
