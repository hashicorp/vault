package command

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testPolicyWriteCommand(tb testing.TB) (*cli.MockUi, *PolicyWriteCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PolicyWriteCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func testPolicyWritePolicyContents(tb testing.TB) []byte {
	return bytes.TrimSpace([]byte(`
path "secret/" {
  capabilities = ["read"]
}
	`))
}

func TestPolicyWriteCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"foo", "bar", "baz"},
			"Too many arguments",
			1,
		},
		{
			"not_enough_args",
			[]string{"foo"},
			"Not enough arguments",
			1,
		},
		{
			"bad_file",
			[]string{"my-policy", "/not/a/real/path.hcl"},
			"Error opening policy file",
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

				ui, cmd := testPolicyWriteCommand(t)
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

	t.Run("file", func(t *testing.T) {
		t.Parallel()

		policy := testPolicyWritePolicyContents(t)
		f, err := ioutil.TempFile("", "vault-policy-write")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write(policy); err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testPolicyWriteCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"my-policy", f.Name(),
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Uploaded policy: my-policy"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		policies, err := client.Sys().ListPolicies()
		if err != nil {
			t.Fatal(err)
		}

		list := []string{"default", "my-policy", "root"}
		if !reflect.DeepEqual(policies, list) {
			t.Errorf("expected %q to be %q", policies, list)
		}
	})

	t.Run("stdin", func(t *testing.T) {
		t.Parallel()

		stdinR, stdinW := io.Pipe()
		go func() {
			policy := testPolicyWritePolicyContents(t)
			stdinW.Write(policy)
			stdinW.Close()
		}()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testPolicyWriteCommand(t)
		cmd.client = client
		cmd.testStdin = stdinR

		code := cmd.Run([]string{
			"my-policy", "-",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Uploaded policy: my-policy"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		policies, err := client.Sys().ListPolicies()
		if err != nil {
			t.Fatal(err)
		}

		list := []string{"default", "my-policy", "root"}
		if !reflect.DeepEqual(policies, list) {
			t.Errorf("expected %q to be %q", policies, list)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPolicyWriteCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"my-policy", "-",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error uploading policy: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPolicyWriteCommand(t)
		assertNoTabs(t, cmd)
	})
}
