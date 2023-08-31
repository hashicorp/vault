// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/go-uuid"
	uberAtomic "go.uber.org/atomic"

	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/identity/mfa"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/quotas"
	"github.com/hashicorp/vault/vault/tokens"
)

const (
	replTimeout                           = 1 * time.Second
	EnvVaultDisableLocalAuthMountEntities = "VAULT_DISABLE_LOCAL_AUTH_MOUNT_ENTITIES"
	// base path to store locked users
	coreLockedUsersPath = "core/login/lockedUsers/"
)

var (
	// DefaultMaxRequestDuration is the amount of time we'll wait for a request
	// to complete, unless overridden on a per-handler basis
	DefaultMaxRequestDuration = 90 * time.Second

	egpDebugLogging bool

	// if this returns an error, the request should be blocked and the error
	// should be returned to the client
	enterpriseBlockRequestIfError = blockRequestIfErrorImpl
)

// HandlerProperties is used to seed configuration into a vaulthttp.Handler.
// It's in this package to avoid a circular dependency
type HandlerProperties struct {
	Core                  *Core
	ListenerConfig        *configutil.Listener
	DisablePrintableCheck bool
	RecoveryMode          bool
	RecoveryToken         *uberAtomic.String
}

// fetchEntityAndDerivedPolicies returns the entity object for the given entity
// ID. If the entity is merged into a different entity object, the entity into
// which the given entity ID is merged into will be returned. This function
// also returns the cumulative list of policies that the entity is entitled to
// if skipDeriveEntityPolicies is set to false. This list includes the policies from the
// entity itself and from all the groups in which the given entity ID is a member of.
func (c *Core) fetchEntityAndDerivedPolicies(ctx context.Context, tokenNS *namespace.Namespace, entityID string, skipDeriveEntityPolicies bool) (*identity.Entity, map[string][]string, error) {
	if entityID == "" || c.identityStore == nil {
		return nil, nil, nil
	}

	// c.logger.Debug("entity set on the token", "entity_id", te.EntityID)

	// Fetch the entity
	entity, err := c.identityStore.MemDBEntityByID(entityID, false)
	if err != nil {
		c.logger.Error("failed to lookup entity using its ID", "error", err)
		return nil, nil, err
	}

	if entity == nil {
		// If there was no corresponding entity object found, it is
		// possible that the entity got merged into another entity. Try
		// finding entity based on the merged entity index.
		entity, err = c.identityStore.MemDBEntityByMergedEntityID(entityID, false)
		if err != nil {
			c.logger.Error("failed to lookup entity in merged entity ID index", "error", err)
			return nil, nil, err
		}
	}

	policies := make(map[string][]string)
	if entity != nil && !skipDeriveEntityPolicies {
		// c.logger.Debug("entity successfully fetched; adding entity policies to token's policies to create ACL")

		// Attach the policies on the entity
		if len(entity.Policies) != 0 {
			policies[entity.NamespaceID] = append(policies[entity.NamespaceID], entity.Policies...)
		}

		groupPolicies, err := c.identityStore.groupPoliciesByEntityID(entity.ID)
		if err != nil {
			c.logger.Error("failed to fetch group policies", "error", err)
			return nil, nil, err
		}

		policyApplicationMode, err := c.GetGroupPolicyApplicationMode(ctx)
		if err != nil {
			return nil, nil, err
		}

		// Filter and add the policies to the resultant set
		for nsID, nsPolicies := range groupPolicies {
			ns, err := NamespaceByID(ctx, nsID, c)
			if err != nil {
				return nil, nil, err
			}
			if ns == nil {
				return nil, nil, namespace.ErrNoNamespace
			}
			// If we're only applying policies to namespaces within the same
			// hierarchy, then skip any policies not found in the same
			// hierarchy
			if policyApplicationMode == groupPolicyApplicationModeWithinNamespaceHierarchy {
				if tokenNS.Path != ns.Path && !ns.HasParent(tokenNS) {
					continue
				}
			}
			nsPolicies = strutil.RemoveDuplicates(nsPolicies, false)
			if len(nsPolicies) != 0 {
				policies[nsID] = append(policies[nsID], nsPolicies...)
			}
		}
	}

	return entity, policies, err
}

func (c *Core) fetchACLTokenEntryAndEntity(ctx context.Context, req *logical.Request) (*ACL, *logical.TokenEntry, *identity.Entity, map[string][]string, error) {
	defer metrics.MeasureSince([]string{"core", "fetch_acl_and_token"}, time.Now())

	// Ensure there is a client token
	if req.ClientToken == "" {
		return nil, nil, nil, nil, logical.ErrPermissionDenied
	}

	if c.tokenStore == nil {
		c.logger.Error("token store is unavailable")
		return nil, nil, nil, nil, ErrInternalError
	}

	// Resolve the token policy
	var te *logical.TokenEntry
	switch req.TokenEntry() {
	case nil:
		var err error
		te, err = c.tokenStore.Lookup(ctx, req.ClientToken)
		if err != nil {
			c.logger.Error("failed to lookup acl token", "error", err)
			return nil, nil, nil, nil, ErrInternalError
		}
		// Set the token entry here since it has not been cached yet
		req.SetTokenEntry(te)
	default:
		te = req.TokenEntry()
	}

	// Ensure the token is valid
	if te == nil {
		return nil, nil, nil, nil, logical.ErrPermissionDenied
	}

	// CIDR checks bind all tokens except non-expiring root tokens
	if te.TTL != 0 && len(te.BoundCIDRs) > 0 {
		var valid bool
		remoteSockAddr, err := sockaddr.NewSockAddr(req.Connection.RemoteAddr)
		if err != nil {
			if c.Logger().IsDebug() {
				c.Logger().Debug("could not parse remote addr into sockaddr", "error", err, "remote_addr", req.Connection.RemoteAddr)
			}
			return nil, nil, nil, nil, logical.ErrPermissionDenied
		}
		for _, cidr := range te.BoundCIDRs {
			if cidr.Contains(remoteSockAddr) {
				valid = true
				break
			}
		}
		if !valid {
			return nil, nil, nil, nil, logical.ErrPermissionDenied
		}
	}

	policyNames := make(map[string][]string)
	// Add tokens policies
	policyNames[te.NamespaceID] = append(policyNames[te.NamespaceID], te.Policies...)

	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, c)
	if err != nil {
		c.logger.Error("failed to fetch token namespace", "error", err)
		return nil, nil, nil, nil, ErrInternalError
	}
	if tokenNS == nil {
		c.logger.Error("failed to fetch token namespace", "error", namespace.ErrNoNamespace)
		return nil, nil, nil, nil, ErrInternalError
	}

	// Add identity policies from all the namespaces
	entity, identityPolicies, err := c.fetchEntityAndDerivedPolicies(ctx, tokenNS, te.EntityID, te.NoIdentityPolicies)
	if err != nil {
		return nil, nil, nil, nil, ErrInternalError
	}
	for nsID, nsPolicies := range identityPolicies {
		policyNames[nsID] = policyutil.SanitizePolicies(append(policyNames[nsID], nsPolicies...), false)
	}

	// Attach token's namespace information to the context. Wrapping tokens by
	// should be able to be used anywhere, so we also special case behavior.
	var tokenCtx context.Context
	if len(policyNames) == 1 &&
		len(policyNames[te.NamespaceID]) == 1 &&
		(policyNames[te.NamespaceID][0] == responseWrappingPolicyName ||
			policyNames[te.NamespaceID][0] == controlGroupPolicyName) &&
		(strings.HasSuffix(req.Path, "sys/wrapping/unwrap") ||
			strings.HasSuffix(req.Path, "sys/wrapping/lookup") ||
			strings.HasSuffix(req.Path, "sys/wrapping/rewrap")) {
		// Use the request namespace; will find the copy of the policy for the
		// local namespace
		tokenCtx = ctx
	} else {
		// Use the token's namespace for looking up policy
		tokenCtx = namespace.ContextWithNamespace(ctx, tokenNS)
	}

	// Add the inline policy if it's set
	policies := make([]*Policy, 0)
	if te.InlinePolicy != "" {
		inlinePolicy, err := ParseACLPolicy(tokenNS, te.InlinePolicy)
		if err != nil {
			return nil, nil, nil, nil, ErrInternalError
		}
		policies = append(policies, inlinePolicy)
	}

	// Construct the corresponding ACL object. ACL construction should be
	// performed on the token's namespace.
	acl, err := c.policyStore.ACL(tokenCtx, entity, policyNames, policies...)
	if err != nil {
		c.logger.Error("failed to construct ACL", "error", err)
		return nil, nil, nil, nil, ErrInternalError
	}

	return acl, te, entity, identityPolicies, nil
}

