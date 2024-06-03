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
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-sockaddr"
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

// testHeaderFormatter is a stub to prevent the need to import the vault package
// to bring in vault.HeadersConfig for testing.
type testHeaderFormatter struct {
	shouldReturnEmpty bool
}

// ApplyConfig satisfies the HeaderFormatter interface for testing.
// It will either return the headers it was supplied or empty headers depending
// on how it is configured.
// ignore-nil-nil-function-check.
func (f *testHeaderFormatter) ApplyConfig(_ context.Context, headers map[string][]string, salter Salter) (result map[string][]string, retErr error) {
	if f.shouldReturnEmpty {
		return make(map[string][]string), nil
	}

	return headers, nil
}

// testTimeProvider is just a test struct used to imitate an AuditEvent's ability
// to provide a formatted time.
type testTimeProvider struct{}

// formattedTime always returns the same value for 22nd March 2024 at 10:00:05 (and 10 nanos).
func (p *testTimeProvider) formattedTime() string {
	return time.Date(2024, time.March, 22, 10, 0o0, 5, 10, time.UTC).UTC().Format(time.RFC3339Nano)
}

// TestNewEntryFormatter ensures we can create new entryFormatter structs.
func TestNewEntryFormatter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Name                 string
		UseStaticSalt        bool
		Logger               hclog.Logger
		Options              map[string]string
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedFormat       format
		ExpectedPrefix       string
	}{
		"empty-name": {
			Name:                 "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "name is required: invalid internal parameter",
		},
		"spacey-name": {
			Name:                 "   ",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "name is required: invalid internal parameter",
		},
		"nil-salter": {
			Name:                 "juan",
			UseStaticSalt:        false,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot create a new audit formatter with nil salter: invalid internal parameter",
		},
		"nil-logger": {
			Name:                 "juan",
			UseStaticSalt:        true,
			Logger:               nil,
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot create a new audit formatter with nil logger: invalid internal parameter",
		},
		"static-salter": {
			Name:            "juan",
			UseStaticSalt:   true,
			Logger:          hclog.NewNullLogger(),
			IsErrorExpected: false,
			Options: map[string]string{
				"format": "json",
			},
			ExpectedFormat: JSONFormat,
		},
		"default": {
			Name:            "juan",
			UseStaticSalt:   true,
			Logger:          hclog.NewNullLogger(),
			IsErrorExpected: false,
			ExpectedFormat:  JSONFormat,
		},
		"config-json": {
			Name:          "juan",
			UseStaticSalt: true,
			Logger:        hclog.NewNullLogger(),
			Options: map[string]string{
				"format": "json",
			},
			IsErrorExpected: false,
			ExpectedFormat:  JSONFormat,
		},
		"config-jsonx": {
			Name:          "juan",
			UseStaticSalt: true,
			Logger:        hclog.NewNullLogger(),
			Options: map[string]string{
				"format": "jsonx",
			},
			IsErrorExpected: false,
			ExpectedFormat:  JSONxFormat,
		},
		"config-json-prefix": {
			Name:          "juan",
			UseStaticSalt: true,
			Logger:        hclog.NewNullLogger(),
			Options: map[string]string{
				"prefix": "foo",
				"format": "json",
			},
			IsErrorExpected: false,
			ExpectedFormat:  JSONFormat,
			ExpectedPrefix:  "foo",
		},
		"config-jsonx-prefix": {
			Name:          "juan",
			UseStaticSalt: true,
			Logger:        hclog.NewNullLogger(),
			Options: map[string]string{
				"prefix": "foo",
				"format": "jsonx",
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

			cfg, err := newFormatterConfig(&testHeaderFormatter{}, tc.Options)
			require.NoError(t, err)
			f, err := newEntryFormatter(tc.Name, cfg, ss, tc.Logger)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, f)
			default:
				require.NoError(t, err)
				require.NotNil(t, f)
				require.Equal(t, tc.ExpectedFormat, f.config.requiredFormat)
				require.Equal(t, tc.ExpectedPrefix, f.config.prefix)
			}
		})
	}
}

