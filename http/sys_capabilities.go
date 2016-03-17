package http

import (
	"net/http"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysCapabilitiesAccessor(core *vault.Core) http.Handler {
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

		switch r.Method {
		case "PUT":
		case "POST":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		var data map[string]interface{}
		if err := parseRequest(r, &data); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		// Perform ACL checking, audit logging and route the request to
		// the system backend for request processing
		resp, ok := request(core, w, r, requestAuth(r, &logical.Request{
			Operation:  logical.UpdateOperation,
			Path:       path,
			Data:       data,
			Connection: getConnection(r),
		}))
		if !ok {
			return
		}

		respondLogical(w, r, path, false, resp)
	})

}

func handleSysCapabilities(core *vault.Core) http.Handler {
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

		switch r.Method {
		case "PUT":
		case "POST":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		var data map[string]interface{}
		if err := parseRequest(r, &data); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		if path == "sys/capabilities-self" {
			// Get the auth for the request so we can access the token directly
			req := requestAuth(r, &logical.Request{})
			path = "sys/capabilities"
			data["token"] = req.ClientToken
		}

		// Perform ACL checking, audit logging and route the request to
		// the system backend for request processing
		resp, ok := request(core, w, r, requestAuth(r, &logical.Request{
			Operation:  logical.UpdateOperation,
			Path:       path,
			Data:       data,
			Connection: getConnection(r),
		}))
		if !ok {
			return
		}

		respondLogical(w, r, path, false, resp)
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
	Accessor string `json:"accessor"`
	Path     string `json:"path"`
}
