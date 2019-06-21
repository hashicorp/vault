package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
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

	var tlsConfig *tls.Config
	switch {
	case req.LeaderCACert != "" && req.LeaderClientCert != "" && req.LeaderClientKey != "":
		// Create TLS Config
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM([]byte(req.LeaderCACert))

		cert, err := tls.X509KeyPair([]byte(req.LeaderClientCert), []byte(req.LeaderClientKey))
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Errorf("invalid key pair: %v", err))
			return
		}

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			MinVersion:   tls.VersionTLS12,
		}

		tlsConfig.BuildNameToCertificate()

	case req.LeaderCACert != "":
		respondError(w, http.StatusBadRequest, errors.New("ca_cert, client_key, client_cert must all be set; or none should be set"))
		return
	case req.LeaderClientCert != "":
		respondError(w, http.StatusBadRequest, errors.New("ca_cert, client_key, client_cert must all be set; or none should be set"))
		return
	case req.LeaderClientKey != "":
		respondError(w, http.StatusBadRequest, errors.New("ca_cert, client_key, client_cert must all be set; or none should be set"))
		return
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
