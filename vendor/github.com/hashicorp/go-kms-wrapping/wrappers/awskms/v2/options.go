// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awskms

import (
	"strconv"

	"github.com/hashicorp/go-hclog"
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
			case "disallow_env_vars":
				disallowEnvVars, err := strconv.ParseBool(v)
				if err != nil {
					return nil, err
				}
				opts.withDisallowEnvVars = disallowEnvVars
			case "key_not_required":
				keyNotRequired, err := strconv.ParseBool(v)
				if err != nil {
					return nil, err
				}
				opts.withKeyNotRequired = keyNotRequired
			case "kms_key_id": // deprecated backend-specific value, set global
				opts.WithKeyId = v
			case "region":
				opts.withRegion = v
			case "endpoint":
				opts.withEndpoint = v
			case "access_key":
				opts.withAccessKey = v
			case "secret_key":
				opts.withSecretKey = v
			case "session_token":
				opts.withSessionToken = v
			case "shared_creds_filename":
				opts.withSharedCredsFilename = v
			case "shared_creds_profile":
				opts.withSharedCredsProfile = v
			case "web_identity_token_file":
				opts.withWebIdentityTokenFile = v
			case "role_session_name":
				opts.withRoleSessionName = v
			case "role_arn":
				opts.withRoleArn = v
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

	withDisallowEnvVars      bool
	withKeyNotRequired       bool
	withRegion               string
	withEndpoint             string
	withAccessKey            string
	withSecretKey            string
	withSessionToken         string
	withSharedCredsFilename  string
	withSharedCredsProfile   string
	withWebIdentityTokenFile string
	withRoleSessionName      string
	withRoleArn              string

	withLogger hclog.Logger
}

func getDefaultOptions() options {
	return options{}
}

// WithDisallowEnvVars provides a way to disable using env vars
func WithDisallowEnvVars(with bool) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withDisallowEnvVars = with
			return nil
		})
	}
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

// WithRegion provides a way to chose the region
func WithRegion(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withRegion = with
			return nil
		})
	}
}

// WithEndpoint provides a way to chose the endpoint
func WithEndpoint(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withEndpoint = with
			return nil
		})
	}
}

// WithAccessKey provides a way to chose the access key
func WithAccessKey(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withAccessKey = with
			return nil
		})
	}
}

// WithSecretKey provides a way to chose the secret key
func WithSecretKey(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withSecretKey = with
			return nil
		})
	}
}

// WithSessionToken provides a way to chose the session token
func WithSessionToken(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withSessionToken = with
			return nil
		})
	}
}

// WithSharedCredsFilename provides a way to chose the shared creds filename
func WithSharedCredsFilename(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withSharedCredsFilename = with
			return nil
		})
	}
}

// WithSharedCredsProfile provides a way to chose the shared creds profile
func WithSharedCredsProfile(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withSharedCredsProfile = with
			return nil
		})
	}
}

// WithWebIdentityTokenFile provides a way to chose the web identity token file
func WithWebIdentityTokenFile(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withWebIdentityTokenFile = with
			return nil
		})
	}
}

// WithRoleSessionName provides a way to chose the role session name
func WithRoleSessionName(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withRoleSessionName = with
			return nil
		})
	}
}

// WithRoleArn provides a way to chose the role ARN
func WithRoleArn(with string) wrapping.Option {
	return func() interface{} {
		return OptionFunc(func(o *options) error {
			o.withRoleArn = with
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
