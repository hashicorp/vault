package http

import (
	"net/http"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func wrapCORSHandler(h http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		corsConf := core.CORSConfig()
		statusCode := corsConf.ApplyHeaders(w, req)
		if statusCode != http.StatusOK {
			respondRaw(w, req, &logical.Response{
				Data: map[string]interface{}{
					logical.HTTPStatusCode: statusCode,
					logical.HTTPRawBody:    []byte(""),
				},
			})
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
