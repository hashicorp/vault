// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestAuditEvent_New exercises the newAudit func to create audit events.
func TestAuditEvent_New(t *testing.T) {
	tests := map[string]struct {
		Options              []AuditOption
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedID           string
		ExpectedFormat       auditFormat
		ExpectedSubtype      auditSubtype
		ExpectedTimestamp    time.Time
		IsNowExpected        bool
	}{
		"nil": {
			Options:              nil,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.newAudit: audit.(audit).validate: audit.(auditSubtype).validate: '' is not a valid event subtype: invalid parameter",
		},
		"empty-option": {
			Options:              []AuditOption{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.newAudit: audit.(audit).validate: audit.(auditSubtype).validate: '' is not a valid event subtype: invalid parameter",
		},
		"bad-id": {
			Options:              []AuditOption{WithID("")},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.newAudit: error applying options: id cannot be empty",
		},
		"good": {
			Options: []AuditOption{
				WithID("audit_123"),
				WithFormat(string(AuditFormatJSON)),
				WithSubtype(string(AuditResponseType)),
				WithNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
			},
			IsErrorExpected:   false,
			ExpectedID:        "audit_123",
			ExpectedTimestamp: time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
			ExpectedSubtype:   AuditResponseType,
			ExpectedFormat:    AuditFormatJSON,
		},
		"good-no-time": {
			Options: []AuditOption{
				WithID("audit_123"),
				WithFormat(string(AuditFormatJSON)),
				WithSubtype(string(AuditResponseType)),
			},
			IsErrorExpected: false,
			ExpectedID:      "audit_123",
			ExpectedSubtype: AuditResponseType,
			ExpectedFormat:  AuditFormatJSON,
			IsNowExpected:   true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			audit, err := newAudit(tc.Options...)
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
		Value                *audit
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"nil": {
			Value:                nil,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(audit).validate: audit is nil: invalid parameter",
		},
		"default": {
			Value:                &audit{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(audit).validate: missing ID: invalid parameter",
		},
		"id-empty": {
			Value: &audit{
				ID:             "",
				Version:        auditVersion,
				Subtype:        AuditRequestType,
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: AuditFormatJSON,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(audit).validate: missing ID: invalid parameter",
		},
		"version-fiddled": {
			Value: &audit{
				ID:             "audit_123",
				Version:        "magic-v2",
				Subtype:        AuditRequestType,
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: AuditFormatJSON,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(audit).validate: audit version unsupported: invalid parameter",
		},
		"subtype-fiddled": {
			Value: &audit{
				ID:             "audit_123",
				Version:        auditVersion,
				Subtype:        auditSubtype("moon"),
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: AuditFormatJSON,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(audit).validate: audit.(auditSubtype).validate: 'moon' is not a valid event subtype: invalid parameter",
		},
		"format-fiddled": {
			Value: &audit{
				ID:             "audit_123",
				Version:        auditVersion,
				Subtype:        AuditResponseType,
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: auditFormat("blah"),
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(audit).validate: audit.(auditFormat).validate: 'blah' is not a valid format: invalid parameter",
		},
		"default-time": {
			Value: &audit{
				ID:             "audit_123",
				Version:        auditVersion,
				Subtype:        AuditResponseType,
				Timestamp:      time.Time{},
				Data:           nil,
				RequiredFormat: AuditFormatJSON,
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(audit).validate: audit timestamp cannot be the zero time instant: invalid parameter",
		},
		"valid": {
			Value: &audit{
				ID:             "audit_123",
				Version:        auditVersion,
				Subtype:        AuditResponseType,
				Timestamp:      time.Now(),
				Data:           nil,
				RequiredFormat: AuditFormatJSON,
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
			ExpectedErrorMessage: "audit.(auditSubtype).validate: '' is not a valid event subtype: invalid parameter",
		},
		"unsupported": {
			Value:                "foo",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditSubtype).validate: 'foo' is not a valid event subtype: invalid parameter",
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

			err := auditSubtype(tc.Value).validate()
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
			ExpectedErrorMessage: "audit.(auditFormat).validate: '' is not a valid format: invalid parameter",
		},
		"unsupported": {
			Value:                "foo",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditFormat).validate: 'foo' is not a valid format: invalid parameter",
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

			err := auditFormat(tc.Value).validate()
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
