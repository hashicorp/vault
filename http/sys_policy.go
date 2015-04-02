package http

import (
	"net/http"

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

type listPolicyResponse struct {
	Policies []string `json:"policies"`
}
