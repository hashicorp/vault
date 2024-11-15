// Copyright 2014-2019 Aerospike, Inc.
//
// Portions may be licensed to Aerospike, Inc. under one or more contributor
// license agreements WHICH ARE COMPATIBLE WITH THE APACHE LICENSE, VERSION 2.0.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package aerospike

import (
	"fmt"

	ParticleType "github.com/aerospike/aerospike-client-go/v5/internal/particle_type"
)

// HyperLogLog (HLL) operations.
// Requires server versions >= 4.9.
//
// HyperLogLog operations on HLL items nested in lists/maps are not currently
// supported by the server.
const (
	_HLL_INIT            = 0
	_HLL_ADD             = 1
	_HLL_SET_UNION       = 2
	_HLL_SET_COUNT       = 3
	_HLL_FOLD            = 4
	_HLL_COUNT           = 50
	_HLL_UNION           = 51
	_HLL_UNION_COUNT     = 52
	_HLL_INTERSECT_COUNT = 53
	_HLL_SIMILARITY      = 54
	_HLL_DESCRIBE        = 55
)

// HLLInitOp creates HLL init operation with minhash bits.
// Server creates a new HLL or resets an existing HLL.
// Server does not return a value.
//
// policy			write policy, use DefaultHLLPolicy for default
// binName			name of bin
// indexBitCount	number of index bits. Must be between 4 and 16 inclusive. Pass -1 for default.
// minHashBitCount  number of min hash bits. Must be between 4 and 58 inclusive. Pass -1 for default.
func HLLInitOp(policy *HLLPolicy, binName string, indexBitCount, minHashBitCount int) *Operation {
	return &Operation{
		opType:   _HLL_MODIFY,
		binName:  binName,
		binValue: ListValue{_HLL_INIT, IntegerValue(indexBitCount), IntegerValue(minHashBitCount), IntegerValue(policy.flags)},
		encoder:  newHLLEncoder,
	}
}

// HLLAddOp creates HLL add operation with minhash bits.
// Server adds values to HLL set. If HLL bin does not exist, use indexBitCount and minHashBitCount
// to create HLL bin. Server returns number of entries that caused HLL to update a register.
//
// policy			write policy, use DefaultHLLPolicy for default
// binName			name of bin
// list				list of values to be added
// indexBitCount	number of index bits. Must be between 4 and 16 inclusive. Pass -1 for default.
// minHashBitCount  number of min hash bits. Must be between 4 and 58 inclusive. Pass -1 for default.
func HLLAddOp(policy *HLLPolicy, binName string, list []Value, indexBitCount, minHashBitCount int) *Operation {
	return &Operation{
		opType:   _HLL_MODIFY,
		binName:  binName,
		binValue: ListValue{_HLL_ADD, ValueArray(list), IntegerValue(indexBitCount), IntegerValue(minHashBitCount), IntegerValue(policy.flags)},
		encoder:  newHLLEncoder,
	}
}

// HLLSetUnionOp creates HLL set union operation.
// Server sets union of specified HLL objects with HLL bin.
// Server does not return a value.
//
// policy			write policy, use DefaultHLLPolicy for default
// binName			name of bin
// list				list of HLL objects
func HLLSetUnionOp(policy *HLLPolicy, binName string, list []HLLValue) *Operation {
	return &Operation{
		opType:   _HLL_MODIFY,
		binName:  binName,
		binValue: ListValue{_HLL_SET_UNION, _HLLValueArray(list), IntegerValue(policy.flags)},
		encoder:  newHLLEncoder,
	}
}

// HLLRefreshCountOp creates HLL refresh operation.
// Server updates the cached count (if stale) and returns the count.
//
// binName			name of bin
func HLLRefreshCountOp(binName string) *Operation {
	return &Operation{
		opType:   _HLL_MODIFY,
		binName:  binName,
		binValue: ListValue{_HLL_SET_COUNT},
		encoder:  newHLLEncoder,
	}
}

// HLLFoldOp creates HLL fold operation.
// Servers folds indexBitCount to the specified value.
// This can only be applied when minHashBitCount on the HLL bin is 0.
// Server does not return a value.
//
// binName			name of bin
// indexBitCount		number of index bits. Must be between 4 and 16 inclusive.
func HLLFoldOp(binName string, indexBitCount int) *Operation {
	return &Operation{
		opType:   _HLL_MODIFY,
		binName:  binName,
		binValue: ListValue{_HLL_FOLD, IntegerValue(indexBitCount)},
		encoder:  newHLLEncoder,
	}

}

// HLLGetCountOp creates HLL getCount operation.
// Server returns estimated number of elements in the HLL bin.
//
// binName			name of bin
func HLLGetCountOp(binName string) *Operation {
	return &Operation{
		opType:   _HLL_READ,
		binName:  binName,
		binValue: ListValue{_HLL_COUNT},
		encoder:  newHLLEncoder,
	}

}

