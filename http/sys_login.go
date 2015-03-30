package http

import (
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/credential"
	"github.com/hashicorp/vault/vault"
)

func handleSysLogin(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Determine the path...
		prefix := "/v1/sys/login/"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			respondError(w, http.StatusNotFound, nil)
			return
		}
		path := r.URL.Path[len(prefix):]
		if path == "" {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		// Parse the IP address
		ipaddr, err := net.ResolveIPAddr("ip", r.RemoteAddr)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		// Do the login request
		resp, err := core.HandleLogin(&credential.Request{
			Path:       path,
			RemoteAddr: ipaddr,
			ConnState:  r.TLS,
		})

		// Determine the response to send. If we were given a secret,
		// then we add the secret to the response.
		var httpResp interface{}
		if resp != nil {
			// TODO: redirect

			logicalResp := &LogicalResponse{Data: resp.Data}
			if resp.Secret != nil {
				logicalResp.VaultId = resp.Secret.VaultID
				logicalResp.Renewable = resp.Secret.Renewable
				logicalResp.LeaseDuration = int(resp.Secret.Lease.Seconds())
			}

			httpResp = logicalResp
		}

		// Respond with the secret and/or data
		respondOk(w, httpResp)
	})
}
