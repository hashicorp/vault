// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubernetes

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type KubernetesAuth struct {
	roleName            string
	mountPath           string
	serviceAccountToken string
}

var _ api.AuthMethod = (*KubernetesAuth)(nil)

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
// The Kubernetes service account token JWT is retrieved from
// /var/run/secrets/kubernetes.io/serviceaccount/token by default. To change this
// path, pass the WithServiceAccountTokenPath option. To instead pass the
// JWT directly as a string, or to read the value from an environment
// variable, use WithServiceAccountToken and WithServiceAccountTokenEnv respectively.
//
// Supported options: WithMountPath, WithServiceAccountTokenPath, WithServiceAccountTokenEnv, WithServiceAccountToken
func NewKubernetesAuth(roleName string, opts ...LoginOption) (*KubernetesAuth, error) {
	if roleName == "" {
		return nil, fmt.Errorf("no role name was provided")
	}

	a := &KubernetesAuth{
		roleName:  roleName,
		mountPath: defaultMountPath,
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

	if a.serviceAccountToken == "" {
		token, err := readTokenFromFile(defaultServiceAccountTokenPath)
		if err != nil {
			return nil, fmt.Errorf("error reading service account token from default location: %w", err)
		}
		a.serviceAccountToken = token
	}

	// return the modified auth struct instance
	return a, nil
}

func (a *KubernetesAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	loginData := map[string]interface{}{
		"jwt":  a.serviceAccountToken,
		"role": a.roleName,
	}

	path := fmt.Sprintf("auth/%s/login", a.mountPath)
	resp, err := client.Logical().WriteWithContext(ctx, path, loginData)
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
		token, err := readTokenFromFile(pathToToken)
		if err != nil {
			return fmt.Errorf("unable to read service account token from file: %w", err)
		}
		a.serviceAccountToken = token
		return nil
	}
}

func WithServiceAccountToken(jwt string) LoginOption {
	return func(a *KubernetesAuth) error {
		a.serviceAccountToken = jwt
		return nil
	}
}

func WithServiceAccountTokenEnv(envVar string) LoginOption {
	return func(a *KubernetesAuth) error {
		token := os.Getenv(envVar)
		if token == "" {
			return fmt.Errorf("service account token was specified with an environment variable with an empty value")
		}
		a.serviceAccountToken = token
		return nil
	}
}

func readTokenFromFile(filepath string) (string, error) {
	jwt, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("unable to read file containing service account token: %w", err)
	}
	return string(jwt), nil
}
