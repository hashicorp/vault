package github

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/password"
)

type CLIHandler struct {
	// for tests
	testStdout io.Writer
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "github"
	}

	// Extract or prompt for token
	token := m["token"]
	if token == "" {
		token = os.Getenv("VAULT_AUTH_GITHUB_TOKEN")
	}
	if token == "" {
		// Override the output
		stdout := h.testStdout
		if stdout == nil {
			stdout = os.Stderr
		}

		var err error
		fmt.Fprintf(stdout, "GitHub Personal Access Token (will be hidden): ")
		token, err = password.Read(os.Stdin)
		fmt.Fprintf(stdout, "\n")
		if err != nil {
			if err == password.ErrInterrupted {
				return nil, fmt.Errorf("user interrupted")
			}

			return nil, errwrap.Wrapf("An error occurred attempting to "+
				"ask for a token. The raw error message is shown below, but usually "+
				"this is because you attempted to pipe a value into the command or "+
				"you are executing outside of a terminal (tty). If you want to pipe "+
				"the value, pass \"-\" as the argument to read from stdin. The raw "+
				"error was: {{err}}", err)
		}
	}

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"token": strings.TrimSpace(token),
	})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("empty response from credential provider")
	}

	return secret, nil
}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=github [CONFIG K=V...]

  The GitHub auth method allows users to authenticate using a GitHub
  personal access token. Users can generate a personal access token from the
  settings page on their GitHub account.

  Authenticate using a GitHub token:

      $ vault login -method=github token=abcd1234

Configuration:

  mount=<string>
      Path where the GitHub credential method is mounted. This is usually
      provided via the -path flag in the "vault login" command, but it can be
      specified here as well. If specified here, it takes precedence over the
      value for -path. The default value is "github".

  token=<string>
      GitHub personal access token to use for authentication. If not provided,
      Vault will prompt for the value.
`

	return strings.TrimSpace(help)
}
