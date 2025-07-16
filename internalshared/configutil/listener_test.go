// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListener_ParseSingleIPTemplate exercises the ParseSingleIPTemplate function to
// ensure that we only attempt to parse templates when the input contains a
// template placeholder (see: go-sockaddr/template).
func TestListener_ParseSingleIPTemplate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		arg             string
		want            string
		isErrorExpected bool
		errorMessage    string
	}{
		"test hostname": {
			arg:             "https://vaultproject.io:8200",
			want:            "https://vaultproject.io:8200",
			isErrorExpected: false,
		},
		"test ipv4": {
			arg:             "https://10.10.1.10:8200",
			want:            "https://10.10.1.10:8200",
			isErrorExpected: false,
		},
		"test ipv6 RFC-5952 4.1 conformance leading zeroes": {
			arg:             "https://[2001:0db8::0001]:8200",
			want:            "https://[2001:db8::1]:8200",
			isErrorExpected: false,
		},
		"test ipv6 RFC-5952 4.2.2 conformance one 16-bit 0 field": {
			arg:             "https://[2001:db8:0:1:1:1:1:1]:8200",
			want:            "https://[2001:db8:0:1:1:1:1:1]:8200",
			isErrorExpected: false,
		},
		"test ipv6 RFC-5952 4.2.3 conformance longest run of 0 bits shortened": {
			arg:             "https://[2001:0:0:1:0:0:0:1]:8200",
			want:            "https://[2001:0:0:1::1]:8200",
			isErrorExpected: false,
		},
		"test ipv6 RFC-5952 4.2.3 conformance equal runs of 0 bits shortened": {
			arg:             "https://[2001:db8:0:0:1:0:0:1]:8200",
			want:            "https://[2001:db8::1:0:0:1]:8200",
			isErrorExpected: false,
		},
		"test ipv6 RFC-5952 4.3 conformance downcase hex letters": {
			arg:             "https://[2001:DB8:AC3:FE4::1]:8200",
			want:            "https://[2001:db8:ac3:fe4::1]:8200",
			isErrorExpected: false,
		},
		"test invalid template func": {
			arg:             "{{ FooBar }}",
			want:            "",
			isErrorExpected: true,
			errorMessage:    "unable to parse address template",
		},
		"test partial template": {
			arg:             "{{FooBar",
			want:            "{{FooBar",
			isErrorExpected: false,
		},
	}
	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseSingleIPTemplate(tc.arg)

			if tc.isErrorExpected {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.want, got)
		})
	}
}

// TestListener_parseType exercises the listener receiver parseType.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		inputType       string
		inputFallback   string
		expectedValue   string
		isErrorExpected bool
		errorMessage    string
	}{
		"empty-all": {
			inputType:       "",
			inputFallback:   "",
			isErrorExpected: true,
			errorMessage:    "listener type must be specified",
		},
		"bad-type": {
			inputType:       "foo",
			isErrorExpected: true,
			errorMessage:    "unsupported listener type",
		},
		"bad-fallback": {
			inputType:       "",
			inputFallback:   "foo",
			isErrorExpected: true,
			errorMessage:    "unsupported listener type",
		},
		"tcp-type-lower": {
			inputType:       "tcp",
			expectedValue:   "tcp",
			isErrorExpected: false,
		},
		"tcp-type-upper": {
			inputType:       "TCP",
			expectedValue:   "tcp",
			isErrorExpected: false,
		},
		"tcp-type-mixed": {
			inputType:       "tCp",
			expectedValue:   "tcp",
			isErrorExpected: false,
		},
		"tcp-fallback-lower": {
			inputType:       "",
			inputFallback:   "tcp",
			expectedValue:   "tcp",
			isErrorExpected: false,
		},
		"tcp-fallback-upper": {
			inputType:       "",
			inputFallback:   "TCP",
			expectedValue:   "tcp",
			isErrorExpected: false,
		},
		"tcp-fallback-mixed": {
			inputType:       "",
			inputFallback:   "tCp",
			expectedValue:   "tcp",
			isErrorExpected: false,
		},
		"unix-type-lower": {
			inputType:       "unix",
			expectedValue:   "unix",
			isErrorExpected: false,
		},
		"unix-type-upper": {
			inputType:       "UNIX",
			expectedValue:   "unix",
			isErrorExpected: false,
		},
		"unix-type-mixed": {
			inputType:       "uNiX",
			expectedValue:   "unix",
			isErrorExpected: false,
		},
		"unix-fallback-lower": {
			inputType:       "",
			inputFallback:   "unix",
			expectedValue:   "unix",
			isErrorExpected: false,
		},
		"unix-fallback-upper": {
			inputType:       "",
			inputFallback:   "UNIX",
			expectedValue:   "unix",
			isErrorExpected: false,
		},
		"unix-fallback-mixed": {
			inputType:       "",
			inputFallback:   "uNiX",
			expectedValue:   "unix",
			isErrorExpected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := &Listener{Type: ListenerType(tc.inputType)}
			err := l.parseType(tc.inputFallback)
			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.expectedValue, l.Type.String())
			}
		})
	}
}

