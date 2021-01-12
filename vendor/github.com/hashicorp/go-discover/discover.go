// Package discover provides functions to get metadata for different
// cloud environments.
package discover

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/go-discover/provider/aliyun"
	"github.com/hashicorp/go-discover/provider/aws"
	"github.com/hashicorp/go-discover/provider/azure"
	"github.com/hashicorp/go-discover/provider/digitalocean"
	"github.com/hashicorp/go-discover/provider/gce"
	"github.com/hashicorp/go-discover/provider/linode"
	"github.com/hashicorp/go-discover/provider/mdns"
	"github.com/hashicorp/go-discover/provider/os"
	"github.com/hashicorp/go-discover/provider/packet"
	"github.com/hashicorp/go-discover/provider/scaleway"
	"github.com/hashicorp/go-discover/provider/softlayer"
	"github.com/hashicorp/go-discover/provider/tencentcloud"
	"github.com/hashicorp/go-discover/provider/triton"
	"github.com/hashicorp/go-discover/provider/vsphere"
)

// Provider has lookup functions for meta data in a
// cloud environment.
type Provider interface {
	// Addrs looks up addresses in the cloud environment according to the
	// configuration provided in args.
	Addrs(args map[string]string, l *log.Logger) ([]string, error)

	// Help provides the configuration help for the command line client.
	Help() string
}

// ProviderWithUserAgent is a provider that declares it's user agent. Not all
// providers support this.
type ProviderWithUserAgent interface {
	// SetUserAgent sets the user agent on the provider to the provided string.
	SetUserAgent(s string)
}

// Providers contains all available providers.
var Providers = map[string]Provider{
	"aliyun":       &aliyun.Provider{},
	"aws":          &aws.Provider{},
	"azure":        &azure.Provider{},
	"digitalocean": &digitalocean.Provider{},
	"gce":          &gce.Provider{},
	"linode":       &linode.Provider{},
	"mdns":         &mdns.Provider{},
	"os":           &os.Provider{},
	"scaleway":     &scaleway.Provider{},
	"softlayer":    &softlayer.Provider{},
	"tencentcloud": &tencentcloud.Provider{},
	"triton":       &triton.Provider{},
	"vsphere":      &vsphere.Provider{},
	"packet":       &packet.Provider{},
}

// Discover looks up metadata in different cloud environments.
type Discover struct {
	// Providers is the list of address lookup providers.
	// If nil, the default list of providers is used.
	Providers map[string]Provider

	// userAgent is the string to use for requests, when supported.
	userAgent string

	// once is used to initialize the actual list of providers.
	once sync.Once
}

// Option is used as an initialization option/
type Option func(*Discover) error

// New creates a new discover client with the given options.
func New(opts ...Option) (*Discover, error) {
	d := new(Discover)

	for _, opt := range opts {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	d.once.Do(d.initProviders)

	return d, nil
}

// WithUserAgent allows specifying a custom user agent option to send with
// requests when the underlying client library supports it.
func WithUserAgent(agent string) Option {
	return func(d *Discover) error {
		d.userAgent = agent
		return nil
	}
}

// WithProviders allows specifying your own set of providers.
func WithProviders(m map[string]Provider) Option {
	return func(d *Discover) error {
		d.Providers = m
		return nil
	}
}

// initProviders sets the list of providers to the
// default list of providers if none are configured.
func (d *Discover) initProviders() {
	if d.Providers == nil {
		d.Providers = Providers
	}
}

// Names returns the names of the configured providers.
func (d *Discover) Names() []string {
	d.once.Do(d.initProviders)

	var names []string
	for n := range d.Providers {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

var globalHelp = `The options for discovering ip addresses are provided as a
  single string value in "key=value key=value ..." format where
  the values are URL encoded.

    provider=aws region=eu-west-1 ...

  The options are provider specific and are listed below.
`

// Help describes the format of the configuration string for address discovery
// and the various provider specific options.
func (d *Discover) Help() string {
	d.once.Do(d.initProviders)

	h := []string{globalHelp}
	for _, name := range d.Names() {
		h = append(h, d.Providers[name].Help())
	}
	return strings.Join(h, "\n")
}

// Addrs discovers ip addresses of nodes that match the given filter criteria.
// The config string must have the format 'provider=xxx key=val key=val ...'
// where the keys and values are provider specific. The values are URL encoded.
func (d *Discover) Addrs(cfg string, l *log.Logger) ([]string, error) {
	d.once.Do(d.initProviders)

	args, err := Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("discover: %s", err)
	}

	name := args["provider"]
	if name == "" {
		return nil, fmt.Errorf("discover: no provider")
	}

	providers := d.Providers
	if providers == nil {
		providers = Providers
	}

	p := providers[name]
	if p == nil {
		return nil, fmt.Errorf("discover: unknown provider " + name)
	}
	l.Printf("[DEBUG] discover: Using provider %q", name)

	if typ, ok := p.(ProviderWithUserAgent); ok {
		typ.SetUserAgent(d.userAgent)
		return p.Addrs(args, l)
	}

	return p.Addrs(args, l)
}
