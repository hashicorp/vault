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

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/aerospike/aerospike-client-go/v5/logger"
	"github.com/aerospike/aerospike-client-go/v5/types"

	ParticleType "github.com/aerospike/aerospike-client-go/v5/internal/particle_type"
	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

const (
	// Flags commented out are not supported by cmd client.
	// Contains a read operation.
	_INFO1_READ int = (1 << 0)
	// Get all bins.
	_INFO1_GET_ALL int = (1 << 1)
	// Batch read or exists.
	_INFO1_BATCH int = (1 << 3)

	// Do not read the bins
	_INFO1_NOBINDATA int = (1 << 5)

	// Involve all replicas in read operation.
	_INFO1_READ_MODE_AP_ALL = (1 << 6)

	// Tell server to compress its response.
	_INFO1_COMPRESS_RESPONSE = (1 << 7)

	// Create or update record
	_INFO2_WRITE int = (1 << 0)
	// Fling a record into the belly of Moloch.
	_INFO2_DELETE int = (1 << 1)
	// Update if expected generation == old.
	_INFO2_GENERATION int = (1 << 2)
	// Update if new generation >= old, good for restore.
	_INFO2_GENERATION_GT int = (1 << 3)
	// Transaction resulting in record deletion leaves tombstone (Enterprise only).
	_INFO2_DURABLE_DELETE int = (1 << 4)
	// Create only. Fail if record already exists.
	_INFO2_CREATE_ONLY int = (1 << 5)

	// Return a result for every operation.
	_INFO2_RESPOND_ALL_OPS int = (1 << 7)

	// This is the last of a multi-part message.
	_INFO3_LAST int = (1 << 0)
	// Commit to master only before declaring success.
	_INFO3_COMMIT_MASTER int = (1 << 1)
	// Partition is complete response in scan.
	_INFO3_PARTITION_DONE int = (1 << 2)
	// Update only. Merge bins.
	_INFO3_UPDATE_ONLY int = (1 << 3)

	// Create or completely replace record.
	_INFO3_CREATE_OR_REPLACE int = (1 << 4)
	// Completely replace existing record only.
	_INFO3_REPLACE_ONLY int = (1 << 5)
	// See Below
	_INFO3_SC_READ_TYPE int = (1 << 6)
	// See Below
	_INFO3_SC_READ_RELAX int = (1 << 7)

	// Interpret SC_READ bits in info3.
	//
	// RELAX   TYPE
	//	                strict
	//	                ------
	//   0      0     sequential (default)
	//   0      1     linearize
	//
	//	                relaxed
	//	                -------
	//   1      0     allow replica
	//   1      1     allow unavailable

	_MSG_TOTAL_HEADER_SIZE     uint8 = 30
	_FIELD_HEADER_SIZE         uint8 = 5
	_OPERATION_HEADER_SIZE     uint8 = 8
	_MSG_REMAINING_HEADER_SIZE uint8 = 22
	_DIGEST_SIZE               uint8 = 20
	_COMPRESS_THRESHOLD        int   = 128
	_CL_MSG_VERSION            int64 = 2
	_AS_MSG_TYPE               int64 = 3
	_AS_MSG_TYPE_COMPRESSED    int64 = 4
)

// command interface describes all commands available
type command interface {
	getPolicy(ifc command) Policy

	writeBuffer(ifc command) Error
	getNode(ifc command) (*Node, Error)
	getConnection(policy Policy) (*Connection, Error)
	putConnection(conn *Connection)
	parseResult(ifc command, conn *Connection) Error
	parseRecordResults(ifc command, receiveSize int) (bool, Error)
	prepareRetry(ifc command, isTimeout bool) bool

	execute(ifc command, isRead bool) Error
	executeAt(ifc command, policy *BasePolicy, isRead bool, deadline time.Time, iterations, commandSentCounter int) Error

	canPutConnBack() bool

	// Executes the command
	Execute() Error
}

// Holds data buffer for the command
type baseCommand struct {
	bufferEx

	node *Node
	conn *Connection

	// dataBufferCompress is not a second buffer; it is just a pointer to
	// the beginning of the dataBuffer.
	// To avoid allocating multiple buffers before compression, the dataBuffer
	// will be referencing to a padded buffer. After the command is written to
	// the buffer, this padding will be used to compress the command in-place,
	// and then the compressed proto header will be written.
	dataBufferCompress []byte
	// oneShot determines if streaming commands like query, scan or queryAggregate
	// are not retried if they error out mid-parsing
	oneShot bool

	// will determine if the buffer will be compressed
	// before being sent to the server
	compressed bool
}

// Writes the command for write operations
func (cmd *baseCommand) setWrite(policy *WritePolicy, operation OperationType, key *Key, bins []*Bin, binMap BinMap) Error {
	cmd.begin()
	fieldCount, err := cmd.estimateKeySize(key, policy.SendKey)
	if err != nil {
		return err
	}

	predSize := 0
	if policy.FilterExpression != nil {
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	if binMap == nil {
		for i := range bins {
			if err := cmd.estimateOperationSizeForBin(bins[i]); err != nil {
				return err
			}
		}
	} else {
		for name, value := range binMap {
			if err := cmd.estimateOperationSizeForBinNameAndValue(name, value); err != nil {
				return err
			}
		}
	}

	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return err
	}

	if binMap == nil {
		cmd.writeHeaderWithPolicy(policy, 0, _INFO2_WRITE, fieldCount, len(bins))
	} else {
		cmd.writeHeaderWithPolicy(policy, 0, _INFO2_WRITE, fieldCount, len(binMap))
	}

	if err := cmd.writeKey(key, policy.SendKey); err != nil {
		return err
	}

	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}

	if binMap == nil {
		for i := range bins {
			if err := cmd.writeOperationForBin(bins[i], operation); err != nil {
				return err
			}
		}
	} else {
		for name, value := range binMap {
			if err := cmd.writeOperationForBinNameAndValue(name, value, operation); err != nil {
				return err
			}
		}
	}

	cmd.end()
	cmd.markCompressed(policy)

	return nil
}

// Writes the command for delete operations
func (cmd *baseCommand) setDelete(policy *WritePolicy, key *Key) Error {
	cmd.begin()
	fieldCount, err := cmd.estimateKeySize(key, false)
	if err != nil {
		return err
	}

	predSize := 0
	if policy.FilterExpression != nil {
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return err
	}
	cmd.writeHeaderWithPolicy(policy, 0, _INFO2_WRITE|_INFO2_DELETE, fieldCount, 0)
	if err := cmd.writeKey(key, false); err != nil {
		return err
	}
	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}
	cmd.end()
	cmd.markCompressed(policy)

	return nil

}

// Writes the command for touch operations
func (cmd *baseCommand) setTouch(policy *WritePolicy, key *Key) Error {
	cmd.begin()
	fieldCount, err := cmd.estimateKeySize(key, policy.SendKey)
	if err != nil {
		return err
	}

	predSize := 0
	if policy.FilterExpression != nil {
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	cmd.estimateOperationSize()
	if err := cmd.sizeBuffer(false); err != nil {
		return err
	}
	cmd.writeHeaderWithPolicy(policy, 0, _INFO2_WRITE, fieldCount, 1)
	if err := cmd.writeKey(key, policy.SendKey); err != nil {
		return err
	}
	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}
	cmd.writeOperationForOperationType(_TOUCH)
	cmd.end()
	return nil

}

// Writes the command for exist operations
func (cmd *baseCommand) setExists(policy *BasePolicy, key *Key) Error {
	cmd.begin()
	fieldCount, err := cmd.estimateKeySize(key, false)
	if err != nil {
		return err
	}

	predSize := 0
	if policy.FilterExpression != nil {
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	if err := cmd.sizeBuffer(false); err != nil {
		return err
	}
	cmd.writeHeader(policy, _INFO1_READ|_INFO1_NOBINDATA, 0, fieldCount, 0)
	if err := cmd.writeKey(key, false); err != nil {
		return err
	}
	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}
	cmd.end()
	return nil

}