// TestListener_parseRequestSettings exercises the listener receiver parseRequestSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseRequestSettings(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		rawMaxRequestSize             any
		expectedMaxRequestSize        int64
		rawMaxRequestDuration         any
		expectedDuration              time.Duration
		rawRequireRequestHeader       any
		expectedRequireRequestHeader  bool
		rawDisableRequestLimiter      any
		expectedDisableRequestLimiter bool
		isErrorExpected               bool
		errorMessage                  string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"max-request-size-bad": {
			rawMaxRequestSize: "juan",
			isErrorExpected:   true,
			errorMessage:      "error parsing max_request_size",
		},
		"max-request-size-good": {
			rawMaxRequestSize:      "5",
			expectedMaxRequestSize: 5,
			isErrorExpected:        false,
		},
		"max-request-duration-bad": {
			rawMaxRequestDuration: "juan",
			isErrorExpected:       true,
			errorMessage:          "error parsing max_request_duration",
		},
		"max-request-duration-good": {
			rawMaxRequestDuration: "30s",
			expectedDuration:      30 * time.Second,
			isErrorExpected:       false,
		},
		"require-request-header-bad": {
			rawRequireRequestHeader:      "juan",
			expectedRequireRequestHeader: false,
			isErrorExpected:              true,
			errorMessage:                 "invalid value for require_request_header",
		},
		"require-request-header-good": {
			rawRequireRequestHeader:      "true",
			expectedRequireRequestHeader: true,
			isErrorExpected:              false,
		},
		"disable-request-limiter-bad": {
			rawDisableRequestLimiter:      "badvalue",
			expectedDisableRequestLimiter: false,
			isErrorExpected:               true,
			errorMessage:                  "invalid value for disable_request_limiter",
		},
		"disable-request-limiter-good": {
			rawDisableRequestLimiter:      "true",
			expectedDisableRequestLimiter: true,
			isErrorExpected:               false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				MaxRequestSizeRaw:        tc.rawMaxRequestSize,
				MaxRequestDurationRaw:    tc.rawMaxRequestDuration,
				RequireRequestHeaderRaw:  tc.rawRequireRequestHeader,
				DisableRequestLimiterRaw: tc.rawDisableRequestLimiter,
			}

			err := l.parseRequestSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Equal(t, tc.expectedMaxRequestSize, l.MaxRequestSize)
				require.Equal(t, tc.expectedDuration, l.MaxRequestDuration)
				require.Equal(t, tc.expectedRequireRequestHeader, l.RequireRequestHeader)
				require.Equal(t, tc.expectedDisableRequestLimiter, l.DisableRequestLimiter)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.MaxRequestSizeRaw)
				require.Nil(t, l.MaxRequestDurationRaw)
				require.Nil(t, l.RequireRequestHeaderRaw)
				require.Nil(t, l.DisableRequestLimiterRaw)
			}
		})
	}
}

