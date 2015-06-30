package audit

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/vault/logical"
)

// FormatJSON is a Formatter implementation that structures data into
// a JSON format.
type FormatJSON struct{}

func (f *FormatJSON) FormatRequest(
	w io.Writer,
	auth *logical.Auth,
	req *logical.Request,
	err error) error {

	// If auth is nil, make an empty one
	if auth == nil {
		auth = new(logical.Auth)
	}
	var errString string
	if err != nil {
		errString = err.Error()
	}

	// Encode!
	enc := json.NewEncoder(w)
	return enc.Encode(&JSONRequestEntry{
		Type:  "request",
		Error: errString,

		Auth: JSONAuth{
			DisplayName: auth.DisplayName,
			Policies:    auth.Policies,
			Metadata:    auth.Metadata,
		},

		Request: JSONRequest{
			Operation:  req.Operation,
			Path:       req.Path,
			Data:       req.Data,
			RemoteAddr: getRemoteAddr(req),
		},
	})
}

func (f *FormatJSON) FormatResponse(
	w io.Writer,
	auth *logical.Auth,
	req *logical.Request,
	resp *logical.Response,
	err error) error {
	// If things are nil, make empty to avoid panics
	if auth == nil {
		auth = new(logical.Auth)
	}
	if resp == nil {
		resp = new(logical.Response)
	}
	var errString string
	if err != nil {
		errString = err.Error()
	}

	var respAuth *JSONAuth
	if resp.Auth != nil {
		respAuth = &JSONAuth{
			ClientToken: resp.Auth.ClientToken,
			DisplayName: resp.Auth.DisplayName,
			Policies:    resp.Auth.Policies,
			Metadata:    resp.Auth.Metadata,
		}
	}

	var respSecret *JSONSecret
	if resp.Secret != nil {
		respSecret = &JSONSecret{
			LeaseID: resp.Secret.LeaseID,
		}
	}

	// Encode!
	enc := json.NewEncoder(w)
	return enc.Encode(&JSONResponseEntry{
		Type:  "response",
		Error: errString,

		Auth: JSONAuth{
			Policies: auth.Policies,
			Metadata: auth.Metadata,
		},

		Request: JSONRequest{
			Operation:  req.Operation,
			Path:       req.Path,
			Data:       req.Data,
			RemoteAddr: getRemoteAddr(req),
		},

		Response: JSONResponse{
			Auth:     respAuth,
			Secret:   respSecret,
			Data:     resp.Data,
			Redirect: resp.Redirect,
		},
	})
}

// JSONRequest is the structure of a request audit log entry in JSON.
type JSONRequestEntry struct {
	Type    string      `json:"type"`
	Auth    JSONAuth    `json:"auth"`
	Request JSONRequest `json:"request"`
	Error   string      `json:"error"`
}

// JSONResponseEntry is the structure of a response audit log entry in JSON.
type JSONResponseEntry struct {
	Type     string       `json:"type"`
	Error    string       `json:"error"`
	Auth     JSONAuth     `json:"auth"`
	Request  JSONRequest  `json:"request"`
	Response JSONResponse `json:"response"`
}

type JSONRequest struct {
	Operation  logical.Operation      `json:"operation"`
	Path       string                 `json:"path"`
	Data       map[string]interface{} `json:"data"`
	RemoteAddr string                 `json:"remote_address"`
}

type JSONResponse struct {
	Auth     *JSONAuth              `json:"auth,omitempty"`
	Secret   *JSONSecret            `json:"secret,emitempty"`
	Data     map[string]interface{} `json:"data"`
	Redirect string                 `json:"redirect"`
}

type JSONAuth struct {
	ClientToken string            `json:"client_token,omitempty"`
	DisplayName string            `json:"display_name"`
	Policies    []string          `json:"policies"`
	Metadata    map[string]string `json:"metadata"`
}

type JSONSecret struct {
	LeaseID string `json:"lease_id"`
}

// getRemoteAddr safely gets the remote address avoiding a nil pointer
func getRemoteAddr(req *logical.Request) string {
	if req != nil && req.Connection != nil {
		return req.Connection.RemoteAddr
	}
	return ""
}
