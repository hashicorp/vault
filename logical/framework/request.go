package framework

import (
	"github.com/hashicorp/vault/logical"
)

// Request is a single request for a backend that wraps a logical.Request
// to provide some extra functionality.
type Request struct {
	Backend        *Backend
	Data           *FieldData
	LogicalRequest *logical.Request
}
