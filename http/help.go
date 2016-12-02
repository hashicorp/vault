package http

import (
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func handleHelpHandler(h http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// If the help parameter is not blank, then show the help
		if v := req.URL.Query().Get("help"); v != "" || req.Method == "HELP" {
			handleHelp(core, w, req)
			return
		}

		h.ServeHTTP(w, req)
		return
	})
}

func handleHelp(core *vault.Core, w http.ResponseWriter, req *http.Request) {
	path, ok := stripPrefix("/v1/", req.URL.Path)
	if !ok {
		respondError(w, http.StatusNotFound, nil)
		return
	}

	lreq := requestAuth(core, req, &logical.Request{
		Operation:  logical.HelpOperation,
		Path:       path,
		Connection: getConnection(req),
	})

	resp, err := core.HandleRequest(lreq)
	if err != nil {
		respondErrorCommon(w, resp, err)
		return
	}

	respondOk(w, resp.Data)
}
