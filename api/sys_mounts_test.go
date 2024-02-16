// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMount(t *testing.T) {
	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultGetMountHandler))
	defer mockVaultServer.Close()

	cfg := DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	mount, err := client.Sys().GetMount("secret")
	if err != nil {
		t.Fatal(err)
	}

	expected := struct {
		Type    string
		Version string
	}{Type: "pki", Version: ""}

	if expected.Type != mount.Type || expected.Version != mount.PluginVersion {
		t.Errorf("Mount did not match: expected %+v but got %+v", expected, mount)
	}
}

func TestListMounts(t *testing.T) {
	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultListMountsHandler))
	defer mockVaultServer.Close()

	cfg := DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Sys().ListMounts()
	if err != nil {
		t.Fatal(err)
	}

	expectedMounts := map[string]struct {
		Type    string
		Version string
	}{
		"cubbyhole/": {Type: "cubbyhole", Version: "v1.0.0"},
		"identity/":  {Type: "identity", Version: ""},
		"secret/":    {Type: "kv", Version: ""},
		"sys/":       {Type: "system", Version: ""},
	}

	for path, mount := range resp {
		expected, ok := expectedMounts[path]
		if !ok {
			t.Errorf("Unexpected mount: %s: %+v", path, mount)
			continue
		}
		if expected.Type != mount.Type || expected.Version != mount.PluginVersion {
			t.Errorf("Mount did not match: %s -> expected %+v but got %+v", path, expected, mount)
		}
	}

	for path, expected := range expectedMounts {
		mount, ok := resp[path]
		if !ok {
			t.Errorf("Expected mount not found mount: %s: %+v", path, expected)
			continue
		}
		if expected.Type != mount.Type || expected.Version != mount.PluginVersion {
			t.Errorf("Mount did not match: %s -> expected %+v but got %+v", path, expected, mount)
		}
	}
}

func mockVaultListMountsHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(listMountsResponse))
}

func mockVaultGetMountHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(getMountResponse))
}

const getMountResponse = `{
 "uuid": "556fc6ef-208c-5ac1-3838-40b0c052eaa8",
 "plugin_version": "",
 "running_plugin_version": "v1.13.9+builtin.vault",
 "deprecation_status": "supported",
 "accessor": "pki_a6c43de6",
 "options": {},
 "local": false,
 "seal_wrap": false,
 "external_entropy_access": false,
 "running_sha256": "",
 "config": {
  "default_lease_ttl": 86400,
  "force_no_cache": false,
  "max_lease_ttl": 86400
 },
 "type": "pki",
 "description": "test",
 "request_id": "285b1a04-821c-686a-8325-cfc70517c3eb",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": {
  "accessor": "pki_a6c43de6",
  "config": {
   "default_lease_ttl": 86400,
   "force_no_cache": false,
   "max_lease_ttl": 86400
  },
  "deprecation_status": "supported",
  "description": "test",
  "external_entropy_access": false,
  "local": false,
  "options": {},
  "plugin_version": "",
  "running_plugin_version": "v1.13.9+builtin.vault",
  "running_sha256": "",
  "seal_wrap": false,
  "type": "pki",
  "uuid": "556fc6ef-208c-5ac1-3838-40b0c052eaa8"
 },
 "wrap_info": null,
 "warnings": null,
 "auth": null
}`

const listMountsResponse = `{
  "request_id": "3cd881e9-ea50-2e06-90b2-5641667485fa",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "cubbyhole/": {
      "accessor": "cubbyhole_2e3fc28d",
      "config": {
        "default_lease_ttl": 0,
        "force_no_cache": false,
        "max_lease_ttl": 0
      },
      "description": "per-token private secret storage",
      "external_entropy_access": false,
      "local": true,
      "options": null,
      "plugin_version": "v1.0.0",
      "running_sha256": "",
      "running_plugin_version": "",
      "seal_wrap": false,
      "type": "cubbyhole",
      "uuid": "575063dc-5ef8-4487-c842-22c494c19a6f"
    },
    "identity/": {
      "accessor": "identity_6e01c327",
      "config": {
        "default_lease_ttl": 0,
        "force_no_cache": false,
        "max_lease_ttl": 0,
        "passthrough_request_headers": [
          "Authorization"
        ]
      },
      "description": "identity store",
      "external_entropy_access": false,
      "local": false,
      "options": null,
      "plugin_version": "",
      "running_sha256": "",
      "running_plugin_version": "",
      "seal_wrap": false,
      "type": "identity",
      "uuid": "187d7eba-3471-554b-c2d9-1479612c8046"
    },
    "secret/": {
      "accessor": "kv_3e2f282f",
      "config": {
        "default_lease_ttl": 0,
        "force_no_cache": false,
        "max_lease_ttl": 0
      },
      "description": "key/value secret storage",
      "external_entropy_access": false,
      "local": false,
      "options": {
        "version": "2"
      },
      "plugin_version": "",
      "running_sha256": "",
      "running_plugin_version": "",
      "seal_wrap": false,
      "type": "kv",
      "uuid": "13375e0f-876e-7e96-0a3e-076f37b6b69d"
    },
    "sys/": {
      "accessor": "system_93503264",
      "config": {
        "default_lease_ttl": 0,
        "force_no_cache": false,
        "max_lease_ttl": 0,
        "passthrough_request_headers": [
          "Accept"
        ]
      },
      "description": "system endpoints used for control, policy and debugging",
      "external_entropy_access": false,
      "local": false,
      "options": null,
      "plugin_version": "",
      "running_sha256": "",
      "running_plugin_version": "",
      "seal_wrap": true,
      "type": "system",
      "uuid": "1373242d-cc4d-c023-410b-7f336e7ba0a8"
    }
  }
}`
