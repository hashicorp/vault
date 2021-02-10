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

import (
	"encoding/base64"

	ParticleType "github.com/aerospike/aerospike-client-go/internal/particle_type"
)

// ExpressionArgument is used for passing arguments to filter expressions.
// The accptable arguments are:
// Value, ExpressionFilter, []*CDTContext
type ExpressionArgument interface {
	pack(BufferEx) (int, error)
}

type ExpType uint

var (
	// ExpTypeNIL is NIL Expression Type
	ExpTypeNIL ExpType = 0
	// ExpTypeBOOL is BOOLEAN Expression Type
	ExpTypeBOOL ExpType = 1
	// ExpTypeINT is INTEGER Expression Type
	ExpTypeINT ExpType = 2
	// ExpTypeSTRING is STRING Expression Type
	ExpTypeSTRING ExpType = 3
	// ExpTypeLIST is LIST Expression Type
	ExpTypeLIST ExpType = 4
	// ExpTypeMAP is MAP Expression Type
	ExpTypeMAP ExpType = 5
	// ExpTypeBLOB is BLOB Expression Type
	ExpTypeBLOB ExpType = 6
	// ExpTypeFLOAT is FLOAT Expression Type
	ExpTypeFLOAT ExpType = 7
	// ExpTypeGEO is GEO String Expression Type
	ExpTypeGEO ExpType = 8
	// ExpTypeHLL is HLL Expression Type
	ExpTypeHLL ExpType = 9
)

type expOp uint

var (
	expOpEQ            expOp = 1
	expOpNE            expOp = 2
	expOpGT            expOp = 3
	expOpGE            expOp = 4
	expOpLT            expOp = 5
	expOpLE            expOp = 6
	expOpREGEX         expOp = 7
	expOpGEO           expOp = 8
	expOpAND           expOp = 16
	expOpOR            expOp = 17
	expOpNOT           expOp = 18
	expOpDIGEST_MODULO expOp = 64
	expOpDEVICE_SIZE   expOp = 65
	expOpLAST_UPDATE   expOp = 66
	expOpSINCE_UPDATE  expOp = 67
	expOpVOID_TIME     expOp = 68
	expOpTTL           expOp = 69
	expOpSET_NAME      expOp = 70
	expOpKEY_EXISTS    expOp = 71
	expOpIS_TOMBSTONE  expOp = 72
	expOpMEMORY_SIZE   expOp = 73
	expOpKEY           expOp = 80
	expOpBIN           expOp = 81
	expOpBIN_TYPE      expOp = 82
	expOpQUOTED        expOp = 126
	expOpCALL          expOp = 127
)

const MODIFY = 0x40

// ExpRegexFlags is used to change the Regex Mode in Expression Filters.
type ExpRegexFlags int

const (
	// ExpRegexFlagNONE uses regex defaults.
	ExpRegexFlagNONE ExpRegexFlags = 0

	// ExpRegexFlagEXTENDED uses POSIX Extended Regular Expression syntax when interpreting regex.
	ExpRegexFlagEXTENDED ExpRegexFlags = 1

	// ExpRegexFlagICASE does not differentiate cases.
	ExpRegexFlagICASE ExpRegexFlags = 2

	// ExpRegexFlagNOSUB does not report position of matches.
	ExpRegexFlagNOSUB ExpRegexFlags = 3

	// ExpRegexFlagNEWLINE does not Match-any-character operators don't match a newline.
	ExpRegexFlagNEWLINE ExpRegexFlags = 8
)

// Filter expression, which can be applied to most commands, to control which records are
// affected by the command. Filter expression are created using the functions in the
// [expressions](crate::expressions) module and its submodules.
type FilterExpression struct {
	// The Operation code
	cmd *expOp
	// The Primary Value of the Operation
	val Value
	// The Bin to use it on (REGEX for example)
	bin *FilterExpression
	// The additional flags for the Operation (REGEX or return_type of Module for example)
	flags *int64
	// The optional Module flag for Module operations or Bin Types
	module *ExpType
	// Sub commands for the CmdExp operation
	exps []*FilterExpression

	arguments []ExpressionArgument
}

func newFilterExpression(
	cmd *expOp,
	val Value,
	bin *FilterExpression,
	flags *int64,
	module *ExpType,
	exps []*FilterExpression,
) *FilterExpression {
	return &FilterExpression{
		cmd:       cmd,
		val:       val,
		bin:       bin,
		flags:     flags,
		module:    module,
		exps:      exps,
		arguments: nil,
	}
}

