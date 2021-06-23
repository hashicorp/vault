package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	multierror "github.com/hashicorp/go-multierror"
	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/quotas"
	uberAtomic "go.uber.org/atomic"
)

const (
	replTimeout = 1 * time.Second
)

var (
	// DefaultMaxRequestDuration is the amount of time we'll wait for a request
	// to complete, unless overridden on a per-handler basis
	DefaultMaxRequestDuration = 90 * time.Second

	egpDebugLogging bool
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
// also returns the cumulative list of policies that the entity is entitled to.
// This list includes the policies from the entity itself and from all the
// groups in which the given entity ID is a member of.
func (c *Core) fetchEntityAndDerivedPolicies(ctx context.Context, tokenNS *namespace.Namespace, entityID string) (*identity.Entity, map[string][]string, error) {
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
	if entity != nil {
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

		// Filter and add the policies to the resultant set
		for nsID, nsPolicies := range groupPolicies {
			ns, err := NamespaceByID(ctx, nsID, c)
			if err != nil {
				return nil, nil, err
			}
			if ns == nil {
				return nil, nil, namespace.ErrNoNamespace
			}
			if tokenNS.Path != ns.Path && !ns.HasParent(tokenNS) {
				continue
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
		return nil, nil, nil, nil, &logical.StatusBadRequest{Err: "missing client token"}
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
			c.logger.Error("failed to lookup token", "error", err)
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

	policies := make(map[string][]string)
	// Add tokens policies
	policies[te.NamespaceID] = append(policies[te.NamespaceID], te.Policies...)

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
	entity, identityPolicies, err := c.fetchEntityAndDerivedPolicies(ctx, tokenNS, te.EntityID)
	if err != nil {
		return nil, nil, nil, nil, ErrInternalError
	}
	for nsID, nsPolicies := range identityPolicies {
		policies[nsID] = append(policies[nsID], nsPolicies...)
	}

	// Attach token's namespace information to the context. Wrapping tokens by
	// should be able to be used anywhere, so we also special case behavior.
	var tokenCtx context.Context
	if len(policies) == 1 &&
		len(policies[te.NamespaceID]) == 1 &&
		(policies[te.NamespaceID][0] == responseWrappingPolicyName ||
			policies[te.NamespaceID][0] == controlGroupPolicyName) &&
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

	// Construct the corresponding ACL object. ACL construction should be
	// performed on the token's namespace.
	acl, err := c.policyStore.ACL(tokenCtx, entity, policies)
	if err != nil {
		if errwrap.ContainsType(err, new(TemplateError)) {
			c.logger.Warn("permission denied due to a templated policy being invalid or containing directives not satisfied by the requestor", "error", err)
			return nil, nil, nil, nil, logical.ErrPermissionDenied
		}
		c.logger.Error("failed to construct ACL", "error", err)
		return nil, nil, nil, nil, ErrInternalError
	}

	return acl, te, entity, identityPolicies, nil
}

func (c *Core) checkToken(ctx context.Context, req *logical.Request, unauth bool) (*logical.Auth, *logical.TokenEntry, error) {
	defer metrics.MeasureSince([]string{"core", "check_token"}, time.Now())

	var acl *ACL
	var te *logical.TokenEntry
	var entity *identity.Entity
	var identityPolicies map[string][]string
	var err error

	// Even if unauth, if a token is provided, there's little reason not to
	// gather as much info as possible for the audit log and to e.g. control
	// trace mode for EGPs.
	if !unauth || (unauth && req.ClientToken != "") {
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
		case checkExists == false:
			// No existence check, so always treat it as an update operation, which is how it is pre 0.5
			req.Operation = logical.UpdateOperation
		case resourceExists == true:
			// It exists, so force an update operation
			req.Operation = logical.UpdateOperation
		case resourceExists == false:
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

	if te != nil {
		auth.IdentityPolicies = identityPolicies[te.NamespaceID]
		auth.TokenPolicies = te.Policies
		auth.Policies = append(te.Policies, identityPolicies[te.NamespaceID]...)
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
	}

	// Check the standard non-root ACLs. Return the token entry if it's not
	// allowed so we can decrement the use count.
	authResults := c.performPolicyChecks(ctx, acl, te, req, entity, &PolicyCheckOpts{
		Unauth:            unauth,
		RootPrivsRequired: rootPath,
	})

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

	resp, err = c.handleCancelableRequest(namespace.ContextWithNamespace(ctx, ns), req)

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
			req.Operation == logical.CreateOperation) {
		return logical.ErrorResponse("cannot write to a path ending in '/'"), nil
	}

	waitGroup, err := waitForReplicationState(ctx, c, req)
	if err != nil {
		return nil, err
	}

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
				return nil, fmt.Errorf("error validating wrapping token: %w", err)
			}
			if !valid {
				return nil, consts.ErrInvalidWrappingToken
			}

		// The -self paths have no meaning outside of the token NS, so
		// requests for these paths always go to the token NS
		case "auth/token/lookup-self", "auth/token/renew-self", "auth/token/revoke-self":
			ctx = newCtx

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
			if !ok {
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
	if c.router.LoginPath(ctx, req.Path) {
		resp, auth, err = c.handleLoginRequest(ctx, req)
	} else {
		resp, auth, err = c.handleRequest(ctx, req)
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
		walState.ClusterID = c.clusterID.Load()
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
		return forward(ctx, c, req)
	}
	atomic.AddUint64(c.counters.requests, 1)
	return resp, err
}

func (c *Core) handleRequest(ctx context.Context, req *logical.Request) (retResp *logical.Response, retAuth *logical.Auth, retErr error) {
	defer metrics.MeasureSince([]string{"core", "handle_request"}, time.Now())

	var nonHMACReqDataKeys []string
	entry := c.router.MatchingMountEntry(ctx, req.Path)
	if entry != nil {
		// Set here so the audit log has it even if authorization fails
		req.MountType = entry.Type
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
	auth, te, ctErr := c.checkToken(ctx, req, false)
	if ctErr == logical.ErrPerfStandbyPleaseForward {
		return nil, nil, ctErr
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

	leaseGenerated := false
	quotaResp, quotaErr := c.applyLeaseCountQuota(&quotas.Request{
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

			leaseID, err := registerFunc(ctx, req, resp)
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

		_, identityPolicies, err := c.fetchEntityAndDerivedPolicies(ctx, tokenNS, resp.Auth.EntityID)
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
			resp.Auth.Policies = policyutil.SanitizePolicies(resp.Auth.TokenPolicies, policyutil.DoNotAddDefaultPolicy)
		} else {
			resp.Auth.TokenPolicies = policyutil.SanitizePolicies(resp.Auth.Policies, policyutil.DoNotAddDefaultPolicy)

			switch resp.Auth.TokenType {
			case logical.TokenTypeBatch:
			case logical.TokenTypeService:
				if err := c.expiration.RegisterAuth(ctx, &logical.TokenEntry{
					TTL:         auth.TTL,
					Policies:    auth.TokenPolicies,
					Path:        resp.Auth.CreationPath,
					NamespaceID: ns.ID,
				}, resp.Auth); err != nil {
					// Best-effort clean up on error, so we log the cleanup error as
					// a warning but still return as internal error.
					if err := c.tokenStore.revokeOrphan(ctx, resp.Auth.ClientToken); err != nil {
						c.logger.Warn("failed to clean up token lease during auth/token/ request", "request_path", req.Path, "error", err)
					}
					c.logger.Error("failed to register token lease during auth/token/ request", "request_path", req.Path, "error", err)
					retErr = multierror.Append(retErr, ErrInternalError)
					return nil, auth, retErr
				}
				leaseGenerated = true
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
		// Get and set ignored HMAC'd value.
		if rawVals, ok := entry.synthesizedConfigCache.Load("audit_non_hmac_request_keys"); ok {
			nonHMACReqDataKeys = rawVals.([]string)
		}
	}

	// Do an unauth check. This will cause EGP policies to be checked
	var auth *logical.Auth
	var ctErr error
	auth, _, ctErr = c.checkToken(ctx, req, true)
	if ctErr == logical.ErrPerfStandbyPleaseForward {
		return nil, nil, ctErr
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

	// If the response generated an authentication, then generate the token
	if resp != nil && resp.Auth != nil {
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			c.logger.Error("failed to get namespace from context", "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
			return
		}

		leaseGenerated := false

		// The request successfully authenticated itself. Run the quota checks
		// before creating lease.
		quotaResp, quotaErr := c.applyLeaseCountQuota(&quotas.Request{
			Path:          req.Path,
			MountPath:     strings.TrimPrefix(req.MountPoint, ns.Path),
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
			!mEntry.Local &&
			c.identityStore != nil {
			// Overwrite the mount type and mount path in the alias
			// information
			auth.Alias.MountType = req.MountType
			auth.Alias.MountAccessor = req.MountAccessor

			if auth.Alias.Name == "" {
				return nil, nil, fmt.Errorf("missing name in alias")
			}

			var err error

			// Fetch the entity for the alias, or create an entity if one
			// doesn't exist.
			entity, err = c.identityStore.CreateOrFetchEntity(ctx, auth.Alias)
			if err != nil {
				entity, err = possiblyForwardAliasCreation(ctx, c, err, auth, entity)
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
			validAliases, err := c.identityStore.refreshExternalGroupMembershipsByEntityID(ctx, auth.EntityID, auth.GroupAliases, req.MountAccessor)
			if err != nil {
				return nil, nil, err
			}
			auth.GroupAliases = validAliases
		}

		// Determine the source of the login
		source := c.router.MatchingMount(ctx, req.Path)
		source = strings.TrimPrefix(source, credentialRoutePrefix)
		source = strings.Replace(source, "/", "-", -1)

		// Prepend the source to the display name
		auth.DisplayName = strings.TrimSuffix(source+auth.DisplayName, "-")

		sysView := c.router.MatchingSystemView(ctx, req.Path)
		if sysView == nil {
			c.logger.Error("unable to look up sys view for login path", "request_path", req.Path)
			return nil, nil, ErrInternalError
		}

		tokenTTL, warnings, err := framework.CalculateTTL(sysView, 0, auth.TTL, auth.Period, auth.MaxTTL, auth.ExplicitMaxTTL, time.Time{})
		if err != nil {
			return nil, nil, err
		}
		for _, warning := range warnings {
			resp.AddWarning(warning)
		}

		_, identityPolicies, err := c.fetchEntityAndDerivedPolicies(ctx, ns, auth.EntityID)
		if err != nil {
			return nil, nil, ErrInternalError
		}

		auth.TokenPolicies = policyutil.SanitizePolicies(auth.Policies, !auth.NoDefaultPolicy)
		allPolicies := policyutil.SanitizePolicies(append(auth.TokenPolicies, identityPolicies[ns.ID]...), policyutil.DoNotAddDefaultPolicy)

		// Prevent internal policies from being assigned to tokens. We check
		// this on auth.Policies including derived ones from Identity before
		// actually making the token.
		for _, policy := range allPolicies {
			if policy == "root" {
				return logical.ErrorResponse("auth methods cannot create root tokens"), nil, logical.ErrInvalidRequest
			}
			if strutil.StrListContains(nonAssignablePolicies, policy) {
				return logical.ErrorResponse(fmt.Sprintf("cannot assign policy %q", policy)), nil, logical.ErrInvalidRequest
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
			retErr = multierror.Append(retErr, funcGetErr)
			return nil, auth, retErr
		}

		err = registerFunc(ctx, tokenTTL, req.Path, auth)
		switch {
		case err == nil:
			if auth.TokenType != logical.TokenTypeBatch {
				leaseGenerated = true
			}
		case err == ErrInternalError:
			return nil, auth, err
		default:
			return logical.ErrorResponse(err.Error()), auth, logical.ErrInvalidRequest
		}

		auth.IdentityPolicies = policyutil.SanitizePolicies(identityPolicies[ns.ID], policyutil.DoNotAddDefaultPolicy)
		delete(identityPolicies, ns.ID)
		auth.ExternalNamespacePolicies = identityPolicies
		auth.Policies = allPolicies

		// Attach the display name, might be used by audit backends
		req.DisplayName = auth.DisplayName

		// Count the successful token creation
		ttl_label := metricsutil.TTLBucket(tokenTTL)
		// Do not include namespace path in mount point; already present as separate label.
		mountPointWithoutNs := ns.TrimmedPath(req.MountPoint)
		c.metricSink.IncrCounterWithLabels(
			[]string{"token", "creation"},
			1,
			[]metrics.Label{
				metricsutil.NamespaceLabel(ns),
				{"auth_method", req.MountType},
				{"mount_point", mountPointWithoutNs},
				{"creation_ttl", ttl_label},
				{"token_type", auth.TokenType.String()},
			},
		)
	}

	return resp, auth, routeErr
}

// RegisterAuth uses a logical.Auth object to create a token entry in the token
// store, and registers a corresponding token lease to the expiration manager.
func (c *Core) RegisterAuth(ctx context.Context, tokenTTL time.Duration, path string, auth *logical.Auth) error {
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
		if err := c.expiration.RegisterAuth(ctx, &te, auth); err != nil {
			if err := c.tokenStore.revokeOrphan(ctx, te.ID); err != nil {
				c.logger.Warn("failed to clean up token lease during login request", "request_path", path, "error", err)
			}
			c.logger.Error("failed to register token lease during login request", "request_path", path, "error", err)
			return ErrInternalError
		}
	}

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
	te, err := c.LookupToken(ctx, token)
	if err != nil {
		dotCount := strings.Count(token, ".")
		// If we have two dots but the second char is a dot it's a vault
		// token of the form s.SOMETHING.nsid, not a JWT
		if dotCount != 2 || (dotCount == 2 && token[1] == '.') {
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
