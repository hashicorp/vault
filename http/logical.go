package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func buildLogicalRequest(core *vault.Core, w http.ResponseWriter, r *http.Request) (*logical.Request, int, error) {
	ns, err := namespace.FromContext(r.Context())
	if err != nil {
		return nil, http.StatusBadRequest, nil
	}
	path := ns.TrimmedPath(r.URL.Path[len("/v1/"):])

	var data map[string]interface{}

	// Determine the operation
	var op logical.Operation
	switch r.Method {
	case "DELETE":
		op = logical.DeleteOperation

	case "GET":
		op = logical.ReadOperation
		queryVals := r.URL.Query()
		var list bool
		var err error
		listStr := queryVals.Get("list")
		if listStr != "" {
			list, err = strconv.ParseBool(listStr)
			if err != nil {
				return nil, http.StatusBadRequest, nil
			}
			if list {
				op = logical.ListOperation
				if !strings.HasSuffix(path, "/") {
					path += "/"
				}
			}
		}

		if !list {
			getData := map[string]interface{}{}

			for k, v := range r.URL.Query() {
				// Skip the help key as this is a reserved parameter
				if k == "help" {
					continue
				}

				switch {
				case len(v) == 0:
				case len(v) == 1:
					getData[k] = v[0]
				default:
					getData[k] = v
				}
			}

			if len(getData) > 0 {
				data = getData
			}
		}

	case "POST", "PUT":
		op = logical.UpdateOperation
		// Parse the request if we can
		if op == logical.UpdateOperation {
			err := parseRequest(r, w, &data)
			if err == io.EOF {
				data = nil
				err = nil
			}
			if err != nil {
				return nil, http.StatusBadRequest, err
			}
		}

	case "LIST":
		op = logical.ListOperation
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}

	case "OPTIONS":
	default:
		return nil, http.StatusMethodNotAllowed, nil
	}

	request_id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, http.StatusBadRequest, errwrap.Wrapf("failed to generate identifier for the request: {{err}}", err)
	}

	req, err := requestAuth(core, r, &logical.Request{
		ID:         request_id,
		Operation:  op,
		Path:       path,
		Data:       data,
		Connection: getConnection(r),
		Headers:    r.Header,
	})
	if err != nil {
		if errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
			return nil, http.StatusForbidden, nil
		}
		return nil, http.StatusBadRequest, errwrap.Wrapf("error performing token check: {{err}}", err)
	}

	req, err = requestWrapInfo(r, req)
	if err != nil {
		return nil, http.StatusBadRequest, errwrap.Wrapf("error parsing X-Vault-Wrap-TTL header: {{err}}", err)
	}

	err = parseMFAHeader(req)
	if err != nil {
		return nil, http.StatusBadRequest, errwrap.Wrapf("failed to parse X-Vault-MFA header: {{err}}", err)
	}

	err = requestPolicyOverride(r, req)
	if err != nil {
		return nil, http.StatusBadRequest, errwrap.Wrapf(fmt.Sprintf(`failed to parse %s header: {{err}}`, PolicyOverrideHeaderName), err)
	}

	return req, 0, nil
}

func handleLogical(core *vault.Core) http.Handler {
	return handleLogicalInternal(core, false)
}

func handleLogicalWithInjector(core *vault.Core) http.Handler {
	return handleLogicalInternal(core, true)
}

