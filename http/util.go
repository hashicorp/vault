package http

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/quotas"
)

var (
	adjustRequest = func(c *vault.Core, r *http.Request) (*http.Request, int) {
		return r.WithContext(namespace.ContextWithNamespace(r.Context(), namespace.RootNamespace)), 0
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

		path, op, status, err := buildLogicalPathAndOp(r)
		if err != nil || status != 0 {
			respondError(w, status, err)
			return
		}

		quotaResp, err := core.ApplyRateLimitQuota(&quotas.Request{
			Type:          quotas.TypeRateLimit,
			Path:          path,
			MountPath:     strings.TrimPrefix(core.MatchingMount(r.Context(), path), ns.Path),
			NamespacePath: ns.Path,
			ClientAddress: parseRemoteIPAddress(r),
		})
		if err != nil {
			core.Logger().Error("failed to apply quota", "path", path, "error", err)
			respondError(w, http.StatusUnprocessableEntity, err)
			return
		}

		if !quotaResp.Allowed {
			quotaErr := errwrap.Wrapf(fmt.Sprintf("request path %q: {{err}}", path), quotas.ErrRateLimitQuotaExceeded)
			respondError(w, http.StatusTooManyRequests, quotaErr)

			if core.Logger().IsTrace() {
				core.Logger().Trace("request rejected due to lease count quota violation", "request_path", path)
			}

			requestId, err := uuid.GenerateUUID()
			if err != nil {
				respondError(w, http.StatusBadRequest, errwrap.Wrapf("failed to generate identifier for the request: {{err}}", err))
				return
			}

			req := &logical.Request{
				ID:         requestId,
				Operation:  op,
				Path:       path,
				Connection: getConnection(r),
				Headers:    r.Header,
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
