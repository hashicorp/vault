// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testFormatJSONReqBasicStrFmt = `
{
  "time": "2015-08-05T13:45:46Z",
  "type": "request",
  "auth": {
    "client_token": "%s",
    "accessor": "bar",
    "display_name": "testtoken",
    "policies": [
      "root"
    ],
    "no_default_policy": true,
    "metadata": null,
    "entity_id": "foobarentity",
    "token_type": "service",
    "token_ttl": 14400,
    "token_issue_time": "2020-05-28T13:40:18-05:00"
  },
  "request": {
    "operation": "update",
    "path": "/foo",
    "data": null,
    "wrap_ttl": 60,
    "remote_address": "127.0.0.1",
    "headers": {
      "foo": [
        "bar"
      ]
    }
  },
  "error": "this is an error"
}
`

// TestNewEntryFormatter ensures we can create new EntryFormatter structs.
func TestNewEntryFormatter(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: cannot audit event (AuditRequest) with no data: invalid parameter",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 nil,
		},
		"json-response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: cannot audit event (AuditResponse) with no data: invalid parameter",
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
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: cannot audit event (AuditRequest) with no data: invalid parameter",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 nil,
		},
		"jsonx-response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(EntryFormatter).Process: cannot audit event (AuditResponse) with no data: invalid parameter",
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
			e := fakeEvent(t, tc.Subtype, tc.Data)
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

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, processed)
			default:
				require.NoError(t, err)
				require.NotNil(t, processed)
				b, found := processed.Format(string(tc.RequiredFormat))
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
	e := fakeEvent(b, RequestType, in)
	require.NotNil(b, e)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e, err = formatter.Process(ctx, e)
			if err != nil {
				panic(err)
			}
			_, err := sink.Process(ctx, e)
			if err != nil {
				panic(err)
			}
		}
	})
}

// TestEntryFormatter_FormatRequest exercises EntryFormatter.FormatRequest with
// varying inputs.
func TestEntryFormatter_FormatRequest(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Input                *logical.LogInput
		IsErrorExpected      bool
		ExpectedErrorMessage string
		RootNamespace        bool
	}{
		"nil": {
			Input:                nil,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "request to request-audit a nil request",
		},
		"basic-input": {
			Input:                &logical.LogInput{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "request to request-audit a nil request",
		},
		"input-and-request-no-ns": {
			Input:                &logical.LogInput{Request: &logical.Request{ID: "123"}},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "no namespace",
			RootNamespace:        false,
		},
		"input-and-request-with-ns": {
			Input:           &logical.LogInput{Request: &logical.Request{ID: "123"}},
			IsErrorExpected: false,
			RootNamespace:   true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := newStaticSalt(t)
			cfg, err := NewFormatterConfig()
			require.NoError(t, err)
			f, err := NewEntryFormatter(cfg, ss)
			require.NoError(t, err)

			var ctx context.Context
			switch {
			case tc.RootNamespace:
				ctx = namespace.RootContext(context.Background())
			default:
				ctx = context.Background()
			}

			entry, err := f.FormatRequest(ctx, tc.Input)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, entry)
			default:
				require.NoError(t, err)
				require.NotNil(t, entry)
			}
		})
	}
}

// TestEntryFormatter_FormatResponse exercises EntryFormatter.FormatResponse with
// varying inputs.
func TestEntryFormatter_FormatResponse(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Input                *logical.LogInput
		IsErrorExpected      bool
		ExpectedErrorMessage string
		RootNamespace        bool
	}{
		"nil": {
			Input:                nil,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "request to response-audit a nil request",
		},
		"basic-input": {
			Input:                &logical.LogInput{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "request to response-audit a nil request",
		},
		"input-and-request-no-ns": {
			Input:                &logical.LogInput{Request: &logical.Request{ID: "123"}},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "no namespace",
			RootNamespace:        false,
		},
		"input-and-request-with-ns": {
			Input:           &logical.LogInput{Request: &logical.Request{ID: "123"}},
			IsErrorExpected: false,
			RootNamespace:   true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := newStaticSalt(t)
			cfg, err := NewFormatterConfig()
			require.NoError(t, err)
			f, err := NewEntryFormatter(cfg, ss)
			require.NoError(t, err)

			var ctx context.Context
			switch {
			case tc.RootNamespace:
				ctx = namespace.RootContext(context.Background())
			default:
				ctx = context.Background()
			}

			entry, err := f.FormatResponse(ctx, tc.Input)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, entry)
			default:
				require.NoError(t, err)
				require.NotNil(t, entry)
			}
		})
	}
}

// TestEntryFormatter_Process_JSON ensures that the JSON output we get matches what
// we expect for the specified LogInput.
func TestEntryFormatter_Process_JSON(t *testing.T) {
	t.Parallel()

	ss := newStaticSalt(t)

	expectedResultStr := fmt.Sprintf(testFormatJSONReqBasicStrFmt, ss.salt.GetIdentifiedHMAC("foo"))

	issueTime, _ := time.Parse(time.RFC3339, "2020-05-28T13:40:18-05:00")
	cases := map[string]struct {
		Auth        *logical.Auth
		Req         *logical.Request
		Err         error
		Prefix      string
		ExpectedStr string
	}{
		"auth, request": {
			&logical.Auth{
				ClientToken:     "foo",
				Accessor:        "bar",
				DisplayName:     "testtoken",
				EntityID:        "foobarentity",
				NoDefaultPolicy: true,
				Policies:        []string{"root"},
				TokenType:       logical.TokenTypeService,
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Hour * 4,
					IssueTime: issueTime,
				},
			},
			&logical.Request{
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
			errors.New("this is an error"),
			"",
			expectedResultStr,
		},
		"auth, request with prefix": {
			&logical.Auth{
				ClientToken:     "foo",
				Accessor:        "bar",
				EntityID:        "foobarentity",
				DisplayName:     "testtoken",
				NoDefaultPolicy: true,
				Policies:        []string{"root"},
				TokenType:       logical.TokenTypeService,
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Hour * 4,
					IssueTime: issueTime,
				},
			},
			&logical.Request{
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
			errors.New("this is an error"),
			"@cee: ",
			expectedResultStr,
		},
	}

	for name, tc := range cases {
		cfg, err := NewFormatterConfig(WithHMACAccessor(false))
		require.NoError(t, err)
		formatter, err := NewEntryFormatter(cfg, ss, WithPrefix(tc.Prefix))
		require.NoError(t, err)

		in := &logical.LogInput{
			Auth:     tc.Auth,
			Request:  tc.Req,
			OuterErr: tc.Err,
		}

		// Create an audit event and more generic eventlogger.event to allow us
		// to process (format).
		auditEvent, err := NewEvent(RequestType)
		require.NoError(t, err)
		auditEvent.Data = in

		e := &eventlogger.Event{
			Type:      eventlogger.EventType(event.AuditType.String()),
			CreatedAt: time.Now(),
			Formatted: make(map[string][]byte),
			Payload:   auditEvent,
		}

		e2, err := formatter.Process(namespace.RootContext(nil), e)
		require.NoErrorf(t, err, "bad: %s\nerr: %s", name, err)

		jsonBytes, ok := e2.Format(JSONFormat.String())
		require.True(t, ok)
		require.Positive(t, len(jsonBytes))

		if !strings.HasPrefix(string(jsonBytes), tc.Prefix) {
			t.Fatalf("no prefix: %s \n log: %s\nprefix: %s", name, expectedResultStr, tc.Prefix)
		}

		expectedJSON := new(RequestEntry)

		if err := jsonutil.DecodeJSON([]byte(expectedResultStr), &expectedJSON); err != nil {
			t.Fatalf("bad json: %s", err)
		}
		expectedJSON.Request.Namespace = &Namespace{ID: "root"}

		actualJSON := new(RequestEntry)
		if err := jsonutil.DecodeJSON(jsonBytes[len(tc.Prefix):], &actualJSON); err != nil {
			t.Fatalf("bad json: %s", err)
		}

		expectedJSON.Time = actualJSON.Time

		expectedBytes, err := json.Marshal(expectedJSON)
		if err != nil {
			t.Fatalf("unable to marshal json: %s", err)
		}

		if !strings.HasSuffix(strings.TrimSpace(string(jsonBytes)), string(expectedBytes)) {
			t.Fatalf("bad: %s\nResult:\n\n%q\n\nExpected:\n\n%q", name, string(jsonBytes), string(expectedBytes))
		}
	}
}

// TestEntryFormatter_Process_JSONx ensures that the JSONx output we get matches what
// we expect for the specified LogInput.
func TestEntryFormatter_Process_JSONx(t *testing.T) {
	t.Parallel()

	s, err := salt.NewSalt(context.Background(), nil, nil)
	require.NoError(t, err)
	tempStaticSalt := &staticSalt{salt: s}

	fooSalted := s.GetIdentifiedHMAC("foo")
	issueTime, _ := time.Parse(time.RFC3339, "2020-05-28T13:40:18-05:00")

	cases := map[string]struct {
		Auth        *logical.Auth
		Req         *logical.Request
		Err         error
		Prefix      string
		Result      string
		ExpectedStr string
	}{
		"auth, request": {
			&logical.Auth{
				ClientToken:     "foo",
				Accessor:        "bar",
				DisplayName:     "testtoken",
				EntityID:        "foobarentity",
				NoDefaultPolicy: true,
				Policies:        []string{"root"},
				TokenType:       logical.TokenTypeService,
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Hour * 4,
					IssueTime: issueTime,
				},
			},
			&logical.Request{
				ID:                  "request",
				ClientToken:         "foo",
				ClientTokenAccessor: "bar",
				Operation:           logical.UpdateOperation,
				Path:                "/foo",
				Connection: &logical.Connection{
					RemoteAddr: "127.0.0.1",
				},
				WrapInfo: &logical.RequestWrapInfo{
					TTL: 60 * time.Second,
				},
				Headers: map[string][]string{
					"foo": {"bar"},
				},
				PolicyOverride: true,
			},
			errors.New("this is an error"),
			"",
			"",
			fmt.Sprintf(`<json:object name="auth"><json:string name="accessor">bar</json:string><json:string name="client_token">%s</json:string><json:string name="display_name">testtoken</json:string><json:string name="entity_id">foobarentity</json:string><json:boolean name="no_default_policy">true</json:boolean><json:array name="policies"><json:string>root</json:string></json:array><json:string name="token_issue_time">2020-05-28T13:40:18-05:00</json:string><json:number name="token_ttl">14400</json:number><json:string name="token_type">service</json:string></json:object><json:string name="error">this is an error</json:string><json:object name="request"><json:string name="client_token">%s</json:string><json:string name="client_token_accessor">bar</json:string><json:object name="headers"><json:array name="foo"><json:string>bar</json:string></json:array></json:object><json:string name="id">request</json:string><json:object name="namespace"><json:string name="id">root</json:string></json:object><json:string name="operation">update</json:string><json:string name="path">/foo</json:string><json:boolean name="policy_override">true</json:boolean><json:string name="remote_address">127.0.0.1</json:string><json:number name="wrap_ttl">60</json:number></json:object><json:string name="type">request</json:string>`,
				fooSalted, fooSalted),
		},
		"auth, request with prefix": {
			&logical.Auth{
				ClientToken:     "foo",
				Accessor:        "bar",
				DisplayName:     "testtoken",
				NoDefaultPolicy: true,
				EntityID:        "foobarentity",
				Policies:        []string{"root"},
				TokenType:       logical.TokenTypeService,
				LeaseOptions: logical.LeaseOptions{
					TTL:       time.Hour * 4,
					IssueTime: issueTime,
				},
			},
			&logical.Request{
				ID:                  "request",
				ClientToken:         "foo",
				ClientTokenAccessor: "bar",
				Operation:           logical.UpdateOperation,
				Path:                "/foo",
				Connection: &logical.Connection{
					RemoteAddr: "127.0.0.1",
				},
				WrapInfo: &logical.RequestWrapInfo{
					TTL: 60 * time.Second,
				},
				Headers: map[string][]string{
					"foo": {"bar"},
				},
				PolicyOverride: true,
			},
			errors.New("this is an error"),
			"",
			"@cee: ",
			fmt.Sprintf(`<json:object name="auth"><json:string name="accessor">bar</json:string><json:string name="client_token">%s</json:string><json:string name="display_name">testtoken</json:string><json:string name="entity_id">foobarentity</json:string><json:boolean name="no_default_policy">true</json:boolean><json:array name="policies"><json:string>root</json:string></json:array><json:string name="token_issue_time">2020-05-28T13:40:18-05:00</json:string><json:number name="token_ttl">14400</json:number><json:string name="token_type">service</json:string></json:object><json:string name="error">this is an error</json:string><json:object name="request"><json:string name="client_token">%s</json:string><json:string name="client_token_accessor">bar</json:string><json:object name="headers"><json:array name="foo"><json:string>bar</json:string></json:array></json:object><json:string name="id">request</json:string><json:object name="namespace"><json:string name="id">root</json:string></json:object><json:string name="operation">update</json:string><json:string name="path">/foo</json:string><json:boolean name="policy_override">true</json:boolean><json:string name="remote_address">127.0.0.1</json:string><json:number name="wrap_ttl">60</json:number></json:object><json:string name="type">request</json:string>`,
				fooSalted, fooSalted),
		},
	}

	for name, tc := range cases {
		cfg, err := NewFormatterConfig(
			WithOmitTime(true),
			WithHMACAccessor(false),
			WithFormat(JSONxFormat.String()),
		)
		require.NoError(t, err)
		formatter, err := NewEntryFormatter(cfg, tempStaticSalt, WithPrefix(tc.Prefix))
		require.NoError(t, err)
		require.NotNil(t, formatter)

		in := &logical.LogInput{
			Auth:     tc.Auth,
			Request:  tc.Req,
			OuterErr: tc.Err,
		}

		// Create an audit event and more generic eventlogger.event to allow us
		// to process (format).
		auditEvent, err := NewEvent(RequestType)
		require.NoError(t, err)
		auditEvent.Data = in

		e := &eventlogger.Event{
			Type:      eventlogger.EventType(event.AuditType.String()),
			CreatedAt: time.Now(),
			Formatted: make(map[string][]byte),
			Payload:   auditEvent,
		}

		e2, err := formatter.Process(namespace.RootContext(nil), e)
		require.NoErrorf(t, err, "bad: %s\nerr: %s", name, err)

		jsonxBytes, ok := e2.Format(JSONxFormat.String())
		require.True(t, ok)
		require.Positive(t, len(jsonxBytes))

		if !strings.HasPrefix(string(jsonxBytes), tc.Prefix) {
			t.Fatalf("no prefix: %s \n log: %s\nprefix: %s", name, tc.Result, tc.Prefix)
		}

		if !strings.HasSuffix(strings.TrimSpace(string(jsonxBytes)), string(tc.ExpectedStr)) {
			t.Fatalf(
				"bad: %s\nResult:\n\n%q\n\nExpected:\n\n%q",
				name, strings.TrimSpace(string(jsonxBytes)), string(tc.ExpectedStr))
		}
	}
}

// TestEntryFormatter_FormatResponse_ElideListResponses ensures that we correctly
// elide data in responses to LIST operations.
func TestEntryFormatter_FormatResponse_ElideListResponses(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		inputData    map[string]any
		expectedData map[string]any
	}{
		"nil data": {
			nil,
			nil,
		},
		"Normal list (keys only)": {
			map[string]any{
				"keys": []string{"foo", "bar", "baz"},
			},
			map[string]any{
				"keys": 3,
			},
		},
		"Enhanced list (has key_info)": {
			map[string]any{
				"keys": []string{"foo", "bar", "baz", "quux"},
				"key_info": map[string]any{
					"foo":  "alpha",
					"bar":  "beta",
					"baz":  "gamma",
					"quux": "delta",
				},
			},
			map[string]any{
				"keys":     4,
				"key_info": 4,
			},
		},
		"Unconventional other values in a list response are not touched": {
			map[string]any{
				"keys":           []string{"foo", "bar"},
				"something_else": "baz",
			},
			map[string]any{
				"keys":           2,
				"something_else": "baz",
			},
		},
		"Conventional values in a list response are not elided if their data types are unconventional": {
			map[string]any{
				"keys": map[string]any{
					"You wouldn't expect keys to be a map": nil,
				},
				"key_info": []string{
					"You wouldn't expect key_info to be a slice",
				},
			},
			map[string]any{
				"keys": map[string]any{
					"You wouldn't expect keys to be a map": nil,
				},
				"key_info": []string{
					"You wouldn't expect key_info to be a slice",
				},
			},
		},
	}

	oneInterestingTestCase := tests["Enhanced list (has key_info)"]

	ss := newStaticSalt(t)
	ctx := namespace.RootContext(context.Background())
	var formatter *EntryFormatter
	var err error

	format := func(t *testing.T, config FormatterConfig, operation logical.Operation, inputData map[string]any) *ResponseEntry {
		formatter, err = NewEntryFormatter(config, ss)
		require.NoError(t, err)
		require.NotNil(t, formatter)

		in := &logical.LogInput{
			Request:  &logical.Request{Operation: operation},
			Response: &logical.Response{Data: inputData},
		}

		resp, err := formatter.FormatResponse(ctx, in)
		require.NoError(t, err)

		return resp
	}

	t.Run("Default case", func(t *testing.T) {
		config, err := NewFormatterConfig(WithElision(true))
		require.NoError(t, err)
		for name, tc := range tests {
			name := name
			tc := tc
			t.Run(name, func(t *testing.T) {
				entry := format(t, config, logical.ListOperation, tc.inputData)
				assert.Equal(t, formatter.hashExpectedValueForComparison(tc.expectedData), entry.Response.Data)
			})
		}
	})

	t.Run("When Operation is not list, eliding does not happen", func(t *testing.T) {
		config, err := NewFormatterConfig(WithElision(true))
		require.NoError(t, err)
		tc := oneInterestingTestCase
		entry := format(t, config, logical.ReadOperation, tc.inputData)
		assert.Equal(t, formatter.hashExpectedValueForComparison(tc.inputData), entry.Response.Data)
	})

	t.Run("When ElideListResponses is false, eliding does not happen", func(t *testing.T) {
		config, err := NewFormatterConfig(WithElision(false), WithFormat(JSONFormat.String()))
		require.NoError(t, err)
		tc := oneInterestingTestCase
		entry := format(t, config, logical.ListOperation, tc.inputData)
		assert.Equal(t, formatter.hashExpectedValueForComparison(tc.inputData), entry.Response.Data)
	})

	t.Run("When Raw is true, eliding still happens", func(t *testing.T) {
		config, err := NewFormatterConfig(WithElision(true), WithRaw(true), WithFormat(JSONFormat.String()))
		require.NoError(t, err)
		tc := oneInterestingTestCase
		entry := format(t, config, logical.ListOperation, tc.inputData)
		assert.Equal(t, tc.expectedData, entry.Response.Data)
	})
}

// TestEntryFormatter_Process_NoMutation tests that the event returned by an
// EntryFormatter.Process method is not the same as the one that it accepted.
func TestEntryFormatter_Process_NoMutation(t *testing.T) {
	t.Parallel()

	// Create the formatter node.
	cfg, err := NewFormatterConfig()
	require.NoError(t, err)
	ss := newStaticSalt(t)
	formatter, err := NewEntryFormatter(cfg, ss)
	require.NoError(t, err)
	require.NotNil(t, formatter)

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

	e := fakeEvent(t, RequestType, in)

	e2, err := formatter.Process(namespace.RootContext(nil), e)
	require.NoError(t, err)
	require.NotNil(t, e2)

	// Ensure the pointers are different.
	require.NotEqual(t, e2, e)

	// Do the same for the audit event in the payload.
	a, ok := e.Payload.(*AuditEvent)
	require.True(t, ok)
	require.NotNil(t, a)

	a2, ok := e2.Payload.(*AuditEvent)
	require.True(t, ok)
	require.NotNil(t, a2)

	require.NotEqual(t, a2, a)
}

// hashExpectedValueForComparison replicates enough of the audit HMAC process on a piece of expected data in a test,
// so that we can use assert.Equal to compare the expected and output values.
func (f *EntryFormatter) hashExpectedValueForComparison(input map[string]any) map[string]any {
	// Copy input before modifying, since we may re-use the same data in another test
	copied, err := copystructure.Copy(input)
	if err != nil {
		panic(err)
	}
	copiedAsMap := copied.(map[string]any)

	s, err := f.salter.Salt(context.Background())
	if err != nil {
		panic(err)
	}

	err = hashMap(s.GetIdentifiedHMAC, copiedAsMap, nil)
	if err != nil {
		panic(err)
	}

	return copiedAsMap
}

// fakeEvent will return a new fake event containing audit data based  on the
// specified subtype, format and logical.LogInput.
func fakeEvent(tb testing.TB, subtype subtype, input *logical.LogInput) *eventlogger.Event {
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

// newStaticSalt returns a new staticSalt for use in testing.
func newStaticSalt(tb testing.TB) *staticSalt {
	s, err := salt.NewSalt(context.Background(), nil, nil)
	require.NoError(tb, err)

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
