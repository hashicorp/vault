// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type expectedValue[T any] struct {
	Value T
	IsNil bool
}

type expectedLogFields struct {
	File           expectedValue[string]
	Format         expectedValue[string]
	Level          expectedValue[string]
	RotateBytes    expectedValue[int]
	RotateDuration expectedValue[string]
	RotateMaxFiles expectedValue[int]
}

// TestSharedConfig_Sanitized_LogFields ensures that 'log related' shared config
// is sanitized as expected.
func TestSharedConfig_Sanitized_LogFields(t *testing.T) {
	tests := map[string]struct {
		Value    *SharedConfig
		IsNil    bool
		Expected expectedLogFields
	}{
		"nil": {
			Value: nil,
			IsNil: true,
		},
		"empty": {
			Value: &SharedConfig{},
			IsNil: false,
			Expected: expectedLogFields{
				File:           expectedValue[string]{IsNil: true},
				Format:         expectedValue[string]{IsNil: false, Value: ""},
				Level:          expectedValue[string]{IsNil: false, Value: ""},
				RotateBytes:    expectedValue[int]{IsNil: true},
				RotateDuration: expectedValue[string]{IsNil: true},
				RotateMaxFiles: expectedValue[int]{IsNil: true},
			},
		},
		"only-log-level-and-format": {
			Value: &SharedConfig{
				LogFormat: "json",
				LogLevel:  "warn",
			},
			IsNil: false,
			Expected: expectedLogFields{
				File:           expectedValue[string]{IsNil: true},
				Format:         expectedValue[string]{Value: "json"},
				Level:          expectedValue[string]{Value: "warn"},
				RotateBytes:    expectedValue[int]{IsNil: true},
				RotateDuration: expectedValue[string]{IsNil: true},
				RotateMaxFiles: expectedValue[int]{IsNil: true},
			},
		},
		"valid-log-fields": {
			Value: &SharedConfig{
				LogFile:           "vault.log",
				LogFormat:         "json",
				LogLevel:          "warn",
				LogRotateBytes:    1024,
				LogRotateDuration: "30m",
				LogRotateMaxFiles: -1,
			},
			IsNil: false,
			Expected: expectedLogFields{
				File:           expectedValue[string]{Value: "vault.log"},
				Format:         expectedValue[string]{Value: "json"},
				Level:          expectedValue[string]{Value: "warn"},
				RotateBytes:    expectedValue[int]{Value: 1024},
				RotateDuration: expectedValue[string]{Value: "30m"},
				RotateMaxFiles: expectedValue[int]{Value: -1},
			},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			cfg := tc.Value.Sanitized()
			switch {
			case tc.IsNil:
				require.Nil(t, cfg)
			default:
				require.NotNil(t, cfg)

				// Log file
				val := cfg["log_file"]
				switch {
				case tc.Expected.File.IsNil:
					require.Nil(t, val)
				default:
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.File.Value, val)
				}

				// Log format
				val = cfg["log_format"]
				switch {
				case tc.Expected.Format.IsNil:
					require.Nil(t, val)
				default:
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.Format.Value, val)
				}

				// Log level
				val = cfg["log_level"]
				switch {
				case tc.Expected.Level.IsNil:
					require.Nil(t, val)
				default:
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.Level.Value, val)
				}

				// Log rotate bytes
				val = cfg["log_rotate_bytes"]
				switch {
				case tc.Expected.RotateBytes.IsNil:
					require.Nil(t, val)
				default:
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.RotateBytes.Value, val)
				}

				// Log rotate duration
				val = cfg["log_rotate_duration"]
				switch {
				case tc.Expected.RotateDuration.IsNil:
					require.Nil(t, val)
				default:
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.RotateDuration.Value, val)
				}

				// Log rotate max files
				val = cfg["log_rotate_max_files"]
				switch {
				case tc.Expected.RotateMaxFiles.IsNil:
					require.Nil(t, val)
				default:
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.RotateMaxFiles.Value, val)
				}
			}
		})
	}
}
