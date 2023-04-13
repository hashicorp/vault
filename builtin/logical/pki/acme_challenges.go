package pki

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// ValidateKeyAuthorization validates that the given keyAuthz from a challenge
// matches our expectation, returning (true, nil) if so, or (false, err) if
// not.
func ValidateKeyAuthorization(keyAuthz string, token string, thumbprint string) (bool, error) {
	parts := strings.Split(keyAuthz, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid authorization: got %v parts, expected 2", len(parts))
	}

	tokenPart := parts[0]
	thumbprintPart := parts[1]

	if token != tokenPart || thumbprint != thumbprintPart {
		return false, fmt.Errorf("key authorization was invalid")
	}

	return true, nil
}

// Validates a given ACME http-01 challenge against the specified domain,
// per RFC 8555.
//
// We attempt to be defensive here against timeouts, extra redirects, &c.
func ValidateHTTP01Challenge(domain string, token string, thumbprint string) (bool, error) {
	path := "http://" + domain + "/.well-known/acme-challenge/" + token

	transport := &http.Transport{
		// Only a single request is sent to this server as we do not do any
		// batching of validation attempts. There is no need to do an HTTP
		// KeepAlive as a result.
		DisableKeepAlives:   true,
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
		MaxConnsPerHost:     1,
		IdleConnTimeout:     1 * time.Second,

		// We'd rather timeout and re-attempt validation later than hang
		// too many validators waiting for slow hosts.
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: -1 * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	maxRedirects := 10
	urlLength := 2000

	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via)+1 >= maxRedirects {
				return fmt.Errorf("http-01: too many redirects: %v", len(via)+1)
			}

			reqUrlLen := len(req.URL.String())
			if reqUrlLen > urlLength {
				return fmt.Errorf("http-01: redirect url length too long: %v", reqUrlLen)
			}

			return nil
		},
	}

	resp, err := client.Get(path)
	if err != nil {
		return false, fmt.Errorf("http-01: failed to fetch path %v: %w", path, err)
	}

	// We provision a buffer which allows for a variable size challenge, some
	// whitespace, and a detection gap for too long of a message.
	minExpected := len(token) + 1 + len(thumbprint)
	maxExpected := 512

	defer resp.Body.Close()

	// Attempt to read the body, but don't do so infinitely.
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(maxExpected+1)))
	if err != nil {
		return false, fmt.Errorf("http-01: unexpected error while reading body: %w", err)
	}

	if len(body) > maxExpected {
		return false, fmt.Errorf("http-01: response too large: received %v > %v bytes", len(body), maxExpected)
	}

	if len(body) < minExpected {
		return false, fmt.Errorf("http-01: response too small: received %v < %v bytes", len(body), minExpected)
	}

	// Per RFC 8555 Section 8.3. HTTP Challenge:
	//
	// > The server SHOULD ignore whitespace characters at the end of the body.
	keyAuthz := string(body)
	keyAuthz = strings.TrimSpace(keyAuthz)

	// If we got here, we got no non-EOF error while reading. Try to validate
	// the token because we're bounded by a reasonable amount of length.
	return ValidateKeyAuthorization(keyAuthz, token, thumbprint)
}
