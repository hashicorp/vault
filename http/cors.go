package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func wrapCORSHandler(h http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		corsConf := core.CORSConfig()

		origin := req.Header.Get("Origin")
		requestMethod := req.Header.Get("Access-Control-Request-Method")

		// If CORS is not enabled or if no Origin header is present (i.e. the request
		// is from the Vault CLI. A browser will always send an Origin header), then
		// just return a 204.
		if !corsConf.IsEnabled() || origin == "" {
			h.ServeHTTP(w, req)
			return
		}

		// Return a 403 if the origin is not
		// allowed to make cross-origin requests.
		if !corsConf.IsValidOrigin(origin) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if req.Method == http.MethodOptions && !corsConf.IsValidMethod(requestMethod) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		corsConf.ApplyHeaders(w, req)

		h.ServeHTTP(w, req)
		return
	})
}
