package http

import (
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysRevoke(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Determine the path...
		prefix := "/v1/sys/revoke/"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			respondError(w, http.StatusNotFound, nil)
			return
		}
		path := r.URL.Path[len(prefix):]
		if path == "" {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		_, err := core.HandleRequest(requestAuth(r, &logical.Request{
			Operation: logical.WriteOperation,
			Path:      "sys/revoke/" + path,
		}))
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		respondOk(w, nil)
	})
}

func handleSysRevokePrefix(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Determine the path...
		prefix := "/v1/sys/revoke-prefix/"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			respondError(w, http.StatusNotFound, nil)
			return
		}
		path := r.URL.Path[len(prefix):]
		if path == "" {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		_, err := core.HandleRequest(requestAuth(r, &logical.Request{
			Operation: logical.WriteOperation,
			Path:      "sys/revoke-prefix/" + path,
		}))
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		respondOk(w, nil)
	})
}
