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
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/jefferai/jsonx"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/pointerstructure"
)

var (
	_ Formatter        = (*EntryFormatter)(nil)
	_ eventlogger.Node = (*EntryFormatter)(nil)
)

// EntryFormatter should be used to format audit requests and responses.
type EntryFormatter struct {
	salter          Salter
	headerFormatter HeaderFormatter
	config          FormatterConfig
	prefix          string
	exclusions      []*exclusion
}

// NewEntryFormatter should be used to create an EntryFormatter.
// Accepted options: WithExclusions, WithHeaderFormatter, WithPrefix.
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
		exclusions:      opts.withExclusions,
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

	// Bail early if the context was cancelled, eventlogger will not carry on asking
	// nodes to process, so any sink node in the pipeline won't be called.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Perform validation on the event, then retrieve the underlying AuditEvent
	// and LogInput (from the AuditEvent Data).
	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, event.ErrInvalidParameter)
	}

	a, ok := e.Payload.(*AuditEvent)
	if !ok {
		return nil, fmt.Errorf("%s: cannot parse event payload: %w", op, event.ErrInvalidParameter)
	}

	if a.Data == nil {
		return nil, fmt.Errorf("%s: cannot audit event (%s) with no data: %w", op, a.Subtype, event.ErrInvalidParameter)
	}

	// Take a copy of the event data before we modify anything.
	data, err := a.Data.Clone()
	if err != nil {
		return nil, fmt.Errorf("%s: unable to copy audit event data: %w", op, err)
	}

	// Ensure that any headers in the request, are formatted as required, and are
	// only present if they have been configured to appear in the audit log.
	// e.g. via: /sys/config/auditing/request-headers/:name
	if f.headerFormatter != nil && data.Request != nil && data.Request.Headers != nil {
		data.Request.Headers, err = f.headerFormatter.ApplyConfig(ctx, data.Request.Headers, f.salter)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to transform headers for auditing: %w", op, err)
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
		entry, err = f.FormatRequest(ctx, data)
	case ResponseType:
		entry, err = f.FormatResponse(ctx, data)
	default:
		return nil, fmt.Errorf("%s: unknown audit event subtype: %q", op, a.Subtype)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: unable to parse %s from audit event: %w", op, a.Subtype.String(), err)
	}

	// If this pipeline has been configured with exclusions then attempt to
	// exclude the fields from the audit entry.
	if len(f.exclusions) > 0 {
		m, err := f.excludeFields(entry)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to exclude audit data from %s: %w", op, a.Subtype.String(), err)
		}

		entry = m
	}

	result, err := jsonutil.EncodeJSON(entry)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to format %s: %w", op, a.Subtype.String(), err)
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

// Salt returns a new salt with default configuration and no storage usage, and no error.
func (s *nonPersistentSalt) Salt(_ context.Context) (*salt.Salt, error) {
	return salt.NewNonpersistentSalt(), nil
}

