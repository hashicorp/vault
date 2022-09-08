package vault

import (
	"context"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

// builtinVersionBackend wraps a (presumably builtin) logical.Backend and adds the current builtin version.
type builtinVersionBackend struct {
	backend logical.Backend
}

func (b builtinVersionBackend) Initialize(ctx context.Context, request *logical.InitializationRequest) error {
	return b.backend.Initialize(ctx, request)
}

func (b builtinVersionBackend) HandleRequest(ctx context.Context, request *logical.Request) (*logical.Response, error) {
	return b.backend.HandleRequest(ctx, request)
}

func (b builtinVersionBackend) SpecialPaths() *logical.Paths {
	return b.backend.SpecialPaths()
}

func (b builtinVersionBackend) System() logical.SystemView {
	return b.backend.System()
}

func (b builtinVersionBackend) Logger() log.Logger {
	return b.backend.Logger()
}

func (b builtinVersionBackend) HandleExistenceCheck(ctx context.Context, request *logical.Request) (bool, bool, error) {
	return b.backend.HandleExistenceCheck(ctx, request)
}

func (b builtinVersionBackend) Cleanup(ctx context.Context) {
	b.backend.Cleanup(ctx)
}

func (b builtinVersionBackend) InvalidateKey(ctx context.Context, s string) {
	b.backend.InvalidateKey(ctx, s)
}

func (b builtinVersionBackend) Setup(ctx context.Context, config *logical.BackendConfig) error {
	return b.backend.Setup(ctx, config)
}

func (b builtinVersionBackend) Type() logical.BackendType {
	return b.backend.Type()
}

// Version returns the builtin version for Vault
func (b builtinVersionBackend) Version() logical.VersionInfo {
	return logical.BuiltinVersion
}

var _ logical.Backend = (*builtinVersionBackend)(nil)

func wrapFactoryAddBuiltinVersion(factory logical.Factory) logical.Factory {
	return func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		b, err := factory(ctx, config)
		if err != nil {
			return nil, err
		}
		return builtinVersionBackend{backend: b}, nil
	}
}

// wrapMapAddBuiltinVersion is a convenience tool to wrap all the logical.Factory interfaces in a
// map with ones that return the builtin version.
func wrapMapAddBuiltinVersion(m map[string]logical.Factory) map[string]logical.Factory {
	newMap := make(map[string]logical.Factory)
	for k, v := range m {
		newMap[k] = wrapFactoryAddBuiltinVersion(v)
	}
	return newMap
}
