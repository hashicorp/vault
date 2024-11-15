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

// List operations support negative indexing.  If the index is negative, the
// resolved index starts backwards from end of list. If an index is out of bounds,
// a parameter error will be returned. If a range is partially out of bounds, the
// valid part of the range will be returned. Index/Range examples:
//
// Index/Range examples:
//
//    Index 0: First item in list.
//    Index 4: Fifth item in list.
//    Index -1: Last item in list.
//    Index -3: Third to last item in list.
//    Index 1 Count 2: Second and third items in list.
//    Index -3 Count 3: Last three items in list.
//    Index -5 Count 4: Range between fifth to last item to second to last item inclusive.
//
// Nested CDT operations are supported by optional Ctx context arguments.  Examples:
//
// bin = [[7,9,5],[1,2,3],[6,5,4,1]]
// Append 11 to last list.
// ListAppendWithPolicyContextOp(DefaultMapPolicy(), "bin", IntegerValue(11), CtxListIndex(-1))
// bin result = [[7,9,5],[1,2,3],[6,5,4,1,11]]
//
// bin = {key1:[[7,9,5],[13]], key2:[[9],[2,4],[6,1,9]], key3:[[6,5]]}
// Append 11 to lowest ranked list in map identified by "key2".
// ListAppendWithPolicyContextOp(DefaultMapPolicy(), "bin", IntegerValue(11), CtxMapKey(StringValue("key2")), CtxListRank(0))
// bin result = {key1:[[7,9,5],[13]], key2:[[9],[2,4,11],[6,1,9]], key3:[[6,5]]}

const (
	_CDT_LIST_SET_TYPE                       = 0
	_CDT_LIST_APPEND                         = 1
	_CDT_LIST_APPEND_ITEMS                   = 2
	_CDT_LIST_INSERT                         = 3
	_CDT_LIST_INSERT_ITEMS                   = 4
	_CDT_LIST_POP                            = 5
	_CDT_LIST_POP_RANGE                      = 6
	_CDT_LIST_REMOVE                         = 7
	_CDT_LIST_REMOVE_RANGE                   = 8
	_CDT_LIST_SET                            = 9
	_CDT_LIST_TRIM                           = 10
	_CDT_LIST_CLEAR                          = 11
	_CDT_LIST_INCREMENT                      = 12
	_CDT_LIST_SORT                           = 13
	_CDT_LIST_SIZE                           = 16
	_CDT_LIST_GET                            = 17
	_CDT_LIST_GET_RANGE                      = 18
	_CDT_LIST_GET_BY_INDEX                   = 19
	_CDT_LIST_GET_BY_RANK                    = 21
	_CDT_LIST_GET_BY_VALUE                   = 22
	_CDT_LIST_GET_BY_VALUE_LIST              = 23
	_CDT_LIST_GET_BY_INDEX_RANGE             = 24
	_CDT_LIST_GET_BY_VALUE_INTERVAL          = 25
	_CDT_LIST_GET_BY_RANK_RANGE              = 26
	_CDT_LIST_GET_BY_VALUE_REL_RANK_RANGE    = 27
	_CDT_LIST_REMOVE_BY_INDEX                = 32
	_CDT_LIST_REMOVE_BY_RANK                 = 34
	_CDT_LIST_REMOVE_BY_VALUE                = 35
	_CDT_LIST_REMOVE_BY_VALUE_LIST           = 36
	_CDT_LIST_REMOVE_BY_INDEX_RANGE          = 37
	_CDT_LIST_REMOVE_BY_VALUE_INTERVAL       = 38
	_CDT_LIST_REMOVE_BY_RANK_RANGE           = 39
	_CDT_LIST_REMOVE_BY_VALUE_REL_RANK_RANGE = 40
)

// ListOrderType determines the order of returned values in CDT list operations.
type ListOrderType int

// Map storage order.
const (
	// ListOrderUnordered signifies that list is not ordered. This is the default.
	ListOrderUnordered ListOrderType = 0

	// ListOrderOrdered signifies that list is Ordered.
	ListOrderOrdered ListOrderType = 1
)

