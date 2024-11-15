// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"io"
	"net/http"
)

// Raw can be used to do raw queries against custom endpoints
type Raw struct {
	c *Client
}

// Raw returns a handle to query endpoints
func (c *Client) Raw() *Raw {
	return &Raw{c}
}

// Query is used to do a GET request against an endpoint
// and deserialize the response into an interface using
// standard Nomad conventions.
func (raw *Raw) Query(endpoint string, out interface{}, q *QueryOptions) (*QueryMeta, error) {
	return raw.c.query(endpoint, out, q)
}

// Response is used to make a GET request against an endpoint and returns the
// response body
func (raw *Raw) Response(endpoint string, q *QueryOptions) (io.ReadCloser, error) {
	return raw.c.rawQuery(endpoint, q)
}

// Write is used to do a PUT request against an endpoint
// and serialize/deserialized using the standard Nomad conventions.
func (raw *Raw) Write(endpoint string, in, out interface{}, q *WriteOptions) (*WriteMeta, error) {
	return raw.c.put(endpoint, in, out, q)
}

// Delete is used to do a DELETE request against an endpoint
// and serialize/deserialized using the standard Nomad conventions.
func (raw *Raw) Delete(endpoint string, out interface{}, q *WriteOptions) (*WriteMeta, error) {
	return raw.c.delete(endpoint, nil, out, q)
}

// Do uses the raw client's internal httpClient to process the request
func (raw *Raw) Do(req *http.Request) (*http.Response, error) {
	return raw.c.httpClient.Do(req)
}
