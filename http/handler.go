package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/duration"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/requestutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	// AuthHeaderName is the name of the header containing the token.
	AuthHeaderName = "X-Vault-Token"

	// WrapHeaderName is the name of the header containing a directive to wrap the
	// response.
	WrapTTLHeaderName = "X-Vault-Wrap-TTL"

	// NoRequestForwardingHeaderName is the name of the header telling Vault
	// not to use request forwarding
	NoRequestForwardingHeaderName = "X-Vault-No-Request-Forwarding"

	// Internal so as not to log a trace message
	intNoForwardingHeaderName = "X-Vault-Internal-No-Request-Forwarding"
)

// Handler returns an http.Handler for the API. This can be used on
// its own to mount the Vault API within another web server.
func Handler(core *vault.Core) http.Handler {
	// Create the muxer to handle the actual endpoints
	mux := http.NewServeMux()
	mux.Handle("/v1/sys/init", handleSysInit(core))
	mux.Handle("/v1/sys/seal-status", handleSysSealStatus(core))
	mux.Handle("/v1/sys/seal", handleSysSeal(core))
	mux.Handle("/v1/sys/step-down", handleSysStepDown(core))
	mux.Handle("/v1/sys/unseal", handleSysUnseal(core))
	mux.Handle("/v1/sys/renew", handleRequestForwarding(core, handleLogical(core, false, nil)))
	mux.Handle("/v1/sys/renew/", handleRequestForwarding(core, handleLogical(core, false, nil)))
	mux.Handle("/v1/sys/leader", handleSysLeader(core))
	mux.Handle("/v1/sys/health", handleSysHealth(core))
	mux.Handle("/v1/sys/generate-root/attempt", handleRequestForwarding(core, handleSysGenerateRootAttempt(core)))
	mux.Handle("/v1/sys/generate-root/update", handleRequestForwarding(core, handleSysGenerateRootUpdate(core)))
	mux.Handle("/v1/sys/rekey/init", handleRequestForwarding(core, handleSysRekeyInit(core, false)))
	mux.Handle("/v1/sys/rekey/update", handleRequestForwarding(core, handleSysRekeyUpdate(core, false)))
	mux.Handle("/v1/sys/rekey-recovery-key/init", handleRequestForwarding(core, handleSysRekeyInit(core, true)))
	mux.Handle("/v1/sys/rekey-recovery-key/update", handleRequestForwarding(core, handleSysRekeyUpdate(core, true)))
	mux.Handle("/v1/sys/capabilities-self", handleRequestForwarding(core, handleLogical(core, true, sysCapabilitiesSelfCallback)))
	mux.Handle("/v1/sys/", handleRequestForwarding(core, handleLogical(core, true, nil)))
	mux.Handle("/v1/", handleRequestForwarding(core, handleLogical(core, false, nil)))

	// Wrap the handler in another handler to trigger all help paths.
	handler := handleHelpHandler(mux, core)

	return handler
}

// ClientToken is required in the handler of sys/capabilities-self endpoint in
// system backend. But the ClientToken gets obfuscated before the request gets
// forwarded to any logical backend. So, setting the ClientToken in the data
// field for this request.
func sysCapabilitiesSelfCallback(req *logical.Request) error {
	if req == nil || req.Data == nil {
		return fmt.Errorf("invalid request")
	}
	req.Data["token"] = req.ClientToken
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

func parseRequest(r *http.Request, out interface{}) error {
	err := jsonutil.DecodeJSONFromReader(r.Body, out)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Failed to parse JSON input: %s", err)
	}
	return err
}

