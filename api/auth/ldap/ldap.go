package ldap

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

type LDAPAuth struct {
	mountPath string
	username  string
	password  string
}

type Password struct {
	// The password as a plaintext string value.
	FromString string
}

var _ api.AuthMethod = (*LDAPAuth)(nil)

type LoginOption func(a *LDAPAuth) error

const (
	defaultMountPath = "ldap"
)

// NewLDAPAuth initializes a new Userpass auth method interface to be
// passed as a parameter to the client.Auth().Login method.
//
// Supported options: WithMountPath
func NewLDAPAuth(username string, password *Password, opts ...LoginOption) (*LDAPAuth, error) {
	if username == "" {
		return nil, fmt.Errorf("no user name provided for login")
	}

	if password == nil {
		return nil, fmt.Errorf("no password provided for login")
	}

	a := &LDAPAuth{
		mountPath: defaultMountPath,
		username:  username,
	}

	if password.FromString != "" {
		a.password = password.FromString
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *LDAPAuth as the argument
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	// return the modified auth struct instance
	return a, nil
}

func (a *LDAPAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	loginData := make(map[string]interface{})

	loginData["password"] = a.password
	path := fmt.Sprintf("auth/%s/login/%s", a.mountPath, a.username)
	resp, err := client.Logical().Write(path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with LDAP auth: %w", err)
	}

	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *LDAPAuth) error {
		a.mountPath = mountPath
		return nil
	}
}
