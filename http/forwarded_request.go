package http

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/hashicorp/vault/helper/compressutil"
	"github.com/hashicorp/vault/helper/jsonutil"
)

type bufCloser struct {
	*bytes.Buffer
}

func (b bufCloser) Close() error {
	b.Reset()
	return nil
}

type forwardedRequest struct {
	// The original method
	Method string `json:"method"`

	// The original path
	RawPath string `json:"raw_path"`

	// The original query string
	RawQuery string `json:"raw_query"`

	// The client token header value
	ClientToken string `json:"client_token"`

	// The wrap TTL header value
	WrapTTL string `json:"wrap_ttl"`

	// The request body
	Body []byte `json:"body"`

	// The specified host

	Host string `json:"host"`
	// The remote address
	RemoteAddr string `json:"remote_addr"`

	// The client's TLS connection state
	ConnectionState *tls.ConnectionState `json:"connection_state"`
}

func generateForwardedRequest(req *http.Request, addr string) (*http.Request, error) {
	fq := forwardedRequest{
		Method:          req.Method,
		RawPath:         req.URL.RawPath,
		RawQuery:        req.URL.RawQuery,
		ClientToken:     req.Header.Get(AuthHeaderName),
		WrapTTL:         req.Header.Get(WrapTTLHeaderName),
		Host:            req.Host,
		RemoteAddr:      req.RemoteAddr,
		ConnectionState: req.TLS,
	}

	buf := bytes.NewBuffer(nil)
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return nil, err
	}
	fq.Body = buf.Bytes()

	newBody, err := jsonutil.EncodeJSONAndCompress(&fq, &compressutil.CompressionConfig{
		Type: compressutil.CompressionTypeLzw,
	})
	if err != nil {
		return nil, err
	}

	ret, err := http.NewRequest("POST", addr, bytes.NewBuffer(newBody))
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// parseForwardedRequest generates a new http.Request that is comprised of the
// values in the given request's body, assuming it correctly parses into a
// forwardedRequest.
func parseForwardedRequest(req *http.Request) (*http.Request, error) {
	buf := bufCloser{
		Buffer: bytes.NewBuffer(nil),
	}
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return nil, err
	}

	var fq forwardedRequest
	err = jsonutil.DecodeJSON(buf.Bytes(), &fq)
	if err != nil {
		return nil, err
	}

	buf.Reset()
	_, err = buf.Write(fq.Body)
	if err != nil {
		return nil, err
	}

	ret := &http.Request{
		Method: fq.Method,
		URL: &url.URL{
			RawPath:  fq.RawPath,
			RawQuery: fq.RawQuery,
		},
		Header: map[string][]string{
			AuthHeaderName: {fq.ClientToken},
		},
		Body:       buf,
		Host:       fq.Host,
		RemoteAddr: fq.RemoteAddr,
		TLS:        fq.ConnectionState,
	}
	if fq.WrapTTL != "" {
		ret.Header.Add(WrapTTLHeaderName, fq.WrapTTL)
	}

	return ret, nil
}
