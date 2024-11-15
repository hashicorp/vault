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

package aerospike

const expMapMODULE int64 = 0

// ExpMapPut creates an expression that writes key/value item to map bin.
func ExpMapPut(
	policy *MapPolicy,
	key *Expression,
	value *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	var args []ExpressionArgument
	op := mapWriteOp(policy, false)
	if op == cdtMapOpTypeReplace {
		args = []ExpressionArgument{
			cdtContextList(ctx),
			IntegerValue(op),
			key,
			value,
		}
	} else {
		args = []ExpressionArgument{
			cdtContextList(ctx),
			IntegerValue(op),
			key,
			value,
			IntegerValue(policy.attributes.attr),
		}
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapPutItems creates an expression that writes each map item to map bin.
func ExpMapPutItems(
	policy *MapPolicy,
	amap *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	var args []ExpressionArgument
	var op = mapWriteOp(policy, true)
	if op == cdtMapOpTypeReplace {
		args = []ExpressionArgument{
			cdtContextList(ctx),
			IntegerValue(op),
			amap,
		}
	} else {
		args = []ExpressionArgument{
			cdtContextList(ctx),
			IntegerValue(op),
			amap,
			IntegerValue(policy.attributes.attr),
		}
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapIncrement creates an expression that increments values by incr for all items identified by key.
// Valid only for numbers.
func ExpMapIncrement(
	policy *MapPolicy,
	key *Expression,
	incr *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeIncrement),
		key,
		incr,
		cdtContextList(ctx),
		IntegerValue(policy.attributes.attr),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapClear creates an expression that removes all items in map.
func ExpMapClear(bin *Expression, ctx ...*CDTContext) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeClear),
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByKey creates an expression that removes map item identified by key.
func ExpMapRemoveByKey(
	key *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByKey),
		IntegerValue(MapReturnType.NONE),
		key,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByKeyList creates an expression that removes map items identified by keys.
func ExpMapRemoveByKeyList(
	keys *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveKeyList),
		IntegerValue(MapReturnType.NONE),
		keys,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByKeyRange creates an expression that removes map items identified by key range (keyBegin inclusive, keyEnd exclusive).
// If keyBegin is null, the range is less than keyEnd.
// If keyEnd is null, the range is greater than equal to keyBegin.
func ExpMapRemoveByKeyRange(
	keyBegin *Expression,
	keyEnd *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	var args = []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(cdtMapOpTypeRemoveByKeyInterval),
		IntegerValue(MapReturnType.NONE),
	}
	if keyBegin != nil {
		args = append(args, keyBegin)
	} else {
		args = append(args, nullValue)
	}
	if keyEnd != nil {
		args = append(args, keyEnd)
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByKeyRelativeIndexRange creates an expression that removes map items nearest to key and greater by index.
//
// Examples for map [{0=17},{4=2},{5=15},{9=10}]:
//
// * (value,index) = [removed items]
// * (5,0) = [{5=15},{9=10}]
// * (5,1) = [{9=10}]
// * (5,-1) = [{4=2},{5=15},{9=10}]
// * (3,2) = [{9=10}]
// * (3,-2) = [{0=17},{4=2},{5=15},{9=10}]
func ExpMapRemoveByKeyRelativeIndexRange(
	key *Expression,
	index *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByKeyRelIndexRange),
		IntegerValue(MapReturnType.NONE),
		key,
		index,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByKeyRelativeIndexRangeCount creates an expression that removes map items nearest to key and greater by index with a count limit.
//
// Examples for map [{0=17},{4=2},{5=15},{9=10}]:
//
// (value,index,count) = [removed items]
// * (5,0,1) = [{5=15}]
// * (5,1,2) = [{9=10}]
// * (5,-1,1) = [{4=2}]
// * (3,2,1) = [{9=10}]
// * (3,-2,2) = [{0=17}]
func ExpMapRemoveByKeyRelativeIndexRangeCount(
	key *Expression,
	index *Expression,
	count *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByKeyRelIndexRange),
		IntegerValue(MapReturnType.NONE),
		key,
		index,
		count,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByValue creates an expression that removes map items identified by value.
func ExpMapRemoveByValue(
	value *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByValue),
		IntegerValue(MapReturnType.NONE),
		value,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByValueList creates an expression that removes map items identified by values.
func ExpMapRemoveByValueList(
	values *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveValueList),
		IntegerValue(MapReturnType.NONE),
		values,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByValueRange creates an expression that removes map items identified by value range (valueBegin inclusive, valueEnd exclusive).
// If valueBegin is null, the range is less than valueEnd.
// If valueEnd is null, the range is greater than equal to valueBegin.
func ExpMapRemoveByValueRange(
	valueBegin *Expression,
	valueEnd *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(cdtMapOpTypeRemoveByValueInterval),
		IntegerValue(MapReturnType.NONE),
	}
	if valueBegin != nil {
		args = append(args, valueBegin)
	} else {
		args = append(args, nullValue)
	}
	if valueEnd != nil {
		args = append(args, valueEnd)
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByValueRelativeRankRange creates an expression that removes map items nearest to value and greater by relative rank.
//
// Examples for map [{4=2},{9=10},{5=15},{0=17}]:
//
// * (value,rank) = [removed items]
// * (11,1) = [{0=17}]
// * (11,-1) = [{9=10},{5=15},{0=17}]
func ExpMapRemoveByValueRelativeRankRange(
	value *Expression,
	rank *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByValueRelRankRange),
		IntegerValue(MapReturnType.NONE),
		value,
		rank,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByValueRelativeRankRangeCount creates an expression that removes map items nearest to value and greater by relative rank with a count limit.
//
// Examples for map [{4=2},{9=10},{5=15},{0=17}]:
//
// * (value,rank,count) = [removed items]
// * (11,1,1) = [{0=17}]
// * (11,-1,1) = [{9=10}]
func ExpMapRemoveByValueRelativeRankRangeCount(
	value *Expression,
	rank *Expression,
	count *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByValueRelRankRange),
		IntegerValue(MapReturnType.NONE),
		value,
		rank,
		count,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByIndex creates an expression that removes map item identified by index.
func ExpMapRemoveByIndex(
	index *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByIndex),
		IntegerValue(MapReturnType.NONE),
		index,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByIndexRange creates an expression that removes map items starting at specified index to the end of map.
func ExpMapRemoveByIndexRange(
	index *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByIndexRange),
		IntegerValue(MapReturnType.NONE),
		index,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByIndexRangeCount creates an expression that removes "count" map items starting at specified index.
func ExpMapRemoveByIndexRangeCount(
	index *Expression,
	count *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByIndexRange),
		IntegerValue(MapReturnType.NONE),
		index,
		count,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByRank creates an expression that removes map item identified by rank.
func ExpMapRemoveByRank(
	rank *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByRank),
		IntegerValue(MapReturnType.NONE),
		rank,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByRankRange creates an expression that removes map items starting at specified rank to the last ranked item.
func ExpMapRemoveByRankRange(
	rank *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByRankRange),
		IntegerValue(MapReturnType.NONE),
		rank,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByRankRangeCount creates an expression that removes "count" map items starting at specified rank.
func ExpMapRemoveByRankRangeCount(
	rank *Expression,
	count *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeRemoveByRankRange),
		IntegerValue(MapReturnType.NONE),
		rank,
		count,
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapSize creates an expression that returns list size.
func ExpMapSize(bin *Expression, ctx ...*CDTContext) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeSize),
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, ExpTypeINT, args)
}

// ExpMapGetByKey creates an expression that selects map item identified by key and returns selected data
// specified by returnType.
func ExpMapGetByKey(
	returnType mapReturnType,
	valueType ExpType,
	key *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByKey),
		IntegerValue(returnType),
		key,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, valueType, args)
}

