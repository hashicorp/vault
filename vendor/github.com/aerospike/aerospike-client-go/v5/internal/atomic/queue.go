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

package atomic

import "sync"

// Queue is a non-blocking FIFO queue.
// If the queue is empty, nil is returned.
// if the queue is full, offer will return false
type Queue struct {
	head, tail uint32
	data       []interface{}
	size       uint32
	wrapped    bool
	mutex      sync.Mutex
}

// NewQueue creates a new queue with initial size.
func NewQueue(size int) *Queue {
	if size <= 0 {
		panic("Queue size cannot be less than 1")
	}

	return &Queue{
		wrapped: false,
		data:    make([]interface{}, uint32(size)),
		size:    uint32(size),
	}
}

// Offer adds an item to the queue unless the queue is full.
// In case the queue is full, the item will not be added to the queue
// and false will be returned
func (q *Queue) Offer(obj interface{}) bool {
	q.mutex.Lock()

	// make sure queue is not full
	if q.tail == q.head && q.wrapped {
		q.mutex.Unlock()
		return false
	}

	if q.head+1 == q.size {
		q.wrapped = true
	}

	q.head = (q.head + 1) % q.size
	q.data[q.head] = obj
	q.mutex.Unlock()
	return true
}

// Poll removes and returns an item from the queue.
// If the queue is empty, nil will be returned.
func (q *Queue) Poll() (res interface{}) {
	q.mutex.Lock()

	// if queue is not empty
	if q.wrapped || (q.tail != q.head) {
		if q.tail+1 == q.size {
			q.wrapped = false
		}
		q.tail = (q.tail + 1) % q.size
		res = q.data[q.tail]
	}

	q.mutex.Unlock()
	return res
}
