// Copyright 2014-2021 Aerospike, Inc.
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

// Unique key map bin operations. Create map operations used by the client operate command.
// The default unique key map is unordered.
//
// All maps maintain an index and a rank.  The index is the item offset from the start of the map,
// for both unordered and ordered maps.  The rank is the sorted index of the value component.
// Map supports negative indexing for index and rank.
//
// Index examples:
//
//  Index 0: First item in map.
//  Index 4: Fifth item in map.
//  Index -1: Last item in map.
//  Index -3: Third to last item in map.
//  Index 1 Count 2: Second and third items in map.
//  Index -3 Count 3: Last three items in map.
//  Index -5 Count 4: Range between fifth to last item to second to last item inclusive.
//
//
// Rank examples:
//
//  Rank 0: Item with lowest value rank in map.
//  Rank 4: Fifth lowest ranked item in map.
//  Rank -1: Item with highest ranked value in map.
//  Rank -3: Item with third highest ranked value in map.
//  Rank 1 Count 2: Second and third lowest ranked items in map.
//  Rank -3 Count 3: Top three ranked items in map.
//
//
// Nested CDT operations are supported by optional CTX context arguments.  Examples:
//
//  bin = {key1:{key11:9,key12:4}, key2:{key21:3,key22:5}}
//  Set map value to 11 for map key "key21" inside of map key "key2".
//  MapOperation.put(MapPolicy.Default, "bin", StringValue("key21"), IntegerValue(11), CtxMapKey(StringValue("key2")))
//  bin result = {key1:{key11:9,key12:4},key2:{key21:11,key22:5}}
//
//  bin : {key1:{key11:{key111:1},key12:{key121:5}}, key2:{key21:{"key211":7}}}
//  Set map value to 11 in map key "key121" for highest ranked map ("key12") inside of map key "key1".
//  MapPutOp(DefaultMapPolicy(), "bin", StringValue("key121"), IntegerValue(11), CtxMapKey(StringValue("key1")), CtxMapRank(-1))
//  bin result = {key1:{key11:{key111:1},key12:{key121:11}}, key2:{key21:{"key211":7}}}

const (
	cdtMapOpTypeSetType                   = 64
	cdtMapOpTypeAdd                       = 65
	cdtMapOpTypeAddItems                  = 66
	cdtMapOpTypePut                       = 67
	cdtMapOpTypePutItems                  = 68
	cdtMapOpTypeReplace                   = 69
	cdtMapOpTypeReplaceItems              = 70
	cdtMapOpTypeIncrement                 = 73
	cdtMapOpTypeDecrement                 = 74
	cdtMapOpTypeClear                     = 75
	cdtMapOpTypeRemoveByKey               = 76
	cdtMapOpTypeRemoveByIndex             = 77
	cdtMapOpTypeRemoveByRank              = 79
	cdtMapOpTypeRemoveKeyList             = 81
	cdtMapOpTypeRemoveByValue             = 82
	cdtMapOpTypeRemoveValueList           = 83
	cdtMapOpTypeRemoveByKeyInterval       = 84
	cdtMapOpTypeRemoveByIndexRange        = 85
	cdtMapOpTypeRemoveByValueInterval     = 86
	cdtMapOpTypeRemoveByRankRange         = 87
	cdtMapOpTypeRemoveByKeyRelIndexRange  = 88
	cdtMapOpTypeRemoveByValueRelRankRange = 89
	cdtMapOpTypeSize                      = 96
	cdtMapOpTypeGetByKey                  = 97
	cdtMapOpTypeGetByIndex                = 98
	cdtMapOpTypeGetByRank                 = 100
	cdtMapOpTypeGetByValue                = 102
	cdtMapOpTypeGetByKeyInterval          = 103
	cdtMapOpTypeGetByIndexRange           = 104
	cdtMapOpTypeGetByValueInterval        = 105
	cdtMapOpTypeGetByRankRange            = 106
	cdtMapOpTypeGetByKeyList              = 107
	cdtMapOpTypeGetByValueList            = 108
	cdtMapOpTypeGetByKeyRelIndexRange     = 109
	cdtMapOpTypeGetByValueRelRankRange    = 110
)

type mapOrderType struct {
	attr int
	flag int
}

