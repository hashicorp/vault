// Copyright 2014-2021 Aerospike, Inc.
//
// Portions may be licensed to Aerospike, Inc. under one or more contributor
// license agreements.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package aerospike

import "fmt"

type privilegeCode string

// Privilege determines user access granularity.
type Privilege struct {
	// Role
	Code privilegeCode

	// Namespace determines namespace scope. Apply permission to this namespace only.
	// If namespace is zero value, the privilege applies to all namespaces.
	Namespace string

	// Set name scope. Apply permission to this set within namespace only.
	// If set is zero value, the privilege applies to all sets within namespace.
	SetName string
}

func (p *Privilege) code() int {
	switch p.Code {
	// User can edit/remove other users.  Global scope only.
	case UserAdmin:
		return 0

	// User can perform systems administration functions on a database that do not involve user
	// administration.  Examples include server configuration.
	// Global scope only.
	case SysAdmin:
		return 1

	// User can perform data administration functions on a database that do not involve user
	// administration.  Examples include index and user defined function management.
	// Global scope only.
	case DataAdmin:
		return 2

	// User can read data only.
	case Read:
		return 10

	// User can read and write data.
	case ReadWrite:
		return 11

	// User can read and write data through user defined functions.
	case ReadWriteUDF:
		return 12

	// User can read and write data through user defined functions.
	case Write:
		return 13
	}

	panic("invalid role: " + p.Code)
}

func privilegeFrom(code uint8) privilegeCode {
	switch code {
	// User can edit/remove other users.  Global scope only.
	case 0:
		return UserAdmin

	// User can perform systems administration functions on a database that do not involve user
	// administration.  Examples include server configuration.
	// Global scope only.
	case 1:
		return SysAdmin

	// User can perform data administration functions on a database that do not involve user
	// administration.  Examples include index and user defined function management.
	// Global scope only.
	case 2:
		return DataAdmin

	// User can read data.
	case 10:
		return Read

	// User can read and write data.
	case 11:
		return ReadWrite

	// User can read and write data through user defined functions.
	case 12:
		return ReadWriteUDF

	// User can only write data.
	case 13:
		return Write
	}

	panic(fmt.Sprintf("invalid privilege code: %v", code))
}

func (p *Privilege) canScope() bool {
	return p.code() >= 10
}
