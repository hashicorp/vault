// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"

	"github.com/hashicorp/eventlogger"

	"github.com/stretchr/testify/require"
)

// TestFileSink_Type ensures that the node is a 'sink' type.
func TestFileSink_Type(t *testing.T) {
	f, err := NewFileSink(t.TempDir(), "json")
	require.NoError(t, err)
	require.NotNil(t, f)
	require.Equal(t, eventlogger.NodeTypeSink, f.Type())
}

// TestNewFileSink tests creation of an AuditFileSink.
func TestNewFileSink(t *testing.T) {
	tests := map[string]struct {
		IsTempDirPath        bool // Path should contain the filename if temp dir is true
		Path                 string
		Format               string
		Options              []Option
		IsErrorExpected      bool
		ExpectedErrorMessage string
		// Expected values of AuditFileSink
		ExpectedFileMode os.FileMode
		ExpectedFormat   string
		ExpectedPath     string
		ExpectedPrefix   string
	}{
		"default-values": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.NewFileSink: path is required",
		},
		"spacey-path": {
			Path:                 "     ",
			Format:               "json",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.NewFileSink: path is required",
		},
		"path-not-exist-valid-format-file-mode": {
			Path:             "qwerty",
			Format:           "json",
			Options:          []Option{WithFileMode("00755")},
			IsErrorExpected:  false,
			ExpectedPath:     "qwerty",
			ExpectedFormat:   "json",
			ExpectedPrefix:   "",
			ExpectedFileMode: os.FileMode(0o755),
		},
		"valid-path-and-format": {
			IsTempDirPath:    true,
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
			IsTempDirPath:    true,
			IsErrorExpected:  false,
			ExpectedFormat:   "json",
			ExpectedPrefix:   "",
			ExpectedFileMode: 0o007,
		},
		"prefix": {
			IsTempDirPath:    true,
			Path:             "vault.log",
			Format:           "json",
			Options:          []Option{WithFileMode("0007"), WithPrefix("bleep")},
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
			if tc.IsTempDirPath {
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
				require.Equal(t, tc.ExpectedPrefix, sink.prefix)
				require.Equal(t, tc.ExpectedFormat, sink.requiredFormat)
				require.Equal(t, tc.ExpectedFileMode, sink.fileMode)

				switch {
				case tc.IsTempDirPath:
					require.Equal(t, tempPath, sink.path)
				default:
					require.Equal(t, tc.ExpectedPath, sink.path)
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
		Path                 string
		IsTempDirPath        bool
		ShouldCreateFile     bool
		Options              []Option
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedFileMode     os.FileMode
	}{
		// Should be ignored by Reopen
		"discard": {
			Path: "discard",
		},
		// Should be ignored by Reopen
		"stdout": {
			Path: "stdout",
		},
		"permission-denied": {
			Path:                 "/tmp/vault/test/foo.log",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(FileSink).open: unable to create file \"/tmp/vault/test/foo.log\": mkdir /tmp/vault/test: permission denied",
		},
		"happy": {
			Path:             "vault.log",
			IsTempDirPath:    true,
			ExpectedFileMode: os.FileMode(defaultFileMode),
		},
		"filemode-existing": {
			Path:             "vault.log",
			IsTempDirPath:    true,
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
			if tc.IsTempDirPath {
				tempDir = t.TempDir()
				tempPath = filepath.Join(tempDir, tc.Path)
			}

			// If the file mode is 0 then we will need a pre-created file to stat.
			// Only do this for paths that are not 'special keywords'
			if tc.ShouldCreateFile && tc.Path != discard && tc.Path != stdout {
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
			case tempPath == discard:
				require.NoError(t, err)
			case tempPath == stdout:
				require.NoError(t, err)
			default:
				require.NoError(t, err)
				info, err := os.Stat(tempPath)
				require.NoError(t, err)
				require.NotNil(t, info)
				require.Equal(t, tc.ExpectedFileMode, info.Mode())
			}
		})
	}
}

// TestFileSink_Process ensures that Process behaves as expected.
func TestFileSink_Process(t *testing.T) {
	tests := map[string]struct {
		Path                 string
		Format               string
		Data                 string
		ShouldIgnoreFormat   bool
		ShouldUseNilEvent    bool
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"discard": {
			Path:            discard,
			Format:          "json",
			Data:            "foo",
			IsErrorExpected: false,
		},
		"stdout": {
			Path:            stdout,
			Format:          "json",
			Data:            "foo",
			IsErrorExpected: false,
		},
		"no-formatted-data": {
			Path:                 "/dev/null",
			Format:               "json",
			Data:                 "foo",
			ShouldIgnoreFormat:   true,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(FileSink).Process: unable to retrieve event formatted as \"json\"",
		},
		"nil": {
			Path:                 "/dev/null",
			Format:               "json",
			Data:                 "foo",
			ShouldUseNilEvent:    true,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(FileSink).Process: event is nil: invalid parameter",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			// Set up a sink
			sink, err := NewFileSink(tc.Path, tc.Format)
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
