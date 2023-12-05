// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testAuthMoveCommand(tb testing.TB) (*cli.MockUi, *AuthMoveCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuthMoveCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestAuthMoveCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			[]string{},
			"Not enough arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"foo", "bar", "baz"},
			"Too many arguments",
			1,
		},
		{
			"non_existent",
			[]string{"not_real", "over_here"},
			"Error moving auth method not_real/ to over_here/",
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

				ui, cmd := testAuthMoveCommand(t)
				cmd.client = client

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

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testAuthMoveCommand(t)
		cmd.client = client

		if err := client.Sys().EnableAuthWithOptions("my-auth", &api.EnableAuthOptions{
			Type: "userpass",
		}); err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			"auth/my-auth/", "auth/my-auth-2/",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Finished moving auth method auth/my-auth/ to auth/my-auth-2/"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		mounts, err := client.Sys().ListAuth()
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := mounts["my-auth-2/"]; !ok {
			t.Errorf("expected mount at my-auth-2/: %#v", mounts)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testAuthMoveCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"auth/my-auth/", "auth/my-auth-2/",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error moving auth method auth/my-auth/ to auth/my-auth-2/:"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuthMoveCommand(t)
		assertNoTabs(t, cmd)
	})
}
