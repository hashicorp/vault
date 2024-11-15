// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubeauth

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configPath = "config"
	rolePrefix = "role/"

	// aliasNameSourceUnset provides backwards compatibility with preexisting roles.
	aliasNameSourceUnset   = ""
	aliasNameSourceSAUid   = "serviceaccount_uid"
	aliasNameSourceSAName  = "serviceaccount_name"
	aliasNameSourceDefault = aliasNameSourceSAUid
	minTLSVersion          = tls.VersionTLS12

	// operationPrefixKubernetes is used as a prefix for OpenAPI operation id's.
	operationPrefixKubernetes = "kubernetes"
)

var (
	// when adding new alias name sources make sure to update the corresponding FieldSchema description in path_role.go
	aliasNameSources          = []string{aliasNameSourceSAUid, aliasNameSourceSAName}
	errInvalidAliasNameSource = fmt.Errorf(`invalid alias_name_source, must be one of: %s`, strings.Join(aliasNameSources, ", "))

	// jwtReloadPeriod is the time period how often the in-memory copy of local
	// service account token can be used, before reading it again from disk.
	//
	// The value is selected according to recommendation in Kubernetes 1.21 changelog:
	// "Clients should reload the token from disk periodically (once per minute
	// is recommended) to ensure they continue to use a valid token."
	jwtReloadPeriod = 1 * time.Minute

	// caReloadPeriod is the time period how often the in-memory copy of local
	// CA cert can be used, before reading it again from disk.
	caReloadPeriod = 1 * time.Hour

	// defaultHorizon provides the default duration to be used
	// in the tlsConfigUpdater's time.Ticker, setup in runTLSConfigUpdater()
	defaultHorizon = time.Second * 30

	// defaultMinHorizon provides the minimum duration that can be specified
	// in the tlsConfigUpdater's time.Ticker, setup in runTLSConfigUpdater()
	defaultMinHorizon = time.Second * 5

	errTLSConfigNotSet  = errors.New("TLSConfig not set")
	errHTTPClientNotSet = errors.New("http.Client not set")
)

// kubeAuthBackend implements logical.Backend
type kubeAuthBackend struct {
	*framework.Backend

	// default HTTP client for connection reuse
	httpClient *http.Client

	// tlsConfig is periodically updated whenever the CA certificate configuration changes.
	tlsConfig *tls.Config

	// reviewFactory is used to configure the strategy for doing a token review.
	// Currently, the only options are using the kubernetes API or mocking the
	// review. Mocks should only be used in tests.
	reviewFactory tokenReviewFactory

	// namespaceValidatorFactory is used to configure the strategy for validating
	// namespace properties (currently labels). Currently, the only options
	// are using the kubernetes API or mocking the validation. Mocks should
	// only be used in tests.
	namespaceValidatorFactory namespaceValidatorFactory

	// serviceAccountGetterFactory is used to configure the strategy for retrieving
	// service account properties (currently metadata). Currently, the only options
	// are using the kubernetes API or mocking the retrieval. Mocks should
	// only be used in tests.
	serviceAccountGetterFactory serviceAccountGetterFactory

	// localSATokenReader caches the service account token in memory.
	// It periodically reloads the token to support token rotation/renewal.
	// Local token is used when running in a pod with following configuration
	// - token_reviewer_jwt is not set
	// - disable_local_ca_jwt is false
	localSATokenReader *cachingFileReader

	// localCACertReader contains the local CA certificate. Local CA certificate is
	// used when running in a pod with following configuration
	// - kubernetes_ca_cert is not set
	// - disable_local_ca_jwt is false
	localCACertReader *cachingFileReader

	// tlsConfigUpdaterRunning is used to signal the current state of the tlsConfig updater routine.
	tlsConfigUpdaterRunning bool

	// tlsConfigUpdateCancelFunc should be called in the backend's Clean(), set in initialize().
	tlsConfigUpdateCancelFunc context.CancelFunc

	l sync.RWMutex

	// tlsMu provides the lock for synchronizing updates to the tlsConfig.
	tlsMu sync.RWMutex
}

// Factory returns a new backend as logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

// getDefaultTLSConfig returns a http.Client with the supported, default
// tls.Config from getDefaultTLSConfig() set on the
// cleanhttp.DefaultPooledTransport() http.Transport.
func getDefaultHTTPClient() *http.Client {
	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = getDefaultTLSConfig()
	return &http.Client{
		Transport: transport,
	}
}

func getDefaultTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: minTLSVersion,
	}
}

