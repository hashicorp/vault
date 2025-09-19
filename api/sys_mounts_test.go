// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestListMounts(t *testing.T) {
	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultMountsHandler))
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

func mockVaultMountsHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(listMountsResponse))
}

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

// TestMountUpdateConfigStructFields ensures that the type MountConfigInput (the struct used to setup a new mount) has
// all the same fields as type TuneMountConfigInput (the struct used to tune an existing mount):
func TestMountUpdateConfigStructFields(t *testing.T) {
	tuneStruct := TuneMountConfigInput{}
	initialStruct := MountConfigInput{}

	tuneReflect := reflect.ValueOf(tuneStruct).Type()
	initialReflect := reflect.ValueOf(initialStruct).Type()

	for i := 0; i < tuneReflect.NumField(); i++ {
		tuneField := tuneReflect.Field(i)
		foundMatch := false
		// Now Find the field in initial Reflect
		for j := 0; j < initialReflect.NumField(); j++ {
			initialField := initialReflect.Field(j)
			if tuneField.Name == initialField.Name {
				jsonTuneFieldTag := tuneField.Tag.Get("json")
				jsonTuneFieldName := strings.Split(jsonTuneFieldTag, ",")[0]
				jsonInitialFieldTag := tuneField.Tag.Get("json")
				jsonInitialFieldName := strings.Split(jsonInitialFieldTag, ",")[0]
				if jsonTuneFieldName != jsonInitialFieldName {
					t.Fatalf("TuneMountConfigInput and MountConfigInput struct fields %v do not have same json names: %v, %v", tuneField.Name, jsonTuneFieldName, jsonInitialFieldName)
				}
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			t.Fatalf("Field %s in TuneMountConfigInput not found in MountConfigInput", tuneField.Name)
		}
	}

	if tuneReflect.NumField() != initialReflect.NumField() {
		for i := 0; i < initialReflect.NumField(); i++ {
			initialField := initialReflect.Field(i)
			foundMatch := false
			for j := 0; j < tuneReflect.NumField(); j++ {
				tuneField := tuneReflect.Field(j)
				if tuneField.Name == initialField.Name {
					foundMatch = true
					break
				}
			}
			if !foundMatch {
				t.Fatalf("Field %s in MountConfigInput not found in TuneMountConfigInput", initialField.Name)
			}
		}
		t.Fatalf("Different number of TuneMountConfigInput fields found in MountConfigInput")
	}
}
