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

	"github.com/aerospike/aerospike-client-go/v5/types"
	xrand "github.com/aerospike/aerospike-client-go/v5/types/rand"
	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

type baseMultiCommand struct {
	baseCommand

	namespace string
	recordset *Recordset

	isOperation bool

	// Used in correct Scans/Queries
	tracker        *partitionTracker
	nodePartitions *nodePartitions

	terminationErrorType types.ResultCode

	resObjType     reflect.Type
	resObjMappings map[string][]int
	selectCases    []reflect.SelectCase

	bc bufferedConn
}

var multiObjectParser func(
	cmd *baseMultiCommand,
	obj reflect.Value,
	opCount int,
	fieldCount int,
	generation uint32,
	expiration uint32,
) Error

var prepareReflectionData func(cmd *baseMultiCommand)

func newMultiCommand(node *Node, recordset *Recordset, isOperation bool) *baseMultiCommand {
	cmd := &baseMultiCommand{
		baseCommand: baseCommand{
			node: node,
		},
		recordset:   recordset,
		isOperation: isOperation,
	}

	if prepareReflectionData != nil {
		prepareReflectionData(cmd)
	}
	return cmd
}

func newStreamingMultiCommand(node *Node, recordset *Recordset, namespace string, isOperation bool) *baseMultiCommand {
	cmd := &baseMultiCommand{
		baseCommand: baseCommand{
			node:    node,
			oneShot: true,
		},
		namespace:   namespace,
		recordset:   recordset,
		isOperation: isOperation,
	}

	if prepareReflectionData != nil {
		prepareReflectionData(cmd)
	}
	return cmd
}

func newCorrectStreamingMultiCommand(recordset *Recordset, namespace string) *baseMultiCommand {
	cmd := &baseMultiCommand{
		baseCommand: baseCommand{
			oneShot: true,
		},
		namespace: namespace,
		recordset: recordset,
	}

	if prepareReflectionData != nil {
		prepareReflectionData(cmd)
	}
	return cmd
}

func (cmd *baseMultiCommand) getNode(ifc command) (*Node, Error) {
	return cmd.node, nil
}

func (cmd *baseMultiCommand) prepareRetry(ifc command, isTimeout bool) bool {
	return false
}

func (cmd *baseMultiCommand) getConnection(policy Policy) (*Connection, Error) {
	return cmd.node.getConnectionWithHint(policy.GetBasePolicy().deadline(), policy.GetBasePolicy().socketTimeout(), byte(xrand.Int64()%256))
}

func (cmd *baseMultiCommand) putConnection(conn *Connection) {
	cmd.node.putConnectionWithHint(conn, byte(xrand.Int64()%256))
}

func (cmd *baseMultiCommand) parseResult(ifc command, conn *Connection) Error {
	// Read socket into receive buffer one record at a time.  Do not read entire receive size
	// because the receive buffer would be too big.
	status := true

	var err Error

	cmd.bc = newBufferedConn(conn, 0)
	for status {
		if err = cmd.conn.initInflater(false, 0); err != nil {
			return newError(types.PARSE_ERROR, "Error setting up zlib inflater:", err.Error())
		}
		cmd.bc.reset(8)

		// Read header.
		if cmd.dataBuffer, err = cmd.bc.read(8); err != nil {
			return err
		}

		proto := Buffer.BytesToInt64(cmd.dataBuffer, 0)
		receiveSize := int(proto & 0xFFFFFFFFFFFF)
		if receiveSize <= 0 {
			continue
		}

		if compressedSize := cmd.compressedSize(); compressedSize > 0 {
			cmd.bc.reset(8)
			// Read header.
			if cmd.dataBuffer, err = cmd.bc.read(8); err != nil {
				return err
			}

			receiveSize = int(Buffer.BytesToInt64(cmd.dataBuffer, 0)) - 8
			if err = cmd.conn.initInflater(true, compressedSize-8); err != nil {
				return newError(types.PARSE_ERROR, fmt.Sprintf("Error setting up zlib inflater for size `%d`: %s", compressedSize-8, err.Error()))
			}

			// read the first 8 bytes
			cmd.bc.reset(8)
			if cmd.dataBuffer, err = cmd.bc.read(8); err != nil {
				return err
			}
		}

		// Validate header to make sure we are at the beginning of a message
		proto = Buffer.BytesToInt64(cmd.dataBuffer, 0)
		if err = cmd.validateHeader(proto); err != nil {
			return err
		}

		if receiveSize > 0 {
			cmd.bc.reset(receiveSize)

			status, err = ifc.parseRecordResults(ifc, receiveSize)
			if err != nil {
				cmd.bc.drainConn()
				return err
			}
		} else {
			status = false
		}
	}

	// if the buffer has been resized, put it back so that it will be reassigned to the connection.
	cmd.dataBuffer = cmd.bc.buf()

	return nil
}

func (cmd *baseMultiCommand) parseKey(fieldCount int) (*Key, Error) {
	var digest [20]byte
	var namespace, setName string
	var userKey Value
	var err Error

	for i := 0; i < fieldCount; i++ {
		if err = cmd.readBytes(4); err != nil {
			return nil, err
		}

		fieldlen := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
		if err = cmd.readBytes(fieldlen); err != nil {
			return nil, err
		}

		fieldtype := FieldType(cmd.dataBuffer[0])
		size := fieldlen - 1

		switch fieldtype {
		case DIGEST_RIPE:
			copy(digest[:], cmd.dataBuffer[1:size+1])
		case NAMESPACE:
			namespace = string(cmd.dataBuffer[1 : size+1])
		case TABLE:
			setName = string(cmd.dataBuffer[1 : size+1])
		case KEY:
			if userKey, err = bytesToKeyValue(int(cmd.dataBuffer[1]), cmd.dataBuffer, 2, size-1); err != nil {
				return nil, err
			}
		}
	}

	return &Key{namespace: namespace, setName: setName, digest: digest, userKey: userKey}, nil
}

