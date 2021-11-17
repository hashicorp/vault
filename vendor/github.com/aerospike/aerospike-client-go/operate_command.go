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

type operateCommand struct {
	readCommand

	policy     *WritePolicy
	operations []*Operation

	hasWrite bool
}

func newOperateCommand(cluster *Cluster, policy *WritePolicy, key *Key, operations []*Operation) (operateCommand, error) {
	hasWrite := hasWriteOp(operations)

	var partition *Partition
	var err error
	if hasWrite {
		partition, err = PartitionForWrite(cluster, &policy.BasePolicy, key)
	} else {
		partition, err = PartitionForRead(cluster, &policy.BasePolicy, key)
	}

	if err != nil {
		return operateCommand{}, err
	}

	readCommand, err := newReadCommand(cluster, &policy.BasePolicy, key, nil, partition)
	if err != nil {
		return operateCommand{}, err
	}

	return operateCommand{
		readCommand: readCommand,
		policy:      policy,
		operations:  operations,

		hasWrite: hasWrite,
	}, nil
}

func (cmd *operateCommand) writeBuffer(ifc command) (err error) {
	cmd.hasWrite, err = cmd.setOperate(cmd.policy, cmd.key, cmd.operations)
	return err
}

func (cmd *operateCommand) getNode(ifc command) (*Node, error) {
	if cmd.hasWrite {
		return cmd.partition.GetNodeWrite(cmd.cluster)
	}

	// this may be affected by Rackaware
	return cmd.partition.GetNodeRead(cmd.cluster)
}

func (cmd *operateCommand) prepareRetry(ifc command, isTimeout bool) bool {
	if cmd.hasWrite {
		cmd.partition.PrepareRetryWrite(isTimeout)
	} else {
		cmd.partition.PrepareRetryRead(isTimeout)
	}
	return true
}

func (cmd *operateCommand) Execute() error {
	return cmd.execute(cmd, !cmd.hasWrite)
}

func hasWriteOp(operations []*Operation) bool {
	for i := range operations {
		switch operations[i].opType {
		case _MAP_READ, _READ, _CDT_READ:
		default:
			// All other cases are a type of write
			return true
		}
	}

	return false
}
