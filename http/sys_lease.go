package http

import (
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysRenew(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Determine the path...
		prefix := "/v1/sys/renew/"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			respondError(w, http.StatusNotFound, nil)
			return
		}
		path := r.URL.Path[len(prefix):]
		if path == "" {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		// Parse the request if we can
		var req RenewRequest
		if err := parseRequest(r, &req); err != nil {
			if err != io.EOF {
				respondError(w, http.StatusBadRequest, err)
				return
			}
		}

		resp, ok := request(core, w, r, requestAuth(r, &logical.Request{
			Operation:  logical.WriteOperation,
			Path:       "sys/renew/" + path,
			Connection: getConnection(r),
			Data: map[string]interface{}{
				"increment": req.Increment,
			},
		}))
		if !ok {
			return
		}

		respondLogical(w, r, path, resp)
	})
}

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
			Operation:  logical.WriteOperation,
			Path:       "sys/revoke/" + path,
			Connection: getConnection(r),
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
			Operation:  logical.WriteOperation,
			Path:       "sys/revoke-prefix/" + path,
			Connection: getConnection(r),
		}))
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		respondOk(w, nil)
	})
}

type RenewRequest struct {
	Increment int `json:"increment"`
}