// TestListener_parseTLSSettings exercises the listener receiver parseTLSSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseTLSSettings(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		rawTLSDisable                         any
		expectedTLSDisable                    bool
		rawTLSCipherSuites                    string
		expectedTLSCipherSuites               []uint16
		rawTLSRequireAndVerifyClientCert      any
		expectedTLSRequireAndVerifyClientCert bool
		rawTLSDisableClientCerts              any
		expectedTLSDisableClientCerts         bool
		isErrorExpected                       bool
		errorMessage                          string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"tls-disable-bad": {
			rawTLSDisable:   "juan",
			isErrorExpected: true,
			errorMessage:    "invalid value for tls_disable",
		},
		"tls-disable-good": {
			rawTLSDisable:      "true",
			expectedTLSDisable: true,
			isErrorExpected:    false,
		},
		"tls-cipher-suites-bad": {
			rawTLSCipherSuites: "juan",
			isErrorExpected:    true,
			errorMessage:       "invalid value for tls_cipher_suites",
		},
		"tls-cipher-suites-good": {
			rawTLSCipherSuites:      "TLS_RSA_WITH_RC4_128_SHA",
			expectedTLSCipherSuites: []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
			isErrorExpected:         false,
		},
		"tls-require-and-verify-client-cert-bad": {
			rawTLSRequireAndVerifyClientCert: "juan",
			isErrorExpected:                  true,
			errorMessage:                     "invalid value for tls_require_and_verify_client_cert",
		},
		"tls-require-and-verify-client-cert-good": {
			rawTLSRequireAndVerifyClientCert:      "true",
			expectedTLSRequireAndVerifyClientCert: true,
			isErrorExpected:                       false,
		},
		"tls-disable-client-certs-bad": {
			rawTLSDisableClientCerts: "juan",
			isErrorExpected:          true,
			errorMessage:             "invalid value for tls_disable_client_certs",
		},
		"tls-disable-client-certs-good": {
			rawTLSDisableClientCerts:      "true",
			expectedTLSDisableClientCerts: true,
			isErrorExpected:               false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				TLSDisableRaw:                    tc.rawTLSDisable,
				TLSCipherSuitesRaw:               tc.rawTLSCipherSuites,
				TLSRequireAndVerifyClientCertRaw: tc.rawTLSRequireAndVerifyClientCert,
				TLSDisableClientCertsRaw:         tc.rawTLSDisableClientCerts,
			}

			err := l.parseTLSSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Equal(t, tc.expectedTLSDisable, l.TLSDisable)
				require.Equal(t, tc.expectedTLSCipherSuites, l.TLSCipherSuites)
				require.Equal(t, tc.expectedTLSRequireAndVerifyClientCert, l.TLSRequireAndVerifyClientCert)
				require.Equal(t, tc.expectedTLSDisableClientCerts, l.TLSDisableClientCerts)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.TLSDisableRaw)
				require.Empty(t, l.TLSCipherSuitesRaw)
				require.Nil(t, l.TLSRequireAndVerifyClientCertRaw)
				require.Nil(t, l.TLSDisableClientCertsRaw)
			}
		})
	}
}

// TestListener_parseHTTPTimeoutSettings exercises the listener receiver parseHTTPTimeoutSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseHTTPTimeoutSettings(t *testing.T) {
	tests := map[string]struct {
		rawHTTPReadTimeout            any
		expectedHTTPReadTimeout       time.Duration
		rawHTTPReadHeaderTimeout      any
		expectedHTTPReadHeaderTimeout time.Duration
		rawHTTPWriteTimeout           any
		expectedHTTPWriteTimeout      time.Duration
		rawHTTPIdleTimeout            any
		expectedHTTPIdleTimeout       time.Duration
		isErrorExpected               bool
		errorMessage                  string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"read-timeout-bad": {
			rawHTTPReadTimeout: "juan",
			isErrorExpected:    true,
			errorMessage:       "error parsing http_read_timeout",
		},
		"read-timeout-good": {
			rawHTTPReadTimeout:      "30s",
			expectedHTTPReadTimeout: 30 * time.Second,
			isErrorExpected:         false,
		},
		"read-header-timeout-bad": {
			rawHTTPReadHeaderTimeout: "juan",
			isErrorExpected:          true,
			errorMessage:             "error parsing http_read_header_timeout",
		},
		"read-header-timeout-good": {
			rawHTTPReadHeaderTimeout:      "30s",
			expectedHTTPReadHeaderTimeout: 30 * time.Second,
			isErrorExpected:               false,
		},
		"write-timeout-bad": {
			rawHTTPWriteTimeout: "juan",
			isErrorExpected:     true,
			errorMessage:        "error parsing http_write_timeout",
		},
		"write-timeout-good": {
			rawHTTPWriteTimeout:      "30s",
			expectedHTTPWriteTimeout: 30 * time.Second,
			isErrorExpected:          false,
		},
		"idle-timeout-bad": {
			rawHTTPIdleTimeout: "juan",
			isErrorExpected:    true,
			errorMessage:       "error parsing http_idle_timeout",
		},
		"idle-timeout-good": {
			rawHTTPIdleTimeout:      "30s",
			expectedHTTPIdleTimeout: 30 * time.Second,
			isErrorExpected:         false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				HTTPReadTimeoutRaw:       tc.rawHTTPReadTimeout,
				HTTPReadHeaderTimeoutRaw: tc.rawHTTPReadHeaderTimeout,
				HTTPWriteTimeoutRaw:      tc.rawHTTPWriteTimeout,
				HTTPIdleTimeoutRaw:       tc.rawHTTPIdleTimeout,
			}

			err := l.parseHTTPTimeoutSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Equal(t, tc.expectedHTTPReadTimeout, l.HTTPReadTimeout)
				require.Equal(t, tc.expectedHTTPReadHeaderTimeout, l.HTTPReadHeaderTimeout)
				require.Equal(t, tc.expectedHTTPWriteTimeout, l.HTTPWriteTimeout)
				require.Equal(t, tc.expectedHTTPIdleTimeout, l.HTTPIdleTimeout)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.HTTPReadTimeoutRaw)
				require.Nil(t, l.HTTPReadHeaderTimeoutRaw)
				require.Nil(t, l.HTTPWriteTimeoutRaw)
				require.Nil(t, l.HTTPIdleTimeoutRaw)
			}
		})
	}
}

