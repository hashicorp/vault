package framework

import (
	"context"
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

	// DefaultDuration is the default value for the duration of the lease for
	// this secret. This can be manually overwritten with the result of
	// Response().
	//
	// If these aren't set, Vault core will set a default lease period which
	// may come from a mount tuning.
	DefaultDuration time.Duration

	// Renew is the callback called to renew this secret. If Renew is
	// not specified then renewable is set to false in the secret.
	// See lease.go for helpers for this value.
	Renew OperationFunc

	// Revoke is the callback called to revoke this secret. This is required.
	Revoke OperationFunc
}

func (s *Secret) Renewable() bool {
	return s.Renew != nil
}

func (s *Secret) Response(
	data, internal map[string]interface{}) *logical.Response {
	internalData := make(map[string]interface{})
	for k, v := range internal {
		internalData[k] = v
	}
	internalData["secret_type"] = s.Type

	return &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       s.DefaultDuration,
				Renewable: s.Renewable(),
			},
			InternalData: internalData,
		},

		Data: data,
	}
}

// HandleRenew is the request handler for renewing this secret.
func (s *Secret) HandleRenew(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	if !s.Renewable() {
		return nil, logical.ErrUnsupportedOperation
	}

	data := &FieldData{
		Raw:    req.Data,
		Schema: s.Fields,
	}

	return s.Renew(ctx, req, data)
}

// HandleRevoke is the request handler for revoking this secret.
func (s *Secret) HandleRevoke(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	data := &FieldData{
		Raw:    req.Data,
		Schema: s.Fields,
	}

	if s.Revoke != nil {
		return s.Revoke(ctx, req, data)
	}

	return nil, logical.ErrUnsupportedOperation
}
