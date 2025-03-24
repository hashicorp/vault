// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cert

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/api"
)

const defaultMountPath = "cert"

type CertAuth struct {
	role      string
	mountPath string
}

var _ api.AuthMethod = (*CertAuth)(nil)

type LoginOption func(a *CertAuth) error

// NewCertAuth initializes a new Cert auth method interface to be
// passed as a parameter to the client.Auth().Login method. The client and other
// TLS configuration should be set up on the passed in client, you can use
// the NewCertAuthClient function to set up the client with the proper TLS client attributes.
//
// Supported options: WithRole, WithMountPath
func NewCertAuth(opts ...LoginOption) (*CertAuth, error) {
	a := &CertAuth{}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *CertAuth as the argument
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	if a.mountPath == "" {
		a.mountPath = defaultMountPath
	}

	// return the modified auth struct instance
	return a, nil
}

// Login sets up the required request body for the Cert auth method's /login
// endpoint, and performs a write to it. We assume the passed in client has the
// proper TLS client certificates set up.
func (a *CertAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if client == nil {
		return nil, fmt.Errorf("client is required for login with the associated client certs initialized")
	}

	loginData := make(map[string]interface{})
	if a.role != "" {
		loginData["name"] = a.role
	}

	certAuthPath := fmt.Sprintf("auth/%s/login", a.mountPath)
	resp, err := client.Logical().WriteWithContext(ctx, certAuthPath, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with cert auth: %w", err)
	}

	return resp, nil
}

// WithRole specifies the role to use for the login request.
func WithRole(roleName string) LoginOption {
	return func(a *CertAuth) error {
		a.role = roleName
		return nil
	}
}

// WithMountPath specifies the mount path to use for the login request.
func WithMountPath(mountPath string) LoginOption {
	return func(a *CertAuth) error {
		a.mountPath = mountPath
		return nil
	}
}

// NewDefaultCertAuthClient initializes a new client with a default configuration
// with the provided address and TLS configuration. The TLSConfig must have the ClientCert and ClientKey
// fields set.
func NewDefaultCertAuthClient(address string, tlsConfig *api.TLSConfig) (*api.Client, error) {
	if tlsConfig == nil {
		return nil, errors.New("tls config is required for cert auth client")
	}

	if tlsConfig.ClientCert == "" || tlsConfig.ClientKey == "" {
		return nil, errors.New("client cert and key are required for cert auth client")
	}

	if len(address) == 0 {
		return nil, errors.New("address is required for cert auth client")
	}

	cfg := api.DefaultConfig()
	cfg.Address = address
	err := cfg.ConfigureTLS(tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed configuring TLS on client config: %w", err)
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return client, nil
}

// NewCertAuthClient initializes a new client based on the passed in client
// with the provided TLS configuration. The TLSConfig must have the ClientCert and ClientKey
// fields set.
func NewCertAuthClient(c *api.Client, config *api.TLSConfig) (*api.Client, error) {
	if c == nil {
		return nil, errors.New("base client is required for cert auth client")
	}
	if config == nil {
		return nil, errors.New("tls config is required for cert auth client")
	}

	if config.ClientCert == "" || config.ClientKey == "" {
		return nil, errors.New("client cert and key are required for cert auth client")
	}

	conf := c.CloneConfig()
	err := conf.ConfigureTLS(config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure TLS on client: %w", err)
	}

	client, err := api.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return client, nil
}