// TestListener_parseProxySettings exercises the listener receiver parseProxySettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseProxySettings(t *testing.T) {
	tests := map[string]struct {
		rawProxyProtocolAuthorizedAddrs any
		expectedNumAddrs                int
		proxyBehavior                   string
		isErrorExpected                 bool
		errorMessage                    string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"bad-addrs": {
			rawProxyProtocolAuthorizedAddrs: "juan",
			isErrorExpected:                 true,
			errorMessage:                    "error parsing proxy_protocol_authorized_addrs",
		},
		"good-addrs": {
			rawProxyProtocolAuthorizedAddrs: "10.0.0.1,10.0.2.1",
			expectedNumAddrs:                2,
			proxyBehavior:                   "",
			isErrorExpected:                 false,
		},
		"behavior-bad": {
			rawProxyProtocolAuthorizedAddrs: "10.0.0.1,10.0.2.1",
			proxyBehavior:                   "juan",
			isErrorExpected:                 true,
			errorMessage:                    "unsupported value supplied for proxy_protocol_behavior",
		},
		"behavior-use-always": {
			rawProxyProtocolAuthorizedAddrs: "10.0.0.1,10.0.2.1",
			expectedNumAddrs:                2,
			proxyBehavior:                   "use_always",
			isErrorExpected:                 false,
		},
		"behavior-empty": {
			rawProxyProtocolAuthorizedAddrs: "10.0.0.1,10.0.2.1",
			expectedNumAddrs:                2,
			proxyBehavior:                   "",
			isErrorExpected:                 false,
		},
		"behavior-allow": {
			rawProxyProtocolAuthorizedAddrs: "10.0.0.1,10.0.2.1",
			expectedNumAddrs:                2,
			proxyBehavior:                   "allow_authorized",
			isErrorExpected:                 false,
		},
		"behavior-deny": {
			rawProxyProtocolAuthorizedAddrs: "10.0.0.1,10.0.2.1",
			expectedNumAddrs:                2,
			proxyBehavior:                   "deny_unauthorized",
			isErrorExpected:                 false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				ProxyProtocolAuthorizedAddrsRaw: tc.rawProxyProtocolAuthorizedAddrs,
				ProxyProtocolBehavior:           tc.proxyBehavior,
			}

			err := l.parseProxySettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Len(t, l.ProxyProtocolAuthorizedAddrs, tc.expectedNumAddrs)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.ProxyProtocolAuthorizedAddrsRaw)
			}
		})
	}
}

// TestListener_parseForwardedForSettings exercises the listener receiver parseForwardedForSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseForwardedForSettings(t *testing.T) {
	tests := map[string]struct {
		rawAuthorizedAddrs          any
		expectedNumAddrs            int
		rawHopSkips                 any
		expectedHopSkips            int64
		rawRejectNotAuthorized      any
		expectedRejectNotAuthorized bool
		rawRejectNotPresent         any
		expectedRejectNotPresent    bool
		isErrorExpected             bool
		errorMessage                string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"authorized-addrs-bad": {
			rawAuthorizedAddrs: "juan",
			isErrorExpected:    true,
			errorMessage:       "error parsing x_forwarded_for_authorized_addrs",
		},
		"authorized-addrs-good": {
			rawAuthorizedAddrs: "10.0.0.1,10.0.2.1",
			expectedNumAddrs:   2,
			isErrorExpected:    false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				XForwardedForAuthorizedAddrsRaw:     tc.rawAuthorizedAddrs,
				XForwardedForHopSkipsRaw:            tc.rawHopSkips,
				XForwardedForRejectNotAuthorizedRaw: tc.rawRejectNotAuthorized,
				XForwardedForRejectNotPresentRaw:    tc.rawRejectNotPresent,
			}

			err := l.parseForwardedForSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)

				require.Len(t, l.XForwardedForAuthorizedAddrs, tc.expectedNumAddrs)
				require.Equal(t, tc.expectedHopSkips, l.XForwardedForHopSkips)
				require.Equal(t, tc.expectedRejectNotAuthorized, l.XForwardedForRejectNotAuthorized)
				require.Equal(t, tc.expectedRejectNotPresent, l.XForwardedForRejectNotPresent)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.XForwardedForAuthorizedAddrsRaw)
				require.Nil(t, l.XForwardedForHopSkipsRaw)
				require.Nil(t, l.XForwardedForRejectNotAuthorizedRaw)
				require.Nil(t, l.XForwardedForRejectNotPresentRaw)
			}
		})
	}
}

