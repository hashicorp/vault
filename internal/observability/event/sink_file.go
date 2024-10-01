// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
)

// defaultFileMode is the default file permissions (read/write for everyone).
const (
	defaultFileMode = 0o600
	devnull         = "/dev/null"
)

var _ eventlogger.Node = (*FileSink)(nil)

// FileSink is a sink node which handles writing events to file.
type FileSink struct {
	file           *os.File
	fileLock       sync.RWMutex
	fileMode       os.FileMode
	path           string
	requiredFormat string
	logger         hclog.Logger
}

// NewFileSink should be used to create a new FileSink.
// Accepted options: WithFileMode.
func NewFileSink(path string, format string, opt ...Option) (*FileSink, error) {
	// Parse and check path
	p := strings.TrimSpace(path)
	if p == "" {
		return nil, fmt.Errorf("path is required: %w", ErrInvalidParameter)
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	mode := os.FileMode(defaultFileMode)
	// If we got an optional file mode supplied and our path isn't a special keyword
	// then we should use the supplied file mode, or maintain the existing file mode.
	switch {
	case path == devnull:
	case opts.withFileMode == nil:
	case *opts.withFileMode == 0: // Maintain the existing file's mode when set to "0000".
		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("unable to determine existing file mode: %w", err)
		}
		mode = fileInfo.Mode()
	default:
		mode = *opts.withFileMode
	}

	sink := &FileSink{
		file:           nil,
		fileLock:       sync.RWMutex{},
		fileMode:       mode,
		requiredFormat: format,
		path:           p,
		logger:         opts.withLogger,
	}

	// Ensure that the file can be successfully opened for writing;
	// otherwise it will be too late to catch later without problems
	// (ref: https://github.com/hashicorp/vault/issues/550)
	if err := sink.open(); err != nil {
		return nil, fmt.Errorf("sanity check failed; unable to open %q for writing: %w", sink.path, err)
	}

	return sink, nil
}

// Process handles writing the event to the file sink.
func (s *FileSink) Process(ctx context.Context, e *eventlogger.Event) (_ *eventlogger.Event, retErr error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	defer func() {
		// If the context is errored (cancelled), and we were planning to return
		// an error, let's also log (if we have a logger) in case the eventlogger's
		// status channel and errors propagated.
		if err := ctx.Err(); err != nil && retErr != nil && s.logger != nil {
			s.logger.Error("file sink error", "context", err, "error", retErr)
		}
	}()

	if e == nil {
		return nil, fmt.Errorf("event is nil: %w", ErrInvalidParameter)
	}

	// '/dev/null' path means we just do nothing and pretend we're done.
	if s.path == devnull {
		return nil, nil
	}

	formatted, found := e.Format(s.requiredFormat)
	if !found {
		return nil, fmt.Errorf("unable to retrieve event formatted as %q: %w", s.requiredFormat, ErrInvalidParameter)
	}

	err := s.log(ctx, formatted)
	if err != nil {
		return nil, fmt.Errorf("error writing file for sink %q: %w", s.path, err)
	}

	// return nil for the event to indicate the pipeline is complete.
	return nil, nil
}

// Reopen handles closing and reopening the file.
func (s *FileSink) Reopen() error {
	// '/dev/null' path means we just do nothing and pretend we're done.
	if s.path == devnull {
		return nil
	}

	s.fileLock.Lock()
	defer s.fileLock.Unlock()

	if s.file == nil {
		return s.open()
	}

	err := s.file.Close()
	// Set to nil here so that even if we error out, on the next access open() will be tried.
	s.file = nil
	if err != nil {
		return fmt.Errorf("unable to close file for re-opening on sink %q: %w", s.path, err)
	}

	return s.open()
}

// Type describes the type of this node (sink).
func (s *FileSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}

// open attempts to open a file at the sink's path, with the sink's fileMode permissions
// if one is not already open.
// It doesn't have any locking and relies on calling functions of FileSink to
// handle this (e.g. log and Reopen methods).
func (s *FileSink) open() error {
	if s.file != nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(s.path), s.fileMode); err != nil {
		return fmt.Errorf("unable to create file %q: %w", s.path, err)
	}

	var err error
	s.file, err = os.OpenFile(s.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, s.fileMode)
	if err != nil {
		return fmt.Errorf("unable to open file for sink %q: %w", s.path, err)
	}

	// Change the file mode in case the log file already existed.
	// We special case '/dev/null' since we can't chmod it, and bypass if the mode is zero.
	switch s.path {
	case devnull:
	default:
		if s.fileMode != 0 {
			err = os.Chmod(s.path, s.fileMode)
			if err != nil {
				return fmt.Errorf("unable to change file permissions '%v' for sink %q: %w", s.fileMode, s.path, err)
			}
		}
	}

	return nil
}

// log writes the buffer to the file.
// NOTE: We attempt to acquire a lock on the file in order to write, but will
// yield if the context is 'done'.
func (s *FileSink) log(ctx context.Context, data []byte) error {
	// Wait for the lock, but ensure we check for a cancelled context as soon as
	// we have it, as there's no point in continuing if we're cancelled.
	s.fileLock.Lock()
	defer s.fileLock.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	reader := bytes.NewReader(data)

	if err := s.open(); err != nil {
		return fmt.Errorf("unable to open file for sink %q: %w", s.path, err)
	}

	if _, err := reader.WriteTo(s.file); err == nil {
		return nil
	}

	// Otherwise, opportunistically try to re-open the FD, once per call (1 retry attempt).
	err := s.file.Close()
	if err != nil {
		return fmt.Errorf("unable to close file for sink %q: %w", s.path, err)
	}

	s.file = nil

	if err := s.open(); err != nil {
		return fmt.Errorf("unable to re-open file for sink %q: %w", s.path, err)
	}

	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("unable to seek to start of file for sink %q: %w", s.path, err)
	}

	_, err = reader.WriteTo(s.file)
	if err != nil {
		return fmt.Errorf("unable to re-write to file for sink %q: %w", s.path, err)
	}

	return nil
}
