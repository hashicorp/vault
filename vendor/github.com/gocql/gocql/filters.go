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

import "fmt"

// HostFilter interface is used when a host is discovered via server sent events.
type HostFilter interface {
	// Called when a new host is discovered, returning true will cause the host
	// to be added to the pools.
	Accept(host *HostInfo) bool
}

// HostFilterFunc converts a func(host HostInfo) bool into a HostFilter
type HostFilterFunc func(host *HostInfo) bool

func (fn HostFilterFunc) Accept(host *HostInfo) bool {
	return fn(host)
}

// AcceptAllFilter will accept all hosts
func AcceptAllFilter() HostFilter {
	return HostFilterFunc(func(host *HostInfo) bool {
		return true
	})
}

func DenyAllFilter() HostFilter {
	return HostFilterFunc(func(host *HostInfo) bool {
		return false
	})
}

// DataCentreHostFilter filters all hosts such that they are in the same data centre
// as the supplied data centre.
func DataCentreHostFilter(dataCentre string) HostFilter {
	return HostFilterFunc(func(host *HostInfo) bool {
		return host.DataCenter() == dataCentre
	})
}

// WhiteListHostFilter filters incoming hosts by checking that their address is
// in the initial hosts whitelist.
func WhiteListHostFilter(hosts ...string) HostFilter {
	hostInfos, err := addrsToHosts(hosts, 9042, nopLogger{})
	if err != nil {
		// dont want to panic here, but rather not break the API
		panic(fmt.Errorf("unable to lookup host info from address: %v", err))
	}

	m := make(map[string]bool, len(hostInfos))
	for _, host := range hostInfos {
		m[host.ConnectAddress().String()] = true
	}

	return HostFilterFunc(func(host *HostInfo) bool {
		return m[host.ConnectAddress().String()]
	})
}
