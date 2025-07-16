// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixConsul,
			OperationSuffix: "roles",
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixConsul,
			OperationSuffix: "role",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},

			// The "policy" and "token_type" parameters were deprecated in Consul back in version 1.4.
			// They have been removed from Consul as of version 1.11. Consider removing them here in the future.
			"policy": {
				Type: framework.TypeString,
				Description: `Policy document, base64 encoded. Required
for 'client' tokens. Required for Consul pre-1.4.`,
				Deprecated: true,
			},

			"token_type": {
				Type:    framework.TypeString,
				Default: "client",
				Description: `Which type of token to create: 'client' or 'management'. If
a 'management' token, the "policy", "policies", and "consul_roles" parameters are not
required. Defaults to 'client'.`,
				Deprecated: true,
			},

			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: `Use "consul_policies" instead.`,
				Deprecated:  true,
			},

			"consul_policies": {
				Type: framework.TypeCommaStringSlice,
				Description: `List of policies to attach to the token. Either "consul_policies"
or "consul_roles" are required for Consul 1.5 and above, or just "consul_policies" if
using Consul 1.4.`,
			},

			"consul_roles": {
				Type: framework.TypeCommaStringSlice,
				Description: `List of Consul roles to attach to the token. Either "policies"
or "consul_roles" are required for Consul 1.5 and above.`,
			},

			"local": {
				Type: framework.TypeBool,
				Description: `Indicates that the token should not be replicated globally 
and instead be local to the current datacenter. Available in Consul 1.4 and above.`,
			},

			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "TTL for the Consul token created from the role.",
			},

			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Max TTL for the Consul token created from the role.",
			},

			"lease": {
				Type:        framework.TypeDurationSecond,
				Description: `Use "ttl" instead.`,
				Deprecated:  true,
			},

			"consul_namespace": {
				Type: framework.TypeString,
				Description: `Indicates which namespace that the token will be
created within. Defaults to 'default'. Available in Consul 1.7 and above.`,
			},

			"partition": {
				Type: framework.TypeString,
				Description: `Indicates which admin partition that the token
will be created within. Defaults to 'default'. Available in Consul 1.11 and above.`,
			},

			"service_identities": {
				Type: framework.TypeStringSlice,
				Description: `List of Service Identities to attach to the
token, separated by semicolons. Available in Consul 1.5 or above.`,
			},

			"node_identities": {
				Type: framework.TypeStringSlice,
				Description: `List of Node Identities to attach to the
token. Available in Consul 1.8.1 or above.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRolesRead,
			logical.UpdateOperation: b.pathRolesWrite,
			logical.DeleteOperation: b.pathRolesDelete,
		},
	}
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "policy/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, "policy/"+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var roleConfigData roleConfig
	if err := entry.DecodeJSON(&roleConfigData); err != nil {
		return nil, err
	}

	if roleConfigData.TokenType == "" {
		roleConfigData.TokenType = "client"
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"lease":            int64(roleConfigData.TTL.Seconds()),
			"ttl":              int64(roleConfigData.TTL.Seconds()),
			"max_ttl":          int64(roleConfigData.MaxTTL.Seconds()),
			"token_type":       roleConfigData.TokenType,
			"local":            roleConfigData.Local,
			"consul_namespace": roleConfigData.ConsulNamespace,
			"partition":        roleConfigData.Partition,
		},
	}
	if roleConfigData.Policy != "" {
		resp.Data["policy"] = base64.StdEncoding.EncodeToString([]byte(roleConfigData.Policy))
	}
	if len(roleConfigData.Policies) > 0 {
		resp.Data["consul_policies"] = roleConfigData.Policies
	}
	if len(roleConfigData.ConsulRoles) > 0 {
		resp.Data["consul_roles"] = roleConfigData.ConsulRoles
	}
	if len(roleConfigData.ServiceIdentities) > 0 {
		resp.Data["service_identities"] = roleConfigData.ServiceIdentities
	}
	if len(roleConfigData.NodeIdentities) > 0 {
		resp.Data["node_identities"] = roleConfigData.NodeIdentities
	}

	return resp, nil
}

func (b *backend) pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	tokenType := d.Get("token_type").(string)
	policy := d.Get("policy").(string)
	consulPolicies := d.Get("consul_policies").([]string)
	policies := d.Get("policies").([]string)
	roles := d.Get("consul_roles").([]string)
	serviceIdentities := d.Get("service_identities").([]string)
	nodeIdentities := d.Get("node_identities").([]string)

	switch tokenType {
	case "client":
		if policy == "" && len(policies) == 0 && len(consulPolicies) == 0 &&
			len(roles) == 0 && len(serviceIdentities) == 0 && len(nodeIdentities) == 0 {
			return logical.ErrorResponse(
				"Use either a policy document, a list of policies or roles, or a set of service or node identities, depending on your Consul version"), nil
		}
	case "management":
	default:
		return logical.ErrorResponse("token_type must be \"client\" or \"management\""), nil
	}

	if len(consulPolicies) == 0 {
		consulPolicies = policies
	}

	policyRaw, err := base64.StdEncoding.DecodeString(policy)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error decoding policy base64: %s", err)), nil
	}

	var ttl time.Duration
	ttlRaw, ok := d.GetOk("ttl")
	if ok {
		ttl = time.Second * time.Duration(ttlRaw.(int))
	} else {
		leaseParamRaw, ok := d.GetOk("lease")
		if ok {
			ttl = time.Second * time.Duration(leaseParamRaw.(int))
		}
	}

	var maxTTL time.Duration
	maxTTLRaw, ok := d.GetOk("max_ttl")
	if ok {
		maxTTL = time.Second * time.Duration(maxTTLRaw.(int))
	}

	name := d.Get("name").(string)
	local := d.Get("local").(bool)
	namespace := d.Get("consul_namespace").(string)
	partition := d.Get("partition").(string)
	entry, err := logical.StorageEntryJSON("policy/"+name, roleConfig{
		Policy:            string(policyRaw),
		Policies:          consulPolicies,
		ConsulRoles:       roles,
		ServiceIdentities: serviceIdentities,
		NodeIdentities:    nodeIdentities,
		TokenType:         tokenType,
		TTL:               ttl,
		MaxTTL:            maxTTL,
		Local:             local,
		ConsulNamespace:   namespace,
		Partition:         partition,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if err := req.Storage.Delete(ctx, "policy/"+name); err != nil {
		return nil, err
	}
	return nil, nil
}

type roleConfig struct {
	Policy            string        `json:"policy"`
	Policies          []string      `json:"policies"`
	ConsulRoles       []string      `json:"consul_roles"`
	ServiceIdentities []string      `json:"service_identities"`
	NodeIdentities    []string      `json:"node_identities"`
	TTL               time.Duration `json:"lease"`
	MaxTTL            time.Duration `json:"max_ttl"`
	TokenType         string        `json:"token_type"`
	Local             bool          `json:"local"`
	ConsulNamespace   string        `json:"consul_namespace"`
	Partition         string        `json:"partition"`
}
