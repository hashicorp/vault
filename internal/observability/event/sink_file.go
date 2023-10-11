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
}

// NewFileSink should be used to create a new FileSink.
// Accepted options: WithFileMode.
func NewFileSink(path string, format string, opt ...Option) (*FileSink, error) {
	const op = "event.NewFileSink"

	// Parse and check path
	p := strings.TrimSpace(path)
	if p == "" {
		return nil, fmt.Errorf("%s: path is required", op)
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("%s: error applying options: %w", op, err)
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
			return nil, fmt.Errorf("%s: unable to determine existing file mode: %w", op, err)
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
	}

	// Ensure that the file can be successfully opened for writing;
	// otherwise it will be too late to catch later without problems
	// (ref: https://github.com/hashicorp/vault/issues/550)
	if err := sink.open(); err != nil {
		return nil, fmt.Errorf("%s: sanity check failed; unable to open %q for writing: %w", op, path, err)
	}

	return sink, nil
}

// Process handles writing the event to the file sink.
func (s *FileSink) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "event.(FileSink).Process"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, ErrInvalidParameter)
	}

	// '/dev/null' path means we just do nothing and pretend we're done.
	if s.path == devnull {
		return nil, nil
	}

	formatted, found := e.Format(s.requiredFormat)
	if !found {
		return nil, fmt.Errorf("%s: unable to retrieve event formatted as %q", op, s.requiredFormat)
	}

	err := s.log(formatted)
	if err != nil {
		return nil, fmt.Errorf("%s: error writing file for sink: %w", op, err)
	}

	// return nil for the event to indicate the pipeline is complete.
	return nil, nil
}

// Reopen handles closing and reopening the file.
func (s *FileSink) Reopen() error {
	const op = "event.(FileSink).Reopen"

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
		return fmt.Errorf("%s: unable to close file for re-opening on sink: %w", op, err)
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
	const op = "event.(FileSink).open"

	if s.file != nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(s.path), s.fileMode); err != nil {
		return fmt.Errorf("%s: unable to create file %q: %w", op, s.path, err)
	}

	var err error
	s.file, err = os.OpenFile(s.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, s.fileMode)
	if err != nil {
		return fmt.Errorf("%s: unable to open file for sink: %w", op, err)
	}

	// Change the file mode in case the log file already existed.
	// We special case '/dev/null' since we can't chmod it, and bypass if the mode is zero.
	switch s.path {
	case devnull:
	default:
		if s.fileMode != 0 {
			err = os.Chmod(s.path, s.fileMode)
			if err != nil {
				return fmt.Errorf("%s: unable to change file %q permissions '%v' for sink: %w", op, s.path, s.fileMode, err)
			}
		}
	}

	return nil
}

// log writes the buffer to the file.
// It acquires a lock on the file to do this.
func (s *FileSink) log(data []byte) error {
	const op = "event.(FileSink).log"

	s.fileLock.Lock()
	defer s.fileLock.Unlock()

	reader := bytes.NewReader(data)

	if err := s.open(); err != nil {
		return fmt.Errorf("%s: unable to open file for sink: %w", op, err)
	}

	if _, err := reader.WriteTo(s.file); err == nil {
		return nil
	}

	// Otherwise, opportunistically try to re-open the FD, once per call (1 retry attempt).
	err := s.file.Close()
	if err != nil {
		return fmt.Errorf("%s: unable to close file for sink: %w", op, err)
	}

	s.file = nil

	if err := s.open(); err != nil {
		return fmt.Errorf("%s: unable to re-open file for sink: %w", op, err)
	}

	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("%s: unable to seek to start of file for sink: %w", op, err)
	}

	_, err = reader.WriteTo(s.file)
	if err != nil {
		return fmt.Errorf("%s: unable to re-write to file for sink: %w", op, err)
	}

	return nil
}
