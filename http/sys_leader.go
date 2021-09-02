package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

// This endpoint is needed to answer queries before Vault unseals
// or becomes the leader.
func handleSysLeader(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Getting custom headers from listener's config
		la := w.Header().Get("X-Vault-Listener-Add")
		lc, err := core.GetCustomResponseHeaders(la)
		if err != nil {
			core.Logger().Debug("failed to get custom headers from listener config")
		}
		switch r.Method {
		case "GET":
			handleSysLeaderGet(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil, lc)
		}
	})
}

func handleSysLeaderGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Getting custom headers from listener's config
	la := w.Header().Get("X-Vault-Listener-Add")
	lc, err := core.GetCustomResponseHeaders(la)
	if err != nil {
		core.Logger().Debug("failed to get custom headers from listener config")
	}
	resp, err := core.GetLeaderStatus()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err, lc)
		return
	}
	respondOk(w, resp, lc)
}
