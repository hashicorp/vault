// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"github.com/hashicorp/vault/sdk/logical"
)

// entry represents an audit entry.
// It could be an entry for a request or response.
type entry struct {
	Auth          *auth     `json:"auth,omitempty"`
	Error         string    `json:"error,omitempty"`
	Forwarded     bool      `json:"forwarded,omitempty"`
	ForwardedFrom string    `json:"forwarded_from,omitempty"` // Populated in Enterprise when a request is forwarded
	Request       *request  `json:"request,omitempty"`
	Response      *response `json:"response,omitempty"`
	Time          string    `json:"time,omitempty"`
	Type          string    `json:"type,omitempty"`
}

type request struct {
	ClientCertificateSerialNumber string                 `json:"client_certificate_serial_number,omitempty"`
	ClientID                      string                 `json:"client_id,omitempty"`
	ClientToken                   string                 `json:"client_token,omitempty"`
	ClientTokenAccessor           string                 `json:"client_token_accessor,omitempty"`
	Data                          map[string]interface{} `json:"data,omitempty"`
	Headers                       map[string][]string    `json:"headers,omitempty"`
	ID                            string                 `json:"id,omitempty"`
	MountAccessor                 string                 `json:"mount_accessor,omitempty"`
	MountClass                    string                 `json:"mount_class,omitempty"`
	MountIsExternalPlugin         bool                   `json:"mount_is_external_plugin,omitempty"`
	MountPoint                    string                 `json:"mount_point,omitempty"`
	MountRunningSha256            string                 `json:"mount_running_sha256,omitempty"`
	MountRunningVersion           string                 `json:"mount_running_version,omitempty"`
	MountType                     string                 `json:"mount_type,omitempty"`
	Namespace                     *namespace             `json:"namespace,omitempty"`
	Operation                     logical.Operation      `json:"operation,omitempty"`
	Path                          string                 `json:"path,omitempty"`
	PolicyOverride                bool                   `json:"policy_override,omitempty"`
	RemoteAddr                    string                 `json:"remote_address,omitempty"`
	RemotePort                    int                    `json:"remote_port,omitempty"`
	ReplicationCluster            string                 `json:"replication_cluster,omitempty"`
	RequestURI                    string                 `json:"request_uri,omitempty"`
	WrapTTL                       int                    `json:"wrap_ttl,omitempty"`
}

type response struct {
	Auth                  *auth                  `json:"auth,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	Headers               map[string][]string    `json:"headers,omitempty"`
	MountAccessor         string                 `json:"mount_accessor,omitempty"`
	MountClass            string                 `json:"mount_class,omitempty"`
	MountIsExternalPlugin bool                   `json:"mount_is_external_plugin,omitempty"`
	MountPoint            string                 `json:"mount_point,omitempty"`
	MountRunningSha256    string                 `json:"mount_running_sha256,omitempty"`
	MountRunningVersion   string                 `json:"mount_running_plugin_version,omitempty"`
	MountType             string                 `json:"mount_type,omitempty"`
	Redirect              string                 `json:"redirect,omitempty"`
	Secret                *secret                `json:"secret,omitempty"`
	WrapInfo              *responseWrapInfo      `json:"wrap_info,omitempty"`
	Warnings              []string               `json:"warnings,omitempty"`
}

type auth struct {
	Accessor                  string              `json:"accessor,omitempty"`
	ClientToken               string              `json:"client_token,omitempty"`
	DisplayName               string              `json:"display_name,omitempty"`
	EntityCreated             bool                `json:"entity_created,omitempty"`
	EntityID                  string              `json:"entity_id,omitempty"`
	ExternalNamespacePolicies map[string][]string `json:"external_namespace_policies,omitempty"`
	IdentityPolicies          []string            `json:"identity_policies,omitempty"`
	Metadata                  map[string]string   `json:"metadata,omitempty"`
	NoDefaultPolicy           bool                `json:"no_default_policy,omitempty"`
	NumUses                   int                 `json:"num_uses,omitempty"`
	Policies                  []string            `json:"policies,omitempty"`
	PolicyResults             *policyResults      `json:"policy_results,omitempty"`
	RemainingUses             int                 `json:"remaining_uses,omitempty"`
	TokenPolicies             []string            `json:"token_policies,omitempty"`
	TokenIssueTime            string              `json:"token_issue_time,omitempty"`
	TokenTTL                  int64               `json:"token_ttl,omitempty"`
	TokenType                 string              `json:"token_type,omitempty"`
}

type policyResults struct {
	Allowed          bool         `json:"allowed"`
	GrantingPolicies []policyInfo `json:"granting_policies,omitempty"`
}

type policyInfo struct {
	Name          string `json:"name,omitempty"`
	NamespaceId   string `json:"namespace_id,omitempty"`
	NamespacePath string `json:"namespace_path,omitempty"`
	Type          string `json:"type"`
}

type secret struct {
	LeaseID string `json:"lease_id,omitempty"`
}

type responseWrapInfo struct {
	Accessor        string `json:"accessor,omitempty"`
	CreationPath    string `json:"creation_path,omitempty"`
	CreationTime    string `json:"creation_time,omitempty"`
	Token           string `json:"token,omitempty"`
	TTL             int    `json:"ttl,omitempty"`
	WrappedAccessor string `json:"wrapped_accessor,omitempty"`
}

type namespace struct {
	ID   string `json:"id,omitempty"`
	Path string `json:"path,omitempty"`
}
