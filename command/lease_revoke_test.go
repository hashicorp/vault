package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testLeaseRevokeCommand(tb testing.TB) (*cli.MockUi, *LeaseRevokeCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &LeaseRevokeCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestLeaseRevokeCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"force_without_prefix",
			[]string{"-force"},
			"requires also specifying -prefix",
			1,
		},
		{
			"single",
			nil,
			"All revocation operations queued successfully",
			0,
		},
		{
			"single_sync",
			[]string{"-sync"},
			"Success",
			0,
		},
		{
			"force_prefix",
			[]string{"-force", "-prefix"},
			"Success",
			0,
		},
		{
			"prefix",
			[]string{"-prefix"},
			"All revocation operations queued successfully",
			0,
		},
		{
			"prefix_sync",
			[]string{"-prefix", "-sync"},
			"Success",
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

				if err := client.Sys().Mount("secret-leased", &api.MountInput{
					Type: "generic-leased",
				}); err != nil {
					t.Fatal(err)
				}

				path := "secret-leased/revoke/" + tc.name
				data := map[string]interface{}{
					"key":   "value",
					"lease": "1m",
				}
				if _, err := client.Logical().Write(path, data); err != nil {
					t.Fatal(err)
				}
				secret, err := client.Logical().Read(path)
				if err != nil {
					t.Fatal(err)
				}

				ui, cmd := testLeaseRevokeCommand(t)
				cmd.client = client

				tc.args = append(tc.args, secret.LeaseID)
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

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testLeaseRevokeCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"foo/bar",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error revoking lease foo/bar: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testLeaseRevokeCommand(t)
		assertNoTabs(t, cmd)
	})
}
