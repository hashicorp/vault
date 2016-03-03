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

		resp, err := core.Capabilities(data.Token, data.Path)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if resp == nil {
			respondOk(w, &capabilitiesResponse{
				Message:      "Token has no capabilities on the path",
				Capabilities: nil,
			})
			return
		}

		var result capabilitiesResponse
		switch resp.Root {
		case true:
			result.Message = "This is a 'root' token. It has all the capabilities on all the 'valid' paths."
			result.Capabilities = nil
		case false:
			if len(resp.Capabilities) == 0 {
				result.Message = "Token has no capabilities on the path"
			} else {
				result.Message = ""
			}
			result.Capabilities = resp.Capabilities
		}

		respondOk(w, result)
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
