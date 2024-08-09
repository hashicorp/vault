// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestOptions_withFormat exercises withFormat option to ensure it performs as expected.
func TestOptions_withFormat(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedValue        format
	}{
		"empty": {
			Value:           "",
			IsErrorExpected: false,
			ExpectedValue:   format(""),
		},
		"whitespace": {
			Value:           "     ",
			IsErrorExpected: false,
			ExpectedValue:   format(""),
		},
		"invalid-test": {
			Value:                "test",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "invalid format \"test\": invalid internal parameter",
		},
		"valid-json": {
			Value:           "json",
			IsErrorExpected: false,
			ExpectedValue:   jsonFormat,
		},
		"valid-jsonx": {
			Value:           "jsonx",
			IsErrorExpected: false,
			ExpectedValue:   jsonxFormat,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &options{}
			applyOption := withFormat(tc.Value)
			err := applyOption(opts)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, opts.withFormat)
			}
		})
	}
}

// TestOptions_withSubtype exercises withSubtype option to ensure it performs as expected.
func TestOptions_withSubtype(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedValue        subtype
	}{
		"empty": {
			Value:                "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "subtype cannot be empty",
		},
		"whitespace": {
			Value:                "     ",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "subtype cannot be empty",
		},
		"valid": {
			Value:           "AuditResponse",
			IsErrorExpected: false,
			ExpectedValue:   ResponseType,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts := &options{}
			applyOption := withSubtype(tc.Value)
			err := applyOption(opts)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, opts.withSubtype)
			}
		})
	}
}

// TestOptions_withNow exercises withNow option to ensure it performs as expected.
func TestOptions_withNow(t *testing.T) {
	t.Parallel()

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

			opts := &options{}
			applyOption := withNow(tc.Value)
			err := applyOption(opts)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, opts.withNow)
			}
		})
	}
}

// TestOptions_withID exercises withID option to ensure it performs as expected.
func TestOptions_withID(t *testing.T) {
	t.Parallel()

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
			opts := &options{}
			applyOption := withID(tc.Value)
			err := applyOption(opts)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, opts.withID)
			}
		})
	}
}

// TestOptions_withPrefix exercises withPrefix option to ensure it performs as expected.
func TestOptions_withPrefix(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedValue        string
	}{
		"empty": {
			Value:           "",
			IsErrorExpected: false,
			ExpectedValue:   "",
		},
		"whitespace": {
			Value:           "     ",
			IsErrorExpected: false,
			ExpectedValue:   "     ",
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
			opts := &options{}
			applyOption := withPrefix(tc.Value)
			err := applyOption(opts)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, opts.withPrefix)
			}
		})
	}
}

// TestOptions_withRaw exercises withRaw option to ensure it performs as expected.
func TestOptions_withRaw(t *testing.T) {
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
			opts := &options{}
			applyOption := withRaw(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withRaw)
		})
	}
}

// TestOptions_withElision exercises withElision option to ensure it performs as expected.
func TestOptions_withElision(t *testing.T) {
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
			opts := &options{}
			applyOption := withElision(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withElision)
		})
	}
}

// TestOptions_withHMACAccessor exercises withHMACAccessor option to ensure it performs as expected.
func TestOptions_withHMACAccessor(t *testing.T) {
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
			opts := &options{}
			applyOption := withHMACAccessor(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withHMACAccessor)
		})
	}
}

// TestOptions_withOmitTime exercises withOmitTime option to ensure it performs as expected.
func TestOptions_withOmitTime(t *testing.T) {
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
			opts := &options{}
			applyOption := withOmitTime(tc.Value)
			err := applyOption(opts)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, opts.withOmitTime)
		})
	}
}

// TestOptions_Default exercises getDefaultOptions to assert the default values.
func TestOptions_Default(t *testing.T) {
	t.Parallel()

	opts := getDefaultOptions()
	require.NotNil(t, opts)
	require.True(t, time.Now().After(opts.withNow))
	require.False(t, opts.withNow.IsZero())
}

// TestOptions_Opts exercises GetOpts with various option values.
func TestOptions_Opts(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts                 []option
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedID           string
		ExpectedSubtype      subtype
		ExpectedFormat       format
		IsNowExpected        bool
		ExpectedNow          time.Time
	}{
		"nil-options": {
			opts:            nil,
			IsErrorExpected: false,
			IsNowExpected:   true,
			ExpectedFormat:  jsonFormat,
		},
		"empty-options": {
			opts:            []option{},
			IsErrorExpected: false,
			IsNowExpected:   true,
			ExpectedFormat:  jsonFormat,
		},
		"with-multiple-valid-id": {
			opts: []option{
				withID("qwerty"),
				withID("juan"),
			},
			IsErrorExpected: false,
			ExpectedID:      "juan",
			IsNowExpected:   true,
			ExpectedFormat:  jsonFormat,
		},
		"with-multiple-valid-subtype": {
			opts: []option{
				withSubtype("AuditRequest"),
				withSubtype("AuditResponse"),
			},
			IsErrorExpected: false,
			ExpectedSubtype: ResponseType,
			IsNowExpected:   true,
			ExpectedFormat:  jsonFormat,
		},
		"with-multiple-valid-format": {
			opts: []option{
				withFormat("json"),
				withFormat("jsonx"),
			},
			IsErrorExpected: false,
			ExpectedFormat:  jsonxFormat,
			IsNowExpected:   true,
		},
		"with-multiple-valid-now": {
			opts: []option{
				withNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
				withNow(time.Date(2023, time.July, 4, 13, 3, 0, 0, time.Local)),
			},
			IsErrorExpected: false,
			ExpectedNow:     time.Date(2023, time.July, 4, 13, 3, 0, 0, time.Local),
			IsNowExpected:   false,
			ExpectedFormat:  jsonFormat,
		},
		"with-multiple-valid-then-invalid-now": {
			opts: []option{
				withNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
				withNow(time.Time{}),
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot specify 'now' to be the zero time instant",
			ExpectedFormat:       jsonFormat,
		},
		"with-multiple-valid-options": {
			opts: []option{
				withID("qwerty"),
				withSubtype("AuditRequest"),
				withFormat("json"),
				withNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
			},
			IsErrorExpected: false,
			ExpectedID:      "qwerty",
			ExpectedSubtype: RequestType,
			ExpectedFormat:  jsonFormat,
			ExpectedNow:     time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts, err := getOpts(tc.opts...)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NotNil(t, opts)
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedID, opts.withID)
				require.Equal(t, tc.ExpectedSubtype, opts.withSubtype)
				require.Equal(t, tc.ExpectedFormat, opts.withFormat)
				switch {
				case tc.IsNowExpected:
					require.True(t, time.Now().After(opts.withNow))
					require.False(t, opts.withNow.IsZero())
				default:
					require.Equal(t, tc.ExpectedNow, opts.withNow)
				}

			}
		})
	}
}
