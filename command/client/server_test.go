// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/version"
	"github.com/stretchr/testify/require"
)

// TestRunCustom_UnavailableCommands checks that commands we leave out stay hidden
// and point users to vault without printing their arguments. Help, typos and
// supported operator subcommands such as key-status should still work.
func TestRunCustom_UnavailableCommands(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		args        []string
		unavailable string
		code        int
		help        bool
		usage       string
		failure     string
	}{
		{args: []string{"server"}, unavailable: "server", code: 1},
		{args: []string{"audit"}, unavailable: "audit", code: 1},
		{args: []string{"policy"}, unavailable: "policy", code: 1},
		{args: []string{"plugin"}, unavailable: "plugin", code: 1},
		{args: []string{"pki"}, unavailable: "pki", code: 1},
		{args: []string{"transit"}, unavailable: "transit", code: 1},
		{args: []string{"transform"}, unavailable: "transform", code: 1},
		{args: []string{"debug"}, unavailable: "debug", code: 1},
		{args: []string{"monitor"}, unavailable: "monitor", code: 1},
		{args: []string{"version-history"}, unavailable: "version-history", code: 1},
		{args: []string{"events"}, unavailable: "events", code: 1},
		{args: []string{"license"}, unavailable: "license", code: 1},
		{args: []string{"operator", "raft"}, unavailable: "operator raft", code: 1},
		{args: []string{"operator", "seal"}, unavailable: "operator seal", code: 1},
		{args: []string{"operator", "init"}, unavailable: "operator init", code: 1},
		{args: []string{"operator", "step-down"}, unavailable: "operator step-down", code: 1},
		{args: []string{"audit", "list"}, unavailable: "audit list", code: 1},
		{args: []string{"policy", "list"}, unavailable: "policy list", code: 1},
		{args: []string{"operator", "raft", "snapshot", "load"}, unavailable: "operator raft snapshot load", code: 1},
		// A typo here falls back to operator help, not the unavailable command.
		{args: []string{"operator", "raft", "snapshotx"}, code: 1, usage: "Usage: vault-client operator"},
		{args: []string{"operator", "raft", "snapshotx", "--help"}, help: true, usage: "Usage: vault-client operator"},
		{args: []string{"operator", "raft", "snapshot", "loadx", "--help"}, help: true, usage: "Usage: vault-client operator"},
		{args: []string{"--help", "operator", "raft", "snapshotx"}, help: true, usage: "Usage: vault-client operator"},
		{args: []string{"operator", "raft", "snapshotx", "-output-curl-string"}, code: 1, usage: "Usage: vault-client operator", failure: "Unable to generate cURL string from command\n"},
		{args: []string{"operator", "raft", "snapshot", "loadx", "-output-policy"}, code: 1, usage: "Usage: vault-client operator", failure: "Unable to generate policy from command\n"},
		{args: []string{"server", "-dev"}, unavailable: "server", code: 1},
		{args: []string{"audit", "--help"}, unavailable: "audit", help: true},
		{args: []string{"operator", "raft", "--help"}, unavailable: "operator raft", help: true},
		{args: []string{"server", "--help"}, unavailable: "server", help: true},
		{args: []string{"--help", "operator", "raft"}, unavailable: "operator raft", help: true},
		{args: []string{"audit", "-output-curl-string"}, unavailable: "audit", code: 1, failure: "Unable to generate cURL string from command\n"},
		{args: []string{"policy", "-output-policy"}, unavailable: "policy", code: 1, failure: "Unable to generate policy from command\n"},
		{args: []string{"audit", "-output-policy", "-output-curl-string"}, unavailable: "audit", code: 1, failure: "Unable to generate cURL string from command\n"},
		{args: []string{"audit", "--help", "-output-curl-string"}, unavailable: "audit", code: 1, help: true, failure: "Could not generate cURL command"},
		{args: []string{"policy", "--help", "-output-policy"}, unavailable: "policy", code: 1, help: true, failure: "Could not generate policy"},
		{args: []string{"operator", "rekey", "dummy-sensitive-operand"}, unavailable: "operator rekey", code: 1},
		{args: []string{"operator", "raft", "dummy sensitive operand"}, unavailable: "operator raft", code: 1},
		{args: []string{"operator", "raft", "", "dummy-sensitive-operand"}, unavailable: "operator raft", code: 1},
		{args: []string{"operator", "raft", "--", "snapshotx", "--help"}, unavailable: "operator raft", code: 1},
		{args: []string{"audit", "frobnicate"}, unavailable: "audit", code: 1},
		{args: []string{"statu"}, code: 127, usage: "Usage: vault-client <command>"},
		{args: []string{"operator", "raftx"}, code: 1, usage: "Usage: vault-client operator"},
		{args: []string{"operator", "raftx", "--help"}, help: true, usage: "Usage: vault-client operator"},
		{args: []string{"operator", "frobnicate", "--help"}, help: true, usage: "Usage: vault-client operator"},
		{args: []string{"status", "--help"}, help: true, usage: "Usage: vault-client status"},
		{args: []string{"ssh", "--help"}, help: true, usage: "Usage: vault-client ssh"},
		// These operator subcommands are supported in vault-client (AC-1.3).
		{args: []string{"operator", "key-status", "--help"}, help: true, usage: "Usage: vault-client operator key-status"},
		{args: []string{"operator", "rotate", "--help"}, help: true, usage: "Usage: vault-client operator rotate"},
		{args: []string{"operator", "unseal", "--help"}, help: true, usage: "Usage: vault-client operator unseal"},
		{args: []string{"operator"}, code: 1, usage: "Usage: vault-client operator"},
		{args: []string{"operator", "--help"}, help: true, usage: "Usage: vault-client operator"},
		{args: []string{"operator", "--help", "raft"}, help: true, usage: "Usage: vault-client operator"},
		{args: []string{"--help"}, help: true, usage: "Usage: vault-client <command>"},
	} {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer
			code := RunCustom(tc.args, &base.RunOptions{
				Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr,
				BinaryName: "vault-client",
			})
			require.Equal(t, tc.code, code)
			out := stdout.String() + stderr.String()
			require.NotContains(t, out, "dummy-sensitive-operand")
			require.NotContains(t, out, "dummy sensitive operand")
			require.NotContains(t, out, "curl ")
			require.NotContains(t, out, "capabilities =")
			if tc.unavailable != "" {
				guide := fmt.Sprintf("%q is not available in vault-client v%s.\nThis command requires the full Vault Server image.\nObtain the full Vault binary: %s\n",
					tc.unavailable, version.GetVersion().VersionNumber(), serverInstallURL)
				if tc.help {
					require.Equal(t, guide, stdout.String())
					require.Equal(t, tc.failure, stderr.String())
				} else {
					require.Empty(t, stdout.String())
					require.Equal(t, guide+tc.failure, stderr.String())
				}
				require.NotContains(t, out, "Usage:")
				require.NotContains(t, out, "Subcommands:")
				return
			}
			if tc.help {
				require.Contains(t, stdout.String(), tc.usage)
				require.Equal(t, tc.failure, stderr.String())
			} else {
				require.Empty(t, stdout.String())
				require.Contains(t, stderr.String(), tc.usage)
				if tc.failure != "" {
					require.True(t, strings.HasSuffix(stderr.String(), tc.failure))
				}
			}
			require.NotContains(t, out, "not available in vault-client")
			if tc.usage == "Usage: vault-client <command>" {
				for _, name := range []string{"server", "audit", "policy", "plugin", "pki", "transit", "transform", "debug", "monitor", "version-history", "events", "license"} {
					require.NotRegexp(t, `(?m)^\s+`+name+`\s`, out, name)
				}
			}
		})
	}
}
