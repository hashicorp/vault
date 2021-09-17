package audit

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strings"
	"time"

	squarejwt "gopkg.in/square/go-jose.v2/jwt"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

type AuditFormatWriter interface {
	// WriteRequest writes the request entry to the writer or returns an error.
	WriteRequest(io.Writer, *AuditRequestEntry) error
	// WriteResponse writes the response entry to the writer or returns an error.
	WriteResponse(io.Writer, *AuditResponseEntry) error
	// Salt returns a non-nil salt or an error.
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
		return fmt.Errorf("error fetching salt: %w", err)
	}

	// Set these to the input values at first
	auth := in.Auth
	req := in.Request
	var connState *tls.ConnectionState
	if auth == nil {
		auth = new(logical.Auth)
	}

	if in.Request.Connection != nil && in.Request.Connection.ConnState != nil {
		connState = in.Request.Connection.ConnState
	}

	if !config.Raw {
		auth, err = HashAuth(salt, auth, config.HMACAccessor)
		if err != nil {
			return err
		}

		req, err = HashRequest(salt, req, config.HMACAccessor, in.NonHMACReqDataKeys)
		if err != nil {
			return err
		}
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
			TokenTTL:                  int64(auth.TTL.Seconds()),
		},

		Request: &AuditRequest{
			ID:                  req.ID,
			ClientToken:         req.ClientToken,
			ClientTokenAccessor: req.ClientTokenAccessor,
			Operation:           req.Operation,
			MountType:           req.MountType,
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

	if !auth.IssueTime.IsZero() {
		reqEntry.Auth.TokenIssueTime = auth.IssueTime.Format(time.RFC3339)
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
		return fmt.Errorf("error fetching salt: %w", err)
	}

	// Set these to the input values at first
	auth, req, resp := in.Auth, in.Request, in.Response
	if auth == nil {
		auth = new(logical.Auth)
	}
	if resp == nil {
		resp = new(logical.Response)
	}
	var connState *tls.ConnectionState

	if in.Request.Connection != nil && in.Request.Connection.ConnState != nil {
		connState = in.Request.Connection.ConnState
	}

	if !config.Raw {
		auth, err = HashAuth(salt, auth, config.HMACAccessor)
		if err != nil {
			return err
		}

		req, err = HashRequest(salt, req, config.HMACAccessor, in.NonHMACReqDataKeys)
		if err != nil {
			return err
		}

		resp, err = HashResponse(salt, resp, config.HMACAccessor, in.NonHMACRespDataKeys)
		if err != nil {
			return err
		}
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
			TokenTTL:                  int64(resp.Auth.TTL.Seconds()),
		}
		if !resp.Auth.IssueTime.IsZero() {
			respAuth.TokenIssueTime = resp.Auth.IssueTime.Format(time.RFC3339)
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
			TokenTTL:                  int64(auth.TTL.Seconds()),
		},

		Request: &AuditRequest{
			ID:                  req.ID,
			ClientToken:         req.ClientToken,
			ClientTokenAccessor: req.ClientTokenAccessor,
			Operation:           req.Operation,
			MountType:           req.MountType,
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
			MountType: req.MountType,
			Auth:      respAuth,
			Secret:    respSecret,
			Data:      resp.Data,
			Warnings:  resp.Warnings,
			Redirect:  resp.Redirect,
			WrapInfo:  respWrapInfo,
			Headers:   resp.Headers,
		},
	}

	if !auth.IssueTime.IsZero() {
		respEntry.Auth.TokenIssueTime = auth.IssueTime.Format(time.RFC3339)
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
	MountType                     string                 `json:"mount_type,omitempty"`
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
	Auth      *AuditAuth             `json:"auth,omitempty"`
	MountType string                 `json:"mount_type,omitempty"`
	Secret    *AuditSecret           `json:"secret,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Warnings  []string               `json:"warnings,omitempty"`
	Redirect  string                 `json:"redirect,omitempty"`
	WrapInfo  *AuditResponseWrapInfo `json:"wrap_info,omitempty"`
	Headers   map[string][]string    `json:"headers,omitempty"`
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
	TokenTTL                  int64               `json:"token_ttl,omitempty"`
	TokenIssueTime            string              `json:"token_issue_time,omitempty"`
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

// Create a formatter not backed by a persistent salt.
func NewTemporaryFormatter(format, prefix string) *AuditFormatter {
	temporarySalt := func(ctx context.Context) (*salt.Salt, error) {
		return salt.NewNonpersistentSalt(), nil
	}
	ret := &AuditFormatter{}

	switch format {
	case "jsonx":
		ret.AuditFormatWriter = &JSONxFormatWriter{
			Prefix:   prefix,
			SaltFunc: temporarySalt,
		}
	default:
		ret.AuditFormatWriter = &JSONFormatWriter{
			Prefix:   prefix,
			SaltFunc: temporarySalt,
		}
	}
	return ret
}
