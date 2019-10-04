package cache

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-test/deep"
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

func TestAPIProxy_requireRequestHeader(t *testing.T) {
	cleanup, client, _, _ := setupClusterAndAgent(namespace.RootContext(nil), t, nil)
	defer cleanup()

	proxier, err := NewAPIProxy(&APIProxyConfig{
		Client:               client,
		Logger:               logging.NewVaultLogger(hclog.Trace),
		RequireRequestHeader: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// test with missing header
	req := client.NewRequest("GET", "/v1/sys/health")

	httpReq, err := req.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: httpReq,
	})
	if diff := deep.Equal(err, errors.New(preconditionFailed)); diff != nil {
		t.Fatal(diff)
	}
	httpResp := resp.Response.Response
	if httpResp.StatusCode != http.StatusPreconditionFailed {
		t.Fatalf("expected response status code %d", http.StatusPreconditionFailed)
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if string(body) != preconditionFailed {
		t.Fatalf("expected response body %s", preconditionFailed)
	}

	// test with invalid header value
	req.Headers = make(http.Header)
	req.Headers[vaultRequestHeader] = []string{"bogus"}

	httpReq, err = req.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}
	resp, err = proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: httpReq,
	})
	if diff := deep.Equal(err, errors.New(preconditionFailed)); diff != nil {
		t.Fatal(diff)
	}
	httpResp = resp.Response.Response
	if httpResp.StatusCode != http.StatusPreconditionFailed {
		t.Fatalf("expected response status code %d", http.StatusPreconditionFailed)
	}
	body, err = ioutil.ReadAll(httpResp.Body)
	if string(body) != preconditionFailed {
		t.Fatalf("expected response body %s", preconditionFailed)
	}

	// test with correct header value
	req.Headers = make(http.Header)
	req.Headers[vaultRequestHeader] = []string{"true"}

	httpReq, err = req.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}
	resp, err = proxier.Send(namespace.RootContext(nil), &SendRequest{
		Request: httpReq,
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
