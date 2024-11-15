// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package wrapping

import (
	"errors"
)

// GetOpts iterates the inbound Options and returns a struct
func GetOpts(opt ...Option) (*Options, error) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o == nil {
			continue
		}
		iface := o()
		switch to := iface.(type) {
		case OptionFunc:
			if err := to(opts); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("option passed into top-level wrapping options handler" +
				" that is not from this package; this is likely due to the wrapper being" +
				" invoked as a plugin but options being sent from a specific wrapper package;" +
				" use WithConfigMap to send options via the plugin interface")
		}
	}
	return opts, nil
}

// Option - a type that wraps an interface for compile-time safety but can
// contain an option for this package or for wrappers implementing this
// interface.
type Option func() interface{}

// OptionFunc - a type for funcs that operate on the shared Options struct. The
// options below explicitly wrap this so that we can switch on it when parsing
// opts for various wrappers.
type OptionFunc func(*Options) error

func getDefaultOptions() *Options {
	return &Options{}
}

// WithAad provides optional additional authenticated data
func WithAad(with []byte) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithAad = with
			return nil
		})
	}
}

// WithKeyId provides a common way to pass in a key identifier
func WithKeyId(with string) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithKeyId = with
			return nil
		})
	}
}

// WithKeyPurpose provides a common way to pass in a key purpose
func WithKeyPurposes(purpose ...KeyPurpose) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithKeyPurposes = purpose
			return nil
		})
	}
}

// WithKeyType provides a common way to pass in a key type
func WithKeyType(keyType KeyType) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithKeyType = keyType
			return nil
		})
	}
}

// WithRandomBytes provides a common way to pass in entropy
func WithRandomBytes(b []byte) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithRandomBytes = b
			return nil
		})
	}
}

// WithConfigMap is an option accepted by wrappers at configuration time
// and/or in other function calls to control wrapper-specific behavior.
func WithConfigMap(with map[string]string) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithConfigMap = with
			return nil
		})
	}
}

// WithIV provides
func WithIV(with []byte) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithIv = with
			return nil
		})
	}
}

// WithKeyEncoding provides a common way to pass in a key encoding
func WithKeyEncoding(encoding KeyEncoding) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithKeyEncoding = encoding
			return nil
		})
	}
}

// WithWrappedKeyEncoding provides a common way to pass in a wrapped_key encoding
func WithWrappedKeyEncoding(encoding KeyEncoding) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithWrappedKeyEncoding = encoding
			return nil
		})
	}
}

// WithDisallowEnvVars provides a common way to configure ignoring environment variables
func WithDisallowEnvVars(disallowEnvVars bool) Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithDisallowEnvVars = disallowEnvVars
			return nil
		})
	}
}

// WithoutHMAC disables the requirement for an HMAC to be included with the mechanism.
func WithoutHMAC() Option {
	return func() interface{} {
		return OptionFunc(func(o *Options) error {
			o.WithoutHmac = true
			return nil
		})
	}
}