// HLLGetUnionOp creates HLL getUnion operation.
// Server returns an HLL object that is the union of all specified HLL objects in the list
// with the HLL bin.
//
// binName			name of bin
// list				list of HLL objects
func HLLGetUnionOp(binName string, list []HLLValue) *Operation {
	return &Operation{
		opType:   _HLL_READ,
		binName:  binName,
		binValue: ListValue{_HLL_UNION, _HLLValueArray(list)},
		encoder:  newHLLEncoder,
	}

}

// HLLGetUnionCountOp creates HLL getUnionCount operation.
// Server returns estimated number of elements that would be contained by the union of these
// HLL objects.
//
// binName			name of bin
// list				list of HLL objects
func HLLGetUnionCountOp(binName string, list []HLLValue) *Operation {
	return &Operation{
		opType:   _HLL_READ,
		binName:  binName,
		binValue: ListValue{_HLL_UNION_COUNT, _HLLValueArray(list)},
		encoder:  newHLLEncoder,
	}
}

// HLLGetIntersectCountOp creates HLL getIntersectCount operation.
// Server returns estimated number of elements that would be contained by the intersection of
// these HLL objects.
//
// binName			name of bin
// list				list of HLL objects
func HLLGetIntersectCountOp(binName string, list []HLLValue) *Operation {
	return &Operation{
		opType:   _HLL_READ,
		binName:  binName,
		binValue: ListValue{_HLL_INTERSECT_COUNT, _HLLValueArray(list)},
		encoder:  newHLLEncoder,
	}
}

// HLLGetSimilarityOp creates HLL getSimilarity operation.
// Server returns estimated similarity of these HLL objects. Return type is a double.
//
// binName			name of bin
// list				list of HLL objects
func HLLGetSimilarityOp(binName string, list []HLLValue) *Operation {
	return &Operation{
		opType:   _HLL_READ,
		binName:  binName,
		binValue: ListValue{_HLL_SIMILARITY, _HLLValueArray(list)},
		encoder:  newHLLEncoder,
	}
}

// HLLDescribeOp creates HLL describe operation.
// Server returns indexBitCount and minHashBitCount used to create HLL bin in a list of longs.
// The list size is 2.
//
// binName			name of bin
func HLLDescribeOp(binName string) *Operation {
	return &Operation{
		opType:   _HLL_READ,
		binName:  binName,
		binValue: ListValue{_HLL_DESCRIBE},
		encoder:  newHLLEncoder,
	}
}

// hllEncoder is used to encode the HLL operations to wire protocol
func newHLLEncoder(op *Operation, packer BufferEx) (int, Error) {
	params := op.binValue.(ListValue)
	opType := params[0].(int)
	if len(op.binValue.(ListValue)) > 1 {
		return packHLLIfcParamsAsArray(packer, int16(opType), params[1:])
	}
	return packHLLIfcParamsAsArray(packer, int16(opType), nil)
}

func packHLLIfcParamsAsArray(packer BufferEx, opType int16, params ListValue) (int, Error) {
	return packHLLIfcVarParamsAsArray(packer, opType, []interface{}(params)...)
}

func packHLLIfcVarParamsAsArray(packer BufferEx, opType int16, params ...interface{}) (int, Error) {
	size := 0
	n, err := packArrayBegin(packer, len(params)+1)
	if err != nil {
		return size + n, err
	}
	size += n

	if n, err = packAInt(packer, int(opType)); err != nil {
		return size + n, err
	}
	size += n

	if len(params) > 0 {
		for i := range params {
			if n, err = packObject(packer, params[i], false); err != nil {
				return size + n, err
			}
			size += n
		}
	}

	return size, nil
}

///////////////////////////////////////////////////////////////////////////////////////

// _HLLValueArray encapsulates an array of Value.
// Supported by Aerospike 3+ servers only.
type _HLLValueArray []HLLValue

func (va _HLLValueArray) EstimateSize() (int, Error) {
	return packHLLValueArray(nil, va)
}

func (va _HLLValueArray) write(cmd BufferEx) (int, Error) {
	return packHLLValueArray(cmd, va)
}

func (va _HLLValueArray) pack(cmd BufferEx) (int, Error) {
	return packHLLValueArray(cmd, va)
}

// GetType returns wire protocol value type.
func (va _HLLValueArray) GetType() int {
	return ParticleType.LIST
}

// GetObject returns original value as an interface{}.
func (va _HLLValueArray) GetObject() interface{} {
	return []HLLValue(va)
}

// String implements Stringer interface.
func (va _HLLValueArray) String() string {
	return fmt.Sprintf("%v", []HLLValue(va))
}

func packHLLValueArray(cmd BufferEx, list _HLLValueArray) (int, Error) {
	size := 0
	n, err := packArrayBegin(cmd, len(list))
	if err != nil {
		return n, err
	}
	size += n

	for i := range list {
		n, err = list[i].pack(cmd)
		if err != nil {
			return 0, err
		}
		size += n
	}

	return size, err
}
