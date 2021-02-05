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
// limitations under the License.package aerospike
package aerospike

const expListMODULE int64 = 0

// ExpListAppend creates an expression that appends value to end of list.
func ExpListAppend(
	policy *ListPolicy,
	value *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_APPEND),
		value,
		IntegerValue(policy.attributes),
		IntegerValue(policy.flags),
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListAppendItems creates an expression that appends list items to end of list.
func ExpListAppendItems(
	policy *ListPolicy,
	list *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_APPEND_ITEMS),
		list,
		IntegerValue(policy.attributes),
		IntegerValue(policy.flags),
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListInsert creates an expression that inserts value to specified index of list.
func ExpListInsert(
	policy *ListPolicy,
	index *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_INSERT),
		index,
		value,
		IntegerValue(policy.flags),
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListInsertItems creates an expression that inserts each input list item starting at specified index of list.
func ExpListInsertItems(
	policy *ListPolicy,
	index *FilterExpression,
	list *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_INSERT_ITEMS),
		index,
		list,
		IntegerValue(policy.flags),
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListIncrement creates an expression that increments `list[index]` by value.
// Value expression should resolve to a number.
func ExpListIncrement(
	policy *ListPolicy,
	index *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_INCREMENT),
		index,
		value,
		IntegerValue(policy.attributes),
		IntegerValue(policy.flags),
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListSet creates an expression that sets item value at specified index in list.
func ExpListSet(
	policy *ListPolicy,
	index *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_SET),
		index,
		value,
		IntegerValue(policy.flags),
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListClear creates an expression that removes all items in list.
func ExpListClear(bin *FilterExpression, ctx ...*CDTContext) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_CLEAR),
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListSort creates an expression that sorts list according to sortFlags.
func ExpListSort(
	sort_flags ListSortFlags,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_SORT),
		IntegerValue(sort_flags),
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByValue creates an expression that removes list items identified by value.
func ExpListRemoveByValue(
	value *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_VALUE),
		IntegerValue(ListReturnTypeNone),
		value,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByValueList creates an expression that removes list items identified by values.
func ExpListRemoveByValueList(
	values *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_VALUE_LIST),
		IntegerValue(ListReturnTypeNone),
		values,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByValueRange creates an expression that removes list items identified by value range (valueBegin inclusive, valueEnd exclusive).
// If valueBegin is null, the range is less than valueEnd. If valueEnd is null, the range is
// greater than equal to valueBegin.
func ExpListRemoveByValueRange(
	value_begin *FilterExpression,
	value_end *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(_CDT_LIST_REMOVE_BY_VALUE_INTERVAL),
		IntegerValue(ListReturnTypeNone),
	}
	if value_begin != nil {
		args = append(args, value_begin)
	} else {
		args = append(args, nullValue)
	}
	if value_end != nil {
		args = append(args, value_end)
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByValueRelativeRankRange creates an expression that removes list items nearest to value and greater by relative rank.
//
// Examples for ordered list \[0, 4, 5, 9, 11, 15\]:
// (value,rank) = [removed items]
// (5,0) = [5,9,11,15]
// (5,1) = [9,11,15]
// (5,-1) = [4,5,9,11,15]
// (3,0) = [4,5,9,11,15]
// (3,3) = [11,15]
// (3,-3) = [0,4,5,9,11,15]
func ExpListRemoveByValueRelativeRankRange(
	value *FilterExpression,
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_VALUE_REL_RANK_RANGE),
		IntegerValue(ListReturnTypeNone),
		value,
		rank,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByValueRelativeRankRangeCount creates an expression that removes list items nearest to value and greater by relative rank with a count limit.
//
// Examples for ordered list \[0, 4, 5, 9, 11, 15\]:
// (value,rank,count) = [removed items]
// (5,0,2) = [5,9]
// (5,1,1) = [9]
// (5,-1,2) = [4,5]
// (3,0,1) = [4]
// (3,3,7) = [11,15]
// (3,-3,2) = []
func ExpListRemoveByValueRelativeRankRangeCount(
	value *FilterExpression,
	rank *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_VALUE_REL_RANK_RANGE),
		IntegerValue(ListReturnTypeNone),
		value,
		rank,
		count,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByIndex creates an expression that removes list item identified by index.
func ExpListRemoveByIndex(
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_INDEX),
		IntegerValue(ListReturnTypeNone),
		index,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByIndexRange creates an expression that removes list items starting at specified index to the end of list.
func ExpListRemoveByIndexRange(
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_INDEX_RANGE),
		IntegerValue(ListReturnTypeNone),
		index,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByIndexRangeCount creates an expression that removes "count" list items starting at specified index.
func ExpListRemoveByIndexRangeCount(
	index *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_INDEX_RANGE),
		IntegerValue(ListReturnTypeNone),
		index,
		count,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByRank creates an expression that removes list item identified by rank.
func ExpListRemoveByRank(
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_RANK),
		IntegerValue(ListReturnTypeNone),
		rank,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByRankRange creates an expression that removes list items starting at specified rank to the last ranked item.
func ExpListRemoveByRankRange(
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_RANK_RANGE),
		IntegerValue(ListReturnTypeNone),
		rank,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListRemoveByRankRangeCount creates an expression that removes "count" list items starting at specified rank.
func ExpListRemoveByRankRangeCount(
	rank *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_REMOVE_BY_RANK_RANGE),
		IntegerValue(ListReturnTypeNone),
		rank,
		count,
		cdtContextList(ctx),
	}
	return cdtListAddWrite(bin, args, ctx...)
}

