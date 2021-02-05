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
	"reflect"

	. "github.com/aerospike/aerospike-client-go/types"
	Buffer "github.com/aerospike/aerospike-client-go/utils/buffer"
)

type batchCommandGet struct {
	batchCommand

	keys         []*Key
	binNames     []string
	records      []*Record
	indexRecords []*BatchRead
	readAttr     int
	index        int
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
) error

func newBatchCommandGet(
	node *Node,
	batch *batchNode,
	policy *BatchPolicy,
	keys []*Key,
	binNames []string,
	records []*Record,
	readAttr int,
) *batchCommandGet {
	res := &batchCommandGet{
		batchCommand: batchCommand{
			baseMultiCommand: *newMultiCommand(node, nil),
			policy:           policy,
			batch:            batch,
		},
		keys:     keys,
		binNames: binNames,
		records:  records,
		readAttr: readAttr,
	}
	res.oneShot = false
	return res
}

func (cmd *batchCommandGet) cloneBatchCommand(batch *batchNode) batcher {
	res := *cmd
	res.node = batch.Node
	res.batch = batch

	return &res
}

func (cmd *batchCommandGet) writeBuffer(ifc command) error {
	return cmd.setBatchIndexReadCompat(cmd.policy, cmd.keys, cmd.batch, cmd.binNames, cmd.readAttr)
}

// On batch operations the key values are not returned from the server
// So we reuse the Key on the batch Object
func (cmd *batchCommandGet) parseKey(fieldCount int) error {
	// var digest [20]byte
	// var namespace, setName string
	// var userKey Value
	var err error

	for i := 0; i < fieldCount; i++ {
		if err = cmd.readBytes(4); err != nil {
			return err
		}

		fieldlen := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
		if err = cmd.readBytes(fieldlen); err != nil {
			return err
		}

		fieldtype := FieldType(cmd.dataBuffer[0])
		size := fieldlen - 1

		switch fieldtype {
		case DIGEST_RIPE:
			copy(cmd.key.digest[:], cmd.dataBuffer[1:size+1])
		case NAMESPACE:
			cmd.key.namespace = string(cmd.dataBuffer[1 : size+1])
		case TABLE:
			cmd.key.setName = string(cmd.dataBuffer[1 : size+1])
		case KEY:
			if cmd.key.userKey, err = bytesToKeyValue(int(cmd.dataBuffer[1]), cmd.dataBuffer, 2, size-1); err != nil {
				return err
			}
		}
	}

	return nil
}

// Parse all results in the batch.  Add records to shared list.
// If the record was not found, the bins will be nil.
func (cmd *batchCommandGet) parseRecordResults(ifc command, receiveSize int) (bool, error) {
	//Parse each message response and add it to the result array
	cmd.dataOffset = 0

	for cmd.dataOffset < receiveSize {
		if err := cmd.readBytes(int(_MSG_REMAINING_HEADER_SIZE)); err != nil {
			return false, err
		}
		resultCode := ResultCode(cmd.dataBuffer[5] & 0xFF)

		// The only valid server return codes are "ok" and "not found" and "filtered out".
		// If other return codes are received, then abort the batch.
		if resultCode != 0 && resultCode != KEY_NOT_FOUND_ERROR {
			if resultCode == FILTERED_OUT {
				cmd.filteredOutCnt++
			} else {
				return false, NewAerospikeError(resultCode)
			}
		}

		info3 := int(cmd.dataBuffer[3])

		// If cmd is the end marker of the response, do not proceed further
		if (info3 & _INFO3_LAST) == _INFO3_LAST {
			return false, nil
		}

		generation := Buffer.BytesToUint32(cmd.dataBuffer, 6)
		expiration := TTL(Buffer.BytesToUint32(cmd.dataBuffer, 10))
		batchIndex := int(Buffer.BytesToUint32(cmd.dataBuffer, 14))
		fieldCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 18))
		opCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 20))
		err := cmd.parseKey(fieldCount)
		if err != nil {
			return false, err
		}

		var offset int
		offset = batchIndex

		if cmd.indexRecords != nil {
			if len(cmd.indexRecords) > 0 {
				if resultCode == 0 {
					if cmd.indexRecords[offset].Record, err = cmd.parseRecord(cmd.indexRecords[offset].Key, opCount, generation, expiration); err != nil {
						return false, err
					}
				}
			}
		} else {
			if resultCode == 0 {
				if cmd.objects == nil {
					if cmd.records[offset], err = cmd.parseRecord(cmd.keys[offset], opCount, generation, expiration); err != nil {
						return false, err
					}
				} else if batchObjectParser != nil {
					// mark it as found
					cmd.objectsFound[offset] = true
					if err := batchObjectParser(cmd, offset, opCount, fieldCount, generation, expiration); err != nil {
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
func (cmd *batchCommandGet) parseRecord(key *Key, opCount int, generation, expiration uint32) (*Record, error) {
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

		bins[name] = value
	}

	return newRecord(cmd.node, key, bins, generation, expiration), nil
}

func (cmd *batchCommandGet) Execute() error {
	return cmd.execute(cmd, true)
}

func (cmd *batchCommandGet) generateBatchNodes(cluster *Cluster) ([]*batchNode, error) {
	return newBatchNodeListKeys(cluster, cmd.policy, cmd.keys, cmd.sequenceAP, cmd.sequenceSC, cmd.batch)
}
