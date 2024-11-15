// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"github.com/hashicorp/go-hclog"
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
			case "mount_path":
				opts.withMountPath = v
			case "key_name":
				opts.withKeyName = v
			case "disable_renewal":
				opts.withDisableRenewal = v
			case "namespace":
				opts.withNamespace = v
			case "address":
				opts.withAddress = v
			case "tls_ca_cert":
				opts.withTlsCaCert = v
			case "tls_ca_path":
				opts.withTlsCaPath = v
			case "tls_client_cert":
				opts.withTlsClientCert = v
			case "tls_client_key":
				opts.withTlsClientKey = v
			case "tls_server_name":
				opts.withTlsServerName = v
			case "tls_skip_verify":
				var err error
				opts.withTlsSkipVerify, err = strconv.ParseBool(v)
				if err != nil {
					return nil, err
				}
			case "key_id_prefix":
				opts.withKeyIdPrefix = v
			case "token":
				opts.withToken = v
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

	withMountPath      string
	withKeyName        string
	withDisableRenewal string
	withNamespace      string
	withAddress        string
	withTlsCaCert      string
	withTlsCaPath      string
	withTlsClientCert  string
	withTlsClientKey   string
	withTlsServerName  string
	withTlsSkipVerify  bool
	withToken          string
	withKeyIdPrefix    string

	withLogger hclog.Logger
}

func getDefaultOptions() options {
	return options{}
}

// WithMountPath provides a way to choose the mount path
func WithMountPath(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withMountPath = with
			return nil
		})
	}
}

// WithKeyName provides a way to choose the key name
func WithKeyName(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withKeyName = with
			return nil
		})
	}
}

// WithDisableRenewal provides a way to disable renewal
func WithDisableRenewal(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withDisableRenewal = with
			return nil
		})
	}
}

// WithNamespace provides a way to choose the namespace
func WithNamespace(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withNamespace = with
			return nil
		})
	}
}

// WithAddress provides a way to choose the address
func WithAddress(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withAddress = with
			return nil
		})
	}
}

// WithTlsCaCert provides a way to choose the CA cert
func WithTlsCaCert(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withTlsCaCert = with
			return nil
		})
	}
}

// WithTlsCaPath provides a way to choose the CA path
func WithTlsCaPath(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withTlsCaPath = with
			return nil
		})
	}
}

// WithTlsClientCert provides a way to choose the client cert
func WithTlsClientCert(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withTlsClientCert = with
			return nil
		})
	}
}

// WithTlsClientKey provides a way to choose the client key
func WithTlsClientKey(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withTlsClientKey = with
			return nil
		})
	}
}

// WithTlsServerName provides a way to choose the server name
func WithTlsServerName(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withTlsServerName = with
			return nil
		})
	}
}

// WithTlsSkipVerify provides a way to skip TLS verification
func WithTlsSkipVerify(with bool) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withTlsSkipVerify = with
			return nil
		})
	}
}

// WithToken provides a way to choose the token
func WithToken(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withToken = with
			return nil
		})
	}
}

// WithLogger provides a way to pass in a logger
func WithLogger(with hclog.Logger) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withLogger = with
			return nil
		})
	}
}

// WithKeyIdPrefix specifies a prefix to prepend to the keyId (key version)
func WithKeyIdPrefix(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withKeyIdPrefix = with
			return nil
		})
	}
}
