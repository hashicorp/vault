package cert

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	var data struct {
		Mount string `mapstructure:"mount"`
		Name  string `mapstructure:"name"`
	}
	if err := mapstructure.WeakDecode(m, &data); err != nil {
		return nil, err
	}

	if data.Mount == "" {
		data.Mount = "cert"
	}

	options := map[string]interface{}{
		"name": data.Name,
	}
	path := fmt.Sprintf("auth/%s/login", data.Mount)
	secret, err := c.Logical().Write(path, options)
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
The "cert" credential provider allows you to authenticate with a
client certificate. No other authentication materials are needed.
Optionally, you may specify the specific certificate role to
authenticate against with the "name" parameter.

    Example: vault auth -method=cert \
                        -client-cert=/path/to/cert.pem \
                        -client-key=/path/to/key.pem
                        name=cert1

	`

	return strings.TrimSpace(help)
}
