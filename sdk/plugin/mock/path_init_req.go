// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package mock

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathInitReq exposes the fields from the last InitializationRequest so that
// gRPC tests can assert that mount metadata is forwarded correctly over the
// wire.
func pathInitReq(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "init-req",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathInitReqRead,
		},
	}
}

func (b *backend) pathInitReqRead(_ context.Context, _ *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if b.lastInitReq == nil {
		return &logical.Response{Data: map[string]interface{}{}}, nil
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"mount_point":    b.lastInitReq.MountPoint,
			"mount_type":     b.lastInitReq.MountType,
			"mount_accessor": b.lastInitReq.MountAccessor,
			"backend_uuid":   b.lastInitReq.BackendUUID,
		},
	}, nil
}
