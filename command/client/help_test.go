// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// TestRunCustom_ClientHelpListings checks that vault-client help only lists and
// shows examples for commands it supports, whether we ask for help or mistype a
// command. Operator help should list key-status, rotate, and unseal and show an
// example for each, but nothing for the operator subcommands we leave out. Help
// for a subcommand we leave out should explain that it isn't available, and
// help for unseal should be its real help.
func TestRunCustom_ClientHelpListings(t *testing.T) {
	t.Parallel()

	excludedOperatorExample := regexp.MustCompile(`(?m)^\s*\$\s+\S+\s+operator\s+` +
		`(seal|init|step-down|raft|generate-root|rekey|members|migrate|diagnose|usage|utilization|import)(\s|$)`)

	for _, tc := range []struct {
		name       string
		binaryName string
		args       []string
		kind       string
		explicit   bool
		shim       string
	}{
		{"root_help", "vault-client", []string{"--help"}, "root", true, ""},
		{"root_fallback", "vault-client", []string{"statu"}, "root", false, ""},
		{"renamed_root_help", "vault", []string{"--help"}, "root", true, ""},
		{"renamed_root_fallback", "vault", []string{"statu"}, "root", false, ""},
		{"auth_help", "vault-client", []string{"auth", "--help"}, "auth", true, ""},
		{"auth_fallback", "vault-client", []string{"auth", "kist"}, "auth", false, ""},
		{"operator_help", "vault-client", []string{"operator", "--help"}, "operator", true, ""},
		{"renamed_operator_help", "vault", []string{"operator", "--help"}, "operator", true, ""},
		{"operator_bare", "vault-client", []string{"operator"}, "operator", false, ""},
		{"operator_fallback", "vault-client", []string{"operator", "not-a-command"}, "operator", false, ""},
		// A typo under an unavailable group falls back to operator help.
		{"operator_partial_help", "vault-client", []string{"operator", "raft", "snapshotx", "--help"}, "operator", true, ""},
		{"operator_raft_help", "vault-client", []string{"operator", "raft", "--help"}, "unavailable", true, "operator raft"},
		{"operator_autopilot_help", "vault-client", []string{"operator", "raft", "autopilot", "--help"}, "unavailable", true, "operator raft autopilot"},
		{"help_operator_raft", "vault-client", []string{"--help", "operator", "raft"}, "unavailable", true, "operator raft"},
		{"operator_seal", "vault-client", []string{"operator", "seal"}, "unavailable", false, "operator seal"},
		{"operator_init", "vault-client", []string{"operator", "init"}, "unavailable", false, "operator init"},
		{"operator_unseal_help", "vault-client", []string{"operator", "unseal", "--help"}, "unseal", true, ""},
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
			out := stderr.String()
			if tc.explicit {
				require.Zero(t, code, "asking for help shouldn't fail")
				require.Empty(t, stderr.String(), "asking for help shouldn't write to stderr")
				out = stdout.String()
			} else {
				require.NotZero(t, code, "mistyped, incomplete, or unavailable commands should fail")
				require.Empty(t, stdout.String(), "output from a failed command should go to stderr")
			}
			require.NotEmpty(t, out, "we should print help or explain why the command isn't available")
			require.NotRegexp(t, excludedOperatorExample, out,
				"we shouldn't show examples for operator commands we don't support")

			var rows []string
			inListing := false
			for _, line := range strings.Split(out, "\n") {
				switch line {
				case "Common commands:", "Other commands:", "Subcommands:":
					inListing = true
				default:
					if inListing && strings.HasPrefix(line, "    ") {
						if fields := strings.Fields(line); len(fields) > 0 {
							rows = append(rows, fields[0])
						}
					}
				}
			}

			switch tc.kind {
			case "root":
				supported := []string{
					"read", "write", "delete", "list", "login", "agent", "proxy", "auth",
					"kv", "lease", "namespace", "secrets", "status", "token", "tpm",
					"print", "path-help", "patch", "unwrap", "hcp", "ssh", "operator",
				}
				supported = append(supported, entClientCommands...)
				require.Subset(t, rows, supported, "we should still list the commands we support")
				for _, hidden := range []string{
					"server", "audit", "policy", "plugin", "pki", "transit", "transform",
					"debug", "monitor", "version-history", "events", "license", "version",
				} {
					require.NotContains(t, rows, hidden, "%s shouldn't show up in the command list", hidden)
				}
			case "auth":
				require.ElementsMatch(t, []string{"disable", "enable", "help", "list", "move", "tune"}, rows,
					"we should still list all auth commands")
				for _, example := range []string{" auth list", " auth enable userpass", " auth help github"} {
					require.Contains(t, out, example, "we should keep examples for auth commands we support")
				}
			case "operator":
				require.Contains(t, out, "Usage: "+tc.binaryName+" operator <subcommand> [options] [args]",
					"operator help should use the binary name")
				require.ElementsMatch(t, []string{"key-status", "rotate", "unseal"}, rows,
					"we should only list the operator subcommands we support")
				var examples []string
				examplePattern := regexp.MustCompile(`(?m)^\s*\$\s+` + regexp.QuoteMeta(tc.binaryName) + `\s+operator\s+(\S+)`)
				for _, match := range examplePattern.FindAllStringSubmatch(out, -1) {
					examples = append(examples, match[1])
				}
				require.ElementsMatch(t, rows, examples,
					"operator examples should match the operator subcommands we list")
				require.Contains(t, out, "full Vault Server image", "operator help should point users to the full server image")
				require.NotContains(t, out, "not available in vault-client", "operator help isn't an unavailable command")
			case "unavailable":
				require.Contains(t, out, fmt.Sprintf("%q is not available in vault-client", tc.shim),
					"we should explain that %s isn't available", tc.shim)
				require.Contains(t, out, "This command requires the full Vault Server image.",
					"we should point users to the full server image")
				require.NotContains(t, out, "Usage:", "we shouldn't fall back to help")
				require.Empty(t, rows, "we shouldn't list subcommands of an unavailable command")
			case "unseal":
				require.Contains(t, out, "Usage: "+tc.binaryName+" operator unseal [options] [KEY]",
					"unseal help should be the real unseal help")
				require.NotContains(t, out, "not available in vault-client", "unseal is supported")
			default:
				require.FailNow(t, "unknown kind", tc.kind)
			}
		})
	}
}

