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

func newCDTCreateOperationEncoder(op *Operation, packer BufferEx) (int, Error) {
	if op.binValue != nil {
		if params := op.binValue.(ListValue); len(params) > 0 {
			return packCDTIfcParamsAsArray(packer, int16(*op.opSubType), op.ctx, op.binValue.(ListValue))
		}
	}
	return packCDTParamsAsArray(packer, int16(*op.opSubType), op.ctx)
}

func newCDTCreateOperationValues2(command int, attributes mapOrderType, binName string, ctx []*CDTContext, value1 interface{}, value2 interface{}) *Operation {
	return &Operation{
		opType:    _MAP_MODIFY,
		opSubType: &command,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{value1, value2, IntegerValue(attributes.attr)}),
		encoder:   newCDTCreateOperationEncoder,
	}
}

func newCDTCreateOperationValues0(command int, typ OperationType, binName string, ctx []*CDTContext) *Operation {
	return &Operation{
		opType:    typ,
		opSubType: &command,
		ctx:       ctx,
		binName:   binName,
		// binValue: NewNullValue(),
		encoder: newCDTCreateOperationEncoder,
	}
}

func newCDTCreateOperationValuesN(command int, typ OperationType, binName string, ctx []*CDTContext, values []interface{}, returnType mapReturnType) *Operation {
	return &Operation{
		opType:    typ,
		opSubType: &command,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{IntegerValue(returnType), ListValue(values)}),
		encoder:   newCDTCreateOperationEncoder,
	}
}

func newCDTCreateOperationValue1(command int, typ OperationType, binName string, ctx []*CDTContext, value interface{}, returnType mapReturnType) *Operation {
	return &Operation{
		opType:    typ,
		opSubType: &command,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{IntegerValue(returnType), value}),
		encoder:   newCDTCreateOperationEncoder,
	}
}

func newCDTCreateOperationIndex(command int, typ OperationType, binName string, ctx []*CDTContext, index int, returnType mapReturnType) *Operation {
	return &Operation{
		opType:    typ,
		opSubType: &command,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{IntegerValue(returnType), index}),
		encoder:   newCDTCreateOperationEncoder,
	}
}

func newCDTCreateOperationIndexCount(command int, typ OperationType, binName string, ctx []*CDTContext, index int, count int, returnType mapReturnType) *Operation {
	return &Operation{
		opType:    typ,
		opSubType: &command,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{IntegerValue(returnType), index, count}),
		encoder:   newCDTCreateOperationEncoder,
	}
}

func newCDTMapCreateOperationRelativeIndex(command int, typ OperationType, binName string, ctx []*CDTContext, key Value, index int, returnType mapReturnType) *Operation {
	return &Operation{
		opType:    typ,
		opSubType: &command,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{IntegerValue(returnType), key, index}),
		encoder:   newCDTCreateOperationEncoder,
	}
}

func newCDTMapCreateOperationRelativeIndexCount(command int, typ OperationType, binName string, ctx []*CDTContext, key Value, index int, count int, returnType mapReturnType) *Operation {
	return &Operation{
		opType:    typ,
		opSubType: &command,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{IntegerValue(returnType), key, index, count}),
		encoder:   newCDTCreateOperationEncoder,
	}
}

func newCDTCreateRangeOperation(command int, typ OperationType, binName string, ctx []*CDTContext, begin interface{}, end interface{}, returnType mapReturnType) *Operation {
	if end == nil {
		return &Operation{
			opType:    typ,
			opSubType: &command,
			ctx:       ctx,
			binName:   binName,
			binValue:  ListValue([]interface{}{IntegerValue(returnType), begin}),
			encoder:   newCDTCreateOperationEncoder,
		}
	}

	return &Operation{
		opType:    typ,
		opSubType: &command,
		ctx:       ctx,
		binName:   binName,
		binValue:  ListValue([]interface{}{IntegerValue(returnType), begin, end}),
		encoder:   newCDTCreateOperationEncoder,
	}
}
