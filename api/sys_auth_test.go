// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAuth(t *testing.T) {
	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultGetAuthHandler))
	defer mockVaultServer.Close()

	cfg := DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	auth, err := client.Sys().GetAuth("oidc")
	if err != nil {
		t.Fatal(err)
	}

	expected := struct {
		Type    string
		Version string
	}{Type: "oidc", Version: ""}

	if expected.Type != auth.Type || expected.Version != auth.PluginVersion {
		t.Errorf("auth mount did not match: expected %+v but got %+v", expected, auth)
	}
}

func mockVaultGetAuthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(getAuthMountResponse))
}

const getAuthMountResponse = `{
 "plugin_version": "",
 "config": {
  "default_lease_ttl": 0,
  "force_no_cache": false,
  "max_lease_ttl": 0,
  "token_type": "default-service"
 },
 "description": "OIDC backend",
 "accessor": "auth_oidc_9015643e",
 "seal_wrap": false,
 "external_entropy_access": false,
 "uuid": "c1194fea-ea22-6a27-4ceb-8b20bcca4728",
 "deprecation_status": "supported",
 "type": "oidc",
 "local": false,
 "options": null,
 "running_plugin_version": "v1.13.9+builtin.vault",
 "running_sha256": "",
 "request_id": "877bf1a1-752c-d1a0-1f53-b95af751b6d6",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": {
  "accessor": "auth_oidc_9015643e",
  "config": {
   "default_lease_ttl": 0,
   "force_no_cache": false,
   "max_lease_ttl": 0,
   "token_type": "default-service"
  },
  "deprecation_status": "supported",
  "description": "OIDC backend",
  "external_entropy_access": false,
  "local": false,
  "options": null,
  "plugin_version": "",
  "running_plugin_version": "v1.13.9+builtin.vault",
  "running_sha256": "",
  "seal_wrap": false,
  "type": "oidc",
  "uuid": "c1194fea-ea22-6a27-4ceb-8b20bcca4728"
 },
 "wrap_info": null,
 "warnings": null,
 "auth": null
}`