// ExpMapGetByKeyRange creates an expression that selects map items identified by key range (keyBegin inclusive, keyEnd exclusive).
// If keyBegin is null, the range is less than keyEnd.
// If keyEnd is null, the range is greater than equal to keyBegin.
// Expression returns selected data specified by returnType.
func ExpMapGetByKeyRange(
	returnType mapReturnType,
	keyBegin *Expression,
	keyEnd *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(cdtMapOpTypeGetByKeyInterval),
		IntegerValue(returnType),
	}
	if keyBegin != nil {
		args = append(args, keyBegin)
	} else {
		args = append(args, nullValue)
	}
	if keyEnd != nil {
		args = append(args, keyEnd)
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByKeyList creates an expression that selects map items identified by keys and returns selected data specified by returnType
func ExpMapGetByKeyList(
	returnType mapReturnType,
	keys *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByKeyList),
		IntegerValue(returnType),
		keys,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByKeyRelativeIndexRange creates an expression that selects map items nearest to key and greater by index.
// Expression returns selected data specified by returnType.
//
// Examples for ordered map [{0=17},{4=2},{5=15},{9=10}]:
//
// * (value,index) = [selected items]
// * (5,0) = [{5=15},{9=10}]
// * (5,1) = [{9=10}]
// * (5,-1) = [{4=2},{5=15},{9=10}]
// * (3,2) = [{9=10}]
// * (3,-2) = [{0=17},{4=2},{5=15},{9=10}]
func ExpMapGetByKeyRelativeIndexRange(
	returnType mapReturnType,
	key *Expression,
	index *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByKeyRelIndexRange),
		IntegerValue(returnType),
		key,
		index,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByKeyRelativeIndexRangeCount creates an expression that selects map items nearest to key and greater by index with a count limit.
// Expression returns selected data specified by returnType.
//
// Examples for ordered map [{0=17},{4=2},{5=15},{9=10}]:
//
// * (value,index,count) = [selected items]
// * (5,0,1) = [{5=15}]
// * (5,1,2) = [{9=10}]
// * (5,-1,1) = [{4=2}]
// * (3,2,1) = [{9=10}]
// * (3,-2,2) = [{0=17}]
func ExpMapGetByKeyRelativeIndexRangeCount(
	returnType mapReturnType,
	key *Expression,
	index *Expression,
	count *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByKeyRelIndexRange),
		IntegerValue(returnType),
		key,
		index,
		count,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByValue creates an expression that selects map items identified by value and returns selected data
// specified by returnType.
func ExpMapGetByValue(
	returnType mapReturnType,
	value *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByValue),
		IntegerValue(returnType),
		value,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByValueRange creates an expression that selects map items identified by value range (valueBegin inclusive, valueEnd exclusive)
// If valueBegin is null, the range is less than valueEnd.
// If valueEnd is null, the range is greater than equal to valueBegin.
//
// Expression returns selected data specified by returnType.
func ExpMapGetByValueRange(
	returnType mapReturnType,
	valueBegin *Expression,
	valueEnd *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(cdtMapOpTypeGetByValueInterval),
		IntegerValue(returnType),
	}
	if valueBegin != nil {
		args = append(args, valueBegin)
	} else {
		args = append(args, nullValue)
	}
	if valueEnd != nil {
		args = append(args, valueEnd)
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByValueList creates an expression that selects map items identified by values and returns selected data specified by returnType.
func ExpMapGetByValueList(
	returnType mapReturnType,
	values *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByValueList),
		IntegerValue(returnType),
		values,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByValueRelativeRankRange creates an expression that selects map items nearest to value and greater by relative rank.
// Expression returns selected data specified by returnType.
//
// Examples for map [{4=2},{9=10},{5=15},{0=17}]:
//
// * (value,rank) = [selected items]
// * (11,1) = [{0=17}]
// * (11,-1) = [{9=10},{5=15},{0=17}]
func ExpMapGetByValueRelativeRankRange(
	returnType mapReturnType,
	value *Expression,
	rank *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByValueRelRankRange),
		IntegerValue(returnType),
		value,
		rank,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByValueRelativeRankRangeCount creates an expression that selects map items nearest to value and greater by relative rank with a count limit.
// Expression returns selected data specified by returnType.
//
// Examples for map [{4=2},{9=10},{5=15},{0=17}]:
//
// * (value,rank,count) = [selected items]
// * (11,1,1) = [{0=17}]
// * (11,-1,1) = [{9=10}]
func ExpMapGetByValueRelativeRankRangeCount(
	returnType mapReturnType,
	value *Expression,
	rank *Expression,
	count *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByValueRelRankRange),
		IntegerValue(returnType),
		value,
		rank,
		count,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByIndex creates an expression that selects map item identified by index and returns selected data specified by returnType.
func ExpMapGetByIndex(
	returnType mapReturnType,
	valueType ExpType,
	index *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByIndex),
		IntegerValue(returnType),
		index,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, valueType, args)
}

