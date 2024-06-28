// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type mapValue[T any] struct {
	Value   T
	IsFound bool
}

type expectedLogFields struct {
	File           mapValue[string]
	Format         mapValue[string]
	Level          mapValue[string]
	RotateBytes    mapValue[int]
	RotateDuration mapValue[string]
	RotateMaxFiles mapValue[int]
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
				Format: mapValue[string]{IsFound: true, Value: ""},
				Level:  mapValue[string]{IsFound: true, Value: ""},
			},
		},
		"only-log-level-and-format": {
			Value: &SharedConfig{
				LogFormat: "json",
				LogLevel:  "warn",
			},
			IsNil: false,
			Expected: expectedLogFields{
				Format: mapValue[string]{IsFound: true, Value: "json"},
				Level:  mapValue[string]{IsFound: true, Value: "warn"},
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
				File:           mapValue[string]{IsFound: true, Value: "vault.log"},
				Format:         mapValue[string]{IsFound: true, Value: "json"},
				Level:          mapValue[string]{IsFound: true, Value: "warn"},
				RotateBytes:    mapValue[int]{IsFound: true, Value: 1024},
				RotateDuration: mapValue[string]{IsFound: true, Value: "30m"},
				RotateMaxFiles: mapValue[int]{IsFound: true, Value: -1},
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
				val, found := cfg["log_file"]
				switch {
				case tc.Expected.File.IsFound:
					require.True(t, found)
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.File.Value, val)
				default:
					require.Nil(t, val)
				}

				// Log format
				val, found = cfg["log_format"]
				switch {
				case tc.Expected.Format.IsFound:
					require.True(t, found)
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.Format.Value, val)
				default:
					require.Nil(t, val)
				}

				// Log level
				val, found = cfg["log_level"]
				switch {
				case tc.Expected.Level.IsFound:
					require.True(t, found)
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.Level.Value, val)
				default:
					require.Nil(t, val)
				}

				// Log rotate bytes
				val, found = cfg["log_rotate_bytes"]
				switch {
				case tc.Expected.RotateBytes.IsFound:
					require.True(t, found)
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.RotateBytes.Value, val)
				default:
					require.Nil(t, val)
				}

				// Log rotate duration
				val, found = cfg["log_rotate_duration"]
				switch {
				case tc.Expected.RotateDuration.IsFound:
					require.True(t, found)
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.RotateDuration.Value, val)
				default:
					require.Nil(t, val)
				}

				// Log rotate max files
				val, found = cfg["log_rotate_max_files"]
				switch {
				case tc.Expected.RotateMaxFiles.IsFound:
					require.True(t, found)
					require.NotNil(t, val)
					require.Equal(t, tc.Expected.RotateMaxFiles.Value, val)
				default:
					require.Nil(t, val)
				}
			}
		})
	}
}
