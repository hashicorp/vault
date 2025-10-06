// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"net/http"
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
