// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	_ Formatter = (*EntryFormatterWriter)(nil)
	_ Writer    = (*EntryFormatterWriter)(nil)
)

// Salter is an interface that provides a way to obtain a Salt for hashing.
type Salter interface {
	// Salt returns a non-nil salt or an error.
	Salt(context.Context) (*salt.Salt, error)
}

// Formatter is an interface that is responsible for formatting a request/response into some format.
// It is recommended that you pass data through Hash prior to formatting it.
type Formatter interface {
	// FormatRequest formats the logical.LogInput into an RequestEntry.
	FormatRequest(context.Context, *logical.LogInput) (*RequestEntry, error)
	// FormatResponse formats the logical.LogInput into an ResponseEntry.
	FormatResponse(context.Context, *logical.LogInput) (*ResponseEntry, error)
}

// Writer is an interface that provides a way to write request and response audit entries.
// Formatters write their output to an io.Writer.
type Writer interface {
	// WriteRequest writes the request entry to the writer or returns an error.
	WriteRequest(io.Writer, *RequestEntry) error
	// WriteResponse writes the response entry to the writer or returns an error.
	WriteResponse(io.Writer, *ResponseEntry) error
}

// EntryFormatter should be used to format audit entries.
type EntryFormatter struct {
	salter Salter
	config FormatterConfig
	prefix string
}

// EntryFormatterWriter should be used to format and write out audit entries.
type EntryFormatterWriter struct {
	Formatter
	Writer
	config FormatterConfig
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

	// The required/target format for the audit entry (supported: JSONFormat and JSONxFormat).
	RequiredFormat format
}

// nonPersistentSalt is used for obtaining a salt that is
type nonPersistentSalt struct{}

type auditEvent struct {
	ID             string            `json:"id"`
	Version        string            `json:"version"`
	Subtype        subtype           `json:"subtype"` // the subtype of the audit event.
	Timestamp      time.Time         `json:"timestamp"`
	Data           *logical.LogInput `json:"data"`
	RequiredFormat format            `json:"format"`
}

// Salt returns a new salt with default configuration and no storage usage, and no error.
func (s *nonPersistentSalt) Salt(_ context.Context) (*salt.Salt, error) {
	return salt.NewNonpersistentSalt(), nil
}

// NewEntryFormatterWriter should be used to create a new EntryFormatterWriter.
// Deprecated: Please move to using eventlogger.Event via EntryFormatter and a sink.
func NewEntryFormatterWriter(config FormatterConfig, formatter Formatter, writer Writer) (*EntryFormatterWriter, error) {
	switch {
	case formatter == nil:
		return nil, errors.New("cannot create a new audit formatter writer with nil formatter")
	case writer == nil:
		return nil, errors.New("cannot create a new audit formatter writer with nil formatter")
	}

	fw := &EntryFormatterWriter{
		Formatter: formatter,
		Writer:    writer,
		config:    config,
	}

	return fw, nil
}

// FormatAndWriteRequest attempts to format the specified logical.LogInput into an RequestEntry,
// and then write the request using the specified io.Writer.
// Deprecated: Please move to using eventlogger.Event via EntryFormatter and a sink.
func (f *EntryFormatterWriter) FormatAndWriteRequest(ctx context.Context, w io.Writer, in *logical.LogInput) error {
	switch {
	case in == nil || in.Request == nil:
		return fmt.Errorf("request to request-audit a nil request")
	case w == nil:
		return fmt.Errorf("writer for audit request is nil")
	case f.Formatter == nil:
		return fmt.Errorf("no formatter specifed")
	case f.Writer == nil:
		return fmt.Errorf("no writer specified")
	}

	reqEntry, err := f.Formatter.FormatRequest(ctx, in)
	if err != nil {
		return err
	}

	return f.Writer.WriteRequest(w, reqEntry)
}

// FormatAndWriteResponse attempts to format the specified logical.LogInput into an ResponseEntry,
// and then write the response using the specified io.Writer.
// Deprecated: Please move to using eventlogger.Event via EntryFormatter and a sink.
func (f *EntryFormatterWriter) FormatAndWriteResponse(ctx context.Context, w io.Writer, in *logical.LogInput) error {
	switch {
	case in == nil || in.Request == nil:
		return errors.New("request to response-audit a nil request")
	case w == nil:
		return errors.New("writer for audit request is nil")
	case f.Formatter == nil:
		return errors.New("no formatter specified")
	case f.Writer == nil:
		return errors.New("no writer specified")
	}

	respEntry, err := f.FormatResponse(ctx, in)
	if err != nil {
		return err
	}

	return f.Writer.WriteResponse(w, respEntry)
}

// RequestEntry is the structure of a request audit log entry in Audit.
type RequestEntry struct {
	Time          string   `json:"time,omitempty"`
	Type          string   `json:"type,omitempty"`
	Auth          *Auth    `json:"auth,omitempty"`
	Request       *Request `json:"request,omitempty"`
	Error         string   `json:"error,omitempty"`
	ForwardedFrom string   `json:"forwarded_from,omitempty"` // Populated in Enterprise when a request is forwarded
}

// ResponseEntry is the structure of a response audit log entry in Audit.
type ResponseEntry struct {
	Time      string    `json:"time,omitempty"`
	Type      string    `json:"type,omitempty"`
	Auth      *Auth     `json:"auth,omitempty"`
	Request   *Request  `json:"request,omitempty"`
	Response  *Response `json:"response,omitempty"`
	Error     string    `json:"error,omitempty"`
	Forwarded bool      `json:"forwarded,omitempty"`
}