// TestOperatorCommand_HelpProfiles checks that vault-client operator help shows
// examples for key-status, rotate, and unseal instead of init and step-down,
// and points users to the full Vault Server image for the rest. Vault's own
// operator help should stay exactly the same, even when we rename either binary.
func TestOperatorCommand_HelpProfiles(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		initFn     base.CommandInitFn
		binaryName string
		client     bool
	}{
		{"shared", InitClientCommands, "vault", false},
		{"shared_client_name", InitClientCommands, "vault-client", false},
		{"client", InitClientBinaryCommands, "vault-client", true},
		{"renamed_client", InitClientBinaryCommands, "vault", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ui := cli.NewMockUi()
			commands := tc.initFn(ui, ui, &base.RunOptions{
				BinaryName: tc.binaryName,
				Stdin:      strings.NewReader(""),
			})
			require.Contains(t, commands, "operator", "we should still be able to ask for operator help")
			cmd, err := commands["operator"]()
			require.NoError(t, err, "we should be able to create the operator command")
			require.IsType(t, &OperatorCommand{}, cmd, "the operator parent should remain a real help command")
			help := cmd.Help()
			require.Contains(t, help, "Usage: "+tc.binaryName+" operator <subcommand> [options] [args]",
				"operator help should use the selected binary name")

			if !tc.client {
				require.Equal(t, strings.TrimSpace(`
Usage: `+tc.binaryName+` operator <subcommand> [options] [args]

  This command groups subcommands for operators interacting with Vault. Most
  users will not need to interact with these commands. Here are a few examples
  of the operator commands:

  Initialize a new Vault cluster:

      $ `+tc.binaryName+` operator init

  Force a Vault to resign leadership in a cluster:

      $ `+tc.binaryName+` operator step-down

  Rotate Vault's underlying encryption key:

      $ `+tc.binaryName+` operator rotate

  Please see the individual subcommand help for detailed usage information.
`), help, "vault's operator help shouldn't change")
				return
			}

			require.Contains(t, help, "full Vault Server image", "operator help should point users to the full server image")
			for _, operation := range []string{"init", "step-down"} {
				require.NotContains(t, help, " operator "+operation,
					"vault-client help shouldn't show an example for %s", operation)
			}
			for _, operation := range []string{"key-status", "rotate", "unseal"} {
				require.Contains(t, help, "$ "+tc.binaryName+" operator "+operation+"\n",
					"operator help should show an example for %s", operation)
			}
		})
	}
}