// TestListener_parseTelemetrySettings exercises the listener receiver parseTelemetrySettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseTelemetrySettings(t *testing.T) {
	tests := map[string]struct {
		rawUnauthenticatedMetricsAccess      any
		expectedUnauthenticatedMetricsAccess bool
		isErrorExpected                      bool
		errorMessage                         string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"unauth-bad": {
			rawUnauthenticatedMetricsAccess: "juan",
			isErrorExpected:                 true,
			errorMessage:                    "invalid value for telemetry.unauthenticated_metrics_access",
		},
		"unauth-good": {
			rawUnauthenticatedMetricsAccess:      "true",
			expectedUnauthenticatedMetricsAccess: true,
			isErrorExpected:                      false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				Telemetry: ListenerTelemetry{
					UnauthenticatedMetricsAccessRaw: tc.rawUnauthenticatedMetricsAccess,
				},
			}

			err := l.parseTelemetrySettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Equal(t, tc.expectedUnauthenticatedMetricsAccess, l.Telemetry.UnauthenticatedMetricsAccess)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.Telemetry.UnauthenticatedMetricsAccessRaw)
			}
		})
	}
}

// TestListener_parseProfilingSettings exercises the listener receiver parseProfilingSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseProfilingSettings(t *testing.T) {
	tests := map[string]struct {
		rawUnauthenticatedPProfAccess      any
		expectedUnauthenticatedPProfAccess bool
		isErrorExpected                    bool
		errorMessage                       string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"bad": {
			rawUnauthenticatedPProfAccess: "juan",
			isErrorExpected:               true,
			errorMessage:                  "invalid value for profiling.unauthenticated_pprof_access",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				Profiling: ListenerProfiling{
					UnauthenticatedPProfAccessRaw: tc.rawUnauthenticatedPProfAccess,
				},
			}

			err := l.parseProfilingSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Equal(t, tc.expectedUnauthenticatedPProfAccess, l.Profiling.UnauthenticatedPProfAccess)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.Profiling.UnauthenticatedPProfAccessRaw)
			}
		})
	}
}

// TestListener_parseInFlightRequestSettings exercises the listener receiver parseInFlightRequestSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseInFlightRequestSettings(t *testing.T) {
	tests := map[string]struct {
		rawUnauthenticatedInFlightAccess      any
		expectedUnauthenticatedInFlightAccess bool
		isErrorExpected                       bool
		errorMessage                          string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"bad": {
			rawUnauthenticatedInFlightAccess: "juan",
			isErrorExpected:                  true,
			errorMessage:                     "invalid value for inflight_requests_logging.unauthenticated_in_flight_requests_access",
		},
		"good": {
			rawUnauthenticatedInFlightAccess:      "true",
			expectedUnauthenticatedInFlightAccess: true,
			isErrorExpected:                       false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				InFlightRequestLogging: ListenerInFlightRequestLogging{
					UnauthenticatedInFlightAccessRaw: tc.rawUnauthenticatedInFlightAccess,
				},
			}

			err := l.parseInFlightRequestSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Equal(t, tc.expectedUnauthenticatedInFlightAccess, l.InFlightRequestLogging.UnauthenticatedInFlightAccess)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.InFlightRequestLogging.UnauthenticatedInFlightAccessRaw)
			}
		})
	}
}

