package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testAuditDisableCommand(tb testing.TB) (*cli.MockUi, *AuditDisableCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuditDisableCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestAuditDisableCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			nil,
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
			"not_real",
			[]string{"not_real"},
			"Success! Disabled audit device (if it was enabled) at: not_real/",
			0,
		},
		{
			"default",
			[]string{"file"},
			"Success! Disabled audit device (if it was enabled) at: file/",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			if err := client.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
				Type: "file",
				Options: map[string]string{
					"file_path": "discard",
				},
			}); err != nil {
				t.Fatal(err)
			}

			ui, cmd := testAuditDisableCommand(t)
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

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuditWithOptions("integration_audit_disable", &api.EnableAuditOptions{
			Type: "file",
			Options: map[string]string{
				"file_path": "discard",
			},
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testAuditDisableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"integration_audit_disable/",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Disabled audit device (if it was enabled) at: integration_audit_disable/"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		mounts, err := client.Sys().ListMounts()
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := mounts["integration_audit_disable"]; ok {
			t.Errorf("expected mount to not exist: %#v", mounts)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testAuditDisableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"file",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error disabling audit device: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuditDisableCommand(t)
		assertNoTabs(t, cmd)
	})
}
