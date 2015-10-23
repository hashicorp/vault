package gocql

import (
	"log"
	"net"
	"time"
)

type HostInfo struct {
	Peer       string
	DataCenter string
	Rack       string
	HostId     string
	Tokens     []string
}

// Polls system.peers at a specific interval to find new hosts
type ringDescriber struct {
	dcFilter        string
	rackFilter      string
	prevHosts       []HostInfo
	prevPartitioner string
	session         *Session
	closeChan       chan bool
}

func (r *ringDescriber) GetHosts() (hosts []HostInfo, partitioner string, err error) {
	// we need conn to be the same because we need to query system.peers and system.local
	// on the same node to get the whole cluster

	iter := r.session.control.query("SELECT data_center, rack, host_id, tokens, partitioner FROM system.local")
	if iter == nil {
		return r.prevHosts, r.prevPartitioner, nil
	}

	conn := r.session.pool.Pick(nil)
	if conn == nil {
		return r.prevHosts, r.prevPartitioner, nil
	}

	host := HostInfo{}
	iter.Scan(&host.DataCenter, &host.Rack, &host.HostId, &host.Tokens, &partitioner)

	if err = iter.Close(); err != nil {
		return nil, "", err
	}

	addr, _, err := net.SplitHostPort(conn.Address())
	if err != nil {
		// this should not happen, ever, as this is the address that was dialed by conn, here
		// a panic makes sense, please report a bug if it occurs.
		panic(err)
	}

	host.Peer = addr

	hosts = []HostInfo{host}

	iter = r.session.control.query("SELECT peer, data_center, rack, host_id, tokens FROM system.peers")
	if iter == nil {
		return r.prevHosts, r.prevPartitioner, nil
	}

	host = HostInfo{}
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

func (h *ringDescriber) run(sleep time.Duration) {
	if sleep == 0 {
		sleep = 30 * time.Second
	}

	for {
		// if we have 0 hosts this will return the previous list of hosts to
		// attempt to reconnect to the cluster otherwise we would never find
		// downed hosts again, could possibly have an optimisation to only
		// try to add new hosts if GetHosts didnt error and the hosts didnt change.
		hosts, partitioner, err := h.GetHosts()
		if err != nil {
			log.Println("RingDescriber: unable to get ring topology:", err)
			continue
		}

		h.session.pool.SetHosts(hosts)
		h.session.pool.SetPartitioner(partitioner)

		select {
		case <-time.After(sleep):
		case <-h.closeChan:
			return
		}
	}
}
