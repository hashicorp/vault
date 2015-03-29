package http

import (
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleLogical(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determine the path...
		if !strings.HasPrefix(r.URL.Path, "/v1/") {
			respondError(w, http.StatusNotFound, nil)
			return
		}
		path := r.URL.Path[len("/v1/"):]
		if path == "" {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		// Determine the operation
		var op logical.Operation
		switch r.Method {
		case "GET":
			op = logical.ReadOperation
		case "PUT":
			op = logical.WriteOperation
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Parse the request if we can
		var req map[string]interface{}
		if op == logical.WriteOperation {
			if err := parseRequest(r, &req); err != nil {
				respondError(w, http.StatusBadRequest, err)
				return
			}
		}

		// Make the internal request
		resp, err := core.HandleRequest(requestAuth(r, &logical.Request{
			Operation: op,
			Path:      path,
			Data:      req,
		}))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if op == logical.ReadOperation && resp == nil {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		var httpResp interface{}
		if resp != nil {
			logicalResp := &LogicalResponse{Data: resp.Data}
			if resp.Secret != nil {
				logicalResp.VaultId = resp.Secret.VaultID
				logicalResp.Renewable = resp.Secret.Renewable
				logicalResp.LeaseDuration = int(resp.Secret.Lease.Seconds())
			}

			httpResp = logicalResp
		}

		// Respond
		respondOk(w, httpResp)
	})
}

type LogicalResponse struct {
	VaultId       string                 `json:"vault_id"`
	Renewable     bool                   `json:"renewable"`
	LeaseDuration int                    `json:"lease_duration"`
	Data          map[string]interface{} `json:"data"`
}
