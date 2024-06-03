// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package userpass

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

type UserpassAuth struct {
	mountPath    string
	username     string
	password     string
	passwordFile string
	passwordEnv  string
}

type Password struct {
	// Path on the file system where the password corresponding to this
	// application's Vault role can be found.
	FromFile string
	// The name of the environment variable containing the password
	// that corresponds to this application's Vault role.
	FromEnv string
	// The password as a plaintext string value.
	FromString string
}

var _ api.AuthMethod = (*UserpassAuth)(nil)

type LoginOption func(a *UserpassAuth) error

const (
	defaultMountPath = "userpass"
)

// NewUserpassAuth initializes a new Userpass auth method interface to be
// passed as a parameter to the client.Auth().Login method.
//
// Supported options: WithMountPath
func NewUserpassAuth(username string, password *Password, opts ...LoginOption) (*UserpassAuth, error) {
	if username == "" {
		return nil, fmt.Errorf("no user name provided for login")
	}

	if password == nil {
		return nil, fmt.Errorf("no password provided for login")
	}

	err := password.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	a := &UserpassAuth{
		mountPath: defaultMountPath,
		username:  username,
	}

	// password will be read in at login time if it comes from a file or environment variable, in case the underlying value changes
	if password.FromFile != "" {
		a.passwordFile = password.FromFile
	}

	if password.FromEnv != "" {
		a.passwordEnv = password.FromEnv
	}

	if password.FromString != "" {
		a.password = password.FromString
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *UserpassAuth as the argument
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	// return the modified auth struct instance
	return a, nil
}

func (a *UserpassAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	loginData := make(map[string]interface{})

	if a.passwordFile != "" {
		passwordValue, err := a.readPasswordFromFile()
		if err != nil {
			return nil, fmt.Errorf("error reading password: %w", err)
		}
		loginData["password"] = passwordValue
	} else if a.passwordEnv != "" {
		passwordValue := os.Getenv(a.passwordEnv)
		if passwordValue == "" {
			return nil, fmt.Errorf("password was specified with an environment variable with an empty value")
		}
		loginData["password"] = passwordValue
	} else {
		loginData["password"] = a.password
	}

	path := fmt.Sprintf("auth/%s/login/%s", a.mountPath, a.username)
	resp, err := client.Logical().WriteWithContext(ctx, path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with userpass auth: %w", err)
	}

	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *UserpassAuth) error {
		a.mountPath = mountPath
		return nil
	}
}

func (a *UserpassAuth) readPasswordFromFile() (string, error) {
	passwordFile, err := os.Open(a.passwordFile)
	if err != nil {
		return "", fmt.Errorf("unable to open file containing password: %w", err)
	}
	defer passwordFile.Close()

	limitedReader := io.LimitReader(passwordFile, 1000)
	passwordBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("unable to read password: %w", err)
	}

	passwordValue := strings.TrimSuffix(string(passwordBytes), "\n")

	return passwordValue, nil
}

func (password *Password) validate() error {
	if password.FromFile == "" && password.FromEnv == "" && password.FromString == "" {
		return fmt.Errorf("password for Userpass auth must be provided with a source file, environment variable, or plaintext string")
	}

	if password.FromFile != "" {
		if password.FromEnv != "" || password.FromString != "" {
			return fmt.Errorf("only one source for the password should be specified")
		}
	}

	if password.FromEnv != "" {
		if password.FromFile != "" || password.FromString != "" {
			return fmt.Errorf("only one source for the password should be specified")
		}
	}

	if password.FromString != "" {
		if password.FromFile != "" || password.FromEnv != "" {
			return fmt.Errorf("only one source for the password should be specified")
		}
	}
	return nil
}
