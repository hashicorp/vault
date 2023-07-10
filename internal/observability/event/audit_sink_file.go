// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/eventlogger"
)

// defaultFileMode is the default file permissions (read/write for everyone).
const defaultFileMode = 0o600

// AuditFileSink is a sink node which handles writing audit events to file.
type AuditFileSink struct {
	file     *os.File
	fileLock sync.RWMutex
	fileMode os.FileMode
	path     string
	format   auditFormat
	prefix   string
}

// AuditFileSinkConfig is the configuration for the AuditFileSink.
type AuditFileSinkConfig struct {
	Path     string
	Prefix   string
	FileMode os.FileMode
	Format   auditFormat
}

// NewAuditFileSink should be used to create a new AuditFileSink.
func NewAuditFileSink(config AuditFileSinkConfig) (*AuditFileSink, error) {
	const op = "event.NewAuditFileSink"

	switch config.Format {
	case AuditFormatJSON:
	case AuditFormatJSONX:
	default:
		return nil, fmt.Errorf("%s: unsupported audit format %q", op, config.Format)
	}

	mode := os.FileMode(defaultFileMode)

	if config.FileMode != defaultFileMode {
		switch config.FileMode {
		case 0:
			// if mode is 0000, then do not modify file mode
			if config.Path != "stdout" && config.Path != "discard" && config.Path != "stderr" {
				fileInfo, err := os.Stat(config.Path)
				if err != nil {
					return nil, err
				}
				mode = fileInfo.Mode()
			}
		default:
			mode = config.FileMode
		}
	}

	return &AuditFileSink{
		file:     nil,
		fileLock: sync.RWMutex{},
		fileMode: mode,
		path:     config.Path,
		format:   config.Format,
	}, nil
}

// Process handles writing the event to the file sink.
func (f *AuditFileSink) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 'discard' path means we just do nothing and pretend we're done.
	if f.path == "discard" {
		return nil, nil
	}

	const op = "event.(AuditFileSink).Process"
	var writer io.Writer
	if f.path == "stdout" {
		writer = os.Stdout
	}
	formatted, exists := e.Format(f.format.String())
	if !exists {
		return nil, fmt.Errorf("%s: unable to retrieve formatted event %q", op, f.format)
	}

	buffer := bytes.NewBuffer(formatted)
	err := f.log(buffer, writer)
	if err != nil {
		return nil, fmt.Errorf("%s: error writing file for audit sink: %w", op, err)
	}

	return nil, nil
}

// Reopen handles closing and reopening the file.
func (f *AuditFileSink) Reopen() error {
	const op = "event.(AuditFileSink).Reopen"

	switch f.path {
	case "stdout", "discard":
		return nil
	}

	f.fileLock.Lock()
	defer f.fileLock.Unlock()

	if f.file == nil {
		return f.open()
	}

	err := f.file.Close()
	// Set to nil here so that even if we error out, on the next access open() will be tried.
	f.file = nil
	if err != nil {
		return fmt.Errorf("%s: unable to close file for re-opening on audit sink: %w", op, err)
	}

	return f.open()
}

// Type is used to define which type of node AuditFileSink is.
func (f *AuditFileSink) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}

// open attempts to open a file at the sink's path, with the sink's fileMode permissions
// if one is not already open.
func (f *AuditFileSink) open() error {
	const op = "event.(AuditFileSink).open"

	if f.file != nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(f.path), f.fileMode); err != nil {
		return fmt.Errorf("%s: unable to create directory %q for audit sink: %w", op, f.path, err)
	}

	var err error
	f.file, err = os.OpenFile(f.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, f.fileMode)
	if err != nil {
		return fmt.Errorf("%s: unable to open file for audit sink: %w", op, err)
	}

	// Change the file mode in case the log file already existed.
	// We special case /dev/null since we can't chmod it and bypass if the mode is zero.
	if f.path != "/dev/null" && f.fileMode != 0 {
		if err = os.Chmod(f.path, f.fileMode); err != nil {
			return fmt.Errorf("%s: unable to change file %q permissions '%v' for audit sink: %w", op, f.path, f.fileMode, err)
		}
	}

	return nil
}

// log writes the buffer to the file.
func (f *AuditFileSink) log(buf *bytes.Buffer, writer io.Writer) error {
	const op = "event.(AuditFileSink).log"

	reader := bytes.NewReader(buf.Bytes())

	f.fileLock.Lock()
	defer f.fileLock.Unlock()

	if writer == nil {
		if err := f.open(); err != nil {
			return fmt.Errorf("%s: unable to open file for audit sink: %w", op, err)
		}
		writer = f.file
	}

	if f.prefix != "" {
		_, err := writer.Write([]byte(f.prefix))
		if err != nil {
			return fmt.Errorf("%s: unable to write prefix %q for audit sink: %w", op, f.prefix, err)
		}
	}

	if _, err := reader.WriteTo(writer); err == nil {
		return nil
	} else if f.path == "stdout" {
		return fmt.Errorf("%s: unable write to %q: %w", op, f.path, err)
	}

	// If writing to stdout there's no real reason to think anything would have changed so return above.
	// Otherwise, opportunistically try to re-open the FD, once per call.
	err := f.file.Close()
	if err != nil {
		return fmt.Errorf("%s: unable to close file for audit sink: %w", op, err)
	}

	f.file = nil

	if err := f.open(); err != nil {
		return fmt.Errorf("%s: unable to re-open file for audit sink: %w", op, err)
	}

	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("%s: unable to seek to start of file for audit sink: %w", op, err)
	}

	_, err = reader.WriteTo(writer)
	return fmt.Errorf("%s: unable to re-write to file for audit sink: %w", op, err)
}