// MapOrder defines map storage order.
var MapOrder = struct {
	// Map is not ordered. This is the default.
	UNORDERED mapOrderType // 0

	// Order map by key.
	KEY_ORDERED mapOrderType // 1

	// Order map by key, then value.
	KEY_VALUE_ORDERED mapOrderType // 3
}{mapOrderType{0, 0x40}, mapOrderType{1, 0x80}, mapOrderType{3, 0xc0}}

type mapReturnType int

// MapReturnType defines the map return type.
// Type of data to return when selecting or removing items from the map.
var MapReturnType = struct {
	// NONE will will not return a result.
	NONE mapReturnType

	// INDEX will return key index order.
	//
	// 0 = first key
	// N = Nth key
	// -1 = last key
	INDEX mapReturnType

	// REVERSE_INDEX will return reverse key order.
	//
	// 0 = last key
	// -1 = first key
	REVERSE_INDEX mapReturnType

	// RANK will return value order.
	//
	// 0 = smallest value
	// N = Nth smallest value
	// -1 = largest value
	RANK mapReturnType

	// REVERSE_RANK will return reverse value order.
	//
	// 0 = largest value
	// N = Nth largest value
	// -1 = smallest value
	REVERSE_RANK mapReturnType

	// COUNT will return count of items selected.
	COUNT mapReturnType

	// KEY will return key for single key read and key list for range read.
	KEY mapReturnType

	// VALUE will return value for single key read and value list for range read.
	VALUE mapReturnType

	// KEY_VALUE will return key/value items. The possible return types are:
	//
	// map[interface{}]interface{} : Returned for unordered maps
	// []MapPair : Returned for range results where range order needs to be preserved.
	KEY_VALUE mapReturnType

	// INVERTED will invert meaning of map command and return values.  For example:
	// MapRemoveByKeyRange(binName, keyBegin, keyEnd, MapReturnType.KEY | MapReturnType.INVERTED)
	// With the INVERTED flag enabled, the keys outside of the specified key range will be removed and returned.
	INVERTED mapReturnType
}{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 0x10000,
}

// Unique key map write type.
type mapWriteMode struct {
	itemCommand  int
	itemsCommand int
}

// MapWriteMode should only be used for server versions < 4.3.
// MapWriteFlags are recommended for server versions >= 4.3.
var MapWriteMode = struct {
	// If the key already exists, the item will be overwritten.
	// If the key does not exist, a new item will be created.
	UPDATE *mapWriteMode

	// If the key already exists, the item will be overwritten.
	// If the key does not exist, the write will fail.
	UPDATE_ONLY *mapWriteMode

	// If the key already exists, the write will fail.
	// If the key does not exist, a new item will be created.
	CREATE_ONLY *mapWriteMode
}{
	&mapWriteMode{cdtMapOpTypePut, cdtMapOpTypePutItems},
	&mapWriteMode{cdtMapOpTypeReplace, cdtMapOpTypeReplaceItems},
	&mapWriteMode{cdtMapOpTypeAdd, cdtMapOpTypeAddItems},
}

/**
 * Map write bit flags.
 * Requires server versions >= 4.3.
 */
const (
	// MapWriteFlagsDefault is the Default. Allow create or update.
	MapWriteFlagsDefault = 0

	// MapWriteFlagsCreateOnly means: If the key already exists, the item will be denied.
	// If the key does not exist, a new item will be created.
	MapWriteFlagsCreateOnly = 1

	// MapWriteFlagsUpdateOnly means: If the key already exists, the item will be overwritten.
	// If the key does not exist, the item will be denied.
	MapWriteFlagsUpdateOnly = 2

	// MapWriteFlagsNoFail means: Do not raise error if a map item is denied due to write flag constraints.
	MapWriteFlagsNoFail = 4

	// MapWriteFlagsNoFail means: Allow other valid map items to be committed if a map item is denied due to
	// write flag constraints.
	MapWriteFlagsPartial = 8
)

// MapPolicy directives when creating a map and writing map items.
type MapPolicy struct {
	attributes   mapOrderType
	flags        int
	itemCommand  int
	itemsCommand int
}

// NewMapPolicy creates a MapPolicy with WriteMode. Use with servers before v4.3.
func NewMapPolicy(order mapOrderType, writeMode *mapWriteMode) *MapPolicy {
	return &MapPolicy{
		attributes:   order,
		flags:        MapWriteFlagsDefault,
		itemCommand:  writeMode.itemCommand,
		itemsCommand: writeMode.itemsCommand,
	}
}

