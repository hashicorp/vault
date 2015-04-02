package http

import (
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysListPolicies(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		resp, err := core.HandleRequest(requestAuth(r, &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "sys/policy",
		}))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		var policies []string
		policiesRaw, ok := resp.Data["keys"]
		if ok {
			policies = policiesRaw.([]string)
		}

		respondOk(w, &listPolicyResponse{Policies: policies})
	})
}

func handleSysPolicy(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysReadPolicy(core, w, r)
		case "PUT":
			fallthrough
		case "POST":
			handleSysWritePolicy(core, w, r)
		case "DELETE":
			handleSysDeletePolicy(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
	})
}

func handleSysDeletePolicy(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Determine the path...
	prefix := "/v1/sys/policy/"
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
		Operation: logical.DeleteOperation,
		Path:      "sys/policy/" + path,
	}))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}

func handleSysReadPolicy(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Determine the path...
	prefix := "/v1/sys/policy/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		respondError(w, http.StatusNotFound, nil)
		return
	}
	path := r.URL.Path[len(prefix):]
	if path == "" {
		respondError(w, http.StatusNotFound, nil)
		return
	}

	resp, err := core.HandleRequest(requestAuth(r, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/policy/" + path,
	}))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, resp.Data)
}

func handleSysWritePolicy(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Determine the path...
	prefix := "/v1/sys/policy/"
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
	var req writePolicyRequest
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	_, err := core.HandleRequest(requestAuth(r, &logical.Request{
		Operation: logical.WriteOperation,
		Path:      "sys/policy/" + path,
		Data: map[string]interface{}{
			"rules": req.Rules,
		},
	}))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}

type listPolicyResponse struct {
	Policies []string `json:"policies"`
}

type writePolicyRequest struct {
	Rules string `json:"rules"`
}
