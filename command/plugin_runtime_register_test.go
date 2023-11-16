// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/cli"
)

func testPluginRuntimeRegisterCommand(tb testing.TB) (*cli.MockUi, *PluginRuntimeRegisterCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PluginRuntimeRegisterCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPluginRuntimeRegisterCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		flags []string
		args  []string
		out   string
		code  int
	}{
		{
			"no type specified",
			[]string{},
			[]string{"foo"},
			"-type is required for plugin runtime registration",
			1,
		},
		{
			"invalid type",
			[]string{"-type", "foo"},
			[]string{"not"},
			"\"foo\" is not a supported plugin runtime type",
			2,
		},
		{
			"not_enough_args",
			[]string{"-type", consts.PluginRuntimeTypeContainer.String()},
			[]string{},
			"Not enough arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"-type", consts.PluginRuntimeTypeContainer.String()},
			[]string{"foo", "bar"},
			"Too many arguments",
			1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testPluginRuntimeRegisterCommand(t)
			cmd.client = client

			args := append(tc.flags, tc.args...)
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

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPluginRuntimeRegisterCommand(t)
		cmd.client = client

		code := cmd.Run([]string{"-type", consts.PluginRuntimeTypeContainer.String(), "my-plugin-runtime"})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error registering plugin runtime my-plugin-runtime"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPluginRuntimeRegisterCommand(t)
		assertNoTabs(t, cmd)
	})
}

// TestPluginRuntimeFlagParsing ensures that flags passed to vault plugin runtime register correctly
// translate into the expected JSON body and request path.
func TestPluginRuntimeFlagParsing(t *testing.T) {
	for name, tc := range map[string]struct {
		runtimeType     api.PluginRuntimeType
		name            string
		ociRuntime      string
		cgroupParent    string
		cpu             int64
		memory          int64
		args            []string
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
			expectedPayload: `{"type":1,"cgroup_parent":"/cpulimit/","memory_bytes":1234,"cpu_nanos":5678,"oci_runtime":"runtime"}`,
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			ui, cmd := testPluginRuntimeRegisterCommand(t)
			var requestLogger *recordingRoundTripper
			cmd.client, requestLogger = mockClient(t)

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

			if tc.runtimeType != api.PluginRuntimeTypeUnsupported {
				args = append(args, "-type="+tc.runtimeType.String())
			}
			args = append(args, tc.name)
			t.Log(args)

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
