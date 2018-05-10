package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
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
	req, err := r.toRetryableHTTP(true)
	if err != nil {
		return nil, err
	}
	return req.Request, nil
}

// legacy indicates whether we want to return a request derived from
// http.NewRequest instead of retryablehttp.NewRequest, so that legacy clents
// that might be using the public ToHTTP method still work
func (r *Request) toRetryableHTTP(legacy bool) (*retryablehttp.Request, error) {
	// Encode the query parameters
	r.URL.RawQuery = r.Params.Encode()

	// Create the HTTP request, defaulting to retryable
	var req *retryablehttp.Request

	if legacy {
		regReq, err := http.NewRequest(r.Method, r.URL.RequestURI(), r.Body)
		if err != nil {
			return nil, err
		}
		req = &retryablehttp.Request{
			Request: regReq,
		}
	} else {
		var buf []byte
		var err error
		if r.Body != nil {
			buf, err = ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
		}
		req, err = retryablehttp.NewRequest(r.Method, r.URL.RequestURI(), bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
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

	return req, nil
}
