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
	Time          string   `bexpr:"time,omitempty" mapstructure:"time,omitempty" json:"time,omitempty"`
	Type          string   `bexpr:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Auth          *Auth    `bexpr:"auth,omitempty" mapstructure:"auth,omitempty" json:"auth,omitempty"`
	Request       *Request `bexpr:"request,omitempty" mapstructure:"request,omitempty" json:"request,omitempty"`
	Error         string   `bexpr:"error,omitempty" mapstructure:"error,omitempty" json:"error,omitempty"`
	ForwardedFrom string   `bexpr:"forwarded_from,omitempty" mapstructure:"forwarded_from,omitempty" json:"forwarded_from,omitempty"` // Populated in Enterprise when a request is forwarded
}

// ResponseEntry is the structure of a response audit log entry.
type ResponseEntry struct {
	Time      string    `bexpr:"time,omitempty" mapstructure:"time,omitempty" json:"time,omitempty"`
	Type      string    `bexpr:"type,omitempty" mapstructure:"type,omitempty" json:"type,omitempty"`
	Auth      *Auth     `bexpr:"auth,omitempty" mapstructure:"auth,omitempty" json:"auth,omitempty"`
	Request   *Request  `bexpr:"request,omitempty" mapstructure:"request,omitempty" json:"request,omitempty"`
	Response  *Response `bexpr:"response,omitempty" mapstructure:"response,omitempty" json:"response,omitempty"`
	Error     string    `bexpr:"error,omitempty" mapstructure:"error,omitempty" json:"error,omitempty"`
	Forwarded bool      `bexpr:"forwarded,omitempty" mapstructure:"forwarded,omitempty" json:"forwarded,omitempty"`
}

type Request struct {
	ID                            string                 `bexpr:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	ClientID                      string                 `bexpr:"client_id,omitempty" mapstructure:"client_id,omitempty" json:"client_id,omitempty"`
	ReplicationCluster            string                 `bexpr:"replication_cluster,omitempty" mapstructure:"replication_cluster,omitempty" json:"replication_cluster,omitempty"`
	Operation                     logical.Operation      `bexpr:"operation,omitempty" mapstructure:"operation,omitempty" json:"operation,omitempty"`
	MountPoint                    string                 `bexpr:"mount_point,omitempty" mapstructure:"mount_point,omitempty" json:"mount_point,omitempty"`
	MountType                     string                 `bexpr:"mount_type,omitempty" mapstructure:"mount_type,omitempty" json:"mount_type,omitempty"`
	MountAccessor                 string                 `bexpr:"mount_accessor,omitempty" mapstructure:"mount_accessor,omitempty" json:"mount_accessor,omitempty"`
	MountRunningVersion           string                 `bexpr:"mount_running_version,omitempty" mapstructure:"mount_running_version,omitempty" json:"mount_running_version,omitempty"`
	MountRunningSha256            string                 `bexpr:"mount_running_sha256,omitempty" mapstructure:"mount_running_sha256,omitempty" json:"mount_running_sha256,omitempty"`
	MountClass                    string                 `bexpr:"mount_class,omitempty" mapstructure:"mount_class,omitempty" json:"mount_class,omitempty"`
	MountIsExternalPlugin         bool                   `bexpr:"mount_is_external_plugin,omitempty" mapstructure:"mount_is_external_plugin,omitempty" json:"mount_is_external_plugin,omitempty"`
	ClientToken                   string                 `bexpr:"client_token,omitempty" mapstructure:"client_token,omitempty" json:"client_token,omitempty"`
	ClientTokenAccessor           string                 `bexpr:"client_token_accessor,omitempty" mapstructure:"client_token_accessor,omitempty" json:"client_token_accessor,omitempty"`
	Namespace                     *Namespace             `bexpr:"namespace,omitempty" mapstructure:"namespace,omitempty"  json:"namespace,omitempty"`
	Path                          string                 `bexpr:"path,omitempty" mapstructure:"path,omitempty" json:"path,omitempty"`
	Data                          map[string]interface{} `bexpr:"data,omitempty" mapstructure:"data,omitempty" json:"data,omitempty"`
	PolicyOverride                bool                   `bexpr:"policy_override,omitempty" mapstructure:"policy_override,omitempty" json:"policy_override,omitempty"`
	RemoteAddr                    string                 `bexpr:"remote_address,omitempty" mapstructure:"remote_address,omitempty" json:"remote_address,omitempty"`
	RemotePort                    int                    `bexpr:"remote_port,omitempty" mapstructure:"remote_port,omitempty" json:"remote_port,omitempty"`
	WrapTTL                       int                    `bexpr:"wrap_ttl,omitempty" mapstructure:"wrap_ttl,omitempty" json:"wrap_ttl,omitempty"`
	Headers                       map[string][]string    `bexpr:"headers,omitempty" mapstructure:"headers,omitempty" json:"headers,omitempty"`
	ClientCertificateSerialNumber string                 `bexpr:"client_certificate_serial_number,omitempty" mapstructure:"client_certificate_serial_number,omitempty" json:"client_certificate_serial_number,omitempty"`
	RequestURI                    string                 `bexpr:"request_uri,omitempty" mapstructure:"request_uri,omitempty" json:"request_uri,omitempty"`
}

