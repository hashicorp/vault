package http

import (
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func wrapCORSHandler(h http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// If the help parameter is not blank, then show the help
		if req.Header.Get(NoCORS) == "" {
			statusCode := handleCORS(core, w, req)
			if statusCode != http.StatusOK {
				respondRaw(w, req, &logical.Response{
					Data: map[string]interface{}{
						logical.HTTPStatusCode: statusCode,
						logical.HTTPRawBody:    []byte(""),
					},
				})
				return
			}
		}

		h.ServeHTTP(w, req)
		return
	})
}

// HandleCORS adds required headers to properly respond to
// requests that require Cross Origin Resource Sharing (CORS) headers.
func handleCORS(core *vault.Core, w http.ResponseWriter, req *http.Request) int {
	corsConf := core.CORSConfig()
	statusCode := corsConf.ApplyHeaders(w, req)

	return statusCode
}
