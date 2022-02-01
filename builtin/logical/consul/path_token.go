package consul

import (
	"context"
	"fmt"
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
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},

			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: `List of policies to attach to the token.`,
			},

			"consul_namespace": {
				Type:        framework.TypeString,
				Description: "Namespace to create the token in.",
			},

			"partition": {
				Type:        framework.TypeString,
				Description: "Admin partition to create the token in.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathTokenRead,
		},
	}
}

func (b *backend) pathTokenRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role := d.Get("role").(string)
	policies := d.Get("policies").([]string)
	namespace := d.Get("consul_namespace").(string)
	partition := d.Get("partition").(string)

	entry, err := req.Storage.Get(ctx, "policy/"+role)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %w", err)
	}
	if entry == nil {
		return logical.ErrorResponse(fmt.Sprintf("role %q not found", role)), nil
	}

	var result roleConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if result.TokenType == "" {
		result.TokenType = "client"
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
	if (result.Policy != "" && result.TokenType == "client") ||
		(result.Policy == "" && result.TokenType == "management") {
		token, _, err := c.ACL().Create(&api.ACLEntry{
			Name:  tokenName,
			Type:  result.TokenType,
			Rules: result.Policy,
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
		s.Secret.TTL = result.TTL
		s.Secret.MaxTTL = result.MaxTTL
		return s, nil
	}

	// Create an ACLToken for Consul 1.4 and above
	// If policies were supplied here, then overwrite the policies
	// that were given when the role was written
	var policyLink []*api.ACLTokenPolicyLink
	if len(policies) > 0 {
		policyLink = getPolicies(policies)
	} else {
		policyLink = getPolicies(result.Policies)
	}

	// If a namespace was supplied here, then overwrite the namespace
	// that was given when the role was written
	if namespace == "" {
		namespace = result.Namespace
	}
	// If a partition was supplied here, then overwrite the partition
	// that was given when the role was written
	if partition == "" {
		partition = result.Partition
	}
	token, _, err := c.ACL().TokenCreate(&api.ACLToken{
		Description: tokenName,
		Policies:    policyLink,
		Local:       result.Local,
		Namespace:   namespace,
		Partition:   partition,
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
	s.Secret.TTL = result.TTL
	s.Secret.MaxTTL = result.MaxTTL

	return s, nil
}

func getPolicies(policies []string) []*api.ACLLink {
	policyLink := []*api.ACLTokenPolicyLink{}

	for _, policyName := range policies {
		policyLink = append(policyLink, &api.ACLTokenPolicyLink{
			Name: policyName,
		})
	}

	return policyLink
}
