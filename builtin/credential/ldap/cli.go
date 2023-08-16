// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"fmt"
	"os"
	"strings"

	pwd "github.com/hashicorp/go-secure-stdlib/password"
	"github.com/hashicorp/vault/api"
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
		password = passwordFromEnv()
		if password == "" {
			fmt.Fprintf(os.Stderr, "Password (will be hidden): ")
			var err error
			password, err = pwd.Read(os.Stdin)
			fmt.Fprintf(os.Stderr, "\n")
			if err != nil {
				return nil, err
			}
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

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=ldap [CONFIG K=V...]

  The LDAP auth method allows users to authenticate using LDAP or
  Active Directory.

  Authenticate as "sally":

      $ vault login -method=ldap username=sally
      Password (will be hidden):

  Authenticate as "bob":

      $ vault login -method=ldap username=bob password=password

Configuration:

  password=<string>
      LDAP password to use for authentication. If not provided, it will use
			the VAULT_LDAP_PASSWORD environment variable. If this is not set, the
			CLI will prompt for this on stdin.

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

func passwordFromEnv() string {
	return os.Getenv("VAULT_LDAP_PASSWORD")
}
