package kv

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// PassthroughDowngrader wraps a normal passthrough backend and downgrades the
// request object from the newer Versioned API to the older Passthrough API.
// This allows us to use the new "vault kv" subcommand with a non-versioned
// instance of the kv store without doing a preflight API version check. The
// CLI will always use the new API definition and this object will make it
// compatible with the passthrough backend. The "X-Vault-Kv-Client" header is
// used to know the request originated from the CLI and uses the newer API.
type PassthroughDowngrader struct {
	next Passthrough
}

func (b *PassthroughDowngrader) handleExistenceCheck() framework.ExistenceFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
		if !b.shouldDowngrade(req) {
			return b.next.handleExistenceCheck()(ctx, req, data)
		}

		respErr := b.invalidPath(req)
		if respErr != nil {
			return false, logical.ErrInvalidRequest
		}

		reqDown := &logical.Request{}
		*reqDown = *req

		reqDown.Path = strings.TrimPrefix(req.Path, "data/")
		return b.next.handleExistenceCheck()(ctx, reqDown, data)
	}
}

func (b *PassthroughDowngrader) handleRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		if !b.shouldDowngrade(req) {
			return b.next.handleRead()(ctx, req, data)
		}

		respErr := b.invalidPath(req)
		if respErr != nil {
			return respErr, logical.ErrInvalidRequest
		}

		if _, ok := data.Raw["version"]; ok {
			return logical.ErrorResponse("retrieving a version is not supported when versioning is disabled"), logical.ErrInvalidRequest
		}

		reqDown := &logical.Request{}
		*reqDown = *req

		reqDown.Path = strings.TrimPrefix(req.Path, "data/")

		resp, err := b.next.handleRead()(ctx, reqDown, data)
		if resp != nil && resp.Data != nil {
			resp.Data = map[string]interface{}{
				"data":     resp.Data,
				"metadata": nil,
			}
		}

		return resp, err
	}
}

func (b *PassthroughDowngrader) handleWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		if !b.shouldDowngrade(req) {
			return b.next.handleWrite()(ctx, req, data)
		}

		respErr := b.invalidPath(req)
		if respErr != nil {
			return respErr, logical.ErrInvalidRequest
		}

		reqDown := &logical.Request{}
		*reqDown = *req
		reqDown.Path = strings.TrimPrefix(req.Path, "data/")

		// Validate the data map is what we expect
		switch req.Data["data"].(type) {
		case map[string]interface{}:
		default:
			return logical.ErrorResponse("could not downgrade request, unexpected data format"), logical.ErrInvalidRequest
		}

		// Move the data object up a level and ignore the options object.
		reqDown.Data = req.Data["data"].(map[string]interface{})

		return b.next.handleWrite()(ctx, reqDown, data)
	}
}

func (b *PassthroughDowngrader) handleDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		if !b.shouldDowngrade(req) {
			return b.next.handleDelete()(ctx, req, data)
		}

		respErr := b.invalidPath(req)
		if respErr != nil {
			return respErr, logical.ErrInvalidRequest
		}

		reqDown := &logical.Request{}
		*reqDown = *req
		reqDown.Path = strings.TrimPrefix(req.Path, "data/")

		return b.next.handleDelete()(ctx, reqDown, data)
	}
}

func (b *PassthroughDowngrader) handleList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		if !b.shouldDowngrade(req) {
			return b.next.handleList()(ctx, req, data)
		}

		reqDown := &logical.Request{}
		*reqDown = *req
		reqDown.Path = strings.TrimPrefix(req.Path, "metadata/")

		return b.next.handleList()(ctx, reqDown, data)
	}
}

func (b *PassthroughDowngrader) shouldDowngrade(req *logical.Request) bool {
	return http.Header(req.Headers).Get(consts.VaultKVCLIClientHeader) != ""
}

// invalidPaths returns an error if we are trying to access an versioned only
// path on a non-versioned kv store.
func (b *PassthroughDowngrader) invalidPath(req *logical.Request) *logical.Response {
	switch {
	case req.Path == "config":
		fallthrough
	case strings.HasPrefix(req.Path, "metadata/"):
		fallthrough
	case strings.HasPrefix(req.Path, "archive/"):
		fallthrough
	case strings.HasPrefix(req.Path, "unarchive/"):
		fallthrough
	case strings.HasPrefix(req.Path, "destroy/"):
		return logical.ErrorResponse("path is not supported when versioning is disabled")
	}

	return nil
}
