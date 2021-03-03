package oidc

import (
	"fmt"
	"time"

	"golang.org/x/text/language"
)

// Request basically represents one OIDC authentication flow for a user. It
// contains the data needed to uniquely represent that one-time flow across the
// multiple interactions needed to complete the OIDC flow the user is
// attempting.
//
// Request() is passed throughout the OIDC interactions to uniquely identify the
// flow's request. The Request.State() and Request.Nonce() cannot be equal, and
// will be used during the OIDC flow to prevent CSRF and replay attacks (see the
// oidc spec for specifics).
//
// Audiences and Scopes are optional overrides of configured provider defaults
// for specific authentication attempts
type Request interface {
	// State is a unique identifier and an opaque value used to maintain request
	// between the oidc request and the callback. State cannot equal the Nonce.
	// See https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest.
	State() string

	// Nonce is a unique nonce and a string value used to associate a Client
	// session with an ID Token, and to mitigate replay attacks. Nonce cannot
	// equal the ID.
	// See https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	// and https://openid.net/specs/openid-connect-core-1_0.html#NonceNotes.
	Nonce() string

	// IsExpired returns true if the request has expired. Implementations should
	// support a time skew (perhaps RequestExpirySkew) when checking expiration.
	IsExpired() bool

	// Audiences is an specific authentication attempt's list of optional
	// case-sensitive strings to use when verifying an id_token's "aud" claim
	// (which is also a list). If provided, the audiences of an id_token must
	// match one of the configured audiences.  If a Request does not have
	// audiences, then the configured list of default audiences will be used.
	Audiences() []string

	// Scopes is a specific authentication attempt's list of optional
	// scopes to request of the provider. The required "oidc" scope is requested
	// by default, and does not need to be part of this optional list. If a
	// Request does not have Scopes, then the configured list of default
	// requested scopes will be used.
	Scopes() []string

	// RedirectURL is a URL where providers will redirect responses to
	// authentication requests.
	RedirectURL() string

	// ImplicitFlow indicates whether or not to use the implicit flow with form
	// post. Getting only an id_token for an implicit flow should be the
	// default for implementations, but at times it's necessary to also request
	// an access_token, so this function and the WithImplicitFlow(...) option
	// allows for those scenarios. Overall, it is recommend to not request
	// access_tokens during the implicit flow.  If you need an access_token,
	// then use the authorization code flows and if you can't secure a client
	// secret then use the authorization code flow with PKCE.
	//
	// The first returned bool represents if the implicit flow has been requested.
	// The second returned bool represents if an access token has been requested
	// during the implicit flow.
	//
	// See: https://openid.net/specs/openid-connect-core-1_0.html#ImplicitFlowAuth
	// See: https://openid.net/specs/oauth-v2-form-post-response-mode-1_0.html
	ImplicitFlow() (useImplicitFlow bool, includeAccessToken bool)

	// PKCEVerifier indicates whether or not to use the authorization code flow
	// with PKCE.  PKCE should be used for any client which cannot secure a
	// client secret (SPA and native apps) or is susceptible to authorization
	// code intercept attacks. When supported by your OIDC provider, PKCE should
	// be used instead of the implicit flow.
	//
	// See: https://tools.ietf.org/html/rfc7636
	PKCEVerifier() CodeVerifier

	// MaxAge: when authAfter is not a zero value (authTime.IsZero()) then the
	// id_token's auth_time claim must be after the specified time.
	//
	// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	MaxAge() (seconds uint, authAfter time.Time)

	// Prompts optionally defines a list of values that specifies whether the
	// Authorization Server prompts the End-User for reauthentication and
	// consent.  See MaxAge() if wish to specify an allowable elapsed time in
	// seconds since the last time the End-User was actively authenticated by
	// the OP.
	//
	// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	Prompts() []Prompt

	// Display optionally specifies how the Authorization Server displays the
	// authentication and consent user interface pages to the End-User.
	//
	// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	Display() Display

	// UILocales optionally specifies End-User's preferred languages via
	// language Tags, ordered by preference.
	//
	// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	UILocales() []language.Tag

	// Claims optionally requests that specific claims be returned using
	// the claims parameter.
	//
	// https://openid.net/specs/openid-connect-core-1_0.html#ClaimsParameter
	Claims() []byte

	// ACRValues() optionally specifies the acr values that the Authorization
	// Server is being requested to use for processing this Authentication
	// Request, with the values appearing in order of preference.
	//
	// NOTE: Requested acr_values are not verified by the Provider.Exchange(...)
	// or Provider.VerifyIDToken() functions, since the request/return values
	// are determined by the provider's implementation. You'll need to verify
	// the claims returned yourself based on values provided by you OIDC
	// Provider's documentation.
	//
	// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	ACRValues() []string
}

