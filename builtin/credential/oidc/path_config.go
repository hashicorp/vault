package oidc

import (
	"fmt"
	"net/url"

	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/coreos/go-oidc"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"net/http"
)

const (
	defaultUsernameClaim = "sub"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		Fields: map[string]*framework.FieldSchema{
			"issuer_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "OIDC issuer URL. This is typically the base URL (no path) to the discovery URL of the provider.",
			},
			"issuer_verify_ca": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "OIDC issuer verification CA. If set, it is a PEM-encoded CA bundle to verify the TLS connections to the OIDC provider. If not provided, system certificates are used.",
			},
			"client_ids": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "OIDC client IDs that are permittable in identity token's `aud` claims.",
			},
			"username_claim": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The OIDC claim inside the identity token that identifies the username. Defaults to `sub`.",
				Default:     defaultUsernameClaim,
			},
			"groups_claim": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The OIDC claim inside the identity token that identifies the groups the token is for. This claim must be a list of strings. If unset, no groups mapping will be used.",
				Default:     "",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.CreateOperation: b.pathConfigWrite,
			logical.UpdateOperation: b.pathConfigWrite,
		},

		ExistenceCheck: b.pathConfigExistenceCheck,

		HelpSynopsis: pathConfigHelp,
	}
}

// Config returns the configuration for this backend.
func (b *backend) Config(s logical.Storage) (*ConfigEntry, error) {
	entry, err := s.Get("config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result ConfigEntry
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (b *backend) pathConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	cfg, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"issuer_url":       cfg.IssuerUrl,
			"issuer_verify_ca": cfg.IssuerCABundle,
			"client_ids":       cfg.ClientIDs,
			"username_claim":   cfg.UsernameClaim,
			"groups_claim":     cfg.GroupsClaim,
		},
	}

	return resp, nil
}

func (b *backend) pathConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	// Due to the existence check, entry will only be nil if it's a create
	// operation, so just create a new one
	if cfg == nil {
		cfg = &ConfigEntry{}
	}

	// Parse issuer_url, required.
	issuerUrlRaw, ok, err := d.GetOkErr("issuer_url")
	if ok {
		_, err = url.Parse(issuerUrlRaw.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error parsing given issuer_url: %s", err)), nil
		}
		cfg.IssuerUrl = issuerUrlRaw.(string)
	} else if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	} else {
		return logical.ErrorResponse("auth/oidc: missing required config value issuer_url"), nil
	}
	// Parse issuer_verify_ca.
	issuerVerifyCaRaw, ok, err := d.GetOkErr("issuer_verify_ca")
	if ok {
		_, err := certPoolFromString(issuerVerifyCaRaw.(string))
		if err != nil {
			return logical.ErrorResponse("auth/oidc: failed parsing PEM of issuer_verify_ca: " + err.Error()), nil
		}
		cfg.IssuerCABundle = issuerVerifyCaRaw.(string)
	} else if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	// Parse client ids, required.
	clientIdsRaw, ok, err := d.GetOkErr("client_ids")
	if ok {
		cfg.ClientIDs = clientIdsRaw.([]string)
	} else if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	} else if !ok {
		return logical.ErrorResponse("auth/oidc: missing required config value client_ids"), nil
	}
	// Read defaults for claim configs.
	cfg.UsernameClaim = d.Get("username_claim").(string)
	cfg.GroupsClaim = d.Get("groups_claim").(string)

	jsonCfg, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if err := req.Storage.Put(jsonCfg); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return nil, nil
}

func (b *backend) pathConfigExistenceCheck(
	req *logical.Request, d *framework.FieldData) (bool, error) {
	cfg, err := b.Config(req.Storage)
	if err != nil {
		return false, err
	}

	return cfg != nil, nil
}

func (b *backend) oidcProviderForConfig(conf *ConfigEntry) (*oidc.Provider, error) {
	// TODO(mwitkow): This creates a new Provider each time, making a request to the Discovery URL of the provider. Add caching per config "hash".
	ctx := context.TODO()
	if conf.IssuerCABundle != "" {
		bundles, err := certPoolFromString(conf.IssuerCABundle)
		if err != nil {
			return nil, err
		}
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: bundles,
			},
		}
		httpClient := &http.Client{
			Transport: transport,
		}
		ctx = oidc.ClientContext(ctx, httpClient)
	}

	return oidc.NewProvider(ctx, conf.IssuerUrl)
}

func (b *backend) validateAndExtractClaims(cfg *ConfigEntry, rawToken string) (userName string, groups []string, token *oidc.IDToken, err error) {
	provider, err := b.oidcProviderForConfig(cfg)
	if err != nil {
		return "", []string{}, nil, err
	}
	if provider == nil {
		return "", []string{}, nil, errors.New("OIDC provider not configured")
	}
	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true}) // note we verify the audience ourselves.
	idToken, err := verifier.Verify(context.TODO(), rawToken)
	if err != nil {
		return "", []string{}, nil, errors.New("OIDC identity token error: " + err.Error())
	}
	if !containsAny(idToken.Audience, cfg.ClientIDs) {
		b.Logger().Info("OIDC token has bad 'aud': %v", idToken.Audience)
		return "", []string{}, nil, errors.New("OIDC identity token issued for unsupported client ID.")
	}
	claims := make(map[string]interface{})
	if err := idToken.Claims(&claims); err != nil {
		b.Logger().Warn("OIDC token claims parsing error: %v", err)
		return "", []string{}, nil, errors.New("OIDC identity token claim parsing error.")
	}
	b.Logger().Debug("OIDC claims %v", claims) // TODO(mwitkow): Remove.

	userNameClaim, ok := claims[cfg.UsernameClaim].(string)
	if !ok {
		b.Logger().Warn("OIDC token doesn't have username under expected claim. Claims: %v", claims)
		return "", []string{}, nil, errors.New("OIDC identity token user claim parsing error.")
	}
	userName = userNameClaim
	if cfg.GroupsClaim != "" {
		groupNameClaim, ok := claims[cfg.GroupsClaim].([]interface{})
		if !ok {
			b.Logger().Warn("OIDC token doesn't have group name under expected claim.", "claimName", cfg.GroupsClaim, "allClaims", claims)
			return "", []string{}, nil, errors.New("OIDC identity token group claim parsing error.")
		}
		for _, gInt := range groupNameClaim {
			g, ok := gInt.(string)
			if ok {
				groups = append(groups, g)
			}
		}
	}
	return userName, groups, idToken, nil
}

func containsAny(existingValues []string, checkedValues []string) bool {
	for _, s := range existingValues {
		for _, c := range checkedValues {
			if s == c {
				return true
			}
		}
	}
	return false
}

func certPoolFromString(caBundle string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM([]byte(caBundle))
	if !ok {
		return nil, errors.New("failed parsing CA bundle")
	}
	return pool, nil
}

// ConfigEntry for OIDC configuration.
type ConfigEntry struct {
	IssuerUrl      string   `json:"issuer_url"`
	IssuerCABundle string   `json:"issuer_verify_ca"`
	ClientIDs      []string `json:"client_ids"`
	UsernameClaim  string   `json:"username_claim"`
	GroupsClaim    string   `json:"group_claim"`
}

// TODO(mwitkow): Update this.
const pathConfigHelp = `
This endpoint allows you to configure the OpenID Connect provider.
`
