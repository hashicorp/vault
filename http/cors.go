package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func wrapCORSHandler(h http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		corsConf := core.CORSConfig()

		// If CORS is not enabled or if no Origin header is present (i.e. the request
		// is from the Vault CLI. A browser will always send an Origin header), then
		// just return a 204.
		if !corsConf.Enabled() || req.Header.Get("Origin") == "" {
			h.ServeHTTP(w, req)
			return
		}

		statusCode := corsConf.ApplyHeaders(w, req)
		if statusCode != http.StatusNoContent {
			h.ServeHTTP(w, req)
			return
		}

		// For pre-flight requests just send back the headers and return.
		if req.Method == http.MethodOptions {
			return
		}

		h.ServeHTTP(w, req)
		return
	})
}
