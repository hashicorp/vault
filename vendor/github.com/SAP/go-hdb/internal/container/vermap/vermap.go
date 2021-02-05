// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Package vermap implements a key value map like used in session variables.
package vermap

import (
	"sync"
	"sync/atomic"
)

// A VerMap is a simple map[string]string keeping track changes via a version number.
// It is safe for concurrent use by multiple goroutines.
type VerMap struct {
	version int64 // atomic access

	mu sync.RWMutex
	m  map[string]string

	rlFlag bool // read lock set flag
}

// NewVerMap returns a new VerMap instance.
func NewVerMap() *VerMap { return &VerMap{m: make(map[string]string)} }

// WithRLock executes function f under a read lock.
func (vm *VerMap) WithRLock(f func()) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	vm.rlFlag = true
	defer func() { vm.rlFlag = false }()
	f()
}

// Version returns the VarMap version.
func (vm *VerMap) Version() int64 { return atomic.LoadInt64(&vm.version) }

// Store stores a string key value map in VarMap.
func (vm *VerMap) Store(m map[string]string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	atomic.AddInt64(&vm.version, 1)
	vm.m = make(map[string]string, len(m))
	for k, v := range m {
		vm.m[k] = v
	}
}

// load returns the content of a VarMap as string key value map.
func (vm *VerMap) load() map[string]string {
	m := make(map[string]string, len(vm.m))
	for k, v := range vm.m {
		m[k] = v
	}
	return m
}

// Load returns the content of a VarMap as string key value map.
func (vm *VerMap) Load() map[string]string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.load()
}

// LoadWithRLock - same as Load but needs to be executed in the context of method WithRLock.
func (vm *VerMap) LoadWithRLock() map[string]string {
	if !vm.rlFlag {
		panic("function LoadWithRLock called outside WithRLock context")
	}
	return vm.load()
}

// compare returns the changes (updates and deletes) in VarMap compared to map parameter m.
func (vm *VerMap) compare(m map[string]string) (upd map[string]string, del map[string]bool) {
	upd = make(map[string]string)
	del = make(map[string]bool)

	// updates
	for k, v := range vm.m {
		v2, ok := m[k]
		if !ok || ok && v != v2 {
			upd[k] = v
		}
	}
	// deletes
	for k := range m {
		if _, ok := vm.m[k]; !ok {
			del[k] = true
		}
	}
	return upd, del
}

// Compare returns the changes (updates and deletes) in VarMap compared to map parameter m.
func (vm *VerMap) Compare(m map[string]string) (upd map[string]string, del map[string]bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.compare(m)
}

// CompareWithRLock - same as Compare but needs to be executed in the context of method WithRLock.
func (vm *VerMap) CompareWithRLock(m map[string]string) (upd map[string]string, del map[string]bool) {
	if !vm.rlFlag {
		panic("function CompareWithRLock called outside WithRLock context")
	}
	return vm.compare(m)
}
