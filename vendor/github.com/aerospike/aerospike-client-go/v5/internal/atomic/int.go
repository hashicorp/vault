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

import "sync/atomic"

// Int implements an int value with atomic semantics
type Int struct {
	val int64
}

// NewInt generates a newVal Int instance.
func NewInt(value int) *Int {
	v := int64(value)
	return &Int{
		val: v,
	}
}

// AddAndGet atomically adds the given value to the current value.
func (ai *Int) AddAndGet(delta int) int {
	res := int(atomic.AddInt64(&ai.val, int64(delta)))
	return res
}

// CompareAndSet atomically sets the value to the given updated value if the current value == expected value.
// Returns true if the expectation was met
func (ai *Int) CompareAndSet(expect int, update int) bool {
	res := atomic.CompareAndSwapInt64(&ai.val, int64(expect), int64(update))
	return res
}

// DecrementAndGet atomically decrements current value by one and returns the result.
func (ai *Int) DecrementAndGet() int {
	res := int(atomic.AddInt64(&ai.val, -1))
	return res
}

// Get atomically retrieves the current value.
func (ai *Int) Get() int {
	res := int(atomic.LoadInt64(&ai.val))
	return res
}

// GetAndAdd atomically adds the given delta to the current value and returns the result.
func (ai *Int) GetAndAdd(delta int) int {
	newVal := atomic.AddInt64(&ai.val, int64(delta))
	res := int(newVal - int64(delta))
	return res
}

// GetAndDecrement atomically decrements the current value by one and returns the result.
func (ai *Int) GetAndDecrement() int {
	newVal := atomic.AddInt64(&ai.val, -1)
	res := int(newVal + 1)
	return res
}

// GetAndIncrement atomically increments current value by one and returns the result.
func (ai *Int) GetAndIncrement() int {
	newVal := atomic.AddInt64(&ai.val, 1)
	res := int(newVal - 1)
	return res
}

// GetAndSet atomically sets current value to the given value and returns the old value.
func (ai *Int) GetAndSet(newValue int) int {
	res := int(atomic.SwapInt64(&ai.val, int64(newValue)))
	return res
}

// IncrementAndGet atomically increments current value by one and returns the result.
func (ai *Int) IncrementAndGet() int {
	res := int(atomic.AddInt64(&ai.val, 1))
	return res
}

// Set atomically sets current value to the given value.
func (ai *Int) Set(newValue int) {
	atomic.StoreInt64(&ai.val, int64(newValue))
}