// ListPolicy directives when creating a list and writing list items.
type ListPolicy struct {
	attributes ListOrderType
	flags      int
}

// NewListPolicy creates a policy with directives when creating a list and writing list items.
// Flags are ListWriteFlags. You can specify multiple by `or`ing them together.
func NewListPolicy(order ListOrderType, flags int) *ListPolicy {
	return &ListPolicy{
		attributes: order,
		flags:      flags,
	}
}

// DefaultListPolicy is the default list policy and can be customized.
var defaultListPolicy = NewListPolicy(ListOrderUnordered, ListWriteFlagsDefault)

// DefaultListPolicy returns the default policy for CDT list operations.
func DefaultListPolicy() *ListPolicy {
	return defaultListPolicy
}

// ListReturnType determines the returned values in CDT List operations.
type ListReturnType int

const (
	// ListReturnTypeNone will not return a result.
	ListReturnTypeNone ListReturnType = 0

	// ListReturnTypeIndex will return index offset order.
	// 0 = first key
	// N = Nth key
	// -1 = last key
	ListReturnTypeIndex ListReturnType = 1

	// ListReturnTypeReverseIndex will return reverse index offset order.
	// 0 = last key
	// -1 = first key
	ListReturnTypeReverseIndex ListReturnType = 2

	// ListReturnTypeRank will return value order.
	// 0 = smallest value
	// N = Nth smallest value
	// -1 = largest value
	ListReturnTypeRank ListReturnType = 3

	// ListReturnTypeReverseRank will return reverse value order.
	// 0 = largest value
	// N = Nth largest value
	// -1 = smallest value
	ListReturnTypeReverseRank ListReturnType = 4

	// ListReturnTypeCount will return count of items selected.
	ListReturnTypeCount ListReturnType = 5

	// ListReturnTypeValue will return value for single key read and value list for range read.
	ListReturnTypeValue ListReturnType = 7

	// ListReturnTypeInverted will invert meaning of list command and return values.  For example:
	// ListOperation.getByIndexRange(binName, index, count, ListReturnType.INDEX | ListReturnType.INVERTED)
	// With the INVERTED flag enabled, the items outside of the specified index range will be returned.
	// The meaning of the list command can also be inverted.  For example:
	// ListOperation.removeByIndexRange(binName, index, count, ListReturnType.INDEX | ListReturnType.INVERTED);
	// With the INVERTED flag enabled, the items outside of the specified index range will be removed and returned.
	ListReturnTypeInverted ListReturnType = 0x10000
)

// ListSortFlags detemines sort flags for CDT lists
type ListSortFlags int

const (
	// ListSortFlagsDefault is the default sort flag for CDT lists, and sort in Ascending order.
	ListSortFlagsDefault ListSortFlags = 0
	// ListSortFlagsDescending will sort the contents of the list in descending order.
	ListSortFlagsDescending ListSortFlags = 1
	// ListSortFlagsDropDuplicates will drop duplicate values in the results of the CDT list operation.
	ListSortFlagsDropDuplicates ListSortFlags = 2
)

// ListWriteFlags detemines write flags for CDT lists
// type ListWriteFlags int

const (
	// ListWriteFlagsDefault is the default behavior. It means:  Allow duplicate values and insertions at any index.
	ListWriteFlagsDefault = 0
	// ListWriteFlagsAddUnique means: Only add unique values.
	ListWriteFlagsAddUnique = 1
	// ListWriteFlagsInsertBounded means: Enforce list boundaries when inserting.  Do not allow values to be inserted
	// at index outside current list boundaries.
	ListWriteFlagsInsertBounded = 2
	// ListWriteFlagsNoFail means: do not raise error if a list item fails due to write flag constraints.
	ListWriteFlagsNoFail = 4
	// ListWriteFlagsPartial means: allow other valid list items to be committed if a list item fails due to
	// write flag constraints.
	ListWriteFlagsPartial = 8
)