// ExpListSize creates an expression that returns list size.
func ExpListSize(bin *FilterExpression, ctx ...*CDTContext) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_SIZE),
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, ExpTypeINT, args)
}

// ExpListGetByValue creates an expression that selects list items identified by value and returns selected
// data specified by returnType.
func ExpListGetByValue(
	returnType ListReturnType,
	value *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_VALUE),
		IntegerValue(returnType),
		value,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, expListGetValueType(returnType), args)
}

// ExpListGetByValueRange creates an expression that selects list items identified by value range and returns selected data
// specified by returnType.
func ExpListGetByValueRange(
	returnType ListReturnType,
	value_begin *FilterExpression,
	value_end *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		cdtContextList(ctx),
		IntegerValue(_CDT_LIST_GET_BY_VALUE_INTERVAL),
		IntegerValue(returnType),
	}
	if value_begin != nil {
		args = append(args, value_begin)
	} else {
		args = append(args, nullValue)
	}
	if value_end != nil {
		args = append(args, value_end)
	}
	return cdtListAddRead(bin, expListGetValueType(returnType), args)
}

// ExpListGetByValueList creates an expression that selects list items identified by values and returns selected data
// specified by returnType.
func ExpListGetByValueList(
	returnType ListReturnType,
	values *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_VALUE_LIST),
		IntegerValue(returnType),
		values,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, expListGetValueType(returnType), args)
}

// ExpListGetByValueRelativeRankRange creates an expression that selects list items nearest to value and greater by relative rank
// and returns selected data specified by returnType.
//
// Examples for ordered list \[0, 4, 5, 9, 11, 15\]:
// (value,rank) = [selected items]
// (5,0) = [5,9,11,15]
// (5,1) = [9,11,15]
// (5,-1) = [4,5,9,11,15]
// (3,0) = [4,5,9,11,15]
// (3,3) = [11,15]
// (3,-3) = [0,4,5,9,11,15]
func ExpListGetByValueRelativeRankRange(
	returnType ListReturnType,
	value *FilterExpression,
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_VALUE_REL_RANK_RANGE),
		IntegerValue(returnType),
		value,
		rank,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, expListGetValueType(returnType), args)
}

