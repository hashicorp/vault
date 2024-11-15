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
	"fmt"
	"reflect"

	"github.com/aerospike/aerospike-client-go/v5/logger"
	"github.com/aerospike/aerospike-client-go/v5/types"

	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

type readCommand struct {
	singleCommand

	policy   *BasePolicy
	binNames []string
	record   *Record

	// pointer to the object that's going to be unmarshalled
	object *reflect.Value

	replicaSequence int
}

// this method uses reflection.
// Will not be set if performance flag is passed for the build.
var objectParser func(
	cmd *readCommand,
	opCount int,
	fieldCount int,
	generation uint32,
	expiration uint32,
) Error

func newReadCommand(cluster *Cluster, policy *BasePolicy, key *Key, binNames []string, partition *Partition) (readCommand, Error) {
	var err Error
	if partition == nil {
		partition, err = PartitionForRead(cluster, policy, key)
		if err != nil {
			return readCommand{}, err
		}
	}

	return readCommand{
		singleCommand: newSingleCommand(cluster, key, partition),
		binNames:      binNames,
		policy:        policy,
	}, nil
}

func (cmd *readCommand) getPolicy(ifc command) Policy {
	return cmd.policy
}

func (cmd *readCommand) writeBuffer(ifc command) Error {
	return cmd.setRead(cmd.policy, cmd.key, cmd.binNames)
}

func (cmd *readCommand) getNode(ifc command) (*Node, Error) {
	return cmd.partition.GetNodeRead(cmd.cluster)
}

func (cmd *readCommand) prepareRetry(ifc command, isTimeout bool) bool {
	cmd.partition.PrepareRetryRead(isTimeout)
	return true
}

func (cmd *readCommand) parseResult(ifc command, conn *Connection) Error {
	// Read proto and check if compressed
	if _, err := conn.Read(cmd.dataBuffer, 8); err != nil {
		logger.Logger.Debug("Connection error reading data for ReadCommand: %s", err.Error())
		return err
	}

	if compressedSize := cmd.compressedSize(); compressedSize > 0 {
		// Read compressed size
		if _, err := conn.Read(cmd.dataBuffer, 8); err != nil {
			logger.Logger.Debug("Connection error reading data for ReadCommand: %s", err.Error())
			return err
		}

		if err := cmd.conn.initInflater(true, compressedSize); err != nil {
			return newError(types.PARSE_ERROR, fmt.Sprintf("Error setting up zlib inflater for size `%d`: %s", compressedSize, err.Error()))
		}

		// Read header.
		if _, err := conn.Read(cmd.dataBuffer, int(_MSG_TOTAL_HEADER_SIZE)); err != nil {
			logger.Logger.Debug("Connection error reading data for ReadCommand: %s", err.Error())
			return err
		}
	} else {
		// Read header.
		if _, err := conn.Read(cmd.dataBuffer[8:], int(_MSG_TOTAL_HEADER_SIZE)-8); err != nil {
			logger.Logger.Debug("Connection error reading data for ReadCommand: %s", err.Error())
			return err
		}
	}

	// A number of these are commented out because we just don't care enough to read
	// that section of the header. If we do care, uncomment and check!
	sz := Buffer.BytesToInt64(cmd.dataBuffer, 0)

	// Validate header to make sure we are at the beginning of a message
	if err := cmd.validateHeader(sz); err != nil {
		return err
	}

	headerLength := int(cmd.dataBuffer[8])
	resultCode := types.ResultCode(cmd.dataBuffer[13] & 0xFF)
	generation := Buffer.BytesToUint32(cmd.dataBuffer, 14)
	expiration := types.TTL(Buffer.BytesToUint32(cmd.dataBuffer, 18))
	fieldCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 26)) // almost certainly 0
	opCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 28))
	receiveSize := int((sz & 0xFFFFFFFFFFFF) - int64(headerLength))

	// Read remaining message bytes.
	if receiveSize > 0 {
		if err := cmd.sizeBufferSz(receiveSize, false); err != nil {
			return err
		}
		if _, err := conn.Read(cmd.dataBuffer, receiveSize); err != nil {
			logger.Logger.Debug("Connection error reading data for ReadCommand: %s", err.Error())
			return err
		}

	}

	if resultCode != 0 {
		if resultCode == types.KEY_NOT_FOUND_ERROR {
			return ErrKeyNotFound.err()
		} else if resultCode == types.FILTERED_OUT {
			return ErrFilteredOut.err()
		} else if resultCode == types.UDF_BAD_RESPONSE {
			cmd.record, _ = cmd.parseRecord(ifc, opCount, fieldCount, generation, expiration)
			err := cmd.handleUdfError(resultCode)
			logger.Logger.Debug("UDF execution error: " + err.Error())
			return err
		}

		return newError(resultCode)
	}

	if cmd.object == nil {
		if opCount == 0 {
			// data Bin was not returned
			cmd.record = newRecord(cmd.node, cmd.key, nil, generation, expiration)
			return nil
		}

		var err Error
		cmd.record, err = cmd.parseRecord(ifc, opCount, fieldCount, generation, expiration)
		if err != nil {
			return err
		}
	} else if objectParser != nil {
		if err := objectParser(cmd, opCount, fieldCount, generation, expiration); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *readCommand) handleUdfError(resultCode types.ResultCode) Error {
	if ret, exists := cmd.record.Bins["FAILURE"]; exists {
		return newError(resultCode, ret.(string))
	}
	return newError(resultCode)
}

