// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventlogger

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	stdout  = "/dev/stdout"
	stderr  = "/dev/stderr"
	devnull = "/dev/null"
)

// FileSink writes the []byte representation of an Event to a file
// as a string.
type FileSink struct {
	// Path is the complete path of the log file directory, excluding FileName
	Path string

	// FileName is the name of the log file
	FileName string

	// Mode is the file's mode and permission bits
	Mode os.FileMode

	// LastCreated represents the creation time of the latest log
	LastCreated time.Time

	// MaxBytes is the maximum number of desired bytes for a log file
	MaxBytes int

	// BytesWritten is the number of bytes written in the current log file
	BytesWritten int64

	// MaxFiles is the maximum number of old files to keep before removing them
	MaxFiles int

	// MaxDuration is the maximum duration allowed between each file rotation
	MaxDuration time.Duration

	// Format specifies the format the []byte representation is formatted in
	// Defaults to JSONFormat
	Format string

	// TimestampOnlyOnRotate specifies the file currently being written
	// should not contain a timestamp in the name even if rotation is
	// enabled.
	//
	// If false (the default) all files, including the currently written
	// one, will contain a timestamp in the filename.
	TimestampOnlyOnRotate bool

	f *os.File
	l sync.Mutex
}

var _ Node = &FileSink{}

const (
	defaultMode = 0600
	dirMode     = 0700
)

// Type describes the type of the node as a Sink.
func (_ *FileSink) Type() NodeType {
	return NodeTypeSink
}

// Process writes the []byte representation of an Event to a file
// as a string.
func (fs *FileSink) Process(_ context.Context, e *Event) (*Event, error) {
	// '/dev/null' should just return success
	if fs.Path == devnull {
		return nil, nil
	}

	format := fs.Format
	if format == "" {
		format = JSONFormat
	}
	val, ok := e.Format(format)
	if !ok {
		return nil, errors.New("event was not marshaled")
	}

	reader := bytes.NewReader(val)

	fs.l.Lock()
	defer fs.l.Unlock()

	var writer io.Writer
	switch fs.Path {
	case stdout:
		writer = os.Stdout
	case stderr:
		writer = os.Stderr
	default:
		if fs.f == nil {
			err := fs.open()
			if err != nil {
				return nil, err
			}
		}
		// Check for last contact, rotate if necessary and able
		if err := fs.rotate(); err != nil {
			return nil, err
		}
		writer = fs.f
	}

	if n, err := reader.WriteTo(writer); err == nil {
		// Sinks are leafs, so do not return the event, since nothing more can
		// happen to it downstream.
		fs.BytesWritten += n
		return nil, nil
	}

	// Since we haven't returned yet, we assume that the attempt to write didn't
	// succeed, and we probably weren't attempting to write to a special path
	// such as: /dev/null, /dev/stdout or /dev/stderr.
	// Attempt a single 'retry' once per call.
	if err := fs.reopen(); err != nil {
		return nil, err
	}

	_, _ = reader.Seek(0, io.SeekStart)
	_, err := reader.WriteTo(fs.f)
	return nil, err
}

// reopen will close, rotate and reopen the Sink's file.
// NOTE: this method is to be called by exported FileSink receivers which must
// handle obtaining the relevant lock on the struct.
func (fs *FileSink) reopen() error {
	switch fs.Path {
	case stdout, stderr, devnull:
		return nil
	}

	if fs.f != nil {
		// Ensure file still exists
		_, err := os.Stat(fs.f.Name())
		if os.IsNotExist(err) {
			fs.f = nil
		}
	}

	if fs.f == nil {
		return fs.open()
	}

	err := fs.f.Close()
	// Set to nil here so that even if we error out, on the next access open()
	// will be tried
	fs.f = nil
	if err != nil {
		return err
	}

	return fs.open()
}

// Reopen will close, rotate and reopen the Sink's file.
func (fs *FileSink) Reopen() error {
	switch fs.Path {
	case stdout, stderr, devnull:
		return nil
	}

	fs.l.Lock()
	defer fs.l.Unlock()

	return fs.reopen()
}

