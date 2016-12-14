package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/duration"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	// AuthHeaderName is the name of the header containing the token.
	AuthHeaderName = "X-Vault-Token"

	// WrapHeaderName is the name of the header containing a directive to wrap the
	// response.
	WrapTTLHeaderName = "X-Vault-Wrap-TTL"

	// NoRequestForwardingHeaderName is the name of the header telling Vault
	// not to use request forwarding
	NoRequestForwardingHeaderName = "X-Vault-No-Request-Forwarding"

	// MaxRequestSize is the maximum accepted request size. This is to prevent
	// a denial of service attack where no Content-Length is provided and the server
	// is fed ever more data until it exhausts memory.
	MaxRequestSize = 32 * 1024 * 1024
)

// Handler returns an http.Handler for the API. This can be used on
// its own to mount the Vault API within another web server.
func Handler(core *vault.Core) http.Handler {
	// Create the muxer to handle the actual endpoints
	mux := http.NewServeMux()
	mux.Handle("/v1/sys/init", handleSysInit(core))
	mux.Handle("/v1/sys/seal-status", handleSysSealStatus(core))
	mux.Handle("/v1/sys/seal", handleSysSeal(core))
	mux.Handle("/v1/sys/step-down", handleSysStepDown(core))
	mux.Handle("/v1/sys/unseal", handleSysUnseal(core))
	mux.Handle("/v1/sys/renew", handleRequestForwarding(core, handleLogical(core, false, nil)))
	mux.Handle("/v1/sys/renew/", handleRequestForwarding(core, handleLogical(core, false, nil)))
	mux.Handle("/v1/sys/leader", handleSysLeader(core))
	mux.Handle("/v1/sys/health", handleSysHealth(core))
	mux.Handle("/v1/sys/generate-root/attempt", handleRequestForwarding(core, handleSysGenerateRootAttempt(core)))
	mux.Handle("/v1/sys/generate-root/update", handleRequestForwarding(core, handleSysGenerateRootUpdate(core)))
	mux.Handle("/v1/sys/rekey/init", handleRequestForwarding(core, handleSysRekeyInit(core, false)))
	mux.Handle("/v1/sys/rekey/update", handleRequestForwarding(core, handleSysRekeyUpdate(core, false)))
	mux.Handle("/v1/sys/rekey-recovery-key/init", handleRequestForwarding(core, handleSysRekeyInit(core, true)))
	mux.Handle("/v1/sys/rekey-recovery-key/update", handleRequestForwarding(core, handleSysRekeyUpdate(core, true)))
	mux.Handle("/v1/sys/wrapping/lookup", handleRequestForwarding(core, handleLogical(core, false, wrappingVerificationFunc)))
	mux.Handle("/v1/sys/wrapping/rewrap", handleRequestForwarding(core, handleLogical(core, false, wrappingVerificationFunc)))
	mux.Handle("/v1/sys/wrapping/unwrap", handleRequestForwarding(core, handleLogical(core, false, wrappingVerificationFunc)))
	mux.Handle("/v1/sys/capabilities-self", handleRequestForwarding(core, handleLogical(core, true, nil)))
	mux.Handle("/v1/sys/", handleRequestForwarding(core, handleLogical(core, true, nil)))
	mux.Handle("/v1/", handleRequestForwarding(core, handleLogical(core, false, nil)))

	// Wrap the handler in another handler to trigger all help paths.
	helpWrappedHandler := wrapHelpHandler(mux, core)

	// Wrap the help wrapped handler with another layer with a generic
	// handler
	genericWrappedHandler := wrapGenericHandler(helpWrappedHandler)

	return genericWrappedHandler
}

// wrapGenericHandler wraps the handler with an extra layer of handler where
// tasks that should be commonly handled for all the requests and/or responses
// are performed.
func wrapGenericHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Cache-Control header for all the responses returned
		// by Vault
		w.Header().Set("Cache-Control", "no-store")
		h.ServeHTTP(w, r)
		return
	})
}

