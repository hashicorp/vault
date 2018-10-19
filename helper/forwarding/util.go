package forwarding

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/golang/protobuf/proto"
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

// GenerateForwardedRequest generates a new http.Request that contains the
// original requests's information in the new request's body.
func GenerateForwardedHTTPRequest(req *http.Request, addr string) (*http.Request, error) {
	fq, err := GenerateForwardedRequest(req)
	if err != nil {
		return nil, err
	}

	var newBody []byte
	switch os.Getenv("VAULT_MESSAGE_TYPE") {
	case "json":
		newBody, err = jsonutil.EncodeJSON(fq)
	case "json_compress":
		newBody, err = jsonutil.EncodeJSONAndCompress(fq, &compressutil.CompressionConfig{
			Type: compressutil.CompressionTypeLZW,
		})
	case "proto3":
		fallthrough
	default:
		newBody, err = proto.Marshal(fq)
	}
	if err != nil {
		return nil, err
	}

	ret, err := http.NewRequest("POST", addr, bytes.NewBuffer(newBody))
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func GenerateForwardedRequest(req *http.Request) (*Request, error) {
	var reader io.Reader = req.Body
	ctx := req.Context()
	maxRequestSize := ctx.Value("max_request_size")
	if maxRequestSize != nil {
		max, ok := maxRequestSize.(int64)
		if !ok {
			return nil, errors.New("could not parse max_request_size from request context")
		}
		if max > 0 {
			reader = io.LimitReader(req.Body, max)
		}
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	fq := Request{
		Method:        req.Method,
		HeaderEntries: make(map[string]*HeaderEntry, len(req.Header)),
		Host:          req.Host,
		RemoteAddr:    req.RemoteAddr,
		Body:          body,
	}

	reqURL := req.URL
	fq.Url = &URL{
		Scheme:   reqURL.Scheme,
		Opaque:   reqURL.Opaque,
		Host:     reqURL.Host,
		Path:     reqURL.Path,
		RawPath:  reqURL.RawPath,
		RawQuery: reqURL.RawQuery,
		Fragment: reqURL.Fragment,
	}

	for k, v := range req.Header {
		fq.HeaderEntries[k] = &HeaderEntry{
			Values: v,
		}
	}

	if req.TLS != nil && req.TLS.PeerCertificates != nil && len(req.TLS.PeerCertificates) > 0 {
		fq.PeerCertificates = make([][]byte, len(req.TLS.PeerCertificates))
		for i, cert := range req.TLS.PeerCertificates {
			fq.PeerCertificates[i] = cert.Raw
		}
	}

	return &fq, nil
}

// ParseForwardedRequest generates a new http.Request that is comprised of the
// values in the given request's body, assuming it correctly parses into a
// ForwardedRequest.
func ParseForwardedHTTPRequest(req *http.Request) (*http.Request, error) {
	buf := bytes.NewBuffer(nil)
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return nil, err
	}

	fq := new(Request)
	switch os.Getenv("VAULT_MESSAGE_TYPE") {
	case "json", "json_compress":
		err = jsonutil.DecodeJSON(buf.Bytes(), fq)
	default:
		err = proto.Unmarshal(buf.Bytes(), fq)
	}
	if err != nil {
		return nil, err
	}

	return ParseForwardedRequest(fq)
}

func ParseForwardedRequest(fq *Request) (*http.Request, error) {
	buf := bufCloser{
		Buffer: bytes.NewBuffer(fq.Body),
	}

	ret := &http.Request{
		Method:     fq.Method,
		Header:     make(map[string][]string, len(fq.HeaderEntries)),
		Body:       buf,
		Host:       fq.Host,
		RemoteAddr: fq.RemoteAddr,
	}

	ret.URL = &url.URL{
		Scheme:   fq.Url.Scheme,
		Opaque:   fq.Url.Opaque,
		Host:     fq.Url.Host,
		Path:     fq.Url.Path,
		RawPath:  fq.Url.RawPath,
		RawQuery: fq.Url.RawQuery,
		Fragment: fq.Url.Fragment,
	}

	for k, v := range fq.HeaderEntries {
		ret.Header[k] = v.Values
	}

	if fq.PeerCertificates != nil && len(fq.PeerCertificates) > 0 {
		ret.TLS = &tls.ConnectionState{
			PeerCertificates: make([]*x509.Certificate, len(fq.PeerCertificates)),
		}
		for i, certBytes := range fq.PeerCertificates {
			cert, err := x509.ParseCertificate(certBytes)
			if err != nil {
				return nil, err
			}
			ret.TLS.PeerCertificates[i] = cert
		}
	}

	return ret, nil
}

type RPCResponseWriter struct {
	statusCode int
	header     http.Header
	body       *bytes.Buffer
}

// NewRPCResponseWriter returns an initialized RPCResponseWriter
func NewRPCResponseWriter() *RPCResponseWriter {
	w := &RPCResponseWriter{
		header:     make(http.Header),
		body:       new(bytes.Buffer),
		statusCode: 200,
	}
	//w.header.Set("Content-Type", "application/octet-stream")
	return w
}

func (w *RPCResponseWriter) Header() http.Header {
	return w.header
}

func (w *RPCResponseWriter) Write(buf []byte) (int, error) {
	w.body.Write(buf)
	return len(buf), nil
}

func (w *RPCResponseWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *RPCResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *RPCResponseWriter) Body() *bytes.Buffer {
	return w.body
}
