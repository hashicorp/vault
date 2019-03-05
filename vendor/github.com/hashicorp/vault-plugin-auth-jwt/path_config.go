package jwtauth

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"

	"context"

	oidc "github.com/coreos/go-oidc"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/oauth2"
)

func pathConfig(b *jwtAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		Fields: map[string]*framework.FieldSchema{
			"oidc_discovery_url": {
				Type:        framework.TypeString,
				Description: `OIDC Discovery URL, without any .well-known component (base path). Cannot be used with "jwt_validation_pubkeys".`,
			},
			"oidc_discovery_ca_pem": {
				Type:        framework.TypeString,
				Description: "The CA certificate or chain of certificates, in PEM format, to use to validate conections to the OIDC Discovery URL. If not set, system certificates are used.",
			},
			"oidc_client_id": {
				Type:        framework.TypeString,
				Description: "The OAuth Client ID configured with your OIDC provider.",
			},
			"oidc_client_secret": {
				Type:             framework.TypeString,
				Description:      "The OAuth Client Secret configured with your OIDC provider.",
				DisplaySensitive: true,
			},
			"default_role": {
				Type:        framework.TypeString,
				Description: "The default role to use if none is provided during login. If not set, a role is required during login.",
			},
			"jwt_validation_pubkeys": {
				Type:        framework.TypeCommaStringSlice,
				Description: `A list of PEM-encoded public keys to use to authenticate signatures locally. Cannot be used with "oidc_discovery_url".`,
			},
			"jwt_supported_algs": {
				Type:        framework.TypeCommaStringSlice,
				Description: `A list of supported signing algorithms. Defaults to RS256.`,
			},
			"bound_issuer": {
				Type:        framework.TypeString,
				Description: "The value against which to match the 'iss' claim in a JWT. Optional.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				Summary:  "Read the current JWT authentication backend configuration.",
			},

			logical.UpdateOperation: &framework.PathOperation{
				Callback:    b.pathConfigWrite,
				Summary:     "Configure the JWT authentication backend.",
				Description: confHelpDesc,
			},
		},

		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}
}

func (b *jwtAuthBackend) config(ctx context.Context, s logical.Storage) (*jwtConfig, error) {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.cachedConfig != nil {
		return b.cachedConfig, nil
	}

	entry, err := s.Get(ctx, configPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	result := &jwtConfig{}
	if entry != nil {
		if err := entry.DecodeJSON(result); err != nil {
			return nil, err
		}
	}

	for _, v := range result.JWTValidationPubKeys {
		key, err := certutil.ParsePublicKeyPEM([]byte(v))
		if err != nil {
			return nil, errwrap.Wrapf("error parsing public key: {{err}}", err)
		}
		result.ParsedJWTPubKeys = append(result.ParsedJWTPubKeys, key)
	}

	b.cachedConfig = result

	return result, nil
}

func (b *jwtAuthBackend) pathConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"oidc_discovery_url":     config.OIDCDiscoveryURL,
			"oidc_discovery_ca_pem":  config.OIDCDiscoveryCAPEM,
			"oidc_client_id":         config.OIDCClientID,
			"default_role":           config.DefaultRole,
			"jwt_validation_pubkeys": config.JWTValidationPubKeys,
			"jwt_supported_algs":     config.JWTSupportedAlgs,
			"bound_issuer":           config.BoundIssuer,
		},
	}

	return resp, nil
}

