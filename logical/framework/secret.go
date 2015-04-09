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

	// Renew is the callback called to renew this secret. If Renew is
	// not specified and RenewExtend is false, then Renewable is set to
	// false in the secret.
	//
	// RenewExtend, if true, will automatically extend the lease of this
	// secret type. You can specify RenewExtendMax to specify the max
	// duration it can be extended, otherwise it will be extended potentially
	// indefinitely.
	Renew          OperationFunc
	RenewExtend    bool
	RenewExtendMax time.Duration

	// Revoke is the callback called to revoke this secret. This is required.
	Revoke OperationFunc
}

func (s *Secret) Renewable() bool {
	return s.Renew != nil || s.RenewExtend
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
				Lease:            s.DefaultDuration,
				LeaseGracePeriod: s.DefaultGracePeriod,
				Renewable:        s.Renewable(),
			},
			InternalData: internalData,
		},

		Data: data,
	}
}

// HandleRenew is the request handler for renewing this secret.
func (s *Secret) HandleRenew(req *logical.Request) (*logical.Response, error) {
	if !s.Renewable() {
		return nil, logical.ErrUnsupportedOperation
	}

	data := &FieldData{
		Raw:    req.Data,
		Schema: s.Fields,
	}

	// If we have a callback, we just call that and that does all the logic.
	if s.Renew != nil {
		return s.Renew(req, data)
	}

	// If we're using RenewExtend, then just automaticaly extend.
	if s.RenewExtend {
		return s.HandleRenewExtend(req, data)
	}

	return nil, logical.ErrUnsupportedOperation
}

// HandleRenewExtend is the OperationFunc that just extends the lease
// of the secret.
func (s *Secret) HandleRenewExtend(
	req *logical.Request, data *FieldData) (*logical.Response, error) {
	// First copy the original secret/data
	var resp logical.Response
	resp.Secret = req.Secret
	resp.Data = req.Data

	// Now extend the lease by the amount specified.
	resp.Secret.Lease = req.Secret.LeaseIncrement

	return &resp, nil
}

// HandleRevoke is the request handler for renewing this secret.
func (s *Secret) HandleRevoke(req *logical.Request) (*logical.Response, error) {
	data := &FieldData{
		Raw:    req.Data,
		Schema: s.Fields,
	}

	if s.Revoke != nil {
		return s.Revoke(req, data)
	}

	return nil, logical.ErrUnsupportedOperation
}