// TestEntryFormatter_Reopen ensures that we do not get an error when calling Reopen.
func TestEntryFormatter_Reopen(t *testing.T) {
	t.Parallel()

	ss := newStaticSalt(t)
	cfg, err := newFormatterConfig(&testHeaderFormatter{}, nil)
	require.NoError(t, err)

	f, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
	require.NoError(t, err)
	require.NotNil(t, f)
	require.NoError(t, f.Reopen())
}

// TestEntryFormatter_Type ensures that the node is a 'formatter' type.
func TestEntryFormatter_Type(t *testing.T) {
	t.Parallel()

	ss := newStaticSalt(t)
	cfg, err := newFormatterConfig(&testHeaderFormatter{}, nil)
	require.NoError(t, err)

	f, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
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
			ExpectedErrorMessage: "cannot audit event (request) with no data: invalid internal parameter",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 nil,
		},
		"json-response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot audit event (response) with no data: invalid internal parameter",
			Subtype:              ResponseType,
			RequiredFormat:       JSONFormat,
			Data:                 nil,
		},
		"json-request-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"json-response-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"json-request-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to parse request from audit event: no namespace",
			Subtype:              RequestType,
			RequiredFormat:       JSONFormat,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"json-response-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to parse response from audit event: no namespace",
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
			ExpectedErrorMessage: "cannot audit event (request) with no data: invalid internal parameter",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 nil,
		},
		"jsonx-response-no-data": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot audit event (response) with no data: invalid internal parameter",
			Subtype:              ResponseType,
			RequiredFormat:       JSONxFormat,
			Data:                 nil,
		},
		"jsonx-request-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to parse request from audit event: request to request-audit a nil request",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"jsonx-response-basic-input": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to parse response from audit event: request to response-audit a nil request",
			Subtype:              ResponseType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Type: "magic"},
		},
		"jsonx-request-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to parse request from audit event: no namespace",
			Subtype:              RequestType,
			RequiredFormat:       JSONxFormat,
			Data:                 &logical.LogInput{Request: &logical.Request{ID: "123"}},
		},
		"jsonx-response-basic-input-and-request-no-ns": {
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to parse response from audit event: no namespace",
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
			cfg, err := newFormatterConfig(&testHeaderFormatter{}, map[string]string{"format": tc.RequiredFormat.String()})
			require.NoError(t, err)

			f, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
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

// BenchmarkAuditFileSink_Process benchmarks the entryFormatter and then event.FileSink calling Process.
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
	cfg, err := newFormatterConfig(&testHeaderFormatter{}, nil)
	require.NoError(b, err)
	ss := newStaticSalt(b)
	formatter, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
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

// TestEntryFormatter_FormatRequest exercises entryFormatter.formatRequest with
// varying inputs.
func TestEntryFormatter_FormatRequest(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Input                *logical.LogInput
		ShouldOmitTime       bool
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
		"omit-time": {
			Input:          &logical.LogInput{Request: &logical.Request{ID: "123"}},
			ShouldOmitTime: true,
			RootNamespace:  true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := newStaticSalt(t)
			cfg, err := newFormatterConfig(&testHeaderFormatter{}, nil)
			cfg.omitTime = tc.ShouldOmitTime
			require.NoError(t, err)
			f, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
			require.NoError(t, err)

			var ctx context.Context
			switch {
			case tc.RootNamespace:
				ctx = namespace.RootContext(context.Background())
			default:
				ctx = context.Background()
			}

			entry, err := f.formatRequest(ctx, tc.Input, &testTimeProvider{})

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, entry)
			case tc.ShouldOmitTime:
				require.NoError(t, err)
				require.NotNil(t, entry)
				require.Zero(t, entry.Time)
			default:
				require.NoError(t, err)
				require.NotNil(t, entry)
				require.NotZero(t, entry.Time)
				require.Equal(t, "2024-03-22T10:00:05.00000001Z", entry.Time)
			}
		})
	}
}

