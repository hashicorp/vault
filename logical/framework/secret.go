package framework

import (
	"time"

	"github.com/hashicorp/vault/logical"
)

// Secret is a type of secret that can be returned from a backend.
type Secret struct {
	// Type is the name of this secret type. This is used to setup the
	// vault ID and to look up the proper secret structure when revocation/
	// renewal happens. Once this is set this should not be changed.
	//
	// The format of this must match (case insensitive): ^a-Z0-9_$
	Type string

	// Fields is the mapping of data fields and schema that comprise
	// the structure of this secret.
	Fields map[string]*FieldSchema

	// DefaultDuration and DefaultGracePeriod are the default values for
	// the duration of the lease for this secret and its grace period. These
	// can be manually overwritten with the result of Response().
	DefaultDuration    time.Duration
	DefaultGracePeriod time.Duration

	// Below are the operations that can be called on the secret.
	//
	// Renew, if not set, will mark the secret as not renewable.
	//
	// Revoke is required.
	Renew  OperationFunc
	Revoke OperationFunc
}

func (s *Secret) Response(data map[string]interface{}) *logical.Response {
	internalData := map[string]interface{}{
		"secret_type": s.Type,
	}

	return &logical.Response{
		Secret: &logical.Secret{
			Renewable:        s.Renew != nil,
			Lease:            s.DefaultDuration,
			LeaseGracePeriod: s.DefaultGracePeriod,
			InternalData:     internalData,
		},

		Data: data,
	}
}