func (c *Core) CheckToken(ctx context.Context, req *logical.Request, unauth bool) (*logical.Auth, *logical.TokenEntry, error) {
	defer metrics.MeasureSince([]string{"core", "check_token"}, time.Now())

	var acl *ACL
	var te *logical.TokenEntry
	var entity *identity.Entity
	var identityPolicies map[string][]string

	// Even if unauth, if a token is provided, there's little reason not to
	// gather as much info as possible for the audit log and to e.g. control
	// trace mode for EGPs.
	if !unauth || (unauth && req.ClientToken != "") {
		var err error
		acl, te, entity, identityPolicies, err = c.fetchACLTokenEntryAndEntity(ctx, req)
		// In the unauth case we don't want to fail the command, since it's
		// unauth, we just have no information to attach to the request, so
		// ignore errors...this was best-effort anyways
		if err != nil && !unauth {
			return nil, te, err
		}
	}

	if entity != nil && entity.Disabled {
		c.logger.Warn("permission denied as the entity on the token is disabled")
		return nil, te, logical.ErrPermissionDenied
	}
	if te != nil && te.EntityID != "" && entity == nil {
		if c.perfStandby {
			return nil, nil, logical.ErrPerfStandbyPleaseForward
		}
		c.logger.Warn("permission denied as the entity on the token is invalid")
		return nil, te, logical.ErrPermissionDenied
	}

	// Check if this is a root protected path
	rootPath := c.router.RootPath(ctx, req.Path)

	if rootPath && unauth {
		return nil, nil, errors.New("cannot access root path in unauthenticated request")
	}

	// At this point we won't be forwarding a raw request; we should delete
	// authorization headers as appropriate
	switch req.ClientTokenSource {
	case logical.ClientTokenFromVaultHeader:
		delete(req.Headers, consts.AuthHeaderName)
	case logical.ClientTokenFromAuthzHeader:
		if headers, ok := req.Headers["Authorization"]; ok {
			retHeaders := make([]string, 0, len(headers))
			for _, v := range headers {
				if strings.HasPrefix(v, "Bearer ") {
					continue
				}
				retHeaders = append(retHeaders, v)
			}
			req.Headers["Authorization"] = retHeaders
		}
	}

	// When we receive a write of either type, rather than require clients to
	// PUT/POST and trust the operation, we ask the backend to give us the real
	// skinny -- if the backend implements an existence check, it can tell us
	// whether a particular resource exists. Then we can mark it as an update
	// or creation as appropriate.
	if req.Operation == logical.CreateOperation || req.Operation == logical.UpdateOperation {
		existsResp, checkExists, resourceExists, err := c.router.RouteExistenceCheck(ctx, req)
		switch err {
		case logical.ErrUnsupportedPath:
			// fail later via bad path to avoid confusing items in the log
			checkExists = false
		case logical.ErrRelativePath:
			return nil, te, errutil.UserError{Err: err.Error()}
		case nil:
			if existsResp != nil && existsResp.IsError() {
				return nil, te, existsResp.Error()
			}
			// Otherwise, continue on
		default:
			c.logger.Error("failed to run existence check", "error", err)
			if _, ok := err.(errutil.UserError); ok {
				return nil, te, err
			} else {
				return nil, te, ErrInternalError
			}
		}

		switch {
		case !checkExists:
			// No existence check, so always treat it as an update operation, which is how it is pre 0.5
			req.Operation = logical.UpdateOperation
		case resourceExists:
			// It exists, so force an update operation
			req.Operation = logical.UpdateOperation
		case !resourceExists:
			// It doesn't exist, force a create operation
			req.Operation = logical.CreateOperation
		default:
			panic("unreachable code")
		}
	}
	// Create the auth response
	auth := &logical.Auth{
		ClientToken: req.ClientToken,
		Accessor:    req.ClientTokenAccessor,
	}

	var clientID string
	var isTWE bool
	if te != nil {
		auth.IdentityPolicies = identityPolicies[te.NamespaceID]
		auth.TokenPolicies = te.Policies
		auth.Policies = policyutil.SanitizePolicies(append(te.Policies, identityPolicies[te.NamespaceID]...), false)
		auth.Metadata = te.Meta
		auth.DisplayName = te.DisplayName
		auth.EntityID = te.EntityID
		delete(identityPolicies, te.NamespaceID)
		auth.ExternalNamespacePolicies = identityPolicies
		// Store the entity ID in the request object
		req.EntityID = te.EntityID
		auth.TokenType = te.Type
		auth.TTL = te.TTL
		if te.CreationTime > 0 {
			auth.IssueTime = time.Unix(te.CreationTime, 0)
		}
		clientID, isTWE = te.CreateClientID()
		req.ClientID = clientID
	}

	// Check the standard non-root ACLs. Return the token entry if it's not
	// allowed so we can decrement the use count.
	authResults := c.performPolicyChecks(ctx, acl, te, req, entity, &PolicyCheckOpts{
		Unauth:            unauth,
		RootPrivsRequired: rootPath,
	})

	auth.PolicyResults = &logical.PolicyResults{
		Allowed: authResults.Allowed,
	}

	if !authResults.Allowed {
		retErr := authResults.Error

		// If we get a control group error and we are a performance standby,
		// restore the client token information to the request so that we can
		// forward this request properly to the active node.
		if retErr.ErrorOrNil() != nil && checkErrControlGroupTokenNeedsCreated(retErr) &&
			c.perfStandby && len(req.ClientToken) != 0 {
			switch req.ClientTokenSource {
			case logical.ClientTokenFromVaultHeader:
				req.Headers[consts.AuthHeaderName] = []string{req.ClientToken}
			case logical.ClientTokenFromAuthzHeader:
				req.Headers["Authorization"] = append(req.Headers["Authorization"], fmt.Sprintf("Bearer %s", req.ClientToken))
			}
			// We also return the appropriate error so that the caller can forward the
			// request to the active node
			return auth, te, logical.ErrPerfStandbyPleaseForward
		}

		if authResults.Error.ErrorOrNil() == nil || authResults.DeniedError {
			retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		}
		return auth, te, retErr
	}

	if authResults.ACLResults != nil && len(authResults.ACLResults.GrantingPolicies) > 0 {
		auth.PolicyResults.GrantingPolicies = authResults.ACLResults.GrantingPolicies
	}
	if authResults.SentinelResults != nil && len(authResults.SentinelResults.GrantingPolicies) > 0 {
		auth.PolicyResults.GrantingPolicies = append(auth.PolicyResults.GrantingPolicies, authResults.SentinelResults.GrantingPolicies...)
	}

	c.activityLogLock.RLock()
	activityLog := c.activityLog
	c.activityLogLock.RUnlock()
	// If it is an authenticated ( i.e. with vault token ) request, increment client count
	if !unauth && activityLog != nil {
		err := activityLog.HandleTokenUsage(ctx, te, clientID, isTWE)
		if err != nil {
			return auth, te, err
		}
	}
	return auth, te, nil
}

// HandleRequest is used to handle a new incoming request
func (c *Core) HandleRequest(httpCtx context.Context, req *logical.Request) (resp *logical.Response, err error) {
	return c.switchedLockHandleRequest(httpCtx, req, true)
}

func (c *Core) switchedLockHandleRequest(httpCtx context.Context, req *logical.Request, doLocking bool) (resp *logical.Response, err error) {
	if doLocking {
		c.stateLock.RLock()
		defer c.stateLock.RUnlock()
	}
	if c.Sealed() {
		return nil, consts.ErrSealed
	}
	if c.standby && !c.perfStandby {
		return nil, consts.ErrStandby
	}

	if c.activeContext == nil || c.activeContext.Err() != nil {
		return nil, errors.New("active context canceled after getting state lock")
	}

	ctx, cancel := context.WithCancel(c.activeContext)
	go func(ctx context.Context, httpCtx context.Context) {
		select {
		case <-ctx.Done():
		case <-httpCtx.Done():
			cancel()
		}
	}(ctx, httpCtx)

	ns, err := namespace.FromContext(httpCtx)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not parse namespace from http context: %w", err)
	}

	ctx = namespace.ContextWithNamespace(ctx, ns)
	inFlightReqID, ok := httpCtx.Value(logical.CtxKeyInFlightRequestID{}).(string)
	if ok {
		ctx = context.WithValue(ctx, logical.CtxKeyInFlightRequestID{}, inFlightReqID)
	}
	requestRole, ok := httpCtx.Value(logical.CtxKeyRequestRole{}).(string)
	if ok {
		ctx = context.WithValue(ctx, logical.CtxKeyRequestRole{}, requestRole)
	}
	resp, err = c.handleCancelableRequest(ctx, req)
	req.SetTokenEntry(nil)
	cancel()
	return resp, err
}

