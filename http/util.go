package http

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/quotas"
)

var (
	adjustRequest = func(c *vault.Core, r *http.Request) (*http.Request, int) {
		return r, 0
	}

	genericWrapping = func(core *vault.Core, in http.Handler, props *vault.HandlerProperties) http.Handler {
		// Wrap the help wrapped handler with another layer with a generic
		// handler
		return wrapGenericHandler(core, in, props)
	}

	additionalRoutes = func(mux *http.ServeMux, core *vault.Core) {}

	nonVotersAllowed = false
)

func rateLimitQuotaWrapping(handler http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ns, err := namespace.FromContext(r.Context())
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		req := w.(*LogicalResponseWriter).request
		quotaResp, err := core.ApplyRateLimitQuota(&quotas.Request{
			Type:          quotas.TypeRateLimit,
			Path:          req.Path,
			MountPath:     strings.TrimPrefix(req.MountPoint, ns.Path),
			NamespacePath: ns.Path,
			ClientAddress: parseRemoteIPAddress(r),
		})
		if err != nil {
			core.Logger().Error("failed to apply quota", "path", req.Path, "error", err)
			respondError(w, http.StatusUnprocessableEntity, err)
			return
		}

		if !quotaResp.Allowed {
			quotaErr := errwrap.Wrapf(fmt.Sprintf("request path %q: {{err}}", req.Path), quotas.ErrRateLimitQuotaExceeded)
			respondError(w, http.StatusTooManyRequests, quotaErr)

			if core.Logger().IsTrace() {
				core.Logger().Trace("request rejected due to lease count quota violation", "request_path", req.Path)
			}

			if core.RateLimitAuditLoggingEnabled() {
				_ = core.AuditLogger().AuditRequest(r.Context(), &logical.LogInput{
					Request:  req,
					OuterErr: quotaErr,
				})
			}

			return
		}

		handler.ServeHTTP(w, r)
		return
	})
}

func parseRemoteIPAddress(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}
