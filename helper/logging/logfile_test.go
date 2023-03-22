// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogFile_openNew(t *testing.T) {
	logFile := &LogFile{
		fileName: "vault.log",
		logPath:  t.TempDir(),
		duration: defaultRotateDuration,
	}

	err := logFile.openNew()
	require.NoError(t, err)

	msg := "[INFO] Something"
	_, err = logFile.Write([]byte(msg))
	require.NoError(t, err)

	content, err := os.ReadFile(logFile.fileInfo.Name())
	require.NoError(t, err)
	require.Contains(t, string(content), msg)
}

func TestLogFile_Rotation_MaxDuration(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	tempDir := t.TempDir()
	logFile := LogFile{
		fileName: "vault.log",
		logPath:  tempDir,
		duration: 50 * time.Millisecond,
	}

	_, err := logFile.Write([]byte("Hello World"))
	assert.NoError(t, err, "error writing rotation max duration part 1")

	time.Sleep(3 * logFile.duration)

	_, err = logFile.Write([]byte("Second File"))
	assert.NoError(t, err, "error writing rotation max duration part 2")

	require.Len(t, listDir(t, tempDir), 2)
}

func TestLogFile_Rotation_MaxBytes(t *testing.T) {
	tempDir := t.TempDir()
	logFile := LogFile{
		fileName: "somefile.log",
		logPath:  tempDir,
		maxBytes: 10,
		duration: defaultRotateDuration,
	}
	_, err := logFile.Write([]byte("Hello World"))
	assert.NoError(t, err, "error writing rotation max bytes part 1")

	_, err = logFile.Write([]byte("Second File"))
	assert.NoError(t, err, "error writing rotation max bytes part 2")

	require.Len(t, listDir(t, tempDir), 2)
}

func TestLogFile_PruneFiles(t *testing.T) {
	tempDir := t.TempDir()
	logFile := LogFile{
		fileName:         "vault.log",
		logPath:          tempDir,
		maxBytes:         10,
		duration:         defaultRotateDuration,
		maxArchivedFiles: 1,
	}
	_, err := logFile.Write([]byte("[INFO] Hello World"))
	assert.NoError(t, err, "error writing during prune files test part 1")

	_, err = logFile.Write([]byte("[INFO] Second File"))
	assert.NoError(t, err, "error writing during prune files test part 1")

	_, err = logFile.Write([]byte("[INFO] Third File"))
	assert.NoError(t, err, "error writing during prune files test part 1")

	logFiles := listDir(t, tempDir)
	sort.Strings(logFiles)
	require.Len(t, logFiles, 2)

	content, err := os.ReadFile(filepath.Join(tempDir, logFiles[0]))
	require.NoError(t, err)
	require.Contains(t, string(content), "Second File")

	content, err = os.ReadFile(filepath.Join(tempDir, logFiles[1]))
	require.NoError(t, err)
	require.Contains(t, string(content), "Third File")
}

func TestLogFile_PruneFiles_Disabled(t *testing.T) {
	tempDir := t.TempDir()
	logFile := LogFile{
		fileName:         "somename.log",
		logPath:          tempDir,
		maxBytes:         10,
		duration:         defaultRotateDuration,
		maxArchivedFiles: 0,
	}

	_, err := logFile.Write([]byte("[INFO] Hello World"))
	assert.NoError(t, err, "error writing during prune files - disabled test part 1")

	_, err = logFile.Write([]byte("[INFO] Second File"))
	assert.NoError(t, err, "error writing during prune files - disabled test part 2")

	_, err = logFile.Write([]byte("[INFO] Third File"))
	assert.NoError(t, err, "error writing during prune files - disabled test part 3")

	require.Len(t, listDir(t, tempDir), 3)
}

func TestLogFile_FileRotation_Disabled(t *testing.T) {
	tempDir := t.TempDir()
	logFile := LogFile{
		fileName:         "vault.log",
		logPath:          tempDir,
		maxBytes:         10,
		maxArchivedFiles: -1,
	}

	_, err := logFile.Write([]byte("[INFO] Hello World"))
	assert.NoError(t, err, "error writing during rotation disabled test part 1")

	_, err = logFile.Write([]byte("[INFO] Second File"))
	assert.NoError(t, err, "error writing during rotation disabled test part 2")

	_, err = logFile.Write([]byte("[INFO] Third File"))
	assert.NoError(t, err, "error writing during rotation disabled test part 3")

	require.Len(t, listDir(t, tempDir), 1)
}

func listDir(t *testing.T, name string) []string {
	t.Helper()
	fh, err := os.Open(name)
	require.NoError(t, err)
	files, err := fh.Readdirnames(100)
	require.NoError(t, err)
	return files
}
