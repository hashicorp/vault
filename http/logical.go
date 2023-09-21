// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"go.uber.org/atomic"
)

// bufferedReader can be used to replace a request body with a buffered
// version. The Close method invokes the original Closer.
type bufferedReader struct {
	*bufio.Reader
	rOrig io.ReadCloser
}

func newBufferedReader(r io.ReadCloser) *bufferedReader {
	return &bufferedReader{
		Reader: bufio.NewReader(r),
		rOrig:  r,
	}
}

func (b *bufferedReader) Close() error {
	return b.rOrig.Close()
}

const MergePatchContentTypeHeader = "application/merge-patch+json"

func buildLogicalRequestNoAuth(perfStandby bool, w http.ResponseWriter, r *http.Request) (*logical.Request, io.ReadCloser, int, error) {
	ns, err := namespace.FromContext(r.Context())
	if err != nil {
		return nil, nil, http.StatusBadRequest, nil
	}
	path := ns.TrimmedPath(r.URL.Path[len("/v1/"):])

	var data map[string]interface{}
	var origBody io.ReadCloser
	var passHTTPReq bool
	var responseWriter http.ResponseWriter

	// Determine the operation
	var op logical.Operation
	switch r.Method {
	case "DELETE":
		op = logical.DeleteOperation
		data = parseQuery(r.URL.Query())
	case "GET":
		op = logical.ReadOperation
		queryVals := r.URL.Query()
		var list bool
		var err error
		listStr := queryVals.Get("list")
		if listStr != "" {
			list, err = strconv.ParseBool(listStr)
			if err != nil {
				return nil, nil, http.StatusBadRequest, nil
			}
			if list {
				queryVals.Del("list")
				op = logical.ListOperation
				if !strings.HasSuffix(path, "/") {
					path += "/"
				}
			}
		}

		data = parseQuery(queryVals)

		switch {
		case strings.HasPrefix(path, "sys/pprof/"):
			passHTTPReq = true
			responseWriter = w
		case path == "sys/storage/raft/snapshot":
			responseWriter = w
		case path == "sys/internal/counters/activity/export":
			responseWriter = w
		case path == "sys/monitor":
			passHTTPReq = true
			responseWriter = w
		}

	case "POST", "PUT":
		op = logical.UpdateOperation

		// Buffer the request body in order to allow us to peek at the beginning
		// without consuming it. This approach involves no copying.
		bufferedBody := newBufferedReader(r.Body)
		r.Body = bufferedBody

		// If we are uploading a snapshot or receiving an ocsp-request (which
		// is der encoded) we don't want to parse it. Instead, we will simply
		// add the HTTP request to the logical request object for later consumption.
		contentType := r.Header.Get("Content-Type")
		if path == "sys/storage/raft/snapshot" || path == "sys/storage/raft/snapshot-force" || isOcspRequest(contentType) {
			passHTTPReq = true
			origBody = r.Body
		} else {
			// Sample the first bytes to determine whether this should be parsed as
			// a form or as JSON. The amount to look ahead (512 bytes) is arbitrary
			// but extremely tolerant (i.e. allowing 511 bytes of leading whitespace
			// and an incorrect content-type).
			head, err := bufferedBody.Peek(512)
			if err != nil && err != bufio.ErrBufferFull && err != io.EOF {
				status := http.StatusBadRequest
				logical.AdjustErrorStatusCode(&status, err)
				return nil, nil, status, fmt.Errorf("error reading data")
			}

			if isForm(head, contentType) {
				formData, err := parseFormRequest(r)
				if err != nil {
					status := http.StatusBadRequest
					logical.AdjustErrorStatusCode(&status, err)
					return nil, nil, status, fmt.Errorf("error parsing form data")
				}

				data = formData
			} else {
				origBody, err = parseJSONRequest(perfStandby, r, w, &data)
				if err == io.EOF {
					data = nil
					err = nil
				}
				if err != nil {
					status := http.StatusBadRequest
					logical.AdjustErrorStatusCode(&status, err)
					return nil, nil, status, fmt.Errorf("error parsing JSON")
				}
			}
		}

	case "PATCH":
		op = logical.PatchOperation

		contentTypeHeader := r.Header.Get("Content-Type")
		contentType, _, err := mime.ParseMediaType(contentTypeHeader)
		if err != nil {
			status := http.StatusBadRequest
			logical.AdjustErrorStatusCode(&status, err)
			return nil, nil, status, err
		}

		if contentType != MergePatchContentTypeHeader {
			return nil, nil, http.StatusUnsupportedMediaType, fmt.Errorf("PATCH requires Content-Type of %s, provided %s", MergePatchContentTypeHeader, contentType)
		}

		origBody, err = parseJSONRequest(perfStandby, r, w, &data)

		if err == io.EOF {
			data = nil
			err = nil
		}

		if err != nil {
			status := http.StatusBadRequest
			logical.AdjustErrorStatusCode(&status, err)
			return nil, nil, status, fmt.Errorf("error parsing JSON")
		}

	case "LIST":
		op = logical.ListOperation
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}

		data = parseQuery(r.URL.Query())
	case "HEAD":
		op = logical.HeaderOperation
		data = parseQuery(r.URL.Query())
	case "OPTIONS":
	default:
		return nil, nil, http.StatusMethodNotAllowed, nil
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to generate identifier for the request: %w", err)
	}

	req := &logical.Request{
		ID:         requestId,
		Operation:  op,
		Path:       path,
		Data:       data,
		Connection: getConnection(r),
		Headers:    r.Header,
	}

	if passHTTPReq {
		req.HTTPRequest = r
	}
	if responseWriter != nil {
		req.ResponseWriter = logical.NewHTTPResponseWriter(responseWriter)
	}

	return req, origBody, 0, nil
}

