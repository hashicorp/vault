// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package internal

import (
	"context"
	"errors"
)

// NewSemaphore creates a new semaphore.
func NewSemaphore(slots uint64) *Semaphore {
	ch := make(chan struct{}, slots)
	for i := uint64(0); i < slots; i++ {
		ch <- struct{}{}
	}

	return &Semaphore{
		permits: ch,
	}
}

// Semaphore is a synchronization primitive that controls access
// to a common resource.
type Semaphore struct {
	permits chan struct{}
}

// Len gets the number of permits available.
func (s *Semaphore) Len() uint64 {
	return uint64(len(s.permits))
}

// Wait waits until a resource is available or until the context
// is done.
func (s *Semaphore) Wait(ctx context.Context) error {
	select {
	case <-s.permits:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release releases a resource back into the pool.
func (s *Semaphore) Release() error {
	select {
	case s.permits <- struct{}{}:
	default:
		return errors.New("internal.Semaphore.Release: attempt to release more resources than are available")
	}

	return nil
}