func (fe *FilterExpression) packExpression(
	exps []*FilterExpression,
	buf BufferEx,
) (int, error) {
	size := 0

	sz, err := packArrayBegin(buf, len(exps)+1)
	size += sz
	if err != nil {
		return size, err
	}

	sz, err = packAInt64(buf, int64(*fe.cmd))
	size += sz
	if err != nil {
		return size, err
	}

	for _, exp := range exps {
		sz, err = exp.pack(buf)
		size += sz
		if err != nil {
			return size, err
		}
	}
	return size, nil
}

func (fe *FilterExpression) packCommand(cmd *expOp, buf BufferEx) (int, error) {
	size := 0

	switch cmd {
	case &expOpREGEX:
		sz, err := packArrayBegin(buf, 4)
		if err != nil {
			return size, err
		}
		size += sz
		// The Operation
		sz, err = packAInt64(buf, int64(*cmd))
		if err != nil {
			return size, err
		}
		size += sz
		// Regex Flags
		sz, err = packAInt64(buf, *fe.flags)
		if err != nil {
			return size, err
		}
		size += sz
		// Raw String is needed instead of the msgpack String that the pack_value method would use.
		sz, err = packRawString(buf, fe.val.String())
		if err != nil {
			return size, err
		}
		size += sz
		// The Bin
		sz, err = fe.bin.pack(buf)
		if err != nil {
			return size, err
		}
		size += sz
	case &expOpCALL:
		// Packing logic for Module
		sz, err := packArrayBegin(buf, 5)
		if err != nil {
			return size, err
		}
		size += sz
		// The Operation
		sz, err = packAInt64(buf, int64(*cmd))
		if err != nil {
			return size, err
		}
		size += sz
		// The Module Operation
		sz, err = packAInt64(buf, int64(*fe.module))
		if err != nil {
			return size, err
		}
		size += sz
		// The Module (List/Map or Bitwise)
		sz, err = packAInt64(buf, *fe.flags)
		if err != nil {
			return size, err
		}
		size += sz
		// Encoding the Arguments
		if args := fe.arguments; len(args) > 0 {
			argLen := 0
			for _, arg := range args {
				// First match to estimate the Size and write the Context
				switch v := arg.(type) {
				case Value, *FilterExpression:
					argLen += 1
				case cdtContextList:
					if len(v) > 0 {
						sz, err = packArrayBegin(buf, 3)
						if err != nil {
							return size, err
						}
						size += sz

						sz, err = packAInt64(buf, 0xff)
						if err != nil {
							return size, err
						}
						size += sz

						sz, err = packArrayBegin(buf, len(v)*2)
						if err != nil {
							return size, err
						}
						size += sz

						for _, c := range v {
							sz, err = c.pack(buf)
							if err != nil {
								return size, err
							}
							size += sz
						}
					}
				default:
					panic("Value `%v` is not acceptable in Expression Filters as an argument")
				}
			}
			sz, err = packArrayBegin(buf, argLen)
			if err != nil {
				return size, err
			}
			size += sz
			// Second match to write the real values
			for _, arg := range args {
				switch val := arg.(type) {
				case Value:
					sz, err := val.pack(buf)
					if err != nil {
						return size, err
					}
					size += sz
				case *FilterExpression:
					sz, err = val.pack(buf)
					if err != nil {
						return size, err
					}
					size += sz
				default:
				}
			}
		} else {
			// No Arguments
			sz, err = fe.val.pack(buf)
			if err != nil {
				return size, err
			}
			size += sz
		}
		// Write the Bin
		sz, err = fe.bin.pack(buf)
		if err != nil {
			return size, err
		}
		size += sz
	case &expOpBIN:
		// Bin Encoder
		sz, err := packArrayBegin(buf, 3)
		if err != nil {
			return size, err
		}
		size += sz
		// The Bin Operation
		sz, err = packAInt64(buf, int64(*cmd))
		if err != nil {
			return size, err
		}
		size += sz
		// The Bin Type (INT/String etc.)
		sz, err = packAInt64(buf, int64(*fe.module))
		if err != nil {
			return size, err
		}
		size += sz
		// The name - Raw String is needed instead of the msgpack String that the pack_value method would use.
		sz, err = packRawString(buf, fe.val.String())
		if err != nil {
			return size, err
		}
		size += sz
	case &expOpBIN_TYPE:
		// BinType encoder
		sz, err := packArrayBegin(buf, 2)
		if err != nil {
			return size, err
		}
		size += sz
		// BinType Operation
		sz, err = packAInt64(buf, int64(*cmd))
		if err != nil {
			return size, err
		}
		size += sz
		// The name - Raw String is needed instead of the msgpack String that the pack_value method would use.
		sz, err = packRawString(buf, fe.val.String())
		if err != nil {
			return size, err
		}
		size += sz
	default:
		// Packing logic for all other Ops
		if value := fe.val; value != nil {
			// Operation has a Value
			sz, err := packArrayBegin(buf, 2)
			if err != nil {
				return size, err
			}
			size += sz
			// Write the Operation
			sz, err = packAInt64(buf, int64(*cmd))
			if err != nil {
				return size, err
			}
			size += sz
			// Write the Value
			sz, err = value.pack(buf)
			if err != nil {
				return size, err
			}
			size += sz
		} else {
			// Operation has no Value
			sz, err := packArrayBegin(buf, 1)
			if err != nil {
				return size, err
			}
			size += sz
			// Write the Operation
			sz, err = packAInt64(buf, int64(*cmd))
			if err != nil {
				return size, err
			}
			size += sz
		}
	}

	return size, nil
}

