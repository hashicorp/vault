package http

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func handleMetricsUnauthenticated(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &logical.Request{Headers: r.Header}

		switch r.Method {
		case "GET":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		format := r.Form.Get("format")
		if format == "" {
			format = metricsutil.FormatFromRequest(req)
		}

		// Define response
		resp := core.MetricsHelper().ResponseForFormat(format)

		// Manually extract the logical response and send back the information
		w.WriteHeader(resp.Data[logical.HTTPStatusCode].(int))
		w.Header().Set("Content-Type", resp.Data[logical.HTTPContentType].(string))
		switch v := resp.Data[logical.HTTPRawBody].(type) {
		case string:
			w.Write([]byte(v))
		case []byte:
			w.Write(v)
		default:
			respondError(w, http.StatusInternalServerError, fmt.Errorf("wrong response returned"))
		}
	})
}
