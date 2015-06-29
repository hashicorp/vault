package http

import (
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysListAudit(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		resp, ok := request(core, w, r, requestAuth(r, &logical.Request{
			Operation:  logical.ReadOperation,
			Path:       "sys/audit",
			Connection: getConnection(r),
		}))
		if !ok {
			return
		}

		respondOk(w, resp.Data)
	})
}

func handleSysAudit(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			fallthrough
		case "PUT":
			handleSysEnableAudit(core, w, r)
		case "DELETE":
			handleSysDisableAudit(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
	})
}

func handleSysDisableAudit(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Determine the path...
	prefix := "/v1/sys/audit/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		respondError(w, http.StatusNotFound, nil)
		return
	}
	path := r.URL.Path[len(prefix):]
	if path == "" {
		respondError(w, http.StatusNotFound, nil)
		return
	}

	_, ok := request(core, w, r, requestAuth(r, &logical.Request{
		Operation:  logical.DeleteOperation,
		Path:       "sys/audit/" + path,
		Connection: getConnection(r),
	}))
	if !ok {
		return
	}

	respondOk(w, nil)
}

func handleSysEnableAudit(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Determine the path...
	prefix := "/v1/sys/audit/"
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
	var req enableAuditRequest
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	_, ok := request(core, w, r, requestAuth(r, &logical.Request{
		Operation:  logical.WriteOperation,
		Path:       "sys/audit/" + path,
		Connection: getConnection(r),
		Data: map[string]interface{}{
			"type":        req.Type,
			"description": req.Description,
			"options":     req.Options,
		},
	}))
	if !ok {
		return
	}

	respondOk(w, nil)
}

type enableAuditRequest struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Options     map[string]string `json:"options"`
}
