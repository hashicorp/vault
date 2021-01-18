// +build !app_engine

// Copyright 2013-2020 Aerospike, Inc.
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
	"sync"

	. "github.com/aerospike/aerospike-client-go/internal/atomic"
	lualib "github.com/aerospike/aerospike-client-go/internal/lua"
	. "github.com/aerospike/aerospike-client-go/types"
	lua "github.com/yuin/gopher-lua"
)

//--------------------------------------------------------
// Query Aggregate functions (Supported by Aerospike 3+ servers only)
//--------------------------------------------------------

// SetLuaPath sets the Lua interpreter path to files
// This path is used to load UDFs for QueryAggregate command
func SetLuaPath(lpath string) {
	lualib.SetPath(lpath)
}

// QueryAggregate executes a Map/Reduce query and returns the results.
// The query executor puts records on the channel from separate goroutines.
// The caller can concurrently pop records off the channel through the
// Recordset.Records channel.
//
// This method is only supported by Aerospike 3+ servers.
// If the policy is nil, the default relevant policy will be used.
func (clnt *Client) QueryAggregate(policy *QueryPolicy, statement *Statement, packageName, functionName string, functionArgs ...Value) (*Recordset, error) {
	statement.SetAggregateFunction(packageName, functionName, functionArgs, true)

	policy = clnt.getUsableQueryPolicy(policy)

	nodes := clnt.cluster.GetNodes()
	if len(nodes) == 0 {
		return nil, NewAerospikeError(SERVER_NOT_AVAILABLE, "QueryAggregate failed because cluster is empty.")
	}

	clusterKey := int64(0)
	if policy.FailOnClusterChange {
		var err error
		clusterKey, err = queryValidateBegin(nodes[0], statement.Namespace)
		if err != nil {
			return nil, err
		}
	}
	first := NewAtomicBool(true)

	// results channel must be async for performance
	recSet := newRecordset(policy.RecordQueueSize, len(nodes), statement.TaskId)

	// get a lua instance
	luaInstance := lualib.LuaPool.Get().(*lua.LState)
	if luaInstance == nil {
		return nil, fmt.Errorf("Error fetching a lua instance from pool")
	}

	// Input Channel
	inputChan := make(chan interface{}, 4096) // 4096 = number of partitions
	istream := lualib.NewLuaStream(luaInstance, inputChan)

	// Output Channe;
	outputChan := make(chan interface{})
	ostream := lualib.NewLuaStream(luaInstance, outputChan)

	// results channel must be async for performance
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for _, node := range nodes {
		// copy policies to avoid race conditions
		newPolicy := *policy
		command := newQueryAggregateCommand(node, &newPolicy, statement, recSet, clusterKey, first.CompareAndToggle(true))
		command.luaInstance = luaInstance
		command.inputChan = inputChan

		go func() {
			defer wg.Done()
			command.Execute()
		}()
	}

	go func() {
		wg.Wait()
		close(inputChan)
	}()

	go func() {
		// we cannot signal end and close the recordset
		// while processing is still going on
		// We will do it only here, after all processing is over
		defer func() {
			for i := 0; i < len(nodes); i++ {
				recSet.signalEnd()
			}
		}()

		for val := range outputChan {
			recSet.Records <- &Record{Bins: BinMap{"SUCCESS": val}}
		}
	}()

	go func() {
		defer close(outputChan)
		defer luaInstance.Close()

		err := luaInstance.DoFile(lualib.LuaPath() + packageName + ".lua")
		if err != nil {
			recSet.Errors <- err
			return
		}

		fn := luaInstance.GetGlobal(functionName)

		luaArgs := []lua.LValue{fn, lualib.NewValue(luaInstance, 2), istream, ostream}
		for _, a := range functionArgs {
			luaArgs = append(luaArgs, lualib.NewValue(luaInstance, unwrapValue(a)))
		}

		if err := luaInstance.CallByParam(lua.P{
			Fn:      luaInstance.GetGlobal("apply_stream"),
			NRet:    1,
			Protect: true,
		},
			luaArgs...,
		); err != nil {
			recSet.Errors <- err
			return
		}

		luaInstance.Get(-1) // returned value
		luaInstance.Pop(1)  // remove received value
	}()

	return recSet, nil
}
