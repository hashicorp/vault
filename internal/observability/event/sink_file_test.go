// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFileSink_Type ensures that the node is a 'sink' type.
func TestFileSink_Type(t *testing.T) {
	f, err := NewFileSink(filepath.Join(t.TempDir(), "vault.log"), "json")
	require.NoError(t, err)
	require.NotNil(t, f)
	require.Equal(t, eventlogger.NodeTypeSink, f.Type())
}

// TestNewFileSink tests creation of an AuditFileSink.
func TestNewFileSink(t *testing.T) {
	tests := map[string]struct {
		ShouldUseAbsolutePath bool // Path should contain the filename if temp dir is true
		Path                  string
		Format                string
		Options               []Option
		IsErrorExpected       bool
		ExpectedErrorMessage  string
		// Expected values of AuditFileSink
		ExpectedFileMode os.FileMode
		ExpectedFormat   string
		ExpectedPath     string
		ExpectedPrefix   string
	}{
		"default-values": {
			ShouldUseAbsolutePath: true,
			IsErrorExpected:       true,
			ExpectedErrorMessage:  "path is required: invalid parameter",
		},
		"spacey-path": {
			ShouldUseAbsolutePath: true,
			Path:                  "     ",
			Format:                "json",
			IsErrorExpected:       true,
			ExpectedErrorMessage:  "path is required: invalid parameter",
		},
		"valid-path-and-format": {
			Path:             "vault.log",
			Format:           "json",
			IsErrorExpected:  false,
			ExpectedFileMode: defaultFileMode,
			ExpectedFormat:   "json",
			ExpectedPrefix:   "",
		},
		"file-mode-not-default-or-zero": {
			Path:             "vault.log",
			Format:           "json",
			Options:          []Option{WithFileMode("0007")},
			IsErrorExpected:  false,
			ExpectedFormat:   "json",
			ExpectedPrefix:   "",
			ExpectedFileMode: 0o007,
		},
		"prefix": {
			Path:             "vault.log",
			Format:           "json",
			Options:          []Option{WithFileMode("0007")},
			IsErrorExpected:  false,
			ExpectedPrefix:   "bleep",
			ExpectedFormat:   "json",
			ExpectedFileMode: 0o007,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			// t.Parallel()

			// If we need a real directory as a path we can use a temp dir.
			// but we should keep track of it for comparison in the new sink.
			var tempDir string
			tempPath := tc.Path
			if !tc.ShouldUseAbsolutePath {
				tempDir = t.TempDir()
				tempPath = filepath.Join(tempDir, tempPath)
			}

			sink, err := NewFileSink(tempPath, tc.Format, tc.Options...)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, sink)
			default:
				require.NoError(t, err)
				require.NotNil(t, sink)

				// Assert properties are correct.
				require.Equal(t, tc.ExpectedFormat, sink.requiredFormat)
				require.Equal(t, tc.ExpectedFileMode, sink.fileMode)

				switch {
				case tc.ShouldUseAbsolutePath:
					require.Equal(t, tc.ExpectedPath, sink.path)
				default:
					require.Equal(t, tempPath, sink.path)
				}
			}
		})
	}
}

// TestFileSink_Reopen tests that the sink reopens files as expected when requested to.
// stdout and discard paths are ignored.
// see: https://developer.hashicorp.com/vault/docs/audit/file#file_path
func TestFileSink_Reopen(t *testing.T) {
	tests := map[string]struct {
		Path                  string
		ShouldUseAbsolutePath bool
		ShouldCreateFile      bool
		ShouldIgnoreFileMode  bool
		Options               []Option
		IsErrorExpected       bool
		ExpectedErrorMessage  string
		ExpectedFileMode      os.FileMode
	}{
		// Should be ignored by Reopen
		"stdout": {
			Path:                  stdout,
			ShouldUseAbsolutePath: true,
			ShouldIgnoreFileMode:  true,
		},
		"stderr": {
			Path:                  stderr,
			ShouldUseAbsolutePath: true,
			ShouldIgnoreFileMode:  true,
		},
		"devnull": {
			Path:                  devnull,
			ShouldUseAbsolutePath: true,
			ShouldIgnoreFileMode:  true,
		},
		"happy": {
			Path:             "vault.log",
			ExpectedFileMode: os.FileMode(defaultFileMode),
		},
		"filemode-existing": {
			Path:             "vault.log",
			ShouldCreateFile: true,
			Options:          []Option{WithFileMode("0000")},
			ExpectedFileMode: os.FileMode(defaultFileMode),
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// If we need a real directory as a path we can use a temp dir.
			// but we should keep track of it for comparison in the new sink.
			var tempDir string
			tempPath := tc.Path
			if !tc.ShouldUseAbsolutePath {
				tempDir = t.TempDir()
				tempPath = filepath.Join(tempDir, tc.Path)
			}

			// If the file mode is 0 then we will need a pre-created file to stat.
			// Only do this for paths that are not 'special keywords'
			if tc.ShouldCreateFile && tc.Path != devnull {
				f, err := os.OpenFile(tempPath, os.O_CREATE, defaultFileMode)
				require.NoError(t, err)
				defer func() {
					err = os.Remove(f.Name())
					require.NoError(t, err)
				}()
			}

			sink, err := NewFileSink(tempPath, "json", tc.Options...)
			require.NoError(t, err)
			require.NotNil(t, sink)

			err = sink.Reopen()

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				info, err := os.Stat(tempPath)
				require.NoError(t, err)
				require.NotNil(t, info)
				if !tc.ShouldIgnoreFileMode {
					require.Equal(t, tc.ExpectedFileMode, info.Mode())
				}
			}
		})
	}
}

