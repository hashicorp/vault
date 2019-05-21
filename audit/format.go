package audit

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	squarejwt "gopkg.in/square/go-jose.v2/jwt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
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

func (f *AuditFormatter) FormatRequest(ctx context.Context, w io.Writer, config FormatterConfig, in *LogInput) error {
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

	req, auth := in.Request, in.Auth
	if auth == nil {
		auth = new(logical.Auth)
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

	reqEntry := &AuditRequestEntry{
		Type:  "request",
		Error: errString,

		Auth: AuditAuth{
			ClientToken:               auth.ClientToken,
			Accessor:                  auth.Accessor,
			DisplayName:               auth.DisplayName,
			Policies:                  auth.Policies,
			TokenPolicies:             auth.TokenPolicies,
			IdentityPolicies:          auth.IdentityPolicies,
			ExternalNamespacePolicies: auth.ExternalNamespacePolicies,
			Metadata:                  auth.Metadata,
			EntityID:                  auth.EntityID,
			RemainingUses:             req.ClientTokenRemainingUses,
			TokenType:                 auth.TokenType.String(),
		},

		Request: AuditRequest{
			ID:                  req.ID,
			ClientToken:         req.ClientToken,
			ClientTokenAccessor: req.ClientTokenAccessor,
			Operation:           req.Operation,
			Namespace: AuditNamespace{
				ID:   ns.ID,
				Path: ns.Path,
			},
			Path:               req.Path,
			Data:               req.Data,
			PolicyOverride:     req.PolicyOverride,
			RemoteAddr:         getRemoteAddr(req),
			ReplicationCluster: req.ReplicationCluster,
			Headers:            req.Headers,
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

func (f *AuditFormatter) FormatResponse(ctx context.Context, w io.Writer, config FormatterConfig, in *LogInput) error {
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
	auth, req, resp := in.Auth, in.Request, in.Response
	if auth == nil {
		auth = new(logical.Auth)
	}
	if resp == nil {
		resp = new(logical.Response)
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

	respEntry := &AuditResponseEntry{
		Type:  "response",
		Error: errString,
		Auth: AuditAuth{
			ClientToken:               auth.ClientToken,
			Accessor:                  auth.Accessor,
			DisplayName:               auth.DisplayName,
			Policies:                  auth.Policies,
			TokenPolicies:             auth.TokenPolicies,
			IdentityPolicies:          auth.IdentityPolicies,
			ExternalNamespacePolicies: auth.ExternalNamespacePolicies,
			Metadata:                  auth.Metadata,
			RemainingUses:             req.ClientTokenRemainingUses,
			EntityID:                  auth.EntityID,
			TokenType:                 auth.TokenType.String(),
		},

		Request: AuditRequest{
			ID:                  req.ID,
			ClientToken:         req.ClientToken,
			ClientTokenAccessor: req.ClientTokenAccessor,
			Operation:           req.Operation,
			Namespace: AuditNamespace{
				ID:   ns.ID,
				Path: ns.Path,
			},
			Path:               req.Path,
			Data:               req.Data,
			PolicyOverride:     req.PolicyOverride,
			RemoteAddr:         getRemoteAddr(req),
			ReplicationCluster: req.ReplicationCluster,
			Headers:            req.Headers,
		},

		Response: AuditResponse{
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
	Time    string       `json:"time,omitempty"`
	Type    string       `json:"type"`
	Auth    AuditAuth    `json:"auth"`
	Request AuditRequest `json:"request"`
	Error   string       `json:"error"`
}

// AuditResponseEntry is the structure of a response audit log entry in Audit.
type AuditResponseEntry struct {
	Time     string        `json:"time,omitempty"`
	Type     string        `json:"type"`
	Auth     AuditAuth     `json:"auth"`
	Request  AuditRequest  `json:"request"`
	Response AuditResponse `json:"response"`
	Error    string        `json:"error"`
}

type AuditRequest struct {
	ID                  string                 `json:"id"`
	ReplicationCluster  string                 `json:"replication_cluster,omitempty"`
	Operation           logical.Operation      `json:"operation"`
	ClientToken         string                 `json:"client_token"`
	ClientTokenAccessor string                 `json:"client_token_accessor"`
	Namespace           AuditNamespace         `json:"namespace"`
	Path                string                 `json:"path"`
	Data                map[string]interface{} `json:"data"`
	PolicyOverride      bool                   `json:"policy_override"`
	RemoteAddr          string                 `json:"remote_address"`
	WrapTTL             int                    `json:"wrap_ttl"`
	Headers             map[string][]string    `json:"headers"`
}

type AuditResponse struct {
	Auth     *AuditAuth             `json:"auth,omitempty"`
	Secret   *AuditSecret           `json:"secret,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Warnings []string               `json:"warnings,omitempty"`
	Redirect string                 `json:"redirect,omitempty"`
	WrapInfo *AuditResponseWrapInfo `json:"wrap_info,omitempty"`
	Headers  map[string][]string    `json:"headers"`
}

type AuditAuth struct {
	ClientToken               string              `json:"client_token"`
	Accessor                  string              `json:"accessor"`
	DisplayName               string              `json:"display_name"`
	Policies                  []string            `json:"policies"`
	TokenPolicies             []string            `json:"token_policies,omitempty"`
	IdentityPolicies          []string            `json:"identity_policies,omitempty"`
	ExternalNamespacePolicies map[string][]string `json:"external_namespace_policies,omitempty"`
	Metadata                  map[string]string   `json:"metadata"`
	NumUses                   int                 `json:"num_uses,omitempty"`
	RemainingUses             int                 `json:"remaining_uses,omitempty"`
	EntityID                  string              `json:"entity_id"`
	TokenType                 string              `json:"token_type"`
}

type AuditSecret struct {
	LeaseID string `json:"lease_id"`
}

type AuditResponseWrapInfo struct {
	TTL             int    `json:"ttl"`
	Token           string `json:"token"`
	Accessor        string `json:"accessor"`
	CreationTime    string `json:"creation_time"`
	CreationPath    string `json:"creation_path"`
	WrappedAccessor string `json:"wrapped_accessor,omitempty"`
}

type AuditNamespace struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// getRemoteAddr safely gets the remote address avoiding a nil pointer
func getRemoteAddr(req *logical.Request) string {
	if req != nil && req.Connection != nil {
		return req.Connection.RemoteAddr
	}
	return ""
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