// NewMapPolicyWithFlags creates a MapPolicy with WriteFlags. Use with servers v4.3+.
// Flags are MapWriteFlags. You can specify multiple flags by 'or'ing them together.
func NewMapPolicyWithFlags(order mapOrderType, flags int) *MapPolicy {
	return &MapPolicy{
		attributes:   order,
		flags:        flags,
		itemCommand:  MapWriteMode.UPDATE.itemCommand,
		itemsCommand: MapWriteMode.UPDATE.itemsCommand,
	}
}

// DefaultMapPolicy returns the default map policy
func DefaultMapPolicy() *MapPolicy {
	return NewMapPolicy(MapOrder.UNORDERED, MapWriteMode.UPDATE)
}

func newMapSetPolicyEncoder(op *Operation, packer BufferEx) (int, Error) {
	return packCDTParamsAsArray(packer, cdtMapOpTypeSetType, op.ctx, op.binValue.(IntegerValue))
}

func newMapSetPolicy(binName string, attributes mapOrderType, ctx []*CDTContext) *Operation {
	return &Operation{
		opType:   _MAP_MODIFY,
		binName:  binName,
		binValue: IntegerValue(attributes.attr),
		ctx:      ctx,
		encoder:  newMapSetPolicyEncoder,
	}
}

func newMapCreatePutEncoder(op *Operation, packer BufferEx) (int, Error) {
	return packCDTIfcParamsAsArray(packer, int16(*op.opSubType), op.ctx, op.binValue.(ListValue))
}

/////////////////////////

// MapCreateOp creates a map create operation.
// Server creates map at given context level.
func MapCreateOp(binName string, order mapOrderType, ctx []*CDTContext) *Operation {
	// If context not defined, the set order for top-level bin map.
	if len(ctx) == 0 {
		return MapSetPolicyOp(NewMapPolicyWithFlags(order, 0), binName)
	}

	return &Operation{
		opType:   _MAP_MODIFY,
		binName:  binName,
		binValue: ListValue([]interface{}{cdtMapOpTypeSetType, order.flag, IntegerValue(order.attr)}),
		ctx:      ctx,
		encoder:  cdtCreateOpEncoder,
	}
}

// MapSetPolicyOp creates set map policy operation.
// Server sets map policy attributes.  Server returns null.
//
// The required map policy attributes can be changed after the map is created.
func MapSetPolicyOp(policy *MapPolicy, binName string, ctx ...*CDTContext) *Operation {
	return newMapSetPolicy(binName, policy.attributes, ctx)
}

// MapPutOp creates map put operation.
// Server writes key/value item to map bin and returns map size.
//
// The required map policy dictates the type of map to create when it does not exist.
// The map policy also specifies the mode used when writing items to the map.
func MapPutOp(policy *MapPolicy, binName string, key interface{}, value interface{}, ctx ...*CDTContext) *Operation {
	if policy.flags != 0 {
		ops := cdtMapOpTypePut

		// Replace doesn't allow map attributes because it does not create on non-existing key.
		return &Operation{
			opType:    _MAP_MODIFY,
			opSubType: &ops,
			ctx:       ctx,
			binName:   binName,
			binValue:  ListValue([]interface{}{key, value, IntegerValue(policy.attributes.attr), IntegerValue(policy.flags)}),
			encoder:   newMapCreatePutEncoder,
		}
	}

	if policy.itemCommand == cdtMapOpTypeReplace {
		// Replace doesn't allow map attributes because it does not create on non-existing key.
		return &Operation{
			opType:    _MAP_MODIFY,
			opSubType: &policy.itemCommand,
			ctx:       ctx,
			binName:   binName,
			binValue:  ListValue([]interface{}{key, value}),
			encoder:   newMapCreatePutEncoder,
		}
	}

	return &Operation{
		opType:    _MAP_MODIFY,
		opSubType: &policy.itemCommand,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{key, value, IntegerValue(policy.attributes.attr)}),
		encoder:   newMapCreatePutEncoder,
	}
}

