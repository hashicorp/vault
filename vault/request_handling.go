package vault

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
)

var (
	// Value for memoizing whether cubbyhole is mounted, e.g. if we are in
	// normal operation and not test mode
	cubbyholeMounted bool

	// mutex to ensure the same
	cubbyholeMountedMutex sync.Mutex
)

// HandleRequest is used to handle a new incoming request
func (c *Core) HandleRequest(req *logical.Request) (resp *logical.Response, err error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.sealed {
		return nil, ErrSealed
	}
	if c.standby {
		return nil, ErrStandby
	}

	// Allowing writing to a path ending in / makes it extremely difficult to
	// understand user intent for the filesystem-like backends (generic,
	// cubbyhole) -- did they want a key named foo/ or did they want to write
	// to a directory foo/ with no (or forgotten) key, or...? It also affects
	// lookup, because paths ending in / are considered prefixes by some
	// backends. Basically, it's all just terrible, so don't allow it.
	if strings.HasSuffix(req.Path, "/") &&
		(req.Operation == logical.UpdateOperation ||
			req.Operation == logical.CreateOperation) {
		return logical.ErrorResponse("cannot write to a path ending in '/'"), nil
	}

	var auth *logical.Auth
	if c.router.LoginPath(req.Path) {
		resp, auth, err = c.handleLoginRequest(req)
	} else {
		resp, auth, err = c.handleRequest(req)
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

	// In order to wrap, we need cubbyhole to be mounted, so we ensure that
	// cubbyhole is actually mounted, as it may not be during tests. We memoize
	// a true response, since cubbyhole cannot be mounted or unmounted during
	// normal operation.
	if !cubbyholeMounted {
		cubbyholeMountedMutex.Lock()
		// Ensure it wasn't changed by another goroutine
		if !cubbyholeMounted {
			if c.router.MatchingMount("cubbyhole/") != "" {
				cubbyholeMounted = true
			}
		}
		cubbyholeMountedMutex.Unlock()
	}

	// We are wrapping if there is anything to wrap (not a nil response) and a
	// TTL was specified for the token, plus if cubbyhole is mounted (which
	// will be the case normally)
	wrapping := cubbyholeMounted && resp != nil && resp.WrapInfo.TTL != 0

	// If we are wrapping, the first part happens before auditing so that
	// resp.WrapInfo.Token can contain the HMAC'd wrapping token ID in the
	// audit logs, so that it can be determined from the audit logs whether the
	// token was ever actually used.
	if wrapping {
		// Create the wrapping token
		te := TokenEntry{
			Path:         req.Path,
			Policies:     []string{"cubbyhole-response-wrapping"},
			CreationTime: time.Now().Unix(),
			TTL:          resp.WrapInfo.TTL,
			NumUses:      1,
		}

		if err := c.tokenStore.create(&te); err != nil {
			c.logger.Printf("[ERR] core: failed to create wrapping token: %v", err)
			return nil, ErrInternalError
		}

		resp.WrapInfo.Token = te.ID

		httpResponse := logical.SanitizeResponse(resp)

		cubbyReq := &logical.Request{
			Operation:   logical.CreateOperation,
			Path:        "cubbyhole/response",
			ClientToken: te.ID,
			Data: map[string]interface{}{
				"response": httpResponse,
			},
		}

		_, err = c.router.Route(cubbyReq)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to store wrapped response information: %v", err)
			return nil, ErrInternalError
		}

		auth := &logical.Auth{
			ClientToken: te.ID,
			Policies:    []string{"cubbyhole-response-wrapping"},
			LeaseOptions: logical.LeaseOptions{
				TTL:       te.TTL,
				Renewable: false,
			},
		}

		// Register the wrapped token with the expiration manager
		if err := c.expiration.RegisterAuth(te.Path, auth); err != nil {
			c.logger.Printf("[ERR] core: failed to register cubbyhole wrapping token lease "+
				"(request path: %s): %v", req.Path, err)
			return nil, ErrInternalError
		}
	}

	// Create an audit trail of the response
	if err := c.auditBroker.LogResponse(auth, req, resp, err); err != nil {
		c.logger.Printf("[ERR] core: failed to audit response (request path: %s): %v",
			req.Path, err)
		return nil, ErrInternalError
	}

	// If we are wrapping, now is when we create a new response object with the
	// wrapped information, since the original response has been audit logged
	if wrapping {
		wrappingResp := &logical.Response{
			WrapInfo: resp.WrapInfo,
		}
		wrappingResp.CloneWarnings(resp)
		resp = wrappingResp
	}

	return
}

