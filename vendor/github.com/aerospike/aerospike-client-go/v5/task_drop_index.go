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

import "strings"

// DropIndexTask is used to poll for long running create index completion.
type DropIndexTask struct {
	*baseTask

	namespace string
	indexName string
}

// NewDropIndexTask initializes a task with fields needed to query server nodes.
func NewDropIndexTask(cluster *Cluster, namespace string, indexName string) *DropIndexTask {
	return &DropIndexTask{
		baseTask:  newTask(cluster),
		namespace: namespace,
		indexName: indexName,
	}
}

// IsDone queries all nodes for task completion status.
func (tski *DropIndexTask) IsDone() (bool, Error) {
	command := "sindex/" + tski.namespace + "/" + tski.indexName
	nodes := tski.cluster.GetNodes()
	complete := false

	for _, node := range nodes {
		responseMap, err := node.requestInfoWithRetry(&tski.cluster.infoPolicy, 5, command)
		if err != nil {
			return false, err
		}

		for _, response := range responseMap {
			if strings.Contains(response, "FAIL:201") {
				complete = true
				continue
			}

			return false, nil
		}
	}
	return complete, nil
}

// OnComplete returns a channel that will be closed as soon as the task is finished.
// If an error is encountered during operation, an error will be sent on the channel.
func (tski *DropIndexTask) OnComplete() chan Error {
	return tski.onComplete(tski)
}
