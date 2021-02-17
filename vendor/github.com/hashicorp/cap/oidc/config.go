package oidc

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/cap/oidc/internal/strutils"
)

// ClientSecret is an oauth client Secret.
type ClientSecret string

// RedactedClientSecret is the redacted string or json for an oauth client secret.
const RedactedClientSecret = "[REDACTED: client secret]"

// String will redact the client secret.
func (t ClientSecret) String() string {
	return RedactedClientSecret
}

// MarshalJSON will redact the client secret.
func (t ClientSecret) MarshalJSON() ([]byte, error) {
	return json.Marshal(RedactedClientSecret)
}

// Config represents the configuration for an OIDC provider used by a relying
// party.
type Config struct {
	// ClientID is the relying party ID.
	ClientID string

	// ClientSecret is the relying party secret.  This may be empty if you only
	// intend to use the provider with the authorization Code with PKCE or the
	// implicit flows.
	ClientSecret ClientSecret

	// Scopes is a list of default oidc scopes to request of the provider. The
	// required "oidc" scope is requested by default, and does not need to be
	// part of this optional list. If a Request has scopes, they will override
	// this configured list for a specific authentication attempt.
	Scopes []string

	// Issuer is a case-sensitive URL string using the https scheme that
	// contains scheme, host, and optionally, port number and path components
	// and no query or fragment components.
	//  See the Issuer Identifier spec: https://openid.net/specs/openid-connect-core-1_0.html#IssuerIdentifier
	//  See the OIDC connect discovery spec: https://openid.net/specs/openid-connect-discovery-1_0.html#IdentifierNormalization
	//  See the id_token spec: https://tools.ietf.org/html/rfc7519#section-4.1.1
	Issuer string

	// SupportedSigningAlgs is a list of supported signing algorithms. List of
	// currently supported algs: RS256, RS384, RS512, ES256, ES384, ES512,
	// PS256, PS384, PS512
	//
	// The list can be used to limit the supported algorithms when verifying
	// id_token signatures, an id_token's at_hash claim against an
	// access_token, etc.
	SupportedSigningAlgs []Alg

	// AllowedRedirectURLs is a list of allowed URLs for the provider to
	// redirect to after a user authenticates.  If AllowedRedirects is empty,
	// the package will not check the Request.RedirectURL() to see if it's
	// allowed, and the check will be left to the OIDC provider's /authorize
	// endpoint.
	AllowedRedirectURLs []string

	// Audiences is an optional default list of case-sensitive strings to use when
	// verifying an id_token's "aud" claim (which is also a list) If provided,
	// the audiences of an id_token must match one of the configured audiences.
	// If a Request has audiences, they will override this configured list for a
	// specific authentication attempt.
	Audiences []string

	// ProviderCA is an optional CA certs (PEM encoded) to use when sending
	// requests to the provider. If you have a list of *x509.Certificates, then
	// see EncodeCertificates(...) to PEM encode them.
	ProviderCA string

	// NowFunc is a time func that returns the current time.
	NowFunc func() time.Time
}

// NewConfig composes a new config for a provider.
//
// The "oidc" scope will always be added to the new configuration's Scopes,
// regardless of what additional scopes are requested via the WithScopes option
// and duplicate scopes are allowed.
//
// Supported options: WithProviderCA, WithScopes, WithAudiences, WithNow
func NewConfig(issuer string, clientID string, clientSecret ClientSecret, supported []Alg, allowedRedirectURLs []string, opt ...Option) (*Config, error) {
	const op = "NewConfig"
	opts := getConfigOpts(opt...)
	c := &Config{
		Issuer:               issuer,
		ClientID:             clientID,
		ClientSecret:         clientSecret,
		SupportedSigningAlgs: supported,
		Scopes:               opts.withScopes,
		ProviderCA:           opts.withProviderCA,
		Audiences:            opts.withAudiences,
		NowFunc:              opts.withNowFunc,
		AllowedRedirectURLs:  allowedRedirectURLs,
	}
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("%s: invalid provider config: %w", op, err)
	}
	return c, nil
}