func (c *Core) handleCancelableRequest(ctx context.Context, req *logical.Request) (resp *logical.Response, err error) {
	// Allowing writing to a path ending in / makes it extremely difficult to
	// understand user intent for the filesystem-like backends (kv,
	// cubbyhole) -- did they want a key named foo/ or did they want to write
	// to a directory foo/ with no (or forgotten) key, or...? It also affects
	// lookup, because paths ending in / are considered prefixes by some
	// backends. Basically, it's all just terrible, so don't allow it.
	if strings.HasSuffix(req.Path, "/") &&
		(req.Operation == logical.UpdateOperation ||
			req.Operation == logical.CreateOperation ||
			req.Operation == logical.PatchOperation) {
		return logical.ErrorResponse("cannot write to a path ending in '/'"), nil
	}
	waitGroup, err := waitForReplicationState(ctx, c, req)
	if err != nil {
		return nil, err
	}

	// MountPoint will not always be set at this point, so we ensure the req contains it
	// as it is depended on by some functionality (e.g. quotas)
	req.MountPoint = c.router.MatchingMount(ctx, req.Path)

	// Decrement the wait group when our request is done
	if waitGroup != nil {
		defer waitGroup.Done()
	}

	if c.MissingRequiredState(req.RequiredState(), c.perfStandby) {
		return nil, logical.ErrMissingRequiredState
	}

	err = c.PopulateTokenEntry(ctx, req)
	if err != nil {
		return nil, err
	}

	// Always forward requests that are using a limited use count token.
	if c.perfStandby && req.ClientTokenRemainingUses > 0 {
		// Prevent forwarding on local-only requests.
		return nil, logical.ErrPerfStandbyPleaseForward
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not parse namespace from http context: %w", err)
	}
	var requestBodyToken string
	var returnRequestAuthToken bool

	// req.Path will be relative by this point. The prefix check is first
	// to fail faster if we're not in this situation since it's a hot path
	switch {
	case strings.HasPrefix(req.Path, "sys/wrapping/"), strings.HasPrefix(req.Path, "auth/token/"):
		// Get the token ns info; if we match the paths below we want to
		// swap in the token context (but keep the relative path)
		te := req.TokenEntry()
		newCtx := ctx
		if te != nil {
			ns, err := NamespaceByID(ctx, te.NamespaceID, c)
			if err != nil {
				c.Logger().Warn("error looking up namespace from the token's namespace ID", "error", err)
				return nil, err
			}
			if ns != nil {
				newCtx = namespace.ContextWithNamespace(ctx, ns)
			}
		}
		switch req.Path {
		// Route the token wrapping request to its respective sys NS
		case "sys/wrapping/lookup", "sys/wrapping/rewrap", "sys/wrapping/unwrap":
			ctx = newCtx
			// A lookup on a token that is about to expire returns nil, which means by the
			// time we can validate a wrapping token lookup will return nil since it will
			// be revoked after the call. So we have to do the validation here.
			valid, err := c.validateWrappingToken(ctx, req)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("error validating wrapping token: %s", err.Error())), logical.ErrPermissionDenied
			}
			if !valid {
				return nil, consts.ErrInvalidWrappingToken
			}

		// The -self paths have no meaning outside of the token NS, so
		// requests for these paths always go to the token NS
		case "auth/token/lookup-self", "auth/token/renew-self", "auth/token/revoke-self":
			ctx = newCtx
			returnRequestAuthToken = true

		// For the following operations, we can set the proper namespace context
		// using the token's embedded nsID if a relative path was provided.
		// The operation will still be gated by ACLs, which are checked later.
		case "auth/token/lookup", "auth/token/renew", "auth/token/revoke", "auth/token/revoke-orphan":
			token, ok := req.Data["token"]
			// If the token is not present (e.g. a bad request), break out and let the backend
			// handle the error
			if !ok {
				// If this is a token lookup request and if the token is not
				// explicitly provided, it will use the client token so we simply set
				// the context to the client token's context.
				if req.Path == "auth/token/lookup" {
					ctx = newCtx
				}
				break
			}
			if token == nil {
				return logical.ErrorResponse("invalid token"), logical.ErrPermissionDenied
			}
			// We don't care if the token is a server side consistent token or not. Either way, we're going
			// to be returning it for these paths instead of the short token stored in vault.
			requestBodyToken = token.(string)
			if IsSSCToken(token.(string)) {
				token, err = c.CheckSSCToken(ctx, token.(string), c.isLoginRequest(ctx, req), c.perfStandby)

				// If we receive an error from CheckSSCToken, we can assume the token is bad somehow, and the client
				// should receive a 403 bad token error like they do for all other invalid tokens, unless the error
				// specifies that we should forward the request or retry the request.
				if err != nil {
					if errors.Is(err, logical.ErrPerfStandbyPleaseForward) || errors.Is(err, logical.ErrMissingRequiredState) {
						return nil, err
					}
					return logical.ErrorResponse("bad token"), logical.ErrPermissionDenied
				}
				req.Data["token"] = token
			}
			_, nsID := namespace.SplitIDFromString(token.(string))
			if nsID != "" {
				ns, err := NamespaceByID(ctx, nsID, c)
				if err != nil {
					c.Logger().Warn("error looking up namespace from the token's namespace ID", "error", err)
					return nil, err
				}
				if ns != nil {
					ctx = namespace.ContextWithNamespace(ctx, ns)
				}
			}
		}

	// The following relative sys/leases/ paths handles re-routing requests
	// to the proper namespace using the lease ID on applicable paths.
	case strings.HasPrefix(req.Path, "sys/leases/"):
		switch req.Path {
		// For the following operations, we can set the proper namespace context
		// using the lease's embedded nsID if a relative path was provided.
		// The operation will still be gated by ACLs, which are checked later.
		case "sys/leases/lookup", "sys/leases/renew", "sys/leases/revoke", "sys/leases/revoke-force":
			leaseID, ok := req.Data["lease_id"]
			// If lease ID is not present, break out and let the backend handle the error
			if !ok || leaseID == nil {
				break
			}
			_, nsID := namespace.SplitIDFromString(leaseID.(string))
			if nsID != "" {
				ns, err := NamespaceByID(ctx, nsID, c)
				if err != nil {
					c.Logger().Warn("error looking up namespace from the lease's namespace ID", "error", err)
					return nil, err
				}
				if ns != nil {
					ctx = namespace.ContextWithNamespace(ctx, ns)
				}
			}
		}

	// Prevent any metrics requests to be forwarded from a standby node.
	// Instead, we return an error since we cannot be sure if we have an
	// active token store to validate the provided token.
	case strings.HasPrefix(req.Path, "sys/metrics"):
		if c.standby && !c.perfStandby {
			return nil, ErrCannotForwardLocalOnly
		}
	}

	ns, err = namespace.FromContext(ctx)
	if err != nil {
		return nil, errwrap.Wrapf("could not parse namespace from http context: {{err}}", err)
	}

	if !hasNamespaces(c) && ns.Path != "" {
		return nil, logical.CodedError(403, "namespaces feature not enabled")
	}

	walState := &logical.WALState{}
	ctx = logical.IndexStateContext(ctx, walState)
	var auth *logical.Auth
	if c.isLoginRequest(ctx, req) {
		resp, auth, err = c.handleLoginRequest(ctx, req)
	} else {
		resp, auth, err = c.handleRequest(ctx, req)
	}

	if err == nil && c.requestResponseCallback != nil {
		c.requestResponseCallback(c.router.MatchingBackend(ctx, req.Path), req, resp)
	}

	// If we saved the token in the request, we should return it in the response
	// data.
	if resp != nil && resp.Data != nil {
		if _, ok := resp.Data["error"]; !ok {
			if requestBodyToken != "" {
				resp.Data["id"] = requestBodyToken
			} else if returnRequestAuthToken && req.InboundSSCToken != "" {
				resp.Data["id"] = req.InboundSSCToken
			}
		}
	}
	if resp != nil && resp.Auth != nil && requestBodyToken != "" {
		// if a client token has already been set and the request body token's internal token
		// is equal to that value, then we can return the original request body token
		tok, _ := c.DecodeSSCToken(requestBodyToken)
		if resp.Auth.ClientToken == tok {
			resp.Auth.ClientToken = requestBodyToken
		}
	}

	// Ensure we don't leak internal data
	if resp != nil {
		if resp.Secret != nil {
			resp.Secret.InternalData = nil
		}
		if resp.Auth != nil {
			resp.Auth.InternalData = nil
		}
	}

	// We are wrapping if there is anything to wrap (not a nil response) and a
	// TTL was specified for the token. Errors on a call should be returned to
	// the caller, so wrapping is turned off if an error is hit and the error
	// is logged to the audit log.
	wrapping := resp != nil &&
		err == nil &&
		!resp.IsError() &&
		resp.WrapInfo != nil &&
		resp.WrapInfo.TTL != 0 &&
		resp.WrapInfo.Token == ""

	if wrapping {
		cubbyResp, cubbyErr := c.wrapInCubbyhole(ctx, req, resp, auth)
		// If not successful, returns either an error response from the
		// cubbyhole backend or an error; if either is set, set resp and err to
		// those and continue so that that's what we audit log. Otherwise
		// finish the wrapping and audit log that.
		if cubbyResp != nil || cubbyErr != nil {
			resp = cubbyResp
			err = cubbyErr
		} else {
			wrappingResp := &logical.Response{
				WrapInfo: resp.WrapInfo,
				Warnings: resp.Warnings,
			}
			resp = wrappingResp
		}
	}

	auditResp := resp
	// When unwrapping we want to log the actual response that will be written
	// out. We still want to return the raw value to avoid automatic updating
	// to any of it.
	if req.Path == "sys/wrapping/unwrap" &&
		resp != nil &&
		resp.Data != nil &&
		resp.Data[logical.HTTPRawBody] != nil {

		// Decode the JSON
		if resp.Data[logical.HTTPRawBodyAlreadyJSONDecoded] != nil {
			delete(resp.Data, logical.HTTPRawBodyAlreadyJSONDecoded)
		} else {
			httpResp := &logical.HTTPResponse{}
			err := jsonutil.DecodeJSON(resp.Data[logical.HTTPRawBody].([]byte), httpResp)
			if err != nil {
				c.logger.Error("failed to unmarshal wrapped HTTP response for audit logging", "error", err)
				return nil, ErrInternalError
			}

			auditResp = logical.HTTPResponseToLogicalResponse(httpResp)
		}
	}

	var nonHMACReqDataKeys []string
	var nonHMACRespDataKeys []string
	entry := c.router.MatchingMountEntry(ctx, req.Path)
	if entry != nil {
		// Get and set ignored HMAC'd value. Reset those back to empty afterwards.
		if rawVals, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
			nonHMACReqDataKeys = rawVals.([]string)
		}

		// Get and set ignored HMAC'd value. Reset those back to empty afterwards.
		if auditResp != nil {
			if rawVals, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_response_keys"); ok {
				nonHMACRespDataKeys = rawVals.([]string)
			}
		}
	}

	// Create an audit trail of the response
	if !isControlGroupRun(req) {
		switch req.Path {
		case "sys/replication/dr/status", "sys/replication/performance/status", "sys/replication/status":
		default:
			logInput := &logical.LogInput{
				Auth:                auth,
				Request:             req,
				Response:            auditResp,
				OuterErr:            err,
				NonHMACReqDataKeys:  nonHMACReqDataKeys,
				NonHMACRespDataKeys: nonHMACRespDataKeys,
			}
			if auditErr := c.auditBroker.LogResponse(ctx, logInput, c.auditedHeaders); auditErr != nil {
				c.logger.Error("failed to audit response", "request_path", req.Path, "error", auditErr)
				return nil, ErrInternalError
			}
		}
	}

	if walState.LocalIndex != 0 || walState.ReplicatedIndex != 0 {
		walState.ClusterID = c.ClusterID()
		if walState.LocalIndex == 0 {
			if c.perfStandby {
				walState.LocalIndex = LastRemoteWAL(c)
			} else {
				walState.LocalIndex = LastWAL(c)
			}
		}
		if walState.ReplicatedIndex == 0 {
			if c.perfStandby {
				walState.ReplicatedIndex = LastRemoteUpstreamWAL(c)
			} else {
				walState.ReplicatedIndex = LastRemoteWAL(c)
			}
		}

		req.SetResponseState(walState)
	}

	return
}

func isControlGroupRun(req *logical.Request) bool {
	return req.ControlGroup != nil
}

func (c *Core) doRouting(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	// If we're replicating and we get a read-only error from a backend, need to forward to primary
	resp, err := c.router.Route(ctx, req)
	if shouldForward(c, resp, err) {
		fwdResp, fwdErr := forward(ctx, c, req)
		if fwdErr != nil && err != logical.ErrReadOnly {
			// When handling the request locally, we got an error that
			// contained ErrReadOnly, but had additional information.
			// Since we've now forwarded this request and got _another_
			// error, we should tell the user about both errors, so
			// they know about both.
			//
			// When there is no error from forwarding, the request
			// succeeded and so no additional context is necessary. When
			// the initial error here was only ErrReadOnly, it's likely
			// the plugin authors intended to forward this request
			// remotely anyway.
			repErr, ok := fwdErr.(*logical.ReplicationCodedError)
			if ok {
				fwdErr = &logical.ReplicationCodedError{
					Msg:  fmt.Sprintf("errors from both primary and secondary; primary error was %s; secondary errors follow: %s", repErr.Error(), err.Error()),
					Code: repErr.Code,
				}
			} else {
				fwdErr = multierror.Append(fwdErr, err)
			}
		}
		return fwdResp, fwdErr
	}
	return resp, err
}

func (c *Core) isLoginRequest(ctx context.Context, req *logical.Request) bool {
	return c.router.LoginPath(ctx, req.Path)
}

