// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/vault"
	"golang.org/x/net/http2"
)

func TestFeatureFlags(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc:             vaulthttp.Handler,
		RequestResponseCallback: schema.ResponseValidatingCallback(t),
	})
	cluster.Start()
	defer cluster.Cleanup()

	// Wait for core to start
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// Create a raw http connection copying the configuration
	// created by NewTestCluster
	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = cluster.Cores[0].TLSConfig()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	callApi := func() map[string]interface{} {
		// Use the normal API client to construct the URL
		req := client.NewRequest("GET", "/v1/sys/internal/ui/feature-flags")
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
		err = json.Unmarshal(httpRespBody, &httpResp)
		if err != nil {
			t.Fatal(err)
		}
		return httpResp
	}

	// First try with no environment variable set
	httpResp := callApi()
	featureFlags, ok := httpResp["feature_flags"]
	if !ok {
		t.Fatal("Missing 'feature_flags' in response")
	}
	if featureFlags != nil {
		t.Fatal("Nonempty 'feature_flags'")
	}

	// Now try with the environment variable temporarily set
	envVar := "VAULT_CLOUD_ADMIN_NAMESPACE"
	os.Setenv(envVar, "1")
	defer os.Unsetenv(envVar)

	httpResp = callApi()
	featureFlags, ok = httpResp["feature_flags"]
	if !ok {
		t.Fatal("Missing 'feature_flags' in response")
	}
	flagList := featureFlags.([]interface{})
	if len(flagList) != 1 {
		t.Fatalf("Bad length for 'feature_flags': %v", flagList)
	}
	flag := flagList[0].(string)
	if flag != envVar {
		t.Fatalf("Bad environment variable in `feature_flags`: %q", flag)
	}
}