// excludeFields takes an entry (*RequestEntry or *ResponseEntry) and attempts to
// exclude fields that have been configured on the formatter node.
// NOTE: Whilst the method accepts 'any' the types must be one of the two specified above.
func (f *EntryFormatter) excludeFields(entry any) (map[string]any, error) {
	const op = "audit.(EntryFormatter).excludeFields"

	// Perform some validation on the entry (and its type).
	if entry == nil {
		return nil, fmt.Errorf("%s: entry cannot be nil: %w", op, event.ErrInvalidParameter)
	}

	switch v := entry.(type) {
	case *RequestEntry, *ResponseEntry:
		// These types are expected.
	default:
		return nil, fmt.Errorf("%s: unexpected type: %T", op, v)
	}

	// Decode the entry into a map which can be manipulated as we exclude fields.
	resultMap := make(map[string]any)
	decoder, err := mapDecoderJSON(&resultMap)
	if err != nil {
		return nil, fmt.Errorf("%s: error creating decoder for entry: %w", op, err)
	}
	err = decoder.Decode(entry)
	if err != nil {
		return nil, fmt.Errorf("%s: error decoding entry: %w", op, err)
	}

	// Take another copy which will be the original we use for all evaluations.
	sourceMap := make(map[string]any, len(resultMap))
	decoder, err = mapDecoderJSON(&sourceMap)
	if err != nil {
		return nil, fmt.Errorf("%s: error creating decoder for (source) entry: %w", op, err)
	}
	err = decoder.Decode(entry)
	if err != nil {
		return nil, fmt.Errorf("%s: error decoding (source) entry: %w", op, err)
	}

	for _, exc := range f.exclusions {
		// By default, we want to remove fields, as expression condition is optional.
		shouldRemoveFields := true

		if exc.Evaluator != nil {
			// Decide if we should/shouldn't remove these fields as we have an
			// optional condition expression configured for evaluation.
			shouldRemoveFields, err = exc.Evaluator.Evaluate(sourceMap)
			switch {
			// There may be cases when the evaluator gives us an error, but it's
			// because the datum doesn't currently have the same structure as
			// the expression it was created with.
			// Examples:
			//
			// 1. (ErrNotFound) RequestEntry doesn't have a 'Response' inside of it.
			// So an expression such as "\"/response/mount_type\" == kv" should not
			// cause a failure when we audit a request, which could block audit.
			//
			// 2. (ErrNotFound) Both Request and Response have a Data property which
			// is flexible as it is described as map[string]interface{}, the
			// following expression is valid but wouldn't run without error on
			// every type of audit entry:
			// "\"response/data/my-key\" is not empty".
			//
			// 3. (ErrOutOfRange) Attempting to evaluate a part of an array/slice
			// that doesn't exist shouldn't stop us from auditing, it should just
			// mean that we don't redact fields as the condition failed.
			// For example if we only have a single auth policy but use:
			// "\"/auth/policies/2\" == bar".
			//
			// 4. (ErrConvert) Attempting to use the wrong type of index for the
			// structure. Auth policies are a slice of strings, so you cannot access
			// them via a key:
			// "\"/auth/policies/my-policy\ == bar".
			//
			// 5. (ErrInvalidKind) Attempting to access something that is a different
			// type, for example, mount type is a string, so we cannot access it
			// using a key as if it were a map:
			// "\"/request/mount_type/test\ == bar".
			case errors.Is(err, pointerstructure.ErrNotFound),
				errors.Is(err, pointerstructure.ErrOutOfRange),
				errors.Is(err, pointerstructure.ErrConvert),
				errors.Is(err, pointerstructure.ErrInvalidKind):
				// We can ignore these errors, we won't attempt to exclude fields
				// as the condition failed.
				// shouldRemoveFields will be set to false as a result of the Evaluate call.
			case err != nil:
				return nil, fmt.Errorf("%s: unable to evaluate conditional expression associated with fields: '%s': %w", op, strings.Join(exc.Fields, ", "), err)
			}
		}

		if !shouldRemoveFields {
			continue
		}

		for _, field := range exc.Fields {
			ptr, err := pointerstructure.Parse(field)
			if err != nil {
				return nil, fmt.Errorf("%s: unable to parse field '%s': %w", op, field, err)
			}

			// We don't need the return value as the map is being modified by Delete.
			_, err = ptr.Delete(resultMap)
			// There are some types of errors that may be returned, which we do not
			// consider to mean redaction has failed. We test for these explicitly
			// so any additional error types that may be added to the pointerstructure
			// library in the future can be considered before also being ignored.
			// We do not check for pointerstructure.ErrParse as we have validated this
			// during the application of the WithExclusions Option.
			switch {
			case errors.Is(err, pointerstructure.ErrNotFound),
				errors.Is(err, pointerstructure.ErrOutOfRange),
				errors.Is(err, pointerstructure.ErrConvert),
				errors.Is(err, pointerstructure.ErrInvalidKind):
				fallthrough
			case err == nil:
				continue
			default:
				return nil, fmt.Errorf("%s: unable to exclude field '%s': %w", op, field, err)
			}
		}
	}

	return resultMap, nil
}

// mapDecoderJSON returns a decoder configured to use JSON struct tags output to the
// specified target.
func mapDecoderJSON(target any) (*mapstructure.Decoder, error) {
	const op = "audit.mapDecoderJSON"

	// Configure the decoder to use JSON struct tags to represent mapstructure ones with the same name.
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &target})
	if err != nil {
		return nil, fmt.Errorf("%s: unable to create JSON map decoder: %w", op, err)
	}

	return d, nil
}
