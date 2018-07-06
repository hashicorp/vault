package http

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/vault"
)

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
			return WrapForwardedForHandler(origHandler, []*sockaddr.SockAddrMarshaler{
				&sockaddr.SockAddrMarshaler{
					SockAddr: goodAddr,
				},
			}, true, false, 0)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: testHandler,
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
			return WrapForwardedForHandler(origHandler, []*sockaddr.SockAddrMarshaler{
				&sockaddr.SockAddrMarshaler{
					SockAddr: badAddr,
				},
			}, true, false, 0)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: testHandler,
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
			return WrapForwardedForHandler(origHandler, []*sockaddr.SockAddrMarshaler{
				&sockaddr.SockAddrMarshaler{
					SockAddr: badAddr,
				},
			}, true, true, 0)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: testHandler,
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
			return WrapForwardedForHandler(origHandler, []*sockaddr.SockAddrMarshaler{
				&sockaddr.SockAddrMarshaler{
					SockAddr: goodAddr,
				},
			}, true, true, 4)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: testHandler,
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
			return WrapForwardedForHandler(origHandler, []*sockaddr.SockAddrMarshaler{
				&sockaddr.SockAddrMarshaler{
					SockAddr: goodAddr,
				},
			}, true, true, 1)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: testHandler,
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
			return WrapForwardedForHandler(origHandler, []*sockaddr.SockAddrMarshaler{
				&sockaddr.SockAddrMarshaler{
					SockAddr: goodAddr,
				},
			}, true, true, 1)
		}

		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc: testHandler,
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
}