// ExpMapGetByIndexRange creates an expression that selects map items starting at specified index to the end of map and returns selected
// data specified by returnType.
func ExpMapGetByIndexRange(
	returnType mapReturnType,
	index *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByIndexRange),
		IntegerValue(returnType),
		index,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByIndexRangeCount creates an expression that selects "count" map items starting at specified index and returns selected data
// specified by returnType.
func ExpMapGetByIndexRangeCount(
	returnType mapReturnType,
	index *Expression,
	count *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByIndexRange),
		IntegerValue(returnType),
		index,
		count,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByRank creates an expression that selects map item identified by rank and returns selected data specified by returnType.
func ExpMapGetByRank(
	returnType mapReturnType,
	valueType ExpType,
	rank *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByRank),
		IntegerValue(returnType),
		rank,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, valueType, args)
}

// ExpMapGetByRankRange creates an expression that selects map items starting at specified rank to the last ranked item and
// returns selected data specified by returnType.
func ExpMapGetByRankRange(
	returnType mapReturnType,
	rank *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByRankRange),
		IntegerValue(returnType),
		rank,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

// ExpMapGetByRankRangeCount creates an expression that selects "count" map items starting at specified rank and returns selected
// data specified by returnType.
func ExpMapGetByRankRangeCount(
	returnType mapReturnType,
	rank *Expression,
	count *Expression,
	bin *Expression,
	ctx ...*CDTContext,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByRankRange),
		IntegerValue(returnType),
		rank,
		count,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(returnType), args)
}

