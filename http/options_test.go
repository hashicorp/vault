// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestOptions_Default ensures that the default values are as expected.
func TestOptions_Default(t *testing.T) {
	opts := getDefaultOptions()
	require.NotNil(t, opts)
	require.Equal(t, "", opts.withRedactionValue)
}

// TestOptions_WithRedactionValue ensures that we set the correct value to use for
// redaction when required.
func TestOptions_WithRedactionValue(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value           string
		ExpectedValue   string
		IsErrorExpected bool
	}{
		"empty": {
			Value:           "",
			ExpectedValue:   "",
			IsErrorExpected: false,
		},
		"whitespace": {
			Value:           "     ",
			ExpectedValue:   "     ",
			IsErrorExpected: false,
		},
		"value": {
			Value:           "*****",
			ExpectedValue:   "*****",
			IsErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &listenerConfigOptions{}
			applyOption := WithRedactionValue(tc.Value)
			err := applyOption(opts)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, opts.withRedactionValue)
			}
		})
	}
}

// TestOptions_WithRedactAddresses ensures that the option works as intended.
func TestOptions_WithRedactAddresses(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value         bool
		ExpectedValue bool
	}{
		"true": {
			Value:         true,
			ExpectedValue: true,
		},
		"false": {
			Value:         false,
			ExpectedValue: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &listenerConfigOptions{}
			applyOption := WithRedactAddresses(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withRedactAddresses)
		})
	}
}

// TestOptions_WithRedactClusterName ensures that the option works as intended.
func TestOptions_WithRedactClusterName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value         bool
		ExpectedValue bool
	}{
		"true": {
			Value:         true,
			ExpectedValue: true,
		},
		"false": {
			Value:         false,
			ExpectedValue: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &listenerConfigOptions{}
			applyOption := WithRedactClusterName(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withRedactClusterName)
		})
	}
}

// TestOptions_WithRedactVersion ensures that the option works as intended.
func TestOptions_WithRedactVersion(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value         bool
		ExpectedValue bool
	}{
		"true": {
			Value:         true,
			ExpectedValue: true,
		},
		"false": {
			Value:         false,
			ExpectedValue: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &listenerConfigOptions{}
			applyOption := WithRedactVersion(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withRedactVersion)
		})
	}
}