func listGenericOpEncoder(op *Operation, packer BufferEx) (int, Error) {
	args := op.binValue.(ListValue)
	if len(args) > 1 {
		return packCDTIfcVarParamsAsArray(packer, int16(args[0].(int)), op.ctx, args[1:]...)
	}
	return packCDTIfcVarParamsAsArray(packer, int16(args[0].(int)), op.ctx)
}

func packCDTParamsAsArray(packer BufferEx, opType int16, ctx []*CDTContext, params ...Value) (int, Error) {
	size := 0
	n := 0
	var err Error
	if len(ctx) > 0 {
		if n, err = packArrayBegin(packer, 3); err != nil {
			return size + n, err
		}
		size += n

		if n, err = packAInt64(packer, 0xff); err != nil {
			return size + n, err
		}
		size += n

		if n, err = packArrayBegin(packer, len(ctx)*2); err != nil {
			return size + n, err
		}
		size += n

		for _, c := range ctx {
			if n, err = packAInt64(packer, int64(c.id)); err != nil {
				return size + n, err
			}
			size += n

			if n, err = c.value.pack(packer); err != nil {
				return size + n, err
			}
			size += n
		}

		if n, err = packArrayBegin(packer, len(params)+1); err != nil {
			return size + n, err
		}
		size += n

		if n, err = packAInt(packer, int(opType)); err != nil {
			return size + n, err
		}
		size += n
	} else {
		if n, err = packShortRaw(packer, opType); err != nil {
			return n, err
		}
		size += n

		if len(params) > 0 {
			if n, err = packArrayBegin(packer, len(params)); err != nil {
				return size + n, err
			}
			size += n
		}
	}

	if len(params) > 0 {
		for i := range params {
			if n, err = params[i].pack(packer); err != nil {
				return size + n, err
			}
			size += n
		}
	}
	return size, nil
}

func packCDTIfcParamsAsArray(packer BufferEx, opType int16, ctx []*CDTContext, params ListValue) (int, Error) {
	return packCDTIfcVarParamsAsArray(packer, opType, ctx, []interface{}(params)...)
}

func packCDTIfcVarParamsAsArray(packer BufferEx, opType int16, ctx []*CDTContext, params ...interface{}) (int, Error) {
	size := 0
	n := 0
	var err Error
	if len(ctx) > 0 {
		if n, err = packArrayBegin(packer, 3); err != nil {
			return size + n, err
		}
		size += n

		if n, err = packAInt64(packer, 0xff); err != nil {
			return size + n, err
		}
		size += n

		if n, err = packArrayBegin(packer, len(ctx)*2); err != nil {
			return size + n, err
		}
		size += n

		for _, c := range ctx {
			if n, err = packAInt64(packer, int64(c.id)); err != nil {
				return size + n, err
			}
			size += n

			if n, err = c.value.pack(packer); err != nil {
				return size + n, err
			}
			size += n
		}

		if n, err = packArrayBegin(packer, len(params)+1); err != nil {
			return size + n, err
		}
		size += n

		if n, err = packAInt(packer, int(opType)); err != nil {
			return size + n, err
		}
		size += n
	} else {
		n, err = packShortRaw(packer, opType)
		if err != nil {
			return n, err
		}
		size += n

		if len(params) > 0 {
			if n, err = packArrayBegin(packer, len(params)); err != nil {
				return size + n, err
			}
			size += n
		}
	}

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

func cdtCreateOpEncoder(op *Operation, packer BufferEx) (int, Error) {
	args := op.binValue.(ListValue)
	if len(args) > 2 {
		return packCDTCreate(packer, int16(args[0].(int)), op.ctx, args[1].(int), args[2:]...)
	}
	return packCDTCreate(packer, int16(args[0].(int)), op.ctx, args[1].(int))
}

func packCDTCreate(packer BufferEx, opType int16, ctx []*CDTContext, flag int, params ...interface{}) (int, Error) {
	size := 0
	n := 0
	var err Error
	if n, err = packArrayBegin(packer, 3); err != nil {
		return size + n, err
	}
	size += n

	if n, err = packAInt64(packer, 0xff); err != nil {
		return size + n, err
	}
	size += n

	if n, err = packArrayBegin(packer, len(ctx)*2); err != nil {
		return size + n, err
	}
	size += n

	var c *CDTContext
	last := len(ctx) - 1

	for i := 0; i < last; i++ {
		c = ctx[i]
		if n, err = packAInt64(packer, int64(c.id)); err != nil {
			return size + n, err
		}
		size += n

		if n, err = c.value.pack(packer); err != nil {
			return size + n, err
		}
		size += n
	}

	c = ctx[last]
	if n, err = packAInt64(packer, int64(c.id|flag)); err != nil {
		return size + n, err
	}
	size += n

	if n, err = c.value.pack(packer); err != nil {
		return size + n, err
	}
	size += n

	if n, err = packArrayBegin(packer, len(params)+1); err != nil {
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

func cdtListOrderFlag(order ListOrderType, pad bool) int {
	if order == 1 {
		return 0xc0
	}

	if pad {
		return 0x80
	}
	return 0x40
}

// ListCreateOp creates list create operation.
// Server creates list at given context level. The context is allowed to be beyond list
// boundaries only if pad is set to true.  In that case, nil list entries will be inserted to
// satisfy the context position.
func ListCreateOp(binName string, listOrder ListOrderType, pad bool, ctx ...*CDTContext) *Operation {
	// If context not defined, the set order for top-level bin list.
	if len(ctx) == 0 {
		return ListSetOrderOp(binName, listOrder)
	}
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_SET_TYPE, cdtListOrderFlag(listOrder, pad), IntegerValue(listOrder)}, encoder: cdtCreateOpEncoder}
}

// ListSetOrderOp creates a set list order operation.
// Server sets list order.  Server returns null.
func ListSetOrderOp(binName string, listOrder ListOrderType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_SET_TYPE, IntegerValue(listOrder)}, encoder: listGenericOpEncoder}
}

