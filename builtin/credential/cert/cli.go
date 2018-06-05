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
Usage: vault login -method=cert [CONFIG K=V...]

  The certificate auth method allows users to authenticate with a
  client certificate passed with the request. The -client-cert and -client-key
  flags are included with the "vault login" command, NOT as configuration to the
  auth method.

  Authenticate using a local client certificate:

      $ vault login -method=cert -client-cert=cert.pem -client-key=key.pem

Configuration:

  name=<string>
      Certificate role to authenticate against.
`

	return strings.TrimSpace(help)
}