// TestFileSink_Process ensures that Process behaves as expected.
func TestFileSink_Process(t *testing.T) {
	tests := map[string]struct {
		ShouldUseAbsolutePath bool
		Path                  string
		ShouldCreateFile      bool
		Format                string
		ShouldIgnoreFormat    bool
		Data                  string
		ShouldUseNilEvent     bool
		IsErrorExpected       bool
		ExpectedErrorMessage  string
	}{
		"stdout": {
			ShouldUseAbsolutePath: true,
			Path:                  stdout,
			Format:                "json",
			Data:                  "foo",
			IsErrorExpected:       false,
		},
		"stderr": {
			ShouldUseAbsolutePath: true,
			Path:                  stderr,
			Format:                "json",
			Data:                  "foo",
			IsErrorExpected:       false,
		},
		"devnull": {
			ShouldUseAbsolutePath: true,
			Path:                  devnull,
			Format:                "json",
			Data:                  "foo",
			IsErrorExpected:       false,
		},
		"no-formatted-data": {
			ShouldCreateFile:     true,
			Path:                 "juan.log",
			Format:               "json",
			Data:                 "foo",
			ShouldIgnoreFormat:   true,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to retrieve event formatted as \"json\": invalid parameter",
		},
		"nil": {
			Path:                 "foo.log",
			Format:               "json",
			Data:                 "foo",
			ShouldUseNilEvent:    true,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event is nil: invalid parameter",
		},
		"happy-path": {
			Path:             "juan.log",
			ShouldCreateFile: true,
			Format:           "json",
			Data:             "{\"foo\": \"bar\"}",
			IsErrorExpected:  false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			// Temp dir for most testing unless we're trying to test an error
			var tempDir string
			tempPath := tc.Path
			if !tc.ShouldUseAbsolutePath {
				tempDir = t.TempDir()
				tempPath = filepath.Join(tempDir, tc.Path)
			}

			// Create a file if we will need it there before Process kicks off.
			if tc.ShouldCreateFile && tc.Path != devnull {
				f, err := os.OpenFile(tempPath, os.O_CREATE, defaultFileMode)
				require.NoError(t, err)
				defer func() {
					err = os.Remove(f.Name())
					require.NoError(t, err)
				}()
			}

			// Set up a sink
			sink, err := NewFileSink(tempPath, tc.Format)
			require.NoError(t, err)
			require.NotNil(t, sink)

			// Generate a fake event
			ctx := namespace.RootContext(nil)

			event := &eventlogger.Event{
				Type:      "audit",
				CreatedAt: time.Now(),
				Formatted: make(map[string][]byte),
				Payload:   struct{ ID string }{ID: "123"},
			}

			if !tc.ShouldIgnoreFormat {
				event.FormattedAs(tc.Format, []byte(tc.Data))
			}

			if tc.ShouldUseNilEvent {
				event = nil
			}

			// The actual exercising of the sink.
			event, err = sink.Process(ctx, event)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, event)
			default:
				require.NoError(t, err)
				require.Nil(t, event)
			}
		})
	}
}

// TestFileSink_log_cancelledContext tests that 'log' is context aware and won't
// just wait for the lock on the file sink forever.
func TestFileSink_log_cancelledContext(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "juan.log")
	data := []byte("{\"foo\": \"bar\"}")

	sink, err := NewFileSink(tempPath, "json")
	require.NoError(t, err)
	require.NotNil(t, sink)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() { cancel() })

	// Manually acquire the lock so that the 'log' method cannot.
	sink.fileLock.Lock()

	var done atomic.Bool
	go func() {
		err = sink.log(ctx, data)
		done.Store(true)
	}()

	// We shouldn't have an error as 'log' will be trying to acquire the lock.
	require.NoError(t, err)

	// Manually cancel the context and unlock to let the waiting 'log' in.
	cancel()
	sink.fileLock.Unlock()

	// Just a little bit of time to make sure that 'log' returned and err was set.
	corehelpers.RetryUntil(t, 3*time.Second, func() error {
		if done.Load() {
			return nil
		}

		return fmt.Errorf("logging still not done")
	})

	// We expect that the error now has context cancelled in it.
	require.True(t, errors.Is(err, context.Canceled))
}

