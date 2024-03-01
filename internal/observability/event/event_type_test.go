// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEventType_Validate exercises the Validate method for EventType.
func TestEventType_Validate(t *testing.T) {
	tests := map[string]struct {
		Value         string
		IsValid       bool
		ExpectedError string
	}{
		"audit": {
			Value:   "audit",
			IsValid: true,
		},
		"empty": {
			Value:         "",
			IsValid:       false,
			ExpectedError: "event.(EventType).Validate: '' is not a valid event type: invalid parameter",
		},
		"random": {
			Value:         "random",
			IsValid:       false,
			ExpectedError: "event.(EventType).Validate: 'random' is not a valid event type: invalid parameter",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			eventType := EventType(tc.Value)
			err := eventType.Validate()
			switch {
			case tc.IsValid:
				require.NoError(t, err)
			case !tc.IsValid:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedError)
			}
		})
	}
}
