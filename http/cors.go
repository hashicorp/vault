package http

import (
	"fmt"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"net/http"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/vault"
)

var allowedMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodOptions,
	http.MethodPost,
	http.MethodPut,
	"LIST", // LIST is not an official HTTP method, but Vault supports it.
}

func wrapCORSHandler(h http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		corsConf := core.CORSConfig()

		// If CORS is not enabled or if no Origin header is present (i.e. the request
		// is from the Vault CLI. A browser will always send an Origin header), then
		// just return a 204.
		if !corsConf.IsEnabled() {
			h.ServeHTTP(w, req)
			return
		}

		origin := req.Header.Get("Origin")
		requestMethod := req.Header.Get("Access-Control-Request-Method")

		if origin == "" {
			h.ServeHTTP(w, req)
			return
		}
		// Getting custom headers from listener's config
		la := w.Header().Get("X-Vault-Listener-Add")
		lc, err := core.GetCustomResponseHeaders(la)
		if err != nil {
			core.Logger().Debug("failed to get custom headers from listener config")
		}
		// Return a 403 if the origin is not allowed to make cross-origin requests.
		if !corsConf.IsValidOrigin(origin) {
			respondError(w, http.StatusForbidden, fmt.Errorf("origin not allowed"), lc)
			return
		}

		if req.Method == http.MethodOptions && !strutil.StrListContains(allowedMethods, requestMethod) {
			status := http.StatusMethodNotAllowed
			listenerutil.SetCustomResponseHeaders(lc, w, status)
			w.WriteHeader(status)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")

		// apply headers for preflight requests
		if req.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConf.AllowedHeaders, ","))
			w.Header().Set("Access-Control-Max-Age", "300")

			return
		}

		h.ServeHTTP(w, req)
		return
	})
}