func isOcspRequest(contentType string) bool {
	contentType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	return contentType == "application/ocsp-request"
}

func buildLogicalPath(r *http.Request) (string, int, error) {
	ns, err := namespace.FromContext(r.Context())
	if err != nil {
		return "", http.StatusBadRequest, nil
	}

	path := ns.TrimmedPath(strings.TrimPrefix(r.URL.Path, "/v1/"))

	switch r.Method {
	case "GET":
		var (
			list bool
			err  error
		)

		queryVals := r.URL.Query()

		listStr := queryVals.Get("list")
		if listStr != "" {
			list, err = strconv.ParseBool(listStr)
			if err != nil {
				return "", http.StatusBadRequest, nil
			}
			if list {
				if !strings.HasSuffix(path, "/") {
					path += "/"
				}
			}
		}

	case "LIST":
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
	}

	return path, 0, nil
}

func buildLogicalRequest(core *vault.Core, w http.ResponseWriter, r *http.Request) (*logical.Request, io.ReadCloser, int, error) {
	req, origBody, status, err := buildLogicalRequestNoAuth(core.PerfStandby(), w, r)
	if err != nil || status != 0 {
		return nil, nil, status, err
	}

	req.SetRequiredState(r.Header.Values(VaultIndexHeaderName))
	requestAuth(r, req)

	req, err = requestWrapInfo(r, req)
	if err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("error parsing X-Vault-Wrap-TTL header: %w", err)
	}

	err = parseMFAHeader(req)
	if err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("failed to parse X-Vault-MFA header: %w", err)
	}

	err = requestPolicyOverride(r, req)
	if err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("failed to parse %s header: %w", PolicyOverrideHeaderName, err)
	}

	return req, origBody, 0, nil
}

// handleLogical returns a handler for processing logical requests. These requests
// may or may not end up getting forwarded under certain scenarios if the node
// is a performance standby. Some of these cases include:
//   - Perf standby and token with limited use count.
//   - Perf standby and token re-validation needed (e.g. due to invalid token).
//   - Perf standby and control group error.
func handleLogical(core *vault.Core) http.Handler {
	return handleLogicalInternal(core, false, false)
}

// handleLogicalWithInjector returns a handler for processing logical requests
// that also have their logical response data injected at the top-level payload.
// All forwarding behavior remains the same as `handleLogical`.
func handleLogicalWithInjector(core *vault.Core) http.Handler {
	return handleLogicalInternal(core, true, false)
}

