package credential

import (
	"crypto/tls"
	"net"

	"github.com/hashicorp/vault/logical"
)

const (
	// PolicyKey is the key in the Secret that is read to determine the
	// associated policies of the user.
	PolicyKey = "policy"

	// MetadataKey is the prefix checked in the InternalData of a Secret
	// to attach to a token. For example "meta_user=armon" is used to add
	// the "user=armon" metadata to a token.
	MetadataKey = "meta_"
)

// Backend interface must be implemented for an authentication
// mechanism to be made available. Requests can flow through credential
// backends to be converted into a token. The logic of each backend is flexible,
// and this is allows for user/password, public/private key, and OAuth schemes
// to all be supported. The credential implementations must also be logical
// backends, allowing them to be mounted and manipulated like procfs.
type Backend interface {
	logical.Backend

	// HandleLogin is used to handle a login request and generate a response.
	// The backend is allowed to ignore this request if it is not applicable.
	HandleLogin(req *Request) (*Response, error)
}

// Factory is the factory function to create a logical backend.
type Factory func(map[string]string) (Backend, error)

// Request is used to provide access to the user parameters of
// a request. This provides more raw access than a logical.Request.
type Request struct {
	// Path is the request path
	Path string

	// Request data is an opaque map that must have string keys.
	Data map[string]interface{}

	// RemoteAddr provides the remote address if applicable
	RemoteAddr net.Addr

	// ConnState provides the TLS connection state if applicable
	ConnState *tls.ConnectionState

	// Storage can be used to durably store and retrieve state.
	Storage logical.Storage
}

// Response is used to tell the core about an authenticated
// user and provide enough information to process a request.
type Response struct {
	// Secret is returned to provide a token with a lease. This should
	// only be returned if the user is authenticated.
	Secret *logical.Secret

	// Response data is an opaque map that must have string keys. For
	// secrets, this data is sent down to the user as-is. To store internal
	// data that you don't want the user to see, store it in
	// Secret.InternalData.
	Data map[string]interface{}

	// Redirect is used to redirect to another location. This can
	// be used for flows that require going through another system.
	Redirect string
}
