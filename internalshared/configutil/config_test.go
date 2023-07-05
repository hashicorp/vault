// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSharedConfig_Sanitized_LogFields ensures that 'log related' shared config
// is sanitized as expected.
func TestSharedConfig_Sanitized_LogFields(t *testing.T) {
	tests := map[string]struct {
		Value                     *SharedConfig
		IsNilConfigExpected       bool
		ExpectedLogFile           string
		ExpectedLogFormat         string
		ExpectedLogLevel          string
		ExpectedLogRotateBytes    int
		ExpectedLogRotateDuration string
		ExpectedLogRotateMaxFiles int
	}{
		"nil": {
			Value:               nil,
			IsNilConfigExpected: true,
		},
		"empty": {
			Value: &SharedConfig{},
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
			ExpectedLogFile:           "vault.log",
			ExpectedLogFormat:         "json",
			ExpectedLogLevel:          "warn",
			ExpectedLogRotateBytes:    1024,
			ExpectedLogRotateDuration: "30m",
			ExpectedLogRotateMaxFiles: -1,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			cfg := tc.Value.Sanitized()
			switch {
			case tc.IsNilConfigExpected:
				require.Nil(t, cfg)
			default:
				require.NotNil(t, cfg)
				require.Equal(t, tc.ExpectedLogFile, cfg["log_file"])
				require.Equal(t, tc.ExpectedLogFormat, cfg["log_format"])
				require.Equal(t, tc.ExpectedLogLevel, cfg["log_level"])
				require.Equal(t, tc.ExpectedLogRotateBytes, cfg["log_rotate_bytes"])
				require.Equal(t, tc.ExpectedLogRotateDuration, cfg["log_rotate_duration"])
				require.Equal(t, tc.ExpectedLogRotateMaxFiles, cfg["log_rotate_max_files"])
			}
		})
	}
}
