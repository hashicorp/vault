// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestOptions_WithNow exercises WithNow option to ensure it performs as expected.
func TestOptions_WithNow(t *testing.T) {
	tests := map[string]struct {
		Value                time.Time
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedValue        time.Time
	}{
		"default-time": {
			Value:                time.Time{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot specify 'now' to be the zero time instant",
		},
		"valid-time": {
			Value:           time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
			IsErrorExpected: false,
			ExpectedValue:   time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			options := &Options{}
			applyOption := WithNow(tc.Value)
			err := applyOption(options)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, options.WithNow)
			}
		})
	}
}

// TestOptions_WithID exercises WithID option to ensure it performs as expected.
func TestOptions_WithID(t *testing.T) {
	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedValue        string
	}{
		"empty": {
			Value:                "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "id cannot be empty",
		},
		"whitespace": {
			Value:                "     ",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "id cannot be empty",
		},
		"valid": {
			Value:           "test",
			IsErrorExpected: false,
			ExpectedValue:   "test",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := &Options{}
			applyOption := WithID(tc.Value)
			err := applyOption(options)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, options.WithID)
			}
		})
	}
}

// TestOptions_Default exercises getDefaultOptions to assert the default values.
func TestOptions_Default(t *testing.T) {
	opts := getDefaultOptions()
	require.NotNil(t, opts)
	require.True(t, time.Now().After(opts.WithNow))
	require.False(t, opts.WithNow.IsZero())
}

// TestOptions_Opts exercises GetOpts with various Option values.
func TestOptions_Opts(t *testing.T) {
	tests := map[string]struct {
		opts                 []Option
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedID           string
		IsNowExpected        bool
		ExpectedNow          time.Time
	}{
		"nil-options": {
			opts:            nil,
			IsErrorExpected: false,
			IsNowExpected:   true,
		},
		"empty-options": {
			opts:            []Option{},
			IsErrorExpected: false,
			IsNowExpected:   true,
		},
		"with-multiple-valid-id": {
			opts: []Option{
				WithID("qwerty"),
				WithID("juan"),
			},
			IsErrorExpected: false,
			ExpectedID:      "juan",
			IsNowExpected:   true,
		},
		//"with-multiple-valid-subtype": {
		//	opts: []Option{
		//		WithSubtype("AuditRequest"),
		//		WithSubtype("AuditResponse"),
		//	},
		//	IsErrorExpected: false,
		//	ExpectedSubtype: AuditResponse,
		//	IsNowExpected:   true,
		//},
		//"with-multiple-valid-format": {
		//	opts: []Option{
		//		WithFormat("json"),
		//		WithFormat("jsonx"),
		//	},
		//	IsErrorExpected: false,
		//	ExpectedFormat:  AuditFormatJSONx,
		//	IsNowExpected:   true,
		//},
		"with-multiple-valid-now": {
			opts: []Option{
				WithNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
				WithNow(time.Date(2023, time.July, 4, 13, 3, 0, 0, time.Local)),
			},
			IsErrorExpected: false,
			ExpectedNow:     time.Date(2023, time.July, 4, 13, 3, 0, 0, time.Local),
			IsNowExpected:   false,
		},
		"with-multiple-valid-then-invalid-now": {
			opts: []Option{
				WithNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
				WithNow(time.Time{}),
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot specify 'now' to be the zero time instant",
		},
		"with-multiple-valid-options": {
			opts: []Option{
				WithID("qwerty"),
				/*	WithSubtype("AuditRequest"),
					WithFormat("json"),*/
				WithNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
			},
			IsErrorExpected: false,
			ExpectedID:      "qwerty",
			// ExpectedSubtype: AuditRequest,
			// ExpectedFormat:  AuditFormatJSON,
			ExpectedNow: time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts, err := GetOpts(tc.opts...)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NotNil(t, opts)
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedID, opts.WithID)
				// require.Equal(t, tc.ExpectedSubtype, opts.withSubtype)
				// require.Equal(t, tc.ExpectedFormat, opts.withFormat)
				switch {
				case tc.IsNowExpected:
					require.True(t, time.Now().After(opts.WithNow))
					require.False(t, opts.WithNow.IsZero())
				default:
					require.Equal(t, tc.ExpectedNow, opts.WithNow)
				}

			}
		})
	}
}

// TODO: Reproduce some of these tests in audit/options.go
