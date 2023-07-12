// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"

	"github.com/hashicorp/vault/sdk/logical"

	vaultaudit "github.com/hashicorp/vault/audit"

	"github.com/hashicorp/vault/sdk/helper/salt"

	"github.com/hashicorp/eventlogger"
	"github.com/stretchr/testify/require"
)

// fakeJSONAuditEvent will return a new fake event containing audit data based
// on the specified auditSubtype and logical.LogInput.
func fakeJSONAuditEvent(t *testing.T, subtype auditSubtype, input *logical.LogInput) *eventlogger.Event {
	t.Helper()

	date := time.Date(2023, time.July, 11, 15, 49, 10, 0o0, time.Local)

	auditEvent, err := newAudit(
		WithID("123"),
		WithSubtype(string(subtype)),
		WithFormat(string(AuditFormatJSON)),
		WithNow(date),
	)
	require.NoError(t, err)
	require.NotNil(t, auditEvent)
	require.Equal(t, "123", auditEvent.ID)
	require.Equal(t, "v0.1", auditEvent.Version)
	require.Equal(t, AuditFormatJSON, auditEvent.RequiredFormat)
	require.Equal(t, subtype, auditEvent.Subtype)
	require.Equal(t, date, auditEvent.Timestamp)

	auditEvent.Data = input

	e := &eventlogger.Event{
		Type:      eventlogger.EventType(AuditType),
		CreatedAt: auditEvent.Timestamp,
		Formatted: make(map[string][]byte),
		Payload:   auditEvent,
	}

	return e
}

// newStaticSalt returns a new staticSalt for use in testing.
func newStaticSalt(t *testing.T) *staticSalt {
	s, err := salt.NewSalt(context.Background(), nil, nil)
	require.NoError(t, err)

	return &staticSalt{salt: s}
}

// staticSalt is a struct which can be used to obtain a static salt.
// a salt must be assigned when the struct is initialized.
type staticSalt struct {
	salt *salt.Salt
}

// Salt returns the static salt and no error.
func (s *staticSalt) Salt(_ context.Context) (*salt.Salt, error) {
	return s.salt, nil
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
			ExpectedErrorMessage: "event.NewAuditFormatterJSON: unable to create new JSON audit formatter: cannot create a new audit formatter with nil salter",
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
			cfg := vaultaudit.FormatterConfig{}
			var ss vaultaudit.Salter
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
	cfg := vaultaudit.FormatterConfig{}

	f, err := NewAuditFormatterJSON(cfg, ss)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.NoError(t, f.Reopen())
}

// TestAuditFormatterJSONX_Type ensures that the node is a 'formatter' type.
func TestAuditFormatterJSON_Type(t *testing.T) {
	ss := newStaticSalt(t)
	cfg := vaultaudit.FormatterConfig{}

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
		Subtype              auditSubtype
		Data                 *logical.LogInput
		RootNamespace        bool
	}{
		"request-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFormatterJSON).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              AuditRequest,
			Data:                 nil,
		},
		"response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFormatterJSON).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              AuditResponse,
			Data:                 nil,
		},
		"request-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFormatterJSON).Process: unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              AuditRequest,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"response-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFormatterJSON).Process: unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              AuditResponse,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"request-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFormatterJSON).Process: unable to parse request from audit event: no namespace",
			Subtype:              AuditRequest,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"response-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "event.(AuditFormatterJSON).Process: unable to parse response from audit event: no namespace",
			Subtype:              AuditResponse,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"request-basic-input-and-request-with-ns": {
			IsErrorExpected: false,
			Subtype:         AuditRequest,
			Data:            &logical.LogInput{Request: &logical.Request{ID: "123"}},
			RootNamespace:   true,
		},
		"response-basic-input-and-request-with-ns": {
			IsErrorExpected: false,
			Subtype:         AuditResponse,
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
			cfg := vaultaudit.FormatterConfig{}

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
			b, found := e.Format(string(AuditFormatJSON))

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
