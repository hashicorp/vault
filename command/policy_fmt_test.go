// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/helper/random"
	"github.com/stretchr/testify/require"
)

func testPolicyFmtCommand(tb testing.TB) (*cli.MockUi, *PolicyFmtCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PolicyFmtCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPolicyFmtCommand_Run(t *testing.T) {
	testCases := map[string]struct {
		args      []string
		envVars   map[string]string
		policyArg string
		out       string
		expected  string
		code      int
	}{
		"not_enough_args": {
			args: []string{},
			out:  "Not enough arguments",
			code: 1,
		},
		"too_many_args": {
			args: []string{"foo", "bar"},
			out:  "Too many arguments",
			code: 1,
		},
		"default": {
			policyArg: `
path "secret" {
  capabilities  =           ["create",    "update","delete"]

}`,
			expected: `
path "secret" {
  capabilities = ["create", "update", "delete"]
}
`,
			code: 0,
		},
		"bad_hcl": {
			policyArg: `dafdaf`,
			out:       "failed to parse policy",
			code:      1,
		},
		"bad_policy": {
			policyArg: `banana "foo" {}`,
			out:       "failed to parse policy",
			code:      1,
		},
		"bad_policy2": {
			policyArg: `path "secret" { capabilities = ["bogus"] }`,
			out:       "failed to parse policy",
			code:      1,
		},
		// TODO (HCL_DUP_KEYS_DEPRECATION): remove this test once deprecation is fully done
		"hcl_duplicate_key_env_set": {
			policyArg: `
path "secret" {
  capabilities = ["create", "update", "delete"]
  capabilities = ["create"]
}
`,
			envVars: map[string]string{random.AllowHclDuplicatesEnvVar: "true"},
			code:    0,
			out:     "WARNING: Duplicate keys found in the provided policy, duplicate keys in HCL files are deprecated and will be forbidden in a future release.",
		},
		"hcl_duplicate_key_env_not_set": {
			policyArg: `
path "secret" {
  capabilities = ["create", "update", "delete"]
  capabilities = ["create"]
}
`,
			code: 1,
			out:  "failed to parse policy: The argument \"capabilities\" at 4:3 was already set. Each argument can only be defined once",
		},
		"hcl_duplicate_key_env_set_to_false": {
			policyArg: `
path "secret" {
  capabilities = ["create", "update", "delete"]
  capabilities = ["create"]
}
`,
			envVars: map[string]string{random.AllowHclDuplicatesEnvVar: "false"},
			code:    1,
			out:     "failed to parse policy: The argument \"capabilities\" at 4:3 was already set. Each argument can only be defined once",
		},
	}

	client, closer := testVaultServer(t)
	t.Cleanup(closer)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := require.New(t)

			for k, v := range tc.envVars {
				t.Setenv(k, v)
			}

			args := tc.args
			if tc.policyArg != "" {
				f := populateTempFile(t, "fmt-test-*.hcl", tc.policyArg)
				args = append(args, f.Name())
			}

			ui, cmd := testPolicyFmtCommand(t)
			cmd.client = client

			code := cmd.Run(args)
			r.Equal(tc.code, code)

			if tc.out != "" {
				t.Log(ui.ErrorWriter.String())
				r.Contains(ui.ErrorWriter.String(), tc.out)
			}

			if tc.expected != "" {
				contents, err := os.ReadFile(args[0])
				r.NoError(err)
				r.Equal(strings.TrimSpace(tc.expected), strings.TrimSpace(string(contents)))
			}
		})
	}
}

// TestPolicyFmtCommandNoTabs asserts the CLI help has no tab characters.
func TestPolicyFmtCommandNoTabs(t *testing.T) {
	t.Parallel()

	_, cmd := testPolicyFmtCommand(t)
	assertNoTabs(t, cmd)
}
