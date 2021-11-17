// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

// ServerKind represents the type of a single server in a topology.
type ServerKind uint32

// These constants are the possible types of servers.
const (
	Standalone   ServerKind = 1
	RSMember     ServerKind = 2
	RSPrimary    ServerKind = 4 + RSMember
	RSSecondary  ServerKind = 8 + RSMember
	RSArbiter    ServerKind = 16 + RSMember
	RSGhost      ServerKind = 32 + RSMember
	Mongos       ServerKind = 256
	LoadBalancer ServerKind = 512
)

// String returns a stringified version of the kind or "Unknown" if the kind is invalid.
func (kind ServerKind) String() string {
	switch kind {
	case Standalone:
		return "Standalone"
	case RSMember:
		return "RSOther"
	case RSPrimary:
		return "RSPrimary"
	case RSSecondary:
		return "RSSecondary"
	case RSArbiter:
		return "RSArbiter"
	case RSGhost:
		return "RSGhost"
	case Mongos:
		return "Mongos"
	case LoadBalancer:
		return "LoadBalancer"
	}

	return "Unknown"
}
