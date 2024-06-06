// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

func getListenerConfigForMarshalerTest(addr sockaddr.IPAddr) *configutil.Listener {
	return &configutil.Listener{
		XForwardedForAuthorizedAddrs: []*sockaddr.SockAddrMarshaler{
			{
				SockAddr: addr,
			},
		},
	}
}

func TestHandler_XForwardedFor(t *testing.T) {
	goodAddr, err := sockaddr.NewIPAddr("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	badAddr, err := sockaddr.NewIPAddr("1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}

	// First: test reject not present
	t.Run("reject_not_present", func(t *testing.T) {
		t.Parallel()
		testHandler := func(props *vault.HandlerProperties) http.Handler {
			origHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.RemoteAddr))
			})
			listenerConfig := getListenerConfigForMarshalerTest(goodAddr)
			listenerConfig.XForwardedForRejectNotPresent = true
			return WrapForwardedForHandler(origHandler, listenerConfig)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: HandlerFunc(testHandler),
		})
		cluster.Start()
		defer cluster.Cleanup()
		client := cluster.Cores[0].Client

		req := client.NewRequest("GET", "/")
		_, err := client.RawRequest(req)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "missing x-forwarded-for") {
			t.Fatalf("bad error message: %v", err)
		}
		req = client.NewRequest("GET", "/")
		req.Headers = make(http.Header)
		req.Headers.Set("x-forwarded-for", "1.2.3.4")
		resp, err := client.RawRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		buf := bytes.NewBuffer(nil)
		buf.ReadFrom(resp.Body)
		if !strings.HasPrefix(buf.String(), "1.2.3.4:") {
			t.Fatalf("bad body: %s", buf.String())
		}
	})

	// Next: test allow unauth
	t.Run("allow_unauth", func(t *testing.T) {
		t.Parallel()
		testHandler := func(props *vault.HandlerProperties) http.Handler {
			origHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.RemoteAddr))
			})
			listenerConfig := getListenerConfigForMarshalerTest(badAddr)
			listenerConfig.XForwardedForRejectNotPresent = true
			return WrapForwardedForHandler(origHandler, listenerConfig)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: HandlerFunc(testHandler),
		})
		cluster.Start()
		defer cluster.Cleanup()
		client := cluster.Cores[0].Client

		req := client.NewRequest("GET", "/")
		req.Headers = make(http.Header)
		req.Headers.Set("x-forwarded-for", "5.6.7.8")
		resp, err := client.RawRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		buf := bytes.NewBuffer(nil)
		buf.ReadFrom(resp.Body)
		if !strings.HasPrefix(buf.String(), "127.0.0.1:") {
			t.Fatalf("bad body: %s", buf.String())
		}
	})

	// Next: test fail unauth
	t.Run("fail_unauth", func(t *testing.T) {
		t.Parallel()
		testHandler := func(props *vault.HandlerProperties) http.Handler {
			origHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.RemoteAddr))
			})
			listenerConfig := getListenerConfigForMarshalerTest(badAddr)
			listenerConfig.XForwardedForRejectNotPresent = true
			listenerConfig.XForwardedForRejectNotAuthorized = true
			return WrapForwardedForHandler(origHandler, listenerConfig)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: HandlerFunc(testHandler),
		})
		cluster.Start()
		defer cluster.Cleanup()
		client := cluster.Cores[0].Client

		req := client.NewRequest("GET", "/")
		req.Headers = make(http.Header)
		req.Headers.Set("x-forwarded-for", "5.6.7.8")
		_, err := client.RawRequest(req)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "not authorized for x-forwarded-for") {
			t.Fatalf("bad error message: %v", err)
		}
	})

	// Next: test bad hops (too many)
	t.Run("too_many_hops", func(t *testing.T) {
		t.Parallel()
		testHandler := func(props *vault.HandlerProperties) http.Handler {
			origHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.RemoteAddr))
			})
			listenerConfig := getListenerConfigForMarshalerTest(goodAddr)
			listenerConfig.XForwardedForRejectNotPresent = true
			listenerConfig.XForwardedForRejectNotAuthorized = true
			listenerConfig.XForwardedForHopSkips = 4
			return WrapForwardedForHandler(origHandler, listenerConfig)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: HandlerFunc(testHandler),
		})
		cluster.Start()
		defer cluster.Cleanup()
		client := cluster.Cores[0].Client

		req := client.NewRequest("GET", "/")
		req.Headers = make(http.Header)
		req.Headers.Set("x-forwarded-for", "2.3.4.5,3.4.5.6")
		_, err := client.RawRequest(req)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "would skip before earliest") {
			t.Fatalf("bad error message: %v", err)
		}
	})

	// Next: test picking correct value
	t.Run("correct_hop_skipping", func(t *testing.T) {
		t.Parallel()
		testHandler := func(props *vault.HandlerProperties) http.Handler {
			origHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.RemoteAddr))
			})
			listenerConfig := getListenerConfigForMarshalerTest(goodAddr)
			listenerConfig.XForwardedForRejectNotPresent = true
			listenerConfig.XForwardedForRejectNotAuthorized = true
			listenerConfig.XForwardedForHopSkips = 1
			return WrapForwardedForHandler(origHandler, listenerConfig)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: HandlerFunc(testHandler),
		})
		cluster.Start()
		defer cluster.Cleanup()
		client := cluster.Cores[0].Client

		req := client.NewRequest("GET", "/")
		req.Headers = make(http.Header)
		req.Headers.Set("x-forwarded-for", "2.3.4.5,3.4.5.6,4.5.6.7,5.6.7.8")
		resp, err := client.RawRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		buf := bytes.NewBuffer(nil)
		buf.ReadFrom(resp.Body)
		if !strings.HasPrefix(buf.String(), "4.5.6.7:") {
			t.Fatalf("bad body: %s", buf.String())
		}
	})

	// Next: multi-header approach
	t.Run("correct_hop_skipping_multi_header", func(t *testing.T) {
		t.Parallel()
		testHandler := func(props *vault.HandlerProperties) http.Handler {
			origHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.RemoteAddr))
			})
			listenerConfig := getListenerConfigForMarshalerTest(goodAddr)
			listenerConfig.XForwardedForRejectNotPresent = true
			listenerConfig.XForwardedForRejectNotAuthorized = true
			listenerConfig.XForwardedForHopSkips = 1
			return WrapForwardedForHandler(origHandler, listenerConfig)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: HandlerFunc(testHandler),
		})
		cluster.Start()
		defer cluster.Cleanup()
		client := cluster.Cores[0].Client

		req := client.NewRequest("GET", "/")
		req.Headers = make(http.Header)
		req.Headers.Add("x-forwarded-for", "2.3.4.5")
		req.Headers.Add("x-forwarded-for", "3.4.5.6,4.5.6.7")
		req.Headers.Add("x-forwarded-for", "5.6.7.8")
		resp, err := client.RawRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		buf := bytes.NewBuffer(nil)
		buf.ReadFrom(resp.Body)
		if !strings.HasPrefix(buf.String(), "4.5.6.7:") {
			t.Fatalf("bad body: %s", buf.String())
		}
	})

	// Next: test an invalid certificate being sent
	t.Run("reject_bad_cert_in_header", func(t *testing.T) {
		t.Parallel()
		testHandler := func(props *vault.HandlerProperties) http.Handler {
			origHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.RemoteAddr))
			})
			listenerConfig := getListenerConfigForMarshalerTest(goodAddr)
			listenerConfig.XForwardedForClientCertHeader = "X-Forwarded-Tls-Client-Cert"
			listenerConfig.XForwardedForClientCertHeaderDecoders = "URL,BASE64"
			return WrapForwardedForHandler(origHandler, listenerConfig)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: HandlerFunc(testHandler),
		})
		cluster.Start()
		defer cluster.Cleanup()
		client := cluster.Cores[0].Client

		req := client.NewRequest("GET", "/")
		req.Headers = make(http.Header)
		req.Headers.Set("x-forwarded-for", "5.6.7.8")
		req.Headers.Set("x-forwarded-tls-client-cert", `BAD_TEXTMIIDtTCCAp2gAwIBAgIUf%2BjhKTFBnqSs34II0WS1L4QsbbAwDQYJKoZIhvcNAQEL%0ABQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMTYwMjI5MDIyNzQxWhcNMjUw%0AMTA1MTAyODExWjAbMRkwFwYDVQQDExBjZXJ0LmV4YW1wbGUuY29tMIIBIjANBgkq%0AhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsZx0Svr82YJpFpIy4fJNW5fKA6B8mhxS%0ATRAVnygAftetT8puHflY0ss7Y6X2OXjsU0PRn%2B1PswtivhKi%2BeLtgWkUF9cFYFGn%0ASgMld6ZWRhNheZhA6ZfQmeM%2FBF2pa5HK2SDF36ljgjL9T%2BnWrru2Uv0BCoHzLAmi%0AYYMiIWplidMmMO5NTRG3k%2B3AN0TkfakB6JVzjLGhTcXdOcVEMXkeQVqJMAuGouU5%0AdonyqtnaHuIJGuUdy54YDnX86txhOQhAv6r7dHXzZxS4pmLvw8UI1rsSf%2FGLcUVG%0AB%2B5%2BAAGF5iuHC3N2DTl4xz3FcN4Cb4w9pbaQ7%2BmCzz%2BanqiJfyr2nwIDAQABo4H1%0AMIHyMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQUm%2B%2Be%0AHpyM3p708bgZJuRYEdX1o%2BUwHwYDVR0jBBgwFoAUncSzT%2F6HMexyuiU9%2F7EgHu%2Bo%0Ak5swOwYIKwYBBQUHAQEELzAtMCsGCCsGAQUFBzAChh9odHRwOi8vMTI3LjAuMC4x%0AOjgyMDAvdjEvcGtpL2NhMCEGA1UdEQQaMBiCEGNlcnQuZXhhbXBsZS5jb22HBH8A%0AAAEwMQYDVR0fBCowKDAmoCSgIoYgaHR0cDovLzEyNy4wLjAuMTo4MjAwL3YxL3Br%0AaS9jcmwwDQYJKoZIhvcNAQELBQADggEBABsuvmPSNjjKTVN6itWzdQy%2BSgMIrwfs%0AX1Yb9Lefkkwmp9ovKFNQxa4DucuCuzXcQrbKwWTfHGgR8ct4rf30xCRoA7dbQWq4%0AaYqNKFWrRaBRAaaYZ%2FO1ApRTOrXqRx9Eqr0H1BXLsoAq%2BmWassL8sf6siae%2BCpwA%0AKqBko5G0dNXq5T4i2LQbmoQSVetIrCJEeMrU%2BidkuqfV2h1BQKgSEhFDABjFdTCN%0AQDAHsEHsi2M4%2FjRW9fqEuhHSDfl2n7tkFUI8wTHUUCl7gXwweJ4qtaSXIwKXYzNj%0AxqKHA8Purc1Yfybz4iE1JCROi9fInKlzr5xABq8nb9Qc%2FJ9DIQM%2BXmk%3D`)
		resp, err := client.RawRequest(req)
		if err == nil {
			t.Fatal("expected error")
		}
		defer resp.Body.Close()
		buf := bytes.NewBuffer(nil)
		buf.ReadFrom(resp.Body)
		if !strings.Contains(buf.String(), "failed to base64 decode the client certificate: ") {
			t.Fatalf("bad body: %v", buf.String())
		}
	})

	// Next: test a valid (unverified) certificate being sent
	t.Run("pass_cert", func(t *testing.T) {
		t.Parallel()
		testHandler := func(props *vault.HandlerProperties) http.Handler {
			origHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(base64.StdEncoding.EncodeToString(r.TLS.PeerCertificates[0].Raw)))
			})
			listenerConfig := getListenerConfigForMarshalerTest(goodAddr)
			listenerConfig.XForwardedForClientCertHeader = "X-Forwarded-Tls-Client-Cert"
			listenerConfig.XForwardedForClientCertHeaderDecoders = "URL,BASE64"
			return WrapForwardedForHandler(origHandler, listenerConfig)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: HandlerFunc(testHandler),
		})
		cluster.Start()
		defer cluster.Cleanup()
		client := cluster.Cores[0].Client

		req := client.NewRequest("GET", "/")
		req.Headers = make(http.Header)
		req.Headers.Set("x-forwarded-for", "5.6.7.8")
		testcertificate := `MIIDtTCCAp2gAwIBAgIUf%2BjhKTFBnqSs34II0WS1L4QsbbAwDQYJKoZIhvcNAQEL%0ABQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMTYwMjI5MDIyNzQxWhcNMjUw%0AMTA1MTAyODExWjAbMRkwFwYDVQQDExBjZXJ0LmV4YW1wbGUuY29tMIIBIjANBgkq%0AhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsZx0Svr82YJpFpIy4fJNW5fKA6B8mhxS%0ATRAVnygAftetT8puHflY0ss7Y6X2OXjsU0PRn%2B1PswtivhKi%2BeLtgWkUF9cFYFGn%0ASgMld6ZWRhNheZhA6ZfQmeM%2FBF2pa5HK2SDF36ljgjL9T%2BnWrru2Uv0BCoHzLAmi%0AYYMiIWplidMmMO5NTRG3k%2B3AN0TkfakB6JVzjLGhTcXdOcVEMXkeQVqJMAuGouU5%0AdonyqtnaHuIJGuUdy54YDnX86txhOQhAv6r7dHXzZxS4pmLvw8UI1rsSf%2FGLcUVG%0AB%2B5%2BAAGF5iuHC3N2DTl4xz3FcN4Cb4w9pbaQ7%2BmCzz%2BanqiJfyr2nwIDAQABo4H1%0AMIHyMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQUm%2B%2Be%0AHpyM3p708bgZJuRYEdX1o%2BUwHwYDVR0jBBgwFoAUncSzT%2F6HMexyuiU9%2F7EgHu%2Bo%0Ak5swOwYIKwYBBQUHAQEELzAtMCsGCCsGAQUFBzAChh9odHRwOi8vMTI3LjAuMC4x%0AOjgyMDAvdjEvcGtpL2NhMCEGA1UdEQQaMBiCEGNlcnQuZXhhbXBsZS5jb22HBH8A%0AAAEwMQYDVR0fBCowKDAmoCSgIoYgaHR0cDovLzEyNy4wLjAuMTo4MjAwL3YxL3Br%0AaS9jcmwwDQYJKoZIhvcNAQELBQADggEBABsuvmPSNjjKTVN6itWzdQy%2BSgMIrwfs%0AX1Yb9Lefkkwmp9ovKFNQxa4DucuCuzXcQrbKwWTfHGgR8ct4rf30xCRoA7dbQWq4%0AaYqNKFWrRaBRAaaYZ%2FO1ApRTOrXqRx9Eqr0H1BXLsoAq%2BmWassL8sf6siae%2BCpwA%0AKqBko5G0dNXq5T4i2LQbmoQSVetIrCJEeMrU%2BidkuqfV2h1BQKgSEhFDABjFdTCN%0AQDAHsEHsi2M4%2FjRW9fqEuhHSDfl2n7tkFUI8wTHUUCl7gXwweJ4qtaSXIwKXYzNj%0AxqKHA8Purc1Yfybz4iE1JCROi9fInKlzr5xABq8nb9Qc%2FJ9DIQM%2BXmk%3D`
		req.Headers.Set("x-forwarded-tls-client-cert", testcertificate)
		resp, err := client.RawRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		buf := bytes.NewBuffer(nil)
		buf.ReadFrom(resp.Body)
		testcertificate, _ = url.QueryUnescape(testcertificate)
		if !strings.Contains(buf.String(), strings.ReplaceAll(testcertificate, "\n", "")) {
			t.Fatalf("bad body: %v vs %v", buf.String(), testcertificate)
		}
	})
}
