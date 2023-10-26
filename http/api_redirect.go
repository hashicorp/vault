package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func handleApiRedirect(core *vault.Core, base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redir, err := core.GetAPIRedirect(r.Context(), r.URL.Path)
		if err != nil {
			core.Logger().Warn("error resolving potential API redirect", "error", err)
		} else {
			if redir != "" {
				w.Header().Set("Location", redir)
				w.WriteHeader(http.StatusFound)
				return
			}
		}
		if base != nil {
			base.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
