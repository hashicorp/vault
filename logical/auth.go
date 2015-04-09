package logical

import "fmt"

// Auth is the resulting authentication information that is part of
// Response for credential backends.
type Auth struct {
	LeaseOptions

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
}

func (a *Auth) GoString() string {
	return fmt.Sprintf("*%#v", *a)
}
