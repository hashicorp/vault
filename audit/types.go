// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"

	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

// Backend interface must be implemented for an audit
// mechanism to be made available. Audit backends can be enabled to
// sink information to different backends such as logs, file, databases,
// or other external services.
type Backend interface {
	// Salter interface must be implemented by anything implementing Backend.
	Salter

	// The PipelineReader interface allows backends to surface information about their
	// nodes for node and pipeline registration.
	event.PipelineReader

	// IsFallback can be used to determine if this audit backend device is intended to
	// be used as a fallback to catch all events that are not written when only using
	// filtered pipelines.
	IsFallback() bool

	// LogTestMessage is used to check an audit backend before adding it
	// permanently. It should attempt to synchronously log the given test
	// message, WITHOUT using the normal Salt (which would require a storage
	// operation on creation, which is currently disallowed.)
	LogTestMessage(context.Context, *logical.LogInput) error

	// Reload is called on SIGHUP for supporting backends.
	Reload(context.Context) error

	// Invalidate is called for path invalidation
	Invalidate(context.Context)
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

// HeaderFormatter is an interface defining the methods of the
// vault.AuditedHeadersConfig structure needed in this package.
type HeaderFormatter interface {
	// ApplyConfig returns a map of header values that consists of the
	// intersection of the provided set of header values with a configured
	// set of headers and will hash headers that have been configured as such.
	ApplyConfig(context.Context, map[string][]string, Salter) (map[string][]string, error)
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
	Time          string   `mapstructure:"time,omitempty"           json:"time,omitempty"`
	Type          string   `mapstructure:"type,omitempty"           json:"type,omitempty"`
	Auth          *Auth    `mapstructure:"auth,omitempty"           json:"auth,omitempty"`
	Request       *Request `mapstructure:"request,omitempty"        json:"request,omitempty"`
	Error         string   `mapstructure:"error,omitempty"          json:"error,omitempty"`
	ForwardedFrom string   `mapstructure:"forwarded_from,omitempty" json:"forwarded_from,omitempty"` // Populated in Enterprise when a request is forwarded
}

// ResponseEntry is the structure of a response audit log entry.
type ResponseEntry struct {
	Time      string    `mapstructure:"time,omitempty" json:"time,omitempty"`
	Type      string    `mapstructure:"type,omitempty" json:"type,omitempty"`
	Auth      *Auth     `mapstructure:"auth,omitempty" json:"auth,omitempty"`
	Request   *Request  `mapstructure:"request,omitempty" json:"request,omitempty"`
	Response  *Response `mapstructure:"response,omitempty" json:"response,omitempty"`
	Error     string    `mapstructure:"error,omitempty" json:"error,omitempty"`
	Forwarded bool      `mapstructure:"forwarded,omitempty" json:"forwarded,omitempty"`
}

type Request struct {
	ID                            string                 `mapstructure:"id,omitempty" json:"id,omitempty"`
	ClientID                      string                 `mapstructure:"client_id,omitempty" json:"client_id,omitempty"`
	ReplicationCluster            string                 `mapstructure:"replication_cluster,omitempty" json:"replication_cluster,omitempty"`
	Operation                     logical.Operation      `mapstructure:"operation,omitempty" json:"operation,omitempty"`
	MountPoint                    string                 `mapstructure:"mount_point,omitempty" json:"mount_point,omitempty"`
	MountType                     string                 `mapstructure:"mount_type,omitempty" json:"mount_type,omitempty"`
	MountAccessor                 string                 `mapstructure:"mount_accessor,omitempty" json:"mount_accessor,omitempty"`
	MountRunningVersion           string                 `mapstructure:"mount_running_version,omitempty" json:"mount_running_version,omitempty"`
	MountRunningSha256            string                 `mapstructure:"mount_running_sha256,omitempty" json:"mount_running_sha256,omitempty"`
	MountClass                    string                 `mapstructure:"mount_class,omitempty" json:"mount_class,omitempty"`
	MountIsExternalPlugin         bool                   `mapstructure:"mount_is_external_plugin,omitempty" json:"mount_is_external_plugin,omitempty"`
	ClientToken                   string                 `mapstructure:"client_token,omitempty" json:"client_token,omitempty"`
	ClientTokenAccessor           string                 `mapstructure:"client_token_accessor,omitempty" json:"client_token_accessor,omitempty"`
	Namespace                     *Namespace             `mapstructure:"namespace,omitempty"  json:"namespace,omitempty"`
	Path                          string                 `mapstructure:"path,omitempty" json:"path,omitempty"`
	Data                          map[string]interface{} `mapstructure:"data,omitempty" json:"data,omitempty"`
	PolicyOverride                bool                   `mapstructure:"policy_override,omitempty" json:"policy_override,omitempty"`
	RemoteAddr                    string                 `mapstructure:"remote_address,omitempty" json:"remote_address,omitempty"`
	RemotePort                    int                    `mapstructure:"remote_port,omitempty" json:"remote_port,omitempty"`
	WrapTTL                       int                    `mapstructure:"wrap_ttl,omitempty" json:"wrap_ttl,omitempty"`
	Headers                       map[string][]string    `mapstructure:"headers,omitempty" json:"headers,omitempty"`
	ClientCertificateSerialNumber string                 `mapstructure:"client_certificate_serial_number,omitempty" json:"client_certificate_serial_number,omitempty"`
	RequestURI                    string                 `mapstructure:"request_uri,omitempty" json:"request_uri,omitempty"`
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
	ClientToken               string              `mapstructure:"client_token,omitempty" json:"client_token,omitempty"`
	Accessor                  string              `mapstructure:"accessor,omitempty" json:"accessor,omitempty"`
	DisplayName               string              `mapstructure:"display_name,omitempty" json:"display_name,omitempty"`
	Policies                  []string            `mapstructure:"policies,omitempty" json:"policies,omitempty"`
	TokenPolicies             []string            `mapstructure:"token_policies,omitempty" json:"token_policies,omitempty"`
	IdentityPolicies          []string            `mapstructure:"identity_policies,omitempty" json:"identity_policies,omitempty"`
	ExternalNamespacePolicies map[string][]string `mapstructure:"external_namespace_policies,omitempty" json:"external_namespace_policies,omitempty"`
	NoDefaultPolicy           bool                `mapstructure:"no_default_policy,omitempty" json:"no_default_policy,omitempty"`
	PolicyResults             *PolicyResults      `mapstructure:"policy_results,omitempty" json:"policy_results,omitempty"`
	Metadata                  map[string]string   `mapstructure:"metadata,omitempty" json:"metadata,omitempty"`
	NumUses                   int                 `mapstructure:"num_uses,omitempty" json:"num_uses,omitempty"`
	RemainingUses             int                 `mapstructure:"remaining_uses,omitempty" json:"remaining_uses,omitempty"`
	EntityID                  string              `mapstructure:"entity_id,omitempty" json:"entity_id,omitempty"`
	EntityCreated             bool                `mapstructure:"entity_created,omitempty" json:"entity_created,omitempty"`
	TokenType                 string              `mapstructure:"token_type,omitempty" json:"token_type,omitempty"`
	TokenTTL                  int64               `mapstructure:"token_ttl,omitempty" json:"token_ttl,omitempty"`
	TokenIssueTime            string              `mapstructure:"token_issue_time,omitempty" json:"token_issue_time,omitempty"`
}

type PolicyResults struct {
	Allowed          bool         `mapstructure:"allowed" json:"allowed"`
	GrantingPolicies []PolicyInfo `mapstructure:"granting_policies,omitempty" json:"granting_policies,omitempty"`
}

type PolicyInfo struct {
	Name          string `mapstructure:"name,omitempty" json:"name,omitempty"`
	NamespaceId   string `mapstructure:"namespace_id,omitempty" json:"namespace_id,omitempty"`
	NamespacePath string `mapstructure:"namespace_path,omitempty" json:"namespace_path,omitempty"`
	Type          string `mapstructure:"type" json:"type"`
}

type Secret struct {
	LeaseID string `mapstructure:"lease_id,omitempty" json:"lease_id,omitempty"`
}

type ResponseWrapInfo struct {
	TTL             int    `mapstructure:"ttl,omitempty" json:"ttl,omitempty"`
	Token           string `mapstructure:"token,omitempty" json:"token,omitempty"`
	Accessor        string `mapstructure:"accessor,omitempty" json:"accessor,omitempty"`
	CreationTime    string `mapstructure:"creation_time,omitempty" json:"creation_time,omitempty"`
	CreationPath    string `mapstructure:"creation_path,omitempty" json:"creation_path,omitempty"`
	WrappedAccessor string `mapstructure:"wrapped_accessor,omitempty" json:"wrapped_accessor,omitempty"`
}

type Namespace struct {
	ID   string `mapstructure:"id,omitempty" json:"id,omitempty"`
	Path string `mapstructure:"path,omitempty" json:"path,omitempty"`
}

// nonPersistentSalt is used for obtaining a salt that is not persisted.
type nonPersistentSalt struct{}

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
type Factory func(context.Context, *BackendConfig, HeaderFormatter) (Backend, error)
