package http

import (
	"net/http"
	"strings"

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

		// Get the auth for the request so we can access the token directly
		req := requestAuth(r, &logical.Request{})

		// Parse the request if we can
		var data capabilitiesRequest
		if err := parseRequest(r, &data); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/v1/sys/capabilities-self") {
			data.Token = req.ClientToken
		}

		capabilities, err := core.Capabilities(data.Token, data.Path)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if capabilities == nil {
			respondOk(w, &capabilitiesResponse{Message: "Token has no capabilities on the given path"})
			return
		}

		var response capabilitiesResponse
		switch capabilities.Root {
		case true:
			response.Message = `Thij is a 'root' token. It has all the capabilities on all the paths.
This token can be used on any valid path.`
			response.Capabilities = nil
		case false:
			response.Message = ""
			response.Capabilities = capabilities.Capabilities
		}

		respondOk(w, response)
	})

}

type capabilitiesResponse struct {
	Message      string   `json:"message"`
	Capabilities []string `json:"capabilities"`
}

type capabilitiesRequest struct {
	Token string `json:"token"`
	Path  string `json:"path"`
}
