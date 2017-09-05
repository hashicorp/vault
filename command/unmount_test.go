package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testUnmountCommand(tb testing.TB) (*cli.MockUi, *UnmountCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &UnmountCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestUnmountCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"empty",
			nil,
			"Missing PATH!",
			1,
		},
		{
			"slash",
			[]string{"/"},
			"Missing PATH!",
			1,
		},
		{
			"not_real",
			[]string{"not_real"},
			"Success! Unmounted the secret backend (if it existed) at: not_real/",
			0,
		},
		{
			"default",
			[]string{"secret"},
			"Success! Unmounted the secret backend (if it existed) at: secret/",
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

				ui, cmd := testUnmountCommand(t)
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

		if err := client.Sys().Mount("integration_unmount/", &api.MountInput{
			Type: "generic",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testUnmountCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"integration_unmount/",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Unmounted the secret backend (if it existed) at: integration_unmount/"
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

		ui, cmd := testUnmountCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"pki/",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error unmounting pki/: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testUnmountCommand(t)
		assertNoTabs(t, cmd)
	})
}