func handleLogicalInternal(core *vault.Core, injectDataIntoTopLevel bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, statusCode, err := buildLogicalRequest(core, w, r)
		if err != nil || statusCode != 0 {
			respondError(w, statusCode, err)
			return
		}

		// Always forward requests that are using a limited use count token
		if core.PerfStandby() && req.ClientTokenRemainingUses > 0 {
			forwardRequest(core, w, r)
			return
		}

		// req.Path will be relative by this point. The prefix check is first
		// to fail faster if we're not in this situation since it's a hot path
		switch {
		case strings.HasPrefix(req.Path, "sys/wrapping/"), strings.HasPrefix(req.Path, "auth/token/"):
			// Get the token ns info; if we match the paths below we want to
			// swap in the token context (but keep the relative path)
			if err != nil {
				core.Logger().Warn("error looking up just-set context", "error", err)
				respondError(w, http.StatusInternalServerError, err)
				return
			}
			te := req.TokenEntry()
			newCtx := r.Context()
			if te != nil {
				ns, err := vault.NamespaceByID(newCtx, te.NamespaceID, core)
				if err != nil {
					core.Logger().Warn("error looking up namespace from the token's namespace ID", "error", err)
					respondError(w, http.StatusInternalServerError, err)
					return
				}
				if ns != nil {
					newCtx = namespace.ContextWithNamespace(newCtx, ns)
				}
			}
			switch req.Path {
			case "sys/wrapping/lookup", "sys/wrapping/rewrap", "sys/wrapping/unwrap":
				r = r.WithContext(newCtx)
				if err := wrappingVerificationFunc(r.Context(), core, req); err != nil {
					if errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
						respondError(w, http.StatusForbidden, err)
					} else {
						respondError(w, http.StatusBadRequest, err)
					}
					return
				}

			// The -self paths have no meaning outside of the token NS, so
			// requests for these paths always go to the token NS
			case "auth/token/lookup-self", "auth/token/renew-self", "auth/token/revoke-self":
				r = r.WithContext(newCtx)

			// For the following operations, we can set the proper namespace context
			// using the token's embedded nsID if a relative path was provided. Since
			// this is done at the HTTP layer, the operation will still be gated by
			// ACLs.
			case "auth/token/lookup", "auth/token/renew", "auth/token/revoke", "auth/token/revoke-orphan":
				token, ok := req.Data["token"]
				// If the token is not present (e.g. a bad request), break out and let the backend
				// handle the error
				if !ok {
					// If this is a token lookup request and if the token is not
					// explicitly provided, it will use the client token so we simply set
					// the context to the client token's context.
					if req.Path == "auth/token/lookup" {
						r = r.WithContext(newCtx)
					}
					break
				}
				_, nsID := namespace.SplitIDFromString(token.(string))
				if nsID != "" {
					ns, err := vault.NamespaceByID(newCtx, nsID, core)
					if err != nil {
						core.Logger().Warn("error looking up namespace from the token's namespace ID", "error", err)
						respondError(w, http.StatusInternalServerError, err)
						return
					}
					if ns != nil {
						newCtx = namespace.ContextWithNamespace(newCtx, ns)
						r = r.WithContext(newCtx)
					}
				}
			}

		// The following relative sys/leases/ paths handles re-routing requests
		// to the proper namespace using the lease ID on applicable paths.
		case strings.HasPrefix(req.Path, "sys/leases/"):
			switch req.Path {
			// For the following operations, we can set the proper namespace context
			// using the lease's embedded nsID if a relative path was provided. Since
			// this is done at the HTTP layer, the operation will still be gated by
			// ACLs.
			case "sys/leases/lookup", "sys/leases/renew", "sys/leases/revoke", "sys/leases/revoke-force":
				leaseID, ok := req.Data["lease_id"]
				// If lease ID is not present, break out and let the backend handle the error
				if !ok {
					break
				}
				_, nsID := namespace.SplitIDFromString(leaseID.(string))
				if nsID != "" {
					newCtx := r.Context()
					ns, err := vault.NamespaceByID(newCtx, nsID, core)
					if err != nil {
						core.Logger().Warn("error looking up namespace from the lease's namespace ID", "error", err)
						respondError(w, http.StatusInternalServerError, err)
						return
					}
					if ns != nil {
						newCtx = namespace.ContextWithNamespace(newCtx, ns)
						r = r.WithContext(newCtx)
					}
				}
			}
		}

		// Make the internal request. We attach the connection info
		// as well in case this is an authentication request that requires
		// it. Vault core handles stripping this if we need to. This also
		// handles all error cases; if we hit respondLogical, the request is a
		// success.
		resp, ok, needsForward := request(core, w, r, req)
		if needsForward {
			forwardRequest(core, w, r)
			return
		}
		if !ok {
			return
		}

		// Build the proper response
		respondLogical(w, r, req, resp, injectDataIntoTopLevel)
	})
}

