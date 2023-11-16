// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegisterPluginRuntime(t *testing.T) {
	mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultHandlerRegister))
	defer mockVaultServer.Close()

	cfg := DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().RegisterPluginRuntime(context.Background(), &RegisterPluginRuntimeInput{
		Name:         "gvisor",
		Type:         PluginRuntimeTypeContainer,
		OCIRuntime:   "runsc",
		CgroupParent: "/cpulimit/",
		CPU:          1,
		Memory:       10000,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPluginRuntime(t *testing.T) {
	for name, tc := range map[string]struct {
		body     string
		expected GetPluginRuntimeResponse
	}{
		"gvisor": {
			body: getPluginRuntimeResponse,
			expected: GetPluginRuntimeResponse{
				Name:         "gvisor",
				Type:         PluginRuntimeTypeContainer.String(),
				OCIRuntime:   "runsc",
				CgroupParent: "/cpulimit/",
				CPU:          1,
				Memory:       10000,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultHandlerInfo(tc.body)))
			defer mockVaultServer.Close()

			cfg := DefaultConfig()
			cfg.Address = mockVaultServer.URL
			client, err := NewClient(cfg)
			if err != nil {
				t.Fatal(err)
			}

			input := GetPluginRuntimeInput{
				Name: "gvisor",
				Type: PluginRuntimeTypeContainer,
			}

			info, err := client.Sys().GetPluginRuntime(context.Background(), &input)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.expected, *info) {
				t.Errorf("expected: %#v\ngot: %#v", tc.expected, info)
			}
		})
	}
}

func TestListPluginRuntimeTyped(t *testing.T) {
	for _, tc := range []struct {
		runtimeType      PluginRuntimeType
		body             string
		expectedResponse *ListPluginRuntimesResponse
		expectedErrNil   bool
	}{
		{
			runtimeType: PluginRuntimeTypeContainer,
			body:        listPluginRuntimeTypedResponse,
			expectedResponse: &ListPluginRuntimesResponse{
				Runtimes: []PluginRuntimeDetails{
					{
						Type:         "container",
						Name:         "gvisor",
						OCIRuntime:   "runsc",
						CgroupParent: "/cpulimit/",
						CPU:          1,
						Memory:       10000,
					},
				},
			},
			expectedErrNil: true,
		},
		{
			runtimeType:      PluginRuntimeTypeUnsupported,
			body:             listPluginRuntimeTypedResponse,
			expectedResponse: nil,
			expectedErrNil:   false,
		},
	} {
		t.Run(tc.runtimeType.String(), func(t *testing.T) {
			mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultHandlerInfo(tc.body)))
			defer mockVaultServer.Close()

			cfg := DefaultConfig()
			cfg.Address = mockVaultServer.URL
			client, err := NewClient(cfg)
			if err != nil {
				t.Fatal(err)
			}

			input := ListPluginRuntimesInput{
				Type: tc.runtimeType,
			}

			list, err := client.Sys().ListPluginRuntimes(context.Background(), &input)
			if tc.expectedErrNil && err != nil {
				t.Fatal(err)
			}

			if (tc.expectedErrNil && !reflect.DeepEqual(tc.expectedResponse, list)) || (!tc.expectedErrNil && list != nil) {
				t.Errorf("expected: %#v\ngot: %#v", tc.expectedResponse, list)
			}
		})
	}
}

func TestListPluginRuntimeUntyped(t *testing.T) {
	for _, tc := range []struct {
		body             string
		expectedResponse *ListPluginRuntimesResponse
		expectedErrNil   bool
	}{
		{
			body: listPluginRuntimeUntypedResponse,
			expectedResponse: &ListPluginRuntimesResponse{
				Runtimes: []PluginRuntimeDetails{
					{
						Type:         "container",
						Name:         "gvisor",
						OCIRuntime:   "runsc",
						CgroupParent: "/cpulimit/",
						CPU:          1,
						Memory:       10000,
					},
					{
						Type:         "container",
						Name:         "foo",
						OCIRuntime:   "otherociruntime",
						CgroupParent: "/memorylimit/",
						CPU:          2,
						Memory:       20000,
					},
					{
						Type:         "container",
						Name:         "bar",
						OCIRuntime:   "otherociruntime",
						CgroupParent: "/cpulimit/",
						CPU:          3,
						Memory:       30000,
					},
				},
			},
			expectedErrNil: true,
		},
	} {
		t.Run("", func(t *testing.T) {
			mockVaultServer := httptest.NewServer(http.HandlerFunc(mockVaultHandlerInfo(tc.body)))
			defer mockVaultServer.Close()

			cfg := DefaultConfig()
			cfg.Address = mockVaultServer.URL
			client, err := NewClient(cfg)
			if err != nil {
				t.Fatal(err)
			}

			info, err := client.Sys().ListPluginRuntimes(context.Background(), nil)
			if tc.expectedErrNil && err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.expectedResponse, info) {
				t.Errorf("expected: %#v\ngot: %#v", tc.expectedResponse, info)
			}
		})
	}
}

const getPluginRuntimeResponse = `{
    "request_id": "e93d3f93-8e4f-8443-a803-f1c97c123456",
    "data": {
        "name": "gvisor",
        "type": "container",
        "oci_runtime": "runsc",
        "cgroup_parent": "/cpulimit/",
        "cpu_nanos": 1,
        "memory_bytes": 10000
    },
    "warnings": null,
    "auth": null
}`

const listPluginRuntimeTypedResponse = `{
    "request_id": "e93d3f93-8e4f-8443-a803-f1c97c123456",
    "data": {
       "runtimes": [
			{
				"name": "gvisor",
				"type": "container",
				"oci_runtime": "runsc",
				"cgroup_parent": "/cpulimit/",
				"cpu_nanos": 1,
				"memory_bytes": 10000
			}
		]
    },
    "warnings": null,
    "auth": null
}
`

const listPluginRuntimeUntypedResponse = `{
    "request_id": "e93d3f93-8e4f-8443-a803-f1c97c123456",
    "data": {
        "runtimes": [
			{
				"name": "gvisor",
				"type": "container",
				"oci_runtime": "runsc",
				"cgroup_parent": "/cpulimit/",
				"cpu_nanos": 1,
				"memory_bytes": 10000
			},
			{
				"name": "foo",
				"type": "container",
				"oci_runtime": "otherociruntime",
				"cgroup_parent": "/memorylimit/",
				"cpu_nanos": 2,
				"memory_bytes": 20000
			},
			{
				"name": "bar",
				"type": "container",
				"oci_runtime": "otherociruntime",
				"cgroup_parent": "/cpulimit/",
				"cpu_nanos": 3,
				"memory_bytes": 30000
			}
		]
    },
    "warnings": null,
    "auth": null
}`
