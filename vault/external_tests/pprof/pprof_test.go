// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pprof

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

func TestSysPprof(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc:             vaulthttp.Handler,
		RequestResponseCallback: schema.ResponseValidatingCallback(t),
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	testSysPprof(t, cluster)
}

// TestSysPprof_Exec is the same as TestSysPprof, but using a Vault binary running as -dev
// instead of a fake single node TestCluster.  There's no particular reason why
// TestSysPprof was chosen to validate that mechanism, other than that it was fast and simple.
func TestSysPprof_Exec(t *testing.T) {
	cluster := testcluster.NewTestExecDevCluster(t, &testcluster.ExecDevClusterOptions{
		ClusterOptions: testcluster.ClusterOptions{
			NumCores: 1,
		},
		// If $VAULT_BINARY is unset, use the `vault` in the system path.  Our CI should
		// populate the env var to use a `vault` built from the current commit.
		BinaryPath: os.Getenv("VAULT_BINARY"),
	})
	defer cluster.Cleanup()

	testSysPprof(t, cluster)
}

func testSysPprof(t *testing.T, cluster testcluster.VaultCluster) {
	nodes := cluster.Nodes()
	client := nodes[0].APIClient()

	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = nodes[0].TLSConfig()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	cases := []struct {
		name    string
		path    string
		seconds string
	}{
		{
			"index",
			"/v1/sys/pprof/",
			"",
		},
		{
			"cmdline",
			"/v1/sys/pprof/cmdline",
			"",
		},
		{
			"goroutine",
			"/v1/sys/pprof/goroutine",
			"",
		},
		{
			"heap",
			"/v1/sys/pprof/heap",
			"",
		},
		{
			"profile",
			"/v1/sys/pprof/profile",
			"1",
		},
		{
			"symbol",
			"/v1/sys/pprof/symbol",
			"",
		},
		{
			"trace",
			"/v1/sys/pprof/trace",
			"1",
		},
	}

	pprofRequest := func(t *testing.T, path string, seconds string) {
		req := client.NewRequest("GET", path)
		if seconds != "" {
			req.Params.Set("seconds", seconds)
		}
		httpReq, err := req.ToHTTP()
		if err != nil {
			t.Fatal(err)
		}
		resp, err := httpClient.Do(httpReq)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		httpRespBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		httpResp := make(map[string]interface{})

		// Skip this error check since some endpoints return binary blobs, we
		// only care about the ok check right after as an existence check.
		_ = json.Unmarshal(httpRespBody, &httpResp)

		// Make sure that we don't get a error response
		if _, ok := httpResp["errors"]; ok {
			t.Fatalf("unexpected error response: %v", httpResp["errors"])
		}

		if len(httpRespBody) == 0 {
			t.Fatal("no pprof index returned")
		}
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pprofRequest(t, tc.path, tc.seconds)
		})
	}
}

func TestSysPprof_MaxRequestDuration(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = cluster.Cores[0].TLSConfig()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	sec := strconv.Itoa(int(vault.DefaultMaxRequestDuration.Seconds()) + 1)

	req := client.NewRequest("GET", "/v1/sys/pprof/profile")
	req.Params.Set("seconds", sec)
	httpReq, err := req.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	httpRespBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	httpResp := make(map[string]interface{})

	// If we error here, it means that profiling likely happened, which is not
	// what we're checking for in this case.
	if err := json.Unmarshal(httpRespBody, &httpResp); err != nil {
		t.Fatalf("expected valid error response, got: %v", err)
	}

	errs, ok := httpResp["errors"].([]interface{})
	if !ok {
		t.Fatalf("expected error response, got: %v", httpResp)
	}
	if len(errs) == 0 || !strings.Contains(errs[0].(string), "exceeds max request duration") {
		t.Fatalf("unexpected error returned: %v", errs)
	}
}

func TestSysPprof_Standby(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		DisablePerformanceStandby: true,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		DefaultHandlerProperties: vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{
				Profiling: configutil.ListenerProfiling{
					UnauthenticatedPProfAccess: true,
				},
			},
		},
	})
	defer cluster.Cleanup()

	testSysPprof_Standby(t, cluster)
}

func TestSysPprof_Standby_Exec(t *testing.T) {
	cluster := testcluster.NewTestExecDevCluster(t, &testcluster.ExecDevClusterOptions{
		ClusterOptions: testcluster.ClusterOptions{
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				DisablePerformanceStandby: true,
			},
		},
	})
	defer cluster.Cleanup()

	testSysPprof_Standby(t, cluster)
}

func testSysPprof_Standby(t *testing.T, cluster testcluster.VaultCluster) {
	pprof := func(client *api.Client) (string, error) {
		req := client.NewRequest("GET", "/v1/sys/pprof/cmdline")
		resp, err := client.RawRequestWithContext(context.Background(), req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		return string(data), err
	}

	cmdline, err := pprof(cluster.Nodes()[0].APIClient())
	require.Nil(t, err)
	require.NotEmpty(t, cmdline)
	t.Log(cmdline)

	cmdline, err = pprof(cluster.Nodes()[1].APIClient())
	require.Nil(t, err)
	require.NotEmpty(t, cmdline)
	t.Log(cmdline)
}
