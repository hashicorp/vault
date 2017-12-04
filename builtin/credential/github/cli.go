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
The GitHub credential provider allows you to authenticate with GitHub.
To use it, specify the "token" parameter. The value should be a personal access
token for your GitHub account. You can generate a personal access token on your
account settings page on GitHub.

    Example: vault auth -method=github token=<token>

Key/Value Pairs:

    mount=github      The mountpoint for the GitHub credential provider.
                      Defaults to "github"

    token=<token>     The GitHub personal access token for authentication.
	`

	return strings.TrimSpace(help)
}
