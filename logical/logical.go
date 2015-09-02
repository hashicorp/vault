package logical

import "log"

// Backend interface must be implemented to be "mountable" at
// a given path. Requests flow through a router which has various mount
// points that flow to a logical backend. The logic of each backend is flexible,
// and this is what allows materialized keys to function. There can be specialized
// logical backends for various upstreams (Consul, PostgreSQL, MySQL, etc) that can
// interact with remote APIs to generate keys dynamically. This interface also
// allows for a "procfs" like interaction, as internal state can be exposed by
// acting like a logical backend and being mounted.
type Backend interface {
	// HandleRequest is used to handle a request and generate a response.
	// The backends must check the operation type and handle appropriately.
	HandleRequest(*Request) (*Response, error)

	// SpecialPaths is a list of paths that are special in some way.
	// See PathType for the types of special paths. The key is the type
	// of the special path, and the value is a list of paths for this type.
	// This is not a regular expression but is an exact match. If the path
	// ends in '*' then it is a prefix-based match. The '*' can only appear
	// at the end.
	SpecialPaths() *Paths

	// System provides an interface to access certain system configuration
	// information, such as globally configured default and max lease TTLs.
	System() SystemView
}

// BackendConfig is provided to the factory to initialize the backend
type BackendConfig struct {
	// View should not be stored, and should only be used for initialization
	View Storage

	// The backend should use this logger. The log should not contain any secrets.
	Logger *log.Logger

	// System provides a view into a subset of safe system information that
	// is useful for backends, such as the default/max lease TTLs
	System SystemView

	// Config is the opaque user configuration provided when mounting
	Config map[string]string
}

// Factory is the factory function to create a logical backend.
type Factory func(*BackendConfig) (Backend, error)

// Paths is the structure of special paths that is used for SpecialPaths.
type Paths struct {
	// Root are the paths that require a root token to access
	Root []string

	// Unauthenticated are the paths that can be accessed without any auth.
	Unauthenticated []string
}
