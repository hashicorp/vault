package http

import (
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysCapabilitiesAccessor(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
		case "POST":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Parse the request if we can
		var data capabilitiesAccessorRequest
		if err := parseRequest(r, &data); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		capabilities, err := core.CapabilitiesAccessor(data.AccessorID, data.Path)
		if err != nil {
			respondErrorStatus(w, err)
			return
		}

		respondOk(w, &capabilitiesResponse{
			Capabilities: capabilities,
		})
	})

}

func handleSysCapabilities(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
		case "POST":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Parse the request if we can
		var data capabilitiesRequest
		if err := parseRequest(r, &data); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/v1/sys/capabilities-self") {
			// Get the auth for the request so we can access the token directly
			req := requestAuth(r, &logical.Request{})
			data.Token = req.ClientToken
		}

		capabilities, err := core.Capabilities(data.Token, data.Path)
		if err != nil {
			respondErrorStatus(w, err)
			return
		}

		respondOk(w, &capabilitiesResponse{
			Capabilities: capabilities,
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

type capabilitiesAccessorRequest struct {
	AccessorID string `json:"accessor_id"`
	Path       string `json:"path"`
}