func (fe *FilterExpression) packValue(buf BufferEx) (int, error) {
	// Packing logic for Value based Ops
	return fe.val.pack(buf)
}

func (fe *FilterExpression) pack(buf BufferEx) (int, error) {
	if len(fe.exps) > 0 {
		return fe.packExpression(fe.exps, buf)
	} else if fe.cmd != nil {
		return fe.packCommand(fe.cmd, buf)
	}
	return fe.packValue(buf)
}

func (fe *FilterExpression) base64() (string, error) {
	sz, err := fe.pack(nil)
	if err != nil {
		return "", err
	}

	input := newBuffer(sz)
	_, err = fe.pack(input)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(input.dataBuffer[:input.dataOffset]), nil
}

// ExpKey creates a record key expression of specified type.
func ExpKey(exp_type ExpType) *FilterExpression {
	return newFilterExpression(
		&expOpKEY,
		IntegerValue(int64(exp_type)),
		nil,
		nil,
		nil,
		nil,
	)
}

// ExpKeyExists creates a function that returns if the primary key is stored in the record meta data
// as a boolean expression. This would occur when `send_key` is true on record write.
func ExpKeyExists() *FilterExpression {
	return newFilterExpression(&expOpKEY_EXISTS, nil, nil, nil, nil, nil)
}

// ExpIntBin creates a 64 bit int bin expression.
func ExpIntBin(name string) *FilterExpression {
	return newFilterExpression(
		&expOpBIN,
		StringValue(name),
		nil,
		nil,
		&ExpTypeINT,
		nil,
	)
}

// ExpStringBin creates a string bin expression.
func ExpStringBin(name string) *FilterExpression {
	return newFilterExpression(
		&expOpBIN,
		StringValue(name),
		nil,
		nil,
		&ExpTypeSTRING,
		nil,
	)
}

// ExpBlobBin creates a blob bin expression.
func ExpBlobBin(name string) *FilterExpression {
	return newFilterExpression(
		&expOpBIN,
		StringValue(name),
		nil,
		nil,
		&ExpTypeBLOB,
		nil,
	)
}

// ExpFloatBin creates a 64 bit float bin expression.
func ExpFloatBin(name string) *FilterExpression {
	return newFilterExpression(
		&expOpBIN,
		StringValue(name),
		nil,
		nil,
		&ExpTypeFLOAT,
		nil,
	)
}

// ExpGeoBin creates a geo bin expression.
func ExpGeoBin(name string) *FilterExpression {
	return newFilterExpression(
		&expOpBIN,
		StringValue(name),
		nil,
		nil,
		&ExpTypeGEO,
		nil,
	)
}

// ExpListBin creates a list bin expression.
func ExpListBin(name string) *FilterExpression {
	return newFilterExpression(
		&expOpBIN,
		StringValue(name),
		nil,
		nil,
		&ExpTypeLIST,
		nil,
	)
}

// ExpMapBin creates a map bin expression.
func ExpMapBin(name string) *FilterExpression {
	return newFilterExpression(
		&expOpBIN,
		StringValue(name),
		nil,
		nil,
		&ExpTypeMAP,
		nil,
	)
}

// ExpHLLBin creates a a HLL bin expression
func ExpHLLBin(name string) *FilterExpression {
	return newFilterExpression(
		&expOpBIN,
		StringValue(name),
		nil,
		nil,
		&ExpTypeHLL,
		nil,
	)
}

