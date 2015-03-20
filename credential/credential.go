package credential

import (
	"crypto/tls"
	"net"

	"github.com/hashicorp/vault/logical"
)

// Backend interface must be implemented for an authentication
// mechanism to be made available. Requests can flow through credential
// backends to be converted into a token. The logic of each backend is flexible,
// and this is allows for user/password, public/private key, and OAuth schemes
// to all be supported. The credential implementations must also be logical
// backends, allowing them to be mounted and manipulated like procfs.
type Backend interface {
	logical.Backend

	// LoginPaths is a list of paths that are unauthenticated and used
	// only for logging in. These paths cannot be reached via HandleRequest,
	// and are sent to HandleLogin instead. Paths are enforced exactly
	// or using a prefix match if they end in '*'
	LoginPaths() []string

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
	// Authenticated is used to indicate if the request has been
	// authenticated. A token will be created with the associated
	// policies.
	Authenticated bool

	// Policies is the named policies that should be applied
	Policies []string

	// Redirect is used to redirect to another location. This can
	// be used for flows that require going through another system.
	Redirect string
}
