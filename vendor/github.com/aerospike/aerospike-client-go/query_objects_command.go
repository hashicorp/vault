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

import . "github.com/aerospike/aerospike-client-go/types"

type queryObjectsCommand struct {
	queryCommand
}

func newQueryObjectsCommand(node *Node, policy *QueryPolicy, statement *Statement, recordset *Recordset, clusterKey int64, first bool) *queryObjectsCommand {
	cmd := &queryObjectsCommand{
		queryCommand: *newQueryCommand(node, policy, nil, statement, nil, recordset, clusterKey, first),
	}

	cmd.terminationErrorType = QUERY_TERMINATED

	return cmd
}

func (cmd *queryObjectsCommand) Execute() error {
	defer cmd.recordset.signalEnd()
	err := cmd.execute(cmd, true)
	if err != nil {
		cmd.recordset.sendError(err)
	}
	return err
}
