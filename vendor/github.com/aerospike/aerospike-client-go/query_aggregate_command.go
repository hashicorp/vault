// +build !app_engine

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
	"fmt"

	. "github.com/aerospike/aerospike-client-go/types"
	Buffer "github.com/aerospike/aerospike-client-go/utils/buffer"
	lua "github.com/yuin/gopher-lua"
)

type queryAggregateCommand struct {
	queryCommand

	luaInstance *lua.LState
	inputChan   chan interface{}
}

func newQueryAggregateCommand(node *Node, policy *QueryPolicy, statement *Statement, recordset *Recordset, clusterKey int64, first bool) *queryAggregateCommand {
	cmd := &queryAggregateCommand{
		queryCommand: *newQueryCommand(node, policy, nil, statement, nil, recordset, clusterKey, first),
	}

	cmd.terminationErrorType = QUERY_TERMINATED

	return cmd
}

func (cmd *queryAggregateCommand) Execute() error {
	// defer cmd.recordset.signalEnd()
	err := cmd.execute(cmd, true)
	if err != nil {
		cmd.recordset.sendError(err)
	}
	return err
}

func (cmd *queryAggregateCommand) parseRecordResults(ifc command, receiveSize int) (bool, error) {
	// Read/parse remaining message bytes one record at a time.
	cmd.dataOffset = 0

	for cmd.dataOffset < receiveSize {
		if err := cmd.readBytes(int(_MSG_REMAINING_HEADER_SIZE)); err != nil {
			err = newNodeError(cmd.node, err)
			return false, err
		}
		resultCode := ResultCode(cmd.dataBuffer[5] & 0xFF)

		if resultCode != 0 {
			if resultCode == KEY_NOT_FOUND_ERROR {
				// consume the rest of the input buffer from the socket
				if cmd.dataOffset < receiveSize {
					if err := cmd.readBytes(receiveSize - cmd.dataOffset); err != nil {
						err = newNodeError(cmd.node, err)
						return false, err
					}
				}
				return false, nil
			}
			err := NewAerospikeError(resultCode)
			err = newNodeError(cmd.node, err)
			return false, err
		}

		info3 := int(cmd.dataBuffer[3])

		// If cmd is the end marker of the response, do not proceed further
		if (info3 & _INFO3_LAST) == _INFO3_LAST {
			return false, nil
		}

		// generation := Buffer.BytesToUint32(cmd.dataBuffer, 6)
		// expiration := TTL(Buffer.BytesToUint32(cmd.dataBuffer, 10))
		fieldCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 18))
		opCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 20))

		if opCount != 1 {
			err := fmt.Errorf("Query aggregate command expects exactly only one bin. Received: %d", opCount)
			err = newNodeError(cmd.node, err)
			return false, err
		}

		_, err := cmd.parseKey(fieldCount)
		if err != nil {
			err = newNodeError(cmd.node, err)
			return false, err
		}

		// if there is a recordset, process the record traditionally
		// otherwise, it is supposed to be a record channel

		// Parse bins.
		var bins BinMap

		for i := 0; i < opCount; i++ {
			if err := cmd.readBytes(8); err != nil {
				err = newNodeError(cmd.node, err)
				return false, err
			}

			opSize := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
			particleType := int(cmd.dataBuffer[5])
			nameSize := int(cmd.dataBuffer[7])

			if err := cmd.readBytes(nameSize); err != nil {
				err = newNodeError(cmd.node, err)
				return false, err
			}
			name := string(cmd.dataBuffer[:nameSize])

			particleBytesSize := opSize - (4 + nameSize)
			if err = cmd.readBytes(particleBytesSize); err != nil {
				err = newNodeError(cmd.node, err)
				return false, err
			}
			value, err := bytesToParticle(particleType, cmd.dataBuffer, 0, particleBytesSize)
			if err != nil {
				err = newNodeError(cmd.node, err)
				return false, err
			}

			if bins == nil {
				bins = make(BinMap, opCount)
			}
			bins[name] = value
		}

		recs, exists := bins["SUCCESS"]
		if !exists {
			if errStr, exists := bins["FAILURE"]; exists {
				err = NewAerospikeError(QUERY_GENERIC, errStr.(string))
				return false, err
			}

			err = NewAerospikeError(QUERY_GENERIC, fmt.Sprintf("QueryAggregate's expected result was not returned. Received: %v", bins))
			return false, err
		}

		// If the channel is full and it blocks, we don't want this command to
		// block forever, or panic in case the channel is closed in the meantime.
		select {
		// send back the result on the async channel
		case cmd.inputChan <- recs:
		case <-cmd.recordset.cancelled:
			return false, NewAerospikeError(QUERY_TERMINATED)
		}
	}

	return true, nil
}
