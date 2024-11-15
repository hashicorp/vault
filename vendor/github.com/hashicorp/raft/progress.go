// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package raft

import (
	"context"
	"io"
	"sync"
	"time"

	hclog "github.com/hashicorp/go-hclog"
)

const (
	snapshotRestoreMonitorInterval = 10 * time.Second
)

type snapshotRestoreMonitor struct {
	logger          hclog.Logger
	cr              CountingReader
	size            int64
	networkTransfer bool

	once   sync.Once
	cancel func()
	doneCh chan struct{}
}

func startSnapshotRestoreMonitor(
	logger hclog.Logger,
	cr CountingReader,
	size int64,
	networkTransfer bool,
) *snapshotRestoreMonitor {
	ctx, cancel := context.WithCancel(context.Background())

	m := &snapshotRestoreMonitor{
		logger:          logger,
		cr:              cr,
		size:            size,
		networkTransfer: networkTransfer,
		cancel:          cancel,
		doneCh:          make(chan struct{}),
	}
	go m.run(ctx)
	return m
}

func (m *snapshotRestoreMonitor) run(ctx context.Context) {
	defer close(m.doneCh)

	ticker := time.NewTicker(snapshotRestoreMonitorInterval)
	defer ticker.Stop()

	ranOnce := false
	for {
		select {
		case <-ctx.Done():
			if !ranOnce {
				m.runOnce()
			}
			return
		case <-ticker.C:
			m.runOnce()
			ranOnce = true
		}
	}
}

func (m *snapshotRestoreMonitor) runOnce() {
	readBytes := m.cr.Count()
	pct := float64(100*readBytes) / float64(m.size)

	message := "snapshot restore progress"
	if m.networkTransfer {
		message = "snapshot network transfer progress"
	}

	m.logger.Info(message,
		"read-bytes", readBytes,
		"percent-complete", hclog.Fmt("%0.2f%%", pct),
	)
}

func (m *snapshotRestoreMonitor) StopAndWait() {
	m.once.Do(func() {
		m.cancel()
		<-m.doneCh
	})
}

type CountingReader interface {
	io.Reader
	Count() int64
}

type countingReader struct {
	reader io.Reader

	mu    sync.Mutex
	bytes int64
}

func (r *countingReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	r.mu.Lock()
	r.bytes += int64(n)
	r.mu.Unlock()
	return n, err
}

func (r *countingReader) Count() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.bytes
}

func newCountingReader(r io.Reader) *countingReader {
	return &countingReader{reader: r}
}

type countingReadCloser struct {
	*countingReader
	readCloser io.ReadCloser
}

func newCountingReadCloser(rc io.ReadCloser) *countingReadCloser {
	return &countingReadCloser{
		countingReader: newCountingReader(rc),
		readCloser:     rc,
	}
}

func (c countingReadCloser) Close() error {
	return c.readCloser.Close()
}

func (c countingReadCloser) WrappedReadCloser() io.ReadCloser {
	return c.readCloser
}

// ReadCloserWrapper allows access to an underlying ReadCloser from a wrapper.
type ReadCloserWrapper interface {
	io.ReadCloser
	WrappedReadCloser() io.ReadCloser
}

var _ ReadCloserWrapper = &countingReadCloser{}
