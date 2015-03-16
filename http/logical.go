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
		resp, err := core.HandleRequest(&logical.Request{
			Operation: op,
			Path:      path,
			Data:      req,
		})
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
			if resp.IsSecret && resp.Lease != nil {
				logicalResp.VaultId = resp.Lease.VaultID
				logicalResp.Renewable = resp.Lease.Renewable
				logicalResp.LeaseDuration = int(resp.Lease.Duration.Seconds())
				logicalResp.LeaseDurationMax = int(resp.Lease.MaxDuration.Seconds())
			}

			httpResp = logicalResp
		}

		// Respond
		respondOk(w, httpResp)
	})
}

type LogicalResponse struct {
	VaultId          string                 `json:"vault_id"`
	Renewable        bool                   `json:"renewable"`
	LeaseDuration    int                    `json:"lease_duration"`
	LeaseDurationMax int                    `json:"lease_duration_max"`
	Data             map[string]interface{} `json:"data"`
}
