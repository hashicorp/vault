package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/pathmanager"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	// WrapTTLHeaderName is the name of the header containing a directive to
	// wrap the response
	WrapTTLHeaderName = "X-Vault-Wrap-TTL"

	// WrapFormatHeaderName is the name of the header containing the format to
	// wrap in; has no effect if the wrap TTL is not set
	WrapFormatHeaderName = "X-Vault-Wrap-Format"

	// NoRequestForwardingHeaderName is the name of the header telling Vault
	// not to use request forwarding
	NoRequestForwardingHeaderName = "X-Vault-No-Request-Forwarding"

	// MFAHeaderName represents the HTTP header which carries the credentials
	// required to perform MFA on any path.
	MFAHeaderName = "X-Vault-MFA"

	// canonicalMFAHeaderName is the MFA header value's format in the request
	// headers. Do not alter the casing of this string.
	canonicalMFAHeaderName = "X-Vault-Mfa"

	// PolicyOverrideHeaderName is the header set to request overriding
	// soft-mandatory Sentinel policies.
	PolicyOverrideHeaderName = "X-Vault-Policy-Override"

	// DefaultMaxRequestSize is the default maximum accepted request size. This
	// is to prevent a denial of service attack where no Content-Length is
	// provided and the server is fed ever more data until it exhausts memory.
	// Can be overridden per listener.
	DefaultMaxRequestSize = 32 * 1024 * 1024
)

var (
	// Set to false by stub_asset if the ui build tag isn't enabled
	uiBuiltIn = true

	// perfStandbyAlwaysForwardPaths is used to check a requested path against
	// the always forward list
	perfStandbyAlwaysForwardPaths = pathmanager.New()

	injectDataIntoTopRoutes = []string{
		"/v1/sys/audit",
		"/v1/sys/audit/",
		"/v1/sys/audit-hash/",
		"/v1/sys/auth",
		"/v1/sys/auth/",
		"/v1/sys/config/cors",
		"/v1/sys/config/auditing/request-headers/",
		"/v1/sys/config/auditing/request-headers",
		"/v1/sys/capabilities",
		"/v1/sys/capabilities-accessor",
		"/v1/sys/capabilities-self",
		"/v1/sys/key-status",
		"/v1/sys/mounts",
		"/v1/sys/mounts/",
		"/v1/sys/policy",
		"/v1/sys/policy/",
		"/v1/sys/rekey/backup",
		"/v1/sys/rekey/recovery-key-backup",
		"/v1/sys/remount",
		"/v1/sys/rotate",
		"/v1/sys/wrapping/wrap",
	}
)