// ExpBinExists creates a function that returns if bin of specified name exists.
func ExpBinExists(name string) *FilterExpression {
	return ExpNotEq(ExpBinType(name), ExpIntVal(ParticleType.NULL))
}

// ExpBinType creates a function that returns bin's integer particle type.
func ExpBinType(name string) *FilterExpression {
	return newFilterExpression(
		&expOpBIN_TYPE,
		StringValue(name),
		nil,
		nil,
		nil,
		nil,
	)
}

// ExpSetName creates a function that returns record set name string.
func ExpSetName() *FilterExpression {
	return newFilterExpression(&expOpSET_NAME, nil, nil, nil, nil, nil)
}

// ExpDeviceSize creates a function that returns record size on disk.
// If server storage-engine is memory, then zero is returned.
func ExpDeviceSize() *FilterExpression {
	return newFilterExpression(&expOpDEVICE_SIZE, nil, nil, nil, nil, nil)
}

// ExpMemorySize creates expression that returns record size in memory. If server storage-engine is
// not memory nor data-in-memory, then zero is returned. This expression usually evaluates
// quickly because record meta data is cached in memory.
func ExpMemorySize() *FilterExpression {
	return newFilterExpression(&expOpMEMORY_SIZE, nil, nil, nil, nil, nil)
}

// ExpLastUpdate creates a function that returns record last update time expressed as 64 bit integer
// nanoseconds since 1970-01-01 epoch.
func ExpLastUpdate() *FilterExpression {
	return newFilterExpression(&expOpLAST_UPDATE, nil, nil, nil, nil, nil)
}

// ExpSinceUpdate creates a expression that returns milliseconds since the record was last updated.
// This expression usually evaluates quickly because record meta data is cached in memory.
func ExpSinceUpdate() *FilterExpression {
	return newFilterExpression(&expOpSINCE_UPDATE, nil, nil, nil, nil, nil)
}

// ExpVoidTime creates a function that returns record expiration time expressed as 64 bit integer
// nanoseconds since 1970-01-01 epoch.
func ExpVoidTime() *FilterExpression {
	return newFilterExpression(&expOpVOID_TIME, nil, nil, nil, nil, nil)
}

// ExpTTL creates a function that returns record expiration time (time to live) in integer seconds.
func ExpTTL() *FilterExpression {
	return newFilterExpression(&expOpTTL, nil, nil, nil, nil, nil)
}

// ExpIsTombstone creates a expression that returns if record has been deleted and is still in tombstone state.
// This expression usually evaluates quickly because record meta data is cached in memory.
func ExpIsTombstone() *FilterExpression {
	return newFilterExpression(&expOpIS_TOMBSTONE, nil, nil, nil, nil, nil)
}

// ExpDigestModulo creates a function that returns record digest modulo as integer.
func ExpDigestModulo(modulo int64) *FilterExpression {
	return newFilterExpression(
		&expOpDIGEST_MODULO,
		NewValue(modulo),
		nil,
		nil,
		nil,
		nil,
	)
}

// ExpRegexCompare creates a function like regular expression string operation.
func ExpRegexCompare(regex string, flags ExpRegexFlags, bin *FilterExpression) *FilterExpression {
	iflags := int64(flags)
	return newFilterExpression(
		&expOpREGEX,
		StringValue(regex),
		bin,
		&iflags,
		nil,
		nil,
	)
}

// ExpGeoCompare creates a compare geospatial operation.
func ExpGeoCompare(left *FilterExpression, right *FilterExpression) *FilterExpression {
	return newFilterExpression(
		&expOpGEO,
		nil,
		nil,
		nil,
		nil,
		[]*FilterExpression{left, right},
	)
}

// createss 64 bit integer value
func ExpIntVal(val int64) *FilterExpression {
	return newFilterExpression(nil, IntegerValue(val), nil, nil, nil, nil)
}

// createss a Boolean value
func ExpBoolVal(val bool) *FilterExpression {
	return newFilterExpression(nil, _BoolValue(val), nil, nil, nil, nil)
}

// createss String bin value
func ExpStringVal(val string) *FilterExpression {
	return newFilterExpression(nil, StringValue(val), nil, nil, nil, nil)
}

// createss 64 bit float bin value
func ExpFloatVal(val float64) *FilterExpression {
	return newFilterExpression(nil, FloatValue(val), nil, nil, nil, nil)
}

// createss Blob bin value
func ExpBlobVal(val []byte) *FilterExpression {
	return newFilterExpression(nil, BytesValue(val), nil, nil, nil, nil)
}

