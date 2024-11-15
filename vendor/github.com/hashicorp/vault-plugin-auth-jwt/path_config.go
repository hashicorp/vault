// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jwtauth

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/cap/jwt"
	"github.com/hashicorp/cap/oidc"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	httputil "github.com/hashicorp/go-secure-stdlib/httputil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
)

const (
	responseTypeCode     = "code"      // Authorization code flow
	responseTypeIDToken  = "id_token"  // ID Token for form post
	responseModeQuery    = "query"     // Response as a redirect with query parameters
	responseModeFormPost = "form_post" // Response as an HTML Form
)

func pathConfig(b *jwtAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixJWT,
		},
		Fields: map[string]*framework.FieldSchema{
			"oidc_discovery_url": {
				Type:        framework.TypeString,
				Description: `OIDC Discovery URL, without any .well-known component (base path). Cannot be used with "jwks_url" or "jwt_validation_pubkeys".`,
			},
			"oidc_discovery_ca_pem": {
				Type:        framework.TypeString,
				Description: "The CA certificate or chain of certificates, in PEM format, to use to validate connections to the OIDC Discovery URL. If not set, system certificates are used.",
			},
			"oidc_client_id": {
				Type:        framework.TypeString,
				Description: "The OAuth Client ID configured with your OIDC provider.",
			},
			"oidc_client_secret": {
				Type:        framework.TypeString,
				Description: "The OAuth Client Secret configured with your OIDC provider.",
				DisplayAttrs: &framework.DisplayAttributes{
					Sensitive: true,
				},
			},
			"oidc_response_mode": {
				Type:        framework.TypeString,
				Description: "The response mode to be used in the OAuth2 request. Allowed values are 'query' and 'form_post'.",
			},
			"oidc_response_types": {
				Type:        framework.TypeCommaStringSlice,
				Description: "The response types to request. Allowed values are 'code' and 'id_token'. Defaults to 'code'.",
			},
			"jwks_url": {
				Type:        framework.TypeString,
				Description: `JWKS URL to use to authenticate signatures. Cannot be used with "oidc_discovery_url" or "jwt_validation_pubkeys".`,
			},
			"jwks_ca_pem": {
				Type:        framework.TypeString,
				Description: "The CA certificate or chain of certificates, in PEM format, to use to validate connections to the JWKS URL. If not set, system certificates are used.",
			},
			"jwks_pairs": {
				Type:        framework.TypeSlice,
				Description: `Set of JWKS Url and CA certificate (or chain of certificates) pairs. CA certificates must be in PEM format. Cannot be used with "jwks_url" or "jwks_ca_pem".`,
			},
			"default_role": {
				Type:        framework.TypeLowerCaseString,
				Description: "The default role to use if none is provided during login. If not set, a role is required during login.",
			},
			"jwt_validation_pubkeys": {
				Type:        framework.TypeCommaStringSlice,
				Description: `A list of PEM-encoded public keys to use to authenticate signatures locally. Cannot be used with "jwks_url" or "oidc_discovery_url".`,
			},
			"jwt_supported_algs": {
				Type:        framework.TypeCommaStringSlice,
				Description: `A list of supported signing algorithms. Defaults to RS256.`,
			},
			"bound_issuer": {
				Type:        framework.TypeString,
				Description: "The value against which to match the 'iss' claim in a JWT. Optional.",
			},
			"provider_config": {
				Type:        framework.TypeMap,
				Description: "Provider-specific configuration. Optional.",
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Provider Config",
				},
			},
			"namespace_in_state": {
				Type:        framework.TypeBool,
				Description: "Pass namespace in the OIDC state parameter instead of as a separate query parameter. With this setting, the allowed redirect URL(s) in Vault and on the provider side should not contain a namespace query parameter. This means only one redirect URL entry needs to be maintained on the provider side for all vault namespaces that will be authenticating against it. Defaults to true for new configs.",
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Namespace in OIDC state",
					Value: true,
				},
			},
			"unsupported_critical_cert_extensions": {
				Type:        framework.TypeCommaStringSlice,
				Description: `A list of ASN1 OIDs of certificate extensions marked Critical that are unsupported by Vault and should be ignored.  This option should very rarely be needed except in specialized PKI environments.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				Summary:  "Read the current JWT authentication backend configuration.",
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "configuration",
				},
			},

			logical.UpdateOperation: &framework.PathOperation{
				Callback:    b.pathConfigWrite,
				Summary:     "Configure the JWT authentication backend.",
				Description: confHelpDesc,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
			},
		},

		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}
}

func (b *jwtAuthBackend) config(ctx context.Context, s logical.Storage) (*jwtConfig, error) {
	b.l.Lock()
	defer b.l.Unlock()

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

	config := &jwtConfig{}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}

	for _, v := range config.JWTValidationPubKeys {
		key, err := certutil.ParsePublicKeyPEM([]byte(v))
		if err != nil {
			return nil, errwrap.Wrapf("error parsing public key: {{err}}", err)
		}
		config.ParsedJWTPubKeys = append(config.ParsedJWTPubKeys, key)
	}

	b.cachedConfig = config

	return config, nil
}

func (b *jwtAuthBackend) pathConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	provider, err := NewProviderConfig(ctx, config, ProviderMap())
	if err != nil {
		return nil, err
	}

	// Omit sensitive keys from provider-specific config
	providerConfig := make(map[string]interface{})
	if provider != nil {
		for k, v := range config.ProviderConfig {
			providerConfig[k] = v
		}

		for _, k := range provider.SensitiveKeys() {
			delete(providerConfig, k)
		}
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"oidc_discovery_url":                   config.OIDCDiscoveryURL,
			"oidc_discovery_ca_pem":                config.OIDCDiscoveryCAPEM,
			"oidc_client_id":                       config.OIDCClientID,
			"oidc_response_mode":                   config.OIDCResponseMode,
			"oidc_response_types":                  config.OIDCResponseTypes,
			"default_role":                         config.DefaultRole,
			"jwt_validation_pubkeys":               config.JWTValidationPubKeys,
			"jwt_supported_algs":                   config.JWTSupportedAlgs,
			"jwks_url":                             config.JWKSURL,
			"jwks_pairs":                           config.JWKSPairs,
			"jwks_ca_pem":                          config.JWKSCAPEM,
			"bound_issuer":                         config.BoundIssuer,
			"provider_config":                      providerConfig,
			"namespace_in_state":                   config.NamespaceInState,
			"unsupported_critical_cert_extensions": config.UnsupportedCriticalCertExtensions,
		},
	}

	return resp, nil
}

func (b *jwtAuthBackend) pathConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config := &jwtConfig{
		OIDCDiscoveryURL:                  d.Get("oidc_discovery_url").(string),
		OIDCDiscoveryCAPEM:                d.Get("oidc_discovery_ca_pem").(string),
		OIDCClientID:                      d.Get("oidc_client_id").(string),
		OIDCClientSecret:                  d.Get("oidc_client_secret").(string),
		OIDCResponseMode:                  d.Get("oidc_response_mode").(string),
		OIDCResponseTypes:                 d.Get("oidc_response_types").([]string),
		JWKSURL:                           d.Get("jwks_url").(string),
		JWKSPairs:                         d.Get("jwks_pairs").([]interface{}),
		JWKSCAPEM:                         d.Get("jwks_ca_pem").(string),
		DefaultRole:                       d.Get("default_role").(string),
		JWTValidationPubKeys:              d.Get("jwt_validation_pubkeys").([]string),
		JWTSupportedAlgs:                  d.Get("jwt_supported_algs").([]string),
		BoundIssuer:                       d.Get("bound_issuer").(string),
		ProviderConfig:                    d.Get("provider_config").(map[string]interface{}),
		UnsupportedCriticalCertExtensions: d.Get("unsupported_critical_cert_extensions").([]string),
	}

	// Check if the config already exists, to determine if this is a create or
	// an update, since req.Operation is always 'update' in this handler, and
	// there's no existence check defined.
	existingConfig, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if nsInState, ok := d.GetOk("namespace_in_state"); ok {
		config.NamespaceInState = nsInState.(bool)
	} else if existingConfig == nil {
		// new configs default to true
		config.NamespaceInState = true
	} else {
		// maintain the existing value
		config.NamespaceInState = existingConfig.NamespaceInState
	}

	// Run checks on values
	methodCount := 0
	if config.OIDCDiscoveryURL != "" {
		methodCount++
	}
	if len(config.JWTValidationPubKeys) != 0 {
		methodCount++
	}
	if config.JWKSURL != "" {
		methodCount++
	}
	if len(config.JWKSPairs) > 0 {
		methodCount++
	}

	var jwksPairs []*JWKSPair
	switch {
	case methodCount != 1:
		return logical.ErrorResponse("exactly one of 'jwt_validation_pubkeys', 'jwks_url', 'jwks_pairs' or 'oidc_discovery_url' must be set"), nil

	case config.OIDCClientID != "" && config.OIDCClientSecret == "",
		config.OIDCClientID == "" && config.OIDCClientSecret != "":
		return logical.ErrorResponse("both 'oidc_client_id' and 'oidc_client_secret' must be set for OIDC"), nil

	case config.OIDCDiscoveryURL != "":
		var err error
		if config.OIDCClientID != "" && config.OIDCClientSecret != "" {
			_, err = b.createProvider(config)
		} else {
			_, err = jwt.NewOIDCDiscoveryKeySet(ctx, config.OIDCDiscoveryURL, config.OIDCDiscoveryCAPEM)
		}
		if err != nil {
			b.Logger().Error("error checking oidc discovery URL", "error", err)
			return logical.ErrorResponse("error checking oidc discovery URL"), nil
		}
	case config.OIDCClientID != "" && config.OIDCDiscoveryURL == "":
		return logical.ErrorResponse("'oidc_discovery_url' must be set for OIDC"), nil

	case config.JWKSCAPEM != "" && len(config.JWKSPairs) > 0:
		return logical.ErrorResponse("CA PEMs must be provided as part of the 'jwks_pairs' when using multiple JWKS URLs"), nil

	case len(config.JWKSPairs) > 0 && config.BoundIssuer != "":
		return logical.ErrorResponse("Bound issuer is not supported for use with 'jwks_pairs'"), nil

	case config.JWKSURL != "":
		if r := b.validateJWKSURL(ctx, config.JWKSURL, config.JWKSCAPEM); r != nil {
			return r, nil
		}

	case len(config.JWKSPairs) > 0:
		if jwksPairs, err = NewJWKSPairsConfig(config); err != nil {
			return logical.ErrorResponse("invalid jwks_pairs: %s", err), nil
		}

		for _, p := range jwksPairs {
			if r := b.validateJWKSURL(ctx, p.JWKSUrl, p.JWKSCAPEM); r != nil {
				return r, nil
			}
		}

	case len(config.JWTValidationPubKeys) != 0:
		for _, v := range config.JWTValidationPubKeys {
			if _, err := certutil.ParsePublicKeyPEM([]byte(v)); err != nil {
				return logical.ErrorResponse(fmt.Errorf("error parsing public key: %w", err).Error()), nil
			}
		}
	case len(config.UnsupportedCriticalCertExtensions) > 0:
		for _, v := range config.UnsupportedCriticalCertExtensions {
			if _, err := certutil.StringToOid(v); err != nil {
				return logical.ErrorResponse(fmt.Errorf("error parsing extension OID: %w", err).Error()), nil
			}
		}

	default:
		return nil, errors.New("unknown condition")
	}

	// NOTE: the OIDC lib states that if nothing is passed into its config, it
	// defaults to "RS256". So in the case of a zero value here it won't
	// default to e.g. "none".
	if err := jwt.SupportedSigningAlgorithm(toAlg(config.JWTSupportedAlgs)...); err != nil {
		return logical.ErrorResponse("invalid jwt_supported_algs: %s", err), nil
	}

	// Validate response_types
	if !strutil.StrListSubset([]string{responseTypeCode, responseTypeIDToken}, config.OIDCResponseTypes) {
		return logical.ErrorResponse("invalid response_types %v. 'code' and 'id_token' are allowed", config.OIDCResponseTypes), nil
	}

	// Validate response_mode
	switch config.OIDCResponseMode {
	case "", responseModeQuery:
		if config.hasType(responseTypeIDToken) {
			return logical.ErrorResponse("query response_mode may not be used with an id_token response_type"), nil
		}
	case responseModeFormPost:
	default:
		return logical.ErrorResponse("invalid response_mode: %q", config.OIDCResponseMode), nil
	}

	// Validate provider_config
	if _, err := NewProviderConfig(ctx, config, ProviderMap()); err != nil {
		return logical.ErrorResponse("invalid provider_config: %s", err), nil
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
	supportedSigAlgs := make([]oidc.Alg, len(config.JWTSupportedAlgs))
	for i, a := range config.JWTSupportedAlgs {
		supportedSigAlgs[i] = oidc.Alg(a)
	}

	if len(supportedSigAlgs) == 0 {
		supportedSigAlgs = []oidc.Alg{oidc.RS256}
	}

	var opts []oidc.Option
	if len(config.UnsupportedCriticalCertExtensions) > 0 {
		var oids []asn1.ObjectIdentifier
		for _, v := range config.UnsupportedCriticalCertExtensions {
			oid, err := certutil.StringToOid(v)
			if err != nil {
				return nil, errwrap.Wrapf("error creating provider: {{err}}", err)
			}
			oids = append(oids, oid)
		}

		var tp http.RoundTripper

		if config.OIDCDiscoveryCAPEM != "" {
			certPool := x509.NewCertPool()
			if ok := certPool.AppendCertsFromPEM([]byte(config.OIDCDiscoveryCAPEM)); ok {
				tp = &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs: certPool,
					},
				}
			}
		}
		ietripper, err := httputil.NewIgnoreUnhandledExtensionsRoundTripper(tp, oids)
		if err != nil {
			return nil, err
		}
		opts = append(opts, oidc.WithRoundTripper(ietripper))
	} else if config.OIDCDiscoveryCAPEM != "" {
		opts = append(opts, oidc.WithProviderCA(config.OIDCDiscoveryCAPEM))
	}
	c, err := oidc.NewConfig(config.OIDCDiscoveryURL, config.OIDCClientID,
		oidc.ClientSecret(config.OIDCClientSecret), supportedSigAlgs, []string{},
		opts...)
	if err != nil {
		return nil, errwrap.Wrapf("error creating provider: {{err}}", err)
	}

	provider, err := oidc.NewProvider(c)
	if err != nil {
		return nil, errwrap.Wrapf("error creating provider with given values: {{err}}", err)
	}

	return provider, nil
}

// createCAContext returns a context with custom TLS client, configured with the root certificates
// from caPEM. If no certificates are configured, the original context is returned.
func (b *jwtAuthBackend) createCAContext(ctx context.Context, caPEM string) (context.Context, error) {
	if caPEM == "" {
		return ctx, nil
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(caPEM)); !ok {
		return nil, errors.New("could not parse CA PEM value successfully")
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

	caCtx := context.WithValue(ctx, oauth2.HTTPClient, tc)

	return caCtx, nil
}

func (b *jwtAuthBackend) validateJWKSURL(ctx context.Context, JWKSURL, JWKSCAPEM string) *logical.Response {
	keyset, err := jwt.NewJSONWebKeySet(ctx, JWKSURL, JWKSCAPEM)
	if err != nil {
		b.Logger().Error("error checking jwks_ca_pem", "error", err)
		return logical.ErrorResponse("error checking jwks_ca_pem")
	}

	// Try to verify a correctly formatted JWT. The signature will fail to match, but other
	// errors with fetching the remote keyset should be reported.
	testJWT := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ"
	_, err = keyset.VerifySignature(ctx, testJWT)
	if err == nil {
		err = errors.New("unexpected verification of JWT")
	}

	if !strings.Contains(err.Error(), "failed to verify id token signature") {
		b.Logger().Error("error checking jwks URL", "url", JWKSURL, "error", err)
		return logical.ErrorResponse("error checking jwks URL %s", JWKSURL)
	}

	return nil
}

type jwtConfig struct {
	OIDCDiscoveryURL                  string                 `json:"oidc_discovery_url"`
	OIDCDiscoveryCAPEM                string                 `json:"oidc_discovery_ca_pem"`
	OIDCClientID                      string                 `json:"oidc_client_id"`
	OIDCClientSecret                  string                 `json:"oidc_client_secret"`
	OIDCResponseMode                  string                 `json:"oidc_response_mode"`
	OIDCResponseTypes                 []string               `json:"oidc_response_types"`
	JWKSURL                           string                 `json:"jwks_url"`
	JWKSCAPEM                         string                 `json:"jwks_ca_pem"`
	JWKSPairs                         []interface{}          `json:"jwks_pairs"`
	JWTValidationPubKeys              []string               `json:"jwt_validation_pubkeys"`
	JWTSupportedAlgs                  []string               `json:"jwt_supported_algs"`
	BoundIssuer                       string                 `json:"bound_issuer"`
	DefaultRole                       string                 `json:"default_role"`
	ProviderConfig                    map[string]interface{} `json:"provider_config"`
	NamespaceInState                  bool                   `json:"namespace_in_state"`
	UnsupportedCriticalCertExtensions []string               `json:"unsupported_critical_cert_extensions"`

	ParsedJWTPubKeys []crypto.PublicKey `json:"-"`
}

const (
	StaticKeys = iota
	JWKS
	MultiJWKS
	OIDCDiscovery
	OIDCFlow
	unconfigured
)

// authType classifies the authorization type/flow based on config parameters.
func (c jwtConfig) authType() int {
	switch {
	case len(c.ParsedJWTPubKeys) > 0:
		return StaticKeys
	case c.JWKSURL != "":
		return JWKS
	case len(c.JWKSPairs) > 0:
		return MultiJWKS
	case c.OIDCDiscoveryURL != "":
		if c.OIDCClientID != "" && c.OIDCClientSecret != "" {
			return OIDCFlow
		}
		return OIDCDiscovery
	}

	return unconfigured
}

// hasType returns whether the list of response types includes the requested
// type. The default type is 'code' so that special case is handled as well.
func (c jwtConfig) hasType(t string) bool {
	if len(c.OIDCResponseTypes) == 0 && t == responseTypeCode {
		return true
	}

	return strutil.StrListContains(c.OIDCResponseTypes, t)
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
