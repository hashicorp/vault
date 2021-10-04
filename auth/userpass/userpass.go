package userpass

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

type UserpassAuth struct {
	MountPath string
	Username  string
	Password  string
}

type LoginOption func(a *UserpassAuth)

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
		MountPath: defaultMountPath,
		Username:  username,
		Password:  password,
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
		"password": a.Password,
	}

	path := fmt.Sprintf("auth/%s/login/%s", a.MountPath, a.Username)
	resp, err := client.Logical().Write(path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with userpass auth: %w", err)
	}

	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *UserpassAuth) {
		a.MountPath = mountPath
	}
}
