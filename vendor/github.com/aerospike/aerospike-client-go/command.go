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
// limitations under the License.

package aerospike

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"time"

	. "github.com/aerospike/aerospike-client-go/logger"
	. "github.com/aerospike/aerospike-client-go/types"

	ParticleType "github.com/aerospike/aerospike-client-go/internal/particle_type"
	Buffer "github.com/aerospike/aerospike-client-go/utils/buffer"
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

	writeBuffer(ifc command) error
	getNode(ifc command) (*Node, error)
	getConnection(policy Policy) (*Connection, error)
	putConnection(conn *Connection)
	parseResult(ifc command, conn *Connection) error
	parseRecordResults(ifc command, receiveSize int) (bool, error)
	prepareRetry(ifc command, isTimeout bool) bool

	execute(ifc command, isRead bool) error
	executeAt(ifc command, policy *BasePolicy, isRead bool, deadline time.Time, iterations, commandSentCounter int) error

	// Executes the command
	Execute() error
}

// Holds data buffer for the command
type baseCommand struct {
	buffer

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
func (cmd *baseCommand) setWrite(policy *WritePolicy, operation OperationType, key *Key, bins []*Bin, binMap BinMap) error {
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

	cmd.writeKey(key, policy.SendKey)

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
func (cmd *baseCommand) setDelete(policy *WritePolicy, key *Key) error {
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
	cmd.writeKey(key, false)
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
func (cmd *baseCommand) setTouch(policy *WritePolicy, key *Key) error {
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
	cmd.writeKey(key, policy.SendKey)
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
func (cmd *baseCommand) setExists(policy *BasePolicy, key *Key) error {
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
	cmd.writeKey(key, false)
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
func (cmd *baseCommand) setReadForKeyOnly(policy *BasePolicy, key *Key) error {
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
	cmd.writeHeader(policy, _INFO1_READ|_INFO1_GET_ALL, 0, fieldCount, 0)
	cmd.writeKey(key, false)
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

// Writes the command for get operations (specified bins)
func (cmd *baseCommand) setRead(policy *BasePolicy, key *Key, binNames []string) (err error) {
	if len(binNames) > 0 {
		cmd.begin()
		fieldCount, err := cmd.estimateKeySize(key, false)
		if err != nil {
			return err
		}

		predSize := 0
		if policy.FilterExpression != nil {
			if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
				return err
			}
		} else if len(policy.PredExp) > 0 {
			predSize = cmd.estimatePredExpSize(policy.PredExp)
			fieldCount++
		}

		for i := range binNames {
			cmd.estimateOperationSizeForBinName(binNames[i])
		}
		if err = cmd.sizeBuffer(false); err != nil {
			return nil
		}
		cmd.writeHeader(policy, _INFO1_READ, 0, fieldCount, len(binNames))
		cmd.writeKey(key, false)

		if policy.FilterExpression != nil {
			if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
				return err
			}
		} else if len(policy.PredExp) > 0 {
			cmd.writePredExp(policy.PredExp, predSize)
		}

		for i := range binNames {
			cmd.writeOperationForBinName(binNames[i], _READ)
		}
		cmd.end()
	} else {
		err = cmd.setReadForKeyOnly(policy, key)
	}

	return err
}

// Writes the command for getting metadata operations
func (cmd *baseCommand) setReadHeader(policy *BasePolicy, key *Key) error {
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

	cmd.writeKey(key, false)
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
	return nil

}

// Implements different command operations
func (cmd *baseCommand) setOperate(policy *WritePolicy, key *Key, operations []*Operation) (bool, error) {
	if len(operations) == 0 {
		return false, NewAerospikeError(PARAMETER_ERROR, "No operations were passed.")
	}

	cmd.begin()
	fieldCount := 0
	readAttr := 0
	writeAttr := 0
	hasWrite := false
	readBin := false
	readHeader := false
	RespondPerEachOp := policy.RespondPerEachOp

	for i := range operations {
		switch operations[i].opType {
		case _BIT_READ:
			fallthrough
		case _HLL_READ:
			fallthrough
		case _MAP_READ:
			// Map operations require RespondPerEachOp to be true.
			RespondPerEachOp = true
			// Fall through to read.
			fallthrough
		case _READ, _CDT_READ:
			if !operations[i].headerOnly {
				readAttr |= _INFO1_READ

				// Read all bins if no bin is specified.
				if operations[i].binName == "" {
					readAttr |= _INFO1_GET_ALL
				}
				readBin = true
			} else {
				readAttr |= _INFO1_READ
				readHeader = true
			}
		case _BIT_MODIFY:
			fallthrough
		case _HLL_MODIFY:
			fallthrough
		case _MAP_MODIFY:
			// Map operations require RespondPerEachOp to be true.
			RespondPerEachOp = true
			// Fall through to default.
			fallthrough
		default:
			writeAttr = _INFO2_WRITE
			hasWrite = true
		}
		cmd.estimateOperationSizeForOperation(operations[i])
	}

	ksz, err := cmd.estimateKeySize(key, policy.SendKey && hasWrite)
	if err != nil {
		return hasWrite, err
	}
	fieldCount += ksz

	predSize := 0
	if policy.FilterExpression != nil {
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return hasWrite, err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(policy.PredExp) > 0 {
		predSize = cmd.estimatePredExpSize(policy.PredExp)
		fieldCount++
	}

	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return hasWrite, err
	}

	if readHeader && !readBin {
		readAttr |= _INFO1_NOBINDATA
	}

	if RespondPerEachOp {
		writeAttr |= _INFO2_RESPOND_ALL_OPS
	}

	if writeAttr != 0 {
		cmd.writeHeaderWithPolicy(policy, readAttr, writeAttr, fieldCount, len(operations))
	} else {
		cmd.writeHeader(&policy.BasePolicy, readAttr, writeAttr, fieldCount, len(operations))
	}
	cmd.writeKey(key, policy.SendKey && hasWrite)

	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return hasWrite, err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return hasWrite, err
		}
	}

	for _, operation := range operations {
		if err := cmd.writeOperationForOperation(operation); err != nil {
			return hasWrite, err
		}
	}

	cmd.end()
	cmd.markCompressed(policy)

	return hasWrite, nil
}

func (cmd *baseCommand) setUdf(policy *WritePolicy, key *Key, packageName string, functionName string, args *ValueArray) error {
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
	cmd.writeKey(key, policy.SendKey)
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
	cmd.writeUdfArgs(args)
	cmd.end()
	cmd.markCompressed(policy)

	return nil
}

func (cmd *baseCommand) setBatchIndexReadCompat(policy *BatchPolicy, keys []*Key, batch *batchNode, binNames []string, readAttr int) error {
	offsets := batch.offsets
	max := len(batch.offsets)
	fieldCountRow := 1
	if policy.SendSetName {
		fieldCountRow = 2
	}

	binNameSize := 0
	operationCount := len(binNames)
	for _, binName := range binNames {
		binNameSize += len(binName) + int(_OPERATION_HEADER_SIZE)
	}

	// Estimate buffer size
	cmd.begin()
	fieldCount := 1
	predSize := 0
	if policy.FilterExpression != nil {
		var err error
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
			cmd.dataOffset += binNameSize
			prev = key
		}
	}

	if err := cmd.sizeBuffer(policy.compress()); err != nil {
		return err
	}

	if policy.ReadModeAP == ReadModeAPAll {
		readAttr |= _INFO1_READ_MODE_AP_ALL
	}

	if len(binNames) == 0 {
		readAttr |= _INFO1_GET_ALL
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
		cmd.Write(key.digest[:])
		// Try reference equality in hope that namespace/set for all keys is set from fixed variables.
		if prev != nil && prev.namespace == key.namespace &&
			(!policy.SendSetName || prev.setName == key.setName) {
			// Can set repeat previous namespace/bin names to save space.
			cmd.WriteByte(1) // repeat
		} else {
			// Write full header, namespace and bin names.
			cmd.WriteByte(0) // do not repeat
			cmd.WriteByte(byte(readAttr))
			cmd.WriteUint16(uint16(fieldCountRow))
			cmd.WriteUint16(uint16(operationCount))
			cmd.writeFieldString(key.namespace, NAMESPACE)

			if policy.SendSetName {
				cmd.writeFieldString(key.setName, TABLE)
			}

			for _, binName := range binNames {
				cmd.writeOperationForBinName(binName, _READ)
			}

			prev = key
		}
	}

	cmd.WriteUint32At(uint32(cmd.dataOffset)-uint32(_MSG_TOTAL_HEADER_SIZE)-4, fieldSizeOffset)
	cmd.end()
	cmd.markCompressed(policy)

	return nil
}

func (cmd *baseCommand) setBatchIndexRead(policy *BatchPolicy, records []*BatchRead, batch *batchNode) error {
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
		var err error
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

		cmd.dataOffset += len(key.digest) + 4

		// Try reference equality in hope that namespace/set for all keys is set from fixed variables.
		if prev != nil && prev.Key.namespace == key.namespace &&
			(!policy.SendSetName || prev.Key.setName == key.setName) &&
			&prev.BinNames == &binNames && prev.ReadAllBins == record.ReadAllBins {
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
	// cmd.writeHeader(&policy.BasePolicy, _INFO1_READ|_INFO1_BATCH, 0, 1, 0)

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
		cmd.Write(key.digest[:])

		// Try reference equality in hope that namespace/set for all keys is set from fixed variables.
		if prev != nil && prev.Key.namespace == key.namespace &&
			(!policy.SendSetName || prev.Key.setName == key.setName) &&
			&prev.BinNames == &binNames && prev.ReadAllBins == record.ReadAllBins {
			// Can set repeat previous namespace/bin names to save space.
			cmd.WriteByte(1) // repeat
		} else {
			// Write full header, namespace and bin names.
			cmd.WriteByte(0) // do not repeat
			if len(binNames) > 0 {
				cmd.WriteByte(byte(readAttr))
				cmd.WriteUint16(uint16(fieldCountRow))
				cmd.WriteUint16(uint16(len(binNames)))
				cmd.writeFieldString(key.namespace, NAMESPACE)

				if policy.SendSetName {
					cmd.writeFieldString(key.setName, TABLE)
				}

				for _, binName := range binNames {
					cmd.writeOperationForBinName(binName, _READ)
				}
			} else {
				attr := byte(readAttr)
				if record.ReadAllBins {
					attr |= byte(_INFO1_GET_ALL)
				} else {
					attr |= byte(_INFO1_NOBINDATA)
				}
				cmd.WriteByte(attr)

				cmd.WriteUint16(uint16(fieldCountRow))
				cmd.WriteUint16(0)
				cmd.writeFieldString(key.namespace, NAMESPACE)

				if policy.SendSetName {
					cmd.writeFieldString(key.setName, TABLE)
				}
			}

			prev = record
		}
	}

	cmd.WriteUint32At(uint32(cmd.dataOffset)-uint32(_MSG_TOTAL_HEADER_SIZE)-4, fieldSizeOffset)
	cmd.end()
	cmd.markCompressed(policy)

	return nil
}

func (cmd *baseCommand) setScan(policy *ScanPolicy, namespace *string, setName *string, binNames []string, taskID uint64) error {
	cmd.begin()
	fieldCount := 0

	predSize := 0
	if policy.FilterExpression != nil {
		var err error
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

	if policy.RecordsPerSecond > 0 {
		cmd.dataOffset += 4 + int(_FIELD_HEADER_SIZE)
		fieldCount++
	}

	// Estimate scan options size.
	cmd.dataOffset += 2 + int(_FIELD_HEADER_SIZE)
	fieldCount++

	// Estimate scan timeout size.
	cmd.dataOffset += 4 + int(_FIELD_HEADER_SIZE)
	fieldCount++

	// Allocate space for TaskId field.
	cmd.dataOffset += 8 + int(_FIELD_HEADER_SIZE)
	fieldCount++

	if binNames != nil {
		for i := range binNames {
			cmd.estimateOperationSizeForBinName(binNames[i])
		}
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

	if policy.FilterExpression != nil {
		if err := cmd.writeFilterExpression(policy.FilterExpression, predSize); err != nil {
			return err
		}
	} else if len(policy.PredExp) > 0 {
		if err := cmd.writePredExp(policy.PredExp, predSize); err != nil {
			return err
		}
	}

	if policy.RecordsPerSecond > 0 {
		cmd.writeFieldInt32(int32(policy.RecordsPerSecond), RECORDS_PER_SECOND)
	}

	cmd.writeFieldHeader(2, SCAN_OPTIONS)
	priority := byte(policy.Priority)
	priority <<= 4

	if policy.FailOnClusterChange {
		priority |= 0x08
	}

	cmd.WriteByte(priority)
	cmd.WriteByte(byte(policy.ScanPercent))

	// Write scan timeout
	cmd.writeFieldHeader(4, SCAN_TIMEOUT)
	cmd.WriteInt32(int32(policy.SocketTimeout / time.Millisecond)) // in milliseconds

	cmd.writeFieldHeader(8, TRAN_ID)
	cmd.WriteUint64(taskID)

	if binNames != nil {
		for i := range binNames {
			cmd.writeOperationForBinName(binNames[i], _READ)
		}
	}

	cmd.end()

	return nil
}

func (cmd *baseCommand) setQuery(policy *QueryPolicy, wpolicy *WritePolicy, statement *Statement, operations []*Operation, write bool) (err error) {
	fieldCount := 0
	filterSize := 0
	binNameSize := 0
	predSize := 0
	predExp := statement.predExps

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
		cmd.dataOffset += (2 + int(_FIELD_HEADER_SIZE))
		fieldCount++

		// Estimate scan timeout size.
		cmd.dataOffset += (4 + int(_FIELD_HEADER_SIZE))
		fieldCount++

		// Estimate records per second size.
		if recordsPerSecond > 0 {
			cmd.dataOffset += 4 + int(_FIELD_HEADER_SIZE)
			fieldCount++
		}
	}

	if len(policy.PredExp) > 0 && len(predExp) == 0 {
		predExp = policy.PredExp
	}

	if policy.FilterExpression != nil {
		var err error
		predSize, err = cmd.estimateExpressionSize(policy.FilterExpression)
		if err != nil {
			return err
		}
		if predSize > 0 {
			fieldCount++
		}
	} else if len(predExp) > 0 {
		predSize = cmd.estimatePredExpSize(predExp)
		fieldCount++
	}

	var functionArgs *ValueArray
	if statement.functionName != "" {
		cmd.dataOffset += int(_FIELD_HEADER_SIZE) + 1 // udf type
		cmd.dataOffset += len(statement.packageName) + int(_FIELD_HEADER_SIZE)
		cmd.dataOffset += len(statement.functionName) + int(_FIELD_HEADER_SIZE)

		fasz := 0
		if len(statement.functionArgs) > 0 {
			functionArgs = NewValueArray(statement.functionArgs)
			fasz, err = functionArgs.EstimateSize()
			if err != nil {
				return err
			}
		}

		cmd.dataOffset += int(_FIELD_HEADER_SIZE) + fasz
		fieldCount += 4
	}

	// Operations (used in query execute) and bin names (used in scan/query) are mutually exclusive.
	if len(operations) > 0 {
		for _, op := range operations {
			cmd.estimateOperationSizeForOperation(op)
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
	cmd.WriteUint64(statement.TaskId)

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
		cmd.writeFieldHeader(2, SCAN_OPTIONS)
		priority := byte(policy.Priority)
		priority <<= 4

		if !write && policy.FailOnClusterChange {
			priority |= 0x08
		}

		cmd.WriteByte(priority)
		cmd.WriteByte(byte(100))

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
	} else if len(predExp) > 0 {
		if err := cmd.writePredExp(predExp, predSize); err != nil {
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
		cmd.writeUdfArgs(functionArgs)
	}

	if len(operations) > 0 {
		for _, op := range operations {
			cmd.writeOperationForOperation(op)
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

func (cmd *baseCommand) estimateKeySize(key *Key, sendKey bool) (int, error) {
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

func (cmd *baseCommand) estimateUdfSize(packageName string, functionName string, args *ValueArray) (int, error) {
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

func (cmd *baseCommand) estimateOperationSizeForBin(bin *Bin) error {
	cmd.dataOffset += len(bin.Name) + int(_OPERATION_HEADER_SIZE)
	sz, err := bin.Value.EstimateSize()
	if err != nil {
		return err
	}
	cmd.dataOffset += sz
	return nil
}

func (cmd *baseCommand) estimateOperationSizeForBinNameAndValue(name string, value interface{}) error {
	cmd.dataOffset += len(name) + int(_OPERATION_HEADER_SIZE)
	sz, err := NewValue(value).EstimateSize()
	if err != nil {
		return err
	}
	cmd.dataOffset += sz
	return nil
}

func (cmd *baseCommand) estimateOperationSizeForOperation(operation *Operation) error {
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

func (cmd *baseCommand) estimateExpressionSize(exp *FilterExpression) (int, error) {
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

func (cmd *baseCommand) writeKey(key *Key, sendKey bool) {
	// Write key into buffer.
	if key.namespace != "" {
		cmd.writeFieldString(key.namespace, NAMESPACE)
	}

	if key.setName != "" {
		cmd.writeFieldString(key.setName, TABLE)
	}

	cmd.writeFieldBytes(key.digest[:], DIGEST_RIPE)

	if sendKey {
		cmd.writeFieldValue(key.userKey, KEY)
	}
}

func (cmd *baseCommand) writeOperationForBin(bin *Bin, operation OperationType) error {
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

func (cmd *baseCommand) writeOperationForBinNameAndValue(name string, val interface{}, operation OperationType) error {
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

func (cmd *baseCommand) writeOperationForOperation(operation *Operation) error {
	nameLength := copy(cmd.dataBuffer[(cmd.dataOffset+int(_OPERATION_HEADER_SIZE)):], operation.binName)

	if operation.used {
		// cahce will set the used flag to false again
		operation.cache()
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

func (cmd *baseCommand) writePredExp(predExp []PredExp, predSize int) error {
	cmd.writeFieldHeader(predSize, FILTER_EXP)
	for i := range predExp {
		if err := predExp[i].marshal(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *baseCommand) writeFilterExpression(exp *FilterExpression, expSize int) error {
	cmd.writeFieldHeader(expSize, FILTER_EXP)
	if _, err := exp.pack(cmd); err != nil {
		return err
	}
	return nil
}

func (cmd *baseCommand) writeFieldValue(value Value, ftype FieldType) error {
	vlen, err := value.EstimateSize()
	if err != nil {
		return err
	}
	cmd.writeFieldHeader(vlen+1, ftype)
	cmd.WriteByte(byte(value.GetType()))

	_, err = value.write(cmd)
	return err
}

func (cmd *baseCommand) writeUdfArgs(value *ValueArray) error {
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
	len := copy(cmd.dataBuffer[(cmd.dataOffset+int(_FIELD_HEADER_SIZE)):], str)
	cmd.writeFieldHeader(len, ftype)
	cmd.dataOffset += len
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

// Int64ToBytes converts an int64 into slice of Bytes.
func (cmd *baseCommand) WriteInt64(num int64) (int, error) {
	return cmd.WriteUint64(uint64(num))
}

// Uint64ToBytes converts an uint64 into slice of Bytes.
func (cmd *baseCommand) WriteUint64(num uint64) (int, error) {
	binary.BigEndian.PutUint64(cmd.dataBuffer[cmd.dataOffset:cmd.dataOffset+8], num)
	cmd.dataOffset += 8
	return 8, nil
}

// Int32ToBytes converts an int32 to a byte slice of size 4
func (cmd *baseCommand) WriteInt32(num int32) (int, error) {
	return cmd.WriteUint32(uint32(num))
}

// Uint32ToBytes converts an uint32 to a byte slice of size 4
func (cmd *baseCommand) WriteUint32(num uint32) (int, error) {
	binary.BigEndian.PutUint32(cmd.dataBuffer[cmd.dataOffset:cmd.dataOffset+4], num)
	cmd.dataOffset += 4
	return 4, nil
}

// Uint32ToBytes converts an uint32 to a byte slice of size 4
func (cmd *baseCommand) WriteUint32At(num uint32, index int) (int, error) {
	binary.BigEndian.PutUint32(cmd.dataBuffer[index:index+4], num)
	return 4, nil
}

// Int16ToBytes converts an int16 to slice of bytes
func (cmd *baseCommand) WriteInt16(num int16) (int, error) {
	return cmd.WriteUint16(uint16(num))
}

// Int16ToBytes converts an int16 to slice of bytes
func (cmd *baseCommand) WriteUint16(num uint16) (int, error) {
	binary.BigEndian.PutUint16(cmd.dataBuffer[cmd.dataOffset:cmd.dataOffset+2], num)
	cmd.dataOffset += 2
	return 2, nil
}

func (cmd *baseCommand) WriteFloat32(float float32) (int, error) {
	bits := math.Float32bits(float)
	binary.BigEndian.PutUint32(cmd.dataBuffer[cmd.dataOffset:cmd.dataOffset+4], bits)
	cmd.dataOffset += 4
	return 4, nil
}

func (cmd *baseCommand) WriteFloat64(float float64) (int, error) {
	bits := math.Float64bits(float)
	binary.BigEndian.PutUint64(cmd.dataBuffer[cmd.dataOffset:cmd.dataOffset+8], bits)
	cmd.dataOffset += 8
	return 8, nil
}

func (cmd *baseCommand) WriteByte(b byte) error {
	cmd.dataBuffer[cmd.dataOffset] = b
	cmd.dataOffset++
	return nil
}

func (cmd *baseCommand) WriteString(s string) (int, error) {
	copy(cmd.dataBuffer[cmd.dataOffset:cmd.dataOffset+len(s)], s)
	cmd.dataOffset += len(s)
	return len(s), nil
}

func (cmd *baseCommand) Write(b []byte) (int, error) {
	copy(cmd.dataBuffer[cmd.dataOffset:cmd.dataOffset+len(b)], b)
	cmd.dataOffset += len(b)
	return len(b), nil
}

func (cmd *baseCommand) begin() {
	cmd.dataOffset = int(_MSG_TOTAL_HEADER_SIZE)
}

func (cmd *baseCommand) sizeBuffer(compress bool) error {
	return cmd.sizeBufferSz(cmd.dataOffset, compress)
}

func (cmd *baseCommand) validateHeader(header int64) error {
	msgVersion := (uint64(header) & 0xFF00000000000000) >> 56
	if msgVersion != 2 {
		return NewAerospikeError(PARSE_ERROR, fmt.Sprintf("Invalid Message Header: Expected version to be 2, but got %v", msgVersion))
	}

	msgType := (uint64(header) & 0x00FF000000000000) >> 49
	if !(msgType == 1 || msgType == 3 || msgType == 4) {
		return NewAerospikeError(PARSE_ERROR, fmt.Sprintf("Invalid Message Header: Expected type to be 1, 3 or 4, but got %v", msgType))
	}

	msgSize := header & 0x0000FFFFFFFFFFFF
	if msgSize > int64(MaxBufferSize) {
		return NewAerospikeError(PARSE_ERROR, fmt.Sprintf("Invalid Message Header: Expected size to be under 10MiB, but got %v", msgSize))
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

func (cmd *baseCommand) sizeBufferSz(size int, willCompress bool) error {

	if willCompress {
		// adds zlib and proto pads to the size of the buffer
		size += msgHeaderPad + zlibHeaderPad
	}

	// Corrupted data streams can result in a huge length.
	// Do a sanity check here.
	if size > MaxBufferSize || size < 0 {
		return NewAerospikeError(PARSE_ERROR, fmt.Sprintf("Invalid size for buffer: %d", size))
	}

	if size <= len(cmd.dataBuffer) {
		// don't touch the buffer
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

func (cmd *baseCommand) compress() error {
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
				return err
			}
		}

		if i < cmd.dataOffset {
			_, err := w.Write(cmd.dataBuffer[i:cmd.dataOffset])
			if err != nil {
				return err
			}
		}

		// flush
		w.Close()

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

func setInDoubt(err error, isRead bool, commandSentCounter int) error {
	// set inDoubt flag
	if ae, ok := err.(AerospikeError); ok {
		ae.SetInDoubt(isRead, commandSentCounter)
		return ae
	}

	return err
}

func (cmd *baseCommand) execute(ifc command, isRead bool) error {
	policy := ifc.getPolicy(ifc).GetBasePolicy()
	deadline := policy.deadline()

	return cmd.executeAt(ifc, policy, isRead, deadline, -1, 0)
}

func (cmd *baseCommand) executeAt(ifc command, policy *BasePolicy, isRead bool, deadline time.Time, iterations, commandSentCounter int) (err error) {
	// for exponential backoff
	interval := policy.SleepBetweenRetries

	notFirstIteration := false
	isClientTimeout := false

	// Execute command until successful, timed out or maximum iterations have been reached.
	for {
		iterations++

		// too many retries
		if (policy.MaxRetries <= 0 && iterations > 0) || (policy.MaxRetries > 0 && iterations > policy.MaxRetries) {
			if ae, ok := err.(AerospikeError); ok {
				err = NewAerospikeError(ae.ResultCode(), fmt.Sprintf("command execution timed out on client: Exceeded number of retries. See `Policy.MaxRetries`. (last error: %s)", err.Error()))
			}

			return setInDoubt(err, isRead, commandSentCounter)
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
			aerr, ok := err.(AerospikeError)
			if !ifc.prepareRetry(ifc, isClientTimeout || (ok && aerr.ResultCode() != SERVER_NOT_AVAILABLE)) {
				if bc, ok := ifc.(batcher); ok {
					// Batch may be retried in separate commands.
					if retry, err := bc.retryBatch(bc, cmd.node.cluster, deadline, iterations, commandSentCounter); retry {
						// Batch was retried in separate commands. Complete this command.
						return err
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
			isClientTimeout = true

			// Node is currently inactive. Retry.
			continue
		}

		// check if node has encountered too many errors
		if err := cmd.node.validateErrorCount(); err != nil {
			isClientTimeout = false

			// Max error rate achieved, try again per policy
			continue
		}

		cmd.conn, err = ifc.getConnection(policy)
		if err != nil {
			isClientTimeout = true

			if err == ErrConnectionPoolEmpty {
				// if the connection pool is empty, we still haven't tried
				// the transaction to increase the iteration count.
				iterations--
			}
			Logger.Debug("Node " + cmd.node.String() + ": " + err.Error())
			continue
		}

		// Assign the connection buffer to the command buffer
		cmd.dataBuffer = cmd.conn.dataBuffer

		// Set command buffer.
		err = ifc.writeBuffer(ifc)
		if err != nil {
			// All runtime exceptions are considered fatal. Do not retry.
			// Close socket to flush out possible garbage. Do not put back in pool.
			cmd.conn.Close()
			cmd.conn = nil
			return err
		}

		// Reset timeout in send buffer (destined for server) and socket.
		binary.BigEndian.PutUint32(cmd.dataBuffer[22:], 0)
		if !deadline.IsZero() {
			serverTimeout := deadline.Sub(time.Now())
			if serverTimeout < time.Millisecond {
				serverTimeout = time.Millisecond
			}
			binary.BigEndian.PutUint32(cmd.dataBuffer[22:], uint32(serverTimeout/time.Millisecond))
		}

		// now that the deadline has been set in the buffer, compress the contents
		if err = cmd.compress(); err != nil {
			return NewAerospikeError(SERIALIZE_ERROR, err.Error())
		}

		// Send command.
		_, err = cmd.conn.Write(cmd.dataBuffer[:cmd.dataOffset])
		if err != nil {
			isClientTimeout = true
			if deviceOverloadError(err) {
				isClientTimeout = false
				cmd.node.incrErrorCount()
			}

			// IO errors are considered temporary anomalies. Retry.
			// Close socket to flush out possible garbage. Do not put back in pool.
			cmd.conn.Close()
			cmd.conn = nil

			Logger.Debug("Node " + cmd.node.String() + ": " + err.Error())
			continue
		}
		commandSentCounter++

		// Parse results.
		err = ifc.parseResult(ifc, cmd.conn)
		if err != nil {
			if networkError(err) {
				isClientTimeout = (err == ErrTimeout)
				if err != ErrTimeout {
					if deviceOverloadError(err) {
						cmd.node.incrErrorCount()
					}
				}

				// IO errors are considered temporary anomalies. Retry.
				// Close socket to flush out possible garbage. Do not put back in pool.
				cmd.conn.Close()

				Logger.Debug("Node " + cmd.node.String() + ": " + err.Error())

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
			if cmd.conn.IsConnected() && KeepConnection(err) {
				// Put connection back in pool.
				cmd.node.PutConnection(cmd.conn)
			} else {
				cmd.conn.Close()
				cmd.conn = nil
			}

			return setInDoubt(err, isRead, commandSentCounter)
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

		// command has completed successfully.  Exit method.
		return nil

	}

	// execution timeout
	return ErrTimeout
}

func (cmd *baseCommand) parseRecordResults(ifc command, receiveSize int) (bool, error) {
	panic("Abstract method. Should not end up here")
}

func networkError(err error) bool {
	_, ok := err.(net.Error)
	return err == ErrTimeout || err == io.EOF || ok
}

func deviceOverloadError(err error) bool {
	aerr, ok := err.(AerospikeError)
	return ok && aerr.ResultCode() == DEVICE_OVERLOAD
}
