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
	"context"
	"sync"
	"time"

	"github.com/aerospike/aerospike-client-go/v5/logger"
	"golang.org/x/sync/semaphore"
)

func (clnt *Client) queryPartitionObjects(policy *QueryPolicy, tracker *partitionTracker, statement *Statement, rs *Recordset) Error {
	defer rs.signalEnd()

	// for exponential backoff
	interval := policy.SleepBetweenRetries

	for {
		rs.resetTaskID()
		list, err := tracker.assignPartitionsToNodes(clnt.Cluster(), statement.Namespace)
		if err != nil {
			return err
		}

		wg := new(sync.WaitGroup)

		// the whole call should be wrapped in a goroutine
		wg.Add(len(list))

		// results channel must be async for performance
		maxConcurrentNodes := policy.MaxConcurrentNodes
		if maxConcurrentNodes <= 0 {
			maxConcurrentNodes = len(list)
		}

		sem := semaphore.NewWeighted(int64(maxConcurrentNodes))
		ctx := context.Background()

		for _, nodePartition := range list {
			if err := sem.Acquire(ctx, 1); err != nil {
				logger.Logger.Error("Constraint Semaphore failed for Query: %s", err.Error())
			}
			go func(nodePartition *nodePartitions) {
				defer sem.Release(1)
				defer wg.Done()
				if err := clnt.queryNodePartitionObjects(policy, rs, tracker, nodePartition, statement); err != nil {
					logger.Logger.Debug("Error while Executing query for node %s: %s", nodePartition.node.String(), err.Error())
				}
			}(nodePartition)
		}

		wg.Wait()

		done, err := tracker.isComplete(&policy.BasePolicy)
		if done || err != nil {
			// Query is complete.
			return err
		}

		if policy.SleepBetweenRetries > 0 {
			// Sleep before trying again.
			time.Sleep(interval)

			if policy.SleepMultiplier > 1 {
				interval = time.Duration(float64(interval) * policy.SleepMultiplier)
			}
		}
	}

}

// QueryNode reads all records in specified namespace and set for one node only.
// If the policy is nil, the default relevant policy will be used.
func (clnt *Client) queryNodePartitionObjects(policy *QueryPolicy, recordset *Recordset, tracker *partitionTracker, nodePartition *nodePartitions, statement *Statement) Error {
	command := newQueryPartitionObjectsCommand(policy, tracker, nodePartition, statement, recordset)
	return command.Execute()
}
