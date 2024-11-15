// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/oracle/oci-go-sdk/v59/common/auth"
)

// operationPrefixOCI is used as a prefix for OpenAPI operation id's.
const operationPrefixOCI = "oci"

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

	// Lock to make changes to authClient entries
	authClientMutex sync.RWMutex

	// The client used to authenticate with OCI Identity
	authenticationClient *AuthenticationClient
}

func Backend() (*backend, error) {
	b := &backend{}

	b.Backend = &framework.Backend{
		Help: backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login/*",
			},
		},
		Paths: []*framework.Path{
			pathLogin(b),
			pathLoginRole(b),
			pathRole(b),
			pathListRoles(b),
			pathConfig(b),
		},
		BackendType: logical.TypeCredential,
	}

	return b, nil
}

// createAuthClient creates an authentication client if one was not already created and stores in the backend.
func (b *backend) createAuthClient() error {

	b.authClientMutex.Lock()
	defer b.authClientMutex.Unlock()

	if b.authenticationClient != nil {
		return nil
	}

	// Create the instance principal provider
	ip, err := auth.InstancePrincipalConfigurationProvider()
	if err != nil {
		b.Logger().Debug("Unable to create InstancePrincipalConfigurationProvider", "err", err)
		return fmt.Errorf("unable to create InstancePrincipalConfigurationProvider")
	}

	// Create the authentication client
	authenticationClient, err := NewAuthenticationClientWithConfigurationProvider(ip)
	if err != nil {
		b.Logger().Debug("Unable to create authenticationClient", "err", err)
		return fmt.Errorf("unable to create authenticationClient")
	}

	b.authenticationClient = &authenticationClient

	return nil
}

const backendHelp = `
The OCI Auth plugin enables authentication and authorization using OCI Identity credentials. 

The OCI Auth plugin authorizes using roles. A role is defined as a set of allowed policies for specific entities. 
When an entity such as a user or instance logs in, it requests a role. 
The OCI Auth plugin checks whether the entity is allowed to use the role and which policies are associated with that role. 
It then assigns the given policies to the request.

The goal of roles is to restrict access to only the subset of secrets that are required, 
even if the entity has access to many more secrets. This conforms to the least-privilege security model.
`
