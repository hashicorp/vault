package gocql

import (
	"sync"
)

type ring struct {
	// endpoints are the set of endpoints which the driver will attempt to connect
	// to in the case it can not reach any of its hosts. They are also used to boot
	// strap the initial connection.
	endpoints []string
	// hosts are the set of all hosts in the cassandra ring that we know of
	mu    sync.RWMutex
	hosts map[string]*HostInfo

	// TODO: we should store the ring metadata here also.
}

func (r *ring) getHost(addr string) *HostInfo {
	r.mu.RLock()
	host := r.hosts[addr]
	r.mu.RUnlock()
	return host
}

func (r *ring) allHosts() []*HostInfo {
	r.mu.RLock()
	hosts := make([]*HostInfo, 0, len(r.hosts))
	for _, host := range r.hosts {
		hosts = append(hosts, host)
	}
	r.mu.RUnlock()
	return hosts
}

func (r *ring) addHost(host *HostInfo) bool {
	r.mu.Lock()
	if r.hosts == nil {
		r.hosts = make(map[string]*HostInfo)
	}

	addr := host.Peer()
	_, ok := r.hosts[addr]
	r.hosts[addr] = host
	r.mu.Unlock()
	return ok
}

func (r *ring) addHostIfMissing(host *HostInfo) (*HostInfo, bool) {
	r.mu.Lock()
	if r.hosts == nil {
		r.hosts = make(map[string]*HostInfo)
	}

	addr := host.Peer()
	existing, ok := r.hosts[addr]
	if !ok {
		r.hosts[addr] = host
		existing = host
	}
	r.mu.Unlock()
	return existing, ok
}

func (r *ring) removeHost(addr string) bool {
	r.mu.Lock()
	if r.hosts == nil {
		r.hosts = make(map[string]*HostInfo)
	}

	_, ok := r.hosts[addr]
	delete(r.hosts, addr)
	r.mu.Unlock()
	return ok
}
