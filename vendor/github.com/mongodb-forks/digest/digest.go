// Copyright 2013 M-Lab, 2020 MongoDB, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The digest package provides an implementation of http.RoundTripper that takes
// care of HTTP Digest Authentication (https://tools.ietf.org/html/rfc7616).

// This only implements the MD5, SHA-256 and "auth" portions of the RFC, which
// is enough to authenticate the client to the majority of available
// server side implementations using Digest Access Authentication (e.g.,
// the Apache Web Server).
//
//
// Example usage:
//
//	t := NewTransport("myUserName", "myP@55w0rd")
//	req, err := http.NewRequest("GET", "http://notreal.com/path?arg=1", nil)
//	if err != nil {
//		return err
//	}
//	resp, err := t.RoundTrip(req)
//	if err != nil {
//		return err
//	}
//
// OR it can be used as a client:
//
//	c, err := t.Client()
//	if err != nil {
//		return err
//	}
//	resp, err := c.Get("http://notreal.com/path?arg=1")
//	if err != nil {
//		return err
//	}
//
// OR if you want fine-grained control over timeouts
//
// t := &digest.Transport{Username: "myUserName", Password: "myP@55w0rd"}
// t.Transport = &http.Transport{
// 	DialContext: (&net.Dialer{
// 		Timeout:   30 * time.Second,
// 		KeepAlive: 10 * time.Second,
// 	}).DialContext,
// 	ExpectContinueTimeout: 10 * time.Second,
// 	IdleConnTimeout:       60 * time.Second,
// 	MaxIdleConns:          100,
// 	MaxIdleConnsPerHost:   4,
// 	Proxy:                 http.ProxyFromEnvironment,
// 	ResponseHeaderTimeout: 30 * time.Second,
// 	TLSHandshakeTimeout:   10 * time.Second,
// }
// c, err := t.Client()
// if err != nil {
// 	return nil, err
// }
// resp, err := c.Get("http://notreal.com/path?arg=1")
// if err != nil {
// 	return nil, err
// }

package digest

import (
	"bytes"
	"crypto/md5" //nolint:gosec // valid for digest
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"
)

const (
	MsgAuth   string = "auth"
	AlgMD5    string = "MD5"
	AlgSha256 string = "SHA-256"
)

var (
	ErrNilTransport      = errors.New("transport is nil")
	ErrBadChallenge      = errors.New("challenge is bad")
	ErrAlgNotImplemented = errors.New("alg not implemented")
)

// Transport is an implementation of http.RoundTripper that takes care of http
// digest authentication.
type Transport struct {
	Username  string
	Password  string
	Transport http.RoundTripper
}

// NewTransport creates a new digest transport using the http.DefaultTransport.
func NewTransport(username, password string) *Transport {
	t := &Transport{
		Username: username,
		Password: password,
	}
	t.Transport = http.DefaultTransport
	return t
}

// NewTransportWithHTTPTransport creates a new digest transport using the supplied http.Transport.
func NewTransportWithHTTPTransport(username, password string, transport *http.Transport) *Transport {
	t := &Transport{
		Username:  username,
		Password:  password,
		Transport: transport,
	}
	return t
}

// NewTransportWithHTTPRoundTripper creates a new digest transport using the supplied http.RoundTripper interface.
func NewTransportWithHTTPRoundTripper(username, password string, transport http.RoundTripper) *Transport {
	t := &Transport{
		Username:  username,
		Password:  password,
		Transport: transport,
	}
	return t
}

type challenge struct {
	Realm     string
	Domain    string
	Nonce     string
	Opaque    string
	Stale     string
	Algorithm string
	Qop       string
}

func parseChallenge(input string) (*challenge, error) {
	const ws = " \n\r\t"
	const qs = `"`
	const n = 2
	s := strings.Trim(input, ws)
	if !strings.HasPrefix(s, "Digest ") {
		return nil, ErrBadChallenge
	}
	s = strings.Trim(s[7:], ws)
	sl := strings.Split(s, ", ")
	c := &challenge{
		Algorithm: AlgMD5,
	}
	var r []string
	for i := range sl {
		r = strings.SplitN(sl[i], "=", n)
		switch r[0] {
		case "realm":
			c.Realm = strings.Trim(r[1], qs)
		case "domain":
			c.Domain = strings.Trim(r[1], qs)
		case "nonce":
			c.Nonce = strings.Trim(r[1], qs)
		case "opaque":
			c.Opaque = strings.Trim(r[1], qs)
		case "stale":
			c.Stale = strings.Trim(r[1], qs)
		case "algorithm":
			c.Algorithm = strings.Trim(r[1], qs)
		case "qop":
			c.Qop = strings.Trim(r[1], qs)
		default:
			return nil, ErrBadChallenge
		}
	}
	return c, nil
}

type credentials struct {
	Username   string
	Realm      string
	Nonce      string
	DigestURI  string
	Algorithm  string
	Cnonce     string
	Opaque     string
	MessageQop string
	NonceCount int
	method     string
	password   string
	impl       hashingFunc
}

type hashingFunc func() hash.Hash

func h(data string, f hashingFunc) (string, error) {
	hf := f()
	if _, err := io.WriteString(hf, data); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hf.Sum(nil)), nil
}

func kd(secret, data string, f hashingFunc) (string, error) {
	return h(fmt.Sprintf("%s:%s", secret, data), f)
}

func (c *credentials) ha1() (string, error) {
	return h(fmt.Sprintf("%s:%s:%s", c.Username, c.Realm, c.password), c.impl)
}

func (c *credentials) ha2() (string, error) {
	return h(fmt.Sprintf("%s:%s", c.method, c.DigestURI), c.impl)
}