// Handler returns an http.Handler for the API. This can be used on
// its own to mount the Vault API within another web server.
func Handler(props *vault.HandlerProperties) http.Handler {
	core := props.Core

	// Create the muxer to handle the actual endpoints
	mux := http.NewServeMux()
	mux.Handle("/v1/sys/init", handleSysInit(core))
	mux.Handle("/v1/sys/seal-status", handleSysSealStatus(core))
	mux.Handle("/v1/sys/seal", handleSysSeal(core))
	mux.Handle("/v1/sys/step-down", handleRequestForwarding(core, handleSysStepDown(core)))
	mux.Handle("/v1/sys/unseal", handleSysUnseal(core))
	mux.Handle("/v1/sys/leader", handleSysLeader(core))
	mux.Handle("/v1/sys/health", handleSysHealth(core))
	mux.Handle("/v1/sys/generate-root/attempt", handleRequestForwarding(core, handleSysGenerateRootAttempt(core, vault.GenerateStandardRootTokenStrategy)))
	mux.Handle("/v1/sys/generate-root/update", handleRequestForwarding(core, handleSysGenerateRootUpdate(core, vault.GenerateStandardRootTokenStrategy)))
	mux.Handle("/v1/sys/rekey/init", handleRequestForwarding(core, handleSysRekeyInit(core, false)))
	mux.Handle("/v1/sys/rekey/update", handleRequestForwarding(core, handleSysRekeyUpdate(core, false)))
	mux.Handle("/v1/sys/rekey/verify", handleRequestForwarding(core, handleSysRekeyVerify(core, false)))
	mux.Handle("/v1/sys/rekey-recovery-key/init", handleRequestForwarding(core, handleSysRekeyInit(core, true)))
	mux.Handle("/v1/sys/rekey-recovery-key/update", handleRequestForwarding(core, handleSysRekeyUpdate(core, true)))
	mux.Handle("/v1/sys/rekey-recovery-key/verify", handleRequestForwarding(core, handleSysRekeyVerify(core, true)))
	for _, path := range injectDataIntoTopRoutes {
		mux.Handle(path, handleRequestForwarding(core, handleLogicalWithInjector(core)))
	}
	mux.Handle("/v1/sys/", handleRequestForwarding(core, handleLogical(core)))
	mux.Handle("/v1/", handleRequestForwarding(core, handleLogical(core)))
	if core.UIEnabled() == true {
		if uiBuiltIn {
			mux.Handle("/ui/", http.StripPrefix("/ui/", gziphandler.GzipHandler(handleUIHeaders(core, handleUI(http.FileServer(&UIAssetWrapper{FileSystem: assetFS()}))))))
			mux.Handle("/robots.txt", gziphandler.GzipHandler(handleUIHeaders(core, handleUI(http.FileServer(&UIAssetWrapper{FileSystem: assetFS()})))))
		} else {
			mux.Handle("/ui/", handleUIHeaders(core, handleUIStub()))
		}
		mux.Handle("/ui", handleUIRedirect())
		mux.Handle("/", handleUIRedirect())
	}

	additionalRoutes(mux, core)

	// Wrap the handler in another handler to trigger all help paths.
	helpWrappedHandler := wrapHelpHandler(mux, core)
	corsWrappedHandler := wrapCORSHandler(helpWrappedHandler, core)

	genericWrappedHandler := genericWrapping(core, corsWrappedHandler, props)

	// Wrap the handler with PrintablePathCheckHandler to check for non-printable
	// characters in the request path.
	printablePathCheckHandler := genericWrappedHandler
	if !props.DisablePrintableCheck {
		printablePathCheckHandler = cleanhttp.PrintablePathCheckHandler(genericWrappedHandler, nil)
	}

	return printablePathCheckHandler
}

// wrapGenericHandler wraps the handler with an extra layer of handler where
// tasks that should be commonly handled for all the requests and/or responses
// are performed.
func wrapGenericHandler(core *vault.Core, h http.Handler, maxRequestSize int64, maxRequestDuration time.Duration) http.Handler {
	if maxRequestDuration == 0 {
		maxRequestDuration = vault.DefaultMaxRequestDuration
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Cache-Control header for all the responses returned
		// by Vault
		w.Header().Set("Cache-Control", "no-store")

		// Start with the request context
		ctx := r.Context()
		var cancelFunc context.CancelFunc
		// Add our timeout
		ctx, cancelFunc = context.WithTimeout(ctx, maxRequestDuration)
		// Add a size limiter if desired
		if maxRequestSize > 0 {
			ctx = context.WithValue(ctx, "max_request_size", maxRequestSize)
		}
		ctx = context.WithValue(ctx, "original_request_path", r.URL.Path)
		r = r.WithContext(ctx)

		switch {
		case strings.HasPrefix(r.URL.Path, "/v1/"):
			newR, status := adjustRequest(core, r)
			if status != 0 {
				respondError(w, status, nil)
				cancelFunc()
				return
			}
			r = newR

		case strings.HasPrefix(r.URL.Path, "/ui"), r.URL.Path == "/robots.txt", r.URL.Path == "/":
		default:
			respondError(w, http.StatusNotFound, nil)
			cancelFunc()
			return
		}

		h.ServeHTTP(w, r)
		cancelFunc()
		return
	})
}

