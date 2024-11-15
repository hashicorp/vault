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
	"regexp"
	"strconv"
	"strings"
)

// IndexTask is used to poll for long running create index completion.
type IndexTask struct {
	*baseTask

	namespace string
	indexName string
}

// NewIndexTask initializes a task with fields needed to query server nodes.
func NewIndexTask(cluster *Cluster, namespace string, indexName string) *IndexTask {
	return &IndexTask{
		baseTask:  newTask(cluster),
		namespace: namespace,
		indexName: indexName,
	}
}

// IsDone queries all nodes for task completion status.
func (tski *IndexTask) IsDone() (bool, Error) {
	command := "sindex/" + tski.namespace + "/" + tski.indexName
	nodes := tski.cluster.GetNodes()
	complete := false

	r := regexp.MustCompile(`\.*load_pct=(\d+)\.*`)

	for _, node := range nodes {
		responseMap, err := node.requestInfoWithRetry(&tski.cluster.infoPolicy, 5, command)
		if err != nil {
			return false, err
		}

		for _, response := range responseMap {
			find := "load_pct="
			index := strings.Index(response, find)

			if index < 0 {
				if tski.retries.Get() > 20 {
					complete = true
				}
				continue
			}

			matchRes := r.FindStringSubmatch(response)
			// we know it exists and is a valid number
			pct, _ := strconv.Atoi(matchRes[1])

			if pct >= 0 && pct < 100 {
				return false, nil
			}
			complete = true
		}
	}
	return complete, nil
}

// OnComplete returns a channel that will be closed as soon as the task is finished.
// If an error is encountered during operation, an error will be sent on the channel.
func (tski *IndexTask) OnComplete() chan Error {
	return tski.onComplete(tski)
}
