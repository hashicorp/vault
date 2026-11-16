// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"sort"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/internalshared/namespace"
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

	return c.capabilitiesAndSubscribeEventTypesForRequest(ctx, &logical.Request{
		ClientToken: token,
		Path:        path,
		Operation:   logical.ListOperation,
	})
}

func (c *Core) capabilitiesAndSubscribeEventTypesForRequest(ctx context.Context, req *logical.Request) ([]string, []string, error) {
	if req == nil || req.Path == "" {
		return nil, nil, &logical.StatusBadRequest{Err: "missing path"}
	}

	acl, err := c.capabilitiesACL(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	capabilities, eventTypes := c.capabilitiesFromACL(ctx, acl, req)
	return capabilities, eventTypes, nil
}

// capabilitiesForPaths returns the capabilities of the request's token on each
// of the given paths. The token is looked up and its ACL is built once, for the
// first path, and every path is then checked against that same ACL.
func (c *Core) capabilitiesForPaths(ctx context.Context, req *logical.Request, paths []string) (map[string][]string, error) {
	ret := make(map[string][]string, len(paths))

	var acl *ACL
	var aclReq *logical.Request
	for _, path := range paths {
		if path == "" {
			return nil, &logical.StatusBadRequest{Err: "missing path"}
		}

		if aclReq == nil {
			r := *req
			r.Path = path
			aclReq = &r

			var err error
			acl, err = c.capabilitiesACL(ctx, aclReq)
			if err != nil {
				return nil, err
			}
		}

		pathReq := *aclReq
		pathReq.Path = path
		ret[path], _ = c.capabilitiesFromACL(ctx, acl, &pathReq)
	}

	return ret, nil
}

// capabilitiesACL looks up the request's token and builds its ACL. It returns a
// nil ACL and a nil error when the token has no policies, which denies every
// path.
func (c *Core) capabilitiesACL(ctx context.Context, req *logical.Request) (*ACL, error) {
	if req.ClientToken == "" {
		return nil, &logical.StatusBadRequest{Err: "missing token"}
	}

	var te *logical.TokenEntry
	var entity *identity.Entity
	var acl *ACL
	var err error
	if IsOAuthJwt(req.ClientToken) || req.JwtIssuer != "" {
		acl, te, entity, _, err = c.fetchACLTokenEntryAndEntity(ctx, req)
		if err != nil {
			return nil, err
		}
	} else {
		te, err = c.tokenStore.Lookup(ctx, req.ClientToken)
		if err != nil {
			return nil, err
		}
	}
	if te == nil {
		return nil, &logical.StatusBadRequest{Err: "invalid token"}
	}

	var tokenNS *namespace.Namespace
	tokenNS, err = NamespaceByID(ctx, te.NamespaceID, c)
	if err != nil {
		return nil, err
	}
	if tokenNS == nil {
		return nil, namespace.ErrNoNamespace
	}

	var policyCount int
	policyNames := make(map[string][]string)
	policyNames[tokenNS.ID] = te.Policies
	policyCount += len(te.Policies)

	var identityPolicies map[string][]string
	if acl == nil {
		entity, identityPolicies, err = c.fetchEntityAndDerivedPolicies(ctx, tokenNS, te.EntityID, te.NoIdentityPolicies)
		if err != nil {
			return nil, err
		}
	}
	if entity != nil && entity.Disabled {
		c.logger.Warn("permission denied as the entity on the token is disabled")
		return nil, logical.ErrPermissionDenied
	}
	if te.EntityID != "" && entity == nil {
		c.logger.Warn("permission denied as the entity on the token is invalid")
		return nil, logical.ErrPermissionDenied
	}

	for nsID, nsPolicies := range identityPolicies {
		policyNames[nsID] = append(policyNames[nsID], nsPolicies...)
		policyCount += len(nsPolicies)
	}

	// Add capabilities of the inline policy if it's set
	policies := make([]*Policy, 0)
	if te.InlinePolicy != "" {
		inlinePolicy, err := ParseACLPolicy(tokenNS, te.InlinePolicy, WithDenySlashInTemplatedPaths(c.denySlashInTemplatedPolicyPaths))
		if err != nil {
			return nil, err
		}
		policies = append(policies, inlinePolicy)
		policyCount++
	}

	if policyCount == 0 && acl == nil {
		return nil, nil
	}

	// Construct the corresponding ACL object. ACL construction should be
	// performed on the token's namespace.
	if acl == nil {
		tokenCtx := namespace.ContextWithNamespace(ctx, tokenNS)
		acl, err = c.policyStore.ACL(tokenCtx, entity, policyNames, policies...)
		if err != nil {
			return nil, err
		}
	}

	return acl, nil
}

// capabilitiesFromACL returns the capabilities and subscribable event types
// that the given ACL grants on req.Path. A nil ACL denies everything.
func (c *Core) capabilitiesFromACL(ctx context.Context, acl *ACL, req *logical.Request) ([]string, []string) {
	if acl == nil {
		return []string{DenyCapability}, nil
	}

	capabilities, eventTypes := acl.CapabilitiesAndSubscribeEventTypes(ctx, req.Path)
	capabilities, eventTypes = c.filterOAuthJWTCapabilities(ctx, req, capabilities, eventTypes)
	sort.Strings(capabilities)
	return capabilities, eventTypes
}
