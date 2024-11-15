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
	"fmt"
	"sync"
	"sync/atomic"
)

type ring struct {
	// endpoints are the set of endpoints which the driver will attempt to connect
	// to in the case it can not reach any of its hosts. They are also used to boot
	// strap the initial connection.
	endpoints []*HostInfo

	mu sync.RWMutex
	// hosts are the set of all hosts in the cassandra ring that we know of.
	// key of map is host_id.
	hosts map[string]*HostInfo
	// hostIPToUUID maps host native address to host_id.
	hostIPToUUID map[string]string

	hostList []*HostInfo
	pos      uint32

	// TODO: we should store the ring metadata here also.
}

func (r *ring) rrHost() *HostInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.hostList) == 0 {
		return nil
	}

	pos := int(atomic.AddUint32(&r.pos, 1) - 1)
	return r.hostList[pos%len(r.hostList)]
}

func (r *ring) getHostByIP(ip string) (*HostInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hi, ok := r.hostIPToUUID[ip]
	return r.hosts[hi], ok
}

func (r *ring) getHost(hostID string) *HostInfo {
	r.mu.RLock()
	host := r.hosts[hostID]
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

func (r *ring) currentHosts() map[string]*HostInfo {
	r.mu.RLock()
	hosts := make(map[string]*HostInfo, len(r.hosts))
	for k, v := range r.hosts {
		hosts[k] = v
	}
	r.mu.RUnlock()
	return hosts
}

func (r *ring) addOrUpdate(host *HostInfo) *HostInfo {
	if existingHost, ok := r.addHostIfMissing(host); ok {
		existingHost.update(host)
		host = existingHost
	}
	return host
}

func (r *ring) addHostIfMissing(host *HostInfo) (*HostInfo, bool) {
	if host.invalidConnectAddr() {
		panic(fmt.Sprintf("invalid host: %v", host))
	}
	hostID := host.HostID()

	r.mu.Lock()
	if r.hosts == nil {
		r.hosts = make(map[string]*HostInfo)
	}
	if r.hostIPToUUID == nil {
		r.hostIPToUUID = make(map[string]string)
	}

	existing, ok := r.hosts[hostID]
	if !ok {
		r.hosts[hostID] = host
		r.hostIPToUUID[host.nodeToNodeAddress().String()] = hostID
		existing = host
		r.hostList = append(r.hostList, host)
	}
	r.mu.Unlock()
	return existing, ok
}

func (r *ring) removeHost(hostID string) bool {
	r.mu.Lock()
	if r.hosts == nil {
		r.hosts = make(map[string]*HostInfo)
	}
	if r.hostIPToUUID == nil {
		r.hostIPToUUID = make(map[string]string)
	}

	h, ok := r.hosts[hostID]
	if ok {
		for i, host := range r.hostList {
			if host.HostID() == hostID {
				r.hostList = append(r.hostList[:i], r.hostList[i+1:]...)
				break
			}
		}
		delete(r.hostIPToUUID, h.nodeToNodeAddress().String())
	}
	delete(r.hosts, hostID)
	r.mu.Unlock()
	return ok
}

type clusterMetadata struct {
	mu          sync.RWMutex
	partitioner string
}

func (c *clusterMetadata) setPartitioner(partitioner string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.partitioner != partitioner {
		// TODO: update other things now
		c.partitioner = partitioner
	}
}