// Writes the command for get operations (all bins)
func (cmd *baseCommand) setReadForKeyOnly(policy *BasePolicy, key *Key) Error {
	cmd.begin()
	fieldCount, err := cmd.estimateKeySize(key, false)
	if err != nil {
		return err
	}
	predSize := 0
	if policy.FilterExpression != nil {
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}
	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return err
	}
	cmd.writeHeader(policy, _INFO1_READ|_INFO1_GET_ALL, 0, fieldCount, 0)
	if err := cmd.writeKey(key, false); err != nil {
		return err
	}
	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}
	cmd.end()
	cmd.markCompressed(policy)

	return nil

}

// Writes the command for get operations (specified bins)
func (cmd *baseCommand) setRead(policy *BasePolicy, key *Key, binNames []string) Error {
	if len(binNames) > 0 {
		cmd.begin()
		fieldCount, err := cmd.estimateKeySize(key, false)
		if err != nil {
			return err
		}

		predSize := 0
		if policy.FilterExpression != nil {
			predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
			if err != nil {
				return err
			}
			if predSize > 0 {
				fieldCount++
			}
		} else if len(policy.PredExp) > 0 {
			predSize = cmd.estimatePredExpSize(policy.PredExp)
			fieldCount++
		}

		for i := range binNames {
			cmd.estimateOperationSizeForBinName(binNames[i])
		}
		if err := cmd.sizeBuffer(policy.compress()); err != nil {
			return nil
		}
		cmd.writeHeader(policy, _INFO1_READ, 0, fieldCount, len(binNames))
		if err := cmd.writeKey(key, false); err != nil {
			return err
		}

		if policy.FilterExpression != nil {
			if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
				return err
			}
		} else if len(policy.PredExp) > 0 {
			if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
				return err
			}
		}

		for i := range binNames {
			cmd.writeOperationForBinName(binNames[i], _READ)
		}
		cmd.end()
		cmd.markCompressed(policy)
		return nil
	}
	return cmd.setReadForKeyOnly(policy, key)
}

// Writes the command for getting metadata operations
func (cmd *baseCommand) setReadHeader(policy *BasePolicy, key *Key) Error {
	cmd.begin()
	fieldCount, err := cmd.estimateKeySize(key, false)
	if err != nil {
		return err
	}

	predSize := 0
	if policy.FilterExpression != nil {
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	cmd.estimateOperationSizeForBinName("")
	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return err
	}

	cmd.writeHeader(policy, _INFO1_READ|_INFO1_NOBINDATA, 0, fieldCount, 1)

	if err := cmd.writeKey(key, false); err != nil {
		return err
	}
	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}
	cmd.writeOperationForBinName("", _READ)
	cmd.end()
	cmd.markCompressed(policy)

	return nil

}

// Implements different command operations
func (cmd *baseCommand) setOperate(policy *WritePolicy, key *Key, args *operateArgs) Error {
	if len(args.operations) == 0 {
		return newError(types.PARAMETER_ERROR, "No operations were passed.")
	}

	cmd.begin()
	fieldCount := 0

	for i := range args.operations {
		if err := cmd.estimateOperationSizeForOperation(args.operations[i], false); err != nil {
			return err
		}
	}

	ksz, err := cmd.estimateKeySize(key, policy.SendKey && args.hasWrite)
	if err != nil {
		return err
	}
	fieldCount += ksz

	predSize := 0
	if policy.FilterExpression != nil {
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return err
	}

	if args.writeAttr != 0 {
		cmd.writeHeaderWithPolicy(policy, args.readAttr, args.writeAttr, fieldCount, len(args.operations))
	} else {
		cmd.writeHeader(&policy.BasePolicy, args.readAttr, args.writeAttr, fieldCount, len(args.operations))
	}

	if err := cmd.writeKey(key, policy.SendKey && args.hasWrite); err != nil {
		return err
	}

	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}

	for _, operation := range args.operations {
		if err := cmd.writeOperationForOperation(operation); err != nil {
			return err
		}
	}

	cmd.end()
	cmd.markCompressed(policy)

	return nil
}

func (cmd *baseCommand) setUdf(policy *WritePolicy, key *Key, packageName string, functionName string, args *ValueArray) Error {
	cmd.begin()
	fieldCount, err := cmd.estimateKeySize(key, policy.SendKey)
	if err != nil {
		return err
	}

	predSize := 0
	if policy.FilterExpression != nil {
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	fc, err := cmd.estimateUdfSize(packageName, functionName, args)
	if err != nil {
		return err
	}
	fieldCount += fc

	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return err
	}

	cmd.writeHeaderWithPolicy(policy, 0, _INFO2_WRITE, fieldCount, 0)
	if err := cmd.writeKey(key, policy.SendKey); err != nil {
		return err
	}
	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}
	cmd.writeFieldString(packageName, UDF_PACKAGE_NAME)
	cmd.writeFieldString(functionName, UDF_FUNCTION)
	if err := cmd.writeUdfArgs(args); err != nil {
		return err
	}
	cmd.end()
	cmd.markCompressed(policy)

	return nil
}

func (cmd *baseCommand) setBatchRead(policy *BatchPolicy, keys []*Key, batch *batchNode, binNames []string, ops []*Operation, readAttr int) Error {
	offsets := batch.offsets
	max := len(batch.offsets)
	fieldCountRow := 1
	if policy.SendSetName {
		fieldCountRow = 2
	}

	// Estimate buffer size
	cmd.begin()
	fieldCount := 1
	predSize := 0
	if policy.FilterExpression != nil {
		var err Error
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}
	cmd.dataOffset += predSize

	cmd.dataOffset += int(_FIELD_HEADER_SIZE) + 5

	var prev *Key
	for i := 0; i < max; i++ {
		key := keys[offsets[i]]
		cmd.dataOffset += len(key.digest) + 4

		// Try reference equality in hope that namespace/set for all keys is set from fixed variables.
		if prev != nil && prev.namespace == key.namespace &&
			(!policy.SendSetName || prev.setName == key.setName) {
			// Can set repeat previous namespace/bin names to save space.
			cmd.dataOffset++
		} else {
			// Must write full header and namespace/set/bin names.
			cmd.dataOffset += len(key.namespace) + int(_FIELD_HEADER_SIZE) + 6

			if policy.SendSetName {
				cmd.dataOffset += len(key.setName) + int(_FIELD_HEADER_SIZE)
			}

			if len(binNames) > 0 {
				for _, binName := range binNames {
					cmd.estimateOperationSizeForBinName(binName)
				}
			} else if len(ops) > 0 {
				for _, op := range ops {
					if err := cmd.estimateOperationSizeForOperation(op, true); err != nil {
						return err
					}
				}
			}

			prev = key
		}
	}

	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return err
	}

	if policy.ReadModeAP == ReadModeAPAll {
		readAttr |= _INFO1_READ_MODE_AP_ALL
	}

	cmd.writeHeader(&policy.BasePolicy, readAttr|_INFO1_BATCH, 0, fieldCount, 0)

	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}

	// Write real field size.
	fieldSizeOffset := cmd.dataOffset
	if policy.SendSetName {
		cmd.writeFieldHeader(0, BATCH_INDEX_WITH_SET)
	} else {
		cmd.writeFieldHeader(0, BATCH_INDEX)
	}

	cmd.WriteUint32(uint32(max))

	if policy.AllowInline {
		cmd.WriteByte(1)
	} else {
		cmd.WriteByte(0)
	}

	prev = nil
	for i := 0; i < max; i++ {
		index := offsets[i]
		cmd.WriteUint32(uint32(index))

		key := keys[index]
		if _, err := cmd.Write(key.digest[:]); err != nil {
			return newCommonError(err)
		}
		// Try reference equality in hope that namespace/set for all keys is set from fixed variables.
		if prev != nil && prev.namespace == key.namespace &&
			(!policy.SendSetName || prev.setName == key.setName) {
			// Can set repeat previous namespace/bin names to save space.
			cmd.WriteByte(1) // repeat
		} else {
			// Write full header, namespace and bin names.
			cmd.WriteByte(0) // do not repeat
			if len(binNames) > 0 {
				cmd.WriteByte(byte(readAttr))
				cmd.writeBatchFields(policy, key, fieldCountRow, len(binNames))
				for _, binName := range binNames {
					cmd.writeOperationForBinName(binName, _READ)
				}
			} else if len(ops) > 0 {
				offset := cmd.dataOffset
				cmd.dataOffset++
				cmd.writeBatchFields(policy, key, fieldCountRow, len(ops))
				cmd.dataBuffer[offset], _ = cmd.writeBatchReadOperations(ops, readAttr)
			} else {
				attr := byte(readAttr)
				if len(binNames) == 0 {
					attr |= byte(_INFO1_GET_ALL)
				} else {
					attr |= byte(_INFO1_NOBINDATA)
				}
				cmd.WriteByte(attr)
				cmd.writeBatchFields(policy, key, fieldCountRow, 0)
			}

			prev = key
		}
	}

	cmd.WriteUint32At(uint32(cmd.dataOffset)-uint32(_MSG_TOTAL_HEADER_SIZE)-4, fieldSizeOffset)
	cmd.end()
	cmd.markCompressed(policy)

	return nil
}