// Name returns a representation of the Sink's name
func (fs *FileSink) Name() string {
	return fmt.Sprintf("sink:%s", fs.Path)
}

func (fs *FileSink) open() error {
	// Return early if the file is open, or we're using a special path.
	switch fs.Path {
	case devnull, stdout, stderr:
		return nil
	default:
		if fs.f != nil {
			return nil
		}
	}

	mode := fs.Mode
	if mode == 0 {
		mode = defaultMode
	}

	if err := os.MkdirAll(fs.Path, dirMode); err != nil {
		return err
	}

	createTime := time.Now()
	// New file name as the format:
	// file rotation enabled: filename-timestamp.extension
	// file rotation disabled: filename.extension
	newFileName := fs.newFileName(createTime)
	newFilePath := filepath.Join(fs.Path, newFileName)

	var err error
	fs.f, err = os.OpenFile(newFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, mode)
	if err != nil {
		return err
	}

	// Change the file mode (if not 0) in case the log file already existed.
	if fs.Mode != 0 {
		err = os.Chmod(newFilePath, fs.Mode)
		if err != nil {
			return err
		}
	}

	// Reset file related statistics
	fs.LastCreated = createTime
	fs.BytesWritten = 0

	return nil
}

func (fs *FileSink) rotate() error {
	switch fs.Path {
	case stdout, stderr, devnull:
		return nil
	}

	// Get the time from the last point of contact
	elapsed := time.Since(fs.LastCreated)
	if (fs.BytesWritten >= int64(fs.MaxBytes) && (fs.MaxBytes > 0)) ||
		((elapsed > fs.MaxDuration) && (fs.MaxDuration > 0)) {

		// Clean up the existing file
		err := fs.f.Close()
		if err != nil {
			return err
		}
		fs.f = nil

		// Move current log file to a timestamped file.
		if fs.TimestampOnlyOnRotate {
			rotateTime := time.Now().UnixNano()
			rotateFileName := fmt.Sprintf(fs.fileNamePattern(), strconv.FormatInt(rotateTime, 10))
			oldPath := filepath.Join(fs.Path, fs.FileName)
			newPath := filepath.Join(fs.Path, rotateFileName)
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("failed to rotate log file: %v", err)
			}
		}

		if err := fs.pruneFiles(); err != nil {
			return fmt.Errorf("failed to prune log files: %w", err)
		}
		return fs.open()
	}

	return nil
}

func (fs *FileSink) pruneFiles() error {
	switch {
	case fs.Path == stdout, fs.Path == stderr, fs.Path == devnull:
		return nil
	case fs.MaxFiles == 0:
		return nil
	}

	// get all the files that match the log file pattern
	pattern := fs.fileNamePattern()
	globExpression := filepath.Join(fs.Path, fmt.Sprintf(pattern, "*"))
	matches, err := filepath.Glob(globExpression)
	if err != nil {
		return err
	}

	// Stort the strings as filepath.Glob does not publicly guarantee that files
	// are sorted, so here we add an extra defensive sort.
	sort.Strings(matches)

	stale := len(matches) - fs.MaxFiles
	for i := 0; i < stale; i++ {
		if err := os.Remove(matches[i]); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileSink) fileNamePattern() string {
	// Extract file extension
	ext := filepath.Ext(fs.FileName)
	if ext == "" {
		ext = ".log"
	}

	// Add format string between file and extension
	return strings.TrimSuffix(fs.FileName, ext) + "-%s" + ext
}

func (fs *FileSink) newFileName(createTime time.Time) string {
	if fs.TimestampOnlyOnRotate {
		return fs.FileName
	}

	if !fs.rotateEnabled() {
		return fs.FileName
	}

	pattern := fs.fileNamePattern()
	return fmt.Sprintf(pattern, strconv.FormatInt(createTime.UnixNano(), 10))
}

func (fs *FileSink) rotateEnabled() bool {
	return fs.MaxBytes > 0 || fs.MaxDuration != 0
}
