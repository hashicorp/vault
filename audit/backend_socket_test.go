// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestSocketBackend_newSocketBackend ensures that we can correctly configure the sink
// node on the Backend, and any incorrect parameters result in the relevant errors.
func TestSocketBackend_newSocketBackend(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		mountPath      string
		address        string
		socketType     string
		writeDuration  string
		format         string
		wantErr        bool
		expectedErrMsg string
		expectedName   string
	}{
		"name-empty": {
			mountPath:      "",
			address:        "wss://foo",
			format:         "json",
			wantErr:        true,
			expectedErrMsg: "mount path cannot be empty: invalid configuration",
		},
		"name-whitespace": {
			mountPath:      "   ",
			address:        "wss://foo",
			format:         "json",
			wantErr:        true,
			expectedErrMsg: "mount path cannot be empty: invalid configuration",
		},
		"address-empty": {
			mountPath:      "foo",
			address:        "",
			format:         "json",
			wantErr:        true,
			expectedErrMsg: "\"address\" cannot be empty: invalid configuration",
		},
		"address-whitespace": {
			mountPath:      "foo",
			address:        "   ",
			format:         "json",
			wantErr:        true,
			expectedErrMsg: "\"address\" cannot be empty: invalid configuration",
		},
		"format-empty": {
			mountPath:      "foo",
			address:        "wss://foo",
			format:         "",
			wantErr:        true,
			expectedErrMsg: "unsupported \"format\": invalid configuration",
		},
		"format-whitespace": {
			mountPath:      "foo",
			address:        "wss://foo",
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "unsupported \"format\": invalid configuration",
		},
		"write-duration-valid": {
			mountPath:     "foo",
			address:       "wss://foo",
			writeDuration: "5s",
			format:        "json",
			wantErr:       false,
			expectedName:  "foo",
		},
		"write-duration-not-valid": {
			mountPath:      "foo",
			address:        "wss://foo",
			writeDuration:  "qwerty",
			format:         "json",
			wantErr:        true,
			expectedErrMsg: "unable to parse max duration: invalid parameter: time: invalid duration \"qwerty\"",
		},
		"happy": {
			mountPath:    "foo",
			address:      "wss://foo",
			format:       "json",
			wantErr:      false,
			expectedName: "foo",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg := &BackendConfig{
				SaltView:   &logical.InmemStorage{},
				SaltConfig: &salt.Config{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"address":       tc.address,
					"format":        tc.format,
					"write_timeout": tc.writeDuration,
					"socket":        tc.socketType,
				},
				MountPath: tc.mountPath,
			}
			b, err := newSocketBackend(cfg, &NoopHeaderFormatter{})

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
				require.Nil(t, b)
			} else {
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 2) // formatter + sink
				require.Len(t, b.nodeMap, 2)
				id := b.nodeIDList[1] // sink is 2nd
				node := b.nodeMap[id]
				require.Equal(t, eventlogger.NodeTypeSink, node.Type())
				mc, ok := node.(*event.MetricsCounter)
				require.True(t, ok)
				require.Equal(t, tc.expectedName, mc.Name)
			}
		})
	}
}
