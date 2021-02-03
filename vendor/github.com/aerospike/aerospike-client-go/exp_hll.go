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

const hllMODULE int64 = 2

var (
	HllExpOpADD            = 1
	HllExpOpCOUNT          = 50
	HllExpOpUNION          = 51
	HllExpOpUNIONCOUNT     = 52
	HllExpOpINTERSECTCOUNT = 53
	HllExpOpSIMILARITY     = 54
	HllExpOpDESCRIBE       = 55
	HllExpOpMAYCONTAIN     = 56
)

// ExpHLLAdd creates an expression that adds list values to a HLL set and returns HLL set.
// The function assumes HLL bin already exists.
func ExpHLLAdd(policy *HLLPolicy, list *FilterExpression, bin *FilterExpression) *FilterExpression {
	return ExpHLLAddWithIndexAndMinHash(policy, list, ExpIntVal(-1), ExpIntVal(-1), bin)
}

// ExpHLLAddWithIndex creates an expression that adds values to a HLL set and returns HLL set.
// If HLL bin does not exist, use `indexBitCount` to create HLL bin.
func ExpHLLAddWithIndex(
	policy *HLLPolicy,
	list *FilterExpression,
	index_bit_count *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	return ExpHLLAddWithIndexAndMinHash(policy, list, index_bit_count, ExpIntVal(-1), bin)
}

// ExpHLLAddWithIndexAndMinHash creates an expression that adds values to a HLL set and returns HLL set. If HLL bin does not
// exist, use `indexBitCount` and `minHashBitCount` to create HLL set.
func ExpHLLAddWithIndexAndMinHash(
	policy *HLLPolicy,
	list *FilterExpression,
	index_bit_count *FilterExpression,
	min_hash_count *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	return expHLLAddWrite(
		bin,
		[]ExpressionArgument{
			IntegerValue(HllExpOpADD),
			list,
			index_bit_count,
			min_hash_count,
			IntegerValue(policy.flags),
		},
	)
}

// ExpHLLGetCount creates an expression that returns estimated number of elements in the HLL bin.
func ExpHLLGetCount(bin *FilterExpression) *FilterExpression {
	return expHLLAddRead(
		bin,
		ExpTypeINT,
		[]ExpressionArgument{
			IntegerValue(HllExpOpCOUNT),
		},
	)
}

// ExpHLLGetUnion creates an expression that returns a HLL object that is the union of all specified HLL objects
// in the list with the HLL bin.
func ExpHLLGetUnion(list *FilterExpression, bin *FilterExpression) *FilterExpression {
	return expHLLAddRead(
		bin,
		ExpTypeHLL,
		[]ExpressionArgument{
			IntegerValue(HllExpOpUNION),
			list,
		},
	)
}

// ExpHLLGetUnionCount creates an expression that returns estimated number of elements that would be contained by
// the union of these HLL objects.
func ExpHLLGetUnionCount(list *FilterExpression, bin *FilterExpression) *FilterExpression {
	return expHLLAddRead(
		bin,
		ExpTypeINT,
		[]ExpressionArgument{
			IntegerValue(HllExpOpUNIONCOUNT),
			list,
		},
	)
}

// ExpHLLGetIntersectCount creates an expression that returns estimated number of elements that would be contained by
// the intersection of these HLL objects.
func ExpHLLGetIntersectCount(list *FilterExpression, bin *FilterExpression) *FilterExpression {
	return expHLLAddRead(
		bin,
		ExpTypeINT,
		[]ExpressionArgument{
			IntegerValue(HllExpOpINTERSECTCOUNT),
			list,
		},
	)
}

// ExpHLLGetSimilarity creates an expression that returns estimated similarity of these HLL objects as a 64 bit float.
func ExpHLLGetSimilarity(list *FilterExpression, bin *FilterExpression) *FilterExpression {
	return expHLLAddRead(
		bin,
		ExpTypeFLOAT,
		[]ExpressionArgument{
			IntegerValue(HllExpOpSIMILARITY),
			list,
		},
	)
}

// ExpHLLDescribe creates an expression that returns `indexBitCount` and `minHashBitCount` used to create HLL bin
// in a list of longs. `list[0]` is `indexBitCount` and `list[1]` is `minHashBitCount`.
func ExpHLLDescribe(bin *FilterExpression) *FilterExpression {
	return expHLLAddRead(
		bin,
		ExpTypeLIST,
		[]ExpressionArgument{
			IntegerValue(HllExpOpDESCRIBE),
		},
	)
}

// ExpHLLMayContain creates an expression that returns one if HLL bin may contain all items in the list.
func ExpHLLMayContain(list *FilterExpression, bin *FilterExpression) *FilterExpression {
	return expHLLAddRead(
		bin,
		ExpTypeINT,
		[]ExpressionArgument{
			IntegerValue(HllExpOpMAYCONTAIN),
			list,
		},
	)
}

func expHLLAddRead(
	bin *FilterExpression,
	return_type ExpType,
	arguments []ExpressionArgument,
) *FilterExpression {
	flags := hllMODULE
	return &FilterExpression{
		cmd:       &expOpCALL,
		val:       nil,
		bin:       bin,
		flags:     &flags,
		module:    &return_type,
		exps:      nil,
		arguments: arguments,
	}
}

func expHLLAddWrite(bin *FilterExpression, arguments []ExpressionArgument) *FilterExpression {
	flags := hllMODULE | MODIFY
	return &FilterExpression{
		cmd:       &expOpCALL,
		val:       nil,
		bin:       bin,
		flags:     &flags,
		module:    &ExpTypeHLL,
		exps:      nil,
		arguments: arguments,
	}
}