// MapPutItemsOp creates map put items operation
// Server writes each map item to map bin and returns map size.
//
// The required map policy dictates the type of map to create when it does not exist.
// The map policy also specifies the mode used when writing items to the map.
func MapPutItemsOp(policy *MapPolicy, binName string, amap map[interface{}]interface{}, ctx ...*CDTContext) *Operation {
	if policy.flags != 0 {
		ops := cdtMapOpTypePutItems

		// Replace doesn't allow map attributes because it does not create on non-existing key.
		return &Operation{
			opType:    _MAP_MODIFY,
			opSubType: &ops,
			ctx:       ctx,
			binName:   binName,
			binValue:  ListValue([]interface{}{amap, IntegerValue(policy.attributes.attr), IntegerValue(policy.flags)}),
			encoder:   newCDTCreateOperationEncoder,
		}
	}

	if policy.itemsCommand == int(cdtMapOpTypeReplaceItems) {
		// Replace doesn't allow map attributes because it does not create on non-existing key.
		return &Operation{
			opType:    _MAP_MODIFY,
			opSubType: &policy.itemsCommand,
			ctx:       ctx,
			binName:   binName,
			binValue:  ListValue([]interface{}{MapValue(amap)}),
			encoder:   newCDTCreateOperationEncoder,
		}
	}

	return &Operation{
		opType:    _MAP_MODIFY,
		opSubType: &policy.itemsCommand,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{MapValue(amap), IntegerValue(policy.attributes.attr)}),
		encoder:   newCDTCreateOperationEncoder,
	}
}

// MapIncrementOp creates map increment operation.
// Server increments values by incr for all items identified by key and returns final result.
// Valid only for numbers.
//
// The required map policy dictates the type of map to create when it does not exist.
// The map policy also specifies the mode used when writing items to the map.
func MapIncrementOp(policy *MapPolicy, binName string, key interface{}, incr interface{}, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValues2(cdtMapOpTypeIncrement, policy.attributes, binName, ctx, key, incr)
}

// MapDecrementOp creates map decrement operation.
// Server decrements values by decr for all items identified by key and returns final result.
// Valid only for numbers.
//
// The required map policy dictates the type of map to create when it does not exist.
// The map policy also specifies the mode used when writing items to the map.
func MapDecrementOp(policy *MapPolicy, binName string, key interface{}, decr interface{}, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValues2(cdtMapOpTypeDecrement, policy.attributes, binName, ctx, key, decr)
}

// MapClearOp creates map clear operation.
// Server removes all items in map.  Server returns null.
func MapClearOp(binName string, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValues0(cdtMapOpTypeClear, _MAP_MODIFY, binName, ctx)
}

// MapRemoveByKeyOp creates map remove operation.
// Server removes map item identified by key and returns removed data specified by returnType.
func MapRemoveByKeyOp(binName string, key interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeRemoveByKey, _MAP_MODIFY, binName, ctx, key, returnType)
}

// MapRemoveByKeyListOp creates map remove operation.
// Server removes map items identified by keys and returns removed data specified by returnType.
func MapRemoveByKeyListOp(binName string, keys []interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeRemoveKeyList, _MAP_MODIFY, binName, ctx, keys, returnType)
}

// MapRemoveByKeyRangeOp creates map remove operation.
// Server removes map items identified by key range (keyBegin inclusive, keyEnd exclusive).
// If keyBegin is null, the range is less than keyEnd.
// If keyEnd is null, the range is greater than equal to keyBegin.
//
// Server returns removed data specified by returnType.
func MapRemoveByKeyRangeOp(binName string, keyBegin interface{}, keyEnd interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateRangeOperation(cdtMapOpTypeRemoveByKeyInterval, _MAP_MODIFY, binName, ctx, keyBegin, keyEnd, returnType)
}

// MapRemoveByValueOp creates map remove operation.
// Server removes map items identified by value and returns removed data specified by returnType.
func MapRemoveByValueOp(binName string, value interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeRemoveByValue, _MAP_MODIFY, binName, ctx, value, returnType)
}

// MapRemoveByValueListOp creates map remove operation.
// Server removes map items identified by values and returns removed data specified by returnType.
func MapRemoveByValueListOp(binName string, values []interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValuesN(cdtMapOpTypeRemoveValueList, _MAP_MODIFY, binName, ctx, values, returnType)
}

