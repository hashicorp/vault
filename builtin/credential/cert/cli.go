package cert

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	var data struct {
		Mount string `mapstructure:"mount"`
	}
	if err := mapstructure.WeakDecode(m, &data); err != nil {
		return "", err
	}

	if data.Mount == "" {
		data.Mount = "cert"
	}

	path := fmt.Sprintf("auth/%s/login", data.Mount)
	secret, err := c.Logical().Write(path, nil)
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
The "cert" credential provider allows you to authenticate with a
client certificate. No other authentication materials are needed.

    Example: vault auth -method=cert \
                        -client-cert=/path/to/cert.pem \
                        -client-key=/path/to/key.pem

	`

	return strings.TrimSpace(help)
}
