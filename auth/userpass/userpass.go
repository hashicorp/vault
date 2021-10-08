package userpass

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

type UserpassAuth struct {
	mountPath string
	username  string
	password  string
}

type LoginOption func(a *UserpassAuth)

// Supported options: WithMountPath
func NewUserpassAuth(username, password string, opts ...LoginOption) (api.AuthMethod, error) {
	if username == "" {
		return nil, fmt.Errorf("no user name provided for login")
	}

	if password == "" {
		return nil, fmt.Errorf("no password provided for login")
	}

	const (
		defaultMountPath = "userpass"
	)

	a := &UserpassAuth{
		mountPath: defaultMountPath,
		username:  username,
		password:  password,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *UserpassAuth as the argument
		opt(a)
	}

	// return the modified auth struct instance
	return a, nil
}

func (a *UserpassAuth) Login(client *api.Client) (*api.Secret, error) {
	loginData := map[string]interface{}{
		"password": a.password,
	}

	path := fmt.Sprintf("auth/%s/login/%s", a.mountPath, a.username)
	resp, err := client.Logical().Write(path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with userpass auth: %w", err)
	}

	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *UserpassAuth) {
		a.mountPath = mountPath
	}
}
