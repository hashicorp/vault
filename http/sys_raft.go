package http

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"

	"github.com/hashicorp/vault/sdk/helper/tlsutil"
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

	var tlsConfig *tls.Config
	var err error
	if len(req.LeaderCACert) != 0 || len(req.LeaderClientCert) != 0 || len(req.LeaderClientKey) != 0 {
		tlsConfig, err = tlsutil.ClientTLSConfig([]byte(req.LeaderCACert), []byte(req.LeaderClientCert), []byte(req.LeaderClientKey))
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}
	}

	joined, err := core.JoinRaftCluster(context.Background(), req.LeaderAPIAddr, tlsConfig, req.Retry)
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
	LeaderAPIAddr    string `json:"leader_api_addr"`
	LeaderCACert     string `json:"leader_ca_cert":`
	LeaderClientCert string `json:"leader_client_cert"`
	LeaderClientKey  string `json:"leader_client_key"`
	Retry            bool   `json:"retry"`
}