// MapRemoveByValueRangeOp creates map remove operation.
// Server removes map items identified by value range (valueBegin inclusive, valueEnd exclusive).
// If valueBegin is null, the range is less than valueEnd.
// If valueEnd is null, the range is greater than equal to valueBegin.
//
// Server returns removed data specified by returnType.
func MapRemoveByValueRangeOp(binName string, valueBegin interface{}, valueEnd interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateRangeOperation(cdtMapOpTypeRemoveByValueInterval, _MAP_MODIFY, binName, ctx, valueBegin, valueEnd, returnType)
}

// MapRemoveByValueRelativeRankRangeOp creates a map remove by value relative to rank range operation.
// Server removes map items nearest to value and greater by relative rank.
// Server returns removed data specified by returnType.
//
// Examples for map [{4=2},{9=10},{5=15},{0=17}]:
//
//  (value,rank) = [removed items]
//  (11,1) = [{0=17}]
//  (11,-1) = [{9=10},{5=15},{0=17}]
func MapRemoveByValueRelativeRankRangeOp(binName string, value interface{}, rank int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateRangeOperation(cdtMapOpTypeRemoveByValueRelRankRange, _MAP_MODIFY, binName, ctx, value, rank, returnType)
}

// MapRemoveByValueRelativeRankRangeCountOp creates a map remove by value relative to rank range operation.
// Server removes map items nearest to value and greater by relative rank with a count limit.
// Server returns removed data specified by returnType (See MapReturnType).
//
// Examples for map [{4=2},{9=10},{5=15},{0=17}]:
//
//  (value,rank,count) = [removed items]
//  (11,1,1) = [{0=17}]
//  (11,-1,1) = [{9=10}]
func MapRemoveByValueRelativeRankRangeCountOp(binName string, value interface{}, rank, count int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTMapCreateOperationRelativeIndexCount(cdtMapOpTypeRemoveByValueRelRankRange, _MAP_MODIFY, binName, ctx, NewValue(value), rank, count, returnType)
}

// MapRemoveByIndexOp creates map remove operation.
// Server removes map item identified by index and returns removed data specified by returnType.
func MapRemoveByIndexOp(binName string, index int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeRemoveByIndex, _MAP_MODIFY, binName, ctx, index, returnType)
}

// MapRemoveByIndexRangeOp creates map remove operation.
// Server removes map items starting at specified index to the end of map and returns removed
// data specified by returnType.
func MapRemoveByIndexRangeOp(binName string, index int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeRemoveByIndexRange, _MAP_MODIFY, binName, ctx, index, returnType)
}

// MapRemoveByIndexRangeCountOp creates map remove operation.
// Server removes "count" map items starting at specified index and returns removed data specified by returnType.
func MapRemoveByIndexRangeCountOp(binName string, index int, count int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationIndexCount(cdtMapOpTypeRemoveByIndexRange, _MAP_MODIFY, binName, ctx, index, count, returnType)
}

// MapRemoveByRankOp creates map remove operation.
// Server removes map item identified by rank and returns removed data specified by returnType.
func MapRemoveByRankOp(binName string, rank int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeRemoveByRank, _MAP_MODIFY, binName, ctx, rank, returnType)
}

// MapRemoveByRankRangeOp creates map remove operation.
// Server removes map items starting at specified rank to the last ranked item and returns removed
// data specified by returnType.
func MapRemoveByRankRangeOp(binName string, rank int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationIndex(cdtMapOpTypeRemoveByRankRange, _MAP_MODIFY, binName, ctx, rank, returnType)
}

// MapRemoveByRankRangeCountOp creates map remove operation.
// Server removes "count" map items starting at specified rank and returns removed data specified by returnType.
func MapRemoveByRankRangeCountOp(binName string, rank int, count int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationIndexCount(cdtMapOpTypeRemoveByRankRange, _MAP_MODIFY, binName, ctx, rank, count, returnType)
}

// MapRemoveByKeyRelativeIndexRangeOp creates a map remove by key relative to index range operation.
// Server removes map items nearest to key and greater by index.
// Server returns removed data specified by returnType.
//
// Examples for map [{0=17},{4=2},{5=15},{9=10}]:
//
//  (value,index) = [removed items]
//  (5,0) = [{5=15},{9=10}]
//  (5,1) = [{9=10}]
//  (5,-1) = [{4=2},{5=15},{9=10}]
//  (3,2) = [{9=10}]
//  (3,-2) = [{0=17},{4=2},{5=15},{9=10}]
func MapRemoveByKeyRelativeIndexRangeOp(binName string, key interface{}, index int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTMapCreateOperationRelativeIndex(cdtMapOpTypeRemoveByKeyRelIndexRange, _MAP_MODIFY, binName, ctx, NewValue(key), index, returnType)
}

