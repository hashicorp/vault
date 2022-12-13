package logical

import (
	"context"

	log "github.com/hashicorp/go-hclog"
)

// BackendType is the type of backend that is being implemented
type BackendType uint32

// The these are the types of backends that can be derived from
// logical.Backend
const (
	TypeUnknown    BackendType = 0 // This is also the zero-value for BackendType
	TypeLogical    BackendType = 1
	TypeCredential BackendType = 2
)

// Stringer implementation
func (b BackendType) String() string {
	switch b {
	case TypeLogical:
		return "secret"
	case TypeCredential:
		return "auth"
	}

	return "unknown"
}

// Backend interface must be implemented to be "mountable" at
// a given path. Requests flow through a router which has various mount
// points that flow to a logical backend. The logic of each backend is flexible,
// and this is what allows materialized keys to function. There can be specialized
// logical backends for various upstreams (Consul, PostgreSQL, MySQL, etc) that can
// interact with remote APIs to generate keys dynamically. This interface also
// allows for a "procfs" like interaction, as internal state can be exposed by
// acting like a logical backend and being mounted.
type Backend interface {
	// Initialize is used to initialize a plugin after it has been mounted.
	Initialize(context.Context, *InitializationRequest) error

	// HandleRequest is used to handle a request and generate a response.
	// The backends must check the operation type and handle appropriately.
	HandleRequest(context.Context, *Request) (*Response, error)

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

	// Logger provides an interface to access the underlying logger. This
	// is useful when a struct embeds a Backend-implemented struct that
	// contains a private instance of logger.
	Logger() log.Logger

	// HandleExistenceCheck is used to handle a request and generate a response
	// indicating whether the given path exists or not; this is used to
	// understand whether the request must have a Create or Update capability
	// ACL applied. The first bool indicates whether an existence check
	// function was found for the backend; the second indicates whether, if an
	// existence check function was found, the item exists or not.
	HandleExistenceCheck(context.Context, *Request) (bool, bool, error)

	// Cleanup is invoked during an unmount of a backend to allow it to
	// handle any cleanup like connection closing or releasing of file handles.
	Cleanup(context.Context)

	// InvalidateKey may be invoked when an object is modified that belongs
	// to the backend. The backend can use this to clear any caches or reset
	// internal state as needed.
	InvalidateKey(context.Context, string)

	// Setup is used to set up the backend based on the provided backend
	// configuration.
	Setup(context.Context, *BackendConfig) error

	// Type returns the BackendType for the particular backend
	Type() BackendType
}

// BackendConfig is provided to the factory to initialize the backend
type BackendConfig struct {
	// View should not be stored, and should only be used for initialization
	StorageView Storage

	// The backend should use this logger. The log should not contain any secrets.
	Logger log.Logger

	// System provides a view into a subset of safe system information that
	// is useful for backends, such as the default/max lease TTLs
	System SystemView

	// BackendUUID is a unique identifier provided to this backend. It's useful
	// when a backend needs a consistent and unique string without using storage.
	BackendUUID string

	// Config is the opaque user configuration provided when mounting
	Config map[string]string
}

// Factory is the factory function to create a logical backend.
type Factory func(context.Context, *BackendConfig) (Backend, error)

// Paths is the structure of special paths that is used for SpecialPaths.
type Paths struct {
	// Root are the API paths that require a root token to access
	Root []string

	// Unauthenticated are the API paths that can be accessed without any auth.
	// These can't be regular expressions, it is either exact match, a prefix
	// match and/or a wildcard match. For prefix match, append '*' as a suffix.
	// For a wildcard match, use '+' in the segment to match any identifier
	// (e.g. 'foo/+/bar'). Note that '+' can't be adjacent to a non-slash.
	Unauthenticated []string

	// LocalStorage are storage paths (prefixes) that are local to this cluster;
	// this indicates that these paths should not be replicated across performance clusters
	// (DR replication is unaffected).
	LocalStorage []string

	// SealWrapStorage are storage paths that, when using a capable seal,
	// should be seal wrapped with extra encryption. It is exact matching
	// unless it ends with '/' in which case it will be treated as a prefix.
	SealWrapStorage []string
}

type Auditor interface {
	AuditRequest(ctx context.Context, input *LogInput) error
	AuditResponse(ctx context.Context, input *LogInput) error
}

// Externaler allows us to check if a backend is running externally (i.e., over GRPC)
type Externaler interface {
	IsExternal() bool
}

type PluginVersion struct {
	Version string
}

// PluginVersioner is an optional interface to return version info.
type PluginVersioner interface {
	// PluginVersion returns the version for the backend
	PluginVersion() PluginVersion
}

var EmptyPluginVersion = PluginVersion{""}
