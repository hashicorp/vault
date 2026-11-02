// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"strings"
	"testing"

	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// TestRunCustom_ServerNotIncluded verifies that vault-client, which does not
// ship the Vault server, answers `server` with an error that links to where
// the vault binary can be obtained and exits 1. It also checks that the
// hidden command never shows up in the usage or help output.
func TestRunCustom_ServerNotIncluded(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := RunCustom([]string{"server"}, &base.RunOptions{
		Stdin:  strings.NewReader(""),
		Stdout: &stdout,
		Stderr: &stderr,
	})
	require.Equal(t, 1, code)
	require.Contains(t, stderr.String(), "does not include the Vault server")
	require.Contains(t, stderr.String(), serverInstallURL)

	for _, args := range [][]string{nil, {"-h"}, {"--help"}} {
		var stdout, stderr bytes.Buffer
		RunCustom(args, &base.RunOptions{
			Stdin:  strings.NewReader(""),
			Stdout: &stdout,
			Stderr: &stderr,
		})
		out := stdout.String() + stderr.String()
		require.Contains(t, out, "Usage:", "args: %v", args)
		require.NotContains(t, out, "server", "args: %v", args)
	}
}
