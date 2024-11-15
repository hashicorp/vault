// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listener

import (
	"net"
	"sync"
)

// WithCloseCallback creates a new Listener with a given wrapped listener and a
// callback that is called when a wrapped Close() method succeeds.
func WithCloseCallback(listener net.Listener, callback func()) net.Listener {
	return &wrapper{
		Listener: listener,
		callback: callback,
		once:     sync.Once{},
	}
}

// wrapper embeds the net.Listener and overrides the Close() method in a way that
// it is synchronously called when a Close() method of a wrapped listener
// succeeds. It is guaranteed that the callback will be called only once.
type wrapper struct {
	net.Listener

	// callback is called when Close is called on the listener.
	callback func()

	// once guarantees that a callback has been called only once.
	once sync.Once
}

// Close implements the net.Listener interface. The call is propagated to the
// wrapped listener, and in case of a successful result a callback is called
// synchronously and only once.
func (w *wrapper) Close() error {
	err := w.Listener.Close()

	if err == nil && w.callback != nil {
		w.once.Do(w.callback)
	}

	return err
}
