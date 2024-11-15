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

func (clnt *Client) scanPartitionObjects(policy *ScanPolicy, tracker *partitionTracker, namespace string, setName string, rs *Recordset, binNames ...string) Error {
	defer rs.signalEnd()

	// for exponential backoff
	interval := policy.SleepBetweenRetries

	for {
		rs.resetTaskID()
		list, err := tracker.assignPartitionsToNodes(clnt.Cluster(), namespace)
		if err != nil {
			return err
		}

		wg := new(sync.WaitGroup)

		// the whole call should be wrapped in a goroutine
		wg.Add(len(list))

		// the whole call should be wrapped in a goroutine
		maxConcurrentNodes := policy.MaxConcurrentNodes
		if maxConcurrentNodes <= 0 {
			maxConcurrentNodes = len(list)
		}

		sem := semaphore.NewWeighted(int64(maxConcurrentNodes))
		ctx := context.Background()

		for _, nodePartition := range list {
			if err := sem.Acquire(ctx, 1); err != nil {
				logger.Logger.Error("Constraint Semaphore failed for Scan: %s", err.Error())
			}
			go func(nodePartition *nodePartitions) {
				defer sem.Release(1)
				defer wg.Done()
				if err := clnt.scanNodePartitionObjects(policy, rs, tracker, nodePartition, namespace, setName, binNames...); err != nil {
					logger.Logger.Debug("Error while Executing scan for node %s: %s", nodePartition.node.String(), err.Error())
				}
			}(nodePartition)
		}

		wg.Wait()

		done, err := tracker.isComplete(&policy.BasePolicy)
		if done || err != nil {
			// Scan is complete.
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

// ScanNode reads all records in specified namespace and set for one node only.
// If the policy is nil, the default relevant policy will be used.
func (clnt *Client) scanNodePartitionObjects(policy *ScanPolicy, recordset *Recordset, tracker *partitionTracker, nodePartition *nodePartitions, namespace string, setName string, binNames ...string) Error {
	command := newScanPartitionObjectsCommand(policy, tracker, nodePartition, namespace, setName, binNames, recordset)
	return command.Execute()
}
