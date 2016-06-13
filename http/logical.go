package http

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

type PrepareRequestFunc func(req *logical.Request) error

func buildLogicalRequest(w http.ResponseWriter, r *http.Request) (*logical.Request, int, error) {
	// Determine the path...
	if !strings.HasPrefix(r.URL.Path, "/v1/") {
		return nil, http.StatusNotFound, nil
	}
	path := r.URL.Path[len("/v1/"):]
	if path == "" {
		return nil, http.StatusNotFound, nil
	}

	// Determine the operation
	var op logical.Operation
	switch r.Method {
	case "DELETE":
		op = logical.DeleteOperation
	case "GET":
		op = logical.ReadOperation
		// Need to call ParseForm to get query params loaded
		queryVals := r.URL.Query()
		listStr := queryVals.Get("list")
		if listStr != "" {
			list, err := strconv.ParseBool(listStr)
			if err != nil {
				return nil, http.StatusBadRequest, nil
			}
			if list {
				op = logical.ListOperation
			}
		}
	case "POST", "PUT":
		op = logical.UpdateOperation
	case "LIST":
		op = logical.ListOperation
	default:
		return nil, http.StatusMethodNotAllowed, nil
	}

	// Parse the request if we can
	var data map[string]interface{}
	if op == logical.UpdateOperation {
		err := parseRequest(r, &data)
		if err == io.EOF {
			data = nil
			err = nil
		}
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	var err error
	req := requestAuth(r, &logical.Request{
		Operation:  op,
		Path:       path,
		Data:       data,
		Connection: getConnection(r),
	})
	req, err = requestWrapTTL(r, req)
	if err != nil {
		return nil, http.StatusBadRequest, errwrap.Wrapf("error parsing X-Vault-Wrap-TTL header: {{err}}", err)
	}

	return req, 0, nil
}

func handleLogical(core *vault.Core, dataOnly bool, prepareRequestCallback PrepareRequestFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, statusCode, err := buildLogicalRequest(w, r)
		if err != nil || statusCode != 0 {
			respondError(w, statusCode, err)
			return
		}

		// Certain endpoints may require changes to the request object. They
		// will have a callback registered to do the needed operations, so
		// invoke it before proceeding.
		if prepareRequestCallback != nil {
			if err := prepareRequestCallback(req); err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return
			}
		}

		// Make the internal request. We attach the connection info
		// as well in case this is an authentication request that requires
		// it. Vault core handles stripping this if we need to.
		resp, ok := request(core, w, r, req)
		if !ok {
			return
		}
		switch {
		case req.Operation == logical.ReadOperation:
			if resp == nil {
				respondError(w, http.StatusNotFound, nil)
				return
			}

		// Basically: if we have empty "keys" or no keys at all, 404. This
		// provides consistency with GET.
		case req.Operation == logical.ListOperation:
			if resp == nil || len(resp.Data) == 0 {
				respondError(w, http.StatusNotFound, nil)
				return
			}
			keysInt, ok := resp.Data["keys"]
			if !ok || keysInt == nil {
				respondError(w, http.StatusNotFound, nil)
				return
			}
			keys, ok := keysInt.([]string)
			if !ok {
				respondError(w, http.StatusInternalServerError, nil)
				return
			}
			if len(keys) == 0 {
				respondError(w, http.StatusNotFound, nil)
				return
			}
		}

		// Build the proper response
		respondLogical(w, r, req.Path, dataOnly, resp)
	})
}

func respondLogical(w http.ResponseWriter, r *http.Request, path string, dataOnly bool, resp *logical.Response) {
	var httpResp interface{}
	if resp != nil {
		if resp.Redirect != "" {
			// If we have a redirect, redirect! We use a 307 code
			// because we don't actually know if its permanent.
			http.Redirect(w, r, resp.Redirect, 307)
			return
		}

		if dataOnly {
			respondOk(w, resp.Data)
			return
		}

		// Check if this is a raw response
		if _, ok := resp.Data[logical.HTTPContentType]; ok {
			respondRaw(w, r, path, resp)
			return
		}

		if resp.WrapInfo != nil && resp.WrapInfo.Token != "" {
			httpResp = logical.HTTPResponse{
				WrapInfo: &logical.HTTPWrapInfo{
					Token:           resp.WrapInfo.Token,
					TTL:             int(resp.WrapInfo.TTL.Seconds()),
					CreationTime:    resp.WrapInfo.CreationTime,
					WrappedAccessor: resp.WrapInfo.WrappedAccessor,
				},
			}
		} else {
			httpResp = logical.SanitizeResponse(resp)
		}
	}

	// Respond
	respondOk(w, httpResp)
	return
}

// respondRaw is used when the response is using HTTPContentType and HTTPRawBody
// to change the default response handling. This is only used for specific things like
// returning the CRL information on the PKI backends.
func respondRaw(w http.ResponseWriter, r *http.Request, path string, resp *logical.Response) {
	// Ensure this is never a secret or auth response
	if resp.Secret != nil || resp.Auth != nil {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	// Get the status code
	statusRaw, ok := resp.Data[logical.HTTPStatusCode]
	if !ok {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}
	status, ok := statusRaw.(int)
	if !ok {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	// Get the header
	contentTypeRaw, ok := resp.Data[logical.HTTPContentType]
	if !ok {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}
	contentType, ok := contentTypeRaw.(string)
	if !ok {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	// Get the body
	bodyRaw, ok := resp.Data[logical.HTTPRawBody]
	if !ok {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}
	body, ok := bodyRaw.([]byte)
	if !ok {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	w.Write(body)
}

// getConnection is used to format the connection information for
// attaching to a logical request
func getConnection(r *http.Request) (connection *logical.Connection) {
	var remoteAddr string

	remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteAddr = ""
	}

	connection = &logical.Connection{
		RemoteAddr: remoteAddr,
		ConnState:  r.TLS,
	}
	return
}
