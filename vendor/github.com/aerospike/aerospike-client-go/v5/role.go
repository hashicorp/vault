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

// Role allows granular access to database entities for users.
type Role struct {
	// Name is role name
	Name string

	// Priviledge is the list of assigned privileges
	Privileges []Privilege

	// While is the list of allowable IP addresses
	Whitelist []string

	// ReadQuota is the maximum reads per second limit for the role
	ReadQuota uint32

	// WriteQuota is the maximum writes per second limit for the role
	WriteQuota uint32
}

// Pre-defined user roles.
const (
	// UserAdmin allows to manages users and their roles.
	UserAdmin privilegeCode = "user-admin"

	// SysAdmin allows to manage indexes, user defined functions and server configuration.
	SysAdmin privilegeCode = "sys-admin"

	// DataAdmin allows to manage indicies and user defined functions.
	DataAdmin privilegeCode = "data-admin"

	// ReadWriteUDF allows read, write and UDF transactions with the database.
	ReadWriteUDF privilegeCode = "read-write-udf"

	// ReadWrite allows read and write transactions with the database.
	ReadWrite privilegeCode = "read-write"

	// Read allows read transactions with the database.
	Read privilegeCode = "read"

	// Write allows write transactions with the database.
	Write privilegeCode = "write"
)