// ExpListVal creates a List bin Value
func ExpListVal(val ...Value) *FilterExpression {
	return newFilterExpression(
		&expOpQUOTED,
		ValueArray(val),
		nil,
		nil,
		nil,
		nil,
	)
}

// ExpValueArrayVal creates a List bin Value
func ExpValueArrayVal(val ValueArray) *FilterExpression {
	return newFilterExpression(
		&expOpQUOTED,
		val,
		nil,
		nil,
		nil,
		nil,
	)
}

// ExpListValueVal creates a List bin Value
func ExpListValueVal(val ...interface{}) *FilterExpression {
	return newFilterExpression(
		&expOpQUOTED,
		NewListValue(val),
		nil,
		nil,
		nil,
		nil,
	)
}

// ExpMapVal creates a Map bin Value
func ExpMapVal(val MapValue) *FilterExpression {
	return newFilterExpression(nil, val, nil, nil, nil, nil)
}

// ExpGeoVal creates a geospatial json string value.
func ExpGeoVal(val string) *FilterExpression {
	return newFilterExpression(nil, GeoJSONValue(val), nil, nil, nil, nil)
}

// ExpNil creates a a Nil Value
func ExpNilValue() *FilterExpression {
	return newFilterExpression(nil, nullValue, nil, nil, nil, nil)
}

// ExpNot creates a "not" operator expression.
func ExpNot(exp *FilterExpression) *FilterExpression {
	return &FilterExpression{
		cmd:       &expOpNOT,
		val:       nil,
		bin:       nil,
		flags:     nil,
		module:    nil,
		exps:      []*FilterExpression{exp},
		arguments: nil,
	}
}

// ExpAnd creates a "and" (&&) operator that applies to a variable number of expressions.
func ExpAnd(exps ...*FilterExpression) *FilterExpression {
	return &FilterExpression{
		cmd:       &expOpAND,
		val:       nil,
		bin:       nil,
		flags:     nil,
		module:    nil,
		exps:      exps,
		arguments: nil,
	}
}

// ExpOr creates a "or" (||) operator that applies to a variable number of expressions.
func ExpOr(exps ...*FilterExpression) *FilterExpression {
	return &FilterExpression{
		cmd:       &expOpOR,
		val:       nil,
		bin:       nil,
		flags:     nil,
		module:    nil,
		exps:      exps,
		arguments: nil,
	}
}

// ExpEq creates a equal (==) expression.
func ExpEq(left *FilterExpression, right *FilterExpression) *FilterExpression {
	return &FilterExpression{
		cmd:       &expOpEQ,
		val:       nil,
		bin:       nil,
		flags:     nil,
		module:    nil,
		exps:      []*FilterExpression{left, right},
		arguments: nil,
	}
}

// ExpNotEq creates a not equal (!=) expression
func ExpNotEq(left *FilterExpression, right *FilterExpression) *FilterExpression {
	return &FilterExpression{
		cmd:       &expOpNE,
		val:       nil,
		bin:       nil,
		flags:     nil,
		module:    nil,
		exps:      []*FilterExpression{left, right},
		arguments: nil,
	}
}

// ExpGreater creates a greater than (>) operation.
func ExpGreater(left *FilterExpression, right *FilterExpression) *FilterExpression {
	return &FilterExpression{
		cmd:       &expOpGT,
		val:       nil,
		bin:       nil,
		flags:     nil,
		module:    nil,
		exps:      []*FilterExpression{left, right},
		arguments: nil,
	}
}

// ExpGreaterEq creates a greater than or equal (>=) operation.
func ExpGreaterEq(left *FilterExpression, right *FilterExpression) *FilterExpression {
	return &FilterExpression{
		cmd:       &expOpGE,
		val:       nil,
		bin:       nil,
		flags:     nil,
		module:    nil,
		exps:      []*FilterExpression{left, right},
		arguments: nil,
	}
}

// ExpLess creates a less than (<) operation.
func ExpLess(left *FilterExpression, right *FilterExpression) *FilterExpression {
	return &FilterExpression{
		cmd:       &expOpLT,
		val:       nil,
		bin:       nil,
		flags:     nil,
		module:    nil,
		exps:      []*FilterExpression{left, right},
		arguments: nil,
	}
}

// ExpLessEq creates a less than or equals (<=) operation.
func ExpLessEq(left *FilterExpression, right *FilterExpression) *FilterExpression {
	return &FilterExpression{
		cmd:       &expOpLE,
		val:       nil,
		bin:       nil,
		flags:     nil,
		module:    nil,
		exps:      []*FilterExpression{left, right},
		arguments: nil,
	}
}
