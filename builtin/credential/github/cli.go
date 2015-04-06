package github

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "github"
	}

	token, ok := m["token"]
	if !ok {
		return "", fmt.Errorf("'token' var must be set")
	}

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"token": token,
	})
	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("empty response from credential provider")
	}

	return secret.Auth.ClientToken, nil
}

func (h *CLIHandler) Help() string {
	help := `
The GitHub credential provider allows you to authenticate with GitHub.
To use it, specify the "token" var with the "-var" flag. The value should
be a personal access token for your GitHub account. You can generate a personal
access token on your account settings page on GitHub.

    Example: vault auth -method=github -var="token=<token>"

	`

	return strings.TrimSpace(help)
}