type Response struct {
	Auth                  *Auth                  `bexpr:"auth,omitempty" mapstructure:"auth,omitempty" json:"auth,omitempty"`
	MountPoint            string                 `bexpr:"mount_point,omitempty" mapstructure:"mount_point,omitempty" json:"mount_point,omitempty"`
	MountType             string                 `bexpr:"mount_type,omitempty" mapstructure:"mount_type,omitempty" json:"mount_type,omitempty"`
	MountAccessor         string                 `bexpr:"mount_accessor,omitempty" mapstructure:"mount_accessor,omitempty" json:"mount_accessor,omitempty"`
	MountRunningVersion   string                 `bexpr:"mount_running_plugin_version,omitempty" mapstructure:"mount_running_plugin_version,omitempty" json:"mount_running_plugin_version,omitempty"`
	MountRunningSha256    string                 `bexpr:"mount_running_sha256,omitempty" mapstructure:"mount_running_sha256,omitempty" json:"mount_running_sha256,omitempty"`
	MountClass            string                 `bexpr:"mount_class,omitempty" mapstructure:"mount_class,omitempty" json:"mount_class,omitempty"`
	MountIsExternalPlugin bool                   `bexpr:"mount_is_external_plugin,omitempty" mapstructure:"mount_is_external_plugin,omitempty" json:"mount_is_external_plugin,omitempty"`
	Secret                *Secret                `bexpr:"secret,omitempty" mapstructure:"secret,omitempty" json:"secret,omitempty"`
	Data                  map[string]interface{} `bexpr:"data,omitempty" mapstructure:"data,omitempty" json:"data,omitempty"`
	Warnings              []string               `bexpr:"warnings,omitempty" mapstructure:"warnings,omitempty" json:"warnings,omitempty"`
	Redirect              string                 `bexpr:"redirect,omitempty" mapstructure:"redirect,omitempty" json:"redirect,omitempty"`
	WrapInfo              *ResponseWrapInfo      `bexpr:"wrap_info,omitempty" mapstructure:"wrap_info,omitempty" json:"wrap_info,omitempty"`
	Headers               map[string][]string    `bexpr:"headers,omitempty" mapstructure:"headers,omitempty" json:"headers,omitempty"`
}

