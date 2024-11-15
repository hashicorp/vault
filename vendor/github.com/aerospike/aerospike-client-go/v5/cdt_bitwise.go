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

// Bit operations. Create bit operations used by client operate command.
// Offset orientation is left-to-right.  Negative offsets are supported.
// If the offset is negative, the offset starts backwards from end of the bitmap.
// If an offset is out of bounds, a parameter error will be returned.
//
//  Nested CDT operations are supported by optional CTX context arguments.  Example:
//  bin = [[0b00000001, 0b01000010],[0b01011010]]
//  Resize first bitmap (in a list of bitmaps) to 3 bytes.
//  BitOperation.resize("bin", 3, BitResizeFlags.DEFAULT, CTX.listIndex(0))
//  bin result = [[0b00000001, 0b01000010, 0b00000000],[0b01011010]]
const (
	_CDT_BITWISE_RESIZE   = 0
	_CDT_BITWISE_INSERT   = 1
	_CDT_BITWISE_REMOVE   = 2
	_CDT_BITWISE_SET      = 3
	_CDT_BITWISE_OR       = 4
	_CDT_BITWISE_XOR      = 5
	_CDT_BITWISE_AND      = 6
	_CDT_BITWISE_NOT      = 7
	_CDT_BITWISE_LSHIFT   = 8
	_CDT_BITWISE_RSHIFT   = 9
	_CDT_BITWISE_ADD      = 10
	_CDT_BITWISE_SUBTRACT = 11
	_CDT_BITWISE_SET_INT  = 12
	_CDT_BITWISE_GET      = 50
	_CDT_BITWISE_COUNT    = 51
	_CDT_BITWISE_LSCAN    = 52
	_CDT_BITWISE_RSCAN    = 53
	_CDT_BITWISE_GET_INT  = 54

	_CDT_BITWISE_INT_FLAGS_SIGNED = 1
)

