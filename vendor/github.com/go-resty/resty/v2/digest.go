// Copyright (c) 2015-2024 Jeevanandam M (jeeva@myjeeva.com)
// 2023 Segev Dagan (https://github.com/segevda)
// 2024 Philipp Wolfer (https://github.com/phw)
// All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"
)

var (
	ErrDigestBadChallenge    = errors.New("digest: challenge is bad")
	ErrDigestCharset         = errors.New("digest: unsupported charset")
	ErrDigestAlgNotSupported = errors.New("digest: algorithm is not supported")
	ErrDigestQopNotSupported = errors.New("digest: no supported qop in list")
	ErrDigestNoQop           = errors.New("digest: qop must be specified")
)

var hashFuncs = map[string]func() hash.Hash{
	"":                 md5.New,
	"MD5":              md5.New,
	"MD5-sess":         md5.New,
	"SHA-256":          sha256.New,
	"SHA-256-sess":     sha256.New,
	"SHA-512-256":      sha512.New,
	"SHA-512-256-sess": sha512.New,
}

type digestCredentials struct {
	username, password string
}

type digestTransport struct {
	digestCredentials
	transport http.RoundTripper
}

func (dt *digestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Copy the request, so we don't modify the input.
	req2 := new(http.Request)
	*req2 = *req
	req2.Header = make(http.Header)
	for k, s := range req.Header {
		req2.Header[k] = s
	}

	// Fix http: ContentLength=xxx with Body length 0
	if req2.Body == nil {
		req2.ContentLength = 0
	} else if req2.GetBody != nil {
		var err error
		req2.Body, err = req2.GetBody()
		if err != nil {
			return nil, err
		}
	}

	// Make a request to get the 401 that contains the challenge.
	resp, err := dt.transport.RoundTrip(req)
	if err != nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}
	chal := resp.Header.Get(hdrWwwAuthenticateKey)
	if chal == "" {
		return resp, ErrDigestBadChallenge
	}

	c, err := parseChallenge(chal)
	if err != nil {
		return resp, err
	}

	// Form credentials based on the challenge
	cr := dt.newCredentials(req2, c)
	auth, err := cr.authorize()
	if err != nil {
		return resp, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Make authenticated request
	req2.Header.Set(hdrAuthorizationKey, auth)
	return dt.transport.RoundTrip(req2)
}

func (dt *digestTransport) newCredentials(req *http.Request, c *challenge) *credentials {
	return &credentials{
		username:   dt.username,
		userhash:   c.userhash,
		realm:      c.realm,
		nonce:      c.nonce,
		digestURI:  req.URL.RequestURI(),
		algorithm:  c.algorithm,
		sessionAlg: strings.HasSuffix(c.algorithm, "-sess"),
		opaque:     c.opaque,
		messageQop: c.qop,
		nc:         0,
		method:     req.Method,
		password:   dt.password,
	}
}

type challenge struct {
	realm     string
	domain    string
	nonce     string
	opaque    string
	stale     string
	algorithm string
	qop       string
	userhash  string
}

func (c *challenge) setValue(k, v string) error {
	switch k {
	case "realm":
		c.realm = v
	case "domain":
		c.domain = v
	case "nonce":
		c.nonce = v
	case "opaque":
		c.opaque = v
	case "stale":
		c.stale = v
	case "algorithm":
		c.algorithm = v
	case "qop":
		c.qop = v
	case "charset":
		if strings.ToUpper(v) != "UTF-8" {
			return ErrDigestCharset
		}
	case "userhash":
		c.userhash = v
	default:
		return ErrDigestBadChallenge
	}
	return nil
}

func parseChallenge(input string) (*challenge, error) {
	const ws = " \n\r\t"
	s := strings.Trim(input, ws)
	if !strings.HasPrefix(s, "Digest ") {
		return nil, ErrDigestBadChallenge
	}
	s = strings.Trim(s[7:], ws)
	c := &challenge{}
	b := strings.Builder{}
	key := ""
	quoted := false
	for _, r := range s {
		switch r {
		case '"':
			quoted = !quoted
		case ',':
			if quoted {
				b.WriteRune(r)
			} else {
				val := strings.Trim(b.String(), ws)
				b.Reset()
				if err := c.setValue(key, val); err != nil {
					return nil, err
				}
				key = ""
			}
		case '=':
			if quoted {
				b.WriteRune(r)
			} else {
				key = strings.Trim(b.String(), ws)
				b.Reset()
			}
		default:
			b.WriteRune(r)
		}
	}
	if quoted || (key == "" && b.Len() > 0) {
		return nil, ErrDigestBadChallenge
	}
	if key != "" {
		val := strings.Trim(b.String(), ws)
		if err := c.setValue(key, val); err != nil {
			return nil, err
		}
	}
	return c, nil
}