func WrapForwardedForHandler(h http.Handler, authorizedAddrs []*sockaddr.SockAddrMarshaler, rejectNotPresent, rejectNonAuthz bool, hopSkips int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers, headersOK := r.Header[textproto.CanonicalMIMEHeaderKey("X-Forwarded-For")]
		if !headersOK || len(headers) == 0 {
			if !rejectNotPresent {
				h.ServeHTTP(w, r)
				return
			}
			respondError(w, http.StatusBadRequest, fmt.Errorf("missing x-forwarded-for header and configured to reject when not present"))
			return
		}

		host, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// If not rejecting treat it like we just don't have a valid
			// header because we can't do a comparison against an address we
			// can't understand
			if !rejectNotPresent {
				h.ServeHTTP(w, r)
				return
			}
			respondError(w, http.StatusBadRequest, errwrap.Wrapf("error parsing client hostport: {{err}}", err))
			return
		}

		addr, err := sockaddr.NewIPAddr(host)
		if err != nil {
			// We treat this the same as the case above
			if !rejectNotPresent {
				h.ServeHTTP(w, r)
				return
			}
			respondError(w, http.StatusBadRequest, errwrap.Wrapf("error parsing client address: {{err}}", err))
			return
		}

		var found bool
		for _, authz := range authorizedAddrs {
			if authz.Contains(addr) {
				found = true
				break
			}
		}
		if !found {
			// If we didn't find it and aren't configured to reject, simply
			// don't trust it
			if !rejectNonAuthz {
				h.ServeHTTP(w, r)
				return
			}
			respondError(w, http.StatusBadRequest, fmt.Errorf("client address not authorized for x-forwarded-for and configured to reject connection"))
			return
		}

		// At this point we have at least one value and it's authorized

		// Split comma separated ones, which are common. This brings it in line
		// to the multiple-header case.
		var acc []string
		for _, header := range headers {
			vals := strings.Split(header, ",")
			for _, v := range vals {
				acc = append(acc, strings.TrimSpace(v))
			}
		}

		indexToUse := len(acc) - 1 - hopSkips
		if indexToUse < 0 {
			// This is likely an error in either configuration or other
			// infrastructure. We could either deny the request, or we
			// could simply not trust the value. Denying the request is
			// "safer" since if this logic is configured at all there may
			// be an assumption it can always be trusted. Given that we can
			// deny accepting the request at all if it's not from an
			// authorized address, if we're at this point the address is
			// authorized (or we've turned off explicit rejection) and we
			// should assume that what comes in should be properly
			// formatted.
			respondError(w, http.StatusBadRequest, fmt.Errorf("malformed x-forwarded-for configuration or request, hops to skip (%d) would skip before earliest chain link (chain length %d)", hopSkips, len(headers)))
			return
		}

		r.RemoteAddr = net.JoinHostPort(acc[indexToUse], port)
		h.ServeHTTP(w, r)
		return
	})
}

// A lookup on a token that is about to expire returns nil, which means by the
// time we can validate a wrapping token lookup will return nil since it will
// be revoked after the call. So we have to do the validation here.
func wrappingVerificationFunc(ctx context.Context, core *vault.Core, req *logical.Request) error {
	if req == nil {
		return fmt.Errorf("invalid request")
	}

	valid, err := core.ValidateWrappingToken(ctx, req)
	if err != nil {
		return errwrap.Wrapf("error validating wrapping token: {{err}}", err)
	}
	if !valid {
		return fmt.Errorf("wrapping token is not valid or does not exist")
	}

	return nil
}

// stripPrefix is a helper to strip a prefix from the path. It will
// return false from the second return value if it the prefix doesn't exist.
func stripPrefix(prefix, path string) (string, bool) {
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}

	path = path[len(prefix):]
	if path == "" {
		return "", false
	}

	return path, true
}

func handleUIHeaders(core *vault.Core, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		header := w.Header()

		userHeaders, err := core.UIHeaders()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if userHeaders != nil {
			for k := range userHeaders {
				v := userHeaders.Get(k)
				header.Set(k, v)
			}
		}
		h.ServeHTTP(w, req)
	})
}

func handleUI(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// The fileserver handler strips trailing slashes and does a redirect.
		// We don't want the redirect to happen so we preemptively trim the slash
		// here.
		req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
		h.ServeHTTP(w, req)
		return
	})
}

