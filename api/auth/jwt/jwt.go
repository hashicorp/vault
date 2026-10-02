// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jwt

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

type JWTAuth struct {
	mountPath string
	jwt       string
	role      string
}

var _ api.AuthMethod = (*JWTAuth)(nil)

type LoginOption func(a *JWTAuth) error

const (
	defaultMountPath = "jwt"
)

// NewJWTAuth initializes a new JWT auth method interface to be
// passed as a parameter to the client.Auth().Login method.
//
// Supported options: WithMountPath
func NewJWTAuth(jwt, role string, opts ...LoginOption) (*JWTAuth, error) {
	if jwt == "" {
		return nil, fmt.Errorf("no jwt provided for login")
	}

	if role == "" {
		return nil, fmt.Errorf("no role provided for login")
	}

	a := &JWTAuth{
		mountPath: defaultMountPath,
		jwt:       jwt,
		role:      role,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *JWTAuth as the argument
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	// return the modified auth struct instance
	return a, nil
}

func (a *JWTAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	loginData := map[string]any{
		"jwt":  a.jwt,
		"role": a.role,
	}

	path := fmt.Sprintf("auth/%s/login", a.mountPath)
	resp, err := client.Logical().WriteWithContext(ctx, path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with jwt auth: %w", err)
	}
	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *JWTAuth) error {
		a.mountPath = mountPath
		return nil
	}
}
