package framework

import (
	"github.com/hashicorp/vault/logical"
)

// Request is a single request for a backend that wraps a logical.Request
// to provide some extra functionality.
type Request struct {
	// Backend is the backend that generated this request.
	Backend *Backend

	// Data is any parameters that were passed into the request according
	// to the path schema. If this request is to a secret operation
	// (revoke, renew), then Data is according to the schema of the
	// secret.
	Data *FieldData

	// The fields below are only set for secret-related requests (renew,
	// revoke).
	//
	// SecretType is the string type of the secret.
	//
	// SecretId is the ID of the secret that was given when generating
	// the secret, or is otherwise just the UUID that was generated.
	SecretType string
	SecretId   string

	// LogicalRequest is the raw logical.Request structure for this request.
	LogicalRequest *logical.Request
}
