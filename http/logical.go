package http

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

type PrepareRequestFunc func(req *logical.Request) error

func handleLogical(core *vault.Core, dataOnly bool, prepareRequestCallback PrepareRequestFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determine the path...
		if !strings.HasPrefix(r.URL.Path, "/v1/") {
			respondError(w, http.StatusNotFound, nil)
			return
		}
		path := r.URL.Path[len("/v1/"):]
		if path == "" {
			respondError(w, http.StatusNotFound, nil)
			return
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
					respondError(w, http.StatusBadRequest, nil)
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
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
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
				respondError(w, http.StatusBadRequest, err)
				return
			}
		}

		req := requestAuth(r, &logical.Request{
			Operation:  op,
			Path:       path,
			Data:       data,
			Connection: getConnection(r),
		})

		// Certain endpoints may require changes to the request object.
		// They will have a callback registered to do the needful.
		// Invoking it before proceeding.
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
		if (op == logical.ReadOperation || op == logical.ListOperation) && resp == nil {
			respondError(w, http.StatusNotFound, nil)
			return
		}

		// Build the proper response
		respondLogical(w, r, path, dataOnly, resp)
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

		logicalResp := &LogicalResponse{
			Data:     resp.Data,
			Warnings: resp.Warnings(),
		}
		if resp.Secret != nil {
			logicalResp.LeaseID = resp.Secret.LeaseID
			logicalResp.Renewable = resp.Secret.Renewable
			logicalResp.LeaseDuration = int(resp.Secret.TTL.Seconds())
		}

		// If we have authentication information, then
		// set up the result structure.
		if resp.Auth != nil {
			logicalResp.Auth = &Auth{
				ClientToken:   resp.Auth.ClientToken,
				Accessor:      resp.Auth.Accessor,
				Policies:      resp.Auth.Policies,
				Metadata:      resp.Auth.Metadata,
				LeaseDuration: int(resp.Auth.TTL.Seconds()),
				Renewable:     resp.Auth.Renewable,
			}
		}

		httpResp = logicalResp
	}

	// Respond
	respondOk(w, httpResp)
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

type LogicalResponse struct {
	LeaseID       string                 `json:"lease_id"`
	Renewable     bool                   `json:"renewable"`
	LeaseDuration int                    `json:"lease_duration"`
	Data          map[string]interface{} `json:"data"`
	Warnings      []string               `json:"warnings"`
	Auth          *Auth                  `json:"auth"`
}

type Auth struct {
	ClientToken   string            `json:"client_token"`
	Accessor      string            `json:"accessor"`
	Policies      []string          `json:"policies"`
	Metadata      map[string]string `json:"metadata"`
	LeaseDuration int               `json:"lease_duration"`
	Renewable     bool              `json:"renewable"`
}
