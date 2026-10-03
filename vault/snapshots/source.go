// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package snapshots

import (
	"context"
	"io"
)

// Source is used to read raw snapshot data
type Source interface {
	ReadCloser(ctx context.Context) (io.ReadCloser, error)
	Type(ctx context.Context) string
}

var _ Source = (*manualUploadSource)(nil)

type manualUploadSource struct {
	r io.ReadCloser
}

// NewManualSnapshotSource creates a new Source that returns the wrapped reader
// as the snapshot data
func NewManualSnapshotSource(r io.ReadCloser) Source {
	return &manualUploadSource{r: r}
}

func (m *manualUploadSource) Type(_ context.Context) string {
	return "manual"
}

func (m *manualUploadSource) ReadCloser(ctx context.Context) (io.ReadCloser, error) {
	return &ctxAwareReadCloser{ctx, m.r}, nil
}

type ctxAwareReadCloser struct {
	ctx context.Context
	io.ReadCloser
}

func (c *ctxAwareReadCloser) Read(p []byte) (n int, err error) {
	if c.ctx.Err() != nil {
		return 0, c.ctx.Err()
	}
	return c.ReadCloser.Read(p)
}
