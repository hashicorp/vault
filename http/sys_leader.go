package http

import (
	"github.com/hashicorp/vault/vault"
	"net/http"
)

// Can we remove this entirely, and only depend on the system backend?
// I suspect this has to remain, to answer queries while Core exists,
// but the system backend does not?
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
