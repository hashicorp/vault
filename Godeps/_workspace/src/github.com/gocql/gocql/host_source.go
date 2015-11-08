package gocql

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type HostInfo struct {
	Peer       string
	DataCenter string
	Rack       string
	HostId     string
	Tokens     []string
}

func (h HostInfo) String() string {
	return fmt.Sprintf("[hostinfo peer=%q data_centre=%q rack=%q host_id=%q num_tokens=%d]", h.Peer, h.DataCenter, h.Rack, h.HostId, len(h.Tokens))
}

// Polls system.peers at a specific interval to find new hosts
type ringDescriber struct {
	dcFilter        string
	rackFilter      string
	prevHosts       []HostInfo
	prevPartitioner string
	session         *Session
	closeChan       chan bool
	// indicates that we can use system.local to get the connections remote address
	localHasRpcAddr bool

	mu sync.Mutex
}

func checkSystemLocal(control *controlConn) (bool, error) {
	iter := control.query("SELECT broadcast_address FROM system.local")
	if err := iter.err; err != nil {
		if errf, ok := err.(*errorFrame); ok {
			if errf.code == errSyntax {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func (r *ringDescriber) GetHosts() (hosts []HostInfo, partitioner string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// we need conn to be the same because we need to query system.peers and system.local
	// on the same node to get the whole cluster

	const (
		legacyLocalQuery = "SELECT data_center, rack, host_id, tokens, partitioner FROM system.local"
		// only supported in 2.2.0, 2.1.6, 2.0.16
		localQuery = "SELECT broadcast_address, data_center, rack, host_id, tokens, partitioner FROM system.local"
	)

	var localHost HostInfo
	if r.localHasRpcAddr {
		iter := r.session.control.query(localQuery)
		if iter == nil {
			return r.prevHosts, r.prevPartitioner, nil
		}

		iter.Scan(&localHost.Peer, &localHost.DataCenter, &localHost.Rack,
			&localHost.HostId, &localHost.Tokens, &partitioner)

		if err = iter.Close(); err != nil {
			return nil, "", err
		}
	} else {
		iter := r.session.control.query(legacyLocalQuery)
		if iter == nil {
			return r.prevHosts, r.prevPartitioner, nil
		}

		iter.Scan(&localHost.DataCenter, &localHost.Rack, &localHost.HostId, &localHost.Tokens, &partitioner)

		if err = iter.Close(); err != nil {
			return nil, "", err
		}

		addr, _, err := net.SplitHostPort(r.session.control.addr())
		if err != nil {
			// this should not happen, ever, as this is the address that was dialed by conn, here
			// a panic makes sense, please report a bug if it occurs.
			panic(err)
		}

		localHost.Peer = addr
	}

	hosts = []HostInfo{localHost}

	iter := r.session.control.query("SELECT rpc_address, data_center, rack, host_id, tokens FROM system.peers")
	if iter == nil {
		return r.prevHosts, r.prevPartitioner, nil
	}

	host := HostInfo{}
	for iter.Scan(&host.Peer, &host.DataCenter, &host.Rack, &host.HostId, &host.Tokens) {
		if r.matchFilter(&host) {
			hosts = append(hosts, host)
		}
		host = HostInfo{}
	}

	if err = iter.Close(); err != nil {
		return nil, "", err
	}

	r.prevHosts = hosts
	r.prevPartitioner = partitioner

	return hosts, partitioner, nil
}

func (r *ringDescriber) matchFilter(host *HostInfo) bool {

	if r.dcFilter != "" && r.dcFilter != host.DataCenter {
		return false
	}

	if r.rackFilter != "" && r.rackFilter != host.Rack {
		return false
	}

	return true
}

func (r *ringDescriber) refreshRing() {
	// if we have 0 hosts this will return the previous list of hosts to
	// attempt to reconnect to the cluster otherwise we would never find
	// downed hosts again, could possibly have an optimisation to only
	// try to add new hosts if GetHosts didnt error and the hosts didnt change.
	hosts, partitioner, err := r.GetHosts()
	if err != nil {
		log.Println("RingDescriber: unable to get ring topology:", err)
		return
	}

	r.session.pool.SetHosts(hosts)
	r.session.pool.SetPartitioner(partitioner)
}

func (r *ringDescriber) run(sleep time.Duration) {
	if sleep == 0 {
		sleep = 30 * time.Second
	}

	for {
		r.refreshRing()

		select {
		case <-time.After(sleep):
		case <-r.closeChan:
			return
		}
	}
}