// Req represents the oidc request used for oidc flows and implements the Request interface.
type Req struct {
	//	state is a unique identifier and an opaque value used to maintain request
	//	between the oidc request and the callback.
	state string

	// nonce is a unique nonce and suitable for use as an oidc nonce.
	nonce string

	// Expiration is the expiration time for the Request.
	expiration time.Time

	// redirectURL is a URL where providers will redirect responses to
	// authentication requests.
	redirectURL string

	// scopes is a specific authentication attempt's list of optional
	// scopes to request of the provider. The required "oidc" scope is requested
	// by default, and does not need to be part of this optional list. If a
	// Request does not have Scopes, then the configured list of default
	// requested scopes will be used.
	scopes []string

	// audiences is an specific authentication attempt's list of optional
	// case-sensitive strings to use when verifying an id_token's "aud" claim
	// (which is also a list). If provided, the audiences of an id_token must
	// match one of the configured audiences.  If a Request does not have
	// audiences, then the configured list of default audiences will be used.
	audiences []string

	// nowFunc is an optional function that returns the current time
	nowFunc func() time.Time

	// withImplicit indicates whether or not to use the implicit flow.  Getting
	// only an id_token for an implicit flow is the default. If an access_token
	// is also required, then withImplicit.withAccessToken will be true. It
	// is recommend to not request access_tokens during the implicit flow.  If
	// you need an access_token, then use the authorization code flows (with
	// optional PKCE).
	withImplicit *implicitFlow

	// withVerifier indicates whether or not to use the authorization code flow
	// with PKCE.  It suppies the required CodeVerifier for PKCE.
	withVerifier CodeVerifier

	// withMaxAge: when withMaxAge.authAfter is not a zero value
	// (authTime.IsZero()) then the id_token's auth_time claim must be after the
	// specified time.
	withMaxAge *maxAge

	// withPrompts optionally defines a list of values that specifies whether
	// the Authorization Server prompts the End-User for reauthentication and
	// consent.
	withPrompts []Prompt

	// withDisplay optionally specifies how the Authorization Server displays the
	// authentication and consent user interface pages to the End-User.
	withDisplay Display

	// withUILocales optionally specifies End-User's preferred languages via
	// language Tags, ordered by preference.
	withUILocales []language.Tag

	// withClaims optionally requests that specific claims be returned
	// using the claims parameter.
	withClaims []byte

	// withACRValues() optionally specifies the acr values that the Authorization
	// Server is being requested to use for processing this Authentication
	// Request, with the values appearing in order of preference.
	withACRValues []string
}

// ensure that Request implements the Request interface.
var _ Request = (*Req)(nil)

