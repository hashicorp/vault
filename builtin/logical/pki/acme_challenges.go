// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	DNSChallengePrefix = "_acme-challenge."
	ALPNProtocol       = "acme-tls/1"
)

// While this should be a constant, there's no way to do a low-level test of
// ValidateTLSALPN01Challenge without spinning up a complicated Docker
// instance to build a custom responder. Because we already have a local
// toolchain, it is far easier to drive this through Go tests with a custom
// (high) port, rather than requiring permission to bind to port 443 (root-run
// tests are even worse).
var ALPNPort = "443"

// OID of the acmeIdentifier X.509 Certificate Extension.
var OIDACMEIdentifier = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 1, 31}

// ValidateKeyAuthorization validates that the given keyAuthz from a challenge
// matches our expectation, returning (true, nil) if so, or (false, err) if
// not.
func ValidateKeyAuthorization(keyAuthz string, token string, thumbprint string) (bool, error) {
	parts := strings.Split(keyAuthz, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("%w: %s", ErrMalformed, fmt.Errorf("invalid authorization: got %v parts, expected 2", len(parts)).Error())
	}

	tokenPart := parts[0]
	thumbprintPart := parts[1]

	if token != tokenPart || thumbprint != thumbprintPart {
		return false, fmt.Errorf("%w: %s", ErrIncorrectResponse, fmt.Errorf("key authorization was invalid").Error())
	}

	return true, nil
}

// ValidateSHA256KeyAuthorization validates that the given keyAuthz from a
// challenge matches our expectation, returning (true, nil) if so, or
// (false, err) if not.
//
// This is for use with DNS challenges, which require base64 encoding.
func ValidateSHA256KeyAuthorization(keyAuthz string, token string, thumbprint string) (bool, error) {
	authzContents := token + "." + thumbprint
	checksum := sha256.Sum256([]byte(authzContents))
	expectedAuthz := base64.RawURLEncoding.EncodeToString(checksum[:])

	if keyAuthz != expectedAuthz {
		return false, fmt.Errorf("sha256 key authorization was invalid")
	}

	return true, nil
}

// ValidateRawSHA256KeyAuthorization validates that the given keyAuthz from a
// challenge matches our expectation, returning (true, nil) if so, or
// (false, err) if not.
//
// This is for use with TLS challenges, which require the raw hash output.
func ValidateRawSHA256KeyAuthorization(keyAuthz []byte, token string, thumbprint string) (bool, error) {
	authzContents := token + "." + thumbprint
	expectedAuthz := sha256.Sum256([]byte(authzContents))

	if len(keyAuthz) != len(expectedAuthz) || subtle.ConstantTimeCompare(expectedAuthz[:], keyAuthz) != 1 {
		return false, fmt.Errorf("sha256 key authorization was invalid")
	}

	return true, nil
}

func buildResolver(config *acmeConfigEntry) (*net.Resolver, error) {
	if len(config.DNSResolver) == 0 {
		return net.DefaultResolver, nil
	}

	return &net.Resolver{
		PreferGo:     true,
		StrictErrors: false,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
			}
			return d.DialContext(ctx, network, config.DNSResolver)
		},
	}, nil
}

func buildDialerConfig(config *acmeConfigEntry) (*net.Dialer, error) {
	resolver, err := buildResolver(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build resolver: %w", err)
	}

	return &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: -1 * time.Second,
		Resolver:  resolver,
	}, nil
}

