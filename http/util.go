// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/quotas"
)

var (
	genericWrapping = func(core *vault.Core, in http.Handler, props *vault.HandlerProperties) http.Handler {
		// Wrap the help wrapped handler with another layer with a generic
		// handler
		return wrapGenericHandler(core, in, props)
	}

	additionalRoutes = func(mux *http.ServeMux, core *vault.Core) {}

	nonVotersAllowed = false

	adjustResponse = func(core *vault.Core, w http.ResponseWriter, req *logical.Request) {}
)

func rateLimitQuotaWrapping(handler http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ns, err := namespace.FromContext(r.Context())
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		// We don't want to do buildLogicalRequestNoAuth here because, if the
		// request gets allowed by the quota, the same function will get called
		// again, which is not desired.
		path, status, err := buildLogicalPath(r)
		if err != nil || status != 0 {
			respondError(w, status, err)
			return
		}
		mountPath := strings.TrimPrefix(core.MatchingMount(r.Context(), path), ns.Path)

		// Clone body, so we do not close the request body reader
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.New("failed to read request body"))
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		quotaReq := &quotas.Request{
			Type:          quotas.TypeRateLimit,
			Path:          path,
			MountPath:     mountPath,
			NamespacePath: ns.Path,
			ClientAddress: parseRemoteIPAddress(r),
		}

		// This checks if any role based quota is required (LCQ or RLQ).
		requiresResolveRole, err := core.ResolveRoleForQuotas(r.Context(), quotaReq)
		if err != nil {
			core.Logger().Error("failed to lookup quotas", "path", path, "error", err)
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		// If any role-based quotas are enabled for this namespace/mount, just
		// do the role resolution once here.
		if requiresResolveRole {
			role := core.DetermineRoleFromLoginRequestFromBytes(r.Context(), mountPath, bodyBytes)
			// add an entry to the context to prevent recalculating request role unnecessarily
			r = r.WithContext(context.WithValue(r.Context(), logical.CtxKeyRequestRole{}, role))
			quotaReq.Role = role
		}

		quotaResp, err := core.ApplyRateLimitQuota(r.Context(), quotaReq)
		if err != nil {
			core.Logger().Error("failed to apply quota", "path", path, "error", err)
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		if core.RateLimitResponseHeadersEnabled() {
			for h, v := range quotaResp.Headers {
				w.Header().Set(h, v)
			}
		}

		if !quotaResp.Allowed {
			quotaErr := fmt.Errorf("request path %q: %w", path, quotas.ErrRateLimitQuotaExceeded)
			respondError(w, http.StatusTooManyRequests, quotaErr)

			if core.Logger().IsTrace() {
				core.Logger().Trace("request rejected due to rate limit quota violation", "request_path", path)
			}

			if core.RateLimitAuditLoggingEnabled() {
				req, _, status, err := buildLogicalRequestNoAuth(core.PerfStandby(), w, r)
				if err != nil || status != 0 {
					respondError(w, status, err)
					return
				}

				err = core.AuditLogger().AuditRequest(r.Context(), &logical.LogInput{
					Request:  req,
					OuterErr: quotaErr,
				})
				if err != nil {
					core.Logger().Warn("failed to audit log request rejection caused by rate limit quota violation", "error", err)
				}
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
