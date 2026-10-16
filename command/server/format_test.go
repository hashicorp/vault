// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"strings"
	"testing"

	base "github.com/hashicorp/vault/command/base"
)

// Test_Format_Parsing verifies that the root command runner applies valid
// output formats and rejects unsupported formats.
func Test_Format_Parsing(t *testing.T) {
	t.Setenv(base.EnvVaultCLINoColor, "")
	t.Setenv(base.EnvVaultFormat, "")

	tests := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{name: "format", args: []string{"token", "renew", "-format", "json"}, out: "{", code: 0},
		{
			name: "format_bad",
			args: []string{"token", "renew", "-format", "nope-not-real"},
			out:  "Invalid output format",
			code: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, closer := testVaultServer(t)
			defer closer()

			stdout := bytes.NewBuffer(nil)
			stderr := bytes.NewBuffer(nil)
			runOpts := &base.RunOptions{Stdout: stdout, Stderr: stderr, Client: client}

			// Login with the token so we can renew-self.
			token, _ := testTokenAndAccessor(t, client)
			client.SetToken(token)

			code := RunCustom(test.args, runOpts)
			if code != test.code {
				t.Errorf("expected %d to be %d", code, test.code)
			}
			combined := stdout.String() + stderr.String()
			if !strings.Contains(combined, test.out) {
				t.Errorf("expected %q to contain %q", combined, test.out)
			}
		})
	}
}