// TestListener_parseCORSSettings exercises the listener receiver parseCORSSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseCORSSettings(t *testing.T) {
	tests := map[string]struct {
		rawCorsEnabled                any
		rawCorsAllowedHeaders         []string
		corsAllowedOrigins            []string
		expectedCorsEnabled           bool
		expectedNumCorsAllowedHeaders int
		isErrorExpected               bool
		errorMessage                  string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"cors-enabled-bad": {
			rawCorsEnabled:      "juan",
			expectedCorsEnabled: false,
			isErrorExpected:     true,
			errorMessage:        "invalid value for cors_enabled",
		},
		"cors-enabled-good": {
			rawCorsEnabled:      "true",
			expectedCorsEnabled: true,
			isErrorExpected:     false,
		},
		"cors-allowed-origins-single-wildcard": {
			corsAllowedOrigins: []string{"*"},
			isErrorExpected:    false,
		},
		"cors-allowed-origins-multi-wildcard": {
			corsAllowedOrigins: []string{"*", "hashicorp.com"},
			isErrorExpected:    true,
			errorMessage:       "cors_allowed_origins must only contain a wildcard or only non-wildcard values",
		},
		"cors-allowed-headers-anything": {
			rawCorsAllowedHeaders:         []string{"foo", "bar"},
			expectedNumCorsAllowedHeaders: 2,
			isErrorExpected:               false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				CorsEnabledRaw:        tc.rawCorsEnabled,
				CorsAllowedHeadersRaw: tc.rawCorsAllowedHeaders,
				CorsAllowedOrigins:    tc.corsAllowedOrigins,
			}

			err := l.parseCORSSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Equal(t, tc.expectedCorsEnabled, l.CorsEnabled)
				require.Len(t, l.CorsAllowedHeaders, tc.expectedNumCorsAllowedHeaders)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.CorsEnabledRaw)
				require.Nil(t, l.CorsAllowedHeadersRaw)
			}
		})
	}
}

// TestListener_parseHTTPHeaderSettings exercises the listener receiver parseHTTPHeaderSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseHTTPHeaderSettings(t *testing.T) {
	tests := map[string]struct {
		listenerType                     ListenerType
		rawCustomResponseHeaders         []map[string]any
		expectedNumCustomResponseHeaders int
		isErrorExpected                  bool
		errorMessage                     string
	}{
		"nil": {
			listenerType:                     TCP,
			isErrorExpected:                  false,
			expectedNumCustomResponseHeaders: 1, // default: Strict-Transport-Security
		},
		"custom-headers-bad": {
			listenerType: TCP,
			rawCustomResponseHeaders: []map[string]any{
				{"juan": false},
			},
			isErrorExpected: true,
			errorMessage:    "failed to parse custom_response_headers",
		},
		"custom-headers-good": {
			listenerType: TCP,
			rawCustomResponseHeaders: []map[string]any{
				{
					"2xx": []map[string]any{
						{"X-Custom-Header": []any{"Custom Header Value 1", "Custom Header Value 2"}},
					},
				},
			},
			expectedNumCustomResponseHeaders: 2,
			isErrorExpected:                  false,
		},
		"unix-no-headers": {
			listenerType: Unix,
			rawCustomResponseHeaders: []map[string]any{
				{
					"2xx": []map[string]any{
						{"X-Custom-Header": []any{"Custom Header Value 1", "Custom Header Value 2"}},
					},
				},
			},
			expectedNumCustomResponseHeaders: 0,
			isErrorExpected:                  false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				Type:                     tc.listenerType,
				CustomResponseHeadersRaw: tc.rawCustomResponseHeaders,
			}

			err := l.parseHTTPHeaderSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Len(t, l.CustomResponseHeaders, tc.expectedNumCustomResponseHeaders)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.CustomResponseHeadersRaw)
			}
		})
	}
}

// TestListener_parseChrootNamespaceSettings exercises the listener receiver parseChrootNamespaceSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseChrootNamespaceSettings(t *testing.T) {
	tests := map[string]struct {
		rawChrootNamespace      any
		expectedChrootNamespace string
		isErrorExpected         bool
		errorMessage            string
	}{
		"nil": {
			isErrorExpected: false,
		},
		"bad": {
			rawChrootNamespace: &Listener{}, // Unsure how we'd ever see this really.
			isErrorExpected:    true,
			errorMessage:       "invalid value for chroot_namespace",
		},
		"good": {
			rawChrootNamespace:      "juan",
			expectedChrootNamespace: "juan/",
			isErrorExpected:         false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				ChrootNamespaceRaw: tc.rawChrootNamespace,
			}

			err := l.parseChrootNamespaceSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Equal(t, tc.expectedChrootNamespace, l.ChrootNamespace)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.ChrootNamespaceRaw)
			}
		})
	}
}

