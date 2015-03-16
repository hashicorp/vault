package http

import (
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysMounts(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		resp, err := core.HandleRequest(&logical.Request{
			Operation: logical.ReadOperation,
			Path:      "sys/mounts",
		})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		respondOk(w, resp.Data)
	})
}
