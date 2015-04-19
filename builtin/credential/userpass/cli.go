package userpass

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	var data struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Mount    string `mapstructure:"mount"`
	}
	if err := mapstructure.WeakDecode(m, &data); err != nil {
		return "", err
	}

	if data.Username == "" || data.Password == "" {
		return "", fmt.Errorf("Both 'username' and 'password' must be specified")
	}
	if data.Mount == "" {
		data.Mount = "userpass"
	}

	path := fmt.Sprintf("auth/%s/login/%s", data.Mount, data.Username)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"password": data.Password,
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
The "userpass" credential provider allows you to authenticate with
a username and password. To use it, specify the "username" and "password"
vars with the "-var" flag.

    Example: vault auth -method=userpass \
      -var="username=<user>"
      -var="password=<password>"

	`

	return strings.TrimSpace(help)
}
