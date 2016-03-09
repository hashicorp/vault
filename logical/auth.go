package logical

import "fmt"

// Auth is the resulting authentication information that is part of
// Response for credential backends.
type Auth struct {
	LeaseOptions

	// InternalData is JSON-encodable data that is stored with the auth struct.
	// This will be sent back during a Renew/Revoke for storing internal data
	// used for those operations.
	InternalData map[string]interface{}

	// DisplayName is a non-security sensitive identifier that is
	// applicable to this Auth. It is used for logging and prefixing
	// of dynamic secrets. For example, DisplayName may be "armon" for
	// the github credential backend. If the client token is used to
	// generate a SQL credential, the user may be "github-armon-uuid".
	// This is to help identify the source without using audit tables.
	DisplayName string

	// Policies is the list of policies that the authenticated user
	// is associated with.
	Policies []string

	// Metadata is used to attach arbitrary string-type metadata to
	// an authenticated user. This metadata will be outputted into the
	// audit log.
	Metadata map[string]string

	// ClientToken is the token that is generated for the authentication.
	// This will be filled in by Vault core when an auth structure is
	// returned. Setting this manually will have no effect.
	ClientToken string

	// Accessor is the identifier for the ClientToken. This can be used
	// to perform management functionalities (especially revocation) when
	// ClientToken in the audit logs are obfuscated. Accessor can be used
	// to revoke a ClientToken and to lookup the capabilities of the ClientToken,
	// both without actually knowing the ClientToken.
	Accessor string
}

func (a *Auth) GoString() string {
	return fmt.Sprintf("*%#v", *a)
}