// ListAppendOp creates a list append operation.
// Server appends values to end of list bin.
// Server returns list size on bin name.
// It will panic is no values have been passed.
func ListAppendOp(binName string, values ...interface{}) *Operation {
	if len(values) == 1 {
		return &Operation{opType: _CDT_MODIFY, binName: binName, binValue: ListValue{_CDT_LIST_APPEND, NewValue(values[0])}, encoder: listGenericOpEncoder}
	}
	return &Operation{opType: _CDT_MODIFY, binName: binName, binValue: ListValue{_CDT_LIST_APPEND_ITEMS, ListValue(values)}, encoder: listGenericOpEncoder}
}

// ListAppendWithPolicyOp creates a list append operation.
// Server appends values to end of list bin.
// Server returns list size on bin name.
// It will panic is no values have been passed.
func ListAppendWithPolicyOp(policy *ListPolicy, binName string, values ...interface{}) *Operation {
	switch len(values) {
	case 1:
		return &Operation{opType: _CDT_MODIFY, binName: binName, binValue: ListValue{_CDT_LIST_APPEND, NewValue(values[0]), IntegerValue(policy.attributes), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
	default:
		return &Operation{opType: _CDT_MODIFY, binName: binName, binValue: ListValue{_CDT_LIST_APPEND_ITEMS, ListValue(values), IntegerValue(policy.attributes), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
	}
}

// ListAppendWithPolicyContextOp creates a list append operation.
// Server appends values to end of list bin.
// Server returns list size on bin name.
// It will panic is no values have been passed.
func ListAppendWithPolicyContextOp(policy *ListPolicy, binName string, ctx []*CDTContext, values ...interface{}) *Operation {
	switch len(values) {
	case 1:
		return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_APPEND, NewValue(values[0]), IntegerValue(policy.attributes), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
	default:
		return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_APPEND_ITEMS, ListValue(values), IntegerValue(policy.attributes), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
	}
}

// ListInsertOp creates a list insert operation.
// Server inserts value to specified index of list bin.
// Server returns list size on bin name.
// It will panic is no values have been passed.
func ListInsertOp(binName string, index int, values ...interface{}) *Operation {
	if len(values) == 1 {
		return &Operation{opType: _CDT_MODIFY, binName: binName, binValue: ListValue{_CDT_LIST_INSERT, IntegerValue(index), NewValue(values[0])}, encoder: listGenericOpEncoder}
	}
	return &Operation{opType: _CDT_MODIFY, binName: binName, binValue: ListValue{_CDT_LIST_INSERT_ITEMS, IntegerValue(index), ListValue(values)}, encoder: listGenericOpEncoder}
}

// ListInsertWithPolicyOp creates a list insert operation.
// Server inserts value to specified index of list bin.
// Server returns list size on bin name.
// It will panic is no values have been passed.
func ListInsertWithPolicyOp(policy *ListPolicy, binName string, index int, values ...interface{}) *Operation {
	if len(values) == 1 {
		return &Operation{opType: _CDT_MODIFY, binName: binName, binValue: ListValue{_CDT_LIST_INSERT, IntegerValue(index), NewValue(values[0]), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
	}
	return &Operation{opType: _CDT_MODIFY, binName: binName, binValue: ListValue{_CDT_LIST_INSERT_ITEMS, IntegerValue(index), ListValue(values), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
}

// ListInsertWithPolicyContextOp creates a list insert operation.
// Server inserts value to specified index of list bin.
// Server returns list size on bin name.
// It will panic is no values have been passed.
func ListInsertWithPolicyContextOp(policy *ListPolicy, binName string, index int, ctx []*CDTContext, values ...interface{}) *Operation {
	if len(values) == 1 {
		return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_INSERT, IntegerValue(index), NewValue(values[0]), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
	}
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_INSERT_ITEMS, IntegerValue(index), ListValue(values), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
}

// ListPopOp creates list pop operation.
// Server returns item at specified index and removes item from list bin.
func ListPopOp(binName string, index int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_POP, IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListPopRangeOp creates a list pop range operation.
// Server returns items starting at specified index and removes items from list bin.
func ListPopRangeOp(binName string, index int, count int, ctx ...*CDTContext) *Operation {
	if count == 1 {
		return ListPopOp(binName, index)
	}

	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_POP_RANGE, IntegerValue(index), IntegerValue(count)}, encoder: listGenericOpEncoder}
}

// ListPopRangeFromOp creates a list pop range operation.
// Server returns items starting at specified index to the end of list and removes items from list bin.
func ListPopRangeFromOp(binName string, index int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_POP_RANGE, IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListRemoveOp creates a list remove operation.
// Server removes item at specified index from list bin.
// Server returns number of items removed.
func ListRemoveOp(binName string, index int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE, IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListRemoveByValueOp creates list remove by value operation.
// Server removes the item identified by value and returns removed data specified by returnType.
func ListRemoveByValueOp(binName string, value interface{}, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_VALUE, IntegerValue(returnType), NewValue(value)}, encoder: listGenericOpEncoder}
}

// ListRemoveByValueListOp creates list remove by value operation.
// Server removes list items identified by value and returns removed data specified by returnType.
func ListRemoveByValueListOp(binName string, values []interface{}, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_VALUE_LIST, IntegerValue(returnType), ListValue(values)}, encoder: listGenericOpEncoder}
}

// ListRemoveByValueRangeOp creates a list remove operation.
// Server removes list items identified by value range (valueBegin inclusive, valueEnd exclusive).
// If valueBegin is nil, the range is less than valueEnd.
// If valueEnd is nil, the range is greater than equal to valueBegin.
// Server returns removed data specified by returnType
func ListRemoveByValueRangeOp(binName string, returnType ListReturnType, valueBegin, valueEnd interface{}, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_VALUE_INTERVAL, IntegerValue(returnType), NewValue(valueBegin), NewValue(valueEnd)}, encoder: listGenericOpEncoder}
}

// ListRemoveByValueRelativeRankRangeOp creates a list remove by value relative to rank range operation.
// Server removes list items nearest to value and greater by relative rank.
// Server returns removed data specified by returnType.
//
// Examples for ordered list [0,4,5,9,11,15]:
//
//  (value,rank) = [removed items]
//  (5,0) = [5,9,11,15]
//  (5,1) = [9,11,15]
//  (5,-1) = [4,5,9,11,15]
//  (3,0) = [4,5,9,11,15]
//  (3,3) = [11,15]
//  (3,-3) = [0,4,5,9,11,15]
func ListRemoveByValueRelativeRankRangeOp(binName string, returnType ListReturnType, value interface{}, rank int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_VALUE_REL_RANK_RANGE, IntegerValue(returnType), NewValue(value), IntegerValue(rank)}, encoder: listGenericOpEncoder}
}

// ListRemoveByValueRelativeRankRangeCountOp creates a list remove by value relative to rank range operation.
// Server removes list items nearest to value and greater by relative rank with a count limit.
// Server returns removed data specified by returnType.
// Examples for ordered list [0,4,5,9,11,15]:
//
//  (value,rank,count) = [removed items]
//  (5,0,2) = [5,9]
//  (5,1,1) = [9]
//  (5,-1,2) = [4,5]
//  (3,0,1) = [4]
//  (3,3,7) = [11,15]
//  (3,-3,2) = []
func ListRemoveByValueRelativeRankRangeCountOp(binName string, returnType ListReturnType, value interface{}, rank, count int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_VALUE_REL_RANK_RANGE, IntegerValue(returnType), NewValue(value), IntegerValue(rank), IntegerValue(count)}, encoder: listGenericOpEncoder}
}

