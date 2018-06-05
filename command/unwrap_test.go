package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testUnwrapCommand(tb testing.TB) (*cli.MockUi, *UnwrapCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &UnwrapCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func testUnwrapWrappedToken(tb testing.TB, client *api.Client, data map[string]interface{}) string {
	tb.Helper()

	wrapped, err := client.Logical().Write("sys/wrapping/wrap", data)
	if err != nil {
		tb.Fatal(err)
	}
	if wrapped == nil || wrapped.WrapInfo == nil || wrapped.WrapInfo.Token == "" {
		tb.Fatalf("missing wrap info: %v", wrapped)
	}
	return wrapped.WrapInfo.Token
}

func TestUnwrapCommand_Run(t *testing.T) {
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
			"default",
			nil, // Token comes in the test func
			"bar",
			0,
		},
		{
			"field",
			[]string{"-field", "foo"},
			"bar",
			0,
		},
		{
			"field_not_found",
			[]string{"-field", "not-a-real-field"},
			"not present in secret",
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

				wrappedToken := testUnwrapWrappedToken(t, client, map[string]interface{}{
					"foo": "bar",
				})

				ui, cmd := testUnwrapCommand(t)
				cmd.client = client

				tc.args = append(tc.args, wrappedToken)
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

		ui, cmd := testUnwrapCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"foo",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error unwrapping: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	// This test needs its own client and server because it modifies the client
	// to the wrapping token
	t.Run("local_token", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		wrappedToken := testUnwrapWrappedToken(t, client, map[string]interface{}{
			"foo": "bar",
		})

		ui, cmd := testUnwrapCommand(t)
		cmd.client = client
		cmd.client.SetToken(wrappedToken)

		// Intentionally don't pass the token here - it should use the local token
		code := cmd.Run([]string{})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, "bar") {
			t.Errorf("expected %q to contain %q", combined, "bar")
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testUnwrapCommand(t)
		assertNoTabs(t, cmd)
	})
}
