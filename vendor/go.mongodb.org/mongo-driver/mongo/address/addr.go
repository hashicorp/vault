// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package address // import "go.mongodb.org/mongo-driver/mongo/address"

import (
	"net"
	"strings"
)

const defaultPort = "27017"

// Address is a network address. It can either be an IP address or a DNS name.
type Address string

// Network is the network protocol for this address. In most cases this will be
// "tcp" or "unix".
func (a Address) Network() string {
	if strings.HasSuffix(string(a), "sock") {
		return "unix"
	}
	return "tcp"
}

// String is the canonical version of this address, e.g. localhost:27017,
// 1.2.3.4:27017, example.com:27017.
func (a Address) String() string {
	// TODO: unicode case folding?
	s := strings.ToLower(string(a))
	if len(s) == 0 {
		return ""
	}
	if a.Network() != "unix" {
		_, _, err := net.SplitHostPort(s)
		if err != nil && strings.Contains(err.Error(), "missing port in address") {
			s += ":" + defaultPort
		}
	}

	return s
}

// Canonicalize creates a canonicalized address.
func (a Address) Canonicalize() Address {
	return Address(a.String())
}
