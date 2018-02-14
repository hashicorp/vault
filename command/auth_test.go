package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"

	credToken "github.com/hashicorp/vault/builtin/credential/token"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/command/token"
)

func testAuthCommand(tb testing.TB) (*cli.MockUi, *AuthCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuthCommand{
		BaseCommand: &BaseCommand{
			UI: ui,

			// Override to our own token helper
			tokenHelper: token.NewTestingTokenHelper(),
		},
		Handlers: map[string]LoginHandler{
			"token":    &credToken.CLIHandler{},
			"userpass": &credUserpass.CLIHandler{},
		},
	}
}

func TestAuthCommand_Run(t *testing.T) {
	t.Parallel()

	// TODO: remove in 0.9.0
	t.Run("deprecated_methods", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testAuthCommand(t)
		cmd.client = client

		// vault auth -methods -> vault auth list
		code := cmd.Run([]string{"-methods"})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}
		stdout, stderr := ui.OutputWriter.String(), ui.ErrorWriter.String()

		if expected := "WARNING!"; !strings.Contains(stderr, expected) {
			t.Errorf("expected %q to contain %q", stderr, expected)
		}

		if expected := "token/"; !strings.Contains(stdout, expected) {
			t.Errorf("expected %q to contain %q", stdout, expected)
		}
	})

	t.Run("deprecated_method_help", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testAuthCommand(t)
		cmd.client = client

		// vault auth -method=foo -method-help -> vault auth help foo
		code := cmd.Run([]string{
			"-method=userpass",
			"-method-help",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}
		stdout, stderr := ui.OutputWriter.String(), ui.ErrorWriter.String()

		if expected := "WARNING!"; !strings.Contains(stderr, expected) {
			t.Errorf("expected %q to contain %q", stderr, expected)
		}

		if expected := "vault login"; !strings.Contains(stdout, expected) {
			t.Errorf("expected %q to contain %q", stdout, expected)
		}
	})

	t.Run("deprecated_login", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("my-auth", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/my-auth/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testAuthCommand(t)
		cmd.client = client

		// vault auth ARGS -> vault login ARGS
		code := cmd.Run([]string{
			"-method", "userpass",
			"-path", "my-auth",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}
		stdout, stderr := ui.OutputWriter.String(), ui.ErrorWriter.String()

		if expected := "WARNING!"; !strings.Contains(stderr, expected) {
			t.Errorf("expected %q to contain %q", stderr, expected)
		}

		if expected := "Success! You are now authenticated."; !strings.Contains(stdout, expected) {
			t.Errorf("expected %q to contain %q", stdout, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuthCommand(t)
		assertNoTabs(t, cmd)
	})
}