func (c *Core) handleRequest(ctx context.Context, req *logical.Request) (retResp *logical.Response, retAuth *logical.Auth, retErr error) {
	defer metrics.MeasureSince([]string{"core", "handle_request"}, time.Now())

	var nonHMACReqDataKeys []string
	entry := c.router.MatchingMountEntry(ctx, req.Path)
	if entry != nil {
		// Set here so the audit log has it even if authorization fails
		req.MountType = entry.Type
		req.SetMountRunningSha256(entry.RunningSha256)
		req.SetMountRunningVersion(entry.RunningVersion)
		req.SetMountIsExternalPlugin(entry.IsExternalPlugin())
		req.SetMountClass(entry.MountClass())

		// Get and set ignored HMAC'd value.
		if rawVals, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
			nonHMACReqDataKeys = rawVals.([]string)
		}
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		c.logger.Error("failed to get namespace from context", "error", err)
		retErr = multierror.Append(retErr, ErrInternalError)
		return
	}

	// Validate the token
	auth, te, ctErr := c.CheckToken(ctx, req, false)
	if ctErr == logical.ErrRelativePath {
		return logical.ErrorResponse(ctErr.Error()), nil, ctErr
	}
	if ctErr == logical.ErrPerfStandbyPleaseForward {
		return nil, nil, ctErr
	}

	// Updating in-flight request data with client/entity ID
	inFlightReqID, ok := ctx.Value(logical.CtxKeyInFlightRequestID{}).(string)
	if ok && req.ClientID != "" {
		c.UpdateInFlightReqData(inFlightReqID, req.ClientID)
	}

	// We run this logic first because we want to decrement the use count even
	// in the case of an error (assuming we can successfully look up; if we
	// need to forward, we exit before now)
	if te != nil && !isControlGroupRun(req) {
		// Attempt to use the token (decrement NumUses)
		var err error
		te, err = c.tokenStore.UseToken(ctx, te)
		if err != nil {
			c.logger.Error("failed to use token", "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, nil, retErr
		}
		if te == nil {
			// Token has been revoked by this point
			retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
			return nil, nil, retErr
		}
		if te.NumUses == tokenRevocationPending {
			// We defer a revocation until after logic has run, since this is a
			// valid request (this is the token's final use). We pass the ID in
			// directly just to be safe in case something else modifies te later.
			defer func(id string) {
				nsActiveCtx := namespace.ContextWithNamespace(c.activeContext, ns)
				leaseID, err := c.expiration.CreateOrFetchRevocationLeaseByToken(nsActiveCtx, te)
				if err == nil {
					err = c.expiration.LazyRevoke(ctx, leaseID)
				}
				if err != nil {
					c.logger.Error("failed to revoke token", "error", err)
					retResp = nil
					retAuth = nil
					retErr = multierror.Append(retErr, ErrInternalError)
				}
				if retResp != nil && retResp.Secret != nil &&
					// Some backends return a TTL even without a Lease ID
					retResp.Secret.LeaseID != "" {
					retResp = logical.ErrorResponse("Secret cannot be returned; token had one use left, so leased credentials were immediately revoked.")
					return
				}
			}(te.ID)
		}
	}

	if ctErr != nil {
		newCtErr, cgResp, cgAuth, cgRetErr := checkNeedsCG(ctx, c, req, auth, ctErr, nonHMACReqDataKeys)
		switch {
		case newCtErr != nil:
			ctErr = newCtErr
		case cgResp != nil || cgAuth != nil:
			if cgRetErr != nil {
				retErr = multierror.Append(retErr, cgRetErr)
			}
			return cgResp, cgAuth, retErr
		}

		// If it is an internal error we return that, otherwise we
		// return invalid request so that the status codes can be correct
		switch {
		case ctErr == ErrInternalError,
			errwrap.Contains(ctErr, ErrInternalError.Error()),
			ctErr == logical.ErrPermissionDenied,
			errwrap.Contains(ctErr, logical.ErrPermissionDenied.Error()):
			switch ctErr.(type) {
			case *multierror.Error:
				retErr = ctErr
			default:
				retErr = multierror.Append(retErr, ctErr)
			}
		default:
			retErr = multierror.Append(retErr, logical.ErrInvalidRequest)
		}

		if !isControlGroupRun(req) {
			logInput := &logical.LogInput{
				Auth:               auth,
				Request:            req,
				OuterErr:           ctErr,
				NonHMACReqDataKeys: nonHMACReqDataKeys,
			}
			if err := c.auditBroker.LogRequest(ctx, logInput, c.auditedHeaders); err != nil {
				c.logger.Error("failed to audit request", "path", req.Path, "error", err)
			}
		}

		if errwrap.Contains(retErr, ErrInternalError.Error()) {
			return nil, auth, retErr
		}
		return logical.ErrorResponse(ctErr.Error()), auth, retErr
	}

	// Attach the display name
	req.DisplayName = auth.DisplayName

	// Create an audit trail of the request
	if !isControlGroupRun(req) {
		logInput := &logical.LogInput{
			Auth:               auth,
			Request:            req,
			NonHMACReqDataKeys: nonHMACReqDataKeys,
		}
		if err := c.auditBroker.LogRequest(ctx, logInput, c.auditedHeaders); err != nil {
			c.logger.Error("failed to audit request", "path", req.Path, "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, auth, retErr
		}
	}

	if err := enterpriseBlockRequestIfError(c, ns.Path, req.Path); err != nil {
		return nil, nil, multierror.Append(retErr, err)
	}

	leaseGenerated := false
	quotaResp, quotaErr := c.applyLeaseCountQuota(ctx, &quotas.Request{
		Path:          req.Path,
		MountPath:     strings.TrimPrefix(req.MountPoint, ns.Path),
		NamespacePath: ns.Path,
	})
	if quotaErr != nil {
		c.logger.Error("failed to apply quota", "path", req.Path, "error", quotaErr)
		retErr = multierror.Append(retErr, quotaErr)
		return nil, auth, retErr
	}

	if !quotaResp.Allowed {
		if c.logger.IsTrace() {
			c.logger.Trace("request rejected due to lease count quota violation", "request_path", req.Path)
		}

		retErr = multierror.Append(retErr, fmt.Errorf("request path %q: %w", req.Path, quotas.ErrLeaseCountQuotaExceeded))
		return nil, auth, retErr
	}

	defer func() {
		if quotaResp.Access != nil {
			quotaAckErr := c.ackLeaseQuota(quotaResp.Access, leaseGenerated)
			if quotaAckErr != nil {
				retErr = multierror.Append(retErr, quotaAckErr)
			}
		}
	}()

	// Route the request
	resp, routeErr := c.doRouting(ctx, req)
	if resp != nil {

		// If wrapping is used, use the shortest between the request and response
		var wrapTTL time.Duration
		var wrapFormat, creationPath string
		var sealWrap bool

		// Ensure no wrap info information is set other than, possibly, the TTL
		if resp.WrapInfo != nil {
			if resp.WrapInfo.TTL > 0 {
				wrapTTL = resp.WrapInfo.TTL
			}
			wrapFormat = resp.WrapInfo.Format
			creationPath = resp.WrapInfo.CreationPath
			sealWrap = resp.WrapInfo.SealWrap
			resp.WrapInfo = nil
		}

		if req.WrapInfo != nil {
			if req.WrapInfo.TTL > 0 {
				switch {
				case wrapTTL == 0:
					wrapTTL = req.WrapInfo.TTL
				case req.WrapInfo.TTL < wrapTTL:
					wrapTTL = req.WrapInfo.TTL
				}
			}
			// If the wrap format hasn't been set by the response, set it to
			// the request format
			if req.WrapInfo.Format != "" && wrapFormat == "" {
				wrapFormat = req.WrapInfo.Format
			}
		}

		if wrapTTL > 0 {
			resp.WrapInfo = &wrapping.ResponseWrapInfo{
				TTL:          wrapTTL,
				Format:       wrapFormat,
				CreationPath: creationPath,
				SealWrap:     sealWrap,
			}
		}
	}

	// If there is a secret, we must register it with the expiration manager.
	// We exclude renewal of a lease, since it does not need to be re-registered
	if resp != nil && resp.Secret != nil && !strings.HasPrefix(req.Path, "sys/renew") &&
		!strings.HasPrefix(req.Path, "sys/leases/renew") {
		// KV mounts should return the TTL but not register
		// for a lease as this provides a massive slowdown
		registerLease := true

		matchingMountEntry := c.router.MatchingMountEntry(ctx, req.Path)
		if matchingMountEntry == nil {
			c.logger.Error("unable to retrieve kv mount entry from router")
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, auth, retErr
		}

		switch matchingMountEntry.Type {
		case "kv", "generic":
			// If we are kv type, first see if we are an older passthrough
			// backend, and otherwise check the mount entry options.
			matchingBackend := c.router.MatchingBackend(ctx, req.Path)
			if matchingBackend == nil {
				c.logger.Error("unable to retrieve kv backend from router")
				retErr = multierror.Append(retErr, ErrInternalError)
				return nil, auth, retErr
			}

			if ptbe, ok := matchingBackend.(*PassthroughBackend); ok {
				if !ptbe.GeneratesLeases() {
					registerLease = false
					resp.Secret.Renewable = false
				}
			} else if matchingMountEntry.Options == nil || matchingMountEntry.Options["leased_passthrough"] != "true" {
				registerLease = false
				resp.Secret.Renewable = false
			}

		case "plugin":
			// If we are a plugin type and the plugin name is "kv" check the
			// mount entry options.
			if matchingMountEntry.Config.PluginName == "kv" && (matchingMountEntry.Options == nil || matchingMountEntry.Options["leased_passthrough"] != "true") {
				registerLease = false
				resp.Secret.Renewable = false
			}
		}

		if registerLease {
			sysView := c.router.MatchingSystemView(ctx, req.Path)
			if sysView == nil {
				c.logger.Error("unable to look up sys view for login path", "request_path", req.Path)
				return nil, nil, ErrInternalError
			}

			ttl, warnings, err := framework.CalculateTTL(sysView, 0, resp.Secret.TTL, 0, resp.Secret.MaxTTL, 0, time.Time{})
			if err != nil {
				return nil, nil, err
			}
			for _, warning := range warnings {
				resp.AddWarning(warning)
			}
			resp.Secret.TTL = ttl

			registerFunc, funcGetErr := getLeaseRegisterFunc(c)
			if funcGetErr != nil {
				retErr = multierror.Append(retErr, funcGetErr)
				return nil, auth, retErr
			}

			leaseID, err := registerFunc(ctx, req, resp, "")
			if err != nil {
				c.logger.Error("failed to register lease", "request_path", req.Path, "error", err)
				retErr = multierror.Append(retErr, ErrInternalError)
				return nil, auth, retErr
			}
			leaseGenerated = true
			resp.Secret.LeaseID = leaseID

			// Count the lease creation
			ttl_label := metricsutil.TTLBucket(resp.Secret.TTL)
			mountPointWithoutNs := ns.TrimmedPath(req.MountPoint)
			c.MetricSink().IncrCounterWithLabels(
				[]string{"secret", "lease", "creation"},
				1,
				[]metrics.Label{
					metricsutil.NamespaceLabel(ns),
					{"secret_engine", req.MountType},
					{"mount_point", mountPointWithoutNs},
					{"creation_ttl", ttl_label},
				},
			)
		}
	}

	// Only the token store is allowed to return an auth block, for any
	// other request this is an internal error.
	if resp != nil && resp.Auth != nil {
		if !strings.HasPrefix(req.Path, "auth/token/") {
			c.logger.Error("unexpected Auth response for non-token backend", "request_path", req.Path)
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, auth, retErr
		}

		// Fetch the namespace to which the token belongs
		tokenNS, err := NamespaceByID(ctx, te.NamespaceID, c)
		if err != nil {
			c.logger.Error("failed to fetch token's namespace", "error", err)
			retErr = multierror.Append(retErr, err)
			return nil, auth, retErr
		}
		if tokenNS == nil {
			c.logger.Error(namespace.ErrNoNamespace.Error())
			retErr = multierror.Append(retErr, namespace.ErrNoNamespace)
			return nil, auth, retErr
		}

		_, identityPolicies, err := c.fetchEntityAndDerivedPolicies(ctx, tokenNS, resp.Auth.EntityID, false)
		if err != nil {
			// Best-effort clean up on error, so we log the cleanup error as a
			// warning but still return as internal error.
			if err := c.tokenStore.revokeOrphan(ctx, resp.Auth.ClientToken); err != nil {
				c.logger.Warn("failed to clean up token lease from entity and policy lookup failure", "request_path", req.Path, "error", err)
			}
			return nil, nil, ErrInternalError
		}

		// We skip expiration manager registration for token renewal since it
		// does not need to be re-registered
		if strings.HasPrefix(req.Path, "auth/token/renew") {
			// We build the "policies" list to be returned by starting with
			// token policies, and add identity policies right after this
			// conditional
			tok, _ := c.DecodeSSCToken(req.InboundSSCToken)
			if resp.Auth.ClientToken == tok {
				resp.Auth.ClientToken = req.InboundSSCToken
			}
			resp.Auth.Policies = policyutil.SanitizePolicies(resp.Auth.TokenPolicies, policyutil.DoNotAddDefaultPolicy)
		} else {
			resp.Auth.TokenPolicies = policyutil.SanitizePolicies(resp.Auth.Policies, policyutil.DoNotAddDefaultPolicy)

			switch resp.Auth.TokenType {
			case logical.TokenTypeBatch:
			case logical.TokenTypeService:
				if !c.perfStandby {
					registeredTokenEntry := &logical.TokenEntry{
						TTL:         auth.TTL,
						Policies:    auth.TokenPolicies,
						Path:        resp.Auth.CreationPath,
						NamespaceID: ns.ID,
					}

					// Only logins apply to role based quotas, so we can omit the role here, as we are not logging in.
					if err := c.expiration.RegisterAuth(ctx, registeredTokenEntry, resp.Auth, ""); err != nil {
						// Best-effort clean up on error, so we log the cleanup error as
						// a warning but still return as internal error.
						if err := c.tokenStore.revokeOrphan(ctx, resp.Auth.ClientToken); err != nil {
							c.logger.Warn("failed to clean up token lease during auth/token/ request", "request_path", req.Path, "error", err)
						}
						c.logger.Error("failed to register token lease during auth/token/ request", "request_path", req.Path, "error", err)
						retErr = multierror.Append(retErr, ErrInternalError)
						return nil, auth, retErr
					}
					if registeredTokenEntry.ExternalID != "" {
						resp.Auth.ClientToken = registeredTokenEntry.ExternalID
					}
					leaseGenerated = true
				}
			}
		}

		// We do these later since it's not meaningful for backends/expmgr to
		// have what is purely a snapshot of current identity policies, and
		// plugins can be confused if they are checking contents of
		// Auth.Policies instead of Auth.TokenPolicies
		resp.Auth.Policies = policyutil.SanitizePolicies(append(resp.Auth.Policies, identityPolicies[te.NamespaceID]...), policyutil.DoNotAddDefaultPolicy)
		resp.Auth.IdentityPolicies = policyutil.SanitizePolicies(identityPolicies[te.NamespaceID], policyutil.DoNotAddDefaultPolicy)
		delete(identityPolicies, te.NamespaceID)
		resp.Auth.ExternalNamespacePolicies = identityPolicies
	}

	if resp != nil &&
		req.Path == "cubbyhole/response" &&
		len(te.Policies) == 1 &&
		te.Policies[0] == responseWrappingPolicyName {
		resp.AddWarning("Reading from 'cubbyhole/response' is deprecated. Please use sys/wrapping/unwrap to unwrap responses, as it provides additional security checks and other benefits.")
	}

	// Return the response and error
	if routeErr != nil {
		retErr = multierror.Append(retErr, routeErr)
	}

	return resp, auth, retErr
}

// handleLoginRequest is used to handle a login request, which is an
// unauthenticated request to the backend.
func (c *Core) handleLoginRequest(ctx context.Context, req *logical.Request) (retResp *logical.Response, retAuth *logical.Auth, retErr error) {
	defer metrics.MeasureSince([]string{"core", "handle_login_request"}, time.Now())

	req.Unauthenticated = true

	var nonHMACReqDataKeys []string
	entry := c.router.MatchingMountEntry(ctx, req.Path)
	if entry != nil {
		// Set here so the audit log has it even if authorization fails
		req.MountType = entry.Type
		req.SetMountRunningSha256(entry.RunningSha256)
		req.SetMountRunningVersion(entry.RunningVersion)
		req.SetMountIsExternalPlugin(entry.IsExternalPlugin())
		req.SetMountClass(entry.MountClass())

		// Get and set ignored HMAC'd value.
		if rawVals, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
			nonHMACReqDataKeys = rawVals.([]string)
		}
	}

	// Do an unauth check. This will cause EGP policies to be checked
	var auth *logical.Auth
	var ctErr error
	auth, _, ctErr = c.CheckToken(ctx, req, true)
	if ctErr == logical.ErrPerfStandbyPleaseForward {
		return nil, nil, ctErr
	}

	// Updating in-flight request data with client/entity ID
	inFlightReqID, ok := ctx.Value(logical.CtxKeyInFlightRequestID{}).(string)
	if ok && req.ClientID != "" {
		c.UpdateInFlightReqData(inFlightReqID, req.ClientID)
	}

	if ctErr != nil {
		// If it is an internal error we return that, otherwise we
		// return invalid request so that the status codes can be correct
		var errType error
		switch ctErr {
		case ErrInternalError, logical.ErrPermissionDenied:
			errType = ctErr
		default:
			errType = logical.ErrInvalidRequest
		}

		logInput := &logical.LogInput{
			Auth:               auth,
			Request:            req,
			OuterErr:           ctErr,
			NonHMACReqDataKeys: nonHMACReqDataKeys,
		}
		if err := c.auditBroker.LogRequest(ctx, logInput, c.auditedHeaders); err != nil {
			c.logger.Error("failed to audit request", "path", req.Path, "error", err)
			return nil, nil, ErrInternalError
		}

		if errType != nil {
			retErr = multierror.Append(retErr, errType)
		}
		if ctErr == ErrInternalError {
			return nil, auth, retErr
		}
		return logical.ErrorResponse(ctErr.Error()), auth, retErr
	}

	switch req.Path {
	case "sys/replication/dr/status", "sys/replication/performance/status", "sys/replication/status":
	default:
		// Create an audit trail of the request. Attach auth if it was returned,
		// e.g. if a token was provided.
		logInput := &logical.LogInput{
			Auth:               auth,
			Request:            req,
			NonHMACReqDataKeys: nonHMACReqDataKeys,
		}
		if err := c.auditBroker.LogRequest(ctx, logInput, c.auditedHeaders); err != nil {
			c.logger.Error("failed to audit request", "path", req.Path, "error", err)
			return nil, nil, ErrInternalError
		}
	}

	// The token store uses authentication even when creating a new token,
	// so it's handled in handleRequest. It should not be reached here.
	if strings.HasPrefix(req.Path, "auth/token/") {
		c.logger.Error("unexpected login request for token backend", "request_path", req.Path)
		return nil, nil, ErrInternalError
	}

	// check if user lockout feature is disabled
	isUserLockoutDisabled, err := c.isUserLockoutDisabled(entry)
	if err != nil {
		return nil, nil, err
	}

	// if user lockout feature is not disabled, check if the user is locked
	if !isUserLockoutDisabled {
		isloginUserLocked, err := c.isUserLocked(ctx, entry, req)
		if err != nil {
			return nil, nil, err
		}
		if isloginUserLocked {
			return nil, nil, logical.ErrPermissionDenied
		}
	}

	// Route the request
	resp, routeErr := c.doRouting(ctx, req)

	// if routeErr has invalid credentials error, update the userFailedLoginMap
	if routeErr != nil && routeErr == logical.ErrInvalidCredentials {
		if !isUserLockoutDisabled {
			err := c.failedUserLoginProcess(ctx, entry, req)
			if err != nil {
				return nil, nil, err
			}
		}
		return resp, nil, routeErr
	}

	if resp != nil {
		// If wrapping is used, use the shortest between the request and response
		var wrapTTL time.Duration
		var wrapFormat, creationPath string
		var sealWrap bool

		// Ensure no wrap info information is set other than, possibly, the TTL
		if resp.WrapInfo != nil {
			if resp.WrapInfo.TTL > 0 {
				wrapTTL = resp.WrapInfo.TTL
			}
			wrapFormat = resp.WrapInfo.Format
			creationPath = resp.WrapInfo.CreationPath
			sealWrap = resp.WrapInfo.SealWrap
			resp.WrapInfo = nil
		}

		if req.WrapInfo != nil {
			if req.WrapInfo.TTL > 0 {
				switch {
				case wrapTTL == 0:
					wrapTTL = req.WrapInfo.TTL
				case req.WrapInfo.TTL < wrapTTL:
					wrapTTL = req.WrapInfo.TTL
				}
			}
			if req.WrapInfo.Format != "" && wrapFormat == "" {
				wrapFormat = req.WrapInfo.Format
			}
		}

		if wrapTTL > 0 {
			resp.WrapInfo = &wrapping.ResponseWrapInfo{
				TTL:          wrapTTL,
				Format:       wrapFormat,
				CreationPath: creationPath,
				SealWrap:     sealWrap,
			}
		}
	}

	// A login request should never return a secret!
	if resp != nil && resp.Secret != nil {
		c.logger.Error("unexpected Secret response for login path", "request_path", req.Path)
		return nil, nil, ErrInternalError
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		c.logger.Error("failed to get namespace from context", "error", err)
		retErr = multierror.Append(retErr, ErrInternalError)
		return
	}
	// If the response generated an authentication, then generate the token
	if resp != nil && resp.Auth != nil && req.Path != "sys/mfa/validate" {
		leaseGenerated := false

		// by placing this after the authorization check, we don't leak
		// information about locked namespaces to unauthenticated clients.
		if err := enterpriseBlockRequestIfError(c, ns.Path, req.Path); err != nil {
			retErr = multierror.Append(retErr, err)
			return
		}

		// Check for request role in context to role based quotas
		var role string
		reqRole := ctx.Value(logical.CtxKeyRequestRole{})
		if reqRole != nil {
			role = reqRole.(string)
		}

		// The request successfully authenticated itself. Run the quota checks
		// before creating lease.
		quotaResp, quotaErr := c.applyLeaseCountQuota(ctx, &quotas.Request{
			Path:          req.Path,
			MountPath:     strings.TrimPrefix(req.MountPoint, ns.Path),
			Role:          role,
			NamespacePath: ns.Path,
		})

		if quotaErr != nil {
			c.logger.Error("failed to apply quota", "path", req.Path, "error", quotaErr)
			retErr = multierror.Append(retErr, quotaErr)
			return
		}

		if !quotaResp.Allowed {
			if c.logger.IsTrace() {
				c.logger.Trace("request rejected due to lease count quota violation", "request_path", req.Path)
			}

			retErr = multierror.Append(retErr, fmt.Errorf("request path %q: %w", req.Path, quotas.ErrLeaseCountQuotaExceeded))
			return
		}

		defer func() {
			if quotaResp.Access != nil {
				quotaAckErr := c.ackLeaseQuota(quotaResp.Access, leaseGenerated)
				if quotaAckErr != nil {
					retErr = multierror.Append(retErr, quotaAckErr)
				}
			}
		}()

		var entity *identity.Entity
		auth = resp.Auth

		mEntry := c.router.MatchingMountEntry(ctx, req.Path)

		if auth.Alias != nil &&
			mEntry != nil &&
			c.identityStore != nil {

			if mEntry.Local && os.Getenv(EnvVaultDisableLocalAuthMountEntities) != "" {
				goto CREATE_TOKEN
			}

			// Overwrite the mount type and mount path in the alias
			// information
			auth.Alias.MountType = req.MountType
			auth.Alias.MountAccessor = req.MountAccessor
			auth.Alias.Local = mEntry.Local

			if auth.Alias.Name == "" {
				return nil, nil, fmt.Errorf("missing name in alias")
			}

			var err error
			// Fetch the entity for the alias, or create an entity if one
			// doesn't exist.
			entity, entityCreated, err := c.identityStore.CreateOrFetchEntity(ctx, auth.Alias)
			if err != nil {
				switch auth.Alias.Local {
				case true:
					// Only create a new entity if the error was a readonly error and the creation flag is true
					// i.e the entity was in the middle of being created
					if entityCreated && errors.Is(err, logical.ErrReadOnly) {
						entity, err = possiblyForwardEntityCreation(ctx, c, err, auth, nil)
						if err != nil {
							if strings.Contains(err.Error(), errCreateEntityUnimplemented) {
								resp.AddWarning("primary cluster doesn't yet issue entities for local auth mounts; falling back to not issuing entities for local auth mounts")
								goto CREATE_TOKEN
							} else {
								return nil, nil, err
							}
						}
					}
					err = updateLocalAlias(ctx, c, auth, entity)
				default:
					entity, entityCreated, err = possiblyForwardAliasCreation(ctx, c, err, auth, entity)
				}
			}
			if err != nil {
				return nil, nil, err
			}
			if entity == nil {
				return nil, nil, fmt.Errorf("failed to create an entity for the authenticated alias")
			}

			if entity.Disabled {
				return nil, nil, logical.ErrPermissionDenied
			}

			auth.EntityID = entity.ID
			auth.EntityCreated = entityCreated
			validAliases, err := c.identityStore.refreshExternalGroupMembershipsByEntityID(ctx, auth.EntityID, auth.GroupAliases, req.MountAccessor)
			if err != nil {
				return nil, nil, err
			}
			auth.GroupAliases = validAliases
		}

	CREATE_TOKEN:
		// Determine the source of the login
		source := c.router.MatchingMount(ctx, req.Path)

		// Login MFA
		entity, _, err := c.fetchEntityAndDerivedPolicies(ctx, ns, auth.EntityID, false)
		if err != nil {
			return nil, nil, ErrInternalError
		}
		// finding the MFAEnforcementConfig that matches the ns and either of
		// entityID, MountAccessor, GroupID, or Auth type.
		matchedMfaEnforcementList, err := c.buildMFAEnforcementConfigList(ctx, entity, req.Path)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to find MFAEnforcement configuration, error: %v", err)
		}

		// (for the context, a response warning above says: "primary cluster
		// doesn't yet issue entities for local auth mounts; falling back
		// to not issuing entities for local auth mounts")
		// based on the above, if the entity is nil, check if MFAEnforcementConfig
		// is configured or not. If not, continue as usual, but if there
		// is something, then report an error indicating that the user is not
		// allowed to login because there is no entity associated with it.
		// This is because an entity is needed to enforce MFA.
		if entity == nil && len(matchedMfaEnforcementList) > 0 {
			// this logic means that an MFAEnforcementConfig was configured with
			// only mount type or mount accessor
			return nil, nil, logical.ErrPermissionDenied
		}

		// The resp.Auth has been populated with the information that is required for MFA validation
		// This is why, the MFA check is placed at this point. The resp.Auth is going to be fully cached
		// in memory so that it would be used to return to the user upon MFA validation is completed.
		if entity != nil {
			if len(matchedMfaEnforcementList) == 0 && len(req.MFACreds) > 0 {
				resp.AddWarning("Found MFA header but failed to find MFA Enforcement Config")
			}

			// If X-Vault-MFA header is supplied to the login request,
			// run single-phase login MFA check, else run two-phase login MFA check
			if len(matchedMfaEnforcementList) > 0 && len(req.MFACreds) > 0 {
				for _, eConfig := range matchedMfaEnforcementList {
					err = c.validateLoginMFA(ctx, eConfig, entity, req.Connection.RemoteAddr, req.MFACreds)
					if err != nil {
						return nil, nil, logical.ErrPermissionDenied
					}
				}
			} else if len(matchedMfaEnforcementList) > 0 && len(req.MFACreds) == 0 {
				mfaRequestID, err := uuid.GenerateUUID()
				if err != nil {
					return nil, nil, err
				}
				// sending back the MFARequirement config
				mfaRequirement := &logical.MFARequirement{
					MFARequestID:   mfaRequestID,
					MFAConstraints: make(map[string]*logical.MFAConstraintAny),
				}
				for _, eConfig := range matchedMfaEnforcementList {
					mfaAny, err := c.buildMfaEnforcementResponse(eConfig)
					if err != nil {
						return nil, nil, err
					}
					mfaRequirement.MFAConstraints[eConfig.Name] = mfaAny
				}

				// for two phased MFA enforcement, we should not return the regular auth
				// response. This flag is indicate to store the auth response for later
				// and return MFARequirement only
				respAuth := &MFACachedAuthResponse{
					CachedAuth:            resp.Auth,
					RequestPath:           req.Path,
					RequestNSID:           ns.ID,
					RequestNSPath:         ns.Path,
					RequestConnRemoteAddr: req.Connection.RemoteAddr, // this is needed for the DUO method
					TimeOfStorage:         time.Now(),
					RequestID:             mfaRequestID,
				}
				err = possiblyForwardSaveCachedAuthResponse(ctx, c, respAuth)
				if err != nil {
					return nil, nil, err
				}
				auth = nil
				resp.Auth = &logical.Auth{
					MFARequirement: mfaRequirement,
				}
				resp.AddWarning("A login request was issued that is subject to MFA validation. Please make sure to validate the login by sending another request to mfa/validate endpoint.")
				// going to return early before generating the token
				// the user receives the mfaRequirement, and need to use the
				// login MFA validate endpoint to get the token
				return resp, auth, nil
			}
		}

		// Attach the display name, might be used by audit backends
		req.DisplayName = auth.DisplayName

		requiresLease := resp.Auth.TokenType != logical.TokenTypeBatch

		// If role was not already determined by http.rateLimitQuotaWrapping
		// and a lease will be generated, calculate a role for the leaseEntry.
		// We can skip this step if there are no pre-existing role-based quotas
		// for this mount and Vault is configured to skip lease role-based lease counting
		// until after they're created. This effectively zeroes out the lease count
		// for new role-based quotas upon creation, rather than counting old leases toward
		// the total.
		if reqRole == nil && requiresLease && !c.impreciseLeaseRoleTracking {
			role = c.DetermineRoleFromLoginRequest(ctx, req.MountPoint, req.Data)
		}

		leaseGen, respTokenCreate, errCreateToken := c.LoginCreateToken(ctx, ns, req.Path, source, role, resp)
		leaseGenerated = leaseGen
		if errCreateToken != nil {
			return respTokenCreate, nil, errCreateToken
		}
		resp = respTokenCreate
	}

	// Successful login, remove any entry from userFailedLoginInfo map
	// if it exists. This is done for batch tokens (for oss & ent)
	// For service tokens on oss it is taken care by core RegisterAuth function.
	// For service tokens on ent it is taken care by registerAuth RPC calls.
	// This update is done as part of registerAuth of RPC calls from standby
	// to active node. This is added there to reduce RPC calls
	if !isUserLockoutDisabled && (auth.TokenType == logical.TokenTypeBatch) {
		loginUserInfoKey := FailedLoginUser{
			aliasName:     auth.Alias.Name,
			mountAccessor: auth.Alias.MountAccessor,
		}

		// We don't need to try to delete the lockedUsers storage entry, since we're
		// processing a login request. If a login attempt is allowed, it means the user is
		// unlocked and we only add storage entry when the user gets locked.
		err = updateUserFailedLoginInfo(ctx, c, loginUserInfoKey, nil, true)
		if err != nil {
			return nil, nil, err
		}
	}

	// if we were already going to return some error from this login, do that.
	// if not, we will then check if the API is locked for the requesting
	// namespace, to avoid leaking locked namespaces to unauthenticated clients.
	if resp != nil && resp.Data != nil {
		if _, ok := resp.Data["error"]; ok {
			return resp, auth, routeErr
		}
	}
	if routeErr != nil {
		return resp, auth, routeErr
	}

	// this check handles the bad login credential case
	if err := enterpriseBlockRequestIfError(c, ns.Path, req.Path); err != nil {
		return nil, nil, multierror.Append(retErr, err)
	}

	return resp, auth, routeErr
}

