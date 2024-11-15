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
	"runtime"
	"sync"
)

// singleConnectionHeap is a non-blocking LIFO heap.
// If the heap is empty, nil is returned.
// if the heap is full, offer will return false
type singleConnectionHeap struct {
	head, tail uint32
	data       []*Connection
	size       uint32
	full       bool
	mutex      sync.Mutex
}

// newSingleConnectionHeap creates a new heap with initial size.
func newSingleConnectionHeap(size int) *singleConnectionHeap {
	if size <= 0 {
		panic("Heap size cannot be less than 1")
	}

	return &singleConnectionHeap{
		full: false,
		data: make([]*Connection, uint32(size)),
		size: uint32(size),
	}
}

func (h *singleConnectionHeap) cleanup() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := range h.data {
		if h.data[i] != nil {
			h.data[i].Close()
		}

		h.data[i] = nil
	}

	// make sure offer and poll both fail
	h.data = nil
	h.full = true
	h.head = 0
	h.tail = 0
}

// Offer adds an item to the heap unless the heap is full.
// In case the heap is full, the item will not be added to the heap
// and false will be returned
func (h *singleConnectionHeap) Offer(conn *Connection) bool {
	h.mutex.Lock()

	// make sure heap is not full or cleaned up
	if h.full || len(h.data) == 0 {
		h.mutex.Unlock()
		return false
	}

	h.head = (h.head + 1) % h.size
	h.full = (h.head == h.tail)
	h.data[h.head] = conn
	h.mutex.Unlock()
	return true
}

// Poll removes and returns an item from the heap.
// If the heap is empty, nil will be returned.
func (h *singleConnectionHeap) Poll() (res *Connection) {
	h.mutex.Lock()

	// the heap has been cleaned up
	if len(h.data) == 0 {
		h.mutex.Unlock()
		return nil
	}

	// if heap is not empty
	if (h.tail != h.head) || h.full {
		res = h.data[h.head]
		h.data[h.head] = nil

		h.full = false
		if h.head == 0 {
			h.head = h.size - 1
		} else {
			h.head--
		}
	}

	h.mutex.Unlock()
	return res
}

// DropIdleTail closes idle connection in tail.
// It will return true if tail connection was idle and dropped
func (h *singleConnectionHeap) DropIdleTail() bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// the heap has been cleaned up
	if h.data == nil {
		return false
	}

	// if heap is not empty
	if h.full || (h.tail != h.head) {
		conn := h.data[(h.tail+1)%h.size]

		if conn.IsConnected() && !conn.isIdle() {
			return false
		}

		h.tail = (h.tail + 1) % h.size
		h.data[h.tail] = nil
		h.full = false
		conn.Close()

		return true
	}

	return false
}

// Len returns the number of connections in the heap
func (h *singleConnectionHeap) Len() int {
	cnt := 0
	h.mutex.Lock()

	if !h.full {
		if h.head >= h.tail {
			cnt = int(h.head) - int(h.tail)
		} else {
			cnt = int(h.size) - (int(h.tail) - int(h.head))
		}
	} else {
		cnt = int(h.size)
	}
	h.mutex.Unlock()
	return cnt
}

// connectionHeap is a non-blocking FIFO heap.
// If the heap is empty, nil is returned.
// if the heap is full, offer will return false
type connectionHeap struct {
	maxSize int
	minSize int
	heaps   []singleConnectionHeap
}

// Close cleans up all the data and removes all the references from
// active objects to ensure GC cleans up everything.
func (h *connectionHeap) cleanup() {
	for i := range h.heaps {
		h.heaps[i].cleanup()
	}
}

func newConnectionHeap(minSize, maxSize int) *connectionHeap {
	if minSize > maxSize {
		panic("minSize is bigger than maxSize for connection heap")
	}

	heapCount := runtime.NumCPU()
	if heapCount > maxSize {
		heapCount = maxSize
	}

	// will be >= 1
	perHeapSize := maxSize / heapCount

	heaps := make([]singleConnectionHeap, heapCount)
	for i := range heaps {
		heaps[i] = *newSingleConnectionHeap(perHeapSize)
	}

	// add a heap for the remainder
	remainder := maxSize - heapCount*perHeapSize
	if remainder > 0 {
		heaps = append(heaps, *newSingleConnectionHeap(remainder))
	}

	return &connectionHeap{
		maxSize: maxSize,
		minSize: minSize,
		heaps:   heaps,
	}
}

// Offer adds an item to the heap unless the heap is full.
// In case the heap is full, the item will not be added to the heap
// and false will be returned
func (h *connectionHeap) Offer(conn *Connection, hint byte) bool {
	idx := int(hint) % len(h.heaps)
	end := idx + len(h.heaps)
	for i := idx; i < end; i++ {
		if h.heaps[i%len(h.heaps)].Offer(conn) {
			// success
			return true
		}
	}
	return false
}

// Poll removes and returns an item from the heap.
// If the heap is empty, nil will be returned.
func (h *connectionHeap) Poll(hint byte) (res *Connection) {
	idx := int(hint)

	end := idx + len(h.heaps)
	for i := idx; i < end; i++ {
		if conn := h.heaps[i%len(h.heaps)].Poll(); conn != nil {
			return conn
		}
	}
	return nil
}

// DropIdle closes all idle connections.
// It will only drop connections if there are
// at least ClientPolicy.MinConnectionPerNode available
func (h *connectionHeap) DropIdle() {
	// decide how many conns are allowed to drop
	// in minSize is 0, up to all connection can
	// be closed if idle
	excessCount := h.LenAll() - h.minSize
	if excessCount <= 0 {
		return
	}

	for i := 0; i < len(h.heaps); i++ {
		for h.heaps[i].DropIdleTail() {
			excessCount--
			if excessCount == 0 {
				return
			}
		}
	}
}

// Cap returns the total capacity of the connectionHeap
func (h *connectionHeap) Cap() int {
	return h.maxSize
}

// Len returns the number of connections in a specific sub-heap.
func (h *connectionHeap) Len(hint byte) (cnt int) {
	return h.heaps[hint].Len()
}

// LenAll returns the number of connections in all sub-heaps.
func (h *connectionHeap) LenAll() int {
	cnt := 0
	for i := range h.heaps {
		cnt += h.heaps[i].Len()
	}

	return cnt
}