func (c *credentials) resp(cnonce string) (resp string, err error) {
	var ha1 string
	var ha2 string
	c.NonceCount++
	if c.MessageQop == MsgAuth {
		if cnonce != "" {
			c.Cnonce = cnonce
		} else {
			const size = 8
			b := make([]byte, size)
			_, err = io.ReadFull(rand.Reader, b)
			if err != nil {
				return "", err
			}
			c.Cnonce = fmt.Sprintf("%x", b)[:16]
		}
		if ha1, err = c.ha1(); err != nil {
			return "", err
		}
		if ha2, err = c.ha2(); err != nil {
			return "", err
		}
		return kd(ha1, fmt.Sprintf("%s:%08x:%s:%s:%s", c.Nonce, c.NonceCount, c.Cnonce, c.MessageQop, ha2), c.impl)
	} else if c.MessageQop == "" {
		if ha1, err = c.ha1(); err != nil {
			return "", err
		}
		if ha2, err = c.ha2(); err != nil {
			return "", err
		}
		return kd(ha1, fmt.Sprintf("%s:%s", c.Nonce, ha2), c.impl)
	}
	return "", ErrAlgNotImplemented
}

func (c *credentials) authorize() (string, error) {
	// Note that this is only implemented for MD5 and NOT MD5-sess.
	// MD5-sess is rarely supported and those that do are a big mess.
	if c.Algorithm != AlgMD5 && c.Algorithm != AlgSha256 {
		return "", ErrAlgNotImplemented
	}
	// Note that this is NOT implemented for "qop=auth-int".  Similarly the
	// auth-int server side implementations that do exist are a mess.
	if c.MessageQop != MsgAuth && c.MessageQop != "" {
		return "", ErrAlgNotImplemented
	}
	resp, err := c.resp("")
	if err != nil {
		return "", ErrAlgNotImplemented
	}
	sl := []string{fmt.Sprintf(`username=%q`, c.Username)}
	sl = append(sl, fmt.Sprintf(`realm=%q`, c.Realm),
		fmt.Sprintf(`nonce=%q`, c.Nonce),
		fmt.Sprintf(`uri=%q`, c.DigestURI),
		fmt.Sprintf(`response=%q`, resp))
	if c.Algorithm != "" {
		sl = append(sl, fmt.Sprintf(`algorithm=%q`, c.Algorithm))
	}
	if c.Opaque != "" {
		sl = append(sl, fmt.Sprintf(`opaque=%q`, c.Opaque))
	}
	if c.MessageQop != "" {
		sl = append(sl, fmt.Sprintf("qop=%s", c.MessageQop),
			fmt.Sprintf("nc=%08x", c.NonceCount),
			fmt.Sprintf(`cnonce=%q`, c.Cnonce))
	}
	return fmt.Sprintf("Digest %s", strings.Join(sl, ", ")), nil
}

func (t *Transport) newCredentials(req *http.Request, c *challenge) (*credentials, error) {
	cred := &credentials{
		Username:   t.Username,
		Realm:      c.Realm,
		Nonce:      c.Nonce,
		DigestURI:  req.URL.RequestURI(),
		Algorithm:  c.Algorithm,
		Opaque:     c.Opaque,
		MessageQop: c.Qop, // "auth" must be a single value
		NonceCount: 0,
		method:     req.Method,
		password:   t.Password,
	}
	switch c.Algorithm {
	case AlgMD5:
		cred.impl = md5.New
	case AlgSha256:
		cred.impl = sha256.New
	default:
		return nil, ErrAlgNotImplemented
	}

	return cred, nil
}

// RoundTrip makes a request expecting a 401 response that will require digest
// authentication.  It creates the credentials it needs and makes a follow-up
// request.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		return nil, ErrNilTransport
	}

	// Copy the request so we don't modify the input.
	origReq := new(http.Request)
	*origReq = *req
	origReq.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		origReq.Header[k] = s
	}

	// We'll need the request body twice. In some cases we can use GetBody
	// to obtain a fresh reader for the second request, which we do right
	// before the RoundTrip(origReq) call. If GetBody is unavailable, read
	// the body into a memory buffer and use it for both requests.
	if req.Body != nil && req.GetBody == nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))
		origReq.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	// Make a request to get the 401 that contains the challenge.
	challenge, resp, err := t.fetchChallenge(req)
	if challenge == "" || err != nil {
		return resp, err
	}

	c, err := parseChallenge(challenge)
	if err != nil {
		return nil, err
	}

	// Form credentials based on the challenge.
	cr, err := t.newCredentials(origReq, c)
	if err != nil {
		return nil, err
	}
	auth, err := cr.authorize()
	if err != nil {
		return nil, err
	}

	// Obtain a fresh body.
	if req.Body != nil && req.GetBody != nil {
		origReq.Body, err = req.GetBody()
		if err != nil {
			return nil, err
		}
	}

	// Make authenticated request.
	origReq.Header.Set("Authorization", auth)
	return t.Transport.RoundTrip(origReq)
}

func (t *Transport) fetchChallenge(req *http.Request) (string, *http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return "", resp, err
	}
	if resp.StatusCode != http.StatusUnauthorized {
		return "", resp, nil
	}

	// We'll no longer use the initial response, so close it
	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCP connection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			_, _ = io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
		}

		resp.Body.Close()
	}()
	return resp.Header.Get("WWW-Authenticate"), resp, nil
}

// Client returns an HTTP client that uses the digest transport.
func (t *Transport) Client() (*http.Client, error) {
	if t.Transport == nil {
		return nil, ErrNilTransport
	}
	return &http.Client{Transport: t}, nil
}
