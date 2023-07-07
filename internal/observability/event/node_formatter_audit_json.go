// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v3/json"

	vaultaudit "github.com/hashicorp/vault/audit" // TODO: this needs to go.

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/sdk/helper/salt"
)

var _ eventlogger.Node = (*AuditFormatterJSON)(nil)

// AuditFormatterJSON represents the formatter node which is used to handle
// formatting audit events as JSON.
type AuditFormatterJSON struct {
	ElideListResponses bool
	HMACAccessor       bool
	OmitTime           bool // for tests
	Raw                bool
	SaltFunc           func(context.Context) (*salt.Salt, error)
	format             auditFormat
}

// AuditFormatterConfig represents configuration that may be required by formatter
// nodes which handle audit events.
type AuditFormatterConfig struct {
	ElideListResponses bool
	HMACAccessor       bool
	OmitTime           bool
	Raw                bool
	SaltFunc           func(context.Context) (*salt.Salt, error)
}

// NewAuditFormatterJSON should be used to create an AuditFormatterJSON.
func NewAuditFormatterJSON(config *AuditFormatterConfig) *AuditFormatterJSON {
	return &AuditFormatterJSON{
		ElideListResponses: config.ElideListResponses,
		HMACAccessor:       config.HMACAccessor,
		OmitTime:           config.OmitTime,
		Raw:                config.Raw,
		SaltFunc:           config.SaltFunc,
		format:             AuditFormatJSON,
	}
}

// Reopen is a no-op for this formatter node.
func (_ *AuditFormatterJSON) Reopen() error {
	return nil
}

// Type describes the type of this node.
func (_ *AuditFormatterJSON) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFormatter
}

// Process will attempt to parse the incoming event data into a corresponding
// audit request/response entry which is serialized to JSON and stored within the event.
func (f *AuditFormatterJSON) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	const op = "event.(AuditFormatterJSON).Process"
	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, ErrInvalidParameter)
	}

	a, ok := e.Payload.(audit)
	if !ok {
		return nil, fmt.Errorf("%s: cannot parse event payload: %w", op, ErrInvalidParameter)
	}

	switch a.Subtype {
	case AuditRequest:
		return f.processRequest(ctx, e)
	case AuditResponse:
		return f.processResponse(ctx, e)
	default:
		return nil, fmt.Errorf("unknown audit event subtype: %q", a.Subtype)
	}
}

// processRequest will parse audit event request data to JSON and store the format
// within the supplied event.
func (f *AuditFormatterJSON) processRequest(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "event.(AuditFormatterJSON).processRequest"

	// TODO: PW: pull out duplication higher up
	a, ok := e.Payload.(audit) // TODO: PW: should this be a pointer, and we can check != nil?
	if !ok {
		return nil, fmt.Errorf("%s: cannot parse event payload: %wd", op, ErrInvalidParameter)
	}

	entry, err := f.parseRequest(ctx, a.Data)
	if err != nil {
		return nil, fmt.Errorf("unable to parse request from audit event: %w", err)
	}

	formatted, err := f.jsonFormat(entry)
	if err != nil {
		return nil, fmt.Errorf("unable to format request: %w", err)
	}

	e.FormattedAs(f.format.String(), formatted)

	return e, nil
}

// processResponse will parse audit event response data to JSON and store the format
// within the supplied event.
func (f *AuditFormatterJSON) processResponse(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "event.(AuditFormatterJSON).processResponse"
	a, ok := e.Payload.(audit) // TODO: PW: should this be a pointer, and we can check != nil?
	if !ok {
		return nil, fmt.Errorf("%s: cannot parse event payload: %wd", op, ErrInvalidParameter)
	}

	entry, err := f.parseResponse(ctx, a.Data)
	if err != nil {
		return nil, fmt.Errorf("unable to parse request from audit event: %w", err)
	}

	formatted, err := f.jsonFormat(entry) // TODO: PW: may be JSONX
	if err != nil {
		return nil, fmt.Errorf("unable to format request: %w", err)
	}

	e.FormattedAs(f.format.String(), formatted)

	return e, nil
}

