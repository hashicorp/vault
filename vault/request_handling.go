package vault

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
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

	// We are wrapping if there is anything to wrap (not a nil response) and a
	// TTL was specified for the token
	wrapping := resp != nil && resp.WrapInfo != nil && resp.WrapInfo.TTL != 0

	if wrapping {
		cubbyResp, err := c.wrapInCubbyhole(req, resp)
		// If not successful, returns either an error response from the
		// cubbyhole backend or an error; if either is set, return
		if cubbyResp != nil || err != nil {
			return cubbyResp, err
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
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, nil, retErr
		}
		if te == nil {
			// Token has been revoked by this point
			retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
			return nil, nil, retErr
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

		if errType != nil {
			retErr = multierror.Append(retErr, errType)
		}
		return logical.ErrorResponse(ctErr.Error()), nil, retErr
	}

	// Attach the display name
	req.DisplayName = auth.DisplayName

	// Create an audit trail of the request
	if err := c.auditBroker.LogRequest(auth, req, nil); err != nil {
		c.logger.Printf("[ERR] core: failed to audit request with path (%s): %v",
			req.Path, err)
		retErr = multierror.Append(retErr, ErrInternalError)
		return nil, auth, retErr
	}

	// Route the request
	resp, err := c.router.Route(req)
	if resp != nil {
		// We don't allow backends to specify this, so ensure it's not set
		resp.WrapInfo = nil

		if req.WrapTTL != 0 {
			resp.WrapInfo = &logical.WrapInfo{
				TTL: req.WrapTTL,
			}
		}
	}

	// If there is a secret, we must register it with the expiration manager.
	// We exclude renewal of a lease, since it does not need to be re-registered
	if resp != nil && resp.Secret != nil && !strings.HasPrefix(req.Path, "sys/renew/") {
		// Get the SystemView for the mount
		sysView := c.router.MatchingSystemView(req.Path)
		if sysView == nil {
			c.logger.Println("[ERR] core: unable to retrieve system view from router")
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, auth, retErr
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
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, auth, retErr
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
				retErr = multierror.Append(retErr, ErrInternalError)
				return nil, auth, retErr
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
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, auth, retErr
		}

		// Register with the expiration manager. We use the token's actual path
		// here because roles allow suffixes.
		te, err := c.tokenStore.Lookup(resp.Auth.ClientToken)
		if err != nil {
			c.logger.Printf("[ERR] core: failed to lookup token: %v", err)
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, nil, retErr
		}

		if err := c.expiration.RegisterAuth(te.Path, resp.Auth); err != nil {
			c.logger.Printf("[ERR] core: failed to register token lease "+
				"(request path: %s): %v", req.Path, err)
			retErr = multierror.Append(retErr, ErrInternalError)
			return nil, auth, retErr
		}
	}

	// Return the response and error
	if err != nil {
		retErr = multierror.Append(retErr, err)
	}
	return resp, auth, retErr
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
	if resp != nil {
		// We don't allow backends to specify this, so ensure it's not set
		resp.WrapInfo = nil

		if req.WrapTTL != 0 {
			resp.WrapInfo = &logical.WrapInfo{
				TTL: req.WrapTTL,
			}
		}
	}

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

func (c *Core) wrapInCubbyhole(req *logical.Request, resp *logical.Response) (*logical.Response, error) {
	// If we are wrapping, the first part (performed in this functions) happens
	// before auditing so that resp.WrapInfo.Token can contain the HMAC'd
	// wrapping token ID in the audit logs, so that it can be determined from
	// the audit logs whether the token was ever actually used.
	creationTime := time.Now()
	te := TokenEntry{
		Path:           req.Path,
		Policies:       []string{"response-wrapping"},
		CreationTime:   creationTime.Unix(),
		TTL:            resp.WrapInfo.TTL,
		NumUses:        1,
		ExplicitMaxTTL: resp.WrapInfo.TTL,
	}

	if err := c.tokenStore.create(&te); err != nil {
		c.logger.Printf("[ERR] core: failed to create wrapping token: %v", err)
		return nil, ErrInternalError
	}

	resp.WrapInfo.Token = te.ID
	resp.WrapInfo.CreationTime = creationTime

	// This will only be non-nil if this response contains a token, so in that
	// case put the accessor in the wrap info.
	if resp.Auth != nil {
		resp.WrapInfo.WrappedAccessor = resp.Auth.Accessor
	}

	httpResponse := logical.SanitizeResponse(resp)

	// Because of the way that JSON encodes (likely just in Go) we actually get
	// mixed-up values for ints if we simply put this object in the response
	// and encode the whole thing; so instead we marshal it first, then store
	// the string response. This actually ends up making it easier on the
	// client side, too, as it becomes a straight read-string-pass-to-unmarshal
	// operation.

	marshaledResponse, err := json.Marshal(httpResponse)
	if err != nil {
		c.logger.Printf("[ERR] core: failed to marshal wrapped response: %v", err)
		return nil, ErrInternalError
	}

	cubbyReq := &logical.Request{
		Operation:   logical.CreateOperation,
		Path:        "cubbyhole/response",
		ClientToken: te.ID,
		Data: map[string]interface{}{
			"response": string(marshaledResponse),
		},
	}

	cubbyResp, err := c.router.Route(cubbyReq)
	if err != nil {
		// Revoke since it's not yet being tracked for expiration
		c.tokenStore.Revoke(te.ID)
		c.logger.Printf("[ERR] core: failed to store wrapped response information: %v", err)
		return nil, ErrInternalError
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		c.tokenStore.Revoke(te.ID)
		c.logger.Printf("[ERR] core: failed to store wrapped response information: %v", cubbyResp.Data["error"])
		return cubbyResp, nil
	}

	auth := &logical.Auth{
		ClientToken: te.ID,
		Policies:    []string{"response-wrapping"},
		LeaseOptions: logical.LeaseOptions{
			TTL:       te.TTL,
			Renewable: false,
		},
	}

	// Register the wrapped token with the expiration manager
	if err := c.expiration.RegisterAuth(te.Path, auth); err != nil {
		// Revoke since it's not yet being tracked for expiration
		c.tokenStore.Revoke(te.ID)
		c.logger.Printf("[ERR] core: failed to register cubbyhole wrapping token lease "+
			"(request path: %s): %v", req.Path, err)
		return nil, ErrInternalError
	}

	return nil, nil
}
