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

// Stream is the data type used as a stream by the lua instances
type Stream struct {
	s chan interface{}
}

const luaLuaStreamTypeName = "LuaStream"

// Registers my luaStream type to given L.
func registerLuaStreamType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaLuaStreamTypeName)

	L.SetGlobal("stream", mt)

	// static attributes
	L.SetField(mt, "__call", L.NewFunction(newStream))
	L.SetField(mt, "read", L.NewFunction(luaStreamRead))
	L.SetField(mt, "write", L.NewFunction(luaStreamWrite))
	L.SetField(mt, "readable", L.NewFunction(luaStreamReadable))
	L.SetField(mt, "writeable", L.NewFunction(luaStreamWriteable))

	// methods
	L.SetFuncs(mt, map[string]lua.LGFunction{
		"__tostring": luaStreamToString,
	})

	L.SetMetatable(mt, mt)
}

// NewStream creates a LuaStream
func NewStream(L *lua.LState, stream chan interface{}) *lua.LUserData {
	luaStream := &Stream{s: stream}
	ud := L.NewUserData()
	ud.Value = luaStream
	L.SetMetatable(ud, L.GetTypeMetatable(luaLuaStreamTypeName))
	return ud
}

func newStream(L *lua.LState) int {
	luaStream := &Stream{s: make(chan interface{}, 64)}
	ud := L.NewUserData()
	ud.Value = luaStream
	L.SetMetatable(ud, L.GetTypeMetatable(luaLuaStreamTypeName))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *LuaStream and returns this *LuaStream.
func checkLuaStream(L *lua.LState) *Stream {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Stream); ok {
		return v
	}
	L.ArgError(1, "luaSteam expected")
	return nil
}

func luaStreamToString(L *lua.LState) int {
	p := checkLuaStream(L)
	if L.GetTop() != 1 {
		L.ArgError(1, "No arguments expected for tostring method")
		return 0
	}
	L.Push(lua.LString(fmt.Sprintf("%v", p.s)))
	return 1
}

func luaStreamRead(L *lua.LState) int {
	p := checkLuaStream(L)
	if L.GetTop() != 1 {
		L.ArgError(1, "No arguments expected for stream:read method")
		return 0
	}

	L.Push(NewValue(L, <-p.s))
	return 1
}

func luaStreamWrite(L *lua.LState) int {
	p := checkLuaStream(L)
	if L.GetTop() != 2 {
		L.ArgError(1, "Only one argument expected for stream:write method")
		return 0
	}

	p.s <- LValueToInterface(L.CheckAny(2))
	return 1
}

func luaStreamReadable(L *lua.LState) int {
	checkLuaStream(L)
	if L.GetTop() != 1 {
		L.ArgError(1, "No arguments expected for readable method")
		return 0
	}
	L.Push(lua.LBool(true))
	return 1
}

func luaStreamWriteable(L *lua.LState) int {
	checkLuaStream(L)
	if L.GetTop() != 1 {
		L.ArgError(1, "No arguments expected for writeable method")
		return 0
	}
	L.Push(lua.LBool(true))
	return 1
}
