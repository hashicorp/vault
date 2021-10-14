package kubernetes

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type KubernetesAuth struct {
	roleName                string
	mountPath               string
	serviceAccountTokenPath string
}

type LoginOption func(a *KubernetesAuth) error

const (
	defaultMountPath               = "kubernetes"
	defaultServiceAccountTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

// NewKubernetesAuth creates a KubernetesAuth struct which can be passed to
// the client.Auth().Login method to authenticate to Vault. The roleName
// parameter should be the name of the role in Vault that was created with
// this app's Kubernetes service account bound to it.
//
// Supported options: WithMountPath, WithServiceAccountTokenPath
func NewKubernetesAuth(roleName string, opts ...LoginOption) (*KubernetesAuth, error) {
	var _ api.AuthMethod = (*KubernetesAuth)(nil)

	if roleName == "" {
		return nil, fmt.Errorf("no role name was provided")
	}

	a := &KubernetesAuth{
		roleName:                roleName,
		mountPath:               defaultMountPath,
		serviceAccountTokenPath: defaultServiceAccountTokenPath,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *KubernetesAuth as the argument
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	// return the modified auth struct instance
	return a, nil
}

func (a *KubernetesAuth) Login(client *api.Client) (*api.Secret, error) {
	jwt, err := os.ReadFile(a.serviceAccountTokenPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read file containing service account token: %w", err)
	}

	loginData := map[string]interface{}{
		"jwt":  string(jwt),
		"role": a.roleName,
	}

	path := fmt.Sprintf("auth/%s/login", a.mountPath)
	resp, err := client.Logical().Write(path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *KubernetesAuth) error {
		a.mountPath = mountPath
		return nil
	}
}

// WithServiceAccountTokenPath allows you to specify a different path to
// where your application's Kubernetes service account token is mounted,
// instead of the default of /var/run/secrets/kubernetes.io/serviceaccount/token
func WithServiceAccountTokenPath(pathToToken string) LoginOption {
	return func(a *KubernetesAuth) error {
		a.serviceAccountTokenPath = pathToToken
		return nil
	}
}
