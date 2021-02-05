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

const expMapMODULE int64 = 0

// ExpMapPut creates an expression that writes key/value item to map bin.
func ExpMapPut(
	policy *MapPolicy,
	key *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{}
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
	amap *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	var args = []ExpressionArgument{}
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
	key *FilterExpression,
	incr *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
func ExpMapClear(bin *FilterExpression, ctx ...*CDTContext) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeClear),
		cdtContextList(ctx),
	}
	return expMapAddWrite(bin, args, ctx...)
}

// ExpMapRemoveByKey creates an expression that removes map item identified by key.
func ExpMapRemoveByKey(
	key *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	keys *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	key_begin *FilterExpression,
	key_end *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	var args = []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(cdtMapOpTypeRemoveByKeyInterval),
		IntegerValue(MapReturnType.NONE),
	}
	if key_begin != nil {
		args = append(args, key_begin)
	} else {
		args = append(args, nullValue)
	}
	if key_end != nil {
		args = append(args, key_end)
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
	key *FilterExpression,
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	key *FilterExpression,
	index *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	value *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	values *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	value_begin *FilterExpression,
	value_end *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(cdtMapOpTypeRemoveByValueInterval),
		IntegerValue(MapReturnType.NONE),
	}
	if value_begin != nil {
		args = append(args, value_begin)
	} else {
		args = append(args, nullValue)
	}
	if value_end != nil {
		args = append(args, value_end)
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
	value *FilterExpression,
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	value *FilterExpression,
	rank *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	index *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
	rank *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
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
func ExpMapSize(bin *FilterExpression, ctx ...*CDTContext) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeSize),
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, ExpTypeINT, args)
}

// ExpMapGetByKey creates an expression that selects map item identified by key and returns selected data
// specified by returnType.
func ExpMapGetByKey(
	return_type mapReturnType,
	value_type ExpType,
	key *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByKey),
		IntegerValue(return_type),
		key,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, value_type, args)
}

// ExpMapGetByKeyRange creates an expression that selects map items identified by key range (keyBegin inclusive, keyEnd exclusive).
// If keyBegin is null, the range is less than keyEnd.
// If keyEnd is null, the range is greater than equal to keyBegin.
// Expression returns selected data specified by returnType.
func ExpMapGetByKeyRange(
	return_type mapReturnType,
	key_begin *FilterExpression,
	key_end *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(cdtMapOpTypeGetByKeyInterval),
		IntegerValue(return_type),
	}
	if key_begin != nil {
		args = append(args, key_begin)
	} else {
		args = append(args, nullValue)
	}
	if key_end != nil {
		args = append(args, key_end)
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
}

// ExpMapGetByKeyList creates an expression that selects map items identified by keys and returns selected data specified by returnType
func ExpMapGetByKeyList(
	return_type mapReturnType,
	keys *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByKeyList),
		IntegerValue(return_type),
		keys,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
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
	return_type mapReturnType,
	key *FilterExpression,
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByKeyRelIndexRange),
		IntegerValue(return_type),
		key,
		index,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
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
	return_type mapReturnType,
	key *FilterExpression,
	index *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByKeyRelIndexRange),
		IntegerValue(return_type),
		key,
		index,
		count,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
}

// ExpMapGetByValue creates an expression that selects map items identified by value and returns selected data
// specified by returnType.
func ExpMapGetByValue(
	return_type mapReturnType,
	value *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByValue),
		IntegerValue(return_type),
		value,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
}

// ExpMapGetByValueRange creates an expression that selects map items identified by value range (valueBegin inclusive, valueEnd exclusive)
// If valueBegin is null, the range is less than valueEnd.
// If valueEnd is null, the range is greater than equal to valueBegin.
//
// Expression returns selected data specified by returnType.
func ExpMapGetByValueRange(
	return_type mapReturnType,
	value_begin *FilterExpression,
	value_end *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(cdtMapOpTypeGetByValueInterval),
		IntegerValue(return_type),
	}
	if value_begin != nil {
		args = append(args, value_begin)
	} else {
		args = append(args, nullValue)
	}
	if value_end != nil {
		args = append(args, value_end)
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
}

// ExpMapGetByValueList creates an expression that selects map items identified by values and returns selected data specified by returnType.
func ExpMapGetByValueList(
	return_type mapReturnType,
	values *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByValueList),
		IntegerValue(return_type),
		values,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
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
	return_type mapReturnType,
	value *FilterExpression,
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByValueRelRankRange),
		IntegerValue(return_type),
		value,
		rank,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
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
	return_type mapReturnType,
	value *FilterExpression,
	rank *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByValueRelRankRange),
		IntegerValue(return_type),
		value,
		rank,
		count,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
}

// ExpMapGetByIndex creates an expression that selects map item identified by index and returns selected data specified by returnType.
func ExpMapGetByIndex(
	return_type mapReturnType,
	value_type ExpType,
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByIndex),
		IntegerValue(return_type),
		index,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, value_type, args)
}

// ExpMapGetByIndexRange creates an expression that selects map items starting at specified index to the end of map and returns selected
// data specified by returnType.
func ExpMapGetByIndexRange(
	return_type mapReturnType,
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByIndexRange),
		IntegerValue(return_type),
		index,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
}

// ExpMapGetByIndexRangeCount creates an expression that selects "count" map items starting at specified index and returns selected data
// specified by returnType.
func ExpMapGetByIndexRangeCount(
	return_type mapReturnType,
	index *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByIndexRange),
		IntegerValue(return_type),
		index,
		count,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
}

// ExpMapGetByRank creates an expression that selects map item identified by rank and returns selected data specified by returnType.
func ExpMapGetByRank(
	return_type mapReturnType,
	value_type ExpType,
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByRank),
		IntegerValue(return_type),
		rank,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, value_type, args)
}

// ExpMapGetByRankRange creates an expression that selects map items starting at specified rank to the last ranked item and
// returns selected data specified by returnType.
func ExpMapGetByRankRange(
	return_type mapReturnType,
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByRankRange),
		IntegerValue(return_type),
		rank,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
}

// ExpMapGetByRankRangeCount creates an expression that selects "count" map items starting at specified rank and returns selected
// data specified by returnType.
func ExpMapGetByRankRangeCount(
	return_type mapReturnType,
	rank *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(cdtMapOpTypeGetByRankRange),
		IntegerValue(return_type),
		rank,
		count,
		cdtContextList(ctx),
	}
	return expMapAddRead(bin, expMapGetValueType(return_type), args)
}

func expMapAddRead(
	bin *FilterExpression,
	return_type ExpType,
	arguments []ExpressionArgument,
) *FilterExpression {
	flags := expMapMODULE
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

func expMapAddWrite(
	bin *FilterExpression,
	arguments []ExpressionArgument,
	ctx ...*CDTContext,
) *FilterExpression {
	var return_type ExpType
	if len(ctx) == 0 {
		return_type = ExpTypeMAP
	} else if (ctx[0].id & ctxTypeListIndex) == 0 {
		return_type = ExpTypeMAP
	} else {
		return_type = ExpTypeLIST
	}

	flags := expMapMODULE | MODIFY
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

func expMapGetValueType(return_type mapReturnType) ExpType {
	var t = return_type & (^MapReturnType.INVERTED)
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
		} else {
			return cdtMapOpTypePut
		}
	case MapWriteFlagsUpdateOnly:
		if multi {
			return cdtMapOpTypeReplaceItems
		} else {
			return cdtMapOpTypeReplace
		}
	case MapWriteFlagsCreateOnly:
		if multi {
			return cdtMapOpTypeAddItems
		} else {
			return cdtMapOpTypeAdd
		}
	}
}
