package logical

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/hashicorp/vault/sdk/helper/wrapping"
)

const (
	// HTTPContentType can be specified in the Data field of a Response
	// so that the HTTP front end can specify a custom Content-Type associated
	// with the HTTPRawBody. This can only be used for non-secrets, and should
	// be avoided unless absolutely necessary, such as implementing a specification.
	// The value must be a string.
	HTTPContentType = "http_content_type"

	// HTTPRawBody is the raw content of the HTTP body that goes with the HTTPContentType.
	// This can only be specified for non-secrets, and should should be similarly
	// avoided like the HTTPContentType. The value must be a byte slice.
	HTTPRawBody = "http_raw_body"

	// HTTPStatusCode is the response code of the HTTP body that goes with the HTTPContentType.
	// This can only be specified for non-secrets, and should should be similarly
	// avoided like the HTTPContentType. The value must be an integer.
	HTTPStatusCode = "http_status_code"

	// For unwrapping we may need to know whether the value contained in the
	// raw body is already JSON-unmarshaled. The presence of this key indicates
	// that it has already been unmarshaled. That way we don't need to simply
	// ignore errors.
	HTTPRawBodyAlreadyJSONDecoded = "http_raw_body_already_json_decoded"

	// If set, HTTPCacheControlHeader will replace the default Cache-Control=no-store header
	// set by the generic wrapping handler. The value must be a string.
	HTTPCacheControlHeader = "http_raw_cache_control"

	// If set, HTTPPragmaHeader will set the Pragma response header.
	// The value must be a string.
	HTTPPragmaHeader = "http_raw_pragma"

	// If set, HTTPWWWAuthenticateHeader will set the WWW-Authenticate response header.
	// The value must be a string.
	HTTPWWWAuthenticateHeader = "http_www_authenticate"
)

// Response is a struct that stores the response of a request.
// It is used to abstract the details of the higher level request protocol.
type Response struct {
	// Secret, if not nil, denotes that this response represents a secret.
	Secret *Secret `json:"secret" structs:"secret" mapstructure:"secret"`

	// Auth, if not nil, contains the authentication information for
	// this response. This is only checked and means something for
	// credential backends.
	Auth *Auth `json:"auth" structs:"auth" mapstructure:"auth"`

	// Response data is an opaque map that must have string keys. For
	// secrets, this data is sent down to the user as-is. To store internal
	// data that you don't want the user to see, store it in
	// Secret.InternalData.
	Data map[string]interface{} `json:"data" structs:"data" mapstructure:"data"`

	// Redirect is an HTTP URL to redirect to for further authentication.
	// This is only valid for credential backends. This will be blanked
	// for any logical backend and ignored.
	Redirect string `json:"redirect" structs:"redirect" mapstructure:"redirect"`

	// Warnings allow operations or backends to return warnings in response
	// to user actions without failing the action outright.
	Warnings []string `json:"warnings" structs:"warnings" mapstructure:"warnings"`

	// Information for wrapping the response in a cubbyhole
	WrapInfo *wrapping.ResponseWrapInfo `json:"wrap_info" structs:"wrap_info" mapstructure:"wrap_info"`

	// Headers will contain the http headers from the plugin that it wishes to
	// have as part of the output
	Headers map[string][]string `json:"headers" structs:"headers" mapstructure:"headers"`
}

// AddWarning adds a warning into the response's warning list
func (r *Response) AddWarning(warning string) {
	if r.Warnings == nil {
		r.Warnings = make([]string, 0, 1)
	}
	r.Warnings = append(r.Warnings, warning)
}

// IsError returns true if this response seems to indicate an error.
func (r *Response) IsError() bool {
	return r != nil && r.Data != nil && len(r.Data) == 1 && r.Data["error"] != nil
}

func (r *Response) Error() error {
	if !r.IsError() {
		return nil
	}
	switch r.Data["error"].(type) {
	case string:
		return errors.New(r.Data["error"].(string))
	case error:
		return r.Data["error"].(error)
	}
	return nil
}

// HelpResponse is used to format a help response
func HelpResponse(text string, seeAlso []string, oapiDoc interface{}) *Response {
	return &Response{
		Data: map[string]interface{}{
			"help":     text,
			"see_also": seeAlso,
			"openapi":  oapiDoc,
		},
	}
}

// ErrorResponse is used to format an error response
func ErrorResponse(text string, vargs ...interface{}) *Response {
	if len(vargs) > 0 {
		text = fmt.Sprintf(text, vargs...)
	}
	return &Response{
		Data: map[string]interface{}{
			"error": text,
		},
	}
}

// ListResponse is used to format a response to a list operation.
func ListResponse(keys []string) *Response {
	resp := &Response{
		Data: map[string]interface{}{},
	}
	if len(keys) != 0 {
		resp.Data["keys"] = keys
	}
	return resp
}