// Validates a given ACME http-01 challenge against the specified domain,
// per RFC 8555.
//
// We attempt to be defensive here against timeouts, extra redirects, &c.
func ValidateHTTP01Challenge(domain string, token string, thumbprint string, config *acmeConfigEntry) (bool, error) {
	path := "http://" + domain + "/.well-known/acme-challenge/" + token
	dialer, err := buildDialerConfig(config)
	if err != nil {
		return false, fmt.Errorf("%w: %s", ErrServerInternal, fmt.Errorf("failed to build dialer: %w", err).Error())
	}

	transport := &http.Transport{
		// Only a single request is sent to this server as we do not do any
		// batching of validation attempts. There is no need to do an HTTP
		// KeepAlive as a result.
		DisableKeepAlives:   true,
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
		MaxConnsPerHost:     1,
		IdleConnTimeout:     1 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},

		// We'd rather timeout and re-attempt validation later than hang
		// too many validators waiting for slow hosts.
		DialContext:           dialer.DialContext,
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
		return false, fmt.Errorf("%w: %s", ErrConnection, fmt.Errorf("http-01: failed to fetch path %v: %w", path, err).Error())
	}

	// We provision a buffer which allows for a variable size challenge, some
	// whitespace, and a detection gap for too long of a message.
	minExpected := len(token) + 1 + len(thumbprint)
	maxExpected := 512

	defer resp.Body.Close()

	// Attempt to read the body, but don't do so infinitely.
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(maxExpected+1)))
	if err != nil {
		return false, fmt.Errorf("%w: %s", ErrIncorrectResponse, fmt.Errorf("http-01: unexpected error while reading body: %w", err).Error())
	}

	if len(body) > maxExpected {
		return false, fmt.Errorf("%w: %s", ErrMalformed, fmt.Errorf("http-01: response too large: received %v > %v bytes", len(body), maxExpected).Error())
	}

	if len(body) < minExpected {
		return false, fmt.Errorf("%w: %s", ErrMalformed, fmt.Errorf("http-01: response too small: received %v < %v bytes", len(body), minExpected).Error())
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

func ValidateDNS01Challenge(domain string, token string, thumbprint string, config *acmeConfigEntry) (bool, error) {
	// Here, domain is the value from the post-wildcard-processed identifier.
	// Per RFC 8555, no difference in validation occurs if a wildcard entry
	// is requested or if a non-wildcard entry is requested.
	//
	// XXX: In this case the DNS server is operator controlled and is assumed
	// to be less malicious so the default resolver is used. In the future,
	// we'll want to use net.Resolver for two reasons:
	//
	// 1. To control the actual resolver via ACME configuration,
	// 2. To use a context to set stricter timeout limits.
	resolver, err := buildResolver(config)
	if err != nil {
		return false, fmt.Errorf("%w: %s", ErrServerInternal, fmt.Errorf("failed to build resolver: %w", err).Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	name := DNSChallengePrefix + domain
	results, err := resolver.LookupTXT(ctx, name)
	if err != nil {
		return false, fmt.Errorf("%w: %s", ErrDNS, fmt.Errorf("dns-01: failed to lookup TXT records for domain (%v) via resolver %v: %w", name, config.DNSResolver, err).Error())
	}

	for _, keyAuthz := range results {
		ok, _ := ValidateSHA256KeyAuthorization(keyAuthz, token, thumbprint)
		if ok {
			return true, nil
		}
	}

	return false, fmt.Errorf("%w: %s", ErrDNS, fmt.Errorf("dns-01: challenge failed against %v records", len(results)).Error())
}

func ValidateTLSALPN01Challenge(domain string, token string, thumbprint string, config *acmeConfigEntry) (bool, error) {
	// This RFC is defined in RFC 8737 Automated Certificate Management
	// Environment (ACME) TLS Application‑Layer Protocol Negotiation
	// (ALPN) Challenge Extension.
	//
	// This is conceptually similar to ValidateHTTP01Challenge, but
	// uses a TLS connection on port 443 with the specified ALPN
	// protocol.

	cfg := &tls.Config{
		// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
		// Negotiation (TLS ALPN) Challenge, the name of the negotiated
		// protocol is "acme-tls/1".
		NextProtos: []string{ALPNProtocol},

		// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
		// Negotiation (TLS ALPN) Challenge:
		//
		// > ... and an SNI extension containing only the domain name
		// > being validated during the TLS handshake.
		//
		// According to the Go docs, setting this option (even though
		// InsecureSkipVerify=true is also specified), allows us to
		// set the SNI extension to this value.
		ServerName: domain,

		VerifyConnection: func(connState tls.ConnectionState) error {
			// We initiated a fresh connection with no session tickets;
			// even if we did have a session ticket, we do not wish to
			// use it. Verify that the server has not inadvertently
			// reused connections between validation attempts or something.
			if connState.DidResume {
				return fmt.Errorf("server under test incorrectly reported that handshake was resumed when no session cache was provided; refusing to continue")
			}

			// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
			// Negotiation (TLS ALPN) Challenge:
			//
			// > The ACME server verifies that during the TLS handshake the
			// > application-layer protocol "acme-tls/1" was successfully
			// > negotiated (and that the ALPN extension contained only the
			// > value "acme-tls/1").
			if connState.NegotiatedProtocol != ALPNProtocol {
				return fmt.Errorf("server under test negotiated unexpected ALPN protocol %v", connState.NegotiatedProtocol)
			}

			// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
			// Negotiation (TLS ALPN) Challenge:
			//
			// > and that the certificate returned
			//
			// Because this certificate MUST be self-signed (per earlier
			// statement in RFC 8737 Section 3), there is no point in sending
			// more than one certificate, and so we will err early here if
			// we got more than one.
			if len(connState.PeerCertificates) > 1 {
				return fmt.Errorf("server under test returned multiple (%v) certificates when we expected only one", len(connState.PeerCertificates))
			}
			cert := connState.PeerCertificates[0]

			// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
			// Negotiation (TLS ALPN) Challenge:
			//
			// > The client prepares for validation by constructing a
			// > self-signed certificate that MUST contain an acmeIdentifier
			// > extension and a subjectAlternativeName extension [RFC5280].
			//
			// Verify that this is a self-signed certificate that isn't signed
			// by another certificate (i.e., with the same key material but
			// different issuer).
			// NOTE: Do not use cert.CheckSignatureFrom(cert) as we need to bypass the
			//       checks for the parent certificate having the IsCA basic constraint set.
			err := cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)
			if err != nil {
				return fmt.Errorf("server under test returned a non-self-signed certificate: %w", err)
			}

			if !bytes.Equal(cert.RawSubject, cert.RawIssuer) {
				return fmt.Errorf("server under test returned a non-self-signed certificate: invalid subject (%v) <-> issuer (%v) match", cert.Subject.String(), cert.Issuer.String())
			}

			// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
			// Negotiation (TLS ALPN) Challenge:
			//
			// > The subjectAlternativeName extension MUST contain a single
			// > dNSName entry where the value is the domain name being
			// > validated.
			//
			// TODO: this does not validate that there are not other SANs
			// with unknown (to Go) OIDs.
			if len(cert.DNSNames) != 1 || len(cert.EmailAddresses) > 0 || len(cert.IPAddresses) > 0 || len(cert.URIs) > 0 {
				return fmt.Errorf("server under test returned a certificate with incorrect SANs")
			}

			// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
			// Negotiation (TLS ALPN) Challenge:
			//
			// > The comparison of dNSNames MUST be case insensitive
			// > [RFC4343]. Note that as ACME doesn't support Unicode
			// > identifiers, all dNSNames MUST be encoded using the rules
			// > of [RFC3492].
			if !strings.EqualFold(cert.DNSNames[0], domain) {
				return fmt.Errorf("server under test returned a certificate with unexpected identifier: %v", cert.DNSNames[0])
			}

			// Per above, verify that the acmeIdentifier extension is present
			// exactly once and has the correct value.
			var foundACMEId bool
			for _, ext := range cert.Extensions {
				if !ext.Id.Equal(OIDACMEIdentifier) {
					continue
				}

				// There must be only a single ACME extension.
				if foundACMEId {
					return fmt.Errorf("server under test returned a certificate with multiple acmeIdentifier extensions")
				}
				foundACMEId = true

				// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
				// Negotiation (TLS ALPN) Challenge:
				//
				// > a critical acmeIdentifier extension
				if !ext.Critical {
					return fmt.Errorf("server under test returned a certificate with an acmeIdentifier extension marked non-Critical")
				}

				var keyAuthz []byte
				remainder, err := asn1.Unmarshal(ext.Value, &keyAuthz)
				if err != nil {
					return fmt.Errorf("server under test returned a certificate with invalid acmeIdentifier extension value: %w", err)
				}
				if len(remainder) > 0 {
					return fmt.Errorf("server under test returned a certificate with invalid acmeIdentifier extension value with additional trailing data")
				}

				ok, err := ValidateRawSHA256KeyAuthorization(keyAuthz, token, thumbprint)
				if !ok || err != nil {
					return fmt.Errorf("server under test returned a certificate with an invalid key authorization (%w)", err)
				}
			}

			// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
			// Negotiation (TLS ALPN) Challenge:
			//
			// > The ACME server verifies that ... the certificate returned
			// > contains: ... a critical acmeIdentifier extension containing
			// > the expected SHA-256 digest computed in step 1.
			if !foundACMEId {
				return fmt.Errorf("server under test returned a certificate without the required acmeIdentifier extension")
			}

			// Remove the handled critical extension and validate that we
			// have no additional critical extensions left unhandled.
			var index int = -1
			for oidIndex, oid := range cert.UnhandledCriticalExtensions {
				if oid.Equal(OIDACMEIdentifier) {
					index = oidIndex
					break
				}
			}
			if index != -1 {
				// Unlike the foundACMEId case, this is not a failure; if Go
				// updates to "understand" this critical extension, we do not
				// wish to fail.
				cert.UnhandledCriticalExtensions = append(cert.UnhandledCriticalExtensions[0:index], cert.UnhandledCriticalExtensions[index+1:]...)
			}
			if len(cert.UnhandledCriticalExtensions) > 0 {
				return fmt.Errorf("server under test returned a certificate with additional unknown critical extensions (%v)", cert.UnhandledCriticalExtensions)
			}

			// All good!
			return nil
		},

		// We never want to resume a connection; do not provide session
		// cache storage.
		ClientSessionCache: nil,

		// Do not trust any system trusted certificates; we're going to be
		// manually validating the chain, so specifying a non-empty pool
		// here could only cause additional, unnecessary work.
		RootCAs: x509.NewCertPool(),

		// Do not bother validating the client's chain; we know it should be
		// self-signed. This also disables hostname verification, but we do
		// this verification as part of VerifyConnection(...) ourselves.
		//
		// Per Go docs, this option is only safe in conjunction with
		// VerifyConnection which we define above.
		InsecureSkipVerify: true,

		// RFC 8737 Section 4. acme-tls/1 Protocol Definition:
		//
		// > ACME servers that implement "acme-tls/1" MUST only negotiate
		// > TLS 1.2 [RFC5246] or higher when connecting to clients for
		// > validation.
		MinVersion: tls.VersionTLS12,

		// While RFC 8737 does not place restrictions around allowed cipher
		// suites, we wish to restrict ourselves to secure defaults. Specify
		// the Intermediate guideline from Mozilla's TLS config generator to
		// disable obviously weak ciphers.
		//
		// See also: https://ssl-config.mozilla.org/#server=go&version=1.14.4&config=intermediate&guideline=5.7
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	// Build a dialer using our custom DNS resolver, to ensure domains get
	// resolved according to configuration.
	dialer, err := buildDialerConfig(config)
	if err != nil {
		return false, fmt.Errorf("%w: %s", ErrServerInternal, fmt.Errorf("failed to build dialer: %w", err).Error())
	}

	// Per RFC 8737 Section 3. TLS with Application-Layer Protocol
	// Negotiation (TLS ALPN) Challenge:
	//
	// > 2. The ACME server resolves the domain name being validated and
	// >    chooses one of the IP addresses returned for validation (the
	// >    server MAY validate against multiple addresses if more than
	// >    one is returned).
	// > 3. The ACME server initiates a TLS connection to the chosen IP
	// >    address. This connection MUST use TCP port 443.
	address := fmt.Sprintf("%v:"+ALPNPort, domain)
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return false, fmt.Errorf("%w: %s", ErrConnection, fmt.Errorf("tls-alpn-01: failed to dial host: %w", err).Error())
	}

	// Initiate the connection to the remote peer.
	client := tls.Client(conn, cfg)

	// We intentionally swallow this error as it isn't useful to the
	// underlying protocol we perform here. Notably, per RFC 8737
	// Section 4. acme-tls/1 Protocol Definition:
	//
	// > Once the handshake is completed, the client MUST NOT exchange
	// > any further data with the server and MUST immediately close the
	// > connection. ... Because of this, an ACME server MAY choose to
	// > withhold authorization if either the certificate signature is
	// > invalid or the handshake doesn't fully complete.
	defer client.Close()

	// We wish to put time bounds on the total time the handshake can
	// stall for, so build a connection context here.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// See note above about why we can allow Handshake to complete
	// successfully.
	if err := client.HandshakeContext(ctx); err != nil {
		return false, fmt.Errorf("%w: %s", ErrTLS, fmt.Errorf("tls-alpn-01: failed to perform handshake: %w", err).Error())
	}
	return true, nil
}
