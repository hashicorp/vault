// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogical_addExtraHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		r           *Request
		headers     http.Header
		wantHeaders http.Header
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "nil headers",
			r: &Request{
				Headers: http.Header{
					"X-Other-Header": []string{"other-value"},
				},
			},
			headers: nil,
			wantHeaders: http.Header{
				"X-Other-Header": []string{"other-value"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "empty headers",
			r: &Request{
				Headers: http.Header{
					"X-Other-Header": []string{"other-value"},
				},
			},
			headers: http.Header{},
			wantHeaders: http.Header{
				"X-Other-Header": []string{"other-value"},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "no headers",
			r:       &Request{},
			wantErr: assert.NoError,
		},
		{
			name: "nil request",
			headers: http.Header{
				"X-Extra-Header": []string{"real-value"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "nil request", i...)
			},
		},
		{
			name: "extra headers",
			r: &Request{
				Headers: http.Header{
					"X-Other-Header": []string{"other-value"},
				},
			},
			headers: http.Header{
				"X-Extra-Header": []string{"real-value"},
			},
			wantHeaders: http.Header{
				"X-Extra-Header": []string{"real-value"},
				"X-Other-Header": []string{"other-value"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "reserved header",
			r: &Request{
				Headers: http.Header{
					"X-Reserved-Header": []string{"reserved-value"},
				},
			},
			headers: http.Header{
				"x-rESERved-hEAder": []string{"other-value"},
			},
			wantHeaders: http.Header{
				"X-Reserved-Header": []string{"reserved-value"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err,
					fmt.Sprintf("cannot set extra header %q, it is reserved", "X-Reserved-Header"), i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Logical{}
			tt.wantErr(t, c.addExtraHeaders(tt.r, tt.headers), fmt.Sprintf("addExtraHeaders(%v, %v)", tt.r, tt.headers))
			if tt.r == nil {
				if tt.wantHeaders != nil {
					require.Fail(t, "invalid test case: nil request with headers", "addExtraHeaders(%v, %v)", tt.r, tt.headers)
				}
				return
			}
			assert.Equalf(t, tt.wantHeaders, tt.r.Headers, "Headers after addExtraHeaders(%v, %v)", tt.r, tt.headers)
		})
	}
}

// TestWriteRawWithContextBodyReadable verifies that WriteRawWithContext and
// PatchRawWithContext return a *Response whose body remains readable after the
// call returns. Previously, both methods routed through writeRaw which installed
// a withConfiguredTimeout cancel and deferred it, cancelling the context (and
// thus the response body) before the caller had a chance to read it.
func TestWriteRawWithContextBodyReadable(t *testing.T) {
	t.Parallel()

	const responseBody = `{"data":"hello"}`

	for _, tc := range []struct {
		name   string
		invoke func(l *Logical, path string, data []byte) (*Response, error)
	}{
		{
			name: "WriteRawWithContext",
			invoke: func(l *Logical, path string, data []byte) (*Response, error) {
				return l.WriteRawWithContext(context.Background(), path, data)
			},
		},
		{
			name: "PatchRawWithContext",
			invoke: func(l *Logical, path string, data []byte) (*Response, error) {
				return l.PatchRawWithContext(context.Background(), path, data)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(responseBody))
			})

			config, ln := testHTTPServer(t, handler)
			defer ln.Close()

			client, err := NewClient(config)
			require.NoError(t, err)
			client.SetToken("test-token")

			resp, err := tc.invoke(client.Logical(), "secret/test", []byte(`{}`))
			require.NoError(t, err)
			require.NotNil(t, resp)
			defer resp.Body.Close()

			// Body must be readable after the call returns — this was broken
			// when writeRaw deferred cancelFunc() before the caller could read.
			got, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.True(t, strings.Contains(string(got), "hello"),
				"expected response body to contain 'hello', got: %s", string(got))
		})
	}
}
