package command

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testPolicyDeleteCommand(tb testing.TB) (*cli.MockUi, *PolicyDeleteCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PolicyDeleteCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPolicyDeleteCommand_Run(t *testing.T) {
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
			[]string{"foo", "bar"},
			"Too many arguments",
			1,
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

				ui, cmd := testPolicyDeleteCommand(t)
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

		policy := `path "secret/" {}`
		if err := client.Sys().PutPolicy("my-policy", policy); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testPolicyDeleteCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"my-policy",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Deleted policy: my-policy"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		policies, err := client.Sys().ListPolicies()
		if err != nil {
			t.Fatal(err)
		}

		list := []string{"default", "root"}
		if !reflect.DeepEqual(policies, list) {
			t.Errorf("expected %q to be %q", policies, list)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPolicyDeleteCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"my-policy",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error deleting my-policy: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPolicyDeleteCommand(t)
		assertNoTabs(t, cmd)
	})
}