func handleUIStub() http.Handler {
	stubHTML := `
	<!DOCTYPE html>
	<html>
	<style>
	body {
	color: #1F2124;
	font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
	}

	.wrapper {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 500px;
	}

	.content ul {
	line-height: 1.5;
	}

	a {
	color: #1563ff;
	text-decoration: none;
	}

	.header {
	display: flex;
	color: #6a7786;
	align-items: center;
	}

	.header svg {
	padding-right: 12px;
	}

	.alert {
	transform: scale(0.07);
	fill: #6a7786;
	}

	h1 {
	font-weight: 500;
	}

	p {
	margin-top: 0px;
	}
	</style>
	<div class="wrapper">
	<div class="content">
	<div class="header">
	<svg width="36px" height="36px" viewBox="0 0 36 36" xmlns="http://www.w3.org/2000/svg">
	<path class="alert" d="M476.7 422.2L270.1 72.7c-2.9-5-8.3-8.7-14.1-8.7-5.9 0-11.3 3.7-14.1 8.7L35.3 422.2c-2.8 5-4.8 13-1.9 17.9 2.9 4.9 8.2 7.9 14 7.9h417.1c5.8 0 11.1-3 14-7.9 3-4.9 1-13-1.8-17.9zM288 400h-64v-48h64v48zm0-80h-64V176h64v144z"/>
	</svg>
	<h1>Vault UI is not available in this binary.</h1>
	</div>
	<p>To get Vault UI do one of the following:</p>
	<ul>
	<li><a href="https://www.vaultproject.io/downloads.html">Download an official release</a></li>
	<li>Run <code>make release</code> to create your own release binaries.
	<li>Run <code>make dev-ui</code> to create a development binary with the UI.
	</ul>
	</div>
	</div>
	</html>
	`
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(stubHTML))
	})
}

func handleUIRedirect() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/ui/", 307)
		return
	})
}

type UIAssetWrapper struct {
	FileSystem *assetfs.AssetFS
}

func (fs *UIAssetWrapper) Open(name string) (http.File, error) {
	file, err := fs.FileSystem.Open(name)
	if err == nil {
		return file, nil
	}
	// serve index.html instead of 404ing
	if err == os.ErrNotExist {
		return fs.FileSystem.Open("index.html")
	}
	return nil, err
}

func parseRequest(r *http.Request, w http.ResponseWriter, out interface{}) error {
	// Limit the maximum number of bytes to MaxRequestSize to protect
	// against an indefinite amount of data being read.
	reader := r.Body
	ctx := r.Context()
	maxRequestSize := ctx.Value("max_request_size")
	if maxRequestSize != nil {
		max, ok := maxRequestSize.(int64)
		if !ok {
			return errors.New("could not parse max_request_size from request context")
		}
		if max > 0 {
			reader = http.MaxBytesReader(w, r.Body, max)
		}
	}
	err := jsonutil.DecodeJSONFromReader(reader, out)
	if err != nil && err != io.EOF {
		return errwrap.Wrapf("failed to parse JSON input: {{err}}", err)
	}
	return err
}

// handleRequestForwarding determines whether to forward a request or not,
// falling back on the older behavior of redirecting the client
func handleRequestForwarding(core *vault.Core, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If we are a performance standby we can handle the request.
		if core.PerfStandby() {
			ns, err := namespace.FromContext(r.Context())
			if err != nil {
				respondError(w, http.StatusBadRequest, err)
				return
			}
			path := ns.TrimmedPath(r.URL.Path[len("/v1/"):])
			switch {
			case !perfStandbyAlwaysForwardPaths.HasPath(path):
				handler.ServeHTTP(w, r)
				return
			case strings.HasPrefix(path, "auth/token/create/"):
				isBatch, err := core.IsBatchTokenCreationRequest(r.Context(), path)
				if err == nil && isBatch {
					handler.ServeHTTP(w, r)
					return
				}
			}
		}

		// Note: in an HA setup, this call will also ensure that connections to
		// the leader are set up, as that happens once the advertised cluster
		// values are read during this function
		isLeader, leaderAddr, _, err := core.Leader()
		if err != nil {
			if err == vault.ErrHANotEnabled {
				// Standalone node, serve request normally
				handler.ServeHTTP(w, r)
				return
			}
			// Some internal error occurred
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if isLeader {
			// No forwarding needed, we're leader
			handler.ServeHTTP(w, r)
			return
		}
		if leaderAddr == "" {
			respondError(w, http.StatusInternalServerError, fmt.Errorf("local node not active but active cluster node not found"))
			return
		}

		forwardRequest(core, w, r)
		return
	})
}

