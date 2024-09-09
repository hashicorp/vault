// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	nshelper "github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/jefferai/jsonx"
	"github.com/mitchellh/copystructure"
)

var _ eventlogger.Node = (*entryFormatter)(nil)

// timeProvider offers a way to supply a pre-configured time.
type timeProvider interface {
	// formatTime provides the pre-configured time in a particular format.
	formattedTime() string
}

// nonPersistentSalt is used for obtaining a salt that is not persisted.
type nonPersistentSalt struct{}

// entryFormatter should be used to format audit requests and responses.
// NOTE: Use newEntryFormatter to initialize the entryFormatter struct.
type entryFormatter struct {
	config formatterConfig
	salter Salter
	logger hclog.Logger
	name   string
}

// newEntryFormatter should be used to create an entryFormatter.
func newEntryFormatter(name string, config formatterConfig, salter Salter, logger hclog.Logger) (*entryFormatter, error) {
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

	return &entryFormatter{
		config: config,
		salter: salter,
		logger: logger,
		name:   name,
	}, nil
}

// Reopen is a no-op for the formatter node.
func (*entryFormatter) Reopen() error {
	return nil
}

// Type describes the type of this node (formatter).
func (*entryFormatter) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFormatter
}

// Process will attempt to parse the incoming event data into a corresponding
// audit request/response which is serialized to JSON/JSONx and stored within the event.
func (f *entryFormatter) Process(ctx context.Context, e *eventlogger.Event) (_ *eventlogger.Event, retErr error) {
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

	a, ok := e.Payload.(*Event)
	if !ok {
		return nil, fmt.Errorf("cannot parse event payload: %w", ErrInvalidParameter)
	}

	if a.Data == nil {
		return nil, fmt.Errorf("cannot audit a '%s' event with no data: %w", a.Subtype, ErrInvalidParameter)
	}

	// Handle panics
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		path := "unknown"
		if a.Data.Request != nil {
			path = a.Data.Request.Path
		}

		f.logger.Error("panic during logging",
			"request_path", path,
			"audit_device_path", f.name,
			"error", r,
			"stacktrace", string(debug.Stack()))

		// Ensure that we add this error onto any pre-existing error that was being returned.
		retErr = errors.Join(retErr, fmt.Errorf("panic generating audit log: %q", f.name))
	}()

	// Using 'any' to make exclusion easier, the JSON encoder doesn't care about types.
	var entry any
	var err error
	entry, err = f.createEntry(ctx, a)
	if err != nil {
		return nil, err
	}

	// If this pipeline has been configured with (Enterprise-only) exclusions then
	// attempt to exclude the fields from the audit entry.
	if f.shouldExclude() {
		m, err := f.excludeFields(entry)
		if err != nil {
			return nil, fmt.Errorf("unable to exclude %s audit data from %q: %w", a.Subtype, f.name, err)
		}

		entry = m
	}

	result, err := jsonutil.EncodeJSON(entry)
	if err != nil {
		return nil, fmt.Errorf("unable to format %s: %w", a.Subtype, err)
	}

	if f.config.requiredFormat == jsonxFormat {
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
	if f.config.prefix != "" {
		result = append([]byte(f.config.prefix), result...)
	}

	// Create a new event, so we can store our formatted data without conflict.
	e2 := &eventlogger.Event{
		Type:      e.Type,
		CreatedAt: e.CreatedAt,
		Formatted: make(map[string][]byte), // we are about to set this ourselves.
		Payload:   a,
	}

	e2.FormattedAs(f.config.requiredFormat.String(), result)

	return e2, nil
}

// remoteAddr safely gets the remote address avoiding a nil pointer.
func remoteAddr(req *logical.Request) string {
	if req != nil && req.Connection != nil {
		return req.Connection.RemoteAddr
	}
	return ""
}

// remotePort safely gets the remote port avoiding a nil pointer.
func remotePort(req *logical.Request) int {
	if req != nil && req.Connection != nil {
		return req.Connection.RemotePort
	}
	return 0
}