func Backend() *kubeAuthBackend {
	b := &kubeAuthBackend{
		localSATokenReader: newCachingFileReader(localJWTPath, jwtReloadPeriod, time.Now),
		localCACertReader:  newCachingFileReader(localCACertPath, caReloadPeriod, time.Now),
		// Set default HTTP client
		httpClient: getDefaultHTTPClient(),
		// Set the default TLSConfig
		tlsConfig: getDefaultTLSConfig(),
		// Set the review factory to default to calling into the kubernetes API.
		reviewFactory:               tokenReviewAPIFactory,
		namespaceValidatorFactory:   newNsValidatorWrapper,
		serviceAccountGetterFactory: newServiceAccountGetterWrapper,
	}

	b.Backend = &framework.Backend{
		AuthRenew:   b.pathLoginRenew(),
		BackendType: logical.TypeCredential,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
			SealWrapStorage: []string{
				configPath,
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathConfig(b),
				pathLogin(b),
			},
			pathsRole(b),
		),
		InitializeFunc: b.initialize,
		Clean:          b.cleanup,
	}

	return b
}

// initialize is used to handle the state of config values just after the K8s plugin has been mounted
func (b *kubeAuthBackend) initialize(ctx context.Context, req *logical.InitializationRequest) error {
	updaterCtx, cancel := context.WithCancel(context.Background())
	if err := b.runTLSConfigUpdater(updaterCtx, req.Storage, defaultHorizon); err != nil {
		cancel()
		return err
	}

	b.tlsConfigUpdateCancelFunc = cancel

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return err
	}

	if config != nil {
		if err := b.updateTLSConfig(config); err != nil {
			return err
		}
	}

	return nil
}

func (b *kubeAuthBackend) cleanup(_ context.Context) {
	b.shutdownTLSConfigUpdater()
}

// validateHTTPClientInit that the Backend's HTTPClient and TLSConfig has been properly instantiated.
func (b *kubeAuthBackend) validateHTTPClientInit() error {
	if b.httpClient == nil {
		return errHTTPClientNotSet
	}
	if b.tlsConfig == nil {
		return errTLSConfigNotSet
	}

	return nil
}