// ExpListGetByValueRelativeRankRangeCount creates an expression that selects list items nearest to value and greater by relative rank with a count limit
// and returns selected data specified by returnType.
//
// Examples for ordered list \[0, 4, 5, 9, 11, 15\]:
// (value,rank,count) = [selected items]
// (5,0,2) = [5,9]
// (5,1,1) = [9]
// (5,-1,2) = [4,5]
// (3,0,1) = [4]
// (3,3,7) = [11,15]
// (3,-3,2) = []
func ExpListGetByValueRelativeRankRangeCount(
	returnType ListReturnType,
	value *FilterExpression,
	rank *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_VALUE_REL_RANK_RANGE),
		IntegerValue(returnType),
		value,
		rank,
		count,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, expListGetValueType(returnType), args)
}

// ExpListGetByIndex creates an expression that selects list item identified by index and returns
// selected data specified by returnType.
func ExpListGetByIndex(
	returnType ListReturnType,
	value_type ExpType,
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_INDEX),
		IntegerValue(returnType),
		index,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, value_type, args)
}

// ExpListGetByIndexRange creates an expression that selects list items starting at specified index to the end of list
// and returns selected data specified by returnType .
func ExpListGetByIndexRange(
	returnType ListReturnType,
	index *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_INDEX_RANGE),
		IntegerValue(returnType),
		index,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, expListGetValueType(returnType), args)
}

// ExpListGetByIndexRangeCount creates an expression that selects "count" list items starting at specified index
// and returns selected data specified by returnType.
func ExpListGetByIndexRangeCount(
	returnType ListReturnType,
	index *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_INDEX_RANGE),
		IntegerValue(returnType),
		index,
		count,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, expListGetValueType(returnType), args)
}

// ExpListGetByRank creates an expression that selects list item identified by rank and returns selected
// data specified by returnType.
func ExpListGetByRank(
	returnType ListReturnType,
	value_type ExpType,
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_RANK),
		IntegerValue(returnType),
		rank,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, value_type, args)
}

// ExpListGetByRankRange creates an expression that selects list items starting at specified rank to the last ranked item
// and returns selected data specified by returnType.
func ExpListGetByRankRange(
	returnType ListReturnType,
	rank *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_RANK_RANGE),
		IntegerValue(returnType),
		rank,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, expListGetValueType(returnType), args)
}

// ExpListGetByRankRangeCount creates an expression that selects "count" list items starting at specified rank and returns
// selected data specified by returnType.
func ExpListGetByRankRangeCount(
	returnType ListReturnType,
	rank *FilterExpression,
	count *FilterExpression,
	bin *FilterExpression,
	ctx ...*CDTContext,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(_CDT_LIST_GET_BY_RANK_RANGE),
		IntegerValue(returnType),
		rank,
		count,
		cdtContextList(ctx),
	}
	return cdtListAddRead(bin, expListGetValueType(returnType), args)
}

func cdtListAddRead(
	bin *FilterExpression,
	returnType ExpType,
	arguments []ExpressionArgument,
) *FilterExpression {
	flags := expListMODULE
	return &FilterExpression{
		cmd:       &expOpCALL,
		val:       nil,
		bin:       bin,
		flags:     &flags,
		module:    &returnType,
		exps:      nil,
		arguments: arguments,
	}
}

func cdtListAddWrite(
	bin *FilterExpression,
	arguments []ExpressionArgument,
	ctx ...*CDTContext,
) *FilterExpression {
	var returnType ExpType
	if len(ctx) == 0 {
		returnType = ExpTypeLIST
	} else if (ctx[0].id & ctxTypeListIndex) == 0 {
		returnType = ExpTypeMAP
	} else {
		returnType = ExpTypeLIST
	}

	flags := expListMODULE | MODIFY
	return &FilterExpression{
		cmd:       &expOpCALL,
		val:       nil,
		bin:       bin,
		flags:     &flags,
		module:    &returnType,
		exps:      nil,
		arguments: arguments,
	}
}

func expListGetValueType(returnType ListReturnType) ExpType {
	if (returnType & (^ListReturnTypeInverted)) == ListReturnTypeValue {
		return ExpTypeLIST
	}
	return ExpTypeINT
}
