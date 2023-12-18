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
	tests := map[string]struct {
		Options              []Option
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
			ExpectedErrorMessage: "audit.newEvent: audit.(AuditEvent).validate: audit.(subtype).validate: '' is not a valid event subtype: invalid parameter",
		},
		"empty-Option": {
			Options:              []Option{},
			Subtype:              subtype(""),
			Format:               format(""),
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.newEvent: audit.(AuditEvent).validate: audit.(subtype).validate: '' is not a valid event subtype: invalid parameter",
		},
		"bad-id": {
			Options:              []Option{WithID("")},
			Subtype:              ResponseType,
			Format:               JSONFormat,
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
			Subtype:           RequestType,
			Format:            JSONxFormat,
			IsErrorExpected:   false,
			ExpectedID:        "audit_123",
			ExpectedTimestamp: time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
			ExpectedSubtype:   RequestType,
			ExpectedFormat:    JSONxFormat,
		},
		"good-no-time": {
			Options: []Option{
				WithID("audit_123"),
				WithFormat(string(JSONFormat)),
				WithSubtype(string(ResponseType)),
			},
			Subtype:         RequestType,
			Format:          JSONxFormat,
			IsErrorExpected: false,
			ExpectedID:      "audit_123",
			ExpectedSubtype: RequestType,
			ExpectedFormat:  JSONxFormat,
			IsNowExpected:   true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			audit, err := NewEvent(tc.Subtype, tc.Options...)
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
	tests := map[string]struct {
		Value                *AuditEvent
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"nil": {
			Value:                nil,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditEvent).validate: event is nil: invalid parameter",
		},
		"default": {
			Value:                &AuditEvent{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditEvent).validate: missing ID: invalid parameter",
		},
		"id-empty": {
			Value: &AuditEvent{
				ID:        "",
				Version:   version,
				Subtype:   RequestType,
				Timestamp: time.Now(),
				Data:      nil,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditEvent).validate: missing ID: invalid parameter",
		},
		"version-fiddled": {
			Value: &AuditEvent{
				ID:        "audit_123",
				Version:   "magic-v2",
				Subtype:   RequestType,
				Timestamp: time.Now(),
				Data:      nil,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditEvent).validate: event version unsupported: invalid parameter",
		},
		"subtype-fiddled": {
			Value: &AuditEvent{
				ID:        "audit_123",
				Version:   version,
				Subtype:   subtype("moon"),
				Timestamp: time.Now(),
				Data:      nil,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditEvent).validate: audit.(subtype).validate: 'moon' is not a valid event subtype: invalid parameter",
		},
		"default-time": {
			Value: &AuditEvent{
				ID:        "audit_123",
				Version:   version,
				Subtype:   ResponseType,
				Timestamp: time.Time{},
				Data:      nil,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditEvent).validate: event timestamp cannot be the zero time instant: invalid parameter",
		},
		"valid": {
			Value: &AuditEvent{
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
