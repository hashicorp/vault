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

	auditEvent, err := NewEvent(subtype,
		WithID("123"),
		WithNow(date),
	)
	require.NoError(tb, err)
	require.NotNil(tb, auditEvent)
	require.Equal(tb, "123", auditEvent.ID)
	require.Equal(tb, "v0.1", auditEvent.Version)
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

// TestNewEntryFormatter ensures we can create new EntryFormatter structs.
func TestNewEntryFormatter(t *testing.T) {
	tests := map[string]struct {
		UseStaticSalt        bool
		Options              []Option // Only supports WithPrefix
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedFormat       format
		ExpectedPrefix       string
	}{
		"nil-salter": {
			UseStaticSalt:        false,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.NewEntryFormatter: cannot create a new audit formatter with nil salter: invalid parameter",
		},
		"static-salter": {
			UseStaticSalt:   true,
			IsErrorExpected: false,
			Options: []Option{
				WithFormat(JSONFormat.String()),
			},
			ExpectedFormat: JSONFormat,
		},
		"default": {
			UseStaticSalt:   true,
			IsErrorExpected: false,
			ExpectedFormat:  JSONFormat,
		},
		"config-json": {
			UseStaticSalt: true,
			Options: []Option{
				WithFormat(JSONFormat.String()),
			},
			IsErrorExpected: false,
			ExpectedFormat:  JSONFormat,
		},
		"config-jsonx": {
			UseStaticSalt: true,
			Options: []Option{
				WithFormat(JSONxFormat.String()),
			},
			IsErrorExpected: false,
			ExpectedFormat:  JSONxFormat,
		},
		"config-json-prefix": {
			UseStaticSalt: true,
			Options: []Option{
				WithPrefix("foo"),
				WithFormat(JSONFormat.String()),
			},
			IsErrorExpected: false,
			ExpectedFormat:  JSONFormat,
			ExpectedPrefix:  "foo",
		},
		"config-jsonx-prefix": {
			UseStaticSalt: true,
			Options: []Option{
				WithPrefix("foo"),
				WithFormat(JSONxFormat.String()),
			},
			IsErrorExpected: false,
			ExpectedFormat:  JSONxFormat,
			ExpectedPrefix:  "foo",
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

			cfg, err := NewFormatterConfig(tc.Options...)
			require.NoError(t, err)
			f, err := NewEntryFormatter(cfg, ss, tc.Options...)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, f)
			default:
				require.NoError(t, err)
				require.NotNil(t, f)
				require.Equal(t, tc.ExpectedFormat, f.config.RequiredFormat)
				require.Equal(t, tc.ExpectedPrefix, f.prefix)
			}
		})
	}
}

// TestEntryFormatter_Reopen ensures that we do not get an error when calling Reopen.
func TestEntryFormatter_Reopen(t *testing.T) {
	ss := newStaticSalt(t)
	cfg, err := NewFormatterConfig()
	require.NoError(t, err)

	f, err := NewEntryFormatter(cfg, ss)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.NoError(t, f.Reopen())
}

// TestEntryFormatter_Type ensures that the node is a 'formatter' type.
func TestEntryFormatter_Type(t *testing.T) {
	ss := newStaticSalt(t)
	cfg, err := NewFormatterConfig()
	require.NoError(t, err)

	f, err := NewEntryFormatter(cfg, ss)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.Equal(t, eventlogger.NodeTypeFormatter, f.Type())
}

// TestEntryFormatter_Process attempts to run the Process method to convert the
// logical.LogInput within an audit event to JSON and JSONx (RequestEntry or ResponseEntry).
func TestEntryFormatter_Process(t *testing.T) {
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
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 nil,
		},
		"json-response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONFormat,
			Data:                 nil,
		},
		"json-request-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"json-response-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"json-request-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse request from audit event: no namespace",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"json-response-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse response from audit event: no namespace",
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
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 nil,
		},
		"jsonx-response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONxFormat,
			Data:                 nil,
		},
		"jsonx-request-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"jsonx-response-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"jsonx-request-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse request from audit event: no namespace",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"jsonx-response-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: unable to parse response from audit event: no namespace",
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
			cfg, err := NewFormatterConfig(WithFormat(tc.RequiredFormat.String()))
			require.NoError(t, err)

			f, err := NewEntryFormatter(cfg, ss)
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

// BenchmarkAuditFileSink_Process benchmarks the EntryFormatter and then event.FileSink calling Process.
// This should replicate the original benchmark testing which used to perform both of these roles together.
func BenchmarkAuditFileSink_Process(b *testing.B) {
	// Base input
	in := &logical.LogInput{
		Auth: &logical.Auth{
			ClientToken:     "foo",
			Accessor:        "bar",
			EntityID:        "foobarentity",
			DisplayName:     "testtoken",
			NoDefaultPolicy: true,
			Policies:        []string{"root"},
			TokenType:       logical.TokenTypeService,
		},
		Request: &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "/foo",
			Connection: &logical.Connection{
				RemoteAddr: "127.0.0.1",
			},
			WrapInfo: &logical.RequestWrapInfo{
				TTL: 60 * time.Second,
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
	}

	ctx := namespace.RootContext(context.Background())

	// Create the formatter node.
	cfg, err := NewFormatterConfig()
	require.NoError(b, err)
	ss := newStaticSalt(b)
	formatter, err := NewEntryFormatter(cfg, ss)
	require.NoError(b, err)
	require.NotNil(b, formatter)

	// Create the sink node.
	sink, err := event.NewFileSink("/dev/null", JSONFormat.String())
	require.NoError(b, err)
	require.NotNil(b, sink)

	// Generate the event
	event := fakeEvent(b, RequestType, JSONFormat, in)
	require.NotNil(b, event)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			event, err = formatter.Process(ctx, event)
			if err != nil {
				panic(err)
			}
			_, err := sink.Process(ctx, event)
			if err != nil {
				panic(err)
			}
		}
	})
}
