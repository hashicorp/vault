package oidc

import (
	"fmt"

	"github.com/hashicorp/cap/oidc/internal/base62"
)

// DefaultIDLength is the default length for generated IDs, which are used for
// state and nonce parameters during OIDC flows.
//
// For ID length requirements see:
// https://tools.ietf.org/html/rfc6749#section-10.10
const DefaultIDLength = 20

// NewID generates a ID with an optional prefix.   The ID generated is suitable
// for a Request's State or Nonce. The ID length will be DefaultIDLen, unless an
// optional prefix is provided which will add the prefix's length + an
// underscore.  The WithPrefix, WithLen options are supported.
//
// For ID length requirements see:
// https://tools.ietf.org/html/rfc6749#section-10.10
func NewID(opt ...Option) (string, error) {
	const op = "NewID"
	opts := getIDOpts(opt...)
	id, err := base62.Random(opts.withLen)
	if err != nil {
		return "", fmt.Errorf("%s: unable to generate id: %w", op, err)
	}
	switch {
	case opts.withPrefix != "":
		return fmt.Sprintf("%s_%s", opts.withPrefix, id), nil
	default:
		return id, nil
	}
}

// idOptions is the set of available options.
type idOptions struct {
	withPrefix string
	withLen    int
}

// idDefaults is a handy way to get the defaults at runtime and
// during unit tests.
func idDefaults() idOptions {
	return idOptions{
		withLen: DefaultIDLength,
	}
}

// getConfigOpts gets the defaults and applies the opt overrides passed
// in.
func getIDOpts(opt ...Option) idOptions {
	opts := idDefaults()
	ApplyOpts(&opts, opt...)
	return opts
}

// WithPrefix provides an optional prefix for an new ID.  When this options is
// provided, NewID will prepend the prefix and an underscore to the new
// identifier.
//
// Valid for: ID
func WithPrefix(prefix string) Option {
	return func(o interface{}) {
		if o, ok := o.(*idOptions); ok {
			o.withPrefix = prefix
		}
	}
}