// ListRemoveRangeOp creates a list remove range operation.
// Server removes "count" items starting at specified index from list bin.
// Server returns number of items removed.
func ListRemoveRangeOp(binName string, index int, count int, ctx ...*CDTContext) *Operation {
	if count == 1 {
		return ListRemoveOp(binName, index)
	}

	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_RANGE, IntegerValue(index), IntegerValue(count)}, encoder: listGenericOpEncoder}
}

// ListRemoveRangeFromOp creates a list remove range operation.
// Server removes all items starting at specified index to the end of list.
// Server returns number of items removed.
func ListRemoveRangeFromOp(binName string, index int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_RANGE, IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListSetOp creates a list set operation.
// Server sets item value at specified index in list bin.
// Server does not return a result by default.
func ListSetOp(binName string, index int, value interface{}, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_SET, IntegerValue(index), NewValue(value)}, encoder: listGenericOpEncoder}
}

// ListTrimOp creates a list trim operation.
// Server removes items in list bin that do not fall into range specified by index
// and count range. If the range is out of bounds, then all items will be removed.
// Server returns number of elemts that were removed.
func ListTrimOp(binName string, index int, count int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_TRIM, IntegerValue(index), IntegerValue(count)}, encoder: listGenericOpEncoder}
}

// ListClearOp creates a list clear operation.
// Server removes all items in list bin.
// Server does not return a result by default.
func ListClearOp(binName string, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_CLEAR}, encoder: listGenericOpEncoder}
}

