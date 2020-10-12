// Copyright 2013-2020 Aerospike, Inc.
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
	Name string

	Privileges []Privilege
	Whitelist  []string
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
