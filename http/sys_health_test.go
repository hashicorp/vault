// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestSysHealth_get(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	// Test without the client first since we want to verify the response code
	raw, err := http.Get(addr + "/v1/sys/health")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, raw, 501)

	// Test with the client because it's a bit easier to work with structs
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Sys().Health()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &api.HealthResponse{
		Enterprise:                 constants.IsEnterprise,
		Initialized:                false,
		Sealed:                     true,
		Standby:                    true,
		PerformanceStandby:         false,
		ReplicationPerformanceMode: consts.ReplicationUnknown.GetPerformanceString(),
		ReplicationDRMode:          consts.ReplicationUnknown.GetDRString(),
	}
	ignore := cmpopts.IgnoreFields(*expected, "ClusterName", "ClusterID", "ServerTimeUTC", "Version")
	if diff := cmp.Diff(resp, expected, ignore); len(diff) > 0 {
		t.Fatal(diff)
	}

	keys, _ := vault.TestCoreInit(t, core)
	raw, err = http.Get(addr + "/v1/sys/health")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, raw, 503)

	resp, err = client.Sys().Health()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected.Initialized = true
	if diff := cmp.Diff(resp, expected, ignore); len(diff) > 0 {
		t.Fatal(diff)
	}

	for _, key := range keys {
		if _, err := vault.TestCoreUnseal(core, vault.TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}
	raw, err = http.Get(addr + "/v1/sys/health")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, raw, 200)

	resp, err = client.Sys().Health()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected.Sealed = false
	expected.Standby = false
	expected.ReplicationPerformanceMode = consts.ReplicationPerformanceDisabled.GetPerformanceString()
	expected.ReplicationDRMode = consts.ReplicationDRDisabled.GetDRString()
	if diff := cmp.Diff(resp, expected, ignore); len(diff) > 0 {
		t.Fatal(diff)
	}
}

func TestSysHealth_customcodes(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	queryurl, err := url.Parse(addr + "/v1/sys/health?uninitcode=581&sealedcode=523&activecode=202")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	raw, err := http.Get(queryurl.String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, raw, 581)

	// Test with the client because it's a bit easier to work with structs
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Sys().Health()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &api.HealthResponse{
		Enterprise:                 constants.IsEnterprise,
		Initialized:                false,
		Sealed:                     true,
		Standby:                    true,
		PerformanceStandby:         false,
		ReplicationPerformanceMode: consts.ReplicationUnknown.GetPerformanceString(),
		ReplicationDRMode:          consts.ReplicationUnknown.GetDRString(),
	}
	ignore := cmpopts.IgnoreFields(*expected, "ClusterName", "ClusterID", "ServerTimeUTC", "Version")
	if diff := cmp.Diff(resp, expected, ignore); len(diff) > 0 {
		t.Fatal(diff)
	}

	keys, _ := vault.TestCoreInit(t, core)
	raw, err = http.Get(queryurl.String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, raw, 523)

	resp, err = client.Sys().Health()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected.Initialized = true
	if diff := cmp.Diff(resp, expected, ignore); len(diff) > 0 {
		t.Fatal(diff)
	}

	for _, key := range keys {
		if _, err := vault.TestCoreUnseal(core, vault.TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}
	raw, err = http.Get(queryurl.String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, raw, 202)

	resp, err = client.Sys().Health()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	expected.Sealed = false
	expected.Standby = false
	expected.ReplicationPerformanceMode = consts.ReplicationPerformanceDisabled.GetPerformanceString()
	expected.ReplicationDRMode = consts.ReplicationDRDisabled.GetDRString()
	if diff := cmp.Diff(resp, expected, ignore); len(diff) > 0 {
		t.Fatal(diff)
	}
}

func TestSysHealth_head(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	testData := []struct {
		uri  string
		code int
	}{
		{"", 200},
		{"?activecode=503", 503},
		{"?activecode=notacode", 400},
	}

	for _, tt := range testData {
		queryurl, err := url.Parse(addr + "/v1/sys/health" + tt.uri)
		if err != nil {
			t.Fatalf("err on %v: %s", queryurl, err)
		}
		resp, err := http.Head(queryurl.String())
		if err != nil {
			t.Fatalf("err on %v: %s", queryurl, err)
		}

		if resp.StatusCode != tt.code {
			t.Fatalf("HEAD %v expected code %d, got %d.", queryurl, tt.code, resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("err on %v: %s", queryurl, err)
		}
		if len(data) > 0 {
			t.Fatalf("HEAD %v expected no body, received \"%v\".", queryurl, data)
		}
	}
}

// TestSysHealth_Removed checks that a removed node returns a 530 and sets
// removed from cluster to be true. The test also checks that the removedcode
// query parameter is respected.
func TestSysHealth_Removed(t *testing.T) {
	core, err := vault.TestCoreWithMockRemovableNodeHABackend(t, true)
	require.NoError(t, err)
	vault.TestCoreInit(t, core)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	raw, err := http.Get(addr + "/v1/sys/health")
	require.NoError(t, err)
	testResponseStatus(t, raw, 530)
	healthResp := HealthResponse{}
	testResponseBody(t, raw, &healthResp)
	require.NotNil(t, healthResp.RemovedFromCluster)
	require.True(t, *healthResp.RemovedFromCluster)

	raw, err = http.Get(addr + "/v1/sys/health?removedcode=299")
	require.NoError(t, err)
	testResponseStatus(t, raw, 299)
	secondHealthResp := HealthResponse{}
	testResponseBody(t, raw, &secondHealthResp)
	require.NotNil(t, secondHealthResp.RemovedFromCluster)
	require.True(t, *secondHealthResp.RemovedFromCluster)
}