func (cmd *baseCommand) setBatchIndexRead(policy *BatchPolicy, records []*BatchRead, batch *batchNode) Error {
	offsets := batch.offsets
	max := len(batch.offsets)
	fieldCountRow := 1
	if policy.SendSetName {
		fieldCountRow = 2
	}

	// Estimate buffer size
	cmd.begin()
	fieldCount := 1
	predSize := 0
	if policy.FilterExpression != nil {
		var err Error
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	cmd.dataOffset += int(_FIELD_HEADER_SIZE) + 5

	var prev *BatchRead
	for i := 0; i < max; i++ {
		record := records[offsets[i]]
		key := record.Key
		binNames := record.BinNames
		ops := record.Ops

		cmd.dataOffset += len(key.digest) + 4

		// Try reference equality in hope that namespace/set for all keys is set from fixed variables.
		if prev != nil && prev.Key.namespace == key.namespace &&
			(!policy.SendSetName || prev.Key.setName == key.setName) &&
			&prev.BinNames == &binNames && prev.ReadAllBins == record.ReadAllBins &&
			&prev.Ops == &ops {
			// Can set repeat previous namespace/bin names to save space.
			cmd.dataOffset++
		} else {
			// Must write full header and namespace/set/bin names.
			cmd.dataOffset += len(key.namespace) + int(_FIELD_HEADER_SIZE) + 6

			if policy.SendSetName {
				cmd.dataOffset += len(key.setName) + int(_FIELD_HEADER_SIZE)
			}

			if len(binNames) != 0 {
				for _, binName := range binNames {
					cmd.estimateOperationSizeForBinName(binName)
				}
			} else if len(ops) != 0 {
				for _, op := range ops {
					cmd.estimateOperationSizeForOperation(op, true)
				}
			}

			prev = record
		}
	}

	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return err
	}

	readAttr := _INFO1_READ
	if policy.ReadModeAP == ReadModeAPAll {
		readAttr |= _INFO1_READ_MODE_AP_ALL
	}

	cmd.writeHeader(&policy.BasePolicy, readAttr|_INFO1_BATCH, 0, fieldCount, 0)

	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}

	// Write real field size.
	fieldSizeOffset := cmd.dataOffset
	if policy.SendSetName {
		cmd.writeFieldHeader(0, BATCH_INDEX_WITH_SET)
	} else {
		cmd.writeFieldHeader(0, BATCH_INDEX)
	}

	cmd.WriteUint32(uint32(max))

	if policy.AllowInline {
		cmd.WriteByte(1)
	} else {
		cmd.WriteByte(0)
	}

	prev = nil
	for i := 0; i < max; i++ {
		index := offsets[i]
		cmd.WriteUint32(uint32(index))

		record := records[index]
		key := record.Key
		binNames := record.BinNames
		ops := record.Ops
		if _, err := cmd.Write(key.digest[:]); err != nil {
			return newCommonError(err)
		}

		// Try reference equality in hope that namespace/set for all keys is set from fixed variables.
		if prev != nil && prev.Key.namespace == key.namespace &&
			(!policy.SendSetName || prev.Key.setName == key.setName) &&
			&prev.BinNames == &binNames && prev.ReadAllBins == record.ReadAllBins &&
			&prev.Ops == &ops {
			// Can set repeat previous namespace/bin names to save space.
			cmd.WriteByte(1) // repeat
		} else {
			// Write full header, namespace and bin names.
			cmd.WriteByte(0) // do not repeat
			if len(binNames) > 0 {
				cmd.WriteByte(byte(readAttr))
				cmd.writeBatchFields(policy, key, fieldCountRow, len(binNames))
				for _, binName := range binNames {
					cmd.writeOperationForBinName(binName, _READ)
				}
			} else if len(ops) > 0 {
				offset := cmd.dataOffset
				cmd.dataOffset++
				cmd.writeBatchFields(policy, key, fieldCountRow, len(ops))
				cmd.dataBuffer[offset], _ = cmd.writeBatchReadOperations(ops, readAttr)
			} else {
				attr := byte(readAttr)
				if record.ReadAllBins {
					attr |= byte(_INFO1_GET_ALL)
				} else {
					attr |= byte(_INFO1_NOBINDATA)
				}
				cmd.WriteByte(attr)
				cmd.writeBatchFields(policy, key, fieldCountRow, 0)
			}

			prev = record
		}
	}

	cmd.WriteUint32At(uint32(cmd.dataOffset)-uint32(_MSG_TOTAL_HEADER_SIZE)-4, fieldSizeOffset)
	cmd.end()
	cmd.markCompressed(policy)

	return nil
}

func (cmd *baseCommand) writeBatchFields(policy *BatchPolicy, key *Key, fieldCount, opCount int) Error {
	cmd.WriteUint16(uint16(fieldCount))
	cmd.WriteUint16(uint16(opCount))
	cmd.writeFieldString(key.namespace, NAMESPACE)

	if policy.SendSetName {
		cmd.writeFieldString(key.setName, TABLE)
	}

	return nil
}

