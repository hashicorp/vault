package userpass

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

type UserpassAuth struct {
	mountPath string
	username  string
	password  string
}

type Password struct {
	// Path on the file system where the password corresponding to this
	// application's Vault role can be found. Other auth backends provide
	// more secure options for authentication.
	FromFile string
	// The name of the environment variable containing the password
	// that corresponds to this application's Vault role. Can be insecure
	// if the environment variable is logged.
	FromEnv string
	// The password as a plaintext string value. Insecure.
	FromString string
}

type LoginOption func(a *UserpassAuth) error

const (
	defaultMountPath = "userpass"
)

// NewUserpassAuth initializes a new Userpass auth method interface to be
// passed as a parameter to the client.Auth().Login method.
//
// Supported options: WithMountPath
func NewUserpassAuth(username string, password *Password, opts ...LoginOption) (*UserpassAuth, error) {
	var _ api.AuthMethod = (*UserpassAuth)(nil)

	if username == "" {
		return nil, fmt.Errorf("no user name provided for login")
	}

	if password == nil {
		return nil, fmt.Errorf("no password provided for login")
	}

	if password.FromFile == "" && password.FromEnv == "" && password.FromString == "" {
		return nil, fmt.Errorf("password must be provided with a source file, environment variable, or plaintext string")
	}

	passwordValue, err := readPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error reading password: %w", err)
	}

	a := &UserpassAuth{
		mountPath: defaultMountPath,
		username:  username,
		password:  passwordValue,
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
	return func(a *UserpassAuth) error {
		a.mountPath = mountPath
		return nil
	}
}

func readPassword(password *Password) (string, error) {
	var parsedPassword string
	if password.FromFile != "" {
		if password.FromEnv != "" || password.FromString != "" {
			return "", fmt.Errorf("only one location for the password should be specified")
		}
		passwordFile, err := os.Open(password.FromFile)
		if err != nil {
			return "", fmt.Errorf("unable to open file containing password: %w", err)
		}
		defer passwordFile.Close()

		limitedReader := io.LimitReader(passwordFile, 1000)
		passwordBytes, err := io.ReadAll(limitedReader)
		if err != nil {
			return "", fmt.Errorf("unable to read password: %w", err)
		}
		parsedPassword = string(passwordBytes)
	}

	if password.FromEnv != "" {
		if password.FromFile != "" || password.FromString != "" {
			return "", fmt.Errorf("only one location for the password should be specified")
		}
		parsedPassword = os.Getenv(password.FromEnv)
		if parsedPassword == "" {
			return "", fmt.Errorf("password was specified with an environment variable with an empty value")
		}
	}

	if password.FromString != "" {
		if password.FromFile != "" || password.FromEnv != "" {
			return "", fmt.Errorf("only one location for the password should be specified")
		}
		parsedPassword = password.FromString
	}

	passwordValue := strings.TrimSuffix(parsedPassword, "\n")

	return passwordValue, nil
}