// Validate the provider configuration.  Among other validations, it verifies
// the issuer is not empty, but it doesn't verify the Issuer is discoverable via
// an http request.  SupportedSigningAlgs are validated against the list of
// currently supported algs: RS256, RS384, RS512, ES256, ES384, ES512, PS256,
// PS384, PS512
func (c *Config) Validate() error {
	const op = "Config.Validate"

	// Note: c.ClientSecret is intentionally not checked for empty, in order to
	// support providers that only use the implicit flow or PKCE.
	if c == nil {
		return fmt.Errorf("%s: provider config is nil: %w", op, ErrNilParameter)
	}
	if c.ClientID == "" {
		return fmt.Errorf("%s: client ID is empty: %w", op, ErrInvalidParameter)
	}
	if c.Issuer == "" {
		return fmt.Errorf("%s: discovery URL is empty: %w", op, ErrInvalidParameter)
	}
	if len(c.AllowedRedirectURLs) > 0 {
		var invalidURLs []string
		for _, allowed := range c.AllowedRedirectURLs {
			if _, err := url.Parse(allowed); err != nil {
				invalidURLs = append(invalidURLs, allowed)
			}
		}
		if len(invalidURLs) > 0 {
			return fmt.Errorf("%s: Invalid AllowedRedirectURLs provided %s: %w", op, strings.Join(invalidURLs, ", "), ErrInvalidParameter)
		}
	}

	u, err := url.Parse(c.Issuer)
	if err != nil {
		return fmt.Errorf("%s: issuer %s is invalid (%s): %w", op, c.Issuer, err, ErrInvalidIssuer)
	}
	if !strutils.StrListContains([]string{"https", "http"}, u.Scheme) {
		return fmt.Errorf("%s: issuer %s schema is not http or https: %w", op, c.Issuer, ErrInvalidIssuer)
	}
	if len(c.SupportedSigningAlgs) == 0 {
		return fmt.Errorf("%s: supported algorithms is empty: %w", op, ErrInvalidParameter)
	}
	for _, a := range c.SupportedSigningAlgs {
		if !supportedAlgorithms[a] {
			return fmt.Errorf("%s: unsupported algorithm %s: %w", op, a, ErrInvalidParameter)
		}
	}
	if c.ProviderCA != "" {
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM([]byte(c.ProviderCA)); !ok {
			return fmt.Errorf("%s: %w", op, ErrInvalidCACert)
		}
	}
	return nil
}

// Now will return the current time which can be overridden by the NowFunc
func (c *Config) Now() time.Time {
	if c.NowFunc != nil {
		return c.NowFunc()
	}
	return time.Now() // fallback to this default
}

// configOptions is the set of available options
type configOptions struct {
	withScopes     []string
	withAudiences  []string
	withProviderCA string
	withNowFunc    func() time.Time
}

// configDefaults is a handy way to get the defaults at runtime and
// during unit tests.
func configDefaults() configOptions {
	return configOptions{
		withScopes: []string{oidc.ScopeOpenID},
	}
}

// getConfigOpts gets the defaults and applies the opt overrides passed
// in.
func getConfigOpts(opt ...Option) configOptions {
	opts := configDefaults()
	ApplyOpts(&opts, opt...)
	return opts
}

// WithProviderCA provides optional CA certs (PEM encoded) for the provider's
// config.  These certs will can be used when making http requests to the
// provider.
//
// Valid for: Config
//
// See EncodeCertificates(...) to PEM encode a number of certs.
func WithProviderCA(cert string) Option {
	return func(o interface{}) {
		if o, ok := o.(*configOptions); ok {
			o.withProviderCA = cert
		}
	}
}

// EncodeCertificates will encode a number of x509 certificates to PEM.  It will
// help encode certs for use with the WithProviderCA(...) option.
func EncodeCertificates(certs ...*x509.Certificate) (string, error) {
	const op = "EncodeCert"
	var buffer bytes.Buffer
	if len(certs) == 0 {
		return "", fmt.Errorf("%s: no certs provided: %w", op, ErrInvalidParameter)
	}
	for _, cert := range certs {
		if cert == nil {
			return "", fmt.Errorf("%s: empty cert: %w", op, ErrNilParameter)
		}
		if err := pem.Encode(&buffer, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		}); err != nil {
			return "", fmt.Errorf("%s: unable to encode cert: %w", op, err)
		}
	}
	return buffer.String(), nil
}