// parseRequest will map logical.LogInput to a vaultaudit.AuditRequestEntry.
func (f *AuditFormatterJSON) parseRequest(ctx context.Context, input *logical.LogInput) (*vaultaudit.AuditRequestEntry, error) {
	// const op = "event.(AuditFormatterJSON).parseRequest"

	if input == nil || input.Request == nil {
		return nil, fmt.Errorf("request to request-audit a nil request")
	}

	s, err := f.SaltFunc(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching salt: %w", err)
	}

	// Set these to the input values at first
	auth := input.Auth
	req := input.Request
	var connState *tls.ConnectionState
	if auth == nil {
		auth = new(logical.Auth)
	}

	if input.Request.Connection != nil && input.Request.Connection.ConnState != nil {
		connState = input.Request.Connection.ConnState
	}

	if !f.Raw {
		auth, err = vaultaudit.HashAuth(s, auth, f.HMACAccessor)
		if err != nil {
			return nil, err
		}

		req, err = vaultaudit.HashRequest(s, req, f.HMACAccessor, input.NonHMACReqDataKeys)
		if err != nil {
			return nil, err
		}
	}

	var errString string
	if input.OuterErr != nil {
		errString = input.OuterErr.Error()
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	reqType := input.Type
	if reqType == "" {
		reqType = "request"
	}
	reqEntry := &vaultaudit.AuditRequestEntry{
		Type:          reqType,
		Error:         errString,
		ForwardedFrom: req.ForwardedFrom,
		Auth: &vaultaudit.AuditAuth{
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

		Request: &vaultaudit.AuditRequest{
			ID:                    req.ID,
			ClientID:              req.ClientID,
			ClientToken:           req.ClientToken,
			ClientTokenAccessor:   req.ClientTokenAccessor,
			Operation:             req.Operation,
			MountPoint:            req.MountPoint,
			MountType:             req.MountType,
			MountAccessor:         req.MountAccessor,
			MountRunningVersion:   req.MountRunningVersion(),
			MountRunningSha256:    req.MountRunningSha256(),
			MountIsExternalPlugin: req.MountIsExternalPlugin(),
			MountClass:            req.MountClass(),
			Namespace: &vaultaudit.AuditNamespace{
				ID:   ns.ID,
				Path: ns.Path,
			},
			Path:                          req.Path,
			Data:                          req.Data,
			PolicyOverride:                req.PolicyOverride,
			RemoteAddr:                    req.GetRemoteAddr(),
			RemotePort:                    req.GetRemotePort(),
			ReplicationCluster:            req.ReplicationCluster,
			Headers:                       req.Headers,
			ClientCertificateSerialNumber: vaultaudit.GetClientCertificateSerialNumber(connState),
		},
	}

	if !auth.IssueTime.IsZero() {
		reqEntry.Auth.TokenIssueTime = auth.IssueTime.Format(time.RFC3339)
	}

	if auth.PolicyResults != nil {
		reqEntry.Auth.PolicyResults = &vaultaudit.AuditPolicyResults{
			Allowed: auth.PolicyResults.Allowed,
		}

		for _, p := range auth.PolicyResults.GrantingPolicies {
			reqEntry.Auth.PolicyResults.GrantingPolicies = append(reqEntry.Auth.PolicyResults.GrantingPolicies, vaultaudit.PolicyInfo{
				Name:          p.Name,
				NamespaceId:   p.NamespaceId,
				NamespacePath: p.NamespacePath,
				Type:          p.Type,
			})
		}
	}

	if req.WrapInfo != nil {
		reqEntry.Request.WrapTTL = int(req.WrapInfo.TTL / time.Second)
	}

	if !f.OmitTime {
		reqEntry.Time = time.Now().UTC().Format(time.RFC3339Nano)
	}
	return reqEntry, nil
}

// parseResponse will map logical.LogInput to a vaultaudit.AuditResponseEntry.
func (f *AuditFormatterJSON) parseResponse(ctx context.Context, in *logical.LogInput) (*vaultaudit.AuditResponseEntry, error) {
	if in == nil || in.Request == nil {
		return nil, fmt.Errorf("request to response-audit a nil request")
	}

	s, err := f.SaltFunc(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching salt: %w", err)
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

	elideListResponseData := f.ElideListResponses && req.Operation == logical.ListOperation

	var respData map[string]interface{}
	if f.Raw {
		// In the non-raw case, elision of list response data occurs inside HashResponse, to avoid redundant deep
		// copies and hashing of data only to elide it later. In the raw case, we need to do it here.
		if elideListResponseData && resp.Data != nil {
			// Copy the data map before making changes, but we only need to go one level deep in this case
			respData = make(map[string]interface{}, len(resp.Data))
			for k, v := range resp.Data {
				respData[k] = v
			}

			vaultaudit.DoElideListResponseData(respData)
		} else {
			respData = resp.Data
		}
	} else {
		auth, err = auth.HashAuth(s, f.HMACAccessor)
		if err != nil {
			return nil, err
		}

		req, err = vaultaudit.HashRequest(s, req, f.HMACAccessor, in.NonHMACReqDataKeys)
		if err != nil {
			return nil, err
		}

		resp, err = vaultaudit.HashResponse(s, resp, f.HMACAccessor, in.NonHMACRespDataKeys, elideListResponseData)
		if err != nil {
			return nil, err
		}

		respData = resp.Data
	}

	var errString string
	if in.OuterErr != nil {
		errString = in.OuterErr.Error()
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var respAuth *vaultaudit.AuditAuth
	if resp.Auth != nil {
		respAuth = &vaultaudit.AuditAuth{
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

	var respSecret *vaultaudit.AuditSecret
	if resp.Secret != nil {
		respSecret = &vaultaudit.AuditSecret{
			LeaseID: resp.Secret.LeaseID,
		}
	}

	var respWrapInfo *vaultaudit.AuditResponseWrapInfo
	if resp.WrapInfo != nil {
		token := resp.WrapInfo.Token
		if jwtToken := vaultaudit.ParseVaultTokenFromJWT(token); jwtToken != nil {
			token = *jwtToken
		}
		respWrapInfo = &vaultaudit.AuditResponseWrapInfo{
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
	respEntry := &vaultaudit.AuditResponseEntry{
		Type:      respType,
		Error:     errString,
		Forwarded: req.ForwardedFrom != "",
		Auth: &vaultaudit.AuditAuth{
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
			EntityCreated:             auth.EntityCreated,
			TokenType:                 auth.TokenType.String(),
			TokenTTL:                  int64(auth.TTL.Seconds()),
		},

		Request: &vaultaudit.AuditRequest{
			ID:                    req.ID,
			ClientToken:           req.ClientToken,
			ClientTokenAccessor:   req.ClientTokenAccessor,
			ClientID:              req.ClientID,
			Operation:             req.Operation,
			MountPoint:            req.MountPoint,
			MountType:             req.MountType,
			MountAccessor:         req.MountAccessor,
			MountRunningVersion:   req.MountRunningVersion(),
			MountRunningSha256:    req.MountRunningSha256(),
			MountIsExternalPlugin: req.MountIsExternalPlugin(),
			MountClass:            req.MountClass(),
			Namespace: &vaultaudit.AuditNamespace{
				ID:   ns.ID,
				Path: ns.Path,
			},
			Path:                          req.Path,
			Data:                          req.Data,
			PolicyOverride:                req.PolicyOverride,
			RemoteAddr:                    req.GetRemoteAddr(),
			RemotePort:                    req.GetRemotePort(),
			ClientCertificateSerialNumber: vaultaudit.GetClientCertificateSerialNumber(connState),
			ReplicationCluster:            req.ReplicationCluster,
			Headers:                       req.Headers,
		},

		Response: &vaultaudit.AuditResponse{
			MountPoint:            req.MountPoint,
			MountType:             req.MountType,
			MountAccessor:         req.MountAccessor,
			MountRunningVersion:   req.MountRunningVersion(),
			MountRunningSha256:    req.MountRunningSha256(),
			MountIsExternalPlugin: req.MountIsExternalPlugin(),
			MountClass:            req.MountClass(),
			Auth:                  respAuth,
			Secret:                respSecret,
			Data:                  respData,
			Warnings:              resp.Warnings,
			Redirect:              resp.Redirect,
			WrapInfo:              respWrapInfo,
			Headers:               resp.Headers,
		},
	}

	if auth.PolicyResults != nil {
		respEntry.Auth.PolicyResults = &vaultaudit.AuditPolicyResults{
			Allowed: auth.PolicyResults.Allowed,
		}

		for _, p := range auth.PolicyResults.GrantingPolicies {
			respEntry.Auth.PolicyResults.GrantingPolicies = append(respEntry.Auth.PolicyResults.GrantingPolicies, vaultaudit.PolicyInfo{
				Name:          p.Name,
				NamespaceId:   p.NamespaceId,
				NamespacePath: p.NamespacePath,
				Type:          p.Type,
			})
		}
	}

	if !auth.IssueTime.IsZero() {
		respEntry.Auth.TokenIssueTime = auth.IssueTime.Format(time.RFC3339)
	}
	if req.WrapInfo != nil {
		respEntry.Request.WrapTTL = int(req.WrapInfo.TTL / time.Second)
	}

	if !f.OmitTime {
		respEntry.Time = time.Now().UTC().Format(time.RFC3339Nano)
	}

	return respEntry, nil
}

// jsonFormat converts the supplied input to JSON encoded bytes.
func (f *AuditFormatterJSON) jsonFormat(data any) ([]byte, error) {
	if data == nil {
		return nil, errors.New("unable to JSON format data, parameter nil")
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(data)

	return buf.Bytes(), err
}
