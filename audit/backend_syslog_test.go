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

// TestSyslogBackend_newSyslogBackend tests the ways we can try to create a new
// SyslogBackend both good and bad.
func TestSyslogBackend_newSyslogBackend(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		mountPath      string
		format         string
		tag            string
		facility       string
		wantErr        bool
		expectedErrMsg string
		expectedName   string
	}{
		"name-empty": {
			mountPath:      "",
			wantErr:        true,
			expectedErrMsg: "mount path cannot be empty: invalid configuration",
		},
		"name-whitespace": {
			mountPath:      "   ",
			wantErr:        true,
			expectedErrMsg: "mount path cannot be empty: invalid configuration",
		},
		"format-empty": {
			mountPath:      "foo",
			format:         "",
			wantErr:        true,
			expectedErrMsg: "unsupported \"format\": invalid configuration",
		},
		"format-whitespace": {
			mountPath:      "foo",
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "unsupported \"format\": invalid configuration",
		},
		"happy": {
			mountPath:    "foo",
			format:       "json",
			wantErr:      false,
			expectedName: "foo",
		},
		"happy-tag": {
			mountPath:    "foo",
			format:       "json",
			tag:          "beep",
			wantErr:      false,
			expectedName: "foo",
		},
		"happy-facility": {
			mountPath:    "foo",
			format:       "json",
			facility:     "daemon",
			wantErr:      false,
			expectedName: "foo",
		},
		"happy-all": {
			mountPath:    "foo",
			format:       "json",
			tag:          "beep",
			facility:     "daemon",
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
					"tag":      tc.tag,
					"facility": tc.facility,
					"format":   tc.format,
				},
				MountPath: tc.mountPath,
			}
			b, err := newSyslogBackend(cfg, &NoopHeaderFormatter{})

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
				require.Nil(t, b)
			} else {
				require.NoError(t, err)
				require.Len(t, b.nodeIDList, 2)
				require.Len(t, b.nodeMap, 2)
				id := b.nodeIDList[1]
				node := b.nodeMap[id]
				require.Equal(t, eventlogger.NodeTypeSink, node.Type())
				mc, ok := node.(*event.MetricsCounter)
				require.True(t, ok)
				require.Equal(t, tc.expectedName, mc.Name)
			}
		})
	}
}