// ListIncrementOp creates a list increment operation.
// Server increments list[index] by value.
// Value should be integer(IntegerValue, LongValue) or float(FloatValue).
// Server returns list[index] after incrementing.
func ListIncrementOp(binName string, index int, value interface{}, ctx ...*CDTContext) *Operation {
	val := NewValue(value)
	switch val.(type) {
	case LongValue, IntegerValue, FloatValue:
	default:
		panic("Increment operation only accepts Integer or Float values")
	}
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_INCREMENT, IntegerValue(index), NewValue(value)}, encoder: listGenericOpEncoder}
}

// ListIncrementByOneOp creates list increment operation with policy.
// Server increments list[index] by 1.
// Server returns list[index] after incrementing.
func ListIncrementByOneOp(binName string, index int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_INCREMENT, IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListIncrementByOneWithPolicyOp creates list increment operation with policy.
// Server increments list[index] by 1.
// Server returns list[index] after incrementing.
func ListIncrementByOneWithPolicyOp(policy *ListPolicy, binName string, index int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_INCREMENT, IntegerValue(index), IntegerValue(1), IntegerValue(policy.attributes), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
}

// ListIncrementWithPolicyOp creates a list increment operation.
// Server increments list[index] by value.
// Server returns list[index] after incrementing.
func ListIncrementWithPolicyOp(policy *ListPolicy, binName string, index int, value interface{}, ctx ...*CDTContext) *Operation {
	val := NewValue(value)
	switch val.(type) {
	case LongValue, IntegerValue, FloatValue:
	default:
		panic("Increment operation only accepts Integer or Float values")
	}
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_INCREMENT, IntegerValue(index), NewValue(value), IntegerValue(policy.flags)}, encoder: listGenericOpEncoder}
}

