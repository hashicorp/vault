// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cf

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/vault-plugin-auth-cf/models"
)

const (
	// These env vars are used frequently to pull the client certificate and private key
	// from CF containers; thus are placed here for ease of discovery and use from
	// outside packages.
	EnvVarInstanceCertificate = "CF_INSTANCE_CERT"
	EnvVarInstanceKey         = "CF_INSTANCE_KEY"

	// operationPrefixCloudFoundry is used as a prefix for OpenAPI operation id's.
	operationPrefixCloudFoundry = "cloud-foundry"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &backend{}
	b.Backend = &framework.Backend{
		AuthRenew: b.pathLoginRenew,
		Help:      backendHelp,
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{"config"},
			Unauthenticated: []string{"login"},
		},
		Paths: []*framework.Path{
			b.pathConfig(),
			b.pathListRoles(),
			b.pathRoles(),
			b.pathLogin(),
		},
		BackendType:    logical.TypeCredential,
		InitializeFunc: b.initialize,
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend
	mu             sync.RWMutex
	cfClient       *cfclient.Client
	cfClientMu     sync.RWMutex
	lastConfigHash *[32]byte
}

const backendHelp = `
The CF auth backend supports logging in using CF's identity service.
Once a CA certificate is configured, and Vault is configured to consume
CF's API, CF's instance identity credentials can be used to authenticate.'
`

var errCFClientNotInitialized = fmt.Errorf("client is not initialized")

func (b *backend) getCFClient(_ context.Context) (*cfclient.Client, error) {
	b.cfClientMu.RLock()
	defer b.cfClientMu.RUnlock()
	if b.cfClient == nil {
		return nil, errCFClientNotInitialized
	}

	return b.cfClient, nil
}

func (b *backend) updateCFClient(ctx context.Context, config *models.Configuration) (bool, error) {
	b.cfClientMu.Lock()
	defer b.cfClientMu.Unlock()

	configHash, err := config.Hash()
	if err != nil {
		return false, err
	}

	if b.lastConfigHash != nil && b.cfClient != nil {
		if err == nil {
			if *b.lastConfigHash == configHash {
				return false, nil
			}
		}
	}

	if b.cfClient != nil {
		if b.cfClient.Config.HttpClient != nil {
			b.cfClient.Config.HttpClient.CloseIdleConnections()
		}
	}

	cfClient, err := b.newCFClient(ctx, config)
	if err != nil {
		return false, err
	}

	b.cfClient = cfClient
	b.lastConfigHash = &configHash

	return true, nil
}

func (b *backend) getCFClientOrRefresh(ctx context.Context, config *models.Configuration) (*cfclient.Client, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration is nil")
	}

	client, err := b.getCFClient(ctx)
	if err != nil {
		if errors.Is(err, errCFClientNotInitialized) {
			if _, err := b.updateCFClient(ctx, config); err != nil {
				return nil, err
			}
			return b.getCFClient(ctx)
		}
		return nil, err
	}

	return client, nil
}

func (b *backend) newCFClient(_ context.Context, config *models.Configuration) (*cfclient.Client, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration is nil")
	}

	clientConf := &cfclient.Config{
		ApiAddress:   config.CFAPIAddr,
		Username:     config.CFUsername,
		Password:     config.CFPassword,
		ClientID:     config.CFClientID,
		ClientSecret: config.CFClientSecret,
		HttpClient:   cleanhttp.DefaultClient(),
	}
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	for idx, certificate := range config.CFAPICertificates {
		if ok := rootCAs.AppendCertsFromPEM([]byte(certificate)); !ok {
			return nil, fmt.Errorf(
				"failed to append CF API cert to cert pool, index=%d, err=%w", idx, err,
			)
		}
	}
	tlsConfig := &tls.Config{
		RootCAs: rootCAs,
	}

	if config.CFMutualTLSCertificate != "" && config.CFMutualTLSKey != "" {
		cert, err := tls.X509KeyPair(
			[]byte(config.CFMutualTLSCertificate),
			[]byte(config.CFMutualTLSKey),
		)

		if err != nil {
			return nil, fmt.Errorf("could not parse X509 key pair for mutual TLS")
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	clientConf.HttpClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// unfortunately, cfclient.NewClient() has a nasty side effect of reaching out
	// to the CF API. That means that the CF API must be reachable at the time of
	// the call. The v3 of go-cfclient does not have this issue. Updating to v3
	// should be a priority.
	return cfclient.NewClient(clientConf)
}

func (b *backend) initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if req == nil {
		return fmt.Errorf("initialization request is nil")
	}

	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		b.Logger().Warn("init: failed to get the config from storage", "error", err)
		return nil
	}

	if config != nil {
		if _, err := b.updateCFClient(ctx, config); err != nil {
			// We only log an error here, since we want the plugin to be able to come up.
			// Subsequent calls to the plugin will attempt to update the client again.
			b.Logger().Warn("init: failed to update CF client", "error", err)
		}
	}
	return nil
}
