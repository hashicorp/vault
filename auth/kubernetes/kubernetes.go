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

type LoginOption func(a *KubernetesAuth)

// NewKubernetesAuth creates a KubernetesAuth struct which can be passed to the client.Auth().Login method to authenticate to Vault.
// The roleName parameter should be the name of the role in Vault that was created with this app's Kubernetes service account bound to it.
//
// Supported options: WithMountPath, WithServiceAccountTokenPath
func NewKubernetesAuth(roleName string, opts ...LoginOption) (api.AuthMethod, error) {
	if roleName == "" {
		return nil, fmt.Errorf("no role name was provided")
	}

	const (
		defaultMountPath               = "kubernetes"
		defaultServiceAccountTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	)

	a := &KubernetesAuth{
		roleName:                roleName,
		mountPath:               defaultMountPath,
		serviceAccountTokenPath: defaultServiceAccountTokenPath,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *KubernetesAuth as the argument
		opt(a)
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
	return func(a *KubernetesAuth) {
		a.mountPath = mountPath
	}
}

// WithServiceAccountTokenPath allows you to specify a different path to where your application's
// Kubernetes service account token is mounted, instead of the default of /var/run/secrets/kubernetes.io/serviceaccount/token
func WithServiceAccountTokenPath(pathToToken string) LoginOption {
	return func(a *KubernetesAuth) {
		a.serviceAccountTokenPath = pathToToken
	}
}
