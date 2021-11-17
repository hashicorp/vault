package jwtauth

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/cap/jwt"
	"github.com/hashicorp/cap/oidc"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/patrickmn/go-cache"
)

const (
	configPath string = "config"
	rolePrefix string = "role/"
)

// Factory is used by framework
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := backend()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

type jwtAuthBackend struct {
	*framework.Backend

	l            sync.RWMutex
	provider     *oidc.Provider
	validator    *jwt.Validator
	cachedConfig *jwtConfig
	oidcRequests *cache.Cache

	providerCtx       context.Context
	providerCtxCancel context.CancelFunc
}

func backend() *jwtAuthBackend {
	b := new(jwtAuthBackend)
	b.providerCtx, b.providerCtxCancel = context.WithCancel(context.Background())
	b.oidcRequests = cache.New(oidcRequestTimeout, oidcRequestCleanupInterval)

	b.Backend = &framework.Backend{
		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
		Invalidate:  b.invalidate,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
				"oidc/auth_url",
				"oidc/callback",

				// Uncomment to mount simple UI handler for local development
				// "ui",
			},
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathLogin(b),
				pathRoleList(b),
				pathRole(b),
				pathConfig(b),

				// Uncomment to mount simple UI handler for local development
				// pathUI(b),
			},
			pathOIDC(b),
		),
		Clean: b.cleanup,
	}

	return b
}

func (b *jwtAuthBackend) cleanup(_ context.Context) {
	b.l.Lock()
	if b.providerCtxCancel != nil {
		b.providerCtxCancel()
	}
	if b.provider != nil {
		b.provider.Done()
	}
	b.l.Unlock()
}

func (b *jwtAuthBackend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config":
		b.reset()
	}
}

func (b *jwtAuthBackend) reset() {
	b.l.Lock()
	if b.provider != nil {
		b.provider.Done()
	}
	b.provider = nil
	b.cachedConfig = nil
	b.validator = nil
	b.l.Unlock()
}

func (b *jwtAuthBackend) getProvider(config *jwtConfig) (*oidc.Provider, error) {
	b.l.Lock()
	defer b.l.Unlock()

	if b.provider != nil {
		return b.provider, nil
	}

	provider, err := b.createProvider(config)
	if err != nil {
		return nil, err
	}

	b.provider = provider
	return provider, nil
}

// jwtValidator returns a new JWT validator based on the provided config.
func (b *jwtAuthBackend) jwtValidator(config *jwtConfig) (*jwt.Validator, error) {
	b.l.Lock()
	defer b.l.Unlock()

	if b.validator != nil {
		return b.validator, nil
	}

	var err error
	var keySet jwt.KeySet

	// Configure the key set for the validator
	switch config.authType() {
	case JWKS:
		keySet, err = jwt.NewJSONWebKeySet(b.providerCtx, config.JWKSURL, config.JWKSCAPEM)
	case StaticKeys:
		keySet, err = jwt.NewStaticKeySet(config.ParsedJWTPubKeys)
	case OIDCDiscovery:
		keySet, err = jwt.NewOIDCDiscoveryKeySet(b.providerCtx, config.OIDCDiscoveryURL, config.OIDCDiscoveryCAPEM)
	default:
		return nil, errors.New("unsupported config type")
	}

	if err != nil {
		return nil, fmt.Errorf("keyset configuration error: %w", err)
	}

	validator, err := jwt.NewValidator(keySet)
	if err != nil {
		return nil, fmt.Errorf("JWT validator configuration error: %w", err)
	}

	b.validator = validator

	return b.validator, nil
}

const (
	backendHelp = `
The JWT backend plugin allows authentication using JWTs (including OIDC).
`
)