// NewRequest creates a new Request (*Req).
//  Supports the options:
//   * WithState
//   * WithNow
//   * WithAudiences
//   * WithScopes
//   * WithImplicit
//   * WithPKCE
//   * WithMaxAge
//   * WithPrompts
//   * WithDisplay
//   * WithUILocales
//   * WithClaims
func NewRequest(expireIn time.Duration, redirectURL string, opt ...Option) (*Req, error) {
	const op = "oidc.NewRequest"
	opts := getReqOpts(opt...)
	if redirectURL == "" {
		return nil, fmt.Errorf("%s: redirect URL is empty: %w", op, ErrInvalidParameter)
	}
	nonce, err := NewID(WithPrefix("n"))
	if err != nil {
		return nil, fmt.Errorf("%s: unable to generate a request's nonce: %w", op, err)
	}

	var state string
	switch {
	case opts.withState != "":
		state = opts.withState
	default:
		var err error
		state, err = NewID(WithPrefix("st"))
		if err != nil {
			return nil, fmt.Errorf("%s: unable to generate a request's state: %w", op, err)
		}
	}

	if expireIn == 0 || expireIn < 0 {
		return nil, fmt.Errorf("%s: expireIn not greater than zero: %w", op, ErrInvalidParameter)
	}
	if opts.withVerifier != nil && opts.withImplicitFlow != nil {
		return nil, fmt.Errorf("%s: requested both implicit flow and authorization code with PKCE: %w", op, ErrInvalidParameter)
	}
	r := &Req{
		state:         state,
		nonce:         nonce,
		redirectURL:   redirectURL,
		nowFunc:       opts.withNowFunc,
		audiences:     opts.withAudiences,
		scopes:        opts.withScopes,
		withImplicit:  opts.withImplicitFlow,
		withVerifier:  opts.withVerifier,
		withPrompts:   opts.withPrompts,
		withDisplay:   opts.withDisplay,
		withUILocales: opts.withUILocales,
		withClaims:    opts.withClaims,
		withACRValues: opts.withACRValues,
	}
	r.expiration = r.now().Add(expireIn)
	if opts.withMaxAge != nil {
		opts.withMaxAge.authAfter = r.now().Add(time.Duration(-opts.withMaxAge.seconds) * time.Second)
		r.withMaxAge = opts.withMaxAge
	}
	return r, nil
}

// State implements the Request.State() interface function.
func (r *Req) State() string { return r.state }

// Nonce implements the Request.Nonce() interface function.
func (r *Req) Nonce() string { return r.nonce }

// Audiences implements the Request.Audiences() interface function and returns a
// copy of the audiences.
func (r *Req) Audiences() []string {
	if r.audiences == nil {
		return nil
	}
	cp := make([]string, len(r.audiences))
	copy(cp, r.audiences)
	return cp
}

// Scopes implements the Request.Scopes() interface function and returns a copy of
// the scopes.
func (r *Req) Scopes() []string {
	if r.scopes == nil {
		return nil
	}
	cp := make([]string, len(r.scopes))
	copy(cp, r.scopes)
	return cp
}

// RedirectURL implements the Request.RedirectURL() interface function.
func (r *Req) RedirectURL() string { return r.redirectURL }

// PKCEVerifier implements the Request.PKCEVerifier() interface function and
// returns a copy of the CodeVerifier
func (r *Req) PKCEVerifier() CodeVerifier {
	if r.withVerifier == nil {
		return nil
	}
	return r.withVerifier.Copy()
}

// Prompts() implements the Request.Prompts() interface function and returns a
// copy of the prompts.
func (r *Req) Prompts() []Prompt {
	if r.withPrompts == nil {
		return nil
	}
	cp := make([]Prompt, len(r.withPrompts))
	copy(cp, r.withPrompts)
	return cp
}

// Display() implements the Request.Display() interface function.
func (r *Req) Display() Display { return r.withDisplay }

// UILocales() implements the Request.UILocales() interface function and returns a
// copy of the UILocales
func (r *Req) UILocales() []language.Tag {
	if r.withUILocales == nil {
		return nil
	}
	cp := make([]language.Tag, len(r.withUILocales))
	copy(cp, r.withUILocales)
	return cp
}

// Claims() implements the Request.Claims() interface function
// and returns a copy of the claims request.
func (r *Req) Claims() []byte {
	if r.withClaims == nil {
		return nil
	}
	cp := make([]byte, len(r.withClaims))
	copy(cp, r.withClaims)
	return cp
}

// ACRValues() implements the Request.ARCValues() interface function and returns a
// copy of the acr values
func (r *Req) ACRValues() []string {
	if len(r.withACRValues) == 0 {
		return nil
	}
	cp := make([]string, len(r.withACRValues))
	copy(cp, r.withACRValues)
	return cp
}