// LoginCreateToken creates a token as a result of a login request.
// If MFA is enforced, mfa/validate endpoint calls this functions
// after successful MFA validation to generate the token.
func (c *Core) LoginCreateToken(ctx context.Context, ns *namespace.Namespace, reqPath, mountPoint, role string, resp *logical.Response) (bool, *logical.Response, error) {
	auth := resp.Auth
	source := strings.TrimPrefix(mountPoint, credentialRoutePrefix)
	source = strings.ReplaceAll(source, "/", "-")

	// Prepend the source to the display name
	auth.DisplayName = strings.TrimSuffix(source+auth.DisplayName, "-")

	// Determine mount type
	mountEntry := c.router.MatchingMountEntry(ctx, reqPath)
	if mountEntry == nil {
		return false, nil, fmt.Errorf("failed to find a matching mount")
	}

	sysView := c.router.MatchingSystemView(ctx, reqPath)
	if sysView == nil {
		c.logger.Error("unable to look up sys view for login path", "request_path", reqPath)
		return false, nil, ErrInternalError
	}

	tokenTTL, warnings, err := framework.CalculateTTL(sysView, 0, auth.TTL, auth.Period, auth.MaxTTL, auth.ExplicitMaxTTL, time.Time{})
	if err != nil {
		return false, nil, err
	}
	for _, warning := range warnings {
		resp.AddWarning(warning)
	}

	_, identityPolicies, err := c.fetchEntityAndDerivedPolicies(ctx, ns, auth.EntityID, false)
	if err != nil {
		return false, nil, ErrInternalError
	}

	auth.TokenPolicies = policyutil.SanitizePolicies(auth.Policies, !auth.NoDefaultPolicy)
	allPolicies := policyutil.SanitizePolicies(append(auth.TokenPolicies, identityPolicies[ns.ID]...), policyutil.DoNotAddDefaultPolicy)

	// Prevent internal policies from being assigned to tokens. We check
	// this on auth.Policies including derived ones from Identity before
	// actually making the token.
	for _, policy := range allPolicies {
		if policy == "root" {
			return false, logical.ErrorResponse("auth methods cannot create root tokens"), logical.ErrInvalidRequest
		}
		if strutil.StrListContains(nonAssignablePolicies, policy) {
			return false, logical.ErrorResponse(fmt.Sprintf("cannot assign policy %q", policy)), logical.ErrInvalidRequest
		}
	}

	var registerFunc RegisterAuthFunc
	var funcGetErr error
	// Batch tokens should not be forwarded to perf standby
	if auth.TokenType == logical.TokenTypeBatch {
		registerFunc = c.RegisterAuth
	} else {
		registerFunc, funcGetErr = getAuthRegisterFunc(c)
	}
	if funcGetErr != nil {
		return false, nil, funcGetErr
	}

	leaseGenerated := false
	err = registerFunc(ctx, tokenTTL, reqPath, auth, role)
	switch {
	case err == nil:
		if auth.TokenType != logical.TokenTypeBatch {
			leaseGenerated = true
		}
	case err == ErrInternalError:
		return false, nil, err
	default:
		return false, logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	auth.IdentityPolicies = policyutil.SanitizePolicies(identityPolicies[ns.ID], policyutil.DoNotAddDefaultPolicy)
	delete(identityPolicies, ns.ID)
	auth.ExternalNamespacePolicies = identityPolicies
	auth.Policies = allPolicies

	// Count the successful token creation
	ttl_label := metricsutil.TTLBucket(tokenTTL)
	// Do not include namespace path in mount point; already present as separate label.
	mountPointWithoutNs := ns.TrimmedPath(mountPoint)
	c.metricSink.IncrCounterWithLabels(
		[]string{"token", "creation"},
		1,
		[]metrics.Label{
			metricsutil.NamespaceLabel(ns),
			{"auth_method", mountEntry.Type},
			{"mount_point", mountPointWithoutNs},
			{"creation_ttl", ttl_label},
			{"token_type", auth.TokenType.String()},
		},
	)

	return leaseGenerated, resp, nil
}

// failedUserLoginProcess updates the userFailedLoginMap with login count and  last failed
// login time for users with failed login attempt
// If the user gets locked for current login attempt, it updates the storage entry too
func (c *Core) failedUserLoginProcess(ctx context.Context, mountEntry *MountEntry, req *logical.Request) error {
	// get the user lockout configuration for the user
	userLockoutConfiguration := c.getUserLockoutConfiguration(mountEntry)

	// determine the key for userFailedLoginInfo map
	loginUserInfoKey, err := c.getLoginUserInfoKey(ctx, mountEntry, req)
	if err != nil {
		return err
	}

	// get entry from userFailedLoginInfo map for the key
	userFailedLoginInfo, err := getUserFailedLoginInfo(ctx, c, loginUserInfoKey)
	if err != nil {
		return err
	}

	// update the last failed login time with current time
	failedLoginInfo := FailedLoginInfo{
		lastFailedLoginTime: int(time.Now().Unix()),
	}

	// set the failed login count value for the entry in userFailedLoginInfo map
	switch userFailedLoginInfo {
	case nil: // entry does not exist in userfailedLoginMap
		failedLoginInfo.count = 1
	default:
		failedLoginInfo.count = userFailedLoginInfo.count + 1

		// if counter reset, set the count value to 1 as this gets counted as new entry
		lastFailedLoginTime := time.Unix(int64(userFailedLoginInfo.lastFailedLoginTime), 0)
		counterResetDuration := userLockoutConfiguration.LockoutCounterReset
		if time.Now().After(lastFailedLoginTime.Add(counterResetDuration)) {
			failedLoginInfo.count = 1
		}
	}

	// update the userFailedLoginInfo map (and/or storage) with the updated/new entry
	err = updateUserFailedLoginInfo(ctx, c, loginUserInfoKey, &failedLoginInfo, false)
	if err != nil {
		return err
	}

	return nil
}

// getLoginUserInfoKey gets failedUserLoginInfo map key for login user
func (c *Core) getLoginUserInfoKey(ctx context.Context, mountEntry *MountEntry, req *logical.Request) (FailedLoginUser, error) {
	userInfo := FailedLoginUser{}
	aliasName, err := c.aliasNameFromLoginRequest(ctx, req)
	if err != nil {
		return userInfo, err
	}
	if aliasName == "" {
		return userInfo, errors.New("failed to determine alias name from login request")
	}

	userInfo.aliasName = aliasName
	userInfo.mountAccessor = mountEntry.Accessor
	return userInfo, nil
}

// isUserLockoutDisabled checks if user lockout feature to prevent brute forcing is disabled
// Auth types userpass, ldap and approle support this feature
// precedence: environment var setting >> auth tune setting >> config file setting >> default (enabled)
func (c *Core) isUserLockoutDisabled(mountEntry *MountEntry) (bool, error) {
	if !strutil.StrListContains(configutil.GetSupportedUserLockoutsAuthMethods(), mountEntry.Type) {
		return true, nil
	}

	// check environment variable
	if disableUserLockoutEnv := os.Getenv(consts.VaultDisableUserLockout); disableUserLockoutEnv != "" {
		var err error
		disableUserLockout, err := strconv.ParseBool(disableUserLockoutEnv)
		if err != nil {
			return false, errors.New("Error parsing the environment variable VAULT_DISABLE_USER_LOCKOUT")
		}
		if disableUserLockout {
			return true, nil
		}
		return false, nil
	}

	// read auth tune for mount entry
	userLockoutConfigFromMount := mountEntry.Config.UserLockoutConfig
	if userLockoutConfigFromMount != nil && userLockoutConfigFromMount.DisableLockout {
		return true, nil
	}

	// read config for auth type from config file
	userLockoutConfiguration := c.getUserLockoutFromConfig(mountEntry.Type)
	if userLockoutConfiguration.DisableLockout {
		return true, nil
	}

	// default
	return false, nil
}

// isUserLocked determines if the login user is locked
func (c *Core) isUserLocked(ctx context.Context, mountEntry *MountEntry, req *logical.Request) (locked bool, err error) {
	// get userFailedLoginInfo map key for login user
	loginUserInfoKey, err := c.getLoginUserInfoKey(ctx, mountEntry, req)
	if err != nil {
		return false, err
	}

	// get entry from userFailedLoginInfo map for the key
	userFailedLoginInfo, err := getUserFailedLoginInfo(ctx, c, loginUserInfoKey)
	if err != nil {
		return false, err
	}
	userLockoutConfiguration := c.getUserLockoutConfiguration(mountEntry)

	switch userFailedLoginInfo {
	case nil:
		// entry not found in userFailedLoginInfo map, check storage to re-verify
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return false, fmt.Errorf("could not parse namespace from http context: %w", err)
		}
		storageUserLockoutPath := fmt.Sprintf(coreLockedUsersPath+"%s/%s/%s", ns.ID, loginUserInfoKey.mountAccessor, loginUserInfoKey.aliasName)
		existingEntry, err := c.barrier.Get(ctx, storageUserLockoutPath)
		if err != nil {
			return false, err
		}
		var lastLoginTime int
		if existingEntry == nil {
			// no storage entry found, user is not locked
			return false, nil
		}

		err = jsonutil.DecodeJSON(existingEntry.Value, &lastLoginTime)
		if err != nil {
			return false, err
		}

		// if time passed from last login time is within lockout duration, the user is locked
		if time.Now().Unix()-int64(lastLoginTime) < int64(userLockoutConfiguration.LockoutDuration.Seconds()) {
			// user locked
			return true, nil
		}

		// else user is not locked. Entry is stale, this will be removed from storage during cleanup
		// by the background thread

	default:
		// entry found in userFailedLoginInfo map, check if the user is locked
		isCountOverLockoutThreshold := userFailedLoginInfo.count >= uint(userLockoutConfiguration.LockoutThreshold)
		isWithinLockoutDuration := time.Now().Unix()-int64(userFailedLoginInfo.lastFailedLoginTime) < int64(userLockoutConfiguration.LockoutDuration.Seconds())

		if isCountOverLockoutThreshold && isWithinLockoutDuration {
			// user locked
			return true, nil
		}
	}
	return false, nil
}

