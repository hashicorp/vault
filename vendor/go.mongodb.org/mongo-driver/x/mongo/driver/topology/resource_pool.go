// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// expiredFunc is the function type used for testing whether or not resources in a resourcePool have stale. It should
// return true if the resource has stale and can be removed from the pool.
type expiredFunc func(interface{}) bool

// closeFunc is the function type used to closeConnection resources in a resourcePool. The pool will always call this function
// asynchronously
type closeFunc func(interface{})

// initFunc is the function used to add a resource to the resource pool to maintain minimum size. It returns a new
// resource each time it is called.
type initFunc func() interface{}

type resourcePoolConfig struct {
	MinSize          uint64
	MaintainInterval time.Duration
	ExpiredFn        expiredFunc
	CloseFn          closeFunc
	InitFn           initFunc
}

// setup sets defaults in the rpc and checks that the given values are valid
func (rpc *resourcePoolConfig) setup() error {
	if rpc.ExpiredFn == nil {
		return fmt.Errorf("an ExpiredFn is required to create a resource pool")
	}
	if rpc.CloseFn == nil {
		return fmt.Errorf("an CloseFn is required to create a resource pool")
	}
	if rpc.MaintainInterval == time.Duration(0) {
		return fmt.Errorf("unable to have MaintainInterval time of %v", rpc.MaintainInterval)
	}
	return nil
}

// resourcePoolElement is a link list element
type resourcePoolElement struct {
	next, prev *resourcePoolElement
	value      interface{}
}

// resourcePool is a concurrent resource pool
type resourcePool struct {
	start, end       *resourcePoolElement
	size, minSize    uint64
	expiredFn        expiredFunc
	closeFn          closeFunc
	initFn           initFunc
	maintainTimer    *time.Timer
	maintainInterval time.Duration

	sync.Mutex
}

// NewResourcePool creates a new resourcePool instance that is capped to maxSize resources.
// If maxSize is 0, the pool size will be unbounded.
func newResourcePool(config resourcePoolConfig) (*resourcePool, error) {
	err := (&config).setup()
	if err != nil {
		return nil, err
	}
	rp := &resourcePool{
		minSize:          config.MinSize,
		expiredFn:        config.ExpiredFn,
		closeFn:          config.CloseFn,
		initFn:           config.InitFn,
		maintainInterval: config.MaintainInterval,
	}

	return rp, nil
}

func (rp *resourcePool) initialize() {
	rp.Lock()
	rp.maintainTimer = time.AfterFunc(rp.maintainInterval, rp.Maintain)
	rp.Unlock()

	rp.Maintain()
}

// add will add a new rpe to the pool, requires that the resource pool is locked
func (rp *resourcePool) add(e *resourcePoolElement) {
	if e == nil {
		e = &resourcePoolElement{
			value: rp.initFn(),
		}
	}

	e.next = rp.start
	if rp.start != nil {
		rp.start.prev = e
	}
	rp.start = e
	if rp.end == nil {
		rp.end = e
	}
	atomic.AddUint64(&rp.size, 1)
}

// Get returns the first un-stale resource from the pool. If no such resource can be found, nil is returned.
func (rp *resourcePool) Get() interface{} {
	rp.Lock()
	defer rp.Unlock()

	for rp.start != nil {
		curr := rp.start
		rp.remove(curr)
		if !rp.expiredFn(curr.value) {
			return curr.value
		}
		rp.closeFn(curr.value)
	}
	return nil
}

// Put puts the resource back into the pool if it will not exceed the max size of the pool
func (rp *resourcePool) Put(v interface{}) bool {
	if rp.expiredFn(v) {
		rp.closeFn(v)
		return false
	}

	rp.Lock()
	defer rp.Unlock()
	rp.add(&resourcePoolElement{value: v})
	return true
}

// remove removes a rpe from the linked list. Requires that the pool be locked
func (rp *resourcePool) remove(e *resourcePoolElement) {
	if e == nil {
		return
	}

	if e.prev != nil {
		e.prev.next = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	}
	if e == rp.start {
		rp.start = e.next
	}
	if e == rp.end {
		rp.end = e.prev
	}
	atomicSubtract1Uint64(&rp.size)
}

// Maintain puts the pool back into a state of having a correct number of resources if possible and removes all stale resources
func (rp *resourcePool) Maintain() {
	rp.Lock()
	defer rp.Unlock()
	for curr := rp.end; curr != nil; curr = curr.prev {
		if rp.expiredFn(curr.value) {
			rp.remove(curr)
			rp.closeFn(curr.value)
		}
	}

	for atomic.LoadUint64(&rp.size) < rp.minSize {
		rp.add(nil)
	}

	// reset the timer for the background cleanup routine
	if rp.maintainTimer == nil {
		rp.maintainTimer = time.AfterFunc(rp.maintainInterval, rp.Maintain)
	}
	if !rp.maintainTimer.Stop() {
		rp.maintainTimer = time.AfterFunc(rp.maintainInterval, rp.Maintain)
		return
	}
	rp.maintainTimer.Reset(rp.maintainInterval)
}

// Close clears the pool and stops the background maintenance of the pool
func (rp *resourcePool) Close() {
	rp.Clear()
	_ = rp.maintainTimer.Stop()
}

// Clear closes all resources in the pool
func (rp *resourcePool) Clear() {
	rp.Lock()
	defer rp.Unlock()
	for ; rp.start != nil; rp.start = rp.start.next {
		rp.closeFn(rp.start.value)
	}
	atomic.StoreUint64(&rp.size, 0)
	rp.end = nil
}

func atomicSubtract1Uint64(p *uint64) {
	if p == nil || atomic.LoadUint64(p) == 0 {
		return
	}

	for {
		expected := atomic.LoadUint64(p)
		if atomic.CompareAndSwapUint64(p, expected, expected-1) {
			return
		}
	}
}
