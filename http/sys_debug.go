package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

// handleSysDebug is a helper that checks for the appropriate permission on the
// path before routing the request to the provided /sys/debug handler.
func handleSysDebug(core *vault.Core, handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _, statusCode, err := buildLogicalRequest(core, w, r)
		if err != nil || statusCode != 0 {
			respondError(w, statusCode, err)
			return
		}

		switch req.Operation {
		case logical.ReadOperation:
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		if err := core.SysDebugTokenCheck(r.Context(), req); err != nil {
			if errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
				respondError(w, http.StatusForbidden, err)
			} else {
				respondError(w, http.StatusBadRequest, err)
			}
			return
		}

		handler(w, r)
	})
}

// handleSysDebugMaxDuration takes care of capping the maximum duration for pprof
// queries that involve probing over a time period, e.g. profile and trace, by
// overwriting the "seconds" query parameter value if greater than the allowed
// duration.
func handleSysDebugMaxDuration(maxDuration time.Duration, handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parseSecs := func(vals url.Values) (time.Duration, error) {
			var secs int
			var err error

			if secsStr := vals.Get("seconds"); secsStr != "" {
				secs, err = strconv.Atoi(secsStr)
				if err != nil {
					return 0, err
				}
			}
			return time.Duration(secs) * time.Second, err
		}

		params := r.URL.Query()
		seconds, err := parseSecs(params)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		if seconds > maxDuration {
			newSecsVal := fmt.Sprintf("%d", int64(maxDuration/time.Second))
			params.Set("seconds", newSecsVal)
			r.URL.RawQuery = params.Encode()
		}

		handler(w, r)
	})
}