// ListSizeOp creates a list size operation.
// Server returns size of list on bin name.
func ListSizeOp(binName string, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_SIZE}, encoder: listGenericOpEncoder}
}

// ListGetOp creates a list get operation.
// Server returns item at specified index in list bin.
func ListGetOp(binName string, index int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET, IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListGetRangeOp creates a list get range operation.
// Server returns "count" items starting at specified index in list bin.
func ListGetRangeOp(binName string, index int, count int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_RANGE, IntegerValue(index), IntegerValue(count)}, encoder: listGenericOpEncoder}
}

// ListGetRangeFromOp creates a list get range operation.
// Server returns items starting at specified index to the end of list.
func ListGetRangeFromOp(binName string, index int, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_RANGE, IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListSortOp creates list sort operation.
// Server sorts list according to sortFlags.
// Server does not return a result by default.
func ListSortOp(binName string, sortFlags ListSortFlags, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_SORT, IntegerValue(sortFlags)}, encoder: listGenericOpEncoder}
}

// ListRemoveByIndexOp creates a list remove operation.
// Server removes list item identified by index and returns removed data specified by returnType.
func ListRemoveByIndexOp(binName string, index int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_INDEX, IntegerValue(returnType), IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListRemoveByIndexRangeOp creates a list remove operation.
// Server removes list items starting at specified index to the end of list and returns removed
// data specified by returnType.
func ListRemoveByIndexRangeOp(binName string, index, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_INDEX_RANGE, IntegerValue(returnType), IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListRemoveByIndexRangeCountOp creates a list remove operation.
// Server removes "count" list items starting at specified index and returns removed data specified by returnType.
func ListRemoveByIndexRangeCountOp(binName string, index, count int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_INDEX_RANGE, IntegerValue(returnType), IntegerValue(index), IntegerValue(count)}, encoder: listGenericOpEncoder}
}

// ListRemoveByRankOp creates a list remove operation.
// Server removes list item identified by rank and returns removed data specified by returnType.
func ListRemoveByRankOp(binName string, rank int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_RANK, IntegerValue(returnType), IntegerValue(rank)}, encoder: listGenericOpEncoder}
}

// ListRemoveByRankRangeOp creates a list remove operation.
// Server removes list items starting at specified rank to the last ranked item and returns removed
// data specified by returnType.
func ListRemoveByRankRangeOp(binName string, rank int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_RANK_RANGE, IntegerValue(returnType), IntegerValue(rank)}, encoder: listGenericOpEncoder}
}

// ListRemoveByRankRangeCountOp creates a list remove operation.
// Server removes "count" list items starting at specified rank and returns removed data specified by returnType.
func ListRemoveByRankRangeCountOp(binName string, rank int, count int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_MODIFY, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_REMOVE_BY_RANK_RANGE, IntegerValue(returnType), IntegerValue(rank), IntegerValue(count)}, encoder: listGenericOpEncoder}
}

// ListGetByValueOp creates a list get by value operation.
// Server selects list items identified by value and returns selected data specified by returnType.
func ListGetByValueOp(binName string, value interface{}, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_VALUE, IntegerValue(returnType), NewValue(value)}, encoder: listGenericOpEncoder}
}

// ListGetByValueListOp creates list get by value list operation.
// Server selects list items identified by values and returns selected data specified by returnType.
func ListGetByValueListOp(binName string, values []interface{}, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_VALUE_LIST, IntegerValue(returnType), ListValue(values)}, encoder: listGenericOpEncoder}
}