type credentials struct {
	username   string
	userhash   string
	realm      string
	nonce      string
	digestURI  string
	algorithm  string
	sessionAlg bool
	cNonce     string
	opaque     string
	messageQop string
	nc         int
	method     string
	password   string
}

func (c *credentials) authorize() (string, error) {
	if _, ok := hashFuncs[c.algorithm]; !ok {
		return "", ErrDigestAlgNotSupported
	}

	if err := c.validateQop(); err != nil {
		return "", err
	}

	resp, err := c.resp()
	if err != nil {
		return "", err
	}

	sl := make([]string, 0, 10)
	if c.userhash == "true" {
		// RFC 7616 3.4.4
		c.username = c.h(fmt.Sprintf("%s:%s", c.username, c.realm))
		sl = append(sl, fmt.Sprintf(`userhash=%s`, c.userhash))
	}
	sl = append(sl, fmt.Sprintf(`username="%s"`, c.username))
	sl = append(sl, fmt.Sprintf(`realm="%s"`, c.realm))
	sl = append(sl, fmt.Sprintf(`nonce="%s"`, c.nonce))
	sl = append(sl, fmt.Sprintf(`uri="%s"`, c.digestURI))
	sl = append(sl, fmt.Sprintf(`response="%s"`, resp))
	sl = append(sl, fmt.Sprintf(`algorithm=%s`, c.algorithm))
	if c.opaque != "" {
		sl = append(sl, fmt.Sprintf(`opaque="%s"`, c.opaque))
	}
	if c.messageQop != "" {
		sl = append(sl, fmt.Sprintf("qop=%s", c.messageQop))
		sl = append(sl, fmt.Sprintf("nc=%08x", c.nc))
		sl = append(sl, fmt.Sprintf(`cnonce="%s"`, c.cNonce))
	}

	return fmt.Sprintf("Digest %s", strings.Join(sl, ", ")), nil
}

func (c *credentials) validateQop() error {
	// Currently only supporting auth quality of protection. TODO: add auth-int support
	// NOTE: cURL support auth-int qop for requests other than POST and PUT (i.e. w/o body) by hashing an empty string
	// is this applicable for resty? see: https://github.com/curl/curl/blob/307b7543ea1e73ab04e062bdbe4b5bb409eaba3a/lib/vauth/digest.c#L774
	if c.messageQop == "" {
		return ErrDigestNoQop
	}
	possibleQops := strings.Split(c.messageQop, ",")
	var authSupport bool
	for _, qop := range possibleQops {
		qop = strings.TrimSpace(qop)
		if qop == "auth" {
			authSupport = true
			break
		}
	}
	if !authSupport {
		return ErrDigestQopNotSupported
	}

	c.messageQop = "auth"

	return nil
}

func (c *credentials) h(data string) string {
	hfCtor := hashFuncs[c.algorithm]
	hf := hfCtor()
	_, _ = hf.Write([]byte(data)) // Hash.Write never returns an error
	return fmt.Sprintf("%x", hf.Sum(nil))
}

func (c *credentials) resp() (string, error) {
	c.nc++

	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	c.cNonce = fmt.Sprintf("%x", b)[:32]

	ha1 := c.ha1()
	ha2 := c.ha2()

	return c.kd(ha1, fmt.Sprintf("%s:%08x:%s:%s:%s",
		c.nonce, c.nc, c.cNonce, c.messageQop, ha2)), nil
}

func (c *credentials) kd(secret, data string) string {
	return c.h(fmt.Sprintf("%s:%s", secret, data))
}

// RFC 7616 3.4.2
func (c *credentials) ha1() string {
	ret := c.h(fmt.Sprintf("%s:%s:%s", c.username, c.realm, c.password))
	if c.sessionAlg {
		return c.h(fmt.Sprintf("%s:%s:%s", ret, c.nonce, c.cNonce))
	}

	return ret
}

// RFC 7616 3.4.3
func (c *credentials) ha2() string {
	// currently no auth-int support
	return c.h(fmt.Sprintf("%s:%s", c.method, c.digestURI))
}
