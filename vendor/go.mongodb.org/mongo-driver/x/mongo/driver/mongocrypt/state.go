// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongocrypt

// State represents a state that a MongocryptContext can be in.
type State int

// These constants are valid values for the State type.
// The values must match the values defined in the mongocrypt_ctx_state_t enum in libmongocrypt.
const (
	StateError         State = 0
	NeedMongoCollInfo  State = 1
	NeedMongoMarkings  State = 2
	NeedMongoKeys      State = 3
	NeedKms            State = 4
	Ready              State = 5
	Done               State = 6
	NeedKmsCredentials State = 7
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
	case NeedKmsCredentials:
		return "NeedKmsCredentials"
	default:
		return "Unknown State"
	}
}
