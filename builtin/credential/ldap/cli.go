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
		fmt.Fprintf(os.Stderr, "Password (will be hidden): ")
		var err error
		password, err = pwd.Read(os.Stdin)
		fmt.Fprintf(os.Stderr, "\n")
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
Usage: vault login -method=ldap [CONFIG K=V...]

  The LDAP auth method allows users to authenticate using LDAP or
  Active Directory.

  If MFA is enabled, a "method" and/or "passcode" may be required depending on
  the MFA method. To check which MFA is in use, run:

      $ vault read auth/<mount>/mfa_config

  Authenticate as "sally":

      $ vault login -method=ldap username=sally
      Password (will be hidden):

  Authenticate as "bob":

      $ vault login -method=ldap username=bob password=password

Configuration:

  method=<string>
      MFA method.

  passcode=<string>
      MFA OTP/passcode.

  password=<string>
      LDAP password to use for authentication. If not provided, the CLI will
      prompt for this on stdin.

  username=<string>
      LDAP username to use for authentication.
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
