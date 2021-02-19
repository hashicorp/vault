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

type queryCommand struct {
	baseMultiCommand

	policy      *QueryPolicy
	writePolicy *WritePolicy
	statement   *Statement
	operations  []*Operation
}

func newQueryCommand(node *Node, policy *QueryPolicy, writePolicy *WritePolicy, statement *Statement, operations []*Operation, recordset *Recordset, clusterKey int64, first bool) *queryCommand {
	return &queryCommand{
		baseMultiCommand: *newCorrectMultiCommand(node, recordset, statement.Namespace, clusterKey, first),
		policy:           policy,
		writePolicy:      writePolicy,
		statement:        statement,
		operations:       operations,
	}
}

func (cmd *queryCommand) getPolicy(ifc command) Policy {
	return cmd.policy
}

func (cmd *queryCommand) writeBuffer(ifc command) (err error) {
	return cmd.setQuery(cmd.policy, cmd.writePolicy, cmd.statement, cmd.operations, cmd.writePolicy != nil)
}

func (cmd *queryCommand) parseResult(ifc command, conn *Connection) error {
	return cmd.baseMultiCommand.parseResult(ifc, conn)
}

// Execute will run the query.
func (cmd *queryCommand) Execute() error {
	defer cmd.recordset.signalEnd()
	err := cmd.execute(cmd, true)
	if err != nil {
		cmd.recordset.sendError(err)
	}
	return err
}