// handleLogicalNoForward returns a handler for processing logical local-only
// requests. These types of requests never forwarded, and return an
// `vault.ErrCannotForwardLocalOnly` error if attempted to do so.
func handleLogicalNoForward(core *vault.Core) http.Handler {
	return handleLogicalInternal(core, false, true)
}

func handleLogicalRecovery(raw *vault.RawBackend, token *atomic.String) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _, statusCode, err := buildLogicalRequestNoAuth(false, w, r)
		if err != nil || statusCode != 0 {
			respondError(w, statusCode, err)
			return
		}
		reqToken := r.Header.Get(consts.AuthHeaderName)
		if reqToken == "" || token.Load() == "" || reqToken != token.Load() {
			respondError(w, http.StatusForbidden, nil)
			return
		}

		resp, err := raw.HandleRequest(r.Context(), req)
		if respondErrorCommon(w, req, resp, err) {
			return
		}

		var httpResp *logical.HTTPResponse
		if resp != nil {
			httpResp = logical.LogicalResponseToHTTPResponse(resp)
			httpResp.RequestID = req.ID
		}
		respondOk(w, httpResp)
	})
}

// handleLogicalInternal is a common helper that returns a handler for
// processing logical requests. The behavior depends on the various boolean
// toggles. Refer to usage on functions for possible behaviors.
func handleLogicalInternal(core *vault.Core, injectDataIntoTopLevel bool, noForward bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, origBody, statusCode, err := buildLogicalRequest(core, w, r)
		if err != nil || statusCode != 0 {
			respondError(w, statusCode, err)
			return
		}

		// Websockets need to be handled at HTTP layer instead of logical requests.
		ns, err := namespace.FromContext(r.Context())
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		nsPath := ns.Path
		if ns.ID == namespace.RootNamespaceID {
			nsPath = ""
		}
		if strings.HasPrefix(r.URL.Path, fmt.Sprintf("/v1/%ssys/events/subscribe/", nsPath)) {
			handler := handleEventsSubscribe(core, req)
			handler.ServeHTTP(w, r)
			return
		}

		// Make the internal request. We attach the connection info
		// as well in case this is an authentication request that requires
		// it. Vault core handles stripping this if we need to. This also
		// handles all error cases; if we hit respondLogical, the request is a
		// success.
		resp, ok, needsForward := request(core, w, r, req)
		switch {
		case needsForward && noForward:
			respondError(w, http.StatusBadRequest, vault.ErrCannotForwardLocalOnly)
			return
		case needsForward && !noForward:
			if origBody != nil {
				r.Body = origBody
			}
			forwardRequest(core, w, r)
			return
		case !ok:
			// If not ok, we simply return. The call on request should have
			// taken care of setting the appropriate response code and payload
			// in this case.
			return
		default:
			// Build and return the proper response if everything is fine.
			respondLogical(core, w, r, req, resp, injectDataIntoTopLevel)
			return
		}
	})
}

func respondLogical(core *vault.Core, w http.ResponseWriter, r *http.Request, req *logical.Request, resp *logical.Response, injectDataIntoTopLevel bool) {
	var httpResp *logical.HTTPResponse
	var ret interface{}

	// If vault's core has already written to the response writer do not add any
	// additional output. Headers have already been sent.
	if req != nil && req.ResponseWriter != nil && req.ResponseWriter.Written() {
		return
	}

	if resp != nil {
		if resp.Redirect != "" {
			// If we have a redirect, redirect! We use a 307 code
			// because we don't actually know if its permanent.
			http.Redirect(w, r, resp.Redirect, 307)
			return
		}

		// Check if this is a raw response
		if _, ok := resp.Data[logical.HTTPStatusCode]; ok {
			respondRaw(w, r, resp)
			return
		}

		if resp.WrapInfo != nil && resp.WrapInfo.Token != "" {
			httpResp = &logical.HTTPResponse{
				WrapInfo: &logical.HTTPWrapInfo{
					Token:           resp.WrapInfo.Token,
					Accessor:        resp.WrapInfo.Accessor,
					TTL:             int(resp.WrapInfo.TTL.Seconds()),
					CreationTime:    resp.WrapInfo.CreationTime.Format(time.RFC3339Nano),
					CreationPath:    resp.WrapInfo.CreationPath,
					WrappedAccessor: resp.WrapInfo.WrappedAccessor,
				},
			}
		} else {
			httpResp = logical.LogicalResponseToHTTPResponse(resp)
			httpResp.RequestID = req.ID
		}

		ret = httpResp

		if injectDataIntoTopLevel {
			injector := logical.HTTPSysInjector{
				Response: httpResp,
			}
			ret = injector
		}
	}

	adjustResponse(core, w, req)

	// Respond
	respondOk(w, ret)
	return
}

