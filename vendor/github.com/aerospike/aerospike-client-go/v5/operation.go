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

// OperationType determines operation type
type OperationType struct {
	op      byte
	isWrite bool
}

type operationSubType *int

// Valid OperationType values that can be used to create custom Operations.
// The names are self-explanatory.
var (
	_READ = OperationType{1, false}
	// READ_HEADER OperationType = {1, false }

	_WRITE      = OperationType{2, true}
	_CDT_READ   = OperationType{3, false}
	_CDT_MODIFY = OperationType{4, true}
	_MAP_READ   = OperationType{3, false}
	_MAP_MODIFY = OperationType{4, true}
	_ADD        = OperationType{5, true}
	_EXP_READ   = OperationType{7, false}
	_EXP_MODIFY = OperationType{8, true}
	_APPEND     = OperationType{9, true}
	_PREPEND    = OperationType{10, true}
	_TOUCH      = OperationType{11, true}
	_BIT_READ   = OperationType{12, false}
	_BIT_MODIFY = OperationType{13, true}
	_DELETE     = OperationType{14, true}
	_HLL_READ   = OperationType{15, false}
	_HLL_MODIFY = OperationType{16, true}
)

// Operation contains operation definition.
// This struct is used in client's operate() method.
type Operation struct {

	// OpType determines type of operation.
	opType OperationType
	// used in CDT commands
	opSubType operationSubType
	// CDT context for nested types
	ctx []*CDTContext

	encoder func(*Operation, BufferEx) (int, Error)

	// binName (Optional) determines the name of bin used in operation.
	binName string

	// binValue (Optional) determines bin value used in operation.
	binValue Value

	// will be true ONLY for GetHeader() operation
	headerOnly bool

	// reused determines if the operation is cached. If so, it will cache the
	// internal bytes in binValue field and remove the encoder for maximum performance
	used bool
}

// cache uses the encoder and caches the packed operation for further use.
func (op *Operation) cache() Error {
	packer := newPacker()

	if _, err := op.encoder(op, packer); err != nil {
		return err
	}

	op.binValue = BytesValue(packer.Bytes())
	op.encoder = nil // do not encode anymore; just use the cache
	op.used = false  // do not encode anymore; just use the cache
	return nil
}

// GetBinOp creates read bin database operation.
func GetBinOp(binName string) *Operation {
	return &Operation{opType: _READ, binName: binName, binValue: NewNullValue()}
}

// GetOp creates read all record bins database operation.
func GetOp() *Operation {
	return &Operation{opType: _READ, binValue: NewNullValue()}
}

// GetHeaderOp creates read record header database operation.
func GetHeaderOp() *Operation {
	return &Operation{opType: _READ, headerOnly: true, binValue: NewNullValue()}
}

// PutOp creates set database operation.
func PutOp(bin *Bin) *Operation {
	return &Operation{opType: _WRITE, binName: bin.Name, binValue: bin.Value}
}

// AppendOp creates string append database operation.
func AppendOp(bin *Bin) *Operation {
	return &Operation{opType: _APPEND, binName: bin.Name, binValue: bin.Value}
}

// PrependOp creates string prepend database operation.
func PrependOp(bin *Bin) *Operation {
	return &Operation{opType: _PREPEND, binName: bin.Name, binValue: bin.Value}
}

// AddOp creates integer add database operation.
func AddOp(bin *Bin) *Operation {
	return &Operation{opType: _ADD, binName: bin.Name, binValue: bin.Value}
}

// TouchOp creates touch record database operation.
func TouchOp() *Operation {
	return &Operation{opType: _TOUCH, binValue: NewNullValue()}
}

// DeleteOp creates delete record database operation.
func DeleteOp() *Operation {
	return &Operation{opType: _DELETE, binValue: NewNullValue()}
}
