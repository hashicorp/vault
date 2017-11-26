package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Request is a raw request configuration structure used to initiate
// API requests to the Vault server.
type Request struct {
	Method        string
	URL           *url.URL
	Params        url.Values
	Headers       http.Header
	ClientToken   string
	MFAHeaderVals []string
	WrapTTL       string
	Obj           interface{}
	Body          io.Reader
	BodySize      int64

	// Whether to request overriding soft-mandatory Sentinel policies (RGPs and
	// EGPs). If set, the override flag will take effect for all policies
	// evaluated during the request.
	PolicyOverride bool

	// ctx is a context to use with this request. It should only be modified using
	// WithContext.
	ctx context.Context
}

// WithContext creates and returns a shallow copy of the request with the given
// context.
func (r *Request) WithContext(ctx context.Context) *Request {
	// Forbid a nil context. The net/http package actually panics here, but that
	// seems like a bad behavior.
	if ctx == nil {
		ctx = context.Background()
	}

	// Copy the request
	r2 := new(Request)
	if r != nil {
		*r2 = *r
	}
	r2.ctx = ctx

	// Deep copy URL because it is mutable by users of WithContext.
	if r.URL != nil {
		r2URL := new(url.URL)
		*r2URL = *r.URL
		r2.URL = r2URL
	}

	return r2
}

// Context returns the context attached to this request, or context.Background()
// if none exists.
func (r *Request) Context() context.Context {
	if r != nil && r.ctx != nil {
		return r.ctx
	}
	return context.Background()
}

// SetJSONBody is used to set a request body that is a JSON-encoded value.
func (r *Request) SetJSONBody(val interface{}) error {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(val); err != nil {
		return err
	}

	r.Obj = val
	r.Body = buf
	r.BodySize = int64(buf.Len())
	return nil
}

// ResetJSONBody is used to reset the body for a redirect
func (r *Request) ResetJSONBody() error {
	if r.Body == nil {
		return nil
	}
	return r.SetJSONBody(r.Obj)
}

// ToHTTP turns this request into a valid *http.Request for use with the
// net/http package.
func (r *Request) ToHTTP() (*http.Request, error) {
	// Encode the query parameters
	r.URL.RawQuery = r.Params.Encode()

	// Create the HTTP request
	req, err := http.NewRequest(r.Method, r.URL.RequestURI(), r.Body)
	if err != nil {
		return nil, err
	}

	req.URL.User = r.URL.User
	req.URL.Scheme = r.URL.Scheme
	req.URL.Host = r.URL.Host
	req.Host = r.URL.Host

	if r.Headers != nil {
		for header, vals := range r.Headers {
			for _, val := range vals {
				req.Header.Add(header, val)
			}
		}
	}

	if len(r.ClientToken) != 0 {
		req.Header.Set("X-Vault-Token", r.ClientToken)
	}

	if len(r.WrapTTL) != 0 {
		req.Header.Set("X-Vault-Wrap-TTL", r.WrapTTL)
	}

	if len(r.MFAHeaderVals) != 0 {
		for _, mfaHeaderVal := range r.MFAHeaderVals {
			req.Header.Add("X-Vault-MFA", mfaHeaderVal)
		}
	}

	if r.PolicyOverride {
		req.Header.Set("X-Vault-Policy-Override", "true")
	}

	// Add the context
	req = req.WithContext(r.Context())

	return req, nil
}