func TestFileSink_Rotation_MaxDuration(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "vault.log")
	options := []Option{WithMaxDurationFile("50ms")}

	sink, err := NewFileSink(tempPath, "json", options...)
	require.NoError(t, err)
	require.NotNil(t, sink)

	ctx := namespace.RootContext(nil)

	event := &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("first entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing rotation max duration part 1")

	time.Sleep(3 * sink.maxDuration)

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("second entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing rotation max duration part 2")

	logFiles := listDir(t, tempDir)
	require.Len(t, logFiles, 2)
	require.Equal(t, "vault.log", logFiles[1])

	content, err := os.ReadFile(filepath.Join(tempDir, logFiles[0]))
	require.NoError(t, err)
	require.Equal(t, []byte("first entry"), content)

	content, err = os.ReadFile(filepath.Join(tempDir, logFiles[1]))
	require.NoError(t, err)
	require.Equal(t, []byte("second entry"), content)
}

func TestFileSink_Rotation_MaxBytes(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "vault.log")
	options := []Option{WithMaxBytes("10"), WithMaxDurationFile("1d")}

	sink, err := NewFileSink(tempPath, "json", options...)
	require.NoError(t, err)
	require.NotNil(t, sink)

	ctx := namespace.RootContext(nil)

	event := &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("first entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing rotation max bytes part 1")

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("second entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing rotation max bytes part 2")

	logFiles := listDir(t, tempDir)
	require.Len(t, logFiles, 2)
	require.Equal(t, "vault.log", logFiles[1])

	content, err := os.ReadFile(filepath.Join(tempDir, logFiles[0]))
	require.NoError(t, err)
	require.Equal(t, []byte("first entry"), content)

	content, err = os.ReadFile(filepath.Join(tempDir, logFiles[1]))
	require.NoError(t, err)
	require.Equal(t, []byte("second entry"), content)
}

func TestFileSink_PruneFiles(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "vault.log")
	options := []Option{WithMaxFiles("1"), WithMaxBytes("10"), WithMaxDurationFile("1d")}

	sink, err := NewFileSink(tempPath, "json", options...)
	require.NoError(t, err)
	require.NotNil(t, sink)

	ctx := namespace.RootContext(nil)

	event := &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("first entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during prune files test part 1")

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("second entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during prune files test part 2")

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("third entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during prune files test part 3")

	logFiles := listDir(t, tempDir)
	require.Len(t, logFiles, 2)
	require.Equal(t, "vault.log", logFiles[1])

	content, err := os.ReadFile(filepath.Join(tempDir, logFiles[0]))
	require.NoError(t, err)
	require.Equal(t, []byte("second entry"), content)

	content, err = os.ReadFile(filepath.Join(tempDir, logFiles[1]))
	require.NoError(t, err)
	require.Equal(t, []byte("third entry"), content)
}

func TestFileSink_PruneFiles_Disabled(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "vault.log")
	options := []Option{WithMaxFiles("0"), WithMaxBytes("10"), WithMaxDurationFile("1d")}

	sink, err := NewFileSink(tempPath, "json", options...)
	require.NoError(t, err)
	require.NotNil(t, sink)

	ctx := namespace.RootContext(nil)

	event := &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("first entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during prune files - disabled test part 1")

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("second entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during prune files - disabled test part 2")

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("third entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during prune files - disabled test part 3")

	require.Len(t, listDir(t, tempDir), 3)
}

func TestFileSink_PruneFiles_Always(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "vault.log")
	options := []Option{WithMaxFiles("-1"), WithMaxBytes("10"), WithMaxDurationFile("1d")}

	sink, err := NewFileSink(tempPath, "json", options...)
	require.NoError(t, err)
	require.NotNil(t, sink)

	ctx := namespace.RootContext(nil)

	event := &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("first entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during prune files - always test part 1")

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("second entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during prune files - always test part 2")

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("third entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during prune files - always test part 3")

	logFiles := listDir(t, tempDir)
	require.Len(t, logFiles, 1)
	require.Equal(t, "vault.log", logFiles[0])

	content, err := os.ReadFile(filepath.Join(tempDir, logFiles[0]))
	require.NoError(t, err)
	require.Equal(t, []byte("third entry"), content)
}

func TestFileSink_FileRotation_Disabled(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "vault.log")
	options := []Option{WithMaxFiles("0"), WithMaxBytes("0"), WithMaxDurationFile("0")}

	sink, err := NewFileSink(tempPath, "json", options...)
	require.NoError(t, err)
	require.NotNil(t, sink)

	ctx := namespace.RootContext(nil)

	event := &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("first entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during rotation - disabled test part 1")

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("second entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during rotation - disabled test part 2")

	event = &eventlogger.Event{
		Formatted: map[string][]byte{"json": []byte("third entry")},
	}
	_, err = sink.Process(ctx, event)
	assert.NoError(t, err, "error writing during rotation - disabled test part 3")

	logFiles := listDir(t, tempDir)
	require.Len(t, logFiles, 1)
	require.Equal(t, "vault.log", logFiles[0])

	content, err := os.ReadFile(filepath.Join(tempDir, logFiles[0]))
	require.NoError(t, err)
	require.Equal(t, []byte("first entrysecond entrythird entry"), content)
}

func listDir(t *testing.T, name string) []string {
	t.Helper()
	fh, err := os.Open(name)
	require.NoError(t, err)
	files, err := fh.Readdirnames(100)
	require.NoError(t, err)
	sort.Strings(files)
	return files
}
