package mock

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// New returns a new backend as an interface. This func
// is only necessary for builtin backend plugins.
func New() (interface{}, error) {
	return Backend(), nil
}

// Factory returns a new backend as logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	b.testNamespace = namespace.Namespace{
		ID:   conf.Config["nsID"],
		Path: conf.Config["nsPath"],
	}

	if err := b.validateCtxNamespace(ctx); err != nil {
		return nil, err
	}
	return b, nil
}

// FactoryType is a wrapper func that allows the Factory func to specify
// the backend type for the mock backend plugin instance.
func FactoryType(backendType logical.BackendType) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		b := Backend()
		b.BackendType = backendType
		if err := b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
}

// Backend returns a private embedded struct of framework.Backend.
func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: "",
		Paths: framework.PathAppend(
			errorPaths(&b),
			kvPaths(&b),
			[]*framework.Path{
				pathInternal(&b),
				pathSpecial(&b),
				pathRaw(&b),
			},
		),
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"special",
			},
		},
		Secrets:        []*framework.Secret{},
		Invalidate:     b.invalidate,
		InitializeFunc: b.initialize,
		Clean:          b.clean,
		BackendType:    logical.TypeLogical,
	}
	b.internal = "bar"
	return &b
}

type backend struct {
	*framework.Backend

	testNamespace namespace.Namespace

	// internal is used to test invalidate
	internal string
}

func (b *backend) invalidate(ctx context.Context, key string) {
	if err := b.validateCtxNamespace(ctx); err != nil {
		panic(err)
	}

	switch key {
	case "internal":
		b.internal = ""
	}
}

func (b *backend) initialize(ctx context.Context, _ *logical.InitializationRequest) error {
	if err := b.validateCtxNamespace(ctx); err != nil {
		return err
	}

	return nil
}

func (b *backend) clean(ctx context.Context) {
	if err := b.validateCtxNamespace(ctx); err != nil {
		panic(err)
	}
}

func (b *backend) validateCtxNamespace(ctx context.Context) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	if *ns != b.testNamespace {
		return fmt.Errorf("expected namespace: %+v, got: %+v", b.testNamespace, ns)
	}

	return nil
}
