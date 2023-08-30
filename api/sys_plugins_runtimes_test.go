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

	err = client.Sys().RegisterPluginRuntimeWithContext(context.Background(), &RegisterPluginRuntimeInput{
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

			info, err := client.Sys().GetPluginRuntimeWithContext(context.Background(), &input)
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
				RuntimesByType: map[PluginRuntimeType][]string{
					PluginRuntimeTypeContainer: {"gvisor"},
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
				RuntimesByType: map[PluginRuntimeType][]string{
					PluginRuntimeTypeContainer: {"gvisor", "foo", "bar"},
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
        "cpu": 1,
        "memory": 10000
    },
    "warnings": null,
    "auth": null
}`

const listPluginRuntimeTypedResponse = `{
    "request_id": "e93d3f93-8e4f-8443-a803-f1c97c123456",
    "data": {
        "container": ["gvisor"]
    },
    "warnings": null,
    "auth": null
}
`

const listPluginRuntimeUntypedResponse = `{
    "request_id": "e93d3f93-8e4f-8443-a803-f1c97c123456",
    "data": {
        "container": ["gvisor", "foo", "bar"]
    },
    "warnings": null,
    "auth": null
}`
