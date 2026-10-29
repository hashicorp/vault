// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/hashicorp/vault/version"
	"github.com/stretchr/testify/require"
)

func testVersionCommand(tb testing.TB) (*cli.MockUi, *clientcmd.VersionCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &clientcmd.VersionCommand{
		VersionInfo: &version.VersionInfo{},
		BaseCommand: &base.BaseCommand{
			UI: ui,
		},
	}
}

func TestVersionCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("output", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testVersionCommand(t)
		cmd.SetClient(client)

		code := cmd.Run(nil)
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Vault"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to equal %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testVersionCommand(t)
		assertNoTabs(t, cmd)
	})
}

// TestRunCustom_VersionBrandingUnchanged verifies that client branding doesn't
// leak into the full vault binary, even when we name it vault-client.
func TestRunCustom_VersionBrandingUnchanged(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		args       []string
		binaryName string
	}{
		{name: "version", args: []string{"version"}, binaryName: "vault"},
		{name: "short_alias_named_vault_client", args: []string{"-v"}, binaryName: "vault-client"},
		{name: "single_dash_alias_renamed", args: []string{"-version"}, binaryName: "renamed-vault"},
		{name: "double_dash_alias", args: []string{"--version"}, binaryName: "vault"},
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
			expected := version.GetVersion().FullVersionNumber(true)
			if version.CgoEnabled {
				expected += " (cgo)"
			}
			require.Equal(t, 0, code, "full CLI version invocations must succeed")
			require.Empty(t, stderr.String(), "version output belongs only on stdout")
			require.Equal(t, expected+"\n", stdout.String(), "full CLI version output must remain unchanged")
			require.NotContains(t, stdout.String(), "Vault Client ", "the full CLI must not acquire client branding")
		})
	}
}
