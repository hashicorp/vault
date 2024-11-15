// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpckms

import (
	"strconv"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

// getOpts iterates the inbound Options and returns a struct
func getOpts(opt ...wrapping.Option) (*options, error) {
	// First, separate out options into local and global
	opts := getDefaultOptions()
	var wrappingOptions []wrapping.Option
	var localOptions []OptionFunc
	for _, o := range opt {
		if o == nil {
			continue
		}
		iface := o()
		switch to := iface.(type) {
		case wrapping.OptionFunc:
			wrappingOptions = append(wrappingOptions, o)
		case OptionFunc:
			localOptions = append(localOptions, to)
		}
	}

	// Parse the global options
	var err error
	opts.Options, err = wrapping.GetOpts(wrappingOptions...)
	if err != nil {
		return nil, err
	}

	// Don't ever return blank options
	if opts.Options == nil {
		opts.Options = new(wrapping.Options)
	}

	// Local options can be provided either via the WithConfigMap field
	// (for over the plugin barrier or embedding) or via local option functions
	// (for embedding). First pull from the option.
	if opts.WithConfigMap != nil {
		for k, v := range opts.WithConfigMap {
			switch k {
			case "key_not_required":
				keyNotRequired, err := strconv.ParseBool(v)
				if err != nil {
					return nil, err
				}
				opts.withKeyNotRequired = keyNotRequired
			case "user_agent":
				opts.withUserAgent = v
			case "credentials":
				opts.withCredentials = v
			case "project":
				opts.withProject = v
			case "region":
				opts.withRegion = v
			case "key_ring":
				opts.withKeyRing = v
			case "crypto_key":
				opts.withCryptoKey = v
			}
		}
	}

	// Now run the local options functions. This may overwrite options set by
	// the options above.
	for _, o := range localOptions {
		if o != nil {
			if err := o(&opts); err != nil {
				return nil, err
			}
		}
	}

	return &opts, nil
}

// OptionFunc holds a function with local options
type OptionFunc func(*options) error

// options = how options are represented
type options struct {
	*wrapping.Options

	withKeyNotRequired bool
	withUserAgent      string
	withCredentials    string
	withProject        string
	withRegion         string
	withKeyRing        string
	withCryptoKey      string
}

func getDefaultOptions() options {
	return options{}
}

// WithKeyNotRequired provides a way to not require a key at config time
func WithKeyNotRequired(with bool) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withKeyNotRequired = with
			return nil
		})
	}
}

// WithUserAgent provides a way to chose the user agent
func WithUserAgent(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withUserAgent = with
			return nil
		})
	}
}

// WithCredentials provides a way to specify credentials
func WithCredentials(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withCredentials = with
			return nil
		})
	}
}

// WithProject provides a way to chose the project
func WithProject(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withProject = with
			return nil
		})
	}
}

// WithRegion provides a way to chose the region
func WithRegion(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withRegion = with
			return nil
		})
	}
}

// WithKeyRing provides a way to chose the key ring
func WithKeyRing(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withKeyRing = with
			return nil
		})
	}
}

// WithCryptoKey provides a way to chose the crypto key
func WithCryptoKey(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withCryptoKey = with
			return nil
		})
	}
}
