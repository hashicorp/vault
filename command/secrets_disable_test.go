package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testSecretsDisableCommand(tb testing.TB) (*cli.MockUi, *SecretsDisableCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &SecretsDisableCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestSecretsDisableCommand_Run(t *testing.T) {
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
			[]string{"foo", "bar"},
			"Too many arguments",
			1,
		},
		{
			"not_real",
			[]string{"not_real"},
			"Success! Disabled the secrets engine (if it existed) at: not_real/",
			0,
		},
		{
			"default",
			[]string{"secret"},
			"Success! Disabled the secrets engine (if it existed) at: secret/",
			0,
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

				ui, cmd := testSecretsDisableCommand(t)
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

		if err := client.Sys().Mount("my-secret/", &api.MountInput{
			Type: "generic",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testSecretsDisableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"my-secret/",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Disabled the secrets engine (if it existed) at: my-secret/"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		mounts, err := client.Sys().ListMounts()
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := mounts["integration_unmount"]; ok {
			t.Errorf("expected mount to not exist: %#v", mounts)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testSecretsDisableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"pki/",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error disabling secrets engine at pki/: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testSecretsDisableCommand(t)
		assertNoTabs(t, cmd)
	})
}
