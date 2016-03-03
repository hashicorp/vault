package http

import (
	"log"
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysCapabilities(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
		case "POST":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		log.Printf("r.URL.Path: %s\n", r.URL.Path)
		// Get the auth for the request so we can access the token directly
		req := requestAuth(r, &logical.Request{})
		log.Printf("handleSysCapabilities req:%#v\n", req)

		// Parse the request if we can
		var data capabilitiesRequest
		if err := parseRequest(r, &data); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		if data.Token == "" {
			data.Token = req.ClientToken
		}

		capabilities, err := core.Capabilities(data.Token, data.Path)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if capabilities == nil {
			respondOk(w, &capabilitiesResponse{Capabilities: nil})
			return
		}

		respondOk(w, &capabilitiesResponse{
			Capabilities: capabilities.Capabilities,
		})
	})

}

type capabilitiesResponse struct {
	Capabilities []string `json:"capabilities"`
}

type capabilitiesRequest struct {
	Token string `json:"token"`
	Path  string `json:"path"`
}
