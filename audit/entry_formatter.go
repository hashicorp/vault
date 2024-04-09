// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/jefferai/jsonx"
)

var (
	_ Formatter        = (*EntryFormatter)(nil)
	_ eventlogger.Node = (*EntryFormatter)(nil)
)

// timeProvider offers a way to supply a pre-configured time.
type timeProvider interface {
	// formatTime provides the pre-configured time in a particular format.
	formattedTime() string
}

// FormatterConfig is used to provide basic configuration to a formatter.
// Use NewFormatterConfig to initialize the FormatterConfig struct.
type FormatterConfig struct {
	Raw          bool
	HMACAccessor bool

	// Vault lacks pagination in its APIs. As a result, certain list operations can return **very** large responses.
	// The user's chosen audit sinks may experience difficulty consuming audit records that swell to tens of megabytes
	// of JSON. The responses of list operations are typically not very interesting, as they are mostly lists of keys,
	// or, even when they include a "key_info" field, are not returning confidential information. They become even less
	// interesting once HMAC-ed by the audit system.
	//
	// Some example Vault "list" operations that are prone to becoming very large in an active Vault installation are:
	//   auth/token/accessors/
	//   identity/entity/id/
	//   identity/entity-alias/id/
	//   pki/certs/
	//
	// This option exists to provide such users with the option to have response data elided from audit logs, only when
	// the operation type is "list". For added safety, the elision only applies to the "keys" and "key_info" fields
	// within the response data - these are conventionally the only fields present in a list response - see
	// logical.ListResponse, and logical.ListResponseWithInfo. However, other fields are technically possible if a
	// plugin author writes unusual code, and these will be preserved in the audit log even with this option enabled.
	// The elision replaces the values of the "keys" and "key_info" fields with an integer count of the number of
	// entries. This allows even the elided audit logs to still be useful for answering questions like
	// "Was any data returned?" or "How many records were listed?".
	ElideListResponses bool

	// This should only ever be used in a testing context
	OmitTime bool

	// The required/target format for the event (supported: JSONFormat and JSONxFormat).
	RequiredFormat format

	// headerFormatter specifies the formatter used for headers that existing in any incoming audit request.
	headerFormatter HeaderFormatter

	// Prefix specifies a Prefix that should be prepended to any formatted request or response before serialization.
	Prefix string
}

// EntryFormatter should be used to format audit requests and responses.
// NOTE: Use NewEntryFormatter to initialize the EntryFormatter struct.
type EntryFormatter struct {
	config FormatterConfig
	salter Salter
	logger hclog.Logger
	name   string
}

// NewFormatterConfig should be used to create a FormatterConfig.
// Accepted options: WithElision, WithFormat, WithHMACAccessor, WithOmitTime, WithPrefix, WithRaw.
func NewFormatterConfig(headerFormatter HeaderFormatter, opt ...Option) (FormatterConfig, error) {
	if headerFormatter == nil || reflect.ValueOf(headerFormatter).IsNil() {
		return FormatterConfig{}, fmt.Errorf("header formatter is required: %w", ErrInvalidParameter)
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return FormatterConfig{}, fmt.Errorf("error applying options: %w", err)
	}

	return FormatterConfig{
		headerFormatter:    headerFormatter,
		ElideListResponses: opts.withElision,
		HMACAccessor:       opts.withHMACAccessor,
		OmitTime:           opts.withOmitTime,
		Prefix:             opts.withPrefix,
		Raw:                opts.withRaw,
		RequiredFormat:     opts.withFormat,
	}, nil
}

