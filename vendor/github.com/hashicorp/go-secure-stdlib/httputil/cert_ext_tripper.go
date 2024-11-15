// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package httputil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type ignoreExtensionsRoundTripper struct {
	base         *http.Transport
	extsToIgnore []asn1.ObjectIdentifier
}

// NewIgnoreUnhandledExtensionsRoundTripper creates a RoundTripper that may be used in an HTTP client which will
// ignore the provided extensions if presently unhandled on a certificate.  If base is nil, the default RoundTripper is used.
func NewIgnoreUnhandledExtensionsRoundTripper(base http.RoundTripper, extsToIgnore []asn1.ObjectIdentifier) (http.RoundTripper, error) {
	if len(extsToIgnore) == 0 {
		return nil, errors.New("no extensions ignored, should use original RoundTripper")
	}
	if base == nil {
		base = http.DefaultTransport
	}

	tp, ok := base.(*http.Transport)
	if !ok {
		// We don't know how to deal with this object, bail
		return base, nil
	}
	if tp != nil && (tp.TLSClientConfig != nil && (tp.TLSClientConfig.InsecureSkipVerify || tp.TLSClientConfig.VerifyConnection != nil)) {
		// Already not verifying or verifying in a custom fashion
		return nil, errors.New("cannot ignore provided extensions, base RoundTripper already handling or skipping verification")
	}
	return &ignoreExtensionsRoundTripper{base: tp, extsToIgnore: extsToIgnore}, nil
}

func (i *ignoreExtensionsRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	domain, _, err := net.SplitHostPort(request.URL.Host)
	if err != nil {
		if strings.Contains(err.Error(), "missing port") {
			domain = request.URL.Host
		} else {
			return nil, fmt.Errorf("error splitting host/port: %w", err)
		}
	}

	var tlsConfig *tls.Config
	perReqTransport := i.base.Clone()
	if perReqTransport.TLSClientConfig != nil {
		tlsConfig = perReqTransport.TLSClientConfig.Clone()
	} else {
		tlsConfig = &tls.Config{}
	}

	// Domain may be an IP address, in which case we shouldn't set ServerName
	var ipBased bool
	if addr := net.ParseIP(domain); addr == nil {
		tlsConfig.ServerName = domain
	} else {
		ipBased = true
	}

	tlsConfig.InsecureSkipVerify = true
	connectionVerifier := i.customVerifyConnection(tlsConfig, ipBased)
	tlsConfig.VerifyConnection = connectionVerifier

	perReqTransport.TLSClientConfig = tlsConfig
	return perReqTransport.RoundTrip(request)
}

func (i *ignoreExtensionsRoundTripper) customVerifyConnection(tc *tls.Config, ipBased bool) func(tls.ConnectionState) error {
	return func(cs tls.ConnectionState) error {
		certs := cs.PeerCertificates

		serverName := cs.ServerName
		if cs.ServerName == "" && !ipBased {
			if tc.ServerName == "" {
				return fmt.Errorf("the ServerName in TLSClientConfig is required to be set when UnhandledExtensionsToIgnore has values")
			}
			serverName = tc.ServerName
		} else if cs.ServerName != tc.ServerName {
			return fmt.Errorf("connection state server name (%s) does not match requested (%s)", cs.ServerName, tc.ServerName)
		}

		for _, cert := range certs {
			if len(cert.UnhandledCriticalExtensions) == 0 {
				continue
			}
			var remainingUnhandled []asn1.ObjectIdentifier
			for _, ext := range cert.UnhandledCriticalExtensions {
				shouldRemove := i.isExtInIgnore(ext)
				if !shouldRemove {
					remainingUnhandled = append(remainingUnhandled, ext)
				}
			}
			cert.UnhandledCriticalExtensions = remainingUnhandled
		}

		// Now verify with the requested extensions removed
		opts := x509.VerifyOptions{
			Roots:         tc.RootCAs,
			DNSName:       serverName,
			Intermediates: x509.NewCertPool(),
		}

		for _, cert := range certs[1:] {
			opts.Intermediates.AddCert(cert)
		}

		_, err := certs[0].Verify(opts)
		if err != nil {
			return &tls.CertificateVerificationError{UnverifiedCertificates: certs, Err: err}
		}

		return nil
	}
}

func (i *ignoreExtensionsRoundTripper) isExtInIgnore(ext asn1.ObjectIdentifier) bool {
	for _, extToIgnore := range i.extsToIgnore {
		if ext.Equal(extToIgnore) {
			return true
		}
	}

	return false
}
