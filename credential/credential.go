package credential

import "github.com/hashicorp/vault/logical"

// Backend interface must be implemented for an authentication
// mechanism to be made available. Requests can flow through credential
// backends to be converted into a token. The logic of each backend is flexible,
// and this is allows for user/password, public/private key, and OAuth schemes
// to all be supported. The credential implementations must also be logical
// backends, allowing them to be mounted and manipulated like procfs.
type Backend interface {
	logical.Backend
}

// Factory is the factory function to create a logical backend.
type Factory func(map[string]string) (Backend, error)
