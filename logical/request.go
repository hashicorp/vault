package logical

import (
	"errors"
	"fmt"
	"time"
)

// Request is a struct that stores the parameters and context
// of a request being made to Vault. It is used to abstract
// the details of the higher level request protocol from the handlers.
type Request struct {
	// Id is the uuid associated with each request
	ID string `json:"id" structs:"id" mapstructure:"id"`

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
	Storage Storage `json:"storage" structs:"storage" mapstructure:"storage"`

	// Secret will be non-nil only for Revoke and Renew operations
	// to represent the secret that was returned prior.
	Secret *Secret `json:"secret" structs:"secret" mapstructure:"secret"`

	// Auth will be non-nil only for Renew operations
	// to represent the auth that was returned prior.
	Auth *Auth `json:"auth" structs:"auth" mapstructure:"auth"`

	// Connection will be non-nil only for credential providers to
	// inspect the connection information and potentially use it for
	// authentication/protection.
	Connection *Connection `json:"connection" structs:"connection" mapstructure:"connection"`

	// ClientToken is provided to the core so that the identity
	// can be verified and ACLs applied. This value is passed
	// through to the logical backends but after being salted and
	// hashed.
	ClientToken string `json:"client_token" structs:"client_token" mapstructure:"client_token"`

	// DisplayName is provided to the logical backend to help associate
	// dynamic secrets with the source entity. This is not a sensitive
	// name, but is useful for operators.
	DisplayName string `json:"display_name" structs:"display_name" mapstructure:"display_name"`

	// MountPoint is provided so that a logical backend can generate
	// paths relative to itself. The `Path` is effectively the client
	// request path with the MountPoint trimmed off.
	MountPoint string `json:"mount_point" structs:"mount_point" mapstructure:"mount_point"`

	// WrapTTL contains the requested TTL of the token used to wrap the
	// response in a cubbyhole.
	WrapTTL time.Duration `json:"wrap_ttl" struct:"wrap_ttl" mapstructure:"wrap_ttl"`
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
	CreateOperation Operation = "create"
	ReadOperation             = "read"
	UpdateOperation           = "update"
	DeleteOperation           = "delete"
	ListOperation             = "list"
	HelpOperation             = "help"

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
