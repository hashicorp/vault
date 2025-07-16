// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mock

import (
	"context"
	"errors"
	"net/rpc"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
)

// pathInternal is used to test viewing internal backend values. In this case,
// it is used to test the invalidate func.
func errorPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "errors/rpc",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathErrorRPCRead,
			},
		},
		{
			Pattern: "errors/kill",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathErrorRPCRead,
			},
		},
		{
			Pattern: "errors/type",
			Fields: map[string]*framework.FieldSchema{
				"err_type": {Type: framework.TypeInt},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathErrorRPCRead,
			},
		},
	}
}

func (b *backend) pathErrorRPCRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	errTypeRaw, ok := data.GetOk("err_type")
	if !ok {
		return nil, rpc.ErrShutdown
	}

	var err error
	switch uint32(errTypeRaw.(int)) {
	case pb.ErrTypeUnknown:
		err = errors.New("test")
	case pb.ErrTypeUserError:
		err = errutil.UserError{Err: "test"}
	case pb.ErrTypeInternalError:
		err = errutil.InternalError{Err: "test"}
	case pb.ErrTypeCodedError:
		err = logical.CodedError(403, "test")
	case pb.ErrTypeStatusBadRequest:
		err = &logical.StatusBadRequest{Err: "test"}
	case pb.ErrTypeUnsupportedOperation:
		err = logical.ErrUnsupportedOperation
	case pb.ErrTypeUnsupportedPath:
		err = logical.ErrUnsupportedPath
	case pb.ErrTypeInvalidRequest:
		err = logical.ErrInvalidRequest
	case pb.ErrTypePermissionDenied:
		err = logical.ErrPermissionDenied
	case pb.ErrTypeMultiAuthzPending:
		err = logical.ErrMultiAuthzPending
	}

	return nil, err
}
