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

type asyncMutex struct {
	lock     sync.Mutex
	waiters  []func(func())
	codeCntr uint64
	curCode  uint64
}

// acquireLocked grabs the lock and returns its unlock code to be passed to .unlock()
func (q *asyncMutex) acquireLocked() uint64 {
	q.codeCntr++
	myCode := q.codeCntr

	if q.curCode != 0 {
		logWarnf("unexpectedly trying to lock asyncMutex while already locked")
	}
	q.curCode = myCode

	return myCode
}

func (q *asyncMutex) Lock(cb func(func())) {
	q.lock.Lock()

	if q.curCode == 0 {
		myCode := q.acquireLocked()
		q.lock.Unlock()

		cb(func() {
			q.unlock(myCode)
		})

		return
	}

	q.waiters = append(q.waiters, cb)

	q.lock.Unlock()
}

func (q *asyncMutex) LockSync() {
	waitCh := make(chan struct{}, 1)

	q.Lock(func(unlockFn func()) {
		waitCh <- struct{}{}
	})

	<-waitCh
}

func (q *asyncMutex) UnlockSync() {
	// We cheat for sync locks and just grab the code
	q.lock.Lock()
	syncCode := q.curCode
	q.lock.Unlock()

	q.unlock(syncCode)
}

func (q *asyncMutex) unlock(myCode uint64) {
	q.lock.Lock()

	if myCode != q.curCode {
		logWarnf("unexpected unlock code for asyncMutex unlock")
		q.lock.Unlock()
		return
	}

	q.curCode = 0

	if len(q.waiters) == 0 {
		q.lock.Unlock()
		return
	}

	nextFn := q.waiters[0]
	q.waiters = q.waiters[1:]

	nextCode := q.acquireLocked()
	q.lock.Unlock()

	nextFn(func() {
		q.unlock(nextCode)
	})
}
