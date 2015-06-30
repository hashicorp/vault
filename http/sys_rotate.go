package http

import (
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysKeyStatus(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		resp, err := core.HandleRequest(requestAuth(r, &logical.Request{
			Operation:  logical.ReadOperation,
			Path:       "sys/key-status",
			Connection: getConnection(r),
		}))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		respondOk(w, resp.Data)
	})
}

func handleSysRotate(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
		case "PUT":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		_, err := core.HandleRequest(requestAuth(r, &logical.Request{
			Operation:  logical.WriteOperation,
			Path:       "sys/rotate",
			Connection: getConnection(r),
		}))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		respondOk(w, nil)
	})
}
