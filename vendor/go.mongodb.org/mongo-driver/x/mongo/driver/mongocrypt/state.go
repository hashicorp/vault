// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongocrypt

// State represents a state that a MongocryptContext can be in.
type State int

// These constants are valid values for the State type.
const (
	StateError State = iota
	NeedMongoCollInfo
	NeedMongoMarkings
	NeedMongoKeys
	NeedKms
	Ready
	Done
)

// String implements the Stringer interface.
func (s State) String() string {
	switch s {
	case StateError:
		return "Error"
	case NeedMongoCollInfo:
		return "NeedMongoCollInfo"
	case NeedMongoMarkings:
		return "NeedMongoMarkings"
	case NeedMongoKeys:
		return "NeedMongoKeys"
	case NeedKms:
		return "NeedKms"
	case Ready:
		return "Ready"
	case Done:
		return "Done"
	default:
		return "Unknown State"
	}
}
