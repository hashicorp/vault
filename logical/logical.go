package logical

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

	// RootPaths is a list of paths that require root level privileges.
	// These paths will be enforced by the router so that backends do
	// not need to handle the authorization. Paths are enforced exactly
	// or using a prefix match if they end in '*'
	RootPaths() []string
}

// Factory is the factory function to create a logical backend.
type Factory func(map[string]string) (Backend, error)