// MapRemoveByKeyRelativeIndexRangeCountOp creates map remove by key relative to index range operation.
// Server removes map items nearest to key and greater by index with a count limit.
// Server returns removed data specified by returnType.
//
// Examples for map [{0=17},{4=2},{5=15},{9=10}]:
//
//  (value,index,count) = [removed items]
//  (5,0,1) = [{5=15}]
//  (5,1,2) = [{9=10}]
//  (5,-1,1) = [{4=2}]
//  (3,2,1) = [{9=10}]
//  (3,-2,2) = [{0=17}]
func MapRemoveByKeyRelativeIndexRangeCountOp(binName string, key interface{}, index, count int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTMapCreateOperationRelativeIndexCount(cdtMapOpTypeRemoveByKeyRelIndexRange, _MAP_MODIFY, binName, ctx, NewValue(key), index, count, returnType)
}

// MapSizeOp creates map size operation.
// Server returns size of map.
func MapSizeOp(binName string, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValues0(cdtMapOpTypeSize, _MAP_READ, binName, ctx)
}

// MapGetByKeyOp creates map get by key operation.
// Server selects map item identified by key and returns selected data specified by returnType.
func MapGetByKeyOp(binName string, key interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeGetByKey, _MAP_READ, binName, ctx, key, returnType)
}

// MapGetByKeyRangeOp creates map get by key range operation.
// Server selects map items identified by key range (keyBegin inclusive, keyEnd exclusive).
// If keyBegin is null, the range is less than keyEnd.
// If keyEnd is null, the range is greater than equal to keyBegin.
//
// Server returns selected data specified by returnType.
func MapGetByKeyRangeOp(binName string, keyBegin interface{}, keyEnd interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateRangeOperation(cdtMapOpTypeGetByKeyInterval, _MAP_READ, binName, ctx, keyBegin, keyEnd, returnType)
}

// MapGetByKeyRelativeIndexRangeOp creates a map get by key relative to index range operation.
// Server selects map items nearest to key and greater by index.
// Server returns selected data specified by returnType.
//
// Examples for ordered map [{0=17},{4=2},{5=15},{9=10}]:
//
//  (value,index) = [selected items]
//  (5,0) = [{5=15},{9=10}]
//  (5,1) = [{9=10}]
//  (5,-1) = [{4=2},{5=15},{9=10}]
//  (3,2) = [{9=10}]
//  (3,-2) = [{0=17},{4=2},{5=15},{9=10}]
func MapGetByKeyRelativeIndexRangeOp(binName string, key interface{}, index int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTMapCreateOperationRelativeIndex(cdtMapOpTypeGetByKeyRelIndexRange, _MAP_READ, binName, ctx, NewValue(key), index, returnType)
}

// MapGetByKeyRelativeIndexRangeCountOp creates a map get by key relative to index range operation.
// Server selects map items nearest to key and greater by index with a count limit.
// Server returns selected data specified by returnType (See MapReturnType).
//
// Examples for ordered map [{0=17},{4=2},{5=15},{9=10}]:
//  (value,index,count) = [selected items]
//  (5,0,1) = [{5=15}]
//  (5,1,2) = [{9=10}]
//  (5,-1,1) = [{4=2}]
//  (3,2,1) = [{9=10}]
//  (3,-2,2) = [{0=17}]
func MapGetByKeyRelativeIndexRangeCountOp(binName string, key interface{}, index, count int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTMapCreateOperationRelativeIndexCount(cdtMapOpTypeGetByKeyRelIndexRange, _MAP_READ, binName, ctx, NewValue(key), index, count, returnType)
}

// MapGetByKeyListOp creates a map get by key list operation.
// Server selects map items identified by keys and returns selected data specified by returnType.
func MapGetByKeyListOp(binName string, keys []interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeGetByKeyList, _MAP_READ, binName, ctx, keys, returnType)
}

