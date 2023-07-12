// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"os"
	"testing"

	"github.com/hashicorp/eventlogger"

	"github.com/stretchr/testify/require"
)

// TestAuditFileSinkConfig_Validate ensures that validation of an AuditFileSinkConfig
// returns the expected results given various values for its fields.
func TestAuditFileSinkConfig_Validate(t *testing.T) {
	tests := map[string]struct {
		IsNil                bool
		Path                 string
		Format               auditFormat
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"default": {
			Path:                 "",
			Format:               "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFileSinkConfig).validate: path cannot be empty: invalid parameter",
		},
		"spacey-path": {
			Path:                 "   ",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFileSinkConfig).validate: path cannot be empty: invalid parameter",
		},
		"no-format": {
			Path:                 "/var/nice-path",
			Format:               "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFileSinkConfig).validate: invalid format: event.(audit).(format).validate: '' is not a valid required format: invalid parameter",
		},
		"bad-format": {
			Path:                 "/var/nice-path",
			Format:               "qwerty",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFileSinkConfig).validate: invalid format: event.(audit).(format).validate: 'qwerty' is not a valid required format: invalid parameter",
		},
		"happy": {
			Path:            "/var/nice-path",
			Format:          AuditFormatJSON,
			IsErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg := &AuditFileSinkConfig{Path: tc.Path, Format: tc.Format}
			err := cfg.validate()
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
			}
		})
	}
}

// TestAuditFileSink_Type ensures that the node is a 'sink' type.
func TestAuditFileSink_Type(t *testing.T) {
	f, err := NewAuditFileSink(AuditFileSinkConfig{
		Path:   t.TempDir(),
		Format: AuditFormatJSON,
	})
	require.NoError(t, err)
	require.NotNil(t, f)
	require.Equal(t, eventlogger.NodeTypeSink, f.Type())
}

// TestNewAuditFileSink tests creation of an AuditFileSink.
func TestNewAuditFileSink(t *testing.T) {
	tests := map[string]struct {
		Config               AuditFileSinkConfig
		IsTempDirPath        bool
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedFileMode     os.FileMode
	}{
		"default": {
			Config:               AuditFileSinkConfig{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.NewAuditFileSink: unable to create new audit file sink: event.(AuditFileSinkConfig).validate: path cannot be empty: invalid parameter",
		},
		"default-path-not-exist-valid-format": {
			Config:               AuditFileSinkConfig{Path: "qwerty", Format: AuditFormatJSON},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.NewAuditFileSink: unable to obtain file info: stat qwerty: no such file or directory",
		},
		"default-valid-path": {
			Config:               AuditFileSinkConfig{},
			IsTempDirPath:        true,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.NewAuditFileSink: unable to create new audit file sink: event.(AuditFileSinkConfig).validate: invalid format: event.(audit).(format).validate: '' is not a valid required format: invalid parameter",
		},
		"default-valid-path-and-format": {
			Config: AuditFileSinkConfig{
				Format: AuditFormatJSON,
			},
			IsTempDirPath:    true,
			IsErrorExpected:  false,
			ExpectedFileMode: 0o20000000755,
		},
		"file-mode-not-default-or-zero": {
			Config:           AuditFileSinkConfig{Format: AuditFormatJSON, FileMode: 0o007},
			IsTempDirPath:    true,
			IsErrorExpected:  false,
			ExpectedFileMode: 0o007,
		},
		"path-stdout": {
			Config: AuditFileSinkConfig{
				Path:     "stdout",
				Prefix:   "",
				FileMode: 0x777,
				Format:   AuditFormatJSON,
			},
			IsTempDirPath:    false,
			IsErrorExpected:  false,
			ExpectedFileMode: 0x777,
		},
		"path-stderr": {
			Config: AuditFileSinkConfig{
				Path:     "stderr",
				Prefix:   "",
				FileMode: 0x777,
				Format:   AuditFormatJSON,
			},
			IsTempDirPath:    false,
			IsErrorExpected:  false,
			ExpectedFileMode: 0x777,
		},
		"path-discard": {
			Config: AuditFileSinkConfig{
				Path:     "discard",
				Prefix:   "",
				FileMode: 0x777,
				Format:   AuditFormatJSON,
			},
			IsTempDirPath:    false,
			IsErrorExpected:  false,
			ExpectedFileMode: 0x777,
		},
		"prefix": {
			Config: AuditFileSinkConfig{
				Format: AuditFormatJSON,
				Prefix: "bleep",
			},
			IsTempDirPath:    true,
			IsErrorExpected:  false,
			ExpectedFileMode: 0o20000000755,
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
			if tc.IsTempDirPath {
				tempDir = t.TempDir()
				tc.Config.Path = tempDir
			}

			sink, err := NewAuditFileSink(tc.Config)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, sink)
			default:
				require.NoError(t, err)
				require.NotNil(t, sink)

				require.Equal(t, tc.ExpectedFileMode, sink.fileMode)
				require.Equal(t, tc.Config.Prefix, sink.prefix)
				require.Equal(t, tc.Config.Format, sink.format)

				switch {
				case tc.IsTempDirPath:
					require.Equal(t, tempDir, sink.path)
				default:
					require.Equal(t, tc.Config.Path, sink.path)
				}
			}
		})
	}
}
