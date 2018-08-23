package alicloud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault-plugin-auth-alicloud/tools"
	"github.com/hashicorp/vault/api"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "alicloud"
	}
	role := m["role"]

	loginData, err := tools.GenerateLoginData(m["access_key"], m["secret_key"], m["security_token"], m["region"])
	if err != nil {
		return nil, err
	}

	loginData["role"] = role
	path := fmt.Sprintf("auth/%s/login", mount)

	secret, err := c.Logical().Write(path, loginData)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.New("empty response from credential provider")
	}
	return secret, nil
}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=alicloud [CONFIG K=V...]

  The AliCloud auth method allows users to authenticate with AliCloud RAM
  credentials.

  The AliCloud RAM credentials may be specified explicitly via the command line:

      $ vault login -method=alicloud access_key=... secret_key=... security_token=... region=...

Configuration:

  access_key=<string>
      Explicit AliCloud access key ID

  secret_key=<string>
      Explicit AliCloud secret access key

  security_token=<string>
      Explicit AliCloud security token

  region=<string>
	  Explicit AliCloud region

  mount=<string>
      Path where the AliCloud credential method is mounted. This is usually provided
      via the -path flag in the "vault login" command, but it can be specified
      here as well. If specified here, it takes precedence over the value for
      -path. The default value is "alicloud".

  role=<string>
      Name of the role to request a token against
`

	return strings.TrimSpace(help)
}
