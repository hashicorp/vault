package github

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "github"
	}

	token, ok := m["token"]
	if !ok {
		if token = os.Getenv("VAULT_AUTH_GITHUB_TOKEN"); token == "" {
			return nil, fmt.Errorf("GitHub token should be provided either as 'value' for 'token' key,\nor via an env var VAULT_AUTH_GITHUB_TOKEN")
		}
	}

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"token": token,
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
Usage: vault auth -method=github [CONFIG K=V...]

  The GitHub authentication provider allows users to authenticate using a
  GitHub personal access token. Users can generate a personal access token
  from the settings page on their GitHub account.

  Authenticate using a GitHub token:

      $ vault auth -method=github token=abcd1234

Configuration:

  mount=<string>
      Path where the GitHub credential provider is mounted. This is usually
      provided via the -path flag in the "vault auth" command, but it can be
      specified here as well. If specified here, it takes precedence over
      the value for -path. The default value is "github".

  token=<string>
      GitHub personal access token to use for authentication.
`

	return strings.TrimSpace(help)
}