func respondLogical(w http.ResponseWriter, r *http.Request, req *logical.Request, resp *logical.Response, injectDataIntoTopLevel bool) {
	var httpResp *logical.HTTPResponse
	var ret interface{}

	if resp != nil {
		if resp.Redirect != "" {
			// If we have a redirect, redirect! We use a 307 code
			// because we don't actually know if its permanent.
			http.Redirect(w, r, resp.Redirect, 307)
			return
		}

		// Check if this is a raw response
		if _, ok := resp.Data[logical.HTTPStatusCode]; ok {
			respondRaw(w, r, resp)
			return
		}

		if resp.WrapInfo != nil && resp.WrapInfo.Token != "" {
			httpResp = &logical.HTTPResponse{
				WrapInfo: &logical.HTTPWrapInfo{
					Token:           resp.WrapInfo.Token,
					Accessor:        resp.WrapInfo.Accessor,
					TTL:             int(resp.WrapInfo.TTL.Seconds()),
					CreationTime:    resp.WrapInfo.CreationTime.Format(time.RFC3339Nano),
					CreationPath:    resp.WrapInfo.CreationPath,
					WrappedAccessor: resp.WrapInfo.WrappedAccessor,
				},
			}
		} else {
			httpResp = logical.LogicalResponseToHTTPResponse(resp)
			httpResp.RequestID = req.ID
		}

		ret = httpResp

		if injectDataIntoTopLevel {
			injector := logical.HTTPSysInjector{
				Response: httpResp,
			}
			ret = injector
		}
	}

	// Respond
	respondOk(w, ret)
	return
}

// respondRaw is used when the response is using HTTPContentType and HTTPRawBody
// to change the default response handling. This is only used for specific things like
// returning the CRL information on the PKI backends.
func respondRaw(w http.ResponseWriter, r *http.Request, resp *logical.Response) {
	retErr := func(w http.ResponseWriter, err string) {
		w.Header().Set("X-Vault-Raw-Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}

	// Ensure this is never a secret or auth response
	if resp.Secret != nil || resp.Auth != nil {
		retErr(w, "raw responses cannot contain secrets or auth")
		return
	}

	// Get the status code
	statusRaw, ok := resp.Data[logical.HTTPStatusCode]
	if !ok {
		retErr(w, "no status code given")
		return
	}

	var status int
	switch statusRaw.(type) {
	case int:
		status = statusRaw.(int)
	case float64:
		status = int(statusRaw.(float64))
	case json.Number:
		s64, err := statusRaw.(json.Number).Float64()
		if err != nil {
			retErr(w, "cannot decode status code")
			return
		}
		status = int(s64)
	default:
		retErr(w, "cannot decode status code")
		return
	}

	nonEmpty := status != http.StatusNoContent

	var contentType string
	var body []byte

	// Get the content type header; don't require it if the body is empty
	contentTypeRaw, ok := resp.Data[logical.HTTPContentType]
	if !ok && nonEmpty {
		retErr(w, "no content type given")
		return
	}
	if ok {
		contentType, ok = contentTypeRaw.(string)
		if !ok {
			retErr(w, "cannot decode content type")
			return
		}
	}

	if nonEmpty {
		// Get the body
		bodyRaw, ok := resp.Data[logical.HTTPRawBody]
		if !ok {
			goto WRITE_RESPONSE
		}

		switch bodyRaw.(type) {
		case string:
			// This is best effort. The value may already be base64-decoded so
			// if it doesn't work we just use as-is
			bodyDec, err := base64.StdEncoding.DecodeString(bodyRaw.(string))
			if err == nil {
				body = bodyDec
			} else {
				body = []byte(bodyRaw.(string))
			}
		case []byte:
			body = bodyRaw.([]byte)
		default:
			retErr(w, "cannot decode body")
			return
		}
	}

WRITE_RESPONSE:
	// Write the response
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	w.WriteHeader(status)
	w.Write(body)
}

// getConnection is used to format the connection information for
// attaching to a logical request
func getConnection(r *http.Request) (connection *logical.Connection) {
	var remoteAddr string

	remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteAddr = ""
	}

	connection = &logical.Connection{
		RemoteAddr: remoteAddr,
		ConnState:  r.TLS,
	}
	return
}
