// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/jefferai/jsonx"
)

var (
	_ Formatter        = (*EntryFormatter)(nil)
	_ eventlogger.Node = (*EntryFormatter)(nil)
)

// NewEntryFormatter should be used to create an EntryFormatter.
// Accepted options: WithHeaderFormatter, WithPrefix.
func NewEntryFormatter(config FormatterConfig, salter Salter, opt ...Option) (*EntryFormatter, error) {
	const op = "audit.NewEntryFormatter"

	if salter == nil {
		return nil, fmt.Errorf("%s: cannot create a new audit formatter with nil salter: %w", op, event.ErrInvalidParameter)
	}

	// We need to ensure that the format isn't just some default empty string.
	if err := config.RequiredFormat.validate(); err != nil {
		return nil, fmt.Errorf("%s: format not valid: %w", op, err)
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("%s: error applying options: %w", op, err)
	}

	return &EntryFormatter{
		salter:          salter,
		config:          config,
		headerFormatter: opts.withHeaderFormatter,
		prefix:          opts.withPrefix,
	}, nil
}

// Reopen is a no-op for the formatter node.
func (*EntryFormatter) Reopen() error {
	return nil
}

// Type describes the type of this node (formatter).
func (*EntryFormatter) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFormatter
}

// Process will attempt to parse the incoming event data into a corresponding
// audit Request/Response which is serialized to JSON/JSONx and stored within the event.
func (f *EntryFormatter) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "audit.(EntryFormatter).Process"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, event.ErrInvalidParameter)
	}

	a, ok := e.Payload.(*AuditEvent)
	if !ok {
		return nil, fmt.Errorf("%s: cannot parse event payload: %w", op, event.ErrInvalidParameter)
	}

	var result []byte
	data := new(logical.LogInput)
	headers := make(map[string][]string)

	if a.Data != nil {
		*data = *a.Data
		if a.Data.Request != nil && a.Data.Request.Headers != nil {
			headers = a.Data.Request.Headers
		}
	}

	if f.headerFormatter != nil {
		adjustedHeaders, err := f.headerFormatter.ApplyConfig(ctx, headers, f.salter)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to transform headers for auditing: %w", op, err)
		}

		data.Request.Headers = adjustedHeaders
	}

	switch a.Subtype {
	case RequestType:
		entry, err := f.FormatRequest(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse request from audit event: %w", op, err)
		}

		result, err = jsonutil.EncodeJSON(entry)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to format request: %w", op, err)
		}
	case ResponseType:
		entry, err := f.FormatResponse(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse response from audit event: %w", op, err)
		}

		result, err = jsonutil.EncodeJSON(entry)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to format response: %w", op, err)
		}
	default:
		return nil, fmt.Errorf("%s: unknown audit event subtype: %q", op, a.Subtype)
	}

	if f.config.RequiredFormat == JSONxFormat {
		var err error
		result, err = jsonx.EncodeJSONBytes(result)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to encode JSONx using JSON data: %w", op, err)
		}
		if result == nil {
			return nil, fmt.Errorf("%s: encoded JSONx was nil: %w", op, err)
		}
	}

	// This makes a bit of a mess of the 'format' since both JSON and XML (JSONx)
	// don't support a prefix just sitting there.
	// However, this would be a breaking change to how Vault currently works to
	// include the prefix as part of the JSON object or XML document.
	if f.prefix != "" {
		result = append([]byte(f.prefix), result...)
	}

	// Store the final format.
	e.FormattedAs(f.config.RequiredFormat.String(), result)

	return e, nil
}

