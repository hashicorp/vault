package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

// AuthCookieName is the name of the cookie containing the token.
const AuthCookieName = "token"

// Handler returns an http.Handler for the API. This can be used on
// its own to mount the Vault API within another web server.
func Handler(core *vault.Core) http.Handler {
	// Create the muxer to handle the actual endpoints
	mux := http.NewServeMux()
	mux.Handle("/v1/sys/init", handleSysInit(core))
	mux.Handle("/v1/sys/seal-status", handleSysSealStatus(core))
	mux.Handle("/v1/sys/seal", handleSysSeal(core))
	mux.Handle("/v1/sys/unseal", handleSysUnseal(core))
	mux.Handle("/v1/sys/mounts/", handleSysMounts(core))
	mux.Handle("/v1/sys/remount", handleSysRemount(core))
	mux.Handle("/v1/sys/policy", handleSysListPolicies(core))
	mux.Handle("/v1/sys/policy/", handleSysPolicy(core))
	mux.Handle("/v1/sys/revoke/", handleSysRevoke(core))
	mux.Handle("/v1/sys/revoke-prefix/", handleSysRevokePrefix(core))
	mux.Handle("/v1/sys/auth/", handleSysAuth(core))
	mux.Handle("/v1/sys/audit", handleSysListAudit(core))
	mux.Handle("/v1/sys/audit/", handleSysAudit(core))
	mux.Handle("/v1/", handleLogical(core))

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
	return dec.Decode(out)
}

// request is a helper to perform a request and properly exit in the
// case of an error.
func request(core *vault.Core, w http.ResponseWriter, r *logical.Request) (*logical.Response, bool) {
	resp, err := core.HandleRequest(r)
	if respondCommon(w, resp) {
		return resp, false
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return resp, false
	}

	return resp, true
}

// requestAuth adds the token to the logical.Request if it exists.
func requestAuth(r *http.Request, req *logical.Request) *logical.Request {
	// Attach the cookie value as the token if we have it
	cookie, err := r.Cookie(AuthCookieName)
	if err == nil {
		req.ClientToken = cookie.Value
	}

	return req
}

func respondError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func respondCommon(w http.ResponseWriter, resp *logical.Response) bool {
	if resp == nil {
		return false
	}

	if resp.IsError() {
		err := fmt.Errorf("%s", resp.Data["error"].(string))
		respondError(w, http.StatusBadRequest, err)
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

type ErrorResponse struct {
	Errors []string `json:"errors"`
}
