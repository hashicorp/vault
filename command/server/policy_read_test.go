// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
)

func testPolicyReadCommand(tb testing.TB) (*cli.MockUi, *clientcmd.PolicyReadCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &clientcmd.PolicyReadCommand{
		BaseCommand: &base.BaseCommand{
			UI: ui,
		},
	}
}

func TestPolicyReadCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"foo", "bar"},
			"Too many arguments",
			1,
		},
		{
			"no_policy_exists",
			[]string{"not-a-real-policy"},
			"No policy named",
			2,
		},
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()

				ui, cmd := testPolicyReadCommand(t)
				cmd.SetClient(client)

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("expected %d to be %d", code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, tc.out) {
					t.Errorf("expected %q to contain %q", combined, tc.out)
				}
			})
		}
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		policy := `path "secret/" {}`
		if err := client.Sys().PutPolicy("my-policy", policy); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testPolicyReadCommand(t)
		cmd.SetClient(client)

		code := cmd.Run([]string{
			"my-policy",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, policy) {
			t.Errorf("expected %q to contain %q", combined, policy)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPolicyReadCommand(t)
		cmd.SetClient(client)

		code := cmd.Run([]string{
			"my-policy",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error reading policy named my-policy: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPolicyReadCommand(t)
		assertNoTabs(t, cmd)
	})
}