// A lookup on a token that is about to expire returns nil, which means by the
// time we can validate a wrapping token lookup will return nil since it will
// be revoked after the call. So we have to do the validation here.
func wrappingVerificationFunc(core *vault.Core, req *logical.Request) error {
	if req == nil {
		return fmt.Errorf("invalid request")
	}

	var token string
	if req.Data != nil && req.Data["token"] != nil {
		if tokenStr, ok := req.Data["token"].(string); !ok {
			return fmt.Errorf("could not decode token in request body")
		} else if tokenStr == "" {
			return fmt.Errorf("empty token in request body")
		} else {
			token = tokenStr
		}
	} else {
		token = req.ClientToken
	}

	valid, err := core.ValidateWrappingToken(token)
	if err != nil {
		return fmt.Errorf("error validating wrapping token: %v", err)
	}
	if !valid {
		return fmt.Errorf("wrapping token is not valid or does not exist")
	}

	return nil
}

// stripPrefix is a helper to strip a prefix from the path. It will
// return false from the second return value if it the prefix doesn't exist.
func stripPrefix(prefix, path string) (string, bool) {
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}

	path = path[len(prefix):]
	if path == "" {
		return "", false
	}

	return path, true
}

func parseRequest(r *http.Request, w http.ResponseWriter, out interface{}) error {
	// Limit the maximum number of bytes to MaxRequestSize to protect
	// against an indefinite amount of data being read.
	limit := http.MaxBytesReader(w, r.Body, MaxRequestSize)
	err := jsonutil.DecodeJSONFromReader(limit, out)
	if err != nil && err != io.EOF {
		return errwrap.Wrapf("failed to parse JSON input: {{err}}", err)
	}
	return err
}

// handleRequestForwarding determines whether to forward a request or not,
// falling back on the older behavior of redirecting the client
func handleRequestForwarding(core *vault.Core, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(vault.IntNoForwardingHeaderName) != "" {
			handler.ServeHTTP(w, r)
			return
		}

		if r.Header.Get(NoRequestForwardingHeaderName) != "" {
			// Forwarding explicitly disabled, fall back to previous behavior
			core.Logger().Trace("http/handleRequestForwarding: forwarding disabled by client request")
			handler.ServeHTTP(w, r)
			return
		}

		// Note: in an HA setup, this call will also ensure that connections to
		// the leader are set up, as that happens once the advertised cluster
		// values are read during this function
		isLeader, leaderAddr, err := core.Leader()
		if err != nil {
			if err == vault.ErrHANotEnabled {
				// Standalone node, serve request normally
				handler.ServeHTTP(w, r)
				return
			}
			// Some internal error occurred
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if isLeader {
			// No forwarding needed, we're leader
			handler.ServeHTTP(w, r)
			return
		}
		if leaderAddr == "" {
			respondError(w, http.StatusInternalServerError, fmt.Errorf("node not active but active node not found"))
			return
		}

		// Attempt forwarding the request. If we cannot forward -- perhaps it's
		// been disabled on the active node -- this will return with an
		// ErrCannotForward and we simply fall back
		statusCode, header, retBytes, err := core.ForwardRequest(r)
		if err != nil {
			if err == vault.ErrCannotForward {
				core.Logger().Trace("http/handleRequestForwarding: cannot forward (possibly disabled on active node), falling back")
			} else {
				core.Logger().Error("http/handleRequestForwarding: error forwarding request", "error", err)
			}

			// Fall back to redirection
			handler.ServeHTTP(w, r)
			return
		}

		if header != nil {
			for k, v := range header {
				for _, j := range v {
					w.Header().Add(k, j)
				}
			}
		}

		w.WriteHeader(statusCode)
		w.Write(retBytes)
		return
	})
}

// request is a helper to perform a request and properly exit in the
// case of an error.
func request(core *vault.Core, w http.ResponseWriter, rawReq *http.Request, r *logical.Request) (*logical.Response, bool) {
	resp, err := core.HandleRequest(r)
	if errwrap.Contains(err, vault.ErrStandby.Error()) {
		respondStandby(core, w, rawReq.URL)
		return resp, false
	}
	if respondErrorCommon(w, resp, err) {
		return resp, false
	}

	return resp, true
}

