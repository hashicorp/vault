package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"

	"github.com/hashicorp/vault/api"
	credToken "github.com/hashicorp/vault/builtin/credential/token"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/command/token"
	"github.com/hashicorp/vault/vault"
)

func testLoginCommand(tb testing.TB) (*cli.MockUi, *LoginCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &LoginCommand{
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

func TestLoginCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("custom_path", func(t *testing.T) {
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

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			"-method", "userpass",
			"-path", "my-auth",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! You are now authenticated."
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to be %q", combined, expected)
		}

		storedToken, err := tokenHelper.Get()
		if err != nil {
			t.Fatal(err)
		}

		if l, exp := len(storedToken), vault.TokenLength+2; l != exp {
			t.Errorf("expected token to be %d characters, was %d: %q", exp, l, storedToken)
		}
	})

	t.Run("no_store", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			TTL:      "30m",
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		_, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		// Ensure we have no token to start
		if storedToken, err := tokenHelper.Get(); err != nil || storedToken != "" {
			t.Errorf("expected token helper to be empty: %s: %q", err, storedToken)
		}

		code := cmd.Run([]string{
			"-no-store",
			token,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		storedToken, err := tokenHelper.Get()
		if err != nil {
			t.Fatal(err)
		}

		if exp := ""; storedToken != exp {
			t.Errorf("expected %q to be %q", storedToken, exp)
		}
	})

	t.Run("stores", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			TTL:      "30m",
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		_, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			token,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		storedToken, err := tokenHelper.Get()
		if err != nil {
			t.Fatal(err)
		}

		if storedToken != token {
			t.Errorf("expected %q to be %q", storedToken, token)
		}
	})

	t.Run("token_only", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			"-token-only",
			"-method", "userpass",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		// Verify only the token was printed
		token := ui.OutputWriter.String()
		if l, exp := len(token), vault.TokenLength+2; l != exp {
			t.Errorf("expected token to be %d characters, was %d: %q", exp, l, token)
		}

		// Verify the token was not stored
		if storedToken, err := tokenHelper.Get(); err != nil || storedToken != "" {
			t.Fatalf("expected token to not be stored: %s: %q", err, storedToken)
		}
	})

	t.Run("failure_no_store", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			"not-a-real-token",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error authenticating: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		if storedToken, err := tokenHelper.Get(); err != nil || storedToken != "" {
			t.Fatalf("expected token to not be stored: %s: %q", err, storedToken)
		}
	})

	t.Run("wrap_auto_unwrap", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		_, cmd := testLoginCommand(t)
		cmd.client = client

		// Set the wrapping ttl to 5s. We can't set this via the flag because we
		// override the client object before that particular flag is parsed.
		client.SetWrappingLookupFunc(func(string, string) string { return "5m" })

		code := cmd.Run([]string{
			"-method", "userpass",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		// Unset the wrapping
		client.SetWrappingLookupFunc(func(string, string) string { return "" })

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}
		token, err := tokenHelper.Get()
		if err != nil || token == "" {
			t.Fatalf("expected token from helper: %s: %q", err, token)
		}
		client.SetToken(token)

		// Ensure the resulting token is unwrapped
		secret, err := client.Auth().Token().LookupSelf()
		if err != nil {
			t.Error(err)
		}
		if secret == nil {
			t.Fatal("secret was nil")
		}

		if secret.WrapInfo != nil {
			t.Errorf("expected to be unwrapped: %#v", secret)
		}
	})

	t.Run("wrap_token_only", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		// Set the wrapping ttl to 5s. We can't set this via the flag because we
		// override the client object before that particular flag is parsed.
		client.SetWrappingLookupFunc(func(string, string) string { return "5m" })

		code := cmd.Run([]string{
			"-token-only",
			"-method", "userpass",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		// Unset the wrapping
		client.SetWrappingLookupFunc(func(string, string) string { return "" })

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}
		storedToken, err := tokenHelper.Get()
		if err != nil || storedToken != "" {
			t.Fatalf("expected token to not be stored: %s: %q", err, storedToken)
		}

		token := strings.TrimSpace(ui.OutputWriter.String())
		if token == "" {
			t.Errorf("expected %q to not be %q", token, "")
		}

		// Ensure the resulting token is, in fact, still wrapped.
		client.SetToken(token)
		secret, err := client.Logical().Unwrap("")
		if err != nil {
			t.Error(err)
		}
		if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
			t.Fatalf("expected secret to have auth: %#v", secret)
		}
	})

	t.Run("wrap_no_store", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		// Set the wrapping ttl to 5s. We can't set this via the flag because we
		// override the client object before that particular flag is parsed.
		client.SetWrappingLookupFunc(func(string, string) string { return "5m" })

		code := cmd.Run([]string{
			"-no-store",
			"-method", "userpass",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		// Unset the wrapping
		client.SetWrappingLookupFunc(func(string, string) string { return "" })

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}
		storedToken, err := tokenHelper.Get()
		if err != nil || storedToken != "" {
			t.Fatalf("expected token to not be stored: %s: %q", err, storedToken)
		}

		expected := "wrapping_token"
		output := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(output, expected) {
			t.Errorf("expected %q to contain %q", output, expected)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"token",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error authenticating: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testLoginCommand(t)
		assertNoTabs(t, cmd)
	})
}
