// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/stretchr/testify/require"
)

var (
	_ eventlogger.Node = (*testEventLoggerNode)(nil)
	_ Labeler          = (*testMetricsCounter)(nil)
)

// TestNewMetricsCounter ensures that NewMetricsCounter operates as intended and
// can validate the input parameters correctly, returning the right error message
// when required.
func TestNewMetricsCounter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name                 string
		node                 eventlogger.Node
		labeler              Labeler
		isErrorExpected      bool
		expectedErrorMessage string
	}{
		"happy": {
			name:            "foo",
			node:            &testEventLoggerNode{},
			labeler:         &testMetricsCounter{},
			isErrorExpected: false,
		},
		"no-name": {
			node:                 nil,
			labeler:              nil,
			isErrorExpected:      true,
			expectedErrorMessage: "event.NewMetricsCounter: name is required: invalid parameter",
		},
		"no-node": {
			name:                 "foo",
			node:                 nil,
			isErrorExpected:      true,
			expectedErrorMessage: "event.NewMetricsCounter: node is required: invalid parameter",
		},
		"no-labeler": {
			name:                 "foo",
			node:                 &testEventLoggerNode{},
			labeler:              nil,
			isErrorExpected:      true,
			expectedErrorMessage: "event.NewMetricsCounter: labeler is required: invalid parameter",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			m, err := NewMetricsCounter(tc.name, tc.node, tc.labeler)

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrorMessage)
			default:
				require.NoError(t, err)
				require.NotNil(t, m)
			}
		})
	}
}

// testEventLoggerNode is for testing and implements the eventlogger.Node interface.
type testEventLoggerNode struct{}

func (t testEventLoggerNode) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	return nil, nil
}

func (t testEventLoggerNode) Reopen() error {
	return nil
}

func (t testEventLoggerNode) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeSink
}

// testMetricsCounter is for testing and implements the event.Labeler interface.
type testMetricsCounter struct{}

func (m *testMetricsCounter) Labels(_ *eventlogger.Event, err error) []string {
	return []string{""}
}