type Auth struct {
	ClientToken               string              `bexpr:"client_token,omitempty" mapstructure:"client_token,omitempty" json:"client_token,omitempty"`
	Accessor                  string              `bexpr:"accessor,omitempty" mapstructure:"accessor,omitempty" json:"accessor,omitempty"`
	DisplayName               string              `bexpr:"display_name,omitempty" mapstructure:"display_name,omitempty" json:"display_name,omitempty"`
	Policies                  []string            `bexpr:"policies,omitempty" mapstructure:"policies,omitempty" json:"policies,omitempty"`
	TokenPolicies             []string            `bexpr:"token_policies,omitempty" mapstructure:"token_policies,omitempty" json:"token_policies,omitempty"`
	IdentityPolicies          []string            `bexpr:"identity_policies,omitempty" mapstructure:"identity_policies,omitempty" json:"identity_policies,omitempty"`
	ExternalNamespacePolicies map[string][]string `bexpr:"external_namespace_policies,omitempty" mapstructure:"external_namespace_policies,omitempty" json:"external_namespace_policies,omitempty"`
	NoDefaultPolicy           bool                `bexpr:"no_default_policy,omitempty" mapstructure:"no_default_policy,omitempty" json:"no_default_policy,omitempty"`
	PolicyResults             *PolicyResults      `bexpr:"policy_results,omitempty" mapstructure:"policy_results,omitempty" json:"policy_results,omitempty"`
	Metadata                  map[string]string   `bexpr:"metadata,omitempty" mapstructure:"metadata,omitempty" json:"metadata,omitempty"`
	NumUses                   int                 `bexpr:"num_uses,omitempty" mapstructure:"num_uses,omitempty" json:"num_uses,omitempty"`
	RemainingUses             int                 `bexpr:"remaining_uses,omitempty" mapstructure:"remaining_uses,omitempty" json:"remaining_uses,omitempty"`
	EntityID                  string              `bexpr:"entity_id,omitempty" mapstructure:"entity_id,omitempty" json:"entity_id,omitempty"`
	EntityCreated             bool                `bexpr:"entity_created,omitempty" mapstructure:"entity_created,omitempty" json:"entity_created,omitempty"`
	TokenType                 string              `bexpr:"token_type,omitempty" mapstructure:"token_type,omitempty" json:"token_type,omitempty"`
	TokenTTL                  int64               `bexpr:"token_ttl,omitempty" mapstructure:"token_ttl,omitempty" json:"token_ttl,omitempty"`
	TokenIssueTime            string              `bexpr:"token_issue_time,omitempty" mapstructure:"token_issue_time,omitempty" json:"token_issue_time,omitempty"`
}

type PolicyResults struct {
	Allowed          bool         `bexpr:"allowed" mapstructure:"allowed" json:"allowed"`
	GrantingPolicies []PolicyInfo `bexpr:"granting_policies,omitempty" mapstructure:"granting_policies,omitempty" json:"granting_policies,omitempty"`
}

type PolicyInfo struct {
	Name          string `bexpr:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	NamespaceId   string `bexpr:"namespace_id,omitempty" mapstructure:"namespace_id,omitempty" json:"namespace_id,omitempty"`
	NamespacePath string `bexpr:"namespace_path,omitempty" mapstructure:"namespace_path,omitempty" json:"namespace_path,omitempty"`
	Type          string `bexpr:"type" mapstructure:"type" json:"type"`
}

type Secret struct {
	LeaseID string `bexpr:"lease_id,omitempty" mapstructure:"lease_id,omitempty" json:"lease_id,omitempty"`
}

type ResponseWrapInfo struct {
	TTL             int    `bexpr:"ttl,omitempty" mapstructure:"ttl,omitempty" json:"ttl,omitempty"`
	Token           string `bexpr:"token,omitempty" mapstructure:"token,omitempty" json:"token,omitempty"`
	Accessor        string `bexpr:"accessor,omitempty" mapstructure:"accessor,omitempty" json:"accessor,omitempty"`
	CreationTime    string `bexpr:"creation_time,omitempty" mapstructure:"creation_time,omitempty" json:"creation_time,omitempty"`
	CreationPath    string `bexpr:"creation_path,omitempty" mapstructure:"creation_path,omitempty" json:"creation_path,omitempty"`
	WrappedAccessor string `bexpr:"wrapped_accessor,omitempty" mapstructure:"wrapped_accessor,omitempty" json:"wrapped_accessor,omitempty"`
}

type Namespace struct {
	ID   string `bexpr:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Path string `bexpr:"path,omitempty" mapstructure:"path,omitempty" json:"path,omitempty"`
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