// TestListener_parseRedactionSettings exercises the listener receiver parseRedactionSettings.
// We check various inputs to ensure we can parse the values as expected and
// assign the relevant value on the SharedConfig struct.
func TestListener_parseRedactionSettings(t *testing.T) {
	tests := map[string]struct {
		listenerType              ListenerType
		rawRedactAddresses        any
		expectedRedactAddresses   bool
		rawRedactClusterName      any
		expectedRedactClusterName bool
		rawRedactVersion          any
		expectedRedactVersion     bool
		isErrorExpected           bool
		errorMessage              string
	}{
		"missing": {
			listenerType:              TCP,
			isErrorExpected:           false,
			expectedRedactAddresses:   false,
			expectedRedactClusterName: false,
			expectedRedactVersion:     false,
		},
		"redact-addresses-bad": {
			listenerType:       TCP,
			rawRedactAddresses: "juan",
			isErrorExpected:    true,
			errorMessage:       "invalid value for redact_addresses",
		},
		"redact-addresses-good": {
			listenerType:            TCP,
			rawRedactAddresses:      "true",
			expectedRedactAddresses: true,
			isErrorExpected:         false,
		},
		"redact-cluster-name-bad": {
			listenerType:         TCP,
			rawRedactClusterName: "juan",
			isErrorExpected:      true,
			errorMessage:         "invalid value for redact_cluster_name",
		},
		"redact-cluster-name-good": {
			listenerType:              TCP,
			rawRedactClusterName:      "true",
			expectedRedactClusterName: true,
			isErrorExpected:           false,
		},
		"redact-version-bad": {
			listenerType:     TCP,
			rawRedactVersion: "juan",
			isErrorExpected:  true,
			errorMessage:     "invalid value for redact_version",
		},
		"redact-version-good": {
			listenerType:          TCP,
			rawRedactVersion:      "true",
			expectedRedactVersion: true,
			isErrorExpected:       false,
		},
		"redact-unix-na": {
			listenerType:              Unix,
			rawRedactAddresses:        "true",
			expectedRedactAddresses:   false,
			rawRedactClusterName:      "true",
			expectedRedactClusterName: false,
			rawRedactVersion:          "true",
			expectedRedactVersion:     false,
			isErrorExpected:           false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Configure listener with raw values
			l := &Listener{
				Type:                 tc.listenerType,
				RedactAddressesRaw:   tc.rawRedactAddresses,
				RedactClusterNameRaw: tc.rawRedactClusterName,
				RedactVersionRaw:     tc.rawRedactVersion,
			}

			err := l.parseRedactionSettings()

			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			default:
				// Assert we got the relevant values.
				require.NoError(t, err)
				require.Equal(t, tc.expectedRedactAddresses, l.RedactAddresses)
				require.Equal(t, tc.expectedRedactClusterName, l.RedactClusterName)
				require.Equal(t, tc.expectedRedactVersion, l.RedactVersion)

				// Ensure the state was modified for the raw values.
				require.Nil(t, l.RedactAddressesRaw)
				require.Nil(t, l.RedactClusterNameRaw)
				require.Nil(t, l.RedactVersionRaw)
			}
		})
	}
}

func TestParseAndClearBool(t *testing.T) {
	testcases := []struct {
		name           string
		raw            interface{}
		rawAssertion   func(assert.TestingT, any, ...any) bool
		expectedParsed bool
		errorAssertion func(assert.TestingT, error, ...any) bool
	}{
		{
			name:           "valid-true-as-string",
			raw:            "true",
			rawAssertion:   assert.Nil,
			expectedParsed: true,
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-false-as-string",
			raw:            "false",
			rawAssertion:   assert.Nil,
			expectedParsed: false,
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-true-as-bool",
			raw:            true,
			rawAssertion:   assert.Nil,
			expectedParsed: true,
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-false-as-bool",
			raw:            false,
			rawAssertion:   assert.Nil,
			expectedParsed: false,
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-true-as-string-mix-case",
			raw:            "True",
			rawAssertion:   assert.Nil,
			expectedParsed: true,
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-false-as-integer",
			raw:            0,
			rawAssertion:   assert.Nil,
			expectedParsed: false,
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-true-as-integer",
			raw:            2,
			rawAssertion:   assert.Nil,
			expectedParsed: true,
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-true-as-float",
			raw:            3.14,
			rawAssertion:   assert.Nil,
			expectedParsed: true,
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-false-as-float",
			raw:            0.0,
			rawAssertion:   assert.Nil,
			expectedParsed: false,
			errorAssertion: assert.NoError,
		},
		{
			name:           "invalid-as-string",
			raw:            "0.0.0.0:8200",
			rawAssertion:   assert.NotNil,
			errorAssertion: assert.Error,
		},
		{
			name:           "invalid-as-struct",
			raw:            struct{}{},
			rawAssertion:   assert.NotNil,
			errorAssertion: assert.Error,
		},
		{
			name:           "not-set",
			raw:            nil,
			rawAssertion:   assert.Nil,
			errorAssertion: assert.NoError,
		},
	}

	for _, testcase := range testcases {
		var parsed bool
		err := parseAndClearBool(&testcase.raw, &parsed)

		testcase.errorAssertion(t, err, testcase.name)
		assert.Equal(t, testcase.expectedParsed, parsed, testcase.name)
		testcase.rawAssertion(t, testcase.raw, testcase.name)
	}
}

