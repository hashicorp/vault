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

// ExpReadFlags is used to change mode in expression reads.
type ExpReadFlags int

const (
	// ExpReadFlagDefault is the default
	ExpReadFlagDefault ExpReadFlags = 0

	// ExpReadFlagEvalNoFail means:
	// Ignore failures caused by the expression resolving to unknown or a non-bin type.
	ExpReadFlagEvalNoFail ExpReadFlags = 1 << 4
)

// ExpWriteFlags is used to change mode in expression writes.
type ExpWriteFlags int

// Expression write Flags
const (
	// ExpWriteFlagDefault is the default. Allows create or update.
	ExpWriteFlagDefault ExpWriteFlags = 0

	// ExpWriteFlagCreateOnly means:
	// If bin does not exist, a new bin will be created.
	// If bin exists, the operation will be denied.
	// If bin exists, fail with Bin Exists
	ExpWriteFlagCreateOnly ExpWriteFlags = 1 << 0

	// ExpWriteFlagUpdateOnly means:
	// If bin exists, the bin will be overwritten.
	// If bin does not exist, the operation will be denied.
	// If bin does not exist, fail with Bin Not Found
	ExpWriteFlagUpdateOnly ExpWriteFlags = 1 << 1

	// ExpWriteFlagAllowDelete means:
	// If expression results in nil value, then delete the bin.
	// Otherwise, return OP Not Applicable when NoFail is not set
	ExpWriteFlagAllowDelete ExpWriteFlags = 1 << 2

	// ExpWriteFlagPolicyNoFail means:
	// Do not raise error if operation is denied.
	ExpWriteFlagPolicyNoFail ExpWriteFlags = 1 << 3

	// ExpWriteFlagEvalNoFail means:
	// Ignore failures caused by the expression resolving to unknown or a non-bin type.
	ExpWriteFlagEvalNoFail ExpWriteFlags = 1 << 4
)

// ExpWriteOp creates an operation with an expression that writes to record bin.
func ExpWriteOp(binName string, exp *Expression, flags ExpWriteFlags) *Operation {
	val, err := encodeExpOperation(exp, int(flags))
	if err != nil {
		panic(err)
	}
	return &Operation{
		opType:   _EXP_MODIFY,
		binName:  binName,
		binValue: NewValue(val),
		encoder:  nil,
	}
}

// ExpReadOp creates an operation with an expression that reads from a record.
func ExpReadOp(name string, exp *Expression, flags ExpReadFlags) *Operation {
	val, err := encodeExpOperation(exp, int(flags))
	if err != nil {
		panic(err)
	}
	return &Operation{
		opType:   _EXP_READ,
		binName:  name,
		binValue: NewValue(val),
		encoder:  nil,
	}
}

// newExpOperationEncoder is used to encode the operation expression wire protocol
func encodeExpOperation(exp *Expression, flags int) ([]byte, Error) {
	// expression is double packed: first normally, and then as a BLOB in operation
	// find the size of packed expression and pack it in temp buffer
	tsz, err := exp.pack(nil)
	if err != nil {
		return nil, err
	}

	// header + packed bytes + flags = 16 + len bytes max
	buf := newBuffer(tsz + 16)

	if _, err := packArrayBegin(buf, 2); err != nil {
		return nil, err
	}

	if _, err = exp.pack(buf); err != nil {
		return nil, err
	}

	if _, err = packAInt(buf, flags); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