// handleRequestForwarding determines whether to forward a request or not,
// falling back on the older behavior of redirecting the client
func handleRequestForwarding(core *vault.Core, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(intNoForwardingHeaderName) != "" {
			handler.ServeHTTP(w, r)
			return
		}

		if r.Header.Get(NoRequestForwardingHeaderName) != "" {
			// Forwarding explicitly disabled, fall back to previous behavior
			core.Logger().Printf("[TRACE] http/handleRequestForwarding: forwarding disabled by client request")
			handler.ServeHTTP(w, r)
			return
		}

		// Note: in an HA setup, this call will also ensure that connections to
		// the leader are set up, as that happens once the advertised cluster
		// values are read during this function
		isLeader, leaderAddr, err := core.Leader()
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
			respondError(w, http.StatusInternalServerError, fmt.Errorf("node not active but active node not found"))
			return
		}

		core.Logger().Printf("[TRACE] http/handleRequestForwarding: forwarding request")

		// Attempt forwarding the request. If we cannot forward -- perhaps it's
		// been disabled on the active node -- this will return with an
		// ErrCannotForward and we simply fall back
		resp, err := core.ForwardRequest(r)
		if err != nil {
			if err == vault.ErrCannotForward {
				// This will use redirects as normal
				core.Logger().Printf("[TRACE] http/handleRequestForwarding: cannot forward, falling back")
				handler.ServeHTTP(w, r)
				return
			}
			core.Logger().Printf("[ERR] http/handleRequestForwarding: error forwarding request: %v", err)
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()

		// Read the body into a buffer so we can write it back out to the
		// original requestor
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			core.Logger().Printf("[ERR] http/handleRequestForwarding: error reading response body: %v", err)
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(buf.Bytes())
		return
	})
}

// request is a helper to perform a request and properly exit in the
// case of an error.
func request(core *vault.Core, w http.ResponseWriter, rawReq *http.Request, r *logical.Request) (*logical.Response, bool) {
	resp, err := core.HandleRequest(r)
	if errwrap.Contains(err, vault.ErrStandby.Error()) {
		respondStandby(core, w, rawReq.URL)
		return resp, false
	}
	if respondErrorCommon(w, resp, err) {
		return resp, false
	}

	return resp, true
}

// respondStandby is used to trigger a redirect in the case that this Vault is currently a hot standby
func respondStandby(core *vault.Core, w http.ResponseWriter, reqURL *url.URL) {
	// Request the leader address
	_, advertise, err := core.Leader()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// If there is no leader, generate a 503 error
	if advertise == "" {
		err = fmt.Errorf("no active Vault instance found")
		respondError(w, http.StatusServiceUnavailable, err)
		return
	}

	// Parse the advertise location
	advertiseURL, err := url.Parse(advertise)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Generate a redirect URL
	redirectURL := url.URL{
		Scheme:   advertiseURL.Scheme,
		Host:     advertiseURL.Host,
		Path:     reqURL.Path,
		RawQuery: reqURL.RawQuery,
	}

	// Ensure there is a scheme, default to https
	if redirectURL.Scheme == "" {
		redirectURL.Scheme = "https"
	}

	// If we have an address, redirect! We use a 307 code
	// because we don't actually know if its permanent and
	// the request method should be preserved.
	w.Header().Set("Location", redirectURL.String())
	w.WriteHeader(307)
}

// requestAuth adds the token to the logical.Request if it exists.
func requestAuth(r *http.Request, req *logical.Request) *logical.Request {
	// Attach the header value if we have it
	if v := r.Header.Get(AuthHeaderName); v != "" {
		req.ClientToken = v
	}

	return req
}

// requestWrapTTL adds the WrapTTL value to the logical.Request if it
// exists.
func requestWrapTTL(r *http.Request, req *logical.Request) (*logical.Request, error) {
	// First try for the header value
	wrapTTL := r.Header.Get(WrapTTLHeaderName)
	if wrapTTL == "" {
		return req, nil
	}

	// If it has an allowed suffix parse as a duration string
	dur, err := duration.ParseDurationSecond(wrapTTL)
	if err != nil {
		return req, err
	}
	if int64(dur) < 0 {
		return req, fmt.Errorf("requested wrap ttl cannot be negative")
	}
	req.WrapTTL = dur

	return req, nil
}