// MaxAge: when authAfter is not a zero value (authTime.IsZero()) then the
// id_token's auth_time claim must be after the specified time.
//
// See: https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
func (r *Req) MaxAge() (uint, time.Time) {
	if r.withMaxAge == nil {
		return 0, time.Time{}
	}
	return r.withMaxAge.seconds, r.withMaxAge.authAfter.Truncate(time.Second)
}

// ImplicitFlow indicates whether or not to use the implicit flow.  Getting
// only an id_token for an implicit flow is the default, but at times
// it's necessary to also request an access_token, so this function and the
// WithImplicitFlow(...) option allows for those scenarios. Overall, it is
// recommend to not request access_tokens during the implicit flow.  If you need
// an access_token, then use the authorization code flows and if you can't
// secure a client secret then use the authorization code flow with PKCE.
//
// The first returned bool represents if the implicit flow has been requested.
// The second returned bool represents if an access token has been requested
// during the implicit flow.
func (r *Req) ImplicitFlow() (bool, bool) {
	if r.withImplicit == nil {
		return false, false
	}
	switch {
	case r.withImplicit.withAccessToken:
		return true, true
	default:
		return true, false
	}
}

// RequestExpirySkew defines a time skew when checking a Request's expiration.
const RequestExpirySkew = 1 * time.Second

// IsExpired returns true if the request has expired.
func (r *Req) IsExpired() bool {
	return r.expiration.Before(time.Now().Add(RequestExpirySkew))
}

// now returns the current time using the optional timeFn
func (r *Req) now() time.Time {
	if r.nowFunc != nil {
		return r.nowFunc()
	}
	return time.Now() // fallback to this default
}

type implicitFlow struct {
	withAccessToken bool
}

type maxAge struct {
	seconds   uint
	authAfter time.Time
}

// reqOptions is the set of available options for Req functions
type reqOptions struct {
	withNowFunc      func() time.Time
	withScopes       []string
	withAudiences    []string
	withImplicitFlow *implicitFlow
	withVerifier     CodeVerifier
	withMaxAge       *maxAge
	withPrompts      []Prompt
	withDisplay      Display
	withUILocales    []language.Tag
	withClaims       []byte
	withACRValues    []string
	withState        string
}

// reqDefaults is a handy way to get the defaults at runtime and during unit
// tests.
func reqDefaults() reqOptions {
	return reqOptions{}
}

// getReqOpts gets the request defaults and applies the opt overrides passed in
func getReqOpts(opt ...Option) reqOptions {
	opts := reqDefaults()
	ApplyOpts(&opts, opt...)
	return opts
}

// WithImplicitFlow provides an option to use an OIDC implicit flow with form
// post. It should be noted that if your OIDC provider supports PKCE, then use
// it over the implicit flow.  Getting only an id_token is the default, and
// optionally passing a true bool will request an access_token as well during
// the flow.  You cannot use WithImplicit and WithPKCE together.  It is
// recommend to not request access_tokens during the implicit flow.  If you need
// an access_token, then use the authorization code flows.
//
// Option is valid for: Request
//
// See: https://openid.net/specs/openid-connect-core-1_0.html#ImplicitFlowAuth
// See: https://openid.net/specs/oauth-v2-form-post-response-mode-1_0.html
func WithImplicitFlow(args ...interface{}) Option {
	withAccessToken := false
	for _, arg := range args {
		switch arg := arg.(type) {
		case bool:
			if arg {
				withAccessToken = true
			}
		}
	}
	return func(o interface{}) {
		if o, ok := o.(*reqOptions); ok {
			o.withImplicitFlow = &implicitFlow{
				withAccessToken: withAccessToken,
			}
		}
	}
}

// WithPKCE provides an option to use a CodeVerifier with the authorization
// code flow with PKCE.  You cannot use WithImplicit and WithPKCE together.
//
// Option is valid for: Request
//
// See: https://tools.ietf.org/html/rfc7636
func WithPKCE(v CodeVerifier) Option {
	return func(o interface{}) {
		if o, ok := o.(*reqOptions); ok {
			o.withVerifier = v
		}
	}
}