func (cmd *baseCommand) setScan(policy *ScanPolicy, namespace *string, setName *string, binNames []string, taskID uint64, nodePartitions *nodePartitions) Error {
	cmd.begin()
	fieldCount := 0
	partsFullSize := 0
	partsPartialSize := 0
	maxRecords := int64(0)

	if nodePartitions != nil {
		partsFullSize = len(nodePartitions.partsFull) * 2
		partsPartialSize = len(nodePartitions.partsPartial) * 20
		maxRecords = nodePartitions.recordMax
	}

	predSize := 0
	if policy.FilterExpression != nil {
		var err Error
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	if namespace != nil {
		cmd.dataOffset += len(*namespace) + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	if setName != nil {
		cmd.dataOffset += len(*setName) + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	if partsFullSize > 0 {
		cmd.dataOffset += partsFullSize + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	if partsPartialSize > 0 {
		cmd.dataOffset += partsPartialSize + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	if maxRecords > 0 {
		cmd.dataOffset += 8 + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	if policy.RecordsPerSecond > 0 {
		cmd.dataOffset += 4 + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	// Estimate scan timeout size.
	cmd.dataOffset += 4 + int(_FIELD_HEADER_SIZE)
	fieldCount++

	// Allocate space for TaskId field.
	cmd.dataOffset += 8 + int(_FIELD_HEADER_SIZE)
	fieldCount++

	for i := range binNames {
		cmd.estimateOperationSizeForBinName(binNames[i])
	}

	if err := cmd.sizeBuffer(false); err != nil {
		return err
	}
	readAttr := _INFO1_READ

	if !policy.IncludeBinData {
		readAttr |= _INFO1_NOBINDATA
	}

	operationCount := 0
	if binNames != nil {
		operationCount = len(binNames)
	}
	cmd.writeHeader(&policy.BasePolicy, readAttr, 0, fieldCount, operationCount)

	if namespace != nil {
		cmd.writeFieldString(*namespace, NAMESPACE)
	}

	if setName != nil {
		cmd.writeFieldString(*setName, TABLE)
	}

	if partsFullSize > 0 {
		cmd.writeFieldHeader(partsFullSize, PID_ARRAY)

		for _, part := range nodePartitions.partsFull {
			cmd.WriteInt16LittleEndian(uint16(part.id))
		}
	}

	if partsPartialSize > 0 {
		cmd.writeFieldHeader(partsPartialSize, DIGEST_ARRAY)

		for _, part := range nodePartitions.partsPartial {
			if _, err := cmd.Write(part.digest[:]); err != nil {
				return newCommonError(err)
			}
		}
	}

	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}

	if maxRecords > 0 {
		cmd.writeFieldInt64(maxRecords, SCAN_MAX_RECORDS)
	}

	if policy.RecordsPerSecond > 0 {
		cmd.writeFieldInt32(int32(policy.RecordsPerSecond), RECORDS_PER_SECOND)
	}

	// Write scan timeout
	cmd.writeFieldHeader(4, SCAN_TIMEOUT)
	cmd.WriteInt32(int32(policy.SocketTimeout / time.Millisecond)) // in milliseconds

	cmd.writeFieldHeader(8, TRAN_ID)
	cmd.WriteUint64(taskID)

	for i := range binNames {
		cmd.writeOperationForBinName(binNames[i], _READ)
	}

	cmd.end()

	return nil
}

func (cmd *baseCommand) setQuery(policy *QueryPolicy, wpolicy *WritePolicy, statement *Statement, taskID uint64, operations []*Operation, write bool, nodePartitions *nodePartitions) Error {
	fieldCount := 0
	filterSize := 0
	binNameSize := 0
	predSize := 0
	partsFullSize := 0
	partsPartialSize := 0
	maxRecords := int64(0)

	recordsPerSecond := 0
	if !write {
		recordsPerSecond = policy.RecordsPerSecond
	}

	cmd.begin()

	if statement.Namespace != "" {
		cmd.dataOffset += len(statement.Namespace) + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	if statement.IndexName != "" {
		cmd.dataOffset += len(statement.IndexName) + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	if statement.SetName != "" {
		cmd.dataOffset += len(statement.SetName) + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	// Allocate space for TaskId field.
	cmd.dataOffset += 8 + int(_FIELD_HEADER_SIZE)
	fieldCount++

	if statement.Filter != nil {
		idxType := statement.Filter.IndexCollectionType()

		if idxType != ICT_DEFAULT {
			cmd.dataOffset += int(_FIELD_HEADER_SIZE) + 1
			fieldCount++
		}

		cmd.dataOffset += int(_FIELD_HEADER_SIZE)
		filterSize++ // num filters

		sz, err := statement.Filter.EstimateSize()
		if err != nil {
			return err
		}
		filterSize += sz

		cmd.dataOffset += filterSize
		fieldCount++

		// Query bin names are specified as a field (Scan bin names are specified later as operations)
		if len(statement.BinNames) > 0 {
			cmd.dataOffset += int(_FIELD_HEADER_SIZE)
			binNameSize++ // num bin names

			for _, binName := range statement.BinNames {
				binNameSize += len(binName) + 1
			}
			cmd.dataOffset += binNameSize
			fieldCount++
		}
	} else {
		// Calling query with no filters is more efficiently handled by a primary index scan.
		// Estimate scan options size.
		if nodePartitions != nil {
			partsFullSize = len(nodePartitions.partsFull) * 2
			partsPartialSize = len(nodePartitions.partsPartial) * 20
			maxRecords = nodePartitions.recordMax
		}

		if partsFullSize > 0 {
			cmd.dataOffset += partsFullSize + int(_FIELD_HEADER_SIZE)
			fieldCount++
		}

		if partsPartialSize > 0 {
			cmd.dataOffset += partsPartialSize + int(_FIELD_HEADER_SIZE)
			fieldCount++
		}

		// Estimate max records size;
		if maxRecords > 0 {
			cmd.dataOffset += 8 + int(_FIELD_HEADER_SIZE)
			fieldCount++
		}

		// Estimate scan timeout size.
		cmd.dataOffset += (4 + int(_FIELD_HEADER_SIZE))
		fieldCount++

		// Estimate records per second size.
		if recordsPerSecond > 0 {
			cmd.dataOffset += 4 + int(_FIELD_HEADER_SIZE)
			fieldCount++
		}
	}

	if policy.FilterExpression != nil {
		var err Error
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	var functionArgs *ValueArray
	if statement.functionName != "" {
		cmd.dataOffset += int(_FIELD_HEADER_SIZE) + 1 // udf type
		cmd.dataOffset += len(statement.packageName) + int(_FIELD_HEADER_SIZE)
		cmd.dataOffset += len(statement.functionName) + int(_FIELD_HEADER_SIZE)

		// function args
		cmd.dataOffset += int(_FIELD_HEADER_SIZE)
		if len(statement.functionArgs) > 0 {
			functionArgs = NewValueArray(statement.functionArgs)
			fasz, err := functionArgs.EstimateSize()
			if err != nil {
				return err
			}

			cmd.dataOffset += fasz
		}

		fieldCount += 4
	}

	// Operations (used in query execute) and bin names (used in scan/query) are mutually exclusive.
	if len(operations) > 0 {
		for _, op := range operations {
			if err := cmd.estimateOperationSizeForOperation(op, false); err != nil {
				return err
			}
		}
	} else if len(statement.BinNames) > 0 && statement.Filter == nil {
		for _, binName := range statement.BinNames {
			cmd.estimateOperationSizeForBinName(binName)
		}
	}

	if err := cmd.sizeBuffer(false); err != nil {
		return err
	}

	operationCount := 0
	if len(operations) > 0 {
		operationCount = len(operations)
	} else if statement.Filter == nil && len(statement.BinNames) > 0 {
		operationCount = len(statement.BinNames)
	}

	if write {
		cmd.writeHeaderWithPolicy(wpolicy, 0, _INFO2_WRITE, fieldCount, operationCount)
	} else {
		readAttr := _INFO1_READ | _INFO1_NOBINDATA
		if policy.IncludeBinData {
			readAttr = _INFO1_READ
		}
		cmd.writeHeader(&policy.BasePolicy, readAttr, 0, fieldCount, operationCount)
	}

	if statement.Namespace != "" {
		cmd.writeFieldString(statement.Namespace, NAMESPACE)
	}

	if statement.IndexName != "" {
		cmd.writeFieldString(statement.IndexName, INDEX_NAME)
	}

	if statement.SetName != "" {
		cmd.writeFieldString(statement.SetName, TABLE)
	}

	cmd.writeFieldHeader(8, TRAN_ID)
	cmd.WriteUint64(taskID)

	if statement.Filter != nil {
		idxType := statement.Filter.IndexCollectionType()

		if idxType != ICT_DEFAULT {
			cmd.writeFieldHeader(1, INDEX_TYPE)
			cmd.WriteByte(byte(idxType))
		}

		cmd.writeFieldHeader(filterSize, INDEX_RANGE)
		cmd.WriteByte(byte(1)) // number of filters

		_, err := statement.Filter.write(cmd)
		if err != nil {
			return err
		}

		if len(statement.BinNames) > 0 {
			cmd.writeFieldHeader(binNameSize, QUERY_BINLIST)
			cmd.WriteByte(byte(len(statement.BinNames)))

			for _, binName := range statement.BinNames {
				len := copy(cmd.dataBuffer[cmd.dataOffset+1:], binName)
				cmd.dataBuffer[cmd.dataOffset] = byte(len)
				cmd.dataOffset += len + 1
			}
		}
	} else {
		// Calling query with no filters is more efficiently handled by a primary index scan.

		if partsFullSize > 0 {
			cmd.writeFieldHeader(partsFullSize, PID_ARRAY)

			for _, part := range nodePartitions.partsFull {
				cmd.WriteInt16LittleEndian(uint16(part.id))
			}
		}

		if partsPartialSize > 0 {
			cmd.writeFieldHeader(partsPartialSize, DIGEST_ARRAY)

			for _, part := range nodePartitions.partsPartial {
				if _, err := cmd.Write(part.digest[:]); err != nil {
					return newCommonError(err)
				}
			}
		}

		if maxRecords > 0 {
			cmd.writeFieldInt64(maxRecords, SCAN_MAX_RECORDS)
		}

		// Write scan timeout
		cmd.writeFieldHeader(4, SCAN_TIMEOUT)
		cmd.WriteInt32(int32(policy.SocketTimeout / time.Millisecond)) // in milliseconds

		// Write records per second.
		if recordsPerSecond > 0 {
			cmd.writeFieldInt32(int32(recordsPerSecond), RECORDS_PER_SECOND)
		}
	}

	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}

	if statement.functionName != "" {
		cmd.writeFieldHeader(1, UDF_OP)
		if statement.returnData {
			cmd.dataBuffer[cmd.dataOffset] = byte(1)
		} else {
			cmd.dataBuffer[cmd.dataOffset] = byte(2)
		}
		cmd.dataOffset++

		cmd.writeFieldString(statement.packageName, UDF_PACKAGE_NAME)
		cmd.writeFieldString(statement.functionName, UDF_FUNCTION)
		if err := cmd.writeUdfArgs(functionArgs); err != nil {
			return err
		}
	}

	if len(operations) > 0 {
		for _, op := range operations {
			if err := cmd.writeOperationForOperation(op); err != nil {
				return err
			}
		}
	} else if len(statement.BinNames) > 0 && statement.Filter == nil {
		// scan binNames come last
		for _, binName := range statement.BinNames {
			cmd.writeOperationForBinName(binName, _READ)
		}
	}

	cmd.end()

	return nil
}

func (cmd *baseCommand) estimateKeySize(key *Key, sendKey bool) (int, Error) {
	fieldCount := 0

	if key.namespace != "" {
		cmd.dataOffset += len(key.namespace) + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	if key.setName != "" {
		cmd.dataOffset += len(key.setName) + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	cmd.dataOffset += int(_DIGEST_SIZE + _FIELD_HEADER_SIZE)
	fieldCount++

	if sendKey {
		// field header size + key size
		sz, err := key.userKey.EstimateSize()
		if err != nil {
			return sz, err
		}
		cmd.dataOffset += sz + int(_FIELD_HEADER_SIZE) + 1
		fieldCount++
	}

	return fieldCount, nil
}

func (cmd *baseCommand) estimateUdfSize(packageName string, functionName string, args *ValueArray) (int, Error) {
	cmd.dataOffset += len(packageName) + int(_FIELD_HEADER_SIZE)
	cmd.dataOffset += len(functionName) + int(_FIELD_HEADER_SIZE)

	sz, err := args.EstimateSize()
	if err != nil {
		return 0, err
	}

	// fmt.Println(args, sz)

	cmd.dataOffset += sz + int(_FIELD_HEADER_SIZE)
	return 3, nil
}

func (cmd *baseCommand) estimateOperationSizeForBin(bin *Bin) Error {
	cmd.dataOffset += len(bin.Name) + int(_OPERATION_HEADER_SIZE)
	sz, err := bin.Value.EstimateSize()
	if err != nil {
		return err
	}
	cmd.dataOffset += sz
	return nil
}

func (cmd *baseCommand) estimateOperationSizeForBinNameAndValue(name string, value interface{}) Error {
	cmd.dataOffset += len(name) + int(_OPERATION_HEADER_SIZE)
	sz, err := NewValue(value).EstimateSize()
	if err != nil {
		return err
	}
	cmd.dataOffset += sz
	return nil
}

func (cmd *baseCommand) estimateOperationSizeForOperation(operation *Operation, isBatch bool) Error {
	if isBatch && operation.opType.isWrite {
		return newError(types.PARAMETER_ERROR, "Write operations not allowed in batch read")
	}

	binLen := len(operation.binName)
	cmd.dataOffset += binLen + int(_OPERATION_HEADER_SIZE)

	if operation.encoder == nil {
		if operation.binValue != nil {
			sz, err := operation.binValue.EstimateSize()
			if err != nil {
				return err
			}
			cmd.dataOffset += sz
		}
	} else {
		sz, err := operation.encoder(operation, nil)
		if err != nil {
			return err
		}
		cmd.dataOffset += sz
	}
	return nil
}

func (cmd *baseCommand) estimateOperationSizeForBinName(binName string) {
	cmd.dataOffset += len(binName) + int(_OPERATION_HEADER_SIZE)
}

func (cmd *baseCommand) estimateOperationSize() {
	cmd.dataOffset += int(_OPERATION_HEADER_SIZE)
}

func (cmd *baseCommand) estimatePredExpSize(predExp []PredExp) int {
	sz := 0
	for _, predexp := range predExp {
		sz += predexp.marshaledSize()
	}
	cmd.dataOffset += sz + int(_FIELD_HEADER_SIZE)
	return sz
}

func (cmd *baseCommand) estimateExpressionSize(exp *Expression) (int, Error) {
	size, err := exp.pack(nil)
	if err != nil {
		return size, err
	}

	cmd.dataOffset += size + int(_FIELD_HEADER_SIZE)
	return size, nil
}

// Generic header write.
func (cmd *baseCommand) writeHeader(policy *BasePolicy, readAttr int, writeAttr int, fieldCount int, operationCount int) {
	infoAttr := 0

	switch policy.ReadModeSC {
	case ReadModeSCSession:
	case ReadModeSCLinearize:
		infoAttr |= _INFO3_SC_READ_TYPE
	case ReadModeSCAllowReplica:
		infoAttr |= _INFO3_SC_READ_RELAX
	case ReadModeSCAllowUnavailable:
		infoAttr |= _INFO3_SC_READ_TYPE | _INFO3_SC_READ_RELAX
	}

	if policy.ReadModeAP == ReadModeAPAll {
		readAttr |= _INFO1_READ_MODE_AP_ALL
	}

	if policy.UseCompression {
		readAttr |= _INFO1_COMPRESS_RESPONSE
	}

	// Write all header data except total size which must be written last.
	cmd.dataBuffer[8] = _MSG_REMAINING_HEADER_SIZE // Message header length.
	cmd.dataBuffer[9] = byte(readAttr)
	cmd.dataBuffer[10] = byte(writeAttr)
	cmd.dataBuffer[11] = byte(infoAttr)

	for i := 12; i < 26; i++ {
		cmd.dataBuffer[i] = 0
	}
	cmd.dataOffset = 26
	cmd.WriteInt16(int16(fieldCount))
	cmd.WriteInt16(int16(operationCount))
	cmd.dataOffset = int(_MSG_TOTAL_HEADER_SIZE)
}

// Header write for write operations.
func (cmd *baseCommand) writeHeaderWithPolicy(policy *WritePolicy, readAttr int, writeAttr int, fieldCount int, operationCount int) {
	// Set flags.
	generation := uint32(0)
	infoAttr := 0

	switch policy.RecordExistsAction {
	case UPDATE:
	case UPDATE_ONLY:
		infoAttr |= _INFO3_UPDATE_ONLY
	case REPLACE:
		infoAttr |= _INFO3_CREATE_OR_REPLACE
	case REPLACE_ONLY:
		infoAttr |= _INFO3_REPLACE_ONLY
	case CREATE_ONLY:
		writeAttr |= _INFO2_CREATE_ONLY
	}

	switch policy.GenerationPolicy {
	case NONE:
	case EXPECT_GEN_EQUAL:
		generation = policy.Generation
		writeAttr |= _INFO2_GENERATION
	case EXPECT_GEN_GT:
		generation = policy.Generation
		writeAttr |= _INFO2_GENERATION_GT
	}

	if policy.CommitLevel == COMMIT_MASTER {
		infoAttr |= _INFO3_COMMIT_MASTER
	}

	if policy.DurableDelete {
		writeAttr |= _INFO2_DURABLE_DELETE
	}

	switch policy.ReadModeSC {
	case ReadModeSCSession:
	case ReadModeSCLinearize:
		infoAttr |= _INFO3_SC_READ_TYPE
	case ReadModeSCAllowReplica:
		infoAttr |= _INFO3_SC_READ_RELAX
	case ReadModeSCAllowUnavailable:
		infoAttr |= _INFO3_SC_READ_TYPE | _INFO3_SC_READ_RELAX
	}

	if policy.ReadModeAP == ReadModeAPAll {
		readAttr |= _INFO1_READ_MODE_AP_ALL
	}

	if policy.UseCompression {
		readAttr |= _INFO1_COMPRESS_RESPONSE
	}

	// Write all header data except total size which must be written last.
	cmd.dataBuffer[8] = _MSG_REMAINING_HEADER_SIZE // Message header length.
	cmd.dataBuffer[9] = byte(readAttr)
	cmd.dataBuffer[10] = byte(writeAttr)
	cmd.dataBuffer[11] = byte(infoAttr)
	cmd.dataBuffer[12] = 0 // unused
	cmd.dataBuffer[13] = 0 // clear the result code
	cmd.dataOffset = 14
	cmd.WriteUint32(generation)
	cmd.dataOffset = 18
	cmd.WriteUint32(policy.Expiration)

	// Initialize timeout. It will be written later.
	cmd.dataBuffer[22] = 0
	cmd.dataBuffer[23] = 0
	cmd.dataBuffer[24] = 0
	cmd.dataBuffer[25] = 0

	cmd.dataOffset = 26
	cmd.WriteInt16(int16(fieldCount))
	cmd.WriteInt16(int16(operationCount))
	cmd.dataOffset = int(_MSG_TOTAL_HEADER_SIZE)
}

func (cmd *baseCommand) writeKey(key *Key, sendKey bool) Error {
	// Write key into buffer.
	if key.namespace != "" {
		cmd.writeFieldString(key.namespace, NAMESPACE)
	}

	if key.setName != "" {
		cmd.writeFieldString(key.setName, TABLE)
	}

	cmd.writeFieldBytes(key.digest[:], DIGEST_RIPE)

	if sendKey {
		if err := cmd.writeFieldValue(key.userKey, KEY); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *baseCommand) writeOperationForBin(bin *Bin, operation OperationType) Error {
	nameLength := copy(cmd.dataBuffer[(cmd.dataOffset+int(_OPERATION_HEADER_SIZE)):], bin.Name)

	valueLength, err := bin.Value.EstimateSize()
	if err != nil {
		return err
	}

	cmd.WriteInt32(int32(nameLength + valueLength + 4))
	cmd.WriteByte((operation.op))
	cmd.WriteByte((byte(bin.Value.GetType())))
	cmd.WriteByte((byte(0)))
	cmd.WriteByte((byte(nameLength)))
	cmd.dataOffset += nameLength
	_, err = bin.Value.write(cmd)
	return err
}

func (cmd *baseCommand) writeOperationForBinNameAndValue(name string, val interface{}, operation OperationType) Error {
	nameLength := copy(cmd.dataBuffer[(cmd.dataOffset+int(_OPERATION_HEADER_SIZE)):], name)

	v := NewValue(val)

	valueLength, err := v.EstimateSize()
	if err != nil {
		return err
	}

	cmd.WriteInt32(int32(nameLength + valueLength + 4))
	cmd.WriteByte((operation.op))
	cmd.WriteByte((byte(v.GetType())))
	cmd.WriteByte((byte(0)))
	cmd.WriteByte((byte(nameLength)))
	cmd.dataOffset += nameLength
	_, err = v.write(cmd)
	return err
}

func (cmd *baseCommand) writeBatchReadOperations(ops []*Operation, readAttr int) (byte, Error) {
	readBin := false
	readHeader := false

	for _, op := range ops {
		switch op.opType {
		case _READ:
			// Read all bins if no bin is specified.
			if len(op.binName) == 0 {
				readAttr |= _INFO1_GET_ALL
			}
			readBin = true

			if op.headerOnly {
				readHeader = true
			}
		default:
		}
		if err := cmd.writeOperationForOperation(op); err != nil {
			return byte(readAttr), err
		}
	}

	if readHeader && !readBin {
		readAttr |= _INFO1_NOBINDATA
	}
	return byte(readAttr), nil
}

func (cmd *baseCommand) writeOperationForOperation(operation *Operation) Error {
	nameLength := copy(cmd.dataBuffer[(cmd.dataOffset+int(_OPERATION_HEADER_SIZE)):], operation.binName)

	if operation.used {
		// cahce will set the used flag to false again
		if err := operation.cache(); err != nil {
			return err
		}
	}

	if operation.encoder == nil {
		valueLength, err := operation.binValue.EstimateSize()
		if err != nil {
			return err
		}

		cmd.WriteInt32(int32(nameLength + valueLength + 4))
		cmd.WriteByte((operation.opType.op))
		cmd.WriteByte((byte(operation.binValue.GetType())))
		cmd.WriteByte((byte(0)))
		cmd.WriteByte((byte(nameLength)))
		cmd.dataOffset += nameLength
		_, err = operation.binValue.write(cmd)
		return err
	}

	valueLength, err := operation.encoder(operation, nil)
	if err != nil {
		return err
	}

	cmd.WriteInt32(int32(nameLength + valueLength + 4))
	cmd.WriteByte((operation.opType.op))
	cmd.WriteByte((byte(ParticleType.BLOB)))
	cmd.WriteByte((byte(0)))
	cmd.WriteByte((byte(nameLength)))
	cmd.dataOffset += nameLength
	_, err = operation.encoder(operation, cmd)
	//mark the operation as used, so that it will be cached the next time it is used
	operation.used = err == nil
	return err
}

func (cmd *baseCommand) writeOperationForBinName(name string, operation OperationType) {
	nameLength := copy(cmd.dataBuffer[(cmd.dataOffset+int(_OPERATION_HEADER_SIZE)):], name)
	cmd.WriteInt32(int32(nameLength + 4))
	cmd.WriteByte((operation.op))
	cmd.WriteByte(byte(0))
	cmd.WriteByte(byte(0))
	cmd.WriteByte(byte(nameLength))
	cmd.dataOffset += nameLength
}

func (cmd *baseCommand) writeOperationForOperationType(operation OperationType) {
	cmd.WriteInt32(int32(4))
	cmd.WriteByte(operation.op)
	cmd.WriteByte(0)
	cmd.WriteByte(0)
	cmd.WriteByte(0)
}

func (cmd *baseCommand) writePredExp(predExp []PredExp, predSize int) Error {
	cmd.writeFieldHeader(predSize, FILTER_EXP)
	for i := range predExp {
		predExp[i].marshal(cmd)
	}
	return nil
}

func (cmd *baseCommand) writeFilterExpression(exp *Expression, expSize int) Error {
	cmd.writeFieldHeader(expSize, FILTER_EXP)
	if _, err := exp.pack(cmd); err != nil {
		return err
	}
	return nil
}

func (cmd *baseCommand) writeFieldValue(value Value, ftype FieldType) Error {
	vlen, err := value.EstimateSize()
	if err != nil {
		return err
	}
	cmd.writeFieldHeader(vlen+1, ftype)
	cmd.WriteByte(byte(value.GetType()))

	_, err = value.write(cmd)
	return err
}

func (cmd *baseCommand) writeUdfArgs(value *ValueArray) Error {
	if value != nil {
		vlen, err := value.EstimateSize()
		if err != nil {
			return err
		}
		cmd.writeFieldHeader(vlen, UDF_ARGLIST)
		_, err = value.pack(cmd)
		return err
	}

	cmd.writeFieldHeader(0, UDF_ARGLIST)
	return nil
}

func (cmd *baseCommand) writeFieldInt32(val int32, ftype FieldType) {
	cmd.writeFieldHeader(4, ftype)
	cmd.WriteInt32(val)
}

func (cmd *baseCommand) writeFieldInt64(val int64, ftype FieldType) {
	cmd.writeFieldHeader(8, ftype)
	cmd.WriteInt64(val)
}

func (cmd *baseCommand) writeFieldString(str string, ftype FieldType) {
	flen := copy(cmd.dataBuffer[(cmd.dataOffset+int(_FIELD_HEADER_SIZE)):], str)
	cmd.writeFieldHeader(flen, ftype)
	cmd.dataOffset += flen
}

func (cmd *baseCommand) writeFieldBytes(bytes []byte, ftype FieldType) {
	copy(cmd.dataBuffer[cmd.dataOffset+int(_FIELD_HEADER_SIZE):], bytes)

	cmd.writeFieldHeader(len(bytes), ftype)
	cmd.dataOffset += len(bytes)
}

func (cmd *baseCommand) writeFieldHeader(size int, ftype FieldType) {
	cmd.WriteInt32(int32(size + 1))
	cmd.WriteByte((byte(ftype)))
}

func (cmd *baseCommand) begin() {
	cmd.dataOffset = int(_MSG_TOTAL_HEADER_SIZE)
}

func (cmd *baseCommand) sizeBuffer(compress bool) Error {
	return cmd.sizeBufferSz(cmd.dataOffset, compress)
}

func (cmd *baseCommand) validateHeader(header int64) Error {
	msgVersion := (uint64(header) & 0xFF00000000000000) >> 56
	if msgVersion != 2 {
		return newCustomNodeError(cmd.node, types.PARSE_ERROR, fmt.Sprintf("Invalid Message Header: Expected version to be 2, but got %v", msgVersion))
	}

	msgType := (uint64(header) & 0x00FF000000000000) >> 49
	if !(msgType == 1 || msgType == 3 || msgType == 4) {
		return newCustomNodeError(cmd.node, types.PARSE_ERROR, fmt.Sprintf("Invalid Message Header: Expected type to be 1, 3 or 4, but got %v", msgType))
	}

	msgSize := header & 0x0000FFFFFFFFFFFF
	if msgSize > int64(MaxBufferSize) {
		return newCustomNodeError(cmd.node, types.PARSE_ERROR, fmt.Sprintf("Invalid Message Header: Expected size to be under 10MiB, but got %v", msgSize))
	}

	return nil
}

var (
	// MaxBufferSize protects against allocating massive memory blocks
	// for buffers. Tweak this number if you are returning a lot of
	// LDT elements in your queries.
	MaxBufferSize = 1024 * 1024 * 120 // 120 MB
)

const (
	msgHeaderPad  = 16
	zlibHeaderPad = 2
)

func (cmd *baseCommand) sizeBufferSz(size int, willCompress bool) Error {

	if willCompress {
		// adds zlib and proto pads to the size of the buffer
		size += msgHeaderPad + zlibHeaderPad
	}

	// Corrupted data streams can result in a huge length.
	// Do a sanity check here.
	if size > MaxBufferSize || size < 0 {
		return newCustomNodeError(cmd.node, types.PARSE_ERROR, fmt.Sprintf("Invalid size for buffer: %d", size))
	}

	if size <= len(cmd.dataBuffer) {
		// don't touch the buffer
		// this is a noop, here to silence the linters
		cmd.dataBuffer = cmd.dataBuffer
	} else if size <= cap(cmd.dataBuffer) {
		cmd.dataBuffer = cmd.dataBuffer[:size]
	} else {
		// not enough space
		cmd.dataBuffer = make([]byte, size)
	}

	// The trick here to keep a ref to the buffer, and set the buffer itself
	// to a padded version of the original:
	// | Proto Header | Original Compressed Size | compressed message |
	// |    8 Bytes   |          8 Bytes         |                    |
	if willCompress {
		cmd.dataBufferCompress = cmd.dataBuffer
		cmd.dataBuffer = cmd.dataBufferCompress[msgHeaderPad+zlibHeaderPad:]
	}

	return nil
}

func (cmd *baseCommand) end() {
	var proto = int64(cmd.dataOffset-8) | (_CL_MSG_VERSION << 56) | (_AS_MSG_TYPE << 48)
	binary.BigEndian.PutUint64(cmd.dataBuffer[0:], uint64(proto))
}

func (cmd *baseCommand) markCompressed(policy Policy) {
	cmd.compressed = policy.compress()
}

func (cmd *baseCommand) compress() Error {
	if cmd.compressed && cmd.dataOffset > _COMPRESS_THRESHOLD {
		b := bytes.NewBuffer(cmd.dataBufferCompress[msgHeaderPad:])
		b.Reset()
		w := zlib.NewWriter(b)

		// There seems to be a bug either in Go's zlib or in zlibc
		// which messes up a single write block of bigger than 64KB to
		// the deflater.
		// Things work in multiple writes of 64KB though, so this is
		// how we're going to do it.
		i := 0
		const step = 64 * 1024
		for i+step < cmd.dataOffset {
			n, err := w.Write(cmd.dataBuffer[i : i+step])
			i += n
			if err != nil {
				return newErrorAndWrap(err, types.SERIALIZE_ERROR)
			}
		}

		if i < cmd.dataOffset {
			if _, err := w.Write(cmd.dataBuffer[i:cmd.dataOffset]); err != nil {
				return newErrorAndWrap(err, types.SERIALIZE_ERROR)
			}
		}

		// flush
		if err := w.Close(); err != nil {
			return newErrorAndWrap(err, types.SERIALIZE_ERROR)
		}

		compressedSz := b.Len()

		// Use compressed buffer if compression completed within original buffer size.
		var proto = int64(compressedSz+8) | (_CL_MSG_VERSION << 56) | (_AS_MSG_TYPE_COMPRESSED << 48)
		binary.BigEndian.PutUint64(cmd.dataBufferCompress[0:], uint64(proto))
		binary.BigEndian.PutUint64(cmd.dataBufferCompress[8:], uint64(cmd.dataOffset))

		cmd.dataBuffer = cmd.dataBufferCompress
		cmd.dataOffset = compressedSz + 16
		cmd.dataBufferCompress = nil
	}

	return nil
}

// isCompressed returns the length of the compressed buffer.
// If the buffer is not compressed, the result will be -1
func (cmd *baseCommand) compressedSize() int {
	proto := Buffer.BytesToInt64(cmd.dataBuffer, 0)
	size := proto & 0xFFFFFFFFFFFF

	msgType := (proto >> 48) & 0xff

	if msgType != _AS_MSG_TYPE_COMPRESSED {
		return -1
	}

	return int(size)
}

////////////////////////////////////

func (cmd *baseCommand) execute(ifc command, isRead bool) Error {
	policy := ifc.getPolicy(ifc).GetBasePolicy()
	deadline := policy.deadline()

	return cmd.executeAt(ifc, policy, isRead, deadline, -1, 0)
}

func (cmd *baseCommand) executeAt(ifc command, policy *BasePolicy, isRead bool, deadline time.Time, iterations, commandSentCounter int) (errChain Error) {
	// for exponential backoff
	interval := policy.SleepBetweenRetries

	notFirstIteration := false
	isClientTimeout := false
	loopCount := 0

	var err Error

	// Execute command until successful, timed out or maximum iterations have been reached.
	for {
		iterations++
		loopCount++

		// too many retries
		if (policy.MaxRetries <= 0 && iterations > 0) || (policy.MaxRetries > 0 && iterations > policy.MaxRetries) {
			return chainErrors(ErrMaxRetriesExceeded.err(), errChain).iter(iterations).setInDoubt(isRead, commandSentCounter)
		}

		// Sleep before trying again, after the first iteration
		if policy.SleepBetweenRetries > 0 && notFirstIteration {
			// Do not sleep if you know you'll wake up after the deadline
			if policy.TotalTimeout > 0 && time.Now().Add(interval).After(deadline) {
				break
			}

			time.Sleep(interval)
			if policy.SleepMultiplier > 1 {
				interval = time.Duration(float64(interval) * policy.SleepMultiplier)
			}
		}

		if notFirstIteration {
			if !ifc.prepareRetry(ifc, isClientTimeout || (err != nil && err.Matches(types.SERVER_NOT_AVAILABLE))) {
				if bc, ok := ifc.(batcher); ok {
					// Batch may be retried in separate commands.
					alreadyRetried, err := bc.retryBatch(bc, cmd.node.cluster, deadline, iterations, commandSentCounter)
					if alreadyRetried {
						// Batch was retried in separate subcommands. Complete this command.
						if err != nil {
							return chainErrors(err, errChain).iter(iterations)
						}
						return nil
					}

					// chain the errors and retry
					if err != nil {
						errChain = chainErrors(err, errChain).iter(iterations)
						continue
					}
				}
			}
		}

		// NOTE: This is important to be after the prepareRetry block above
		isClientTimeout = false

		notFirstIteration = true

		// check for command timeout
		if policy.TotalTimeout > 0 && time.Now().After(deadline) {
			break
		}

		// set command node, so when you return a record it has the node
		cmd.node, err = ifc.getNode(ifc)
		if cmd.node == nil || !cmd.node.IsActive() || err != nil {
			isClientTimeout = false

			// chain the errors
			if err != nil {
				errChain = chainErrors(err, errChain).iter(iterations)
			}

			// Node is currently inactive. Retry.
			continue
		}

		// check if node has encountered too many errors
		if err = cmd.node.validateErrorCount(); err != nil {
			isClientTimeout = false

			// chain the errors
			errChain = chainErrors(err, errChain).iter(iterations).setNode(cmd.node)

			// Max error rate achieved, try again per policy
			continue
		}

		cmd.conn, err = ifc.getConnection(policy)
		if err != nil {
			isClientTimeout = false

			// chain the errors
			errChain = chainErrors(err, errChain).iter(iterations).setNode(cmd.node)

			// exit immediately if connection pool is exhausted and the corresponding policy option is set
			if policy.ExitFastOnExhaustedConnectionPool && errors.Is(err, ErrConnectionPoolExhausted) {
				break
			}

			if errors.Is(err, ErrConnectionPoolEmpty) || errors.Is(err, ErrConnectionPoolExhausted) {
				if errors.Is(err, ErrConnectionPoolExhausted) || (errors.Is(err, ErrConnectionPoolEmpty) && loopCount == 1) {
					isClientTimeout = true
				}
				// if the connection pool is empty, we still haven't tried
				// the transaction to increase the iteration count.
				iterations--
			}
			logger.Logger.Debug("Node " + cmd.node.String() + ": " + err.Error())
			continue
		}

		// Assign the connection buffer to the command buffer
		cmd.dataBuffer = cmd.conn.dataBuffer

		// Set command buffer.
		err = ifc.writeBuffer(ifc)
		if err != nil {
			// chain the errors
			errChain = chainErrors(err, errChain).iter(iterations).setNode(cmd.node)

			// All runtime exceptions are considered fatal. Do not retry.
			// Close socket to flush out possible garbage. Do not put back in pool.
			cmd.conn.Close()
			cmd.conn = nil
			return err
		}

		// Reset timeout in send buffer (destined for server) and socket.
		binary.BigEndian.PutUint32(cmd.dataBuffer[22:], 0)
		if !deadline.IsZero() {
			serverTimeout := time.Until(deadline)
			if serverTimeout < time.Millisecond {
				serverTimeout = time.Millisecond
			}
			binary.BigEndian.PutUint32(cmd.dataBuffer[22:], uint32(serverTimeout/time.Millisecond))
		}

		// now that the deadline has been set in the buffer, compress the contents
		if err = cmd.compress(); err != nil {
			return chainErrors(err, errChain).iter(iterations).setNode(cmd.node)
		}

		// if cmd, ok := ifc.(*operateCommand); ok {
		// 	println("Writing...")
		// 	ioutil.WriteFile("dump_not_ok"+strconv.Itoa(int(time.Now().UnixNano())), []byte(hex.Dump(cmd.dataBuffer[:cmd.dataOffset])), 0644)
		// }

		// Send command.
		_, err = cmd.conn.Write(cmd.dataBuffer[:cmd.dataOffset])
		if err != nil {
			// chain the errors
			errChain = chainErrors(err, errChain).iter(iterations).setNode(cmd.node)

			isClientTimeout = false
			if deviceOverloadError(err) {
				cmd.node.incrErrorCount()
			}

			// IO errors are considered temporary anomalies. Retry.
			// Close socket to flush out possible garbage. Do not put back in pool.
			cmd.conn.Close()
			cmd.conn = nil

			logger.Logger.Debug("Node " + cmd.node.String() + ": " + err.Error())
			continue
		}
		commandSentCounter++

		// Parse results.
		err = ifc.parseResult(ifc, cmd.conn)
		if err != nil {
			// chain the errors
			errChain = chainErrors(err, errChain).iter(iterations).setNode(cmd.node)

			if networkError(err) {
				isTimeout := errors.Is(err, ErrTimeout)
				isClientTimeout = isTimeout
				if !isTimeout {
					if deviceOverloadError(err) {
						cmd.node.incrErrorCount()
					}
				}

				// IO errors are considered temporary anomalies. Retry.
				// Close socket to flush out possible garbage. Do not put back in pool.
				cmd.conn.Close()

				logger.Logger.Debug("Node " + cmd.node.String() + ": " + err.Error())

				// retry only for non-streaming commands
				if !cmd.oneShot {
					cmd.conn = nil
					continue
				}
			}

			// close the connection
			// cancelling/closing the batch/multi commands will return an error, which will
			// close the connection to throw away its data and signal the server about the
			// situation. We will not put back the connection in the buffer.
			if ifc.canPutConnBack() && cmd.conn.IsConnected() && KeepConnection(err) {
				// Put connection back in pool.
				cmd.node.PutConnection(cmd.conn)
			} else {
				cmd.conn.Close()
				cmd.conn = nil
			}

			return errChain.setInDoubt(isRead, commandSentCounter)
		}

		// in case it has grown and re-allocated
		if len(cmd.dataBufferCompress) > len(cmd.dataBuffer) {
			cmd.conn.dataBuffer = cmd.dataBufferCompress
		} else {
			cmd.conn.dataBuffer = cmd.dataBuffer
		}

		// Put connection back in pool.
		// cmd.node.PutConnection(cmd.conn)
		ifc.putConnection(cmd.conn)

		// command has completed successfully. Exit method.
		return nil

	}

	// execution timeout
	errChain = chainErrors(ErrTimeout.err(), errChain).iter(iterations).setNode(cmd.node)
	return errChain
}

func (cmd *baseCommand) canPutConnBack() bool {
	return true
}

func (cmd *baseCommand) parseRecordResults(ifc command, receiveSize int) (bool, Error) {
	panic("Abstract method. Should not end up here")
}

func networkError(err Error) bool {
	return err.Matches(types.NETWORK_ERROR, types.TIMEOUT)
}

func deviceOverloadError(err Error) bool {
	return err.Matches(types.DEVICE_OVERLOAD)
}
