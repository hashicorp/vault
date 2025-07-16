// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	tokenPolicyType = "token"
)

func pathToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("role"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixConsul,
			OperationVerb:   "generate",
			OperationSuffix: "credentials",
		},

		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathTokenRead,
		},
	}
}

func (b *backend) pathTokenRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role := d.Get("role").(string)
	entry, err := req.Storage.Get(ctx, "policy/"+role)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %w", err)
	}
	if entry == nil {
		return logical.ErrorResponse(fmt.Sprintf("role %q not found", role)), nil
	}

	var roleConfigData roleConfig
	if err := entry.DecodeJSON(&roleConfigData); err != nil {
		return nil, err
	}

	if roleConfigData.TokenType == "" {
		roleConfigData.TokenType = "client"
	}

	// Get the consul client
	c, userErr, intErr := b.client(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}

	// Generate a name for the token
	tokenName := fmt.Sprintf("Vault %s %s %d", role, req.DisplayName, time.Now().UnixNano())

	writeOpts := &api.WriteOptions{}
	writeOpts = writeOpts.WithContext(ctx)

	// Create an ACLEntry for Consul pre 1.4
	if (roleConfigData.Policy != "" && roleConfigData.TokenType == "client") ||
		(roleConfigData.Policy == "" && roleConfigData.TokenType == "management") {
		token, _, err := c.ACL().Create(&api.ACLEntry{
			Name:  tokenName,
			Type:  roleConfigData.TokenType,
			Rules: roleConfigData.Policy,
		}, writeOpts)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		// Use the helper to create the secret
		s := b.Secret(SecretTokenType).Response(map[string]interface{}{
			"token": token,
		}, map[string]interface{}{
			"token": token,
			"role":  role,
		})
		s.Secret.TTL = roleConfigData.TTL
		s.Secret.MaxTTL = roleConfigData.MaxTTL
		return s, nil
	}

	// Create an ACLToken for Consul 1.4 and above
	policyLinks := []*api.ACLTokenPolicyLink{}
	for _, policyName := range roleConfigData.Policies {
		policyLinks = append(policyLinks, &api.ACLTokenPolicyLink{
			Name: policyName,
		})
	}

	roleLinks := []*api.ACLTokenRoleLink{}
	for _, roleName := range roleConfigData.ConsulRoles {
		roleLinks = append(roleLinks, &api.ACLTokenRoleLink{
			Name: roleName,
		})
	}

	aclServiceIdentities := parseServiceIdentities(roleConfigData.ServiceIdentities)
	aclNodeIdentities := parseNodeIdentities(roleConfigData.NodeIdentities)

	token, _, err := c.ACL().TokenCreate(&api.ACLToken{
		Description:       tokenName,
		Policies:          policyLinks,
		Roles:             roleLinks,
		ServiceIdentities: aclServiceIdentities,
		NodeIdentities:    aclNodeIdentities,
		Local:             roleConfigData.Local,
		Namespace:         roleConfigData.ConsulNamespace,
		Partition:         roleConfigData.Partition,
	}, writeOpts)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Use the helper to create the secret
	s := b.Secret(SecretTokenType).Response(map[string]interface{}{
		"token":            token.SecretID,
		"accessor":         token.AccessorID,
		"local":            token.Local,
		"consul_namespace": token.Namespace,
		"partition":        token.Partition,
	}, map[string]interface{}{
		"token":   token.AccessorID,
		"role":    role,
		"version": tokenPolicyType,
	})
	s.Secret.TTL = roleConfigData.TTL
	s.Secret.MaxTTL = roleConfigData.MaxTTL

	return s, nil
}

func parseServiceIdentities(data []string) []*api.ACLServiceIdentity {
	aclServiceIdentities := []*api.ACLServiceIdentity{}

	for _, serviceIdentity := range data {
		entry := &api.ACLServiceIdentity{}
		components := strings.Split(serviceIdentity, ":")
		entry.ServiceName = components[0]
		if len(components) == 2 {
			entry.Datacenters = strings.Split(components[1], ",")
		}
		aclServiceIdentities = append(aclServiceIdentities, entry)
	}

	return aclServiceIdentities
}

func parseNodeIdentities(data []string) []*api.ACLNodeIdentity {
	aclNodeIdentities := []*api.ACLNodeIdentity{}

	for _, nodeIdentity := range data {
		entry := &api.ACLNodeIdentity{}
		components := strings.Split(nodeIdentity, ":")
		entry.NodeName = components[0]
		if len(components) > 1 {
			entry.Datacenter = components[1]
		}
		aclNodeIdentities = append(aclNodeIdentities, entry)
	}

	return aclNodeIdentities
}
