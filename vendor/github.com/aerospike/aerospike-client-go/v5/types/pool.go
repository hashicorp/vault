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

package types

import (
	"github.com/aerospike/aerospike-client-go/v5/internal/atomic"
)

// Pool implements a general purpose fixed-size pool.
type Pool struct {
	pool *atomic.Queue

	// New will create a new object
	New func(params ...interface{}) interface{}
	// IsUsable checks if the object polled from the pool is still fresh and usable
	IsUsable func(obj interface{}, params ...interface{}) bool
	// CanReturn checkes if the object is eligible to go back to the pool
	CanReturn func(obj interface{}) bool
	// Finalize will be called when an object is not eligible to go back to the pool.
	// Usable to close connections, file handles, ...
	Finalize func(obj interface{})
}

// NewPool creates a new fixed size pool.
func NewPool(poolSize int) *Pool {
	return &Pool{
		pool: atomic.NewQueue(poolSize),
	}
}

// Get returns an element from the pool.
// If the pool is empty, or the returned element is not usable,
// nil or the result of the New function will be returned
func (bp *Pool) Get(params ...interface{}) interface{} {
	res := bp.pool.Poll()
	if res == nil || (bp.IsUsable != nil && !bp.IsUsable(res, params...)) {
		// not usable, so finalize
		if res != nil && bp.Finalize != nil {
			bp.Finalize(res)
		}

		if bp.New != nil {
			res = bp.New(params...)
		}
	}

	return res
}

// Put will add the elem back to the pool, unless the pool is full.
func (bp *Pool) Put(obj interface{}) {
	finalize := true
	if bp.CanReturn == nil || bp.CanReturn(obj) {
		finalize = !bp.pool.Offer(obj)
	}

	if finalize && bp.Finalize != nil {
		bp.Finalize(obj)
	}
}
