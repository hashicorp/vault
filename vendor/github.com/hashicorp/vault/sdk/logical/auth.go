package logical

import (
	"fmt"
	"time"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

// Auth is the resulting authentication information that is part of
// Response for credential backends.
type Auth struct {
	LeaseOptions

	// InternalData is JSON-encodable data that is stored with the auth struct.
	// This will be sent back during a Renew/Revoke for storing internal data
	// used for those operations.
	InternalData map[string]interface{} `json:"internal_data" mapstructure:"internal_data" structs:"internal_data"`

	// DisplayName is a non-security sensitive identifier that is
	// applicable to this Auth. It is used for logging and prefixing
	// of dynamic secrets. For example, DisplayName may be "armon" for
	// the github credential backend. If the client token is used to
	// generate a SQL credential, the user may be "github-armon-uuid".
	// This is to help identify the source without using audit tables.
	DisplayName string `json:"display_name" mapstructure:"display_name" structs:"display_name"`

	// Policies is the list of policies that the authenticated user
	// is associated with.
	Policies []string `json:"policies" mapstructure:"policies" structs:"policies"`

	// TokenPolicies and IdentityPolicies break down the list in Policies to
	// help determine where a policy was sourced
	TokenPolicies    []string `json:"token_policies" mapstructure:"token_policies" structs:"token_policies"`
	IdentityPolicies []string `json:"identity_policies" mapstructure:"identity_policies" structs:"identity_policies"`

	// ExternalNamespacePolicies represent the policies authorized from
	// different namespaces indexed by respective namespace identifiers
	ExternalNamespacePolicies map[string][]string `json:"external_namespace_policies" mapstructure:"external_namespace_policies" structs:"external_namespace_policies"`

	// Indicates that the default policy should not be added by core when
	// creating a token. The default policy will still be added if it's
	// explicitly defined.
	NoDefaultPolicy bool `json:"no_default_policy" mapstructure:"no_default_policy" structs:"no_default_policy"`

	// Metadata is used to attach arbitrary string-type metadata to
	// an authenticated user. This metadata will be outputted into the
	// audit log.
	Metadata map[string]string `json:"metadata" mapstructure:"metadata" structs:"metadata"`

	// ClientToken is the token that is generated for the authentication.
	// This will be filled in by Vault core when an auth structure is
	// returned. Setting this manually will have no effect.
	ClientToken string `json:"client_token" mapstructure:"client_token" structs:"client_token"`

	// Accessor is the identifier for the ClientToken. This can be used
	// to perform management functionalities (especially revocation) when
	// ClientToken in the audit logs are obfuscated. Accessor can be used
	// to revoke a ClientToken and to lookup the capabilities of the ClientToken,
	// both without actually knowing the ClientToken.
	Accessor string `json:"accessor" mapstructure:"accessor" structs:"accessor"`

	// Period indicates that the token generated using this Auth object
	// should never expire. The token should be renewed within the duration
	// specified by this period.
	Period time.Duration `json:"period" mapstructure:"period" structs:"period"`

	// ExplicitMaxTTL is the max TTL that constrains periodic tokens. For normal
	// tokens, this value is constrained by the configured max ttl.
	ExplicitMaxTTL time.Duration `json:"explicit_max_ttl" mapstructure:"explicit_max_ttl" structs:"explicit_max_ttl"`

	// Number of allowed uses of the issued token
	NumUses int `json:"num_uses" mapstructure:"num_uses" structs:"num_uses"`

	// EntityID is the identifier of the entity in identity store to which the
	// identity of the authenticating client belongs to.
	EntityID string `json:"entity_id" mapstructure:"entity_id" structs:"entity_id"`

	// Alias is the information about the authenticated client returned by
	// the auth backend
	Alias *Alias `json:"alias" mapstructure:"alias" structs:"alias"`

	// GroupAliases are the informational mappings of external groups which an
	// authenticated user belongs to. This is used to check if there are
	// mappings groups for the group aliases in identity store. For all the
	// matching groups, the entity ID of the user will be added.
	GroupAliases []*Alias `json:"group_aliases" mapstructure:"group_aliases" structs:"group_aliases"`

	// The set of CIDRs that this token can be used with
	BoundCIDRs []*sockaddr.SockAddrMarshaler `json:"bound_cidrs"`

	// CreationPath is a path that the backend can return to use in the lease.
	// This is currently only supported for the token store where roles may
	// change the perceived path of the lease, even though they don't change
	// the request path itself.
	CreationPath string `json:"creation_path"`

	// TokenType is the type of token being requested
	TokenType TokenType `json:"token_type"`

	// Orphan is set if the token does not have a parent
	Orphan bool `json:"orphan"`
}

func (a *Auth) GoString() string {
	return fmt.Sprintf("*%#v", *a)
}
