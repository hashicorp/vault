// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/sdk/helper/salt"

	"github.com/hashicorp/vault/sdk/logical"
)

// Audit subtypes.
const (
	RequestType  subtype = "AuditRequest"
	ResponseType subtype = "AuditResponse"
)

// Audit formats.
const (
	JSONFormat  format = "json"
	JSONxFormat format = "jsonx"
)

// version defines the version of audit events.
const version = "v0.1"

// subtype defines the type of audit event.
type subtype string

// format defines types of format audit events support.
type format string

// auditEvent is the audit event.
type auditEvent struct {
	ID        string            `json:"id"`
	Version   string            `json:"version"`
	Subtype   subtype           `json:"subtype"` // the subtype of the audit event.
	Timestamp time.Time         `json:"timestamp"`
	Data      *logical.LogInput `json:"data"`
}

// Option is how options are passed as arguments.
type Option func(*options) error

// options are used to represent configuration for a audit related nodes.
type options struct {
	withID              string
	withNow             time.Time
	withSubtype         subtype
	withFormat          format
	withPrefix          string
	withRaw             bool
	withElision         bool
	withOmitTime        bool
	withHMACAccessor    bool
	withHeaderFormatter HeaderFormatter
}

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

// HeaderFormatter is an interface defining the methods of the
// vault.AuditedHeadersConfig structure needed in this package.
type HeaderFormatter interface {
	// ApplyConfig returns a map of header values that consists of the
	// intersection of the provided set of header values with a configured
	// set of headers and will hash headers that have been configured as such.
	ApplyConfig(context.Context, map[string][]string, Salter) (map[string][]string, error)
}

// EntryFormatter should be used to format audit requests and responses.
type EntryFormatter struct {
	salter          Salter
	headerFormatter HeaderFormatter
	config          FormatterConfig
	prefix          string
}

// EntryFormatterWriter should be used to format and write out audit requests and responses.
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

	// The required/target format for the event (supported: JSONFormat and JSONxFormat).
	RequiredFormat format
}

// RequestEntry is the structure of a request audit log entry.
type RequestEntry struct {
	Time          string   `json:"time,omitempty"`
	Type          string   `json:"type,omitempty"`
	Auth          *Auth    `json:"auth,omitempty"`
	Request       *Request `json:"request,omitempty"`
	Error         string   `json:"error,omitempty"`
	ForwardedFrom string   `json:"forwarded_from,omitempty"` // Populated in Enterprise when a request is forwarded
}

// ResponseEntry is the structure of a response audit log entry.
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

// nonPersistentSalt is used for obtaining a salt that is not persisted.
type nonPersistentSalt struct{}

// Backend interface must be implemented for an audit
// mechanism to be made available. Audit backends can be enabled to
// sink information to different backends such as logs, file, databases,
// or other external services.
type Backend interface {
	// Salter interface must be implemented by anything implementing Backend.
	Salter

	// LogRequest is used to synchronously log a request. This is done after the
	// request is authorized but before the request is executed. The arguments
	// MUST not be modified in any way. They should be deep copied if this is
	// a possibility.
	LogRequest(context.Context, *logical.LogInput) error

	// LogResponse is used to synchronously log a response. This is done after
	// the request is processed but before the response is sent. The arguments
	// MUST not be modified in any way. They should be deep copied if this is
	// a possibility.
	LogResponse(context.Context, *logical.LogInput) error

	// LogTestMessage is used to check an audit backend before adding it
	// permanently. It should attempt to synchronously log the given test
	// message, WITHOUT using the normal Salt (which would require a storage
	// operation on creation, which is currently disallowed.)
	LogTestMessage(context.Context, *logical.LogInput, map[string]string) error

	// Reload is called on SIGHUP for supporting backends.
	Reload(context.Context) error

	// Invalidate is called for path invalidation
	Invalidate(context.Context)

	// RegisterNodesAndPipeline provides an eventlogger.Broker pointer so that
	// the Backend can call its RegisterNode and RegisterPipeline methods with
	// the nodes and the pipeline that were created in the corresponding
	// Factory function.
	RegisterNodesAndPipeline(*eventlogger.Broker, string) error
}

// BackendConfig contains configuration parameters used in the factory func to
// instantiate audit backends
type BackendConfig struct {
	// The view to store the salt
	SaltView logical.Storage

	// The salt config that should be used for any secret obfuscation
	SaltConfig *salt.Config

	// Config is the opaque user configuration provided when mounting
	Config map[string]string

	// MountPath is the path where this Backend is mounted
	MountPath string
}

// Factory is the factory function to create an audit backend.
type Factory func(context.Context, *BackendConfig, bool, HeaderFormatter) (Backend, error)
