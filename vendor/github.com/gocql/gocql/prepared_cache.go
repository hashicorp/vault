/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * Content before git sha 34fdeebefcbf183ed7f916f931aa0586fdaa1b40
 * Copyright (c) 2016, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

package gocql

import (
	"bytes"
	"sync"

	"github.com/gocql/gocql/internal/lru"
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

func (p *preparedLRU) keyFor(hostID, keyspace, statement string) string {
	// TODO: we should just use a struct for the key in the map
	return hostID + keyspace + statement
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
