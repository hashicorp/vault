// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/version"
	"github.com/stretchr/testify/require"
)

// TestVersionCommand_RunBranding verifies that branding leaves the build details
// and cgo suffix alone, still ignores extra arguments, and only writes to stdout.
func TestVersionCommand_RunBranding(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		info version.VersionInfo
		want string
	}{
		{
			name: "shared_default",
			info: version.VersionInfo{Version: "2.3.0"},
			want: "Vault v2.3.0",
		},
		{
			name: "client_all_fields",
			info: version.VersionInfo{
				BinaryName:        "Vault Client",
				Version:           "2.3.0",
				VersionPrerelease: "beta1",
				VersionMetadata:   "ent.hsm.Vault.custom",
				Revision:          "build-Vault-123",
				BuildDate:         "2026-10-27T12:34:56Z",
			},
			want: "Vault Client v2.3.0-beta1+ent.hsm.Vault.custom (build-Vault-123), built 2026-10-27T12:34:56Z",
		},
		{
			name: "custom_name_git_describe",
			info: version.VersionInfo{BinaryName: "Custom Vault", Version: "v2.3.0-beta1-7-gabc123"},
			want: "Custom Vault vv2.3.0-beta1-7-gabc123",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ui := cli.NewMockUi()
			info := tc.info
			cmd := &VersionCommand{
				BaseCommand: &base.BaseCommand{UI: ui},
				VersionInfo: &info,
			}
			expected := tc.want
			if version.CgoEnabled {
				expected += " (cgo)"
			}

			code := cmd.Run([]string{"ignored", "-unknown-flag", "--unknown-flag=value"})
			require.Equal(t, 0, code, "version must succeed despite extra arguments and flags")
			require.Empty(t, ui.ErrorWriter.String(), "version output belongs only on stdout")
			require.Equal(t, expected+"\n", ui.OutputWriter.String(), "branding must preserve the complete formatted version")
			require.Equal(t, tc.info, info, "the command must not rewrite the supplied version information")
		})
	}
}

// TestInitClientBinaryCommands_VersionBranding verifies that we pick branding
// from the command set, not the binary's name.
func TestInitClientBinaryCommands_VersionBranding(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		binaryName string
		init       base.CommandInitFn
		product    string
	}{
		{name: "client_named_vault", binaryName: "vault", init: initClientBinaryCommands, product: "Vault Client"},
		{name: "renamed_client", binaryName: "renamed-client", init: initClientBinaryCommands, product: "Vault Client"},
		{name: "shared_named_vault_client", binaryName: "vault-client", init: InitClientCommands},
		{name: "renamed_shared", binaryName: "renamed-vault", init: InitClientCommands},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ui := cli.NewMockUi()
			serverUI := cli.NewMockUi()
			commands := tc.init(ui, serverUI, &base.RunOptions{
				BinaryName: tc.binaryName,
				Stdin:      strings.NewReader(""),
			})
			factory, ok := commands["version"]
			require.True(t, ok, "both command profiles must register version")
			cmd, err := factory()
			require.NoError(t, err, "the real version factory must construct a command")

			info := version.GetVersion()
			info.BinaryName = tc.product
			expected := info.FullVersionNumber(true)
			if version.CgoEnabled {
				expected += " (cgo)"
			}
			require.Equal(t, 0, cmd.Run(nil), "the registered version command must succeed")
			require.Empty(t, ui.ErrorWriter.String(), "version output belongs only on stdout")
			require.Equal(t, expected+"\n", ui.OutputWriter.String(), "the profile must select branding while preserving the compiled build details")
			require.Empty(t, serverUI.OutputWriter.String(), "version must not use the server UI")
			require.Empty(t, serverUI.ErrorWriter.String(), "version must not write server diagnostics")
		})
	}
}

// TestRunCustom_VersionBranding verifies that version and its aliases print
// "Vault Client" with the same build details even after we rename the binary.
func TestRunCustom_VersionBranding(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		args       []string
		binaryName string
	}{
		{name: "version", args: []string{"version"}, binaryName: "vault-client"},
		{name: "short_alias_named_vault", args: []string{"-v"}, binaryName: "vault"},
		{name: "single_dash_alias_renamed", args: []string{"-version"}, binaryName: "renamed-client"},
		{name: "double_dash_alias", args: []string{"--version"}, binaryName: "vault-client"},
		{name: "ignored_arguments", args: []string{"version", "ignored", "-unknown-flag"}, binaryName: "vault-client"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer
			code := RunCustom(tc.args, &base.RunOptions{
				BinaryName: tc.binaryName,
				Stdin:      strings.NewReader(""),
				Stdout:     &stdout,
				Stderr:     &stderr,
			})
			info := version.GetVersion()
			info.BinaryName = "Vault Client"
			expected := info.FullVersionNumber(true)
			if version.CgoEnabled {
				expected += " (cgo)"
			}
			require.Equal(t, 0, code, "version invocations must succeed")
			require.Empty(t, stderr.String(), "version output belongs only on stdout")
			require.Equal(t, expected+"\n", stdout.String(), "client branding must preserve the complete build version")
		})
	}
}
