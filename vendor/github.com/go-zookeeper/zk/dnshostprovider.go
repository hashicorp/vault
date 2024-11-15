package zk

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

const _defaultLookupTimeout = 3 * time.Second

type lookupHostFn func(context.Context, string) ([]string, error)

// DNSHostProviderOption is an option for the DNSHostProvider.
type DNSHostProviderOption interface {
	apply(*DNSHostProvider)
}

type lookupTimeoutOption struct {
	timeout time.Duration
}

// WithLookupTimeout returns a DNSHostProviderOption that sets the lookup timeout.
func WithLookupTimeout(timeout time.Duration) DNSHostProviderOption {
	return lookupTimeoutOption{
		timeout: timeout,
	}
}

func (o lookupTimeoutOption) apply(provider *DNSHostProvider) {
	provider.lookupTimeout = o.timeout
}

// DNSHostProvider is the default HostProvider. It currently matches
// the Java StaticHostProvider, resolving hosts from DNS once during
// the call to Init.  It could be easily extended to re-query DNS
// periodically or if there is trouble connecting.
type DNSHostProvider struct {
	mu            sync.Mutex // Protects everything, so we can add asynchronous updates later.
	servers       []string
	curr          int
	last          int
	lookupTimeout time.Duration
	lookupHost    lookupHostFn // Override of net.LookupHost, for testing.
}

// NewDNSHostProvider creates a new DNSHostProvider with the given options.
func NewDNSHostProvider(options ...DNSHostProviderOption) *DNSHostProvider {
	var provider DNSHostProvider
	for _, option := range options {
		option.apply(&provider)
	}
	return &provider
}

// Init is called first, with the servers specified in the connection
// string. It uses DNS to look up addresses for each server, then
// shuffles them all together.
func (hp *DNSHostProvider) Init(servers []string) error {
	hp.mu.Lock()
	defer hp.mu.Unlock()

	lookupHost := hp.lookupHost
	if lookupHost == nil {
		var resolver net.Resolver
		lookupHost = resolver.LookupHost
	}

	timeout := hp.lookupTimeout
	if timeout == 0 {
		timeout = _defaultLookupTimeout
	}

	// TODO: consider using a context from the caller.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	found := []string{}
	for _, server := range servers {
		host, port, err := net.SplitHostPort(server)
		if err != nil {
			return err
		}
		addrs, err := lookupHost(ctx, host)
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			found = append(found, net.JoinHostPort(addr, port))
		}
	}

	if len(found) == 0 {
		return fmt.Errorf("No hosts found for addresses %q", servers)
	}

	// Randomize the order of the servers to avoid creating hotspots
	stringShuffle(found)

	hp.servers = found
	hp.curr = -1
	hp.last = -1

	return nil
}

// Len returns the number of servers available
func (hp *DNSHostProvider) Len() int {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	return len(hp.servers)
}

// Next returns the next server to connect to. retryStart will be true
// if we've looped through all known servers without Connected() being
// called.
func (hp *DNSHostProvider) Next() (server string, retryStart bool) {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	hp.curr = (hp.curr + 1) % len(hp.servers)
	retryStart = hp.curr == hp.last
	if hp.last == -1 {
		hp.last = 0
	}
	return hp.servers[hp.curr], retryStart
}

// Connected notifies the HostProvider of a successful connection.
func (hp *DNSHostProvider) Connected() {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	hp.last = hp.curr
}
