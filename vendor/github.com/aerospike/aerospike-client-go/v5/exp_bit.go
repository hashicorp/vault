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
// limitations under the License.package aerospike

package aerospike

const bitwiseMODULE = 1
const bitwiseINT_FLAGS_SIGNED = 1

const (
	_BitExpOpRESIZE   = 0
	_BitExpOpINSERT   = 1
	_BitExpOpREMOVE   = 2
	_BitExpOpSET      = 3
	_BitExpOpOR       = 4
	_BitExpOpXOR      = 5
	_BitExpOpAND      = 6
	_BitExpOpNOT      = 7
	_BitExpOpLSHIFT   = 8
	_BitExpOpRSHIFT   = 9
	_BitExpOpADD      = 10
	_BitExpOpSUBTRACT = 11
	_BitExpOpSETINT   = 12
	_BitExpOpGET      = 50
	_BitExpOpCOUNT    = 51
	_BitExpOpLSCAN    = 52
	_BitExpOpRSCAN    = 53
	_BitExpOpGETINT   = 54
)

// ExpBitResize creates an expression that resizes []byte to byteSize according to resizeFlags
// and returns []byte.
func ExpBitResize(
	policy *BitPolicy,
	byteSize *Expression,
	resizeFlags BitResizeFlags,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpRESIZE),
		byteSize,
		IntegerValue(policy.flags),
		IntegerValue(resizeFlags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitInsert creates an expression that inserts value bytes into []byte bin at byteOffset and returns []byte.
func ExpBitInsert(
	policy *BitPolicy,
	byteOffset *Expression,
	value *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpINSERT),
		byteOffset,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitRemove creates an expression that removes bytes from []byte bin at byteOffset for byteSize and returns []byte.
func ExpBitRemove(
	policy *BitPolicy,
	byteOffset *Expression,
	byteSize *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpREMOVE),
		byteOffset,
		byteSize,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitSet creates an expression that sets value on []byte bin at bitOffset for bitSize and returns []byte.
func ExpBitSet(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	value *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpSET),
		bitOffset,
		bitSize,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitOr creates an expression that performs bitwise "or" on value and []byte bin at bitOffset for bitSize
// and returns []byte.
func ExpBitOr(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	value *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpOR),
		bitOffset,
		bitSize,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitXor creates an expression that performs bitwise "xor" on value and []byte bin at bitOffset for bitSize
// and returns []byte.
func ExpBitXor(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	value *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpXOR),
		bitOffset,
		bitSize,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitAnd creates an expression that performs bitwise "and" on value and []byte bin at bitOffset for bitSize
// and returns []byte.
//
//  bin = [0b00000001, 0b01000010, 0b00000011, 0b00000100, 0b00000101]
//  bitOffset = 23
//  bitSize = 9
//  value = [0b00111100, 0b10000000]
//  bin result = [0b00000001, 0b01000010, 0b00000010, 0b00000000, 0b00000101]
func ExpBitAnd(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	value *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpAND),
		bitOffset,
		bitSize,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitNot creates an expression that negates []byte bin starting at bitOffset for bitSize and returns []byte.
func ExpBitNot(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpNOT),
		bitOffset,
		bitSize,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitLShift creates an expression that shifts left []byte bin starting at bitOffset for bitSize and returns []byte.
func ExpBitLShift(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	shift *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpLSHIFT),
		bitOffset,
		bitSize,
		shift,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitRShift creates an expression that shifts right []byte bin starting at bitOffset for bitSize and returns []byte.
func ExpBitRShift(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	shift *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpRSHIFT),
		bitOffset,
		bitSize,
		shift,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitAdd creates an expression that adds value to []byte bin starting at bitOffset for bitSize and returns []byte.
// `BitSize` must be <= 64. Signed indicates if bits should be treated as a signed number.
// If add overflows/underflows, `BitOverflowAction` is used.
func ExpBitAdd(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	value *Expression,
	signed bool,
	action BitOverflowAction,
	bin *Expression,
) *Expression {
	flags := byte(action)
	if signed {
		flags |= bitwiseINT_FLAGS_SIGNED
	}
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpADD),
		bitOffset,
		bitSize,
		value,
		IntegerValue(policy.flags),
		IntegerValue(flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitSubtract creates an expression that subtracts value from []byte bin starting at bitOffset for bitSize and returns []byte.
// `BitSize` must be <= 64. Signed indicates if bits should be treated as a signed number.
// If add overflows/underflows, `BitOverflowAction` is used.
func ExpBitSubtract(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	value *Expression,
	signed bool,
	action BitOverflowAction,
	bin *Expression,
) *Expression {
	flags := byte(action)
	if signed {
		flags |= bitwiseINT_FLAGS_SIGNED
	}
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpSUBTRACT),
		bitOffset,
		bitSize,
		value,
		IntegerValue(policy.flags),
		IntegerValue(flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitSetInt creates an expression that sets value to []byte bin starting at bitOffset for bitSize and returns []byte.
// `BitSize` must be <= 64.
func ExpBitSetInt(
	policy *BitPolicy,
	bitOffset *Expression,
	bitSize *Expression,
	value *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpSETINT),
		bitOffset,
		bitSize,
		value,
		IntegerValue(policy.flags),
	}

	return expBitAddWrite(bin, args)
}

// ExpBitGet creates an expression that returns bits from []byte bin starting at bitOffset for bitSize.
func ExpBitGet(
	bitOffset *Expression,
	bitSize *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpGET),
		bitOffset,
		bitSize,
	}

	return expBitAddRead(bin, ExpTypeBLOB, args)
}

// ExpBitCount creates an expression that returns integer count of set bits from []byte bin starting at
// bitOffset for bitSize.
func ExpBitCount(
	bitOffset *Expression,
	bitSize *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpCOUNT),
		bitOffset,
		bitSize,
	}

	return expBitAddRead(bin, ExpTypeINT, args)
}

// ExpBitLScan creates an expression that returns integer bit offset of the first specified value bit in []byte bin
// starting at bitOffset for bitSize.
func ExpBitLScan(
	bitOffset *Expression,
	bitSize *Expression,
	value *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpLSCAN),
		bitOffset,
		bitSize,
		value,
	}

	return expBitAddRead(bin, ExpTypeINT, args)
}

// ExpBitRScan creates an expression that returns integer bit offset of the last specified value bit in []byte bin
// starting at bitOffset for bitSize.
func ExpBitRScan(
	bitOffset *Expression,
	bitSize *Expression,
	value *Expression,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpRSCAN),
		bitOffset,
		bitSize,
		value,
	}

	return expBitAddRead(bin, ExpTypeINT, args)
}

// ExpBitGetInt Create expression that returns integer from []byte bin starting at bitOffset for bitSize.
// Signed indicates if bits should be treated as a signed number.
func ExpBitGetInt(
	bitOffset *Expression,
	bitSize *Expression,
	signed bool,
	bin *Expression,
) *Expression {
	args := []ExpressionArgument{
		IntegerValue(_BitExpOpGETINT),
		bitOffset,
		bitSize,
	}
	if signed {
		args = append(args, IntegerValue(bitwiseINT_FLAGS_SIGNED))
	}

	return expBitAddRead(bin, ExpTypeINT, args)
}

func expBitAddWrite(bin *Expression, arguments []ExpressionArgument) *Expression {
	flags := int64(bitwiseMODULE | _MODIFY)
	return &Expression{
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
	bin *Expression,
	returnType ExpType,
	arguments []ExpressionArgument,
) *Expression {
	flags := int64(bitwiseMODULE)
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
