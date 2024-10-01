// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

// ListenerConfigOption is how listenerConfigOptions are passed as arguments.
type ListenerConfigOption func(*listenerConfigOptions) error

// listenerConfigOptions are used to represent configuration of listeners for http handlers.
type listenerConfigOptions struct {
	withRedactionValue    string
	withRedactAddresses   bool
	withRedactClusterName bool
	withRedactVersion     bool
}

// getDefaultOptions returns listenerConfigOptions with their default values.
func getDefaultOptions() listenerConfigOptions {
	return listenerConfigOptions{
		withRedactionValue: "", // Redacted values will be set to an empty string by default.
	}
}

// getOpts applies each supplied ListenerConfigOption and returns the fully configured listenerConfigOptions.
// Each ListenerConfigOption is applied in the order it appears in the argument list, so it is
// possible to supply the same ListenerConfigOption numerous times and the 'last write wins'.
func getOpts(opt ...ListenerConfigOption) (listenerConfigOptions, error) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o == nil {
			continue
		}
		if err := o(&opts); err != nil {
			return listenerConfigOptions{}, err
		}
	}
	return opts, nil
}

// WithRedactionValue provides an ListenerConfigOption to represent the value used to redact
// values which require redaction.
func WithRedactionValue(r string) ListenerConfigOption {
	return func(o *listenerConfigOptions) error {
		o.withRedactionValue = r
		return nil
	}
}

// WithRedactAddresses provides an ListenerConfigOption to represent whether redaction of addresses is required.
func WithRedactAddresses(r bool) ListenerConfigOption {
	return func(o *listenerConfigOptions) error {
		o.withRedactAddresses = r
		return nil
	}
}

// WithRedactClusterName provides an ListenerConfigOption to represent whether redaction of cluster names is required.
func WithRedactClusterName(r bool) ListenerConfigOption {
	return func(o *listenerConfigOptions) error {
		o.withRedactClusterName = r
		return nil
	}
}

// WithRedactVersion provides an ListenerConfigOption to represent whether redaction of version is required.
func WithRedactVersion(r bool) ListenerConfigOption {
	return func(o *listenerConfigOptions) error {
		o.withRedactVersion = r
		return nil
	}
}
