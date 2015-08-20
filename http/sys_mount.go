package http

import (
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysRemount(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT", "POST":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Parse the request if we can
		var req RemountRequest
		if err := parseRequest(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		_, err := core.HandleRequest(requestAuth(r, &logical.Request{
			Operation:  logical.WriteOperation,
			Path:       "sys/remount",
			Connection: getConnection(r),
			Data: map[string]interface{}{
				"from": req.From,
				"to":   req.To,
			},
		}))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		respondOk(w, nil)
	})
}

type MountRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type RemountRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}
