package http

import (
	"io"
	"net/http"
	"strings"
	"time"

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
		case "DELETE":
			op = logical.DeleteOperation
		case "GET":
			op = logical.ReadOperation
		case "POST":
			fallthrough
		case "PUT":
			op = logical.WriteOperation
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Parse the request if we can
		var req map[string]interface{}
		if op == logical.WriteOperation {
			err := parseRequest(r, &req)
			if err == io.EOF {
				req = nil
				err = nil
			}
			if err != nil {
				respondError(w, http.StatusBadRequest, err)
				return
			}
		}

		// Make the internal request. We attach the connection info
		// as well in case this is an authentication request that requires
		// it. Vault core handles stripping this if we need to.
		resp, ok := request(core, w, r, requestAuth(r, &logical.Request{
			Operation: op,
			Path:      path,
			Data:      req,
			Connection: &logical.Connection{
				RemoteAddr: r.RemoteAddr,
				ConnState:  r.TLS,
			},
		}))
		if !ok {
			return
		}
		if op == logical.ReadOperation && resp == nil {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		// Build the proper response
		respondLogical(w, r, path, resp)
	})
}

func respondLogical(w http.ResponseWriter, r *http.Request, path string, resp *logical.Response) {
	var httpResp interface{}
	if resp != nil {
		if resp.Redirect != "" {
			// If we have a redirect, redirect! We use a 302 code
			// because we don't actually know if its permanent.
			http.Redirect(w, r, resp.Redirect, 302)
			return
		}

		logicalResp := &LogicalResponse{Data: resp.Data}
		if resp.Secret != nil {
			logicalResp.LeaseID = resp.Secret.LeaseID
			logicalResp.Renewable = resp.Secret.Renewable
			logicalResp.LeaseDuration = int(resp.Secret.Lease.Seconds())
		}

		// If we have authentication information, then set the cookie
		// and setup the result structure.
		if resp.Auth != nil {
			expireDuration := 365 * 24 * time.Hour
			if logicalResp.LeaseDuration != 0 {
				expireDuration =
					time.Duration(logicalResp.LeaseDuration) * time.Second
			}

			// Do not set the token as the auth cookie if the endpoint
			// is the token store. Otherwise, attempting to create a token
			// will cause the client to be authenticated as that token.
			if !strings.HasPrefix(path, "auth/token/") {
				http.SetCookie(w, &http.Cookie{
					Name:    AuthCookieName,
					Value:   resp.Auth.ClientToken,
					Path:    "/",
					Expires: time.Now().UTC().Add(expireDuration),
				})
			}

			logicalResp.Auth = &Auth{
				ClientToken:   resp.Auth.ClientToken,
				Policies:      resp.Auth.Policies,
				Metadata:      resp.Auth.Metadata,
				LeaseDuration: int(resp.Auth.Lease.Seconds()),
				Renewable:     resp.Auth.Renewable,
			}
		}

		httpResp = logicalResp
	}

	// Respond
	respondOk(w, httpResp)
}

type LogicalResponse struct {
	LeaseID       string                 `json:"lease_id"`
	Renewable     bool                   `json:"renewable"`
	LeaseDuration int                    `json:"lease_duration"`
	Data          map[string]interface{} `json:"data"`
	Auth          *Auth                  `json:"auth"`
}

type Auth struct {
	ClientToken   string            `json:"client_token"`
	Policies      []string          `json:"policies"`
	Metadata      map[string]string `json:"metadata"`
	LeaseDuration int               `json:"lease_duration"`
	Renewable     bool              `json:"renewable"`
}
