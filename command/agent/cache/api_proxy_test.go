package cache

import (
	"net/http"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func TestAPIProxy(t *testing.T) {
	cleanup, client, _, _ := setupClusterAndAgent(namespace.RootContext(nil), t, nil)
	defer cleanup()

	proxier, err := NewAPIProxy(&APIProxyConfig{
		Client: client,
		Logger: logging.NewVaultLogger(hclog.Trace),
	})
	if err != nil {
		t.Fatal(err)
	}

	r := client.NewRequest("GET", "/v1/sys/health")
	req, err := r.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: req,
	})
	if err != nil {
		t.Fatal(err)
	}

	var result api.HealthResponse
	err = jsonutil.DecodeJSONFromReader(resp.Response.Body, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Initialized || result.Sealed || result.Standby {
		t.Fatalf("bad sys/health response: %#v", result)
	}
}

func TestAPIProxy_queryParams(t *testing.T) {
	// Set up an agent that points to a standby node for this particular test
	// since it needs to proxy a /sys/health?standbyok=true request to a standby
	cleanup, client, _, _ := setupClusterAndAgentOnStandby(namespace.RootContext(nil), t, nil)
	defer cleanup()

	proxier, err := NewAPIProxy(&APIProxyConfig{
		Client: client,
		Logger: logging.NewVaultLogger(hclog.Trace),
	})
	if err != nil {
		t.Fatal(err)
	}

	r := client.NewRequest("GET", "/v1/sys/health")
	req, err := r.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}

	// Add a query parameter for testing
	q := req.URL.Query()
	q.Add("standbyok", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: req,
	})
	if err != nil {
		t.Fatal(err)
	}

	var result api.HealthResponse
	err = jsonutil.DecodeJSONFromReader(resp.Response.Body, &result)
	if err != nil {
		t.Fatal(err)
	}

	if !result.Initialized || result.Sealed || !result.Standby {
		t.Fatalf("bad sys/health response: %#v", result)
	}

	if resp.Response.StatusCode != http.StatusOK {
		t.Fatalf("exptected standby to return 200, got: %v", resp.Response.StatusCode)
	}
}