// NewEntryFormatter should be used to create an EntryFormatter.
func NewEntryFormatter(name string, config FormatterConfig, salter Salter, logger hclog.Logger) (*EntryFormatter, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name is required: %w", ErrInvalidParameter)
	}

	if salter == nil {
		return nil, fmt.Errorf("cannot create a new audit formatter with nil salter: %w", ErrInvalidParameter)
	}

	if logger == nil || reflect.ValueOf(logger).IsNil() {
		return nil, fmt.Errorf("cannot create a new audit formatter with nil logger: %w", ErrInvalidParameter)
	}

	return &EntryFormatter{
		config: config,
		salter: salter,
		logger: logger,
		name:   name,
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
func (f *EntryFormatter) Process(ctx context.Context, e *eventlogger.Event) (_ *eventlogger.Event, retErr error) {
	// Return early if the context was cancelled, eventlogger will not carry on
	// asking nodes to process, so any sink node in the pipeline won't be called.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Perform validation on the event, then retrieve the underlying AuditEvent
	// and LogInput (from the AuditEvent Data).
	if e == nil {
		return nil, fmt.Errorf("event is nil: %w", ErrInvalidParameter)
	}

	a, ok := e.Payload.(*AuditEvent)
	if !ok {
		return nil, fmt.Errorf("cannot parse event payload: %w", ErrInvalidParameter)
	}

	if a.Data == nil {
		return nil, fmt.Errorf("cannot audit event (%s) with no data: %w", a.Subtype, ErrInvalidParameter)
	}

	// Handle panics
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		f.logger.Error("panic during logging",
			"request_path", a.Data.Request.Path,
			"audit_device_path", f.name,
			"error", r,
			"stacktrace", string(debug.Stack()))

		// Ensure that we add this error onto any pre-existing error that was being returned.
		retErr = multierror.Append(retErr, fmt.Errorf("panic generating audit log: %q", f.name)).ErrorOrNil()
	}()

	// Take a copy of the event data before we modify anything.
	data, err := a.Data.Clone()
	if err != nil {
		return nil, fmt.Errorf("unable to clone audit event data: %w", err)
	}

	// If the request is present in the input data, apply header configuration
	// regardless. We shouldn't be in a situation where the header formatter isn't
	// present as it's required.
	if data.Request != nil {
		// Ensure that any headers in the request, are formatted as required, and are
		// only present if they have been configured to appear in the audit log.
		// e.g. via: /sys/config/auditing/request-headers/:name
		data.Request.Headers, err = f.config.headerFormatter.ApplyConfig(ctx, data.Request.Headers, f.salter)
		if err != nil {
			return nil, fmt.Errorf("unable to transform headers for auditing: %w", err)
		}
	}

	// If the request contains a Server-Side Consistency Token (SSCT), and we
	// have an auth response, overwrite the existing client token with the SSCT,
	// so that the SSCT appears in the audit log for this entry.
	if data.Request != nil && data.Request.InboundSSCToken != "" && data.Auth != nil {
		data.Auth.ClientToken = data.Request.InboundSSCToken
	}

	// Using 'any' as we have two different types that we can get back from either
	// FormatRequest or FormatResponse, but the JSON encoder doesn't care about types.
	var entry any

	switch a.Subtype {
	case RequestType:
		entry, err = f.FormatRequest(ctx, data, a)
	case ResponseType:
		entry, err = f.FormatResponse(ctx, data, a)
	default:
		return nil, fmt.Errorf("unknown audit event subtype: %q", a.Subtype)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s from audit event: %w", a.Subtype, err)
	}

	result, err := jsonutil.EncodeJSON(entry)
	if err != nil {
		return nil, fmt.Errorf("unable to format %s: %w", a.Subtype, err)
	}

	if f.config.RequiredFormat == JSONxFormat {
		var err error
		result, err = jsonx.EncodeJSONBytes(result)
		if err != nil {
			return nil, fmt.Errorf("unable to encode JSONx using JSON data: %w", err)
		}
		if result == nil {
			return nil, fmt.Errorf("encoded JSONx was nil: %w", err)
		}
	}

	// This makes a bit of a mess of the 'format' since both JSON and XML (JSONx)
	// don't support a prefix just sitting there.
	// However, this would be a breaking change to how Vault currently works to
	// include the prefix as part of the JSON object or XML document.
	if f.config.Prefix != "" {
		result = append([]byte(f.config.Prefix), result...)
	}

	// Copy some properties from the event (and audit event) and store the
	// format for the next (sink) node to Process.
	a2 := &AuditEvent{
		ID:        a.ID,
		Version:   a.Version,
		Subtype:   a.Subtype,
		Timestamp: a.Timestamp,
		Data:      data, // Use the cloned data here rather than a pointer to the original.
	}

	e2 := &eventlogger.Event{
		Type:      e.Type,
		CreatedAt: e.CreatedAt,
		Formatted: make(map[string][]byte), // we are about to set this ourselves.
		Payload:   a2,
	}

	e2.FormattedAs(f.config.RequiredFormat.String(), result)

	return e2, nil
}

// FormatRequest attempts to format the specified logical.LogInput into a RequestEntry.
func (f *EntryFormatter) FormatRequest(ctx context.Context, in *logical.LogInput, provider timeProvider) (*RequestEntry, error) {
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

	if req.HTTPRequest != nil && req.HTTPRequest.RequestURI != req.Path {
		reqEntry.Request.RequestURI = req.HTTPRequest.RequestURI
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
		// Use the time provider to supply the time for this entry.
		reqEntry.Time = provider.formattedTime()
	}

	return reqEntry, nil
}

// FormatResponse attempts to format the specified logical.LogInput into a ResponseEntry.
func (f *EntryFormatter) FormatResponse(ctx context.Context, in *logical.LogInput, provider timeProvider) (*ResponseEntry, error) {
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

	if req.HTTPRequest != nil && req.HTTPRequest.RequestURI != req.Path {
		respEntry.Request.RequestURI = req.HTTPRequest.RequestURI
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
		// Use the time provider to supply the time for this entry.
		respEntry.Time = provider.formattedTime()
	}

	return respEntry, nil
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
		salter: &nonPersistentSalt{},
		config: n.config,
	}
}

// Salt returns a new salt with default configuration and no storage usage, and no error.
func (s *nonPersistentSalt) Salt(_ context.Context) (*salt.Salt, error) {
	return salt.NewNonpersistentSalt(), nil
}
