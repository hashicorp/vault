// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSinkType_Validate exercises the validation for a sink type.
func TestSinkType_Validate(t *testing.T) {
	tests := map[string]struct {
		Value         string
		IsValid       bool
		ExpectedError string
	}{
		"file": {
			Value:   "file",
			IsValid: true,
		},
		"syslog": {
			Value:   "syslog",
			IsValid: true,
		},
		"socket": {
			Value:   "socket",
			IsValid: true,
		},
		"empty": {
			Value:         "",
			IsValid:       false,
			ExpectedError: "event.(SinkType).Validate: '' is not a valid sink type: invalid parameter",
		},
		"random": {
			Value:         "random",
			IsValid:       false,
			ExpectedError: "event.(SinkType).Validate: 'random' is not a valid sink type: invalid parameter",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			sinkType := SinkType(tc.Value)
			err := sinkType.Validate()
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