// MapGetByValueOp creates map get by value operation.
// Server selects map items identified by value and returns selected data specified by returnType.
func MapGetByValueOp(binName string, value interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeGetByValue, _MAP_READ, binName, ctx, value, returnType)
}

// MapGetByValueRangeOp creates map get by value range operation.
// Server selects map items identified by value range (valueBegin inclusive, valueEnd exclusive)
// If valueBegin is null, the range is less than valueEnd.
// If valueEnd is null, the range is greater than equal to valueBegin.
//
// Server returns selected data specified by returnType.
func MapGetByValueRangeOp(binName string, valueBegin interface{}, valueEnd interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateRangeOperation(cdtMapOpTypeGetByValueInterval, _MAP_READ, binName, ctx, valueBegin, valueEnd, returnType)
}

// MapGetByValueRelativeRankRangeOp creates a map get by value relative to rank range operation.
// Server selects map items nearest to value and greater by relative rank.
// Server returns selected data specified by returnType.
//
// Examples for map [{4=2},{9=10},{5=15},{0=17}]:
//
//  (value,rank) = [selected items]
//  (11,1) = [{0=17}]
//  (11,-1) = [{9=10},{5=15},{0=17}]
func MapGetByValueRelativeRankRangeOp(binName string, value interface{}, rank int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTMapCreateOperationRelativeIndex(cdtMapOpTypeGetByValueRelRankRange, _MAP_READ, binName, ctx, NewValue(value), rank, returnType)
}

// MapGetByValueRelativeRankRangeCountOp creates a map get by value relative to rank range operation.
// Server selects map items nearest to value and greater by relative rank with a count limit.
// Server returns selected data specified by returnType.
//
// Examples for map [{4=2},{9=10},{5=15},{0=17}]:
//
//  (value,rank,count) = [selected items]
//  (11,1,1) = [{0=17}]
//  (11,-1,1) = [{9=10}]
func MapGetByValueRelativeRankRangeCountOp(binName string, value interface{}, rank, count int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTMapCreateOperationRelativeIndexCount(cdtMapOpTypeGetByValueRelRankRange, _MAP_READ, binName, ctx, NewValue(value), rank, count, returnType)
}

// MapGetByValueListOp creates map get by value list operation.
// Server selects map items identified by values and returns selected data specified by returnType.
func MapGetByValueListOp(binName string, values []interface{}, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeGetByValueList, _MAP_READ, binName, ctx, values, returnType)
}

// MapGetByIndexOp creates map get by index operation.
// Server selects map item identified by index and returns selected data specified by returnType.
func MapGetByIndexOp(binName string, index int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeGetByIndex, _MAP_READ, binName, ctx, index, returnType)
}

// MapGetByIndexRangeOp creates map get by index range operation.
// Server selects map items starting at specified index to the end of map and returns selected
// data specified by returnType.
func MapGetByIndexRangeOp(binName string, index int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeGetByIndexRange, _MAP_READ, binName, ctx, index, returnType)
}

// MapGetByIndexRangeCountOp creates map get by index range operation.
// Server selects "count" map items starting at specified index and returns selected data specified by returnType.
func MapGetByIndexRangeCountOp(binName string, index int, count int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationIndexCount(cdtMapOpTypeGetByIndexRange, _MAP_READ, binName, ctx, index, count, returnType)
}

// MapGetByRankOp creates map get by rank operation.
// Server selects map item identified by rank and returns selected data specified by returnType.
func MapGetByRankOp(binName string, rank int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeGetByRank, _MAP_READ, binName, ctx, rank, returnType)
}

// MapGetByRankRangeOp creates map get by rank range operation.
// Server selects map items starting at specified rank to the last ranked item and returns selected
// data specified by returnType.
func MapGetByRankRangeOp(binName string, rank int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationValue1(cdtMapOpTypeGetByRankRange, _MAP_READ, binName, ctx, rank, returnType)
}

// MapGetByRankRangeCountOp creates map get by rank range operation.
// Server selects "count" map items starting at specified rank and returns selected data specified by returnType.
func MapGetByRankRangeCountOp(binName string, rank int, count int, returnType mapReturnType, ctx ...*CDTContext) *Operation {
	return newCDTCreateOperationIndexCount(cdtMapOpTypeGetByRankRange, _MAP_READ, binName, ctx, rank, count, returnType)
}
