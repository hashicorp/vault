// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ocikms

import (
	"fmt"
	"strconv"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

const (
	// keyId config
	KmsConfigKeyId = "key_id"
	// cryptoEndpoint config
	KmsConfigCryptoEndpoint = "crypto_endpoint"
	// managementEndpoint config
	KmsConfigManagementEndpoint = "management_endpoint"
	// authTypeApiKey config
	KmsConfigAuthTypeApiKey = "auth_type_api_key"
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
			case KmsConfigCryptoEndpoint:
				opts.withCryptoEndpoint = v
			case KmsConfigManagementEndpoint:
				opts.withManagementEndpoint = v
			case KmsConfigAuthTypeApiKey:
				var err error
				opts.withAuthTypeApiKey, err = strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("failed parsing "+KmsConfigAuthTypeApiKey+" parameter: %w", err)
				}
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

	withCryptoEndpoint     string
	withManagementEndpoint string
	withAuthTypeApiKey     bool
}

func getDefaultOptions() options {
	return options{}
}

// WithCryptoEndpoint provides a way to chose the endpoint
func WithCryptoEndpoint(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withCryptoEndpoint = with
			return nil
		})
	}
}

// WithManagementEndpoint provides a way to chose the management endpoint
func WithManagementEndpoint(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withManagementEndpoint = with
			return nil
		})
	}
}

// WithAuthTypeApiKey provides a way to say to use api keys for auth
func WithAuthTypeApiKey(with bool) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withAuthTypeApiKey = with
			return nil
		})
	}
}