// respondStandby is used to trigger a redirect in the case that this Vault is currently a hot standby
func respondStandby(core *vault.Core, w http.ResponseWriter, reqURL *url.URL) {
	// Request the leader address
	_, redirectAddr, err := core.Leader()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// If there is no leader, generate a 503 error
	if redirectAddr == "" {
		err = fmt.Errorf("no active Vault instance found")
		respondError(w, http.StatusServiceUnavailable, err)
		return
	}

	// Parse the redirect location
	redirectURL, err := url.Parse(redirectAddr)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Generate a redirect URL
	finalURL := url.URL{
		Scheme:   redirectURL.Scheme,
		Host:     redirectURL.Host,
		Path:     reqURL.Path,
		RawQuery: reqURL.RawQuery,
	}

	// Ensure there is a scheme, default to https
	if finalURL.Scheme == "" {
		finalURL.Scheme = "https"
	}

	// If we have an address, redirect! We use a 307 code
	// because we don't actually know if its permanent and
	// the request method should be preserved.
	w.Header().Set("Location", finalURL.String())
	w.WriteHeader(307)
}

// requestAuth adds the token to the logical.Request if it exists.
func requestAuth(core *vault.Core, r *http.Request, req *logical.Request) *logical.Request {
	// Attach the header value if we have it
	if v := r.Header.Get(AuthHeaderName); v != "" {
		req.ClientToken = v

		// Also attach the accessor if we have it. This doesn't fail if it
		// doesn't exist because the request may be to an unauthenticated
		// endpoint/login endpoint where a bad current token doesn't matter, or
		// a token from a Vault version pre-accessors.
		te, err := core.LookupToken(v)
		if err == nil && te != nil {
			req.ClientTokenAccessor = te.Accessor
		}
	}

	return req
}

// requestWrapTTL adds the WrapTTL value to the logical.Request if it
// exists.
func requestWrapTTL(r *http.Request, req *logical.Request) (*logical.Request, error) {
	// First try for the header value
	wrapTTL := r.Header.Get(WrapTTLHeaderName)
	if wrapTTL == "" {
		return req, nil
	}

	// If it has an allowed suffix parse as a duration string
	dur, err := duration.ParseDurationSecond(wrapTTL)
	if err != nil {
		return req, err
	}
	if int64(dur) < 0 {
		return req, fmt.Errorf("requested wrap ttl cannot be negative")
	}
	req.WrapTTL = dur

	return req, nil
}

func respondError(w http.ResponseWriter, status int, err error) {
	// Adjust status code when sealed
	if errwrap.Contains(err, vault.ErrSealed.Error()) {
		status = http.StatusServiceUnavailable
	}

	// Adjust status code on
	if errwrap.Contains(err, "http: request body too large") {
		status = http.StatusRequestEntityTooLarge
	}

	// Allow HTTPCoded error passthrough to specify a code
	if t, ok := err.(logical.HTTPCodedError); ok {
		status = t.Code()
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func respondErrorCommon(w http.ResponseWriter, resp *logical.Response, err error) bool {
	// If there are no errors return
	if err == nil && (resp == nil || !resp.IsError()) {
		return false
	}

	// Start out with internal server error since in most of these cases there
	// won't be a response so this won't be overridden
	statusCode := http.StatusInternalServerError
	// If we actually have a response, start out with bad request
	if resp != nil {
		statusCode = http.StatusBadRequest
	}

	// Now, check the error itself; if it has a specific logical error, set the
	// appropriate code
	if err != nil {
		switch {
		case errwrap.ContainsType(err, new(vault.StatusBadRequest)):
			statusCode = http.StatusBadRequest
		case errwrap.Contains(err, logical.ErrPermissionDenied.Error()):
			statusCode = http.StatusForbidden
		case errwrap.Contains(err, logical.ErrUnsupportedOperation.Error()):
			statusCode = http.StatusMethodNotAllowed
		case errwrap.Contains(err, logical.ErrUnsupportedPath.Error()):
			statusCode = http.StatusNotFound
		case errwrap.Contains(err, logical.ErrInvalidRequest.Error()):
			statusCode = http.StatusBadRequest
		}
	}

	if resp != nil && resp.IsError() {
		err = fmt.Errorf("%s", resp.Data["error"].(string))
	}

	respondError(w, statusCode, err)
	return true
}

func respondOk(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")

	if body == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(body)
	}
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}
