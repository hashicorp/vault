// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newLogicalRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		values  url.Values
		data    map[string]interface{}
		headers http.Header
		want    LogicalRequest
	}{
		{
			name:    "read request",
			path:    "sys/health",
			values:  url.Values{"standbyok": []string{"true"}},
			data:    nil,
			headers: http.Header{"X-Read-Request": []string{"true"}},
			want: &defaultLogicalRequest{
				path:    "sys/health",
				values:  url.Values{"standbyok": []string{"true"}},
				headers: http.Header{"X-Read-Request": []string{"true"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, newLogicalRequest(tt.path, tt.values, tt.data, tt.headers), "newLogicalRequest(%v, %v, %v, %v)", tt.path, tt.values, tt.data, tt.headers)
		})
	}
}

func Test_defaultLogicalRequest(t *testing.T) {
	t.Parallel()

	type want struct {
		path    string
		headers http.Header
		values  url.Values
		data    map[string]interface{}
	}

	tests := []struct {
		name    string
		want    want
		request *defaultLogicalRequest
	}{
		{
			name: "all",
			want: want{
				path:    "sys/health",
				values:  url.Values{"standbyok": []string{"true"}},
				data:    map[string]interface{}{"key": "value"},
				headers: http.Header{"X-Read-Request": []string{"true"}},
			},
			request: &defaultLogicalRequest{
				path:    "sys/health",
				values:  url.Values{"standbyok": []string{"true"}},
				data:    map[string]interface{}{"key": "value"},
				headers: http.Header{"X-Read-Request": []string{"true"}},
			},
		},
		{
			name:    "empty",
			want:    want{},
			request: &defaultLogicalRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.request
			assert.Equalf(t, tt.want.path, r.Path(), "Path()")

			gotData := r.Data()
			if assert.Equalf(t, tt.want.data, gotData, "Data()") {
				if tt.want.data != nil {
					gotData["other"] = "value" // Ensure that Data() returns a copy
					assert.NotEqualf(t, gotData, r.Data(), "Data() should return a copy")
				}
			}

			gotHeaders := r.Headers()
			if assert.Equalf(t, tt.want.headers, gotHeaders, "Headers()") {
				if tt.want.headers != nil {
					gotHeaders["X-other"] = []string{"false"} // Ensure that Headers() returns a copy
					assert.NotEqualf(t, gotHeaders, r.Headers(), "Headers() should return a copy")
				}
			}

			gotValues := r.Values()
			if assert.Equalf(t, tt.want.values, gotValues, "Values()") {
				if gotValues != nil {
					gotValues.Add("other", "value") // Ensure that Values() returns a copy
					assert.NotEqualf(t, gotValues, r.Values(), "Values() should return a copy")
				}
			}
		})
	}
}

func TestNewLogicalReadRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		values  url.Values
		headers http.Header
		want    LogicalReadRequest
	}{
		{
			name: "basic read request",
			path: "sys/health",
			want: &defaultLogicalRequest{
				path: "sys/health",
			},
		},
		{
			name:   "read request with values",
			path:   "sys/health",
			values: url.Values{"standbyok": []string{"true"}},
			want: &defaultLogicalRequest{
				path:   "sys/health",
				values: url.Values{"standbyok": []string{"true"}},
			},
		},
		{
			name:    "read request with headers",
			path:    "sys/health",
			headers: http.Header{"X-Read-Request": []string{"true"}},
			want: &defaultLogicalRequest{
				path:    "sys/health",
				headers: http.Header{"X-Read-Request": []string{"true"}},
			},
		},
		{
			name:    "read request with values and headers",
			path:    "sys/health",
			headers: http.Header{"X-Read-Request": []string{"true"}},
			values:  url.Values{"standbyok": []string{"true"}},
			want: &defaultLogicalRequest{
				path:    "sys/health",
				headers: http.Header{"X-Read-Request": []string{"true"}},
				values:  url.Values{"standbyok": []string{"true"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewLogicalReadRequest(tt.path, tt.values, tt.headers), "NewLogicalReadRequest(%v, %v, %v)", tt.path, tt.values, tt.headers)
		})
	}
}

func TestNewLogicalWriteRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		data    map[string]interface{}
		headers http.Header
		want    LogicalWriteRequest
	}{
		{
			name: "basic read request",
			path: "sys/health",
			want: &defaultLogicalRequest{
				path: "sys/health",
			},
		},
		{
			name: "read request with data",
			path: "sys/health",
			data: map[string]interface{}{"key": "value"},
			want: &defaultLogicalRequest{
				path: "sys/health",
				data: map[string]interface{}{"key": "value"},
			},
		},
		{
			name:    "read request with headers",
			path:    "sys/health",
			headers: http.Header{"X-Read-Request": []string{"true"}},
			want: &defaultLogicalRequest{
				path:    "sys/health",
				headers: http.Header{"X-Read-Request": []string{"true"}},
			},
		},
		{
			name:    "read request with data and headers",
			path:    "sys/health",
			headers: http.Header{"X-Read-Request": []string{"true"}},
			data:    map[string]interface{}{"key": "value"},
			want: &defaultLogicalRequest{
				path:    "sys/health",
				headers: http.Header{"X-Read-Request": []string{"true"}},
				data:    map[string]interface{}{"key": "value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewLogicalWriteRequest(tt.path, tt.data, tt.headers), "NewLogicalWriteRequest(%v, %v, %v)", tt.path, tt.data, tt.headers)
		})
	}
}
