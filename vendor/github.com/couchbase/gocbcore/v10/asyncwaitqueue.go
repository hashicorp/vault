// Copyright 2021 Couchbase
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gocbcore

import (
	"sync"
)

type asyncWaitGroup struct {
	lock    sync.Mutex
	count   int
	waiters []func()
}

func (q *asyncWaitGroup) IsEmpty() bool {
	q.lock.Lock()
	isEmpty := q.count == 0
	q.lock.Unlock()

	return isEmpty
}

func (q *asyncWaitGroup) Add(n int) {
	var waiters []func()

	q.lock.Lock()
	q.count += n
	if q.count == 0 {
		waiters = q.waiters
		q.waiters = nil
	}
	q.lock.Unlock()

	for _, waiter := range waiters {
		waiter()
	}
}

func (q *asyncWaitGroup) Done() {
	q.Add(-1)
}

func (q *asyncWaitGroup) Wait(fn func()) {
	q.lock.Lock()
	if q.count == 0 {
		q.lock.Unlock()

		fn()
		return
	}

	q.waiters = append(q.waiters, fn)
	q.lock.Unlock()
}