// clientCertSerialNumber attempts the retrieve the serial number of the peer
// certificate from the specified tls.ConnectionState.
func clientCertSerialNumber(req *logical.Request) string {
	if req == nil || req.Connection == nil {
		return ""
	}

	connState := req.Connection.ConnState

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

// newTemporaryEntryFormatter creates a cloned entryFormatter instance with a non-persistent Salter.
func newTemporaryEntryFormatter(n *entryFormatter) *entryFormatter {
	return &entryFormatter{
		salter: &nonPersistentSalt{},
		config: n.config,
	}
}

// Salt returns a new salt with default configuration and no storage usage, and no error.
func (s *nonPersistentSalt) Salt(_ context.Context) (*salt.Salt, error) {
	return salt.NewNonpersistentSalt(), nil
}

// clone can be used to deep clone the specified type.
func clone[V any](s V) (V, error) {
	s2, err := copystructure.Copy(s)

	return s2.(V), err
}

// newAuth takes a logical.Auth and the number of remaining client token uses
// (which should be supplied from the logical.Request's client token), and creates
// an audit auth.
// tokenRemainingUses should be the client token remaining uses to include in auth.
// This usually can be found in logical.Request.ClientTokenRemainingUses.
// NOTE: supplying a nil value for auth will result in a nil return value and
// (nil) error. The caller should check the return value before attempting to use it.
// ignore-nil-nil-function-check.
func newAuth(input *logical.Auth, tokenRemainingUses int) (*auth, error) {
	if input == nil {
		return nil, nil
	}

	extNSPolicies, err := clone(input.ExternalNamespacePolicies)
	if err != nil {
		return nil, fmt.Errorf("unable to clone logical auth: external namespace policies: %w", err)
	}

	identityPolicies, err := clone(input.IdentityPolicies)
	if err != nil {
		return nil, fmt.Errorf("unable to clone logical auth: identity policies: %w", err)
	}

	metadata, err := clone(input.Metadata)
	if err != nil {
		return nil, fmt.Errorf("unable to clone logical auth: metadata: %w", err)
	}

	policies, err := clone(input.Policies)
	if err != nil {
		return nil, fmt.Errorf("unable to clone logical auth: policies: %w", err)
	}

	var polRes *policyResults
	if input.PolicyResults != nil {
		polRes = &policyResults{
			Allowed:          input.PolicyResults.Allowed,
			GrantingPolicies: make([]policyInfo, len(input.PolicyResults.GrantingPolicies)),
		}

		for _, p := range input.PolicyResults.GrantingPolicies {
			polRes.GrantingPolicies = append(polRes.GrantingPolicies, policyInfo{
				Name:          p.Name,
				NamespaceId:   p.NamespaceId,
				NamespacePath: p.NamespacePath,
				Type:          p.Type,
			})
		}
	}

	tokenPolicies, err := clone(input.TokenPolicies)
	if err != nil {
		return nil, fmt.Errorf("unable to clone logical auth: token policies: %w", err)
	}

	var tokenIssueTime string
	if !input.IssueTime.IsZero() {
		tokenIssueTime = input.IssueTime.Format(time.RFC3339)
	}

	return &auth{
		Accessor:                  input.Accessor,
		ClientToken:               input.ClientToken,
		DisplayName:               input.DisplayName,
		EntityCreated:             input.EntityCreated,
		EntityID:                  input.EntityID,
		ExternalNamespacePolicies: extNSPolicies,
		IdentityPolicies:          identityPolicies,
		Metadata:                  metadata,
		NoDefaultPolicy:           input.NoDefaultPolicy,
		NumUses:                   input.NumUses,
		Policies:                  policies,
		PolicyResults:             polRes,
		RemainingUses:             tokenRemainingUses,
		TokenPolicies:             tokenPolicies,
		TokenIssueTime:            tokenIssueTime,
		TokenTTL:                  int64(input.TTL.Seconds()),
		TokenType:                 input.TokenType.String(),
	}, nil
}

// newRequest takes a logical.Request and namespace.Namespace, transforms and
// aggregates them into an audit request.
func newRequest(req *logical.Request, ns *nshelper.Namespace) (*request, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	remoteAddr := remoteAddr(req)
	remotePort := remotePort(req)
	clientCertSerial := clientCertSerialNumber(req)

	data, err := clone(req.Data)
	if err != nil {
		return nil, fmt.Errorf("unable to clone logical request: data: %w", err)
	}

	headers, err := clone(req.Headers)
	if err != nil {
		return nil, fmt.Errorf("unable to clone logical request: headers: %w", err)
	}

	var reqURI string
	if req.HTTPRequest != nil && req.HTTPRequest.RequestURI != req.Path {
		reqURI = req.HTTPRequest.RequestURI
	}
	var wrapTTL int
	if req.WrapInfo != nil {
		wrapTTL = int(req.WrapInfo.TTL / time.Second)
	}

	return &request{
		ClientCertificateSerialNumber: clientCertSerial,
		ClientID:                      req.ClientID,
		ClientToken:                   req.ClientToken,
		ClientTokenAccessor:           req.ClientTokenAccessor,
		Data:                          data,
		Headers:                       headers,
		ID:                            req.ID,
		MountAccessor:                 req.MountAccessor,
		MountClass:                    req.MountClass(),
		MountIsExternalPlugin:         req.MountIsExternalPlugin(),
		MountPoint:                    req.MountPoint,
		MountRunningSha256:            req.MountRunningSha256(),
		MountRunningVersion:           req.MountRunningVersion(),
		MountType:                     req.MountType,
		Namespace: &namespace{
			ID:   ns.ID,
			Path: ns.Path,
		},
		Operation:          req.Operation,
		Path:               req.Path,
		PolicyOverride:     req.PolicyOverride,
		RemoteAddr:         remoteAddr,
		RemotePort:         remotePort,
		ReplicationCluster: req.ReplicationCluster,
		RequestURI:         reqURI,
		WrapTTL:            wrapTTL,
	}, nil
}

// newResponse takes a logical.Response and logical.Request, transforms and
// aggregates them into an audit response.
// isElisionRequired is used to indicate that response 'Data' should be elided.
// NOTE: supplying a nil value for response will result in a nil return value and
// (nil) error. The caller should check the return value before attempting to use it.
// ignore-nil-nil-function-check.
func newResponse(resp *logical.Response, req *logical.Request, isElisionRequired bool) (*response, error) {
	if resp == nil {
		return nil, nil
	}

	if req == nil {
		// Request should never be nil, even for a response.
		return nil, fmt.Errorf("request cannot be nil")
	}

	auth, err := newAuth(resp.Auth, req.ClientTokenRemainingUses)
	if err != nil {
		return nil, fmt.Errorf("unable to convert logical auth response: %w", err)
	}

	var data map[string]any
	if resp.Data != nil {
		data = make(map[string]any, len(resp.Data))

		if isElisionRequired {
			// Performs the actual elision (ideally for list operations) of response data,
			// once surrounding code has determined it should apply to a particular request.
			// If the value for a key should not be elided, then it will be cloned.
			for k, v := range resp.Data {
				isCloneRequired := true
				switch k {
				case "keys":
					if vSlice, ok := v.([]string); ok {
						data[k] = len(vSlice)
						isCloneRequired = false
					}
				case "key_info":
					if vMap, ok := v.(map[string]any); ok {
						data[k] = len(vMap)
						isCloneRequired = false
					}
				}

				// Clone values if they weren't legitimate keys or key_info.
				if isCloneRequired {
					v2, err := clone(v)
					if err != nil {
						return nil, fmt.Errorf("unable to clone response data while eliding: %w", err)
					}
					data[k] = v2
				}
			}
		} else {
			// Deep clone all values, no shortcuts here.
			data, err = clone(resp.Data)
			if err != nil {
				return nil, fmt.Errorf("unable to clone response data: %w", err)
			}
		}
	}

	headers, err := clone(resp.Headers)
	if err != nil {
		return nil, fmt.Errorf("unable to clone logical response: headers: %w", err)
	}

	var s *secret
	if resp.Secret != nil {
		s = &secret{LeaseID: resp.Secret.LeaseID}
	}

	var wrapInfo *responseWrapInfo
	if resp.WrapInfo != nil {
		token := resp.WrapInfo.Token
		if jwtToken := parseVaultTokenFromJWT(token); jwtToken != nil {
			token = *jwtToken
		}

		ttl := int(resp.WrapInfo.TTL / time.Second)
		wrapInfo = &responseWrapInfo{
			TTL:             ttl,
			Token:           token,
			Accessor:        resp.WrapInfo.Accessor,
			CreationTime:    resp.WrapInfo.CreationTime.UTC().Format(time.RFC3339Nano),
			CreationPath:    resp.WrapInfo.CreationPath,
			WrappedAccessor: resp.WrapInfo.WrappedAccessor,
		}
	}

	warnings, err := clone(resp.Warnings)
	if err != nil {
		return nil, fmt.Errorf("unable to clone logical response: warnings: %w", err)
	}

	return &response{
		Auth:                  auth,
		Data:                  data,
		Headers:               headers,
		MountAccessor:         req.MountAccessor,
		MountClass:            req.MountClass(),
		MountIsExternalPlugin: req.MountIsExternalPlugin(),
		MountPoint:            req.MountPoint,
		MountRunningSha256:    req.MountRunningSha256(),
		MountRunningVersion:   req.MountRunningVersion(),
		MountType:             req.MountType,
		Redirect:              resp.Redirect,
		Secret:                s,
		WrapInfo:              wrapInfo,
		Warnings:              warnings,
	}, nil
}

// createEntry takes the AuditEvent and builds an audit entry.
// The entry will be HMAC'd and elided where required.
func (f *entryFormatter) createEntry(ctx context.Context, a *Event) (*entry, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}

	data := a.Data

	if data.Request == nil {
		// Request should never be nil, even for a response.
		return nil, fmt.Errorf("unable to parse request from '%s' audit event: request cannot be nil", a.Subtype)
	}

	ns, err := nshelper.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve namespace from context: %w", err)
	}

	auth, err := newAuth(data.Auth, data.Request.ClientTokenRemainingUses)
	if err != nil {
		return nil, fmt.Errorf("cannot convert auth: %w", err)
	}

	req, err := newRequest(data.Request, ns)
	if err != nil {
		return nil, fmt.Errorf("cannot convert request: %w", err)
	}

	var resp *response
	if a.Subtype == ResponseType {
		shouldElide := f.config.elideListResponses && req.Operation == logical.ListOperation
		resp, err = newResponse(data.Response, data.Request, shouldElide)
		if err != nil {
			return nil, fmt.Errorf("cannot convert response: %w", err)
		}
	}

	var outerErr string
	if data.OuterErr != nil {
		outerErr = data.OuterErr.Error()
	}

	entryType := data.Type
	if entryType == "" {
		entryType = a.Subtype.String()
	}

	entry := &entry{
		Auth:          auth,
		Error:         outerErr,
		Forwarded:     false,
		ForwardedFrom: data.Request.ForwardedFrom,
		Request:       req,
		Response:      resp,
		Type:          entryType,
	}

	if !f.config.omitTime {
		// Use the time provider to supply the time for this entry.
		entry.Time = a.timeProvider().formattedTime()
	}

	// If the request is present in the input data, apply header configuration
	// regardless. We shouldn't be in a situation where the header formatter isn't
	// present as it's required.
	if entry.Request != nil {
		// Ensure that any headers in the request, are formatted as required, and are
		// only present if they have been configured to appear in the audit log.
		// e.g. via: /sys/config/auditing/request-headers/:name
		entry.Request.Headers, err = f.config.headerFormatter.ApplyConfig(ctx, entry.Request.Headers, f.salter)
		if err != nil {
			return nil, fmt.Errorf("unable to transform headers for auditing: %w", err)
		}
	}

	// If the request contains a Server-Side Consistency Token (SSCT), and we
	// have an auth response, overwrite the existing client token with the SSCT,
	// so that the SSCT appears in the audit log for this entry.
	if data.Request != nil && data.Request.InboundSSCToken != "" && entry.Auth != nil {
		entry.Auth.ClientToken = data.Request.InboundSSCToken
	}

	// Hash the entry if we aren't expecting raw output.
	if !f.config.raw {
		// Requests and responses have auth and request.
		err = hashAuth(ctx, f.salter, entry.Auth, f.config.hmacAccessor)
		if err != nil {
			return nil, err
		}

		err = hashRequest(ctx, f.salter, entry.Request, f.config.hmacAccessor, data.NonHMACReqDataKeys)
		if err != nil {
			return nil, err
		}

		if a.Subtype == ResponseType {
			if err = hashResponse(ctx, f.salter, entry.Response, f.config.hmacAccessor, data.NonHMACRespDataKeys); err != nil {
				return nil, err
			}
		}
	}

	return entry, nil
}