// TestEntryFormatter_FormatResponse exercises entryFormatter.formatResponse with
// varying inputs.
func TestEntryFormatter_FormatResponse(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Input                *logical.LogInput
		ShouldOmitTime       bool
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
		"omit-time": {
			Input:           &logical.LogInput{Request: &logical.Request{ID: "123"}},
			ShouldOmitTime:  true,
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
			cfg, err := newFormatterConfig(&testHeaderFormatter{}, nil)
			cfg.omitTime = tc.ShouldOmitTime
			require.NoError(t, err)
			f, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
			require.NoError(t, err)

			var ctx context.Context
			switch {
			case tc.RootNamespace:
				ctx = namespace.RootContext(context.Background())
			default:
				ctx = context.Background()
			}

			entry, err := f.formatResponse(ctx, tc.Input, &testTimeProvider{})

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
				require.Nil(t, entry)
			case tc.ShouldOmitTime:
				require.NoError(t, err)
				require.NotNil(t, entry)
				require.Zero(t, entry.Time)
			default:
				require.NoError(t, err)
				require.NotNil(t, entry)
				require.NotZero(t, entry.Time)
				require.Equal(t, "2024-03-22T10:00:05.00000001Z", entry.Time)
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
		cfg, err := newFormatterConfig(&testHeaderFormatter{}, map[string]string{
			"hmac_accessor": "false",
			"prefix":        tc.Prefix,
		})
		require.NoError(t, err)
		formatter, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
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
			Type:      event.AuditType.AsEventType(),
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
		cfg, err := newFormatterConfig(
			&testHeaderFormatter{},
			map[string]string{
				"format":        "jsonx",
				"hmac_accessor": "false",
				"prefix":        tc.Prefix,
			})
		cfg.omitTime = true
		require.NoError(t, err)
		formatter, err := newEntryFormatter("juan", cfg, tempStaticSalt, hclog.NewNullLogger())
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
			Type:      event.AuditType.AsEventType(),
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
	var formatter *entryFormatter
	var err error

	format := func(t *testing.T, config formatterConfig, operation logical.Operation, inputData map[string]any) *ResponseEntry {
		formatter, err = newEntryFormatter("juan", config, ss, hclog.NewNullLogger())
		require.NoError(t, err)
		require.NotNil(t, formatter)

		in := &logical.LogInput{
			Request:  &logical.Request{Operation: operation},
			Response: &logical.Response{Data: inputData},
		}

		resp, err := formatter.formatResponse(ctx, in, &testTimeProvider{})
		require.NoError(t, err)

		return resp
	}

	t.Run("Default case", func(t *testing.T) {
		config, err := newFormatterConfig(&testHeaderFormatter{}, map[string]string{"elide_list_responses": "true"})
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
		config, err := newFormatterConfig(&testHeaderFormatter{}, map[string]string{"elide_list_responses": "true"})
		require.NoError(t, err)
		tc := oneInterestingTestCase
		entry := format(t, config, logical.ReadOperation, tc.inputData)
		assert.Equal(t, formatter.hashExpectedValueForComparison(tc.inputData), entry.Response.Data)
	})

	t.Run("When elideListResponses is false, eliding does not happen", func(t *testing.T) {
		config, err := newFormatterConfig(&testHeaderFormatter{}, map[string]string{
			"elide_list_responses": "false",
			"format":               "json",
		})
		require.NoError(t, err)
		tc := oneInterestingTestCase
		entry := format(t, config, logical.ListOperation, tc.inputData)
		assert.Equal(t, formatter.hashExpectedValueForComparison(tc.inputData), entry.Response.Data)
	})

	t.Run("When raw is true, eliding still happens", func(t *testing.T) {
		config, err := newFormatterConfig(&testHeaderFormatter{}, map[string]string{
			"elide_list_responses": "true",
			"format":               "json",
			"log_raw":              "true",
		})
		require.NoError(t, err)
		tc := oneInterestingTestCase
		entry := format(t, config, logical.ListOperation, tc.inputData)
		assert.Equal(t, tc.expectedData, entry.Response.Data)
	})
}

// TestEntryFormatter_Process_NoMutation tests that the event returned by an
// entryFormatter.Process method is not the same as the one that it accepted.
func TestEntryFormatter_Process_NoMutation(t *testing.T) {
	t.Parallel()

	// Create the formatter node.
	cfg, err := newFormatterConfig(&testHeaderFormatter{}, nil)
	require.NoError(t, err)
	ss := newStaticSalt(t)
	formatter, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
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

// TestEntryFormatter_Process_Panic tries to send data into the entryFormatter
// which will currently cause a panic when a response is formatted due to the
// underlying hashing that is done with reflectwalk.
func TestEntryFormatter_Process_Panic(t *testing.T) {
	t.Parallel()

	// Create the formatter node.
	cfg, err := newFormatterConfig(&testHeaderFormatter{}, nil)
	require.NoError(t, err)
	ss := newStaticSalt(t)
	formatter, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
	require.NoError(t, err)
	require.NotNil(t, formatter)

	// The secret sauce, create a bad addr.
	// see: https://github.com/hashicorp/vault/issues/16462
	badAddr, err := sockaddr.NewSockAddr("10.10.10.2/32 10.10.10.3/32")
	require.NoError(t, err)

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
			Data: map[string]interface{}{},
		},
		Response: &logical.Response{
			Data: map[string]any{
				"token_bound_cidrs": []*sockaddr.SockAddrMarshaler{
					{SockAddr: badAddr},
				},
			},
		},
	}

	e := fakeEvent(t, ResponseType, in)

	e2, err := formatter.Process(namespace.RootContext(nil), e)
	require.Error(t, err)
	require.Contains(t, err.Error(), "panic generating audit log: \"juan\"")
	require.Nil(t, e2)
}

