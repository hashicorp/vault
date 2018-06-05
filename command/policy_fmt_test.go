package command

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testPolicyFmtCommand(tb testing.TB) (*cli.MockUi, *PolicyFmtCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PolicyFmtCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPolicyFmtCommand_Run(t *testing.T) {
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
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()

				ui, cmd := testPolicyFmtCommand(t)
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

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		policy := strings.TrimSpace(`
path "secret" {
  capabilities  =           ["create",    "update","delete"]

}
`)

		f, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		if _, err := f.Write([]byte(policy)); err != nil {
			t.Fatal(err)
		}
		f.Close()

		client, closer := testVaultServer(t)
		defer closer()

		_, cmd := testPolicyFmtCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			f.Name(),
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := strings.TrimSpace(`
path "secret" {
  capabilities = ["create", "update", "delete"]
}
`) + "\n"

		contents, err := ioutil.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		if string(contents) != expected {
			t.Errorf("expected %q to be %q", string(contents), expected)
		}
	})

	t.Run("bad_hcl", func(t *testing.T) {
		t.Parallel()

		policy := `dafdaf`

		f, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		if _, err := f.Write([]byte(policy)); err != nil {
			t.Fatal(err)
		}
		f.Close()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testPolicyFmtCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			f.Name(),
		})
		if exp := 1; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		stderr := ui.ErrorWriter.String()
		expected := "failed to parse policy"
		if !strings.Contains(stderr, expected) {
			t.Errorf("expected %q to include %q", stderr, expected)
		}
	})

	t.Run("bad_policy", func(t *testing.T) {
		t.Parallel()

		policy := `banana "foo" {}`

		f, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		if _, err := f.Write([]byte(policy)); err != nil {
			t.Fatal(err)
		}
		f.Close()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testPolicyFmtCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			f.Name(),
		})
		if exp := 1; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		stderr := ui.ErrorWriter.String()
		expected := "failed to parse policy"
		if !strings.Contains(stderr, expected) {
			t.Errorf("expected %q to include %q", stderr, expected)
		}
	})

	t.Run("bad_policy", func(t *testing.T) {
		t.Parallel()

		policy := `path "secret/" { capabilities = ["bogus"] }`

		f, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		if _, err := f.Write([]byte(policy)); err != nil {
			t.Fatal(err)
		}
		f.Close()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testPolicyFmtCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			f.Name(),
		})
		if exp := 1; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		stderr := ui.ErrorWriter.String()
		expected := "failed to parse policy"
		if !strings.Contains(stderr, expected) {
			t.Errorf("expected %q to include %q", stderr, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPolicyFmtCommand(t)
		assertNoTabs(t, cmd)
	})
}
