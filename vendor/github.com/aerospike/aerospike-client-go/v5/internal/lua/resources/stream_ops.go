// +build !app_engine

package luaLib

// LibStreamOps is the source code for the stream library in the lua instance
const LibStreamOps = `
-- Lua Interface for Aerospike Record Stream Support
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

local function check_limit(v)
    return type(v) == 'number' and v >= 1000
end

--
-- clone a table. creates a shallow copy of the table.
--
local function clone_table(t)
    local out = {}
    for k,v in pairs(t) do
        out[k] = v
    end
    return out
end

--
-- Clone a value.
--
local function clone(v)

    local t = type(v)

    if t == 'number' then
        return v
    elseif t == 'string' then
        return v
    elseif t == 'boolean' then
        return v
    elseif t == 'table' then
        return clone_table(v)
    elseif t == 'userdata' then
        if v.__index == Map then
            return map.clone(v)
        elseif v.__index == List then
            return list.clone(v)
        end
        return nil
    end

    return v
end

--
-- Filter values
-- @param next - a generator that produces the next value from a stream
-- @param f - the function to transform each value
--
function filter( next, p )
    -- done indicates if we exhausted the 'next' stream
    local done = false

    -- return a closure which the caller can use to get the next value
    return function()
        
        -- we bail if we already exhausted the stream
        if done then return nil end

        -- find the first value which satisfies the predicate
        for a in next do
            if p(a) then
                return a
            end
        end

        done = true

        return nil
    end
end

--
-- Transform values
-- @param next - a generator that produces the next value from a stream
-- @param f - the tranfomation operation
--
function transform( next, f )
    -- done indicates if we exhausted the 'next' stream
    local done = false

    -- return a closure which the caller can use to get the next value
    return function()
        
        -- we bail if we already exhausted the stream
        if done then return nil end
        
        -- get the first value
        local a = next()

        -- apply the transformation
        if a ~= nil then
            return f(a)
        end

        done = true;

        return nil
    end
end

--
-- Combines two values from an istream into a single value.
-- @param next - a generator that produces the next value from a stream
-- @param f - the reduction operation
--
function reduce( next, f )
    -- done indicates if we exhausted the 'next' stream
    local done = false

    -- return a closure which the caller can use to get the next value
    return function()


        -- we bail if we already exhausted the stream
        if done then return nil end
        
        -- get the first value
        local a = next()

        if a ~= nil then
            -- get each subsequent value and reduce them
            for b in next do
                a = f(a,b)
            end
        end

        -- we are done!
        done = true

        return a
    end
end

--
-- Aggregate values into a single value.
-- @param next - a generator that produces the next value from a stream
-- @param f - the aggregation operation
--
function aggregate( next, init, f )
    -- done indicates if we exhausted the 'next' stream
    local done = false

    -- return a closure which the caller can use to get the next value
    return function()

        -- we bail if we already exhausted the stream
        if done then return nil end

        -- get the initial value
        local a = clone(init)
        
        -- get each subsequent value and aggregate them
        for b in next do
            a = f(a,b)

            -- check the size limit, if it is exceeded,
            -- then return the value
            if check_limit(a) then
                return a
            end
        end

        -- we are done!
        done = true

        return a
    end
end

--
-- as_stream iterator
--
function stream_iterator(s)
    local done = false
    return function()
        if done then return nil end
        local v = stream.read(s)
        if v == nil then
            done = true
        end
        return v;
    end
end



-- ######################################################################################
--
-- StreamOps
-- Builds a sequence of operations to be applied to a stream of values.
--
-- ######################################################################################

StreamOps = {}
StreamOps_mt = { __index = StreamOps }

-- Op only executes on server
local SCOPE_SERVER = 1

-- Op only executes on client
local SCOPE_CLIENT = 2

-- Op can execute on either client or server
local SCOPE_EITHER = 3

-- Op executes on both client and server
local SCOPE_BOTH = 4

--
-- Creates a new StreamOps using an array of ops
-- 
-- @param ops an array of operations
--
function StreamOps_create()
    local self = {}
    setmetatable(self, StreamOps_mt)
    self.ops = {}
    return self
end

function StreamOps_apply(stream, ops, i, n)

    -- if nil, then use default values
    i = i or 1
    n = n or #ops
    
    -- if index in list > size of list, then return the stream
    if i > n then return stream end
    
    -- get the current operation
    local op = ops[i]

    -- apply the operation and get a stream or use provided stream
    local s = op.func(stream, unpack(op.args)) or stream

    -- move to the next operation
    return StreamOps_apply(s, ops, i + 1, n)
end


--
-- This selects the operations appropriate for a given scope.
-- For the SERVER scope, it will select the first n ops until one of the ops
-- is a CLIENT scope op.
-- For the CLIENT scope, it will skip the first n ops that are SERVER scope 
-- ops, then it will take the remaining ops, including SERVER scoped ops.
--
function StreamOps_select(stream_ops, scope)
    local server_ops = {}
    local client_ops = {}
    
    local phase = SCOPE_SERVER
    for i,op in ipairs(stream_ops) do
        if phase == SCOPE_SERVER then
            if op.scope == SCOPE_SERVER then
                table.insert(server_ops, op)
            elseif op.scope == SCOPE_EITHER then
                table.insert(server_ops, op)
            elseif op.scope == SCOPE_BOTH then
                table.insert(server_ops, op)
                table.insert(client_ops, op)
                phase = SCOPE_CLIENT
            end
        elseif phase == SCOPE_CLIENT then
            table.insert(client_ops, op)
        end 
    end
    
    if scope == SCOPE_CLIENT then
        return client_ops
    else
        return server_ops
    end
end



-- 
-- OPS: [ OP, ... ]
-- OP: {scope=SCOPE, name=NAME, func=FUNC, args=ARGS}
-- SCOPE: ANY(0) | SERVER(1) | CLIENT(2) | 
-- NAME: FUNCTION NAME
-- FUNC: FUNCTION POINTER
-- ARGS: ARRAY OF ARGUMENTS
--


function StreamOps:aggregate(...)
    table.insert(self.ops, { scope = SCOPE_SERVER, name = "aggregate", func = aggregate, args = {...}})
    return self
end

function StreamOps:reduce(...)
    table.insert(self.ops, { scope = SCOPE_BOTH, name = "reduce", func = reduce, args = {...}})
    return self
end

function StreamOps:map(...)
    table.insert(self.ops, { scope = SCOPE_EITHER, name = "map", func = transform, args = {...}})
    return self
end

function StreamOps:filter(...)
    table.insert(self.ops, { scope = SCOPE_EITHER, name = "filter", func = filter, args = {...}})
    return self
end

-- stream : group(f)
--
-- Group By will return a Map of keys to a list of values. The key is determined by applying the 
-- function 'f' to each element in the stream.
--
function StreamOps:groupby(f)

    local function _aggregate(m, v)
        local k = f and f(v) or nil;
        local l = m[k] or list()
        list.append(l, v)
        m[k] = l;
        return m;
    end

    local function _merge(l1, l2)
        local l = list.clone(l1)
        for v in list.iterator(l2) do
            list.append(l, v)
        end
        return l
    end

    function _reduce(m1, m2)
        return map.merge(m1, m2, _merge)
    end

    return self : aggregate(map(), _aggregate) : reduce(_reduce)
end
`