// getUserLockoutConfiguration gets the user lockout configuration for a mount entry
// it checks the config file and auth tune values
// precedence: auth tune >> config file values for auth type >> config file values for all type
// >> default user lockout values
// getUserLockoutFromConfig call in this function takes care of config file precedence
func (c *Core) getUserLockoutConfiguration(mountEntry *MountEntry) (userLockoutConfig UserLockoutConfig) {
	// get user configuration values from config file
	userLockoutConfig = c.getUserLockoutFromConfig(mountEntry.Type)

	authTuneUserLockoutConfig := mountEntry.Config.UserLockoutConfig
	// if user lockout is not configured using auth tune, return values from config file
	if authTuneUserLockoutConfig == nil {
		return userLockoutConfig
	}
	// replace values in return with config file configuration
	// for fields that are not configured using auth tune
	if authTuneUserLockoutConfig.LockoutThreshold != 0 {
		userLockoutConfig.LockoutThreshold = authTuneUserLockoutConfig.LockoutThreshold
	}
	if authTuneUserLockoutConfig.LockoutDuration != 0 {
		userLockoutConfig.LockoutDuration = authTuneUserLockoutConfig.LockoutDuration
	}
	if authTuneUserLockoutConfig.LockoutCounterReset != 0 {
		userLockoutConfig.LockoutCounterReset = authTuneUserLockoutConfig.LockoutCounterReset
	}
	if authTuneUserLockoutConfig.DisableLockout {
		userLockoutConfig.DisableLockout = authTuneUserLockoutConfig.DisableLockout
	}
	return userLockoutConfig
}

