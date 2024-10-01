// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"sort"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// Capabilities is used to fetch the capabilities of the given token on the
// given path.
func (c *Core) Capabilities(ctx context.Context, token, path string) ([]string, error) {
	capabilities, _, err := c.CapabilitiesAndSubscribeEventTypes(ctx, token, path)
	return capabilities, err
}

// CapabilitiesAndSubscribeEventTypes is used to fetch the capabilities and event types that are allowed to
// be subscribed to by given token on the given path.
func (c *Core) CapabilitiesAndSubscribeEventTypes(ctx context.Context, token, path string) ([]string, []string, error) {
	if path == "" {
		return nil, nil, &logical.StatusBadRequest{Err: "missing path"}
	}

	if token == "" {
		return nil, nil, &logical.StatusBadRequest{Err: "missing token"}
	}

	te, err := c.tokenStore.Lookup(ctx, token)
	if err != nil {
		return nil, nil, err
	}
	if te == nil {
		return nil, nil, &logical.StatusBadRequest{Err: "invalid token"}
	}

	var tokenNS *namespace.Namespace
	tokenNS, err = NamespaceByID(ctx, te.NamespaceID, c)
	if err != nil {
		return nil, nil, err
	}
	if tokenNS == nil {
		return nil, nil, namespace.ErrNoNamespace
	}

	var policyCount int
	policyNames := make(map[string][]string)
	policyNames[tokenNS.ID] = te.Policies
	policyCount += len(te.Policies)

	entity, identityPolicies, err := c.fetchEntityAndDerivedPolicies(ctx, tokenNS, te.EntityID, te.NoIdentityPolicies)
	if err != nil {
		return nil, nil, err
	}
	if entity != nil && entity.Disabled {
		c.logger.Warn("permission denied as the entity on the token is disabled")
		return nil, nil, logical.ErrPermissionDenied
	}
	if te.EntityID != "" && entity == nil {
		c.logger.Warn("permission denied as the entity on the token is invalid")
		return nil, nil, logical.ErrPermissionDenied
	}

	for nsID, nsPolicies := range identityPolicies {
		policyNames[nsID] = append(policyNames[nsID], nsPolicies...)
		policyCount += len(nsPolicies)
	}

	// Add capabilities of the inline policy if it's set
	policies := make([]*Policy, 0)
	if te.InlinePolicy != "" {
		inlinePolicy, err := ParseACLPolicy(tokenNS, te.InlinePolicy)
		if err != nil {
			return nil, nil, err
		}
		policies = append(policies, inlinePolicy)
		policyCount++
	}

	if policyCount == 0 {
		return []string{DenyCapability}, nil, nil
	}

	// Construct the corresponding ACL object. ACL construction should be
	// performed on the token's namespace.
	tokenCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	acl, err := c.policyStore.ACL(tokenCtx, entity, policyNames, policies...)
	if err != nil {
		return nil, nil, err
	}

	capabilities, eventTypes := acl.CapabilitiesAndSubscribeEventTypes(ctx, path)
	sort.Strings(capabilities)
	return capabilities, eventTypes, nil
}
