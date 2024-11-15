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
	"fmt"
	"net"
	"strconv"

	"github.com/aerospike/aerospike-client-go/v5/types"
)

// Host name/port of database server.
type Host struct {

	// Host name or IP address of database server.
	Name string

	//TLSName defines the TLS certificate name used for secure connections.
	TLSName string

	// Port of database server.
	Port int
}

// NewHost initializes new host instance.
func NewHost(name string, port int) *Host {
	return &Host{Name: name, Port: port}
}

// Implements stringer interface
func (h *Host) String() string {
	return net.JoinHostPort(h.Name, strconv.Itoa(h.Port))
}

// Implements stringer interface
func (h *Host) equals(other *Host) bool {
	return h.Name == other.Name && h.Port == other.Port
}

// NewHosts initializes new host instances by a passed slice of addresses.
func NewHosts(addresses ...string) ([]*Host, Error) {
	aerospikeHosts := make([]*Host, 0, len(addresses))
	for _, address := range addresses {
		hostStr, portStr, err := net.SplitHostPort(address)
		if err != nil {
			return nil, newErrorAndWrap(err, types.PARAMETER_ERROR, fmt.Sprintf("error parsing address %s: %s", address, err))
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, newErrorAndWrap(err, types.PARAMETER_ERROR, fmt.Sprintf("error converting port %s: %s", address, err))
		}

		host := NewHost(hostStr, port)
		aerospikeHosts = append(aerospikeHosts, host)
	}

	return aerospikeHosts, nil
}
