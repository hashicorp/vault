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
	"github.com/aerospike/aerospike-client-go/v5/types"
	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

type serverCommand struct {
	queryCommand
}

func newServerCommand(node *Node, policy *QueryPolicy, writePolicy *WritePolicy, statement *Statement, taskId uint64, operations []*Operation) *serverCommand {
	return &serverCommand{
		queryCommand: *newQueryCommand(node, policy, writePolicy, statement, operations, nil),
	}
}

func (cmd *serverCommand) writeBuffer(ifc command) (err Error) {
	return cmd.setQuery(cmd.policy, cmd.writePolicy, cmd.statement, cmd.statement.TaskId, cmd.operations, cmd.writePolicy != nil, nil)
}

func (cmd *serverCommand) parseRecordResults(ifc command, receiveSize int) (bool, Error) {
	// Server commands (Query/Execute UDF) should only send back a return code.
	// Keep parsing logic to empty socket buffer just in case server does
	// send records back.
	cmd.dataOffset = 0

	for cmd.dataOffset < receiveSize {
		if err := cmd.readBytes(int(_MSG_REMAINING_HEADER_SIZE)); err != nil {
			return false, err
		}
		resultCode := types.ResultCode(cmd.dataBuffer[5] & 0xFF)

		if resultCode != 0 {
			if resultCode == types.KEY_NOT_FOUND_ERROR {
				return false, nil
			}
			return false, newError(resultCode)
		}

		info3 := int(cmd.dataBuffer[3])

		// If cmd is the end marker of the response, do not proceed further
		if (info3 & _INFO3_LAST) == _INFO3_LAST {
			return false, nil
		}

		fieldCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 18))
		opCount := int(Buffer.BytesToUint16(cmd.dataBuffer, 20))

		if _, err := cmd.parseKey(fieldCount); err != nil {
			return false, err
		}

		for i := 0; i < opCount; i++ {
			if err := cmd.readBytes(8); err != nil {
				return false, err
			}
			opSize := int(Buffer.BytesToUint32(cmd.dataBuffer, 0))
			nameSize := int(cmd.dataBuffer[7])

			if err := cmd.readBytes(nameSize); err != nil {
				return false, err
			}

			particleBytesSize := opSize - (4 + nameSize)
			if err := cmd.readBytes(particleBytesSize); err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func (cmd *serverCommand) Execute() Error {
	return cmd.execute(cmd, false)
}
