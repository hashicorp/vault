// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestAuditEvent_new exercises the newEvent func to create audit events.
func TestAuditEvent_new(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Options              []option
		Subtype              subtype
		Format               format
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedID           string
		ExpectedFormat       format
		ExpectedSubtype      subtype
		ExpectedTimestamp    time.Time
		IsNowExpected        bool
	}{
		"nil": {
			Options:              nil,
			Subtype:              subtype(""),
			Format:               format(""),
			IsErrorExpected:      true,
			ExpectedErrorMessage: "invalid event subtype \"\": invalid internal parameter",
		},
		"empty-option": {
			Options:              []option{},
			Subtype:              subtype(""),
			Format:               format(""),
			IsErrorExpected:      true,
			ExpectedErrorMessage: "invalid event subtype \"\": invalid internal parameter",
		},
		"bad-id": {
			Options:              []option{withID("")},
			Subtype:              ResponseType,
			Format:               jsonFormat,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "id cannot be empty",
		},
		"good": {
			Options: []option{
				withID("audit_123"),
				withFormat(string(jsonFormat)),
				withSubtype(string(ResponseType)),
				withNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
			},
			Subtype:           RequestType,
			Format:            jsonxFormat,
			IsErrorExpected:   false,
			ExpectedID:        "audit_123",
			ExpectedTimestamp: time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
			ExpectedSubtype:   RequestType,
			ExpectedFormat:    jsonxFormat,
		},
		"good-no-time": {
			Options: []option{
				withID("audit_123"),
				withFormat(string(jsonFormat)),
				withSubtype(string(ResponseType)),
			},
			Subtype:         RequestType,
			Format:          jsonxFormat,
			IsErrorExpected: false,
			ExpectedID:      "audit_123",
			ExpectedSubtype: RequestType,
			ExpectedFormat:  jsonxFormat,
			IsNowExpected:   true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			audit, err := newEvent(tc.Subtype, tc.Options...)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, audit)
			default:
				require.NoError(t, err)
				require.NotNil(t, audit)
				require.Equal(t, tc.ExpectedID, audit.ID)
				require.Equal(t, tc.ExpectedSubtype, audit.Subtype)
				switch {
				case tc.IsNowExpected:
					require.True(t, time.Now().After(audit.Timestamp))
					require.False(t, audit.Timestamp.IsZero())
				default:
					require.Equal(t, tc.ExpectedTimestamp, audit.Timestamp)
				}
			}
		})
	}
}

// TestAuditEvent_Validate exercises the validation for an audit event.
func TestAuditEvent_Validate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value                *Event
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"nil": {
			Value:                nil,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event is nil: invalid internal parameter",
		},
		"default": {
			Value:                &Event{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "missing ID: invalid internal parameter",
		},
		"id-empty": {
			Value: &Event{
				ID:        "",
				Version:   version,
				Subtype:   RequestType,
				Timestamp: time.Now(),
				Data:      nil,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "missing ID: invalid internal parameter",
		},
		"version-fiddled": {
			Value: &Event{
				ID:        "audit_123",
				Version:   "magic-v2",
				Subtype:   RequestType,
				Timestamp: time.Now(),
				Data:      nil,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event version unsupported: invalid internal parameter",
		},
		"subtype-fiddled": {
			Value: &Event{
				ID:        "audit_123",
				Version:   version,
				Subtype:   subtype("moon"),
				Timestamp: time.Now(),
				Data:      nil,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "invalid event subtype \"moon\": invalid internal parameter",
		},
		"default-time": {
			Value: &Event{
				ID:        "audit_123",
				Version:   version,
				Subtype:   ResponseType,
				Timestamp: time.Time{},
				Data:      nil,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event timestamp cannot be the zero time instant: invalid internal parameter",
		},
		"valid": {
			Value: &Event{
				ID:        "audit_123",
				Version:   version,
				Subtype:   ResponseType,
				Timestamp: time.Now(),
				Data:      nil,
			},
			IsErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.Value.validate()
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
			}
		})
	}
}

// TestAuditEvent_Validate_Subtype exercises the validation for an audit event's subtype.
func TestAuditEvent_Validate_Subtype(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"empty": {
			Value:                "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "invalid event subtype \"\": invalid internal parameter",
		},
		"unsupported": {
			Value:                "foo",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "invalid event subtype \"foo\": invalid internal parameter",
		},
		"request": {
			Value:           "AuditRequest",
			IsErrorExpected: false,
		},
		"response": {
			Value:           "AuditResponse",
			IsErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := subtype(tc.Value).validate()
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
			}
		})
	}
}

