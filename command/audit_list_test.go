package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testAuditListCommand(tb testing.TB) (*cli.MockUi, *AuditListCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuditListCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestAuditListCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"foo"},
			"Too many arguments",
			1,
		},
		{
			"lists",
			nil,
			"Path",
			0,
		},
		{
			"detailed",
			[]string{"-detailed"},
			"Options",
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

			ui, cmd := testAuditListCommand(t)
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

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testAuditListCommand(t)
		cmd.client = client

		code := cmd.Run([]string{})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error listing audits: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuditListCommand(t)
		assertNoTabs(t, cmd)
	})
}
