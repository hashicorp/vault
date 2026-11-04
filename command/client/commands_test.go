// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"io"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// TestInitClientCommands verifies that InitClientCommands registers exactly
// the vault-client command set from PRD VLT-2015, in both CE and enterprise
// builds, including the TPM commands and enterprise recover commands. Server-only
// commands such as audit, policy, license or operator init must never appear
// here; the vault server binary registers them on top of this set. It also checks that every factory builds a command that reads
// input from RunOptions.Stdin.
func TestInitClientCommands(t *testing.T) {
	t.Parallel()

	expected := []string{
		"agent",
		"agent generate-config",
		"auth",
		"auth disable",
		"auth enable",
		"auth help",
		"auth list",
		"auth move",
		"auth tune",
		"delete",
		"hcp connect",
		"hcp disconnect",
		"kv",
		"kv delete",
		"kv destroy",
		"kv enable-versioning",
		"kv get",
		"kv list",
		"kv metadata",
		"kv metadata delete",
		"kv metadata get",
		"kv metadata patch",
		"kv metadata put",
		"kv patch",
		"kv put",
		"kv rollback",
		"kv undelete",
		"lease",
		"lease lookup",
		"lease renew",
		"lease revoke",
		"list",
		"login",
		"namespace",
		"namespace create",
		"namespace delete",
		"namespace list",
		"namespace lock",
		"namespace lookup",
		"namespace patch",
		"namespace unlock",
		"operator",
		"operator key-status",
		"operator rotate",
		"operator unseal",
		"patch",
		"path-help",
		"print",
		"print token",
		"proxy",
		"read",
		"secrets",
		"secrets disable",
		"secrets enable",
		"secrets list",
		"secrets move",
		"secrets tune",
		"ssh",
		"status",
		"token",
		"token capabilities",
		"token create",
		"token lookup",
		"token renew",
		"token revoke",
		"tpm",
		"tpm attest",
		"tpm cleanup",
		"tpm ek",
		"tpm enroll",
		"tpm inspect",
		"tpm reattest",
		"unwrap",
		"version",
		"write",
	}
	expected = append(expected, entClientCommands...)
	slices.Sort(expected)

	ui := cli.NewMockUi()
	stdin := strings.NewReader("")
	commands := InitClientCommands(ui, ui, &base.RunOptions{Stdin: stdin})

	var names []string
	for name := range commands {
		names = append(names, name)
	}
	slices.Sort(names)
	require.Equal(t, expected, names)

	for name, factory := range commands {
		cmd, err := factory()
		require.NoError(t, err, name)
		require.NotNil(t, cmd, name)

		// agent and proxy use the server UI and do not take RunOptions.
		if bc, ok := cmd.(interface{ Stdin() io.Reader }); ok && name != "agent" && name != "proxy" {
			require.Same(t, stdin, bc.Stdin(), name)
		}
	}
}

// TestInitClientBinaryCommands checks that commands we leave out get placeholders
// without changing the shared commands or stdin wiring. It also checks that the
// supported operator subcommands (AC-1.3) keep their real implementations.
func TestInitClientBinaryCommands(t *testing.T) {
	t.Parallel()

	ui := cli.NewMockUi()
	stdin := strings.NewReader("")
	opts := &base.RunOptions{Stdin: stdin}
	shared := InitClientCommands(ui, ui, opts)
	binary := InitClientBinaryCommands(ui, ui, opts)
	expectedKeys := make(map[string]bool)
	for name := range shared {
		expectedKeys[name] = true
	}
	for _, name := range unavailableClientCommands {
		expectedKeys[name] = true
	}
	require.ElementsMatch(t, slices.Collect(maps.Keys(expectedKeys)), slices.Collect(maps.Keys(binary)))

	for name, factory := range binary {
		cmd, err := factory()
		require.NoError(t, err, name)
		require.NotNil(t, cmd, name)
		if slices.Contains(unavailableClientCommands, name) {
			require.IsType(t, &UnavailableCommand{}, cmd, name)
			unavailable := cmd.(*UnavailableCommand)
			require.Equal(t, name, unavailable.name)
			require.Same(t, ui, unavailable.ui, name)
			next, nextErr := factory()
			require.NoError(t, nextErr, name)
			require.NotSame(t, cmd, next, name)
			require.Equal(t, "Not included in this binary", cmd.Synopsis(), name)
			require.Contains(t, cmd.Help(), `"`+name+`"`, name)
			require.Empty(t, ui.OutputWriter.String(), "help must not write to the UI")
			require.Empty(t, ui.ErrorWriter.String(), "help must not write to the UI")
			continue
		}
		sharedCmd, sharedErr := shared[name]()
		require.NoError(t, sharedErr, name)
		require.IsType(t, sharedCmd, cmd, name)
		if bc, ok := cmd.(interface{ Stdin() io.Reader }); ok && name != "agent" && name != "proxy" {
			require.Same(t, stdin, bc.Stdin(), name)
		}
	}

	for name, want := range map[string]cli.Command{
		"operator":            &OperatorCommand{},
		"operator key-status": &OperatorKeyStatusCommand{},
		"operator rotate":     &OperatorRotateCommand{},
		"operator unseal":     &OperatorUnsealCommand{},
	} {
		factory, ok := binary[name]
		require.True(t, ok, name)
		cmd, err := factory()
		require.NoError(t, err, name)
		require.IsType(t, want, cmd, name)
	}
}
