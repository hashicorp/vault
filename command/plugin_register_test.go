// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/cli"
)

func testPluginRegisterCommand(tb testing.TB) (*cli.MockUi, *PluginRegisterCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginRegisterCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginRegisterCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			nil,
			"Not enough arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"foo", "bar", "fizz"},
			"Too many arguments",
			1,
		},
		{
			"not_a_plugin",
			[]string{consts.PluginTypeCredential.String(), "nope_definitely_never_a_plugin_nope"},
			"",
			2,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testPluginRegisterCommand(t)
			cmd.client = client

			args := append([]string{"-sha256", "abcd1234"}, tc.args...)
			code := cmd.Run(args)
			if code != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, tc.out) {
				t.Errorf("expected %q to contain %q", combined, tc.out)
			}
		})
	}

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		pluginDir, cleanup := corehelpers.MakeTestPluginDir(t)
		defer cleanup(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		pluginName := "my-plugin"
		_, sha256Sum := testPluginCreate(t, pluginDir, pluginName)

		ui, cmd := testPluginRegisterCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-sha256", sha256Sum,
			consts.PluginTypeCredential.String(), pluginName,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Registered plugin: my-plugin"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
			Type: api.PluginTypeCredential,
		})
		if err != nil {
			t.Fatal(err)
		}

		found := false
		for _, plugins := range resp.PluginsByType {
			for _, p := range plugins {
				if p == pluginName {
					found = true
				}
			}
		}
		if !found {
			t.Errorf("expected %q to be in %q", pluginName, resp.PluginsByType)
		}
	})

	t.Run("integration with version", func(t *testing.T) {
		t.Parallel()

		pluginDir, cleanup := corehelpers.MakeTestPluginDir(t)
		defer cleanup(t)

		client, _, closer := testVaultServerPluginDir(t, pluginDir)
		defer closer()

		const pluginName = "my-plugin"
		versions := []string{"v1.0.0", "v2.0.1"}
		_, sha256Sum := testPluginCreate(t, pluginDir, pluginName)
		types := []api.PluginType{api.PluginTypeCredential, api.PluginTypeDatabase, api.PluginTypeSecrets}

		for _, typ := range types {
			for _, version := range versions {
				ui, cmd := testPluginRegisterCommand(t)
				cmd.client = client

				code := cmd.Run([]string{
					"-version=" + version,
					"-sha256=" + sha256Sum,
					typ.String(),
					pluginName,
				})
				if exp := 0; code != exp {
					t.Errorf("expected %d to be %d", code, exp)
				}

				expected := "Success! Registered plugin: my-plugin"
				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, expected) {
					t.Errorf("expected %q to contain %q", combined, expected)
				}
			}
		}

		resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
			Type: api.PluginTypeUnknown,
		})
		if err != nil {
			t.Fatal(err)
		}

		found := make(map[api.PluginType]int)
		versionsFound := make(map[api.PluginType][]string)
		for _, p := range resp.Details {
			if p.Name == pluginName {
				typ, err := api.ParsePluginType(p.Type)
				if err != nil {
					t.Fatal(err)
				}
				found[typ]++
				versionsFound[typ] = append(versionsFound[typ], p.Version)
			}
		}

		for _, typ := range types {
			if found[typ] != 2 {
				t.Fatalf("expected %q to be found 2 times, but found it %d times for %s type in %#v", pluginName, found[typ], typ.String(), resp.Details)
			}
			sort.Strings(versions)
			sort.Strings(versionsFound[typ])
			if !reflect.DeepEqual(versions, versionsFound[typ]) {
				t.Fatalf("expected %v versions but got %v", versions, versionsFound[typ])
			}
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPluginRegisterCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-sha256", "abcd1234",
			consts.PluginTypeCredential.String(), "my-plugin",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error registering plugin my-plugin:"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginRegisterCommand(t)
		assertNoTabs(t, cmd)
	})
}

// TestFlagParsing ensures that flags passed to vault plugin register correctly
// translate into the expected JSON body and request path.
func TestFlagParsing(t *testing.T) {
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
			expectedPayload: `{"type":0,"command":"foo","sha256":"abc123"}`,
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
			expectedPayload: `{"type":1,"args":["--a=b","--b=c","positional"],"command":"cmd","sha256":"abc123","version":"v1.0.0","oci_image":"image","runtime":"runtime","env":["x=1","y=2"]}`,
		},
		"command remains empty if oci_image specified": {
			pluginType:      api.PluginTypeCredential,
			name:            "name",
			ociImage:        "image",
			sha256:          "abc123",
			expectedPayload: `{"type":1,"sha256":"abc123","oci_image":"image"}`,
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			ui, cmd := testPluginRegisterCommand(t)
			var requestLogger *recordingRoundTripper
			cmd.client, requestLogger = mockClient(t)

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
			t.Log(args)

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

func mockClient(t *testing.T) (*api.Client, *recordingRoundTripper) {
	t.Helper()

	config := api.DefaultConfig()
	httpClient := cleanhttp.DefaultClient()
	roundTripper := &recordingRoundTripper{}
	httpClient.Transport = roundTripper
	config.HttpClient = httpClient
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}

	return client, roundTripper
}

var _ http.RoundTripper = (*recordingRoundTripper)(nil)

type recordingRoundTripper struct {
	path string
	body []byte
}

func (r *recordingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.path = req.URL.Path
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	r.body = body
	return &http.Response{
		StatusCode: 200,
	}, nil
}
