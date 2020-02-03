package centrify

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	pwd "github.com/hashicorp/vault/sdk/helper/password"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "centrify"
	}

	username, ok := m["username"]
	if !ok {
		return nil, fmt.Errorf("'username' not supplied")
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
		"username": username,
		"password": password,
	}

	mode, ok := m["mode"]
	if ok {
		data["mode"] = mode
	}

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, data)
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
The "centrify" credential provider allows you to authenticate with
a username and password. To use it, specify the "username" and "password"
parameters. If password is not provided on the command line, it will be
read from stdin.`

	return strings.TrimSpace(help)
}
