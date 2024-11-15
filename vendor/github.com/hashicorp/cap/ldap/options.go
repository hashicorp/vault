// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ldap

// Option defines a common functional options type which can be used in a
// variadic parameter pattern.
type Option func(interface{})

type configOptions struct {
	withURLs                      []string
	withInsecureTLS               bool
	withTLSMinVersion             string
	withTLSMaxVersion             string
	withCertificates              []string
	withClientTLSCert             string
	withClientTLSKey              string
	withGroups                    bool
	withUserAttributes            bool
	withLowerUserAttributeKeys    bool
	withEmptyAnonymousGroupSearch bool
}

func configDefaults() configOptions {
	return configOptions{}
}

// getConfigOpts gets the defaults and applies the opt overrides passed
// in.
func getConfigOpts(opt ...Option) configOptions {
	opts := configDefaults()
	ApplyOpts(&opts, opt...)
	return opts
}

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

// WithURLs provides a set of optional ldap URLs for directory services
func WithURLs(urls ...string) Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withURLs = urls
		}
	}
}

// WithGroups requests that the groups be included in the response.
func WithGroups() Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withGroups = true
		}
	}
}

// WithUserAttributes requests that authenticating user's DN and attributes be
// included in the response. Note: the default password attribute for both
// openLDAP (userPassword) and AD (unicodePwd) will always be excluded.  To
// exclude additional attributes see: Config.ExcludedUserAttributes.
func WithUserAttributes() Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withUserAttributes = true
		}
	}
}

// WithLowerUserAttributeKeys returns a User Attribute map where the keys
// are all cast to lower case. This is necessary for some clients, such as Vault,
// where user configured user attribute key names have always been stored lower case.
func WithLowerUserAttributeKeys() Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withLowerUserAttributeKeys = true
		}
	}
}

// WithEmptyAnonymousGroupSearch removes userDN from anonymous group searches.
func WithEmptyAnonymousGroupSearch() Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withEmptyAnonymousGroupSearch = true
		}
	}
}

func withTLSMinVersion(version string) Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withTLSMinVersion = version
		}
	}
}

func withTLSMaxVersion(version string) Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withTLSMaxVersion = version
		}
	}
}

func withInsecureTLS(withInsecure bool) Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withInsecureTLS = withInsecure
		}
	}
}

func withCertificates(cert ...string) Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withCertificates = cert
		}
	}
}

func withClientTLSKey(key string) Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withClientTLSKey = key
		}
	}
}

func withClientTLSCert(cert string) Option {
	return func(o interface{}) {
		switch v := o.(type) {
		case *configOptions:
			v.withClientTLSCert = cert
		}
	}
}
