package http

import (
	"fmt"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"net/http"

	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func handleMetricsUnauthenticated(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Getting custom headers from listener's config
		la := w.Header().Get("X-Vault-Listener-Add")
		lc, err := core.GetCustomResponseHeaders(la)
		if err != nil {
			core.Logger().Debug("failed to get custom headers from listener config")
		}
		req := &logical.Request{Headers: r.Header}

		switch r.Method {
		case "GET":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil, lc)
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			respondError(w, http.StatusBadRequest, err, lc)
			return
		}

		format := r.Form.Get("format")
		if format == "" {
			format = metricsutil.FormatFromRequest(req)
		}

		// Define response
		resp := core.MetricsHelper().ResponseForFormat(format)

		// Manually extract the logical response and send back the information
		status := resp.Data[logical.HTTPStatusCode].(int)
		listenerutil.SetCustomResponseHeaders(lc, w, status)
		w.WriteHeader(status)
		w.Header().Set("Content-Type", resp.Data[logical.HTTPContentType].(string))
		switch v := resp.Data[logical.HTTPRawBody].(type) {
		case string:
			w.Write([]byte(v))
		case []byte:
			w.Write(v)
		default:
			respondError(w, http.StatusInternalServerError, fmt.Errorf("wrong response returned"), lc)
		}
	})
}
