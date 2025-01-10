// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var _ logical.ExtendedSystemView = (*extendedSystemViewImpl)(nil)

type extendedSystemViewImpl struct {
	dynamicSystemView
}

func (e extendedSystemViewImpl) Auditor() logical.Auditor {
	return genericAuditor{
		mountType: e.mountEntry.Type,
		namespace: e.mountEntry.Namespace(),
		c:         e.core,
	}
}

func (e extendedSystemViewImpl) ForwardGenericRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	// Forward the request if allowed
	if couldForward(e.core) {
		ctx = namespace.ContextWithNamespace(ctx, e.mountEntry.Namespace())
		ctx = logical.IndexStateContext(ctx, &logical.WALState{})
		ctx = context.WithValue(ctx, ctxKeyForwardedRequestMountAccessor{}, e.mountEntry.Accessor)
		return forward(ctx, e.core, req)
	}

	return nil, logical.ErrReadOnly
}

// SudoPrivilege returns true if given path has sudo privileges
// for the given client token
func (e extendedSystemViewImpl) SudoPrivilege(ctx context.Context, path string, token string) bool {
	// Resolve the token policy
	te, err := e.core.tokenStore.Lookup(ctx, token)
	if err != nil {
		e.core.logger.Error("failed to lookup sudo token", "error", err)
		return false
	}

	// Ensure the token is valid
	if te == nil {
		e.core.logger.Error("entry not found for given token")
		return false
	}

	policyNames := make(map[string][]string)
	// Add token policies
	policyNames[te.NamespaceID] = append(policyNames[te.NamespaceID], te.Policies...)

	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, e.core)
	if err != nil {
		e.core.logger.Error("failed to lookup token namespace", "error", err)
		return false
	}
	if tokenNS == nil {
		e.core.logger.Error("failed to lookup token namespace", "error", namespace.ErrNoNamespace)
		return false
	}

	// Add identity policies from all the namespaces
	entity, identityPolicies, err := e.core.fetchEntityAndDerivedPolicies(ctx, tokenNS, te.EntityID, te.NoIdentityPolicies)
	if err != nil {
		e.core.logger.Error("failed to fetch identity policies", "error", err)
		return false
	}
	for nsID, nsPolicies := range identityPolicies {
		policyNames[nsID] = append(policyNames[nsID], nsPolicies...)
	}

	tokenCtx := namespace.ContextWithNamespace(ctx, tokenNS)

	// Add the inline policy if it's set
	policies := make([]*Policy, 0)
	if te.InlinePolicy != "" {
		inlinePolicy, err := ParseACLPolicy(tokenNS, te.InlinePolicy)
		if err != nil {
			e.core.logger.Error("failed to parse the token's inline policy", "error", err)
			return false
		}
		policies = append(policies, inlinePolicy)
	}

	// Construct the corresponding ACL object. Derive and use a new context that
	// uses the req.ClientToken's namespace
	acl, err := e.core.policyStore.ACL(tokenCtx, entity, policyNames, policies...)
	if err != nil {
		e.core.logger.Error("failed to retrieve ACL for token's policies", "token_policies", te.Policies, "error", err)
		return false
	}

	// The operation type isn't important here as this is run from a path the
	// user has already been given access to; we only care about whether they
	// have sudo. Note that we use root context because the path that comes in
	// must be fully-qualified already so we don't want AllowOperation to
	// prepend a namespace prefix onto it.
	req := new(logical.Request)
	req.Operation = logical.ReadOperation
	req.Path = path
	authResults := acl.AllowOperation(namespace.RootContext(ctx), req, true)
	return authResults.RootPrivs
}

func (e extendedSystemViewImpl) APILockShouldBlockRequest() (bool, error) {
	mountEntry := e.mountEntry
	if mountEntry == nil {
		return false, fmt.Errorf("no mount entry")
	}
	ns := mountEntry.Namespace()

	if err := e.core.entBlockRequestIfError(ns.Path, mountEntry.Path); err != nil {
		return true, nil
	}

	return false, nil
}

func (e extendedSystemViewImpl) RequestWellKnownRedirect(ctx context.Context, src, dest string) error {
	return e.core.WellKnownRedirects.TryRegister(ctx, e.core, e.mountEntry.UUID, src, dest)
}

func (e extendedSystemViewImpl) DeregisterWellKnownRedirect(ctx context.Context, src string) bool {
	return e.core.WellKnownRedirects.DeregisterSource(e.mountEntry.UUID, src)
}

// GetPinnedPluginVersion implements logical.ExtendedSystemView.
func (e extendedSystemViewImpl) GetPinnedPluginVersion(ctx context.Context, pluginType consts.PluginType, pluginName string) (*pluginutil.PinnedVersion, error) {
	return e.core.pluginCatalog.GetPinnedVersion(ctx, pluginType, pluginName)
}