// runTLSConfigUpdater sets up a routine that periodically calls b.updateTLSConfig(). This ensures that the
// httpClient's TLS configuration is consistent with the backend's stored configuration.
func (b *kubeAuthBackend) runTLSConfigUpdater(ctx context.Context, s logical.Storage, horizon time.Duration) error {
	b.tlsMu.Lock()
	defer b.tlsMu.Unlock()

	if b.tlsConfigUpdaterRunning {
		return nil
	}

	if horizon < defaultMinHorizon {
		return fmt.Errorf("update horizon must be equal to or greater than %s", defaultMinHorizon)
	}

	if err := b.validateHTTPClientInit(); err != nil {
		return err
	}

	updateTLSConfig := func(ctx context.Context, s logical.Storage) error {
		config, err := b.config(ctx, s)
		if err != nil {
			return fmt.Errorf("failed config read, err=%w", err)
		}

		if config == nil {
			b.Logger().Trace("Skipping TLSConfig update, no configuration set")
			return nil
		}

		if err := b.updateTLSConfig(config); err != nil {
			return err
		}

		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	ticker := time.NewTicker(horizon)
	go func(ctx context.Context, s logical.Storage) {
		defer func() {
			b.tlsMu.Lock()
			defer b.tlsMu.Unlock()
			ticker.Stop()
			b.tlsConfigUpdaterRunning = false
			b.Logger().Trace("TLSConfig updater shutdown completed")
		}()

		b.Logger().Trace("TLSConfig updater starting", "horizon", horizon)
		b.tlsConfigUpdaterRunning = true
		wg.Done()
		for {
			select {
			case <-ctx.Done():
				b.Logger().Trace("TLSConfig updater shutting down")
				return
			case <-ticker.C:
				if err := updateTLSConfig(ctx, s); err != nil {
					b.Logger().Warn("TLSConfig update failed, retrying",
						"horizon", defaultHorizon.String(), "err", err)
				}
			}
		}
	}(ctx, s)
	wg.Wait()

	return nil
}

func (b *kubeAuthBackend) shutdownTLSConfigUpdater() {
	if b.tlsConfigUpdateCancelFunc != nil {
		b.Logger().Debug("TLSConfig updater shutdown requested")
		b.tlsConfigUpdateCancelFunc()
		b.tlsConfigUpdateCancelFunc = nil
	}
}

// config takes a storage object and returns a kubeConfig object.
// It does not return local token and CA file which are specific to the pod we run in.
func (b *kubeAuthBackend) config(ctx context.Context, s logical.Storage) (*kubeConfig, error) {
	raw, err := s.Get(ctx, configPath)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	conf := &kubeConfig{}
	if err := json.Unmarshal(raw.Value, conf); err != nil {
		return nil, err
	}

	// Parse the public keys from the CertificatesBytes
	conf.PublicKeys = make([]crypto.PublicKey, len(conf.PEMKeys))
	for i, cert := range conf.PEMKeys {
		conf.PublicKeys[i], err = parsePublicKeyPEM([]byte(cert))
		if err != nil {
			return nil, err
		}
	}

	return conf, nil
}

// loadConfig fetches the kubeConfig from storage and optionally decorates it with
// local token and CA certificate. Since loadConfig does not return an error if the kubeConfig reference
// is nil, we should nil-check. This behavior exists to allow loadConfig's caller to
// make a decision based on the returned reference.
func (b *kubeAuthBackend) loadConfig(ctx context.Context, s logical.Storage) (*kubeConfig, error) {
	config, err := b.config(ctx, s)
	if err != nil {
		return nil, err
	}
	// We know the config is empty so exit early
	if config == nil {
		return config, nil
	}
	// Nothing more to do if loading local CA cert and JWT token is disabled.
	if config.DisableLocalCAJwt {
		return config, nil
	}

	// Read local JWT token unless it was not stored in config.
	if config.TokenReviewerJWT == "" {
		config.TokenReviewerJWT, err = b.localSATokenReader.ReadFile()
		if err != nil {
			// Ignore error: make the best effort trying to load local JWT,
			// otherwise the JWT submitted in login payload will be used.
			b.Logger().Debug("failed to read local service account token, will use client token", "error", err)
		}
	}

	// Read local CA cert unless it was stored in config.
	// Else build the TLSConfig with the trusted CA cert and load into client
	if config.CACert == "" {
		config.CACert, err = b.localCACertReader.ReadFile()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

// role takes a storage backend and the name and returns the role's storage
// entry
func (b *kubeAuthBackend) role(ctx context.Context, s logical.Storage, name string) (*roleStorageEntry, error) {
	raw, err := s.Get(ctx, fmt.Sprintf("%s%s", rolePrefix, strings.ToLower(name)))
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	role := &roleStorageEntry{}
	if err := json.Unmarshal(raw.Value, role); err != nil {
		return nil, err
	}

	if role.TokenTTL == 0 && role.TTL > 0 {
		role.TokenTTL = role.TTL
	}
	if role.TokenMaxTTL == 0 && role.MaxTTL > 0 {
		role.TokenMaxTTL = role.MaxTTL
	}
	if role.TokenPeriod == 0 && role.Period > 0 {
		role.TokenPeriod = role.Period
	}
	if role.TokenNumUses == 0 && role.NumUses > 0 {
		role.TokenNumUses = role.NumUses
	}
	if len(role.TokenPolicies) == 0 && len(role.Policies) > 0 {
		role.TokenPolicies = role.Policies
	}
	if len(role.TokenBoundCIDRs) == 0 && len(role.BoundCIDRs) > 0 {
		role.TokenBoundCIDRs = role.BoundCIDRs
	}

	return role, nil
}

// getHTTPClient return the backend's HTTP client for connecting to the Kubernetes API.
func (b *kubeAuthBackend) getHTTPClient() (*http.Client, error) {
	b.tlsMu.RLock()
	defer b.tlsMu.RUnlock()

	if err := b.validateHTTPClientInit(); err != nil {
		return nil, err
	}

	return b.httpClient, nil
}

// updateTLSConfig ensures that the httpClient's TLS configuration is consistent
// with the backend's stored configuration.
func (b *kubeAuthBackend) updateTLSConfig(config *kubeConfig) error {
	b.tlsMu.Lock()
	defer b.tlsMu.Unlock()

	if err := b.validateHTTPClientInit(); err != nil {
		return err
	}

	// attempt to read the CA certificates from the config directly or from the filesystem.
	var caCertBytes []byte
	if config.CACert != "" {
		caCertBytes = []byte(config.CACert)
	} else if !config.DisableLocalCAJwt && b.localCACertReader != nil {
		data, err := b.localCACertReader.ReadFile()
		if err != nil {
			return err
		}
		caCertBytes = []byte(data)
	}

	var certPool *x509.CertPool
	if len(caCertBytes) == 0 {
		// since the CA chain is not configured, we use the system's cert pool.
		var err error
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return err
		}
	} else {
		// since we have a CA chain configured, we create a new x509.CertPool with its
		// contents.
		certPool = x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(caCertBytes); !ok {
			b.Logger().Warn("Configured CA PEM data contains no valid certificates, TLS verification will fail")
		}
	}

	setTLSClientConfig := func() error {
		transport, ok := b.httpClient.Transport.(*http.Transport)
		if !ok {
			// should never happen
			return fmt.Errorf("type assertion failed for %T", b.httpClient.Transport)
		}

		b.tlsConfig.RootCAs = certPool
		transport.TLSClientConfig = b.tlsConfig
		return nil
	}

	// only refresh the Root CAs if they have changed since the last full update.
	if b.tlsConfig.RootCAs == nil {
		return setTLSClientConfig()
	} else if !b.tlsConfig.RootCAs.Equal(certPool) {
		b.Logger().Trace("Root CA certificate pool has changed, updating the client's transport")
		return setTLSClientConfig()
	} else {
		b.Logger().Trace("Root CA certificate pool is unchanged, no update required")
	}

	return nil
}

func validateAliasNameSource(source string) error {
	for _, s := range aliasNameSources {
		if s == source {
			return nil
		}
	}
	return errInvalidAliasNameSource
}

var backendHelp string = `
The Kubernetes Auth Backend allows authentication for Kubernetes service accounts.
`