func (b *jwtAuthBackend) pathConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config := &jwtConfig{
		OIDCDiscoveryURL:     d.Get("oidc_discovery_url").(string),
		OIDCDiscoveryCAPEM:   d.Get("oidc_discovery_ca_pem").(string),
		OIDCClientID:         d.Get("oidc_client_id").(string),
		OIDCClientSecret:     d.Get("oidc_client_secret").(string),
		DefaultRole:          d.Get("default_role").(string),
		JWTValidationPubKeys: d.Get("jwt_validation_pubkeys").([]string),
		JWTSupportedAlgs:     d.Get("jwt_supported_algs").([]string),
		BoundIssuer:          d.Get("bound_issuer").(string),
	}

	// Run checks on values
	switch {
	case config.OIDCDiscoveryURL == "" && len(config.JWTValidationPubKeys) == 0,
		config.OIDCDiscoveryURL != "" && len(config.JWTValidationPubKeys) != 0:
		return logical.ErrorResponse("exactly one of 'oidc_discovery_url' and 'jwt_validation_pubkeys' must be set"), nil

	case config.OIDCClientID != "" && config.OIDCClientSecret == "",
		config.OIDCClientID == "" && config.OIDCClientSecret != "":
		return logical.ErrorResponse("both 'oidc_client_id' and 'oidc_client_secret' must be set for OIDC"), nil

	case config.OIDCDiscoveryURL != "":
		_, err := b.createProvider(config)
		if err != nil {
			return logical.ErrorResponse(errwrap.Wrapf("error checking discovery URL: {{err}}", err).Error()), nil
		}

	case config.OIDCClientID != "" && config.OIDCDiscoveryURL == "":
		return logical.ErrorResponse("'oidc_discovery_url' must be set for OIDC"), nil

	case len(config.JWTValidationPubKeys) != 0:
		for _, v := range config.JWTValidationPubKeys {
			if _, err := certutil.ParsePublicKeyPEM([]byte(v)); err != nil {
				return logical.ErrorResponse(errwrap.Wrapf("error parsing public key: {{err}}", err).Error()), nil
			}
		}

	default:
		return nil, errors.New("unknown condition")
	}

	for _, a := range config.JWTSupportedAlgs {
		switch a {
		case oidc.RS256, oidc.RS384, oidc.RS512, oidc.ES256, oidc.ES384, oidc.ES512, oidc.PS256, oidc.PS384, oidc.PS512:
		default:
			return logical.ErrorResponse(fmt.Sprintf("Invalid supported algorithm: %s", a)), nil
		}
	}

	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	b.reset()

	return nil, nil
}

func (b *jwtAuthBackend) createProvider(config *jwtConfig) (*oidc.Provider, error) {
	var certPool *x509.CertPool
	if config.OIDCDiscoveryCAPEM != "" {
		certPool = x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM([]byte(config.OIDCDiscoveryCAPEM)); !ok {
			return nil, errors.New("could not parse 'oidc_discovery_ca_pem' value successfully")
		}
	}

	tr := cleanhttp.DefaultPooledTransport()
	if certPool != nil {
		tr.TLSClientConfig = &tls.Config{
			RootCAs: certPool,
		}
	}
	tc := &http.Client{
		Transport: tr,
	}
	oidcCtx := context.WithValue(b.providerCtx, oauth2.HTTPClient, tc)

	provider, err := oidc.NewProvider(oidcCtx, config.OIDCDiscoveryURL)
	if err != nil {
		return nil, errwrap.Wrapf("error creating provider with given values: {{err}}", err)
	}

	return provider, nil
}

type jwtConfig struct {
	OIDCDiscoveryURL     string   `json:"oidc_discovery_url"`
	OIDCDiscoveryCAPEM   string   `json:"oidc_discovery_ca_pem"`
	OIDCClientID         string   `json:"oidc_client_id"`
	OIDCClientSecret     string   `json:"oidc_client_secret"`
	JWTValidationPubKeys []string `json:"jwt_validation_pubkeys"`
	JWTSupportedAlgs     []string `json:"jwt_supported_algs"`
	BoundIssuer          string   `json:"bound_issuer"`
	DefaultRole          string   `json:"default_role"`

	ParsedJWTPubKeys []interface{} `json:"-"`
}

const (
	confHelpSyn = `
Configures the JWT authentication backend.
`
	confHelpDesc = `
The JWT authentication backend validates JWTs (or OIDC) using the configured
credentials. If using OIDC Discovery, the URL must be provided, along
with (optionally) the CA cert to use for the connection. If performing JWT
validation locally, a set of public keys must be provided.
`
)
