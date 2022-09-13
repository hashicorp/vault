package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

// This endpoint is needed to answer queries before Vault unseals
// or becomes the leader.
func handleSysLeader(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysLeaderGet(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysLeaderGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	resp, err := core.GetLeaderStatus()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, resp)
}