// BitResizeOp creates byte "resize" operation.
// Server resizes []byte to byteSize according to resizeFlags (See BitResizeFlags).
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010]
//  byteSize = 4
//  resizeFlags = 0
//  bin result = [0b00000001, 0b01000010, 0b00000000, 0b00000000]
func BitResizeOp(policy *BitPolicy, binName string, byteSize int, resizeFlags BitResizeFlags, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_RESIZE, IntegerValue(byteSize), IntegerValue(policy.flags), IntegerValue(resizeFlags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitInsertOp creates byte "insert" operation.
// Server inserts value bytes into []byte bin at byteOffset.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  byteOffset = 1
//  value = [0b11111111, 0b11000111]
//  bin result = [0b00000001, 0b11111111, 0b11000111, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
func BitInsertOp(policy *BitPolicy, binName string, byteOffset int, value []byte, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_INSERT, IntegerValue(byteOffset), BytesValue(value), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitRemoveOp creates byte "remove" operation.
// Server removes bytes from []byte bin at byteOffset for byteSize.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  byteOffset = 2
//  byteSize = 3
//  bin result = [0b00000001, 0b01000010]
func BitRemoveOp(policy *BitPolicy, binName string, byteOffset int, byteSize int, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_REMOVE, IntegerValue(byteOffset), IntegerValue(byteSize), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitSetOp creates bit "set" operation.
// Server sets value on []byte bin at bitOffset for bitSize.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 13
//  bitSize = 3
//  value = [0b11100000]
//  bin result = [0b00000001, 0b01000111, 0b00000011, 0b00000100, 0b00000101]
func BitSetOp(policy *BitPolicy, binName string, bitOffset int, bitSize int, value []byte, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_SET, IntegerValue(bitOffset), IntegerValue(bitSize), BytesValue(value), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitOrOp creates bit "or" operation.
// Server performs bitwise "or" on value and []byte bin at bitOffset for bitSize.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 17
//  bitSize = 6
//  value = [0b10101000]
//  bin result = [0b00000001, 0b01000010, 0b01010111, 0b00000100, 0b00000101]
func BitOrOp(policy *BitPolicy, binName string, bitOffset int, bitSize int, value []byte, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_OR, IntegerValue(bitOffset), IntegerValue(bitSize), BytesValue(value), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitXorOp creates bit "exclusive or" operation.
// Server performs bitwise "xor" on value and []byte bin at bitOffset for bitSize.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 17
//  bitSize = 6
//  value = [0b10101100]
//  bin result = [0b00000001, 0b01000010, 0b01010101, 0b00000100, 0b00000101]
func BitXorOp(policy *BitPolicy, binName string, bitOffset int, bitSize int, value []byte, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_XOR, IntegerValue(bitOffset), IntegerValue(bitSize), BytesValue(value), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitAndOp creates bit "and" operation.
// Server performs bitwise "and" on value and []byte bin at bitOffset for bitSize.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 23
//  bitSize = 9
//  value = [0b00111100, 0b10000000]
//  bin result = [0b00000001, 0b01000010, 0b00000010, 0b00000000, 0b00000101]
func BitAndOp(policy *BitPolicy, binName string, bitOffset int, bitSize int, value []byte, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_AND, IntegerValue(bitOffset), IntegerValue(bitSize), BytesValue(value), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitNotOp creates bit "not" operation.
// Server negates []byte bin starting at bitOffset for bitSize.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 25
//  bitSize = 6
//  bin result = [0b00000001, 0b01000010, 0b00000011, 0b01111010, 0b00000101]
func BitNotOp(policy *BitPolicy, binName string, bitOffset int, bitSize int, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_NOT, IntegerValue(bitOffset), IntegerValue(bitSize), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitLShiftOp creates bit "left shift" operation.
// Server shifts left []byte bin starting at bitOffset for bitSize.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 32
//  bitSize = 8
//  shift = 3
//  bin result = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00101000]
func BitLShiftOp(policy *BitPolicy, binName string, bitOffset int, bitSize int, shift int, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_LSHIFT, IntegerValue(bitOffset), IntegerValue(bitSize), IntegerValue(shift), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitRShiftOp creates bit "right shift" operation.
// Server shifts right []byte bin starting at bitOffset for bitSize.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 0
//  bitSize = 9
//  shift = 1
//  bin result = [0b00000000, 0b11000010, 0b00000011, 0b00000100, 0b00000101]
func BitRShiftOp(policy *BitPolicy, binName string, bitOffset int, bitSize int, shift int, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_RSHIFT, IntegerValue(bitOffset), IntegerValue(bitSize), IntegerValue(shift), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitAddOp creates bit "add" operation.
// Server adds value to []byte bin starting at bitOffset for bitSize. BitSize must be <= 64.
// Signed indicates if bits should be treated as a signed number.
// If add overflows/underflows, BitOverflowAction is used.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 24
//  bitSize = 16
//  value = 128
//  signed = false
//  bin result = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b10000101]
func BitAddOp(
	policy *BitPolicy,
	binName string,
	bitOffset int,
	bitSize int,
	value int64,
	signed bool,
	action BitOverflowAction,
	ctx ...*CDTContext,
) *Operation {
	// return createMathOperation(ADD, policy, binName, ctx, bitOffset, bitSize, value, signed, action)

	actionFlags := action
	if signed {
		actionFlags |= _CDT_BITWISE_INT_FLAGS_SIGNED
	}
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_ADD, IntegerValue(bitOffset), IntegerValue(bitSize), IntegerValue(value), IntegerValue(policy.flags), IntegerValue(actionFlags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitSubtractOp creates bit "subtract" operation.
// Server subtracts value from []byte bin starting at bitOffset for bitSize. BitSize must be <= 64.
// Signed indicates if bits should be treated as a signed number.
// If add overflows/underflows, BitOverflowAction is used.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 24
//  bitSize = 16
//  value = 128
//  signed = false
//  bin result = [0b00000001, 0b01000010, 0b00000011, 0b0000011, 0b10000101]
func BitSubtractOp(
	policy *BitPolicy,
	binName string,
	bitOffset int,
	bitSize int,
	value int64,
	signed bool,
	action BitOverflowAction,
	ctx ...*CDTContext,
) *Operation {
	// return createMathOperation(SUBTRACT, policy, binName, ctx, bitOffset, bitSize, value, signed, action)

	actionFlags := action
	if signed {
		actionFlags |= _CDT_BITWISE_INT_FLAGS_SIGNED
	}
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_SUBTRACT, IntegerValue(bitOffset), IntegerValue(bitSize), IntegerValue(value), IntegerValue(policy.flags), IntegerValue(actionFlags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitSetIntOp creates bit "setInt" operation.
// Server sets value to []byte bin starting at bitOffset for bitSize. Size must be <= 64.
// Server does not return a value.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 1
//  bitSize = 8
//  value = 127
//  bin result = [0b00111111, 0b11000010, 0b00000011, 0b0000100, 0b00000101]
func BitSetIntOp(policy *BitPolicy, binName string, bitOffset int, bitSize int, value int64, ctx ...*CDTContext) *Operation {
	// Packer packer = new Packer();
	// init(packer, ctx, SET_INT, 4)
	// packer.packInt(bitOffset)
	// packer.packInt(bitSize)
	// packer.packLong(value)
	// packer.packInt(policy.flags)
	// return newOperation(_BIT_MODIFY, binName, Value.get(packer.toByteArray()))
	return &Operation{
		opType:   _BIT_MODIFY,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_SET_INT, IntegerValue(bitOffset), IntegerValue(bitSize), IntegerValue(value), IntegerValue(policy.flags)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitGetOp creates bit "get" operation.
// Server returns bits from []byte bin starting at bitOffset for bitSize.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 9
//  bitSize = 5
//  returns [0b1000000]
func BitGetOp(binName string, bitOffset int, bitSize int, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_READ,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_GET, IntegerValue(bitOffset), IntegerValue(bitSize)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitCountOp creates bit "count" operation.
// Server returns integer count of set bits from []byte bin starting at bitOffset for bitSize.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 20
//  bitSize = 4
//  returns 2
func BitCountOp(binName string, bitOffset int, bitSize int, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_READ,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_COUNT, IntegerValue(bitOffset), IntegerValue(bitSize)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitLScanOp creates bit "left scan" operation.
// Server returns integer bit offset of the first specified value bit in []byte bin
// starting at bitOffset for bitSize.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 24
//  bitSize = 8
//  value = true
//  returns 5
func BitLScanOp(binName string, bitOffset int, bitSize int, value bool, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_READ,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_LSCAN, IntegerValue(bitOffset), IntegerValue(bitSize), BoolValue(value)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitRScanOp creates bit "right scan" operation.
// Server returns integer bit offset of the last specified value bit in []byte bin
// starting at bitOffset for bitSize.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 32
//  bitSize = 8
//  value = true
//  returns 7
func BitRScanOp(binName string, bitOffset int, bitSize int, value bool, ctx ...*CDTContext) *Operation {
	return &Operation{
		opType:   _BIT_READ,
		ctx:      ctx,
		binName:  binName,
		binValue: ListValue{_CDT_BITWISE_RSCAN, IntegerValue(bitOffset), IntegerValue(bitSize), BoolValue(value)},
		encoder:  newCDTBitwiseEncoder,
	}
}

// BitGetIntOp creates bit "get integer" operation.
// Server returns integer from []byte bin starting at bitOffset for bitSize.
// Signed indicates if bits should be treated as a signed number.
// Example:
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 8
//  bitSize = 16
//  signed = false
//  returns 16899
func BitGetIntOp(binName string, bitOffset int, bitSize int, signed bool, ctx ...*CDTContext) *Operation {
	binValue := ListValue{_CDT_BITWISE_GET_INT, IntegerValue(bitOffset), IntegerValue(bitSize)}
	if signed {
		binValue = append(binValue, IntegerValue(_CDT_BITWISE_INT_FLAGS_SIGNED))
	}
	return &Operation{
		opType:   _BIT_READ,
		ctx:      ctx,
		binName:  binName,
		binValue: binValue,
		encoder:  newCDTBitwiseEncoder,
	}
}

func newCDTBitwiseEncoder(op *Operation, packer BufferEx) (int, Error) {
	params := op.binValue.(ListValue)
	opType := params[0].(int)
	if len(op.binValue.(ListValue)) > 1 {
		return packCDTBitIfcParamsAsArray(packer, int16(opType), op.ctx, params[1:])
	}
	return packCDTBitIfcParamsAsArray(packer, int16(opType), op.ctx, nil)
}

func packCDTBitIfcParamsAsArray(packer BufferEx, opType int16, ctx []*CDTContext, params ListValue) (int, Error) {
	return packCDTBitIfcVarParamsAsArray(packer, opType, ctx, []interface{}(params)...)
}

func packCDTBitIfcVarParamsAsArray(packer BufferEx, opType int16, ctx []*CDTContext, params ...interface{}) (int, Error) {
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
	}

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