type Request struct {
	ID                            string                 `json:"id,omitempty"`
	ClientID                      string                 `json:"client_id,omitempty"`
	ReplicationCluster            string                 `json:"replication_cluster,omitempty"`
	Operation                     logical.Operation      `json:"operation,omitempty"`
	MountPoint                    string                 `json:"mount_point,omitempty"`
	MountType                     string                 `json:"mount_type,omitempty"`
	MountAccessor                 string                 `json:"mount_accessor,omitempty"`
	MountRunningVersion           string                 `json:"mount_running_version,omitempty"`
	MountRunningSha256            string                 `json:"mount_running_sha256,omitempty"`
	MountClass                    string                 `json:"mount_class,omitempty"`
	MountIsExternalPlugin         bool                   `json:"mount_is_external_plugin,omitempty"`
	ClientToken                   string                 `json:"client_token,omitempty"`
	ClientTokenAccessor           string                 `json:"client_token_accessor,omitempty"`
	Namespace                     *Namespace             `json:"namespace,omitempty"`
	Path                          string                 `json:"path,omitempty"`
	Data                          map[string]interface{} `json:"data,omitempty"`
	PolicyOverride                bool                   `json:"policy_override,omitempty"`
	RemoteAddr                    string                 `json:"remote_address,omitempty"`
	RemotePort                    int                    `json:"remote_port,omitempty"`
	WrapTTL                       int                    `json:"wrap_ttl,omitempty"`
	Headers                       map[string][]string    `json:"headers,omitempty"`
	ClientCertificateSerialNumber string                 `json:"client_certificate_serial_number,omitempty"`
}

type Response struct {
	Auth                  *Auth                  `json:"auth,omitempty"`
	MountPoint            string                 `json:"mount_point,omitempty"`
	MountType             string                 `json:"mount_type,omitempty"`
	MountAccessor         string                 `json:"mount_accessor,omitempty"`
	MountRunningVersion   string                 `json:"mount_running_plugin_version,omitempty"`
	MountRunningSha256    string                 `json:"mount_running_sha256,omitempty"`
	MountClass            string                 `json:"mount_class,omitempty"`
	MountIsExternalPlugin bool                   `json:"mount_is_external_plugin,omitempty"`
	Secret                *Secret                `json:"secret,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	Warnings              []string               `json:"warnings,omitempty"`
	Redirect              string                 `json:"redirect,omitempty"`
	WrapInfo              *ResponseWrapInfo      `json:"wrap_info,omitempty"`
	Headers               map[string][]string    `json:"headers,omitempty"`
}

type Auth struct {
	ClientToken               string              `json:"client_token,omitempty"`
	Accessor                  string              `json:"accessor,omitempty"`
	DisplayName               string              `json:"display_name,omitempty"`
	Policies                  []string            `json:"policies,omitempty"`
	TokenPolicies             []string            `json:"token_policies,omitempty"`
	IdentityPolicies          []string            `json:"identity_policies,omitempty"`
	ExternalNamespacePolicies map[string][]string `json:"external_namespace_policies,omitempty"`
	NoDefaultPolicy           bool                `json:"no_default_policy,omitempty"`
	PolicyResults             *PolicyResults      `json:"policy_results,omitempty"`
	Metadata                  map[string]string   `json:"metadata,omitempty"`
	NumUses                   int                 `json:"num_uses,omitempty"`
	RemainingUses             int                 `json:"remaining_uses,omitempty"`
	EntityID                  string              `json:"entity_id,omitempty"`
	EntityCreated             bool                `json:"entity_created,omitempty"`
	TokenType                 string              `json:"token_type,omitempty"`
	TokenTTL                  int64               `json:"token_ttl,omitempty"`
	TokenIssueTime            string              `json:"token_issue_time,omitempty"`
}

type PolicyResults struct {
	Allowed          bool         `json:"allowed"`
	GrantingPolicies []PolicyInfo `json:"granting_policies,omitempty"`
}

type PolicyInfo struct {
	Name          string `json:"name,omitempty"`
	NamespaceId   string `json:"namespace_id,omitempty"`
	NamespacePath string `json:"namespace_path,omitempty"`
	Type          string `json:"type"`
}

type Secret struct {
	LeaseID string `json:"lease_id,omitempty"`
}

type ResponseWrapInfo struct {
	TTL             int    `json:"ttl,omitempty"`
	Token           string `json:"token,omitempty"`
	Accessor        string `json:"accessor,omitempty"`
	CreationTime    string `json:"creation_time,omitempty"`
	CreationPath    string `json:"creation_path,omitempty"`
	WrappedAccessor string `json:"wrapped_accessor,omitempty"`
}

type Namespace struct {
	ID   string `json:"id,omitempty"`
	Path string `json:"path,omitempty"`
}

// NewTemporaryFormatter creates a formatter not backed by a persistent salt
func NewTemporaryFormatter(requiredFormat, prefix string) (*EntryFormatterWriter, error) {
	cfg, err := NewFormatterConfig(WithFormat(requiredFormat))
	if err != nil {
		return nil, err
	}

	eventFormatter, err := NewEntryFormatter(cfg, &nonPersistentSalt{}, WithPrefix(prefix))
	if err != nil {
		return nil, err
	}

	var w Writer

	switch {
	case strings.EqualFold(requiredFormat, JSONxFormat.String()):
		w = &JSONxWriter{Prefix: prefix}
	default:
		w = &JSONWriter{Prefix: prefix}
	}

	fw, err := NewEntryFormatterWriter(cfg, eventFormatter, w)
	if err != nil {
		return nil, err
	}

	return fw, nil
}
