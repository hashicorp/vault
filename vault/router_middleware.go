package vault

import (
	"context"

	"github.com/hashicorp/vault/logical"
)

type LogicalHandlerFunc func(context.Context, *logical.Request) (*logical.Response, error)

type RouterMiddleware interface {
	Update(map[string]interface{}) error
	Handler(LogicalHandlerFunc) LogicalHandlerFunc
}
