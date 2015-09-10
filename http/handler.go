package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

// AuthHeaderName is the name of the header containing the token.
const AuthHeaderName = "X-Vault-Token"

// Handler returns an http.Handler for the API. This can be used on
// its own to mount the Vault API within another web server.
func Handler(core *vault.Core) http.Handler {
	// Create the muxer to handle the actual endpoints
	mux := http.NewServeMux()
	mux.Handle("/v1/sys/init", handleSysInit(core))
	mux.Handle("/v1/sys/seal-status", handleSysSealStatus(core))
	mux.Handle("/v1/sys/seal", handleSysSeal(core))
	mux.Handle("/v1/sys/unseal", handleSysUnseal(core))
	mux.Handle("/v1/sys/mounts", proxySysRequest(core))
	mux.Handle("/v1/sys/mounts/", proxySysRequest(core))
	mux.Handle("/v1/sys/remount", proxySysRequest(core))
	mux.Handle("/v1/sys/policy", handleSysListPolicies(core))
	mux.Handle("/v1/sys/policy/", handleSysPolicy(core))
	mux.Handle("/v1/sys/renew/", proxySysRequest(core))
	mux.Handle("/v1/sys/revoke/", handleSysRevoke(core))
	mux.Handle("/v1/sys/revoke-prefix/", handleSysRevokePrefix(core))
	mux.Handle("/v1/sys/auth", proxySysRequest(core))
	mux.Handle("/v1/sys/auth/", proxySysRequest(core))
	mux.Handle("/v1/sys/audit", handleSysListAudit(core))
	mux.Handle("/v1/sys/audit/", handleSysAudit(core))
	mux.Handle("/v1/sys/leader", handleSysLeader(core))
	mux.Handle("/v1/sys/health", handleSysHealth(core))
	mux.Handle("/v1/sys/rotate", proxySysRequest(core))
	mux.Handle("/v1/sys/key-status", proxySysRequest(core))
	mux.Handle("/v1/sys/rekey/init", handleSysRekeyInit(core))
	mux.Handle("/v1/sys/rekey/update", handleSysRekeyUpdate(core))
	mux.Handle("/v1/", handleLogical(core, false))

	// Wrap the handler in another handler to trigger all help paths.
	handler := handleHelpHandler(mux, core)

	return handler
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

func parseRequest(r *http.Request, out interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(out)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Failed to parse JSON input: %s", err)
	}
	return err
}

// request is a helper to perform a request and properly exit in the
// case of an error.
func request(core *vault.Core, w http.ResponseWriter, rawReq *http.Request, r *logical.Request) (*logical.Response, bool) {
	resp, err := core.HandleRequest(r)
	if err == vault.ErrStandby {
		respondStandby(core, w, rawReq.URL)
		return resp, false
	}
	if respondCommon(w, resp, err) {
		return resp, false
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return resp, false
	}

	return resp, true
}

// respondStandby is used to trigger a redirect in the case that this Vault is currently a hot standby
func respondStandby(core *vault.Core, w http.ResponseWriter, reqURL *url.URL) {
	// Request the leader address
	_, advertise, err := core.Leader()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// If there is no leader, generate a 503 error
	if advertise == "" {
		err = fmt.Errorf("no active Vault instance found")
		respondError(w, http.StatusServiceUnavailable, err)
		return
	}

	// Parse the advertise location
	advertiseURL, err := url.Parse(advertise)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Generate a redirect URL
	redirectURL := url.URL{
		Scheme:   advertiseURL.Scheme,
		Host:     advertiseURL.Host,
		Path:     reqURL.Path,
		RawQuery: reqURL.RawQuery,
	}

	// Ensure there is a scheme, default to https
	if redirectURL.Scheme == "" {
		redirectURL.Scheme = "https"
	}

	// If we have an address, redirect! We use a 307 code
	// because we don't actually know if its permanent and
	// the request method should be preserved.
	w.Header().Set("Location", redirectURL.String())
	w.WriteHeader(307)
}

// requestAuth adds the token to the logical.Request if it exists.
func requestAuth(r *http.Request, req *logical.Request) *logical.Request {
	// Attach the header value if we have it
	if v := r.Header.Get(AuthHeaderName); v != "" {
		req.ClientToken = v
	}

	return req
}

func respondError(w http.ResponseWriter, status int, err error) {
	// Adjust status code when sealed
	if err == vault.ErrSealed {
		status = http.StatusServiceUnavailable
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

func respondCommon(w http.ResponseWriter, resp *logical.Response, err error) bool {
	if resp == nil {
		return false
	}

	if resp.IsError() {
		var statusCode int

		switch err {
		case logical.ErrPermissionDenied:
			statusCode = http.StatusForbidden
		case logical.ErrUnsupportedOperation:
			statusCode = http.StatusMethodNotAllowed
		case logical.ErrUnsupportedPath:
			statusCode = http.StatusNotFound
		case logical.ErrInvalidRequest:
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusBadRequest
		}

		err := fmt.Errorf("%s", resp.Data["error"].(string))
		respondError(w, statusCode, err)
		return true
	}

	return false
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

func proxySysRequest(core *vault.Core) http.Handler {
	return handleLogical(core, true)
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}
