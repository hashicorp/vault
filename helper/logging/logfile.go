// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

	"github.com/hashicorp/go-multierror"
)

var now = time.Now

type LogFile struct {
	// Name of the log file
	fileName string

	// Path to the log file
	logPath string

	// duration between each file rotation operation
	duration time.Duration

	// lastCreated represents the creation time of the latest log
	lastCreated time.Time

	// fileInfo is the pointer to the current file being written to
	fileInfo *os.File

	// maxBytes is the maximum number of desired bytes for a log file
	maxBytes int

	// bytesWritten is the number of bytes written in the current log file
	bytesWritten int64

	// Max rotated files to keep before removing them.
	maxArchivedFiles int

	// acquire is the mutex utilized to ensure we have no concurrency issues
	acquire sync.Mutex
}

// Write is used to implement io.Writer
func (l *LogFile) Write(b []byte) (n int, err error) {
	l.acquire.Lock()
	defer l.acquire.Unlock()

	// Create a new file if we have no file to write to
	if l.fileInfo == nil {
		if err := l.openNew(); err != nil {
			return 0, err
		}
	}
	if err := l.rotate(); err != nil { // Check for the last contact and rotate if necessary
		return 0, err
	}

	bytesWritten, err := l.fileInfo.Write(b)

	if bytesWritten > 0 {
		l.bytesWritten += int64(bytesWritten)
	}

	return bytesWritten, err
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
	newFileName := l.fileName
	newFilePath := filepath.Join(l.logPath, newFileName)

	// Try creating or opening the active log file. Since the active log file
	// always has the same name, append log entries to prevent overwriting
	// previous log data.
	filePointer, err := os.OpenFile(newFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}

	// New file, new bytes tracker, new creation time :)
	l.fileInfo = filePointer
	l.lastCreated = now()
	l.bytesWritten = 0
	return nil
}

func (l *LogFile) rotate() error {
	// Get the time from the last point of contact
	timeElapsed := time.Since(l.lastCreated)
	// Rotate if we hit the byte file limit or the time limit
	if (l.bytesWritten >= int64(l.maxBytes) && (l.maxBytes > 0)) || timeElapsed >= l.duration {
		if err := l.fileInfo.Close(); err != nil {
			return err
		}
		if err := l.renameCurrentFile(); err != nil {
			return err
		}
		if err := l.pruneFiles(); err != nil {
			return err
		}
		return l.openNew()
	}
	return nil
}

func (l *LogFile) pruneFiles() error {
	if l.maxArchivedFiles == 0 {
		return nil
	}

	pattern := filepath.Join(l.logPath, fmt.Sprintf(l.fileNamePattern(), "*"))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	switch {
	case l.maxArchivedFiles < 0:
		return removeFiles(matches)
	case len(matches) < l.maxArchivedFiles:
		return nil
	}

	sort.Strings(matches)
	last := len(matches) - l.maxArchivedFiles
	return removeFiles(matches[:last])
}

func removeFiles(files []string) (err error) {
	for _, file := range files {
		if fileError := os.Remove(file); fileError != nil {
			err = multierror.Append(err, fmt.Errorf("error removing file %s: %v", file, fileError))
		}
	}
	return err
}

func (l *LogFile) renameCurrentFile() error {
	fileNamePattern := l.fileNamePattern()
	createTime := now()
	currentFilePath := filepath.Join(l.logPath, l.fileName)
	oldFileName := fmt.Sprintf(fileNamePattern, strconv.FormatInt(createTime.UnixNano(), 10))
	oldFilePath := filepath.Join(l.logPath, oldFileName)

	return os.Rename(currentFilePath, oldFilePath)
}