func (cmd *readCommand) parseRecord(
	ifc command,
	opCount int,
	fieldCount int,
	generation uint32,
	expiration uint32,
) (*Record, Error) {
	var bins BinMap
	receiveOffset := 0

	type opList []interface{}
	_, isOperate := ifc.(*operateCommand)
	var binNamesSet []string

	// There can be fields in the response (setname etc).
	// But for now, ignore them. Expose them to the API if needed in the future.
	//logger.Logger.Debug("field count: %d, databuffer: %v", fieldCount, cmd.dataBuffer)
	if fieldCount > 0 {
		// Just skip over all the fields
		for i := 0; i < fieldCount; i++ {
			//logger.Logger.Debug("%d", receiveOffset)
			fieldSize := int(Buffer.BytesToUint32(cmd.dataBuffer, receiveOffset))
			receiveOffset += (4 + fieldSize)
		}
	}

	if opCount > 0 {
		bins = make(BinMap, opCount)
	}

	for i := 0; i < opCount; i++ {
		opSize := int(Buffer.BytesToUint32(cmd.dataBuffer, receiveOffset))
		particleType := int(cmd.dataBuffer[receiveOffset+5])
		nameSize := int(cmd.dataBuffer[receiveOffset+7])
		name := string(cmd.dataBuffer[receiveOffset+8 : receiveOffset+8+nameSize])
		receiveOffset += 4 + 4 + nameSize

		particleBytesSize := opSize - (4 + nameSize)
		value, _ := bytesToParticle(particleType, cmd.dataBuffer, receiveOffset, particleBytesSize)
		receiveOffset += particleBytesSize

		if bins == nil {
			bins = make(BinMap, opCount)
		}

		if isOperate {
			// for operate list command results
			if prev, exists := bins[name]; exists {
				if res, ok := prev.(opList); ok {
					// List already exists.  Add to it.
					bins[name] = append(res, value)
				} else {
					// Make a list to store all values.
					bins[name] = opList{prev, value}
					binNamesSet = append(binNamesSet, name)
				}
			} else {
				bins[name] = value
			}
		} else {
			bins[name] = value
		}
	}

	if isOperate {
		for i := range binNamesSet {
			bins[binNamesSet[i]] = []interface{}(bins[binNamesSet[i]].(opList))
		}
	}

	return newRecord(cmd.node, cmd.key, bins, generation, expiration), nil
}

func (cmd *readCommand) GetRecord() *Record {
	return cmd.record
}

func (cmd *readCommand) Execute() Error {
	return cmd.execute(cmd, true)
}
