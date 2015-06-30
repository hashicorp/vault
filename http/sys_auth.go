package http

import (
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysAuth(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysListAuth(core).ServeHTTP(w, r)
		case "POST", "PUT", "DELETE":
			handleSysEnableDisableAuth(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
	})
}

func handleSysListAuth(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		resp, err := core.HandleRequest(requestAuth(r, &logical.Request{
			Operation:  logical.ReadOperation,
			Path:       "sys/auth",
			Connection: getConnection(r),
		}))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		respondOk(w, resp.Data)
	})
}

func handleSysEnableDisableAuth(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Determine the path...
	prefix := "/v1/sys/auth/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		respondError(w, http.StatusNotFound, nil)
		return
	}
	path := r.URL.Path[len(prefix):]
	if path == "" {
		respondError(w, http.StatusNotFound, nil)
		return
	}

	switch r.Method {
	case "PUT", "POST":
		handleSysEnableAuth(core, w, r, path)
	case "DELETE":
		handleSysDisableAuth(core, w, r, path)
	default:
		panic("should never happen")
	}
}

func handleSysEnableAuth(
	core *vault.Core,
	w http.ResponseWriter,
	r *http.Request,
	path string) {
	// Parse the request if we can
	var req EnableAuthRequest
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	_, err := core.HandleRequest(requestAuth(r, &logical.Request{
		Operation:  logical.WriteOperation,
		Path:       "sys/auth/" + path,
		Connection: getConnection(r),
		Data: map[string]interface{}{
			"type":        req.Type,
			"description": req.Description,
		},
	}))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}

func handleSysDisableAuth(
	core *vault.Core,
	w http.ResponseWriter,
	r *http.Request,
	path string) {
	_, err := core.HandleRequest(requestAuth(r, &logical.Request{
		Operation:  logical.DeleteOperation,
		Path:       "sys/auth/" + path,
		Connection: getConnection(r),
	}))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}

type EnableAuthRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}
