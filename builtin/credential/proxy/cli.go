package proxy

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
		Role  string `mapstructure:"role"`
	}
	if err := mapstructure.WeakDecode(m, &data); err != nil {
		return nil, err
	}

	if data.Mount == "" {
		data.Mount = "proxy"
	}

	options := map[string]interface{}{
		"role": data.Role,
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
Usage: vault login -method=proxy [CONFIG K=V...]

  The proxy auth method allows users to authenticate to a vault cluster that's
  fronted by a proxy that handles authentication.  Note that the -client-cert
  and -client-key flags are included with the "vault login" command, NOT as
  configuration to the auth method.

  Sample invocation to login to the "users" role:

      $ vault login -method=proxy -client-cert=cert.pem -client-key=key.pem role=users

Configuration:

  role=<string>
      Role name to authenticate against.
`

	return strings.TrimSpace(help)
}
