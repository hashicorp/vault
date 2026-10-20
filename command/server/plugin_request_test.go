// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

type recordingRoundTripper struct {
	path string
	body []byte
}

func (r *recordingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.path = req.URL.Path
	if req.Body != nil {
		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		r.body = body
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"warnings":[]}`))),
	}, nil
}

func mockClient(t *testing.T) (*api.Client, *recordingRoundTripper) {
	t.Helper()

	recorder := &recordingRoundTripper{}
	config := api.DefaultConfig()
	config.Address = "http://127.0.0.1"
	config.HttpClient = &http.Client{Transport: recorder}
	client, err := api.NewClient(config)
	require.NoError(t, err)
	return client, recorder
}

// TestFlagParsing ensures that flags passed to vault plugin register correctly
// translate into the expected JSON body and request path.
func TestFlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		pluginType      api.PluginType
		name            string
		command         string
		ociImage        string
		runtime         string
		version         string
		sha256          string
		args            []string
		env             []string
		expectedPayload string
	}{
		"minimal": {
			pluginType:      api.PluginTypeUnknown,
			name:            "foo",
			sha256:          "abc123",
			expectedPayload: `{"type":"unknown","command":"foo","sha256":"abc123"}`,
		},
		"full": {
			pluginType:      api.PluginTypeCredential,
			name:            "name",
			command:         "cmd",
			ociImage:        "image",
			runtime:         "runtime",
			version:         "v1.0.0",
			sha256:          "abc123",
			args:            []string{"--a=b", "--b=c", "positional"},
			env:             []string{"x=1", "y=2"},
			expectedPayload: `{"type":"auth","args":["--a=b","--b=c","positional"],"command":"cmd","sha256":"abc123","version":"v1.0.0","oci_image":"image","runtime":"runtime","env":["x=1","y=2"]}`,
		},
		"command remains empty if oci_image specified": {
			pluginType:      api.PluginTypeCredential,
			name:            "name",
			ociImage:        "image",
			sha256:          "abc123",
			expectedPayload: `{"type":"auth","sha256":"abc123","oci_image":"image"}`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ui := cli.NewMockUi()
			cmd := &PluginRegisterCommand{BaseCommand: &base.BaseCommand{UI: ui}}
			client, requestLogger := mockClient(t)
			cmd.SetClient(client)

			var args []string
			if tc.command != "" {
				args = append(args, "-command="+tc.command)
			}
			if tc.ociImage != "" {
				args = append(args, "-oci_image="+tc.ociImage)
			}
			if tc.runtime != "" {
				args = append(args, "-runtime="+tc.runtime)
			}
			if tc.sha256 != "" {
				args = append(args, "-sha256="+tc.sha256)
			}
			if tc.version != "" {
				args = append(args, "-version="+tc.version)
			}
			for _, arg := range tc.args {
				args = append(args, "-args="+arg)
			}
			for _, env := range tc.env {
				args = append(args, "-env="+env)
			}
			if tc.pluginType != api.PluginTypeUnknown {
				args = append(args, tc.pluginType.String())
			}
			args = append(args, tc.name)

			code := cmd.Run(args)
			if exp := 0; code != exp {
				t.Fatalf("expected %d to be %d\nstdout: %s\nstderr: %s", code, exp, ui.OutputWriter.String(), ui.ErrorWriter.String())
			}

			actual := &api.RegisterPluginInput{}
			expected := &api.RegisterPluginInput{}
			err := json.Unmarshal(requestLogger.body, actual)
			if err != nil {
				t.Fatal(err)
			}
			err = json.Unmarshal([]byte(tc.expectedPayload), expected)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %s\ngot: %s", tc.expectedPayload, requestLogger.body)
			}
			expectedPath := fmt.Sprintf("/v1/sys/plugins/catalog/%s/%s", tc.pluginType.String(), tc.name)
			if tc.pluginType == api.PluginTypeUnknown {
				expectedPath = fmt.Sprintf("/v1/sys/plugins/catalog/%s", tc.name)
			}
			if requestLogger.path != expectedPath {
				t.Errorf("Expected path %s, got %s", expectedPath, requestLogger.path)
			}
		})
	}
}

// TestPluginRuntimeFlagParsing ensures that flags passed to vault plugin runtime register correctly
// translate into the expected JSON body and request path.
func TestPluginRuntimeFlagParsing(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		runtimeType     api.PluginRuntimeType
		name            string
		ociRuntime      string
		cgroupParent    string
		cpu             int64
		memory          int64
		rootless        bool
		expectedPayload string
	}{
		"minimal": {
			runtimeType:     api.PluginRuntimeTypeContainer,
			name:            "foo",
			expectedPayload: `{"type":1,"name":"foo"}`,
		},
		"full": {
			runtimeType:     api.PluginRuntimeTypeContainer,
			name:            "foo",
			cgroupParent:    "/cpulimit/",
			ociRuntime:      "runtime",
			cpu:             5678,
			memory:          1234,
			rootless:        true,
			expectedPayload: `{"type":1,"cgroup_parent":"/cpulimit/","memory_bytes":1234,"cpu_nanos":5678,"oci_runtime":"runtime","rootless":true}`,
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ui := cli.NewMockUi()
			cmd := &PluginRuntimeRegisterCommand{BaseCommand: &base.BaseCommand{UI: ui}}
			client, requestLogger := mockClient(t)
			cmd.SetClient(client)

			var args []string
			if tc.cgroupParent != "" {
				args = append(args, "-cgroup_parent="+tc.cgroupParent)
			}
			if tc.ociRuntime != "" {
				args = append(args, "-oci_runtime="+tc.ociRuntime)
			}
			if tc.memory != 0 {
				args = append(args, fmt.Sprintf("-memory_bytes=%d", tc.memory))
			}
			if tc.cpu != 0 {
				args = append(args, fmt.Sprintf("-cpu_nanos=%d", tc.cpu))
			}
			if tc.rootless {
				args = append(args, "-rootless=true")
			}
			if tc.runtimeType != api.PluginRuntimeTypeUnsupported {
				args = append(args, "-type="+tc.runtimeType.String())
			}
			args = append(args, tc.name)

			code := cmd.Run(args)
			if exp := 0; code != exp {
				t.Fatalf("expected %d to be %d\nstdout: %s\nstderr: %s", code, exp, ui.OutputWriter.String(), ui.ErrorWriter.String())
			}

			actual := &api.RegisterPluginRuntimeInput{}
			expected := &api.RegisterPluginRuntimeInput{}
			err := json.Unmarshal(requestLogger.body, actual)
			if err != nil {
				t.Fatal(err)
			}
			err = json.Unmarshal([]byte(tc.expectedPayload), expected)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %s\ngot: %s", tc.expectedPayload, requestLogger.body)
			}
			expectedPath := fmt.Sprintf("/v1/sys/plugins/runtimes/catalog/%s/%s", tc.runtimeType.String(), tc.name)
			if requestLogger.path != expectedPath {
				t.Errorf("Expected path %s, got %s", expectedPath, requestLogger.path)
			}
		})
	}
}
