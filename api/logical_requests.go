// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"maps"
	"net/http"
	"net/url"
)

var _ LogicalRequest = (*defaultLogicalRequest)(nil)

// NewLogicalReadRequest creates a new LogicalReadRequest with the given path,
// values, and headers.
func NewLogicalReadRequest(path string, values url.Values, headers http.Header) LogicalReadRequest {
	return newLogicalRequest(path, values, nil, headers)
}

// NewLogicalWriteRequest creates a new LogicalWriteRequest with the given path,
// data, and headers.
func NewLogicalWriteRequest(path string, data map[string]interface{}, headers http.Header) LogicalWriteRequest {
	return newLogicalRequest(path, nil, data, headers)
}

// NewDeleteRequest creates a new LogicalDeleteRequest with the given path and values.
func NewDeleteRequest(path string, values url.Values, headers http.Header) LogicalDeleteRequest {
	return newLogicalRequest(path, values, nil, headers)
}

// newLogicalRequest creates a new LogicalRequest with the given path, values,
// data, and headers.
func newLogicalRequest(path string, values url.Values, data map[string]interface{}, headers http.Header) LogicalRequest {
	return &defaultLogicalRequest{
		path:    path,
		values:  values,
		data:    data,
		headers: headers,
	}
}

// BaseLogicalRequest is the interface for requests to Vault's logical backend
// that do not include data or values.
type BaseLogicalRequest interface {
	// Path returns the path to write to in Vault, without the "/v1/" prefix.
	Path() string
	// Headers returns the headers to be included in the request to Vault. All
	// headers are additive, and must not collide with any of the reserved headers.
	Headers() http.Header
}

// LogicalRequest is the interface for requests to Vault's logical backend.
type LogicalRequest interface {
	BaseLogicalRequest
	// Values returns the query parameters to be used in the request.
	// Values are only used in read and delete requests.
	Values() url.Values
	// Data returns the data to be written to the path. It is marshaled to JSON.
	// Data is only used in write requests.
	Data() map[string]interface{}
}

// LogicalWriteRequest is the interface for requests that write data to Vault's
// logical backend.
type LogicalWriteRequest interface {
	BaseLogicalRequest
	Data() map[string]interface{}
}

// LogicalReadRequest is the interface for requests that read data from Vault's
// logical backend.
type LogicalReadRequest interface {
	BaseLogicalRequest
	Values() url.Values
}

// LogicalDeleteRequest is the interface for requests that delete data from Vault's
// logical backend. It is semantically similar the same as a read request,
type LogicalDeleteRequest interface {
	BaseLogicalRequest
	Values() url.Values
}

// defaultLogicalRequest is the default implementation of LogicalRequest.
type defaultLogicalRequest struct {
	path    string
	values  url.Values
	headers http.Header
	data    map[string]interface{}
}

// Path returns the path to write to in Vault, without the "/v1/" prefix.
func (r *defaultLogicalRequest) Path() string {
	return r.path
}

// Headers returns a copy of the headers to be included in the request to Vault.
func (r *defaultLogicalRequest) Headers() http.Header {
	if r.headers == nil {
		return nil
	}
	return maps.Clone(r.headers)
}

// Data returns a copy of the data to be written to the path.
func (r *defaultLogicalRequest) Data() map[string]interface{} {
	if r.data == nil {
		return nil
	}
	return maps.Clone(r.data)
}

// Values returns a copy of the query parameters to be used in the request.
func (r *defaultLogicalRequest) Values() url.Values {
	if r.values == nil {
		return nil
	}
	return maps.Clone(r.values)
}