func respondError(w http.ResponseWriter, status int, err error) {
	// Adjust status code when sealed
	if errwrap.Contains(err, vault.ErrSealed.Error()) {
		status = http.StatusServiceUnavailable
	}

	// Allow HTTPCoded error passthrough to specify a code
	if t, ok := err.(logical.HTTPCodedError); ok {
		status = t.Code()
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func respondErrorCommon(w http.ResponseWriter, resp *logical.Response, err error) bool {
	// If there are no errors return
	if err == nil && (resp == nil || !resp.IsError()) {
		return false
	}

	// Start out with internal server error since in most of these cases there
	// won't be a response so this won't be overridden
	statusCode := http.StatusInternalServerError
	// If we actually have a response, start out with bad request
	if resp != nil {
		statusCode = http.StatusBadRequest
	}

	// Now, check the error itself; if it has a specific logical error, set the
	// appropriate code
	if err != nil {
		switch {
		case errwrap.ContainsType(err, new(vault.StatusBadRequest)):
			statusCode = http.StatusBadRequest
		case errwrap.Contains(err, logical.ErrPermissionDenied.Error()):
			statusCode = http.StatusForbidden
		case errwrap.Contains(err, logical.ErrUnsupportedOperation.Error()):
			statusCode = http.StatusMethodNotAllowed
		case errwrap.Contains(err, logical.ErrUnsupportedPath.Error()):
			statusCode = http.StatusNotFound
		case errwrap.Contains(err, logical.ErrInvalidRequest.Error()):
			statusCode = http.StatusBadRequest
		}
	}

	if resp != nil && resp.IsError() {
		err = fmt.Errorf("%s", resp.Data["error"].(string))
	}

	respondError(w, statusCode, err)
	return true
}

func respondOk(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")

	if body == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(body)
	}
}

// WrapListenersForClustering takes in Vault's listeners and original HTTP
// handler, creates a new handler that handles forwarded requests, and returns
// the cluster setup function that creates the new listners and assigns to the
// new handler
func WrapListenersForClustering(lns []net.Listener, handler http.Handler, logger *log.Logger) func() ([]net.Listener, http.Handler, error) {
	// This mux handles cluster functions (right now, only forwarded requests)
	mux := http.NewServeMux()
	mux.HandleFunc("/cluster/forwarded-request", func(w http.ResponseWriter, req *http.Request) {
		freq, err := requestutil.ParseForwardedRequest(req)
		if err != nil {
			logger.Printf("[ERR] http/ForwardedRequestHandler: error parsing forwarded request: %v", err)
			respondError(w, http.StatusInternalServerError, fmt.Errorf("unable to parse forwarded request"))
			return
		}

		// To avoid the risk of a forward loop in some pathological condition,
		// set the no-forward header
		freq.Header.Set(intNoForwardingHeaderName, "true")
		handler.ServeHTTP(w, freq)
	})

	return func() ([]net.Listener, http.Handler, error) {
		ret := make([]net.Listener, 0, len(lns))
		// Loop over the existing listeners and start listeners on appropriate ports
		for _, ln := range lns {
			tcpAddr, ok := ln.Addr().(*net.TCPAddr)
			if !ok {
				logger.Printf("[TRACE] http/WrapClusterListener: %s not a candidate for cluster request handling", ln.Addr().String())
				continue
			}
			logger.Printf("[TRACE] http/WrapClusterListener: %s is a candidate for cluster request handling at addr %s and port %d", tcpAddr.String(), tcpAddr.IP.String(), tcpAddr.Port+1)

			ipStr := tcpAddr.IP.String()
			if len(tcpAddr.IP) == net.IPv6len {
				ipStr = fmt.Sprintf("[%s]", ipStr)
			}
			ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ipStr, tcpAddr.Port+1))
			if err != nil {
				return nil, nil, err
			}
			ret = append(ret, ln)
		}

		return ret, mux, nil
	}
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}
