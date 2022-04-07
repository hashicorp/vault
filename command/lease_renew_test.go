package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testLeaseRenewCommand(tb testing.TB) (*cli.MockUi, *LeaseRenewCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &LeaseRenewCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

// testLeaseRenewCommandMountAndLease mounts a leased secret backend and returns
// the leaseID of an item.
func testLeaseRenewCommandMountAndLease(tb testing.TB, client *api.Client) string {
	if err := client.Sys().Mount("testing", &api.MountInput{
		Type: "generic-leased",
	}); err != nil {
		tb.Fatal(err)
	}

	if _, err := client.Logical().Write("testing/foo", map[string]interface{}{
		"key":   "value",
		"lease": "5m",
	}); err != nil {
		tb.Fatal(err)
	}

	// Read the secret back to get the leaseID
	secret, err := client.Logical().Read("testing/foo")
	if err != nil {
		tb.Fatal(err)
	}
	if secret == nil || secret.LeaseID == "" {
		tb.Fatalf("missing secret or lease: %#v", secret)
	}

	return secret.LeaseID
}

func TestLeaseRenewCommand_Run(t *testing.T) {
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
			"Missing ID!",
			1,
		},
		{
			"increment",
			[]string{"-increment", "60s"},
			"foo",
			0,
		},
		{
			"increment_no_suffix",
			[]string{"-increment", "60"},
			"foo",
			0,
		},
	}

	t.Run("group", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()

				leaseID := testLeaseRenewCommandMountAndLease(t, client)

				ui, cmd := testLeaseRenewCommand(t)
				cmd.client = client

				if tc.args != nil {
					tc.args = append(tc.args, leaseID)
				}
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

		leaseID := testLeaseRenewCommandMountAndLease(t, client)

		_, cmd := testLeaseRenewCommand(t)
		cmd.client = client

		code := cmd.Run([]string{leaseID})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testLeaseRenewCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"foo/bar",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error renewing foo/bar: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testLeaseRenewCommand(t)
		assertNoTabs(t, cmd)
	})
}
