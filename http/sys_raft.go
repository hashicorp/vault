// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
)

func handleSysRaftBootstrap(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST", "PUT":
			if core.Sealed() {
				respondError(w, http.StatusBadRequest, errors.New("node must be unsealed to bootstrap"))
			}

			if err := core.RaftBootstrap(context.Background(), false); err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return
			}

		default:
			respondError(w, http.StatusBadRequest, nil)
		}
	})
}

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
	if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil && err != io.EOF {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if req.NonVoter && !nonVotersAllowed {
		respondError(w, http.StatusBadRequest, errors.New("non-voting nodes not allowed"))
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
		tlsConfig.ServerName = req.LeaderTLSServerName
	}

	if req.AutoJoinScheme != "" && (req.AutoJoinScheme != "http" && req.AutoJoinScheme != "https") {
		respondError(w, http.StatusBadRequest, fmt.Errorf("invalid scheme %q; must either be http or https", req.AutoJoinScheme))
		return
	}

	leaderInfos := []*raft.LeaderJoinInfo{
		{
			AutoJoin:       req.AutoJoin,
			AutoJoinScheme: req.AutoJoinScheme,
			AutoJoinPort:   req.AutoJoinPort,
			LeaderAPIAddr:  configutil.NormalizeAddr(req.LeaderAPIAddr),
			TLSConfig:      tlsConfig,
			Retry:          req.Retry,
		},
	}

	joined, err := core.JoinRaftCluster(context.Background(), leaderInfos, req.NonVoter)
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
	AutoJoin            string `json:"auto_join"`
	AutoJoinScheme      string `json:"auto_join_scheme"`
	AutoJoinPort        uint   `json:"auto_join_port"`
	LeaderAPIAddr       string `json:"leader_api_addr"`
	LeaderCACert        string `json:"leader_ca_cert"`
	LeaderClientCert    string `json:"leader_client_cert"`
	LeaderClientKey     string `json:"leader_client_key"`
	LeaderTLSServerName string `json:"leader_tls_servername"`
	Retry               bool   `json:"retry"`
	NonVoter            bool   `json:"non_voter"`
}
