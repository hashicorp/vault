// +build !app_engine

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

	"github.com/aerospike/aerospike-client-go/v5/types"

	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
	lua "github.com/yuin/gopher-lua"
)

type queryAggregateCommand struct {
	queryCommand

	luaInstance *lua.LState
	inputChan   chan interface{}
}

func newQueryAggregateCommand(node *Node, policy *QueryPolicy, statement *Statement, recordset *Recordset) *queryAggregateCommand {
	cmd := &queryAggregateCommand{
		queryCommand: *newQueryCommand(node, policy, nil, statement, nil, recordset),
	}

	cmd.terminationErrorType = types.QUERY_TERMINATED

	return cmd
}

func (cmd *queryAggregateCommand) Execute() Error {
	cmd.policy.MaxRetries = 0
	err := cmd.execute(cmd, true)
	if err != nil {
		cmd.recordset.sendError(err)
	}
	return err
}

func (cmd *queryAggregateCommand) parseRecordResults(ifc command, receiveSize int) (bool, Error) {
	// Read/parse remaining message bytes one record at a time.
	cmd.dataOffset = 0

	for cmd.dataOffset < receiveSize {
		if err := cmd.readBytes(int(_MSG_REMAINING_HEADER_SIZE)); err != nil {
			err = newNodeError(cmd.node, err)
			return false, err
		}
		resultCode := types.ResultCode(cmd.dataBuffer[5] & 0xFF)

		if resultCode != 0 {
			if resultCode == types.KEY_NOT_FOUND_ERROR {
				// consume the rest of the input buffer from the socket
				if cmd.dataOffset < receiveSize {
					if err := cmd.readBytes(receiveSize - cmd.dataOffset); err != nil {
						err = newNodeError(cmd.node, err)
						return false, err
					}
				}
				return false, nil
			}
			return false, newCustomNodeError(cmd.node, resultCode)
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
			return false, newCustomNodeError(cmd.node, types.PARSE_ERROR, fmt.Sprintf("Query aggregate command expects exactly only one bin. Received: %d", opCount))
		}

		if _, err := cmd.parseKey(fieldCount); err != nil {
			return false, newNodeError(cmd.node, err)
		}

		// if there is a recordset, process the record traditionally
		// otherwise, it is supposed to be a record channel

		// Parse bins.
		var bins BinMap

		for i := 0; i < opCount; i++ {
			if err := cmd.readBytes(8); err != nil {
				return false, newNodeError(cmd.node, err)
			}

			opSize := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
			particleType := int(cmd.dataBuffer[5])
			nameSize := int(cmd.dataBuffer[7])

			if err := cmd.readBytes(nameSize); err != nil {
				return false, newNodeError(cmd.node, err)
			}
			name := string(cmd.dataBuffer[:nameSize])

			particleBytesSize := opSize - (4 + nameSize)
			if err := cmd.readBytes(particleBytesSize); err != nil {
				return false, newNodeError(cmd.node, err)
			}
			value, err := bytesToParticle(particleType, cmd.dataBuffer, 0, particleBytesSize)
			if err != nil {
				return false, newNodeError(cmd.node, err)
			}

			if bins == nil {
				bins = make(BinMap, opCount)
			}
			bins[name] = value
		}

		recs, exists := bins["SUCCESS"]
		if !exists {
			if errStr, exists := bins["FAILURE"]; exists {
				return false, newError(types.QUERY_GENERIC, errStr.(string))
			}

			return false, newError(types.QUERY_GENERIC, fmt.Sprintf("QueryAggregate's expected result was not returned. Received: %v", bins))
		}

		// If the channel is full and it blocks, we don't want this command to
		// block forever, or panic in case the channel is closed in the meantime.
		select {
		// send back the result on the async channel
		case cmd.inputChan <- recs:
		case <-cmd.recordset.cancelled:
			return false, newError(types.QUERY_TERMINATED)
		}
	}

	return true, nil
}
