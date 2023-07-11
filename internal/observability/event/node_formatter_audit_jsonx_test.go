// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"

	"github.com/hashicorp/eventlogger"

	"github.com/stretchr/testify/require"
)

// fakeEvent will return a new fake event containing audit data.
// The audit event is for a response and should be JSONx format.
func fakeEvent(t *testing.T) *eventlogger.Event {
	t.Helper()

	date := time.Date(2023, time.July, 11, 15, 49, 10, 0o0, time.Local)

	auditEvent, err := newAudit(
		WithID("123"),
		WithSubtype(string(AuditResponse)),
		WithFormat(string(AuditFormatJSONX)),
		WithNow(date),
	)
	require.NoError(t, err)
	require.NotNil(t, auditEvent)
	require.Equal(t, "123", auditEvent.ID)
	require.Equal(t, "v0.1", auditEvent.Version)
	require.Equal(t, AuditFormatJSONX, auditEvent.RequiredFormat)
	require.Equal(t, AuditResponse, auditEvent.Subtype)
	require.Equal(t, date, auditEvent.Timestamp)

	e := &eventlogger.Event{
		Type:      eventlogger.EventType(AuditType),
		CreatedAt: auditEvent.Timestamp,
		Formatted: make(map[string][]byte),
		Payload:   auditEvent,
	}

	return e
}

// TestNewAuditFormatterJSONX ensures we can create new AuditFormatterJSONX structs.
func TestNewAuditFormatterJSONX(t *testing.T) {
	f := NewAuditFormatterJSONX()
	require.NotNil(t, f)
}

// TestAuditFormatterJSONX_Reopen ensures that we do no get an error when calling Reopen.
func TestAuditFormatterJSONX_Reopen(t *testing.T) {
	require.NoError(t, NewAuditFormatterJSONX().Reopen())
}

// TestAuditFormatterJSONX_Type ensures that the node is a 'formatter' type.
func TestAuditFormatterJSONX_Type(t *testing.T) {
	require.Equal(t, eventlogger.NodeTypeFormatter, NewAuditFormatterJSONX().Type())
}

// TestAuditFormatterJSONX_Process attempts to run the Process method to convert
// pre-formatted JSON to XML (JSONx).
func TestAuditFormatterJSONX_Process(t *testing.T) {
	tests := map[string]struct {
		IsErrorExpected      bool
		ExpectedErrorMessage string
		Data                 any
	}{
		"no-formatted-json": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFormatterJSONX).Process: pre-formatted JSON required but not found: invalid parameter",
			Data:                 nil,
		},
		"basic-json": {
			IsErrorExpected: false,
			Data:            &logical.LogInput{Type: "magic"},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			e := fakeEvent(t)
			require.NotNil(t, e)

			// If we have data specified, then encode it and store as a format.
			if tc.Data != nil {
				jsonBytes, err := jsonutil.EncodeJSON(tc.Data)
				require.NoError(t, err)
				require.NotNil(t, jsonBytes)
				e.FormattedAs(string(AuditFormatJSON), jsonBytes)
			}

			processed, err := NewAuditFormatterJSONX().Process(context.Background(), e)
			b, found := e.Format(string(AuditFormatJSONX))

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, processed)
				require.False(t, found)
				require.Nil(t, b)
			default:
				require.NoError(t, err)
				require.NotNil(t, processed)
				require.True(t, found)
				require.NotNil(t, b)
			}
		})
	}
}