func forwardRequest(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(vault.IntNoForwardingHeaderName) != "" {
		respondStandby(core, w, r.URL)
		return
	}

	if r.Header.Get(NoRequestForwardingHeaderName) != "" {
		// Forwarding explicitly disabled, fall back to previous behavior
		core.Logger().Debug("handleRequestForwarding: forwarding disabled by client request")
		respondStandby(core, w, r.URL)
		return
	}

	// Attempt forwarding the request. If we cannot forward -- perhaps it's
	// been disabled on the active node -- this will return with an
	// ErrCannotForward and we simply fall back
	statusCode, header, retBytes, err := core.ForwardRequest(r)
	if err != nil {
		if err == vault.ErrCannotForward {
			core.Logger().Debug("cannot forward request (possibly disabled on active node), falling back")
		} else {
			core.Logger().Error("forward request error", "error", err)
		}

		// Fall back to redirection
		respondStandby(core, w, r.URL)
		return
	}

	if header != nil {
		for k, v := range header {
			w.Header()[k] = v
		}
	}

	w.WriteHeader(statusCode)
	w.Write(retBytes)
}

// request is a helper to perform a request and properly exit in the
// case of an error.
func request(core *vault.Core, w http.ResponseWriter, rawReq *http.Request, r *logical.Request) (*logical.Response, bool, bool) {
	resp, err := core.HandleRequest(rawReq.Context(), r)
	if r.LastRemoteWAL() > 0 && !vault.WaitUntilWALShipped(rawReq.Context(), core, r.LastRemoteWAL()) {
		if resp == nil {
			resp = &logical.Response{}
		}
		resp.AddWarning("Timeout hit while waiting for local replicated cluster to apply primary's write; this client may encounter stale reads of values written during this operation.")
	}
	if errwrap.Contains(err, consts.ErrStandby.Error()) {
		respondStandby(core, w, rawReq.URL)
		return resp, false, false
	}
	if err != nil && errwrap.Contains(err, logical.ErrPerfStandbyPleaseForward.Error()) {
		return nil, false, true
	}

	if resp != nil && len(resp.Headers) > 0 {
		// Set this here so it will take effect regardless of any other type of
		// response processing
		header := w.Header()
		for k, v := range resp.Headers {
			for _, h := range v {
				header.Add(k, h)
			}
		}

		switch {
		case resp.Secret != nil,
			resp.Auth != nil,
			len(resp.Data) > 0,
			resp.Redirect != "",
			len(resp.Warnings) > 0,
			resp.WrapInfo != nil:
			// Nothing, resp has data

		default:
			// We have an otherwise totally empty response except for headers,
			// so nil out the response now that the headers are written out
			resp = nil
		}
	}

	if respondErrorCommon(w, r, resp, err) {
		return resp, false, false
	}

	return resp, true, false
}

