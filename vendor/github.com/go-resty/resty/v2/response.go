// Copyright (c) 2015-2024 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Response struct and methods
//_______________________________________________________________________

// Response struct holds response values of executed requests.
type Response struct {
	Request     *Request
	RawResponse *http.Response

	body       []byte
	size       int64
	receivedAt time.Time
}

// Body method returns the HTTP response as `[]byte` slice for the executed request.
//
// NOTE: [Response.Body] might be nil if [Request.SetOutput] is used.
// Also see [Request.SetDoNotParseResponse], [Client.SetDoNotParseResponse]
func (r *Response) Body() []byte {
	if r.RawResponse == nil {
		return []byte{}
	}
	return r.body
}

// SetBody method sets [Response] body in byte slice. Typically,
// It is helpful for test cases.
//
//	resp.SetBody([]byte("This is test body content"))
//	resp.SetBody(nil)
func (r *Response) SetBody(b []byte) *Response {
	r.body = b
	return r
}

// Status method returns the HTTP status string for the executed request.
//
//	Example: 200 OK
func (r *Response) Status() string {
	if r.RawResponse == nil {
		return ""
	}
	return r.RawResponse.Status
}

// StatusCode method returns the HTTP status code for the executed request.
//
//	Example: 200
func (r *Response) StatusCode() int {
	if r.RawResponse == nil {
		return 0
	}
	return r.RawResponse.StatusCode
}

// Proto method returns the HTTP response protocol used for the request.
func (r *Response) Proto() string {
	if r.RawResponse == nil {
		return ""
	}
	return r.RawResponse.Proto
}

// Result method returns the response value as an object if it has one
//
// See [Request.SetResult]
func (r *Response) Result() interface{} {
	return r.Request.Result
}

// Error method returns the error object if it has one
//
// See [Request.SetError], [Client.SetError]
func (r *Response) Error() interface{} {
	return r.Request.Error
}

// Header method returns the response headers
func (r *Response) Header() http.Header {
	if r.RawResponse == nil {
		return http.Header{}
	}
	return r.RawResponse.Header
}

// Cookies method to returns all the response cookies
func (r *Response) Cookies() []*http.Cookie {
	if r.RawResponse == nil {
		return make([]*http.Cookie, 0)
	}
	return r.RawResponse.Cookies()
}

// String method returns the body of the HTTP response as a `string`.
// It returns an empty string if it is nil or the body is zero length.
func (r *Response) String() string {
	if len(r.body) == 0 {
		return ""
	}
	return strings.TrimSpace(string(r.body))
}

// Time method returns the duration of HTTP response time from the request we sent
// and received a request.
//
// See [Response.ReceivedAt] to know when the client received a response and see
// `Response.Request.Time` to know when the client sent a request.
func (r *Response) Time() time.Duration {
	if r.Request.clientTrace != nil {
		return r.Request.TraceInfo().TotalTime
	}
	return r.receivedAt.Sub(r.Request.Time)
}

// ReceivedAt method returns the time we received a response from the server for the request.
func (r *Response) ReceivedAt() time.Time {
	return r.receivedAt
}

// Size method returns the HTTP response size in bytes. Yeah, you can rely on HTTP `Content-Length`
// header, however it won't be available for chucked transfer/compressed response.
// Since Resty captures response size details when processing the response body
// when possible. So that users get the actual size of response bytes.
func (r *Response) Size() int64 {
	return r.size
}

// RawBody method exposes the HTTP raw response body. Use this method in conjunction with
// [Client.SetDoNotParseResponse] or [Request.SetDoNotParseResponse]
// option; otherwise, you get an error as `read err: http: read on closed response body.`
//
// Do not forget to close the body, otherwise you might get into connection leaks, no connection reuse.
// You have taken over the control of response parsing from Resty.
func (r *Response) RawBody() io.ReadCloser {
	if r.RawResponse == nil {
		return nil
	}
	return r.RawResponse.Body
}

// IsSuccess method returns true if HTTP status `code >= 200 and <= 299` otherwise false.
func (r *Response) IsSuccess() bool {
	return r.StatusCode() > 199 && r.StatusCode() < 300
}

// IsError method returns true if HTTP status `code >= 400` otherwise false.
func (r *Response) IsError() bool {
	return r.StatusCode() > 399
}

func (r *Response) setReceivedAt() {
	r.receivedAt = time.Now()
	if r.Request.clientTrace != nil {
		r.Request.clientTrace.endTime = r.receivedAt
	}
}

func (r *Response) fmtBodyString(sl int64) string {
	if r.Request.client.notParseResponse || r.Request.notParseResponse {
		return "***** DO NOT PARSE RESPONSE - Enabled *****"
	}
	if len(r.body) > 0 {
		if int64(len(r.body)) > sl {
			return fmt.Sprintf("***** RESPONSE TOO LARGE (size - %d) *****", len(r.body))
		}
		ct := r.Header().Get(hdrContentTypeKey)
		if IsJSONType(ct) {
			out := acquireBuffer()
			defer releaseBuffer(out)
			err := json.Indent(out, r.body, "", "   ")
			if err != nil {
				return fmt.Sprintf("*** Error: Unable to format response body - \"%s\" ***\n\nLog Body as-is:\n%s", err, r.String())
			}
			return out.String()
		}
		return r.String()
	}

	return "***** NO CONTENT *****"
}
