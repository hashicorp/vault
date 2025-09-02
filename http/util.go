// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/limits"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/quotas"
)

var nonVotersAllowed = false

func wrapMaxRequestSizeHandler(handler http.Handler, props *vault.HandlerProperties) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var maxRequestSize, maxJSONDepth, maxStringValueLength, maxObjectEntryCount, maxArrayElementCount int64

		if props.ListenerConfig != nil {
			maxRequestSize = props.ListenerConfig.MaxRequestSize
			maxJSONDepth = props.ListenerConfig.CustomMaxJSONDepth
			maxStringValueLength = props.ListenerConfig.CustomMaxJSONStringValueLength
			maxObjectEntryCount = props.ListenerConfig.CustomMaxJSONObjectEntryCount
			maxArrayElementCount = props.ListenerConfig.CustomMaxJSONArrayElementCount
		}

		if maxRequestSize == 0 {
			maxRequestSize = DefaultMaxRequestSize
		}
		if maxJSONDepth == 0 {
			maxJSONDepth = CustomMaxJSONDepth
		}
		if maxStringValueLength == 0 {
			maxStringValueLength = CustomMaxJSONStringValueLength
		}
		if maxObjectEntryCount == 0 {
			maxObjectEntryCount = CustomMaxJSONObjectEntryCount
		}
		if maxArrayElementCount == 0 {
			maxArrayElementCount = CustomMaxJSONArrayElementCount
		}

		jsonLimits := jsonutil.JSONLimits{
			MaxDepth:             int(maxJSONDepth),
			MaxStringValueLength: int(maxStringValueLength),
			MaxObjectEntryCount:  int(maxObjectEntryCount),
			MaxArrayElementCount: int(maxArrayElementCount),
		}

		// If the payload is JSON, the VerifyMaxDepthStreaming function will perform validations.
		buf, err := jsonLimitsValidation(w, r, maxRequestSize, jsonLimits)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		// Replace the body and update the context.
		// This ensures the request object is in a consistent state for all downstream handlers.
		// Because the original request body stream has been fully consumed by io.ReadAll,
		// we must replace it so that subsequent handlers can read the content.
		r.Body = newMultiReaderCloser(buf, r.Body)
		contextBody := r.Body
		ctx := logical.CreateContextOriginalBody(r.Context(), contextBody)

		if maxRequestSize > 0 {
			r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
		}
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

func jsonLimitsValidation(w http.ResponseWriter, r *http.Request, maxRequestSize int64, jsonLimits jsonutil.JSONLimits) (*bytes.Buffer, error) {
	// The TeeReader reads from the original body and writes a copy to our buffer.
	// We wrap the original body with a MaxBytesReader first to enforce the hard size limit.
	var limitedTeeReader io.Reader
	buf := &bytes.Buffer{}
	bodyReader := r.Body
	if maxRequestSize > 0 {
		bodyReader = http.MaxBytesReader(w, r.Body, maxRequestSize)
	}
	limitedTeeReader = io.TeeReader(bodyReader, buf)
	_, err := jsonutil.VerifyMaxDepthStreaming(limitedTeeReader, jsonLimits)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func wrapRequestLimiterHandler(handler http.Handler, props *vault.HandlerProperties) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := r.WithContext(
			context.WithValue(
				r.Context(),
				limits.CtxKeyDisableRequestLimiter{},
				props.ListenerConfig.DisableRequestLimiter,
			),
		)
		handler.ServeHTTP(w, request)
	})
}

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

		quotaReq := &quotas.Request{
			Type:          quotas.TypeRateLimit,
			Path:          path,
			MountPath:     mountPath,
			NamespacePath: ns.Path,
			ClientAddress: parseRemoteIPAddress(r),
		}

		entRlqRequestFields(core, r, quotaReq)

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
			buf := bytes.Buffer{}
			teeReader := io.TeeReader(r.Body, &buf)
			role := core.DetermineRoleFromLoginRequestFromReader(r.Context(), mountPath, teeReader)

			// Reset the body if it was read
			if buf.Len() > 0 {
				r.Body = io.NopCloser(&buf)
				originalBody, ok := logical.ContextOriginalBodyValue(r.Context())
				if ok {
					r = r.WithContext(logical.CreateContextOriginalBody(r.Context(), newMultiReaderCloser(&buf, originalBody)))
				}
			}
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
				req, _, status, err := buildLogicalRequestNoAuth(core.PerfStandby(), core.RouterAccess(), w, r)
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

func disableReplicationStatusEndpointWrapping(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := r.WithContext(logical.CreateContextDisableReplicationStatusEndpoints(r.Context(), true))

		h.ServeHTTP(w, request)
	})
}

func redactionSettingsWrapping(h http.Handler, redactVersion, redactAddresses, redactClusterName bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := r.WithContext(logical.CreateContextRedactionSettings(r.Context(), redactVersion, redactAddresses, redactClusterName))

		h.ServeHTTP(w, request)
	})
}

func parseRemoteIPAddress(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}

type multiReaderCloser struct {
	readers []io.Reader
	io.Reader
}

func newMultiReaderCloser(readers ...io.Reader) *multiReaderCloser {
	return &multiReaderCloser{
		readers: readers,
		Reader:  io.MultiReader(readers...),
	}
}

func (m *multiReaderCloser) Close() error {
	var err error
	for _, r := range m.readers {
		if c, ok := r.(io.Closer); ok {
			err = multierror.Append(err, c.Close())
		}
	}
	return err
}