// ListResponseWithInfo is used to format a response to a list operation and
// return the keys as well as a map with corresponding key info.
func ListResponseWithInfo(keys []string, keyInfo map[string]interface{}) *Response {
	resp := ListResponse(keys)

	keyInfoData := make(map[string]interface{})
	for _, key := range keys {
		val, ok := keyInfo[key]
		if ok {
			keyInfoData[key] = val
		}
	}

	if len(keyInfoData) > 0 {
		resp.Data["key_info"] = keyInfoData
	}

	return resp
}

// RespondWithStatusCode takes a response and converts it to a raw response with
// the provided Status Code.
func RespondWithStatusCode(resp *Response, req *Request, code int) (*Response, error) {
	ret := &Response{
		Data: map[string]interface{}{
			HTTPContentType: "application/json",
			HTTPStatusCode:  code,
		},
	}

	if resp != nil {
		httpResp := LogicalResponseToHTTPResponse(resp)

		if req != nil {
			httpResp.RequestID = req.ID
		}

		body, err := json.Marshal(httpResp)
		if err != nil {
			return nil, err
		}

		// We default to string here so that the value is HMAC'd via audit.
		// Since this function is always marshaling to JSON, this is
		// appropriate.
		ret.Data[HTTPRawBody] = string(body)
	}

	return ret, nil
}

// HTTPResponseWriter is optionally added to a request object and can be used to
// write directly to the HTTP response writer.
type HTTPResponseWriter struct {
	http.ResponseWriter
	written *uint32
}

// NewHTTPResponseWriter creates a new HTTPResponseWriter object that wraps the
// provided io.Writer.
func NewHTTPResponseWriter(w http.ResponseWriter) *HTTPResponseWriter {
	return &HTTPResponseWriter{
		ResponseWriter: w,
		written:        new(uint32),
	}
}

// Write will write the bytes to the underlying io.Writer.
func (w *HTTPResponseWriter) Write(bytes []byte) (int, error) {
	atomic.StoreUint32(w.written, 1)
	return w.ResponseWriter.Write(bytes)
}

// Written tells us if the writer has been written to yet.
func (w *HTTPResponseWriter) Written() bool {
	return atomic.LoadUint32(w.written) == 1
}

type WrappingResponseWriter interface {
	http.ResponseWriter
	Wrapped() http.ResponseWriter
}

type StatusHeaderResponseWriter struct {
	wrapped     http.ResponseWriter
	wroteHeader bool
	StatusCode  int
	headers     map[string][]*CustomHeader
}

func NewStatusHeaderResponseWriter(w http.ResponseWriter, h map[string][]*CustomHeader) *StatusHeaderResponseWriter {
	return &StatusHeaderResponseWriter{
		wrapped:     w,
		wroteHeader: false,
		StatusCode:  200,
		headers:     h,
	}
}

func (w *StatusHeaderResponseWriter) Wrapped() http.ResponseWriter {
	return w.wrapped
}

func (w *StatusHeaderResponseWriter) Header() http.Header {
	return w.wrapped.Header()
}

func (w *StatusHeaderResponseWriter) Write(buf []byte) (int, error) {
	// It is allowed to only call ResponseWriter.Write and skip
	// ResponseWriter.WriteHeader. An example of such a situation is
	// "handleUIStub". The Write function will internally set the status code
	// 200 for the response for which that call might invoke other
	// implementations of the WriteHeader function. So, we still need to set
	// the custom headers. In cases where both WriteHeader and Write of
	// statusHeaderResponseWriter struct are called the internal call to the
	// WriterHeader invoked from inside Write method won't change the headers.
	if !w.wroteHeader {
		w.setCustomResponseHeaders(w.StatusCode)
	}

	return w.wrapped.Write(buf)
}

func (w *StatusHeaderResponseWriter) WriteHeader(statusCode int) {
	w.setCustomResponseHeaders(statusCode)
	w.wrapped.WriteHeader(statusCode)
	w.StatusCode = statusCode
	// in cases where Write is called after WriteHeader, let's prevent setting
	// ResponseWriter headers twice
	w.wroteHeader = true
}

func (w *StatusHeaderResponseWriter) setCustomResponseHeaders(status int) {
	sch := w.headers
	if sch == nil {
		return
	}

	// Checking the validity of the status code
	if status >= 600 || status < 100 {
		return
	}

	// setter function to set the headers
	setter := func(hvl []*CustomHeader) {
		for _, hv := range hvl {
			w.Header().Set(hv.Name, hv.Value)
		}
	}

	// Setting the default headers first
	setter(sch["default"])

	// setting the Xyy pattern first
	d := fmt.Sprintf("%vxx", status/100)
	if val, ok := sch[d]; ok {
		setter(val)
	}

	// Setting the specific headers
	if val, ok := sch[strconv.Itoa(status)]; ok {
		setter(val)
	}

	return
}

var _ WrappingResponseWriter = &StatusHeaderResponseWriter{}

// ResolveRoleResponse returns a standard response to be returned by functions handling a ResolveRoleOperation
func ResolveRoleResponse(roleName string) (*Response, error) {
	return &Response{
		Data: map[string]interface{}{
			"role": roleName,
		},
	}, nil
}
