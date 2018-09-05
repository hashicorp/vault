package vault

import (
	"context"
	"sort"

	"github.com/hashicorp/vault/logical"
)

// Capabilities is used to fetch the capabilities of the given token on the given path
func (c *Core) Capabilities(ctx context.Context, token, path string) ([]string, error) {
	if path == "" {
		return nil, &logical.StatusBadRequest{Err: "missing path"}
	}

	if token == "" {
		return nil, &logical.StatusBadRequest{Err: "missing token"}
	}

	te, err := c.tokenStore.Lookup(ctx, token)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return nil, &logical.StatusBadRequest{Err: "invalid token"}
	}

	// Start with token entry policies
	policies := te.Policies

	// Fetch entity and entity group policies
	entity, derivedPolicies, err := c.fetchEntityAndDerivedPolicies(te.EntityID)
	if err != nil {
		return nil, err
	}
	if entity != nil && entity.Disabled {
		c.logger.Warn("permission denied as the entity on the token is disabled")
		return nil, logical.ErrPermissionDenied
	}
	if te.EntityID != "" && entity == nil {
		c.logger.Warn("permission denied as the entity on the token is invalid")
		return nil, logical.ErrPermissionDenied
	}
	policies = append(policies, derivedPolicies...)

	if len(policies) == 0 {
		return []string{DenyCapability}, nil
	}

	acl, err := c.policyStore.ACL(ctx, entity, policies...)
	if err != nil {
		return nil, err
	}

	capabilities := acl.Capabilities(path)
	sort.Strings(capabilities)
	return capabilities, nil
}