func TestParseAndClearString(t *testing.T) {
	testcases := []struct {
		name           string
		raw            any
		rawAssertion   func(assert.TestingT, any, ...any) bool
		expectedParsed string
		errorAssertion func(assert.TestingT, error, ...any) bool
	}{
		{
			name:           "valid-empty-string",
			raw:            "",
			rawAssertion:   assert.Nil,
			expectedParsed: "",
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-some-string",
			raw:            "blah blah",
			rawAssertion:   assert.Nil,
			expectedParsed: "blah blah",
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-as-integer",
			raw:            8,
			rawAssertion:   assert.Nil,
			expectedParsed: "8",
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-as-bool",
			raw:            true,
			rawAssertion:   assert.Nil,
			expectedParsed: "1",
			errorAssertion: assert.NoError,
		},
		{
			name:           "not-set",
			raw:            nil,
			rawAssertion:   assert.Nil,
			expectedParsed: "",
			errorAssertion: assert.NoError,
		},
		{
			name:           "invalid-as-struct",
			raw:            struct{}{},
			rawAssertion:   assert.NotNil,
			errorAssertion: assert.Error,
		},
	}
	for _, testcase := range testcases {
		var parsed string
		err := parseAndClearString(&testcase.raw, &parsed)

		testcase.errorAssertion(t, err, testcase.name)
		assert.Equal(t, testcase.expectedParsed, parsed, testcase.name)
		testcase.rawAssertion(t, testcase.raw, testcase.name)
	}
}

func TestParseAndClearInt(t *testing.T) {
	testcases := []struct {
		name           string
		raw            any
		rawAssertion   func(assert.TestingT, any, ...any) bool
		expectedParsed int64
		errorAssertion func(assert.TestingT, error, ...any) bool
	}{
		{
			name:           "valid-as-int",
			raw:            200,
			rawAssertion:   assert.Nil,
			expectedParsed: int64(200),
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-as-string",
			raw:            "53",
			rawAssertion:   assert.Nil,
			expectedParsed: int64(53),
			errorAssertion: assert.NoError,
		},
		{
			name:           "invalid-as-hex-string",
			raw:            "0xa",
			rawAssertion:   assert.NotNil,
			errorAssertion: assert.Error,
		},
		{
			name:           "not-set",
			raw:            nil,
			rawAssertion:   assert.Nil,
			errorAssertion: assert.NoError,
		},
	}

	for _, testcase := range testcases {
		var parsed int64
		err := parseAndClearInt(&testcase.raw, &parsed)

		testcase.errorAssertion(t, err, testcase.name)
		assert.Equal(t, testcase.expectedParsed, parsed, testcase.name)
		testcase.rawAssertion(t, testcase.raw, testcase.name)
	}
}

func TestParseAndClearDurationSecond(t *testing.T) {
	testcases := []struct {
		name           string
		raw            any
		rawAssertion   func(assert.TestingT, any, ...any) bool
		expectedParsed time.Duration
		errorAssertion func(assert.TestingT, error, ...any) bool
	}{
		{
			name:           "valid-as-string",
			raw:            "30s",
			rawAssertion:   assert.Nil,
			expectedParsed: time.Duration(30 * time.Second),
			errorAssertion: assert.NoError,
		},
		{
			name:           "valid-as-string-more-complex",
			raw:            "29h24m49s",
			rawAssertion:   assert.Nil,
			expectedParsed: time.Duration((29 * time.Hour) + (24 * time.Minute) + (49 * time.Second)),
			errorAssertion: assert.NoError,
		},
		{
			name:           "invalid-as-string-using-days",
			raw:            "1d3s",
			rawAssertion:   assert.NotNil,
			errorAssertion: assert.Error,
		},
		{
			name:           "valid-as-integer",
			raw:            87,
			rawAssertion:   assert.Nil,
			expectedParsed: time.Duration(87 * time.Second),
			errorAssertion: assert.NoError,
		},
		{
			name:           "not-set",
			raw:            nil,
			rawAssertion:   assert.Nil,
			errorAssertion: assert.NoError,
		},
		{
			name:           "invalid-as-struct",
			raw:            struct{}{},
			rawAssertion:   assert.NotNil,
			errorAssertion: assert.Error,
		},
	}

	for _, testcase := range testcases {
		var parsed time.Duration

		err := parseAndClearDurationSecond(&testcase.raw, &parsed)
		testcase.errorAssertion(t, err, testcase.name)
		assert.Equal(t, testcase.expectedParsed, parsed)
		testcase.rawAssertion(t, testcase.raw, testcase.name)
	}
}
