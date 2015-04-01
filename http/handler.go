package http

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

// AuthCookieName is the name of the cookie containing the token.
const AuthCookieName = "token"

// Handler returns an http.Handler for the API. This can be used on
// its own to mount the Vault API within another web server.
func Handler(core *vault.Core) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/v1/sys/init", handleSysInit(core))
	mux.Handle("/v1/sys/seal-status", handleSysSealStatus(core))
	mux.Handle("/v1/sys/seal", handleSysSeal(core))
	mux.Handle("/v1/sys/unseal", handleSysUnseal(core))
	mux.Handle("/v1/sys/mounts/", handleSysMounts(core))
	mux.Handle("/v1/sys/revoke/", handleSysRevoke(core))
	mux.Handle("/v1/sys/revoke-prefix/", handleSysRevokePrefix(core))
	mux.Handle("/v1/", handleLogical(core))
	return mux
}

func parseRequest(r *http.Request, out interface{}) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(out)
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
