package approle

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/vault/api"
)

type AppRoleAuth struct {
	mountPath      string
	roleID         string
	pathToSecretID string
	unwrap         bool
}

type LoginOption func(a *AppRoleAuth)

// NewAppRoleAuth initializes a new AppRole auth method interface to be passed as a parameter to the client.Auth().Login method.
//
// For a secret ID, the recommended secure pattern is to unwrap a one-time-use response-wrapping token that was placed here by a trusted orchestrator (https://learn.hashicorp.com/tutorials/vault/approle-best-practices?in=vault/auth-methods#secretid-delivery-best-practices)
// To indicate that the filepath points to this wrapping token and not just a plaintext secret ID, initialize NewAppRoleAuth with the WithWrappingToken LoginOption.
//
// Supported options: WithMountPath, WithWrappingToken
func NewAppRoleAuth(roleID, pathToSecretID string, opts ...LoginOption) (api.AuthMethod, error) {
	if roleID == "" {
		return nil, fmt.Errorf("no role ID provided for login")
	}

	if pathToSecretID == "" {
		return nil, fmt.Errorf("no path to secret ID provided for login")
	}

	const (
		defaultMountPath = "approle"
	)

	a := &AppRoleAuth{
		mountPath:      defaultMountPath,
		roleID:         roleID,
		pathToSecretID: pathToSecretID,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *AppRoleAuth as the argument
		opt(a)
	}

	// return the modified auth struct instance
	return a, nil
}

func (a *AppRoleAuth) Login(client *api.Client) (*api.Secret, error) {
	loginData := map[string]interface{}{
		"role_id": a.roleID,
	}

	secretIDBytes, err := ioutil.ReadFile(a.pathToSecretID)
	if err != nil {
		return nil, fmt.Errorf("unable to read file containing secret ID: %w", err)
	}
	secretID := string(secretIDBytes)

	// if it was indicated that the value in the file was actually a wrapping token, unwrap it first
	if a.unwrap {
		unwrappedToken, err := client.Logical().Unwrap(strings.TrimSuffix(secretID, "\n"))
		if err != nil {
			return nil, fmt.Errorf("unable to unwrap token: %w. If the AppRoleAuth struct was initialized with the WithWrappingToken LoginOption, then the filepath used should be a path to a response-wrapping token", err)
		}
		loginData["secret_id"] = unwrappedToken.Data["secret_id"]
	} else {
		loginData["secret_id"] = secretID
	}

	path := fmt.Sprintf("auth/%s/login", a.mountPath)
	resp, err := client.Logical().Write(path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with app role auth: %w", err)
	}

	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *AppRoleAuth) {
		a.mountPath = mountPath
	}
}

func WithWrappingToken() LoginOption {
	return func(a *AppRoleAuth) {
		a.unwrap = true
	}
}