// WithMaxAge provides an optional maximum authentication age, which is the
// allowable elapsed time in seconds since the last time the user was actively
// authenticated by the provider.  When a max age is specified, the provider
// must include a auth_time claim in the returned id_token.  This makes it
// preferable to prompt=login, where you have no way to verify when an
// authentication took place.
//
// Option is valid for: Request
//
// See: https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
func WithMaxAge(seconds uint) Option {
	return func(o interface{}) {
		if o, ok := o.(*reqOptions); ok {
			// authAfter will be a zero value, since it's not set until the
			// NewRequest() factory, when it can determine it's nowFunc
			o.withMaxAge = &maxAge{
				seconds: seconds,
			}
		}
	}
}

// WithPrompts provides an optional list of values that specifies whether the
// Authorization Server prompts the End-User for reauthentication and consent.
//
// See MaxAge() if wish to specify an allowable elapsed time in seconds since
// the last time the End-User was actively authenticated by the OP.
//
// Option is valid for: Request
//
// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
func WithPrompts(prompts ...Prompt) Option {
	return func(o interface{}) {
		if o, ok := o.(*reqOptions); ok {
			o.withPrompts = prompts
		}
	}
}

// WithDisplay optionally specifies how the Authorization Server displays the
// authentication and consent user interface pages to the End-User.
//
// Option is valid for: Request
//
// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
func WithDisplay(d Display) Option {
	return func(o interface{}) {
		if o, ok := o.(*reqOptions); ok {
			o.withDisplay = d
		}
	}
}

// WithUILocales optionally specifies End-User's preferred languages via
// language Tags, ordered by preference.
//
// Option is valid for: Request
//
// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
func WithUILocales(locales ...language.Tag) Option {
	return func(o interface{}) {
		if o, ok := o.(*reqOptions); ok {
			o.withUILocales = locales
		}
	}
}

// WithClaims optionally requests that specific claims be returned using
// the claims parameter.
//
// Option is valid for: Request
//
// https://openid.net/specs/openid-connect-core-1_0.html#ClaimsParameter
func WithClaims(json []byte) Option {
	return func(o interface{}) {
		if o, ok := o.(*reqOptions); ok {
			o.withClaims = json
		}
	}
}

// WithACRValues optionally specifies the acr values that the Authorization
// Server is being requested to use for processing this Authentication
// Request, with the values appearing in order of preference.
//
// NOTE: Requested acr_values are not verified by the Provider.Exchange(...)
// or Provider.VerifyIDToken() functions, since the request/return values
// are determined by the provider's implementation. You'll need to verify
// the claims returned yourself based on values provided by you OIDC
// Provider's documentation.
//
// Option is valid for: Request
//
// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
func WithACRValues(values ...string) Option {
	return func(o interface{}) {
		if o, ok := o.(*reqOptions); ok {
			o.withACRValues = values
		}
	}
}

// WithState optionally specifies a value to use for the request's state.
// Typically, state is a random string generated for you when you create
// a new Request. This option allows you to override that auto-generated value
// with a specific value of your own choosing.
//
// The primary reason for using the state parameter is to mitigate CSRF attacks
// by using a unique and non-guessable value associated with each authentication
// request about to be initiated. That value allows you to prevent the attack by
// confirming that the value coming from the response matches the one you sent.
// Since the state parameter is a string, you can encode any other information
// in it.
//
// Some care must be taken to not use a state which is longer than your OIDC
// Provider allows.  The specification places no limit on the length, but there
// are many practical limitations placed on the length by browsers, proxies and
// of course your OIDC provider.
//
// State should be at least 20 chars long (see:
// https://tools.ietf.org/html/rfc6749#section-10.10).
//
// See NewID(...) for a function that generates a sufficiently
// random string and supports the WithPrefix(...) option, which can be used
// prefix your custom state payload.
//
// Neither a max or min length is enforced when you use the WithState option.
//
// Option is valid for: Request
//
func WithState(s string) Option {
	return func(o interface{}) {
		if o, ok := o.(*reqOptions); ok {
			o.withState = s
		}
	}
}
