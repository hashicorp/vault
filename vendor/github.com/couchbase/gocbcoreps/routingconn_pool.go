package gocbcoreps

import (
	"sync/atomic"
)

type routingConnPool struct {
	conns []*routingConn

	idx  uint32
	size uint32
}

func newRoutingConnPool(conns []*routingConn) *routingConnPool {
	return &routingConnPool{
		conns: conns,
		size:  uint32(len(conns)),
	}
}

func (pool *routingConnPool) Conn() *routingConn {
	idx := atomic.AddUint32(&pool.idx, 1)
	return pool.conns[idx%pool.Size()]
}

func (pool *routingConnPool) Size() uint32 {
	return pool.size
}

func (pool *routingConnPool) Close() error {
	var err error
	for _, conn := range pool.conns {
		closeErr := conn.Close()
		if closeErr != nil {
			err = closeErr
		}
	}

	return err
}

func (pool *routingConnPool) State() ConnState {
	var numOnline uint32
	var numOffline uint32
	for _, conn := range pool.conns {
		switch conn.State() {
		case ConnStateOffline:
			numOffline++
		case ConnStateOnline:
			numOnline++
		}
	}
	if numOffline == pool.Size() {
		// If all connections are offline then our state is offline.
		return ConnStateOffline
	} else if numOnline == pool.Size() {
		// If all connections are online then our state is online.
		return ConnStateOnline
	} else {
		// If we have some connections online and some offline then we're degraded.
		return ConnStateDegraded
	}
}
