package ldap

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	pwd "github.com/hashicorp/vault/helper/password"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "ldap"
	}

	username, ok := m["username"]
	if !ok {
		username = usernameFromEnv()
		if username == "" {
			return nil, fmt.Errorf("'username' not supplied and neither 'LOGNAME' nor 'USER' env vars set")
		}
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

	mfa_method, ok := m["method"]
	if ok {
		data["method"] = mfa_method
	}
	mfa_passcode, ok := m["passcode"]
	if ok {
		data["passcode"] = mfa_passcode
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

func (h *CLIHandler) Help() string {
	help := `
The LDAP credential provider allows you to authenticate with LDAP.
To use it, first configure it through the "config" endpoint, and then
login by specifying username and password. If password is not provided
on the command line, it will be read from stdin.

If multi-factor authentication (MFA) is enabled, a "method" and/or "passcode"
may be provided depending on the MFA backend enabled. To check
which MFA backend is in use, read "auth/[mount]/mfa_config".

    Example: vault auth -method=ldap username=john

    `

	return strings.TrimSpace(help)
}

func usernameFromEnv() string {
	if logname := os.Getenv("LOGNAME"); logname != "" {
		return logname
	}
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	return ""
}