func (c *Core) handleRequest(req *logical.Request) (retResp *logical.Response, retAuth *logical.Auth, retErr error) {
	defer metrics.MeasureSince([]string{"core", "handle_request"}, time.Now())

	// Validate the token
	auth, te, ctErr := c.checkToken(req)
	// We run this logic first because we want to decrement the use count even in the case of an error
	if te != nil {
		// Attempt to use the token (decrement NumUses)
		var err error
		te, err = c.tokenStore.UseToken(te)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to use token: %v", err)
			return nil, nil, ErrInternalError
		}
		if te == nil {
			// Token has been revoked by this point
			return nil, nil, logical.ErrPermissionDenied
		}
		if te.NumUses == -1 {
			// We defer a revocation until after logic has run, since this is a
			// valid request (this is the token's final use). We pass the ID in
			// directly just to be safe in case something else modifies te later.
			defer func(id string) {
				err = c.tokenStore.Revoke(id)
				if err != nil {
					c.logger.Printf("[ERR] core: failed to revoke token: %v", err)
					retResp = nil
					retAuth = nil
					retErr = ErrInternalError
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
		// If it is an internal error we return that, otherwise we
		// return invalid request so that the status codes can be correct
		var errType error
		switch ctErr {
		case ErrInternalError, logical.ErrPermissionDenied:
			errType = ctErr
		default:
			errType = logical.ErrInvalidRequest
		}

		if err := c.auditBroker.LogRequest(auth, req, ctErr); err != nil {
			c.logger.Printf("[ERR] core: failed to audit request with path (%s): %v",
				req.Path, err)
		}

		return logical.ErrorResponse(ctErr.Error()), nil, errType
	}

	// Attach the display name
	req.DisplayName = auth.DisplayName

	// Create an audit trail of the request
	if err := c.auditBroker.LogRequest(auth, req, nil); err != nil {
		c.logger.Printf("[ERR] core: failed to audit request with path (%s): %v",
			req.Path, err)
		return nil, auth, ErrInternalError
	}

	// Route the request
	resp, err := c.router.Route(req)

	// If there is a secret, we must register it with the expiration manager.
	// We exclude renewal of a lease, since it does not need to be re-registered
	if resp != nil && resp.Secret != nil && !strings.HasPrefix(req.Path, "sys/renew/") {
		// Get the SystemView for the mount
		sysView := c.router.MatchingSystemView(req.Path)
		if sysView == nil {
			c.logger.Println("[ERR] core: unable to retrieve system view from router")
			return nil, auth, ErrInternalError
		}

		// Apply the default lease if none given
		if resp.Secret.TTL == 0 {
			resp.Secret.TTL = sysView.DefaultLeaseTTL()
		}

		// Limit the lease duration
		maxTTL := sysView.MaxLeaseTTL()
		if resp.Secret.TTL > maxTTL {
			resp.Secret.TTL = maxTTL
		}

		// Generic mounts should return the TTL but not register
		// for a lease as this provides a massive slowdown
		registerLease := true
		matchingBackend := c.router.MatchingBackend(req.Path)
		if matchingBackend == nil {
			c.logger.Println("[ERR] core: unable to retrieve generic backend from router")
			return nil, auth, ErrInternalError
		}
		if ptbe, ok := matchingBackend.(*PassthroughBackend); ok {
			if !ptbe.GeneratesLeases() {
				registerLease = false
				resp.Secret.Renewable = false
			}
		}

		if registerLease {
			leaseID, err := c.expiration.Register(req, resp)
			if err != nil {
				c.logger.Printf(
					"[ERR] core: failed to register lease "+
						"(request path: %s): %v", req.Path, err)
				return nil, auth, ErrInternalError
			}
			resp.Secret.LeaseID = leaseID
		}
	}

	// Only the token store is allowed to return an auth block, for any
	// other request this is an internal error. We exclude renewal of a token,
	// since it does not need to be re-registered
	if resp != nil && resp.Auth != nil && !strings.HasPrefix(req.Path, "auth/token/renew") {
		if !strings.HasPrefix(req.Path, "auth/token/") {
			c.logger.Printf(
				"[ERR] core: unexpected Auth response for non-token backend "+
					"(request path: %s)", req.Path)
			return nil, auth, ErrInternalError
		}

		// Register with the expiration manager. We use the token's actual path
		// here because roles allow suffixes.
		te, err := c.tokenStore.Lookup(resp.Auth.ClientToken)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to lookup token: %v", err)
			return nil, nil, ErrInternalError
		}

		if err := c.expiration.RegisterAuth(te.Path, resp.Auth); err != nil {
			c.logger.Printf("[ERR] core: failed to register token lease "+
				"(request path: %s): %v", req.Path, err)
			return nil, auth, ErrInternalError
		}
	}

	// Return the response and error
	return resp, auth, err
}

// handleLoginRequest is used to handle a login request, which is an
// unauthenticated request to the backend.
func (c *Core) handleLoginRequest(req *logical.Request) (*logical.Response, *logical.Auth, error) {
	defer metrics.MeasureSince([]string{"core", "handle_login_request"}, time.Now())

	// Create an audit trail of the request, auth is not available on login requests
	if err := c.auditBroker.LogRequest(nil, req, nil); err != nil {
		c.logger.Printf("[ERR] core: failed to audit request with path %s: %v",
			req.Path, err)
		return nil, nil, ErrInternalError
	}

	// Route the request
	resp, err := c.router.Route(req)

	// A login request should never return a secret!
	if resp != nil && resp.Secret != nil {
		c.logger.Printf("[ERR] core: unexpected Secret response for login path"+
			"(request path: %s)", req.Path)
		return nil, nil, ErrInternalError
	}

	// If the response generated an authentication, then generate the token
	var auth *logical.Auth
	if resp != nil && resp.Auth != nil {
		auth = resp.Auth

		// Determine the source of the login
		source := c.router.MatchingMount(req.Path)
		source = strings.TrimPrefix(source, credentialRoutePrefix)
		source = strings.Replace(source, "/", "-", -1)

		// Prepend the source to the display name
		auth.DisplayName = strings.TrimSuffix(source+auth.DisplayName, "-")

		sysView := c.router.MatchingSystemView(req.Path)
		if sysView == nil {
			c.logger.Printf("[ERR] core: unable to look up sys view for login path"+
				"(request path: %s)", req.Path)
			return nil, nil, ErrInternalError
		}

		// Set the default lease if non-provided, root tokens are exempt
		if auth.TTL == 0 && !strutil.StrListContains(auth.Policies, "root") {
			auth.TTL = sysView.DefaultLeaseTTL()
		}

		// Limit the lease duration
		if auth.TTL > sysView.MaxLeaseTTL() {
			auth.TTL = sysView.MaxLeaseTTL()
		}

		// Generate a token
		te := TokenEntry{
			Path:         req.Path,
			Policies:     auth.Policies,
			Meta:         auth.Metadata,
			DisplayName:  auth.DisplayName,
			CreationTime: time.Now().Unix(),
			TTL:          auth.TTL,
		}

		if strutil.StrListSubset(te.Policies, []string{"root"}) {
			te.Policies = []string{"root"}
		} else {
			// Use a map to filter out/prevent duplicates
			policyMap := map[string]bool{}
			for _, policy := range te.Policies {
				if policy == "" {
					// Don't allow a policy with no name, even though it is a valid
					// slice member
					continue
				}
				policyMap[policy] = true
			}

			// Add the default policy
			policyMap["default"] = true

			te.Policies = []string{}
			for k, _ := range policyMap {
				te.Policies = append(te.Policies, k)
			}

			sort.Strings(te.Policies)
		}

		if err := c.tokenStore.create(&te); err != nil {
			c.logger.Printf("[ERR] core: failed to create token: %v", err)
			return nil, auth, ErrInternalError
		}

		// Populate the client token and accessor
		auth.ClientToken = te.ID
		auth.Accessor = te.Accessor
		auth.Policies = te.Policies

		// Register with the expiration manager
		if err := c.expiration.RegisterAuth(te.Path, auth); err != nil {
			c.logger.Printf("[ERR] core: failed to register token lease "+
				"(request path: %s): %v", req.Path, err)
			return nil, auth, ErrInternalError
		}

		// Attach the display name, might be used by audit backends
		req.DisplayName = auth.DisplayName
	}

	return resp, auth, err
}
