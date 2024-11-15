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

type readHeaderCommand struct {
	singleCommand

	policy *BasePolicy
	record *Record
}

func newReadHeaderCommand(cluster *Cluster, policy *BasePolicy, key *Key) (readHeaderCommand, Error) {
	partition, err := PartitionForRead(cluster, policy, key)
	if err != nil {
		return readHeaderCommand{}, err
	}

	newReadHeaderCmd := readHeaderCommand{
		singleCommand: newSingleCommand(cluster, key, partition),
		policy:        policy,
	}

	return newReadHeaderCmd, nil
}

func (cmd *readHeaderCommand) getPolicy(ifc command) Policy {
	return cmd.policy
}

func (cmd *readHeaderCommand) writeBuffer(ifc command) Error {
	return cmd.setReadHeader(cmd.policy, cmd.key)
}

func (cmd *readHeaderCommand) getNode(ifc command) (*Node, Error) {
	return cmd.partition.GetNodeRead(cmd.cluster)
}

func (cmd *readHeaderCommand) prepareRetry(ifc command, isTimeout bool) bool {
	cmd.partition.PrepareRetryRead(isTimeout)
	return true
}

func (cmd *readHeaderCommand) parseResult(ifc command, conn *Connection) Error {
	// Read header.
	if _, err := conn.Read(cmd.dataBuffer, int(_MSG_TOTAL_HEADER_SIZE)); err != nil {
		return err
	}

	header := Buffer.BytesToInt64(cmd.dataBuffer, 0)

	// Validate header to make sure we are at the beginning of a message
	if err := cmd.validateHeader(header); err != nil {
		return err
	}

	resultCode := cmd.dataBuffer[13] & 0xFF

	if resultCode == 0 {
		generation := Buffer.BytesToUint32(cmd.dataBuffer, 14)
		expiration := types.TTL(Buffer.BytesToUint32(cmd.dataBuffer, 18))
		cmd.record = newRecord(cmd.node, cmd.key, nil, generation, expiration)
	} else {
		if types.ResultCode(resultCode) == types.KEY_NOT_FOUND_ERROR {
			cmd.record = nil
		} else if types.ResultCode(resultCode) == types.FILTERED_OUT {
			return ErrFilteredOut.err()
		} else {
			return newError(types.ResultCode(resultCode))
		}
	}
	return cmd.emptySocket(conn)
}

func (cmd *readHeaderCommand) GetRecord() *Record {
	return cmd.record
}

func (cmd *readHeaderCommand) Execute() Error {
	return cmd.execute(cmd, true)
}
