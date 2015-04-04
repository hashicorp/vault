package logical

import (
	"log"
)

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

	// SetLogger is called to set the logger for the backend. The backend
	// should use this logger. The log should not contain any secrets.
	// It should not be assumed that this function will be called every time.
	//
	// SetLogger will not be called by Vault core in parallel, and
	// therefore doesn't need any lock protection.
	SetLogger(*log.Logger)
}

// Factory is the factory function to create a logical backend.
type Factory func(map[string]string) (Backend, error)

// Paths is the structure of special paths that is used for SpecialPaths.
type Paths struct {
	// Root are the paths that require a root token to access
	Root []string

	// Unauthenticated are the paths that can be accessed without any auth.
	Unauthenticated []string
}
