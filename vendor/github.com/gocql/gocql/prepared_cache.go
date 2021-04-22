package gocql

import (
	"bytes"
	"github.com/gocql/gocql/internal/lru"
	"sync"
)

const defaultMaxPreparedStmts = 1000

// preparedLRU is the prepared statement cache
type preparedLRU struct {
	mu  sync.Mutex
	lru *lru.Cache
}

func (p *preparedLRU) clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for p.lru.Len() > 0 {
		p.lru.RemoveOldest()
	}
}

func (p *preparedLRU) add(key string, val *inflightPrepare) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lru.Add(key, val)
}

func (p *preparedLRU) remove(key string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.lru.Remove(key)
}

func (p *preparedLRU) execIfMissing(key string, fn func(lru *lru.Cache) *inflightPrepare) (*inflightPrepare, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	val, ok := p.lru.Get(key)
	if ok {
		return val.(*inflightPrepare), true
	}

	return fn(p.lru), false
}

func (p *preparedLRU) keyFor(addr, keyspace, statement string) string {
	// TODO: we should just use a struct for the key in the map
	return addr + keyspace + statement
}

func (p *preparedLRU) evictPreparedID(key string, id []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()

	val, ok := p.lru.Get(key)
	if !ok {
		return
	}

	ifp, ok := val.(*inflightPrepare)
	if !ok {
		return
	}

	select {
	case <-ifp.done:
		if bytes.Equal(id, ifp.preparedStatment.id) {
			p.lru.Remove(key)
		}
	default:
	}

}
