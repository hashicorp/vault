// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

// Option is how options are passed as arguments.
type Option func(*options) error

// options are used to represent configuration for http handlers.
type options struct {
	withRedactionValue    string
	withRedactAddresses   bool
	withRedactClusterName bool
	withRedactVersion     bool
}

// getDefaultOptions returns options with their default values.
func getDefaultOptions() options {
	return options{
		withRedactionValue: "", // Redact using empty string.
	}
}

// getOpts applies each supplied Option and returns the fully configured options.
// Each Option is applied in the order it appears in the argument list, so it is
// possible to supply the same Option numerous times and the 'last write wins'.
func getOpts(opt ...Option) (options, error) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o == nil {
			continue
		}
		if err := o(&opts); err != nil {
			return options{}, err
		}
	}
	return opts, nil
}

// WithRedactionValue provides an Option to represent the value used to redact
// values which require redaction.
func WithRedactionValue(r string) Option {
	return func(o *options) error {
		o.withRedactionValue = r
		return nil
	}
}

// WithRedactAddresses provides an Option to represent whether redaction of addresses is required.
func WithRedactAddresses(r bool) Option {
	return func(o *options) error {
		o.withRedactAddresses = r
		return nil
	}
}

// WithRedactClusterName provides an Option to represent whether redaction of cluster names is required.
func WithRedactClusterName(r bool) Option {
	return func(o *options) error {
		o.withRedactClusterName = r
		return nil
	}
}

// WithRedactVersion provides an Option to represent whether redaction of version is required.
func WithRedactVersion(r bool) Option {
	return func(o *options) error {
		o.withRedactVersion = r
		return nil
	}
}
