// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestNewSocketSink ensures that we validate the input arguments and can create
// the SocketSink if everything goes to plan.
func TestNewSocketSink(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		address        string
		format         string
		opts           []Option
		want           *SocketSink
		wantErr        bool
		expectedErrMsg string
	}{
		"address-empty": {
			address:        "",
			wantErr:        true,
			expectedErrMsg: "address is required: invalid parameter",
		},
		"address-whitespace": {
			address:        "    ",
			wantErr:        true,
			expectedErrMsg: "address is required: invalid parameter",
		},
		"format-empty": {
			address:        "addr",
			format:         "",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"format-whitespace": {
			address:        "addr",
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"bad-max-duration": {
			address:        "addr",
			format:         "json",
			opts:           []Option{WithMaxDuration("bar")},
			wantErr:        true,
			expectedErrMsg: "unable to parse max duration: invalid parameter: time: invalid duration \"bar\"",
		},
		"happy": {
			address: "wss://foo",
			format:  "json",
			want: &SocketSink{
				requiredFormat: "json",
				address:        "wss://foo",
				socketType:     "tcp",           // defaults to tcp
				maxDuration:    2 * time.Second, // defaults to 2 secs
			},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := NewSocketSink(tc.address, tc.format, tc.opts...)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}