func (cmd *baseMultiCommand) skipKey(fieldCount int) (err Error) {
	for i := 0; i < fieldCount; i++ {
		if err = cmd.readBytes(4); err != nil {
			return err
		}

		fieldlen := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
		if err = cmd.readBytes(fieldlen); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *baseMultiCommand) readBytes(length int) (err Error) {
	// Corrupted data streams can result in a huge length.
	// Do a sanity check here.
	if length > MaxBufferSize || length < 0 {
		return newError(types.PARSE_ERROR, fmt.Sprintf("Invalid readBytes length: %d", length))
	}

	cmd.dataBuffer, err = cmd.bc.read(length)
	if err != nil {
		return err
	}
	cmd.dataOffset += length

	return nil
}

func (cmd *baseMultiCommand) parseRecordResults(ifc command, receiveSize int) (bool, Error) {
	// Read/parse remaining message bytes one record at a time.
	cmd.dataOffset = 0

	for cmd.dataOffset < receiveSize {
		if err := cmd.readBytes(int(_MSG_REMAINING_HEADER_SIZE)); err != nil {
			err = newNodeError(cmd.node, err)
			return false, err
		}
		resultCode := types.ResultCode(cmd.dataBuffer[5] & 0xFF)

		if resultCode != 0 && resultCode != types.PARTITION_UNAVAILABLE {
			if resultCode == types.KEY_NOT_FOUND_ERROR || resultCode == types.FILTERED_OUT {
				return false, nil
			}
			err := newError(resultCode)
			err = newNodeError(cmd.node, err)
			return false, err
		}

		info3 := int(cmd.dataBuffer[3])

		// If cmd is the end marker of the response, do not proceed further
		if (info3 & _INFO3_LAST) == _INFO3_LAST {
			return false, nil
		}

		generation := Buffer.BytesToUint32(cmd.dataBuffer, 6)
		expiration := types.TTL(Buffer.BytesToUint32(cmd.dataBuffer, 10))
		fieldCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 18))
		opCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 20))

		key, err := cmd.parseKey(fieldCount)
		if err != nil {
			err = newNodeError(cmd.node, err)
			return false, err
		}

		// Partition is done, don't go further
		if (info3 & _INFO3_PARTITION_DONE) != 0 {
			// When nodes are removed, the tracker is getting set to nil
			if resultCode == 0 && cmd.tracker != nil {
				cmd.tracker.partitionDone(cmd.nodePartitions, int(generation))
			}
			continue
		}

		// if there is a recordset, process the record traditionally
		// otherwise, it is supposed to be a record channel
		if cmd.selectCases == nil {
			// Parse bins.
			var bins BinMap

			for i := 0; i < opCount; i++ {
				if err = cmd.readBytes(8); err != nil {
					return false, newNodeError(cmd.node, err)
				}

				opSize := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
				particleType := int(cmd.dataBuffer[5])
				nameSize := int(cmd.dataBuffer[7])

				if err = cmd.readBytes(nameSize); err != nil {
					return false, newNodeError(cmd.node, err)
				}
				name := string(cmd.dataBuffer[:nameSize])

				particleBytesSize := opSize - (4 + nameSize)
				if err = cmd.readBytes(particleBytesSize); err != nil {
					return false, newNodeError(cmd.node, err)
				}
				value, err := bytesToParticle(particleType, cmd.dataBuffer, 0, particleBytesSize)
				if err != nil {
					return false, newNodeError(cmd.node, err)
				}

				if bins == nil {
					bins = make(BinMap, opCount)
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

			// If the channel is full and it blocks, we don't want this command to
			// block forever, or panic in case the channel is closed in the meantime.
			select {
			// send back the result on the async channel
			case cmd.recordset.records <- &Result{Record: newRecord(cmd.node, key, bins, generation, expiration), Err: nil}:
			case <-cmd.recordset.cancelled:
				switch cmd.terminationErrorType {
				case types.SCAN_TERMINATED:
					return false, ErrScanTerminated.err()
				case types.QUERY_TERMINATED:
					return false, ErrQueryTerminated.err()
				default:
					return false, newError(cmd.terminationErrorType)
				}
			}
		} else if multiObjectParser != nil {
			obj := reflect.New(cmd.resObjType)
			if err := multiObjectParser(cmd, obj, opCount, fieldCount, generation, expiration); err != nil {
				err = newNodeError(cmd.node, err)
				return false, err
			}

			// set the object to send
			cmd.selectCases[0].Send = obj

			chosen, _, _ := reflect.Select(cmd.selectCases)
			switch chosen {
			case 0: // object sent
			case 1: // cancel channel is closed
				return false, newError(cmd.terminationErrorType)
			}
		}

		if cmd.tracker != nil {
			cmd.tracker.setDigest(cmd.nodePartitions, key)
		}
	}

	return true, nil
}

func (cmd *baseMultiCommand) canPutConnBack() bool {
	return false
}

func (cmd *baseMultiCommand) execute(ifc command, isRead bool) Error {

	/***************************************************************************
	IMPORTANT: 	No need to send the error here to the recordset.Error channel.
				It is being sent from the downstream command from the result
				returned from the function.
	****************************************************************************/

	return cmd.baseCommand.execute(ifc, isRead)
}
