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
	"time"

	"github.com/aerospike/aerospike-client-go/v5/internal/atomic"
	"github.com/aerospike/aerospike-client-go/v5/types"
)

// Task interface defines methods for asynchronous tasks.
type Task interface {
	IsDone() (bool, Error)

	onComplete(ifc Task) chan Error
	OnComplete() chan Error
}

// baseTask is used to poll for server task completion.
type baseTask struct {
	retries atomic.Int
	cluster *Cluster
}

// newTask initializes task with fields needed to query server nodes.
func newTask(cluster *Cluster) *baseTask {
	return &baseTask{
		cluster: cluster,
	}
}

// Wait for asynchronous task to complete using default sleep interval.
func (btsk *baseTask) onComplete(ifc Task) chan Error {
	ch := make(chan Error, 1)

	// goroutine will loop every <interval> until IsDone() returns true or error
	go func() {
		// always close the channel on return
		defer close(ch)

		var interval = 100 * time.Millisecond

		for {
			time.Sleep(interval)

			done, err := ifc.IsDone()
			// Every 5 failed retries increase the interval
			if btsk.retries.IncrementAndGet()%5 == 0 {
				interval *= 2

				if interval > 5*time.Second {
					interval = 5 * time.Second
				}
			}
			if err != nil {
				if err.Matches(types.TIMEOUT) {
					err.markInDoubt()
				}
				ch <- err
				return
			} else if done {
				ch <- nil
				return
			}
		} // for
	}()

	return ch
}
