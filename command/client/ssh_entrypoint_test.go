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

// TestRunCustom_SSH verifies that SSH reaches its own help and validation through
// the client registry, with executable names suitable for either CLI.
func TestRunCustom_SSH(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name       string
		binaryName string
		args       []string
		code       int
	}{
		{"client long help", "vault-client", []string{"ssh", "--help"}, 0},
		{"client short help", "vault-client", []string{"ssh", "-h"}, 0},
		{"server help branding", "vault", []string{"ssh", "--help"}, 0},
		{"missing destination", "vault-client", []string{"ssh"}, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer
			code := RunCustom(test.args, &base.RunOptions{
				BinaryName: test.binaryName,
				Stdin:      strings.NewReader(""),
				Stdout:     &stdout,
				Stderr:     &stderr,
			})
			out := stdout.String() + stderr.String()
			require.Equal(t, test.code, code, "SSH entry-point output: %s", out)
			require.NotContains(t, out, "Common commands:", "SSH must not fall back to root usage")
			require.NotContains(t, out, "Other commands:", "SSH must not fall back to root usage")
			if test.code != 0 {
				require.Contains(t, stderr.String(), "Not enough arguments", "SSH must validate its destination")
				return
			}

			require.Contains(t, out, "Usage: "+test.binaryName+" ssh [options]", "SSH usage must name this CLI")
			for _, option := range []string{"SSH Options", "-role", "-mode", "-mount-point"} {
				require.Contains(t, out, option, "SSH help must describe %s", option)
			}
			examples := 0
			for _, line := range strings.Split(out, "\n") {
				if strings.HasPrefix(strings.TrimSpace(line), "$ ") {
					require.Contains(t, line, "$ "+test.binaryName+" ssh", "SSH examples must name this CLI")
					examples++
				}
			}
			require.Equal(t, 3, examples, "all SSH examples must be checked")
		})
	}
}
