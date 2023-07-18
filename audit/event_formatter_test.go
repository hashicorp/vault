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

// fakeEvent will return a new fake event containing audit data based  on the
// specified subtype, format and logical.LogInput.
func fakeEvent(tb testing.TB, subtype subtype, format format, input *logical.LogInput) *eventlogger.Event {
	tb.Helper()

	date := time.Date(2023, time.July, 11, 15, 49, 10, 0o0, time.Local)

	auditEvent, err := newEvent(subtype, format,
		WithID("123"),
		WithNow(date),
	)
	require.NoError(tb, err)
	require.NotNil(tb, auditEvent)
	require.Equal(tb, "123", auditEvent.ID)
	require.Equal(tb, "v0.1", auditEvent.Version)
	require.Equal(tb, format, auditEvent.RequiredFormat)
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

// TestNewEventFormatter ensures we can create new EventFormatter structs.
func TestNewEventFormatter(t *testing.T) {
	tests := map[string]struct {
		UseStaticSalt        bool
		Config               FormatterConfig
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedFormat       format
	}{
		"nil-salter": {
			UseStaticSalt:        false,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.NewEventFormatter: cannot create a new audit formatter with nil salter: invalid parameter",
		},
		"static-salter": {
			UseStaticSalt:   true,
			IsErrorExpected: false,
			ExpectedFormat:  JSONFormat,
		},
		"default": {
			UseStaticSalt:   true,
			Config:          FormatterConfig{},
			IsErrorExpected: false,
			ExpectedFormat:  JSONFormat,
		},
		"config-json": {
			UseStaticSalt:   true,
			Config:          FormatterConfig{RequiredFormat: JSONFormat},
			IsErrorExpected: false,
			ExpectedFormat:  JSONFormat,
		},
		"config-jsonx": {
			UseStaticSalt:   true,
			Config:          FormatterConfig{RequiredFormat: JSONxFormat},
			IsErrorExpected: false,
			ExpectedFormat:  JSONxFormat,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var ss Salter
			if tc.UseStaticSalt {
				ss = newStaticSalt(t)
			}

			f, err := NewEventFormatter(tc.Config, ss)

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

// TestEventFormatter_Reopen ensures that we do not get an error when calling Reopen.
func TestEventFormatter_Reopen(t *testing.T) {
	ss := newStaticSalt(t)
	cfg := FormatterConfig{}

	f, err := NewEventFormatter(cfg, ss)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.NoError(t, f.Reopen())
}

// TestEventFormatter_Type ensures that the node is a 'formatter' type.
func TestEventFormatter_Type(t *testing.T) {
	ss := newStaticSalt(t)
	cfg := FormatterConfig{}

	f, err := NewEventFormatter(cfg, ss)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.Equal(t, eventlogger.NodeTypeFormatter, f.Type())
}

// TestEventFormatter_Process attempts to run the Process method to convert the
// logical.LogInput within an audit event to JSON and JSONx (RequestEntry or ResponseEntry).
func TestEventFormatter_Process(t *testing.T) {
	tests := map[string]struct {
		IsErrorExpected      bool
		ExpectedErrorMessage string
		Subtype              subtype
		RequiredFormat       format
		Data                 *logical.LogInput
		RootNamespace        bool
	}{
		"json-request-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 nil,
		},
		"json-response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONFormat,
			Data:                 nil,
		},
		"json-request-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"json-response-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"json-request-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse request from audit event: no namespace",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"json-response-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse response from audit event: no namespace",
			Subtype:              ResponseType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"json-request-basic-input-and-request-with-ns": {
			IsErrorExpected: false,
			Subtype:         RequestType,
			RequiredFormat:  JSONFormat,
			Data:            &logical.LogInput{Request: &logical.Request{ID: "123"}},
			RootNamespace:   true,
		},
		"json-response-basic-input-and-request-with-ns": {
			IsErrorExpected: false,
			Subtype:         ResponseType,
			RequiredFormat:  JSONFormat,
			Data:            &logical.LogInput{Request: &logical.Request{ID: "123"}},
			RootNamespace:   true,
		},
		"jsonx-request-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 nil,
		},
		"jsonx-response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONxFormat,
			Data:                 nil,
		},
		"jsonx-request-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"jsonx-response-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"jsonx-request-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse request from audit event: no namespace",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"jsonx-response-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EventFormatter).Process: unable to parse response from audit event: no namespace",
			Subtype:              ResponseType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"jsonx-request-basic-input-and-request-with-ns": {
			IsErrorExpected: false,
			Subtype:         RequestType,
			RequiredFormat:  JSONxFormat,
			Data:            &logical.LogInput{Request: &logical.Request{ID: "123"}},
			RootNamespace:   true,
		},
		"jsonx-response-basic-input-and-request-with-ns": {
			IsErrorExpected: false,
			Subtype:         ResponseType,
			RequiredFormat:  JSONxFormat,
			Data:            &logical.LogInput{Request: &logical.Request{ID: "123"}},
			RootNamespace:   true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := fakeEvent(t, tc.Subtype, tc.RequiredFormat, tc.Data)
			require.NotNil(t, e)

			ss := newStaticSalt(t)
			cfg := FormatterConfig{
				RequiredFormat: tc.RequiredFormat,
			}

			f, err := NewEventFormatter(cfg, ss)
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
			b, found := e.Format(string(tc.RequiredFormat))

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
