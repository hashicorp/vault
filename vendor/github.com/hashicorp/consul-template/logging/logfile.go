// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/logutils"
)

type LogFile struct {
	// Name of the log file
	fileName string

	// Path to the log file
	logPath string

	// Duration between each file rotation operation
	duration time.Duration

	// LastCreated represents the creation time of the latest log
	LastCreated time.Time

	// FileInfo is the pointer to the current file being written to
	FileInfo *os.File

	// MaxBytes is the maximum number of desired bytes for a log file
	MaxBytes int

	// BytesWritten is the number of bytes written in the current log file
	BytesWritten int64

	// Max rotated files to keep before removing them.
	MaxFiles int

	// filt is used to filter log messages depending on their level
	filt *logutils.LevelFilter

	// acquire is the mutex utilized to ensure we have no concurrency issues
	acquire sync.Mutex
}

func (l *LogFile) fileNamePattern() string {
	// Extract the file extension
	fileExt := filepath.Ext(l.fileName)
	// If we have no file extension we append .log
	if fileExt == "" {
		fileExt = ".log"
	}
	// Remove the file extension from the filename
	return strings.TrimSuffix(l.fileName, fileExt) + "-%s" + fileExt
}

func (l *LogFile) openNew() error {
	fileNamePattern := l.fileNamePattern()

	createTime := time.Now()
	newfileName := fmt.Sprintf(fileNamePattern, strconv.FormatInt(createTime.UnixNano(), 10))
	newfilePath := filepath.Join(l.logPath, newfileName)

	// Try creating a file. We truncate the file because we are the only authority to write the logs
	filePointer, err := os.OpenFile(newfilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}

	l.FileInfo = filePointer
	// New file, new bytes tracker, new creation time :)
	l.LastCreated = createTime
	l.BytesWritten = 0
	return nil
}

func (l *LogFile) rotate() error {
	// Get the time from the last point of contact
	timeElapsed := time.Since(l.LastCreated)
	// Rotate if we hit the byte file limit or the time limit
	if (l.BytesWritten >= int64(l.MaxBytes) && (l.MaxBytes > 0)) || timeElapsed >= l.duration {
		l.FileInfo.Close()
		if err := l.pruneFiles(); err != nil {
			return err
		}
		return l.openNew()
	}
	return nil
}

func (l *LogFile) pruneFiles() error {
	if l.MaxFiles == 0 {
		return nil
	}

	pattern := filepath.Join(l.logPath, fmt.Sprintf(l.fileNamePattern(), "*"))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	switch {
	case l.MaxFiles < 0:
		return removeFiles(matches)
	case len(matches) < l.MaxFiles:
		return nil
	}

	sort.Strings(matches)
	last := len(matches) - l.MaxFiles
	return removeFiles(matches[:last])
}

func removeFiles(files []string) error {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}

// Write is used to implement io.Writer.
func (l *LogFile) Write(b []byte) (int, error) {
	l.acquire.Lock()
	defer l.acquire.Unlock()

	// Skip if the log level doesn't apply
	if l.filt != nil && !l.filt.Check(b) {
		return 0, nil
	}

	// Create a new file if we have no file to write to
	if l.FileInfo == nil {
		if err := l.openNew(); err != nil {
			return 0, err
		}
	}
	// Check for the last contact and rotate if necessary
	if err := l.rotate(); err != nil {
		return 0, err
	}
	l.BytesWritten += int64(len(b))
	return l.FileInfo.Write(b)
}
