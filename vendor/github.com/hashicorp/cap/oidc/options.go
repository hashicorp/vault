package oidc

import (
	"time"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/cap/oidc/internal/strutils"
)

// Option defines a common functional options type which can be used in a
// variadic parameter pattern.
type Option func(interface{})

// ApplyOpts takes a pointer to the options struct as a set of default options
// and applies the slice of opts as overrides.
func ApplyOpts(opts interface{}, opt ...Option) {
	for _, o := range opt {
		if o == nil { // ignore any nil Options
			continue
		}
		o(opts)
	}
}

// WithNow provides an optional func for determining what the current time it
// is.
//
// Valid for: Config, Tk and Request
func WithNow(now func() time.Time) Option {
	return func(o interface{}) {
		if now == nil {
			return
		}
		switch v := o.(type) {
		case *configOptions:
			v.withNowFunc = now
		case *tokenOptions:
			v.withNowFunc = now
		case *reqOptions:
			v.withNowFunc = now
		}
	}
}

// WithScopes provides an optional list of scopes.
//
// Valid for: Config and Request
func WithScopes(scopes ...string) Option {
	return func(o interface{}) {
		if len(scopes) == 0 {
			return
		}
		switch v := o.(type) {
		case *configOptions:
			// configOptions already has the oidc.ScopeOpenID in its defaults.
			scopes = strutils.RemoveDuplicatesStable(scopes, false)
			v.withScopes = append(v.withScopes, scopes...)
		case *reqOptions:
			// need to prepend the oidc.ScopeOpenID
			ts := append([]string{oidc.ScopeOpenID}, scopes...)
			scopes = strutils.RemoveDuplicatesStable(ts, false)
			v.withScopes = append(v.withScopes, scopes...)
		}
	}
}

// WithAudiences provides an optional list of audiences.
//
//Valid for: Config and Request
func WithAudiences(auds ...string) Option {
	return func(o interface{}) {
		if len(auds) == 0 {
			return
		}
		auds := strutils.RemoveDuplicatesStable(auds, false)
		switch v := o.(type) {
		case *configOptions:
			v.withAudiences = append(v.withAudiences, auds...)
		case *reqOptions:
			v.withAudiences = append(v.withAudiences, auds...)
		case *userInfoOptions:
			v.withAudiences = append(v.withAudiences, auds...)
		}
	}
}