// respondStandby is used to trigger a redirect in the case that this Vault is currently a hot standby
func respondStandby(core *vault.Core, w http.ResponseWriter, reqURL *url.URL) {
	// Request the leader address
	_, redirectAddr, _, err := core.Leader()
	if err != nil {
		if err == vault.ErrHANotEnabled {
			// Standalone node, serve 503
			err = errors.New("node is not active")
			respondError(w, http.StatusServiceUnavailable, err)
			return
		}

		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// If there is no leader, generate a 503 error
	if redirectAddr == "" {
		err = errors.New("no active Vault instance found")
		respondError(w, http.StatusServiceUnavailable, err)
		return
	}

	// Parse the redirect location
	redirectURL, err := url.Parse(redirectAddr)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Generate a redirect URL
	finalURL := url.URL{
		Scheme:   redirectURL.Scheme,
		Host:     redirectURL.Host,
		Path:     reqURL.Path,
		RawQuery: reqURL.RawQuery,
	}

	// Ensure there is a scheme, default to https
	if finalURL.Scheme == "" {
		finalURL.Scheme = "https"
	}

	// If we have an address, redirect! We use a 307 code
	// because we don't actually know if its permanent and
	// the request method should be preserved.
	w.Header().Set("Location", finalURL.String())
	w.WriteHeader(307)
}

// getTokenFromReq parse headers of the incoming request to extract token if
// present it accepts Authorization Bearer (RFC6750) and X-Vault-Token header.
// Returns true if the token was sourced from a Bearer header.
func getTokenFromReq(r *http.Request) (string, bool) {
	if token := r.Header.Get(consts.AuthHeaderName); token != "" {
		return token, false
	}
	if headers, ok := r.Header["Authorization"]; ok {
		// Reference for Authorization header format: https://tools.ietf.org/html/rfc7236#section-3

		// If string does not start by 'Bearer ', it is not one we would use,
		// but might be used by plugins
		for _, v := range headers {
			if !strings.HasPrefix(v, "Bearer ") {
				continue
			}
			return strings.TrimSpace(v[7:]), true
		}
	}
	return "", false
}

// requestAuth adds the token to the logical.Request if it exists.
func requestAuth(core *vault.Core, r *http.Request, req *logical.Request) (*logical.Request, error) {
	// Attach the header value if we have it
	token, fromAuthzHeader := getTokenFromReq(r)
	if token != "" {
		req.ClientToken = token
		req.ClientTokenSource = logical.ClientTokenFromVaultHeader
		if fromAuthzHeader {
			req.ClientTokenSource = logical.ClientTokenFromAuthzHeader
		}

		// Also attach the accessor if we have it. This doesn't fail if it
		// doesn't exist because the request may be to an unauthenticated
		// endpoint/login endpoint where a bad current token doesn't matter, or
		// a token from a Vault version pre-accessors.
		te, err := core.LookupToken(r.Context(), token)
		if err != nil && strings.Count(token, ".") != 2 {
			return req, err
		}
		if err == nil && te != nil {
			req.ClientTokenAccessor = te.Accessor
			req.ClientTokenRemainingUses = te.NumUses
			req.SetTokenEntry(te)
		}
	}

	return req, nil
}

func requestPolicyOverride(r *http.Request, req *logical.Request) error {
	raw := r.Header.Get(PolicyOverrideHeaderName)
	if raw == "" {
		return nil
	}

	override, err := parseutil.ParseBool(raw)
	if err != nil {
		return err
	}

	req.PolicyOverride = override
	return nil
}

// requestWrapInfo adds the WrapInfo value to the logical.Request if wrap info exists
func requestWrapInfo(r *http.Request, req *logical.Request) (*logical.Request, error) {
	// First try for the header value
	wrapTTL := r.Header.Get(WrapTTLHeaderName)
	if wrapTTL == "" {
		return req, nil
	}

	// If it has an allowed suffix parse as a duration string
	dur, err := parseutil.ParseDurationSecond(wrapTTL)
	if err != nil {
		return req, err
	}
	if int64(dur) < 0 {
		return req, fmt.Errorf("requested wrap ttl cannot be negative")
	}

	req.WrapInfo = &logical.RequestWrapInfo{
		TTL: dur,
	}

	wrapFormat := r.Header.Get(WrapFormatHeaderName)
	switch wrapFormat {
	case "jwt":
		req.WrapInfo.Format = "jwt"
	}

	return req, nil
}

// parseMFAHeader parses the MFAHeaderName in the request headers and organizes
// them with MFA method name as the index.
func parseMFAHeader(req *logical.Request) error {
	if req == nil {
		return fmt.Errorf("request is nil")
	}

	if req.Headers == nil {
		return nil
	}

	// Reset and initialize the credentials in the request
	req.MFACreds = make(map[string][]string)

	for _, mfaHeaderValue := range req.Headers[canonicalMFAHeaderName] {
		// Skip the header with no value in it
		if mfaHeaderValue == "" {
			continue
		}

		// Handle the case where only method name is mentioned and no value
		// is supplied
		if !strings.Contains(mfaHeaderValue, ":") {
			// Mark the presence of method name, but set an empty set to it
			// indicating that there were no values supplied for the method
			if req.MFACreds[mfaHeaderValue] == nil {
				req.MFACreds[mfaHeaderValue] = []string{}
			}
			continue
		}

		shardSplits := strings.SplitN(mfaHeaderValue, ":", 2)
		if shardSplits[0] == "" {
			return fmt.Errorf("invalid data in header %q; missing method name", MFAHeaderName)
		}

		if shardSplits[1] == "" {
			return fmt.Errorf("invalid data in header %q; missing method value", MFAHeaderName)
		}

		req.MFACreds[shardSplits[0]] = append(req.MFACreds[shardSplits[0]], shardSplits[1])
	}

	return nil
}

func respondError(w http.ResponseWriter, status int, err error) {
	logical.RespondError(w, status, err)
}

func respondErrorCommon(w http.ResponseWriter, req *logical.Request, resp *logical.Response, err error) bool {
	statusCode, newErr := logical.RespondErrorCommon(req, resp, err)
	if newErr == nil && statusCode == 0 {
		return false
	}

	respondError(w, statusCode, newErr)
	return true
}

func respondOk(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if body == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(body)
	}
}
