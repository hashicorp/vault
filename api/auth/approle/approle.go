package approle

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

type AppRoleAuth struct {
	mountPath    string
	roleID       string
	secretID     string
	secretIDFile string
	secretIDEnv  string
	unwrap       bool
}

var _ api.AuthMethod = (*AppRoleAuth)(nil)

// SecretID is a struct that allows you to specify where your application is
// storing the secret ID required for login to the AppRole auth method.
type SecretID struct {
	// Path on the file system where a trusted orchestrator has placed the
	// application's secret ID. The recommended secure pattern is to use
	// response-wrapping tokens rather than a plaintext value, by passing
	// WithWrappingToken() to NewAppRoleAuth.
	// https://learn.hashicorp.com/tutorials/vault/approle-best-practices?in=vault/auth-methods#secretid-delivery-best-practices
	FromFile string
	// The name of the environment variable containing the application's
	// secret ID.
	FromEnv string
	// The secret ID as a plaintext string value.
	FromString string
}

type LoginOption func(a *AppRoleAuth) error

const (
	defaultMountPath = "approle"
)

// NewAppRoleAuth initializes a new AppRole auth method interface to be
// passed as a parameter to the client.Auth().Login method.
//
// For a secret ID, the recommended secure pattern is to unwrap a one-time-use
// response-wrapping token that was placed here by a trusted orchestrator
// (https://learn.hashicorp.com/tutorials/vault/approle-best-practices?in=vault/auth-methods#secretid-delivery-best-practices)
// To indicate that the filepath points to this wrapping token and not just
// a plaintext secret ID, initialize NewAppRoleAuth with the
// WithWrappingToken LoginOption.
//
// Supported options: WithMountPath, WithWrappingToken
func NewAppRoleAuth(roleID string, secretID *SecretID, opts ...LoginOption) (*AppRoleAuth, error) {
	if roleID == "" {
		return nil, fmt.Errorf("no role ID provided for login")
	}

	if secretID == nil {
		return nil, fmt.Errorf("no secret ID provided for login")
	}

	err := secretID.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid secret ID: %w", err)
	}

	a := &AppRoleAuth{
		mountPath: defaultMountPath,
		roleID:    roleID,
	}

	// secret ID will be read in at login time if it comes from a file or environment variable, in case the underlying value changes
	if secretID.FromFile != "" {
		a.secretIDFile = secretID.FromFile
	}

	if secretID.FromEnv != "" {
		a.secretIDEnv = secretID.FromEnv
	}

	if secretID.FromString != "" {
		a.secretID = secretID.FromString
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *AppRoleAuth as the argument
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	// return the modified auth struct instance
	return a, nil
}

func (a *AppRoleAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	loginData := map[string]interface{}{
		"role_id": a.roleID,
	}

	if a.secretIDFile != "" {
		secretIDValue, err := a.readSecretIDFromFile()
		if err != nil {
			return nil, fmt.Errorf("error reading secret ID: %w", err)
		}

		// if it was indicated that the value in the file was actually a wrapping
		// token, unwrap it first
		if a.unwrap {
			unwrappedToken, err := client.Logical().Unwrap(secretIDValue)
			if err != nil {
				return nil, fmt.Errorf("unable to unwrap token: %w. If the AppRoleAuth struct was initialized with the WithWrappingToken LoginOption, then the secret ID's filepath should be a path to a response-wrapping token", err)
			}
			loginData["secret_id"] = unwrappedToken.Data["secret_id"]
		} else {
			loginData["secret_id"] = secretIDValue
		}
	} else if a.secretIDEnv != "" {
		secretIDValue := os.Getenv(a.secretIDEnv)
		if secretIDValue == "" {
			return nil, fmt.Errorf("secret ID was specified with an environment variable with an empty value")
		}
		loginData["secret_id"] = secretIDValue
	} else {
		loginData["secret_id"] = a.secretID
	}

	path := fmt.Sprintf("auth/%s/login", a.mountPath)
	resp, err := client.Logical().Write(path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with app role auth: %w", err)
	}

	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *AppRoleAuth) error {
		a.mountPath = mountPath
		return nil
	}
}

func WithWrappingToken() LoginOption {
	return func(a *AppRoleAuth) error {
		a.unwrap = true
		return nil
	}
}

func (a *AppRoleAuth) readSecretIDFromFile() (string, error) {
	secretIDFile, err := os.Open(a.secretIDFile)
	if err != nil {
		return "", fmt.Errorf("unable to open file containing secret ID: %w", err)
	}
	defer secretIDFile.Close()

	limitedReader := io.LimitReader(secretIDFile, 1000)
	secretIDBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("unable to read secret ID: %w", err)
	}

	secretIDValue := strings.TrimSuffix(string(secretIDBytes), "\n")

	return secretIDValue, nil
}

func (secretID *SecretID) validate() error {
	if secretID.FromFile == "" && secretID.FromEnv == "" && secretID.FromString == "" {
		return fmt.Errorf("secret ID for AppRole must be provided with a source file, environment variable, or plaintext string")
	}

	if secretID.FromFile != "" {
		if secretID.FromEnv != "" || secretID.FromString != "" {
			return fmt.Errorf("only one source for the secret ID should be specified")
		}
	}

	if secretID.FromEnv != "" {
		if secretID.FromFile != "" || secretID.FromString != "" {
			return fmt.Errorf("only one source for the secret ID should be specified")
		}
	}

	if secretID.FromString != "" {
		if secretID.FromFile != "" || secretID.FromEnv != "" {
			return fmt.Errorf("only one source for the secret ID should be specified")
		}
	}
	return nil
}
