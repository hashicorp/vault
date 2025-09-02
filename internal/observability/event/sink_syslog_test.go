// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewSyslogSink ensures that we validate the input arguments and can create
// the SyslogSink if everything goes to plan.
func TestNewSyslogSink(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		format         string
		opts           []Option
		want           *SyslogSink
		wantErr        bool
		expectedErrMsg string
	}{
		"format-empty": {
			format:         "",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"format-whitespace": {
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"happy": {
			format: "json",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := NewSyslogSink(tc.format, tc.opts...)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}
