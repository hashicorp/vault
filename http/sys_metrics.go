package http

import (
	"net/http"

	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func handleMetricsUnauthenticated(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &logical.Request{Headers: r.Header}
		format := r.Form.Get("format")
		if format == "" {
			format = metricsutil.FormatFromRequest(req)
		}

		// Define response
		resp, err := core.MetricsHelper().ResponseForFormat(format)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Manually extract the logical response and send back the information
		w.WriteHeader(resp.Data[logical.HTTPStatusCode].(int))
		w.Header().Set("Content-Type", resp.Data[logical.HTTPContentType].(string))
		w.Write(resp.Data[logical.HTTPRawBody].([]byte))
	})
}
