package oidc

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "oidc"
	}

	token, ok := m["token"]
	if !ok {
		if token = os.Getenv("VAULT_AUTH_OIDC_TOKEN"); token == "" {
			return "", fmt.Errorf("OpenID Connect (OIDC) token should be provided either as 'value' for 'token' key,\nor via an env var VAULT_AUTH_OIDC_TOKEN")
		}
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
The OpenID Connect credential provider allows you to authenticate with OIDC providers.
To use it, specify the "token" parameter. The value should be the user's identity
token for the OIDC provider. Usually you should get a new OIDC identity token with a third
party CLI tool.

    Example: vault auth -method=oidc token=<token>

Key/Value Pairs:

    mount=oidc      The mountpoint for the OIDC credential provider.
                      Defaults to "oidc"

    token=<token>     The OIDC identity token for authentication.
	`

	return strings.TrimSpace(help)
}
