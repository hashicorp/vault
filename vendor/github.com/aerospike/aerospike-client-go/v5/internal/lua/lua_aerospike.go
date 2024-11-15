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
	"github.com/aerospike/aerospike-client-go/v5/logger"
	lua "github.com/yuin/gopher-lua"
)

const luaLuaAerospikeTypeName = "LuaAerospike"

// Registers my luaAerospike type to given L.
func registerLuaAerospikeType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaLuaAerospikeTypeName)

	L.SetGlobal("aerospike", mt)

	// static attributes
	L.SetField(mt, "log", L.NewFunction(luaAerospikeLog))

	L.SetMetatable(mt, mt)
}

func luaAerospikeLog(L *lua.LState) int {
	if L.GetTop() < 2 || L.GetTop() > 3 {
		L.ArgError(1, "2 arguments are expected for aerospike:log method")
		return 0
	}

	// account for calling it on a table
	paramIdx := 1
	if L.GetTop() == 3 {
		paramIdx = 2
	}

	level := L.CheckInt(paramIdx)
	str := L.CheckString(paramIdx + 1)

	switch level {
	case 1:
		logger.Logger.Warn(str)
	case 2:
		logger.Logger.Info(str)
	case 3, 4:
		logger.Logger.Debug(str)
	}

	return 0
}
