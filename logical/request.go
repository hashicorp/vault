package logical

import (
	"errors"
	"fmt"
	"time"
)

// RequestWrapInfo is a struct that stores information about desired response
// wrapping behavior
type RequestWrapInfo struct {
	// Setting to non-zero specifies that the response should be wrapped.
	// Specifies the desired TTL of the wrapping token.
	TTL time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`

	// The format to use for the wrapped response; if not specified it's a bare
	// token
	Format string `json:"format" structs:"format" mapstructure:"format"`
}

// Request is a struct that stores the parameters and context
// of a request being made to Vault. It is used to abstract
// the details of the higher level request protocol from the handlers.
type Request struct {
	// Id is the uuid associated with each request
	ID string `json:"id" structs:"id" mapstructure:"id"`

	// If set, the name given to the replication secondary where this request
	// originated
	ReplicationCluster string `json:"replication_cluster" structs:"replication_cluster", mapstructure:"replication_cluster"`

	// Operation is the requested operation type
	Operation Operation `json:"operation" structs:"operation" mapstructure:"operation"`

	// Path is the part of the request path not consumed by the
	// routing. As an example, if the original request path is "prod/aws/foo"
	// and the AWS logical backend is mounted at "prod/aws/", then the
	// final path is "foo" since the mount prefix is trimmed.
	Path string `json:"path" structs:"path" mapstructure:"path"`

	// Request data is an opaque map that must have string keys.
	Data map[string]interface{} `json:"map" structs:"data" mapstructure:"data"`

	// Storage can be used to durably store and retrieve state.
	Storage Storage `json:"-"`

	// Secret will be non-nil only for Revoke and Renew operations
	// to represent the secret that was returned prior.
	Secret *Secret `json:"secret" structs:"secret" mapstructure:"secret"`

	// Auth will be non-nil only for Renew operations
	// to represent the auth that was returned prior.
	Auth *Auth `json:"auth" structs:"auth" mapstructure:"auth"`

	// Headers will contain the http headers from the request. This value will
	// be used in the audit broker to ensure we are auditing only the allowed
	// headers.
	Headers map[string][]string `json:"headers" structs:"headers" mapstructure:"headers"`

	// Connection will be non-nil only for credential providers to
	// inspect the connection information and potentially use it for
	// authentication/protection.
	Connection *Connection `json:"connection" structs:"connection" mapstructure:"connection"`

	// ClientToken is provided to the core so that the identity
	// can be verified and ACLs applied. This value is passed
	// through to the logical backends but after being salted and
	// hashed.
	ClientToken string `json:"client_token" structs:"client_token" mapstructure:"client_token"`

	// ClientTokenAccessor is provided to the core so that the it can get
	// logged as part of request audit logging.
	ClientTokenAccessor string `json:"client_token_accessor" structs:"client_token_accessor" mapstructure:"client_token_accessor"`

	// DisplayName is provided to the logical backend to help associate
	// dynamic secrets with the source entity. This is not a sensitive
	// name, but is useful for operators.
	DisplayName string `json:"display_name" structs:"display_name" mapstructure:"display_name"`

	// MountPoint is provided so that a logical backend can generate
	// paths relative to itself. The `Path` is effectively the client
	// request path with the MountPoint trimmed off.
	MountPoint string `json:"mount_point" structs:"mount_point" mapstructure:"mount_point"`

	// MountType is provided so that a logical backend can make decisions
	// based on the specific mount type (e.g., if a mount type has different
	// aliases, generating different defaults depending on the alias)
	MountType string `json:"mount_type" structs:"mount_type" mapstructure:"mount_type"`

	// WrapInfo contains requested response wrapping parameters
	WrapInfo *RequestWrapInfo `json:"wrap_info" structs:"wrap_info" mapstructure:"wrap_info"`

	// ClientTokenRemainingUses represents the allowed number of uses left on the
	// token supplied
	ClientTokenRemainingUses int `json:"client_token_remaining_uses" structs:"client_token_remaining_uses" mapstructure:"client_token_remaining_uses"`

	// For replication, contains the last WAL on the remote side after handling
	// the request, used for best-effort avoidance of stale read-after-write
	lastRemoteWAL uint64
}

// Get returns a data field and guards for nil Data
func (r *Request) Get(key string) interface{} {
	if r.Data == nil {
		return nil
	}
	return r.Data[key]
}

// GetString returns a data field as a string
func (r *Request) GetString(key string) string {
	raw := r.Get(key)
	s, _ := raw.(string)
	return s
}

func (r *Request) GoString() string {
	return fmt.Sprintf("*%#v", *r)
}

func (r *Request) LastRemoteWAL() uint64 {
	return r.lastRemoteWAL
}

func (r *Request) SetLastRemoteWAL(last uint64) {
	r.lastRemoteWAL = last
}

// RenewRequest creates the structure of the renew request.
func RenewRequest(
	path string, secret *Secret, data map[string]interface{}) *Request {
	return &Request{
		Operation: RenewOperation,
		Path:      path,
		Data:      data,
		Secret:    secret,
	}
}

// RenewAuthRequest creates the structure of the renew request for an auth.
func RenewAuthRequest(
	path string, auth *Auth, data map[string]interface{}) *Request {
	return &Request{
		Operation: RenewOperation,
		Path:      path,
		Data:      data,
		Auth:      auth,
	}
}

// RevokeRequest creates the structure of the revoke request.
func RevokeRequest(
	path string, secret *Secret, data map[string]interface{}) *Request {
	return &Request{
		Operation: RevokeOperation,
		Path:      path,
		Data:      data,
		Secret:    secret,
	}
}

// RollbackRequest creates the structure of the revoke request.
func RollbackRequest(path string) *Request {
	return &Request{
		Operation: RollbackOperation,
		Path:      path,
		Data:      make(map[string]interface{}),
	}
}

// Operation is an enum that is used to specify the type
// of request being made
type Operation string

const (
	// The operations below are called per path
	CreateOperation         Operation = "create"
	ReadOperation                     = "read"
	UpdateOperation                   = "update"
	DeleteOperation                   = "delete"
	ListOperation                     = "list"
	HelpOperation                     = "help"
	AliasLookaheadOperation           = "alias-lookahead"

	// The operations below are called globally, the path is less relevant.
	RevokeOperation   Operation = "revoke"
	RenewOperation              = "renew"
	RollbackOperation           = "rollback"
)

var (
	// ErrUnsupportedOperation is returned if the operation is not supported
	// by the logical backend.
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrUnsupportedPath is returned if the path is not supported
	// by the logical backend.
	ErrUnsupportedPath = errors.New("unsupported path")

	// ErrInvalidRequest is returned if the request is invalid
	ErrInvalidRequest = errors.New("invalid request")

	// ErrPermissionDenied is returned if the client is not authorized
	ErrPermissionDenied = errors.New("permission denied")
)
