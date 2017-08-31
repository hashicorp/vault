package okta

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	pwd "github.com/hashicorp/vault/helper/password"
)

// CLIHandler struct
type CLIHandler struct{}

// Auth cli method
func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "okta"
	}

	username, ok := m["username"]
	if !ok {
		return nil, fmt.Errorf("'username' var must be set")
	}
	password, ok := m["password"]
	if !ok {
		fmt.Printf("Password (will be hidden): ")
		var err error
		password, err = pwd.Read(os.Stdin)
		fmt.Println()
		if err != nil {
			return nil, err
		}
	}

	data := map[string]interface{}{
		"password": password,
	}

	path := fmt.Sprintf("auth/%s/login/%s", mount, username)
	secret, err := c.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("empty response from credential provider")
	}

	return secret, nil
}

// Help method for okta cli
func (h *CLIHandler) Help() string {
	help := `
The Okta credential provider allows you to authenticate with Okta.
To use it, first configure it through the "config" endpoint, and then
login by specifying username and password. If password is not provided
on the command line, it will be read from stdin.

    Example: vault auth -method=okta username=john

    `

	return strings.TrimSpace(help)
}
