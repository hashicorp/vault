// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/internal/observability/event"

	"github.com/hashicorp/vault/helper/namespace"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/eventlogger"
	"github.com/stretchr/testify/require"
)

// fakeJSONAuditEvent will return a new fake event containing audit data based
// on the specified subtype and logical.LogInput.
func fakeJSONAuditEvent(tb testing.TB, subtype subtype, input *logical.LogInput) *eventlogger.Event {
	tb.Helper()

	date := time.Date(2023, time.July, 11, 15, 49, 10, 0o0, time.Local)

	auditEvent, err := newEvent(subtype, JSONFormat,
		WithID("123"),
		WithNow(date),
	)
	require.NoError(tb, err)
	require.NotNil(tb, auditEvent)
	require.Equal(tb, "123", auditEvent.ID)
	require.Equal(tb, "v0.1", auditEvent.Version)
	require.Equal(tb, JSONFormat, auditEvent.RequiredFormat)
	require.Equal(tb, subtype, auditEvent.Subtype)
	require.Equal(tb, date, auditEvent.Timestamp)

	auditEvent.Data = input

	e := &eventlogger.Event{
		Type:      eventlogger.EventType(event.AuditType),
		CreatedAt: auditEvent.Timestamp,
		Formatted: make(map[string][]byte),
		Payload:   auditEvent,
	}

	return e
}

// TestNewAuditFormatterJSON ensures we can create new AuditFormatterJSONX structs.
func TestNewAuditFormatterJSON(t *testing.T) {
	tests := map[string]struct {
		UseStaticSalt        bool
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"nil-salter": {
			UseStaticSalt:        false,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.NewAuditFormatterJSON: unable to create new JSON audit formatter: cannot create a new audit formatter with nil salter",
		},
		"static-salter": {
			UseStaticSalt:   true,
			IsErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cfg := FormatterConfig{}
			var ss Salter
			if tc.UseStaticSalt {
				ss = newStaticSalt(t)
			}

			f, err := NewAuditFormatterJSON(cfg, ss)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, f)
			default:
				require.NoError(t, err)
				require.NotNil(t, f)
			}
		})
	}
}

// TestAuditFormatterJSONX_Reopen ensures that we do no get an error when calling Reopen.
func TestAuditFormatterJSON_Reopen(t *testing.T) {
	ss := newStaticSalt(t)
	cfg := FormatterConfig{}

	f, err := NewAuditFormatterJSON(cfg, ss)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.NoError(t, f.Reopen())
}

// TestAuditFormatterJSONX_Type ensures that the node is a 'formatter' type.
func TestAuditFormatterJSON_Type(t *testing.T) {
	ss := newStaticSalt(t)
	cfg := FormatterConfig{}

	f, err := NewAuditFormatterJSON(cfg, ss)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.Equal(t, eventlogger.NodeTypeFormatter, f.Type())
}

// TestAuditFormatterJSON_Process attempts to run the Process method to convert
// the logical.LogInput within an audit event to JSON (AuditRequestEntry or AuditResponseEntry).
func TestAuditFormatterJSON_Process(t *testing.T) {
	tests := map[string]struct {
		IsErrorExpected      bool
		ExpectedErrorMessage string
		Subtype              subtype
		Data                 *logical.LogInput
		RootNamespace        bool
	}{
		"request-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditFormatterJSON).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			Data:                 nil,
		},
		"response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditFormatterJSON).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			Data:                 nil,
		},
		"request-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditFormatterJSON).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"response-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditFormatterJSON).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"request-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditFormatterJSON).Process: unable to parse request from audit event: no namespace",
			Subtype:              RequestType,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"response-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(AuditFormatterJSON).Process: unable to parse response from audit event: no namespace",
			Subtype:              ResponseType,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"request-basic-input-and-request-with-ns": {
			IsErrorExpected: false,
			Subtype:         RequestType,
			Data:            &logical.LogInput{Request: &logical.Request{ID: "123"}},
			RootNamespace:   true,
		},
		"response-basic-input-and-request-with-ns": {
			IsErrorExpected: false,
			Subtype:         ResponseType,
			Data:            &logical.LogInput{Request: &logical.Request{ID: "123"}},
			RootNamespace:   true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := fakeJSONAuditEvent(t, tc.Subtype, tc.Data)
			require.NotNil(t, e)

			ss := newStaticSalt(t)
			cfg := FormatterConfig{}

			f, err := NewAuditFormatterJSON(cfg, ss)
			require.NoError(t, err)
			require.NotNil(t, f)

			var ctx context.Context
			switch {
			case tc.RootNamespace:
				ctx = namespace.RootContext(context.Background())
			default:
				ctx = context.Background()
			}

			processed, err := f.Process(ctx, e)
			b, found := e.Format(string(JSONFormat))

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
