// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

// TopologyKind represents a specific topology configuration.
type TopologyKind uint32

// These constants are the available topology configurations.
const (
	Single                TopologyKind = 1
	ReplicaSet            TopologyKind = 2
	ReplicaSetNoPrimary   TopologyKind = 4 + ReplicaSet
	ReplicaSetWithPrimary TopologyKind = 8 + ReplicaSet
	Sharded               TopologyKind = 256
	LoadBalanced          TopologyKind = 512
)

// String implements the fmt.Stringer interface.
func (kind TopologyKind) String() string {
	switch kind {
	case Single:
		return "Single"
	case ReplicaSet:
		return "ReplicaSet"
	case ReplicaSetNoPrimary:
		return "ReplicaSetNoPrimary"
	case ReplicaSetWithPrimary:
		return "ReplicaSetWithPrimary"
	case Sharded:
		return "Sharded"
	case LoadBalanced:
		return "LoadBalanced"
	}

	return "Unknown"
}
