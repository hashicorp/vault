package util

import (
	"context"

	"github.com/hashicorp/vault/logical/framework"
)

func Invalidator(invalidationFuncs ...framework.InvalidateFunc) *invalidator {
	return &invalidator{invalidationFuncs}
}

type invalidator struct {
	toCall []framework.InvalidateFunc
}

func (v *invalidator) Invalidate(ctx context.Context, key string) {
	for _, f := range v.toCall {
		f(ctx, key)
	}
}
