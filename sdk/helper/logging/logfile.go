package logging

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var now = time.Now

type LogFile struct {
	// Name of the log file
	fileName string

	// Path to the log file
	logPath string

	// LastCreated represents the creation time of the latest log
	LastCreated time.Time

	// FileInfo is the pointer to the current file being written to
	FileInfo *os.File

	// BytesWritten is the number of bytes written in the current log file
	BytesWritten int64

	// acquire is the mutex utilized to ensure we have no concurrency issues
	acquire sync.Mutex
}

// Write is used to implement io.Writer
func (l *LogFile) Write(b []byte) (n int, err error) {
	l.acquire.Lock()
	defer l.acquire.Unlock()
	// Create a new file if we have no file to write to
	if l.FileInfo == nil {
		if err := l.openNew(); err != nil {
			return 0, err
		}
	}

	l.BytesWritten += int64(len(b))
	return l.FileInfo.Write(b)
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
	createTime := now()
	newFilePath := filepath.Join(l.logPath, l.fileName)

	// Try creating a file. We truncate the file because we are the only authority to write the logs
	filePointer, err := os.OpenFile(newFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}

	l.FileInfo = filePointer
	// New file, new bytes tracker, new creation time :)
	l.LastCreated = createTime
	l.BytesWritten = 0
	return nil
}
