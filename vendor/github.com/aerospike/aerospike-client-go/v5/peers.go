// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"sync"

	"github.com/aerospike/aerospike-client-go/v5/internal/atomic"
)

type peers struct {
	_peers       map[string]*peer
	_hosts       map[Host]struct{}
	_nodes       map[string]*Node
	refreshCount atomic.Int
	genChanged   atomic.Bool

	mutex sync.RWMutex
}

func newPeers(peerCapacity int, addCapacity int) *peers {
	return &peers{
		_peers:     make(map[string]*peer, peerCapacity),
		_hosts:     make(map[Host]struct{}, addCapacity),
		_nodes:     make(map[string]*Node, addCapacity),
		genChanged: *atomic.NewBool(true),
	}
}

func (ps *peers) hostExists(host Host) bool {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()
	_, exists := ps._hosts[host]
	return exists
}

func (ps *peers) addHost(host Host) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	ps._hosts[host] = struct{}{}
}

func (ps *peers) addNode(name string, node *Node) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	ps._nodes[name] = node
}

func (ps *peers) nodeByName(name string) *Node {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()
	return ps._nodes[name]
}

func (ps *peers) appendPeers(peers []*peer) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	for _, peer := range peers {
		ps._peers[peer.nodeName] = peer
	}

}

func (ps *peers) peers() []*peer {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	res := make([]*peer, 0, len(ps._peers))
	for _, peer := range ps._peers {
		res = append(res, peer)
	}
	return res
}

func (ps *peers) nodes() map[string]*Node {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()
	return ps._nodes
}

type peer struct {
	nodeName string
	tlsName  string
	hosts    []*Host
}
