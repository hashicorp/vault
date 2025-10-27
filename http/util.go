// Copyright IBM Corp. 2016, 2025
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

// ctxKeyRoleBasedQuota is used to signal that role-based quota resolution
// is needed for the request.
type ctxKeyRoleBasedQuota struct{}

func (c ctxKeyRoleBasedQuota) String() string {
	return "role-based-quota"
}

// resetBodyIfRead deals with creating a new body for the request if a portion
// of it has been read into buf. This function handles both the main body and
// the full, non-limited original body stored in the context.
func resetBodyIfRead(r *http.Request, buf *bytes.Buffer) *http.Request {
	if buf.Len() > 0 {
		r.Body = newMultiReaderCloser(buf, r.Body)
		originalBody, ok := logical.ContextOriginalBodyValue(r.Context())
		if ok {
			r = r.WithContext(logical.CreateContextOriginalBody(r.Context(), newMultiReaderCloser(buf, originalBody)))
		}
	}
	return r
}

// wrapMaxRequestSizeHandler limits the size of the request body to the
// configured size
func wrapMaxRequestSizeHandler(handler http.Handler, props *vault.HandlerProperties) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var maxRequestSize int64
		if props.ListenerConfig != nil {
			maxRequestSize = props.ListenerConfig.MaxRequestSize
		}
		if maxRequestSize == 0 {
			maxRequestSize = DefaultMaxRequestSize
		}
		ctx := r.Context()
		originalBody := r.Body
		if maxRequestSize > 0 {
			r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
		}
		ctx = logical.CreateContextOriginalBody(ctx, originalBody)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

// wrapJSONLimitsHandler enforces limits on JSON request bodies
func wrapJSONLimitsHandler(handler http.Handler, props *vault.HandlerProperties) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var maxRequestSize, maxJSONDepth, maxStringValueLength, maxObjectEntryCount, maxArrayElementCount, maxToken int64

		if props.ListenerConfig != nil {
			maxRequestSize = props.ListenerConfig.MaxRequestSize
			maxJSONDepth = props.ListenerConfig.CustomMaxJSONDepth
			maxStringValueLength = props.ListenerConfig.CustomMaxJSONStringValueLength
			maxObjectEntryCount = props.ListenerConfig.CustomMaxJSONObjectEntryCount
			maxArrayElementCount = props.ListenerConfig.CustomMaxJSONArrayElementCount
			maxToken = props.ListenerConfig.CustomMaxJSONToken
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
		if maxToken == 0 {
			maxToken = CustomMaxJSONToken
		}
		jsonLimits := jsonutil.JSONLimits{
			MaxDepth:             int(maxJSONDepth),
			MaxStringValueLength: int(maxStringValueLength),
			MaxObjectEntryCount:  int(maxObjectEntryCount),
			MaxArrayElementCount: int(maxArrayElementCount),
			MaxTokens:            int(maxToken),
		}

		// If the payload is JSON, the VerifyMaxDepthStreaming function will perform validations.
		buf, err := jsonLimitsValidation(r, maxRequestSize, jsonLimits)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		r = resetBodyIfRead(r, buf)

		handler.ServeHTTP(w, r)
	})
}

func jsonLimitsValidation(r *http.Request, maxRequestSize int64, jsonLimits jsonutil.JSONLimits) (*bytes.Buffer, error) {
	// The TeeReader reads from the original body and writes a copy to our buffer.
	var limitedTeeReader io.Reader
	buf := &bytes.Buffer{}
	bodyReader := r.Body
	limitedTeeReader = io.TeeReader(bodyReader, buf)
	var maxSize *int64
	if maxRequestSize > 0 {
		maxSize = &maxRequestSize
	}
	_, err := jsonutil.VerifyMaxDepthStreaming(limitedTeeReader, jsonLimits, maxSize)
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

// withRoleRateLimitQuotaWrapping performs any quota checking for request
// that require role resolution
func withRoleRateLimitQuotaWrapping(handler http.Handler, core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if there's a quota that requires role resolution, it will be in the
		// context
		// if the context value is nil, then no role-based quota resolution is
		// needed and we can continue handling the request
		quotaReqValue := r.Context().Value(ctxKeyRoleBasedQuota{})
		if quotaReqValue == nil {
			handler.ServeHTTP(w, r)
			return
		}

		quotaReq := quotaReqValue.(*quotas.Request)

		buf := bytes.Buffer{}
		teeReader := io.TeeReader(r.Body, &buf)
		role := core.DetermineRoleFromLoginRequestFromReader(r.Context(), quotaReq.MountPath, teeReader, getConnection(r), r.Header)

		// Reset the body if it was read
		r = resetBodyIfRead(r, &buf)

		// add an entry to the context to prevent recalculating request role unnecessarily
		r = r.WithContext(context.WithValue(r.Context(), logical.CtxKeyRequestRole{}, role))
		quotaReq.Role = role

		if hitRateLimitQuota(core, r, quotaReq, w) {
			return
		}

		handler.ServeHTTP(w, r)
	})
}

// rateLimitQuotaWrapping performs quota checking for requests that do not
// require role resolution. If the request does require role resolution, the
// quota request is added to the context and handled by withRoleRateLimitQuotaWrapping
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

		// If any role-based quotas are enabled for this namespace/mount, wait until after the
		// json limits checks to perform it
		if requiresResolveRole {
			// Add a context entry to indicate that role-based quota resolution
			// is needed downstream
			r = r.WithContext(context.WithValue(r.Context(), ctxKeyRoleBasedQuota{}, quotaReq))
			handler.ServeHTTP(w, r)
			return
		}

		if hitRateLimitQuota(core, r, quotaReq, w) {
			return
		}

		handler.ServeHTTP(w, r)
		return
	})
}

// hitRateLimitQuota checks the applies the rate limit quota and handles writing
// headers, as well as any auditing and logging if the quota is exceeded. The
// function returns true if the quota was exceeded and the request should not
// be processed further.
func hitRateLimitQuota(core *vault.Core, r *http.Request, quotaReq *quotas.Request, w http.ResponseWriter) bool {
	path := quotaReq.Path
	quotaResp, err := core.ApplyRateLimitQuota(r.Context(), quotaReq)
	if err != nil {
		core.Logger().Error("failed to apply quota", "path", path, "error", err)
		respondError(w, http.StatusInternalServerError, err)
		return true
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
				return true
			}

			err = core.AuditLogger().AuditRequest(r.Context(), &logical.LogInput{
				Request:  req,
				OuterErr: quotaErr,
			})
			if err != nil {
				core.Logger().Warn("failed to audit log request rejection caused by rate limit quota violation", "error", err)
			}
		}

		return true
	}
	return false
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
