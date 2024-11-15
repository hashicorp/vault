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

// batchExecute Uses werrGroup to run commands using multiple goroutines,
// and waits for their return
func (clnt *Client) batchExecute(policy *BatchPolicy, batchNodes []*batchNode, cmd batcher) (int, Error) {
	maxConcurrentNodes := policy.ConcurrentNodes
	if maxConcurrentNodes <= 0 {
		maxConcurrentNodes = len(batchNodes)
	}

	// we need this list to count the number of filtered out records
	list := make([]batcher, 0, len(batchNodes))

	weg := newWeightedErrGroup(maxConcurrentNodes)
	for _, batchNode := range batchNodes {
		newCmd := cmd.cloneBatchCommand(batchNode)
		list = append(list, newCmd)
		weg.execute(newCmd)
	}

	errs := weg.wait()

	// count the filtered out records
	filteredOut := 0
	for i := range list {
		filteredOut += list[i].filteredOut()
	}

	return filteredOut, errs
}
