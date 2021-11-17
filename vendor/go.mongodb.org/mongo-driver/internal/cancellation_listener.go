// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package internal

import "context"

// CancellationListener listens for context cancellation in a loop until the context expires or the listener is aborted.
type CancellationListener struct {
	aborted bool
	done    chan struct{}
}

// NewCancellationListener constructs a CancellationListener.
func NewCancellationListener() *CancellationListener {
	return &CancellationListener{
		done: make(chan struct{}),
	}
}

// Listen blocks until the provided context is cancelled or listening is aborted via the StopListening function. If this
// detects that the context has been cancelled (i.e. ctx.Err() == context.Canceled), the provided callback is called to
// abort in-progress work. Even if the context expires, this function will block until StopListening is called.
func (c *CancellationListener) Listen(ctx context.Context, abortFn func()) {
	c.aborted = false

	select {
	case <-ctx.Done():
		if ctx.Err() == context.Canceled {
			c.aborted = true
			abortFn()
		}

		<-c.done
	case <-c.done:
	}
}

// StopListening stops the in-progress Listen call. This blocks if there is no in-progress Listen call. This function
// will return true if the provided abort callback was called when listening for cancellation on the previous context.
func (c *CancellationListener) StopListening() bool {
	c.done <- struct{}{}
	return c.aborted
}
