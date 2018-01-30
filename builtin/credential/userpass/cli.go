package userpass

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	pwd "github.com/hashicorp/vault/helper/password"
	"github.com/mitchellh/mapstructure"
)

type CLIHandler struct {
	DefaultMount string
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	var data struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Mount    string `mapstructure:"mount"`
		Method   string `mapstructure:"method"`
		Passcode string `mapstructure:"passcode"`
	}
	if err := mapstructure.WeakDecode(m, &data); err != nil {
		return nil, err
	}

	if data.Username == "" {
		return nil, fmt.Errorf("'username' must be specified")
	}
	if data.Password == "" {
		fmt.Fprintf(os.Stderr, "Password (will be hidden): ")
		password, err := pwd.Read(os.Stdin)
		fmt.Fprintf(os.Stderr, "\n")
		if err != nil {
			return nil, err
		}
		data.Password = password
	}
	if data.Mount == "" {
		data.Mount = h.DefaultMount
	}

	options := map[string]interface{}{
		"password": data.Password,
	}
	if data.Method != "" {
		options["method"] = data.Method
	}
	if data.Passcode != "" {
		options["passcode"] = data.Passcode
	}

	path := fmt.Sprintf("auth/%s/login/%s", data.Mount, data.Username)
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
Usage: vault login -method=userpass [CONFIG K=V...]

  The userpass auth method allows users to authenticate using Vault's
  internal user database.

  If MFA is enabled, a "method" and/or "passcode" may be required depending on
  the MFA method. To check which MFA is in use, run:

      $ vault read auth/<mount>/mfa_config

  Authenticate as "sally":

      $ vault login -method=userpass username=sally
      Password (will be hidden):

  Authenticate as "bob":

      $ vault login -method=userpass username=bob password=password

Configuration:

  method=<string>
      MFA method.

  passcode=<string>
      MFA OTP/passcode.

  password=<string>
      Password to use for authentication. If not provided, the CLI will prompt
      for this on stdin.

  username=<string>
      Username to use for authentication.
`

	return strings.TrimSpace(help)
}
