package http

import (
	"encoding/json"
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
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	// AuthHeaderName is the name of the header containing the token.
	AuthHeaderName = "X-Vault-Token"

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

	// MaxRequestSize is the maximum accepted request size. This is to prevent
	// a denial of service attack where no Content-Length is provided and the server
	// is fed ever more data until it exhausts memory.
	MaxRequestSize = 32 * 1024 * 1024
)

var (
	ReplicationStaleReadTimeout = 2 * time.Second

	// Set to false by stub_asset if the ui build tag isn't enabled
	uiBuiltIn = true
)

// Handler returns an http.Handler for the API. This can be used on
// its own to mount the Vault API within another web server.
func Handler(core *vault.Core) http.Handler {
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
	mux.Handle("/v1/sys/wrapping/lookup", handleRequestForwarding(core, handleLogical(core, false, wrappingVerificationFunc)))
	mux.Handle("/v1/sys/wrapping/rewrap", handleRequestForwarding(core, handleLogical(core, false, wrappingVerificationFunc)))
	mux.Handle("/v1/sys/wrapping/unwrap", handleRequestForwarding(core, handleLogical(core, false, wrappingVerificationFunc)))
	for _, path := range injectDataIntoTopRoutes {
		mux.Handle(path, handleRequestForwarding(core, handleLogical(core, true, nil)))
	}
	mux.Handle("/v1/sys/", handleRequestForwarding(core, handleLogical(core, false, nil)))
	mux.Handle("/v1/", handleRequestForwarding(core, handleLogical(core, false, nil)))
	if core.UIEnabled() == true {
		if uiBuiltIn {
			mux.Handle("/ui/", http.StripPrefix("/ui/", gziphandler.GzipHandler(handleUIHeaders(core, handleUI(http.FileServer(&UIAssetWrapper{FileSystem: assetFS()}))))))
		} else {
			mux.Handle("/ui/", handleUIHeaders(core, handleUIStub()))
		}
		mux.Handle("/", handleRootRedirect())
	}

	// Wrap the handler in another handler to trigger all help paths.
	helpWrappedHandler := wrapHelpHandler(mux, core)
	corsWrappedHandler := wrapCORSHandler(helpWrappedHandler, core)

	// Wrap the help wrapped handler with another layer with a generic
	// handler
	genericWrappedHandler := wrapGenericHandler(corsWrappedHandler)

	// Wrap the handler with PrintablePathCheckHandler to check for non-printable
	// characters in the request path.
	printablePathCheckHandler := cleanhttp.PrintablePathCheckHandler(genericWrappedHandler, nil)

	return printablePathCheckHandler
}

// wrapGenericHandler wraps the handler with an extra layer of handler where
// tasks that should be commonly handled for all the requests and/or responses
// are performed.
func wrapGenericHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Cache-Control header for all the responses returned
		// by Vault
		w.Header().Set("Cache-Control", "no-store")
		h.ServeHTTP(w, r)
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
func wrappingVerificationFunc(core *vault.Core, req *logical.Request) error {
	if req == nil {
		return fmt.Errorf("invalid request")
	}

	valid, err := core.ValidateWrappingToken(req)
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
	<p>Vault UI is not available in this binary. To get Vault UI do one of the following:</p>
	<ul>
	<li><a href="https://www.vaultproject.io/downloads.html">Download an official release</a></li>
	<li>Run <code>make release</code> to create your own release binaries.
	<li>Run <code>make dev-ui</code> to create a development binary with the UI.
	</ul>
	</html>
	`
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(stubHTML))
	})
}

func handleRootRedirect() http.Handler {
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
	limit := http.MaxBytesReader(w, r.Body, MaxRequestSize)
	err := jsonutil.DecodeJSONFromReader(limit, out)
	if err != nil && err != io.EOF {
		return errwrap.Wrapf("failed to parse JSON input: {{err}}", err)
	}
	return err
}

// handleRequestForwarding determines whether to forward a request or not,
// falling back on the older behavior of redirecting the client
func handleRequestForwarding(core *vault.Core, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(vault.IntNoForwardingHeaderName) != "" {
			handler.ServeHTTP(w, r)
			return
		}

		if r.Header.Get(NoRequestForwardingHeaderName) != "" {
			// Forwarding explicitly disabled, fall back to previous behavior
			core.Logger().Debug("handleRequestForwarding: forwarding disabled by client request")
			handler.ServeHTTP(w, r)
			return
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

		// Attempt forwarding the request. If we cannot forward -- perhaps it's
		// been disabled on the active node -- this will return with an
		// ErrCannotForward and we simply fall back
		statusCode, header, retBytes, err := core.ForwardRequest(r)
		if err != nil {
			if err == vault.ErrCannotForward {
				core.Logger().Debug("handleRequestForwarding: cannot forward (possibly disabled on active node), falling back")
			} else {
				core.Logger().Error("handleRequestForwarding: error forwarding request", "error", err)
			}

			// Fall back to redirection
			handler.ServeHTTP(w, r)
			return
		}

		if header != nil {
			for k, v := range header {
				w.Header()[k] = v
			}
		}

		w.WriteHeader(statusCode)
		w.Write(retBytes)
		return
	})
}

// request is a helper to perform a request and properly exit in the
// case of an error.
func request(core *vault.Core, w http.ResponseWriter, rawReq *http.Request, r *logical.Request) (*logical.Response, bool) {
	resp, err := core.HandleRequest(r)
	if errwrap.Contains(err, consts.ErrStandby.Error()) {
		respondStandby(core, w, rawReq.URL)
		return resp, false
	}
	if respondErrorCommon(w, r, resp, err) {
		return resp, false
	}

	return resp, true
}

// respondStandby is used to trigger a redirect in the case that this Vault is currently a hot standby
func respondStandby(core *vault.Core, w http.ResponseWriter, reqURL *url.URL) {
	// Request the leader address
	_, redirectAddr, _, err := core.Leader()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// If there is no leader, generate a 503 error
	if redirectAddr == "" {
		err = fmt.Errorf("no active Vault instance found")
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

// requestAuth adds the token to the logical.Request if it exists.
func requestAuth(core *vault.Core, r *http.Request, req *logical.Request) *logical.Request {
	// Attach the header value if we have it
	if v := r.Header.Get(AuthHeaderName); v != "" {
		req.ClientToken = v

		// Also attach the accessor if we have it. This doesn't fail if it
		// doesn't exist because the request may be to an unauthenticated
		// endpoint/login endpoint where a bad current token doesn't matter, or
		// a token from a Vault version pre-accessors.
		te, err := core.LookupToken(v)
		if err == nil && te != nil {
			req.ClientTokenAccessor = te.Accessor
			req.ClientTokenRemainingUses = te.NumUses
			req.SetTokenEntry(te)
		}
	}

	return req
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

func respondError(w http.ResponseWriter, status int, err error) {
	logical.AdjustErrorStatusCode(&status, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
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

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

var injectDataIntoTopRoutes = []string{
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
