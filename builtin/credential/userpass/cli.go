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
		fmt.Printf("Password (will be hidden): ")
		password, err := pwd.Read(os.Stdin)
		fmt.Println()
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
The "userpass"/"radius" credential provider allows you to authenticate with
a username and password. To use it, specify the "username" and "password"
parameters. If password is not provided on the command line, it will be
read from stdin.

If multi-factor authentication (MFA) is enabled, a "method" and/or "passcode"
may be provided depending on the MFA backend enabled. To check
which MFA backend is in use, read "auth/[mount]/mfa_config".

    Example: vault auth -method=userpass \
      username=<user> \
      password=<password>

	`

	return strings.TrimSpace(help)
}