// TestEntryFormatter_NewFormatterConfig_NilHeaderFormatter ensures we cannot
// create a formatterConfig using NewFormatterConfig if we supply a nil formatter.
func TestEntryFormatter_NewFormatterConfig_NilHeaderFormatter(t *testing.T) {
	_, err := newFormatterConfig(nil, nil)
	require.Error(t, err)
}

// TestEntryFormatter_Process_NeverLeaksHeaders ensures that if we never accidentally
// leak headers if applying them means we don't have any. This is more like a sense
// check to ensure the returned event doesn't somehow end up with the headers 'back'.
func TestEntryFormatter_Process_NeverLeaksHeaders(t *testing.T) {
	t.Parallel()

	// Create the formatter node.
	cfg, err := newFormatterConfig(&testHeaderFormatter{shouldReturnEmpty: true}, nil)
	require.NoError(t, err)
	ss := newStaticSalt(t)
	formatter, err := newEntryFormatter("juan", cfg, ss, hclog.NewNullLogger())
	require.NoError(t, err)
	require.NotNil(t, formatter)

	// Set up the input and verify we have a single foo:bar header.
	var input *logical.LogInput
	err = json.Unmarshal([]byte(testFormatJSONReqBasicStrFmt), &input)
	require.NoError(t, err)
	require.NotNil(t, input)
	require.ElementsMatch(t, input.Request.Headers["foo"], []string{"bar"})

	e := fakeEvent(t, RequestType, input)

	// Process the node.
	ctx := namespace.RootContext(context.Background())
	e2, err := formatter.Process(ctx, e)
	require.NoError(t, err)
	require.NotNil(t, e2)

	// Now check we can retrieve the formatted JSON.
	jsonFormatted, b2 := e2.Format(JSONFormat.String())
	require.True(t, b2)
	require.NotNil(t, jsonFormatted)
	var input2 *logical.LogInput
	err = json.Unmarshal(jsonFormatted, &input2)
	require.NoError(t, err)
	require.NotNil(t, input2)
	require.Len(t, input2.Request.Headers, 0)
}

// hashExpectedValueForComparison replicates enough of the audit HMAC process on a piece of expected data in a test,
// so that we can use assert.Equal to compare the expected and output values.
func (f *entryFormatter) hashExpectedValueForComparison(input map[string]any) map[string]any {
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
