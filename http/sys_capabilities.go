package http

import (
	"log"
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysCapabilitiesSelf(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" && r.Method != "PUT" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
		log.Printf("vishal: handleSysCapabilitiesSelf: r:%#v, r.URL:%s r.URL.Path:%s\n", r, r.URL, r.URL.Path)
		// Parse the request if we can
		var req capabilitiesRequest
		if err := parseRequest(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		resp, ok := request(core, w, r, requestAuth(r, &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "sys/capabilities-self",
			//			Connection: getConnection(r),
			Data: map[string]interface{}{
				"path": req.Path,
			},
		}))
		if !ok {
			return
		}
		if resp == nil {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		var capabilities []string
		capabilitiesRaw, ok := resp.Data["keys"]
		if ok {
			capabilities = capabilitiesRaw.([]string)
		}

		respondOk(w, &capabilitiesResponse{Capabilities: capabilities})
	})
}

func handleSysCapabilities(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" && r.Method != "PUT" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
		log.Printf("vishal: handleSysCapabilities: r: %#v\n", r)
		// Parse the request if we can
		var req capabilitiesRequest
		if err := parseRequest(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		resp, ok := request(core, w, r, requestAuth(r, &logical.Request{
			Operation:  logical.UpdateOperation,
			Path:       "sys/capabilities",
			Connection: getConnection(r),
			Data: map[string]interface{}{
				"token": req.Token,
				"path":  req.Path,
			},
		}))
		if !ok {
			return
		}
		if resp == nil {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		var capabilities []string
		capabilitiesRaw, ok := resp.Data["keys"]
		if ok {
			capabilities = capabilitiesRaw.([]string)
		}

		respondOk(w, &capabilitiesResponse{Capabilities: capabilities})
	})
}

type capabilitiesResponse struct {
	Capabilities []string `json:"policies"`
}

type capabilitiesRequest struct {
	Token string `json:"token"`
	Path  string `json:"path"`
}
