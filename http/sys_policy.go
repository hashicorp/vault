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

		resp, ok := request(core, w, r, requestAuth(r, &logical.Request{
			Operation:  logical.ReadOperation,
			Path:       "sys/policy",
			Connection: getConnection(r),
		}))
		if !ok {
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

	_, ok := request(core, w, r, requestAuth(r, &logical.Request{
		Operation:  logical.DeleteOperation,
		Path:       "sys/policy/" + path,
		Connection: getConnection(r),
	}))
	if !ok {
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

	resp, ok := request(core, w, r, requestAuth(r, &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "sys/policy/" + path,
		Connection: getConnection(r),
	}))
	if !ok {
		return
	}
	if resp == nil {
		respondError(w, http.StatusNotFound, nil)
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

	_, ok := request(core, w, r, requestAuth(r, &logical.Request{
		Operation:  logical.WriteOperation,
		Path:       "sys/policy/" + path,
		Connection: getConnection(r),
		Data: map[string]interface{}{
			"rules": req.Rules,
		},
	}))
	if !ok {
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
