// +build !app_engine

package luaLib

// LibAerospike is the source code for the Aerospike library in the lua instance
const LibAerospike = `

-- The Lua Interface to Aerospike
--
-- ======================================================================
-- Copyright [2014] Aerospike, Inc.. Portions may be licensed
-- to Aerospike, Inc. under one or more contributor license agreements.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--  http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.
-- ======================================================================

-- A table to track whether we had sandboxed a function
sandboxed = {}

ldebug = debug;

-- ############################################################################
--
-- LOG FUNCTIONS
--
-- ############################################################################

function trace(m, ...)
    return aerospike:log(4, string.format(m, ...))
end

function debug(m, ...)
    return aerospike:log(3, string.format(m, ...))
end

function info(m, ...)
    return aerospike:log(2, string.format(m, ...))
end

function warn(m, ...)
    return aerospike:log(1, string.format(m, ...))
end

-- ############################################################################
--
-- APPLY FUNCTIONS
--
-- ############################################################################

--
-- Creates a new environment for use in apply_record functions
--
function env_record()
    return {

        -- aerospike types
        ["record"] = record,
        ["iterator"] = iterator,
        ["list"] = list,
        ["map"] = map,
        ["bytes"] = bytes,
        ["aerospike"] = aerospike,

        ["putX"] = putX,

        -- logging functions
        ["trace"] = trace,
        ["debug"] = debug,
        ["info"] = info,
        ["warn"] = warn,
        
        -- standard lua functions
        ["collectgarbage"] = collectgarbage,
        ["error"] = error,
        ["getmetatable"] = getmetatable,
        ["ipairs"] = ipairs,
        ["load"] = loadstring,
        ["module"] = module,
        ["next"] = next,
        ["pairs"] = pairs,
        ["print"] = print,
        ["pcall"] = pcall,
        ["rawequal"] = rawequal,
        ["rawget"] = rawget,
        ["rawset"] = rawset,
        ["require"] = require,
        ["require"] = require,
        ["select"] = select,
        ["setmetatable"] = setmetatable,
        ["setfenv"] = setfenv,
        ["tonumber"] = tonumber,
        ["tostring"] = tostring,
        ["type"] = type,
        ["unpack"] = unpack,
        ["xpcall"] = xpcall,

        -- standard lua objects
        ["math"] = math,
        ["io"] = io,
        ["os"] = {
            ['clock'] = os.clock,
            ['date'] = os.date,
            ['difftime'] = os.difftime,
            ['getenv'] = os.getenv,
            ['setlocale'] = os.setlocale,
            ['time'] = os.time,
            ['tmpname'] = os.tmpname
        },
        ["package"] = package,
        ["string"] = string,
        ["table"] = table,

        -- standard lua variables
        ["_G"] = {}
    }
end

--
-- Apply function to a record and arguments.
--
-- @param f the fully-qualified name of the function.
-- @param r the record to be applied to the function.
-- @param ... additional arguments to be applied to the function.
-- @return result of the called function or nil.
-- 
function apply_record(f, r, ...)

    if f == nil then
        error("function not found", 2)
    end
    
    if not sandboxed[f] then
        setfenv(f,env_record())
        sandboxed[f] = true
    end

    success, result = pcall(f, r, ...)
    if success then
        return result
    else
        error(result, 2)
        return nil
    end
end

--
-- Apply function to an iterator and arguments.
--
-- @param f the fully-qualified name of the function.
-- @param s the iterator to be applied to the function.
-- @param ... additional arguments to be applied to the function.
-- @return 0 on success, otherwise failure.
-- 
function apply_stream(f, scope, istream, ostream, ...)

    if f == nil then
        error("function not found", 2)
        return 2
    end
    
    --require("stream_ops")

    if not sandboxed[f] then
        setfenv(f,env_record())
        sandboxed[f] = true
    end

    local stream_ops = StreamOps_create();
    success, result = pcall(f, stream_ops, ...)

    -- info("apply_stream: success=%s, result=%s", tostring(success), tostring(result))

    if success then
        local ops = StreamOps_select(result.ops, scope);
        
        -- Apply server operations to the stream
        -- result => a stream_ops object
        local values = StreamOps_apply(stream_iterator(istream), ops);

        -- Iterate the stream of values from the computation
        -- then pipe it to the ostream
        for value in values do
            -- info("value = %s", tostring(value))
            stream.write(ostream, value)
        end

        -- 0 is success
        return 0
    else
        error(result, 2)
        return 2
    end
end
`
