package http

import (
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func wrapHelpHandler(h http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		// If the help parameter is not blank, then show the help. We request
		// forward because standby nodes do not have mounts and other state.
		if v := req.URL.Query().Get("help"); v != "" || req.Method == "HELP" {
			handleRequestForwarding(core,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handleHelp(core, w, r)
				})).ServeHTTP(writer, req)
			return
		}

		h.ServeHTTP(writer, req)
		return
	})
}

func handleHelp(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	ns, err := namespace.FromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusBadRequest, nil)
		return
	}
	path := ns.TrimmedPath(r.URL.Path[len("/v1/"):])

	req, err := requestAuth(core, r, &logical.Request{
		Operation:  logical.HelpOperation,
		Path:       path,
		Connection: getConnection(r),
	})
	if err != nil {
		if errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
			respondError(w, http.StatusForbidden, nil)
			return
		}
		respondError(w, http.StatusBadRequest, errwrap.Wrapf("error performing token check: {{err}}", err))
		return
	}

	resp, err := core.HandleRequest(r.Context(), req)
	if err != nil {
		respondErrorCommon(w, req, resp, err)
		return
	}

	respondOk(w, resp.Data)
}
