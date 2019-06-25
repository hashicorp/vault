// Copyright © 2019, Oracle and/or its affiliates. All rights reserved.
package ociauth

import (
	"context"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/oracle/oci-go-sdk/common/auth"
	"sync"
	"fmt"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend()
	if err != nil {
		return nil, err
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend

	// Lock to make changes to role entries
	roleMutex sync.RWMutex

	// Lock to make changes to config entries
	configMutex sync.RWMutex

	authenticationClient *AuthenticationClient
}

func Backend() (*backend, error) {
	b := &backend{}

	//Create the instance principal provider
	ip, err := auth.InstancePrincipalConfigurationProvider()
	if err != nil {
		return nil, fmt.Errorf("Unable to create InstancePrincipalConfigurationProvider")
	}

	//Create the authentication client
	authenticationClient, err := NewAuthenticationClientWithConfigurationProvider(ip)
	if err != nil {
		return nil, fmt.Errorf("Unable to create authenticationClient")
	}
	b.authenticationClient = &authenticationClient

	b.Backend = &framework.Backend{
		Help: backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login/*",
			},
		},
		Paths: []*framework.Path{
			pathLogin(b),
			pathRole(b),
			pathListRoles(b),
			pathConfig(b),
			pathListConfigs(b),
		},
		BackendType: logical.TypeCredential,
	}

	return b, nil
}

const backendHelp = `
OCI Auth Plugin
`
