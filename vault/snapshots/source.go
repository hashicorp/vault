// Copyright (c) HashiCorp, Inc.
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

func (m *manualUploadSource) ReadCloser(_ context.Context) (io.ReadCloser, error) {
	return m.r, nil
}