// ListGetByValueRangeOp creates a list get by value range operation.
// Server selects list items identified by value range (valueBegin inclusive, valueEnd exclusive)
// If valueBegin is null, the range is less than valueEnd.
// If valueEnd is null, the range is greater than equal to valueBegin.
// Server returns selected data specified by returnType.
func ListGetByValueRangeOp(binName string, beginValue, endValue interface{}, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_VALUE_INTERVAL, IntegerValue(returnType), NewValue(beginValue), NewValue(endValue)}, encoder: listGenericOpEncoder}
}

// ListGetByIndexOp creates list get by index operation.
// Server selects list item identified by index and returns selected data specified by returnType
func ListGetByIndexOp(binName string, index int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_INDEX, IntegerValue(returnType), IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListGetByIndexRangeOp creates list get by index range operation.
// Server selects list items starting at specified index to the end of list and returns selected
// data specified by returnType.
func ListGetByIndexRangeOp(binName string, index int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_INDEX_RANGE, IntegerValue(returnType), IntegerValue(index)}, encoder: listGenericOpEncoder}
}

// ListGetByIndexRangeCountOp creates list get by index range operation.
// Server selects "count" list items starting at specified index and returns selected data specified
// by returnType.
func ListGetByIndexRangeCountOp(binName string, index, count int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_INDEX_RANGE, IntegerValue(returnType), IntegerValue(index), count}, encoder: listGenericOpEncoder}
}

// ListGetByRankOp creates a list get by rank operation.
// Server selects list item identified by rank and returns selected data specified by returnType.
func ListGetByRankOp(binName string, rank int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_RANK, IntegerValue(returnType), IntegerValue(rank)}, encoder: listGenericOpEncoder}
}

// ListGetByRankRangeOp creates a list get by rank range operation.
// Server selects list items starting at specified rank to the last ranked item and returns selected
// data specified by returnType.
func ListGetByRankRangeOp(binName string, rank int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_RANK_RANGE, IntegerValue(returnType), IntegerValue(rank)}, encoder: listGenericOpEncoder}
}

// ListGetByRankRangeCountOp creates a list get by rank range operation.
// Server selects "count" list items starting at specified rank and returns selected data specified by returnType.
func ListGetByRankRangeCountOp(binName string, rank, count int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_RANK_RANGE, IntegerValue(returnType), IntegerValue(rank), IntegerValue(count)}, encoder: listGenericOpEncoder}
}

// ListGetByValueRelativeRankRangeOp creates a list get by value relative to rank range operation.
// Server selects list items nearest to value and greater by relative rank.
// Server returns selected data specified by returnType.
//
// Examples for ordered list [0,4,5,9,11,15]:
//
//  (value,rank) = [selected items]
//  (5,0) = [5,9,11,15]
//  (5,1) = [9,11,15]
//  (5,-1) = [4,5,9,11,15]
//  (3,0) = [4,5,9,11,15]
//  (3,3) = [11,15]
//  (3,-3) = [0,4,5,9,11,15]
func ListGetByValueRelativeRankRangeOp(binName string, value interface{}, rank int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_VALUE_REL_RANK_RANGE, IntegerValue(returnType), NewValue(value), IntegerValue(rank)}, encoder: listGenericOpEncoder}
}

// ListGetByValueRelativeRankRangeCountOp creates a list get by value relative to rank range operation.
// Server selects list items nearest to value and greater by relative rank with a count limit.
// Server returns selected data specified by returnType.
//
// Examples for ordered list [0,4,5,9,11,15]:
//
//  (value,rank,count) = [selected items]
//  (5,0,2) = [5,9]
//  (5,1,1) = [9]
//  (5,-1,2) = [4,5]
//  (3,0,1) = [4]
//  (3,3,7) = [11,15]
//  (3,-3,2) = []
func ListGetByValueRelativeRankRangeCountOp(binName string, value interface{}, rank, count int, returnType ListReturnType, ctx ...*CDTContext) *Operation {
	return &Operation{opType: _CDT_READ, ctx: ctx, binName: binName, binValue: ListValue{_CDT_LIST_GET_BY_VALUE_REL_RANK_RANGE, IntegerValue(returnType), NewValue(value), IntegerValue(rank), IntegerValue(count)}, encoder: listGenericOpEncoder}
}
