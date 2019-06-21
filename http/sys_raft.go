package http

import (
	"context"
	"io"
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func handleSysRaftJoin(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST", "PUT":
			handleSysRaftJoinPost(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysRaftJoinPost(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req JoinRequest
	if _, err := parseRequest(core, r, w, &req); err != nil && err != io.EOF {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	joined, err := core.JoinRaftCluster(context.Background(), req.LeaderAddr, nil, req.Retry)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	resp := JoinResponse{
		Joined: joined,
	}
	respondOk(w, resp)
}

type JoinResponse struct {
	Joined bool `json:"joined"`
}

type JoinRequest struct {
	LeaderAddr string `json:"leader_api_addr"`
	CACert     string `json:"ca_cert":`
	Retry      bool   `json:"retry"`
}