// respondRaw is used when the response is using HTTPContentType and HTTPRawBody
// to change the default response handling. This is only used for specific things like
// returning the CRL information on the PKI backends.
func respondRaw(w http.ResponseWriter, r *http.Request, resp *logical.Response) {
	retErr := func(w http.ResponseWriter, err string) {
		w.Header().Set("X-Vault-Raw-Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}

	// Ensure this is never a secret or auth response
	if resp.Secret != nil || resp.Auth != nil {
		retErr(w, "raw responses cannot contain secrets or auth")
		return
	}

	// Get the status code
	statusRaw, ok := resp.Data[logical.HTTPStatusCode]
	if !ok {
		retErr(w, "no status code given")
		return
	}

	var status int
	switch statusRaw.(type) {
	case int:
		status = statusRaw.(int)
	case float64:
		status = int(statusRaw.(float64))
	case json.Number:
		s64, err := statusRaw.(json.Number).Float64()
		if err != nil {
			retErr(w, "cannot decode status code")
			return
		}
		status = int(s64)
	default:
		retErr(w, "cannot decode status code")
		return
	}

	nonEmpty := status != http.StatusNoContent

	var contentType string
	var body []byte

	// Get the content type header; don't require it if the body is empty
	contentTypeRaw, ok := resp.Data[logical.HTTPContentType]
	if !ok && nonEmpty {
		retErr(w, "no content type given")
		return
	}
	if ok {
		contentType, ok = contentTypeRaw.(string)
		if !ok {
			retErr(w, "cannot decode content type")
			return
		}
	}

	if nonEmpty {
		// Get the body
		bodyRaw, ok := resp.Data[logical.HTTPRawBody]
		if !ok {
			goto WRITE_RESPONSE
		}

		switch bodyRaw.(type) {
		case string:
			// This is best effort. The value may already be base64-decoded so
			// if it doesn't work we just use as-is
			bodyDec, err := base64.StdEncoding.DecodeString(bodyRaw.(string))
			if err == nil {
				body = bodyDec
			} else {
				body = []byte(bodyRaw.(string))
			}
		case []byte:
			body = bodyRaw.([]byte)
		default:
			retErr(w, "cannot decode body")
			return
		}
	}

WRITE_RESPONSE:
	// Write the response
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	if cacheControl, ok := resp.Data[logical.HTTPCacheControlHeader].(string); ok {
		w.Header().Set("Cache-Control", cacheControl)
	}

	if pragma, ok := resp.Data[logical.HTTPPragmaHeader].(string); ok {
		w.Header().Set("Pragma", pragma)
	}

	if wwwAuthn, ok := resp.Data[logical.HTTPWWWAuthenticateHeader].(string); ok {
		w.Header().Set("WWW-Authenticate", wwwAuthn)
	}

	w.WriteHeader(status)
	w.Write(body)
}

// getConnection is used to format the connection information for
// attaching to a logical request
func getConnection(r *http.Request) (connection *logical.Connection) {
	var remoteAddr string
	var remotePort int

	remoteAddr, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteAddr = ""
	} else {
		remotePort, err = strconv.Atoi(port)
		if err != nil {
			remotePort = 0
		}
	}

	connection = &logical.Connection{
		RemoteAddr: remoteAddr,
		RemotePort: remotePort,
		ConnState:  r.TLS,
	}
	return
}
