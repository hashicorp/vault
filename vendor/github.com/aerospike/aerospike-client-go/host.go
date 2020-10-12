// Copyright 2013-2020 Aerospike, Inc.
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
	"net"
	"strconv"
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