// getUserLockoutFromConfig gets the userlockout configuration for given mount type from config file
// it reads the user lockout configuration from server config
// it has values for "all" type and any mountType that is configured using config file
// "all" type values are updated in shared config with default values i.e; if "all" type is
// not configured in config file, it is updated in shared config with default configuration
// If "all" type is configured in config file, any missing fields are updated with default values
// similarly missing values for a given mount type in config file are updated with "all" type
// default values
// If user_lockout configuration is not configured using config file at all, defaults are returned
func (c *Core) getUserLockoutFromConfig(mountType string) UserLockoutConfig {
	defaultUserLockoutConfig := UserLockoutConfig{
		LockoutThreshold:    configutil.UserLockoutThresholdDefault,
		LockoutDuration:     configutil.UserLockoutDurationDefault,
		LockoutCounterReset: configutil.UserLockoutCounterResetDefault,
		DisableLockout:      configutil.DisableUserLockoutDefault,
	}
	conf := c.rawConfig.Load()
	if conf == nil {
		return defaultUserLockoutConfig
	}
	userlockouts := conf.(*server.Config).UserLockouts
	if userlockouts == nil {
		return defaultUserLockoutConfig
	}
	for _, userLockoutConfig := range userlockouts {
		switch userLockoutConfig.Type {
		case "all":
			defaultUserLockoutConfig = UserLockoutConfig{
				LockoutThreshold:    userLockoutConfig.LockoutThreshold,
				LockoutDuration:     userLockoutConfig.LockoutDuration,
				LockoutCounterReset: userLockoutConfig.LockoutCounterReset,
				DisableLockout:      userLockoutConfig.DisableLockout,
			}
		case mountType:
			return UserLockoutConfig{
				LockoutThreshold:    userLockoutConfig.LockoutThreshold,
				LockoutDuration:     userLockoutConfig.LockoutDuration,
				LockoutCounterReset: userLockoutConfig.LockoutCounterReset,
				DisableLockout:      userLockoutConfig.DisableLockout,
			}

		}
	}
	return defaultUserLockoutConfig
}

func (c *Core) buildMfaEnforcementResponse(eConfig *mfa.MFAEnforcementConfig) (*logical.MFAConstraintAny, error) {
	mfaAny := &logical.MFAConstraintAny{
		Any: []*logical.MFAMethodID{},
	}
	for _, methodID := range eConfig.MFAMethodIDs {
		mConfig, err := c.loginMFABackend.MemDBMFAConfigByID(methodID)
		if err != nil {
			return nil, fmt.Errorf("failed to get methodID %s from MFA config table, error: %v", methodID, err)
		}
		var duoUsePasscode bool
		if mConfig.Type == mfaMethodTypeDuo {
			duoConf, ok := mConfig.Config.(*mfa.Config_DuoConfig)
			if !ok {
				return nil, fmt.Errorf("invalid MFA configuration type")
			}
			duoUsePasscode = duoConf.DuoConfig.UsePasscode
		}
		mfaMethod := &logical.MFAMethodID{
			Type:         mConfig.Type,
			ID:           methodID,
			UsesPasscode: mConfig.Type == mfaMethodTypeTOTP || duoUsePasscode,
			Name:         mConfig.Name,
		}
		mfaAny.Any = append(mfaAny.Any, mfaMethod)
	}
	return mfaAny, nil
}

func blockRequestIfErrorImpl(_ *Core, _, _ string) error { return nil }

