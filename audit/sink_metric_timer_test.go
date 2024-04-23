// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/stretchr/testify/require"
)

// TestNewSinkMetricTimer ensures that parameters are checked correctly and errors
// reported as expected when attempting to create a SinkMetricTimer.
func TestNewSinkMetricTimer(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name                 string
		node                 eventlogger.Node
		isErrorExpected      bool
		expectedErrorMessage string
	}{
		"happy": {
			name:            "foo",
			node:            &event.FileSink{},
			isErrorExpected: false,
		},
		"no-name": {
			name:                 "",
			isErrorExpected:      true,
			expectedErrorMessage: "audit.NewSinkMetricTimer: name is required: invalid parameter",
		},
		"no-node": {
			name:                 "foo",
			node:                 nil,
			isErrorExpected:      true,
			expectedErrorMessage: "audit.NewSinkMetricTimer: sink node is required: invalid parameter",
		},
		"bad-node": {
			name:                 "foo",
			node:                 &EntryFormatter{},
			isErrorExpected:      true,
			expectedErrorMessage: "audit.NewSinkMetricTimer: sink node must be of type 'sink': invalid parameter",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m, err := NewSinkMetricTimer(tc.name, tc.node)

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrorMessage)
				require.Nil(t, m)
			default:
				require.NoError(t, err)
				require.NotNil(t, m)
			}
		})
	}
}