// TestAuditEvent_Validate_Format exercises the validation for an audit event's format.
func TestAuditEvent_Validate_Format(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"empty": {
			Value:                "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "invalid format \"\": invalid internal parameter",
		},
		"unsupported": {
			Value:                "foo",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "invalid format \"foo\": invalid internal parameter",
		},
		"json": {
			Value:           "json",
			IsErrorExpected: false,
		},
		"jsonx": {
			Value:           "jsonx",
			IsErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := format(tc.Value).validate()
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
			}
		})
	}
}

// TestAuditEvent_Subtype_MetricTag is used to ensure that we get the string value
// we expect for a subtype when we want to use it as a metrics tag.
// In some strange scenario where the subtype was never validated, it is technically
// possible to get a value that isn't related to request/response, but this shouldn't
// really be happening, so we will return it as is.
func TestAuditEvent_Subtype_MetricTag(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input          string
		expectedOutput string
	}{
		"request": {
			input:          "AuditRequest",
			expectedOutput: "log_request",
		},
		"response": {
			input:          "AuditResponse",
			expectedOutput: "log_response",
		},
		"non-validated": {
			input:          "juan",
			expectedOutput: "juan",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			st := subtype(tc.input)
			tag := st.MetricTag()
			require.Equal(t, tc.expectedOutput, tag)
		})
	}
}

// TestAuditEvent_Subtype_String is used to ensure that we get the string value
// we expect for a subtype when it is used with the Stringer interface.
// e.g. an AuditRequest subtype is 'request'
func TestAuditEvent_Subtype_String(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input          string
		expectedOutput string
	}{
		"request": {
			input:          "AuditRequest",
			expectedOutput: "request",
		},
		"response": {
			input:          "AuditResponse",
			expectedOutput: "response",
		},
		"non-validated": {
			input:          "juan",
			expectedOutput: "juan",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			st := subtype(tc.input)
			require.Equal(t, tc.expectedOutput, st.String())
		})
	}
}

// TestAuditEvent_formattedTime is used to check the output from the formattedTime
// method returns the correct format.
func TestAuditEvent_formattedTime(t *testing.T) {
	theTime := time.Date(2024, time.March, 22, 10, 0o0, 5, 10, time.UTC)
	a, err := newEvent(ResponseType, withNow(theTime))
	require.NoError(t, err)
	require.NotNil(t, a)
	require.Equal(t, "2024-03-22T10:00:05.00000001Z", a.formattedTime())
}

// TestEvent_IsValidFormat ensures that we can correctly determine valid and
// invalid formats.
func TestEvent_IsValidFormat(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    string
		expected bool
	}{
		"empty": {
			input:    "",
			expected: false,
		},
		"whitespace": {
			input:    "     ",
			expected: false,
		},
		"invalid-test": {
			input:    "test",
			expected: false,
		},
		"valid-json": {
			input:    "json",
			expected: true,
		},
		"upper-json": {
			input:    "JSON",
			expected: true,
		},
		"mixed-json": {
			input:    "Json",
			expected: true,
		},
		"spacey-json": {
			input:    "  json  ",
			expected: true,
		},
		"valid-jsonx": {
			input:    "jsonx",
			expected: true,
		},
		"upper-jsonx": {
			input:    "JSONX",
			expected: true,
		},
		"mixed-jsonx": {
			input:    "JsonX",
			expected: true,
		},
		"spacey-jsonx": {
			input:    "  jsonx  ",
			expected: true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := isValidFormat(tc.input)
			require.Equal(t, tc.expected, res)
		})
	}
}
