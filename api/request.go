package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

// Request is a raw request configuration structure used to initiate
// API requests to the Vault server.
type Request struct {
	Method      string
	URL         *url.URL
	Params      url.Values
	ClientToken string
	WrapTTL     string
	Obj         interface{}
	Body        io.Reader
	BodySize    int64
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

// ToHTTP turns this request into a *retryablehttp.Request
func (r *Request) ToHTTP() (*retryablehttp.Request, error) {
	// Encode the query parameters
	r.URL.RawQuery = r.Params.Encode()

	// Create the HTTP request; retryable needs a ReadSeeker
	body := bytes.NewBuffer(nil)
	if r.Body != nil {
		n, err := body.ReadFrom(r.Body)
		if err != nil {
			return nil, err
		}
		if n != r.BodySize {
			return nil, fmt.Errorf("Could not read full body size from Request")
		}
	}
	req, err := retryablehttp.NewRequest(r.Method, r.URL.RequestURI(), bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, err
	}

	req.URL.Scheme = r.URL.Scheme
	req.URL.Host = r.URL.Host
	req.Host = r.URL.Host

	if len(r.ClientToken) != 0 {
		req.Header.Set("X-Vault-Token", r.ClientToken)
	}

	if len(r.WrapTTL) != 0 {
		req.Header.Set("X-Vault-Wrap-TTL", r.WrapTTL)
	}

	return req, nil
}
