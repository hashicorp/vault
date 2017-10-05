package logical

import (
	"fmt"
	"time"
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

	// Number of allowed uses of the issued token
	NumUses int `json:"num_uses" mapstructure:"num_uses" structs:"num_uses"`

	// Alias is the information about the authenticated client returned by
	// the auth backend
	Alias *Alias `json:"alias" structs:"alias" mapstructure:"alias"`
}

func (a *Auth) GoString() string {
	return fmt.Sprintf("*%#v", *a)
}
