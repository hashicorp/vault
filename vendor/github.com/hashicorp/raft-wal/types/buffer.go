// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package types

// PooledBuffer is a wrapper that allows WAL to return read buffers to segment
// implementations when we're done decoding.
type PooledBuffer struct {
	Bs      []byte
	CloseFn func()
}

// Close implements io.Closer and returns the buffer to the pool. It should be
// called exactly once for each buffer when it's no longer needed. It's no
// longer safe to access Bs or any slice taken from it after the call.
func (b *PooledBuffer) Close() error {
	if b.CloseFn != nil {
		b.CloseFn()
	}
	return nil
}
