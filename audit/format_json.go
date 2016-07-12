package audit

import (
	"encoding/json"
	"io"
	"time"

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
		Time:  time.Now().UTC().Format(time.RFC3339),
		Type:  "request",
		Error: errString,

		Auth: JSONAuth{
			DisplayName: auth.DisplayName,
			Policies:    auth.Policies,
			Metadata:    auth.Metadata,
		},

		Request: JSONRequest{
			ClientToken: req.ClientToken,
			Operation:   req.Operation,
			Path:        req.Path,
			Data:        req.Data,
			RemoteAddr:  getRemoteAddr(req),
			WrapTTL:     int(req.WrapTTL / time.Second),
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
			Accessor:    resp.Auth.Accessor,
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

	var respWrapInfo *JSONWrapInfo
	if resp.WrapInfo != nil {
		respWrapInfo = &JSONWrapInfo{
			TTL:             int(resp.WrapInfo.TTL / time.Second),
			Token:           resp.WrapInfo.Token,
			CreationTime:    resp.WrapInfo.CreationTime,
			WrappedAccessor: resp.WrapInfo.WrappedAccessor,
		}
	}

	// Encode!
	enc := json.NewEncoder(w)
	return enc.Encode(&JSONResponseEntry{
		Time:  time.Now().UTC().Format(time.RFC3339),
		Type:  "response",
		Error: errString,

		Auth: JSONAuth{
			DisplayName: auth.DisplayName,
			Policies:    auth.Policies,
			Metadata:    auth.Metadata,
		},

		Request: JSONRequest{
			ClientToken: req.ClientToken,
			Operation:   req.Operation,
			Path:        req.Path,
			Data:        req.Data,
			RemoteAddr:  getRemoteAddr(req),
			WrapTTL:     int(req.WrapTTL / time.Second),
		},

		Response: JSONResponse{
			Auth:     respAuth,
			Secret:   respSecret,
			Data:     resp.Data,
			Redirect: resp.Redirect,
			WrapInfo: respWrapInfo,
		},
	})
}

// JSONRequest is the structure of a request audit log entry in JSON.
type JSONRequestEntry struct {
	Time    string      `json:"time"`
	Type    string      `json:"type"`
	Auth    JSONAuth    `json:"auth"`
	Request JSONRequest `json:"request"`
	Error   string      `json:"error"`
}

// JSONResponseEntry is the structure of a response audit log entry in JSON.
type JSONResponseEntry struct {
	Time     string       `json:"time"`
	Type     string       `json:"type"`
	Error    string       `json:"error"`
	Auth     JSONAuth     `json:"auth"`
	Request  JSONRequest  `json:"request"`
	Response JSONResponse `json:"response"`
}

type JSONRequest struct {
	Operation   logical.Operation      `json:"operation"`
	ClientToken string                 `json:"client_token"`
	Path        string                 `json:"path"`
	Data        map[string]interface{} `json:"data"`
	RemoteAddr  string                 `json:"remote_address"`
	WrapTTL     int                    `json:"wrap_ttl"`
}

type JSONResponse struct {
	Auth     *JSONAuth              `json:"auth,omitempty"`
	Secret   *JSONSecret            `json:"secret,emitempty"`
	Data     map[string]interface{} `json:"data"`
	Redirect string                 `json:"redirect"`
	WrapInfo *JSONWrapInfo          `json:"wrap_info,omitempty"`
}

type JSONAuth struct {
	ClientToken string            `json:"client_token,omitempty"`
	Accessor    string            `json:"accessor,omitempty"`
	DisplayName string            `json:"display_name"`
	Policies    []string          `json:"policies"`
	Metadata    map[string]string `json:"metadata"`
}

type JSONSecret struct {
	LeaseID string `json:"lease_id"`
}

type JSONWrapInfo struct {
	TTL             int       `json:"ttl"`
	Token           string    `json:"token"`
	CreationTime    time.Time `json:"creation_time"`
	WrappedAccessor string    `json:"wrapped_accessor,omitempty"`
}

// getRemoteAddr safely gets the remote address avoiding a nil pointer
func getRemoteAddr(req *logical.Request) string {
	if req != nil && req.Connection != nil {
		return req.Connection.RemoteAddr
	}
	return ""
}
