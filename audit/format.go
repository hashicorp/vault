package audit

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	squarejwt "gopkg.in/square/go-jose.v2/jwt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
)

type AuditFormatWriter interface {
	WriteRequest(io.Writer, *AuditRequestEntry) error
	WriteResponse(io.Writer, *AuditResponseEntry) error
	Salt(context.Context) (*salt.Salt, error)
}

// AuditFormatter implements the Formatter interface, and allows the underlying
// marshaller to be swapped out
type AuditFormatter struct {
	AuditFormatWriter
}

var _ Formatter = (*AuditFormatter)(nil)

func (f *AuditFormatter) FormatRequest(ctx context.Context, w io.Writer, config FormatterConfig, in *logical.LogInput) error {
	if in == nil || in.Request == nil {
		return fmt.Errorf("request to request-audit a nil request")
	}

	if w == nil {
		return fmt.Errorf("writer for audit request is nil")
	}

	if f.AuditFormatWriter == nil {
		return fmt.Errorf("no format writer specified")
	}

	salt, err := f.Salt(ctx)
	if err != nil {
		return errwrap.Wrapf("error fetching salt: {{err}}", err)
	}

	// Set these to the input values at first
	auth := in.Auth
	req := in.Request
	var connState *tls.ConnectionState

	if in.Request.Connection != nil && in.Request.Connection.ConnState != nil {
		connState = in.Request.Connection.ConnState
	}

	if !config.Raw {
		// Before we copy the structure we must nil out some data
		// otherwise we will cause reflection to panic and die
		if connState != nil {
			in.Request.Connection.ConnState = nil
			defer func() {
				in.Request.Connection.ConnState = connState
			}()
		}

		// Copy the auth structure
		if in.Auth != nil {
			cp, err := copystructure.Copy(in.Auth)
			if err != nil {
				return err
			}
			auth = cp.(*logical.Auth)
		}

		cp, err := copystructure.Copy(in.Request)
		if err != nil {
			return err
		}
		req = cp.(*logical.Request)
		for k, v := range req.Data {
			if o, ok := v.(logical.OptMarshaler); ok {
				marshaled, err := o.MarshalJSONWithOptions(&logical.MarshalOptions{
					ValueHasher: salt.GetIdentifiedHMAC,
				})
				if err != nil {
					return err
				}
				req.Data[k] = json.RawMessage(marshaled)
			}
		}

		// Hash any sensitive information
		if auth != nil {
			// Cache and restore accessor in the auth
			var authAccessor string
			if !config.HMACAccessor && auth.Accessor != "" {
				authAccessor = auth.Accessor
			}
			if err := Hash(salt, auth, nil); err != nil {
				return err
			}
			if authAccessor != "" {
				auth.Accessor = authAccessor
			}
		}

		// Cache and restore accessor in the request
		var clientTokenAccessor string
		if !config.HMACAccessor && req != nil && req.ClientTokenAccessor != "" {
			clientTokenAccessor = req.ClientTokenAccessor
		}
		if err := Hash(salt, req, in.NonHMACReqDataKeys); err != nil {
			return err
		}
		if clientTokenAccessor != "" {
			req.ClientTokenAccessor = clientTokenAccessor
		}
	}

	// If auth is nil, make an empty one
	if auth == nil {
		auth = new(logical.Auth)
	}
	var errString string
	if in.OuterErr != nil {
		errString = in.OuterErr.Error()
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	reqType := in.Type
	if reqType == "" {
		reqType = "request"
	}
	reqEntry := &AuditRequestEntry{
		Type:  reqType,
		Error: errString,

		Auth: &AuditAuth{
			ClientToken:               auth.ClientToken,
			Accessor:                  auth.Accessor,
			DisplayName:               auth.DisplayName,
			Policies:                  auth.Policies,
			TokenPolicies:             auth.TokenPolicies,
			IdentityPolicies:          auth.IdentityPolicies,
			ExternalNamespacePolicies: auth.ExternalNamespacePolicies,
			NoDefaultPolicy:           auth.NoDefaultPolicy,
			Metadata:                  auth.Metadata,
			EntityID:                  auth.EntityID,
			RemainingUses:             req.ClientTokenRemainingUses,
			TokenType:                 auth.TokenType.String(),
		},

		Request: &AuditRequest{
			ID:                  req.ID,
			ClientToken:         req.ClientToken,
			ClientTokenAccessor: req.ClientTokenAccessor,
			Operation:           req.Operation,
			Namespace: &AuditNamespace{
				ID:   ns.ID,
				Path: ns.Path,
			},
			Path:                          req.Path,
			Data:                          req.Data,
			PolicyOverride:                req.PolicyOverride,
			RemoteAddr:                    getRemoteAddr(req),
			ReplicationCluster:            req.ReplicationCluster,
			Headers:                       req.Headers,
			ClientCertificateSerialNumber: getClientCertificateSerialNumber(connState),
		},
	}

	if req.WrapInfo != nil {
		reqEntry.Request.WrapTTL = int(req.WrapInfo.TTL / time.Second)
	}

	if !config.OmitTime {
		reqEntry.Time = time.Now().UTC().Format(time.RFC3339Nano)
	}

	return f.AuditFormatWriter.WriteRequest(w, reqEntry)
}

func (f *AuditFormatter) FormatResponse(ctx context.Context, w io.Writer, config FormatterConfig, in *logical.LogInput) error {
	if in == nil || in.Request == nil {
		return fmt.Errorf("request to response-audit a nil request")
	}

	if w == nil {
		return fmt.Errorf("writer for audit request is nil")
	}

	if f.AuditFormatWriter == nil {
		return fmt.Errorf("no format writer specified")
	}

	salt, err := f.Salt(ctx)
	if err != nil {
		return errwrap.Wrapf("error fetching salt: {{err}}", err)
	}

	// Set these to the input values at first
	auth := in.Auth
	req := in.Request
	resp := in.Response
	var connState *tls.ConnectionState

	if in.Request.Connection != nil && in.Request.Connection.ConnState != nil {
		connState = in.Request.Connection.ConnState
	}

	if !config.Raw {
		// Before we copy the structure we must nil out some data
		// otherwise we will cause reflection to panic and die
		if connState != nil {
			in.Request.Connection.ConnState = nil
			defer func() {
				in.Request.Connection.ConnState = connState
			}()
		}

		// Copy the auth structure
		if in.Auth != nil {
			cp, err := copystructure.Copy(in.Auth)
			if err != nil {
				return err
			}
			auth = cp.(*logical.Auth)
		}

		cp, err := copystructure.Copy(in.Request)
		if err != nil {
			return err
		}
		req = cp.(*logical.Request)
		for k, v := range req.Data {
			if o, ok := v.(logical.OptMarshaler); ok {
				marshaled, err := o.MarshalJSONWithOptions(&logical.MarshalOptions{
					ValueHasher: salt.GetIdentifiedHMAC,
				})
				if err != nil {
					return err
				}
				req.Data[k] = json.RawMessage(marshaled)
			}
		}

		if in.Response != nil {
			cp, err := copystructure.Copy(in.Response)
			if err != nil {
				return err
			}
			resp = cp.(*logical.Response)
			for k, v := range resp.Data {
				if o, ok := v.(logical.OptMarshaler); ok {
					marshaled, err := o.MarshalJSONWithOptions(&logical.MarshalOptions{
						ValueHasher: salt.GetIdentifiedHMAC,
					})
					if err != nil {
						return err
					}
					resp.Data[k] = json.RawMessage(marshaled)
				}
			}
		}

		// Hash any sensitive information

		// Cache and restore accessor in the auth
		if auth != nil {
			var accessor string
			if !config.HMACAccessor && auth.Accessor != "" {
				accessor = auth.Accessor
			}
			if err := Hash(salt, auth, nil); err != nil {
				return err
			}
			if accessor != "" {
				auth.Accessor = accessor
			}
		}

		// Cache and restore accessor in the request
		var clientTokenAccessor string
		if !config.HMACAccessor && req != nil && req.ClientTokenAccessor != "" {
			clientTokenAccessor = req.ClientTokenAccessor
		}
		if err := Hash(salt, req, in.NonHMACReqDataKeys); err != nil {
			return err
		}
		if clientTokenAccessor != "" {
			req.ClientTokenAccessor = clientTokenAccessor
		}

		// Cache and restore accessor in the response
		if resp != nil {
			var accessor, wrappedAccessor, wrappingAccessor string
			if !config.HMACAccessor && resp != nil && resp.Auth != nil && resp.Auth.Accessor != "" {
				accessor = resp.Auth.Accessor
			}
			if !config.HMACAccessor && resp != nil && resp.WrapInfo != nil && resp.WrapInfo.WrappedAccessor != "" {
				wrappedAccessor = resp.WrapInfo.WrappedAccessor
				wrappingAccessor = resp.WrapInfo.Accessor
			}
			if err := Hash(salt, resp, in.NonHMACRespDataKeys); err != nil {
				return err
			}
			if accessor != "" {
				resp.Auth.Accessor = accessor
			}
			if wrappedAccessor != "" {
				resp.WrapInfo.WrappedAccessor = wrappedAccessor
			}
			if wrappingAccessor != "" {
				resp.WrapInfo.Accessor = wrappingAccessor
			}
		}
	}

	// If things are nil, make empty to avoid panics
	if auth == nil {
		auth = new(logical.Auth)
	}
	if resp == nil {
		resp = new(logical.Response)
	}
	var errString string
	if in.OuterErr != nil {
		errString = in.OuterErr.Error()
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	var respAuth *AuditAuth
	if resp.Auth != nil {
		respAuth = &AuditAuth{
			ClientToken:               resp.Auth.ClientToken,
			Accessor:                  resp.Auth.Accessor,
			DisplayName:               resp.Auth.DisplayName,
			Policies:                  resp.Auth.Policies,
			TokenPolicies:             resp.Auth.TokenPolicies,
			IdentityPolicies:          resp.Auth.IdentityPolicies,
			ExternalNamespacePolicies: resp.Auth.ExternalNamespacePolicies,
			NoDefaultPolicy:           resp.Auth.NoDefaultPolicy,
			Metadata:                  resp.Auth.Metadata,
			NumUses:                   resp.Auth.NumUses,
			EntityID:                  resp.Auth.EntityID,
			TokenType:                 resp.Auth.TokenType.String(),
		}
	}

	var respSecret *AuditSecret
	if resp.Secret != nil {
		respSecret = &AuditSecret{
			LeaseID: resp.Secret.LeaseID,
		}
	}

	var respWrapInfo *AuditResponseWrapInfo
	if resp.WrapInfo != nil {
		token := resp.WrapInfo.Token
		if jwtToken := parseVaultTokenFromJWT(token); jwtToken != nil {
			token = *jwtToken
		}
		respWrapInfo = &AuditResponseWrapInfo{
			TTL:             int(resp.WrapInfo.TTL / time.Second),
			Token:           token,
			Accessor:        resp.WrapInfo.Accessor,
			CreationTime:    resp.WrapInfo.CreationTime.UTC().Format(time.RFC3339Nano),
			CreationPath:    resp.WrapInfo.CreationPath,
			WrappedAccessor: resp.WrapInfo.WrappedAccessor,
		}
	}

	respType := in.Type
	if respType == "" {
		respType = "response"
	}
	respEntry := &AuditResponseEntry{
		Type:  respType,
		Error: errString,
		Auth: &AuditAuth{
			ClientToken:               auth.ClientToken,
			Accessor:                  auth.Accessor,
			DisplayName:               auth.DisplayName,
			Policies:                  auth.Policies,
			TokenPolicies:             auth.TokenPolicies,
			IdentityPolicies:          auth.IdentityPolicies,
			ExternalNamespacePolicies: auth.ExternalNamespacePolicies,
			NoDefaultPolicy:           auth.NoDefaultPolicy,
			Metadata:                  auth.Metadata,
			RemainingUses:             req.ClientTokenRemainingUses,
			EntityID:                  auth.EntityID,
			TokenType:                 auth.TokenType.String(),
		},

		Request: &AuditRequest{
			ID:                  req.ID,
			ClientToken:         req.ClientToken,
			ClientTokenAccessor: req.ClientTokenAccessor,
			Operation:           req.Operation,
			Namespace: &AuditNamespace{
				ID:   ns.ID,
				Path: ns.Path,
			},
			Path:                          req.Path,
			Data:                          req.Data,
			PolicyOverride:                req.PolicyOverride,
			RemoteAddr:                    getRemoteAddr(req),
			ClientCertificateSerialNumber: getClientCertificateSerialNumber(connState),
			ReplicationCluster:            req.ReplicationCluster,
			Headers:                       req.Headers,
		},

		Response: &AuditResponse{
			Auth:     respAuth,
			Secret:   respSecret,
			Data:     resp.Data,
			Warnings: resp.Warnings,
			Redirect: resp.Redirect,
			WrapInfo: respWrapInfo,
			Headers:  resp.Headers,
		},
	}

	if req.WrapInfo != nil {
		respEntry.Request.WrapTTL = int(req.WrapInfo.TTL / time.Second)
	}

	if !config.OmitTime {
		respEntry.Time = time.Now().UTC().Format(time.RFC3339Nano)
	}

	return f.AuditFormatWriter.WriteResponse(w, respEntry)
}

// AuditRequestEntry is the structure of a request audit log entry in Audit.
type AuditRequestEntry struct {
	Time    string        `json:"time,omitempty"`
	Type    string        `json:"type,omitempty"`
	Auth    *AuditAuth    `json:"auth,omitempty"`
	Request *AuditRequest `json:"request,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// AuditResponseEntry is the structure of a response audit log entry in Audit.
type AuditResponseEntry struct {
	Time     string         `json:"time,omitempty"`
	Type     string         `json:"type,omitempty"`
	Auth     *AuditAuth     `json:"auth,omitempty"`
	Request  *AuditRequest  `json:"request,omitempty"`
	Response *AuditResponse `json:"response,omitempty"`
	Error    string         `json:"error,omitempty"`
}

type AuditRequest struct {
	ID                            string                 `json:"id,omitempty"`
	ReplicationCluster            string                 `json:"replication_cluster,omitempty"`
	Operation                     logical.Operation      `json:"operation,omitempty"`
	ClientToken                   string                 `json:"client_token,omitempty"`
	ClientTokenAccessor           string                 `json:"client_token_accessor,omitempty"`
	Namespace                     *AuditNamespace        `json:"namespace,omitempty"`
	Path                          string                 `json:"path,omitempty"`
	Data                          map[string]interface{} `json:"data,omitempty"`
	PolicyOverride                bool                   `json:"policy_override,omitempty"`
	RemoteAddr                    string                 `json:"remote_address,omitempty"`
	WrapTTL                       int                    `json:"wrap_ttl,omitempty"`
	Headers                       map[string][]string    `json:"headers,omitempty"`
	ClientCertificateSerialNumber string                 `json:"client_certificate_serial_number,omitempty"`
}

type AuditResponse struct {
	Auth     *AuditAuth             `json:"auth,omitempty"`
	Secret   *AuditSecret           `json:"secret,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Warnings []string               `json:"warnings,omitempty"`
	Redirect string                 `json:"redirect,omitempty"`
	WrapInfo *AuditResponseWrapInfo `json:"wrap_info,omitempty"`
	Headers  map[string][]string    `json:"headers,omitempty"`
}

type AuditAuth struct {
	ClientToken               string              `json:"client_token,omitempty"`
	Accessor                  string              `json:"accessor,omitempty"`
	DisplayName               string              `json:"display_name,omitempty"`
	Policies                  []string            `json:"policies,omitempty"`
	TokenPolicies             []string            `json:"token_policies,omitempty"`
	IdentityPolicies          []string            `json:"identity_policies,omitempty"`
	ExternalNamespacePolicies map[string][]string `json:"external_namespace_policies,omitempty"`
	NoDefaultPolicy           bool                `json:"no_default_policy,omitempty"`
	Metadata                  map[string]string   `json:"metadata,omitempty"`
	NumUses                   int                 `json:"num_uses,omitempty"`
	RemainingUses             int                 `json:"remaining_uses,omitempty"`
	EntityID                  string              `json:"entity_id,omitempty"`
	TokenType                 string              `json:"token_type,omitempty"`
}

type AuditSecret struct {
	LeaseID string `json:"lease_id,omitempty"`
}

type AuditResponseWrapInfo struct {
	TTL             int    `json:"ttl,omitempty"`
	Token           string `json:"token,omitempty"`
	Accessor        string `json:"accessor,omitempty"`
	CreationTime    string `json:"creation_time,omitempty"`
	CreationPath    string `json:"creation_path,omitempty"`
	WrappedAccessor string `json:"wrapped_accessor,omitempty"`
}

type AuditNamespace struct {
	ID   string `json:"id,omitempty"`
	Path string `json:"path,omitempty"`
}

// getRemoteAddr safely gets the remote address avoiding a nil pointer
func getRemoteAddr(req *logical.Request) string {
	if req != nil && req.Connection != nil {
		return req.Connection.RemoteAddr
	}
	return ""
}

func getClientCertificateSerialNumber(connState *tls.ConnectionState) string {
	if connState == nil || len(connState.VerifiedChains) == 0 || len(connState.VerifiedChains[0]) == 0 {
		return ""
	}

	return connState.VerifiedChains[0][0].SerialNumber.String()
}

// parseVaultTokenFromJWT returns a string iff the token was a JWT and we could
// extract the original token ID from inside
func parseVaultTokenFromJWT(token string) *string {
	if strings.Count(token, ".") != 2 {
		return nil
	}

	parsedJWT, err := squarejwt.ParseSigned(token)
	if err != nil {
		return nil
	}

	var claims squarejwt.Claims
	if err = parsedJWT.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return nil
	}

	return &claims.ID
}
