// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/eventlogger"

	"github.com/stretchr/testify/require"
)

// TestAuditFileSink_Type ensures that the node is a 'sink' type.
func TestAuditFileSink_Type(t *testing.T) {
	f, err := NewAuditFileSink(t.TempDir(), AuditFormatJSON)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.Equal(t, eventlogger.NodeTypeSink, f.Type())
}

// TestNewAuditFileSink tests creation of an AuditFileSink.
func TestNewAuditFileSink(t *testing.T) {
	tests := map[string]struct {
		IsTempDirPath        bool // Path should contain the filename if temp dir is true
		Path                 string
		Format               auditFormat
		Options              []Option
		IsErrorExpected      bool
		ExpectedErrorMessage string
		// Expected values of AuditFileSink
		ExpectedFileMode os.FileMode
		ExpectedFormat   auditFormat
		ExpectedPath     string
		ExpectedPrefix   string
	}{
		"default-values": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.NewAuditFileSink: path is required",
		},
		"spacey-path": {
			Path:                 "     ",
			Format:               AuditFormatJSON,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.NewAuditFileSink: path is required",
		},
		"bad-format": {
			Path:                 "qwerty",
			Format:               "squirrels",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.NewAuditFileSink: invalid format: event.(auditFormat).validate: 'squirrels' is not a valid format: invalid parameter",
		},
		"path-not-exist-valid-format-file-mode": {
			Path:             "qwerty",
			Format:           AuditFormatJSON,
			Options:          []Option{WithFileMode("00755")},
			IsErrorExpected:  false,
			ExpectedPath:     "qwerty",
			ExpectedFormat:   AuditFormatJSON,
			ExpectedPrefix:   "",
			ExpectedFileMode: os.FileMode(0o755),
		},
		"valid-path-no-format": {
			IsTempDirPath:        true,
			Path:                 "vault.log",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.NewAuditFileSink: invalid format: event.(auditFormat).validate: '' is not a valid format: invalid parameter",
		},
		"valid-path-and-format": {
			IsTempDirPath:    true,
			Path:             "vault.log",
			Format:           AuditFormatJSON,
			IsErrorExpected:  false,
			ExpectedFileMode: defaultFileMode,
			ExpectedFormat:   AuditFormatJSON,
			ExpectedPrefix:   "",
		},
		"file-mode-not-default-or-zero": {
			Path:             "vault.log",
			Format:           AuditFormatJSON,
			Options:          []Option{WithFileMode("0007")},
			IsTempDirPath:    true,
			IsErrorExpected:  false,
			ExpectedFormat:   AuditFormatJSON,
			ExpectedPrefix:   "",
			ExpectedFileMode: 0o007,
		},
		"path-stdout": {
			Path:             "stdout",
			Format:           AuditFormatJSON,
			Options:          []Option{WithFileMode("0007")}, // Will be ignored as stdout
			IsTempDirPath:    false,
			IsErrorExpected:  false,
			ExpectedPath:     "stdout",
			ExpectedFormat:   AuditFormatJSON,
			ExpectedPrefix:   "",
			ExpectedFileMode: defaultFileMode,
		},
		"path-discard": {
			Path:             "discard",
			Format:           AuditFormatJSON,
			Options:          []Option{WithFileMode("0007")},
			IsTempDirPath:    false,
			IsErrorExpected:  false,
			ExpectedPath:     "discard",
			ExpectedFormat:   AuditFormatJSON,
			ExpectedPrefix:   "",
			ExpectedFileMode: defaultFileMode,
		},
		"prefix": {
			IsTempDirPath:    true,
			Path:             "vault.log",
			Format:           AuditFormatJSON,
			Options:          []Option{WithFileMode("0007"), WithPrefix("bleep")},
			IsErrorExpected:  false,
			ExpectedPrefix:   "bleep",
			ExpectedFormat:   AuditFormatJSON,
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

			sink, err := NewAuditFileSink(tempPath, tc.Format, tc.Options...)

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
				require.Equal(t, tc.ExpectedFormat, sink.format)
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

// TestAuditFileSink_Reopen tests that the sink reopens files as expected when requested to.
// stdout and discard paths are ignored.
// see: https://developer.hashicorp.com/vault/docs/audit/file#file_path
func TestAuditFileSink_Reopen(t *testing.T) {
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
			ExpectedErrorMessage: "event.(AuditFileSink).open: unable to create file \"/tmp/vault/test/foo.log\": mkdir /tmp/vault/test: permission denied",
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

			sink, err := NewAuditFileSink(tempPath, AuditFormatJSON, tc.Options...)
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
