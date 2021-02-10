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

const bitwiseMODULE = 1
const bitwiseINT_FLAGS_SIGNED = 1

const (
	BitExpOpRESIZE   = 0
	BitExpOpINSERT   = 1
	BitExpOpREMOVE   = 2
	BitExpOpSET      = 3
	BitExpOpOR       = 4
	BitExpOpXOR      = 5
	BitExpOpAND      = 6
	BitExpOpNOT      = 7
	BitExpOpLSHIFT   = 8
	BitExpOpRSHIFT   = 9
	BitExpOpADD      = 10
	BitExpOpSUBTRACT = 11
	BitExpOpSETINT   = 12
	BitExpOpGET      = 50
	BitExpOpCOUNT    = 51
	BitExpOpLSCAN    = 52
	BitExpOpRSCAN    = 53
	BitExpOpGETINT   = 54
)

// Create expression that resizes byte[] to byteSize according to resizeFlags
// and returns byte[].
func ExpBitResize(
	policy *BitPolicy,
	byte_size *FilterExpression,
	resize_flags BitResizeFlags,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpRESIZE),
		byte_size,
		IntegerValue(policy.flags),
		IntegerValue(resize_flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that inserts value bytes into byte[] bin at byteOffset and returns byte[].
func ExpBitInsert(
	policy *BitPolicy,
	byte_offset *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpINSERT),
		byte_offset,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that removes bytes from byte[] bin at byteOffset for byteSize and returns byte[].
func ExpBitRemove(
	policy *BitPolicy,
	byte_offset *FilterExpression,
	byte_size *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpREMOVE),
		byte_offset,
		byte_size,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that sets value on byte[] bin at bitOffset for bitSize and returns byte[].
func ExpBitSet(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpSET),
		bit_offset,
		bit_size,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that performs bitwise "or" on value and byte[] bin at bitOffset for bitSize
// and returns byte[].
func ExpBitOr(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpOR),
		bit_offset,
		bit_size,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that performs bitwise "xor" on value and byte[] bin at bitOffset for bitSize
// and returns byte[].
func ExpBitXor(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpXOR),
		bit_offset,
		bit_size,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that performs bitwise "and" on value and byte[] bin at bitOffset for bitSize
// and returns byte[].
//
// bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
// bitOffset = 23
// bitSize = 9
// value = [0b00111100, 0b10000000]
// bin result = [0b00000001, 0b01000010, 0b00000010, 0b00000000, 0b00000101]
func ExpBitAnd(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpAND),
		bit_offset,
		bit_size,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that negates byte[] bin starting at bitOffset for bitSize and returns byte[].
func ExpBitNot(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpNOT),
		bit_offset,
		bit_size,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that shifts left byte[] bin starting at bitOffset for bitSize and returns byte[].
func ExpBitLShift(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	shift *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpLSHIFT),
		bit_offset,
		bit_size,
		shift,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that shifts right byte[] bin starting at bitOffset for bitSize and returns byte[].
func ExpBitRShift(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	shift *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpRSHIFT),
		bit_offset,
		bit_size,
		shift,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that adds value to byte[] bin starting at bitOffset for bitSize and returns byte[].
// `BitSize` must be <= 64. Signed indicates if bits should be treated as a signed number.
// If add overflows/underflows, `BitOverflowAction` is used.
func ExpBitAdd(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	value *FilterExpression,
	signed bool,
	action BitOverflowAction,
	bin *FilterExpression,
) *FilterExpression {
	flags := byte(action)
	if signed {
		flags |= bitwiseINT_FLAGS_SIGNED
	}
	args := []ExpressionArgument{
		IntegerValue(BitExpOpADD),
		bit_offset,
		bit_size,
		value,
		IntegerValue(policy.flags),
		IntegerValue(flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that subtracts value from byte[] bin starting at bitOffset for bitSize and returns byte[].
// `BitSize` must be <= 64. Signed indicates if bits should be treated as a signed number.
// If add overflows/underflows, `BitOverflowAction` is used.
func ExpBitSubtract(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	value *FilterExpression,
	signed bool,
	action BitOverflowAction,
	bin *FilterExpression,
) *FilterExpression {
	flags := byte(action)
	if signed {
		flags |= bitwiseINT_FLAGS_SIGNED
	}
	args := []ExpressionArgument{
		IntegerValue(BitExpOpSUBTRACT),
		bit_offset,
		bit_size,
		value,
		IntegerValue(policy.flags),
		IntegerValue(flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that sets value to byte[] bin starting at bitOffset for bitSize and returns byte[].
// `BitSize` must be <= 64.
func ExpBitSetInt(
	policy *BitPolicy,
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpSETINT),
		bit_offset,
		bit_size,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// Create expression that returns bits from byte[] bin starting at bitOffset for bitSize.
func ExpBitGet(
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpGET),
		bit_offset,
		bit_size,
	}

	return expBitAddRead(bin, ExpTypeBLOB, args)
}

// Create expression that returns integer count of set bits from byte[] bin starting at
// bitOffset for bitSize.
func ExpBitCount(
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpCOUNT),
		bit_offset,
		bit_size,
	}

	return expBitAddRead(bin, ExpTypeINT, args)
}

// Create expression that returns integer bit offset of the first specified value bit in byte[] bin
// starting at bitOffset for bitSize.
func ExpBitLScan(
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpLSCAN),
		bit_offset,
		bit_size,
		value,
	}

	return expBitAddRead(bin, ExpTypeINT, args)
}

// Create expression that returns integer bit offset of the last specified value bit in byte[] bin
// starting at bitOffset for bitSize.
func ExpBitRScan(
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	value *FilterExpression,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpRSCAN),
		bit_offset,
		bit_size,
		value,
	}

	return expBitAddRead(bin, ExpTypeINT, args)
}

// Create expression that returns integer from byte[] bin starting at bitOffset for bitSize.
// Signed indicates if bits should be treated as a signed number.
func ExpBitGetInt(
	bit_offset *FilterExpression,
	bit_size *FilterExpression,
	signed bool,
	bin *FilterExpression,
) *FilterExpression {
	args := []ExpressionArgument{
		IntegerValue(BitExpOpGETINT),
		bit_offset,
		bit_size,
	}
	if signed {
		args = append(args, IntegerValue(bitwiseINT_FLAGS_SIGNED))
	}

	return expBitAddRead(bin, ExpTypeINT, args)
}

func expBitAddWrite(bin *FilterExpression, arguments []ExpressionArgument) *FilterExpression {
	flags := int64(bitwiseMODULE | MODIFY)
	return &FilterExpression{
		cmd:       &expOpCALL,
		val:       nil,
		bin:       bin,
		flags:     &flags,
		module:    &ExpTypeBLOB,
		exps:      nil,
		arguments: arguments,
	}
}

func expBitAddRead(
	bin *FilterExpression,
	return_type ExpType,
	arguments []ExpressionArgument,
) *FilterExpression {
	flags := int64(bitwiseMODULE)
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
