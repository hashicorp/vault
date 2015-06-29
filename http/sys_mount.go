package http

import (
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysMounts(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysListMounts(core).ServeHTTP(w, r)
		case "PUT", "POST":
			fallthrough
		case "DELETE":
			handleSysMountUnmount(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
	})
}

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

func handleSysListMounts(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		resp, err := core.HandleRequest(requestAuth(r, &logical.Request{
			Operation:  logical.ReadOperation,
			Path:       "sys/mounts",
			Connection: getConnection(r),
		}))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		respondOk(w, resp.Data)
	})
}

func handleSysMountUnmount(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT", "POST":
	case "DELETE":
	default:
		respondError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// Determine the path...
	prefix := "/v1/sys/mounts/"
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
		handleSysMount(core, w, r, path)
	case "DELETE":
		handleSysUnmount(core, w, r, path)
	default:
		panic("should never happen")
	}
}

func handleSysMount(
	core *vault.Core,
	w http.ResponseWriter,
	r *http.Request,
	path string) {
	// Parse the request if we can
	var req MountRequest
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	_, err := core.HandleRequest(requestAuth(r, &logical.Request{
		Operation:  logical.WriteOperation,
		Path:       "sys/mounts/" + path,
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

func handleSysUnmount(
	core *vault.Core,
	w http.ResponseWriter,
	r *http.Request,
	path string) {
	_, err := core.HandleRequest(requestAuth(r, &logical.Request{
		Operation:  logical.DeleteOperation,
		Path:       "sys/mounts/" + path,
		Connection: getConnection(r),
	}))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}

type MountRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type RemountRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}