func expMapAddRead(
	bin *Expression,
	returnType ExpType,
	arguments []ExpressionArgument,
) *Expression {
	flags := expMapMODULE
	return &Expression{
		cmd:       &expOpCALL,
		val:       nil,
		bin:       bin,
		flags:     &flags,
		module:    &returnType,
		exps:      nil,
		arguments: arguments,
	}
}

func expMapAddWrite(
	bin *Expression,
	arguments []ExpressionArgument,
	ctx ...*CDTContext,
) *Expression {
	var returnType ExpType
	if len(ctx) == 0 {
		returnType = ExpTypeMAP
	} else if (ctx[0].id & ctxTypeListIndex) == 0 {
		returnType = ExpTypeMAP
	} else {
		returnType = ExpTypeLIST
	}

	flags := expMapMODULE | _MODIFY
	return &Expression{
		cmd:       &expOpCALL,
		val:       nil,
		bin:       bin,
		flags:     &flags,
		module:    &returnType,
		exps:      nil,
		arguments: arguments,
	}
}

func expMapGetValueType(returnType mapReturnType) ExpType {
	var t = returnType & (^MapReturnType.INVERTED)
	if t == MapReturnType.KEY || t == MapReturnType.VALUE {
		return ExpTypeLIST
	} else if t == MapReturnType.KEY_VALUE {
		return ExpTypeMAP
	}
	return ExpTypeINT
}

// Determines the correct operation to use when setting one or more map values, depending on the
// map policy.
func mapWriteOp(policy *MapPolicy, multi bool) int {
	switch policy.flags {
	default:
		fallthrough
	case MapWriteFlagsDefault:
		if multi {
			return cdtMapOpTypePutItems
		}
		return cdtMapOpTypePut
	case MapWriteFlagsUpdateOnly:
		if multi {
			return cdtMapOpTypeReplaceItems
		}
		return cdtMapOpTypeReplace
	case MapWriteFlagsCreateOnly:
		if multi {
			return cdtMapOpTypeAddItems
		}
		return cdtMapOpTypeAdd
	}
}
