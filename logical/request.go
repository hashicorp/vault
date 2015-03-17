package logical

import (
	"errors"
)

// Request is a struct that stores the parameters and context
// of a request being made to Vault. It is used to abstract
// the details of the higher level request protocol from the handlers.
type Request struct {
	// Operation is the requested operation type
	Operation Operation

	// Path is the part of the request path not consumed by the
	// routing. As an example, if the original request path is "prod/aws/foo"
	// and the AWS logical backend is mounted at "prod/aws/", then the
	// final path is "foo" since the mount prefix is trimmed.
	Path string

	// Request data is an opaque map that must have string keys.
	Data map[string]interface{}

	// Storage can be used to durably store and retrieve state.
	Storage Storage
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

// Operation is an enum that is used to specify the type
// of request being made
type Operation string

const (
	ReadOperation     Operation = "read"
	WriteOperation              = "write"
	DeleteOperation             = "delete"
	ListOperation               = "list"
	RevokeOperation             = "revoke"
	RenewOperation              = "renew"
	RollbackOperation           = "rollback"
	HelpOperation               = "help"
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
)