// RegisterAuth uses a logical.Auth object to create a token entry in the token
// store, and registers a corresponding token lease to the expiration manager.
// role is the login role used as part of the creation of the token entry. If not
// relevant, can be omitted (by being provided as "").
func (c *Core) RegisterAuth(ctx context.Context, tokenTTL time.Duration, path string, auth *logical.Auth, role string) error {
	// We first assign token policies to what was returned from the backend
	// via auth.Policies. Then, we get the full set of policies into
	// auth.Policies from the backend + entity information -- this is not
	// stored in the token, but we perform sanity checks on it and return
	// that information to the user.

	// Generate a token
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	te := logical.TokenEntry{
		Path:           path,
		Meta:           auth.Metadata,
		DisplayName:    auth.DisplayName,
		CreationTime:   time.Now().Unix(),
		TTL:            tokenTTL,
		NumUses:        auth.NumUses,
		EntityID:       auth.EntityID,
		BoundCIDRs:     auth.BoundCIDRs,
		Policies:       auth.TokenPolicies,
		NamespaceID:    ns.ID,
		ExplicitMaxTTL: auth.ExplicitMaxTTL,
		Period:         auth.Period,
		Type:           auth.TokenType,
	}

	if te.TTL == 0 && (len(te.Policies) != 1 || te.Policies[0] != "root") {
		c.logger.Error("refusing to create a non-root zero TTL token")
		return ErrInternalError
	}

	if err := c.tokenStore.create(ctx, &te); err != nil {
		c.logger.Error("failed to create token", "error", err)
		return ErrInternalError
	}

	// Populate the client token, accessor, and TTL
	auth.ClientToken = te.ID
	auth.Accessor = te.Accessor
	auth.TTL = te.TTL
	auth.Orphan = te.Parent == ""

	switch auth.TokenType {
	case logical.TokenTypeBatch:
		// Ensure it's not marked renewable since it isn't
		auth.Renewable = false
	case logical.TokenTypeService:
		// Register with the expiration manager
		if err := c.expiration.RegisterAuth(ctx, &te, auth, role); err != nil {
			if err := c.tokenStore.revokeOrphan(ctx, te.ID); err != nil {
				c.logger.Warn("failed to clean up token lease during login request", "request_path", path, "error", err)
			}
			c.logger.Error("failed to register token lease during login request", "request_path", path, "error", err)
			return ErrInternalError
		}
		if te.ExternalID != "" {
			auth.ClientToken = te.ExternalID
		}
		// Successful login, remove any entry from userFailedLoginInfo map
		// if it exists. This is done for service tokens (for oss) here.
		// For ent it is taken care by registerAuth RPC calls.
		if auth.Alias != nil {
			loginUserInfoKey := FailedLoginUser{
				aliasName:     auth.Alias.Name,
				mountAccessor: auth.Alias.MountAccessor,
			}

			// We don't need to try to delete the lockedUsers storage entry, since we're
			// processing a login request. If a login attempt is allowed, it means the user is
			// unlocked and we only add storage entry when the user gets locked.
			err = updateUserFailedLoginInfo(ctx, c, loginUserInfoKey, nil, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// LocalGetUserFailedLoginInfo gets the failed login information for a user based on alias name and mountAccessor
func (c *Core) LocalGetUserFailedLoginInfo(ctx context.Context, userKey FailedLoginUser) *FailedLoginInfo {
	c.userFailedLoginInfoLock.Lock()
	value, exists := c.userFailedLoginInfo[userKey]
	c.userFailedLoginInfoLock.Unlock()
	if exists {
		return value
	}
	return nil
}

// LocalUpdateUserFailedLoginInfo updates the failed login information for a user based on alias name and mountAccessor
func (c *Core) LocalUpdateUserFailedLoginInfo(ctx context.Context, userKey FailedLoginUser, failedLoginInfo *FailedLoginInfo, deleteEntry bool) error {
	c.userFailedLoginInfoLock.Lock()
	switch deleteEntry {
	case false:
		// update entry in the map
		c.userFailedLoginInfo[userKey] = failedLoginInfo

		// get the user lockout configuration for the user
		mountEntry := c.router.MatchingMountByAccessor(userKey.mountAccessor)
		if mountEntry == nil {
			mountEntry = &MountEntry{}
			mountEntry.NamespaceID = namespace.RootNamespaceID
		}
		userLockoutConfiguration := c.getUserLockoutConfiguration(mountEntry)

		// if failed login count has reached threshold, create a storage entry as the user got locked
		if failedLoginInfo.count >= uint(userLockoutConfiguration.LockoutThreshold) {
			// user locked
			storageUserLockoutPath := fmt.Sprintf(coreLockedUsersPath+"%s/%s/%s", mountEntry.NamespaceID, userKey.mountAccessor, userKey.aliasName)

			compressedBytes, err := jsonutil.EncodeJSONAndCompress(failedLoginInfo.lastFailedLoginTime, nil)
			if err != nil {
				c.logger.Error("failed to encode or compress failed login user entry", "error", err)
				return err
			}

			// Create an entry
			entry := &logical.StorageEntry{
				Key:   storageUserLockoutPath,
				Value: compressedBytes,
			}

			// Write to the physical backend
			if err := c.barrier.Put(ctx, entry); err != nil {
				c.logger.Error("failed to persist failed login user entry", "error", err)
				return err
			}

		}

	default:
		// delete the entry from the map, if no key exists it is no-op
		delete(c.userFailedLoginInfo, userKey)
	}
	c.userFailedLoginInfoLock.Unlock()
	return nil
}

// PopulateTokenEntry looks up req.ClientToken in the token store and uses
// it to set other fields in req.  Does nothing if ClientToken is empty
// or a JWT token, or for service tokens that don't exist in the token store.
// Should be called with read stateLock held.
func (c *Core) PopulateTokenEntry(ctx context.Context, req *logical.Request) error {
	if req.ClientToken == "" {
		return nil
	}

	// Also attach the accessor if we have it. This doesn't fail if it
	// doesn't exist because the request may be to an unauthenticated
	// endpoint/login endpoint where a bad current token doesn't matter, or
	// a token from a Vault version pre-accessors. We ignore errors for
	// JWTs.
	token := req.ClientToken
	var err error
	req.InboundSSCToken = token
	decodedToken := token
	if IsSSCToken(token) {
		// If ForwardToActive is set to ForwardSSCTokenToActive, we ignore
		// whether the endpoint is a login request, as since we have the token
		// forwarded to us, we should treat it as an unauthenticated endpoint
		// and ensure the token is populated too regardless.
		// Notably, this is important for some endpoints, such as endpoints
		// such as sys/ui/mounts/internal, which is unauthenticated but a token
		// may be provided to be used.
		// Without the check to see if
		// c.ForwardToActive() == ForwardSSCTokenToActive unauthenticated
		// requests that do not use a token but were provided one anyway
		// could fail with a 412.
		// We only follow this behaviour if we're a perf standby, as
		// this behaviour only makes sense in that case as only they
		// could be missing the token population.
		// Without ForwardToActive being set to ForwardSSCTokenToActive,
		// behaviours that rely on this functionality also wouldn't make
		// much sense, as they would fail with 412 required index not present
		// as perf standbys aren't guaranteed to have the WAL state
		// for new tokens.
		unauth := c.isLoginRequest(ctx, req)
		if c.ForwardToActive() == ForwardSSCTokenToActive && c.perfStandby {
			unauth = false
		}
		decodedToken, err = c.CheckSSCToken(ctx, token, unauth, c.perfStandby)
		// If we receive an error from CheckSSCToken, we can assume the token is bad somehow, and the client
		// should receive a 403 bad token error like they do for all other invalid tokens, unless the error
		// specifies that we should forward the request or retry the request.
		if err != nil {
			if errors.Is(err, logical.ErrPerfStandbyPleaseForward) || errors.Is(err, logical.ErrMissingRequiredState) {
				return err
			}
			return logical.ErrPermissionDenied
		}
	}
	req.ClientToken = decodedToken
	// We ignore the token returned from CheckSSCToken here as Lookup also
	// decodes the SSCT, and it may need the original SSCT to check state.
	te, err := c.LookupToken(ctx, token)
	if err != nil {
		// If we're missing required state, return that error
		// as-is to the client
		if errors.Is(err, logical.ErrPerfStandbyPleaseForward) || errors.Is(err, logical.ErrMissingRequiredState) {
			return err
		}
		// If we have two dots but the second char is a dot it's a vault
		// token of the form s.SOMETHING.nsid, not a JWT
		if !IsJWT(token) {
			return fmt.Errorf("error performing token check: %w", err)
		}
	}
	if err == nil && te != nil {
		req.ClientTokenAccessor = te.Accessor
		req.ClientTokenRemainingUses = te.NumUses
		req.SetTokenEntry(te)
	}
	return nil
}

func (c *Core) CheckSSCToken(ctx context.Context, token string, unauth bool, isPerfStandby bool) (string, error) {
	if unauth && token != "" {
		// This token shouldn't really be here, but alas it was sent along with the request
		// Since we're already knee deep in the token checking code pre-existing token checking
		// code, we have to deal with this token whether we like it or not. So, we'll just try
		// to get the inner token, and if that fails, return the token as-is. We intentionally
		// will skip any token checks, because this is an unauthenticated paths and the token
		// is just a nuisance rather than a means of auth.

		// We cannot return whatever we like here, because if we do then CheckToken, which looks up
		// the corresponding lease, will not find the token entry and lease. There are unauth'ed
		// endpoints that use the token entry (such as sys/ui/mounts/internal) to do custom token
		// checks, which would then fail. Therefore, we must try to get whatever thing is tied to
		// token entries, but we must explicitly not do any SSC Token checks.
		tok, err := c.DecodeSSCToken(token)
		if err != nil || tok == "" {
			return token, nil
		}
		return tok, nil
	}
	return c.checkSSCTokenInternal(ctx, token, isPerfStandby)
}

// DecodeSSCToken returns the random part of an SSCToken without
// performing any signature or WAL checks.
func (c *Core) DecodeSSCToken(token string) (string, error) {
	// Skip batch and old style service tokens. These can have the prefix "b.",
	// "s." (for old tokens) or "hvb."
	if !IsSSCToken(token) {
		return token, nil
	}
	tok, err := DecodeSSCTokenInternal(token)
	if err != nil {
		return "", err
	}
	return tok.Random, nil
}

// DecodeSSCTokenInternal is a helper used to get the inner part of a SSC token without
// checking the token signature or the WAL index.
func DecodeSSCTokenInternal(token string) (*tokens.Token, error) {
	signedToken := &tokens.SignedToken{}

	// Skip batch and old style service tokens. These can have the prefix "b.",
	// "s." (for old tokens) or "hvb."
	if !strings.HasPrefix(token, consts.ServiceTokenPrefix) {
		return nil, fmt.Errorf("not service token")
	}

	// Consider the suffix of the token only when unmarshalling
	suffixToken := token[4:]

	tokenBytes, err := base64.RawURLEncoding.DecodeString(suffixToken)
	if err != nil {
		return nil, fmt.Errorf("can't decode token")
	}

	err = proto.Unmarshal(tokenBytes, signedToken)
	if err != nil {
		return nil, err
	}
	plainToken := &tokens.Token{}
	err2 := proto.Unmarshal([]byte(signedToken.Token), plainToken)
	if err2 != nil {
		return nil, err2
	}
	return plainToken, nil
}

func (c *Core) checkSSCTokenInternal(ctx context.Context, token string, isPerfStandby bool) (string, error) {
	signedToken := &tokens.SignedToken{}

	// Skip batch and old style service tokens. These can have the prefix "b.",
	// "s." (for old tokens) or "hvb."
	if !strings.HasPrefix(token, consts.ServiceTokenPrefix) {
		return token, nil
	}

	// Check token length to guess if this is an server side consistent token or not.
	// Note that even when the DisableSSCTokens flag is set, index
	// bearing tokens that have already been given out may still be used.
	if !IsSSCToken(token) {
		return token, nil
	}

	// Consider the suffix of the token only when unmarshalling
	suffixToken := token[4:]

	tokenBytes, err := base64.RawURLEncoding.DecodeString(suffixToken)
	if err != nil {
		c.logger.Warn("cannot decode token", "error", err)
		return token, nil
	}

	err = proto.Unmarshal(tokenBytes, signedToken)
	if err != nil {
		return "", fmt.Errorf("error occurred when unmarshalling ssc token: %w", err)
	}
	hm, err := c.tokenStore.CalculateSignedTokenHMAC(signedToken.Token)
	if !hmac.Equal(hm, signedToken.Hmac) {
		return "", fmt.Errorf("token mac for %+v is incorrect: err %w", signedToken, err)
	}
	plainToken := &tokens.Token{}
	err = proto.Unmarshal([]byte(signedToken.Token), plainToken)
	if err != nil {
		return "", err
	}

	// Disregard SSCT on perf-standbys for non-raft storage
	if c.perfStandby && c.getRaftBackend() == nil {
		return plainToken.Random, nil
	}

	ep := int(plainToken.IndexEpoch)
	if ep < c.tokenStore.GetSSCTokensGenerationCounter() {
		return plainToken.Random, nil
	}

	requiredWalState := &logical.WALState{ClusterID: c.ClusterID(), LocalIndex: plainToken.LocalIndex, ReplicatedIndex: 0}
	if c.HasWALState(requiredWalState, isPerfStandby) {
		return plainToken.Random, nil
	}
	// Make sure to forward the request instead of checking the token if the flag
	// is set and we're on a perf standby
	if c.ForwardToActive() == ForwardSSCTokenToActive && isPerfStandby {
		return "", logical.ErrPerfStandbyPleaseForward
	}
	// In this case, the server side consistent token cannot be used on this node. We return the appropriate
	// status code.
	return "", logical.ErrMissingRequiredState
}
