package logging

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type LogFile struct {
	// Name of the log file
	fileName string

	// Path to the log file
	logPath string

	// fileInfo is the pointer to the current file being written to
	fileInfo *os.File

	// acquire is the mutex utilized to ensure we have no concurrency issues
	acquire sync.Mutex
}

func NewLogFile(logPath string, fileName string) *LogFile {
	return &LogFile{
		fileName: strings.TrimSpace(fileName),
		logPath:  strings.TrimSpace(logPath),
	}
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

	return l.fileInfo.Write(b)
}

func (l *LogFile) openNew() error {
	newFilePath := filepath.Join(l.logPath, l.fileName)

	// Try to open an existing file or create a new one if it doesn't exist.
	filePointer, err := os.OpenFile(newFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}

	l.fileInfo = filePointer
	return nil
}