// FormatRequest attempts to format the specified logical.LogInput into a RequestEntry.
func (f *EntryFormatter) FormatRequest(ctx context.Context, in *logical.LogInput) (*RequestEntry, error) {
	switch {
	case in == nil || in.Request == nil:
		return nil, errors.New("request to request-audit a nil request")
	case f.salter == nil:
		return nil, errors.New("salt func not configured")
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

	if !f.config.Raw {
		var err error
		auth, err = HashAuth(ctx, f.salter, auth, f.config.HMACAccessor)
		if err != nil {
			return nil, err
		}

		req, err = HashRequest(ctx, f.salter, req, f.config.HMACAccessor, in.NonHMACReqDataKeys)
		if err != nil {
			return nil, err
		}
	}

	var errString string
	if in.OuterErr != nil {
		errString = in.OuterErr.Error()
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	reqType := in.Type
	if reqType == "" {
		reqType = "request"
	}
	reqEntry := &RequestEntry{
		Type:          reqType,
		Error:         errString,
		ForwardedFrom: req.ForwardedFrom,
		Auth: &Auth{
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

		Request: &Request{
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
			Namespace: &Namespace{
				ID:   ns.ID,
				Path: ns.Path,
			},
			Path:                          req.Path,
			Data:                          req.Data,
			PolicyOverride:                req.PolicyOverride,
			RemoteAddr:                    getRemoteAddr(req),
			RemotePort:                    getRemotePort(req),
			ReplicationCluster:            req.ReplicationCluster,
			Headers:                       req.Headers,
			ClientCertificateSerialNumber: getClientCertificateSerialNumber(connState),
		},
	}

	if !auth.IssueTime.IsZero() {
		reqEntry.Auth.TokenIssueTime = auth.IssueTime.Format(time.RFC3339)
	}

	if auth.PolicyResults != nil {
		reqEntry.Auth.PolicyResults = &PolicyResults{
			Allowed: auth.PolicyResults.Allowed,
		}

		for _, p := range auth.PolicyResults.GrantingPolicies {
			reqEntry.Auth.PolicyResults.GrantingPolicies = append(reqEntry.Auth.PolicyResults.GrantingPolicies, PolicyInfo{
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

	if !f.config.OmitTime {
		reqEntry.Time = time.Now().UTC().Format(time.RFC3339Nano)
	}

	return reqEntry, nil
}

// FormatResponse attempts to format the specified logical.LogInput into a ResponseEntry.
func (f *EntryFormatter) FormatResponse(ctx context.Context, in *logical.LogInput) (*ResponseEntry, error) {
	switch {
	case f == nil:
		return nil, errors.New("formatter is nil")
	case in == nil || in.Request == nil:
		return nil, errors.New("request to response-audit a nil request")
	case f.salter == nil:
		return nil, errors.New("salt func not configured")
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

	elideListResponseData := f.config.ElideListResponses && req.Operation == logical.ListOperation

	var respData map[string]interface{}
	if f.config.Raw {
		// In the non-raw case, elision of list response data occurs inside HashResponse, to avoid redundant deep
		// copies and hashing of data only to elide it later. In the raw case, we need to do it here.
		if elideListResponseData && resp.Data != nil {
			// Copy the data map before making changes, but we only need to go one level deep in this case
			respData = make(map[string]interface{}, len(resp.Data))
			for k, v := range resp.Data {
				respData[k] = v
			}

			doElideListResponseData(respData)
		} else {
			respData = resp.Data
		}
	} else {
		var err error
		auth, err = HashAuth(ctx, f.salter, auth, f.config.HMACAccessor)
		if err != nil {
			return nil, err
		}

		req, err = HashRequest(ctx, f.salter, req, f.config.HMACAccessor, in.NonHMACReqDataKeys)
		if err != nil {
			return nil, err
		}

		resp, err = HashResponse(ctx, f.salter, resp, f.config.HMACAccessor, in.NonHMACRespDataKeys, elideListResponseData)
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

	var respAuth *Auth
	if resp.Auth != nil {
		respAuth = &Auth{
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

	var respSecret *Secret
	if resp.Secret != nil {
		respSecret = &Secret{
			LeaseID: resp.Secret.LeaseID,
		}
	}

	var respWrapInfo *ResponseWrapInfo
	if resp.WrapInfo != nil {
		token := resp.WrapInfo.Token
		if jwtToken := parseVaultTokenFromJWT(token); jwtToken != nil {
			token = *jwtToken
		}
		respWrapInfo = &ResponseWrapInfo{
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
	respEntry := &ResponseEntry{
		Type:      respType,
		Error:     errString,
		Forwarded: req.ForwardedFrom != "",
		Auth: &Auth{
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

		Request: &Request{
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
			Namespace: &Namespace{
				ID:   ns.ID,
				Path: ns.Path,
			},
			Path:                          req.Path,
			Data:                          req.Data,
			PolicyOverride:                req.PolicyOverride,
			RemoteAddr:                    getRemoteAddr(req),
			RemotePort:                    getRemotePort(req),
			ClientCertificateSerialNumber: getClientCertificateSerialNumber(connState),
			ReplicationCluster:            req.ReplicationCluster,
			Headers:                       req.Headers,
		},

		Response: &Response{
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
		respEntry.Auth.PolicyResults = &PolicyResults{
			Allowed: auth.PolicyResults.Allowed,
		}

		for _, p := range auth.PolicyResults.GrantingPolicies {
			respEntry.Auth.PolicyResults.GrantingPolicies = append(respEntry.Auth.PolicyResults.GrantingPolicies, PolicyInfo{
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

	if !f.config.OmitTime {
		respEntry.Time = time.Now().UTC().Format(time.RFC3339Nano)
	}

	return respEntry, nil
}

// NewFormatterConfig should be used to create a FormatterConfig.
// Accepted options: WithElision, WithHMACAccessor, WithOmitTime, WithRaw, WithFormat.
func NewFormatterConfig(opt ...Option) (FormatterConfig, error) {
	const op = "audit.NewFormatterConfig"

	opts, err := getOpts(opt...)
	if err != nil {
		return FormatterConfig{}, fmt.Errorf("%s: error applying options: %w", op, err)
	}

	return FormatterConfig{
		ElideListResponses: opts.withElision,
		HMACAccessor:       opts.withHMACAccessor,
		OmitTime:           opts.withOmitTime,
		Raw:                opts.withRaw,
		RequiredFormat:     opts.withFormat,
	}, nil
}

// getRemoteAddr safely gets the remote address avoiding a nil pointer
func getRemoteAddr(req *logical.Request) string {
	if req != nil && req.Connection != nil {
		return req.Connection.RemoteAddr
	}
	return ""
}

// getRemotePort safely gets the remote port avoiding a nil pointer
func getRemotePort(req *logical.Request) int {
	if req != nil && req.Connection != nil {
		return req.Connection.RemotePort
	}
	return 0
}

// getClientCertificateSerialNumber attempts the retrieve the serial number of
// the peer certificate from the specified tls.ConnectionState.
func getClientCertificateSerialNumber(connState *tls.ConnectionState) string {
	if connState == nil || len(connState.VerifiedChains) == 0 || len(connState.VerifiedChains[0]) == 0 {
		return ""
	}

	return connState.VerifiedChains[0][0].SerialNumber.String()
}

// parseVaultTokenFromJWT returns a string iff the token was a JWT, and we could
// extract the original token ID from inside
func parseVaultTokenFromJWT(token string) *string {
	if strings.Count(token, ".") != 2 {
		return nil
	}

	parsedJWT, err := jwt.ParseSigned(token)
	if err != nil {
		return nil
	}

	var claims jwt.Claims
	if err = parsedJWT.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return nil
	}

	return &claims.ID
}

// doElideListResponseData performs the actual elision of list operation response data, once surrounding code has
// determined it should apply to a particular request. The data map that is passed in must be a copy that is safe to
// modify in place, but need not be a full recursive deep copy, as only top-level keys are changed.
//
// See the documentation of the controlling option in FormatterConfig for more information on the purpose.
func doElideListResponseData(data map[string]interface{}) {
	for k, v := range data {
		if k == "keys" {
			if vSlice, ok := v.([]string); ok {
				data[k] = len(vSlice)
			}
		} else if k == "key_info" {
			if vMap, ok := v.(map[string]interface{}); ok {
				data[k] = len(vMap)
			}
		}
	}
}

// newTemporaryEntryFormatter creates a cloned EntryFormatter instance with a non-persistent Salter.
func newTemporaryEntryFormatter(n *EntryFormatter) *EntryFormatter {
	return &EntryFormatter{
		salter:          &nonPersistentSalt{},
		headerFormatter: n.headerFormatter,
		config:          n.config,
		prefix:          n.prefix,
	}
}
