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

import "github.com/aerospike/aerospike-client-go/v5/types"

type scanObjectsCommand struct {
	baseMultiCommand

	policy    *ScanPolicy
	namespace string
	setName   string
	binNames  []string
	taskID    uint64
}

func newScanObjectsCommand(
	node *Node,
	policy *ScanPolicy,
	namespace string,
	setName string,
	binNames []string,
	recordset *Recordset,
) *scanObjectsCommand {
	cmd := &scanObjectsCommand{
		baseMultiCommand: *newStreamingMultiCommand(node, recordset, namespace, false),
		policy:           policy,
		namespace:        namespace,
		setName:          setName,
		binNames:         binNames,
	}

	cmd.terminationErrorType = types.SCAN_TERMINATED

	return cmd
}

func (cmd *scanObjectsCommand) getPolicy(ifc command) Policy {
	return cmd.policy
}

func (cmd *scanObjectsCommand) writeBuffer(ifc command) Error {
	return cmd.setScan(cmd.policy, &cmd.namespace, &cmd.setName, cmd.binNames, cmd.taskID, nil)
}

func (cmd *scanObjectsCommand) parseResult(ifc command, conn *Connection) Error {
	return cmd.baseMultiCommand.parseResult(ifc, conn)
}

func (cmd *scanObjectsCommand) Execute() Error {
	defer cmd.recordset.signalEnd()
	err := cmd.execute(cmd, true)
	if err != nil {
		cmd.recordset.sendError(err)
	}
	return err
}
