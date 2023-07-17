// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestAuditEvent_New exercises the newEvent func to create audit events.
func TestAuditEvent_New(t *testing.T) {
	tests := map[string]struct {
		Options              []Option
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
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.newEvent: audit.(auditEvent).validate: audit.(subtype).validate: '' is not a valid event subtype: invalid parameter",
		},
		"empty-Option": {
			Options:              []Option{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.newEvent: audit.(auditEvent).validate: audit.(subtype).validate: '' is not a valid event subtype: invalid parameter",
		},
		"bad-id": {
			Options:              []Option{WithID("")},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.newEvent: error applying options: id cannot be empty",
		},
		"good": {
			Options: []Option{
				WithID("audit_123"),
				WithFormat(string(JSONFormat)),
				WithSubtype(string(ResponseType)),
				WithNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
			},
			IsErrorExpected:   false,
			ExpectedID:        "audit_123",
			ExpectedTimestamp: time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
			ExpectedSubtype:   ResponseType,
			ExpectedFormat:    JSONFormat,
		},
		"good-no-time": {
			Options: []Option{
				WithID("audit_123"),
				WithFormat(string(JSONFormat)),
				WithSubtype(string(ResponseType)),
			},
			IsErrorExpected: false,
			ExpectedID:      "audit_123",
			ExpectedSubtype: ResponseType,
			ExpectedFormat:  JSONFormat,
			IsNowExpected:   true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			audit, err := newEvent(tc.Options...)
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
				require.Equal(t, tc.ExpectedFormat, audit.RequiredFormat)
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
	tests := map[string]struct {
		Value                *auditEvent
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"nil": {
			Value:                nil,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditEvent).validate: event is nil: invalid parameter",
		},
		"default": {
			Value:                &auditEvent{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditEvent).validate: missing ID: invalid parameter",
		},
		"id-empty": {
			Value: &auditEvent{
				ID:             "",
				Version:        version,
				Subtype:        RequestType,
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: JSONFormat,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditEvent).validate: missing ID: invalid parameter",
		},
		"version-fiddled": {
			Value: &auditEvent{
				ID:             "audit_123",
				Version:        "magic-v2",
				Subtype:        RequestType,
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: JSONFormat,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditEvent).validate: event version unsupported: invalid parameter",
		},
		"subtype-fiddled": {
			Value: &auditEvent{
				ID:             "audit_123",
				Version:        version,
				Subtype:        subtype("moon"),
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: JSONFormat,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditEvent).validate: audit.(subtype).validate: 'moon' is not a valid event subtype: invalid parameter",
		},
		"format-fiddled": {
			Value: &auditEvent{
				ID:             "audit_123",
				Version:        version,
				Subtype:        ResponseType,
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: format("blah"),
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditEvent).validate: audit.(format).validate: 'blah' is not a valid format: invalid parameter",
		},
		"default-time": {
			Value: &auditEvent{
				ID:             "audit_123",
				Version:        version,
				Subtype:        ResponseType,
				Timestamp:      time.Time{},
				Data:           nil,
				RequiredFormat: JSONFormat,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditEvent).validate: event timestamp cannot be the zero time instant: invalid parameter",
		},
		"valid": {
			Value: &auditEvent{
				ID:             "audit_123",
				Version:        version,
				Subtype:        ResponseType,
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: JSONFormat,
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
	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"empty": {
			Value:                "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(subtype).validate: '' is not a valid event subtype: invalid parameter",
		},
		"unsupported": {
			Value:                "foo",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(subtype).validate: 'foo' is not a valid event subtype: invalid parameter",
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
	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"empty": {
			Value:                "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(format).validate: '' is not a valid format: invalid parameter",
		},
		"unsupported": {
			Value:                "foo",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(format).validate: 'foo' is not a valid format: invalid parameter",
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
