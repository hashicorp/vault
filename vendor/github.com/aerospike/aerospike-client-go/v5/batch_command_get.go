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
	"reflect"

	"github.com/aerospike/aerospike-client-go/v5/types"
	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

type batchCommandGet struct {
	batchCommand

	keys         []*Key
	binNames     []string     // binNames are mutually exclusive with ops
	ops          []*Operation // ops are mutually exclusive with binNames
	records      []*Record
	indexRecords []*BatchRead
	readAttr     int
	key          Key

	// pointer to the object that's going to be unmarshalled
	objects      []*reflect.Value
	objectsFound []bool
}

// this method uses reflection.
// Will not be set if performance flag is passed for the build.
var batchObjectParser func(
	cmd *batchCommandGet,
	offset int,
	opCount int,
	fieldCount int,
	generation uint32,
	expiration uint32,
) Error

func newBatchCommandGet(
	node *Node,
	batch *batchNode,
	policy *BatchPolicy,
	keys []*Key,
	binNames []string,
	ops []*Operation,
	records []*Record,
	readAttr int,
	isOperation bool,
) *batchCommandGet {
	res := &batchCommandGet{
		batchCommand: batchCommand{
			baseMultiCommand: *newMultiCommand(node, nil, isOperation),
			policy:           policy,
			batch:            batch,
		},
		keys:     keys,
		ops:      ops,
		binNames: binNames,
		records:  records,
		readAttr: readAttr,
	}
	return res
}

func (cmd *batchCommandGet) cloneBatchCommand(batch *batchNode) batcher {
	res := *cmd
	res.node = batch.Node
	res.batch = batch

	return &res
}

func (cmd *batchCommandGet) writeBuffer(ifc command) Error {
	return cmd.setBatchRead(cmd.policy, cmd.keys, cmd.batch, cmd.binNames, cmd.ops, cmd.readAttr)
}

// Parse all results in the batch.  Add records to shared list.
// If the record was not found, the bins will be nil.
func (cmd *batchCommandGet) parseRecordResults(ifc command, receiveSize int) (bool, Error) {
	//Parse each message response and add it to the result array
	cmd.dataOffset = 0

	for cmd.dataOffset < receiveSize {
		if err := cmd.readBytes(int(_MSG_REMAINING_HEADER_SIZE)); err != nil {
			return false, err
		}
		resultCode := types.ResultCode(cmd.dataBuffer[5] & 0xFF)

		// The only valid server return codes are "ok" and "not found" and "filtered out".
		// If other return codes are received, then abort the batch.
		if resultCode != 0 && resultCode != types.KEY_NOT_FOUND_ERROR {
			if resultCode == types.FILTERED_OUT {
				cmd.filteredOutCnt++
			} else {
				return false, newCustomNodeError(cmd.node, resultCode)
			}
		}

		info3 := int(cmd.dataBuffer[3])

		// If cmd is the end marker of the response, do not proceed further
		if (info3 & _INFO3_LAST) == _INFO3_LAST {
			return false, nil
		}

		generation := Buffer.BytesToUint32(cmd.dataBuffer, 6)
		expiration := types.TTL(Buffer.BytesToUint32(cmd.dataBuffer, 10))
		batchIndex := int(Buffer.BytesToUint32(cmd.dataBuffer, 14))
		fieldCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 18))
		opCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 20))
		err := cmd.skipKey(fieldCount)
		if err != nil {
			return false, err
		}

		if cmd.indexRecords != nil {
			if len(cmd.indexRecords) > 0 {
				if resultCode == 0 {
					if cmd.indexRecords[batchIndex].Record, err = cmd.parseRecord(cmd.indexRecords[batchIndex].Key, opCount, generation, expiration); err != nil {
						return false, err
					}
				}
			}
		} else {
			if resultCode == 0 {
				if cmd.objects == nil {
					if cmd.records[batchIndex], err = cmd.parseRecord(cmd.keys[batchIndex], opCount, generation, expiration); err != nil {
						return false, err
					}
				} else if batchObjectParser != nil {
					// mark it as found
					cmd.objectsFound[batchIndex] = true
					if err := batchObjectParser(cmd, batchIndex, opCount, fieldCount, generation, expiration); err != nil {
						return false, err

					}
				}
			}
		}
	}
	return true, nil
}

// Parses the given byte buffer and populate the result object.
// Returns the number of bytes that were parsed from the given buffer.
func (cmd *batchCommandGet) parseRecord(key *Key, opCount int, generation, expiration uint32) (*Record, Error) {
	bins := make(BinMap, opCount)

	for i := 0; i < opCount; i++ {
		if err := cmd.readBytes(8); err != nil {
			return nil, err
		}
		opSize := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
		particleType := int(cmd.dataBuffer[5])
		nameSize := int(cmd.dataBuffer[7])

		if err := cmd.readBytes(nameSize); err != nil {
			return nil, err
		}
		name := string(cmd.dataBuffer[:nameSize])

		particleBytesSize := opSize - (4 + nameSize)
		if err := cmd.readBytes(particleBytesSize); err != nil {
			return nil, err
		}
		value, err := bytesToParticle(particleType, cmd.dataBuffer, 0, particleBytesSize)
		if err != nil {
			return nil, err
		}

		if cmd.isOperation {
			if prev, ok := bins[name]; ok {
				if prev2, ok := prev.(OpResults); ok {
					bins[name] = append(prev2, value)
				} else {
					bins[name] = OpResults{prev, value}
				}
			} else {
				bins[name] = value
			}
		} else {
			bins[name] = value
		}
	}

	return newRecord(cmd.node, key, bins, generation, expiration), nil
}

func (cmd *batchCommandGet) Execute() Error {
	return cmd.execute(cmd, true)
}

func (cmd *batchCommandGet) generateBatchNodes(cluster *Cluster) ([]*batchNode, Error) {
	return newBatchNodeListKeys(cluster, cmd.policy, cmd.keys, cmd.sequenceAP, cmd.sequenceSC, cmd.batch)
}
